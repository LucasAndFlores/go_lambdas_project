package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/LucasAndFlores/go_lambdas_project/config"
	"github.com/LucasAndFlores/go_lambdas_project/constant"
	"github.com/LucasAndFlores/go_lambdas_project/internal/dto"
	"github.com/LucasAndFlores/go_lambdas_project/internal/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type HttpRequest = events.APIGatewayProxyRequest

type HttpResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type handler struct {
	service service.IMetadataService
}

func (h *handler) handleRequest(ctx context.Context, request HttpRequest) (HttpResponse, error) {
	var parsedBody dto.MetadataDTOInput

	err := json.Unmarshal([]byte(request.Body), &parsedBody)

	if err != nil {
		return HttpResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "Unable to process the body. Please, review the content",
		}, nil
	}

	validatonErr := parsedBody.Validate()

	if validatonErr != nil {
		bytes, err := json.Marshal(map[string]interface{}{"errors": validatonErr})

		if err != nil {
			return HttpResponse{
				StatusCode: http.StatusBadRequest,
				Body:       constant.INTERNAL_SERVER_ERROR,
			}, nil
		}
		return HttpResponse{
			StatusCode: http.StatusBadRequest,
			Body:       string(bytes),
		}, nil
	}

	err = h.service.CreateItem(ctx, parsedBody)

	if err != nil {
		switch {
		case errors.Is(err, service.ConfilctErr):
			return HttpResponse{
				StatusCode: http.StatusConflict,
				Body:       err.Error(),
			}, nil

		case errors.Is(err, service.FileNotFoundErr):
			return HttpResponse{
				StatusCode: http.StatusUnprocessableEntity,
				Body:       err.Error(),
			}, nil

		default:
			return HttpResponse{
				StatusCode: http.StatusInternalServerError,
				Body:       constant.INTERNAL_SERVER_ERROR,
			}, nil
		}
	}

	return HttpResponse{
		StatusCode: http.StatusCreated,
		Body:       "successfully stored",
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
