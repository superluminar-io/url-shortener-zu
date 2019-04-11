package main

import (
	"encoding/json"
	"hash/fnv"
	neturl "net/url"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

func shorten(u string) (string, error) {
	if _, err := neturl.ParseRequestURI(u); err != nil {
		return "", err
	}

	hash := fnv.New64a()
	if _, err := hash.Write([]byte(u)); err != nil {
		return "", err
	}

	return strconv.FormatUint(hash.Sum64(), 36), nil
}

type handler struct {
	DynamoDBTableName string
	DynamoDBClient    dynamodbiface.DynamoDBAPI
}

func (h *handler) run(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var data map[string]string
	err := json.Unmarshal([]byte(request.Body), &data)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Unable to parse payload"}, nil
	}

	url, ok := data["url"]
	if !ok {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Unable to access URL in payload"}, nil
	}

	id, err := shorten(url)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Unable to shorten provided URL in payload"}, nil
	}

	_, err = h.DynamoDBClient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(h.DynamoDBTableName),
		Item: map[string]*dynamodb.AttributeValue{
			"id":  {S: aws.String(id)},
			"url": {S: aws.String(url)},
		},
	})

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Unable to store URL mapping"}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 201, Body: id}, nil
}
