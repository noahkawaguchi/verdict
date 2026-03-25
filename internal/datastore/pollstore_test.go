package datastore_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/noahkawaguchi/verdict/internal/datastore"
	"github.com/noahkawaguchi/verdict/internal/voting"
)

func TestPutPoll_Error(t *testing.T) {
	t.Parallel()
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		PutItemMock: func(
			context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
		) (*dynamodb.PutItemOutput, error) {
			return nil, errors.New("mocked error")
		},
	})

	tests := []struct {
		name string
		poll *voting.Poll
	}{
		{
			"three choices",
			voting.NewPoll(
				"What is the best programming language?",
				[]string{"Go", "Rust", "C++"},
			),
		},
		{
			"four choices",
			voting.NewPoll(
				"What is the best int size?",
				[]string{"32", "64", "8", "anything unsigned"},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tableStore.PutPoll(
				tt.poll,
			); err == nil ||
				err.Error() != "mocked error" {
				t.Error(`expected "mocked error", got:`, err)
			}
		})
	}
}

func TestPutPoll_Success(t *testing.T) {
	t.Parallel()
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		PutItemMock: func(
			context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
		) (*dynamodb.PutItemOutput, error) {
			return &dynamodb.PutItemOutput{}, nil
		},
	})

	tests := []struct {
		name string
		poll *voting.Poll
	}{
		{
			"three choices",
			voting.NewPoll(
				"What is the best programming language?",
				[]string{"Go", "Rust", "C++"},
			),
		},
		{
			"four choices",
			voting.NewPoll(
				"What is the best int size?",
				[]string{"32", "64", "8", "anything unsigned"},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tableStore.PutPoll(tt.poll); err != nil {
				t.Error("expected success, got:", err)
			}
		})
	}
}

func TestGetPoll_Error(t *testing.T) {
	t.Parallel()
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		GetItemMock: func(
			context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
		) (*dynamodb.GetItemOutput, error) {
			return nil, errors.New("mocked error")
		},
	})
	if _, err := tableStore.GetPoll("any poll"); err == nil || err.Error() != "mocked error" {
		t.Error(`expected "mocked error", got:`, err)
	}
}

func TestGetPoll_Success(t *testing.T) {
	t.Parallel()
	tableStore := datastore.New(context.TODO(), &mockDynamo{
		GetItemMock: func(
			context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
		) (*dynamodb.GetItemOutput, error) {
			return &dynamodb.GetItemOutput{}, nil
		},
	})
	if _, err := tableStore.GetPoll("any poll"); err != nil {
		t.Error("expected success, got:", err)
	}
}
