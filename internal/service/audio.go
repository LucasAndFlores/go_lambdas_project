package service

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/LucasAndFlores/go_lambdas_project/constant"
	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var BUCKET_NAME = os.Getenv("BUCKET_NAME")

const (
	FILENAME = "filename"
	AUTHOR   = "author"
	LABEL    = "label"
	TYPE     = "type"
	WORDS    = "words"
)

type S3URLPresigner interface {
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
	PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type AudioService struct {
	s3PresignedAPI S3URLPresigner
}

type IAudioService interface {
	GeneratePreSignedPutURL(filename string, ctx context.Context) (string, error)
	GeneratePreSignedGetURL(string, context.Context) (int, responseBody)
}

func NewAudioService(s S3URLPresigner) IAudioService {
	return &AudioService{
		s3PresignedAPI: s,
	}
}

type responseBody = map[string]interface{}

func (s *AudioService) GeneratePreSignedPutURL(filename string, ctx context.Context) (string, error) {
	request, err := s.s3PresignedAPI.PresignPutObject(ctx, &s3.PutObjectInput{Bucket: aws.String(BUCKET_NAME), Key: aws.String(filename)})

	if err != nil {
		log.Println("An error happened when tried to pre sign a PUT URL", err)
		return "", err
	}

	return request.URL, nil
}

func (s *AudioService) GeneratePreSignedGetURL(param string, ctx context.Context) (int, responseBody) {

	if param == "" {
		return http.StatusBadRequest, responseBody{"message": constant.MISSING_PARAM_ERROR}
	}

	request, err := s.s3PresignedAPI.PresignGetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(BUCKET_NAME), Key: aws.String(param)})

	if err != nil {
		log.Println("An error happened when tried to pre sign a GET URL", err)
		return http.StatusInternalServerError, responseBody{"message": constant.INTERNAL_SERVER_ERROR}
	}

	return http.StatusOK, responseBody{"url": request.URL}

}
