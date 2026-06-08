-- +goose Up

CREATE TABLE users (
    id              UUID            PRIMARY KEY,
    email           TEXT            NOT NULL UNIQUE CHECK (length(email) <= 254),
    password_hash   TEXT            NOT NULL,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS users;