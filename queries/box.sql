-- name: CreateBox :one
INSERT INTO box (
    fingerprint_id,
    container_id,
    status,
    last_active
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetBoxByFingerprint :one
SELECT * FROM box
WHERE fingerprint_id = $1 LIMIT 1;

-- name: GetBoxByContainerID :one
SELECT * FROM box
WHERE container_id = $1 LIMIT 1;

-- name: UpdateBoxStatus :exec
UPDATE box
SET status = $2
WHERE fingerprint_id = $1;

-- name: TouchBox :exec
-- Updates last_active and ensures status is 'active'
UPDATE box
SET last_active = $2,
    status = 'active'
WHERE fingerprint_id = $1;

-- name: DeleteBox :exec
DELETE FROM box
WHERE fingerprint_id = $1;

-- name: GetExpiredBoxes :many
-- Used by the 24h cleanup worker
SELECT * FROM box
WHERE last_active < $1;

-- name: ListBoxesByStatus :many
SELECT * FROM box
WHERE status = $1;
