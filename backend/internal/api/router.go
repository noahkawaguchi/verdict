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
	// Match the method and path
	switch request.HTTPMethod {
	case http.MethodPost:
		switch request.Path {
		case "/poll":
			return createPollHandler(ctx, request)
		case "/ballot":
			return castBallotHandler(ctx, request)
		default:
			return response404, nil
		}
	case http.MethodGet:
		if matched, _ := regexp.MatchString("^/poll/.*$", request.Path); matched {
			return createBallotHandler(ctx, request)
		} else {
			return response404, nil
		}
	default:
		return response404, nil
	}
}
