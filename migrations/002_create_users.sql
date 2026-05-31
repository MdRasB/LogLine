-- +goose Up
CREATE TABLE users (
    id              UUID            PRIMARY KEY,
    email           TEXT            NOT NULL UNIQUE CHECK (length(email) <= 254),
    password_hash   TEXT            NOT NULL CHECK (length(password_hash) >= 4 AND length(password_hash) <= 60),
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS users;