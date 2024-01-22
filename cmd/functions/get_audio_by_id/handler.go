package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func handler(request events.APIGatewayProxyRequest) (Response, error) {
	return Response{
		StatusCode: 200,
		Body:       "ok",
	}, nil
}

func main() {
	lambda.Start(handler)
}

