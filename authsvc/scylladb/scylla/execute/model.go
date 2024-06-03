package execute

import (
	"github.com/gocql/gocql"
	"time"
)

type Users struct {
	Id       gocql.UUID `json:"id"`
	Username string     `json:"username"`
	Password string     `json:"-"`
	Fullname string     `json:"fullname"`
	Email    string     `json:"email"`
	Phone    string     `json:"phone"`
	Role     string     `json:"role"`
	CreateAt time.Time  `json:"create_at"`
	UpdateAt time.Time  `json:"update_at"`
}
