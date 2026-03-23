package api

import (
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/internal/models"
)

func (h *handler) createPoll() events.APIGatewayProxyResponse {
	// Unmarshal the request
	var poll *models.Poll
	if err := json.Unmarshal([]byte(h.req.Body), &poll); err != nil {
		return resp400("invalid JSON")
	}
	// Validate the fields
	if err := poll.Validate(); err != nil {
		return resp400(err.Error())
	}
	// Put the poll in the database
	if err := h.store.PutPoll(poll); err != nil {
		slog.Error("failed to put poll in DB", "ID", poll.ID(), "err", err)
		return resp500
	}
	// Send the poll ID back in the response
	return resp201(`{"pollId":"` + poll.ID() + `"}`)
}

func (h *handler) getPollInfo() events.APIGatewayProxyResponse {
	// Check for the poll ID (redundant with the path check in the router in most cases)
	pollID := h.req.PathParameters["pollId"]
	if pollID == "" {
		return resp400("missing poll ID")
	}
	// Retrieve the poll from the database
	poll, err := h.store.GetPoll(pollID)
	if err != nil {
		slog.Error("failed to get poll from DB", "ID", pollID, "err", err)
		return resp500
	}
	// Handle nonexistent polls
	if err = poll.Validate(); err != nil {
		return resp404("no poll found for poll ID " + pollID)
	}
	// Marshal the response
	body, err := json.Marshal(poll)
	if err != nil {
		slog.Error("failed to marshal response", "err", err)
		return resp500
	}
	return resp200(string(body))
}

func (h *handler) castBallot() events.APIGatewayProxyResponse {
	// Unmarshal the request
	var ballot *models.Ballot
	if err := json.Unmarshal([]byte(h.req.Body), &ballot); err != nil {
		return resp400("invalid JSON")
	}
	// Validate the fields
	if err := ballot.Validate(); err != nil {
		return resp400(err.Error())
	}
	// Put the ballot in the database
	if err := h.store.PutBallot(ballot); err != nil {
		slog.Error("failed to put ballot in DB", "err", err)
		return resp500
	}
	// Send a success message back in the response
	return resp201(`{"message":"successfully cast ballot"}`)
}

func (h *handler) getResult() events.APIGatewayProxyResponse {
	// Check for the poll ID
	pollID := h.req.PathParameters["pollId"]
	if pollID == "" {
		return resp400("missing poll ID")
	}
	// Get the poll from the database
	poll, err := h.store.GetPoll(pollID)
	if err != nil {
		slog.Error("failed to get poll from DB", "err", err)
		return resp500
	}
	// Handle nonexistent polls
	if err = poll.Validate(); err != nil {
		return resp404("no poll found for poll ID " + pollID)
	}
	// Get the poll's ballots from the database
	ballots, err := h.store.GetBallots(pollID)
	if err != nil {
		slog.Error("failed to get poll's ballots from DB", "err", err)
		return resp500
	}
	// Handle the case where no ballots are found
	if len(ballots) == 0 {
		return resp404("no ballots found for the specified poll")
	}
	// Calculate the result
	result, err := models.NewResult(poll, ballots)
	if err != nil {
		slog.Error("failed to calculate result", "err", err)
		return resp500
	}
	// Marshal the response
	body, err := json.Marshal(result)
	if err != nil {
		slog.Error("failed to marshal response", "err", err)
		return resp500
	}
	return resp200(string(body))
}
