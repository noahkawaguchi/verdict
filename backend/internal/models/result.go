package models

import "fmt"

type choiceStats struct {
	choice     string
	votes      int
	percentage int
}

type resultQuestion struct {
	prompt  string
	choices []choiceStats
}

type result struct {
	title      string
	numBallots int
	results    []resultQuestion
}

func newResult(title string, numBallots, numQuestions int) *result {
	return &result{
		title:      title,
		numBallots: numBallots,
		results:    make([]resultQuestion, numQuestions),
	}
}

func (r *result) newChoiceStats(choice string, votes int) choiceStats {
	return choiceStats{
		choice: choice,
		votes: votes,
		percentage: votes * 100 / r.numBallots,
	}
}

func (r *result) String() string {
	ret := fmt.Sprintf("\nResults for %q:\n", r.title)
	for _, q := range r.results {
		ret += fmt.Sprintf("  %v\n", q.prompt)
		for _, c := range q.choices {
			ret += fmt.Sprintf("    %d vote(s) (%d%%): %v\n", c.votes, c.percentage, c.choice)
		}
	}
	return ret
}
