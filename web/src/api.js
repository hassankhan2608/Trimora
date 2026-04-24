const baseURL = (import.meta.env.VITE_API_BASE_URL || "").replace(/\/$/, "");

export async function shorten({ url, alias }) {
  const body = { url };
  if (alias) body.alias = alias;

  const res = await fetch(`${baseURL}/api/shorten`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });

  let data = null;
  try {
    data = await res.json();
  } catch {
    // empty body
  }

  if (!res.ok) {
    const message = data?.error || `Request failed (${res.status})`;
    throw new Error(message);
  }
  return data;
}
