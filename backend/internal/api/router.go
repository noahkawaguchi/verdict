package api

import (
	"context"
	"net/http"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

// Router directs API Gateway requests to the correct handler.
func Router(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	// Define a reusable path not found response
	pathNotFound := events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Body:       `{"error": "Path not found"}`,
	}
	// Match the method and path
	switch request.HTTPMethod {
	case http.MethodPost:
		switch request.Path {
		case "/poll":
			return createPollHandler(ctx, request)
		default:
			return pathNotFound, nil
		}
	case http.MethodGet:
		if matched, _ := regexp.MatchString("^/poll/.*$", request.Path); matched {
			return createBallotHandler(ctx, request)
		} else {
			return pathNotFound, nil
		}
	default:
		return pathNotFound, nil
	}
}
