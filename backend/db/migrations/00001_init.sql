-- +goose Up
CREATE TABLE users (
    id         INTEGER PRIMARY KEY,
    name       TEXT    NOT NULL,
    friendcode TEXT    NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
