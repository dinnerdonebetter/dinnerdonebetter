package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	_ = "https://sqs.us-east-1.amazonaws.com/966107642521/data_changes.fifo"
)

func buildHandler() func(ctx context.Context, sqsEvent events.SQSEvent) error {
	return func(ctx context.Context, sqsEvent events.SQSEvent) error {
		for i := 0; i < len(sqsEvent.Records); i++ {
			message := sqsEvent.Records[i]
			log.Println(message.Body)
		}

		return nil
	}
}

func main() {
	lambda.Start(buildHandler())
}
