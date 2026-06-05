export PATH := $(PATH):$(HOME)/go/bin

APP_NAME=backend/server
MAIN=cmd/server/main.go
DOCS_DIR=docs

.PHONY: env-up env-down swag run dev init-env

# Check if backend/.env exists, if not - copy from .env.local
init-env:
	@test -f backend/.env || (cp backend/.env.local backend/.env && echo "Created backend/.env from .env.local")

# Start infrastructure (Databases + Keycloak) in Docker
env-up:
	docker compose up -d postgres_keycloak keycloak postgres_app mongo mongo_express

# Stop infrastructure
env-down:
	docker compose down

# Generate fresh Swagger docs
swag:
	cd backend && swag init -g cmd/server/main.go -o docs

# Run Go server locally (clean and safe)
run: init-env
	cd backend && go run $(MAIN)

# Run full development workflow
dev: swag run