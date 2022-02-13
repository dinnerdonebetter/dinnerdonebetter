package main

import (
	"context"

	ps "cloud.google.com/go/pubsub"
	"google.golang.org/api/option"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

const (
	serviceAccountFilepath = "cmd/playground/service-account.json"
)

func main() {
	ctx := context.Background()
	logger := zerolog.NewZerologLogger()

	pubsubClient, newClientErr := ps.NewClient(ctx, "prixfixe-dev", option.WithCredentialsFile(serviceAccountFilepath))
	if newClientErr != nil {
		logger.Fatal(newClientErr)
	}

	//tracerProvider := tracing.NewNoopTracerProvider()
	//pp := pubsub.ProvidePubSubPublisherProvider(logger, tracerProvider, pubsubClient)
	//publisher, publisherProviderErr := pp.ProviderPublisher("data_changes")
	//if publisherProviderErr != nil {
	//	logger.Fatal(publisherProviderErr)
	//}

	t := pubsubClient.Topic("data_changes")
	msg := &ps.Message{Data: []byte("fart")}

	result := t.Publish(ctx, msg)
	if _, publishErr := result.Get(ctx); publishErr != nil {
		logger.Fatal(publishErr)
	}

	println()
}
