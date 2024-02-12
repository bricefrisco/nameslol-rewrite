package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bricefrisco/nameslol/shared"
	"log"
	"os"
)

type summonersService interface {
	Fetch(region string, name string) (*shared.SummonerDTO, error)
	Save(summoner *shared.SummonerDTO) error
	Delete(region string, name string) error
}

var summoners summonersService

func init() {
	log.SetFlags(0)

	var err error
	summoners, err = shared.NewSummoners(os.Getenv("DYNAMODB_TABLE"), os.Getenv("RIOT_API_TOKEN"))
	if err != nil {
		log.Fatalf("could not create summoners service, %v", err)
	}
}

type SQSMessage struct {
	Region string `json:"region"`
	Name   string `json:"name"`
}

func HandleRequest(_ context.Context, event events.SQSEvent) error {
	for _, message := range event.Records {
		var sqsMessage SQSMessage
		err := json.Unmarshal([]byte(message.Body), &sqsMessage)
		if err != nil {
			return err
		}

		summoner, err := summoners.Fetch(sqsMessage.Region, sqsMessage.Name)
		if err != nil {
			if err.Error() == "summoner not found" {
				log.Printf("summoner '%v' was not found in region '%v', deleting...", sqsMessage.Name, sqsMessage.Region)
				err = summoners.Delete(sqsMessage.Region, sqsMessage.Name)
				if err != nil {
					return err
				}
			}

			return err
		}

		err = summoners.Save(summoner)
		if err != nil {
			return err
		}

		log.Printf("summoner '%v' updated in region '%v'", summoner.Name, summoner.Region)
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
