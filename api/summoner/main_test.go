package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/bricefrisco/nameslol/shared"
	"strings"
	"testing"
)

type SummonersServiceMock struct {
	ShouldFetchFail      bool
	ShouldSaveFail       bool
	ShouldReturnNotFound bool
	Calls                []struct {
		Region string
		Name   string
	}
}

type RegionsServiceMock struct {
	ShouldFail bool
	Calls      []string
}

var summonerDto *shared.SummonerDTO
var corsOrigins = "test-origin"
var corsMethods = "test-methods"
var expectedHeaders = map[string]string{
	"Access-Control-Allow-Origin":  corsOrigins,
	"Access-Control-Allow-Methods": corsMethods,
	"Access-Control-Allow-Headers": "*",
	"Content-Type":                 "application/json",
}

func setup() {
	summonerDto = &shared.SummonerDTO{
		Name:             "Test",
		Region:           "NA",
		AccountID:        "test-aid",
		RevisionDate:     1234567890,
		AvailabilityDate: 1234567891,
		Level:            30,
		LastUpdated:      1234567892,
		SummonerIcon:     123,
	}

	summoners = &SummonersServiceMock{}
	regions = &RegionsServiceMock{}
	responses = shared.NewHttpResponses(corsOrigins, corsMethods)
}

func (s *SummonersServiceMock) Fetch(region string, name string) (*shared.SummonerDTO, error) {
	s.Calls = append(s.Calls, struct {
		Region string
		Name   string
	}{Region: region, Name: name})

	if s.ShouldFetchFail {
		return nil, fmt.Errorf("error")
	}

	if s.ShouldReturnNotFound {
		return nil, fmt.Errorf("summoner not found")
	}

	return summonerDto, nil
}

func (s *SummonersServiceMock) Save(_ *shared.SummonerDTO) error {
	if s.ShouldSaveFail {
		return fmt.Errorf("error")
	}

	return nil
}

func (s *RegionsServiceMock) Validate(region string) bool {
	s.Calls = append(s.Calls, region)

	return !s.ShouldFail
}

func TestHandleRequest_ValidatesNameTooShort(t *testing.T) {
	setup()

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"region": "NA",
			"name":   "T",
		},
	}

	res, err := HandleRequest(nil, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", res.StatusCode)
	}

	if !strings.Contains(res.Body, "must be at least 3") {
		t.Errorf("Expected body to contain 'must be at least 3', got %s", res.Body)
	}
}

func TestHandleRequest_ValidatesNameTooLong(t *testing.T) {
	setup()

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"region": "NA",
			"name":   "TestTestTestTestT",
		},
	}

	res, err := HandleRequest(nil, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", res.StatusCode)
	}

	if !strings.Contains(res.Body, "must be at most 16") {
		t.Errorf("Expected body to contain 'must be at most 16', got %s", res.Body)
	}
}

func TestHandleRequest_ValidatesRegion(t *testing.T) {
	setup()

	regions = &RegionsServiceMock{ShouldFail: true}

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"region": "NA",
			"name":   "Test",
		},
	}

	res, err := HandleRequest(nil, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if res.StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", res.StatusCode)
	}

	if !strings.Contains(res.Body, "Invalid 'region'") {
		t.Errorf("Expected body to contain 'Invalid 'region'', got %s", res.Body)
	}

	regionMock := regions.(*RegionsServiceMock)

	if len(regionMock.Calls) != 1 {
		t.Errorf("Expected 1 call to Validate, got %d", len(regionMock.Calls))
	}

	if regionMock.Calls[0] != "NA" {
		t.Errorf("Expected call to Validate with region 'NA', got %s", regionMock.Calls[0])
	}
}

func TestHandleRequest_CallsFetchWithCorrectParameters(t *testing.T) {
	setup()

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"region": "NA",
			"name":   "Test",
		},
	}

	mockSummoners := summoners.(*SummonersServiceMock)

	_, err := HandleRequest(nil, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(mockSummoners.Calls) != 1 {
		t.Errorf("Expected 1 call to Fetch, got %d", len(mockSummoners.Calls))
	}

	if mockSummoners.Calls[0].Region != "NA" {
		t.Errorf("Expected call to Fetch with region 'NA', got %s", mockSummoners.Calls[0].Region)
	}

	if mockSummoners.Calls[0].Name != "Test" {
		t.Errorf("Expected call to Fetch with name 'Test', got %s", mockSummoners.Calls[0].Name)
	}
}

func TestHandleRequest_ReturnsErrorOnFetchError(t *testing.T) {
	setup()

	mockSummoners := summoners.(*SummonersServiceMock)
	mockSummoners.ShouldFetchFail = true

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"region": "NA",
			"name":   "Test",
		},
	}

	res, err := HandleRequest(nil, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if res.StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", res.StatusCode)
	}

	if !strings.Contains(res.Body, "Internal server error") {
		t.Errorf("Expected body to contain 'Internal server error', got %s", res.Body)
	}
}

func TestHandleRequest_Returns404WhenSummonerNotFound(t *testing.T) {
	setup()

	mockSummoners := summoners.(*SummonersServiceMock)
	mockSummoners.ShouldReturnNotFound = true

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"region": "NA",
			"name":   "Test",
		},
	}

	res, err := HandleRequest(nil, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if res.StatusCode != 404 {
		t.Errorf("Expected status code 404, got %d", res.StatusCode)
	}

	if !strings.Contains(res.Body, "Summoner not found") {
		t.Errorf("Expected body to contain 'Summoner not found', got %s", res.Body)
	}
}

func TestHandleRequest_ReturnsSummonerOnSuccess(t *testing.T) {
	setup()

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"region": "NA",
			"name":   "Test",
		},
	}

	res, err := HandleRequest(nil, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}

	jsonBody, err := json.Marshal(summonerDto)
	if err != nil {
		t.Errorf("Error marshalling successful response: %v\n", err)
	}

	if res.Body != string(jsonBody) {
		t.Errorf("Expected body to contain summoner, got %s", res.Body)
	}
}

func TestHandleRequest_Returns200EvenWhenSaveFails(t *testing.T) {
	setup()

	mockSummoners := summoners.(*SummonersServiceMock)
	mockSummoners.ShouldSaveFail = true

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"region": "NA",
			"name":   "Test",
		},
	}

	res, err := HandleRequest(nil, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if res.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
}

func TestHandleRequest_HasCorrectHeadersOnSuccess(t *testing.T) {
	setup()

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"region": "NA",
			"name":   "Test",
		},
	}

	res, err := HandleRequest(nil, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	for k, v := range expectedHeaders {
		if res.Headers[k] != v {
			t.Errorf("Expected header %s to be %s, got %s", k, v, res.Headers[k])
		}
	}
}

func TestHandleRequest_HasCorrectHeadersOnFailure(t *testing.T) {
	setup()

	mockSummoners := summoners.(*SummonersServiceMock)
	mockSummoners.ShouldFetchFail = true

	request := events.APIGatewayProxyRequest{
		QueryStringParameters: map[string]string{
			"region": "NA",
			"name":   "Test",
		},
	}

	res, err := HandleRequest(nil, request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	for k, v := range expectedHeaders {
		if res.Headers[k] != v {
			t.Errorf("Expected header %s to be %s, got %s", k, v, res.Headers[k])
		}
	}
}
