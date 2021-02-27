package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var scraper Scraper
var storer Storer

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

	scraper = &CollyKiwiBuildWebScraper{
		url: os.Getenv("KIWIBUILD_URL"),
	}
	storer = &DynamoDBStorer{
		svc:       dynamodb.New(sess, cfg),
		tableName: aws.String(os.Getenv("DYNAMODB_TABLE_NAME")),
	}
}

func handler(scraper Scraper, storer Storer) {
	properties, err := scraper.Scrape()
	if err != nil {
		log.Fatalf("failed to scrape properties: %v", err)
	}

	err = storer.Store(properties)
	if err != nil {
		log.Fatalf("failed to store properties: %v", err)
	}
}

func main() {
	lambda.Start(func() {
		handler(scraper, storer)
	})
}
