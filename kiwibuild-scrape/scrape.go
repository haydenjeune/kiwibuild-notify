package main

import (
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

// Property represents a single kiwibuild property tile
// fields must be named to match relevant html class, eg Title -> .card__title
type Property struct {
	Title    string
	Location string
	Price    string
	Type     string
	Bed      string
	Bath     string
	Car      string
}

// newFromPropertiesCard creates a new Property given a colly element representing the root
// of it's tile on the KiwiBuild website
func newFromPropertiesCard(e *colly.HTMLElement) *Property {
	var prop Property
	// iterare over fields in Property, extracting the values from the HTML element
	iterOverStruct(&prop, func(f reflect.StructField, v reflect.Value) {
		selector := "div.card__content .card__" + strings.ToLower(f.Name)
		value := toSingleSpaces(e.ChildText(selector))
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
