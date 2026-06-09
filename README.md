# ArtoFino Backend

Core REST API service for the ArtoFino project, built with Go. This service manages dynamic asset pricing, fractional ownership transactions, temporary art hosting logistics, and two-stage creator onboarding.

## Tech Stack
- **Language:** Go 1.22
- **Web Framework:** Gin-Gonic
- **Databases:** PostgreSQL (Transactional State), MongoDB (Dynamic Art Metadata)
- **Identity & Auth:** Keycloak 25.0
- **API Docs:** Swagger (swaggo)
- **Task Runner:** Makefile

---

## Core API Modules

* **`objects` (Public):** Art object directory displaying live, dynamically calculated prices based on asset lifecycle, linear growth formulas, and social activity coefficients.
* **`transactions` (Protected):** Secure fractional ownership purchases (ranging from 1% to 10% shares) with atomic price-locking mechanism.
* **`transfers` (Protected):** Logistics management enabling partners to request and approve physical relocation of art objects for pop-up events and temporary exhibitions.
* **`authors` (Protected):** Two-stage creator application system for verifying artist portfolios and professional social media handles.
* **`admin` & `system`:** Platform monitoring, database health checks, and global system metrics.

---

## Quick Start

### 1. Prerequisites
Ensure you have the following installed:
- Docker & Docker Compose
- Go 1.22+
- swag CLI

### 2. Start Infrastructure
Spin up the development databases (PostgreSQL, MongoDB), database management dashboards, and the Keycloak IAM container:
> make env-up

### 3. Launch Backend Workflow
Run the local development pipeline. This command automatically sets up your local configuration variables, builds fresh Swagger API endpoints, and starts the live server:
> make dev

Once the service initializes successfully, access the open documentation portal at:
- Swagger UI: http://localhost:9000/swagger/index.html

---

## Technical Operations & Verification

Our test suite handles automated verification of the mathematical growth engines, API payload boundaries, and route protection structures.

- Run the complete integration and unit test suite:
  make test

- Target a highly specific test block inside the project tree:
  make test-target name=TestGetArtObject_PriceCalculation

---

## Makefile Reference

- make env-up         -> Spin up all core databases, management GUIs, and Keycloak in Docker.
- make env-down       -> Gracefully stop and tear down infrastructure containers.
- make swag           -> Re-generate Swagger documentation using unified camelCase strategy.
- make run            -> Boot up the Go HTTP engine locally using active configuration profiles.
- make dev            -> Standard workflow: runs swag followed immediately by run.
- make test           -> Discover and run all automated unit and integration tests across target packages.
- make test-target    -> Isolate and run a single specific test by its descriptor flag name.