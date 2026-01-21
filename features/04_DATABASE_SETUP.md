We want to create a database that will hold the configuration of our application.

- It should have a `users` table.  It should have a one to many relationship with characters
- It should have a `characters` table.  It should have an `isNPC` boolean
- It should have a `rooms` table.  It should have an `exits` field which should represent the direction, and roomId of the room in that direction.
- We should start with a simple cross way called "the hole" which has 5 rooms, in a cross pattern.