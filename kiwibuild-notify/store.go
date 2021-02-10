package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Storer is used to store properties and checks whether they previously existed or not
type Storer interface {
	Store(p *Property) (bool, error)
}

// DynamoDBStorer implements the Storer interface with a DynamoDB backend
type DynamoDBStorer struct {
	svc       *dynamodb.DynamoDB
	tableName *string
}

// Store saves a property to DynamoDB and returns whether it previously existed
func (s *DynamoDBStorer) Store(p *Property) (bool, error) {
	av, err := dynamodbattribute.MarshalMap(p)
	if err != nil {
		fmt.Println("Got error marshalling new item:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: s.tableName,
	}

	_, err = s.svc.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: &v", err)
	}

	fmt.Printf("Successfully added %v\n", p)

	return false, nil
}
