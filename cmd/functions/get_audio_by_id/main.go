package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/LucasAndFlores/go_lambdas_project/config"
	"github.com/LucasAndFlores/go_lambdas_project/constant"
	"github.com/LucasAndFlores/go_lambdas_project/internal/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type HttpRequest = events.APIGatewayProxyRequest

type HttpBodyResponse struct {
	Url string `json:"url"`
}

type HttpResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type handler struct {
	service service.IAudioService
}

func (h *handler) handleRequest(ctx context.Context, request HttpRequest) (HttpResponse, error) {
	param := request.PathParameters["filename"]

	if param == "" {
		return HttpResponse{
			StatusCode: http.StatusBadRequest,
			Body:       constant.MISSING_PARAM_ERROR,
		}, nil
	}

	url, err := h.service.GeneratePreSignedGetURL(param, ctx)

	if err != nil {
		return HttpResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       constant.INTERNAL_SERVER_ERROR,
		}, nil
	}

	bytes, err := json.Marshal(HttpBodyResponse{Url: url})

	if err != nil {
		return HttpResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       constant.INTERNAL_SERVER_ERROR,
		}, nil
	}

	return HttpResponse{
		StatusCode: http.StatusOK,
		Body:       string(bytes),
	}, nil
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Fatalf("An error occurred when tried to load AWS config. Error: %v", err)
	}

	bucket := s3.NewFromConfig(cfg)
	preSigned := s3.NewPresignClient(bucket)

	s := service.NewAudioService(preSigned)
	h := handler{service: s}

	lambda.Start(h.handleRequest)
}
