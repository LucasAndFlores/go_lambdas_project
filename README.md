# Go lambdas
These go lambdas are responsible for returning pre-signed PUT and GET URLs from AWS S3, storing their metadata in a DynamoDB, and listing all this stored metadata.

## Motivation to build this project
The project is being built to deliver audio for a mobile app. This mobile application is called "Medo e delirio em brasilia", and it was created by persons that hear a podcast with the same name, very famous for "audio insertions" inside of podcast. The app creators would like all people who hear the podcast to have a chance to share these "audio insertions" with their friends and have some fun.

## Usage and setup for environment
Before you start, be sure that you have AWS CLI, AWS SAM, Docker, and Go installed. 

As an initial step, you should log into your AWS account and create a bucket, which will store your code.

Inside the Makefile, fill the local variable `MY_AWS_PROFILE` with your local AWS profile name, and the variable `CODE_BUCKET` with the bucket that you created to store your code. You can define a new project name, if you want, at the variable `PROJECT_NAME`.

After this, go to the file `template.yaml` and fill in the parameters `BucketName` and `DynamoTableName` to create a new bucket and dynamodb table. These two resources will be used in this API to store the files and the metadata.

You should run this command to deploy the API and also to create the dynamodb table and S3 bucket: 
```bash
make deploy
```

### Local environment
To run and test the AWS Lambdas locally, you can run:

```bash
make dev
```

#### Unit testing
If you want to run the unit testing, run:
```bash
make unit-test-internal
```

You can also check the coverage report with this command: 
```bash
make coverage-report
```

#### Routes 

`POST /audio`

This route generates a pre-signed S3 URL to store the object.

A body object is required, example:
```json
{
	"fileName": "test", 
}
```

Request: 
```bash
curl -X POST -H "Content-Type: application/json" -d '{
  "fileName": "test",
}' http://localhost:3000/audio
```

Expected responses:

Status Code: 201 <br>
Body:
```json
{
	"url": "http://aws.url"

}
```

Status Code: 400 <br>
Reason: Invalid JSON or the property `fileName` is missing <br>
Body:
```json
{
	"message": "Unable to process the body. Please, review the content"
}
```

Status Code: 500 <br>
Reason: An error occurred with S3. <br>
Body:
```json
{
	"message": "internal server error"
}
```

`GET /audio/:filename`

This route will return a S3 pre-signed URL, where you can download the file. 

Request: 
```bash
curl http://localhost:3000/audio/test
```

Expected responses:

Status: 200 <br>
Body:
```json
{
	"url": "http://aws.url"
}
```

Status Code: 400 <br>
Reason: The `fileName` parameter is missing <br>
Body:
```json
{
	"message": "Missing required parameter"
}
```

Status Code: 500 <br>
Reason: An error occurred with S3. <br>
Body:
```json
{
	"message": "internal server error"
}
```

`POST /metadata`

This route will store the metadata related to a file stored in S3. The `fileName` property should be stored on S3 with the same value and it is unique.

A body object is required, example: 
```json
{
	"fileName": "test", 
	"author": "test",
	"label": "test",
	"type": "test",
	"words": "test"
}
```

Request: 
```bash
curl -X POST -H "Content-Type: application/json" -d '{
	"fileName": "test", 
	"author": "test",
	"label": "test",
	"type": "test",
	"words": "test"
}' http://localhost:3000/metadata
```

Expected responses:

Status: 201 <br>
Body:
```json
{
	"message": "successfully stored"
}
```

Status Code: 400 <br>
Reason: The request body probably has a bad syntax <br>
Body:
```json
{
	"message": "Unable to process the body. Please, review the content"
}
```

Status Code: 400 <br>
Reason: Invalid type or field is missing on JSON property <br>
Body:
```json
{
   "errors":[
      {
         "field":"Words",
         "tag":"required",
         "value":""
      }
   ]
}
```

Status Code: 409 <br>
Reason: The metadata already exists on database <br>
Body:
```json
{
   "message": "The object already exists"
}
```

Status Code: 422 <br>
Reason: The JSON fields and types are valid, but the audio is not on S3. <br>
Body:
```json
{
   "message": "Filename not found. Unable to complete the operation"
}
```

Status Code: 500 <br>
Reason: An internal error happened. <br>
Body:
```json
{
	"message": "internal server error"
}
```

`GET /metadata`

A body object is required, example:This route will list all metadata stored on dynamoDB. 

Request: 
```bash
curl http://localhost:3000/metadata
```

Expected responses:

Status: 200 <br>
Body:
```json
{
   "metadata":[
      {
         "filename":"test",
         "author":"test",
         "label":"test",
         "type":"string",
         "words":"test"
      }
   ]
}
```

Status Code: 500 <br>
Reason: An internal error happened. <br>
Body:
```json
{
	"message": "internal server error"
}
```



 


