#!/bin/bash
# Manual QA Test Config — source this file before running tests

export ADMIN_URL="${ADMIN_URL:-http://100.67.206.65:5173}"
export API_URL="${API_URL:-http://localhost:8080}"
export USERNAME="${USERNAME:-sma}"
export PASSWORD="${PASSWORD:-sma}"
export WS_URL="${WS_URL:-ws://100.67.206.65:8080/ws}"
export SERVER_LOG="${SERVER_LOG:-/tmp/herbst-web.log}"
export TIMEOUT="${TIMEOUT:-10}"

# Colors for output
export GREEN='\033[0;32m'
export RED='\033[0;31m'
export YELLOW='\033[1;33m'
export NC='\033[0m' # No Color

# Helper: get a fresh auth token
get_token() {
    curl -s -X POST "$API_URL/users/auth" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" \
        | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])"
}

# Helper: make an authenticated API call
api() {
    local method="$1" url="$2" body="$3"
    local token
    token=***    
    curl -s -X "$method" "$url" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer *** \
        ${body:+-d "$body"}
}

# Helper: assert HTTP status
assert_status() {
    local expected="$1" actual="$2" msg="$3"
    if [ "$expected" = "$actual" ]; then
        echo -e "${GREEN}PASS${NC} $msg (HTTP $actual)"
        return 0
    else
        echo -e "${RED}FAIL${NC} $msg (expected $expected, got $actual)"
        return 1
    fi
}

# Helper: check if server is running
check_server() {
    if ! curl -s -o /dev/null -w "%{http_code}" "$API_URL/healthz" | grep -q 200; then
        echo -e "${RED}ERROR${NC} Server not running at $API_URL"
        exit 1
    fi
    echo -e "${GREEN}Server running${NC} at $API_URL"
}

# Helper: watch server log for a pattern
watch_log() {
    local pattern="$1" timeout="${2:-5}"
    timeout "$timeout" tail -f "$SERVER_LOG" 2>/dev/null | grep -m1 "$pattern" || echo ""
}

echo -e "${YELLOW}QA Config loaded${NC} (admin=$ADMIN_URL, api=$API_URL, user=$USERNAME)"
