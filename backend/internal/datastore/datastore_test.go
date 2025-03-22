package datastore_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// mockDynamo implements the dynamoClient interface for testing purposes.
type mockDynamo struct {
	PutItemMock func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItemMock func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	QueryMock   func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

func (md *mockDynamo) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if md.PutItemMock != nil {
		return md.PutItemMock(ctx, params, optFns...)
	}
	return nil, nil
}

func (md *mockDynamo) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if md.GetItemMock != nil {
		return md.GetItemMock(ctx, params, optFns...)
	}
	return nil, nil
}

func (md *mockDynamo) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if md.QueryMock != nil {
		return md.QueryMock(ctx, params, optFns...)
	}
	return nil, nil
}
