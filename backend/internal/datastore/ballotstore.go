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

var ballotsTableInfo = &tableInfo{
	name:         "Ballots",
	partitionKey: "PollID",
	sortKey:      "UserID",
}

// PutBallot creates a new ballot entry in the database.
func PutBallot(ctx context.Context, ballot *models.Ballot) error {
	// Marshal the struct into a DynamoDB-compatible map
	av, err := attributevalue.MarshalMap(ballot)
	if err != nil {
		return err
	}
	// Put the ballot into DynamoDB
	_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(ballotsTableInfo.name),
		Item:      av,
	})
	return err
}

// getBallot retrieves a ballot from the database by its poll ID and user ID.
func getBallot(ctx context.Context, pollID, userID string) (*models.Ballot, error) {
	out, err := dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(ballotsTableInfo.name),
		Key: map[string]types.AttributeValue{
			ballotsTableInfo.partitionKey: &types.AttributeValueMemberS{Value: pollID},
			ballotsTableInfo.sortKey:      &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, err
	}
	if out.Item == nil {
		return nil, fmt.Errorf("ballot with poll ID %s and user ID %s not found in the database",
			pollID, userID)
	}
	// Unmarshal the retrieved ballot into a struct
	var ballot models.Ballot
	err = attributevalue.UnmarshalMap(out.Item, &ballot)
	return &ballot, err
}

// GetPollWithBallots retrieves a poll and all of its ballots from the database.
func GetPollWithBallots(ctx context.Context, pollID string) (*models.Poll, []*models.Ballot, error) {
	// Get the poll
	poll, err := getPoll(ctx, pollID)
	if err != nil {
		return nil, nil, err
	}
	// Define input to query by pollID
	input := &dynamodb.QueryInput{
		TableName:              aws.String(ballotsTableInfo.name),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :pk", pollsTableInfo.partitionKey)),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pollID},
		},
	}
	// Perform the query
	out, err := dbClient.Query(ctx, input)
	if err != nil {
		return nil, nil, err
	}
	// Unmarshal the retrieved items
	var ballots []*models.Ballot
	err = attributevalue.UnmarshalListOfMaps(out.Items, &ballots)
	return poll, ballots, err
}
