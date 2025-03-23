package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/backend/internal/api"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

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
		handler := api.Handler{Store: &mockDatastore{
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
		handler := api.Handler{Store: &mockDatastore{}, Req: req}
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

func TestGetPollInfoHandler_Invalid(t *testing.T) {
	tests := []struct {
		statusCode     int
		errMsg         string
		path           string
		pathParameters map[string]string
		getPollMock    func(pollID string) (*models.Poll, error)
	}{
		// Not testing for the poll ID being missing from the path specifically (not the
		// parameters map) because that is handled in the router
		{
			http.StatusBadRequest,
			"missing poll ID",
			"/poll/da932fe1-9a4c-4e07-adb3-9f66b4767050",
			map[string]string{"pollId": ""},
			nil,
		},
		{
			http.StatusBadRequest,
			"missing poll ID",
			"/poll/da932fe1-9a4c-4e07-adb3-9f66b4767050",
			map[string]string{},
			nil,
		},
		{
			http.StatusInternalServerError,
			"failed to get the poll from the database",
			"/poll/da932fe1-9a4c-4e07-adb3-9f66b4767050",
			map[string]string{"pollId": "da932fe1-9a4c-4e07-adb3-9f66b4767050"},
			func(pollID string) (*models.Poll, error) {
				return nil, errors.New("mock error")
			},
		},
		{
			http.StatusNotFound,
			"no poll found for the specified ID",
			"/poll/da932fe1-9a4c-4e07-adb3-9f66b4767050",
			map[string]string{"pollId": "da932fe1-9a4c-4e07-adb3-9f66b4767050"},
			func(pollID string) (*models.Poll, error) {
				return models.NewPoll("", []string{""}), nil
			},
		},
	}

	for _, test := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod:     http.MethodGet,
			Path:           test.path,
			PathParameters: test.pathParameters,
		}
		handler := api.Handler{Store: &mockDatastore{GetPollMock: test.getPollMock}, Req: req}
		resp := handler.Route()
		if resp.StatusCode != test.statusCode {
			t.Errorf("unexpected status code: expected %d, got %d", test.statusCode, resp.StatusCode)
		}
		if resp.Body != `{"error":"`+test.errMsg+`"}` {
			t.Error("unexpected response body:", resp.Body)
			t.Error(`expected: {"error":"` + test.errMsg + `"}`)
		}
	}
}

func TestGetPollInfoHandler_Valid(t *testing.T) {
	tests := []*models.Poll{
		models.NewPoll("What is the best day of the week?",
			[]string{"Wednesday", "Tuesday", "None of the above"}),
		models.NewPoll("What is the worst day of the week?",
			[]string{"Monday", "Thursday", "Either Monday or Thursday"}),
	}

	for _, test := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod:     http.MethodGet,
			Path:           "/poll/" + test.GetPollID(),
			PathParameters: map[string]string{"pollId": test.GetPollID()},
		}
		handler := api.Handler{Store: &mockDatastore{
			GetPollMock: func(pollID string) (*models.Poll, error) { return test, nil },
		}, Req: req}
		resp := handler.Route()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: expected %d, got %d", http.StatusOK, resp.StatusCode)
		}
		body, err := json.Marshal(test)
		if err != nil {
			t.Error("unexpected error marshaling JSON")
		}
		if resp.Body != string(body) {
			t.Error("unexpected response body:", resp.Body)
			t.Error("expected:", string(body))
		}
	}
}
