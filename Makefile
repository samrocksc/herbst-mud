# Makefile for the MUD server

BINARY=mudserver
MAIN_DIR=cmd/mudserver

# Build the binary
build:
	go build -o ${BINARY} ./${MAIN_DIR}

# Run the server
run:
	go run ./${MAIN_DIR}

# Run the server with debug mode enabled
run-debug:
	DEBUG=true go run ./${MAIN_DIR}

# Install dependencies
deps:
	go mod tidy

# Clean build artifacts
clean:
	rm -f ${BINARY}

# Test the code
test:
	go test ./...

# Format the code
fmt:
	go fmt ./...

.PHONY: build run run-debug deps clean test fmt