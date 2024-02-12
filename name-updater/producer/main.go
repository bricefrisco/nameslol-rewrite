package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/bricefrisco/nameslol/shared"
	"log"
	"os"
	"time"
)

type Event struct {
	RefreshType string `json:"refreshType"`
}

type SQSMessage struct {
	Region string `json:"region"`
	Name   string `json:"name"`
}

type summonerService interface {
	GetBetweenDate(region string, limit int32, start int64, end int64) ([]*shared.SummonerDTO, error)
}

type regionService interface {
	GetAll() map[string]string
}

type sqsService interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

var summoners summonerService
var regions regionService
var queue sqsService
var queueUrl string

func init() {
	log.SetFlags(0)

	var err error

	summoners, err = shared.NewSummoners(os.Getenv("DYNAMODB_TABLE"), os.Getenv("RIOT_API_TOKEN"))
	if err != nil {
		log.Fatalf("could not create summoners service, %v", err)
	}

	regions = shared.NewRegions()

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("could not load AWS config, %v", err)
	}

	queue = sqs.NewFromConfig(cfg)
	queueUrl = os.Getenv("QUEUE_URL")
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

func sendToQueue(ctx context.Context, region string, name string) error {
	body := &SQSMessage{
		Region: region,
		Name:   name,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	messageBody := string(jsonBytes)

	_, err = queue.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: &messageBody,
	})

	if err != nil {
		return err
	}

	return nil
}

func HandleRequest(ctx context.Context, event *Event) error {
	start, end, err := getUpdateBetweenDates(event.RefreshType)
	if err != nil {
		return err
	}

	for region := range regions.GetAll() {
		summonersToUpdate, err := summoners.GetBetweenDate(region, 8000, start, end)
		if err != nil {
			return err
		}

		for _, s := range summonersToUpdate {
			err := sendToQueue(ctx, region, s.Name)
			if err != nil {
				return err
			}

			log.Printf("sent name: %s, region: %s to queue", s.Name, s.Region)
		}
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
