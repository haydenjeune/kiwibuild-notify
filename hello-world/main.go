package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
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

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	c.OnHTML(`div.properties__card`, func(e *colly.HTMLElement) {
		prop := NewFromPropertiesCard(e)
		// Print text
		fmt.Printf("Found Property: %v\n", prop)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on kiwibuild.govt.nz
	c.Visit("https://kiwibuild.govt.nz/available-homes/")

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hi"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
