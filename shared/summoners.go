package shared

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"strconv"
)

type DynamoDBService interface {
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
}

type RegionsService interface {
	Validate(region string) bool
}

type Summoners struct {
	dynamodb  DynamoDBService
	regions   RegionsService
	tableName string
}

func NewSummoners(dynamoDbTableName string) (*Summoners, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Summoners{
		dynamodb:  dynamodb.NewFromConfig(cfg),
		regions:   NewRegions(),
		tableName: dynamoDbTableName,
	}, nil
}

func (s *Summoners) GetBetweenDate(region string, limit int32, t1 int64, t2 int64) ([]*SummonerDTO, error) {
	valid := s.regions.Validate(region)
	if !valid {
		return nil, fmt.Errorf("invalid region '%s'", region)
	}

	output, err := s.dynamodb.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(s.tableName),
		Limit:                  aws.Int32(limit),
		KeyConditionExpression: aws.String("r = :region and ad between :t1 and :t2"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":region": &types.AttributeValueMemberS{Value: region},
			":t1":     &types.AttributeValueMemberN{Value: strconv.FormatInt(t1, 10)},
			":t2":     &types.AttributeValueMemberN{Value: strconv.FormatInt(t2, 10)},
		},
		IndexName: aws.String("region-availability-date-index"),
	})

	if err != nil {
		return nil, err
	}

	return SummonersFromQueryOutput(output)
}
