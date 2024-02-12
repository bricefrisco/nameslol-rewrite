package shared

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
)

type DynamoDBServiceMock struct {
	ShouldReturnError bool
}

func (d *DynamoDBServiceMock) Query(_ context.Context, _ *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
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

type RegionsServiceMock struct {
	ShouldReturnFalse bool
}

func (r *RegionsServiceMock) Validate(_ string) bool {
	return !r.ShouldReturnFalse
}

func TestGetBetweenDate_ReturnsErrorIfRegionValidationFails(t *testing.T) {
	s := &Summoners{
		dynamodb:  &DynamoDBServiceMock{},
		regions:   &RegionsServiceMock{ShouldReturnFalse: true},
		tableName: "test-table",
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

func TestGetBetweenDate_ReturnsSummonerDTOs(t *testing.T) {
	s := &Summoners{
		dynamodb:  &DynamoDBServiceMock{},
		regions:   &RegionsServiceMock{},
		tableName: "test-table",
	}

	summoners, err := s.GetBetweenDate("region", 10, 0, 0)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}

	if len(summoners) != 1 {
		t.Errorf("expected 1, got %d", len(summoners))
	}
}
