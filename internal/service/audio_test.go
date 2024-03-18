package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/LucasAndFlores/go_lambdas_project/internal/mocks"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestGeneratePreSignedPutURLSuccessfulResponse(t *testing.T) {
	expected := "test.com/audio.mp3"

	preSigned := mocks.MockedPresignedClient{}

	preSigned.PresignPutObjectFuncMock = func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
		return &v4.PresignedHTTPRequest{URL: expected, SignedHeader: http.Header{}, Method: "PUT"}, nil
	}

	serviceHandler := NewAudioService(preSigned)

	filename := "audio"

	url, err := serviceHandler.GeneratePreSignedPutURL(filename, context.TODO())

	if err != nil {
		t.Errorf("An error occurred when tried to test sucesss scenario. Result: %v, Expected: %v", err.Error(), expected)
	}

	if url != expected {
		t.Errorf("The result is different from the expected. Result: %v, Expected: %v", url, expected)
	}
}

func TestGeneratePreSignedPutURLErrorFromAWS(t *testing.T) {
	awsErr := errors.New("AWS Error")

	preSigned := mocks.MockedPresignedClient{}

	preSigned.PresignPutObjectFuncMock = func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
		return nil, awsErr
	}

	serviceHandler := NewAudioService(preSigned)

	filename := "audio"

	url, err := serviceHandler.GeneratePreSignedPutURL(filename, context.TODO())

	if url != "" {
		t.Errorf("An error occurred when tried to test error scenario. Result: %v, Expected: %v", url, nil)
	}

	if err.Error() != awsErr.Error() {
		t.Errorf("The result is different from the expected. Result: %v, Expected: %v", err, awsErr)
	}
}

func TestGeneratePreSignedGetURLSuccessfulResponse(t *testing.T) {
	expected := "test.com/audio.mp3"

	preSigned := mocks.MockedPresignedClient{}

	preSigned.PresignGetObjectFuncMock = func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
		return &v4.PresignedHTTPRequest{URL: expected, SignedHeader: http.Header{}, Method: "GET"}, nil
	}

	serviceHandler := NewAudioService(preSigned)

	filename := "audio"

	url, err := serviceHandler.GeneratePreSignedGetURL(filename, context.TODO())

	if err != nil {
		t.Errorf("An error occurred when tried to test sucesss scenario. Result: %v, Expected: %v", err.Error(), expected)
	}

	if url != expected {
		t.Errorf("The result is different from the expected. Result: %v, Expected: %v", url, expected)
	}
}

func TestGeneratePreSignedGetURLErrorResponse(t *testing.T) {
	awsErr := errors.New("AWS Error")

	preSigned := mocks.MockedPresignedClient{}

	preSigned.PresignGetObjectFuncMock = func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
		return nil, awsErr
	}

	serviceHandler := NewAudioService(preSigned)

	filename := "audio"

	url, err := serviceHandler.GeneratePreSignedGetURL(filename, context.TODO())

	if url != "" {
		t.Errorf("An error occurred when tried to test error scenario. Result: %v, Expected: %v", url, nil)
	}

	if err.Error() != awsErr.Error() {
		t.Errorf("The result is different from the expected. Result: %v, Expected: %v", err, awsErr)
	}

}
