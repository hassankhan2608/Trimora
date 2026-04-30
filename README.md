<div align="center">

<img src="docs/screenshots/desktop-hero.png" alt="Trimora — a minimal open-source URL shortener built with Go, React, and PostgreSQL" width="880">

# Trimora

**A minimal, self-hosted URL shortener.**
Paste a URL, pick an optional alias, set an optional expiry, share something shorter.

Built with **Go**, **React**, and **PostgreSQL**. No accounts, no analytics, no clutter.

<p>
  <a href="#-quick-start"><b>Quick Start</b></a>
  &nbsp;·&nbsp;
  <a href="#-features"><b>Features</b></a>
  &nbsp;·&nbsp;
  <a href="#-screenshots"><b>Screenshots</b></a>
  &nbsp;·&nbsp;
  <a href="#-api-reference"><b>API</b></a>
  &nbsp;·&nbsp;
  <a href="#-configuration"><b>Configuration</b></a>
  &nbsp;·&nbsp;
  <a href="#-development"><b>Development</b></a>
</p>

<p>
  <a href="LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/license-MIT-111111?style=flat-square"></a>
  <img alt="Go" src="https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white">
  <img alt="React" src="https://img.shields.io/badge/React-19-61DAFB?style=flat-square&logo=react&logoColor=white">
  <img alt="Vite" src="https://img.shields.io/badge/Vite-6-646CFF?style=flat-square&logo=vite&logoColor=white">
  <img alt="TypeScript" src="https://img.shields.io/badge/TypeScript-5-3178C6?style=flat-square&logo=typescript&logoColor=white">
  <img alt="PostgreSQL" src="https://img.shields.io/badge/PostgreSQL-16-4169E1?style=flat-square&logo=postgresql&logoColor=white">
  <img alt="Docker" src="https://img.shields.io/badge/Docker-ready-2496ED?style=flat-square&logo=docker&logoColor=white">
</p>

</div>

---

## ✨ Features

- 🔗 **Short links in one click** — paste a URL, get a short link.
- ✏️ **Custom aliases** — claim memorable names like `trimora.app/launch-notes`.
- ⏳ **Optional expiry** — links can expire after `1h`, `1d`, `7d`, or `30d`. Default is forever.
- ♻️ **Smart de-duplication** — the same URL with no alias and no expiry returns the same short code.
- 🪶 **Calm, accessible UI** — single page, mobile-first, warm editorial style, no tracking.
- 🩺 **Production-ready health checks** — `/livez` and `/readyz` for Kubernetes / load balancers.
- 🐳 **One-command deploy** — `docker compose up`, with optional bundled Postgres or your own (Supabase, Neon, RDS).
- 🧪 **Tested** — Go unit tests, ESLint, strict TypeScript, Playwright UI checks.

---

## 📸 Screenshots

<table>
  <tr>
    <td align="center"><b>Desktop · Empty form</b></td>
    <td align="center"><b>Desktop · Result with expiry</b></td>
    <td align="center"><b>Mobile</b></td>
  </tr>
  <tr>
    <td><img src="docs/screenshots/desktop-hero.png" alt="Trimora desktop UI showing URL shortener form with alias and expiry inputs"></td>
    <td><img src="docs/screenshots/desktop-result.png" alt="Trimora desktop UI showing a generated short link with copy button and expiry pill"></td>
    <td><img src="docs/screenshots/mobile-result.png" alt="Trimora responsive mobile layout with stacked form fields and short link result"></td>
  </tr>
</table>

---

## 🚀 Quick Start

