---
title: Creating ROOM CRUD Operations
id: 06
requires_bdd: false
is_completed: true
---

## Summary

We need to implement CRUD (Create, Read, Update, Delete) operations for managing rooms in our application. This will allow users to create new rooms, view existing rooms, update room information, and delete rooms as needed.

## Acceptance Criteria

- [x] Implement a POST endpoint to create a new room. The endpoint should accept room details such as name and description.
- [x] Implement a GET endpoint to retrieve a list of all rooms. The endpoint should return room details including name and description.
- [x] Implement a GET endpoint to retrieve a single room by its ID. The endpoint should return room details including name and description.
- [x] Implement a PUT endpoint to update an existing room's information. The endpoint should accept room details such as name and description.
- [x] Implement a DELETE endpoint to remove a room by its ID.
- [x] Create a seed file that creates four rooms in a cross pattern:
  - The Hole (starting room)
  - North Room
  - South Room
  - East Room
  - West Room
