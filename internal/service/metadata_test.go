package service

import (
	"context"
	"errors"
	"testing"

	"github.com/LucasAndFlores/go_lambdas_project/internal/mocks"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestReactToStoredObjectSuccessfulResponse(t *testing.T) {
	mockedS3 := mocks.MockedS3{}

	mockedS3.HeadObjectFuncMock = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
		return &s3.HeadObjectOutput{
			Metadata: map[string]string{
				"author":   "test",
				"filename": "test",
				"label":    "123",
				"words":    "test",
			}}, nil
	}

	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.PutItemFuncMock = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{
			Attributes: map[string]types.AttributeValue{},
		}, nil
	}

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	s3Event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				EventSource: "aws:s3",
				EventName:   "ObjectCreated:Put",
				S3: events.S3Entity{
					Bucket: events.S3Bucket{Name: "test"},
					Object: events.S3Object{Key: "test"},
				},
			},
		},
	}

	err := serviceHandler.ReactToStoredObject(context.TODO(), s3Event)

	if err != nil {
		t.Errorf("The result is different from expected. Expected nil. Result: %v", err)
	}

}

func TestReactToStoredObjectS3Error(t *testing.T) {
	s3Error := errors.New("S3 Bucket error")

	mockedS3 := mocks.MockedS3{}

	mockedS3.HeadObjectFuncMock = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
		return nil, s3Error
	}

	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.PutItemFuncMock = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{
			Attributes: map[string]types.AttributeValue{},
		}, nil
	}

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	s3Event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				EventSource: "aws:s3",
				EventName:   "ObjectCreated:Put",
				S3: events.S3Entity{
					Bucket: events.S3Bucket{Name: "test"},
					Object: events.S3Object{Key: "test"},
				},
			},
		},
	}

	err := serviceHandler.ReactToStoredObject(context.TODO(), s3Event)

	if err != s3Error {
		t.Errorf("The result is different from expected. Expected: %v. Result: %v", s3Error, err)
	}

}

func TestReactToStoredObjectDynamoError(t *testing.T) {
	dynamoError := errors.New("DynamoDB error")

	mockedS3 := mocks.MockedS3{}

	mockedS3.HeadObjectFuncMock = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
		return &s3.HeadObjectOutput{
			Metadata: map[string]string{
				"author":   "test",
				"filename": "test",
				"label":    "123",
				"words":    "test",
			}}, nil
	}

	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.PutItemFuncMock = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		return nil, dynamoError
	}

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	s3Event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				EventSource: "aws:s3",
				EventName:   "ObjectCreated:Put",
				S3: events.S3Entity{
					Bucket: events.S3Bucket{Name: "test"},
					Object: events.S3Object{Key: "test"},
				},
			},
		},
	}

	err := serviceHandler.ReactToStoredObject(context.TODO(), s3Event)

	if err != dynamoError {
		t.Errorf("The result is different from expected. Expected: %v. Result: %v", dynamoError, err)
	}

}

func TestReactToStoredObjectMetadataError(t *testing.T) {
	metadataError := errors.New("One of the metadata fields is missing: words")

	mockedS3 := mocks.MockedS3{}

	mockedS3.HeadObjectFuncMock = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
		return &s3.HeadObjectOutput{
			Metadata: map[string]string{
				"author":   "test",
				"filename": "test",
				"label":    "123",
			}}, nil
	}

	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.PutItemFuncMock = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{
			Attributes: map[string]types.AttributeValue{},
		}, nil
	}

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	s3Event := events.S3Event{
		Records: []events.S3EventRecord{
			{
				EventSource: "aws:s3",
				EventName:   "ObjectCreated:Put",
				S3: events.S3Entity{
					Bucket: events.S3Bucket{Name: "test"},
					Object: events.S3Object{Key: "test"},
				},
			},
		},
	}

	err := serviceHandler.ReactToStoredObject(context.TODO(), s3Event)

	if err.Error() != metadataError.Error() {
		t.Errorf("The result is different from expected. Expected: %v. Result: %v", metadataError, err)
	}

}
