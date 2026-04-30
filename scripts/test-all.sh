#!/usr/bin/env bash
# Trimora full-stack test runner.
#
# What it does (in order, fails fast on any error):
#   1. Go: gofmt check, go vet, go test, go build
#   2. Web: npm run typecheck, npm run lint, npm run build
#   3. Boot a local Postgres (docker compose --profile local-db) if not already up
#   4. Start the freshly-built API binary against that DB
#   5. Start `vite preview` against the freshly-built frontend
#   6. Run scripts/smoke-api.sh (curl-driven endpoint matrix)
#   7. Run scripts/smoke-ui.mjs (Playwright desktop + mobile)
#   8. Tear everything down (only what we started)
#
# Usage:
#   scripts/test-all.sh              # run everything
#   scripts/test-all.sh --skip-ui    # skip Playwright (faster)
#   scripts/test-all.sh --keep-up    # leave services running for debugging

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

# ---- args ----
SKIP_UI=0
KEEP_UP=0
for arg in "$@"; do
  case "$arg" in
    --skip-ui) SKIP_UI=1 ;;
    --keep-up) KEEP_UP=1 ;;
    -h|--help)
      sed -n '2,20p' "$0"
      exit 0 ;;
    *) echo "unknown flag: $arg" >&2; exit 2 ;;
  esac
done

# ---- pretty output ----
if [ -t 1 ]; then
  C_RESET=$'\e[0m'; C_BOLD=$'\e[1m'; C_GREEN=$'\e[32m'; C_RED=$'\e[31m'; C_BLUE=$'\e[34m'; C_YELLOW=$'\e[33m'
else
  C_RESET=""; C_BOLD=""; C_GREEN=""; C_RED=""; C_BLUE=""; C_YELLOW=""
fi
section() { printf "\n%s── %s ──%s\n" "$C_BOLD$C_BLUE" "$1" "$C_RESET"; }
ok()      { printf "%s✓%s %s\n" "$C_GREEN" "$C_RESET" "$1"; }
warn()    { printf "%s!%s %s\n" "$C_YELLOW" "$C_RESET" "$1"; }
die()     { printf "%s✗ %s%s\n" "$C_RED" "$1" "$C_RESET" >&2; exit 1; }

# ---- runtime state we may need to clean up ----
TMP_DIR="$(mktemp -d -t trimora-test.XXXXXX)"
API_PID=""
WEB_PID=""
STARTED_DB=0
STARTED_API=0
STARTED_WEB=0

cleanup() {
  local exit_code=$?
  if [ "$KEEP_UP" -eq 1 ]; then
    warn "leaving services running (--keep-up)"
    warn "tmp dir: $TMP_DIR"
    exit "$exit_code"
  fi
  section "cleanup"
  if [ "$STARTED_API" -eq 1 ] && [ -n "$API_PID" ] && kill -0 "$API_PID" 2>/dev/null; then
    kill "$API_PID" 2>/dev/null || true; ok "stopped api ($API_PID)"
  fi
  if [ "$STARTED_WEB" -eq 1 ] && [ -n "$WEB_PID" ] && kill -0 "$WEB_PID" 2>/dev/null; then
    kill "$WEB_PID" 2>/dev/null || true; ok "stopped web ($WEB_PID)"
  fi
  if [ "$STARTED_DB" -eq 1 ]; then
    docker compose --profile local-db down -v >/dev/null 2>&1 || true
    ok "stopped local-db"
  fi
  rm -rf "$TMP_DIR"
  exit "$exit_code"
}
trap cleanup EXIT INT TERM

wait_for_http() {
  local url="$1" tries="${2:-40}"
  for _ in $(seq 1 "$tries"); do
    if curl -sf -o /dev/null "$url"; then return 0; fi
    sleep 0.25
  done
  return 1
}

# =============================================================
section "1/7  go: fmt, vet, test, build"
# =============================================================
pushd api >/dev/null
unformatted="$(gofmt -l .)"
[ -z "$unformatted" ] || die "gofmt issues:\n$unformatted"
ok "gofmt clean"
go vet ./... && ok "go vet"
go test ./... && ok "go test"
go build -o "$TMP_DIR/trimora-api" ./cmd/server && ok "go build"
popd >/dev/null

# =============================================================
section "2/7  web: typecheck, lint, build"
# =============================================================
pushd web >/dev/null
npm run typecheck --silent && ok "typecheck"
npm run lint --silent && ok "lint"
npm run build --silent && ok "build"
popd >/dev/null

# =============================================================
section "3/7  postgres (local-db profile)"
# =============================================================
if docker exec trimora-pg pg_isready -U trimora >/dev/null 2>&1; then
  ok "reusing existing trimora-pg container"
elif docker ps -a --format '{{.Names}}' | grep -qx trimora-pg; then
  ok "starting existing trimora-pg container"
  docker start trimora-pg >/dev/null
else
  ok "starting docker compose local-db profile"
  docker compose --profile local-db up -d db >/dev/null
  STARTED_DB=1
fi
# Wait for pg to accept connections
for _ in $(seq 1 40); do
  docker exec trimora-pg pg_isready -U trimora >/dev/null 2>&1 && break
  sleep 0.5
done
docker exec trimora-pg pg_isready -U trimora >/dev/null || die "postgres did not become ready"
ok "postgres ready on :5432"

# Reset schema for a clean run
docker exec trimora-pg psql -U trimora -d trimora -c "DROP TABLE IF EXISTS links CASCADE;" >/dev/null
ok "schema reset"

# =============================================================
section "4/7  api server"
# =============================================================
if curl -sf -o /dev/null http://localhost:8080/livez 2>/dev/null; then
  warn "something already on :8080 — stop it first or use --keep-up after this run"
  die "port 8080 in use"
fi

PORT=8080 \
  BASE_URL=http://localhost:8080 \
  ALLOWED_ORIGINS=http://localhost:4173 \
  DATABASE_URL="postgres://trimora:trimora@localhost:5432/trimora?sslmode=disable" \
  "$TMP_DIR/trimora-api" >"$TMP_DIR/api.log" 2>&1 &
API_PID=$!
STARTED_API=1
wait_for_http "http://localhost:8080/livez" 60 || { tail -40 "$TMP_DIR/api.log"; die "api failed to start"; }
ok "api up (pid $API_PID)"

# =============================================================
section "5/7  web preview"
# =============================================================
if curl -sf -o /dev/null http://localhost:4173/ 2>/dev/null; then
  die "port 4173 in use"
fi

(cd web && npx --no-install vite preview --port 4173 --host >"$TMP_DIR/web.log" 2>&1) &
WEB_PID=$!
STARTED_WEB=1
wait_for_http "http://localhost:4173/" 60 || { tail -40 "$TMP_DIR/web.log"; die "web failed to start"; }
ok "web up (pid $WEB_PID)"

# =============================================================
section "6/7  api smoke (curl matrix)"
# =============================================================
API_BASE=http://localhost:8080 bash "$ROOT/scripts/smoke-api.sh"

# =============================================================
section "7/7  ui smoke (playwright)"
# =============================================================
if [ "$SKIP_UI" -eq 1 ]; then
  warn "skipped (--skip-ui)"
else
  WEB_BASE=http://localhost:4173 API_BASE=http://localhost:8080 \
    node "$ROOT/scripts/smoke-ui.mjs"
fi

printf "\n%s%s ALL CHECKS PASSED %s\n\n" "$C_BOLD$C_GREEN" "✔" "$C_RESET"
