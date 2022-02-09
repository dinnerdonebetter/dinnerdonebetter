package main

import (
	"context"
	"log"
)

// PubSubMessage is the payload of a Pub/Sub event. See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

// ProcessDataChange handles a data change
func ProcessDataChange(ctx context.Context, m PubSubMessage) error {
	log.Println("invoked")

	return nil
}
