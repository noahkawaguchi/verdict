package datastore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

var pollsTableInfo = &tableInfo{
	name:         "Polls",
	partitionKey: "PollID",
	// (No sort key)
}

// PutPoll creates a new poll entry in the database.
func (ts *TableStore) PutPoll(ctx context.Context, poll *models.Poll) error {
	// Marshal the struct into a DynamoDB-compatible map
	av, err := attributevalue.MarshalMap(poll)
	if err != nil {
		return err
	}
	// Put the poll into DynamoDB
	_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(pollsTableInfo.name),
		Item:      av,
	})
	return err
}

// getPoll retrieves a poll from the database by its poll ID.
func (ts *TableStore) getPoll(ctx context.Context, pollID string) (*models.Poll, error) {
	out, err := dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(pollsTableInfo.name),
		Key: map[string]types.AttributeValue{
			pollsTableInfo.partitionKey: &types.AttributeValueMemberS{Value: pollID},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, fmt.Errorf("poll with id %s not found in the database", pollID)
	}
	// Unmarshal the retrieved poll into a struct
	var poll models.Poll
	err = attributevalue.UnmarshalMap(out.Item, &poll)
	return &poll, err
}

// GetPollData retrieves a poll's information from the database in JSON string format.
func (ts *TableStore) GetPollData(ctx context.Context, pollID string) (string, error) {
	poll, err := ts.getPoll(ctx, pollID)
	if err != nil {
		return "", err
	}
	body, err := json.Marshal(poll)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
