package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gocolly/colly/v2"
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
	s := &DynamoDBStorer{svc: dynamoClient, tableName: aws.String(os.Getenv("DYNAMODB_TABLE_NAME"))}
	n := &SNSNotifier{svc: snsClient, topic: sns.Topic{TopicArn: aws.String(os.Getenv("SNS_NOTIFICATION_TOPIC"))}}

	c := colly.NewCollector()
	c.OnHTML(`div.properties__card`, func(e *colly.HTMLElement) {
		ProcessPropertyCard(e, s, n)
	})
	c.Visit("https://kiwibuild.govt.nz/available-homes/")
}

func main() {
	lambda.Start(handler)
}
