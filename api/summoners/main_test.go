package main

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/bricefrisco/nameslol/shared"
	"strings"
	"testing"
)

type RegionMock struct {
	ValidateCalls []string
	IsValid       bool
}

type HttpResponsesMock struct {
	SuccessCalls []any
	ErrorCalls   []struct {
		StatusCode int
		Message    string
	}
}

type SummonersMock struct {
	GetByNameLengthCalls []struct {
		Region     string
		Limit      int32
		NameLength int32
		T1         int64
		Backwards  bool
	}
	GetAfterCalls []struct {
		Region    string
		Limit     int32
		T1        int64
		Backwards bool
	}
	ReturnError bool
}

func (r *RegionMock) Validate(region string) bool {
	r.ValidateCalls = append(r.ValidateCalls, region)
	return r.IsValid
}

func (h *HttpResponsesMock) Success(responseObj any) events.APIGatewayProxyResponse {
	h.SuccessCalls = append(h.SuccessCalls, responseObj)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}
}

func (h *HttpResponsesMock) Error(statusCode int, message string) events.APIGatewayProxyResponse {
	h.ErrorCalls = append(h.ErrorCalls, struct {
		StatusCode int
		Message    string
	}{statusCode, message})

	return events.APIGatewayProxyResponse{
		StatusCode: 500,
	}
}

func (s *SummonersMock) GetByNameLength(region string, limit int32, nameLength int32, t1 int64, backwards bool) ([]*shared.SummonerDTO, error) {
	s.GetByNameLengthCalls = append(s.GetByNameLengthCalls, struct {
		Region     string
		Limit      int32
		NameLength int32
		T1         int64
		Backwards  bool
	}{region, limit, nameLength, t1, backwards})

	if s.ReturnError {
		return nil, errors.New("error")
	}

	return []*shared.SummonerDTO{}, nil
}

func (s *SummonersMock) GetAfter(region string, limit int32, t1 int64, backwards bool) ([]*shared.SummonerDTO, error) {
	s.GetAfterCalls = append(s.GetAfterCalls, struct {
		Region    string
		Limit     int32
		T1        int64
		Backwards bool
	}{region, limit, t1, backwards})

	if s.ReturnError {
		return nil, errors.New("error")
	}

	return []*shared.SummonerDTO{}, nil
}

func setup() {
	regions = &RegionMock{IsValid: true}
	responses = &HttpResponsesMock{}
	summoners = &SummonersMock{}
}

func TestHandleRequest_Returns400ErrorWhenRegionIsInvalid(t *testing.T) {
	setup()
	regions.(*RegionMock).IsValid = false

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region": "invalid",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(regions.(*RegionMock).ValidateCalls) != 1 {
		t.Errorf("Expected 1 call to Validate, got %d", len(regions.(*RegionMock).ValidateCalls))
	}

	if regions.(*RegionMock).ValidateCalls[0] != "INVALID" {
		t.Errorf("Expected 'INVALID' to be passed to Validate, got %s", regions.(*RegionMock).ValidateCalls[0])
	}

	if len(responses.(*HttpResponsesMock).ErrorCalls) != 1 {
		t.Errorf("Expected 1 error response, got %d", len(responses.(*HttpResponsesMock).ErrorCalls))
	}

	if responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode)
	}

	if !strings.Contains(responses.(*HttpResponsesMock).ErrorCalls[0].Message, "Invalid 'region'") {
		t.Errorf("Expected message to contain \"Invalid 'region'\", got %s", responses.(*HttpResponsesMock).ErrorCalls[0].Message)
	}
}

func TestHandleRequest_Returns400ErrorWhenTimestampIsInvalid(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":    "na",
			"timestamp": "invalid",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(responses.(*HttpResponsesMock).ErrorCalls) != 1 {
		t.Errorf("Expected 1 error response, got %d", len(responses.(*HttpResponsesMock).ErrorCalls))
	}

	if responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode)
	}

	if !strings.Contains(responses.(*HttpResponsesMock).ErrorCalls[0].Message, "Invalid 'timestamp'") {
		t.Errorf("Expected message 'Invalid 'timestamp' query parameter', got %s", responses.(*HttpResponsesMock).ErrorCalls[0].Message)
	}
}

