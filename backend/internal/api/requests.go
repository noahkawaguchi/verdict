package api

import "slices"

type createPollRequest struct {
	Prompt  string   `json:"prompt"`
	Choices []string `json:"choices"`
}

// allValidFields validates that the prompt and choices are non-empty and that there are at least
// two choices.
func (cpr *createPollRequest) allValidFields() bool {
	if cpr.Prompt == "" { // Empty prompt
		return false
	}
	if len(cpr.Choices) < 2 { // Too few choices
		return false
	}
	if slices.Contains(cpr.Choices, "") { // Empty choice
		return false
	}
	return true
}
