-- name: CreateRoom :one
INSERT INTO rooms (id, name, bed_number, patient_number, create_at, update_at)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;

-- name: GetRoom :one
SELECT * FROM rooms WHERE id = $1;

-- name: ExistByName :one
SELECT EXISTS (
    SELECT 1
    FROM rooms
    WHERE name = $1
) AS exists;

-- name: UpdateBedNumber :exec
UPDATE rooms
SET bed_number = bed_number + 1
WHERE id = $1;

-- name: UpdatePatientNumber :exec
UPDATE rooms
SET patient_number = patient_number + 1
WHERE id = $1;

-- name: UpdateInformation :exec
UPDATE rooms
SET name = $1
WHERE id = $2;

-- name: GetAllRoomByUser :many
SELECT *
FROM rooms
JOIN use_room ON rooms.id = use_room.id_room
WHERE id_doctor = $1;

-- name: GetAllByAdmin :many
SELECT * FROM rooms;
