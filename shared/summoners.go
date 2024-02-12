package shared

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type dynamoDbService interface {
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

type regionsService interface {
	Validate(region string) bool
	Get(region string) (string, error)
}

type httpService interface {
	Do(req *http.Request) (*http.Response, error)
}

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

type RiotSummonerDTO struct {
	AccountId     string `json:"accountId"`
	ProfileIconId int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	Name          string `json:"name"`
	Id            string `json:"id"`
	Puuid         string `json:"puuid"`
	SummonerLevel int    `json:"summonerLevel"`
}

type Summoners struct {
	dynamodb   dynamoDbService
	regions    regionsService
	http       httpService
	tableName  string
	riotApiKey string
}

func NewSummoners(dynamoDbTableName string, riotApiKey string) (*Summoners, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	return &Summoners{
		dynamodb:   dynamodb.NewFromConfig(cfg),
		regions:    NewRegions(),
		http:       http.DefaultClient,
		tableName:  dynamoDbTableName,
		riotApiKey: riotApiKey,
	}, nil
}

func (s *Summoners) Fetch(region string, summonerName string) (*SummonerDTO, error) {
	riotRegion, err := s.regions.Get(region)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-name/%s", riotRegion, summonerName)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("X-Riot-Token", s.riotApiKey)

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("summoner not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("riot api returned status code %d with body %s", resp.StatusCode, resp.Body)
	}

	var riotSummoner RiotSummonerDTO
	err = json.NewDecoder(resp.Body).Decode(&riotSummoner)
	if err != nil {
		return nil, err
	}

	return s.summonerFromRiotSummoner(&riotSummoner, region)
}

func (s *Summoners) Save(summoner *SummonerDTO) error {
	_, err := s.dynamodb.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item: map[string]types.AttributeValue{
			"n":   &types.AttributeValueMemberS{Value: summoner.Region + "#" + strings.ToUpper(summoner.Name)},
			"r":   &types.AttributeValueMemberS{Value: summoner.Region},
			"ad":  &types.AttributeValueMemberN{Value: strconv.FormatInt(summoner.AvailabilityDate, 10)},
			"aid": &types.AttributeValueMemberS{Value: summoner.AccountID},
			"rd":  &types.AttributeValueMemberN{Value: strconv.FormatInt(summoner.RevisionDate, 10)},
			"l":   &types.AttributeValueMemberN{Value: strconv.Itoa(summoner.Level)},
			"nl":  &types.AttributeValueMemberN{Value: summoner.Region + "#" + strconv.Itoa(len(summoner.Name))},
			"ld":  &types.AttributeValueMemberN{Value: strconv.FormatInt(summoner.LastUpdated, 10)},
			"si":  &types.AttributeValueMemberN{Value: strconv.Itoa(summoner.SummonerIcon)},
		},
	})

	return err
}

func (s *Summoners) Delete(region string, summonerName string) error {
	_, err := s.dynamodb.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			"n": &types.AttributeValueMemberS{Value: region + "#" + strings.ToUpper(summonerName)},
		},
	})

	return err
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

func CalcAvailabilityDate(revisionDate int64, level int32) int64 {
	monthsToAdd := math.Min(30, math.Max(6, float64(level)))
	return time.UnixMilli(revisionDate).AddDate(0, int(monthsToAdd), 0).UnixMilli()
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

func (s *Summoners) summonerFromRiotSummoner(riotSummoner *RiotSummonerDTO, region string) (*SummonerDTO, error) {
	ok := s.regions.Validate(region)
	if !ok {
		return nil, fmt.Errorf("invalid region '%s'", region)
	}

	return &SummonerDTO{
		Name:             strings.ToLower(riotSummoner.Name),
		Region:           region,
		AccountID:        riotSummoner.AccountId,
		RevisionDate:     riotSummoner.RevisionDate,
		AvailabilityDate: CalcAvailabilityDate(riotSummoner.RevisionDate, int32(riotSummoner.SummonerLevel)),
		Level:            riotSummoner.SummonerLevel,
		LastUpdated:      time.Now().UnixMilli(),
		SummonerIcon:     riotSummoner.ProfileIconId,
	}, nil
}
