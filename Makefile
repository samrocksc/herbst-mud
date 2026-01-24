.PHONY: help start stop run start-web stop-web dev test test-bdd test-server-bdd

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

start: ## Start the SSH server in the background
	@echo "Starting SSH server on port 4444..."
	@cd herbst && go run main.go &
	@echo $$! > .herbst.pid
	@echo "SSH server started with PID $$(cat .herbst.pid)"

stop: ## Stop both SSH and web servers
	@if [ -f .herbst.pid ]; then \
		echo "Stopping SSH server with PID $$(cat .herbst.pid)..."; \
		kill $$(cat .herbst.pid) 2>/dev/null || true; \
		rm -f .herbst.pid; \
		echo "SSH server stopped."; \
	else \
		echo "No running SSH server found."; \
	fi
	@if [ -f .web.pid ]; then \
		echo "Stopping web server with PID $$(cat .web.pid)..."; \
		kill $$(cat .web.pid) 2>/dev/null || true; \
		rm -f .web.pid; \
		echo "Web server stopped."; \
	else \
		echo "No running web server found."; \
	fi

run: ## Start the SSH server in the foreground
	@cd herbst && go run main.go

start-web: ## Start the web server in the background
	@echo "Starting web server on port 8080..."
	@cd server && go run main.go &
	@echo $$! > .web.pid
	@echo "Web server started with PID $$(cat .web.pid)"

stop-web: ## Stop the web server
	@if [ -f .web.pid ]; then \
		echo "Stopping web server with PID $$(cat .web.pid)..."; \
		kill $$(cat .web.pid) 2>/dev/null || true; \
		rm -f .web.pid; \
		echo "Web server stopped."; \
	else \
		echo "No running web server found."; \
	fi

dev: ## Start both SSH and web servers in the background
	@echo "Starting both SSH and web servers..."
	@cd herbst && go run main.go &
	@echo $$! > .herbst.pid
	@cd server && go run main.go &
	@echo $$! > .web.pid
	@echo "SSH server started with PID $$(cat .herbst.pid)"
	@echo "Web server started with PID $$(cat .web.pid)"

test: ## Run tests
	@cd herbst && go test ./...

test-bdd: ## Run BDD tests
	@cd herbst && go test -v ./... -run TestBDD

test-server: ## Run server tests
	@cd server && go test -v

test-server-bdd: ## Run server Gherkin BDD tests
	@cd server && go test -v -run TestFeatures