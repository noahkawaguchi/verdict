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

type result []resultQuestion

func NewResult(numQuestions int) result {
	return result(make([]resultQuestion, numQuestions))
}

func (r result) String() string {
	ret := "\nResults:\n"
	for _, q := range r {
		ret += q.prompt + "\n"
		for _, c := range q.choices {
			ret += fmt.Sprintf("  %v: %d\n", c.choice, c.votes)
		}
	}
	return ret
}
