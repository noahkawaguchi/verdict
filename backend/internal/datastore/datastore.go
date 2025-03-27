package datastore

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
	"github.com/noahkawaguchi/verdict/backend/internal/utils"
)

type dynamoClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

type dynamoStore struct {
	ctx    context.Context
	client dynamoClient
}

type tableInfo struct{ name, partitionKey, sortKey string }

var ballotsTableInfo = &tableInfo{"Ballots", "PollID", "UserID"}

var pollsTableInfo = &tableInfo{name: "Polls", partitionKey: "PollID"} // No sort key

func NewDynamoStore(ctx context.Context, client dynamoClient) *dynamoStore {
	return &dynamoStore{ctx, client}
}

// PutPoll creates a new poll entry in the database.
func (ds *dynamoStore) PutPoll(poll *models.Poll) error { return storeItem(ds, poll) }

// PutBallot creates a new ballot entry in the database.
func (ds *dynamoStore) PutBallot(ballot *models.Ballot) error { return storeItem(ds, ballot) }

// GetPoll retrieves a poll from the database by its poll ID.
func (ds *dynamoStore) GetPoll(pollID string) (*models.Poll, error) {
	// Define the key to get the poll by ID
	key := map[string]types.AttributeValue{
		pollsTableInfo.partitionKey: &types.AttributeValueMemberS{Value: pollID},
	}
	return retrieveItem[models.Poll](ds, key)
}

// GetBallots retrieves all of the ballots for the specified poll from the database.
func (ds *dynamoStore) GetBallots(pollID string) ([]*models.Ballot, error) {
	// Define the key condition expression and expression attribute values to query by poll ID
	keyConExp := utils.Ref(fmt.Sprintf("%s = :pk", pollsTableInfo.partitionKey))
	expAttVals := map[string]types.AttributeValue{":pk": &types.AttributeValueMemberS{Value: pollID}}
	return retrieveItems[models.Ballot](ds, keyConExp, expAttVals)
}
