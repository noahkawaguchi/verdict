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
	// Unmarshal and validate the request
	poll, pollID, err := models.ValidatedPollFromJSON(request.Body)
	if err != nil {
		return response400(err.Error())
	}
	// Put the poll in the database
	if err := datastore.PutPoll(ctx, poll); err != nil {
		return response500("failed to put the poll in the database")
	}
	// Send the poll ID back in the response
	return response201(`{"pollId": "` + pollID + `"}`)
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
	// Retrieve the poll data from the database
	if body, err := datastore.GetPollData(ctx, pollID); err != nil {
		return response500(err.Error())
	} else {
		return response200(body)
	}
}

func castBallotHandler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) events.APIGatewayProxyResponse {
	// Unmarshal and validate the request
	ballot, err := models.ValidatedBallotFromJSON(request.Body)
	if err != nil {
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
	// Get the poll and its ballots from the database
	poll, ballots, err := datastore.GetPollWithBallots(ctx, pollID)
	if err != nil {
		return response500(err.Error())
	}
	// Handle the case where no ballots are found
	if len(ballots) == 0 {
		return response404("no ballots found for the specified poll")
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
