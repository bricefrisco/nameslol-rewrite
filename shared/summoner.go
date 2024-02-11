package shared

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strconv"
	"strings"
)

type SummonerDTO struct {
	Name             string `json:"name"`
	Region           string `json:"region"`
	AccountID        string `json:"accountId"`
	RevisionDate     int64  `json:"revisionDate"`
	AvailabilityDate int64  `json:"availabilityDate"`
	Level            int    `json:"level"`
	LastUpdated      int64  `json:"lastUpdated"`
	SummonerIcon     int    `json:"summonerIcon"`
}

func SummonersFromQueryOutput(output *dynamodb.QueryOutput) ([]*SummonerDTO, error) {
	summoners := make([]*SummonerDTO, len(output.Items))

	for i, item := range output.Items {
		revisionDate, err := strconv.ParseInt(item["rd"].(*types.AttributeValueMemberN).Value, 10, 64)
		if err != nil {
			return nil, err
		}

		availabilityDate, err := strconv.ParseInt(item["ad"].(*types.AttributeValueMemberN).Value, 10, 64)
		if err != nil {
			return nil, err
		}

		level, err := strconv.Atoi(item["l"].(*types.AttributeValueMemberN).Value)
		if err != nil {
			return nil, err
		}

		lastUpdated, err := strconv.ParseInt(item["ld"].(*types.AttributeValueMemberN).Value, 10, 64)
		if err != nil {
			return nil, err
		}

		summonerIcon, err := strconv.Atoi(item["si"].(*types.AttributeValueMemberN).Value)
		if err != nil {
			return nil, err
		}

		summoners[i] = &SummonerDTO{
			Name:             strings.ToLower(strings.SplitAfter(item["n"].(*types.AttributeValueMemberS).Value, "#")[1]),
			Region:           item["r"].(*types.AttributeValueMemberS).Value,
			AccountID:        item["aid"].(*types.AttributeValueMemberS).Value,
			RevisionDate:     revisionDate,
			AvailabilityDate: availabilityDate,
			Level:            level,
			LastUpdated:      lastUpdated,
			SummonerIcon:     summonerIcon,
		}
	}

	return summoners, nil
}
