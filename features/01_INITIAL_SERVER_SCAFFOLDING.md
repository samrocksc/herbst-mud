---
title: Initial SSH MUD Server Scaffolding
id: 01
requires_bdd: true
is_completed: true
---

## Summary

The initial scaffolding for the MUD ssh server needs to be set up.

## Acceptance Criteria

- [x] Create a golang project in a folder called _herbst_
- [x] initialize a Makefile that will start projects.
- [x] It should start an ssh server with the charmland wish library. the port should be 4444, it should not require ANY authentication, and should log connections.
- [x] Utilize bubbletea for the ssh server UI
- [x] Ensure that the server can be started and stopped using the Makefile commands.
- [x] Write BDD tests to verify the server can be ssh'd into
