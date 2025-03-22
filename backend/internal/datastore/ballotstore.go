package datastore

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

// PutBallot creates a new ballot entry in the database.
func (ts *TableStore) PutBallot(ballot *models.Ballot) error {
	// Marshal the struct into a DynamoDB-compatible map
	av, err := attributevalue.MarshalMap(ballot)
	if err != nil {
		return err
	}
	// Put the ballot into DynamoDB
	_, err = ts.Client.PutItem(ts.Ctx, &dynamodb.PutItemInput{
		TableName: aws.String(ballotsTableInfo.name),
		Item:      av,
	})
	return err
}

// GetBallots retrieves all of the ballots for the specified poll from the database.
func (ts *TableStore) GetBallots(pollID string) ([]*models.Ballot, error) {
	// Define input to query by pollID
	input := &dynamodb.QueryInput{
		TableName:              aws.String(ballotsTableInfo.name),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = :pk", pollsTableInfo.partitionKey)),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pollID},
		},
	}
	// Perform the query
	out, err := ts.Client.Query(ts.Ctx, input)
	if err != nil {
		return nil, err
	}
	// Unmarshal the retrieved ballots
	var ballots []*models.Ballot
	err = attributevalue.UnmarshalListOfMaps(out.Items, &ballots)
	return ballots, err
}
