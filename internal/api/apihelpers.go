package api

import (
	"errors"
	"maps"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// getFirstSegment extracts the first path segment for routing purposes.
// For example: `/poll/abcdefg12345` => `/poll`, `/health` => `/health`.
// It returns an error if the path is malformed, i.e. doesn't have at least one `/`.
func getFirstSegment(fullPath string) (string, error) {
	if parts := strings.SplitN(fullPath, "/", 3); len(parts) >= 2 {
		return "/" + parts[1], nil
	}
	return "", errors.New("malformed path")
}

var defaultHeaders = map[string]string{
	"Content-Type":                 "application/json",
	"Access-Control-Allow-Origin":  os.Getenv("FRONTEND_URL"),
	"Access-Control-Allow-Methods": "OPTIONS,GET,POST",
	"Access-Control-Allow-Headers": "Content-Type,Authorization",
}

// resp200 creates a 200 OK HTTP response with the provided body.
func resp200(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    defaultHeaders,
		Body:       body,
	}
}

// resp201 creates a 201 Created HTTP response with the provided body.
func resp201(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    defaultHeaders,
		Body:       body,
	}
}

// resp400 creates a 400 Bad Request HTTP response with a custom error message.
func resp400(errMsg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Headers:    defaultHeaders,
		Body:       `{"error":"` + errMsg + `"}`,
	}
}

// resp404 creates a 404 Not Found HTTP response with a custom error message.
func resp404(errMsg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Headers:    defaultHeaders,
		Body:       `{"error":"` + errMsg + `"}`,
	}
}

// resp405 creates a 405 Method Not Allowed HTTP response with a custom error message and a
// custom header specifying the allowed methods.
func resp405(receivedMethod string, allowedMethods ...string) events.APIGatewayProxyResponse {
	headers405 := maps.Clone(defaultHeaders)
	headers405["Allow"] = strings.Join(allowedMethods, ", ")
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Headers:    headers405,
		Body:       `{"error":"method ` + receivedMethod + ` not allowed"}`,
	}
}

// resp500 is a 500 Internal Server Error HTTP response with an "internal server error" error
// message.
var resp500 events.APIGatewayProxyResponse = events.APIGatewayProxyResponse{
	StatusCode: http.StatusInternalServerError,
	Headers:    defaultHeaders,
	Body:       `{"error":"internal server error"}`,
}
