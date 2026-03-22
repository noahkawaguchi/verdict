package models_test

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func TestValidateBallot_Invalid(t *testing.T) {
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
		ballot := models.NewBallot(test.pollID, test.userID, test.rankOrder)
		if err := ballot.Validate(); err == nil || err.Error() != test.errMsg {
			t.Errorf("expected error with message %q, got %v", test.errMsg, err)
		}
	}
}

func TestValidateBallot_Valid(t *testing.T) {
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
		ballot := models.NewBallot(test.pollID, test.userID, test.rankOrder)
		if err := ballot.Validate(); err != nil {
			t.Errorf("expected success, got %v", err)
		}
	}
}

func TestBallotUnmarshalJSON(t *testing.T) {
	tests := []struct {
		pollID     string
		userID     string
		rankOrder  []int
		jsonString string
	}{
		{"poll1", "user1", []int{0, 1, 2},
			`{"pollId": "poll1", "userId": "user1", "rankOrder": [0, 1, 2]}`},
		{"poll2", "user2", []int{3, 0, 1, 2, 4},
			`{"pollId": "poll2", "userId": "user2", "rankOrder": [3, 0, 1, 2, 4]}`},
		// Omitting user ID is valid
		{
			pollID:     "poll2",
			rankOrder:  []int{4, 0, 3, 2, 1},
			jsonString: `{"pollId": "poll2", "rankOrder": [4, 0, 3, 2, 1]}`,
		},
		{
			pollID:     "poll3",
			rankOrder:  []int{0, 3, 2, 1},
			jsonString: `{"pollId": "poll3", "rankOrder": [0, 3, 2, 1]}`,
		},
	}
	for _, test := range tests {
		var unmarshaledBallot *models.Ballot
		if err := json.Unmarshal([]byte(test.jsonString), &unmarshaledBallot); err != nil {
			t.Errorf("expected success, got %v", err)
		}
		if test.userID != "" { // User ID provided cases
			constructedBallot := models.NewBallot(test.pollID, test.userID, test.rankOrder)
			if !cmp.Equal(
				unmarshaledBallot,
				constructedBallot,
				cmp.AllowUnexported(models.Ballot{}),
			) {
				t.Error("unexpected unmarshaled ballot:", unmarshaledBallot)
				t.Error("expected ballot:", constructedBallot)
			}
		} else { // User ID automatically generated cases
			userID := "dummy user ID"
			constructedBallot := models.NewBallot(test.pollID, userID, test.rankOrder)
			if !cmp.Equal(
				unmarshaledBallot,
				constructedBallot,
				cmp.AllowUnexported(models.Ballot{}),
				cmpopts.IgnoreFields(models.Ballot{}, "userID"),
			) {
				t.Error("unexpected unmarshaled ballot:", unmarshaledBallot)
			}
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
		inputBallot := models.NewBallot(test.pollID, test.userID, test.rankOrder)
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
