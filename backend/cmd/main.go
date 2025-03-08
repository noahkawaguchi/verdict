package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/noahkawaguchi/verdict/backend/internal/api"
	"github.com/noahkawaguchi/verdict/backend/internal/datastore"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func main() {
	if false {
		datastore.DatastoreDemo()
		models.ResultDemo()
	}
	lambda.Start(api.Router)
}
