package scylladb

import (
	"context"
	"datn-microservice/scylladb/cql"
	"datn-microservice/scylladb/scylla"
	"datn-microservice/scylladb/scylla/execute"
	"github.com/scylladb/gocqlx/v2/migrate"
)

var queries *execute.Queries

func Connect(host string, ksname string) *execute.Queries {
	ctx := context.Background()
	manager := scylla.NewManager(host, ksname)
	err := manager.CreateKeyspace()
	if err != nil {
		panic(err)
	}
	session, err := manager.Connect()
	if err != nil {
		panic(err)
	}
	err = migrate.FromFS(ctx, session, cql.Files)
	if err != nil {
		panic(err)
	}
	queries = execute.New(session, manager.ScyllaKeyspace)
	return queries
}

func Queries() *execute.Queries { return queries }
