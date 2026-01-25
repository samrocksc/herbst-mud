---
title: Creating User CRUD Operations
id: 08
requires_bdd: false
is_completed: false
---

## Summary

We need to implement CRUD (Create, Read, Update, Delete) operations for managing users in our application. This will allow us to create new users, view existing users, update user information, and delete users as needed.

## Acceptance Criteria

- [ ] Ensure that the user has an `is_admin` boolean field to indicate if they are an admin or not.
- [ ] Implement a POST endpoint to create a new user. The endpoint should accept user details such as username and password.
- [ ] Implement a GET endpoint to retrieve a list of all users. The endpoint should return user details including username.
- [ ] Implement a GET endpoint to retrieve a single user by its ID. The endpoint should return user details including username.
- [ ] Implement a PUT endpoint to update an existing user's information. The endpoint should accept user details such as username and password.
- [ ] Implement a DELETE endpoint to remove a user by its ID.
- [ ] The user should have a one to many relationship with characters.
- [ ]
