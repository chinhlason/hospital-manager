package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"roomsvc/pkg/service"
)

type UserEndpoints struct {
	GetUserInformationEndpoint endpoint.Endpoint
	ValidatUserEndpoint        endpoint.Endpoint
}

func (u UserEndpoints) GetUserInformation(ctx context.Context, value, option string) (service.User, error) {
	resp, err := u.GetUserInformationEndpoint(ctx, GetUserInforReq{
		Option: option,
		Value:  value,
	})
	if err != nil {
		return service.User{}, err
	}
	res := resp.(service.User)
	return res, nil
}

func (u UserEndpoints) ValidateUser(ctx context.Context, token string) (string, error) {
	resp, err := u.ValidatUserEndpoint(ctx, GetCurrRoomReq{
		Token: token,
	})
	if err != nil {
		return "", err
	}
	temp := resp.(service.ValidateRes)
	return temp.Subject, nil
}

type GetUserInforReq struct {
	Option string
	Value  string
}

type GetCurrRoomReq struct {
	Token string
}
