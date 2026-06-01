-- name: GetUserByAPIKey :one
SELECT u.* FROM users u
JOIN api_keys ak ON u.id = ak.user_id
WHERE ak.key_value = $1 AND ak.is_active = true;

-- name: CreateUser :one
INSERT INTO users (username) VALUES ($1) RETURNING *;

-- name: CreateAPIKey :one
INSERT INTO api_keys (user_id, key_value) VALUES ($1, $2) RETURNING *;