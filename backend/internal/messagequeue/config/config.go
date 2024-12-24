package config

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/messagequeue/pubsub"
	"github.com/dinnerdonebetter/backend/internal/messagequeue/redis"
	"github.com/dinnerdonebetter/backend/internal/messagequeue/sqs"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	ps "cloud.google.com/go/pubsub"
)

const (
	// ProviderRedis is used to refer to redis.
	ProviderRedis provider = "redis"
	// ProviderSQS is used to refer to sqs.
	ProviderSQS provider = "sqs"
	// ProviderPubSub is used to refer to GCP Pub/Sub.
	ProviderPubSub provider = "pubsub"
)

var (
	ErrNilConfig = errors.New("nil config provided")
)

type (
	// provider is used to indicate what messaging provider we'll use.
	provider string

	// MessageQueueConfig is used to indicate how the messaging provider should be configured.
	MessageQueueConfig struct {
		_ struct{} `json:"-"`

		Provider provider      `env:"PROVIDER"      json:"provider,omitempty"`
		SQS      sqs.Config    `envPrefix:"SQS_"    json:"sqs,omitempty"`
		PubSub   pubsub.Config `envPrefix:"PUBSUB_" json:"pubSub,omitempty"`
		Redis    redis.Config  `envPrefix:"REDIS_"  json:"redis,omitempty"`
	}

	// Config is used to indicate how the messaging provider should be configured.
	Config struct {
		_ struct{} `json:"-"`

		Consumer  MessageQueueConfig `envPrefix:"CONSUMER_"  json:"consumers,omitempty"`
		Publisher MessageQueueConfig `envPrefix:"PUBLISHER_" json:"publishers,omitempty"`
	}
)

func cleanString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ProvideConsumerProvider provides a PublisherProvider.
func ProvideConsumerProvider(ctx context.Context, logger logging.Logger, c *Config) (messagequeue.ConsumerProvider, error) {
	if c == nil {
		return nil, ErrNilConfig
	}

	switch cleanString(string(c.Consumer.Provider)) {
	case string(ProviderRedis):
		return redis.ProvideRedisConsumerProvider(logger, c.Consumer.Redis), nil
	case string(ProviderPubSub):
		client, err := ps.NewClient(ctx, c.Consumer.PubSub.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("establishing PubSub client: %w", err)
		}

		return pubsub.ProvidePubSubConsumerProvider(logger, client), nil
	default:
		return nil, fmt.Errorf("invalid provider: %q", c.Consumer.Provider)
	}
}

// ProvidePublisherProvider provides a PublisherProvider.
func ProvidePublisherProvider(ctx context.Context, logger logging.Logger, tracerProvider tracing.TracerProvider, c *Config) (messagequeue.PublisherProvider, error) {
	if c == nil {
		return nil, ErrNilConfig
	}

	switch cleanString(string(c.Publisher.Provider)) {
	case string(ProviderRedis):
		return redis.ProvideRedisPublisherProvider(logger, tracerProvider, c.Publisher.Redis), nil
	case string(ProviderSQS):
		return sqs.ProvideSQSPublisherProvider(logger, tracerProvider), nil
	case string(ProviderPubSub):
		client, err := ps.NewClient(ctx, c.Publisher.PubSub.ProjectID)
		if err != nil {
			return nil, fmt.Errorf("establishing PubSub client: %w", err)
		}

		return pubsub.ProvidePubSubPublisherProvider(logger, tracerProvider, client), nil
	default:
		return nil, fmt.Errorf("invalid publisher provider: %q", c.Publisher.Provider)
	}
}
