package service

import (
	"context"
	"roomsvc/postgres/execute"
	"time"
)

type UserService interface {
	GetUserInformation(ctx context.Context, value, option string) (User, error)
	ValidateUser(ctx context.Context, token string) (string, error)
}

type User struct {
	Id       string    `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"`
	Fullname string    `json:"fullname"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	Role     string    `json:"role"`
	CreateAt time.Time `json:"create_at"`
	UpdateAt time.Time `json:"update_at"`
}

type GetAllRoomsRes struct {
	User User
	Room []execute.Room
}

type ValidateRes struct {
	Subject string
	Extra   User
	Err     error
}
