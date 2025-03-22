package datastore_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/noahkawaguchi/verdict/backend/internal/datastore"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func TestPutPoll_Error(t *testing.T) {
	tableStore := datastore.TableStore{Ctx: context.TODO(), Client: &mockDynamo{
		PutItemMock: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("mocked error")
		},
	}}

	tests := []*models.Poll{
		models.NewPoll("What is the best programming language?", []string{"Go", "Rust", "C++"}),
		models.NewPoll("What is the best int size?", []string{"32", "64", "8", "anything unsigned"}),
	}

	for _, test := range tests {
		if err := tableStore.PutPoll(test); err == nil || err.Error() != "mocked error" {
			t.Error(`expected "mocked error", got:`, err)
		}
	}
}

func TestPutPoll_Success(t *testing.T) {
	tableStore := datastore.TableStore{Ctx: context.TODO(), Client: &mockDynamo{
		PutItemMock: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	}}

	tests := []*models.Poll{
		models.NewPoll("What is the best programming language?", []string{"Go", "Rust", "C++"}),
		models.NewPoll("What is the best int size?", []string{"32", "64", "8", "anything unsigned"}),
	}

	for _, test := range tests {
		if err := tableStore.PutPoll(test); err != nil {
			t.Error("expected success, got:", err)
		}
	}
}

func TestGetPoll_Error(t *testing.T) {
	tableStore := datastore.TableStore{Ctx: context.TODO(), Client: &mockDynamo{
		GetItemMock: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return nil, errors.New("mocked error")
		},
	}}
	if _, err := tableStore.GetPoll("any poll"); err == nil || err.Error() != "mocked error" {
		t.Error(`expected "mocked error", got:`, err)
	}
}

func TestGetPoll_Success(t *testing.T) {
	tableStore := datastore.TableStore{Ctx: context.TODO(), Client: &mockDynamo{
		GetItemMock: func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{}, nil
		},
	}}
	if _, err := tableStore.GetPoll("any poll"); err != nil {
		t.Error("expected success, got:", err)
	}
}
