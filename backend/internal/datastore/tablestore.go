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

var dbClient *dynamodb.Client

type tableInfo struct {
	name, partitionKey, sortKey string
}

type TableStore struct {}

// GetPollWithBallots retrieves a poll and all of its ballots from the database.
func (ts *TableStore) GetPollWithBallots(ctx context.Context, pollID string) (*models.Poll, []*models.Ballot, error) {
	// Get the poll
	poll, err := ts.getPoll(ctx, pollID)
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
