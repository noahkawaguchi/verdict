package models

func dummyPoll() (*Poll, string) {
	poll, pollID := NewPoll("Flowers, Fruits, and DST")
	poll.AddQuestion(
		"What is the best flower?",
		[]string{"rose", "lily", "tulip", "carnation", "iris"},
	)
	poll.AddQuestion(
		"What is the best fruit?",
		[]string{"apple", "orange", "banana"},
	)
	poll.AddQuestion(
		"Should the US end the Daylight Savings Time system?",
		[]string{"no", "yes"},
	)
	return poll, pollID
}

func dummyBallot1(pollID, userID string) *Ballot {
	return NewBallot(pollID, userID, []int{1, 2, 1}) // lily, banana, yes
}

func dummyBallot2(pollID, userID string) *Ballot {
	return NewBallot(pollID, userID, []int{4, 0, 1}) // iris, apple, yes
}

func dummyBallot3(pollID, userID string) *Ballot {
	return NewBallot(pollID, userID, []int{0, 2, 1}) // rose, banana, yes
}

func DummyData(userID1, userID2, userID3 string) (*Poll, string, []*Ballot) {
	poll, pollID := dummyPoll()
	ballots := []*Ballot{dummyBallot1(pollID, userID1),
		dummyBallot2(pollID, userID2), dummyBallot3(pollID, userID3)}
	return poll, pollID, ballots
}
