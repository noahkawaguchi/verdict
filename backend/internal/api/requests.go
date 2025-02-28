package api

import (
	"slices"

	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

type createPollRequest struct {
	Title     string                `json:"title"`
	Questions []models.PollQuestion `json:"questions"`
}

// anyEmptyFields checks if the title, questions, or any of the questions'
// prompts or choices are empty.
func (cpr *createPollRequest) anyEmptyFields() bool {
	if cpr.Title == "" { // Empty title
		return true
	}
	if len(cpr.Questions) == 0 { // No questions
		return true
	}
	for _, q := range cpr.Questions {
		if q.Prompt == "" { // Empty prompt
			return true
		}
		if len(q.Choices) == 0 { // No choices
			return true
		}
		if slices.Contains(q.Choices, "") { // Empty choice
			return true
		}
	}
	return false
}
