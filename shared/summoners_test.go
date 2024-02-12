package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

type DynamoDBServiceMock struct {
	ShouldReturnError bool
	QueryCalls        []struct {
		Input *dynamodb.QueryInput
	}
	PutItemCalls []struct {
		Input *dynamodb.PutItemInput
	}
	DeleteItemCalls []struct {
		Input *dynamodb.DeleteItemInput
	}
}

func (d *DynamoDBServiceMock) Query(_ context.Context, input *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	d.QueryCalls = append(d.QueryCalls, struct {
		Input *dynamodb.QueryInput
	}{input})

	if d.ShouldReturnError {
		return nil, fmt.Errorf("error")
	}

	return &dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{
			{
				"n":   &types.AttributeValueMemberS{Value: "NA#TEST"},
				"r":   &types.AttributeValueMemberS{Value: "region"},
				"aid": &types.AttributeValueMemberS{Value: "123"},
				"ad":  &types.AttributeValueMemberN{Value: "123"},
				"rd":  &types.AttributeValueMemberN{Value: "123"},
				"l":   &types.AttributeValueMemberN{Value: "123"},
				"ld":  &types.AttributeValueMemberN{Value: "123"},
				"si":  &types.AttributeValueMemberN{Value: "123"},
			},
		},
	}, nil
}

func (d *DynamoDBServiceMock) PutItem(_ context.Context, input *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	d.PutItemCalls = append(d.PutItemCalls, struct {
		Input *dynamodb.PutItemInput
	}{input})

	if d.ShouldReturnError {
		return nil, fmt.Errorf("error")
	}

	return &dynamodb.PutItemOutput{}, nil
}

func (d *DynamoDBServiceMock) DeleteItem(_ context.Context, input *dynamodb.DeleteItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	d.DeleteItemCalls = append(d.DeleteItemCalls, struct {
		Input *dynamodb.DeleteItemInput
	}{input})

	if d.ShouldReturnError {
		return nil, fmt.Errorf("error")
	}

	return &dynamodb.DeleteItemOutput{}, nil
}

type RegionsServiceMock struct {
	IsInvalid bool
}

func (r *RegionsServiceMock) Validate(_ string) bool {
	return !r.IsInvalid
}

func (r *RegionsServiceMock) Get(region string) (string, error) {
	if r.IsInvalid {
		return "", fmt.Errorf("error")
	}

	return region, nil
}

type MockHttpClient struct {
	ShouldFail      bool
	ShouldReturn404 bool
	ShouldReturn500 bool
	Calls           []struct {
		Request *http.Request
	}
}

type MockReadCloser struct {
	io.Reader
	closed bool
}

func (m *MockReadCloser) Close() error {
	m.closed = true
	return nil
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	m.Calls = append(m.Calls, struct {
		Request *http.Request
	}{req})

	if m.ShouldFail {
		return nil, fmt.Errorf("error")
	}

	if m.ShouldReturn404 {
		return &http.Response{
			StatusCode: 404,
		}, nil
	}

	if m.ShouldReturn500 {
		return &http.Response{
			StatusCode: 500,
			Body: &MockReadCloser{
				Reader: strings.NewReader("an error occurred"),
				closed: false,
			},
		}, nil
	}

	riotSummoner := &RiotSummonerDTO{
		AccountId:     "test-aid",
		ProfileIconId: 1,
		RevisionDate:  time.Date(2024, time.February, 12, 8, 21, 30, 0, time.UTC).UnixMilli(),
		Name:          "test",
		Id:            "test-id",
		Puuid:         "test-puuid",
		SummonerLevel: 32,
	}

	jsonBytes, err := json.Marshal(riotSummoner)
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: 200,
		Body: &MockReadCloser{
			Reader: strings.NewReader(string(jsonBytes)),
			closed: false,
		},
	}, nil
}

var tableName string
var riotApiKey string
var summoners *Summoners

