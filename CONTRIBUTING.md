# Contributing to Trimora

Thanks for your interest in Trimora! This project is intentionally small — a minimal, self-hosted URL shortener — so contributions that keep it focused are most likely to land.

## Project scope

Trimora deliberately stays narrow:

**In scope:** shortening a URL, optional alias, optional expiry, redirect, copy in UI, health checks.

**Out of scope:** authentication, dashboards, analytics, billing, QR codes, custom domains, queues, Redis, ORMs, migration frameworks, marketing pages.

If you're unsure whether an idea fits, please open an issue first to discuss.

## Getting started

See the [Quick Start in the README](README.md#-quick-start) for setup. The fastest path is:

```bash
git clone https://github.com/hassankhan2608/Trimora.git
cd Trimora
cp .env.example .env
docker compose --profile local-db up --build
```

## Repository layout

```
api/    Go backend (Chi + database/sql + pgx + slog)
web/    Vite + React + TypeScript frontend
docs/   Screenshots and assets
```

## Development workflow

1. Fork the repo and create a feature branch:
   ```bash
   git checkout -b feat/your-thing
   ```
2. Make small, focused changes.
3. Run the checks listed below.
4. Open a pull request describing **what** changed and **why**.

## Required checks before opening a PR

**Backend (`api/`):**
```bash
cd api
gofmt -l .
go vet ./...
go test ./...
go build ./...
```

**Frontend (`web/`):**
```bash
cd web
npm run typecheck
npm run lint
npm run build
```

## Coding style

- **Go:** standard `gofmt`, `slog` for logging, `database/sql` (no ORM), small focused packages under `internal/`.
- **TypeScript/React:** strict TypeScript, plain CSS (no styling libraries), keep components small.
- **Commits:** small, conventional commits (`feat:`, `fix:`, `chore:`, `docs:`, `refactor:`).

## Pull request guidelines

- Keep PRs small and focused. One concern per PR.
- Include screenshots for any UI change.
- Update the README if you change behavior, environment variables, or endpoints.
- Make sure CI / local checks pass.

## Reporting bugs

Open a [GitHub issue](https://github.com/hassankhan2608/Trimora/issues) with:

- What you did
- What you expected
- What actually happened
- Trimora version / commit, OS, browser (if UI-related)
- Logs from the API or browser console

## Reporting security issues

Please **do not** open a public issue for security vulnerabilities. See [SECURITY.md](SECURITY.md) for the private disclosure process.

## Code of conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md). By participating, you agree to abide by its terms.

## License

By contributing, you agree that your contributions will be licensed under the [MIT License](LICENSE).
