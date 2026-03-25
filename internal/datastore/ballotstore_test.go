package datastore_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/noahkawaguchi/verdict/internal/datastore"
	"github.com/noahkawaguchi/verdict/internal/voting"
)

func TestPutBallot_Error(t *testing.T) {
	t.Parallel()
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		PutItemMock: func(
			context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
		) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("mocked error")
		},
	})

	tests := []struct {
		name   string
		ballot *voting.Ballot
	}{
		{"poll1 user1", voting.NewBallot("poll1", "user1", []int{0, 2, 3, 1})},
		{"poll1 user2", voting.NewBallot("poll1", "user2", []int{1, 3, 0, 2})},
		{"poll2 user4", voting.NewBallot("poll2", "user4", []int{1, 0, 2})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tableStore.PutBallot(
				tt.ballot,
			); err == nil ||
				err.Error() != "mocked error" {
				t.Error(`expected "mocked error", got:`, err)
			}
		})
	}
}

func TestPutBallot_Success(t *testing.T) {
	t.Parallel()
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		PutItemMock: func(
			context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
		) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	})

	tests := []struct {
		name   string
		ballot *voting.Ballot
	}{
		{"poll1 user1", voting.NewBallot("poll1", "user1", []int{0, 2, 3, 1})},
		{"poll1 user2", voting.NewBallot("poll1", "user2", []int{1, 3, 0, 2})},
		{"poll2 user4", voting.NewBallot("poll2", "user4", []int{1, 0, 2})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tableStore.PutBallot(tt.ballot); err != nil {
				t.Error("expected success, got:", err)
			}
		})
	}
}

func TestGetBallots_Error(t *testing.T) {
	t.Parallel()
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		QueryMock: func(
			context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options),
		) (*dynamodb.QueryOutput, error) {
			return nil, errors.New("mocked error")
		},
	})
	if _, err := tableStore.GetBallots(
		"any poll",
	); err == nil ||
		err.Error() != "mocked error" {
		t.Error(`expected "mocked error", got:`, err)
	}
}

func TestGetBallots_Success(t *testing.T) {
	t.Parallel()
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		QueryMock: func(
			context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options),
		) (*dynamodb.QueryOutput, error) {
			return &dynamodb.QueryOutput{}, nil
		},
	})
	if _, err := tableStore.GetBallots("any poll"); err != nil {
		t.Error("expected success, got:", err)
	}
}
