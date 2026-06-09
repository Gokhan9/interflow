-- name: GetUserByAPIKey :one
SELECT u.*, ak.id as api_key_id FROM users u
JOIN api_keys ak ON u.id = ak.user_id
WHERE ak.key_value = $1 AND ak.is_active = true;

-- name: CreateUser :one
INSERT INTO users (username) 
VALUES ($1) 
RETURNING *;

-- name: CreateAPIKey :one
INSERT INTO api_keys (user_id, key_value) 
VALUES ($1, $2) 
RETURNING *;

-- name: CreateUsageLog :exec
INSERT INTO usage_logs (
    api_key_id, 
    provider, 
    model, 
    prompt_tokens,
    completion_tokens, 
    total_tokens, 
    latency_ms, 
    status_code
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
);
