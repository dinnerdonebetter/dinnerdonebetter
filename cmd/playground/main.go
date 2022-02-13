package main

import (
	"context"
	"time"

	"github.com/prixfixeco/api_server/internal/observability/logging"

	"github.com/prixfixeco/api_server/internal/messagequeue/pubsub"
	"github.com/prixfixeco/api_server/internal/observability/tracing"

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
	logger.SetLevel(logging.DebugLevel)

	pubsubClient, newClientErr := ps.NewClient(ctx, "prixfixe-dev", option.WithCredentialsFile(serviceAccountFilepath))
	if newClientErr != nil {
		logger.Fatal(newClientErr)
	}

	tracerProvider := tracing.NewNoopTracerProvider()
	pp := pubsub.ProvidePubSubPublisherProvider(logger, tracerProvider, pubsubClient)
	publisher, publisherProviderErr := pp.ProviderPublisher("data_changes")
	if publisherProviderErr != nil {
		logger.Fatal(publisherProviderErr)
	}

	//publisher := pubsubClient.Topic("data_changes")
	//msg := &ps.Message{Data: []byte("fart")}
	//result := publisher.Publish(ctx, msg)
	//if _, publishErr := result.Get(ctx); publishErr != nil {
	//	logger.Fatal(publishErr)
	//}

	for {
		select {
		case <-time.NewTicker(time.Second).C:
			if publishErr := publisher.Publish(ctx, []byte("fart")); publishErr != nil {
				logger.Fatal(publishErr)
			}
		}
	}

	println()
}
