package service

import (
	"context"
	"time"
)

type RoomService interface {
	GetRoomByCurr(ctx context.Context, token string) (GetRoomsRes, error)
}

type GetRoomsRes []Room

type Room struct {
	ID            string    `json:"ID"`
	Name          string    `json:"Name"`
	BedNumber     int       `json:"BedNumber"`
	PatientNumber int       `json:"PatientNumber"`
	CreateAt      time.Time `json:"CreateAt"`
	UpdateAt      time.Time `json:"UpdateAt"`
}
