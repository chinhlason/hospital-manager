package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"io"
	"net/http"
	"net/url"
	endpoint2 "roomsvc/pkg/endpoint"
	"roomsvc/pkg/service"
	"strings"
)

func NewUserClient(instance string) (service.UserService, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	var getUserInforEndpoint endpoint.Endpoint
	{
		getUserInforEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/profile"),
			encodeHTTPGenericRequest,
			decodeHTTPGetResponse,
		).Endpoint()
	}

	var validateUserEndpoint endpoint.Endpoint
	{
		validateUserEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/validate-token"),
			encodeHTTPValidateRequest,
			decodeHTTPValidateResponse,
		).Endpoint()
	}

	return endpoint2.UserEndpoints{
		GetUserInformationEndpoint: getUserInforEndpoint,
		ValidatUserEndpoint:        validateUserEndpoint,
	}, nil
}

func decodeHTTPGetResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp service.User
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPValidateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp service.ValidateRes
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = io.NopCloser(&buf)
	return nil
}

func encodeHTTPValidateRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	req := request.(endpoint2.GetCurrRoomReq)
	token := "Bearer " + req.Token
	r.Header.Set("Authorization", token)
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = io.NopCloser(&buf)
	return nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
