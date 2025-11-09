# makeathing-mud

## 0.1.3

### Minor Changes

- feat: Update command prompt to show hitpoints, mana, and experience

  The command prompt has been updated to display the player's current hitpoints (HP), mana (M), and experience (XP) instead of just the basic ">" prompt. This provides players with immediate visibility into their character's stats during gameplay.

  The prompt format is now: `HP:<health> M:<mana> XP:<experience>>`

  For example: `HP:30 M:0 XP:0>`

  If character data cannot be retrieved, the prompt falls back to default values: `HP:30 M:0 XP:0>`

## 0.1.2

### Patch Changes

- Added changesets support for managing changelog and versioning.

## 0.1.1

### Patch Changes

- Initial changeset for the MUD server with database persistence, authentication and global state tracking.
