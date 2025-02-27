package models

func dummyPoll(pollID string) *Poll {
	poll := NewPoll(pollID)
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
	return poll
}

func dummyBallot1(pollID, userID string) *Ballot {
	return NewBallot(pollID, userID, []int{1, 2, 1}) // lily, banana, yes
}

func dummyBallot2(pollID, userID string) *Ballot {
	return NewBallot(pollID, userID, []int{4, 0, 1}) // iris, apple, yes
}

func DummyData(pollID, userID1, userID2 string) (*Poll, *Ballot, *Ballot) {
	return dummyPoll(pollID), dummyBallot1(pollID, userID1), dummyBallot2(pollID, userID2)
}
