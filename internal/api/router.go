package api

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/internal/voting"
)

type datastore interface {
	PutPoll(poll *voting.Poll) error
	GetPoll(pollID string) (*voting.Poll, error)
	PutBallot(ballot *voting.Ballot) error
	GetBallots(pollID string) ([]*voting.Ballot, error)
}

type handler struct {
	store datastore
	req   events.APIGatewayProxyRequest
}

func NewHandler(store datastore, req events.APIGatewayProxyRequest) *handler {
	return &handler{store, req}
}

// Route matches the method and path of the request and calls the relevant method or returns an
// appropriate error response.
func (h *handler) Route() events.APIGatewayProxyResponse {
	switch h.req.HTTPMethod {

	case http.MethodPost:
		switch h.req.Path {
		case "/polls":
			return h.createPoll()
		case "/ballots":
			return h.castBallot()
		default:
			return resp404("path not found for method POST: " + h.req.Path)
		}

	case http.MethodGet:
		firstSegment, err := getFirstSegment(h.req.Path)
		if err != nil {
			return resp400(err.Error())
		}

		switch firstSegment {
		case "/health":
			return resp200Text("The function is available!")
		case "/polls":
			return h.getPollInfo()
		case "/results":
			return h.getResult()
		default:
			return resp404("path not found for method GET: " + h.req.Path)
		}

	default:
		return resp405(h.req.HTTPMethod, "OPTIONS", "GET", "POST")
	}
}
