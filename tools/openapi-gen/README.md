# OpenAPI Generator

The canonical OpenAPI spec lives at `server/static/openapi.json`. It is the **single source of truth** — all other artifacts (TypeScript types, Go clients, request validators) are generated from it.

## Regenerate the spec

After modifying routes, run:

```bash
./tools/openapi-gen/generate.sh
```

This script:
1. Scans `server/routes/*.go` for all registered endpoints
2. Merges any hand-written additions from `server/static/openapi_additions.json`
3. Writes the final spec to `server/static/openapi.json`
4. Prints a diff summary (paths added/removed)

## Long-term: oapi-codegen workflow (Phase 2)

Currently the spec is hand-written as a static JSON file. The Phase 2 approach uses `oapi-codegen` so the spec lives alongside Go types:

```bash
# Extract handlers from inline closures to named functions (SVC-001)
# Then annotate with Go doc comments:
go generate ./...

# Generate TypeScript client for admin web UI
oapi-codegen -generate typescript -package api \
  http://localhost:8080/openapi.json > admin/src/api/types.gen.ts

# Generate Go client SDK for admin TUI
oapi-codegen -generate gin -package client \
  http://localhost:8080/openapi.json > admin-tui/client.gen.go

# Generate request validators (middleware)
oapi-codegen -generate gin -package middleware \
  http://localhost:8080/openapi.json > server/middleware/validators.gen.go
```

## Adding new endpoints

When adding a new endpoint, update `server/static/openapi.json` directly (Phase 1) or add doc comments to the handler function (Phase 2).

Any new endpoint MUST:
- Be added to the spec BEFORE merging the PR
- Include request/response schema references
- Include tags for grouping
- Have a one-line summary

## OpenAPI spec at a glance

- **91 paths** across 12 tag groups
- **22 schemas** (request/response models)
- Visit `/docs` in a browser for Swagger UI
