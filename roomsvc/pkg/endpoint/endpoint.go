package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"roomsvc/pkg/service"
)

type RoomEndpoints struct {
	CreateRoomEndpoint        endpoint.Endpoint
	CreateListRoomsEndpoint   endpoint.Endpoint
	GetRoomByIdEndpoint       endpoint.Endpoint
	HandoverRoomEndpoint      endpoint.Endpoint
	GetAllRoomByUserEndpoint  endpoint.Endpoint
	UpdateNumberEndpoint      endpoint.Endpoint
	UpdateInformationEndpoint endpoint.Endpoint
	UpdateUseRoomEndpoint     endpoint.Endpoint
	GetAllByAdmin             endpoint.Endpoint
	GetAllByCurrent           endpoint.Endpoint
}

func MakeRoomServerEndpoints(s service.RoomService) RoomEndpoints {
	return RoomEndpoints{
		CreateRoomEndpoint:        makeCreateRoomEndpoint(s),
		CreateListRoomsEndpoint:   makeCreateListRoomsEndpoint(s),
		GetRoomByIdEndpoint:       makeGetRoomById(s),
		HandoverRoomEndpoint:      makeHandoverRoom(s),
		GetAllRoomByUserEndpoint:  makeGetAllRoomByUser(s),
		UpdateNumberEndpoint:      makeUpdateNumber(s),
		UpdateInformationEndpoint: makeUpdateInformation(s),
		UpdateUseRoomEndpoint:     makeUpdateUseRoom(s),
		GetAllByAdmin:             makeGetAllByAdmin(s),
		GetAllByCurrent:           makeGetAllByCurrent(s),
	}
}

func makeCreateRoomEndpoint(s service.RoomService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CreateRoomReq)
		err = s.CreateRoom(ctx, req.Name)
		if err != nil {
			return "create room fail", err
		}
		return "create room success", nil
	}
}

func makeGetRoomById(s service.RoomService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRoomByIdReq)
		room, err := s.GetRoomById(ctx, req.Id)
		if err != nil {
			return nil, err
		}
		return room, nil
	}
}

func makeHandoverRoom(s service.RoomService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(HandoverReq)
		err = s.HandoverRoom(ctx, req.IdRoom, req.IdDoctor)
		if err != nil {
			return "handover fail", err
		}
		return "handover success", nil
	}
}

func makeGetAllRoomByUser(s service.RoomService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRoomByUser)
		res, err := s.GetAllRoomByUser(ctx, req.IdDoctor)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func makeUpdateNumber(s service.RoomService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateNumberReq)
		err = s.UpdateNumber(ctx, req.IdRoom, req.Option)
		if err != nil {
			return nil, err
		}
		return "update success", nil
	}
}

func makeUpdateInformation(s service.RoomService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateInforReq)
		err = s.UpdateInformation(ctx, req.IdRoom, req.Name)
		if err != nil {
			return nil, err
		}
		return "update success", err
	}
}

func makeUpdateUseRoom(s service.RoomService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UpdateUseRoom)
		err = s.UpdateUseRoom(ctx, req.IdRoom, req.IdDoctor)
		if err != nil {
			return nil, err
		}
		return "update success", err
	}
}

func makeCreateListRoomsEndpoint(s service.RoomService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var bugs BugsRes
		reqs := request.(CreateListRoomReq)
		for _, req := range reqs.Rooms {
			check := s.ExistName(ctx, req.Name)
			if check {
				bugs.Bugs = append(bugs.Bugs, req.Name)
			}
		}
		if len(bugs.Bugs) > 0 {
			return BugsRes{
				Message: "this name is unvailble",
				Bugs:    bugs.Bugs,
			}, errors.New("can not insert room")
		}
		for _, req := range reqs.Rooms {
			err := s.CreateRoom(ctx, req.Name)
			if err != nil {
				return nil, err
			}
		}
		return "insert success", nil
	}
}

func makeGetAllByAdmin(s service.RoomService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		rooms, err := s.GetAllByAdmin(ctx)
		if err != nil {
			return nil, err
		}
		return rooms, nil
	}
}

func makeGetAllByCurrent(s service.RoomService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetCurrRoomReq)
		rooms, err := s.GetRoomByCurrent(ctx, req.Token)
		if err != nil {
			return nil, err
		}
		return rooms, nil
	}
}

type CreateRoomReq struct {
	Name string
}

type CreateListRoomReq struct {
	Rooms []CreateRoomReq
}

type BugsRes struct {
	Message string
	Bugs    []string
}

type GetRoomByIdReq struct {
	Id string
}

type HandoverReq struct {
	IdRoom   string
	IdDoctor string
}

type GetRoomByUser struct {
	IdDoctor string
}

type UpdateNumberReq struct {
	IdRoom string
	Option string
}

type UpdateInforReq struct {
	Name   string
	IdRoom string
}

type UpdateUseRoom struct {
	IdRoom   string
	IdDoctor string
}
