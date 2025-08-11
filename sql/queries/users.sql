-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
    VALUES ( gen_random_uuid (), NOW(), NOW(),  $1, $2 )
    RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UpdateUserByID :one
UPDATE users
    SET updated_at = NOW(), email = $2, hashed_password = $3
    WHERE id = $1
    RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;
