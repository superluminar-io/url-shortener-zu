package main

import (
	"fmt"
	"github.com/badoux/goscraper"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type Response events.APIGatewayProxyResponse

type handler struct {
	DynamoDBTableName string
	DynamoDBClient    dynamodbiface.DynamoDBAPI
}

func (h *handler) run(event events.DynamoDBEvent) (Response, error) {

	for _, r := range event.Records {
		logrus.Infof("%v", r.Change.NewImage)
		url, ok := r.Change.NewImage["url"]
		if !ok {
			return Response{StatusCode: 501}, fmt.Errorf("cant handle event: %v", event)
		}

		s, err := goscraper.Scrape(url.String(), 5)
		if err != nil {
			logrus.WithField("error", err).Errorf("failed to scrape '%s'", url)
		}

		_, err = h.DynamoDBClient.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String(h.DynamoDBTableName),
			Item: map[string]*dynamodb.AttributeValue{
				"url":   {S: aws.String(url.String())},
				"image": {S: aws.String(s.Preview.Images[0])},
				"name":  {S: aws.String(s.Preview.Name)},
				"title": {S: aws.String(s.Preview.Title)},
			}})

		if err != nil {
			logrus.WithField("error", err).Error("Couldn't save Preview")
		}
	}

	resp := Response{
		StatusCode: 201,
		Body:       "s",
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
	}

	return resp, nil
}
