//go:build dev

package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/noahkawaguchi/verdict/backend/internal/datastore"
)

// init sets up the dbClient before main executes, once per cold start.
func init() {
	var endpoint string
	if os.Getenv("AWS_SAM_LOCAL") == "true" { // Running the Lambda function locally with SAM
		endpoint = "http://host.docker.internal:8000"
	} else { // Just running the code in development, not using SAM
		endpoint = "http://localhost:8000"
	}
	configFunctions := []func(*config.LoadOptions) error{
		config.WithRegion("us-east-2"),    // (Ohio) Required but not used locally
		config.WithBaseEndpoint(endpoint), // Local DynamoDB in Docker
		config.WithCredentialsProvider( // Required but not checked locally
			credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
		),
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), configFunctions...)
	if err != nil {
		log.Printf("Unable to load SDK config (development).\ndbClient will be nil:\n%v\n", err)
	}
	// Set the DynamoDB client
	dbClient = dynamodb.NewFromConfig(cfg)
	// Create the tables if they don't exist
	datastore.EnsureBothLocalTablesExist(dbClient)
}
