---
title: Initial Backend Server Scaffolding
id: 05
requires_bdd: false
is_completed: false
---

## Summary

We want to create a database that will hold the configuration of our application. We want to utilize [ent](https://entgo.io/) as our ORM and a relational database postgres PostgreSQL.

## Acceptance Criteria

- [ ] We will use ent in a folder named `db`, it will generate the ORM client for the server, and also for herbst, so they will need to have clients generated for each.
- [ ] We will use PostgreSQL as our database.
- [ ] The database will be set up with the `docker-compose.yml` file to run a postgres instance.
- [ ] Run migrations on the server start to ensure the status of the database.
- [ ] It should have a `users` table. It should have a one to many relationship with characters
- [ ] It should have a `characters` table. It should have an `isNPC` boolean
- [ ] It should have a `rooms` table. It should have an `exits` field which should represent the direction, and roomId of the room in that direction.
- [ ] We should start with a simple cross way called "the hole" which has 5 rooms, in a cross pattern.
- [ ] Each room should have a description field.
- [ ] Each room should have a name field.
- [ ] Each room should have a `isStartingRoom` boolean field to indicate if it's a starting room.
- [ ] Each character should have a `currentRoomId` field to indicate which room they are

## Technical Guidance

Utilize the ent library to create the client for herbst
