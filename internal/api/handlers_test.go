package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/internal/api"
	"github.com/noahkawaguchi/verdict/internal/voting"
)

func quickJSON(anyStruct any) string {
	jsonBytes, _ := json.Marshal(anyStruct)
	return string(jsonBytes)
}

func TestCreatePollHandler_Error(t *testing.T) {
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
				Prompt: "What is the best day of the week?",
				Choices: []string{
					"Wednesday",
					"Tuesday",
					"None of the above",
					"Tuesday",
				},
			}),
		},
		{
			http.StatusInternalServerError,
			"internal server error",
			quickJSON(struct {
				Prompt  string   `json:"prompt"`
				Choices []string `json:"choices"`
			}{
				Prompt:  "What is the best day of the week?",
				Choices: []string{"Wednesday", "Tuesday", "None of the above"},
			}),
		},
	}

	for _, tt := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Path:       "/poll",
			Body:       tt.body,
		}
		handler := api.NewHandler(&mockDatastore{
			PutPollMock: func(poll *voting.Poll) error { return errors.New("mock error") },
		}, req)
		resp := handler.Route()
		if resp.StatusCode != tt.statusCode {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		if resp.Body != `{"error":"`+tt.errMsg+`"}` {
			t.Error("unexpected response body:", resp.Body)
		}
	}
}

func TestCreatePollHandler_Success(t *testing.T) {
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

	for _, tt := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Path:       "/poll",
			Body:       tt,
		}
		handler := api.NewHandler(&mockDatastore{}, req)
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

func TestGetPollInfoHandler_Error(t *testing.T) {
	tests := []struct {
		statusCode     int
		errMsg         string
		pathParameters map[string]string
		getPollMock    func(pollID string) (*voting.Poll, error)
	}{
		// Not testing for the poll ID being missing from the path specifically (not the
		// parameters map) because that is handled in the router
		{
			http.StatusBadRequest,
			"missing poll ID",
			map[string]string{"pollId": ""},
			nil,
		},
		{
			http.StatusBadRequest,
			"missing poll ID",
			map[string]string{},
			nil,
		},
		{
			http.StatusInternalServerError,
			"internal server error",
			map[string]string{"pollId": "da932fe1-9a4c-4e07-adb3-9f66b4767050"},
			func(pollID string) (*voting.Poll, error) {
				return nil, errors.New("mock error")
			},
		},
		{
			http.StatusNotFound,
			"no poll found for poll ID da932fe1-9a4c-4e07-adb3-9f66b4767050",
			map[string]string{"pollId": "da932fe1-9a4c-4e07-adb3-9f66b4767050"},
			func(pollID string) (*voting.Poll, error) {
				return voting.NewPoll("", []string{""}), nil
			},
		},
	}

	for _, tt := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod:     http.MethodGet,
			Path:           "/poll/da932fe1-9a4c-4e07-adb3-9f66b4767050",
			PathParameters: tt.pathParameters,
		}
		handler := api.NewHandler(&mockDatastore{GetPollMock: tt.getPollMock}, req)
		resp := handler.Route()
		if resp.StatusCode != tt.statusCode {
			t.Errorf(
				"unexpected status code: expected %d, got %d",
				tt.statusCode,
				resp.StatusCode,
			)
		}
		if resp.Body != `{"error":"`+tt.errMsg+`"}` {
			t.Error("unexpected response body:", resp.Body)
			t.Error(`expected: {"error":"` + tt.errMsg + `"}`)
		}
	}
}

func TestGetPollInfoHandler_Success(t *testing.T) {
	tests := []*voting.Poll{
		voting.NewPoll("What is the best day of the week?",
			[]string{"Wednesday", "Tuesday", "None of the above"}),
		voting.NewPoll("What is the worst day of the week?",
			[]string{"Monday", "Thursday", "Either Monday or Thursday"}),
	}

	for _, tt := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod:     http.MethodGet,
			Path:           "/poll/" + tt.ID(),
			PathParameters: map[string]string{"pollId": tt.ID()},
		}
		handler := api.NewHandler(&mockDatastore{
			GetPollMock: func(pollID string) (*voting.Poll, error) { return tt, nil },
		}, req)
		resp := handler.Route()
		if resp.StatusCode != http.StatusOK {
			t.Errorf(
				"unexpected status code: expected %d, got %d",
				http.StatusOK,
				resp.StatusCode,
			)
		}
		body, err := json.Marshal(tt)
		if err != nil {
			t.Error("unexpected error marshaling JSON")
		}
		if resp.Body != string(body) {
			t.Error("unexpected response body:", resp.Body)
			t.Error("expected:", string(body))
		}
	}
}

