---
title: Initial Backend Server Scaffolding
id: 02
requires_bdd: true
is_completed: true
---

## Summary

The initial scaffolding for the project needs to be set up.

## Acceptance Criteria

- [x] Create a golang project in a folder called _server_
- [x] Modify the makefile to start the server.
- [x] Modify the makefile to have a `make dev` command that starts both the ssh and web servers.
- [x] It should start a _GIN_ web server with a simple `healthz` endpoint that returns a 200 status code if the ssh server is running.
- [x] It should include an openapi.json endpoint using a common openapi specification library.
- [x] Ensure that both servers can be started and stopped using the Makefile commands.
- [x] Update the Makefile to start both servers in 'make dev'
- [x] Write BDD tests to verify that the `healthz` endpoint returns the expected status code.

## Open Implementation Notes

- utilize <https://github.com/oapi-codegen/gin-middleware> for the openapi spec, but make sure the endpoint is still `openapi.json`