package api_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-cmp/cmp"
	"github.com/noahkawaguchi/verdict/backend/internal/api"
)

func TestRouter_MethodNotAllowed(t *testing.T) {
	tests := []events.APIGatewayProxyRequest{
		{
			Path:       "/poll",
			HTTPMethod: http.MethodPut,
		},
		{
			Path:       "/ballot",
			HTTPMethod: http.MethodPut,
		},
		{
			Path:       "/poll",
			HTTPMethod: http.MethodPatch,
		},
		{
			Path:       "/ballot",
			HTTPMethod: http.MethodPatch,
		},
		{
			Path:       "/poll",
			HTTPMethod: http.MethodDelete,
		},
		{
			Path:       "/ballot",
			HTTPMethod: http.MethodDelete,
		},
	}

	for _, test := range tests {
		handler := api.NewHandler(&mockDatastore{}, test)
		resp := handler.Route()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		expectedHeaders := map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  os.Getenv("FRONTEND_URL"),
			"Access-Control-Allow-Methods": "OPTIONS,GET,POST",
			"Access-Control-Allow-Headers": "Content-Type,Authorization",
			"Allow":                        "OPTIONS, GET, POST",
		}
		if !cmp.Equal(resp.Headers, expectedHeaders) {
			t.Error("unexpected headers:", resp.Headers)
		}
		if resp.Body != `{"error":"method `+test.HTTPMethod+` not allowed"}` {
			t.Error("unexpected response body:", resp.Body)
		}
	}
}

func TestRouter_PathNotFound(t *testing.T) {
	tests := []events.APIGatewayProxyRequest{
		{
			Path:       "/pole",
			HTTPMethod: http.MethodPost,
		},
		{
			Path:       "/ballot-cast",
			HTTPMethod: http.MethodPost,
		},
		{
			Path:       "/election",
			HTTPMethod: http.MethodGet,
		},
		{
			Path:       "/poll-voting",
			HTTPMethod: http.MethodGet,
		},
	}

	for _, test := range tests {
		handler := api.NewHandler(&mockDatastore{}, test)
		resp := handler.Route()
		if resp.StatusCode != http.StatusNotFound {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		if resp.Body != `{"error":"path not found for method `+test.HTTPMethod+`: `+test.Path+`"}` {
			t.Error("unexpected response body:", resp.Body)
		}
	}
}
