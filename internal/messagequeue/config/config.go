package config

import (
	"fmt"
	"strings"

	"github.com/prixfixeco/api_server/internal/messagequeue/consumers"
	redisconsumer "github.com/prixfixeco/api_server/internal/messagequeue/consumers/redis"
	sqsconsumer "github.com/prixfixeco/api_server/internal/messagequeue/consumers/sqs"
	"github.com/prixfixeco/api_server/internal/messagequeue/publishers"
	redispublisher "github.com/prixfixeco/api_server/internal/messagequeue/publishers/redis"
	sqspublisher "github.com/prixfixeco/api_server/internal/messagequeue/publishers/sqs"
	"github.com/prixfixeco/api_server/internal/observability/logging"
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

	// MessageQueueAddress is a simple string alias for the location of our event queue server.
	MessageQueueAddress string

	// RedisConfig configures a Redis-backed consumer.
	RedisConfig struct {
		QueueAddress MessageQueueAddress `json:"message_queue_address,omitempty" mapstructure:"message_queue_address" toml:"message_queue_address,omitempty"`
	}

	// SQSConfig configures a SQS-backed consumer.
	SQSConfig struct {
		QueueAddress MessageQueueAddress `json:"message_queue_address,omitempty" mapstructure:"message_queue_address" toml:"message_queue_address,omitempty"`
	}

	// Config is used to indicate how the messaging provider should be configured.
	Config struct {
		_ struct{}

		Provider    Provider    `json:"provider,omitempty" mapstructure:"provider" toml:"provider,omitempty"`
		RedisConfig RedisConfig `json:"redis,omitempty" mapstructure:"redis" toml:"redis,omitempty"`
		SQSConfig   SQSConfig   `json:"sqs,omitempty" mapstructure:"sqs" toml:"sqs,omitempty"`
	}
)

func cleanString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ProvideConsumerProvider provides a PublisherProvider.
func ProvideConsumerProvider(logger logging.Logger, c *Config) (consumers.ConsumerProvider, error) {
	switch cleanString(string(c.Provider)) {
	case ProviderRedis:
		return redisconsumer.ProvideRedisConsumerProvider(logger, string(c.RedisConfig.QueueAddress)), nil
	case ProviderSQS:
		return sqsconsumer.ProvideSQSConsumerProvider(logger), nil
	default:
		return nil, fmt.Errorf("invalid provider: %q", c.Provider)
	}
}

// ProvidePublisherProvider provides a PublisherProvider.
func ProvidePublisherProvider(logger logging.Logger, c *Config) (publishers.PublisherProvider, error) {
	switch cleanString(string(c.Provider)) {
	case ProviderRedis:
		return redispublisher.ProvideRedisPublisherProvider(logger, string(c.RedisConfig.QueueAddress)), nil
	case ProviderSQS:
		return sqspublisher.ProvideSQSPublisherProvider(logger, string(c.SQSConfig.QueueAddress)), nil
	default:
		return nil, fmt.Errorf("invalid publisher provider: %q", c.Provider)
	}
}
