package models

type Ballot struct {
	PollID     string
	UserID     string
	Selections []int // The indices of the voter's choices
}

func NewBallot(pollID, userID string, selections []int) *Ballot {
	return &Ballot{
		PollID: pollID,
		UserID: userID,
		Selections: selections,
	}
}
