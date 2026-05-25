.PHONY: help start stop run start-web stop-web dev start-admin stop-admin start-client stop-client dev-all test test-bdd test-server-bdd logs-ssh logs-web logs-client build build-all reload token

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

stop: ## Stop SSH, web, and admin servers
	@fuser -k 4444/tcp 2>/dev/null; \
	rm -f .herbst.pid; \
	echo "SSH server stopped."
	@fuser -k 8080/tcp 2>/dev/null; \
	rm -f .web.pid; \
	echo "Web server stopped."
	@[ -f .admin.pid ] && kill $$(cat .admin.pid) 2>/dev/null; \
	rm -f .admin.pid; \
	echo "Admin frontend stopped."
	@[ -f web-client/.client.pid ] && kill $$(cat web-client/.client.pid) 2>/dev/null; \
	rm -f web-client/.client.pid; \
	echo "Web client stopped."

run: ## Start the SSH server in the foreground (uses pre-built binary)
	@[ -f herbst/herbst ] || $(MAKE) build
	@herbst/herbst

start-web: ## Start the web server in the background (uses pre-built binary)
	@echo "Starting web server on port 8080..."
	@[ -f server/herbst-web ] || $(MAKE) build-web
	@fuser -k 8080/tcp 2>/dev/null; \
	(cd /home/sam/GitHub/herbst-mud && server/herbst-web > /tmp/herbst-web.log 2>&1) & \
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

start-client: ## Start the web client in the background
	@echo "Starting web client on port 5174..."
	@cd web-client && npm run dev > /tmp/herbst-client.log 2>&1 & \
		echo $$! > .client.pid && echo "Web client started with PID $$(cat .client.pid)"

stop-client: ## Stop the web client
	@[ -f web-client/.client.pid ] && kill $$(cat web-client/.client.pid) 2>/dev/null; \
	rm -f web-client/.client.pid; \
	echo "Web client stopped."

dev: ## Build and start both SSH + web servers + web client
	@echo "Building..."
	@$(MAKE) build-all
	@echo "Starting services..."
	@fuser -k 4444/tcp 8080/tcp 5174/tcp 2>/dev/null; sleep 1
	@herbst/herbst > /tmp/herbst-ssh.log 2>&1 & echo $$! > .herbst.pid
	@server/herbst-web > /tmp/herbst-web.log 2>&1 & echo $$! > .web.pid
	@cd web-client && npm run dev > /tmp/herbst-client.log 2>&1 & echo $$! > .client.pid
	@sleep 2
	@echo "SSH: $$(cat .herbst.pid) | Web: $$(cat .web.pid) | Client: $$(cat .client.pid)"
	@echo "Logs: make logs-ssh / make logs-web / make logs-client"

dev-all: ## Build and start all services (SSH + web + admin + web client)
	@echo "Building..."
	@$(MAKE) build-all
	@echo "Starting all services..."
	@fuser -k 4444/tcp 8080/tcp 5173/tcp 5174/tcp 2>/dev/null; sleep 1
	@herbst/herbst > /tmp/herbst-ssh.log 2>&1 & echo $$! > .herbst.pid
	@server/herbst-web > /tmp/herbst-web.log 2>&1 & echo $$! > .web.pid
	@cd admin && npm run dev > /tmp/herbst-admin.log 2>&1 & echo $$! > .admin.pid
	@cd web-client && npm run dev > /tmp/herbst-client.log 2>&1 & echo $$! > .client.pid
	@sleep 2
	@echo "SSH: $$(cat .herbst.pid) | Web: $$(cat .web.pid) | Admin: $$(cat .admin.pid) | Client: $$(cat .client.pid)"
	@echo "Logs: make logs-ssh / make logs-web / make logs-admin / make logs-client"

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

logs-admin: ## Tail admin frontend logs
	@tail -f /tmp/herbst-admin.log

logs-client: ## Tail web client logs
	@tail -f /tmp/herbst-client.log

token: ## Generate a Bearer token for API debugging
	@echo "Usage: make token USER_ID=1 EMAIL=example@test.com IS_ADMIN=true"
	@echo ""
	@cd server && go run cmd/token/main.go "${USER_ID}" "${EMAIL}" "${IS_ADMIN}"

# ── Production Deployment ────────────────────────────────────────

deploy-build: ## Build all production containers
	@echo "Building production containers..."
	@docker compose -f docker-compose.prod.yml build --no-cache
	@echo "Build complete"

deploy-up: ## Start production stack
	@docker compose -f docker-compose.prod.yml up -d
	@echo "Production stack running"
	@docker compose -f docker-compose.prod.yml ps

deploy-down: ## Stop production stack
	@docker compose -f docker-compose.prod.yml down
	@echo "Production stack stopped"

deploy-logs: ## Tail production logs
	@docker compose -f docker-compose.prod.yml logs -f

deploy-push: ## Push setup script to a droplet and run it (set DROPLET_IP)
	@if [ -z "${DROPLET_IP}" ]; then \
	  echo "Usage: make deploy-push DROPLET_IP=1.2.3.4"; \
	  exit 1; \
	fi
	@echo "Deploying to ${DROPLET_IP}..."
	@scp deploy/setup.sh root@${DROPLET_IP}:/root/
	@ssh root@${DROPLET_IP} 'bash /root/setup.sh'

deploy-update: ## Pull latest code and rebuild production stack
	@git fetch origin main
	@git reset --hard origin/main
	@docker compose -f docker-compose.prod.yml build --no-cache
	@docker compose -f docker-compose.prod.yml up -d
	@echo "Updated and restarted"
