package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.Body == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "Greeting parameter missing"}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello %s", request.Body),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
