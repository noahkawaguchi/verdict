package datastore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/noahkawaguchi/verdict/internal/voting"
)

type dynamoClient interface {
	PutItem(
		ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)

	GetItem(
		ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options),
	) (*dynamodb.GetItemOutput, error)

	Query(
		ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options),
	) (*dynamodb.QueryOutput, error)
}

type dynamoStore struct {
	ctx    context.Context
	client dynamoClient
}

type tableInfo struct{ name, partitionKey, sortKey string }

var ballotsTableInfo = &tableInfo{"Ballots", "PollID", "UserID"}

var pollsTableInfo = &tableInfo{name: "Polls", partitionKey: "PollID"} // No sort key

func New(ctx context.Context, client dynamoClient) *dynamoStore { return &dynamoStore{ctx, client} }

// PutPoll marshals a poll and puts it into the database.
func (d *dynamoStore) PutPoll(poll *voting.Poll) error {
	av, err := attributevalue.MarshalMap(poll)
	if err != nil {
		return err
	}
	_, err = d.client.PutItem(
		d.ctx,
		&dynamodb.PutItemInput{TableName: &pollsTableInfo.name, Item: av},
	)
	return err
}

// PutBallot marshals a ballot and puts it into the database.
func (d *dynamoStore) PutBallot(ballot *voting.Ballot) error {
	av, err := attributevalue.MarshalMap(ballot)
	if err != nil {
		return err
	}
	_, err = d.client.PutItem(
		d.ctx,
		&dynamodb.PutItemInput{TableName: &ballotsTableInfo.name, Item: av},
	)
	return err
}

// GetPoll retrieves and unmarshals a poll from the database.
func (d *dynamoStore) GetPoll(pollID string) (*voting.Poll, error) {
	var out *voting.Poll
	input := dynamodb.GetItemInput{
		TableName: &pollsTableInfo.name,
		Key: map[string]types.AttributeValue{
			pollsTableInfo.partitionKey: &types.AttributeValueMemberS{Value: pollID},
		},
	}
	dbOut, err := d.client.GetItem(d.ctx, &input)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalMap(dbOut.Item, &out)
	return out, err
}

// GetBallots retrieves and unmarshals all of the ballots for the specified poll from the database.
func (d *dynamoStore) GetBallots(pollID string) ([]*voting.Ballot, error) {
	var out []*voting.Ballot
	// Query for ballots by poll ID
	q := dynamodb.QueryInput{
		TableName: &ballotsTableInfo.name,
		KeyConditionExpression: aws.String(
			fmt.Sprintf("%s = :pk", pollsTableInfo.partitionKey),
		),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: pollID},
		},
	}
	dbOut, err := d.client.Query(d.ctx, &q)
	if err != nil {
		return nil, err
	}
	err = attributevalue.UnmarshalListOfMaps(dbOut.Items, &out)
	return out, err
}
