package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bricefrisco/nameslol/shared"
	"log"
	"os"
	"strings"
)

type ErrorResponse struct {
	Message string
}

type SummonersService interface {
	Fetch(region string, name string) (*shared.SummonerDTO, error)
	Save(summoner *shared.SummonerDTO) error
}

type RegionsService interface {
	Validate(region string) bool
}

type HttpHeadersService interface {
	CreateHeaders() map[string]string
}

var summoners SummonersService
var regions RegionsService
var headers HttpHeadersService

func init() {
	log.SetFlags(0)

	var err error
	summoners, err = shared.NewSummoners(os.Getenv("DYNAMODB_TABLE"), os.Getenv("RIOT_API_TOKEN"))
	if err != nil {
		log.Fatalf("Error creating summoners: %v\n", err)
	}

	regions = shared.NewRegions()
	headers = shared.NewHttpHeaders(os.Getenv("CORS_ORIGINS"), os.Getenv("CORS_METHODS"))
}

func errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	jsonBody, err := json.Marshal(&ErrorResponse{
		Message: message,
	})
	if err != nil {
		log.Fatalf("Error marshalling error response: %v\n", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(jsonBody),
		Headers:    headers.CreateHeaders(),
	}
}

func successfulResponse(summoner *shared.SummonerDTO) events.APIGatewayProxyResponse {
	jsonBody, err := json.Marshal(summoner)
	if err != nil {
		log.Fatalf("Error marshalling successful response: %v\n", err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(jsonBody),
		Headers:    headers.CreateHeaders(),
	}
}

func HandleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := request.QueryStringParameters["name"]
	if len(name) < 3 {
		return errorResponse(400, "Query parameter 'name' must be at least 3 characters"), nil
	}

	if len(name) > 16 {
		return errorResponse(400, "Query parameter 'name' must be at most 16 characters"), nil
	}

	region := strings.ToUpper(request.QueryStringParameters["region"])
	if !regions.Validate(region) {
		return errorResponse(400, "Invalid 'region' query parameter"), nil
	}

	result, err := summoners.Fetch(region, name)
	if err != nil {
		if err.Error() == "summoner not found" {
			return errorResponse(404, "Summoner not found"), nil
		}

		log.Printf("Error fetching summoner: %v\n", err)
		return errorResponse(500, "Internal server error"), nil
	}

	err = summoners.Save(result)
	if err != nil {
		log.Printf("Error saving summoner: %v\n", err)
	} else {
		log.Printf("Successfully saved summoner: %v\n", result)
	}

	return successfulResponse(result), nil
}

func main() {
	lambda.Start(HandleRequest)
}
