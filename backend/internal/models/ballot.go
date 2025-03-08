package models

import "fmt"

type Ballot struct {
	PollID string
	UserID string
	// The indices of the voter's choices. For example, if the voter's first choice is at index 2
	// in the poll's choices, then RankOrder[0] = 2.
	RankOrder []int
}

func NewBallot(pollID, userID string, rankOrder []int) *Ballot {
	return &Ballot{
		PollID:    pollID,
		UserID:    userID,
		RankOrder: rankOrder,
	}
}

func (b *Ballot) String() string {
	shortPollID := b.PollID[:5] + "... "
	return fmt.Sprintf("Ballot from user %v for poll %v with choices %v",
		b.UserID, shortPollID, b.RankOrder)
}
