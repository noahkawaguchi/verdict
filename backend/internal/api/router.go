package api

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

// Router directs API Gateway requests to the correct handler.
func Router(
	ctx context.Context,
	request events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	switch request.Path {
	case "/poll/create":
		return createPollHandler(ctx, request)
	default:
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       `{"error": "Path not found"}`,
		}, nil
	}
}
