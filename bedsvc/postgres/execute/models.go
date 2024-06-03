// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package execute

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Bed struct {
	ID       uuid.UUID
	IDRoom   uuid.UUID
	Name     string
	Status   string
	CreateAt time.Time
	UpdateAt time.Time
}


type Patient struct {
	ID          uuid.UUID
	PatientCode sql.NullString
	Name        sql.NullString
	Phone       sql.NullString
	Address     sql.NullString
	CreateAt    sql.NullTime
	UpdateAt    sql.NullTime
}

type Record struct {
	ID        uuid.UUID
	IDPatient uuid.NullUUID
	IDDoctor  uuid.NullUUID
	Status    sql.NullString
	CreateAt  sql.NullTime
	UpdateAt  sql.NullTime
}

type Room struct {
	ID            uuid.UUID
	Name          sql.NullString
	BedNumber     sql.NullInt16
	PatientNumber sql.NullInt16
	CreateAt      sql.NullTime
	UpdateAt      sql.NullTime
}

type UseBed struct {
	ID       uuid.UUID
	IDBed    uuid.UUID
	IDRecord uuid.UUID
	Status   string
	CreateAt time.Time
	EndAt    time.Time
}

type UseDevice struct {
	ID       uuid.UUID
	IDDevice uuid.NullUUID
	IDRecord uuid.NullUUID
	Status   sql.NullString
	CreateAt sql.NullTime
	EndAt    sql.NullTime
}

type UseRoom struct {
	ID       uuid.UUID
	IDDoctor sql.NullString
	IDRoom   uuid.NullUUID
}
