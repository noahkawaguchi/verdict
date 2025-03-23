//go:build dev

package datastore

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// EnsureBothLocalTablesExist ensure that both the Ballots and Polls tables exist for local
// development purposes.
func EnsureBothLocalTablesExist(client *dynamodb.Client) {
	if !localTableExists(client, ballotsTableInfo) {
		createLocalTable(client, ballotsTableInfo, createBallotsTableInput)
	}
	if !localTableExists(client, pollsTableInfo) {
		createLocalTable(client, pollsTableInfo, createPollsTableInput)
	}
	printLocalTables(client)
}

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

// localTableExists checks if the specified table exists in the local DynamoDB in Docker.
func localTableExists(client *dynamodb.Client, table *tableInfo) bool {
	_, err := client.DescribeTable(
		context.TODO(),
		&dynamodb.DescribeTableInput{TableName: aws.String(table.name)},
	)
	return err == nil
}

// createLocalTable creates the specified table in the local DynamoDB in Docker.
func createLocalTable(client *dynamodb.Client, table *tableInfo, input *dynamodb.CreateTableInput) {
	// Attempt to create the table
	if _, err := client.CreateTable(context.TODO(), input); err != nil {
		log.Printf("failed to create table %s: %v", table.name, err)
		return
	}
	// Wait for the table to be created
	for {
		out, err := client.DescribeTable(
			context.TODO(),
			&dynamodb.DescribeTableInput{TableName: aws.String(table.name)},
		)
		if err == nil && out.Table.TableStatus == types.TableStatusActive {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// printLocalTables prints a list of the tables in the local database.
func printLocalTables(client *dynamodb.Client) {
	result, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		log.Printf("failed to print local tables: %v", err)
	} else {
		fmt.Println("Local tables:", result.TableNames)
	}
}
