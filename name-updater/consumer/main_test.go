package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/bricefrisco/nameslol/shared"
	"testing"
	"time"
)

type SummonersServiceMock struct {
	ShouldFail       bool
	SummonerNotFound bool
	FetchCalls       []struct {
		Region string
		Name   string
	}
	SaveCalls []struct {
		Summoner *shared.SummonerDTO
	}
	DeleteCalls []struct {
		Region string
		Name   string
	}
}

func (s *SummonersServiceMock) Fetch(region string, name string) (*shared.SummonerDTO, error) {
	s.FetchCalls = append(s.FetchCalls, struct {
		Region string
		Name   string
	}{region, name})

	if s.ShouldFail {
		return nil, fmt.Errorf("error")
	}

	if s.SummonerNotFound {
		return nil, fmt.Errorf("summoner not found")
	}

	return summonerDto, nil
}

func (s *SummonersServiceMock) Save(summoner *shared.SummonerDTO) error {
	s.SaveCalls = append(s.SaveCalls, struct {
		Summoner *shared.SummonerDTO
	}{summoner})

	if s.ShouldFail {
		return fmt.Errorf("error")
	}

	return nil
}

func (s *SummonersServiceMock) Delete(region string, name string) error {
	s.DeleteCalls = append(s.DeleteCalls, struct {
		Region string
		Name   string
	}{region, name})

	if s.ShouldFail {
		return fmt.Errorf("error")
	}

	return nil
}

var summonerDto *shared.SummonerDTO

func setup() {
	summonerDto = &shared.SummonerDTO{
		Name:             "test",
		Region:           "NA",
		AccountID:        "test-aid",
		RevisionDate:     time.Date(2023, time.February, 12, 1, 2, 3, 0, time.UTC).UnixMilli(),
		AvailabilityDate: time.Date(2023, time.August, 12, 0, 2, 3, 0, time.UTC).UnixMilli(),
		Level:            5,
		LastUpdated:      time.Now().UnixMilli(),
		SummonerIcon:     1,
	}

	summoners = &SummonersServiceMock{}
}

func TestHandleRequest_CallsFetchWithCorrectParams(t *testing.T) {
	setup()

	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: `{"region":"NA","name":"test"}`,
			},
		},
	}

	err := HandleRequest(context.Background(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if summoners.(*SummonersServiceMock).FetchCalls[0].Region != "NA" {
		t.Errorf("expected region to be 'NA', got %v", summoners.(*SummonersServiceMock).FetchCalls[0].Region)
	}

	if summoners.(*SummonersServiceMock).FetchCalls[0].Name != "test" {
		t.Errorf("expected name to be 'test', got %v", summoners.(*SummonersServiceMock).FetchCalls[0].Name)
	}
}

func TestHandleRequest_CallsSaveWithCorrectParams(t *testing.T) {
	setup()

	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: `{"region":"NA","name":"test"}`,
			},
		},
	}

	err := HandleRequest(context.Background(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if summoners.(*SummonersServiceMock).SaveCalls[0].Summoner != summonerDto {
		t.Errorf("expected summoner to be %v, got %v", summonerDto, summoners.(*SummonersServiceMock).SaveCalls[0].Summoner)
	}
}

func TestHandleRequest_CallsDeleteWithCorrectParams(t *testing.T) {
	summoners = &SummonersServiceMock{
		SummonerNotFound: true,
	}

	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: `{"region":"NA","name":"test"}`,
			},
		},
	}

	err := HandleRequest(context.Background(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(summoners.(*SummonersServiceMock).DeleteCalls) == 0 {
		t.Errorf("expected delete to be called")
	}

	if summoners.(*SummonersServiceMock).DeleteCalls[0].Region != "NA" {
		t.Errorf("expected region to be 'NA', got %v", summoners.(*SummonersServiceMock).DeleteCalls[0].Region)
	}

	if summoners.(*SummonersServiceMock).DeleteCalls[0].Name != "test" {
		t.Errorf("expected name to be 'test', got %v", summoners.(*SummonersServiceMock).DeleteCalls[0].Name)
	}
}

func TestHandleRequest_IteratesOverAllMessages(t *testing.T) {
	setup()

	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: `{"region":"NA","name":"test"}`,
			},
			{
				Body: `{"region":"EUW","name":"test"}`,
			},
		},
	}

	err := HandleRequest(context.Background(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(summoners.(*SummonersServiceMock).FetchCalls) != 2 {
		t.Errorf("expected 2 fetch calls, got %v", len(summoners.(*SummonersServiceMock).FetchCalls))
	}
}

func TestHandleRequest_ReturnsError_WhenFail(t *testing.T) {
	summoners = &SummonersServiceMock{
		ShouldFail: true,
	}

	event := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: `{"region":"NA","name":"test"}`,
			},
		},
	}

	err := HandleRequest(context.Background(), event)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
