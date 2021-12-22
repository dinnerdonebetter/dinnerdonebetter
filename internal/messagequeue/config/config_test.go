package config

import (
	"testing"

	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"

	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"
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
			Consumers: ProviderConfig{
				Provider: ProviderRedis,
			},
		}

		provider, err := ProvideConsumerProvider(logger, trace.NewNoopTracerProvider(), cfg)
		assert.NoError(t, err)
		assert.NotNil(t, provider)
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		logger := zerolog.NewZerologLogger()
		cfg := &Config{}

		provider, err := ProvideConsumerProvider(logger, trace.NewNoopTracerProvider(), cfg)
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
			Publishers: ProviderConfig{
				Provider: ProviderRedis,
			},
		}

		provider, err := ProvidePublisherProvider(logger, trace.NewNoopTracerProvider(), cfg)
		assert.NoError(t, err)
		assert.NotNil(t, provider)
	})

	T.Run("with invalid provider", func(t *testing.T) {
		t.Parallel()

		logger := zerolog.NewZerologLogger()
		cfg := &Config{}

		provider, err := ProvidePublisherProvider(logger, trace.NewNoopTracerProvider(), cfg)
		assert.Error(t, err)
		assert.Nil(t, provider)
	})
}
