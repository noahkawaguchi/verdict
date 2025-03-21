//go:build test

package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func TestMain(m *testing.M) {
	// Make sure the tables don't already exist before testing
	deleteLocalTable(ballotsTableInfo)
	deleteLocalTable(pollsTableInfo)
	code := m.Run()
	os.Exit(code)
}

func TestBallotStore(t *testing.T) {
	createLocalTable(ballotsTableInfo, createBallotsTableInput)
	t.Cleanup(func() { deleteLocalTable(ballotsTableInfo) })

	tests := []struct {
		pollID    string
		userID    string
		rankOrder []int
	}{
		{"poll1", "user1", []int{1, 0, 2, 3}},
		{"poll1", "user2", []int{2, 1, 0, 3}},
		{"poll2", "user3", []int{3, 4, 2, 0, 1, 5}},
	}

	ts := &TableStore{}
	for _, test := range tests {
		inputBallot := models.NewBallot(test.pollID, test.userID, test.rankOrder)
		if err := ts.PutBallot(context.TODO(), inputBallot); err != nil {
			t.Error("unexpected error putting ballot:", err)
		}
		gotBallot, err := ts.getBallot(context.TODO(), test.pollID, test.userID)
		if err != nil {
			t.Error("unexpected error getting ballot:", err)
		}
		if !cmp.Equal(gotBallot, inputBallot, cmp.AllowUnexported(models.Ballot{})) {
			t.Error("got ballot did not match input ballot:", gotBallot)
		}
	}
}

func TestPollStore(t *testing.T) {
	createLocalTable(pollsTableInfo, createPollsTableInput)
	t.Cleanup(func() { deleteLocalTable(pollsTableInfo) })

	tests := []struct {
		prompt  string
		choices []string
	}{
		{"What is the best number?", []string{"1", "99", "100", "42"}},
		{"Why is the best number 42?", []string{"wabi-sabi", "randomness", "trick question"}},
		{"How big is 42?", []string{"at least 42", "at most 42", "7"}},
	}

	ts := &TableStore{}
	for _, test := range tests {
		inputPoll := models.NewPoll(test.prompt, test.choices)
		if err := ts.PutPoll(context.TODO(), inputPoll); err != nil {
			t.Error("unexpected error putting poll:", err)
		}
		gotPoll, err := ts.getPoll(context.TODO(), inputPoll.GetPollID())
		if err != nil {
			t.Error("unexpected error getting poll:", err)
		}
		if !cmp.Equal(gotPoll, inputPoll, cmp.AllowUnexported(models.Poll{})) {
			t.Error("got poll did not match input poll:", gotPoll)
		}
	}
}

func expectedPollJSON(prompt string, choices []string) string {
	choicesJSON, err := json.Marshal(choices)
	if err != nil {
		return ""
	}
	return fmt.Sprintf(`{"prompt":%q,"choices":%s}`, prompt, string(choicesJSON))
}

func TestGetPollData(t *testing.T) {
	createLocalTable(pollsTableInfo, createPollsTableInput)
	t.Cleanup(func() { deleteLocalTable(pollsTableInfo) })

	tests := []struct {
		prompt  string
		choices []string
	}{
		{"What is the best number?", []string{"1", "99", "100", "42"}},
		{"Why is the best number 42?", []string{"wabi-sabi", "randomness", "trick question"}},
		{"How big is 42?", []string{"at least 42", "at most 42", "7"}},
	}

	ts := &TableStore{}
	for _, test := range tests {
		inputPoll := models.NewPoll(test.prompt, test.choices)
		if err := ts.PutPoll(context.TODO(), inputPoll); err != nil {
			t.Error("unexpected error putting poll:", err)
		}
		gotPollData, err := ts.GetPollData(context.TODO(), inputPoll.GetPollID())
		if err != nil {
			t.Error("unexpected error getting poll data:", err)
		}
		if gotPollData != expectedPollJSON(test.prompt, test.choices) {
			t.Error("got poll data did not match input:", gotPollData)
		}
	}
}

func ballotSetEquality(ballotSlice1, ballotSlice2 []*models.Ballot) bool {
	if len(ballotSlice1) != len(ballotSlice2) {
		return false
	}
	// Convert ballots to their string representations so they can be used as map keys
	slice1Strings := make([]string, len(ballotSlice1))
	slice2Strings := make([]string, len(ballotSlice2))
	for i := range ballotSlice1 {
		slice1Strings[i] = fmt.Sprint(ballotSlice1[i])
		slice2Strings[i] = fmt.Sprint(ballotSlice2[i])
	}
	ballotCounts := make(map[string]bool)
	for _, ballotString := range slice1Strings {
		ballotCounts[ballotString] = true
	}
	for _, ballotString := range slice2Strings {
		if !ballotCounts[ballotString] {
			return false
		}
		ballotCounts[ballotString] = false
	}
	return true
}

func TestGetPollWithBallots(t *testing.T) {
	createLocalTable(ballotsTableInfo, createBallotsTableInput)
	createLocalTable(pollsTableInfo, createPollsTableInput)
	t.Cleanup(func() {
		deleteLocalTable(ballotsTableInfo)
		deleteLocalTable(pollsTableInfo)
	})

	tests := []struct {
		prompt  string
		choices []string
		ballots []struct {
			userID    string
			rankOrder []int
		}
	}{
		{
			"What is the best color?",
			[]string{"red", "blue", "green"},
			[]struct {
				userID    string
				rankOrder []int
			}{
				{"user10", []int{1, 0, 2}},
				{"user20", []int{0, 2, 1}},
			},
		},
		{
			"What is the best flavor of ice cream?",
			[]string{"vanilla", "chocolate", "strawberry", "mint"},
			[]struct {
				userID    string
				rankOrder []int
			}{
				{"user30", []int{3, 2, 0, 1}},
				{"user25", []int{2, 3, 1, 0}},
				{"user21", []int{0, 3, 2, 1}},
			},
		},
	}

	ts := &TableStore{}
	for _, test := range tests {
		inputPoll := models.NewPoll(test.prompt, test.choices)
		if err := ts.PutPoll(context.TODO(), inputPoll); err != nil {
			t.Error("unexpected error putting poll:", err)
		}
		inputBallots := make([]*models.Ballot, len(test.ballots))
		for i, b := range test.ballots {
			inputBallot := models.NewBallot(inputPoll.GetPollID(), b.userID, b.rankOrder)
			inputBallots[i] = inputBallot
			if err := ts.PutBallot(context.TODO(), inputBallot); err != nil {
				t.Error("unexpected error putting ballot:", err)
			}
		}
		gotPoll, gotBallots, err := ts.GetPollWithBallots(context.TODO(), inputPoll.GetPollID())
		if err != nil {
			t.Error("unexpected error getting poll with ballots:", err)
		}
		if !cmp.Equal(gotPoll, inputPoll, cmp.AllowUnexported(models.Poll{})) {
			t.Error("got poll did not match input:", gotPoll)
		}
		// Use sets because the ballots may have been retrieved in a different order
		if !ballotSetEquality(gotBallots, inputBallots) {
			t.Error("got ballots did not match input:", gotBallots)
		}
	}
}
