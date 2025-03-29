//go:build (!dev && !test) || all

package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// init sets up the dbClient before main executes, once per cold start.
func init() {
	// Load AWS config for production (region and credentials automatically detected from
	// environment variables, `context.TODO()` because the persistent client does not need the
	// context of each invocation)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("Unable to load SDK config:", err)
	}
	// Set the DynamoDB client
	dbClient = dynamodb.NewFromConfig(cfg)
}
