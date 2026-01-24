# MUD - Multi-User Dungeon

A text-based multiplayer game with SSH connectivity and React admin panel.

## 🚀 Quick Start

### Option 1: Automated Setup (Recommended)
```bash
make install && make dev
```

### Option 2: Docker Development
```bash
docker-compose up
```

### Option 3: Manual Setup
```bash
# Install dependencies
make install

# Start services
make dev
```

## 📋 Available Makefile Commands

### Setup & Installation
- `make help` - Show all available commands
- `make install` - Install all dependencies (Go and Node.js)
- `make build` - Build all components
- `make clean` - Clean build artifacts

### Development
- `make dev` - Start SSH server + Admin panel
- `make ssh-server` - Start SSH server only
- `make admin` - Start Admin panel only
- `make stop` - Stop all development services
- `make status` - Check development status

### Testing
- `make test` - Run all tests

### Docker
- `make docker-dev` - Start with Docker Compose
- `make docker-build` - Build Docker images
- `make docker-clean` - Clean Docker resources

## 🔌 Connection Details

### SSH Server
- **Port**: 4444
- **Connection**: `ssh localhost -p 4444`
- **Password**: No password required
- **Authentication**: Passwordless SSH (authentication handled by game)

### Admin Panel
- **URL**: http://localhost:3000
- **Technology**: Vite + React + TypeScript
- **Features**: TanStack Router/Query/Forms

## 📁 Project Structure

```
a-mud/
├── Makefile              # Development automation
├── docker-compose.yml    # Docker development
├── Dockerfile            # MUD SSH server image
├── cmd/ssh/              # SSH server implementation
│   ├── main.go           # Main SSH server with Bubbletea
│   ├── main_test.go      # Comprehensive tests
│   └── host_key          # SSH host key
├── admin/                # React admin panel
│   ├── package.json      # Dependencies and scripts
│   ├── tsconfig.json     # TypeScript configuration
│   └── src/              # TypeScript source files
├── bin/                  # Compiled binaries
└── docs/                 # Project documentation
```

## 🛠️ Technology Stack

### Backend (SSH Server)
- **Go 1.21+** - Systems programming language
- **Bubbletea** - Terminal UI framework
- **golang.org/x/crypto/ssh** - SSH server implementation
- **Standard testing** - Go testing library

### Frontend (Admin Panel)
- **TypeScript 5.0+** - Type-safe JavaScript
- **Vite** - Fast development build tool
- **React 18** - UI library
- **TanStack Router** - Modern routing
- **TanStack Query** - Data fetching
- **TanStack Forms** - Form management

## 🧪 Testing

### SSH Server Tests
```bash
make test
cd cmd/ssh && go test -v
```

### Features Implemented

✅ SSH server on port 4444
✅ Passwordless SSH authentication
✅ Welcome screen on connection
✅ Bubbletea terminal interface
✅ Tests for SSH functionality
✅ Docker development environment
✅ React admin panel with TypeScript
✅ TanStack Router/Query/Forms
✅ Makefile automation
✅ Comprehensive documentation
