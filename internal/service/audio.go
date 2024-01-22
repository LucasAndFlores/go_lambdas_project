package service

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/LucasAndFlores/go_lambdas_project/constant"
	"github.com/LucasAndFlores/go_lambdas_project/internal/entity"
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
	PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

type AudioService struct {
	s3PresignedAPI S3URLPresigner
}

type IAudioService interface {
	GeneratePreSignedPutURL(string, context.Context) (int, responseBody)
}

func NewAudioService(s S3URLPresigner) IAudioService {
	return &AudioService{
		s3PresignedAPI: s,
	}
}

type responseBody = map[string]interface{}

func (s *AudioService) GeneratePreSignedPutURL(body string, ctx context.Context) (int, responseBody) {
	var parsedBody entity.AudioDTOInput

	err := json.Unmarshal([]byte(body), &parsedBody)

	if err != nil {
		return http.StatusUnprocessableEntity, responseBody{"message": constant.INTERNAL_SERVER_ERROR}

	}

	validationErr := parsedBody.Validate()

	if validationErr != nil {
		return http.StatusUnprocessableEntity, responseBody{"errors": validationErr}
	}

	request, err := s.s3PresignedAPI.PresignPutObject(ctx, &s3.PutObjectInput{Bucket: aws.String(BUCKET_NAME), Key: aws.String(parsedBody.FileName), Metadata: map[string]string{
		FILENAME: parsedBody.FileName, AUTHOR: parsedBody.Author, LABEL: parsedBody.Label, TYPE: parsedBody.Type, WORDS: parsedBody.Words,
	}})

	if err != nil {
		log.Println("An error happened when tried to pre sign a PUT URL", err)
		return http.StatusInternalServerError, responseBody{"message": constant.INTERNAL_SERVER_ERROR}
	}

	return http.StatusCreated, responseBody{"url": request.URL}
}
