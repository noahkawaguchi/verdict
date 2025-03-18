package models

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
)

func TestBallotValidateFields_EmptyFields(t *testing.T) {
	tests := []struct {
		errMsg    string
		pollID    string
		userID    string
		rankOrder []int
	}{
		{"poll ID cannot be empty", "", "user1", []int{0, 1, 2}},
		{"poll ID cannot be empty", "", "", []int{0, 1, 2}},
		{"poll ID cannot be empty", "", "", []int{}},
		{"user ID cannot be empty", "poll1", "", []int{0, 1, 2}},
		{"user ID cannot be empty", "poll1", "", []int{}},
		{"there must be at least two rankings", "poll1", "user1", []int{}},
		{"there must be at least two rankings", "poll1", "user1", []int{0}},
	}
	for _, test := range tests {
		ballot := NewBallot(test.pollID, test.userID, test.rankOrder)
		if err := ballot.ValidateFields(); err == nil || err.Error() != test.errMsg {
			t.Errorf("expected error with message %q, got %v", test.errMsg, err)
		}
	}
}

func TestBallotValidateFields_InvalidRankOrder(t *testing.T) {
	tests := [][]int{
		{0, 0, 1, 2},
		{1, 2, 3, 4},
		{4, 3, 2, 0},
		{00, 11, 22, 33},
		{3, 2, 3, 1},
		{0, 10, 20, 30},
	}
	for _, test := range tests {
		ballot := NewBallot("poll1", "user1", test)
		if err := ballot.ValidateFields(); err == nil || err.Error() != "not a valid rank order" {
			t.Errorf("expected error with message \"not a valid rank order,\", got %v", err)
		}
	}
}

func TestBallotValidateFields_Valid(t *testing.T) {
	tests := []Ballot{
		{"poll1", "user1", []int{0, 2, 1, 3}},
		{"poll1", "user2", []int{1, 3, 2, 0}},
		{"poll2", "user3", []int{1, 0, 2}},
		{"poll2", "user4", []int{0, 1, 2}},
		{"poll3", "user5", []int{0, 1, 2, 4, 3, 5}},
		{"poll3", "user6", []int{3, 4, 2, 1, 0, 5}},
	}
	for _, test := range tests {
		ballot := NewBallot(test.pollID, test.userID, test.rankOrder)
		if err := ballot.ValidateFields(); err != nil {
			t.Errorf("expected success, got %v", err)
		}
	}
}

func TestBallotUnmarshalJSON_AllThreeFieldsProvided(t *testing.T) {
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
	}
	for _, test := range tests {
		var b Ballot
		if err := json.Unmarshal([]byte(test.jsonString), &b); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if b.pollID != test.pollID || b.userID != test.userID ||
			!reflect.DeepEqual(b.rankOrder, test.rankOrder) {
			t.Errorf("unexpected unmarshaled result: %+v", b)
		}
	}
}

func TestBallotUnmarshalJSON_AutomaticUserID(t *testing.T) {
	tests := []struct {
		pollID     string
		rankOrder  []int
		jsonString string
	}{
		{"poll1", []int{0, 1, 2}, `{"pollId": "poll1", "rankOrder": [0, 1, 2]}`},
		{"poll2", []int{3, 0, 1, 2, 4}, `{"pollId": "poll2", "rankOrder": [3, 0, 1, 2, 4]}`},
	}
	for _, test := range tests {
		var b Ballot
		if err := json.Unmarshal([]byte(test.jsonString), &b); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if b.pollID != test.pollID || !reflect.DeepEqual(b.rankOrder, test.rankOrder) {
			t.Errorf("unexpected unmarshaled result: %+v", b)
		}
		if b.userID == "" {
			t.Error("failed to automatically generate user ID")
		}
	}
}

func TestBallotMarshalUnmarshalDynamoDBAttributeValue(t *testing.T) {
	tests := []*Ballot{
		{"poll1", "user1", []int{0, 1}},
		{"poll2", "user3", []int{0, 1, 3, 2}},
		{"poll5", "user5", []int{3, 2, 0, 1, 5, 4}},
	}
	for _, test := range tests {
		av, err := attributevalue.MarshalMap(test)
		if err != nil {
			t.Errorf("failed to marshal map: %v", err)
		}
		var b *Ballot
		if err = attributevalue.UnmarshalMap(av, &b); err != nil {
			t.Errorf("failed to unmarshal map: %v", err)
		}
		if b.pollID != test.pollID || b.userID != test.userID ||
			!reflect.DeepEqual(b.rankOrder, test.rankOrder) {
			t.Errorf("unexpected unmarshaled result: %+v", b)
		}
	}
}
