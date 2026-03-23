//go:build dev

package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/noahkawaguchi/verdict/internal/datastore"
)

// devDBContainerName is the name for the local DynamoDB container. It must match the one used in
// the `justfile`.
const devDBContainerName string = "verdict-dev-db"

// init sets up the dbClient before main executes, once per cold start.
func init() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithBaseEndpoint("http://"+devDBContainerName+":8000"),
	)
	if err != nil {
		log.Fatal("Unable to load SDK config:", err)
	}
	// Set the DynamoDB client
	dbClient = dynamodb.NewFromConfig(cfg)
	// Create the tables if they don't exist
	datastore.EnsureBothLocalTablesExist(dbClient)
}