func setup() {
	tableName = "test-table"
	riotApiKey = "riot-api-key"

	summoners = &Summoners{
		dynamodb:   &DynamoDBServiceMock{},
		regions:    &RegionsServiceMock{},
		http:       &MockHttpClient{},
		tableName:  tableName,
		riotApiKey: riotApiKey,
	}
}

func TestCalcAvailabilityDate_WhenLevelIsLessThan6_ReturnsRevisionDatePlus6Months(t *testing.T) {
	setup()

	revisionDate := time.Date(2024, time.February, 12, 8, 21, 30, 0, time.UTC)
	level := int32(3)

	// Daylight savings time, subtract one hour
	expected := time.Date(2024, time.August, 12, 7, 21, 30, 0, time.UTC).UnixMilli()
	result := CalcAvailabilityDate(revisionDate.UnixMilli(), level)

	if result != expected {
		t.Errorf("expected %s, got %s", time.UnixMilli(expected), time.UnixMilli(result))
	}
}

func TestCalcAvailabilityDate_WhenLevelIsBetween6And30_ReturnsRevisionDatePlusLMonths(t *testing.T) {
	setup()

	revisionDate := time.Date(2024, time.February, 12, 8, 21, 30, 0, time.UTC)
	level := int32(12)

	expected := time.Date(2025, time.February, 12, 8, 21, 30, 0, time.UTC).UnixMilli()
	result := CalcAvailabilityDate(revisionDate.UnixMilli(), level)

	if result != expected {
		t.Errorf("expected %s, got %s", time.UnixMilli(expected), time.UnixMilli(result))
	}
}

func TestCalcAvailabilityDate_WhenLevelIsGreaterThan30_ReturnsRevisionDatePlus30Months(t *testing.T) {
	setup()

	revisionDate := time.Date(2024, time.February, 12, 8, 21, 30, 0, time.UTC)
	level := int32(35)

	// Daylight savings time, subtract one hour
	expected := time.Date(2026, time.August, 12, 7, 21, 30, 0, time.UTC).UnixMilli()
	result := CalcAvailabilityDate(revisionDate.UnixMilli(), level)

	if result != expected {
		t.Errorf("expected %s, got %s", time.UnixMilli(expected), time.UnixMilli(result))
	}
}

