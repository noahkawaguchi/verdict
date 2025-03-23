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

type Handler struct {
	Store datastore
	Req   events.APIGatewayProxyRequest
}

// Route matches the method and path of the request and calls the relevant method.
func (h *Handler) Route() events.APIGatewayProxyResponse {
	switch h.Req.HTTPMethod {
	case http.MethodPost:
		switch h.Req.Path {
		case "/poll":
			return h.createPoll()
		case "/ballot":
			return h.castBallot()
		default:
			return response404("path not found for method POST: " + h.Req.Path)
		}
	case http.MethodGet:
		switch getShortPath(h.Req.Path) {
		case "/poll":
			return h.getPollInfo()
		case "/result":
			return h.getResult()
		default:
			return response404("path not found for method GET: " + h.Req.Path)
		}
	default:
		return response405(h.Req.HTTPMethod, "OPTIONS", "GET", "POST")
	}
}
