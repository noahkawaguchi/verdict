package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/noahkawaguchi/verdict/backend/internal/api"
	"github.com/noahkawaguchi/verdict/backend/internal/datastore"
)

func main() {
	if false {
		datastore.DatastoreDemo()
	}
	lambda.Start(api.Router)
}
