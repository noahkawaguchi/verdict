package api

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/backend/internal/datastore"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func createPollHandler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	// Unmarshal the request
	var req pollRequestResponse
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response400("invalid request"), nil
	}
	// Validate the request
	if err := req.validateFields(); err != nil {
		return response400(err.Error()), nil
	}
	// Create the new poll
	poll, pollId := models.NewPoll(req.Prompt, req.Choices)
	// Put the new poll in the database
	if err := datastore.PutPoll(ctx, poll); err != nil {
		return response500("failed to put the new poll in the database"), nil
	}
	// Send the poll ID back in the response
	return response200(`{"pollId": "` + pollId + `"}`), nil
}

func createBallotHandler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	// Check for the poll ID
	pollId := request.PathParameters["pollId"]
	if pollId == "" {
		return response400("missing pollId"), nil
	}
	// Get the poll from the database
	poll, err := datastore.GetPoll(ctx, pollId)
	if err != nil {
		return response500(err.Error()), nil
	}
	// Marshal the struct into JSON
	body, err := json.Marshal(poll)
	if err != nil {
		return response500("failed to marshal response"), nil
	}
	// Send the poll information back in the response
	return response200(string(body)), nil
}
