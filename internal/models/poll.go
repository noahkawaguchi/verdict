package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/noahkawaguchi/verdict/backend/internal/utils"
)

type Poll struct {
	pollID, prompt string
	choices        []string
}

// NewPoll creates a new poll with a newly generated poll ID.
func NewPoll(prompt string, choices []string) *Poll {
	return &Poll{uuid.New().String(), prompt, choices}
}

// ID gets the poll's poll ID.
func (p *Poll) ID() string { return p.pollID }

// Validate ensures that the prompt and all choices are non-empty, that there are at least two
// choices, and that all choices are unique.
func (p *Poll) Validate() error {
	if p.prompt == "" {
		return errors.New("prompt cannot be empty")
	}
	if len(p.choices) < 2 {
		return errors.New("there must be at least two choices")
	}
	if slices.Contains(p.choices, "") {
		return errors.New("none of the choices can be empty")
	}
	if len(p.choices) != utils.NewSet(p.choices...).Len() {
		return errors.New("choices must be unique")
	}
	return nil
}

func (p *Poll) String() string {
	ret := fmt.Sprintf("Poll with ID %s:\n%s\n", p.pollID[:5]+"... ", p.prompt)
	for _, c := range p.choices {
		ret += fmt.Sprintf("  %s\n", c)
	}
	return ret
}

// MarshalJSON is a custom marshaler that omits the poll ID.
func (p *Poll) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Prompt  string   `json:"prompt"`
		Choices []string `json:"choices"`
	}{p.prompt, p.choices})
}

// UnmarshalJSON is a custom JSON unmarshaler. It generates a poll ID for the new poll.
func (p *Poll) UnmarshalJSON(data []byte) error {
	// Create an auxiliary struct with exported fields to unmarshal the data
	var aux struct {
		Prompt  string   `json:"prompt"`
		Choices []string `json:"choices"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// Create a new poll ID
	p.pollID = uuid.New().String()
	// Set the other unmarshaled values back to the main struct
	p.prompt, p.choices = aux.Prompt, aux.Choices
	return nil
}

// MarshalDynamoDBAttributeValue is a custom marshaler to control how the struct is serialized
// to DynamoDB.
func (p *Poll) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	m, err := attributevalue.MarshalMap(struct {
		PollID, Prompt string
		Choices        []string
	}{p.pollID, p.prompt, p.choices})
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
	var aux struct {
		PollID, Prompt string
		Choices        []string
	}
	// Try to unmarshal using the custom struct
	if err := attributevalue.UnmarshalMap(m.Value, &aux); err != nil {
		return err
	}
	// Set the unmarshaled values back to the main struct
	p.pollID, p.prompt, p.choices = aux.PollID, aux.Prompt, aux.Choices
	return nil
}
