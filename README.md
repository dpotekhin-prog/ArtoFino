# ArtoFino Backend

Core REST API service for the ArtoFino project, built with Go.

## Tech Stack
- Language: Go 1.22
- Web Framework: Gin-Gonic
- Databases: PostgreSQL, MongoDB
- Identity & Auth: Keycloak 25.0
- API Docs: Swagger (swaggo)
- Task Runner: Makefile

---

## Quick Start

### 1. Prerequisites
Ensure you have the following installed:
- Docker & Docker Compose
- Go 1.22+
- swag CLI (install via: go install github.com/swaggo/swag/cmd/swag@latest)

### 2. Start Infrastructure
Run the database containers and Keycloak:

make env-up

### 3. Launch Backend
Run the development environment. This command will automatically initialize your .env file from .env.local if it is missing, generate fresh Swagger docs, and boot up the server:

make dev

Once running, access the API documentation at:
- Swagger UI: http://localhost:9000/swagger/index.html

---

## Makefile Commands
- make env-up  - Spin up all databases and Keycloak in Docker.
- make env-down - Stop and remove infrastructure containers.
- make swag     - Re-generate Swagger documentation.
- make run      - Run the Go server locally.
- make dev      - Run the full cycle (swag + run).