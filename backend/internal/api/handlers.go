package api

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func (h *Handler) createPoll() events.APIGatewayProxyResponse {
	// Unmarshal the request
	var poll *models.Poll
	if err := json.Unmarshal([]byte(h.Req.Body), &poll); err != nil {
		return response400("invalid JSON")
	}
	// Validate the fields
	if err := poll.Validate(); err != nil {
		return response400(err.Error())
	}
	// Put the poll in the database
	if err := h.Store.PutPoll(poll); err != nil {
		return response500("failed to put the poll in the database")
	}
	// Send the poll ID back in the response
	return response201(`{"pollId": "` + poll.GetPollID() + `"}`)
}

func (h *Handler) getPollInfo() events.APIGatewayProxyResponse {
	// Check for the poll ID
	pollID := h.Req.PathParameters["pollId"]
	if pollID == "" {
		return response400("missing poll ID")
	}
	// Retrieve the poll from the database
	poll, err := h.Store.GetPoll(pollID)
	if err != nil {
		return response500("failed to get the poll from the database")
	}
	// Handle nonexistent polls
	if err = poll.Validate(); err != nil {
		return response404("no poll found for the specified ID")
	}
	// Marshal the response
	body, err := json.Marshal(poll)
	if err != nil {
		return response500("failed to marshal response")
	}
	return response200(string(body))
}

func (h *Handler) castBallot() events.APIGatewayProxyResponse {
	// Unmarshal the request
	var ballot *models.Ballot
	if err := json.Unmarshal([]byte(h.Req.Body), &ballot); err != nil {
		return response400("invalid JSON")
	}
	// Validate the fields
	if err := ballot.Validate(); err != nil {
		return response400(err.Error())
	}
	// Put the ballot in the database
	if err := h.Store.PutBallot(ballot); err != nil {
		return response500("failed to put the ballot in the database")
	}
	// Send a success message back in the response
	return response201(`{"message": "successfully cast ballot"}`)
}

func (h *Handler) getResult() events.APIGatewayProxyResponse {
	// Check for the poll ID
	pollID := h.Req.PathParameters["pollId"]
	if pollID == "" {
		return response400("missing poll ID")
	}
	// Get the poll from the database
	poll, err := h.Store.GetPoll(pollID)
	if err != nil {
		return response500("failed to get the poll from the database")
	}
	// Handle nonexistent polls
	if err = poll.Validate(); err != nil {
		return response404("no poll found for the specified ID")
	}
	// Get the poll's ballots from the database
	ballots, err := h.Store.GetBallots(pollID)
	if err != nil {
		return response500("failed to get the poll's ballots from the database")
	}
	// Handle the case where no ballots are found
	if len(ballots) == 0 {
		return response404("no ballots found for the specified poll")
	}
	// Calculate the result
	result, err := models.NewResult(poll, ballots)
	if err != nil {
		return response500(err.Error())
	}
	// Marshal the response
	body, err := json.Marshal(result)
	if err != nil {
		return response500("failed to marshal response")
	}
	return response200(string(body))
}
