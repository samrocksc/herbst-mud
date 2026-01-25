---
title: OPENAPI Spec Scaffolding
id: 04
requires_bdd: true
is_completed: true
---

## Summary

We need to use the openapi generator from the api in order to create a `client` in our frontend instead of manually building queries. When the frontend server starts it should introspect the open apispec from the backend

## Acceptance Criteria

- [ ] As a user I should be able to run make dev and have the frontend start, and then generate the client from the backend openapi spec
- [ ] It should be typesafe.

## Technical Guidance

Here is an example of how we should use the makefile to accomplish this, but utilize the `admin` file instead:

```makefile
dev-frontend:
 @echo "Starting backend server..."
 cd backend && . venv/bin/activate && uvicorn main:app --reload & \
 BACKEND_PID=$$! && \
 echo "Backend started with PID $$BACKEND_PID" && \
 sleep 5 && \
 echo "Generating frontend types from API..." && \
 cd frontend && npx @hey-api/openapi-ts -i http://localhost:8000/openapi.json -o src/client && \
 echo "Generated frontend types" && \
 echo "Starting frontend development server..." && \
 cd ../frontend && npm run dev
```
