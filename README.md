# ArtoFino Backend Core

Backend service for **ArtoFino** — a platform for digital art tokenization, fractional ownership management, and physical art logistics. It acts as the central orchestrator for state transitions, multi-database persistence, and identity context verification.

---

## Architecture & Storage

Dual-storage strategy balancing transactional integrity with metadata flexibility:

- **PostgreSQL (GORM):** Source of truth for ledger states, physical asset logistics, user profiles, creator applications, and share balances. Enforces relational FK constraints and ACID transactions.
- **MongoDB:** Stores semi-structured art asset metadata, dynamic price evaluation configs, and history timelines.

### Data Flow

```
[ Client (Frontend / Mobile) ]
              │
              ▼
  [ Nginx / Reverse Proxy ]
              │
              ▼
 [ Keycloak Auth Middleware ]
   (JWT validation + User Context)
              │
              ▼
      [ Gin Router ]
              │
    ┌─────────┴─────────┐
    ▼                   ▼
[ PostgreSQL ]      [ MongoDB ]
 - User profiles     - Art object specs
 - Creator states    - Media asset arrays
 - Escrow ledgers    - Pricing matrices
```

---

## Features

- **Onboarding Pipeline:** Artists submit applications that flow through an admin review state machine.
- **Fractional Shares Engine:** Atomic `ExecutePurchase` transfers share percentages between users in a single DB transaction — locks the seller's `ArtShareBalance` row with `SELECT FOR UPDATE`, deducts from the seller, upserts to the buyer, and appends a `Transaction` history record.
- **Logistics Tracking:** State-machine for monitoring high-value physical goods in transit to/from art vaults.
- **OpenAPI Docs:** Auto-refreshing Swagger spec linked to handler annotations.

---

## Tech Stack

| Component | Technology |
| :--- | :--- |
| Language | Go 1.22+ |
| HTTP Router | Gin |
| ORM | GORM (PostgreSQL) |
| NoSQL | MongoDB Go Driver |
| Auth | Keycloak (JWT / JWKS) |
| Testing | Testify |

---

## Configuration

Create a `.env` file in the backend root:

```env
# Server
SERVER_PORT=9000
GIN_MODE=debug

# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=artofino_admin
POSTGRES_PASSWORD=secure_postgres_pass
POSTGRES_DB=artofino_ledger
POSTGRES_SSLMODE=disable

# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DB_NAME=artofino_metadata

# Keycloak
KEYCLOAK_REALM=artofino-platform
KEYCLOAK_AUTH_SERVER_URL=http://localhost:8080/auth
KEYCLOAK_CLIENT_ID=artofino-backend-api
```

---

## Local Development

Requires running instances of PostgreSQL, MongoDB, and a configured Keycloak realm.

```bash
# 1. Clone and enter the directory
git clone https://github.com/dpotekhin-prog/ArtoFino.git
cd ArtoFino/backend

# 2. Start dev server (runs GORM migrations, opens connection pools, live reload)
make dev

# 3. Run tests (unit tests with mocks across handlers, validators, repositories)
make test

# 4. Rebuild Swagger spec (run after adding or modifying routes)
make swagger
```

---

## API Reference

Swagger UI: [http://localhost:9000/swagger/index.html](http://localhost:9000/swagger/index.html)

| Module | Route | Method | Auth | Description |
| :--- | :--- | :---: | :--- | :--- |
| Admin | `/admin/applications/:id/approve` | `POST` | Admin | Approve a creator onboarding application |
| Admin | `/admin/objects` | `POST` | Admin | Create and persist a tokenized art object |
| Admin | `/admin/stats` | `GET` | Admin | Platform ledger metrics and statistics |
| Admin | `/admin/ping` | `GET` | Admin | Admin middleware health check |
| Authors | `/authors/apply` | `POST` | Token | Submit portfolio and bio for verification |
| Transactions | `/transactions/buy` | `POST` | Token | Purchase fractional asset shares |
| Transfers | `/transfers/request` | `POST` | Token | Register a physical delivery request (requires active share ownership — 403 if none) |
| Transfers | `/transfers/:id/approve` | `POST` | Admin | Advance asset escrow milestone |
| Users | `/users/me` | `GET` | Token | Get current user profile from token context |
| System | `/hc` | `GET` | Public | Health check — database pool status |
| System | `/lc` | `GET` | Public | Liveness probe |
| System | `/rc` | `GET` | Public | Readiness probe — database initialization |

---

## Security

1. **JWT Middleware:** Validates bearer tokens against Keycloak JWKS endpoints and injects a verified user UUID into request context.
2. **Share Ownership Guard:** `POST /transfers/request` verifies the requester holds an active `ArtShareBalance > 0` for the target object before creating a logistics entry. Returns `403` otherwise.
3. **CORS Policy:** Restricts cross-origin requests to configured frontend origins.
4. **Panic Recovery:** Intercepts panics, logs structured output, and prevents internal errors from leaking to clients.
