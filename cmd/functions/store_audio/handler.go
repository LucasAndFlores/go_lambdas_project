package main

import (
	"context"
	"log"

	"github.com/LucasAndFlores/go_lambdas_project/config"
	"github.com/LucasAndFlores/go_lambdas_project/internal/service"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type HttpRequest = events.APIGatewayProxyRequest

type HttpBodyResponse = map[string]interface{}

type HttpResponse struct {
	StatusCode int              `json:"statusCode"`
	Body       HttpBodyResponse `json:"body"`
}

type handler struct {
	service service.IAudioService
}

func (h *handler) handleRequest(ctx context.Context, request HttpRequest) (HttpResponse, error) {
	statusCode, body := h.service.GeneratePreSignedPutURL(request.Body, ctx)

	return HttpResponse{
		StatusCode: statusCode,
		Body:       body,
	}, nil
}

func main() {
	preSigned, err := config.LoadPreSignedClient(context.Background())

	if err != nil {
		log.Fatalf("An error occurred when tried to load AWS config. Error: %v", err)
	}

	s := service.NewAudioService(preSigned)
	h := handler{service: s}

	lambda.Start(h.handleRequest)
}

