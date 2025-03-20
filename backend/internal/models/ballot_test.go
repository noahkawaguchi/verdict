package models_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/google/go-cmp/cmp"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func TestNewValidatedBallot_Invalid(t *testing.T) {
	tests := []struct {
		errMsg    string
		pollID    string
		userID    string
		rankOrder []int
	}{
		{
			errMsg:    "poll ID cannot be empty",
			userID:    "user1",
			rankOrder: []int{0, 1, 2},
		},
		{
			errMsg: "poll ID cannot be empty",
			userID: "user1",
		},
		{
			errMsg:    "poll ID cannot be empty",
			rankOrder: []int{0, 1, 2},
		},
		{
			errMsg:    "poll ID cannot be empty",
			pollID:    "",
			userID:    "user1",
			rankOrder: []int{0, 1, 2},
		},
		{
			errMsg: "there must be at least two rankings",
			pollID: "poll3",
			userID: "user3",
		},
		{"there must be at least two rankings", "poll2", "user2", []int{0}},
		{"there must be at least two rankings", "poll2", "user2", []int{}},
		{"not a valid rank order", "poll2", "user2", []int{3, 5, 1, 2, 4}},
		{"not a valid rank order", "poll3", "user3", []int{0, 1, 1, 1}},
	}
	for _, test := range tests {
		_, err := models.NewValidatedBallot(test.pollID, test.userID, test.rankOrder)
		if err == nil || err.Error() != test.errMsg {
			t.Errorf("expected error with message %q, got %v", test.errMsg, err)
		}
	}
}

func TestNewValidatedBallot_Valid(t *testing.T) {
	tests := []struct {
		pollID    string
		userID    string
		rankOrder []int
	}{
		{"poll1", "user1", []int{0, 1, 2}},
		{"poll2", "user2", []int{3, 0, 1, 2, 4}},
		{"poll2", "user4", []int{4, 1, 0, 3, 2}},
	}
	for _, test := range tests {
		if _, err := models.NewValidatedBallot(
			test.pollID, test.userID, test.rankOrder,
		); err != nil {
			t.Errorf("expected success, got %v", err)
		}
	}
}

func TestValidatedBallotFromJSON_Invalid(t *testing.T) {
	tests := []struct {
		errMsg     string
		jsonString string
	}{
		{"invalid JSON", `{"pollId": "poll1", "userId": "user1", "rankOrder": 0, 1, 2]}`},
		{"invalid JSON", `{"pollId": "poll1", "userId: "user1", "rankOrder": [0, 1, 2]}`},
		{"poll ID cannot be empty", `{"userId": "user1", "rankOrder": [0, 1, 2]}`},
		{"poll ID cannot be empty", `{"userId": "user1"}`},
		{"poll ID cannot be empty", `{"rankOrder": [0, 1, 2]}`},
		{"poll ID cannot be empty", `{"pollId": "", "userId": "user1", "rankOrder": [0, 1, 2]}`},
		{"there must be at least two rankings",
			`{"pollId": "poll2", "userId": "user2", "rankOrder": [0]}`},
		{"there must be at least two rankings",
			`{"pollId": "poll2", "userId": "user2", "rankOrder": []}`},
		{"there must be at least two rankings", `{"pollId": "poll3", "userId": "user3"}`},
		{"not a valid rank order",
			`{"pollId": "poll2", "userId": "user2", "rankOrder": [3, 5, 1, 2, 4]}`},
		{"not a valid rank order",
			`{"pollId": "poll3", "userId": "user3", "rankOrder": [0, 1, 1, 1]}`},
	}
	for _, test := range tests {
		_, err := models.ValidatedBallotFromJSON(test.jsonString)
		if err == nil || err.Error() != test.errMsg {
			t.Errorf("expected error with message %q, got %v", test.errMsg, err)
		}
	}
}

func TestValidatedBallotFromJSON_Valid(t *testing.T) {
	tests := []string{
		`{"pollId": "poll1", "userId": "user1", "rankOrder": [0, 1, 2]}`,
		`{"pollId": "poll2", "userId": "user2", "rankOrder": [3, 0, 1, 2, 4]}`,
		// Omitting user ID is valid
		`{"pollId": "poll2", "rankOrder": [4, 0, 3, 2, 1]}`,
		`{"pollId": "poll3", "rankOrder": [0, 3, 2, 1]}`,
	}
	for _, test := range tests {
		if _, err := models.ValidatedBallotFromJSON(test); err != nil {
			t.Errorf("expected success, got %v", err)
		}
	}
}

func TestBallotMarshalUnmarshalDynamoDBAttributeValue(t *testing.T) {
	tests := []struct {
		pollID    string
		userID    string
		rankOrder []int
	}{
		{"poll1", "user1", []int{0, 1}},
		{"poll2", "user3", []int{0, 1, 3, 2}},
		{"poll5", "user5", []int{3, 2, 0, 1, 5, 4}},
	}
	for _, test := range tests {
		inputBallot, _ := models.NewValidatedBallot(test.pollID, test.userID, test.rankOrder)
		av, err := attributevalue.MarshalMap(inputBallot)
		if err != nil {
			t.Errorf("failed to marshal map: %v", err)
		}
		var b *models.Ballot
		if err = attributevalue.UnmarshalMap(av, &b); err != nil {
			t.Errorf("failed to unmarshal map: %v", err)
		}
		if !cmp.Equal(inputBallot, b, cmp.AllowUnexported(models.Ballot{})) {
			t.Errorf("unexpected unmarshaled result: %+v", b)
		}
	}
}
