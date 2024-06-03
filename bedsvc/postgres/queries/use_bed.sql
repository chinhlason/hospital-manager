-- name: CreateUsageBed :exec
INSERT INTO use_bed(id, id_bed, id_record, status, create_at, end_at)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateUsageBedStatus :exec
UPDATE use_bed
SET status = $1
WHERE id = $2;

-- name: GetUsageBedInUse :one
SELECT * FROM use_bed
WHERE id_bed = $1
AND status = 'IN_USE';

-- name: GetUsageBedByBedId :many
SELECT * FROM use_bed
WHERE id_bed = $1;