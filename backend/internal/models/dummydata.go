package models

import "math/rand"

func dummyPoll() (*Poll, string) {
	poll, pollID := NewPoll(
		"What is the best fruit?",
		[]string{"apple", "banana", "clementine", "durian"},
	)
	return poll, pollID
}

func dummyBallot(pollID, userID string) *Ballot {
	ranks := []int{0, 1, 2, 3}
	rand.Shuffle(len(ranks), func(i, j int) {
		ranks[i], ranks[j] = ranks[j], ranks[i]
	})
	return NewBallot(pollID, userID, ranks)
}

func DummyData(userIDs []string) (*Poll, string, []*Ballot) {
	poll, pollID := dummyPoll()
	ballots := make([]*Ballot, len(userIDs))
	for i, userID := range userIDs {
		ballots[i] = dummyBallot(pollID, userID)
	}
	return poll, pollID, ballots
}
