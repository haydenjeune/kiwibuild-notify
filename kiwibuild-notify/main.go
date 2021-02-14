package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

var notifier SNSNotifier

func init() {
	sess := session.Must(session.NewSession())
	notifier = SNSNotifier{svc: sns.New(sess), topic: sns.Topic{TopicArn: aws.String(os.Getenv("SNS_NOTIFICATION_TOPIC"))}}
}

// handler sends a formatted notification when the given event is an insert
func handler(e events.DynamoDBEvent, n *SNSNotifier) {
	for _, record := range e.Records {
		if record.EventName == "INSERT" {
			fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)
			msg := fmt.Sprintf("A new KiwiBuild property '%v' in %v has been listed! %v priced %v",
				record.Change.NewImage["Title"],
				record.Change.NewImage["Location"],
				record.Change.NewImage["Type"],
				record.Change.NewImage["Price"],
			)
			n.Notify(msg)
		}
	}
}

func main() {
	lambda.Start(func(e events.DynamoDBEvent) {
		handler(e, &notifier)
	})
}
