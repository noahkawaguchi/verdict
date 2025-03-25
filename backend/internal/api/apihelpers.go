package api

import (
	"maps"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// getShortPath extracts the base path without parameters for routing purposes.
// For example: `/poll/abcdefg12345` => `/poll`. If no match is found, it returns "default".
func getShortPath(longPath string) string {
	matches := regexp.MustCompile(`^(/.+)/.+$`).FindStringSubmatch(longPath)
	if len(matches) > 1 {
		return matches[1]
	}
	return "default"
}

var defaultHeaders = map[string]string{
	"Content-Type":                 "application/json",
	"Access-Control-Allow-Origin":  os.Getenv("FRONTEND_URL"),
	"Access-Control-Allow-Methods": "OPTIONS,GET,POST",
	"Access-Control-Allow-Headers": "Content-Type,Authorization",
}

// response200 creates a 200 OK HTTP response with the provided body.
func response200(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    defaultHeaders,
		Body:       body,
	}
}

// response201 creates a 201 Created HTTP response with the provided body.
func response201(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    defaultHeaders,
		Body:       body,
	}
}

// response400 creates a 400 Bad Request HTTP response with a custom error message.
func response400(errMsg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Headers:    defaultHeaders,
		Body:       `{"error":"` + errMsg + `"}`,
	}
}

// response404 creates a 404 Not Found HTTP response with a custom error message.
func response404(errMsg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Headers:    defaultHeaders,
		Body:       `{"error":"` + errMsg + `"}`,
	}
}

// response405 creates a 405 Method Not Allowed HTTP response with a custom error message and a
// custom header specifying the allowed methods.
func response405(receivedMethod string, allowedMethods ...string) events.APIGatewayProxyResponse {
	headers405 := maps.Clone(defaultHeaders)
	headers405["Allow"] = strings.Join(allowedMethods, ", ")
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusMethodNotAllowed,
		Headers:    headers405,
		Body:       `{"error":"method ` + receivedMethod + ` not allowed"}`,
	}
}

// response500 creates a 500 Internal Server Error HTTP response with a custom error message.
func response500(errMsg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Headers:    defaultHeaders,
		Body:       `{"error":"` + errMsg + `"}`,
	}
}
