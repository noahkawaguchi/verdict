package datastore

import "github.com/aws/aws-sdk-go-v2/service/dynamodb"

type tableInfo struct {
	name, partitionKey, sortKey string
}

var dbClient *dynamodb.Client
