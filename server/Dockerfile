# HerbSt REST API Server Dockerfile
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . ./

# Build the binary
RUN go build -o herbst-web .
RUN go build -o herbst-migrate ./cmd/migrate

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/herbst-web .
COPY --from=builder /app/herbst-migrate .

# Expose API port
EXPOSE 8080

# Run with environment variables
CMD ["./herbst-web"]
