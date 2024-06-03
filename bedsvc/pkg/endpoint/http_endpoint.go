package endpoint

import (
	"bedsvc/pkg/service"
	"context"
	"github.com/go-kit/kit/endpoint"
)

type RoomsEndpoints struct {
	GetRoomByCurrEndpoint endpoint.Endpoint
}

func (r RoomsEndpoints) GetRoomByCurr(ctx context.Context, token string) (service.GetRoomsRes, error) {
	resp, err := r.GetRoomByCurrEndpoint(ctx, GetRoomCurr{
		Token: token,
	})
	if err != nil {
		return service.GetRoomsRes{}, err
	}
	res := resp.(service.GetRoomsRes)
	return res, nil
}

type GetRoomCurr struct {
	Token string
}
