package datastore

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
	"github.com/noahkawaguchi/verdict/backend/internal/utils"
)

type tableModel interface{ models.Ballot | models.Poll }

// tableNameFor determines the appropriate table name based on the type of the item.
func tableNameFor[T tableModel](item *T) *string {
	switch any(item).(type) {
	case *models.Ballot:
		return &ballotsTableInfo.name
	case *models.Poll:
		return &pollsTableInfo.name
	}
	// This will never be reached because all type terms in the type set are covered
	return utils.Ref("")
}

// storeItem marshals and puts an item in the database or returns an error.
func storeItem[T tableModel](ds *dynamoStore, item *T) error {
	// Marshal the item into a DynamoDB-compatible map
	av, err := attributevalue.MarshalMap(item)
	if err == nil {
		// Put the item in the database
		_, err = ds.client.PutItem(ds.ctx, &dynamodb.PutItemInput{
			TableName: tableNameFor(item), Item: av,
		})
	}
	return err
}

// retrieveItem gets and unmarshals an item from the database or returns an error.
func retrieveItem[T tableModel](ds *dynamoStore, key map[string]types.AttributeValue) (*T, error) {
	var out *T
	// Get the item from the database using the provided key
	dbOut, err := ds.client.GetItem(ds.ctx, &dynamodb.GetItemInput{
		TableName: tableNameFor(out), Key: key,
	})
	if err == nil {
		// Unmarshal the retrieved item into the struct to return
		err = attributevalue.UnmarshalMap(dbOut.Item, &out)
	}
	return out, err
}

// retrieveItems queries and unmarshals a slice of items from the database or returns an error.
func retrieveItems[T tableModel](
	ds *dynamoStore, keyConExp *string, expAttVals map[string]types.AttributeValue,
) ([]*T, error) {
	var out []*T
	// Query the database using the provided KCE and EAV
	dbOut, err := ds.client.Query(ds.ctx, &dynamodb.QueryInput{
		TableName:                 tableNameFor(new(T)),
		KeyConditionExpression:    keyConExp,
		ExpressionAttributeValues: expAttVals,
	})
	if err == nil {
		// Unmarshal the result into the slice to return
		err = attributevalue.UnmarshalListOfMaps(dbOut.Items, &out)
	}
	return out, err
}
