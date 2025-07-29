-- +goose Up
CREATE TABLE users (
    id TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    email TEXT
);

-- +goose Down
DROP TABLE users;
