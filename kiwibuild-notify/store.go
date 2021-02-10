package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Storer is used to store properties and check whether the record is new
type Storer interface {
	Store(p *Property) (bool, error)
}

// DynamoDBStorer implements the Storer interface with a DynamoDB backend
type DynamoDBStorer struct {
	svc       *dynamodb.DynamoDB
	tableName *string
}

// Store saves a property to DynamoDB and returns whether the record is new
func (s *DynamoDBStorer) Store(p *Property) (bool, error) {
	av, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		log.Fatalf("error marshalling new item: %v", err)
	}

	input := &dynamodb.PutItemInput{
		Item:         av,
		TableName:    s.tableName,
		ReturnValues: aws.String("ALL_OLD"),
	}

	resp, err := s.svc.PutItem(input)
	if err != nil {
		log.Fatalf("error calling PutItem: %v", err)
	}

	// If the old item had no attributes then the item is new
	new := len(resp.Attributes) == 0
	if new {
		log.Printf("New Property: %v (%v)\n", p.Title, p.Location)
	}

	return new, nil
}
