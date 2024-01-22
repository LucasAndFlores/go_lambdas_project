# Go lambdas
These go lambdas are responsible for returning pre-signed PUT and GET URLs from AWS S3, reacting when an object is stored and storing its metadata in a Postgres, and listing all this stored metadata.

## Motivation to build this project
The project is being built to deliver audio for a mobile app. This mobile application is called "Medo e delirio em brasilia", and it was created by persons that hear a podcast with the same name, very famous for "audio insertions" inside of podcast. The app creators would like all people who hear the podcast to have a chance to share these "audio insertions" with their friends and have some fun.

## Usage and setup for environment
Before you start, be sure that you have AWS CLI, Docker, and Go installed. 

As an initial step, you should log into your AWS account and create two buckets, one of them will store your code and the other will contain all audios.

Inside the Makefile, fill the local variable `MY_AWS_PROFILE` with your local AWS profile name, and define which bucket you will store the code in the variable `CODE_BUCKET`. You can define a new project name, if you want, at the variable `PROJECT_NAME`.

### Local environment
To run and test the AWS Lambdas locally, you can run:
```bash
make dev
```

#### Routes 

`POST /audio`

This route is responsible for storing audio and metadata in an AWS bucket. 

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

Expected response: 
```json
{
	"url": "http://aws.url"

}
```

`GET /audio/:filename`

This route will return a S3 pre-signed URL, where you can download the file. 

Expected response: 
```json
{
	"url": "http://aws.url"
}
```

##### Next steps
- [ ] Build a lambda for reacting and storing all metadata from the object stored in S3 in a database
- [ ] Build a lambda to return all the metadata stored in the database.



 


