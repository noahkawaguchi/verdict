package api

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/google/uuid"
	"github.com/noahkawaguchi/verdict/backend/internal/datastore"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func createPollHandler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) events.APIGatewayProxyResponse {
	// Unmarshal the request
	var req pollRequestResponse
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response400("invalid request")
	}
	// Validate the request
	if err := req.validateFields(); err != nil {
		return response400(err.Error())
	}
	// Create the new poll
	poll, pollId := models.NewPoll(req.Prompt, req.Choices)
	// Put the poll in the database
	if err := datastore.PutPoll(ctx, poll); err != nil {
		return response500("failed to put the poll in the database")
	}
	// Send the poll ID back in the response
	return response201(`{"pollId": "` + pollId + `"}`)
}

func createBallotHandler(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) events.APIGatewayProxyResponse {
	// Check for the poll ID
	pollId := request.PathParameters["pollId"]
	if pollId == "" {
		return response400("missing pollId")
	}
	// Get the poll from the database
	poll, err := datastore.GetPoll(ctx, pollId)
	if err != nil {
		return response500(err.Error())
	}
	// Omit the poll ID
	resp := responseFromPoll(poll)
	// Marshal the struct into JSON
	body, err := json.Marshal(resp)
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
	var req ballotRequest
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return response400("invalid request")
	}
	// Validate the request
	if err := req.validateFields(); err != nil {
		return response400(err.Error())
	}
	// Create the new ballot
	userID := uuid.New().String() // Eventually will only do this for non-authenticated polls
	ballot := models.NewBallot(req.PollID, userID, req.RankOrder)
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
	pollId := request.PathParameters["pollId"]
	if pollId == "" {
		return response400("missing pollId")
	}
	// Get the poll from the database
	poll, err := datastore.GetPoll(ctx, pollId)
	if err != nil {
		return response500(err.Error())
	}
	// Get all the ballots from the database for this poll
	ballots, err := datastore.GetPollBallots(ctx, pollId)
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
