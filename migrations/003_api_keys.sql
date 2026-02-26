CREATE TABLE IF NOT EXISTS api_keys (
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    environment  TEXT NOT NULL CHECK(environment IN ('live', 'test', 'staging')),
    prefix       TEXT NOT NULL,
    hashed_key   TEXT NOT NULL UNIQUE,
    revoked      BOOLEAN NOT NULL DEFAULT 0,
    created_at   DATETIME NOT NULL DEFAULT (datetime('now')),
    last_used_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_api_keys_hashed ON api_keys(hashed_key) WHERE revoked = 0;
