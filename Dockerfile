# HerbSt SSH Server Dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY herbst/go.mod herbst/go.sum ./
RUN go mod download

# Copy source code
COPY herbst/ ./

# Build the binary
RUN go build -o herbst .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS (needed for API calls)
RUN apk --no-cache add ca-certificates openssh

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/herbst .

# SSH host key will be generated at runtime or mounted
RUN mkdir -p .ssh

# Expose SSH port
EXPOSE 4444

# Run with environment variables
CMD ["./herbst"]
