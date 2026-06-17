# ArtoFino Backend Core 🎨🚀

Welcome to the backend repository for **ArtoFino** — a high-performance platform designed for digital art tokenization, fractional ownership management, and secure physical art logistics. This service acts as the orchestrator for the entire ArtoFino ecosystem, handling state transitions, secure multi-database persistence, and identity context verification.

---

## 🏗️ Architecture & Storage Model

The system utilizes a dual-storage strategy to balance strict transactional integrity with metadata flexibility:

1. **PostgreSQL (via GORM):** Acts as the single source of truth for all operational ledger states, physical asset logistics, user profiles, and multi-stage creator applications. Enforces strict relational foreign key constraints and atomic ACID transactions.
2. **MongoDB:** Stores complex, deeply nested, and semi-structured art asset metadata, dynamic price evaluation configurations, and history timelines.

### Data Flow Overview

```text
       [ Client Request (Frontend / Mobile) ]
                         │
                         ▼
             [ Nginx / Reverse Proxy ]
                         │
                         ▼
        [ Keycloak Authentication Middleware ]
         (Checks JWT & populates User Context)
                         │
                         ▼
               [ Gin Engine Router ]
          (Matches URL to HTTP Handlers)
                         │
        ┌────────────────┴────────────────┐
        ▼                                 ▼
 [ PostgreSQL (via GORM) ]       [ MongoDB Go Driver ]
  - User Profiles Metadata        - Art Object Dynamic Specs
  - Creator Application States    - Complex Media Asset Arrays
  - Relational Escrow Ledgers     - Performance Pricing Matrices

## 🌟 Core Features & Workflow Chains

- **Onboarding Chain:** A structured application pipeline enabling digital and physical artists to submit applications directly to database ledgers for administrative review.
- **Fractional Shares Engine:** Executes ledger operations tracking granular ownership tokens for physical artwork, guaranteeing exact fractional balances.
- **Logistics Security Layer:** A state-machine tracking system mapping out the transport of high-value physical goods to and from certified art vaults.
- **Interactive Spec Documentation:** Auto-refreshing OpenAPI / Swagger compliance mapping linked directly to handler-level annotations.

---

## 🛠️ Tech Stack & Project Dependencies

- **Core Engine:** Go (Golang) v1.22+
- **Routing Framework:** Gin Gonic (Highly optimized HTTP router multiplexer)
- **Object-Relational Wrapper:** GORM (PostgreSQL Dialect)
- **NoSQL Integration:** MongoDB Go Official Driver
- **Security Framework:** Keycloak Identity Provider (JWT Token parsing and validation)
- **Testing Toolkit:** Testify (Structured Assertions and isolated interface mocking)

---

## 🔧 Environment Configuration (.env)

Before spinning up the service, populate a `.env` file within the backend root folder with the following variables:

# Server Run Profile
SERVER_PORT=9000

GIN_MODE=debug

# Relational Infrastructure
POSTGRES_HOST=localhost

POSTGRES_PORT=5432

POSTGRES_USER=artofino_admin

POSTGRES_PASSWORD=secure_postgres_pass

POSTGRES_DB=artofino_ledger

POSTGRES_SSLMODE=disable

# Document Infrastructure
MONGO_URI=mongodb://localhost:27017

MONGO_DB_NAME=artofino_metadata

# Identity Management (Keycloak Integration)
KEYCLOAK_REALM=artofino-platform

KEYCLOAK_AUTH_SERVER_URL=http://localhost:8080/auth

KEYCLOAK_CLIENT_ID=artofino-backend-api


---

## 📦 Local Development Lifecycle

Ensure your local development environment has running instances of PostgreSQL, MongoDB, and a configured Keycloak realm before executing the steps below.

### 1. Repository Setup & Directory Focus
START_CODE_BASH
git clone https://github.com/dpotekhin-prog/ArtoFino.git
cd ArtoFino/backend
END_CODE

### 2. Launch Development Server
Executes automated migrations for PostgreSQL tables via GORM, opens storage connection pools, and activates the live-reloading server:
START_CODE_BASH
make dev
END_CODE

### 3. Run Verification Tests
Triggers unit test suites across handlers, validators, and database interface abstraction layers using localized mock engines:
START_CODE_BASH
make test
END_CODE

### 4. Recompile API Spec Blueprints
If you alter or introduce Gin context routes, execute this recipe to rebuild Swagger components:
START_CODE_BASH
make swagger
END_CODE

---

## 📖 API Documentation Matrix

Active Swagger documentation is rendered natively on local development interfaces via the built-in swagger handler.

🔗 **Interactive Portal Access Point:** [http://localhost:9000/swagger/index.html](http://localhost:9000/swagger/index.html)

### Production REST Endpoint Schema

| Domain Module | URI Route Template | Method | Security Policy | Business Process Target |
| :--- | :--- | :--- | :--- | :--- |
| **Admin** | `/admin/applications/:id/approve` | `POST` | Admin Role | Approves an author onboarding application and updates Postgres |
| **Admin** | `/admin/objects` | `POST` | Admin Role | Generates and persists a new tokenized art item document into MongoDB |
| **Admin** | `/admin/stats` | `GET` | Admin Role | Compiles relational ledger metrics and operational platform statistics |
| **Admin** | `/admin/ping` | `GET` | Admin Role | Inspects admin context routing paths and middleware vitality states |
| **Authors** | `/authors/apply` | `POST` | Verified Token | Submits biographical metadata and portfolio assets for verification |
| **Transactions** | `/transactions/buy` | `POST` | Verified Token | Finalizes a micro-transaction purchasing fractional asset shares |
| **Transfers** | `/transfers/request` | `POST` | Verified Token | Registers a physical delivery request inside the logistics chain |
| **Transfers** | `/transfers/:id/approve` | `POST` | Admin Role | Updates physical asset escrow verification milestones |
| **Users** | `/users/me` | `GET` | Verified Token | Decodes current Keycloak token context and retrieves profile datasets |
| **System** | `/hc` | `GET` | Public Access | Evaluates systemic health status of database pools and connection runtimes |
| **System** | `/lc` | `GET` | Public Access | Emits platform environment runtime liveness signals |
| **System** | `/rc` | `GET` | Public Access | Emits readiness signals confirming database initialization states |

---

## 🔒 Security & Middleware Layers

1. **Identity Context Hydration:** Intercepts bearer signatures, decodes standard JWT properties against public Keycloak JWKS endpoints, and establishes a context payload containing the validated UUID.
2. **Strict CORS Guard:** Limits resource consumption strictly to configured frontend endpoints while defending against unsafe multi-origin request attempts.
3. **Structured Output Interceptor:** Intercepts standard panic sequences, outputs structured logging parameters to console environments, and ensures zero runtime leaks of inner database errors to clients.