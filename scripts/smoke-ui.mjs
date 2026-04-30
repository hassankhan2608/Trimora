// Trimora UI smoke test (Playwright). Assumes the web preview is reachable
// at WEB_BASE and the API at API_BASE. Tests both desktop and mobile viewports.
import { chromium, devices } from "playwright";

const WEB = process.env.WEB_BASE || "http://localhost:4173";
const API = process.env.API_BASE || "http://localhost:8080";

const C = {
  reset: "\x1b[0m",
  green: "\x1b[32m",
  red: "\x1b[31m",
  dim: "\x1b[2m",
  bold: "\x1b[1m",
};

async function runOn(label, contextOpts) {
  console.log(`\n  ${C.bold}— ${label} —${C.reset}`);
  const browser = await chromium.launch();
  const ctx = await browser.newContext(contextOpts);
  await ctx.grantPermissions(["clipboard-read", "clipboard-write"], { origin: WEB });
  const page = await ctx.newPage();

  const failures = [];
  const ok = (name, detail = "") =>
    console.log(`    ${C.green}✓${C.reset} ${name.padEnd(34)} ${C.dim}${detail}${C.reset}`);
  const fail = (name, msg) => {
    failures.push(`${label}/${name}: ${msg}`);
    console.log(`    ${C.red}✗${C.reset} ${name.padEnd(34)} ${msg}`);
  };

  try {
    // 1. Page loads
    const resp = await page.goto(WEB, { waitUntil: "networkidle" });
    resp?.status() === 200 ? ok("homepage 200") : fail("homepage", `status=${resp?.status()}`);

    // 2. Hero rendered
    const title = (await page.locator(".hero__title").innerText()).replace(/\s+/g, " ").trim();
    /short links/i.test(title) ? ok("hero title", title) : fail("hero", title);

    // 3. Empty submit blocked by HTML5 validation
    await page.locator("button[type=submit]").click();
    await page.waitForTimeout(150);
    page.url().startsWith(WEB) ? ok("empty submit blocked") : fail("empty submit", page.url());

    // 4. Submit a real URL with 1h expiry
    const targetUrl = `https://example.com/ui-${label}-${Date.now()}`;
    await page.locator('input[type="url"]').fill(targetUrl);
    await page.locator("select").selectOption("1h");
    await page.locator("button[type=submit]").click();

    await page.locator(".result").waitFor({ timeout: 5000 }).catch(() => {});
    (await page.locator(".result").isVisible())
      ? ok("result panel visible")
      : fail("result panel", "not visible");

    const shortHref = await page.locator(".result__link").getAttribute("href");
    shortHref?.startsWith(API) ? ok("short link", shortHref) : fail("short link", String(shortHref));

    const expiryText = (await page.locator(".result__expiry").innerText().catch(() => "")).trim();
    /expires/i.test(expiryText) ? ok("expiry pill", expiryText) : fail("expiry pill", expiryText);

    // 5. Copy button toggles
    await page.locator(".result button").click();
    await page.waitForTimeout(150);
    const copyText = (await page.locator(".result button").innerText()).trim();
    /copied/i.test(copyText) ? ok("copy → Copied") : fail("copy", copyText);

    // 6. Redirect via API works
    const code = shortHref?.split("/").pop();
    const redirectResp = await page.request.get(`${API}/${code}`, { maxRedirects: 0 });
    redirectResp.status() === 302 && redirectResp.headers().location === targetUrl
      ? ok("API redirect 302", targetUrl)
      : fail("API redirect", `status=${redirectResp.status()} loc=${redirectResp.headers().location}`);

    // 7. Custom alias
    const alias = `ui-${label}-${Math.random().toString(36).slice(2, 8)}`;
    await page.locator('input[type="url"]').fill(`https://example.com/aliased-${alias}`);
    await page.locator('input[placeholder="my-link"]').fill(alias);
    await page.locator("select").selectOption("");
    await page.locator("button[type=submit]").click();
    await page.waitForTimeout(800);
    const newHref = await page.locator(".result__link").getAttribute("href");
    newHref?.endsWith(`/${alias}`) ? ok("alias used", newHref) : fail("alias", String(newHref));

    // 8. Reserved alias error
    await page.locator('input[type="url"]').fill("https://example.com/reserved-test");
    await page.locator('input[placeholder="my-link"]').fill("livez");
    await page.locator("button[type=submit]").click();
    await page.waitForTimeout(500);
    const errText = (await page.locator(".form__error").innerText().catch(() => "")).trim();
    /reserved/i.test(errText) ? ok("reserved alias error", errText) : fail("reserved", errText);

    // 9. Mobile-specific: no horizontal overflow
    if (contextOpts.viewport && contextOpts.viewport.width < 500) {
      const docW = await page.evaluate(() => document.documentElement.scrollWidth);
      const winW = await page.evaluate(() => window.innerWidth);
      docW <= winW + 1
        ? ok("no horizontal overflow", `${docW}≤${winW}`)
        : fail("overflow", `${docW} > ${winW}`);
    }
  } finally {
    await browser.close();
  }
  return failures;
}

const desktopFails = await runOn("desktop", { viewport: { width: 1280, height: 900 } });
const mobileFails = await runOn("mobile", { ...devices["iPhone 13"] });

const total = desktopFails.length + mobileFails.length;
if (total === 0) {
  console.log(`\n  ${C.green}${C.bold}all UI checks passed${C.reset}`);
  process.exit(0);
} else {
  console.log(`\n  ${C.red}${C.bold}${total} UI check(s) failed${C.reset}`);
  for (const f of [...desktopFails, ...mobileFails]) console.log(`    - ${f}`);
  process.exit(1);
}
