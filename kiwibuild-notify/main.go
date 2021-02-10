package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gocolly/colly/v2"
)

// TODO: make a struct to inject into to avoid magical access across files
var svc *dynamodb.DynamoDB

func init() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cfg := aws.NewConfig()

	if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
		log.Println("Using local endpoint: " + endpoint)
		cfg.WithEndpoint(endpoint)
	}

	svc = dynamodb.New(sess, cfg)
}

func handler() {
	s := &DynamoDBStorer{svc: svc, tableName: aws.String(os.Getenv("DYNAMODB_TABLE_NAME"))}
	c := colly.NewCollector()
	c.OnHTML(`div.properties__card`, func(e *colly.HTMLElement) {
		ProcessPropertyCard(e, s)
	})
	c.Visit("https://kiwibuild.govt.nz/available-homes/")
}

func main() {
	lambda.Start(handler)
}
