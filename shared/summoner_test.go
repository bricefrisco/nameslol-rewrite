package shared

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"testing"
)

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

	summoners, err := SummonersFromQueryOutput(&output)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if len(summoners) != 2 {
		t.Errorf("Expected 1 summoner, got %d", len(summoners))
	}

	if summoners[0].Name != "name1" {
		t.Errorf("Expected name1, got %s", summoners[0].Name)
	}

	if summoners[0].Region != "NA" {
		t.Errorf("Expected NA, got %s", summoners[0].Region)
	}

	if summoners[0].AccountID != "1234567" {
		t.Errorf("Expected 1234567, got %s", summoners[0].AccountID)
	}

	if summoners[0].RevisionDate != 12345 {
		t.Errorf("Expected 12345, got %d", summoners[0].RevisionDate)
	}

	if summoners[0].AvailabilityDate != 123456 {
		t.Errorf("Expected 123456, got %d", summoners[0].AvailabilityDate)
	}

	if summoners[0].Level != 30 {
		t.Errorf("Expected 30, got %d", summoners[0].Level)
	}

	if summoners[0].LastUpdated != 1612345678 {
		t.Errorf("Expected 1612345678, got %d", summoners[0].LastUpdated)
	}

	if summoners[0].SummonerIcon != 123 {
		t.Errorf("Expected 123, got %d", summoners[0].SummonerIcon)
	}

	if summoners[1].Name != "name2" {
		t.Errorf("Expected name2, got %s", summoners[1].Name)
	}

	if summoners[1].Region != "NA" {
		t.Errorf("Expected na1, got %s", summoners[1].Region)
	}

	if summoners[1].AccountID != "1234567" {
		t.Errorf("Expected 1234567, got %s", summoners[1].AccountID)
	}

	if summoners[1].RevisionDate != 12345 {
		t.Errorf("Expected 12345, got %d", summoners[1].RevisionDate)
	}

	if summoners[1].AvailabilityDate != 123456 {
		t.Errorf("Expected 123456, got %d", summoners[1].AvailabilityDate)
	}

	if summoners[1].Level != 30 {
		t.Errorf("Expected 30, got %d", summoners[1].Level)
	}

	if summoners[1].LastUpdated != 1612345678 {
		t.Errorf("Expected 1612345678, got %d", summoners[1].LastUpdated)
	}

	if summoners[1].SummonerIcon != 123 {
		t.Errorf("Expected 123, got %d", summoners[1].SummonerIcon)
	}
}
