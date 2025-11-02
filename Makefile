# Makefile for the MUD server

BINARY=mudserver
MAIN_DIR=cmd/mudserver
VERSION=$(shell cat VERSION)

# Build the binary for current platform
build:
	go build -o ${BINARY} ./${MAIN_DIR}

# Build for all supported platforms
build-all: build-linux build-darwin-arm64 build-darwin-amd64

# Build for Linux (Ubuntu)
build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/${BINARY}-linux-amd64 ./${MAIN_DIR}

# Build for macOS ARM64 (Apple Silicon)
build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o dist/${BINARY}-darwin-arm64 ./${MAIN_DIR}

# Build for macOS AMD64 (Intel)
build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o dist/${BINARY}-darwin-amd64 ./${MAIN_DIR}

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
	rm -rf dist/

# Test the code
test:
	go test ./...

# Format the code
fmt:
	go fmt ./...

.PHONY: build build-all build-linux build-darwin-arm64 build-darwin-amd64 run run-debug deps clean test fmt