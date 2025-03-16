package datastore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

var pollsTableInfo = &tableInfo{
	name:         "Polls",
	partitionKey: "pollId",
	// (No sort key)
}

// PutPoll creates a new poll entry in the database.
func PutPoll(ctx context.Context, poll *models.Poll) error {
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

// GetPoll retrieves a poll from the database by its PollID.
func GetPoll(ctx context.Context, id string) (*models.Poll, error) {
	out, err := dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(pollsTableInfo.name),
		Key: map[string]types.AttributeValue{
			pollsTableInfo.partitionKey: &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, fmt.Errorf("poll with id %s not found in the database", id)
	}
	// Unmarshal the retrieved poll into a struct
	var poll models.Poll
	err = attributevalue.UnmarshalMap(out.Item, &poll)
	return &poll, err
}
