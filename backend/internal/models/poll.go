package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type Poll struct {
	pollID  string
	prompt  string
	choices []string
}

func NewPoll(prompt string, choices []string) (*Poll, string) {
	pollID := uuid.New().String()
	poll := &Poll{
		pollID:  pollID,
		prompt:  prompt,
		choices: choices,
	}
	return poll, pollID
}

func (p *Poll) GetPollID() string {
	return p.pollID
}

func (p *Poll) GetPrompt() string {
	return p.prompt
}

// ValidateFields ensures that all fields are non-empty and that there are at least two choices.
func (p *Poll) ValidateFields() error {
	if p.prompt == "" {
		return errors.New("prompt cannot be empty")
	}
	if len(p.choices) < 2 {
		return errors.New("there must be at least two choices")
	}
	if slices.Contains(p.choices, "") {
		return errors.New("none of the choices can be empty")
	}
	return nil
}

func (p *Poll) String() string {
	shortID := p.pollID[:5] + "... "
	ret := fmt.Sprintf("Poll with ID %v:\n%v\n", shortID, p.prompt)
	for _, c := range p.choices {
		ret += fmt.Sprintf("  %v\n", c)
	}
	return ret
}

// MarshalJSON is a custom marshaler that omits the poll ID.
func (p *Poll) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Prompt  string   `json:"prompt"`
		Choices []string `json:"choices"`
	}{
		Prompt:  p.prompt,
		Choices: p.choices,
	})
}

// UnmarshalJSON is custom unmarshaler that handles new poll creation when the data comes directly
// from JSON.
func (p *Poll) UnmarshalJSON(data []byte) error {
	// Create an auxiliary struct to unmarshal only the prompt and choices
	aux := &struct {
		Prompt  string   `json:"prompt"`
		Choices []string `json:"choices"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// Create a new poll ID here
	p.pollID = uuid.New().String()
	// Set the unmarshaled values back to the main struct
	p.prompt = aux.Prompt
	p.choices = aux.Choices
	return nil
}

// MarshalDynamoDBAttributeValue is a custom marshaler to control how the struct is serialized
// to DynamoDB.
func (p *Poll) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	m, err := attributevalue.MarshalMap(struct {
		PollID  string   `dynamodbav:"pollId"`
		Prompt  string   `dynamodbav:"prompt"`
		Choices []string `dynamodbav:"choices"`
	}{
		PollID:  p.pollID,
		Prompt:  p.prompt,
		Choices: p.choices,
	})
	if err != nil {
		return nil, err
	}
	return &types.AttributeValueMemberM{Value: m}, nil
}

// UnmarshalDynamoDBAttributeValue is a custom unmarshaler to control how the struct is
// deserialized from DynamoDB.
func (p *Poll) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	// Assert that av is of the correct type
	m, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return fmt.Errorf("expected *types.AttributeValueMemberM, got %T", av)
	}
	// Create a struct for custom unmarshaling
	var result struct {
		PollID  string   `dynamodbav:"pollId"`
		Prompt  string   `dynamodbav:"prompt"`
		Choices []string `dynamodbav:"choices"`
	}
	// Try to unmarshal using the custom struct
	if err := attributevalue.UnmarshalMap(m.Value, &result); err != nil {
		return err
	}
	// Set the unmarshaled values back to the main struct
	p.pollID = result.PollID
	p.prompt = result.Prompt
	p.choices = result.Choices
	return nil
}
