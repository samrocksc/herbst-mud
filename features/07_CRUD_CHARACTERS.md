---
title: Creating Characters CRUD Operations
id: 06
requires_bdd: false
is_completed: false
---

## Summary

We need to implement CRUD (Create, Read, Update, Delete) operations for managing characters in our application. This will allow users to create new characters, view existing characters, update character information, and delete characters as needed.

## Acceptance Criteria

- [ ] Implement a POST endpoint to create a new character. The endpoint should accept character details such as name, isNPC status, and currentRoomId, startingRoomId.
- [ ] Implement a GET endpoint to retrieve a list of all characters. The endpoint should return character details including name, isNPC status, and currentRoomId.
- [ ] Implement a GET endpoint to retrieve a single character by its ID. The endpoint should return character details including name, isNPC status, and currentRoomId.
- [ ] Implement a PUT endpoint to update an existing character's information. The endpoint should accept character details such as name, isNPC status, and currentRoomId.
- [ ] Implement a DELETE endpoint to remove a character by its ID.
- [ ] The character should have a many to one relationship with users.
- [ ] The character should have a currentRoomId field to indicate which room they are currently
- [ ] modify the seed file with five new characters, and put them in rnadom rooms from the current seed file.
