package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bricefrisco/nameslol/shared"
	"log"
	"os"
	"strings"
)

type SummonersService interface {
	Fetch(region string, name string) (*shared.SummonerDTO, error)
	Save(summoner *shared.SummonerDTO) error
}

type RegionsService interface {
	Validate(region string) bool
}

type HttpResponsesService interface {
	Success(responseObj any) events.APIGatewayProxyResponse
	Error(statusCode int, message string) events.APIGatewayProxyResponse
}

var summoners SummonersService
var regions RegionsService
var responses HttpResponsesService

func init() {
	log.SetFlags(0)

	var err error
	summoners, err = shared.NewSummoners(os.Getenv("DYNAMODB_TABLE"), os.Getenv("RIOT_API_TOKEN"))
	if err != nil {
		log.Fatalf("Error creating summoners: %v\n", err)
	}

	regions = shared.NewRegions()
	responses = shared.NewHttpResponses(os.Getenv("CORS_ORIGINS"), os.Getenv("CORS_METHODS"))
}

func HandleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "OPTIONS" {
		return responses.Success(nil), nil
	}

	if request.HTTPMethod != "GET" {
		return responses.Error(405, "Method not allowed"), nil
	}

	name := request.QueryStringParameters["name"]
	if len(name) < 3 {
		return responses.Error(400, "Query parameter 'name' must be at least 3 characters"), nil
	}

	if len(name) > 16 {
		return responses.Error(400, "Query parameter 'name' must be at most 16 characters"), nil
	}

	region := strings.ToUpper(request.QueryStringParameters["region"])
	if !regions.Validate(region) {
		return responses.Error(400, "Invalid 'region' query parameter"), nil
	}

	result, err := summoners.Fetch(region, name)
	if err != nil {
		if err.Error() == "summoner not found" {
			return responses.Error(404, "Summoner not found"), nil
		}

		log.Printf("Error fetching summoner: %v\n", err)
		return responses.Error(500, "Internal server error"), nil
	}

	err = summoners.Save(result)
	if err != nil {
		log.Printf("Error saving summoner: %v\n", err)
	} else {
		log.Printf("Successfully saved summoner: %v\n", result)
	}

	return responses.Success(result), nil
}

func main() {
	lambda.Start(HandleRequest)
}
