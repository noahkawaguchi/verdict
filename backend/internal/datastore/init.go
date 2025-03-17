package datastore

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type tableInfo struct {
	name, partitionKey, sortKey string
}

var dbClient *dynamodb.Client

// init sets up the dbClient before main executes, once per cold start
func init() {
	if os.Getenv("USE_LOCAL_DYNAMO") == "true" { // Local development without SAM
		developmentSetup("http://localhost:8000") 
	} else if os.Getenv("AWS_SAM_LOCAL") == "true" { // Local development with SAM
		developmentSetup("http://host.docker.internal:8000") 
	} else {
		productionSetup()
	}
}

func productionSetup() {
	// Load AWS config for production (region and credentials automatically detected from
	// environment variables)
	if cfg, err := config.LoadDefaultConfig(context.TODO()); err != nil {
		log.Printf("Unable to load SDK config (production).\ndbClient will be nil:\n%s\n",
			err.Error())
	} else { // Set the DynamoDB client
		dbClient = dynamodb.NewFromConfig(cfg)
	}
}

// developmentSetup connects to the local Docker DynamoDB for development purposes.
func developmentSetup(endpoint string) {
	// Set up AWS config for local development
	configFunctions := []func(*config.LoadOptions) error{
		config.WithRegion("us-east-2"),    // (Ohio) Required but not used locally
		config.WithBaseEndpoint(endpoint), // Local DynamoDB in Docker
		config.WithCredentialsProvider( // Required but not checked locally
			credentials.NewStaticCredentialsProvider("dummy", "dummy", ""),
		),
	}
	if cfg, err := config.LoadDefaultConfig(context.TODO(), configFunctions...); err != nil {
		log.Printf("Unable to load SDK config (development).\ndbClient will be nil:\n%s\n",
			err.Error())
		return
	} else { // Set the DynamoDB client
		dbClient = dynamodb.NewFromConfig(cfg)
	}
	// Create the tables if they don't exist
	ballotsTableInput := &dynamodb.CreateTableInput{
		TableName: aws.String(ballotsTableInfo.name),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(ballotsTableInfo.partitionKey),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String(ballotsTableInfo.sortKey),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(ballotsTableInfo.partitionKey),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String(ballotsTableInfo.sortKey),
				KeyType:       types.KeyTypeRange,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}
	if err := ensureLocalTableExists(ballotsTableInfo, ballotsTableInput); err != nil {
		log.Printf("Failed to ensure Ballots table exists: %s\n", err.Error())
	}
	pollsTableInput := &dynamodb.CreateTableInput{
		TableName: aws.String(pollsTableInfo.name),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(pollsTableInfo.partitionKey),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(pollsTableInfo.partitionKey),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	}
	if err := ensureLocalTableExists(pollsTableInfo, pollsTableInput); err != nil {
		log.Printf("Failed to ensure Polls table exists: %s\n", err.Error())
	}
	printLocalTables()
}

// ensureLocalTableExists creates the table if it doesn't exist. This is only for local
// development. For production, create the table from the AWS console instead of in code.
func ensureLocalTableExists(table *tableInfo, input *dynamodb.CreateTableInput) error {
	// Check if the table exists
	if _, err := dbClient.DescribeTable(
		context.TODO(),
		&dynamodb.DescribeTableInput{TableName: aws.String(table.name)},
	); err == nil { // Already exists
		fmt.Printf("The table %s already existed\n", table.name)
		printLocalTables()
		return nil
	}
	// Need to create it
	fmt.Printf("Creating the table %s...\n", table.name)
	if _, err := dbClient.CreateTable(context.TODO(), input); err != nil {
		return err
	}
	// Wait (forever) for the table to be created (only used in local development)
	fmt.Printf("Waiting for the table %s to be created...\n", table.name)
	for {
		if out, err := dbClient.DescribeTable(
			context.TODO(),
			&dynamodb.DescribeTableInput{TableName: aws.String(table.name)},
		); err == nil && out.Table.TableStatus == types.TableStatusActive {
			break
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("The table %s has now been created\n", table.name)
	return nil
}

// printLocalTables prints a list of the tables in the local database.
func printLocalTables() {
	if result, err := dbClient.ListTables(
		context.TODO(),
		&dynamodb.ListTablesInput{},
	); err != nil {
		log.Printf("Failed to list tables: %s\n", err.Error())
	} else {
		fmt.Println("Local tables:", result.TableNames)
	}
}
