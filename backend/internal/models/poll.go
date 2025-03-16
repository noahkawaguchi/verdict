package models

import (
	"fmt"

	"github.com/google/uuid"
)

type Poll struct {
	PollID  string
	Prompt  string
	Choices []string
}

func NewPoll(prompt string, choices []string) (*Poll, string) {
	pollID := uuid.New().String()
	poll := &Poll{
		PollID:  pollID,
		Prompt:  prompt,
		Choices: choices,
	}
	return poll, pollID
}

func (p *Poll) String() string {
	shortID := p.PollID[:5] + "... "
	ret := fmt.Sprintf("Poll with ID %v:\n%v\n", shortID, p.Prompt)
	for _, c := range p.Choices {
		ret += fmt.Sprintf("  %v\n", c)
	}
	return ret
}
