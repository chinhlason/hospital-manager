-- name: CreateRoom :one
INSERT INTO beds (id, id_room, name, status, create_at, update_at)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;

-- name: ExistByName :one
SELECT EXISTS (
    SELECT 1
    FROM beds
    WHERE name = $1 AND id_room = $2
) AS exists;

-- name: GetBedByStatus :many
SELECT * FROM beds
WHERE status = $1 AND id_room = $2;

-- name: UpdateBedName :exec
UPDATE beds
SET name = $1
WHERE id = $2;

-- name: UpdateBedStatus :exec
UPDATE beds
SET status = $1
WHERE id = $2;

-- name: GetAllByAdmin :many
SELECT * FROM beds;
