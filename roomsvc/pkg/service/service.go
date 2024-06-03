package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"roomsvc/postgres/execute"
	"time"
)

type RoomService interface {
	CreateRoom(ctx context.Context, name string) error
	GetRoomById(ctx context.Context, id string) (execute.Room, error)
	GetAllRoomByUser(ctx context.Context, idDoctor string) (GetAllRoomsRes, error)
	HandoverRoom(ctx context.Context, idRoom, idDoctor string) error
	UpdateNumber(ctx context.Context, id, option string) error
	UpdateInformation(ctx context.Context, id, name string) error
	UpdateUseRoom(ctx context.Context, idRoom, idDoctor string) error
	ExistName(ctx context.Context, name string) bool
	GetAllByAdmin(ctx context.Context) ([]execute.Room, error)
	GetRoomByCurrent(ctx context.Context, token string) ([]execute.Room, error)
}

type roomService struct {
	Queries     *execute.Queries
	UserService UserService
}

func NewRoomService(queries *execute.Queries, usersvc UserService) RoomService {
	return &roomService{
		Queries:     queries,
		UserService: usersvc,
	}
}

func (r *roomService) CreateRoom(ctx context.Context, name string) error {
	id := uuid.New()
	check, _ := r.Queries.ExistByName(ctx, name)
	if check {
		return errors.New("duplicate room's name")
	}

	_, err := r.Queries.CreateRoom(ctx, execute.CreateRoomParams{
		ID:            id,
		Name:          name,
		BedNumber:     0,
		PatientNumber: 0,
		CreateAt:      time.Now(),
		UpdateAt:      time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *roomService) GetRoomById(ctx context.Context, id string) (execute.Room, error) {
	room, err := r.Queries.GetRoom(ctx, id)
	if err != nil {
		return execute.Room{}, err
	}
	return room, nil
}

func (r *roomService) HandoverRoom(ctx context.Context, idRoom, idDoctor string) error {
	id := uuid.New()
	err := r.Queries.CreateUsageRoom(ctx, execute.CreateUsageRoomParams{
		ID:       id,
		IDRoom:   idRoom,
		IDDoctor: idDoctor,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *roomService) GetAllRoomByUser(ctx context.Context, idDoctor string) (GetAllRoomsRes, error) {
	var result GetAllRoomsRes
	var temp []execute.Room
	doctor, err := r.UserService.GetUserInformation(ctx, idDoctor, "id")
	if err != nil {
		return GetAllRoomsRes{}, err
	}
	rooms, err := r.Queries.GetAllRoomByUser(ctx, idDoctor)
	if err != nil {
		return GetAllRoomsRes{}, err
	}
	for _, room := range rooms {
		temp2 := execute.Room{
			ID:            room.ID,
			Name:          room.Name,
			BedNumber:     room.BedNumber,
			PatientNumber: room.PatientNumber,
			CreateAt:      room.CreateAt,
			UpdateAt:      room.UpdateAt,
		}
		temp = append(temp, temp2)
	}
	result.User = doctor
	result.Room = temp
	return result, nil
}

func (r *roomService) UpdateNumber(ctx context.Context, id, option string) error {
	if option == "patient" {
		err := r.Queries.UpdatePatientNumber(ctx, id)
		if err != nil {
			return err
		}
		return nil
	}
	err := r.Queries.UpdateBedNumber(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *roomService) UpdateInformation(ctx context.Context, id, name string) error {
	check, _ := r.Queries.ExistByName(ctx, name)
	if check {
		return errors.New("duplicate room's name")
	}

	err := r.Queries.UpdateInformation(ctx, execute.UpdateInformationParams{
		ID:   id,
		Name: name,
	})

	if err != nil {
		return err
	}
	return nil
}

func (r *roomService) UpdateUseRoom(ctx context.Context, idRoom, idDoctor string) error {
	_, err := r.UserService.GetUserInformation(ctx, idDoctor, "id")
	if err != nil {
		return err
	}
	err = r.Queries.UpdateUsageRoom(ctx, execute.UpdateUsageRoomParams{
		IDDoctor: idDoctor,
		IDRoom:   idRoom,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *roomService) ExistName(ctx context.Context, name string) bool {
	check, _ := r.Queries.ExistByName(ctx, name)
	return check
}

func (r *roomService) GetAllByAdmin(ctx context.Context) ([]execute.Room, error) {
	rooms, err := r.Queries.GetAllByAdmin(ctx)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomService) GetRoomByCurrent(ctx context.Context, token string) ([]execute.Room, error) {
	userId, err := r.UserService.ValidateUser(ctx, token)
	var result []execute.Room
	if err != nil {
		return nil, err
	}
	rooms, err := r.Queries.GetAllRoomByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	for _, room := range rooms {
		temp := execute.Room{
			ID:            room.ID,
			Name:          room.Name,
			BedNumber:     room.BedNumber,
			PatientNumber: room.PatientNumber,
			CreateAt:      room.CreateAt,
			UpdateAt:      room.UpdateAt,
		}
		result = append(result, temp)
	}
	return result, nil
}
