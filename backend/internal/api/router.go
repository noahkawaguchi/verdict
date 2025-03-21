package api

import (
	"context"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/noahkawaguchi/verdict/backend/internal/models"
)

type datastore interface {
	PutPoll(ctx context.Context, poll *models.Poll) error
	GetPoll(ctx context.Context, pollID string) (*models.Poll, error)
	PutBallot(ctx context.Context, ballot *models.Ballot) error
	GetPollWithBallots(ctx context.Context, pollID string) (*models.Poll, []*models.Ballot, error)
}

type handler struct {
	ctx context.Context
	req events.APIGatewayProxyRequest
	ds  datastore
}

// Router returns a function that creates a handler to handle the request.
func Router(ds datastore) func(
	ctx context.Context, request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	return func(
		ctx context.Context, request events.APIGatewayProxyRequest,
	) (events.APIGatewayProxyResponse, error) {
		h := &handler{ctx, request, ds}
		return h.route(), nil
	}
}

// route matches the method and path of the request and calls the relevant method.
func (h *handler) route() events.APIGatewayProxyResponse {
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
		if matched, _ := regexp.MatchString("^/poll/.*$", h.req.Path); matched {
			return h.createBallot()
		} else if matched, _ := regexp.MatchString("^/result/.*$", h.req.Path); matched {
			return h.getResult()
		} else {
			return response404("path not found for method GET: " + h.req.Path)
		}
	default:
		return response405(h.req.HTTPMethod, "OPTIONS", "GET", "POST")
	}
}
