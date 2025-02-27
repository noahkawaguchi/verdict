package models

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
