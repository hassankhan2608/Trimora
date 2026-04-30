# Security Policy

The Trimora maintainers take security seriously. Thank you for helping keep the project and its users safe.

## Supported versions

Trimora is a small project and only the **latest commit on `main`** receives security fixes. Older tags and forks are not supported.

| Version          | Supported          |
| ---------------- | ------------------ |
| `main` (latest)  | :white_check_mark: |
| Older commits    | :x:                |

## Reporting a vulnerability

**Please do not open a public GitHub issue for security vulnerabilities.**

Instead, report privately via one of:

1. **GitHub Security Advisories** — preferred. Go to the [Security tab](https://github.com/hassankhan2608/Trimora/security/advisories/new) and click *Report a vulnerability*. This keeps the discussion private until a fix is ready.
2. **Email** — open a regular issue asking for a private contact channel and a maintainer will reach out, or contact the repository owner directly through their GitHub profile.

When reporting, please include:

- A clear description of the vulnerability and the affected component (`api/`, `web/`, deployment).
- Steps to reproduce, ideally with a minimal proof of concept.
- The commit SHA or version you tested against.
- The impact you believe it has (e.g., open redirect, SSRF, data exposure, denial of service).
- Any suggested mitigation, if you have one.

## What to expect

- **Acknowledgement:** within 5 business days.
- **Initial assessment:** within 10 business days, with a planned timeline for a fix.
- **Disclosure:** coordinated with you. Once a fix is released on `main`, we'll publish a GitHub Security Advisory crediting you (unless you prefer to remain anonymous).

## Out of scope

- Vulnerabilities in third-party dependencies that have not been publicly disclosed yet — please report those upstream first.
- Issues that require physical access to the server, a compromised database, or stolen credentials.
- Best-practice suggestions without a concrete exploit (please open a regular issue or PR for those).

## Hardening recommendations for self-hosters

- Always set `ALLOWED_ORIGINS` to the exact origins you serve from — never `*`.
- Use TLS in front of Trimora (e.g., a reverse proxy like Caddy or nginx).
- Use a managed Postgres provider with `sslmode=require`, or a private network for self-hosted Postgres.
- Rotate `DATABASE_URL` credentials if `.env` is ever exposed.
- Keep Trimora and its base Docker images up to date.

Thanks again for helping make Trimora safer.
