package transport

import (
	"context"
	"datn-microservice/pkg/endpoint"
	"encoding/json"
	endpoint2 "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	httpTransport "github.com/go-kit/kit/transport/http"
	"net/http"
	"strings"
)

func NewHTTPHandler(endpoints endpoint.AuthEndpoints, logger log.Logger) http.Handler {
	mux := http.NewServeMux()
	options := []httpTransport.ServerOption{
		httpTransport.ServerErrorEncoder(errorEncoder),
		httpTransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}
	mux.Handle("/register", httpTransport.NewServer(
		endpoints.RegisterEndpoint,
		decodeHTTPRegisterRequest,
		encodeHTTPAuthResponse,
		options...),
	)
	mux.Handle("/login", httpTransport.NewServer(
		endpoints.LoginEndpoint,
		decodeHTTPLoginRequest,
		encodeHTTPAuthResponse,
		options...),
	)
	mux.Handle("/validate-token", httpTransport.NewServer(
		endpoints.ValidateEndpoint,
		decodeHTTPValidateRequest,
		encodeHTTPValidateResponse,
		options...),
	)
	mux.Handle("/refresh-token", httpTransport.NewServer(
		endpoints.RefreshTokenEndpoint,
		decodeHTTPRefreshTokenRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	mux.Handle("/profile", httpTransport.NewServer(
		endpoints.GetProfileUserEndpoint,
		decodeHTTPProfileRequest,
		encodeHTTPGenericResponse,
		options...),
	)
	return mux
}

// ================================ //DECODER// ================================ //

func decodeHTTPRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.RegisterReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.LoginReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPValidateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.ValidateReq
	token := r.URL.Query().Get("token")
	if token == "" {
		token = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	}
	req.AccessToken = token
	return req, nil
}

func decodeHTTPRefreshTokenRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.RefreshTokenReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPProfileRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.GetProfileReq
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

// ================================ //ENCODER// ================================ //

func encodeHTTPAuthResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeHTTPValidateResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint2.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	res := response.(endpoint.ValidateRes)
	if res.Err != nil {
		w.WriteHeader(http.StatusUnauthorized)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

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
