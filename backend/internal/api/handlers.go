package api

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func (h *handler) createPoll() events.APIGatewayProxyResponse {
	// Unmarshal the request
	var poll *models.Poll
	if err := json.Unmarshal([]byte(h.req.Body), &poll); err != nil {
		return response400("invalid JSON")
	}
	// Validate the fields
	if err := poll.Validate(); err != nil {
		return response400(err.Error())
	}
	// Put the poll in the database
	if err := h.ds.PutPoll(h.ctx, poll); err != nil {
		return response500("failed to put the poll in the database")
	}
	// Send the poll ID back in the response
	return response201(`{"pollId": "` + poll.GetPollID() + `"}`)
}

func (h *handler) createBallot() events.APIGatewayProxyResponse {
	// Check for the poll ID
	pollID := h.req.PathParameters["pollId"]
	if pollID == "" {
		return response400("missing poll ID")
	}
	// Retrieve the poll from the database
	poll, err := h.ds.GetPoll(h.ctx, pollID)
	if err != nil {
		return response500(err.Error())
	}
	// Marshal the response
	if body, err := json.Marshal(poll); err != nil {
		return response500("failed to marshal response")
	} else {
		return response200(string(body))
	}
}

func (h *handler) castBallot() events.APIGatewayProxyResponse {
	// Unmarshal the request
	var ballot *models.Ballot
	if err := json.Unmarshal([]byte(h.req.Body), ballot); err != nil {
		return response400("invalid JSON")
	}
	// Validate the fields
	if err := ballot.Validate(); err != nil {
		return response400(err.Error())
	}
	// Put the ballot in the database
	if err := h.ds.PutBallot(h.ctx, ballot); err != nil {
		return response500("failed to put the ballot in the database")
	}
	// Send a success message back in the response
	return response201(`{"message": "successfully cast ballot"}`)
}

func (h *handler) getResult() events.APIGatewayProxyResponse {
	// Check for the poll ID
	pollID := h.req.PathParameters["pollId"]
	if pollID == "" {
		return response400("missing poll ID")
	}
	// Get the poll and its ballots from the database
	poll, ballots, err := h.ds.GetPollWithBallots(h.ctx, pollID)
	if err != nil {
		return response500(err.Error())
	}
	// Handle the case where no ballots are found
	if len(ballots) == 0 {
		return response404("no ballots found for the specified poll")
	}
	// Calculate the result
	if body, err := models.CalculateResultData(poll, ballots); err != nil {
		return response500(err.Error())
	} else {
		return response200(body)
	}
}
