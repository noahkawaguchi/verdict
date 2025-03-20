//go:build dev || test

package datastore

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var createBallotsTableInput = &dynamodb.CreateTableInput{
	TableName: aws.String(ballotsTableInfo.name),
	AttributeDefinitions: []types.AttributeDefinition{
		{AttributeName: aws.String(ballotsTableInfo.partitionKey), AttributeType: types.ScalarAttributeTypeS},
		{AttributeName: aws.String(ballotsTableInfo.sortKey), AttributeType: types.ScalarAttributeTypeS},
	},
	KeySchema: []types.KeySchemaElement{
		{AttributeName: aws.String(ballotsTableInfo.partitionKey), KeyType: types.KeyTypeHash},
		{AttributeName: aws.String(ballotsTableInfo.sortKey), KeyType: types.KeyTypeRange},
	},
	BillingMode: types.BillingModePayPerRequest,
}

var createPollsTableInput = &dynamodb.CreateTableInput{
	TableName: aws.String(pollsTableInfo.name),
	AttributeDefinitions: []types.AttributeDefinition{
		{AttributeName: aws.String(pollsTableInfo.partitionKey), AttributeType: types.ScalarAttributeTypeS},
	},
	KeySchema: []types.KeySchemaElement{
		{AttributeName: aws.String(pollsTableInfo.partitionKey), KeyType: types.KeyTypeHash},
	},
	BillingMode: types.BillingModePayPerRequest,
}

// localClientSetup connects to the local Docker DynamoDB for development purposes.
func localClientSetup(endpoint string) {
	// Set up AWS config for local development
	configFunctions := []func(*config.LoadOptions) error{
		config.WithRegion("us-east-2"),    // (Ohio) Required but not used locally
		config.WithBaseEndpoint(endpoint), // Local DynamoDB in Docker
		config.WithCredentialsProvider( // Required but not checked locally
			credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
		),
	}
	if cfg, err := config.LoadDefaultConfig(context.TODO(), configFunctions...); err != nil {
		log.Printf("Unable to load SDK config (development).\ndbClient will be nil:\n%v\n", err)
	} else { // Set the DynamoDB client
		dbClient = dynamodb.NewFromConfig(cfg)
	}
}

// localTableExists checks if the specified table exists in the local DynamoDB in Docker.
func localTableExists(table *tableInfo) bool {
	if _, err := dbClient.DescribeTable(
		context.TODO(),
		&dynamodb.DescribeTableInput{TableName: aws.String(table.name)},
	); err == nil {
		return true
	}
	return false
}

// createLocalTable creates the specified table in the local DynamoDB in Docker.
func createLocalTable(table *tableInfo, input *dynamodb.CreateTableInput) {
	// Attempt to create the table
	if _, err := dbClient.CreateTable(context.TODO(), input); err != nil {
		log.Printf("failed to create table %s: %v", table.name, err)
		return
	}
	// Wait for the table to be created
	for {
		if out, err := dbClient.DescribeTable(
			context.TODO(),
			&dynamodb.DescribeTableInput{TableName: aws.String(table.name)},
		); err == nil && out.Table.TableStatus == types.TableStatusActive {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// deleteLocalTable deletes the specified table from the local DynamoDB in Docker if it exists.
func deleteLocalTable(table *tableInfo) {
	// If the table doesn't exist, don't attempt to delete it
	if !localTableExists(table) {
		return
	}
	// Attempt to delete the table
	if _, err := dbClient.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(table.name),
	}); err != nil {
		log.Printf("failed to delete %s table: %v", table.name, err)
		return
	}
	// Wait for the table to be deleted
	for {
		if _, err := dbClient.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(table.name),
		}); err != nil {
			// Should be a resource not found error
			var notFoundError *types.ResourceNotFoundException
			if errors.As(err, &notFoundError) {
				return
			}
			log.Printf("unexpected error waiting for %s table to be deleted: %v", table.name, err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// printLocalTables prints a list of the tables in the local database.
func printLocalTables() {
	if result, err := dbClient.ListTables(
		context.TODO(),
		&dynamodb.ListTablesInput{},
	); err != nil {
		log.Printf("failed to print local tables: %v", err)
	} else {
		fmt.Println("Local tables:", result.TableNames)
	}
}
