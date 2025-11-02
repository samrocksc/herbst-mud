# CRUSH.md - Development Guide for Makeathing MUD Server

This document provides essential information for agents working with the Makeathing MUD (Multi-User Dungeon) server codebase.

## Project Overview

A MUD game server written in Go with SSH capability using the charm `wish` library. The server allows multiple players to connect via SSH and navigate a text-based world with rooms, items, characters, and combat.

## Essential Commands

```bash
# Build the server for current platform
make build

# Run the server (with database storage - default)
make run

# Run the server with JSON storage only
make run-json

# Run the server with database storage explicitly
make run-db

# Run the server with debug mode enabled
make run-debug

# Run the server with debug mode and JSON storage only
make run-debug-json

# Build for all supported platforms
make build-all

# Run tests
make test

# Format code
make fmt

# Install/update dependencies
make deps

# Clean build artifacts
make clean
```

## Code Organization

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
│   ├── actions/            # Actions system
│   ├── configuration/      # Configuration system
│   └── database/           # Database implementation
├── data/
│   ├── items/              # JSON item definitions
│   ├── rooms/              # JSON room definitions
│   ├── characters/         # JSON character definitions
│   ├── configuration.json  # JSON configuration
│   └── schemas/            # JSON schemas for validation
├── .ssh/                   # SSH keys for the server
├── Makefile                # Build automation
└── README.md               # Project documentation
```

## Architecture Patterns

### Core Components

1. **Game Engine**: Central coordinator in `cmd/mudserver/main.go`
2. **Adapters**: Connection handling (SSH) in `internal/adapters/`
3. **Domain Objects**:
   - Rooms in `internal/rooms/`
   - Characters in `internal/characters/`
   - Items in `internal/items/`
   - Combat in `internal/combat/`

### Data Persistence

- Uses JSON files for initial data loading in the `data/` directory
- JSON Schema validation for data integrity (schemas in `data/schemas/`)
- Each entity type (rooms, items, characters, configuration) has its own directory or file with JSON files
- References between entities are resolved at load time
- SQLite database for runtime persistence of configuration, sessions, and users
- Automated migration system for database schema updates

### Session Management

- Sessions are managed by `SessionManager` in `internal/adapters/session.go`
- Each connected player has a `PlayerSession` tracking their state
- Thread-safe with mutex protection for concurrent access

### Adapter Pattern

- `Adapter` interface in `internal/adapters/adapter.go` defines connection handling
- `SSHAdapter` implements the interface for SSH connections
- Processes commands and manages user interaction

### Database Implementation

The MUD server now includes a SQLite-based database implementation for runtime persistence:

- **Database Package**: Located in `internal/database/` with repository pattern implementation
- **Migration System**: Automated schema management with version tracking
- **Core Tables**:
  - `configuration`: Game configuration settings (e.g., MUD name)
  - `sessions`: Active user sessions with session IDs, user IDs, character IDs, and room IDs
  - `users`: User accounts linking characters to rooms
  - `rooms`: Room data mirroring JSON structure with JSON-serialized complex fields
- **Database Adapter**: Integration layer between database and game logic in `internal/database/adapter.go`
- **Repository Pattern**: Clean data access layer with separate repositories for each entity type
- **Configuration Loading**: Initial configuration loaded from `data/configuration.json` at startup

## Key Data Structures

### Room

- Defined in `internal/rooms/room.go`
- Contains description, exits, objects, NPCs, and smells
- Exits are directional (north, south, east, west, etc.)
- Supports both movable and immovable objects

### Character

- Defined in `internal/characters/character.go`
- Has race, class, stats (strength, intelligence, dexterity), health, mana
- Supports levels, vendor status, and NPC flags

### Item

- Defined in `internal/items/item.go`
- Has ID, name, description, and movability flag

### Configuration

- Defined in `internal/configuration/configuration.go`
- Has ID and name fields that mirror the database configuration table
- Stored as JSON file in `data/configuration.json`
- Loaded at server startup for initial configuration

## Testing Approach

- Unit tests in `*_test.go` files
- Integration tests in `*_integration_test.go` files
- Uses Go's standard testing package
- Run all tests with `make test`

## Naming Conventions

- Go standard conventions: PascalCase for exported, camelCase for unexported
- JSON field names use snake_case
- Interface names are often suffixed with "er" (e.g., Adapter)
- Constants use SCREAMING_SNAKE_CASE

## Code Style Guidelines

From `agents.md`:

- **Imports**: Use goimports formatting, group stdlib, external, internal packages
- **Formatting**: Use gofumpt (stricter than gofmt)
- **Error handling**: Return errors explicitly, use `fmt.Errorf` for wrapping
- **Context**: Always pass context.Context as first parameter for operations
- **Interfaces**: Define interfaces in consuming packages, keep them small and focused
- **JSON tags**: Use snake_case for JSON field names
- **Comments**: End comments in periods unless comments are at the end of the line

## Important Gotchas

1. **SSH Input Handling**: Custom line reading logic in `SSHAdapter.HandleConnection` to properly handle both Unix (\n) and Windows (\r\n) line endings

2. **Session Management**: Thread-safe access to sessions using read/write mutexes

3. **Data Loading**: Items and characters are loaded first, then referenced in rooms at load time

4. **Command Abbreviations**: Single-letter movement commands (n, s, e, w) are expanded to full directions

5. **Debug Mode**: Enable with `DEBUG=true` environment variable for verbose logging

## Common Development Tasks

### Adding a New Room

1. Create a new JSON file in `data/rooms/` following the room schema
2. Define exits to other rooms by their IDs
3. Reference existing items/characters or create new ones
4. Test by running the server and navigating to the new room

### Adding a New Command

1. Modify the `processCommand` function in `internal/adapters/adapter.go`
2. Add command help text in the "help" case
3. Implement command logic
4. Test with `make run` and connecting via SSH

### Adding a New Item/Character

1. Create a new JSON file in the appropriate directory (`data/items/` or `data/characters/`)
2. Follow the respective schema in `data/schemas/`
3. Reference the item/character in rooms as needed

### Adding a New Configuration

1. Modify `data/configuration.json` with the desired settings
2. The configuration is loaded at server startup
3. For runtime changes, use the database configuration system

### Uploading JSON Rooms to Database

1. All existing JSON rooms in `data/rooms/` have been uploaded to the database
2. New rooms added as JSON files can be uploaded using database adapter methods
3. The rooms table stores complex nested data as JSON-serialized strings
4. Use `DBAdapter.CreateRoom()` to upload individual rooms
5. Existing rooms can be updated with `DBAdapter.UpdateRoom()`

### Adding a New Database Migration

1. Add a new entry to the `migrations` slice in `internal/database/migrations.go`
2. Follow the naming pattern `XXX_description` where XXX is a sequential number
3. Write the SQL for the migration
4. The migration will be automatically applied when the server starts

### Adding a New Database Entity

1. Create a new table in a migration file
2. Add a new repository file in `internal/database/` (e.g., `entity.go`)
3. Implement repository methods for Create, Get, Update, Delete operations
4. Add the repository to the `DBAdapter` struct and initialization
5. Test with comprehensive unit tests
6. Add new CLI methods to `DBAdapter` for easy access

## Deployment

The server listens on port 2222 for SSH connections. Connect with:
```bash
ssh localhost -p 2222
```

Cross-platform binaries can be built with:
```bash
make build-all
```

## Debugging

Enable debug mode to get verbose logging:
```bash
make run-debug
# or
DEBUG=true go run ./cmd/mudserver
```
