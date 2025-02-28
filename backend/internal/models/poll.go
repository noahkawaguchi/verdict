package models

import "github.com/google/uuid"

type PollQuestion struct {
	Prompt  string   `json:"prompt"`
	Choices []string `json:"choices"`
}

type Poll struct {
	PollID    string
	Title     string
	Questions []PollQuestion
}

func NewPoll(title string, questions []PollQuestion) (*Poll, string) {
	pollID := uuid.New().String()
	poll := &Poll{
		PollID:    pollID,
		Title:     title,
		Questions: questions,
	}
	return poll, pollID
}

// TallyVotes accumulates the results of the poll.
func (p *Poll) TallyVotes(ballots ...*Ballot) *result {
	// Create a 2D slice of 0s
	results := make([][]int, len(p.Questions))
	for i := range results {
		results[i] = make([]int, len(p.Questions[i].Choices))
	}
	// Accumulate the results
	for _, b := range ballots {
		if b.PollID == p.PollID { // Make sure the ballot was for this poll
			for i, choice := range b.Selections {
				results[i][choice]++
			}
		}
	}
	// Convert the results to a result struct
	pollResult := newResult(p.Title, len(ballots), len(p.Questions))
	for i, q := range p.Questions {
		pollResult.results[i].prompt = q.Prompt
		pollResult.results[i].choices = make([]choiceStats, len(q.Choices))
		for j, c := range q.Choices {
			pollResult.results[i].choices[j] = pollResult.newChoiceStats(c, results[i][j])
		}
	}
	pollResult.sortChoices()
	return pollResult
}
