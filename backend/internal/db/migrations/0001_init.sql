-- Users table
CREATE TABLE IF NOT EXISTS users (
                                     id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email           TEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL,
    keycloak_id     TEXT NOT NULL UNIQUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Authors validation requests
CREATE TABLE IF NOT EXISTS author_validation_requests (
                                                          id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status          TEXT NOT NULL, -- e.g. 'pending', 'approved', 'rejected'
    comment         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Ownerships (user owns some object stored in Mongo)
CREATE TABLE IF NOT EXISTS ownerships (
                                          id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    object_id       TEXT NOT NULL, -- Mongo ObjectID as string
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Transactions (money flow between users)
CREATE TABLE IF NOT EXISTS transactions (
                                            id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    from_user_id    UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    to_user_id      UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    amount_cents    BIGINT NOT NULL,
    currency        TEXT NOT NULL DEFAULT 'EUR',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );

-- Indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_author_validation_requests_user_id ON author_validation_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_ownerships_user_id ON ownerships(user_id);
CREATE INDEX IF NOT EXISTS idx_ownerships_object_id ON ownerships(object_id);
CREATE INDEX IF NOT EXISTS idx_transactions_from_user_id ON transactions(from_user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_to_user_id ON transactions(to_user_id);
