-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
    VALUES ( gen_random_uuid (), NOW(), NOW(),  $1, $2 )
    RETURNING *;

-- name: GetChirpByID :one
SELECT * FROM chirps
WHERE id = $1  LIMIT 1;

-- name: GetChirpsByUserID :many
SELECT * FROM chirps
WHERE user_id = $1 ORDER BY created_at ASC LIMIT $2;

-- name: DeleteAllChirps :exec
DELETE FROM chirps;

-- name: DeleteOwnChirpByID :many
DELETE FROM chirps WHERE id = $1 AND user_id = $2 RETURNING *;
   -- LIMIT 1 

-- name: GetAllChirps :many
SELECT * FROM chirps ORDER BY created_at ASC;
