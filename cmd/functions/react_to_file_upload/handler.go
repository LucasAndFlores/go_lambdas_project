package main

import (
	"context"
	"log"

	"github.com/LucasAndFlores/go_lambdas_project/config"
	"github.com/LucasAndFlores/go_lambdas_project/internal/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type handler struct {
	service service.IMetadataService
}

func (h *handler) handleRequest(ctx context.Context, s3Event events.S3Event) error {
	e := h.service.ReactToStoredObject(ctx, s3Event)
	return e
}

func main() {

	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Fatalf("An error occurred when tried to load AWS config. Error: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)

	dynamo := dynamodb.NewFromConfig(cfg)

	s := service.NewMetadataService(s3Client, dynamo)
	h := handler{service: s}

	lambda.Start(h.handleRequest)
}
