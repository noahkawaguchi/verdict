//go:build (!dev && !test) || all

package datastore

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
	if cfg, err := config.LoadDefaultConfig(context.TODO()); err != nil {
		log.Printf("Unable to load SDK config (production).\ndbClient will be nil:\n%s\n",
			err.Error())
	} else { // Set the DynamoDB client
		dbClient = dynamodb.NewFromConfig(cfg)
	}
}
