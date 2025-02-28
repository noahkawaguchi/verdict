package models

import "fmt"

type choiceStats struct {
	choice string
	votes  int
}

type resultQuestion struct {
	prompt  string
	choices []choiceStats
}

type result struct {
	title   string
	results []resultQuestion
}

func NewResult(title string, numQuestions int) *result {
	return &result{
		title: title,
		results: make([]resultQuestion, numQuestions),
	}
}

func (r *result) String() string {
	ret := fmt.Sprintf("\nResults for %q:\n", r.title)
	for _, q := range r.results {
		ret += fmt.Sprintf("  %v\n", q.prompt)
		for _, c := range q.choices {
			ret += fmt.Sprintf("    %v: %d\n", c.choice, c.votes)
		}
	}
	return ret
}
