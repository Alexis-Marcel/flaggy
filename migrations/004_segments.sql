CREATE TABLE IF NOT EXISTS segments (
    key         TEXT PRIMARY KEY,
    description TEXT NOT NULL DEFAULT '',
    created_at  DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS segment_conditions (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    segment_key TEXT NOT NULL REFERENCES segments(key) ON DELETE CASCADE,
    attribute   TEXT NOT NULL,
    operator    TEXT NOT NULL,
    value       TEXT NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_segment_conditions_key ON segment_conditions(segment_key);

-- Junction table: links rules <-> segments (many-to-many)
CREATE TABLE IF NOT EXISTS rule_segments (
    rule_id     INTEGER NOT NULL REFERENCES rules(id) ON DELETE CASCADE,
    segment_key TEXT NOT NULL REFERENCES segments(key) ON DELETE CASCADE,
    PRIMARY KEY (rule_id, segment_key)
);

CREATE INDEX IF NOT EXISTS idx_rule_segments_rule ON rule_segments(rule_id);
CREATE INDEX IF NOT EXISTS idx_rule_segments_segment ON rule_segments(segment_key);
