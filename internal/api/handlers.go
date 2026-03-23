package api

import (
	"encoding/json"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/internal/voting"
)

// createPoll validates and stores a new poll, responding with the ID of the newly created poll.
func (h *handler) createPoll() events.APIGatewayProxyResponse {
	var poll *voting.Poll
	if err := json.Unmarshal([]byte(h.req.Body), &poll); err != nil {
		return resp400("invalid JSON")
	}
	if err := poll.Validate(); err != nil {
		return resp400(err.Error())
	}
	if err := h.store.PutPoll(poll); err != nil {
		slog.Error("failed to put poll in DB", "ID", poll.ID(), "err", err)
		return resp500
	}
	return resp201(`{"pollId":"` + poll.ID() + `"}`)
}

// getPollInfo retrieves the poll specified in the path parameter by ID.
func (h *handler) getPollInfo() events.APIGatewayProxyResponse {
	// This check for the poll ID is redundant with the path check in the router in most cases
	pollID := h.req.PathParameters["pollId"]
	if pollID == "" {
		return resp400("missing poll ID")
	}
	poll, err := h.store.GetPoll(pollID)
	if err != nil {
		slog.Error("failed to get poll from DB", "ID", pollID, "err", err)
		return resp500
	}
	// Handle nonexistent polls here
	if err = poll.Validate(); err != nil {
		return resp404("no poll found for poll ID " + pollID)
	}
	body, err := json.Marshal(poll)
	if err != nil {
		slog.Error("failed to marshal response", "err", err)
		return resp500
	}
	return resp200(string(body))
}

// castBallot validates and stores a ballot.
func (h *handler) castBallot() events.APIGatewayProxyResponse {
	var ballot *voting.Ballot
	if err := json.Unmarshal([]byte(h.req.Body), &ballot); err != nil {
		return resp400("invalid JSON")
	}
	if err := ballot.Validate(); err != nil {
		return resp400(err.Error())
	}
	if err := h.store.PutBallot(ballot); err != nil {
		slog.Error("failed to put ballot in DB", "err", err)
		return resp500
	}
	return resp201(`{"message":"successfully cast ballot"}`)
}

// getResult retrieves a poll and its ballots by poll ID, then calculates and responds with the
// result.
func (h *handler) getResult() events.APIGatewayProxyResponse {
	// This check for the poll ID is redundant with the path check in the router in most cases
	pollID := h.req.PathParameters["pollId"]
	if pollID == "" {
		return resp400("missing poll ID")
	}
	poll, err := h.store.GetPoll(pollID)
	if err != nil {
		slog.Error("failed to get poll from DB", "err", err)
		return resp500
	}
	// Handle nonexistent polls here
	if err = poll.Validate(); err != nil {
		return resp404("no poll found for poll ID " + pollID)
	}
	ballots, err := h.store.GetBallots(pollID)
	if err != nil {
		slog.Error("failed to get poll's ballots from DB", "err", err)
		return resp500
	}
	// Handle the case where no ballots are found here
	if len(ballots) == 0 {
		return resp404("no ballots found for the specified poll")
	}
	result, err := voting.NewResult(poll, ballots)
	if err != nil {
		slog.Error("failed to calculate result", "err", err)
		return resp500
	}
	body, err := json.Marshal(result)
	if err != nil {
		slog.Error("failed to marshal response", "err", err)
		return resp500
	}
	return resp200(string(body))
}
