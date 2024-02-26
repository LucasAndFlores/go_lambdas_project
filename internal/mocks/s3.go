package mocks

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type MockedS3 struct {
	HeadObjectFuncMock func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

func (m MockedS3) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return m.HeadObjectFuncMock(ctx, params, optFns...)
}
