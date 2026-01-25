---
title: Creating a login screen during the ssh initiation
id: 08
requires_bdd: false
is_completed: false
---

## Summary

We need to implement way for the users to be able to login during the ssh initiation. This will allow us to authenticate users before they can access the game.

## Acceptance Criteria

- [ ] Implement a POST endpoint to create a new user. The endpoint should accept user details such as username and password.
- [ ] Implement a GET endpoint to retrieve a list of all users. The endpoint should return user details including username.
- [ ] Implement a GET endpoint to retrieve a single user by its ID. The endpoint should return user details including username.
- [ ] Implement a PUT endpoint to update an existing user's information. The endpoint should accept user details such as username and password.
- [ ] Implement a DELETE endpoint to remove a user by its ID.
- [ ] The user should have a one to many relationship with characters.
