package transport

import (
	endpoint2 "bedsvc/pkg/endpoint"
	"bedsvc/pkg/service"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NewRoomClient(instance string) (service.RoomService, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	var getRoomByCurrEndpoint endpoint.Endpoint
	{
		getRoomByCurrEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/room/get-all-by-current"),
			encodeHTTPValidateRequest,
			decodeHTTPGetResponse,
		).Endpoint()
	}

	return endpoint2.RoomsEndpoints{
		GetRoomByCurrEndpoint: getRoomByCurrEndpoint,
	}, nil
}

func decodeHTTPGetResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp service.GetRoomsRes
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	fmt.Println("gh", r.Header.Get("Authorization"))
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = io.NopCloser(&buf)
	return nil
}

func encodeHTTPValidateRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	req := request.(endpoint2.GetRoomCurr)
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
