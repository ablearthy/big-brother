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