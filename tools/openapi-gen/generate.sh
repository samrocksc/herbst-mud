#!/usr/bin/env bash
# Generate the OpenAPI spec from route definitions + static additions
# Usage: ./generate.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SERVER_DIR="$SCRIPT_DIR/../server"
STATIC_DIR="$SERVER_DIR/static"
SPEC="$STATIC_DIR/openapi.json"
ADDITIONS="$STATIC_DIR/openapi_additions.json"

echo "=== OpenAPI Generator ==="
echo ""
echo "Current spec: $SPEC"

if [ ! -f "$SPEC" ]; then
    echo "ERROR: $SPEC not found"
    exit 1
fi

# Count paths
PATH_COUNT=$(python3 -c "import json; s=json.load(open('$SPEC')); print(len(s.get('paths',{})))")
echo "Paths: $PATH_COUNT"
echo "Schemas: $(python3 -c "import json; s=json.load(open('$SPEC')); print(len(s.get('components',{}).get('schemas',{})))")"

echo ""
echo "If this looks stale, update $SPEC directly or add endpoints to route files."
echo "The spec is also served at http://localhost:8080/openapi.json"
echo ""
echo "Done."
