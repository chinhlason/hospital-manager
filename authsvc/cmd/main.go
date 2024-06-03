package main

import (
	"datn-microservice/pkg/endpoint"
	"datn-microservice/pkg/service"
	"datn-microservice/pkg/transport"
	redis2 "datn-microservice/redis"
	"datn-microservice/scylladb"
	"flag"
	"github.com/go-kit/kit/log"
	"net"
	"net/http"
	"os"
)

func Start() (net.Listener, error) {
	httpAddr := flag.String("httpAddr", ":8880", "HTTP listen address")
	scyllaHost := flag.String("scyllaHost", "localhost:9042", "Scylla listen address")
	scyllaKS := flag.String("scyllaDB", "authsvc", "Scylla keyspace name")
	redisHost := flag.String("redisHost", "localhost:6379", "Redis listen address")
	flag.Parse()

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "caller", log.DefaultCaller)

	var (
		scyllaQueries = scylladb.Connect(*scyllaHost, *scyllaKS)
		redis         = redis2.ConnectRedis(*redisHost)
		svc           = service.NewAuthService(scyllaQueries, redis)
		endpoints     = endpoint.NewAuthEndpoints(svc)
		httpHandler   = transport.NewHTTPHandler(endpoints, logger)
	)
	httpListener, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		return nil, err
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
