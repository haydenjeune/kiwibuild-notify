package main

import (
	"encoding/json"
	"log"
	"reflect"
	"regexp"
	"sync"

	"github.com/gocolly/colly/v2"
)

const selectorTag string = "selector"

// Property represents a single KiwiBuild property tile
//
// Fields must be tagged with a colly selector that can be used to retrieve the value from
// the root of the tile.
type Property struct {
	Title    string `selector:".card__title"`
	Location string `selector:".card__location"`
	Price    string `selector:".card__price"`
	Type     string `selector:".card__type"`
	Bed      string `selector:".card__bed"`
	Bath     string `selector:".card__bath"`
	Car      string `selector:".card__car"`
	Status   string `selector:".card__status-label"`
}

// newFromPropertiesCard creates a new Property given a colly element representing the root
// of it's tile on the KiwiBuild website
func newFromPropertiesCard(e *colly.HTMLElement) *Property {
	var prop Property
	// iterare over fields in Property, extracting the values from the HTML element
	iterOverStruct(&prop, func(f reflect.StructField, v reflect.Value) {
		value := toSingleSpaces(e.ChildText(f.Tag.Get(selectorTag)))
		v.SetString(value)
	})
	return &prop
}

// Scraper is the interface that scraping code must satisfy
type Scraper interface {
	Scrape() ([]*Property, error)
}

// CollyKiwiBuildWebScraper uses Colly to extract a list of properties from the KiwiBuild website
type CollyKiwiBuildWebScraper struct {
	url string
}

// Scrape gets list of properties
func (s *CollyKiwiBuildWebScraper) Scrape() ([]*Property, error) {
	properties := make([]*Property, 0, 64)
	mu := new(sync.Mutex) // protects properties
	c := colly.NewCollector()

	c.OnHTML(`div.properties__card`, func(e *colly.HTMLElement) {
		prop := newFromPropertiesCard(e)
		mu.Lock()
		properties = append(properties, prop)
		mu.Unlock()

		propJSON, _ := json.MarshalIndent(prop, "", "  ")
		log.Printf("Found property: %s\n", propJSON)
	})

	c.Visit(s.url)

	return properties, nil
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