> **Prerequisites:** [Docker](https://docs.docker.com/get-docker/) **or** [Go 1.25+](https://go.dev/dl/) with [Node.js 20+](https://nodejs.org/) and a PostgreSQL 14+ database.

<details open>
<summary><b>🐳 Docker Compose (recommended)</b></summary>

<br>

Spin up the full stack — API, web UI, and a local Postgres — with one command:

```bash
git clone https://github.com/your-org/trimora.git
cd trimora
cp .env.example .env
docker compose --profile local-db up --build
```

Then open **<http://localhost:5173>**.

> Want to point at an external database (Supabase, Neon, RDS)?
> Drop the `--profile local-db` flag and set `DATABASE_URL` in `.env`.

</details>

<details>
<summary><b>☁️ Docker + Supabase (or any managed Postgres)</b></summary>

<br>

```bash
cp .env.example .env
# Edit .env and set DATABASE_URL to your managed Postgres connection string.
# For Supabase, use the direct (non-pooling) connection on port 5432
# with sslmode=require.

docker compose up --build
```

Trimora opens persistent connections and runs `CREATE TABLE IF NOT EXISTS` on
boot, which works best on the **direct** port (`5432`) rather than the
pgbouncer pooler (`6543`).

</details>

<details>
<summary><b>🛠️ Local development (no Docker)</b></summary>

<br>

**1. Start Postgres** (any local Postgres works):

```bash
docker run -d --rm --name trimora-pg \
  -e POSTGRES_USER=trimora -e POSTGRES_PASSWORD=trimora -e POSTGRES_DB=trimora \
  -p 5432:5432 postgres:16-alpine
```

**2. Start the Go API:**

```bash
cd api
DATABASE_URL='postgres://trimora:trimora@localhost:5432/trimora?sslmode=disable' \
BASE_URL='http://localhost:8080' \
ALLOWED_ORIGINS='http://localhost:5173' \
PORT=8080 \
go run ./cmd/server
```

**3. Start the Vite web app:**

```bash
cd web
npm install
VITE_API_BASE_URL='http://localhost:8080' npm run dev
```

Open **<http://localhost:5173>**.

</details>

---

## 🔌 API Reference

Base URL: `http://localhost:8080` (defaults; override with `BASE_URL`).

<details open>
<summary><b><code>POST /api/shorten</code> — create a short link</b></summary>

<br>

**Request body**

| Field         | Type                            | Required | Notes                                                          |
| ------------- | ------------------------------- | -------- | -------------------------------------------------------------- |
| `url`         | string                          | yes      | `http://` or `https://`, max 2048 chars.                       |
| `alias`       | string                          | no       | 3–32 chars, `[A-Za-z0-9_-]`. Reserved: `api`, `livez`, `readyz`, `healthz`, `health`, `static`, `admin`, `assets`, `favicon.ico`, `robots.txt`. |
| `expires_in`  | `"1h"` \| `"1d"` \| `"7d"` \| `"30d"` | no | Omit for permanent links. |

**Example**

```bash
curl -X POST http://localhost:8080/api/shorten \
  -H 'Content-Type: application/json' \
  -d '{
    "url": "https://example.com/long/article",
    "alias": "launch-notes",
    "expires_in": "1d"
  }'
```

**Response — `201 Created`**

```json
{
  "code": "launch-notes",
  "short_url": "http://localhost:8080/launch-notes",
  "url": "https://example.com/long/article",
  "expires_at": "2026-04-28T15:33:09+05:30"
}
```

**Errors**

| Status | When |
| ------ | ---- |
| `400 Bad Request` | invalid URL, invalid/reserved alias, invalid `expires_in` |
| `409 Conflict`    | alias already taken |

</details>

<details>
<summary><b><code>GET /{code}</code> — follow a short link</b></summary>

<br>

```bash
curl -i http://localhost:8080/launch-notes
# HTTP/1.1 302 Found
# Location: https://example.com/long/article
```

| Status | When |
| ------ | ---- |
| `302 Found` | redirect to the original URL |
| `404 Not Found` | unknown code (HTML page for browsers, JSON for API clients) |
| `410 Gone` | link expired (HTML page for browsers, JSON for API clients) |

</details>

<details>
<summary><b><code>GET /livez</code>, <code>/readyz</code>, <code>/healthz</code> — health checks</b></summary>

<br>

| Endpoint    | Purpose                                                  |
| ----------- | -------------------------------------------------------- |
| `/livez`    | Process is alive. No DB call. Use as Kubernetes liveness probe. |
| `/readyz`   | Database is reachable (`PING` with 2s timeout). Use as readiness probe. |
| `/healthz`  | Backwards-compatible alias for `/livez`.                 |

```bash
curl http://localhost:8080/livez
# {"status":"ok"}
```

</details>

---

## ⚙️ Configuration

Trimora is configured entirely through environment variables. Copy `.env.example` to `.env` to get started.

<details open>
<summary><b>Server (API)</b></summary>

<br>

| Variable          | Default                  | Description                                            |
| ----------------- | ------------------------ | ------------------------------------------------------ |
| `PORT`            | `8080`                   | API listen port.                                       |
| `BASE_URL`        | `http://localhost:8080`  | Public origin used to build `short_url` in responses.  |
| `ALLOWED_ORIGINS` | `http://localhost:5173`  | Comma-separated CORS allow-list.                       |
| `DATABASE_URL`    | _required_               | PostgreSQL connection string.                          |

</details>

<details>
<summary><b>Database</b></summary>

<br>

```env
# Local Postgres via docker compose --profile local-db
DATABASE_URL=postgres://trimora:trimora@db:5432/trimora?sslmode=disable

# Managed Postgres (Supabase, Neon, RDS, …) — sslmode=require
DATABASE_URL=postgres://USER:PASSWORD@HOST:5432/DB?sslmode=require
```

The schema is auto-created on boot (`links` table + indexes). No migration tool required.

</details>

<details>
<summary><b>Web (Vite)</b></summary>

<br>

| Variable                | Default                  | Description                                                                       |
| ----------------------- | ------------------------ | --------------------------------------------------------------------------------- |
| `VITE_API_BASE_URL`     | _empty_                  | Absolute API URL for production builds. Leave empty in dev to use the proxy.      |
| `VITE_API_PROXY_TARGET` | `http://localhost:8080`  | Vite dev-server proxy target for `/api`, `/livez`, `/readyz`.                     |

</details>

---

## 🧱 Project Structure

```txt
Trimora/
├── api/                     # Go backend
│   ├── cmd/server/          # main entrypoint
│   └── internal/
│       ├── config/          # env loading
│       ├── httpapi/         # handlers, router, CORS, branded 404/410
│       ├── links/           # service + repository (database/sql + pgx)
│       ├── shortcode/       # base62 short-code generator (crypto/rand)
│       ├── storage/         # connection + schema bootstrap
│       └── validate/        # URL, alias, expiry rules
├── web/                     # Vite + React + TypeScript frontend
│   └── src/
│       ├── components/      # ShortenForm, Result
│       ├── types/api.ts     # shared API types
│       ├── api.ts           # tiny fetch client
│       └── App.tsx
├── docs/screenshots/        # README assets
├── docker-compose.yml
├── README.md
└── LICENSE
```

---

## 🧑‍💻 Development

<details open>
<summary><b>Backend checks</b></summary>

<br>

```bash
cd api
gofmt -l .
go vet ./...
go test ./...
go build ./...
```

</details>

<details>
<summary><b>Frontend checks</b></summary>

<br>

```bash
cd web
npm run typecheck
npm run lint
npm run build
```

</details>

<details>
<summary><b>End-to-end UI testing (Playwright)</b></summary>

<br>

The web UI is verified manually with Playwright on both desktop (1280×820) and mobile (390×844) viewports. The screenshots in this README are captured headlessly via the same flow.

</details>

---

## 🧭 Design Choices

- **No accounts, no analytics, no dashboard.** Trimora is a tool, not a product.
- **`database/sql` + pgx** — no ORM, no migration framework. The schema is small enough to manage in code.
- **`crypto/rand`-backed base62 codes** — short, URL-safe, collision-retried at the service layer.
- **Optional expiry** is a single nullable `expires_at` column with an index. Expired links return `410 Gone`.
- **Same URL = same short code** when no alias and no expiry are requested, to avoid pointless duplicates.
- **Branded 404 / 410 pages** for browsers, JSON for API clients — content-negotiated by `Accept`.

---

## 🤝 Contributing

Issues and pull requests are welcome. Please keep changes small and focused, and run the checks listed above before opening a PR.

## 📄 License

[MIT](LICENSE) © Trimora contributors
