# MUD - Multi-User Dungeon

A text-based multiplayer game with SSH connectivity and React admin panel.

## Quick Start

### Option 1: Run directly
```bash
./bin/mud-ssh
```

### Option 2: Docker Compose
```bash
docker-compose up
```

## SSH Connection

Connect using:
```bash
ssh mud@localhost -p 4444
```
Password: mud

## Development

### SSH Server
- Location: cmd/ssh/
- Tests: cmd/ssh/main_test.go
- Compile: go build -o bin/mud-ssh ./cmd/ssh

### Admin Panel
- Location: admin/
- Framework: Vite + React
- Port: 3000

### Testing
```bash
# SSH server tests
cd cmd/ssh && go test -v

# Build and run
go build -o bin/mud-ssh ./cmd/ssh
./bin/mud-ssh
```

## Architecture

- Go SSH server using native crypto library
- Bubbletea for terminal UI
- Simple authentication (username: mud, password: mud)
- Clean separation between cmd/, internal/, and pkg/

## Features Implemented

✅ SSH server on port 4444
✅ Bubbletea terminal interface
✅ Basic menu system
✅ Tests for SSH functionality
✅ Docker development environment
✅ React admin panel scaffolding
