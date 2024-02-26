package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var DYNAMO_TABLE = os.Getenv("DYNAMO_TABLE")

type S3Bucket interface {
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

type DynamoDB interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type MetadataService struct {
	s3     S3Bucket
	dynamo DynamoDB
}

type IMetadataService interface {
	ReactToStoredObject(context.Context, events.S3Event) error
}

func NewMetadataService(s S3Bucket, d DynamoDB) IMetadataService {
	return &MetadataService{
		s3:     s,
		dynamo: d,
	}
}

func (s *MetadataService) ReactToStoredObject(ctx context.Context, events events.S3Event) error {
	for _, record := range events.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.URLDecodedKey
		headOutput, err := s.s3.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})

		if err != nil {
			log.Printf("Error getting head of object %s/%s: %s", bucket, key, err)
			return err
		}

		err = checkMetadataFields(headOutput.Metadata)

		if err != nil {
			log.Printf("Error when checking the metadata fields: %v", err)
			return err
		}

		m, err := attributevalue.MarshalMap(headOutput.Metadata)

		if err != nil {
			log.Printf("Error when tried to marshall: %s", err)
			return err
		}

		_, err = s.dynamo.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(DYNAMO_TABLE), Item: m,
		})

		if err != nil {
			log.Printf("Error when trying to use putItem method: %s", err)
			return err
		}

		log.Println("The object metadata was successfully stored in dynamodb")
	}

	return nil
}

func checkMetadataFields(inputMap map[string]string) error {
	requiredFields := []string{"author", "filename", "label", "words"}

	for _, field := range requiredFields {
		_, ok := inputMap[field]

		if !ok {
			return fmt.Errorf("One of the metadata fields is missing: %s", field)
		}
	}

	return nil
}
