package main

import (
	"context"
	"log"
	"net/http"

	"github.com/LucasAndFlores/go_lambdas_project/config"
	"github.com/LucasAndFlores/go_lambdas_project/constant"
	"github.com/LucasAndFlores/go_lambdas_project/internal/service"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type handler struct {
	service service.IMetadataService
}

type HttpBodyResponse = map[string]interface{}

type Response struct {
	StatusCode int              `json:"statusCode"`
	Body       HttpBodyResponse `json:"body"`
}

func (h *handler) handleRequest(ctx context.Context) (Response, error) {
	metadata, err := h.service.ListAllItems(ctx)

	if err != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Body:       HttpBodyResponse{"message": constant.INTERNAL_SERVER_ERROR},
		}, nil
	}

	return Response{
		StatusCode: http.StatusOK,
		Body:       HttpBodyResponse{"metadata": metadata},
	}, nil
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
