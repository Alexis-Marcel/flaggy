CREATE TABLE IF NOT EXISTS flags (
    key          TEXT PRIMARY KEY,
    type         TEXT NOT NULL CHECK(type IN ('boolean', 'string', 'number', 'json')),
    description  TEXT NOT NULL DEFAULT '',
    enabled      BOOLEAN NOT NULL DEFAULT 0,
    default_value TEXT NOT NULL,
    created_at   DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at   DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS rules (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    flag_key    TEXT NOT NULL REFERENCES flags(key) ON DELETE CASCADE,
    description TEXT NOT NULL DEFAULT '',
    value       TEXT NOT NULL,
    priority    INTEGER NOT NULL DEFAULT 0,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_rules_flag_key ON rules(flag_key);
CREATE INDEX IF NOT EXISTS idx_rules_priority ON rules(flag_key, priority);

CREATE TABLE IF NOT EXISTS conditions (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    rule_id    INTEGER NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    attribute  TEXT NOT NULL,
    operator   TEXT NOT NULL,
    value      TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_conditions_rule_id ON conditions(rule_id);
