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
	pollID, userID string
	// The indices of the voter's choices. For example, if the voter's first choice is at index 2
	// in the poll's choices, then rankOrder[0] = 2.
	rankOrder []int
}

func NewBallot(pollID, userID string, rankOrder []int) *Ballot {
	return &Ballot{pollID, userID, rankOrder}
}

// Validate ensures that none of the fields are empty, there are at least two rankings, and the
// rank order is a permutation of its indices.
func (b *Ballot) Validate() error {
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

func (b *Ballot) String() string {
	return fmt.Sprintf("Ballot from user %s for poll %s with choices %v",
		b.userID, b.pollID[:5]+"... ", b.rankOrder)
}

// UnmarshalJSON is a custom JSON unmarshaler. If no user ID is provided, a new one is generated.
func (b *Ballot) UnmarshalJSON(data []byte) error {
	// Create an auxiliary struct with exported fields to unmarshal the data
	var aux struct {
		PollID    string `json:"pollId"`
		UserID    string `json:"userId"`
		RankOrder []int  `json:"rankOrder"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	// Create a new user ID if it's not provided
	if aux.UserID == "" {
		b.userID = uuid.New().String()
	} else {
		b.userID = aux.UserID
	}
	// Set the other unmarshaled values back to the main struct
	b.pollID, b.rankOrder = aux.PollID, aux.RankOrder
	return nil
}

// MarshalDynamoDBAttributeValue is a custom marshaler to control how the struct is serialized
// to DynamoDB.
func (b *Ballot) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	m, err := attributevalue.MarshalMap(struct {
		PollID, UserID string
		RankOrder      []int
	}{b.pollID, b.userID, b.rankOrder})
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
	var aux struct {
		PollID, UserID string
		RankOrder      []int
	}
	// Try to unmarshal using the custom struct
	if err := attributevalue.UnmarshalMap(m.Value, &aux); err != nil {
		return err
	}
	// Set the unmarshaled values back to the main struct
	b.pollID, b.userID, b.rankOrder = aux.PollID, aux.UserID, aux.RankOrder
	return nil
}
