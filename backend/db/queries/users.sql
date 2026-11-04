-- name: CreateUser :one
INSERT INTO users (name, friendcode)
VALUES (?, ?)
RETURNING *;
