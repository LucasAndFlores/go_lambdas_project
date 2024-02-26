package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func LoadDefaultConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		log.Fatal("An error occurred when trying to set AWS Config", err)
		return aws.Config{}, err
	}

	return cfg, nil
}
