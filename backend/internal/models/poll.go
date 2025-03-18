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
	return &Poll{
		pollID:  pollID,
		prompt:  prompt,
		choices: choices,
	}, pollID
}

func (p *Poll) String() string {
	shortID := p.pollID[:5] + "... "
	ret := fmt.Sprintf("Poll with ID %v:\n%v\n", shortID, p.prompt)
	for _, c := range p.choices {
		ret += fmt.Sprintf("  %v\n", c)
	}
	return ret
}

// ValidatedPollFromJSON unmarshals the provided JSON into a new poll. If the prompt or any of the
// choices are empty, there are fewer than two choices, or choices are not unique, it returns an
// error.
func ValidatedPollFromJSON(jsonString string) (*Poll, string, error) {
	// Create an auxiliary struct with exported fields to unmarshal the data
	aux := &struct {
		PollID  string   `json:"pollId"`
		Prompt  string   `json:"prompt"`
		Choices []string `json:"choices"`
	}{}
	if err := json.Unmarshal([]byte(jsonString), &aux); err != nil {
		return nil, "", errors.New("invalid JSON")
	}
	// Create a new poll ID here
	aux.PollID = uuid.New().String()
	// Validate the other fields
	if aux.Prompt == "" {
		return nil, "", errors.New("prompt cannot be empty")
	}
	if len(aux.Choices) < 2 {
		return nil, "", errors.New("there must be at least two choices")
	}
	if slices.Contains(aux.Choices, "") {
		return nil, "", errors.New("none of the choices can be empty")
	}
	// Use a "set" to validate uniqueness
	seen := make(map[string]struct{})
	for _, choice := range aux.Choices {
		if _, exists := seen[choice]; exists {
			return nil, "", errors.New("choices must be unique")
		}
		seen[choice] = struct{}{}
	}
	return &Poll{
		pollID:  aux.PollID,
		prompt:  aux.Prompt,
		choices: aux.Choices,
	}, aux.PollID, nil
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

// MarshalDynamoDBAttributeValue is a custom marshaler to control how the struct is serialized
// to DynamoDB.
func (p *Poll) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	m, err := attributevalue.MarshalMap(struct {
		PollID  string
		Prompt  string
		Choices []string
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
		PollID  string
		Prompt  string
		Choices []string
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
