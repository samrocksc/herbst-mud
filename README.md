# Makeathing - A MUD Game Server

A Multi-User Dungeon (MUD) game server written in Go with SSH capability using the charm `wish` library.

## Why do this?

Basically, I want to learn Go, and I love MUDs. I think this is a great way to spend my time. I am heavily using `crush` to build this.

## Features

- SSH-based multiplayer access
- Character system with races and classes
- Room navigation with cardinal directions
- Item system with movable and immovable objects
- Combat system with stats-based mechanics
- Extensible architecture with adapter pattern

## Technical Specifications

Based on the requirements in `rules.md`:

- Server built with Go and charm's `wish` library
- SSH accessible with lipgloss for UI
- Adapter-based system for different connection types
- BDD testing approach (planned)

## Character System

- Player, Admin, and NPC character types
- Races: Human, Rat People, Dwarf
- Classes: Warrior, Mage, Rogue
- Stats: Strength, Intelligence, Dexterity (1-25 range)

## Room System

- Cardinal direction navigation
- Immovable and movable objects
- Smell descriptions
- First 4 rooms as specified in requirements

## Getting Started

### Prerequisites

- Go 1.19 or higher

### Building and Running

```bash
# Clone the repository
git clone <repository-url>
cd makeathing

# Build the server
make build

# Run the server
make run
```

### Connecting

Connect to the server using any SSH client:

```bash
ssh localhost -p 2222
```

When connecting for the first time, you may see a security warning about the host key. This is normal for a new SSH server.

Available commands:

- `help` - Show available commands
- `look` - Look around the room
- `quit`/`exit` - Exit the game

## Project Structure

```
.
├── cmd/
│   └── mudserver/          # Main server application
├── internal/
│   ├── adapters/           # Connection adapters (SSH, etc.)
│   ├── characters/         # Character system
│   ├── rooms/              # Room system
│   ├── items/              # Item system
│   ├── combat/             # Combat system
│   └── actions/            # Actions system
├── .ssh/                   # SSH keys for the server
├── rules.md                # Game requirements
├── Makefile                # Build automation
└── README.md               # This file
```

## Testing

Run tests with:

```bash
make test
```

## Future Enhancements

- Inventory management
- Combat system implementation
- Character creation
- World expansion
- BDD testing with Gherkin
- Database persistence
- More character races and classes
- Quest system
- Vendor NPCs

