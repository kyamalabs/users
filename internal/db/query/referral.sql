-- name: CreateReferral :one
INSERT INTO referrals (
    referrer,
    referee
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetReferrer :one
SELECT * FROM referrals
WHERE referee = $1
LIMIT 1;

-- name: ListReferrals :many
SELECT * FROM referrals
WHERE referrer = $1
ORDER BY referred_at DESC
LIMIT $2
OFFSET $3;
