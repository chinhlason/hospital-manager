package execute

import (
	"context"
	"datn-microservice/utils"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/scylladb/gocqlx/v2/qb"
	"time"
)

func (q *Queries) CreateAdminAccount() error {
	tableName := fmt.Sprintf("%s.users", q.keyspace)
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, err := gocql.ParseUUID(uuid.New().String())
	if err != nil {
		panic(err)
	}
	password := utils.HashAndSalt([]byte("son"))
	insert := &Users{
		Id:       id,
		Username: "Admin",
		Password: password,
		Fullname: "Admin",
		Email:    "sonnvt05@gmail.com",
		Phone:    "0923151911",
		Role:     "ADMIN",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	stmt := qb.Insert(tableName).
		Columns("id", "username", "password",
			"fullname", "email", "phone", "role", "create_at", "update_at").
		Query(q.session)
	stmt.BindStruct(insert)
	if err := stmt.ExecRelease(); err != nil {
		return err
	}
	return nil
}

func (q *Queries) GetUserById(ctx context.Context, id string) (Users, error) {
	tableName := fmt.Sprintf("%s.users", q.keyspace)
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var user Users
	stmt, names := qb.Select(tableName).
		Where(qb.Eq(id)).
		ToCql()
	query := q.session.Query(stmt, names).BindMap(qb.M{
		"id": id,
	})
	if err := query.GetRelease(&user); err != nil {
		return Users{}, err
	}
	return user, nil
}

func (q *Queries) GetUserByOption(ctx context.Context, value string, option string) ([]Users, error) {
	tableName := fmt.Sprintf("%s.users", q.keyspace)
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var users []Users
	stmt, names := qb.Select(tableName).
		Where(qb.Eq(option)).
		AllowFiltering().
		ToCql()
	query := q.session.Query(stmt, names).BindMap(qb.M{
		option: value,
	})
	if err := query.SelectRelease(&users); err != nil {
		return nil, err
	}
	return users, nil
}

func (q *Queries) InsertUser(ctx context.Context, username string, password string, fullname string, email string, phone string) error {
	tableName := fmt.Sprintf("%s.users", q.keyspace)
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	checkExistUsername, err := q.GetUserByOption(ctx, username, "username")
	if err != nil {
		return err
	}
	if len(checkExistUsername) > 0 {
		return errors.New("Username existed")
	}

	checkExistEmail, err := q.GetUserByOption(ctx, email, "email")
	if err != nil {
		return err
	}
	if len(checkExistEmail) > 0 {
		return errors.New("Email existed")
	}

	id, err := gocql.ParseUUID(uuid.New().String())
	encodedPassword := utils.HashAndSalt([]byte(password))
	insert := &Users{
		Id:       id,
		Username: username,
		Password: encodedPassword,
		Fullname: fullname,
		Email:    email,
		Phone:    phone,
		Role:     "USER",
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
	}
	stmt := qb.Insert(tableName).
		Columns("id", "username", "password",
			"fullname", "email", "phone", "role", "create_at", "update_at").
		Query(q.session)
	stmt.BindStruct(insert)
	if err := stmt.ExecRelease(); err != nil {
		return err
	}
	return nil
}

func (q *Queries) Validate(ctx context.Context, username, password string) (Users, error) {
	_, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	users, err := q.GetUserByOption(ctx, username, "username")
	if err != nil {
		return Users{}, err
	}
	if len(users) == 0 {
		return Users{}, errors.New("username or password is incorrect")
	}
	isValid := utils.ComparePassword(users[0].Password, []byte(password))
	if !isValid {
		return Users{}, errors.New("username or Password is incorrect")
	}
	return users[0], nil
}
