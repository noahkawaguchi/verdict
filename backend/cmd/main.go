package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/noahkawaguchi/verdict/backend/internal/api"
	"github.com/noahkawaguchi/verdict/backend/internal/datastore"
)

var dbClient *dynamodb.Client // Set in the init function

func main() {
	lambda.Start(func(
		ctx context.Context, request events.APIGatewayProxyRequest,
	) (events.APIGatewayProxyResponse, error) {
		store := datastore.New(ctx, dbClient)
		handler := api.NewHandler(store, request)
		return handler.Route(), nil
	})
}
