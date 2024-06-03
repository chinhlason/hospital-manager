package service

import (
	"bedsvc/constant"
	"bedsvc/postgres/execute"
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
)

type BedServices interface {
	CreateBed(ctx context.Context, roomName, bedName, token string) error
	GetBedInRoom(ctx context.Context, idRoom, token string) ([]Room, error)
	//GetBedByStatus(ctx context.Context, status, roomName string) ([]execute.Bed, error)
	//UpdateBedName(ctx context.Context, id, name string) error
	//ChangeBedStatus(ctx context.Context, id, status string) error
	//HandoverBed(ctx context.Context, idRecord, idBed string) error
}

type bedservice struct {
	Queries     *execute.Queries
	RoomService RoomService
}

func NewBedService(q *execute.Queries, r RoomService) BedServices {
	return &bedservice{
		Queries:     q,
		RoomService: r,
	}
}

func (b *bedservice) checkRoomPermissionById(ctx context.Context, idRoom, token string) bool {
	rooms, err := b.RoomService.GetRoomByCurr(ctx, token)
	var checkPermissionInRoom = true
	if err != nil {
		return false
	}
	for _, room := range rooms {
		if idRoom == room.ID {
			checkPermissionInRoom = true
			idRoom = room.ID
			break
		}
		checkPermissionInRoom = false
	}
	return checkPermissionInRoom
}

func (b *bedservice) CreateBed(ctx context.Context, roomName, bedName, token string) error {
	rooms, err := b.RoomService.GetRoomByCurr(ctx, token)
	var checkPermissionInRoom = true
	var idRoom string
	if err != nil {
		return err
	}
	for _, room := range rooms {
		if roomName == room.Name {
			checkPermissionInRoom = true
			idRoom = room.ID
			break
		}
		checkPermissionInRoom = false
	}
	if !checkPermissionInRoom {
		return errors.New("do not have permission in this room")
	}
	checkExist, err := b.Queries.ExistByName(ctx, execute.ExistByNameParams{
		Name:   bedName,
		IDRoom: idRoom,
	})
	if checkExist {
		return errors.New("bed name already exist")
	}
	id := uuid.New()
	_, err = b.Queries.CreateBed(ctx, execute.CreateBedParams{
		ID:       id,
		IDRoom:   idRoom,
		Name:     bedName,
		Status:   constant.AVAILABLE_BED_STATUS,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *bedservice) GetBedInRoom(ctx context.Context, idRoom, token string) ([]Room, error) {
	//check := b.checkRoomPermissionById(ctx, idRoom, token)
	//if check {
	//	beds, err := b.Queries.GetBedByStatus()
	//}
	//if err != nil {
	//	return nil, err
	//}
	//return rooms, nil
	return nil, nil
}

type CreateBedReq struct {
	BedName  string
	RoomName string
	Token    string
}

type GetInRoomReq struct {
	IdRoom string
	Token  string
}
