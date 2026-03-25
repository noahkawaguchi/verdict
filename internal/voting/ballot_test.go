package voting_test

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/noahkawaguchi/verdict/internal/voting"
)

func TestValidateBallot_Invalid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, errMsg, pollID, userID string
		rankOrder                    []int
	}{
		{
			name:      "empty poll ID with rank order",
			errMsg:    "poll ID cannot be empty",
			userID:    "user1",
			rankOrder: []int{0, 1, 2},
		},
		{
			name:   "empty poll ID no rank order",
			errMsg: "poll ID cannot be empty",
			userID: "user1",
		},
		{
			name:      "empty poll ID no user ID",
			errMsg:    "poll ID cannot be empty",
			rankOrder: []int{0, 1, 2},
		},
		{
			name:      "explicit empty string poll ID",
			errMsg:    "poll ID cannot be empty",
			pollID:    "",
			userID:    "user1",
			rankOrder: []int{0, 1, 2},
		},
		{
			name:   "no rank order",
			errMsg: "there must be at least two rankings",
			pollID: "poll3",
			userID: "user3",
		},
		{
			name:      "single rank",
			errMsg:    "there must be at least two rankings",
			pollID:    "poll2",
			userID:    "user2",
			rankOrder: []int{0},
		},
		{
			name:      "empty rank order",
			errMsg:    "there must be at least two rankings",
			pollID:    "poll2",
			userID:    "user2",
			rankOrder: []int{},
		},
		{
			name:      "rank out of range",
			errMsg:    "not a valid rank order",
			pollID:    "poll2",
			userID:    "user2",
			rankOrder: []int{3, 5, 1, 2, 4},
		},
		{
			name:      "duplicate rank values",
			errMsg:    "not a valid rank order",
			pollID:    "poll3",
			userID:    "user3",
			rankOrder: []int{0, 1, 1, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ballot := voting.NewBallot(tt.pollID, tt.userID, tt.rankOrder)
			if err := ballot.Validate(); err == nil || err.Error() != tt.errMsg {
				t.Errorf("expected error with message %q, got %v", tt.errMsg, err)
			}
		})
	}
}

func TestValidateBallot_Valid(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, pollID, userID string
		rankOrder            []int
	}{
		{"three rankings", "poll1", "user1", []int{0, 1, 2}},
		{"five rankings", "poll2", "user2", []int{3, 0, 1, 2, 4}},
		{"five rankings alt", "poll2", "user4", []int{4, 1, 0, 3, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ballot := voting.NewBallot(tt.pollID, tt.userID, tt.rankOrder)
			if err := ballot.Validate(); err != nil {
				t.Errorf("expected success, got %v", err)
			}
		})
	}
}

func TestBallotUnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, pollID, userID string
		rankOrder            []int
		jsonString           string
	}{
		{
			name:       "with user ID",
			pollID:     "poll1",
			userID:     "user1",
			rankOrder:  []int{0, 1, 2},
			jsonString: `{"pollId": "poll1", "userId": "user1", "rankOrder": [0, 1, 2]}`,
		},
		{
			name:       "with user ID alt",
			pollID:     "poll2",
			userID:     "user2",
			rankOrder:  []int{3, 0, 1, 2, 4},
			jsonString: `{"pollId": "poll2", "userId": "user2", "rankOrder": [3, 0, 1, 2, 4]}`,
		},
		// Omitting user ID is valid
		{
			name:       "without user ID",
			pollID:     "poll2",
			rankOrder:  []int{4, 0, 3, 2, 1},
			jsonString: `{"pollId": "poll2", "rankOrder": [4, 0, 3, 2, 1]}`,
		},
		{
			name:       "without user ID alt",
			pollID:     "poll3",
			rankOrder:  []int{0, 3, 2, 1},
			jsonString: `{"pollId": "poll3", "rankOrder": [0, 3, 2, 1]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var unmarshaledBallot *voting.Ballot
			if err := json.Unmarshal(
				[]byte(tt.jsonString),
				&unmarshaledBallot,
			); err != nil {
				t.Errorf("expected success, got %v", err)
			}
			if tt.userID != "" { // User ID provided cases
				constructedBallot := voting.NewBallot(
					tt.pollID,
					tt.userID,
					tt.rankOrder,
				)
				if !cmp.Equal(
					unmarshaledBallot,
					constructedBallot,
					cmp.AllowUnexported(voting.Ballot{}),
				) {
					t.Error("unexpected unmarshaled ballot:", unmarshaledBallot)
					t.Error("expected ballot:", constructedBallot)
				}
			} else { // User ID automatically generated cases
				userID := "dummy user ID"
				constructedBallot := voting.NewBallot(
					tt.pollID,
					userID,
					tt.rankOrder,
				)
				if !cmp.Equal(
					unmarshaledBallot,
					constructedBallot,
					cmp.AllowUnexported(voting.Ballot{}),
					cmpopts.IgnoreFields(voting.Ballot{}, "userID"),
				) {
					t.Error("unexpected unmarshaled ballot:", unmarshaledBallot)
				}
			}
		})
	}
}

func TestBallotMarshalUnmarshalDynamoDBAttributeValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		pollID, userID string
		rankOrder      []int
	}{
		{"poll1", "user1", []int{0, 1}},
		{"poll2", "user3", []int{0, 1, 3, 2}},
		{"poll5", "user5", []int{3, 2, 0, 1, 5, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.pollID, func(t *testing.T) {
			t.Parallel()
			inputBallot := voting.NewBallot(tt.pollID, tt.userID, tt.rankOrder)
			av, err := attributevalue.MarshalMap(inputBallot)
			if err != nil {
				t.Errorf("failed to marshal map: %v", err)
			}
			var b *voting.Ballot
			if err = attributevalue.UnmarshalMap(av, &b); err != nil {
				t.Errorf("failed to unmarshal map: %v", err)
			}
			if !cmp.Equal(inputBallot, b, cmp.AllowUnexported(voting.Ballot{})) {
				t.Errorf("unexpected unmarshaled result: %+v", b)
			}
		})
	}
}
