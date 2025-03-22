package datastore

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

// PutPoll creates a new poll entry in the database.
func (ts *TableStore) PutPoll(poll *models.Poll) error {
	// Marshal the struct into a DynamoDB-compatible map
	av, err := attributevalue.MarshalMap(poll)
	if err != nil {
		return err
	}
	// Put the poll into DynamoDB
	_, err = ts.Client.PutItem(ts.Ctx, &dynamodb.PutItemInput{
		TableName: aws.String(pollsTableInfo.name),
		Item:      av,
	})
	return err
}

// GetPoll retrieves a poll from the database by its poll ID.
func (ts *TableStore) GetPoll(pollID string) (*models.Poll, error) {
	// Define the input to get the poll by ID
	input := &dynamodb.GetItemInput{
		TableName: aws.String(pollsTableInfo.name),
		Key: map[string]types.AttributeValue{
			pollsTableInfo.partitionKey: &types.AttributeValueMemberS{Value: pollID},
		},
	}
	out, err := ts.Client.GetItem(ts.Ctx, input)
	if err != nil {
		return nil, err
	}
	// Unmarshal the retrieved poll into a struct
	var poll models.Poll
	err = attributevalue.UnmarshalMap(out.Item, &poll)
	return &poll, err
}
