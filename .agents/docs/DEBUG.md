# Debug Mode

The MUD client includes a debug mode that provides verbose logging and additional information during gameplay. When enabled, debug mode displays extra diagnostic messages for commands like `look` and shows the current Room ID in the status bar.

## Commands

| Command | Description |
|---------|-------------|
| `debug` | Toggle or check debug mode status |
| `debug on` | Enable debug mode |
| `debug off` | Disable debug mode |

Alternate values: `true/false`, `1/0`, `yes/no`

## Features

When debug mode is active:
- **Room ID** is displayed in the status bar (e.g., `[Room: 5]`)
- **Look command** outputs verbose matching information showing:
  - The target being searched
  - Number of room items and characters loaded
  - Each item/character being checked
  - Match results

## Example

```
> debug on
✓ Debug mode: ON (Room ID will show in status bar)

> look gizmo
[DEBUG] Looking for: 'gizmo'
[DEBUG] Room items: 2, Room characters: 1
[DEBUG] Checking item: 'Rusty Sword'
[DEBUG] Checking item: 'Old Scroll'
[DEBUG] Checking character: 'Gizmo' (IsNPC: true)
[DEBUG] Matched character: 'Gizmo'
[Gizmo]
A small, curious creature...
```