func TestHandleRequest_Returns400ErrorWhenTimestampIsZero(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":    "na",
			"timestamp": "0",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(responses.(*HttpResponsesMock).ErrorCalls) != 1 {
		t.Errorf("Expected 1 error response, got %d", len(responses.(*HttpResponsesMock).ErrorCalls))
	}

	if responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode)
	}

	if !strings.Contains(responses.(*HttpResponsesMock).ErrorCalls[0].Message, "Invalid 'timestamp'") {
		t.Errorf("Expected message 'Invalid 'timestamp' query parameter', got %s", responses.(*HttpResponsesMock).ErrorCalls[0].Message)
	}
}

func TestHandleRequest_Returns400ErrorWhenNameLengthIsInvalid(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":     "na",
			"timestamp":  "1",
			"nameLength": "invalid",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(responses.(*HttpResponsesMock).ErrorCalls) != 1 {
		t.Errorf("Expected 1 error response, got %d", len(responses.(*HttpResponsesMock).ErrorCalls))
	}

	if responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode)
	}

	if !strings.Contains(responses.(*HttpResponsesMock).ErrorCalls[0].Message, "Invalid 'nameLength'") {
		t.Errorf("Expected message 'Invalid 'nameLength' query parameter', got %s", responses.(*HttpResponsesMock).ErrorCalls[0].Message)
	}
}

func TestHandleRequest_Returns400ErrorWhenNameLengthIsTooLong(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":     "na",
			"timestamp":  "1",
			"nameLength": "17",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(responses.(*HttpResponsesMock).ErrorCalls) != 1 {
		t.Errorf("Expected 1 error response, got %d", len(responses.(*HttpResponsesMock).ErrorCalls))
	}

	if responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode)
	}

	if !strings.Contains(responses.(*HttpResponsesMock).ErrorCalls[0].Message, "Invalid 'nameLength'") {
		t.Errorf("Expected message 'Invalid 'nameLength' query parameter', got %s", responses.(*HttpResponsesMock).ErrorCalls[0].Message)
	}
}

func TestHandleRequest_Returns400ErrorWhenBackwardsIsInvalid(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":    "na",
			"timestamp": "1",
			"backwards": "invalid",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(responses.(*HttpResponsesMock).ErrorCalls) != 1 {
		t.Errorf("Expected 1 error response, got %d", len(responses.(*HttpResponsesMock).ErrorCalls))
	}

	if responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode != 400 {
		t.Errorf("Expected status code 400, got %d", responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode)
	}

	if !strings.Contains(responses.(*HttpResponsesMock).ErrorCalls[0].Message, "Invalid 'backwards'") {
		t.Errorf("Expected message 'Invalid 'backwards' query parameter', got %s", responses.(*HttpResponsesMock).ErrorCalls[0].Message)
	}
}

func TestHandleRequest_CallsGetAfterWhenNameLengthIsNotPassed(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":    "na",
			"timestamp": "1",
			"backwards": "false",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(summoners.(*SummonersMock).GetAfterCalls) != 1 {
		t.Errorf("Expected 1 call to GetAfter, got %d", len(summoners.(*SummonersMock).GetAfterCalls))
	}
}

func TestHandleRequest_CallsGetByNameLengthWhenNameLengthIsPassed(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":     "na",
			"timestamp":  "1",
			"nameLength": "3",
			"backwards":  "false",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(summoners.(*SummonersMock).GetByNameLengthCalls) != 1 {
		t.Errorf("Expected 1 call to GetByNameLength, got %d", len(summoners.(*SummonersMock).GetByNameLengthCalls))
	}
}

func TestHandleRequest_DefaultsBackwardsToFalseWhenNotPassed(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":    "na",
			"timestamp": "1",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(summoners.(*SummonersMock).GetAfterCalls) != 1 {
		t.Errorf("Expected 1 call to GetAfter, got %d", len(summoners.(*SummonersMock).GetAfterCalls))
	}

	if summoners.(*SummonersMock).GetAfterCalls[0].Backwards {
		t.Errorf("Expected backwards to be false, got true")
	}
}

