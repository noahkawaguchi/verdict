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

// validateFields validates that the prompt and choices are non-empty and that there are at least
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
