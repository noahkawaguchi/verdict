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

type Ballot struct {
	pollID string
	userID string
	// The indices of the voter's choices. For example, if the voter's first choice is at index 2
	// in the poll's choices, then rankOrder[0] = 2.
	rankOrder []int
}

// NewValidatedBallot creates a new ballot. If the poll ID or user ID are empty, if the rank order 
// has fewer than two rankings, or if the rank order is not a permutation of its indices, it 
// returns an error.
func NewValidatedBallot(pollID, userID string, rankOrder []int) (*Ballot, error) {
	ballot := &Ballot{
		pollID:    pollID,
		userID:    userID,
		rankOrder: rankOrder,
	}
	if err := ballot.validate(); err != nil {
		return nil, err
	}
	return ballot, nil
}

// ValidatedBallotFromJSON unmarshals the provided JSON into a new ballot. If no user ID is
// provided, a new one is generated. If the poll ID or rank order are missing, if rank order has
// fewer than two rankings, or if the rank order is not a permutation of its indices, it returns an
// error.
func ValidatedBallotFromJSON(jsonString string) (*Ballot, error) {
	// Create an auxiliary struct with exported fields to unmarshal the data
	aux := &struct {
		PollID    string `json:"pollId"`
		UserID    string `json:"userId"`
		RankOrder []int  `json:"rankOrder"`
	}{}
	if err := json.Unmarshal([]byte(jsonString), &aux); err != nil {
		return nil, errors.New("invalid JSON")
	}
	// Create a new user ID if it's not provided
	if aux.UserID == "" {
		aux.UserID = uuid.New().String()
	}
	ballot := &Ballot{
		pollID:    aux.PollID,
		userID:    aux.UserID,
		rankOrder: aux.RankOrder,
	}
	// Validate the fields
	if err := ballot.validate(); err != nil {
		return nil, err
	}
	return ballot, nil
}

func (b *Ballot) String() string {
	shortPollID := b.pollID[:5] + "... "
	return fmt.Sprintf("Ballot from user %s for poll %s with choices %v",
		b.userID, shortPollID, b.rankOrder)
}

// validate ensures that none of the fields are empty, there are at least two rankings, and the
// rank order is a permutation of its indices.
func (b *Ballot) validate() error {
	if b.pollID == "" {
		return errors.New("poll ID cannot be empty")
	}
	if b.userID == "" {
		return errors.New("user ID cannot be empty")
	}
	if len(b.rankOrder) < 2 {
		return errors.New("there must be at least two rankings")
	}
	// Copy the slice to avoid changing the original underlying array
	sortCopy := slices.Clone(b.rankOrder)
	slices.Sort(sortCopy)
	for i, v := range sortCopy {
		if v != i {
			return errors.New("not a valid rank order")
		}
	}
	return nil
}

// MarshalDynamoDBAttributeValue is a custom marshaler to control how the struct is serialized
// to DynamoDB.
func (b *Ballot) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	m, err := attributevalue.MarshalMap(struct {
		PollID    string
		UserID    string
		RankOrder []int
	}{
		PollID:    b.pollID,
		UserID:    b.userID,
		RankOrder: b.rankOrder,
	})
	if err != nil {
		return nil, err
	}
	return &types.AttributeValueMemberM{Value: m}, nil
}

// UnmarshalDynamoDBAttributeValue is a custom unmarshaler to control how the struct is
// deserialized from DynamoDB.
func (b *Ballot) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	// Assert that av is of the correct type
	m, ok := av.(*types.AttributeValueMemberM)
	if !ok {
		return fmt.Errorf("expected *types.AttributeValueMemberM, got %T", av)
	}
	// Create a struct for custom unmarshaling
	var result struct {
		PollID    string
		UserID    string
		RankOrder []int
	}
	// Try to unmarshal using the custom struct
	if err := attributevalue.UnmarshalMap(m.Value, &result); err != nil {
		return err
	}
	// Set the unmarshaled values back to the main struct
	b.pollID = result.PollID
	b.userID = result.UserID
	b.rankOrder = result.RankOrder
	return nil
}
