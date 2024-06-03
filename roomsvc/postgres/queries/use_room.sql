-- name: CreateUsageRoom :exec
INSERT INTO use_room (id, id_doctor, id_room)
VALUES ($1, $2, $3);

-- name: UpdateUsageRoom :exec
UPDATE use_room
SET id_doctor = $1
WHERE id_room = $2;
