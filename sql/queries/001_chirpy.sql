-- name: CreateChirpy :one
INSERT INTO chirpies (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid (),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllChirpies :exec
DELETE FROM chirpies;