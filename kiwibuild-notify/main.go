package main

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gocolly/colly/v2"
)

// Property represents a single kiwibuild property tile
type Property struct {
	Title    string
	Location string
	Price    string
	Type     string
	Bed      string
	Bath     string
	Car      string
}

// NewFromPropertiesCard creates a new Property given a colly element representing the root
// of it's tile on the KiwiBuild website
func NewFromPropertiesCard(e *colly.HTMLElement) *Property {
	var prop Property
	iterOverStruct(&prop, func(f reflect.StructField, v reflect.Value) {
		selector := "div.card__content .card__" + strings.ToLower(f.Name)
		value := toSingleSpaces(e.ChildText(selector))
		v.SetString(value)
	})
	return &prop
}

// iterOverStruct iterates over each field in a given struct, executing a callback for each
// field. The callback is passed the relevant reflect.StructField and reflect.Value
func iterOverStruct(i interface{}, f func(field reflect.StructField, value reflect.Value)) {
	// TODO: validate structure of interface
	v := reflect.ValueOf(i).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f(t.Field(i), v.Field(i))
	}
}

// toSingleSpaces removes replaces all continuous whitespace with a single space
func toSingleSpaces(s string) string {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(s, " ")
}

func persist(p *Property) {
	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cfg := aws.NewConfig()

	if endpoint := os.Getenv("DYNAMODB_ENDPOINT"); endpoint != "" {
		fmt.Println("Using local endpoint: " + endpoint)
		cfg.WithEndpoint(endpoint)
	}

	// Create DynamoDB client
	svc := dynamodb.New(sess, cfg)

	av, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		fmt.Println("Got error marshalling new item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Create item in table Property
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println("Got error calling PutItem:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Printf("Successfully added %v\n", p)
}

func handler() error {

	c := colly.NewCollector()

	// On every property card call callback
	c.OnHTML(`div.properties__card`, func(e *colly.HTMLElement) {
		prop := NewFromPropertiesCard(e)
		// Print text
		fmt.Printf("Found Property: %v\n", prop)
		persist(prop)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on kiwibuild.govt.nz
	c.Visit("https://kiwibuild.govt.nz/available-homes/")

	return nil
}

func main() {
	lambda.Start(handler)
}
