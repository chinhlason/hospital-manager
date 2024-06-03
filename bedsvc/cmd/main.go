package main

import (
	"bedsvc/pkg/endpoint"
	"bedsvc/pkg/service"
	"bedsvc/pkg/transport"
	"bedsvc/postgres"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"net"
	"net/http"
	"os"
)

func Start() (net.Listener, error) {
	httpAddr := flag.String("httpAddr", ":8882", "HTTP listen address")
	roomAddr := flag.String("roomAddr", "localhost:8881", "HTTP logs address")
	postgresURL := flag.String("postgresURL", "postgresql://sonnvt:sonnvt@localhost:5432/demo?sslmode=disable", "Postgres connection string")
	flag.Parse()
	fmt.Println("httpAddr", *httpAddr)
	fmt.Println("userAddr", *roomAddr)
	fmt.Println("postgresURL", *postgresURL)

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "caller", log.DefaultCaller)

	var (
		queries       = postgres.Connect(*postgresURL)
		roomClient, _ = transport.NewRoomClient(*roomAddr)
		bedSvc        = service.NewBedService(queries, roomClient)
		endpoints     = endpoint.NewBedServerEndpoint(bedSvc)
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
