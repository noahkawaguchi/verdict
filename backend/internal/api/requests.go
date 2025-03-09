package api

import (
	"errors"
	"slices"
)

type createPollRequest struct {
	Prompt  string   `json:"prompt"`
	Choices []string `json:"choices"`
}

// validateFields validates that the prompt and choices are non-empty and that there are at least
// two choices.
func (cpr *createPollRequest) validateFields() error {
	if cpr.Prompt == "" {
		return errors.New("prompt cannot be empty")
	}
	if len(cpr.Choices) < 2 {
		return errors.New("there must be at least two choices")
	}
	if slices.Contains(cpr.Choices, "") {
		return errors.New("none of the choices can be empty")
	}
	return nil
}
