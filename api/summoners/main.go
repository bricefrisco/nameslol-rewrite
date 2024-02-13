package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bricefrisco/nameslol/shared"
	"log"
	"os"
	"strconv"
	"strings"
)

type RegionsService interface {
	Validate(region string) bool
}

type HttpResponsesService interface {
	Success(responseObj any) events.APIGatewayProxyResponse
	Error(statusCode int, message string) events.APIGatewayProxyResponse
}

type SummonersService interface {
	GetByNameLength(region string, limit int32, nameLength int32, t1 int64, backwards bool) ([]*shared.SummonerDTO, error)
	GetAfter(region string, limit int32, t1 int64, backwards bool) ([]*shared.SummonerDTO, error)
}

var regions RegionsService
var responses HttpResponsesService
var summoners SummonersService

func init() {
	log.SetFlags(0)
	var err error

	regions = shared.NewRegions()
	responses = shared.NewHttpResponses(os.Getenv("CORS_ORIGINS"), os.Getenv("CORS_METHODS"))
	summoners, err = shared.NewSummoners(os.Getenv("DYNAMODB_TABLE"), os.Getenv("RIOT_API_TOKEN"))
	if err != nil {
		log.Fatalf("Error creating summoners service: %v\n", err)
	}
}

func HandleRequest(_ context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	region := strings.ToUpper(request.QueryStringParameters["region"])
	if !regions.Validate(region) {
		return responses.Error(400, "Invalid 'region' query parameter"), nil
	}

	t1, err := strconv.Atoi(request.QueryStringParameters["timestamp"])
	if err != nil {
		return responses.Error(400, "Invalid 'timestamp' query parameter"), nil
	}

	if t1 <= 0 {
		return responses.Error(400, "Invalid 'timestamp' query parameter"), nil
	}

	var nameLength int
	nameLengthStr := request.QueryStringParameters["nameLength"]
	if nameLengthStr != "" {
		nameLength, err = strconv.Atoi(nameLengthStr)
		if err != nil {
			return responses.Error(400, "Invalid 'nameLength' query parameter"), nil
		}

		if nameLength < 3 || nameLength > 16 {
			return responses.Error(400, "Invalid 'nameLength' query parameter"), nil
		}
	}

	var backwards bool
	if request.QueryStringParameters["backwards"] != "" {
		backwards, err = strconv.ParseBool(request.QueryStringParameters["backwards"])
		if err != nil {
			return responses.Error(400, "Invalid 'backwards' query parameter"), nil
		}
	}

	var response []*shared.SummonerDTO
	if nameLength == 0 {
		response, err = summoners.GetAfter(region, 35, int64(t1), backwards)
	} else {
		response, err = summoners.GetByNameLength(region, 35, int32(nameLength), int64(t1), backwards)
	}

	if err != nil {
		return responses.Error(500, "Internal server error"), nil
	}

	return responses.Success(response), nil
}

func main() {
	lambda.Start(HandleRequest)
}
