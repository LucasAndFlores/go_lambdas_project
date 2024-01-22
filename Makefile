MY_AWS_PROFILE = "your-local-aws-profile"
CODE_BUCKET="your-bucket"
PROJECT_NAME="goLambdasProject"

GOOS=linux
GOARCH=amd64
CGO_ENABLED=0

run:
	AWS_PROFILE=${MY_AWS_PROFILE} sam local start-api --env-vars env.json

dev:
	make build && make run

fmt:
	gofmt -w ./cmd ./config ./internal 

vet:
	go vet ./cmd/... ./config/... ./internal/...  

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 sam build

unit-test-internal: 
	go test ./internal/... -coverprofile=cover.out

coverage-report:
	go tool cover -html=cover.out

package:
	AWS_PROFILE=${MY_AWS_PROFILE} sam package \
		--s3-bucket ${CODE_BUCKET} \
		--output-template-file ./.aws-sam/packaged.yaml

deploy:
	AWS_PROFILE=${MY_AWS_PROFILE} make build && make package && sam deploy --no-confirm-changeset \
		--no-fail-on-empty-changeset \
		--s3-bucket ${CODE_BUCKET} \
		--stack-name ${PROJECT_NAME} \
		--capabilities CAPABILITY_IAM
