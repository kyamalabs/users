-- name: CreateProfile :one
INSERT INTO profiles (
    wallet_address,
    gamer_tag
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetProfile :one
SELECT * FROM profiles
WHERE wallet_address = $1
LIMIT 1;

-- name: GetProfilesCount :one
SELECT COUNT(*) as total_profiles FROM profiles;

-- name: ListProfiles :many
SELECT * FROM profiles
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: UpdateProfile :one
UPDATE profiles
SET
    gamer_tag = COALESCE(sqlc.narg(gamer_tag), gamer_tag)
WHERE
    wallet_address = sqlc.arg(wallet_address)
RETURNING *;

-- name: DeleteProfile :exec
DELETE FROM profiles
WHERE wallet_address = $1;
