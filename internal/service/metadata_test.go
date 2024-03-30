package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/LucasAndFlores/go_lambdas_project/internal/dto"
	"github.com/LucasAndFlores/go_lambdas_project/internal/mocks"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestCreateItemSuccessfulResponse(t *testing.T) {
	mockedS3 := mocks.MockedS3{}

	mockedS3.HeadObjectFuncMock = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
		time := time.Now()

		return &s3.HeadObjectOutput{
			LastModified: &time,
		}, nil
	}

	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.PutItemFuncMock = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{
			Attributes: map[string]types.AttributeValue{},
		}, nil
	}

	mockedDynamodb.GetItemFuncMock = func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
		return &dynamodb.GetItemOutput{
			Item: nil,
		}, nil
	}

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	metadata := dto.MetadataDTOInput{
		FileName: "test",
		Author:   "test",
		Label:    "123",
		Words:    "test",
		Type:     "test",
	}

	err := serviceHandler.CreateItem(context.TODO(), metadata)

	if err != nil {
		t.Errorf("The result is different from expected. Expected nil. Result: %v", err)
	}

}

func TestCreateItemS3HeadObjectError(t *testing.T) {
	mockedS3 := mocks.MockedS3{}

	mockedS3.HeadObjectFuncMock = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
		return nil, errors.New("AWS Error")
	}

	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.PutItemFuncMock = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{
			Attributes: map[string]types.AttributeValue{},
		}, nil
	}

	mockedDynamodb.GetItemFuncMock = func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
		return &dynamodb.GetItemOutput{
			Item: nil,
		}, nil
	}

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	metadata := dto.MetadataDTOInput{
		FileName: "test",
		Author:   "test",
		Label:    "123",
		Words:    "test",
		Type:     "test",
	}

	err := serviceHandler.CreateItem(context.TODO(), metadata)

	expected := errors.New("Filename not found. Unable to complete the operation")

	if err.Error() != expected.Error() {
		t.Errorf("Result is different from expected. Expected: %v. Result: %v", expected, err)
	}

}

func TestCreateItemDynamoDBGetItemError(t *testing.T) {
	mockedS3 := mocks.MockedS3{}

	mockedS3.HeadObjectFuncMock = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
		time := time.Now()

		return &s3.HeadObjectOutput{
			LastModified: &time,
		}, nil
	}

	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.PutItemFuncMock = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{
			Attributes: map[string]types.AttributeValue{},
		}, nil
	}

	mockedDynamodb.GetItemFuncMock = func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
		return nil, errors.New("Dynamodb error")
	}

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	metadata := dto.MetadataDTOInput{
		FileName: "test",
		Author:   "test",
		Label:    "123",
		Words:    "test",
		Type:     "test",
	}

	err := serviceHandler.CreateItem(context.TODO(), metadata)

	expected := errors.New("Dynamodb error")

	if err.Error() != expected.Error() {
		t.Errorf("Result is different from expected. Expected: %v. Result: %v", expected, err)
	}

}

func TestCreateItemConflictError(t *testing.T) {
	mockedS3 := mocks.MockedS3{}

	mockedS3.HeadObjectFuncMock = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
		time := time.Now()

		return &s3.HeadObjectOutput{
			LastModified: &time,
		}, nil
	}

	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.PutItemFuncMock = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		return &dynamodb.PutItemOutput{
			Attributes: map[string]types.AttributeValue{},
		}, nil
	}

	mockedDynamodb.GetItemFuncMock = func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
		return &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"filename": &types.AttributeValueMemberS{Value: "test"},
				"author":   &types.AttributeValueMemberS{Value: "test"},
				"label":    &types.AttributeValueMemberS{Value: "123"},
				"words":    &types.AttributeValueMemberS{Value: "test"},
				"type":     &types.AttributeValueMemberS{Value: "test"},
			},
		}, nil
	}

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	metadata := dto.MetadataDTOInput{
		FileName: "test",
		Author:   "test",
		Label:    "123",
		Words:    "test",
		Type:     "test",
	}

	err := serviceHandler.CreateItem(context.TODO(), metadata)

	expected := errors.New("The object already exists")

	if err.Error() != expected.Error() {
		t.Errorf("Result is different from expected. Expected: %v. Result: %v", expected, err)
	}

}

func TestCreateDynamoDBPutItemError(t *testing.T) {
	mockedS3 := mocks.MockedS3{}

	mockedS3.HeadObjectFuncMock = func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
		time := time.Now()

		return &s3.HeadObjectOutput{
			LastModified: &time,
		}, nil
	}

	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.PutItemFuncMock = func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
		return nil, errors.New("Dynamodb error")
	}

	mockedDynamodb.GetItemFuncMock = func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
		return &dynamodb.GetItemOutput{
			Item: nil,
		}, nil
	}

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	metadata := dto.MetadataDTOInput{
		FileName: "test",
		Author:   "test",
		Label:    "123",
		Words:    "test",
		Type:     "test",
	}

	err := serviceHandler.CreateItem(context.TODO(), metadata)

	expected := errors.New("Dynamodb error")

	if err.Error() != expected.Error() {
		t.Errorf("Result is different from expected. Expected: %v. Result: %v", expected, err)
	}

}

func TestListAllItemsSuccessfulResponse(t *testing.T) {
	mockedS3 := mocks.MockedS3{}
	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.ScanFuncMock = func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
		return &dynamodb.ScanOutput{
			Items: []map[string]types.AttributeValue{
				{
					"filename": &types.AttributeValueMemberS{Value: "test"},
					"author":   &types.AttributeValueMemberS{Value: "test"},
					"label":    &types.AttributeValueMemberS{Value: "123"},
					"words":    &types.AttributeValueMemberS{Value: "test"},
					"type":     &types.AttributeValueMemberS{Value: "test"},
				},
			},
		}, nil
	}

	expected := []dto.MetadataDTOOutput{
		{
			FileName: "test",
			Author:   "test",
			Label:    "123",
			Words:    "test",
			Type:     "test",
		},
	}

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	metadata, err := serviceHandler.ListAllItems(context.TODO())

	if err != nil {
		t.Errorf("Expected nil but received an error. Error: %v", err)
	}

	if metadata[0] != expected[0] {
		t.Errorf("The result is different from expected.Result: %v. Expected: %v", metadata[0], expected[0])
	}
}

func TestListAllItemsDynamoDBError(t *testing.T) {
	mockedS3 := mocks.MockedS3{}
	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.ScanFuncMock = func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
		return nil, errors.New("Dynamodb error")
	}

	expected := errors.New("Dynamodb error")

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	metadata, err := serviceHandler.ListAllItems(context.TODO())

	if len(metadata) != 0 {
		t.Errorf("Expected an empty array but received an item. Item: %v", metadata)
	}

	if err.Error() != expected.Error() {
		t.Errorf("The result is different from expected.Result: %v. Expected: %v", err, expected)
	}
}

func TestListAllItemsEmptyDynamoResponse(t *testing.T) {
	mockedS3 := mocks.MockedS3{}
	mockedDynamodb := mocks.MockedDynamoDB{}

	mockedDynamodb.ScanFuncMock = func(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
		return &dynamodb.ScanOutput{
			Items: []map[string]types.AttributeValue{},
		}, nil
	}

	serviceHandler := NewMetadataService(mockedS3, mockedDynamodb)

	metadata, err := serviceHandler.ListAllItems(context.TODO())

	if err != nil {
		t.Errorf("Expected nil but received an error. Error: %v", err)
	}

	if len(metadata) != 0 {
		t.Errorf("Expected an empty array but received an item. Item: %v", metadata)
	}
}
