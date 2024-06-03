package service

import (
	"context"
	"datn-microservice/scylladb/scylla/execute"
	"datn-microservice/utils"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"time"
)

type AuthServices interface {
	Register(ctx context.Context, username string, password string, fullname string, email string, phone string) error
	Login(ctx context.Context, username string, password string) (execute.Users, error)
	SetRefreshTokenToRedis(ctx context.Context, token string, userid string) error
	ValidateAccessToken(ctx context.Context, accessToken string) (execute.Users, error)
	RefreshToken(ctx context.Context, refreshToken string, accessToken string) (string, error)
	GetProfile(ctx context.Context, option, value string) (execute.Users, error)

	//Logout() error
}

type authservice struct {
	queries *execute.Queries
	redis   *redis.Client
}

func NewAuthService(queries *execute.Queries, redis *redis.Client) AuthServices {
	return &authservice{queries: queries, redis: redis}
}

func (s *authservice) Register(ctx context.Context, username string, password string, fullname string, email string, phone string) error {
	err := s.queries.InsertUser(ctx, username, password, fullname, email, phone)
	if err != nil {
		return err
	}
	return nil
}

func (s *authservice) Login(ctx context.Context, username string, password string) (execute.Users, error) {
	user, err := s.queries.Validate(ctx, username, password)
	if err != nil {
		return execute.Users{}, err
	}
	return user, nil
}

func (s *authservice) SetRefreshTokenToRedis(ctx context.Context, token string, userid string) error {
	key := fmt.Sprintf("refresh_token:%s", userid)
	expiration := time.Hour * 7 * 24
	err := s.redis.Set(ctx, key, token, expiration)
	if err.Err() != nil {
		return err.Err()
	}
	return nil
}

func (s *authservice) ValidateAccessToken(ctx context.Context, accessToken string) (execute.Users, error) {
	tokens, err := utils.ValidateToken(accessToken)
	if err != nil {
		return execute.Users{}, err
	}
	if tokens.Valid {
		claims := tokens.Claims.(jwt.MapClaims)
		userid := fmt.Sprintf("%v", claims["userid"])
		user, err := s.queries.GetUserByOption(ctx, userid, "id")
		if err != nil {
			return execute.Users{}, err
		}
		return user[0], nil
	}
	return execute.Users{}, errors.New("invalid token")
}

func (s *authservice) RefreshToken(ctx context.Context, refreshToken string, accessToken string) (string, error) {
	accessTokenParse, err := utils.ValidateToken(accessToken)
	if err != nil {
		return "", err
	}
	claims := accessTokenParse.Claims.(jwt.MapClaims)
	userid := fmt.Sprintf("%v", claims["userid"])
	role := fmt.Sprintf("%v", claims["role"])

	refreshTokenParse, err := utils.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}
	rfClaims := refreshTokenParse.Claims.(jwt.MapClaims)
	rfUserid := fmt.Sprintf("%v", rfClaims["userid"])
	if rfUserid == userid {
		key := fmt.Sprintf("refresh_token:%s", userid)
		storedRefreshToken := s.redis.Get(ctx, key)
		if refreshToken == storedRefreshToken.Val() {
			newAccessToken, err := utils.GenToken(userid, role, time.Hour*7*24)
			if err != nil {
				return "", err
			}
			return newAccessToken, nil
		}
		return "", errors.New("refresh token is invalid")
	}
	return "", errors.New("access token is invalid")
}

func (s *authservice) GetProfile(ctx context.Context, option, value string) (execute.Users, error) {
	user, err := s.queries.GetUserByOption(ctx, value, option)
	if err != nil {
		return execute.Users{}, err
	}
	if len(user) == 0 {
		return execute.Users{}, errors.New("no user data found")
	}
	return user[0], nil
}
