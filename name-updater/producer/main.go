package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bricefrisco/nameslol/shared"
	"log"
	"os"
	"time"
)

type Event struct {
	RefreshType string `json:"refreshType"`
}

type SummonerService interface {
	GetBetweenDate(region string, limit int32, start int64, end int64) ([]*shared.SummonerDTO, error)
}

type RegionService interface {
	GetAll() map[string]string
}

var summoners SummonerService
var regions RegionService

func init() {
	var err error

	summoners, err = shared.NewSummoners(os.Getenv("DYNAMODB_TABLE"))
	if err != nil {
		log.Fatalf("could not create summoners service, %v", err)
	}

	regions = shared.NewRegions()
}

func getUpdateBetweenDates(refreshType string) (int64, int64, error) {
	if refreshType == "hourly" {
		return time.Now().Add(-3 * 24 * time.Hour).UnixMilli(), time.Now().Add(3 * 24 * time.Hour).UnixMilli(), nil
	} else if refreshType == "weekly" {
		return time.Now().Add(-30 * 24 * time.Hour).UnixMilli(), time.Now().Add(30 * 24 * time.Hour).UnixMilli(), nil
	} else if refreshType == "monthly" {
		return time.Now().Add(-90 * 24 * time.Hour).UnixMilli(), time.Now().Add(90 * 24 * time.Hour).UnixMilli(), nil
	}

	return 0, 0, fmt.Errorf("invalid refreshType '%s'", refreshType)
}

func HandleRequest(_ context.Context, event *Event) error {
	start, end, err := getUpdateBetweenDates(event.RefreshType)
	if err != nil {
		return err
	}

	for region := range regions.GetAll() {
		summonersToUpdate, err := summoners.GetBetweenDate(region, 10, start, end)
		if err != nil {
			return err
		}

		for _, s := range summonersToUpdate {
			log.Println("region:", s.Region, "name:", s.Name)
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
