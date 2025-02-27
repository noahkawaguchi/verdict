package models

import (
	"fmt"
)

type Question struct {
	Prompt  string
	Choices []string
}

type Poll struct {
	PollID    string
	Questions []Question
}

func NewPoll(pollID string) *Poll {
	return &Poll{
		PollID: pollID,
		Questions: make([]Question, 0),
	}
}

func (p *Poll) AddQuestion(prompt string, choices []string) {
	p.Questions = append(p.Questions, Question{
		Prompt: prompt,
		Choices: choices,
	})
}

// TallyVotes creates a formatted string displaying the results of the poll.
func (p *Poll) TallyVotes(ballots ...*Ballot) string {
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
	// Create a readable string
	ret := "\nResults:\n"
	for i, q := range p.Questions {
		ret += q.Prompt + "\n"
		for j, c := range q.Choices {
			ret += fmt.Sprintf("  %v: %d\n", c, results[i][j])
		}
	}
	return ret
}
