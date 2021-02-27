package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/cenkalti/backoff/v4"
)

// Storer is used to store properties in some backend
type Storer interface {
	Store(p []*Property) error
}

// DynamoDBStorer implements the Storer interface with a DynamoDB backend
type DynamoDBStorer struct {
	svc       *dynamodb.DynamoDB
	tableName *string
}

// Store saves list of properties to DynamoDB
func (s *DynamoDBStorer) Store(p []*Property) error {
	items, err := s.buildRequestItems(p)
	if err != nil {
		return fmt.Errorf("failed to build request: %v", err)
	}

	retryPolicy := backoff.NewExponentialBackOff()
	retryPolicy.MaxElapsedTime = 5 * time.Second

	err = backoff.Retry(func() error {
		req := &dynamodb.BatchWriteItemInput{
			RequestItems: items,
		}
		resp, err := s.svc.BatchWriteItem(req)
		if err != nil {
			// unrecoverable error
			return backoff.Permanent(err)
		}
		if len(resp.UnprocessedItems) > 0 {
			// retry unprocessed items
			items = resp.UnprocessedItems
			return fmt.Errorf("unprocessed items remaining [%v]", items)
		}
		return nil
	}, retryPolicy)
	if err != nil {
		return fmt.Errorf("failed to write items to dynamo: %v", err)
	}

	return nil
}

// buildRequestItems builds a dynamodb batch write request to put each property to the table
func (s *DynamoDBStorer) buildRequestItems(p []*Property) (map[string][]*dynamodb.WriteRequest, error) {
	reqs := make([]*dynamodb.WriteRequest, len(p))
	for i, val := range p {
		item, err := dynamodbattribute.MarshalMap(val)
		if err != nil {
			return nil, fmt.Errorf("failed to marshall property '%s' to attribute value: %v", val.Title, err)
		}
		reqs[i] = &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{Item: item},
		}
	}
	var reqItems = map[string][]*dynamodb.WriteRequest{*s.tableName: reqs}
	return reqItems, nil
}
