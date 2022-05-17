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

-- name: GetAllUserTokens :many
SELECT user_tokens.user_id,
       user_tokens.access_token,
       vk_tokens.vk_user_id
FROM user_tokens
LEFT JOIN vk_tokens ON user_tokens.access_token = vk_tokens.access_token;

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

-- name: GetVkUserIdByUserId :one
SELECT vk_tokens.vk_user_id
FROM user_tokens
LEFT JOIN vk_tokens ON user_tokens.access_token = vk_tokens.access_token
WHERE user_tokens.user_id = $1;

-- name: SaveVkMessage :one
INSERT INTO vk_messages (
    vk_owner_id, message_id, message
) VALUES (
    $1, $2, $3
)
RETURNING id;

-- name: SaveMessageEvent :exec
INSERT INTO vk_message_events (
    internal_message_id, m_type, created_at
) VALUES (
    $1, $2, $3
);

-- name: GetLastSavedVKMessage :one
SELECT max(id)
FROM vk_messages
WHERE vk_owner_id = $1 AND message_id = $2;

-- name: GetMessageById :one
SELECT * FROM vk_messages WHERE id = $1;

-- name: SaveActivityEvent :exec
INSERT INTO vk_activity_events (
    vk_owner_id, target_id, activity, platform, kicked_by_timeout, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetLastMessages :many
SELECT vk_message_events.id, m_type, msgs.message
FROM vk_message_events
LEFT JOIN vk_messages msgs ON internal_message_id = msgs.id
WHERE vk_owner_id = $1 AND vk_message_events.id < $2
ORDER BY vk_message_events.id DESC
LIMIT $3;