#!/usr/bin/env bash
# API endpoint smoke matrix. Assumes API is reachable at $API_BASE
# and a local Postgres named `trimora-pg` for the expiry-mutation case.
set -euo pipefail

API="${API_BASE:-http://localhost:8080}"

if [ -t 1 ]; then
  C_RESET=$'\e[0m'; C_GREEN=$'\e[32m'; C_RED=$'\e[31m'; C_DIM=$'\e[2m'
else
  C_RESET=""; C_GREEN=""; C_RED=""; C_DIM=""
fi
PASS=0; FAIL=0
check() { # name, expected_substr, actual
  local name="$1" expected="$2" actual="$3"
  if [[ "$actual" == *"$expected"* ]]; then
    printf "  %s✓%s %-32s %s%s%s\n" "$C_GREEN" "$C_RESET" "$name" "$C_DIM" "$actual" "$C_RESET"
    PASS=$((PASS + 1))
  else
    printf "  %s✗%s %-32s expected %q got %q\n" "$C_RED" "$C_RESET" "$name" "$expected" "$actual"
    FAIL=$((FAIL + 1))
  fi
}

post() { # endpoint json — returns body even on 4xx (no -f)
  curl -s -X POST "$API$1" -H "Content-Type: application/json" -d "$2"
}
post_status() { # endpoint json -> http code
  curl -s -o /dev/null -w "%{http_code}" -X POST "$API$1" -H "Content-Type: application/json" -d "$2"
}
get_status() { # url, accept
  curl -s -o /dev/null -w "%{http_code}" -H "Accept: $2" "$1"
}
get_body_json() { # url
  curl -s -H "Accept: application/json" "$1"
}
html_title() { # url
  curl -s -H "Accept: text/html" "$1" | grep -oE '<title>[^<]+</title>' | head -1
}

# 1. Health endpoints
check "GET /livez"            "200" "$(get_status "$API/livez" "*/*")"
check "GET /readyz"           "200" "$(get_status "$API/readyz" "*/*")"
check "GET /healthz"          "200" "$(get_status "$API/healthz" "*/*")"

# 2. Create basic
ts="$(date +%s%N)"
basic="$(post /api/shorten "{\"url\":\"https://example.com/smoke-${ts}\"}")"
check "POST basic"            "short_url" "$basic"
basic_code="$(printf '%s' "$basic" | sed -n 's/.*"code":"\([^"]*\)".*/\1/p')"

# 3. Custom alias
alias="smoke-${ts}"
aliased="$(post /api/shorten "{\"url\":\"https://example.com/aliased\",\"alias\":\"${alias}\"}")"
check "POST alias"            "\"code\":\"${alias}\"" "$aliased"

# 4. Dedupe (same URL no expiry => same code)
dup="$(post /api/shorten "{\"url\":\"https://example.com/smoke-${ts}\"}")"
check "POST dedupe"           "\"code\":\"${basic_code}\"" "$dup"

# 5. Invalid expiry
bad_exp="$(post /api/shorten '{"url":"https://example.com/x","expires_in":"5m"}')"
check "POST invalid expiry"   "expiry must be one of" "$bad_exp"

# 6. Reserved alias
reserved="$(post /api/shorten '{"url":"https://example.com/x","alias":"livez"}')"
check "POST reserved alias"   "reserved" "$reserved"

# 7. Conflicting alias
conflict="$(post /api/shorten "{\"url\":\"https://example.com/y\",\"alias\":\"${alias}\"}")"
check "POST conflicting alias" "alias unavailable" "$conflict"

# 8. Expiring create
exp_resp="$(post /api/shorten "{\"url\":\"https://example.com/expiring-${ts}\",\"expires_in\":\"1h\"}")"
check "POST with 1h expiry"   "expires_at" "$exp_resp"
exp_code="$(printf '%s' "$exp_resp" | sed -n 's/.*"code":"\([^"]*\)".*/\1/p')"

# 9. Redirect
redir_status="$(curl -s -o /dev/null -w "%{http_code}" "$API/${basic_code}")"
redir_loc="$(curl -s -o /dev/null -w "%{redirect_url}" "$API/${basic_code}")"
check "GET redirect 302"      "302" "$redir_status"
check "GET redirect Location" "smoke-${ts}" "$redir_loc"

# 10. 404 JSON
check "GET 404 JSON status"   "404" "$(get_status "$API/no-such-${ts}" "application/json")"
check "GET 404 JSON body"     "link not found" "$(get_body_json "$API/no-such-${ts}")"

# 11. 404 HTML
check "GET 404 HTML status"   "404" "$(get_status "$API/no-such-${ts}" "text/html")"
check "GET 404 HTML title"    "Link not found" "$(html_title "$API/no-such-${ts}")"

# 12. Force-expire then 410 (requires local pg)
if docker exec trimora-pg pg_isready -U trimora >/dev/null 2>&1; then
  docker exec trimora-pg psql -U trimora -d trimora \
    -c "UPDATE links SET expires_at = NOW() - INTERVAL '1 minute' WHERE code = '${exp_code}';" >/dev/null
  check "GET 410 JSON status" "410" "$(get_status "$API/${exp_code}" "application/json")"
  check "GET 410 JSON body"   "link expired" "$(get_body_json "$API/${exp_code}")"
  check "GET 410 HTML status" "410" "$(get_status "$API/${exp_code}" "text/html")"
  check "GET 410 HTML title"  "Link expired" "$(html_title "$API/${exp_code}")"
else
  printf "  %s!%s skipping 410 tests (no trimora-pg)\n" "$C_DIM" "$C_RESET"
fi

printf "\n  %d passed, %d failed\n" "$PASS" "$FAIL"
[ "$FAIL" -eq 0 ] || exit 1
