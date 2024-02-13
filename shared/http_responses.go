package shared

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
)

type HttpResponses struct {
	corsOrigins string
	corsMethods string
}

type ErrResponse struct {
	Message string `json:"message"`
}

func NewHttpResponses(corsOrigins string, corsMethods string) *HttpResponses {
	return &HttpResponses{corsOrigins: corsOrigins, corsMethods: corsMethods}
}

func (h *HttpResponses) Success(responseObj any) events.APIGatewayProxyResponse {
	jsonBody, err := json.Marshal(responseObj)
	if err != nil {
		log.Fatalf("Error marshalling successful response: %v\n", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonBody),
		Headers:    h.createHeaders(),
	}
}

func (h *HttpResponses) Error(statusCode int, message string) events.APIGatewayProxyResponse {
	jsonBody, err := json.Marshal(&ErrResponse{
		Message: message,
	})
	if err != nil {
		log.Fatalf("Error marshalling error response: %v\n", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(jsonBody),
		Headers:    h.createHeaders(),
	}
}

func (h *HttpResponses) createHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = h.corsOrigins
	headers["Access-Control-Allow-Methods"] = h.corsMethods
	headers["Access-Control-Allow-Headers"] = "*"
	headers["Content-Type"] = "application/json"
	return headers
}
