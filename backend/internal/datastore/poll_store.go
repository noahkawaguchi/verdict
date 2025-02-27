package datastore

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

var pollsTableInfo = &tableInfo{
	name:         "Polls",
	partitionKey: "PollID",
	// No sort key
}

// PutPoll creates a new poll entry in the database.
func PutPoll(ctx context.Context, poll *models.Poll) error {
	// Marshal the struct into a DynamoDB-compatible map
	av, err := attributevalue.MarshalMap(poll)
	if err != nil {
		return err
	}
	// Put the poll into DynamoDB
	_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(pollsTableInfo.name),
		Item:      av,
	})
	return err
}

// getPoll retrieves a poll from the database by its PollID.
func getPoll(ctx context.Context, id string) (*models.Poll, error) {
	out, err := dbClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(pollsTableInfo.name),
		Key: map[string]types.AttributeValue{
			pollsTableInfo.partitionKey: &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return nil, err
	}
	// Unmarshal the retrieved poll into a struct
	var poll models.Poll
	err = attributevalue.UnmarshalMap(out.Item, &poll)
	return &poll, err
}

// // getTilesDue retrieves the TileIDs for all the user's tiles with a NextReview earlier than or
// // equal to the date provided (as a string in the format YYYY-MM-DD).
// func getTilesDue(ctx context.Context, userID, date string) ([]string, error) {
// 	// Query using the GSI
// 	input := &dynamodb.QueryInput{
// 		TableName:              aws.String(tableName),
// 		IndexName:              aws.String(indexName),
// 		KeyConditionExpression: aws.String("UserID = :u AND NextReview <= :d"),
// 		ExpressionAttributeValues: map[string]types.AttributeValue{
// 			":u": &types.AttributeValueMemberS{Value: userID},
// 			":d": &types.AttributeValueMemberS{Value: date}, // String in the format YYYY-MM-DD
// 		},
// 	}
// 	out, err := dbClient.Query(ctx, input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Unmarshal the keys
// 	var keysSlice []tileKeysOnly
// 	err = attributevalue.UnmarshalListOfMaps(out.Items, &keysSlice)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Extract only the TileIDs
// 	tilesDue := make([]string, 0, len(keysSlice))
// 	for _, keys := range keysSlice {
// 		tilesDue = append(tilesDue, keys.TileID)
// 	}
// 	return tilesDue, nil
// }
