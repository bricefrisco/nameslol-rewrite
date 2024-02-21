package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/bricefrisco/nameslol/shared"
	"testing"
	"time"
)

type MockSummonerService struct {
	ShouldFail bool
	Calls      []struct {
		Region string
		Limit  int32
		Start  int64
		End    int64
	}
}

func (m *MockSummonerService) GetBetweenDate(region string, limit int32, start int64, end int64) ([]*shared.SummonerDTO, error) {
	if m.ShouldFail {
		return nil, fmt.Errorf("error")
	}

	m.Calls = append(m.Calls, struct {
		Region string
		Limit  int32
		Start  int64
		End    int64
	}{Region: region, Limit: limit, Start: start, End: end})

	summonerDtos := make([]*shared.SummonerDTO, 1)
	summonerDtos[0] = &shared.SummonerDTO{
		Name:             "Testing",
		Region:           "NA",
		AccountID:        "123",
		RevisionDate:     123,
		AvailabilityDate: 123,
		Level:            123,
		LastUpdated:      123,
		SummonerIcon:     123,
	}

	return summonerDtos, nil
}

type MockRegionService struct {
	NumCalls int
}

func (m *MockRegionService) GetAll() map[string]string {
	m.NumCalls++
	return map[string]string{
		"NA":  "na1",
		"EUW": "euw1",
	}
}

type MockSQSService struct {
	ShouldFail bool
	Calls      []struct {
		Ctx    context.Context
		params *sqs.SendMessageInput
	}
}

func (m *MockSQSService) SendMessage(ctx context.Context, params *sqs.SendMessageInput, _ ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	if m.ShouldFail {
		return nil, fmt.Errorf("error")
	}

	m.Calls = append(m.Calls, struct {
		Ctx    context.Context
		params *sqs.SendMessageInput
	}{Ctx: ctx, params: params})

	return &sqs.SendMessageOutput{}, nil
}

func setup() {
	summoners = &MockSummonerService{}
	regions = &MockRegionService{}
	queue = &MockSQSService{}
	queueUrl = "test.queue.url"
}

func TestHandleRequest_CorrectStartAndEndDateHourly(t *testing.T) {
	setup()

	event := &Event{RefreshType: "hourly"}
	err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockSummoners := summoners.(*MockSummonerService)

	leeway := 1000
	start := time.Now().Add(-3 * 24 * time.Hour).UnixMilli()
	end := time.Now().Add(3 * 24 * time.Hour).UnixMilli()

	if mockSummoners.Calls[0].Start < start-int64(leeway) || mockSummoners.Calls[0].Start > start+int64(leeway) {
		t.Errorf("expected start to be within 3 days ago, got %d", mockSummoners.Calls[0].Start)
	}

	if mockSummoners.Calls[0].End < end-int64(leeway) || mockSummoners.Calls[0].End > end+int64(leeway) {
		t.Errorf("expected end to be within 3 days from now, got %d", mockSummoners.Calls[0].End)
	}
}

func TestHandleRequest_CorrectStartAndEndDateWeekly(t *testing.T) {
	setup()

	event := &Event{RefreshType: "weekly"}
	err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockSummoners := summoners.(*MockSummonerService)

	leeway := 1000
	start := time.Now().Add(-30 * 24 * time.Hour).UnixMilli()
	end := time.Now().Add(30 * 24 * time.Hour).UnixMilli()

	if mockSummoners.Calls[0].Start < start-int64(leeway) || mockSummoners.Calls[0].Start > start+int64(leeway) {
		t.Errorf("expected start to be within 3 days ago, got %d", mockSummoners.Calls[0].Start)
	}

	if mockSummoners.Calls[0].End < end-int64(leeway) || mockSummoners.Calls[0].End > end+int64(leeway) {
		t.Errorf("expected end to be within 3 days from now, got %d", mockSummoners.Calls[0].End)
	}
}

func TestHandleRequest_CorrectStartAndEndDateMonthly(t *testing.T) {
	setup()

	event := &Event{RefreshType: "monthly"}
	err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockSummoners := summoners.(*MockSummonerService)

	leeway := 1000
	start := time.Now().Add(-90 * 24 * time.Hour).UnixMilli()
	end := time.Now().Add(90 * 24 * time.Hour).UnixMilli()

	if mockSummoners.Calls[0].Start < start-int64(leeway) || mockSummoners.Calls[0].Start > start+int64(leeway) {
		t.Errorf("expected start to be within 3 days ago, got %d", mockSummoners.Calls[0].Start)
	}

	if mockSummoners.Calls[0].End < end-int64(leeway) || mockSummoners.Calls[0].End > end+int64(leeway) {
		t.Errorf("expected end to be within 3 days from now, got %d", mockSummoners.Calls[0].End)
	}
}