func TestHandleRequest_CallsGetAfterWithCorrectParameters(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":    "na",
			"timestamp": "1",
			"backwards": "true",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(summoners.(*SummonersMock).GetAfterCalls) != 1 {
		t.Errorf("Expected 1 call to GetAfter, got %d", len(summoners.(*SummonersMock).GetAfterCalls))
	}

	if summoners.(*SummonersMock).GetAfterCalls[0].Region != "NA" {
		t.Errorf("Expected region to be 'NA', got %s", summoners.(*SummonersMock).GetAfterCalls[0].Region)
	}

	if summoners.(*SummonersMock).GetAfterCalls[0].Limit != 35 {
		t.Errorf("Expected limit to be 35, got %d", summoners.(*SummonersMock).GetAfterCalls[0].Limit)
	}

	if summoners.(*SummonersMock).GetAfterCalls[0].T1 != 1 {
		t.Errorf("Expected t1 to be 1, got %d", summoners.(*SummonersMock).GetAfterCalls[0].T1)
	}

	if !summoners.(*SummonersMock).GetAfterCalls[0].Backwards {
		t.Errorf("Expected backwards to be true, got false")
	}
}

func TestHandleRequest_CallsGetByNameLengthWithCorrectParameters(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":     "na",
			"timestamp":  "1",
			"nameLength": "3",
			"backwards":  "false",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(summoners.(*SummonersMock).GetByNameLengthCalls) != 1 {
		t.Errorf("Expected 1 call to GetByNameLength, got %d", len(summoners.(*SummonersMock).GetByNameLengthCalls))
	}

	if summoners.(*SummonersMock).GetByNameLengthCalls[0].Region != "NA" {
		t.Errorf("Expected region to be 'NA', got %s", summoners.(*SummonersMock).GetByNameLengthCalls[0].Region)
	}

	if summoners.(*SummonersMock).GetByNameLengthCalls[0].Limit != 35 {
		t.Errorf("Expected limit to be 35, got %d", summoners.(*SummonersMock).GetByNameLengthCalls[0].Limit)
	}

	if summoners.(*SummonersMock).GetByNameLengthCalls[0].NameLength != 3 {
		t.Errorf("Expected nameLength to be 3, got %d", summoners.(*SummonersMock).GetByNameLengthCalls[0].NameLength)
	}

	if summoners.(*SummonersMock).GetByNameLengthCalls[0].T1 != 1 {
		t.Errorf("Expected t1 to be 1, got %d", summoners.(*SummonersMock).GetByNameLengthCalls[0].T1)
	}

	if summoners.(*SummonersMock).GetByNameLengthCalls[0].Backwards {
		t.Errorf("Expected backwards to be false, got true")
	}
}

func TestHandleRequest_Returns500ErrorWhenGetAfterReturnsError(t *testing.T) {
	setup()
	summoners.(*SummonersMock).ReturnError = true
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":    "na",
			"timestamp": "1",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(responses.(*HttpResponsesMock).ErrorCalls) != 1 {
		t.Errorf("Expected 1 error response, got %d", len(responses.(*HttpResponsesMock).ErrorCalls))
	}

	if responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode != 500 {
		t.Errorf("Expected status code 500, got %d", responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode)
	}

	if !strings.Contains(responses.(*HttpResponsesMock).ErrorCalls[0].Message, "Internal server error") {
		t.Errorf("Expected message 'Internal server error', got %s", responses.(*HttpResponsesMock).ErrorCalls[0].Message)
	}
}

func TestHandleRequest_Returns200SuccessWhenGetAfterReturnsNoError(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		QueryStringParameters: map[string]string{
			"region":    "na",
			"timestamp": "1",
		},
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(responses.(*HttpResponsesMock).SuccessCalls) != 1 {
		t.Errorf("Expected 1 success response, got %d", len(responses.(*HttpResponsesMock).SuccessCalls))
	}

	if responses.(*HttpResponsesMock).SuccessCalls[0] == nil {
		t.Errorf("Expected response to not be nil, got nil")
	}
}

func TestHandleRequest_Returns200OnOptionsRequest(t *testing.T) {
	setup()
	request := events.APIGatewayProxyRequest{
		HTTPMethod: "OPTIONS",
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(responses.(*HttpResponsesMock).SuccessCalls) != 1 {
		t.Errorf("Expected 1 success response, got %d", len(responses.(*HttpResponsesMock).SuccessCalls))
	}
}

func TestHandleRequest_Returns405ErrorWhenMethodIsNotGet(t *testing.T) {
	setup()

	request := events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
	}

	_, err := HandleRequest(context.TODO(), request)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(responses.(*HttpResponsesMock).ErrorCalls) != 1 {
		t.Errorf("Expected 1 error response, got %d", len(responses.(*HttpResponsesMock).ErrorCalls))
	}

	if responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode != 405 {
		t.Errorf("Expected status code 405, got %d", responses.(*HttpResponsesMock).ErrorCalls[0].StatusCode)
	}
}
