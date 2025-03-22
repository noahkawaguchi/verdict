package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/go-cmp/cmp"
	"github.com/noahkawaguchi/verdict/backend/internal/api"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

// mockDatastore implements the datastore interface for testing purposes.
type mockDatastore struct {
	PutPollMock    func(poll *models.Poll) error
	GetPollMock    func(pollID string) (*models.Poll, error)
	PutBallotMock  func(ballot *models.Ballot) error
	GetBallotsMock func(pollID string) ([]*models.Ballot, error)
}

func (m *mockDatastore) PutPoll(poll *models.Poll) error {
	if m.PutPollMock != nil {
		return m.PutPollMock(poll)
	}
	return nil
}

func (m *mockDatastore) GetPoll(pollID string) (*models.Poll, error) {
	if m.GetPollMock != nil {
		return m.GetPollMock(pollID)
	}
	return nil, nil
}

func (m *mockDatastore) PutBallot(ballot *models.Ballot) error {
	if m.PutBallotMock != nil {
		return m.PutBallotMock(ballot)
	}
	return nil
}

func (m *mockDatastore) GetBallots(pollID string) ([]*models.Ballot, error) {
	if m.GetBallotsMock != nil {
		return m.GetBallotsMock(pollID)
	}
	return nil, nil
}

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
		handler := api.Handler{DS: &mockDatastore{}, Req: test}
		resp := handler.Route()
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
		handler := api.Handler{DS: &mockDatastore{}, Req: test}
		resp := handler.Route()
		if resp.StatusCode != http.StatusNotFound {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		if resp.Body != `{"error":"path not found for method `+test.HTTPMethod+`: `+test.Path+`"}` {
			t.Error("unexpected response body:", resp.Body)
		}
	}
}

func quickJSON(anyStruct any) string {
	jsonBytes, _ := json.Marshal(anyStruct)
	return string(jsonBytes)
}

func TestCreatePollHandler_Invalid(t *testing.T) {
	tests := []struct {
		statusCode int
		errMsg     string
		body       string
	}{
		{
			http.StatusBadRequest,
			"invalid JSON",
			`{"prompt":"What is the best day of the week?,` +
				`"choices":["Wednesday", "Tuesday", "None of the above"]}`,
		},
		{
			http.StatusBadRequest,
			"choices must be unique",
			quickJSON(struct {
				Prompt  string   `json:"prompt"`
				Choices []string `json:"choices"`
			}{
				Prompt:  "What is the best day of the week?",
				Choices: []string{"Wednesday", "Tuesday", "None of the above", "Tuesday"},
			}),
		},
		{
			http.StatusInternalServerError,
			"failed to put the poll in the database",
			quickJSON(struct {
				Prompt  string   `json:"prompt"`
				Choices []string `json:"choices"`
			}{
				Prompt:  "What is the best day of the week?",
				Choices: []string{"Wednesday", "Tuesday", "None of the above"},
			}),
		},
	}

	for _, test := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Path:       "/poll",
			Body:       test.body,
		}
		handler := api.Handler{DS: &mockDatastore{
			PutPollMock: func(poll *models.Poll) error { return errors.New("mock error") },
		}, Req: req}
		resp := handler.Route()
		if resp.StatusCode != test.statusCode {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		if resp.Body != `{"error":"`+test.errMsg+`"}` {
			t.Error("unexpected response body:", resp.Body)
		}
	}
}

func TestCreatePollHandler_Valid(t *testing.T) {
	tests := []string{
		quickJSON(struct {
			Prompt  string   `json:"prompt"`
			Choices []string `json:"choices"`
		}{
			Prompt:  "What is the best day of the week?",
			Choices: []string{"Wednesday", "Tuesday", "None of the above"},
		}),
		quickJSON(struct {
			Prompt  string   `json:"prompt"`
			Choices []string `json:"choices"`
		}{
			Prompt:  "What is the worst day of the week?",
			Choices: []string{"Monday", "Thursday", "Either Monday or Thursday"},
		}),
	}

	for _, test := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Path:       "/poll",
			Body:       test,
		}
		handler := api.Handler{DS: &mockDatastore{}, Req: req}
		resp := handler.Route()
		if resp.StatusCode != http.StatusCreated {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		var respStruct struct {
			PollID string `json:"pollId"`
		}
		if err := json.Unmarshal([]byte(resp.Body), &respStruct); err != nil {
			t.Error("unexpected error unmarshaling JSON:", err)
		}
		if respStruct.PollID == "" {
			t.Error("unexpectedly empty poll ID in response body:", resp.Body)
		}
	}
}