func TestCastBallotHandler_Error(t *testing.T) {
	tests := []struct {
		statusCode int
		errMsg     string
		body       string
	}{
		{
			http.StatusBadRequest,
			"invalid JSON",
			`{"pollId":"poll1,"rankOrder":[1, 0, 3, 2]}`,
		},
		{
			http.StatusBadRequest,
			"not a valid rank order",
			quickJSON(struct {
				PollID    string `json:"pollId"`
				RankOrder []int  `json:"rankOrder"`
			}{
				PollID:    "poll22",
				RankOrder: []int{2, 4, 3, 1},
			}),
		},
		{
			http.StatusInternalServerError,
			"internal server error",
			quickJSON(struct {
				PollID    string `json:"pollId"`
				RankOrder []int  `json:"rankOrder"`
			}{
				PollID:    "poll23",
				RankOrder: []int{2, 0, 3, 1},
			}),
		},
	}

	for _, tt := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Path:       "/ballot",
			Body:       tt.body,
		}
		handler := api.NewHandler(&mockDatastore{
			PutBallotMock: func(ballot *voting.Ballot) error { return errors.New("mock error") },
		}, req)
		resp := handler.Route()
		if resp.StatusCode != tt.statusCode {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		if resp.Body != `{"error":"`+tt.errMsg+`"}` {
			t.Error("unexpected response body:", resp.Body)
		}
	}
}

func TestCastBallotHandler_Success(t *testing.T) {
	tests := []string{
		quickJSON(struct {
			PollID    string `json:"pollId"`
			RankOrder []int  `json:"rankOrder"`
		}{
			PollID:    "poll23",
			RankOrder: []int{2, 0, 3, 1},
		}),
		quickJSON(struct {
			PollID    string `json:"pollId"`
			UserID    string `json:"userId"`
			RankOrder []int  `json:"rankOrder"`
		}{
			PollID:    "poll24",
			UserID:    "user123",
			RankOrder: []int{0, 3, 1, 4, 2, 5},
		}),
	}

	for _, tt := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod: http.MethodPost,
			Path:       "/ballot",
			Body:       tt,
		}
		handler := api.NewHandler(&mockDatastore{}, req)
		resp := handler.Route()
		if resp.StatusCode != http.StatusCreated {
			t.Error("unexpected status code:", resp.StatusCode)
		}
		if resp.Body != `{"message":"successfully cast ballot"}` {
			t.Error("unexpected response body:", resp.Body)
		}
	}
}

