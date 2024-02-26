package mocks

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type MockedDynamoDB struct {
	PutItemFuncMock func(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

func (m MockedDynamoDB) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return m.PutItemFuncMock(ctx, params, optFns...)
}
