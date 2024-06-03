package endpoint

import (
	"context"
	"datn-microservice/pkg/service"
	"datn-microservice/scylladb/scylla/execute"
	"datn-microservice/utils"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type AuthEndpoints struct {
	RegisterEndpoint       endpoint.Endpoint
	LoginEndpoint          endpoint.Endpoint
	ValidateEndpoint       endpoint.Endpoint
	RefreshTokenEndpoint   endpoint.Endpoint
	GetProfileUserEndpoint endpoint.Endpoint
}

func NewAuthEndpoints(s service.AuthServices) AuthEndpoints {
	return AuthEndpoints{
		RegisterEndpoint:       makeRegisterEndpoint(s),
		LoginEndpoint:          makeLoginEndpoint(s),
		ValidateEndpoint:       makeValidateEndpoint(s),
		RefreshTokenEndpoint:   makeRefreshTokenEndpoint(s),
		GetProfileUserEndpoint: makeGetProfileEndpoint(s),
	}
}

func makeRegisterEndpoint(s service.AuthServices) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(RegisterReq)
		err = s.Register(ctx, req.Username, req.Password, req.Fullname, req.Email, req.Phone)
		if err != nil {
			return nil, err
		}
		return "register successfully", nil
	}
}

func makeLoginEndpoint(s service.AuthServices) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(LoginReq)
		user, err := s.Login(ctx, req.Username, req.Password)
		if err != nil {
			return LoginRes{}, err
		}
		accessToken, err := utils.GenToken(user.Id.String(), user.Role, time.Hour*7*24)
		if err != nil {
			return LoginRes{}, err
		}
		refreshToken, err := utils.GenToken(user.Id.String(), user.Role, time.Hour*24*7)
		if err != nil {
			return LoginRes{}, err
		}
		err = s.SetRefreshTokenToRedis(ctx, refreshToken, user.Id.String())
		if err != nil {
			return LoginRes{}, err
		}
		return LoginRes{
			User:         user,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}, nil
	}
}

func makeValidateEndpoint(s service.AuthServices) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ValidateReq)
		user, err := s.ValidateAccessToken(ctx, req.AccessToken)
		if err != nil {
			return ValidateRes{
				Subject: "",
				Extra:   execute.Users{},
				Err:     err,
			}, nil
		}
		return ValidateRes{
			Subject: user.Id.String(),
			Extra: execute.Users{
				Id:       user.Id,
				Username: user.Username,
				Role:     user.Role,
			},
			Err: nil,
		}, nil
	}
}

func makeRefreshTokenEndpoint(s service.AuthServices) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(RefreshTokenReq)
		newAccessToken, err := s.RefreshToken(ctx, req.RefreshToken, req.AccessToken)
		if err != nil {
			return RefreshTokenRes{}, err
		}
		return RefreshTokenRes{
			NewAccessToken: newAccessToken,
		}, nil
	}
}

func makeGetProfileEndpoint(s service.AuthServices) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetProfileReq)
		user, err := s.GetProfile(ctx, req.Option, req.Value)
		if err != nil {
			return execute.Users{}, err
		}
		return user, nil
	}
}

type GetProfileReq struct {
	Option string
	Value  string
}

type RefreshTokenReq struct {
	RefreshToken string
	AccessToken  string
}

type RefreshTokenRes struct {
	NewAccessToken string
}

type ValidateReq struct {
	AccessToken string
}

type ValidateRes struct {
	Subject string
	Extra   execute.Users
	Err     error
}

type LoginReq struct {
	Username string
	Password string
}

type LoginRes struct {
	User         execute.Users
	AccessToken  string
	RefreshToken string
}

type RegisterReq struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Fullname string `validate:"required"`
	Email    string `validate:"required"`
	Phone    string `validate:"required"`
}
