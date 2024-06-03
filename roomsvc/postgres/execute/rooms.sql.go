// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: rooms.sql

package execute

import (
	"context"
	"errors"
	"time"
	"github.com/google/uuid"
)

const createRoom = `-- name: CreateRoom :one
INSERT INTO rooms (id, name, bed_number, patient_number, create_at, update_at)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, name, bed_number, patient_number, create_at, update_at
`

type CreateRoomParams struct {
	ID            uuid.UUID
	Name          string
	BedNumber     int64
	PatientNumber int64
	CreateAt      time.Time
	UpdateAt      time.Time
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) (Room, error) {
	row := q.db.QueryRowContext(ctx, createRoom,
		arg.ID,
		arg.Name,
		arg.BedNumber,
		arg.PatientNumber,
		arg.CreateAt,
		arg.UpdateAt,
	)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.BedNumber,
		&i.PatientNumber,
		&i.CreateAt,
		&i.UpdateAt,
	)
	return i, err
}

const existByName = `-- name: ExistByName :one
SELECT EXISTS (
    SELECT 1
    FROM rooms
    WHERE name = $1
) AS exists
`

func (q *Queries) ExistByName(ctx context.Context, name string) (bool, error) {
	row := q.db.QueryRowContext(ctx, existByName, name)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const getAllByAdmin = `-- name: GetAllByAdmin :many
SELECT id, name, bed_number, patient_number, create_at, update_at FROM rooms
`

func (q *Queries) GetAllByAdmin(ctx context.Context) ([]Room, error) {
	rows, err := q.db.QueryContext(ctx, getAllByAdmin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Room
	for rows.Next() {
		var i Room
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.BedNumber,
			&i.PatientNumber,
			&i.CreateAt,
			&i.UpdateAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllRoomByUser = `-- name: GetAllRoomByUser :many
SELECT rooms.id, name, bed_number, patient_number, create_at, update_at, use_room.id, id_doctor, id_room
FROM rooms
JOIN use_room ON rooms.id = use_room.id_room
WHERE id_doctor = $1
`

type GetAllRoomByUserRow struct {
	ID            uuid.UUID
	Name          string
	BedNumber     int64
	PatientNumber int64
	CreateAt      time.Time
	UpdateAt      time.Time
	ID_2          uuid.UUID
	IDDoctor      string
	IDRoom        string
}

func (q *Queries) GetAllRoomByUser(ctx context.Context, idDoctor string) ([]GetAllRoomByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllRoomByUser, idDoctor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllRoomByUserRow
	for rows.Next() {
		var i GetAllRoomByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.BedNumber,
			&i.PatientNumber,
			&i.CreateAt,
			&i.UpdateAt,
			&i.ID_2,
			&i.IDDoctor,
			&i.IDRoom,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoom = `-- name: GetRoom :one
SELECT id, name, bed_number, patient_number, create_at, update_at FROM rooms WHERE id = $1
`

func (q *Queries) GetRoom(ctx context.Context, id string) (Room, error) {
	row := q.db.QueryRowContext(ctx, getRoom, id)
	var i Room
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.BedNumber,
		&i.PatientNumber,
		&i.CreateAt,
		&i.UpdateAt,
	)
	return i, err
}

const updateBedNumber = `-- name: UpdateBedNumber :exec
UPDATE rooms
SET bed_number = bed_number + 1
WHERE id = $1
`

func (q *Queries) UpdateBedNumber(ctx context.Context, id string) error {
	check, err := q.db.ExecContext(ctx, updateBedNumber, id)
	checkId, _ := check.RowsAffected()
	if checkId == 0 {
		return errors.New("no id data found")
	}
	return err
}

const updateInformation = `-- name: UpdateInformation :exec
UPDATE rooms
SET name = $1
WHERE id = $2
`

type UpdateInformationParams struct {
	Name string
	ID   string
}

func (q *Queries) UpdateInformation(ctx context.Context, arg UpdateInformationParams) error {
	check, err := q.db.ExecContext(ctx, updateInformation, arg.Name, arg.ID)
	checkId, _ := check.RowsAffected()
	if checkId == 0 {
		return errors.New("no id data found")
	}
	return err
}

const updatePatientNumber = `-- name: UpdatePatientNumber :exec
UPDATE rooms
SET patient_number = patient_number + 1
WHERE id = $1
`

func (q *Queries) UpdatePatientNumber(ctx context.Context, id string) error {
	check, err := q.db.ExecContext(ctx, updatePatientNumber, id)
	checkId, _ := check.RowsAffected()
	if checkId == 0 {
		return errors.New("no id data found")
	}
	return err
}
