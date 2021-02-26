package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
)

var dynamoClient *dynamodb.DynamoDB
var snsClient *sns.SNS

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cfg := aws.NewConfig()

	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		endpoint := os.Getenv("DYNAMODB_TEST_ENDPOINT")
		log.Println("Using local endpoint: " + endpoint)
		cfg.WithEndpoint(endpoint)
	}

	dynamoClient = dynamodb.New(sess, cfg)
	snsClient = sns.New(sess)
}

func handler() {
	// TODO: init these in exec context and inject into handler?
	//s := &DynamoDBStorer{svc: dynamoClient, tableName: aws.String(os.Getenv("DYNAMODB_TABLE_NAME"))}

	s := &CollyKiwiBuildWebScraper{url: "https://kiwibuild.govt.nz/available-homes/"}
	p, _ := s.Scrape()
	for _, v := range p {
		fmt.Println(*v)
	}
}

func main() {
	lambda.Start(handler)
}
