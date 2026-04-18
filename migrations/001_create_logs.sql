-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE logs (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    level       VARCHAR(10)  NOT NULL,
    message     TEXT         NOT NULL,
    service     VARCHAR(100) NOT NULL,
    timestamp   TIMESTAMPTZ  NOT NULL,
    metadata    JSONB,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_logs_service   ON logs(service);
CREATE INDEX idx_logs_level     ON logs(level);
CREATE INDEX idx_logs_timestamp ON logs(timestamp DESC);

-- +goose Down
DROP TABLE IF EXISTS logs;
