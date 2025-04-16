-- +goose Up
CREATE TABLE auth (
    id SERIAL PRIMARY KEY,
    refresh_token VARCHAR(255) UNIQUE NOT NULL,
    access_token_uid VARCHAR(255) UNIQUE NOT NULL,
    is_used INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE auth;