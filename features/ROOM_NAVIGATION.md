---
title: Room Navigation
id: ROOM_NAVIGATION
requires_bdd: false
is_completed: true
---

## Summary

Players need to be able to navigate between rooms in the MUD using cardinal directions (north, south, east, west).

## Background

The MUD has a cross-shaped room layout with 5 rooms:
- Room 1: North room
- Room 2: South room
- Room 3: East room
- Room 4: West room
- Room 5 (Starting): The Hole - center hub

Each room has exits defined in the database.

## Implementation

### SSH Server Changes

Update `herbst/sshserver.go` to handle movement commands:

1. Track player's current room in session state
2. Add command handlers for: `n`, `north`, `s`, `south`, `e`, `east`, `w`, `west`
3. Look up exit from current room
4. Move player to new room if exit exists
5. Display new room description on entry

### Database

No changes needed - exits already stored in Room entity as JSON map.

## Acceptance Criteria

- [x] Player can type "north" or "n" to move north
- [x] Player can type "south" or "s" to move south
- [x] Player can type "east" or "e" to move east
- [x] Player can type "west" or "w" to move west
- [x] Moving displays the new room name and description
- [x] Trying to move in a direction with no exit shows "You can't go that way."
- [x] Player starts in the starting room (The Hole - ID 5)

## Implementation Notes

- Implemented in `herbst/main.go`
- Uses BubbleTea model with room state tracking
- Commands: n/north, s/south, e/east, w/west, look, exits, help