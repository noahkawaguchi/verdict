package datastore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/noahkawaguchi/verdict/internal/utils"
	"github.com/noahkawaguchi/verdict/internal/voting"
)

type dynamoClient interface {
	PutItem(
		ctx context.Context,
		params *dynamodb.PutItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
	GetItem(
		ctx context.Context,
		params *dynamodb.GetItemInput,
		optFns ...func(*dynamodb.Options),
	) (*dynamodb.GetItemOutput, error)
	Query(
		ctx context.Context,
		params *dynamodb.QueryInput,
		optFns ...func(*dynamodb.Options),
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

// PutPoll creates a new poll entry in the database.
func (ds *dynamoStore) PutPoll(poll *voting.Poll) error { return storeItem(ds, poll) }

// PutBallot creates a new ballot entry in the database.
func (ds *dynamoStore) PutBallot(ballot *voting.Ballot) error { return storeItem(ds, ballot) }

// GetPoll retrieves a poll from the database by its poll ID.
func (ds *dynamoStore) GetPoll(pollID string) (*voting.Poll, error) {
	// Define the key to get the poll by ID
	key := map[string]types.AttributeValue{
		pollsTableInfo.partitionKey: &types.AttributeValueMemberS{Value: pollID},
	}
	return retrieveItem[voting.Poll](ds, key)
}

// GetBallots retrieves all of the ballots for the specified poll from the database.
func (ds *dynamoStore) GetBallots(pollID string) ([]*voting.Ballot, error) {
	// Define the key condition expression and expression attribute values to query by poll ID
	keyConExp := utils.Ref(fmt.Sprintf("%s = :pk", pollsTableInfo.partitionKey))
	expAttVals := map[string]types.AttributeValue{
		":pk": &types.AttributeValueMemberS{Value: pollID},
	}
	return retrieveItems[voting.Ballot](ds, keyConExp, expAttVals)
}
