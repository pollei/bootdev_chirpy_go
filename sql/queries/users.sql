-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
    VALUES ( gen_random_uuid (), NOW(), NOW(),  $1 )
    RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;
