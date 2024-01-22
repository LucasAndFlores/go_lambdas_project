package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func LoadPreSignedClient(ctx context.Context) (*s3.PresignClient,
	error) {

	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatal("An error occurred when trying to set AWS Config", err)
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	preSigned := s3.NewPresignClient(client)

	return preSigned, nil
}
