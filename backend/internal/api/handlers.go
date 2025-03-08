package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/backend/internal/datastore"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func createPollHandler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	// Unmarshal the request
	var req createPollRequest
	err := json.Unmarshal([]byte(request.Body), &req)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Invalid request"}`,
		}, nil
	}
	// Validate the request
	if !req.allValidFields() {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       `{"error": "Fields cannot be empty"}`,
		}, nil
	}
	// Create the new poll
	poll, pollID := models.NewPoll(req.Prompt, req.Choices)
	// Put the new poll in the database
	if err = datastore.PutPoll(ctx, poll); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       `{"error": "Failed to put the new poll in the database"}`,
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       `{"pollID": "` + pollID + `"}`,
	}, nil
}
