package service

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/LucasAndFlores/go_lambdas_project/constant"
	"github.com/LucasAndFlores/go_lambdas_project/internal/dto"
	"github.com/LucasAndFlores/go_lambdas_project/internal/mocks"
	"github.com/aws/aws-lambda-go/events"
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

	request := events.APIGatewayProxyRequest{Body: "{\"fileName\":\"test\",\"author\":\"test\",\"label\":\"test\",\"type\":\"test\",\"words\":\"test\"}"}

	statusCode, body := serviceHandler.GeneratePreSignedPutURL(request.Body, context.TODO())

	if statusCode != http.StatusCreated {
		t.Errorf("An error occurred when tried to test sucesss scenario. Result: %v, Expected: %v", statusCode, http.StatusCreated)
	}

	if body["url"] != expected {
		t.Errorf("The result is different from the expected. Result: %v, Expected: %v", statusCode, expected)
	}
}

func TestGeneratePreSignedPutURLWrongJSONFormat(t *testing.T) {
	awsResult := "test.com/audio.mp3"

	preSigned := mocks.MockedPresignedClient{}

	preSigned.PresignPutObjectFuncMock = func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
		return &v4.PresignedHTTPRequest{URL: awsResult, SignedHeader: http.Header{}, Method: "PUT"}, nil
	}

	serviceHandler := NewAudioService(preSigned)

	wrongJSON := events.APIGatewayProxyRequest{Body: "{fileName:\"test\",author:\"test\",label:\"test\",type:\"test\",words:\"test\"}"}

	statusCode, body := serviceHandler.GeneratePreSignedPutURL(wrongJSON.Body, context.TODO())

	if statusCode != http.StatusUnprocessableEntity {
		t.Errorf("An error occurred when tried to test error scenario. Result: %v, Expected: %v", statusCode, http.StatusUnprocessableEntity)
	}

	if body["message"] != constant.INTERNAL_SERVER_ERROR {
		t.Errorf("The result is different from the expected. Result: %v, Expected: %v", body["message"], constant.INTERNAL_SERVER_ERROR)
	}
}

func TestGeneratePreSignedPutURLBodyValidation(t *testing.T) {
	awsResult := "test.com/audio.mp3"

	preSigned := mocks.MockedPresignedClient{}

	preSigned.PresignPutObjectFuncMock = func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
		return &v4.PresignedHTTPRequest{URL: awsResult, SignedHeader: http.Header{}, Method: "PUT"}, nil
	}

	serviceHandler := NewAudioService(preSigned)

	type testCase struct {
		payload  events.APIGatewayProxyRequest
		expected dto.AudioInputError
	}

	tCases := []testCase{
		{
			payload: events.APIGatewayProxyRequest{Body: "{\"author\":\"test\",\"label\":\"test\",\"type\":\"test\",\"words\":\"test\"}"},
			expected: dto.AudioInputError{
				Field: "FileName",
				Tag:   "required",
				Value: "",
			},
		},
		{
			payload: events.APIGatewayProxyRequest{Body: "{\"fileName\":\"test\",\"label\":\"test\",\"type\":\"test\",\"words\":\"test\"}"},
			expected: dto.AudioInputError{
				Field: "Author",
				Tag:   "required",
				Value: "",
			},
		},
		{
			payload: events.APIGatewayProxyRequest{Body: "{\"fileName\":\"test\",\"author\":\"test\",\"type\":\"test\",\"words\":\"test\"}"},
			expected: dto.AudioInputError{
				Field: "Label",
				Tag:   "required",
				Value: "",
			},
		},
		{
			payload: events.APIGatewayProxyRequest{Body: "{\"fileName\":\"test\",\"author\":\"test\",\"label\":\"test\",\"words\":\"test\"}"},
			expected: dto.AudioInputError{
				Field: "Type",
				Tag:   "required",
				Value: "",
			},
		},
		{
			payload: events.APIGatewayProxyRequest{Body: "{\"fileName\":\"test\",\"author\":\"test\",\"label\":\"test\",\"type\":\"test\"}"},
			expected: dto.AudioInputError{
				Field: "Words",
				Tag:   "required",
				Value: "",
			},
		},
	}

	for _, value := range tCases {
		statusCode, body := serviceHandler.GeneratePreSignedPutURL(value.payload.Body, context.TODO())

		if statusCode != http.StatusUnprocessableEntity {
			t.Errorf("The status code is differend from expected. Result: %v, Expected: %v", statusCode, http.StatusUnprocessableEntity)
		}

		var rawResponse interface{}

		rawResponse = body["errors"]

		if parsed, ok := rawResponse.([]dto.AudioInputError); ok {
			if parsed[0] != value.expected {
				t.Errorf("The result is different from expected. Result: %v. Expected: %v", parsed[0], value.expected)
			}
		} else {
			t.Error("An error occurred when tried to parse the returned result")
		}

	}
}

