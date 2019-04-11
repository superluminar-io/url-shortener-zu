package main

import (
	"github.com/aws/aws-xray-sdk-go/xray"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	sess := session.Must(session.NewSession())

	dbClient := dynamodb.New(sess)
	xray.AWS(dbClient.Client)

	h := &handler{
		DynamoDBTableName: os.Getenv("DYNAMODB_TABLE_NAME"),
		DynamoDBClient:    dbClient,
	}

	lambda.Start(h.run)
}
