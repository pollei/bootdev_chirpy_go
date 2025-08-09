-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
    VALUES ( gen_random_uuid (), NOW(), NOW(),  $1, $2 )
    RETURNING *;

-- name: GetChirpsByUserID :many
SELECT * FROM chirps
WHERE user_id = $1 LIMIT $2;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;

-- name: GetAllChirps :many
SELECT * FROM chirps;
