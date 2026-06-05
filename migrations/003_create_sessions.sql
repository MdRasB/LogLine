-- +goose Up

CREATE TABLE sessions (
    id          UUID   PRIMARY KEY,
    user_id     UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT        NOT NULL UNIQUE CHECK(length(token_hash) = 64),
    expires_at  TIMESTAMPTZ NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_user_id   ON  sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON  sessions(expires_at);

-- +goose Down
DROP TABLE IF EXISTS sessions; 