package datastore_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/noahkawaguchi/verdict/backend/internal/datastore"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func TestPutBallot_Error(t *testing.T) {
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		PutItemMock: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("mocked error")
		},
	})

	tests := []*models.Ballot{
		models.NewBallot("poll1", "user1", []int{0, 2, 3, 1}),
		models.NewBallot("poll1", "user2", []int{1, 3, 0, 2}),
		models.NewBallot("poll2", "user4", []int{1, 0, 2}),
	}

	for _, test := range tests {
		if err := tableStore.PutBallot(test); err == nil || err.Error() != "mocked error" {
			t.Error(`expected "mocked error", got:`, err)
		}
	}
}

func TestPutBallot_Success(t *testing.T) {
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		PutItemMock: func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	})

	tests := []*models.Ballot{
		models.NewBallot("poll1", "user1", []int{0, 2, 3, 1}),
		models.NewBallot("poll1", "user2", []int{1, 3, 0, 2}),
		models.NewBallot("poll2", "user4", []int{1, 0, 2}),
	}

	for _, test := range tests {
		if err := tableStore.PutBallot(test); err != nil {
			t.Error("expected success, got:", err)
		}
	}
}

func TestGetBallots_Error(t *testing.T) {
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		QueryMock: func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("mocked error")
		},
	})
	if _, err := tableStore.GetBallots("any poll"); err == nil || err.Error() != "mocked error" {
		t.Error(`expected "mocked error", got:`, err)
	}
}

func TestGetBallots_Success(t *testing.T) {
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		QueryMock: func(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{}, nil
		},
	})
	if _, err := tableStore.GetBallots("any poll"); err != nil {
		t.Error("expected success, got:", err)
	}
}
