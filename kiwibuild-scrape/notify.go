package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

// Notifier provides an interface to send a notification with a given message
type Notifier interface {
	Notify(string) error
}

// SNSNotifier allows notifications to be sent to a particular SNS topic
type SNSNotifier struct {
	svc   *sns.SNS
	topic sns.Topic
}

// Notify publishes a message to the SNSNotifiers internal topic
func (n *SNSNotifier) Notify(message string) error {
	input := &sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: n.topic.TopicArn,
	}

	_, err := n.svc.Publish(input)
	if err != nil {
		log.Fatalf("Publish error: %v", err)
	}

	return nil
}
