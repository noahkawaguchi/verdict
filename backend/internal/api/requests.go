package api

import (
	"errors"
	"slices"
)

// pollRequestResponse serves as the body for both POST requests and GET responses.
type pollRequestResponse struct {
	Prompt  string   `json:"prompt"`
	Choices []string `json:"choices"`
}

// validateFields ensures that the prompt and choices are non-empty and that there are at least
// two choices.
func (prr *pollRequestResponse) validateFields() error {
	if prr.Prompt == "" {
		return errors.New("prompt cannot be empty")
	}
	if len(prr.Choices) < 2 {
		return errors.New("there must be at least two choices")
	}
	if slices.Contains(prr.Choices, "") {
		return errors.New("none of the choices can be empty")
	}
	return nil
}

type ballotRequest struct {
	PollID    string `json:"pollId"`
	RankOrder []int  `json:"rankOrder"`
}

// validateFields ensures that the poll ID is non-empty, the rank order has at least two rankings,
// and the rank order is a permutation of its indices.
func (br *ballotRequest) validateFields() error {
	if br.PollID == "" {
		return errors.New("missing poll ID")
	}
	if len(br.RankOrder) < 2 {
		return errors.New("there must be at least two rankings")
	}
	// Copy the slice to avoid changing the original underlying array
	sortCopy := slices.Clone(br.RankOrder)
	slices.Sort(sortCopy)
	for i, v := range sortCopy {
		if v != i {
			return errors.New("not a valid rank order")
		}
	}
	return nil
}