func TestGeneratePreSignedPutURLErrorFromAWS(t *testing.T) {
	awsErr := errors.New("AWS Error")

	preSigned := mocks.MockedPresignedClient{}

	preSigned.PresignPutObjectFuncMock = func(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
		return nil, awsErr
	}

	serviceHandler := NewAudioService(preSigned)

	request := events.APIGatewayProxyRequest{Body: "{\"fileName\":\"test\",\"author\":\"test\",\"label\":\"test\",\"type\":\"test\",\"words\":\"test\"}"}

	statusCode, body := serviceHandler.GeneratePreSignedPutURL(request.Body, context.TODO())

	if statusCode != http.StatusInternalServerError {
		t.Errorf("An error occurred when tried to test sucesss scenario. Result: %v, Expected: %v", statusCode, http.StatusInternalServerError)
	}

	if body["message"] != constant.INTERNAL_SERVER_ERROR {
		t.Errorf("The result is different from the expected. Result: %v, Expected: %v", body["message"], awsErr)
	}
}

func TestGeneratePreSignedGetURLSuccessfulResponse(t *testing.T) {
	expected := "test.com/audio.mp3"

	preSigned := mocks.MockedPresignedClient{}

	preSigned.PresignGetObjectFuncMock = func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
		return &v4.PresignedHTTPRequest{URL: expected, SignedHeader: http.Header{}, Method: "GET"}, nil
	}

	serviceHandler := NewAudioService(preSigned)

	statusCode, body := serviceHandler.GeneratePreSignedGetURL("test", context.TODO())

	if body["url"] != expected {
		t.Errorf("An error occurred when tried to test sucesss scenario. Result: %v, Expected: %v", body["url"], expected)
	}

	if statusCode != http.StatusOK {
		t.Errorf("An error occurred when tried to test sucesss scenario. Result: %v, Expected: %v", statusCode, http.StatusOK)
	}
}

func TestGeneratePreSignedGetURLMissingParam(t *testing.T) {
	expected := "test.com/audio.mp3"

	preSigned := mocks.MockedPresignedClient{}

	preSigned.PresignGetObjectFuncMock = func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
		return &v4.PresignedHTTPRequest{URL: expected, SignedHeader: http.Header{}, Method: "GET"}, nil
	}

	serviceHandler := NewAudioService(preSigned)

	statusCode, body := serviceHandler.GeneratePreSignedGetURL("", context.TODO())

	if body["message"] != constant.MISSING_PARAM_ERROR {
		t.Errorf("An error occurred when tried to test error scenario. Result: %v, Expected: %v", body["url"], constant.INTERNAL_SERVER_ERROR)
	}

	if statusCode != http.StatusBadRequest {
		t.Errorf("An error occurred when tried to test error scenario. Result: %v, Expected: %v", statusCode, http.StatusOK)
	}

}

func TestGeneratePreSignedGetURLErrorResponse(t *testing.T) {
	awsErr := errors.New("AWS Error")

	preSigned := mocks.MockedPresignedClient{}

	preSigned.PresignGetObjectFuncMock = func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
		return nil, awsErr
	}

	serviceHandler := NewAudioService(preSigned)

	statusCode, body := serviceHandler.GeneratePreSignedGetURL("test", context.TODO())

	if body["message"] != constant.INTERNAL_SERVER_ERROR {
		t.Errorf("An error occurred when tried to test error scenario. Result: %v, Expected: %v", body["url"], constant.INTERNAL_SERVER_ERROR)
	}

	if statusCode != http.StatusInternalServerError {
		t.Errorf("An error occurred when tried to test error scenario. Result: %v, Expected: %v", statusCode, http.StatusOK)
	}

}
