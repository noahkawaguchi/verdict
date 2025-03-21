package models_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func threeOptionPoll() (*models.Poll, string) {
	poll := models.NewPoll("What is the best fruit?",
		[]string{"apple", "banana", "clementine"})
	return poll, poll.GetPollID()
}

func fourOptionPoll() (*models.Poll, string) {
	poll := models.NewPoll("What is the best fruit?",
		[]string{"apple", "banana", "clementine", "durian"})
	return poll, poll.GetPollID()
}

func ballotClosure(pollID string) func([]int) *models.Ballot {
	userID := 0
	return func(rankOrder []int) *models.Ballot {
		userID++
		ballot := models.NewBallot(pollID, strconv.Itoa(userID), rankOrder)
		return ballot
	}
}

func expectedResultJSON(totalVotes, winningVotes, winningChoiceIdx, winningRound int) string {
	winningChoice := []string{"apple", "banana", "clementine", "durian"}[winningChoiceIdx]
	return fmt.Sprintf(
		`{"prompt":"What is the best fruit?","totalVotes":%d,"winningVotes":%d,`+
			`"winningChoice":%q,"winningRound":%d}`,
		totalVotes, winningVotes, winningChoice, winningRound,
	)
}

func TestResult_SimpleMajority(t *testing.T) {
	poll, pollID := threeOptionPoll()
	ballotWithRanks := ballotClosure(pollID)
	ballots := []*models.Ballot{
		ballotWithRanks([]int{1, 0, 2}),
		ballotWithRanks([]int{2, 0, 1}),
		ballotWithRanks([]int{2, 1, 0}),
	}
	body, err := models.CalculateResultData(poll, ballots)
	if err != nil {
		t.Error(err.Error())
	}
	if body != expectedResultJSON(3, 2, 2, 1) {
		t.Errorf("unexpected result: %s", body)
	}
}

func TestResult_Runoff(t *testing.T) {
	poll, pollID := threeOptionPoll()
	ballotWithRanks := ballotClosure(pollID)
	ballots := []*models.Ballot{
		ballotWithRanks([]int{0, 1, 2}),
		ballotWithRanks([]int{1, 0, 2}),
		ballotWithRanks([]int{1, 0, 2}),
		ballotWithRanks([]int{2, 0, 1}),
		ballotWithRanks([]int{2, 1, 0}),
	}
	body, err := models.CalculateResultData(poll, ballots)
	if err != nil {
		t.Error(err.Error())
	}
	if body != expectedResultJSON(5, 3, 1, 2) {
		t.Errorf("unexpected result: %s", body)
	}
}

func TestResult_TieForLast(t *testing.T) {
	poll, pollID := fourOptionPoll()
	ballotWithRanks := ballotClosure(pollID)
	ballots := []*models.Ballot{
		ballotWithRanks([]int{0, 2, 3, 1}),
		ballotWithRanks([]int{1, 3, 0, 2}),
		ballotWithRanks([]int{1, 3, 0, 2}),
		ballotWithRanks([]int{2, 0, 1, 3}),
		ballotWithRanks([]int{2, 3, 0, 1}),
		ballotWithRanks([]int{3, 2, 0, 1}),
	}
	body, err := models.CalculateResultData(poll, ballots)
	if err != nil {
		t.Error(err.Error())
	}
	/*
		Round 1:
			- 0 and 3 tie for last.
			- The tie-breaking algorithm chooses 0 for elimination because it has fewer
			  second-choice votes amongst the ballots counting toward choices not under
			  consideration for elimination.
			- The ballot with 0 as the first choice is redistributed to its second choice, 2.
		Round 2:
			- 2 has exactly 3/6 votes, not a strict majority, so the algorithm continues.
			- 3 is now in last place, so the ballot with 3 as the first choice is redistributed to
			  its second choice, 2.
		Round 3:
			- 2 now has 4/6 votes, a strict majority, and wins.
	*/
	if body != expectedResultJSON(6, 4, 2, 3) {
		t.Errorf("unexpected result: %s", body)
	}
}

func TestResult_InfiniteTieForLast(t *testing.T) {
	poll, pollID := fourOptionPoll()
	ballotWithRanks := ballotClosure(pollID)
	ballots := []*models.Ballot{
		ballotWithRanks([]int{0, 2, 3, 1}),
		ballotWithRanks([]int{0, 3, 1, 2}),
		ballotWithRanks([]int{1, 3, 0, 2}),
		ballotWithRanks([]int{2, 0, 1, 3}),
		ballotWithRanks([]int{2, 3, 0, 1}),
		ballotWithRanks([]int{3, 0, 1, 2}),
		ballotWithRanks([]int{3, 1, 2, 0}),
		ballotWithRanks([]int{3, 2, 0, 1}),
	}
	body, err := models.CalculateResultData(poll, ballots)
	if err != nil {
		t.Error(err.Error())
	}
	/*
		Round 1:
			- Ballot state: [0, 0, 1, 2, 2, 3, 3, 3]
			- 1 is in last place and is eliminated.
			- The ballot with 1 as its first choice is redistributed to its second choice, 3.
		Round 2:
			- Ballot state: [0, 0, 3, 2, 2, 3, 3, 3]
			- 0 and 2 tie for last place, each with two votes.
			- In the tie-breaking algorithm, 0 and 2 each receive four votes. In other words, in a
			  theoretical poll with the same voters and only these two choices, the voters would be
			  split exactly.
			- The tie-breaking algorithm randomly chooses 0 or 2. (Ideally, this case is rare.)
			- The ballots counting for the eliminated choice are redistributed to their next
			  highest choice.
		Round 3:
			- Ballot state if 0 was eliminated: [2, 3, 3, 2, 2, 3, 3, 3]
			- Ballot state if 2 was eliminated: [0, 0, 3, 0, 3, 3, 3, 3]
			- In either case, 3 now has 5 out of 8 votes, a strict majority, and wins.
	*/
	if body != expectedResultJSON(8, 5, 3, 3) {
		t.Errorf("unexpected result: %s", body)
	}
}
