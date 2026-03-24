package api_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-cmp/cmp"
	"github.com/noahkawaguchi/verdict/internal/api"
)

func TestRouter_MethodNotAllowed(t *testing.T) {
	tests := []events.APIGatewayProxyRequest{
		{
			Path:       "/polls",
			HTTPMethod: http.MethodPut,
		},
		{
			Path:       "/ballots",
			HTTPMethod: http.MethodPut,
		},
		{
			Path:       "/polls",
			HTTPMethod: http.MethodPatch,
		},
		{
			Path:       "/ballots",
			HTTPMethod: http.MethodPatch,
		},
		{
			Path:       "/polls",
			HTTPMethod: http.MethodDelete,
		},
		{
			Path:       "/ballots",
			HTTPMethod: http.MethodDelete,
		},
	}

	for _, tt := range tests {
		handler := api.NewHandler(&mockDatastore{}, tt)
		resp := handler.Route()
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		expectedHeaders := map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  os.Getenv("FRONTEND_URL"),
			"Access-Control-Allow-Methods": "OPTIONS,GET,POST",
			"Access-Control-Allow-Headers": "Content-Type",
			"Allow":                        "OPTIONS, GET, POST",
		}
		if !cmp.Equal(resp.Headers, expectedHeaders) {
			t.Error("unexpected headers:", resp.Headers)
		}
		if resp.Body != `{"error":"method `+tt.HTTPMethod+` not allowed"}` {
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

	for _, tt := range tests {
		handler := api.NewHandler(&mockDatastore{}, tt)
		resp := handler.Route()
		if resp.StatusCode != http.StatusNotFound {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		if resp.Body != `{"error":"path not found for method `+tt.HTTPMethod+`: `+tt.Path+`"}` {
			t.Error("unexpected response body:", resp.Body)
		}
	}
}

func TestRouter_HealthCheck(t *testing.T) {
	req := events.APIGatewayProxyRequest{
		Path:       "/health",
		HTTPMethod: http.MethodGet,
	}
	resp := api.NewHandler(&mockDatastore{}, req).Route()
	if resp.StatusCode != http.StatusOK {
		t.Error("unexpected status code:", resp.StatusCode)
	}
	if resp.Body != "The function is available!" {
		t.Error("unexpected response body:", resp.Body)
	}
}
