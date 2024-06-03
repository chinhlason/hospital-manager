package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"net"
	"net/http"
	"os"
	"roomsvc/pkg/endpoint"
	"roomsvc/pkg/service"
	"roomsvc/pkg/transport"
	"roomsvc/postgres"
)

func Start() (net.Listener, error) {
	httpAddr := flag.String("httpAddr", ":8881", "HTTP listen address")
	userAddr := flag.String("userAddr", "localhost:8880", "HTTP logs address")
	postgresURL := flag.String("postgresURL", "postgresql://sonnvt:sonnvt@localhost:5432/demo?sslmode=disable", "Postgres connection string")
	flag.Parse()
	fmt.Println("httpAddr", *httpAddr)
	fmt.Println("userAddr", *userAddr)
	fmt.Println("postgresURL", *postgresURL)

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "caller", log.DefaultCaller)

	var (
		queries       = postgres.Connect(*postgresURL)
		userClient, _ = transport.NewUserClient(*userAddr)
		roomSvc       = service.NewRoomService(queries, userClient)
		endpoints     = endpoint.MakeRoomServerEndpoints(roomSvc)
		httpHandler   = transport.NewHTTPHandler(endpoints, logger)
	)
	httpListener, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		return httpListener, err
	}
	logger.Log("transport", "Http", "addr", *httpAddr)
	return httpListener, http.Serve(httpListener, httpHandler)
}
func main() {
	httpListener, err := Start()
	if err != nil {
		httpListener.Close()
	}
}
