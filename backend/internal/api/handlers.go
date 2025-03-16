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
) events.APIGatewayProxyResponse {
	// Unmarshal the request
	var poll *models.Poll
	if err := json.Unmarshal([]byte(request.Body), &poll); err != nil {
		return response400("invalid request")
	}
	// Validate the fields
	if err := poll.ValidateFields(); err != nil {
		return response400(err.Error())
	}
	// Put the poll in the database
	if err := datastore.PutPoll(ctx, poll); err != nil {
		return response500("failed to put the poll in the database")
	}
	// Send the poll ID back in the response
	return response201(`{"pollId": "` + poll.GetPollID() + `"}`)
}

func createBallotHandler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) events.APIGatewayProxyResponse {
	// Check for the poll ID
	pollID := request.PathParameters["pollId"]
	if pollID == "" {
		return response400("missing poll ID")
	}
	// Get the poll from the database
	poll, err := datastore.GetPoll(ctx, pollID)
	if err != nil {
		return response500(err.Error())
	}
	// Marshal the struct into JSON
	body, err := json.Marshal(poll)
	if err != nil {
		return response500("failed to marshal response")
	}
	// Send the poll information back in the response
	return response200(string(body))
}

func castBallotHandler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) events.APIGatewayProxyResponse {
	// Unmarshal the request
	var ballot *models.Ballot
	if err := json.Unmarshal([]byte(request.Body), &ballot); err != nil {
		return response400("invalid request")
	}
	// Validate the request
	if err := ballot.ValidateFields(); err != nil {
		return response400(err.Error())
	}
	// Put the ballot in the database
	if err := datastore.PutBallot(ctx, ballot); err != nil {
		return response500("failed to put the ballot in the database")
	}
	// Send a success message back in the response
	return response201(`{"message": "successfully cast ballot"}`)
}

func getResultHandler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) events.APIGatewayProxyResponse {
	// Check for the poll ID
	pollID := request.PathParameters["pollId"]
	if pollID == "" {
		return response400("missing poll ID")
	}
	// Get the poll from the database
	poll, err := datastore.GetPoll(ctx, pollID)
	if err != nil {
		return response500(err.Error())
	}
	// Get all the ballots from the database for this poll
	ballots, err := datastore.GetPollBallots(ctx, pollID)
	if err != nil {
		return response500(err.Error())
	}
	// Calculate the result
	result := models.NewResult(poll, ballots)
	// Marshal the struct into JSON (using its custom MarshalJSON method)
	body, err := json.Marshal(result)
	if err != nil {
		return response500("failed to marshal response")
	}
	// Send the poll information back in the response
	return response200(string(body))
}
