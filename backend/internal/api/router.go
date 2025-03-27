package api

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

type datastore interface {
	PutPoll(poll *models.Poll) error
	GetPoll(pollID string) (*models.Poll, error)
	PutBallot(ballot *models.Ballot) error
	GetBallots(pollID string) ([]*models.Ballot, error)
}

type handler struct {
	store datastore
	req   events.APIGatewayProxyRequest
}

func NewHandler(store datastore, req events.APIGatewayProxyRequest) *handler {
	return &handler{store, req}
}

// Route matches the method and path of the request and calls the relevant method.
func (h *handler) Route() events.APIGatewayProxyResponse {
	switch h.req.HTTPMethod {
	case http.MethodPost:
		switch h.req.Path {
		case "/poll":
			return h.createPoll()
		case "/ballot":
			return h.castBallot()
		default:
			return response404("path not found for method POST: " + h.req.Path)
		}
	case http.MethodGet:
		switch getShortPath(h.req.Path) {
		case "/poll":
			return h.getPollInfo()
		case "/result":
			return h.getResult()
		default:
			return response404("path not found for method GET: " + h.req.Path)
		}
	default:
		return response405(h.req.HTTPMethod, "OPTIONS", "GET", "POST")
	}
}
