
.PHONY: build run test clean docker-build docker-run

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run:
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Docker build
docker-build:
	docker build -t ddd-crud-server .

# Docker run
docker-run:
	docker run -p 8080:8080 --env-file .env ddd-crud-server

# Run database migrations
migrate:
	./scripts/migrate.sh

# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/server .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./server"]