func TestFetch_WhenRegionIsInvalid_ReturnsError(t *testing.T) {
	setup()

	summoners = &Summoners{
		dynamodb:   &DynamoDBServiceMock{},
		regions:    &RegionsServiceMock{IsInvalid: true},
		http:       &MockHttpClient{},
		tableName:  tableName,
		riotApiKey: riotApiKey,
	}

	_, err := summoners.Fetch("invalid", "test")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestFetch_WhenHttpClientFails_ReturnsError(t *testing.T) {
	setup()

	summoners = &Summoners{
		dynamodb:   &DynamoDBServiceMock{},
		regions:    &RegionsServiceMock{},
		http:       &MockHttpClient{ShouldFail: true},
		tableName:  tableName,
		riotApiKey: riotApiKey,
	}

	_, err := summoners.Fetch("na1", "test")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestFetch_CallsHttpClientWithCorrectUrl(t *testing.T) {
	setup()

	region := "na1"
	summonerName := "test"

	_, _ = summoners.Fetch(region, summonerName)

	expectedUrl := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s", region, summonerName)
	actualUrl := summoners.http.(*MockHttpClient).Calls[0].Request.URL.String()

	if actualUrl != expectedUrl {
		t.Errorf("expected %s, got %s", expectedUrl, actualUrl)
	}
}

func TestFetch_CallsHttpClientWithCorrectAuthToken(t *testing.T) {
	setup()

	_, _ = summoners.Fetch("na1", "test")

	expectedAuthToken := riotApiKey
	actualAuthToken := summoners.http.(*MockHttpClient).Calls[0].Request.Header.Get("X-Riot-Token")

	if actualAuthToken != expectedAuthToken {
		t.Errorf("expected %s, got %s", expectedAuthToken, actualAuthToken)
	}
}

func TestFetch_WhenHttpClientReturns404_ReturnsError(t *testing.T) {
	setup()

	summoners = &Summoners{
		dynamodb:   &DynamoDBServiceMock{},
		regions:    &RegionsServiceMock{},
		http:       &MockHttpClient{ShouldReturn404: true},
		tableName:  tableName,
		riotApiKey: riotApiKey,
	}

	_, err := summoners.Fetch("na1", "test")

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	expectedErrorMessage := "summoner not found"
	actualErrorMessage := err.Error()
	if actualErrorMessage != expectedErrorMessage {
		t.Errorf("expected %s, got %s", expectedErrorMessage, actualErrorMessage)
	}
}

func TestFetch_WhenHttpClientReturns500_ReturnsError(t *testing.T) {
	setup()

	summoners = &Summoners{
		dynamodb:   &DynamoDBServiceMock{},
		regions:    &RegionsServiceMock{},
		http:       &MockHttpClient{ShouldReturn500: true},
		tableName:  tableName,
		riotApiKey: riotApiKey,
	}

	_, err := summoners.Fetch("na1", "test")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestFetch_ReturnsCorrectSummonerWhenResponseSuccessful(t *testing.T) {
	setup()

	s, err := summoners.Fetch("NA", "test")
	if err != nil {
		t.Errorf("expected nil, got %s", err)
	}

	if s.Name != "test" {
		t.Errorf("expected test, got %s", s.Name)
	}

	if s.Region != "NA" {
		t.Errorf("expected NA, got %s", s.Region)
	}

	if s.AccountID != "test-aid" {
		t.Errorf("expected test-aid, got %s", s.AccountID)
	}

	expectedDate := time.Date(2024, time.February, 12, 8, 21, 30, 0, time.UTC).UnixMilli()
	if s.RevisionDate != expectedDate {
		t.Errorf("expected %d, got %d", expectedDate, s.RevisionDate)
	}

	// Daylight savings time, subtract one hour
	expectedAvailabilityDate := time.Date(2026, time.August, 12, 7, 21, 30, 0, time.UTC).UnixMilli()
	if s.AvailabilityDate != expectedAvailabilityDate {
		t.Errorf("expected %d, got %d", expectedAvailabilityDate, s.AvailabilityDate)
	}

	if s.Level != 32 {
		t.Errorf("expected 32, got %d", s.Level)
	}

	if s.LastUpdated == 0 {
		t.Errorf("expected non-zero, got 0")
	}

	if s.SummonerIcon != 1 {
		t.Errorf("expected 1, got %d", s.SummonerIcon)
	}
}

func TestGetBetweenDate_ReturnsErrorIfRegionValidationFails(t *testing.T) {
	tableName = "test-table"
	riotApiKey = "riot-api-key"

	s := &Summoners{
		dynamodb:   &DynamoDBServiceMock{},
		regions:    &RegionsServiceMock{IsInvalid: true},
		http:       &MockHttpClient{},
		tableName:  tableName,
		riotApiKey: riotApiKey,
	}

	_, err := s.GetBetweenDate("invalid", 10, 0, 0)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetBetweenDate_ReturnsErrorIfDynamoDBQueryFails(t *testing.T) {
	s := &Summoners{
		dynamodb:  &DynamoDBServiceMock{ShouldReturnError: true},
		regions:   &RegionsServiceMock{},
		tableName: "test-table",
	}

	_, err := s.GetBetweenDate("region", 10, 0, 0)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestGetBetweenDate_CallsDynamoDBQueryWithCorrectInput(t *testing.T) {
	setup()

	region := "region"
	limit := int32(30)
	revisionDate := int64(10)
	availabilityDate := int64(20)

	_, err := summoners.GetBetweenDate(region, limit, availabilityDate, revisionDate)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	expectedInput := &dynamodb.QueryInput{
		TableName:              &tableName,
		Limit:                  &limit,
		KeyConditionExpression: aws.String("r = :region and ad between :t1 and :t2"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":region": &types.AttributeValueMemberS{Value: region},
			":t1":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", availabilityDate)},
			":t2":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", revisionDate)},
		}, IndexName: aws.String("region-availability-date-index"),
	}

	actualInput := summoners.dynamodb.(*DynamoDBServiceMock).QueryCalls[0].Input

	if *actualInput.TableName != *expectedInput.TableName {
		t.Errorf("expected %s, got %s", *expectedInput.TableName, *actualInput.TableName)
	}

	if *actualInput.Limit != *expectedInput.Limit {
		t.Errorf("expected %d, got %d", *expectedInput.Limit, *actualInput.Limit)
	}

	if *actualInput.KeyConditionExpression != *expectedInput.KeyConditionExpression {
		t.Errorf("expected %s, got %s", *expectedInput.KeyConditionExpression, *actualInput.KeyConditionExpression)
	}

	if *actualInput.IndexName != *expectedInput.IndexName {
		t.Errorf("expected %s, got %s", *expectedInput.IndexName, *actualInput.IndexName)
	}
}

func TestGetBetweenDate_ReturnsCorrectSummonerDTOs(t *testing.T) {
	s := &Summoners{
		dynamodb:  &DynamoDBServiceMock{},
		regions:   &RegionsServiceMock{},
		tableName: "test-table",
	}

	result, err := s.GetBetweenDate("region", 10, 0, 0)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if len(result) != 1 {
		t.Errorf("expected 1, got %d", len(result))
	}
}

func TestSummonersFromQueryOutput(t *testing.T) {
	output := dynamodb.QueryOutput{
		Items: []map[string]types.AttributeValue{
			{
				"rd":  &types.AttributeValueMemberN{Value: "12345"},
				"ad":  &types.AttributeValueMemberN{Value: "123456"},
				"l":   &types.AttributeValueMemberN{Value: "30"},
				"ld":  &types.AttributeValueMemberN{Value: "1612345678"},
				"si":  &types.AttributeValueMemberN{Value: "123"},
				"n":   &types.AttributeValueMemberS{Value: "NA#NAME1"},
				"r":   &types.AttributeValueMemberS{Value: "NA"},
				"aid": &types.AttributeValueMemberS{Value: "1234567"},
			},
			{
				"rd":  &types.AttributeValueMemberN{Value: "12345"},
				"ad":  &types.AttributeValueMemberN{Value: "123456"},
				"l":   &types.AttributeValueMemberN{Value: "30"},
				"ld":  &types.AttributeValueMemberN{Value: "1612345678"},
				"si":  &types.AttributeValueMemberN{Value: "123"},
				"n":   &types.AttributeValueMemberS{Value: "EUW#NAME2"},
				"r":   &types.AttributeValueMemberS{Value: "NA"},
				"aid": &types.AttributeValueMemberS{Value: "1234567"},
			},
		},
	}

	result, err := SummonersFromQueryOutput(&output)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 1 summoner, got %d", len(result))
	}

	if result[0].Name != "name1" {
		t.Errorf("Expected name1, got %s", result[0].Name)
	}

	if result[0].Region != "NA" {
		t.Errorf("Expected NA, got %s", result[0].Region)
	}

	if result[0].AccountID != "1234567" {
		t.Errorf("Expected 1234567, got %s", result[0].AccountID)
	}

	if result[0].RevisionDate != 12345 {
		t.Errorf("Expected 12345, got %d", result[0].RevisionDate)
	}

	if result[0].AvailabilityDate != 123456 {
		t.Errorf("Expected 123456, got %d", result[0].AvailabilityDate)
	}

	if result[0].Level != 30 {
		t.Errorf("Expected 30, got %d", result[0].Level)
	}

	if result[0].LastUpdated != 1612345678 {
		t.Errorf("Expected 1612345678, got %d", result[0].LastUpdated)
	}

	if result[0].SummonerIcon != 123 {
		t.Errorf("Expected 123, got %d", result[0].SummonerIcon)
	}

	if result[1].Name != "name2" {
		t.Errorf("Expected name2, got %s", result[1].Name)
	}

	if result[1].Region != "NA" {
		t.Errorf("Expected na1, got %s", result[1].Region)
	}

	if result[1].AccountID != "1234567" {
		t.Errorf("Expected 1234567, got %s", result[1].AccountID)
	}

	if result[1].RevisionDate != 12345 {
		t.Errorf("Expected 12345, got %d", result[1].RevisionDate)
	}

	if result[1].AvailabilityDate != 123456 {
		t.Errorf("Expected 123456, got %d", result[1].AvailabilityDate)
	}

	if result[1].Level != 30 {
		t.Errorf("Expected 30, got %d", result[1].Level)
	}

	if result[1].LastUpdated != 1612345678 {
		t.Errorf("Expected 1612345678, got %d", result[1].LastUpdated)
	}

	if result[1].SummonerIcon != 123 {
		t.Errorf("Expected 123, got %d", result[1].SummonerIcon)
	}
}

func TestSave_ReturnsErrorIfDynamoDBPutItemFails(t *testing.T) {
	s := &Summoners{
		dynamodb:   &DynamoDBServiceMock{ShouldReturnError: true},
		regions:    &RegionsServiceMock{},
		tableName:  tableName,
		riotApiKey: riotApiKey,
	}

	err := s.Save(&SummonerDTO{})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestSave_CallsDynamoDBPutItemWithCorrectInput(t *testing.T) {
	setup()

	s := &SummonerDTO{
		Name:             "test",
		Region:           "NA",
		AccountID:        "123",
		RevisionDate:     123,
		AvailabilityDate: 1234,
		Level:            30,
		LastUpdated:      123456,
		SummonerIcon:     123,
	}

	err := summoners.Save(s)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	expectedInput := &dynamodb.PutItemInput{
		TableName: &tableName,
		Item: map[string]types.AttributeValue{
			"n":   &types.AttributeValueMemberS{Value: "NA#TEST"},
			"r":   &types.AttributeValueMemberS{Value: "NA"},
			"ad":  &types.AttributeValueMemberN{Value: "1234"},
			"aid": &types.AttributeValueMemberS{Value: "123"},
			"rd":  &types.AttributeValueMemberN{Value: "123"},
			"l":   &types.AttributeValueMemberN{Value: "30"},
			"nl":  &types.AttributeValueMemberN{Value: "NA#4"},
			"ld":  &types.AttributeValueMemberN{Value: "123456"},
			"si":  &types.AttributeValueMemberN{Value: "123"},
		},
	}

	actualInput := summoners.dynamodb.(*DynamoDBServiceMock).PutItemCalls[0].Input

	if *actualInput.TableName != *expectedInput.TableName {
		t.Errorf("expected %s, got %s", *expectedInput.TableName, *actualInput.TableName)
	}

	log.Println(actualInput.Item["n"])

	if actualInput.Item["n"].(*types.AttributeValueMemberS).Value != expectedInput.Item["n"].(*types.AttributeValueMemberS).Value {
		t.Errorf("expected %s, got %s", expectedInput.Item["n"].(*types.AttributeValueMemberS).Value, actualInput.Item["n"].(*types.AttributeValueMemberS).Value)
	}

	if actualInput.Item["r"].(*types.AttributeValueMemberS).Value != expectedInput.Item["r"].(*types.AttributeValueMemberS).Value {
		t.Errorf("expected %s, got %s", expectedInput.Item["r"].(*types.AttributeValueMemberS).Value, actualInput.Item["r"].(*types.AttributeValueMemberS).Value)
	}

	if actualInput.Item["ad"].(*types.AttributeValueMemberN).Value != expectedInput.Item["ad"].(*types.AttributeValueMemberN).Value {
		t.Errorf("expected %s, got %s", expectedInput.Item["ad"].(*types.AttributeValueMemberN).Value, actualInput.Item["ad"].(*types.AttributeValueMemberN).Value)
	}

	if actualInput.Item["aid"].(*types.AttributeValueMemberS).Value != expectedInput.Item["aid"].(*types.AttributeValueMemberS).Value {
		t.Errorf("expected %s, got %s", expectedInput.Item["aid"].(*types.AttributeValueMemberS).Value, actualInput.Item["aid"].(*types.AttributeValueMemberS).Value)
	}

	if actualInput.Item["rd"].(*types.AttributeValueMemberN).Value != expectedInput.Item["rd"].(*types.AttributeValueMemberN).Value {
		t.Errorf("expected %s, got %s", expectedInput.Item["rd"].(*types.AttributeValueMemberN).Value, actualInput.Item["rd"].(*types.AttributeValueMemberN).Value)
	}

	if actualInput.Item["l"].(*types.AttributeValueMemberN).Value != expectedInput.Item["l"].(*types.AttributeValueMemberN).Value {
		t.Errorf("expected %s, got %s", expectedInput.Item["l"].(*types.AttributeValueMemberN).Value, actualInput.Item["l"].(*types.AttributeValueMemberN).Value)
	}

	if actualInput.Item["nl"].(*types.AttributeValueMemberN).Value != expectedInput.Item["nl"].(*types.AttributeValueMemberN).Value {
		t.Errorf("expected %s, got %s", expectedInput.Item["nl"].(*types.AttributeValueMemberN).Value, actualInput.Item["nl"].(*types.AttributeValueMemberN).Value)
	}

	if actualInput.Item["ld"].(*types.AttributeValueMemberN).Value != expectedInput.Item["ld"].(*types.AttributeValueMemberN).Value {
		t.Errorf("expected %s, got %s", expectedInput.Item["ld"].(*types.AttributeValueMemberN).Value, actualInput.Item["ld"].(*types.AttributeValueMemberN).Value)
	}

	if actualInput.Item["si"].(*types.AttributeValueMemberN).Value != expectedInput.Item["si"].(*types.AttributeValueMemberN).Value {
		t.Errorf("expected %s, got %s", expectedInput.Item["si"].(*types.AttributeValueMemberN).Value, actualInput.Item["si"].(*types.AttributeValueMemberN).Value)
	}
}

func TestDelete_ReturnsErrorWhenDynamoDBDeleteItemFails(t *testing.T) {
	s := &Summoners{
		dynamodb:   &DynamoDBServiceMock{ShouldReturnError: true},
		regions:    &RegionsServiceMock{},
		tableName:  tableName,
		riotApiKey: riotApiKey,
	}

	err := s.Delete("region", "test")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestDelete_CallsDynamoDBDeleteItemWithCorrectInput(t *testing.T) {
	setup()

	err := summoners.Delete("NA", "test")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	expectedInput := &dynamodb.DeleteItemInput{
		TableName: &tableName,
		Key: map[string]types.AttributeValue{
			"n": &types.AttributeValueMemberS{Value: "NA#TEST"},
		},
	}

	actualInput := summoners.dynamodb.(*DynamoDBServiceMock).DeleteItemCalls[0].Input

	if *actualInput.TableName != *expectedInput.TableName {
		t.Errorf("expected %s, got %s", *expectedInput.TableName, *actualInput.TableName)
	}

	if actualInput.Key["n"].(*types.AttributeValueMemberS).Value != expectedInput.Key["n"].(*types.AttributeValueMemberS).Value {
		t.Errorf("expected %s, got %s", expectedInput.Key["n"].(*types.AttributeValueMemberS).Value, actualInput.Key["n"].(*types.AttributeValueMemberS).Value)
	}
}
