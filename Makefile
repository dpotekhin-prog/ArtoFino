export PATH := $(PATH):$(HOME)/go/bin

APP_NAME=server
MAIN=cmd/server/main.go
DOCS_DIR=docs

.PHONY: build run swag docker-up docker-down lint air

# Build Go binary
build:
	go build -o $(APP_NAME) $(MAIN)

# Run backend locally
run:
	go run $(MAIN)

# Generate Swagger docs
swag:
	swag init -g $(MAIN) -o $(DOCS_DIR)

# Start docker-compose
docker-up:
	docker compose up -d --build

# Stop docker-compose
docker-down:
	docker compose down

# Lint (requires golangci-lint)
lint:
	golangci-lint run

# Hot reload with Air
air:
	air
