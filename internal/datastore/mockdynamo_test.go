package datastore_test

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// mockDynamo implements the dynamoClient interface for testing purposes.
type mockDynamo struct {
	PutItemMock func(
		context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)

	GetItemMock func(
		context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.GetItemOutput, error)

	QueryMock func(
		context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options),
	) (*dynamodb.QueryOutput, error)
}

func (m *mockDynamo) PutItem(
	ctx context.Context,
	params *dynamodb.PutItemInput,
	optFns ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	if m.PutItemMock != nil {
		return m.PutItemMock(ctx, params, optFns...)
	}
	return nil, nil
}

func (m *mockDynamo) GetItem(
	ctx context.Context,
	params *dynamodb.GetItemInput,
	optFns ...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	if m.GetItemMock != nil {
		return m.GetItemMock(ctx, params, optFns...)
	}
	return nil, nil
}

func (m *mockDynamo) Query(
	ctx context.Context,
	params *dynamodb.QueryInput,
	optFns ...func(*dynamodb.Options),
) (*dynamodb.QueryOutput, error) {
	if m.QueryMock != nil {
		return m.QueryMock(ctx, params, optFns...)
	}
	return nil, nil
}