func TestGetResultHandler_Error(t *testing.T) {
	tests := []struct {
		statusCode     int
		errMsg         string
		pathParameters map[string]string
		getPollMock    func(pollID string) (*voting.Poll, error)
		getBallotsMock func(pollID string) ([]*voting.Ballot, error)
	}{
		// Not testing for the poll ID being missing from the path specifically (not the
		// parameters map) because that is handled in the router
		{
			http.StatusBadRequest,
			"missing poll ID",
			map[string]string{"pollId": ""},
			nil,
			nil,
		},
		{
			http.StatusBadRequest,
			"missing poll ID",
			map[string]string{},
			nil,
			nil,
		},
		{
			http.StatusInternalServerError,
			"internal server error",
			map[string]string{"pollId": "da932fe1-9a4c-4e07-adb3-9f66b4767050"},
			func(pollID string) (*voting.Poll, error) {
				return nil, errors.New("mock error")
			},
			nil,
		},
		{
			http.StatusNotFound,
			"no poll found for poll ID da932fe1-9a4c-4e07-adb3-9f66b4767050",
			map[string]string{"pollId": "da932fe1-9a4c-4e07-adb3-9f66b4767050"},
			func(pollID string) (*voting.Poll, error) {
				return voting.NewPoll("", []string{""}), nil
			},
			nil,
		},
		{
			http.StatusInternalServerError,
			"internal server error",
			map[string]string{"pollId": "da932fe1-9a4c-4e07-adb3-9f66b4767050"},
			func(pollID string) (*voting.Poll, error) {
				return voting.NewPoll(
					"What is the best day of the week?",
					[]string{"Wednesday", "Tuesday", "None of the above"},
				), nil
			},
			func(pollID string) ([]*voting.Ballot, error) {
				return nil, errors.New("mock error")
			},
		},
		{
			http.StatusNotFound,
			"no ballots found for the specified poll",
			map[string]string{"pollId": "da932fe1-9a4c-4e07-adb3-9f66b4767050"},
			func(pollID string) (*voting.Poll, error) {
				return voting.NewPoll(
					"What is the best day of the week?",
					[]string{"Wednesday", "Tuesday", "None of the above"},
				), nil
			},
			func(pollID string) ([]*voting.Ballot, error) {
				return []*voting.Ballot{}, nil
			},
		},
	}

	for _, tt := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod:     http.MethodGet,
			Path:           "/result/da932fe1-9a4c-4e07-adb3-9f66b4767050",
			PathParameters: tt.pathParameters,
		}
		handler := api.NewHandler(&mockDatastore{
			GetPollMock:    tt.getPollMock,
			GetBallotsMock: tt.getBallotsMock,
		}, req)
		resp := handler.Route()
		if resp.StatusCode != tt.statusCode {
			t.Errorf(
				"unexpected status code: expected %d, got %d",
				tt.statusCode,
				resp.StatusCode,
			)
		}
		if resp.Body != `{"error":"`+tt.errMsg+`"}` {
			t.Error("unexpected response body:", resp.Body)
			t.Error(`expected: {"error":"` + tt.errMsg + `"}`)
		}
	}
}

func TestGetResultHandler_Success(t *testing.T) {
	tests := []struct {
		poll    *voting.Poll
		ballots []struct {
			userID    string
			rankOrder []int
		}
	}{
		{
			voting.NewPoll("What is the best day of the week?",
				[]string{"Wednesday", "Tuesday", "None of the above"}),
			[]struct {
				userID    string
				rankOrder []int
			}{{"user1", []int{0, 2, 1}}, {"user2", []int{0, 1, 2}}},
		},
		{
			voting.NewPoll(
				"What is the worst day of the week?",
				[]string{
					"Monday",
					"Thursday",
					"Either Monday or Thursday",
					"None of these",
				},
			),
			[]struct {
				userID    string
				rankOrder []int
			}{{"user4", []int{2, 3, 1, 0}}, {"user7", []int{2, 1, 0, 3}}},
		},
	}

	for _, tt := range tests {
		req := events.APIGatewayProxyRequest{
			HTTPMethod:     http.MethodGet,
			Path:           "/result/" + tt.poll.ID(),
			PathParameters: map[string]string{"pollId": tt.poll.ID()},
		}
		ballots := make([]*voting.Ballot, len(tt.ballots))
		for i, ballot := range tt.ballots {
			ballots[i] = voting.NewBallot(
				tt.poll.ID(),
				ballot.userID,
				ballot.rankOrder,
			)
		}
		handler := api.NewHandler(&mockDatastore{
			GetPollMock:    func(pollID string) (*voting.Poll, error) { return tt.poll, nil },
			GetBallotsMock: func(pollID string) ([]*voting.Ballot, error) { return ballots, nil },
		}, req)
		resp := handler.Route()
		if resp.StatusCode != http.StatusOK {
			t.Errorf(
				"unexpected status code: expected %d, got %d",
				http.StatusOK,
				resp.StatusCode,
			)
		}
		result, err := voting.NewResult(tt.poll, ballots)
		if err != nil {
			t.Error("unexpected error calculating result:", err)
		}
		body, err := json.Marshal(result)
		if err != nil {
			t.Error("unexpected error marshaling JSON")
		}
		if resp.Body != string(body) {
			t.Error("unexpected response body:", resp.Body)
			t.Error("expected:", string(body))
		}
	}
}
