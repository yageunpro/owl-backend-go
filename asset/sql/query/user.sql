-- name: CreateUser :exec
INSERT INTO auth."user"(id, email, username)
VALUES ($1, $2, $3);

-- name: CreatePassword :exec
INSERT INTO auth.password (id, password_hash)
VALUES ($1, $2);

-- name: CreateOAuth :exec
INSERT INTO auth.oauth (id, open_id, access_token, refresh_token, allow_sync, valid_until)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetUser :one
SELECT id, email, username
FROM auth."user"
WHERE id = $1;

-- name: GetUserPassword :one
SELECT id, password_hash
FROM auth.password
WHERE id = $1;

-- name: GetAllOAuthUserIds :many
SELECT id
FROM auth.oauth
WHERE deleted_at IS NULL;

-- name: GetUserOAuthToken :one
SELECT id, open_id, access_token, refresh_token, valid_until
FROM auth.oauth
WHERE id = $1;

-- name: FindUser :one
SELECT u.id AS id, u.email AS email, p.password_hash AS password_hash
FROM auth."user" AS u
         INNER JOIN auth.password AS p ON u.id = p.id
WHERE email = $1;

-- name: FindOAuth :one
SELECT id, open_id
FROM auth.oauth
WHERE open_id = $1;

-- name: UpdateOAuthToken :exec
UPDATE auth.oauth
SET access_token  = COALESCE(sqlc.narg(u_access), access_token),
    refresh_token = COALESCE(sqlc.narg(u_refresh), refresh_token),
    valid_until   = COALESCE(sqlc.narg(u_valid_until), valid_until),
    allow_sync    = COALESCE(sqlc.narg(u_allow_sync), allow_sync),
    updated_at    = NOW()
WHERE id = $1;
