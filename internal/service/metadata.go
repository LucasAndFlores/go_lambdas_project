package service

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/LucasAndFlores/go_lambdas_project/internal/dto"
	"github.com/LucasAndFlores/go_lambdas_project/internal/entity"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var DYNAMO_TABLE = os.Getenv("DYNAMO_TABLE")

var FileNotFoundErr = errors.New("Filename not found. Unable to complete the operation")
var ConfilctErr = errors.New("The object already exists")

type S3Bucket interface {
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

type DynamoDB interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type MetadataService struct {
	s3     S3Bucket
	dynamo DynamoDB
}

type IMetadataService interface {
	CreateItem(context.Context, dto.MetadataDTOInput) error
	ListAllItems(context.Context) ([]dto.MetadataDTOOutput, error)
}

func NewMetadataService(s S3Bucket, d DynamoDB) IMetadataService {
	return &MetadataService{
		s3:     s,
		dynamo: d,
	}
}

func (s *MetadataService) CreateItem(ctx context.Context, metadata dto.MetadataDTOInput) error {
	_, err := s.s3.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(BUCKET_NAME),
		Key:    aws.String(metadata.FileName),
	})

	if err != nil {
		log.Printf("Error getting head of object %s/%s: %s", BUCKET_NAME, metadata.FileName, err.Error())
		return FileNotFoundErr
	}

	output, err := s.dynamo.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(DYNAMO_TABLE),
		Key: map[string]types.AttributeValue{
			"filename": &types.AttributeValueMemberS{Value: metadata.FileName},
		},
	})

	if err != nil {
		log.Printf("Error when tried to getItem from dynamoDB: %s", err)
		return err
	}

	if output.Item != nil {
		return ConfilctErr
	}

	item := map[string]types.AttributeValue{
		"filename": &types.AttributeValueMemberS{Value: metadata.FileName},
		"author":   &types.AttributeValueMemberS{Value: metadata.Author},
		"label":    &types.AttributeValueMemberS{Value: metadata.Label},
		"type":     &types.AttributeValueMemberS{Value: metadata.Type},
		"words":    &types.AttributeValueMemberS{Value: metadata.Words},
	}

	putItemInput := &dynamodb.PutItemInput{
		TableName: aws.String(DYNAMO_TABLE),
		Item:      item,
	}

	_, err = s.dynamo.PutItem(ctx, putItemInput)

	if err != nil {
		log.Printf("Error when trying to use putItem method: %s", err)
		return err
	}

	return nil

}

func (s *MetadataService) ListAllItems(ctx context.Context) ([]dto.MetadataDTOOutput, error) {
	output, err := s.dynamo.Scan(context.Background(), &dynamodb.ScanInput{TableName: aws.String(DYNAMO_TABLE)})

	if err != nil {
		log.Printf("An error occurred when tried to scan all items. Error: %v", err)
		return []dto.MetadataDTOOutput{}, err
	}

	var listOfAllItems []entity.Metadata

	err = attributevalue.UnmarshalListOfMaps(output.Items, &listOfAllItems)

	if err != nil {
		log.Printf("An error occurred when tried use attributevalue. Error: %v", err)
		return []dto.MetadataDTOOutput{}, err
	}

	if len(listOfAllItems) == 0 {
		return []dto.MetadataDTOOutput{}, nil
	}

	var metadataOutput []dto.MetadataDTOOutput

	for _, val := range listOfAllItems {
		converted := val.ConvertToDTO()

		metadataOutput = append(metadataOutput, converted)
	}

	return metadataOutput, nil
}
