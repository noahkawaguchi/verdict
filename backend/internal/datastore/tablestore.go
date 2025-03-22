package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type dynamoClient interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

type TableStore struct {
	Ctx    context.Context
	Client dynamoClient
}

type tableInfo struct {
	name         string
	partitionKey string
	sortKey      string
}

var ballotsTableInfo = &tableInfo{
	name:         "Ballots",
	partitionKey: "PollID",
	sortKey:      "UserID",
}

var pollsTableInfo = &tableInfo{
	name:         "Polls",
	partitionKey: "PollID",
	// (No sort key)
}
