.PHONY: help start stop run

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

stop: ## Stop the SSH server
	@if [ -f .herbst.pid ]; then \
		echo "Stopping SSH server with PID $$(cat .herbst.pid)..."; \
		kill $$(cat .herbst.pid) 2>/dev/null || true; \
		rm -f .herbst.pid; \
		echo "SSH server stopped."; \
	else \
		echo "No running server found."; \
	fi

run: ## Start the SSH server in the foreground
	@cd herbst && go run main.go

test: ## Run tests
	@cd herbst && go test ./...

test-bdd: ## Run BDD tests
	@cd herbst && go test -v ./... -run TestBDD