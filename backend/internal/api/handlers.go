package api

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

func (h *handler) createPoll() events.APIGatewayProxyResponse {
	// Unmarshal and validate the request
	poll, pollID, err := models.ValidatedPollFromJSON(h.req.Body)
	if err != nil {
		return response400(err.Error())
	}
	// Put the poll in the database
	if err = h.ds.PutPoll(h.ctx, poll); err != nil {
		return response500("failed to put the poll in the database")
	}
	// Send the poll ID back in the response
	return response201(`{"pollId": "` + pollID + `"}`)
}

func (h *handler) createBallot() events.APIGatewayProxyResponse {
	// Check for the poll ID
	pollID := h.req.PathParameters["pollId"]
	if pollID == "" {
		return response400("missing poll ID")
	}
	// Retrieve the poll data from the database
	if body, err := h.ds.GetPollData(h.ctx, pollID); err != nil {
		return response500(err.Error())
	} else {
		return response200(body)
	}
}

func (h *handler) castBallot() events.APIGatewayProxyResponse {
	// Unmarshal and validate the request
	ballot, err := models.ValidatedBallotFromJSON(h.req.Body)
	if err != nil {
		return response400(err.Error())
	}
	// Put the ballot in the database
	if err = h.ds.PutBallot(h.ctx, ballot); err != nil {
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
