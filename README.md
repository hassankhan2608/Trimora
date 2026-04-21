<div align="center">

# Trimora

Small URL shortener. Go API. Vite frontend. PostgreSQL storage.

<p>
  <a href="#quick-start"><b>Quick Start</b></a>
  <span>&nbsp;·&nbsp;</span>
  <a href="#api"><b>API</b></a>
  <span>&nbsp;·&nbsp;</span>
  <a href="#environment"><b>Environment</b></a>
  <span>&nbsp;·&nbsp;</span>
  <a href="#structure"><b>Structure</b></a>
  <span>&nbsp;·&nbsp;</span>
  <a href="#checks"><b>Checks</b></a>
</p>

<p>
  <a href="LICENSE"><img alt="License" src="https://img.shields.io/badge/license-MIT-111111?style=flat-square"></a>
  <img alt="Backend" src="https://img.shields.io/badge/backend-Go-00ADD8?style=flat-square">
  <img alt="Frontend" src="https://img.shields.io/badge/frontend-Vite-646CFF?style=flat-square">
  <img alt="Database" src="https://img.shields.io/badge/database-PostgreSQL-4169E1?style=flat-square">
</p>

</div>

---

Trimora does one thing: turns long links into short ones. It has no accounts, dashboards, analytics, billing, or marketing screens.

## Stack

| Part | Tooling |
| --- | --- |
| API | Go, Chi, slog |
| Database | PostgreSQL |
| SQL | `database/sql`, pgx |
| Web | Vite, React |
| Runtime | Docker, Docker Compose |
| Quality | gofmt, go vet, golangci-lint, npm lint/build |

## Quick Start

<table>
<tr>
<td><a href="#docker"><b>Docker</b></a></td>
<td><a href="#api-only"><b>API only</b></a></td>
<td><a href="#web-only"><b>Web only</b></a></td>
</tr>
</table>

### Docker

```bash
cp .env.example .env
docker compose up --build
```

Open:

```txt
http://localhost:5173
```

### API only

```bash
cd backend
cp .env.example .env
go mod download
go run ./cmd/api
```

API runs on:

```txt
http://localhost:8080
```

### Web only

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Web runs on:

```txt
http://localhost:5173
```

## API

### Shorten a URL

```http
POST /api/shorten
Content-Type: application/json
```

```json
{
  "url": "https://example.com/very/long/page",
  "customCode": "notes"
}
```

```json
{
  "code": "notes",
  "shortUrl": "http://localhost:8080/notes",
  "originalUrl": "https://example.com/very/long/page"
}
```

### Redirect

```http
GET /notes
```

### Health

```http
GET /healthz
```

## Environment

### Root `.env`

```env
POSTGRES_DB=trimora
POSTGRES_USER=trimora
POSTGRES_PASSWORD=trimora
POSTGRES_PORT=5432

API_PORT=8080
WEB_PORT=5173
PUBLIC_BASE_URL=http://localhost:8080
```

### Backend `.env`

```env
APP_ENV=local
PORT=8080
DATABASE_URL=postgres://trimora:trimora@localhost:5432/trimora?sslmode=disable
PUBLIC_BASE_URL=http://localhost:8080
```

### Frontend `.env`

```env
VITE_API_URL=http://localhost:8080
```

## Structure

```txt
Trimora/
├── backend/
│   ├── cmd/api/
│   ├── internal/
│   ├── pkg/shortid/
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   ├── public/
│   ├── Dockerfile
│   └── package.json
├── docker-compose.yml
├── README.md
└── LICENSE
```

## Checks

### Backend

```bash
cd backend
gofmt -w .
go test ./...
go vet ./...
golangci-lint run
```

### Frontend

```bash
cd frontend
npm run lint
npm run build
```

## Commit Format

```txt
chore: initialize project metadata
feat(api): add shorten endpoint
feat(web): add shortener form
fix(api): handle duplicate aliases
refactor(api): simplify repository errors
docs: update docker setup
```

## License

[MIT](LICENSE)
