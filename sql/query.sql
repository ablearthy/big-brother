-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
    username, password, inviter_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetUserIdByInviteCode :one
SELECT user_id FROM invite_codes
WHERE invite_code = $1 LIMIT 1;

-- name: CreateInviteCode :one
INSERT INTO invite_codes (
    user_id, invite_code
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetCountOfUsedInviteCodes :one
SELECT count(*) FROM users
WHERE inviter_id = $1;

-- name: GetUserByUsername :one
SELECT id, password FROM users
WHERE username = $1;

-- name: GetTokenById :one
SELECT access_token FROM user_tokens
WHERE user_id = $1;

-- name: DeleteTokenById :one
DELETE FROM user_tokens WHERE user_id = $1
RETURNING *;

-- name: CreateUserToken :one
INSERT INTO user_tokens (
    user_id, access_token
) VALUES (
    $1, $2
)
RETURNING *;

-- name: CreateVkToken :exec
INSERT INTO vk_tokens (
    access_token, vk_user_id
) VALUES (
    $1, $2
)
ON CONFLICT DO NOTHING;

-- name: GetVkToken :one
SELECT access_token, vk_user_id
FROM vk_tokens
WHERE access_token = $1;
