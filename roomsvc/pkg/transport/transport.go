package transport

import (
	"context"
	"encoding/json"
	endpoint2 "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httpTransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"roomsvc/pkg/endpoint"
	"strings"
)

func NewHTTPHandler(endpoints endpoint.RoomEndpoints, logger log.Logger) http.Handler {
	mux := http.NewServeMux()
	options := []httpTransport.ServerOption{
		httpTransport.ServerErrorEncoder(errorEncoder),
		httpTransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}
	mux.Handle("/room/create", httpTransport.NewServer(
		endpoints.CreateRoomEndpoint,
		decodeHTTPCreateRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	mux.Handle("/room/get", httpTransport.NewServer(
		endpoints.GetRoomByIdEndpoint,
		decodeHTTPGetRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	mux.Handle("/room/get-all-by-user", httpTransport.NewServer(
		endpoints.GetAllRoomByUserEndpoint,
		decodeHTTPGetAllByUserRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	mux.Handle("/room/handover", httpTransport.NewServer(
		endpoints.HandoverRoomEndpoint,
		decodeHTTPHandoverRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	mux.Handle("/room/update-number", httpTransport.NewServer(
		endpoints.UpdateNumberEndpoint,
		decodeHTTPUpdateNumberRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	mux.Handle("/room/update-information", httpTransport.NewServer(
		endpoints.UpdateInformationEndpoint,
		decodeHTTPUpdateInforRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	mux.Handle("/room/update-use-room", httpTransport.NewServer(
		endpoints.UpdateUseRoomEndpoint,
		decodeHTTPUseRoomRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	mux.Handle("/room/insert-list", httpTransport.NewServer(
		endpoints.CreateListRoomsEndpoint,
		decodeHTTPListRoomRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	mux.Handle("/room/get-all-by-admin", httpTransport.NewServer(
		endpoints.GetAllByAdmin,
		decodeHTTPNoBodyRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	mux.Handle("/room/get-all-by-current", httpTransport.NewServer(
		endpoints.GetAllByCurrent,
		decodeHTTPGetCurrRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	return mux
}

// ================================ //DECODER// ================================ //

func decodeHTTPCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.CreateRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.GetRoomByIdReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPGetAllByUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.GetRoomByUser
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPHandoverRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.HandoverReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPUpdateNumberRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.UpdateNumberReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPUpdateInforRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.UpdateInforReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPUseRoomRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.UpdateUseRoom
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPListRoomRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req []endpoint.CreateRoomReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return endpoint.CreateListRoomReq{Rooms: req}, err
}

func decodeHTTPNoBodyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeHTTPGetCurrRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.GetCurrRoomReq
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
