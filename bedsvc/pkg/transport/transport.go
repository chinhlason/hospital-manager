package transport

import (
	"bedsvc/pkg/endpoint"
	"bedsvc/pkg/service"
	"context"
	"encoding/json"
	endpoint2 "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httpTransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"strings"
)

func NewHTTPHandler(endpoints endpoint.BedEndpoints, logger log.Logger) http.Handler {
	mux := http.NewServeMux()
	options := []httpTransport.ServerOption{
		httpTransport.ServerErrorEncoder(errorEncoder),
		httpTransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}
	mux.Handle("/bed/create", httpTransport.NewServer(
		endpoints.CreateBedEndpoint,
		decodeHTTPCreateRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	mux.Handle("/bed/get", httpTransport.NewServer(
		endpoints.GetBedInRoomEndpoint,
		decodeHTTPGetRequest,
		encodeHTTPGenericResponse,
		options...),
	)

	return mux
}

// ================================ //DECODER// ================================ //

func decodeHTTPCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req service.CreateBedReq
	err := json.NewDecoder(r.Body).Decode(&req)
	token := r.URL.Query().Get("token")
	if token == "" {
		token = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	}
	req.Token = token
	return req, err
}

func decodeHTTPGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.GetRoomCurr
	token := r.URL.Query().Get("token")
	if token == "" {
		token = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	}
	req.Token = token
	return req, nil
}

// ================================ //ENCODER// ================================ //

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint2.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errorWrapper{
		Error: err.Error(),
		Code:  http.StatusBadRequest,
	})
}

type errorWrapper struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
