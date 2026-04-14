.PHONY: help start stop run start-web stop-web dev start-admin stop-admin dev-all test test-bdd test-server-bdd logs-ssh logs-web build build-all reload build-admin-tui run-admin-tui

PATH := $(PATH):/usr/local/go/bin
export PATH

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build SSH server binary
	@echo "Building SSH server..."
	@cd herbst && go build -o herbst . && echo "SSH binary built"

build-web: ## Build web server binary
	@echo "Building web server..."
	@cd server && go build -o herbst-web . && echo "Web binary built"

build-all: ## Build all server binaries
	@$(MAKE) build
	@$(MAKE) build-web

start: ## Start the SSH server in the background (uses pre-built binary)
	@echo "Starting SSH server on port 4444..."
	@[ -f herbst/herbst ] || $(MAKE) build
	@fuser -k 4444/tcp 2>/dev/null; \
	herbst/herbst > /tmp/herbst-ssh.log 2>&1 & \
	echo $$! > .herbst.pid && echo "SSH server started with PID $$(cat .herbst.pid)"

stop: ## Stop both SSH and web servers
	@fuser -k 4444/tcp 2>/dev/null; \
	rm -f .herbst.pid; \
	echo "SSH server stopped."
	@fuser -k 8080/tcp 2>/dev/null; \
	rm -f .web.pid; \
	echo "Web server stopped."

run: ## Start the SSH server in the foreground (uses pre-built binary)
	@[ -f herbst/herbst ] || $(MAKE) build
	@herbst/herbst

start-web: ## Start the web server in the background (uses pre-built binary)
	@echo "Starting web server on port 8080..."
	@[ -f server/herbst-web ] || $(MAKE) build-web
	@fuser -k 8080/tcp 2>/dev/null; \
	server/herbst-web > /tmp/herbst-web.log 2>&1 & \
	echo $$! > .web.pid && echo "Web server started with PID $$(cat .web.pid)"

stop-web: ## Stop the web server
	@fuser -k 8080/tcp 2>/dev/null; \
	rm -f .web.pid; \
	echo "Web server stopped."

start-admin: ## Start the admin frontend in the background
	@echo "Starting admin frontend..."
	@cd admin && npm run dev &
	@echo $$! > .admin.pid && echo "Admin frontend started with PID $$(cat .admin.pid)"

stop-admin: ## Stop the admin frontend
	@[ -f .admin.pid ] && kill $$(cat .admin.pid) 2>/dev/null; \
	rm -f .admin.pid; \
	echo "Admin frontend stopped."

dev: ## Build and start both SSH + web servers (uses pre-built binaries)
	@echo "Building..."
	@$(MAKE) build-all
	@echo "Starting services..."
	@fuser -k 4444/tcp 8080/tcp 2>/dev/null; sleep 1
	@herbst/herbst > /tmp/herbst-ssh.log 2>&1 & echo $$! > .herbst.pid
	@server/herbst-web > /tmp/herbst-web.log 2>&1 & echo $$! > .web.pid
	@sleep 2
	@echo "SSH: $$(cat .herbst.pid) | Web: $$(cat .web.pid)"
	@echo "Logs: make logs-ssh / make logs-web"

dev-all: ## Build and start all services (SSH + web + admin)
	@echo "Building..."
	@$(MAKE) build-all
	@echo "Starting all services..."
	@fuser -k 4444/tcp 8080/tcp 2>/dev/null; sleep 1
	@herbst/herbst > /tmp/herbst-ssh.log 2>&1 & echo $$! > .herbst.pid
	@server/herbst-web > /tmp/herbst-web.log 2>&1 & echo $$! > .web.pid
	@cd admin && npm run dev &
	@echo $$! > .admin.pid
	@sleep 2
	@echo "SSH: $$(cat .herbst.pid) | Web: $$(cat .web.pid) | Admin: $$(cat .admin.pid)"

build-admin-tui: ## Build the admin TUI binary
	@echo "Building admin TUI..."
	@cd admin-tui && go build -o admin-tui . && echo "Admin TUI binary built"

run-admin-tui: ## Run the admin TUI (requires dev stack running on localhost:8080)
	@echo "Starting admin TUI..."
	@API_BASE_URL=http://localhost:8080 ./admin-tui/admin-tui

reload: ## Rebuild SSH binary and restart (hot reload)
	@echo "Building..."
	@$(MAKE) build
	@echo "Restarting SSH server..."
	@fuser -k 4444/tcp 2>/dev/null; sleep 1
	@herbst/herbst > /tmp/herbst-ssh.log 2>&1 & echo $$! > .herbst.pid
	@sleep 2 && tail -2 /tmp/herbst-ssh.log

reload-web: ## Rebuild web binary and restart
	@echo "Building..."
	@$(MAKE) build-web
	@echo "Restarting web server..."
	@fuser -k 8080/tcp 2>/dev/null; sleep 1
	@server/herbst-web > /tmp/herbst-web.log 2>&1 & echo $$! > .web.pid
	@sleep 2 && curl -s http://localhost:8080/healthz

test: ## Run tests
	@cd herbst && go test ./...

test-bdd: ## Run BDD tests
	@cd herbst && go test -v ./... -run TestBDD

test-server: ## Run server tests
	@cd server && go test -v

test-server-bdd: ## Run server Gherkin BDD tests
	@cd server && go test -v -run TestFeatures

logs-ssh: ## Tail SSH server logs
	@tail -f /tmp/herbst-ssh.log

logs-web: ## Tail web server logs
	@tail -f /tmp/herbst-web.log
