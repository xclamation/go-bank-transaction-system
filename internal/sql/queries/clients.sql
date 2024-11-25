-- name: CreateClient :one
INSERT INTO clients (name, balance)
VALUES ($1, $2)
RETURNING id, name, balance;

-- name: GetClient :one
SELECT id, name, balance FROM clients
WHERE id = $1;

-- name: UpdateClientBalance :exec
UPDATE clients
SET balance = $2
WHERE id = $1;