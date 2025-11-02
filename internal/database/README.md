# Database Implementation

This directory contains the SQLite database implementation for the MUD server.

## Overview

The database implementation provides:

1. **Migration System** - Automated database schema migrations
2. **Repository Pattern** - Clean data access layer for each entity
3. **Database Adapter** - Integration with the existing game codebase

## Entities

### Configuration
- Stores game configuration settings
- Currently only stores the name of the MUD

### Sessions
- Tracks connected users
- Stores session ID, user ID, character ID, and room ID

### Users
- Stores user information
- Links characters to rooms

### Rooms
- Stores room data mirroring the JSON structure
- Complex nested data (exits, items, NPCs) stored as JSON-serialized strings
- All existing JSON rooms have been uploaded to the database

### Characters
- Stores character data mirroring the JSON structure
- Complex nested data (stats, inventory, skills) stored as JSON-serialized strings
- All existing JSON characters have been uploaded to the database

### Items
- Stores item data mirroring the JSON structure
- Complex nested data (stats) stored as JSON-serialized strings
- All existing JSON items have been uploaded to the database

### Actions
- Stores action data for available game actions
- Complex nested data (requirements) stored as JSON-serialized strings

## Usage

To use the database in your code:

```go
// Create a new database adapter
dbAdapter, err := database.NewDBAdapter("./data/mud.db")
if err != nil {
    log.Fatal(err)
}
defer dbAdapter.Close()

// Set a configuration value
err = dbAdapter.SetConfiguration("mud_name", "My Awesome MUD")
if err != nil {
    log.Fatal(err)
}

// Get a configuration value
config, err := dbAdapter.GetConfiguration("My Awesome MUD")
if err != nil {
    log.Fatal(err)
}

// Create a room from JSON
err = dbAdapter.CreateRoom(roomJSON)
if err != nil {
    log.Fatal(err)
}
```

## Migration System

The database automatically applies migrations when connecting. Migrations are defined in `migrations.go` and are applied in order. Each migration is only applied once.

To add a new migration:
1. Add a new entry to the `migrations` slice in `migrations.go`
2. Give it a unique name following the pattern `XXX_description`
3. Write the SQL for the migration

## Testing

Run database tests with:

```bash
go test ./internal/database
```