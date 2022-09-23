package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

func Test_cleanString(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		assert.NotEmpty(t, cleanString(t.Name()))
	})
}

func TestProvideConsumerProvider(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := zerolog.NewZerologLogger()
		cfg := &Config{
			Consumers: MessageQueueConfig{
				Provider: ProviderRedis,
			},
		}

		provider, err := ProvideConsumerProvider(logger, tracing.NewNoopTracerProvider(), cfg)
		assert.NoError(t, err)
		assert.NotNil(t, provider)
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		logger := zerolog.NewZerologLogger()
		cfg := &Config{}

		provider, err := ProvideConsumerProvider(logger, tracing.NewNoopTracerProvider(), cfg)
		assert.Error(t, err)
		assert.Nil(t, provider)
	})
}

func TestProvidePublisherProvider(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := zerolog.NewZerologLogger()
		cfg := &Config{
			Publishers: MessageQueueConfig{
				Provider: ProviderRedis,
			},
		}

		provider, err := ProvidePublisherProvider(logger, tracing.NewNoopTracerProvider(), cfg)
		assert.NoError(t, err)
		assert.NotNil(t, provider)
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		logger := zerolog.NewZerologLogger()
		cfg := &Config{}

		provider, err := ProvidePublisherProvider(logger, tracing.NewNoopTracerProvider(), cfg)
		assert.Error(t, err)
		assert.Nil(t, provider)
	})
}
