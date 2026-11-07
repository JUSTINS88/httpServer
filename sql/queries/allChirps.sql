-- name: GetAllChirps :many
SELECT * FROM chirps;

-- name: GetAllChirpsByAuthorID :many
SELECT * FROM chirps WHERE user_id = $1
ORDER BY created_at ASC;