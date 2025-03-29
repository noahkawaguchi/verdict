//go:build dev

package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/noahkawaguchi/verdict/backend/internal/datastore"
)

// init sets up the dbClient before main executes, once per cold start.
func init() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		// (Ohio) Required but not used locally
		config.WithRegion("us-east-2"),
		// Local DynamoDB running in Docker, connected via SAM
		config.WithBaseEndpoint("http://host.docker.internal:8000"),
		// Required but not checked locally
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		log.Fatal("Unable to load SDK config:", err)
	}
	// Set the DynamoDB client
	dbClient = dynamodb.NewFromConfig(cfg)
	// Create the tables if they don't exist
	datastore.EnsureBothLocalTablesExist(dbClient)
}
