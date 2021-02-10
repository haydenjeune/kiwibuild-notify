package main

import (
	"log"
	"reflect"
	"regexp"
	"strings"

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
	iterOverStruct(&prop, func(f reflect.StructField, v reflect.Value) {
		selector := "div.card__content .card__" + strings.ToLower(f.Name)
		value := toSingleSpaces(e.ChildText(selector))
		v.SetString(value)
	})
	return &prop
}

// ProcessPropertyCard coordinates the parsing, storing, and possible notification associated with each property card
func ProcessPropertyCard(e *colly.HTMLElement, s Storer) {
	prop := newFromPropertiesCard(e)
	log.Printf("Found Property: %v (%v)\n", prop.Title, prop.Location)
	s.Store(prop)
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
