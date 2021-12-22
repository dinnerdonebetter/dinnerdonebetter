package config

import (
	"fmt"
	"strings"

	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/messagequeue/redis"
	"github.com/prixfixeco/api_server/internal/messagequeue/sqs"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	// ProviderRedis is used to refer to redis.
	ProviderRedis = "redis"
	// ProviderSQS is used to refer to sqs.
	ProviderSQS = "sqs"
)

type (
	// Provider is used to indicate what messaging provider we'll use.
	Provider string

	// RedisConfig configures a Redis-backed consumer.
	RedisConfig struct {
		QueueAddress string `json:"message_queue_address,omitempty" mapstructure:"message_queue_address" toml:"message_queue_address,omitempty"`
	}

	// SQSConfig configures a SQS-backed consumer.
	SQSConfig struct {
		QueueAddress string `json:"message_queue_address,omitempty" mapstructure:"message_queue_address" toml:"message_queue_address,omitempty"`
	}

	// ProviderConfig is used to indicate how the messaging provider should be configured.
	ProviderConfig struct {
		_           struct{}
		Provider    Provider     `json:"provider,omitempty" mapstructure:"provider" toml:"provider,omitempty"`
		SQSConfig   sqs.Config   `json:"sqs,omitempty" mapstructure:"sqs" toml:"sqs,omitempty"`
		RedisConfig redis.Config `json:"redis,omitempty" mapstructure:"redis" toml:"redis,omitempty"`
	}

	// Config is used to indicate how the messaging provider should be configured.
	Config struct {
		_          struct{}
		Consumers  ProviderConfig `json:"consumers,omitempty" mapstructure:"consumers" toml:"consumers,omitempty"`
		Publishers ProviderConfig `json:"publishers,omitempty" mapstructure:"publishers" toml:"publishers,omitempty"`
	}
)

func cleanString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ProvideConsumerProvider provides a PublisherProvider.
func ProvideConsumerProvider(logger logging.Logger, tracerProvider tracing.TracerProvider, c *Config) (messagequeue.ConsumerProvider, error) {
	switch cleanString(string(c.Consumers.Provider)) {
	case ProviderRedis:
		return redis.ProvideRedisConsumerProvider(logger, tracerProvider, c.Consumers.RedisConfig), nil
	case ProviderSQS:
		return sqs.ProvideSQSConsumerProvider(logger, tracerProvider), nil
	default:
		return nil, fmt.Errorf("invalid provider: %q", c.Consumers.Provider)
	}
}

// ProvidePublisherProvider provides a PublisherProvider.
func ProvidePublisherProvider(logger logging.Logger, tracerProvider tracing.TracerProvider, c *Config) (messagequeue.PublisherProvider, error) {
	switch cleanString(string(c.Publishers.Provider)) {
	case ProviderRedis:
		return redis.ProvideRedisPublisherProvider(logger, tracerProvider, c.Publishers.RedisConfig), nil
	case ProviderSQS:
		return sqs.ProvideSQSPublisherProvider(logger, tracerProvider), nil
	default:
		return nil, fmt.Errorf("invalid publisher provider: %q", c.Publishers.Provider)
	}
}
