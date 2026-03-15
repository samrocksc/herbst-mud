# Herbst MUD

A modern MUD (Multi-User Dungeon) game built with Go, TypeScript, and PostgreSQL.

## Features

- SSH-based multiplayer gameplay
- REST API for administration and game management
- Web-based admin panel
- PostgreSQL database with ent ORM
- Docker-based deployment

## Project Structure

- `herbst/` - SSH client implementation
- `server/` - REST API server
- `admin/` - Web-based admin panel
- `db/` - Database setup and migrations
- `docs/` - Documentation
- `features/` - Feature specifications

## Quick Start

1. Start all services:
   ```bash
   make dev-all
   ```

2. Connect to the MUD via SSH:
   ```bash
   ssh -p 4444 localhost
   ```

3. Access the admin panel at http://localhost:3000

4. Access the API at http://localhost:8080

## Database Setup

The project uses PostgreSQL with ent ORM. The database is automatically initialized when the servers start.

## Development

### Prerequisites

- Go 1.25+
- Node.js 18+
- Docker and Docker Compose
- PostgreSQL client

### Available Make Commands

- `make dev-all` - Start all services
- `make dev` - Start SSH and web servers
- `make start` - Start SSH server
- `make start-web` - Start web server
- `make start-admin` - Start admin frontend
- `make test` - Run tests
- `make test-bdd` - Run BDD tests

## Documentation

See the `docs/` directory for detailed documentation on various aspects of the project.# User CRUD Feature
