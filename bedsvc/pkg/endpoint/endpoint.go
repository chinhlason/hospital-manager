package endpoint

import (
	"bedsvc/pkg/service"
	"context"
	"github.com/go-kit/kit/endpoint"
)

type BedEndpoints struct {
	CreateBedEndpoint    endpoint.Endpoint
	GetBedInRoomEndpoint endpoint.Endpoint
}

func NewBedServerEndpoint(b service.BedServices) BedEndpoints {
	return BedEndpoints{
		CreateBedEndpoint:    makeCreateBedEndpoint(b),
		GetBedInRoomEndpoint: makeGetBedInRoomEndpoint(b),
	}
}

func makeCreateBedEndpoint(b service.BedServices) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(service.CreateBedReq)
		err = b.CreateBed(ctx, req.RoomName, req.BedName, req.Token)
		if err != nil {
			return nil, err
		}
		return "create success", nil
	}
}

func makeGetBedInRoomEndpoint(b service.BedServices) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRoomCurr)
		room, err := b.GetBedInRoom(ctx, "", req.Token)
		if err != nil {
			return nil, err
		}
		return room, nil
	}
}