func TestHandleRequest_ReturnsErrorForInvalidRefreshType(t *testing.T) {
	setup()

	event := &Event{RefreshType: "invalid"}
	err := HandleRequest(context.TODO(), event)
	if err == nil {
		t.Error("expected error for invalid refresh type")
	}
}

func TestHandleRequest_ValidNumberOfRegionCalls(t *testing.T) {
	setup()

	event := &Event{RefreshType: "hourly"}
	err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockRegions := regions.(*MockRegionService)
	if mockRegions.NumCalls != 1 {
		t.Errorf("expected 1 call to GetAll, got %d", mockRegions.NumCalls)
	}
}

func TestHandleRequest_CallsSummonersWithCorrectRegion(t *testing.T) {
	setup()

	event := &Event{RefreshType: "hourly"}
	err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockSummoners := summoners.(*MockSummonerService)

	// map doesn't retain order, so we need to check both calls
	expectedRegions := make(map[string]bool)
	expectedRegions["NA"] = true
	expectedRegions["EUW"] = true

	for i, call := range mockSummoners.Calls {
		if !expectedRegions[call.Region] {
			t.Errorf("expected call %d to GetBetweenDate to be for NA or EUW, got %s", i, call.Region)
		}
	}
}

func TestHandleRequest_CallsSummonersWithCorrectLimit(t *testing.T) {
	setup()

	event := &Event{RefreshType: "hourly"}
	err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockSummoners := summoners.(*MockSummonerService)
	if mockSummoners.Calls[0].Limit != 8000 {
		t.Errorf("expected first call to GetBetweenDate to have limit 8000, got %d", mockSummoners.Calls[0].Limit)
	}
}

func TestHandleRequest_ReturnsErrorWhenSummonersReturnsError(t *testing.T) {
	setup()

	summoners = &MockSummonerService{
		ShouldFail: true,
	}

	event := &Event{RefreshType: "hourly"}
	err := HandleRequest(context.TODO(), event)
	if err == nil {
		t.Error("expected error from HandleRequest")
	}
}

func TestHandleRequest_ValidNumberOfSummonerCalls(t *testing.T) {
	setup()

	event := &Event{RefreshType: "hourly"}
	err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockSummoners := summoners.(*MockSummonerService)
	if len(mockSummoners.Calls) != 2 {
		t.Errorf("expected 2 calls to GetBetweenDate, got %d", len(mockSummoners.Calls))
	}
}

func TestHandleRequest_CallsSendToQueueCorrectNumberOfTimes(t *testing.T) {
	setup()

	event := &Event{RefreshType: "hourly"}
	err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockQueue := queue.(*MockSQSService)
	if len(mockQueue.Calls) != 2 {
		t.Errorf("expected 2 calls to SendMessage, got %d", len(mockQueue.Calls))
	}
}

func TestHandleRequest_CallsSendToQueueWithCorrectQueueUrl(t *testing.T) {
	setup()

	event := &Event{RefreshType: "hourly"}
	err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	mockQueue := queue.(*MockSQSService)
	if *mockQueue.Calls[0].params.QueueUrl != queueUrl {
		t.Errorf("expected first call to SendMessage to have QueueUrl %s, got %s", queueUrl, *mockQueue.Calls[0].params.QueueUrl)
	}
}

func TestHandleRequest_CallsSendToQueueWithCorrectBody(t *testing.T) {
	setup()

	event := &Event{RefreshType: "hourly"}
	err := HandleRequest(context.TODO(), event)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// map doesn't retain order, so we need to check both calls
	mockQueue := queue.(*MockSQSService)
	body := &SQSMessage{
		Region: "NA",
		Name:   "Testing",
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		t.Errorf("error marshalling body: %v", err)
	}

	match := false
	for _, calls := range mockQueue.Calls {
		if *calls.params.MessageBody == string(jsonBytes) {
			match = true
		}
	}

	if !match {
		t.Errorf("expected call to SendMessage to have MessageBody %s", string(jsonBytes))
	}
}

func TestHandleRequest_ReturnsErrorWhenSendToQueueReturnsError(t *testing.T) {
	setup()

	queue = &MockSQSService{
		ShouldFail: true,
	}

	event := &Event{RefreshType: "hourly"}
	err := HandleRequest(context.TODO(), event)
	if err == nil {
		t.Error("expected error from HandleRequest")
	}
}
