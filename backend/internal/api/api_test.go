package api_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-cmp/cmp"
	"github.com/noahkawaguchi/verdict/backend/internal/api"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

// mockDatastore implements the datastore interface for testing purposes.
type mockDatastore struct{}

func (m *mockDatastore) PutPoll(ctx context.Context, poll *models.Poll) error {
	return nil
}

func (m *mockDatastore) GetPoll(ctx context.Context, pollID string) (*models.Poll, error) {
	return nil, nil
}

func (m *mockDatastore) PutBallot(ctx context.Context, ballot *models.Ballot) error {
	return nil
}

func (m *mockDatastore) GetPollWithBallots(ctx context.Context, pollID string) (
	*models.Poll, []*models.Ballot, error,
) {
	return nil, nil, nil
}

var routerFunc = api.Router(&mockDatastore{})

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
		resp, _ := routerFunc(context.TODO(), test)
		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		if !cmp.Equal(
			resp.Headers,
			map[string]string{"Allow": "OPTIONS, GET, POST", "Content-Type": "application/json"},
		) {
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
		resp, _ := routerFunc(context.TODO(), test)
		if resp.StatusCode != http.StatusNotFound {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		if resp.Body != `{"error":"path not found for method `+test.HTTPMethod+`: `+test.Path+`"}` {
			t.Error("unexpected response body:", resp.Body)
		}
	}
}
