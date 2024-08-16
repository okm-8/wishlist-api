CREATE TABLE IF NOT EXISTS migrations (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    executed_at TIMESTAMP DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS migrations_name_idx ON migrations (filename);
CREATE INDEX IF NOT EXISTS migrations_executed_at_idx ON migrations (executed_at);
