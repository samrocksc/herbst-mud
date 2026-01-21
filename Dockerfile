FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/mud-ssh ./cmd/ssh

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/mud-ssh .
COPY --from=builder /app/cmd/ssh/host_key ./host_key

EXPOSE 4444
CMD ["./mud-ssh"]
