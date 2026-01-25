---
title: Initial Backend Server Scaffolding
id: 05
requires_bdd: false
is_completed: true
---

## Summary

We want to create a database that will hold the configuration of our application. We want to utilize [ent](https://entgo.io/) as our ORM and a relational database postgres PostgreSQL.
l

## Acceptance Criteria

- [x] We will use ent in a folder named `db`, it will generate the ORM client for the server, and also for herbst, so they will need to have clients generated for each.
- [x] We will use PostgreSQL as our database.
- [x] The database will be set up with the `docker-compose.yml` file to run a postgres instance.
- [x] Run migrations on the server start to ensure the status of the database.
- [x] It should have a `users` table. It should have a one to many relationship with characters
- [x] It should have a `characters` table. It should have an `isNPC` boolean
- [x] It should have a `rooms` table. It should have an `exits` field which should represent the direction, and roomId of the room in that direction.
- [x] We should start with a simple cross way called "the hole" which has 5 rooms, in a cross pattern.
- [x] Each room should have a description field.
- [x] Each room should have a name field.
- [x] Each room should have a `isStartingRoom` boolean field to indicate if it's a starting room.
- [x] Each character should have a `currentRoomId` field to indicate which room they are

## Technical Guidance

Utilize the ent library to create the client for herbst and server.
