package postgres

import (
	"bedsvc/postgres/execute"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func Connect(connStr string) *execute.Queries {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return execute.New(db)
}
