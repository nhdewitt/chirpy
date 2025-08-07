-- name: PostChirp :one
INSERT INTO CHIRPS(id, user_id, created_at, updated_at, body)
VALUES (
    gen_random_uuid(),
    $1,
    NOW(),
    NOW(),
    $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM CHIRPS
ORDER BY created_at ASC;

-- name: GetOneChirp :one
SELECT * FROM CHIRPS
WHERE id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;