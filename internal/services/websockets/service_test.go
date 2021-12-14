package websockets

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/encoding"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockconsumers "github.com/prixfixeco/api_server/internal/messagequeue/consumers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func buildTestService() *service {
	return &service{
		logger:         logging.NewNoopLogger(),
		encoderDecoder: mockencoding.NewMockEncoderDecoder(),
		tracer:         tracing.NewTracerForTest("test"),
		connections:    map[string][]websocketConnection{},
		authConfig: &authservice.Config{
			Cookies: authservice.CookieConfig{
				Name: "cookie",
			},
		},
	}
}

func TestProvideService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		authCfg := &authservice.Config{}
		cfg := Config{}
		logger := logging.NewNoopLogger()
		encoder := encoding.ProvideServerEncoderDecoder(logger, trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		consumer := &mockconsumers.Consumer{}
		consumer.On("Consume", chan bool(nil), chan error(nil))

		consumerProvider := &mockconsumers.ConsumerProvider{}
		consumerProvider.On(
			"ProvideConsumer",
			testutils.ContextMatcher,
			cfg.DataChangesTopicName,
			mock.Anything,
		).Return(consumer, nil)

		actual, err := ProvideService(
			ctx,
			authCfg,
			cfg,
			logger,
			encoder,
			consumerProvider,
			trace.NewNoopTracerProvider(),
		)

		require.NoError(t, err)
		require.NotNil(t, actual)

		mock.AssertExpectationsForObjects(t, consumerProvider)
	})

	T.Run("with consumer provider error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		authCfg := &authservice.Config{}
		cfg := Config{}
		logger := logging.NewNoopLogger()
		encoder := encoding.ProvideServerEncoderDecoder(logger, trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		consumerProvider := &mockconsumers.ConsumerProvider{}
		consumerProvider.On(
			"ProvideConsumer",
			testutils.ContextMatcher,
			cfg.DataChangesTopicName,
			mock.Anything,
		).Return(&mockconsumers.Consumer{}, errors.New("blah"))

		actual, err := ProvideService(
			ctx,
			authCfg,
			cfg,
			logger,
			encoder,
			consumerProvider,
			trace.NewNoopTracerProvider(),
		)

		require.Error(t, err)
		require.Nil(t, actual)
	})
}

func Test_buildWebsocketErrorFunc(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		encoder := encoding.ProvideServerEncoderDecoder(nil, trace.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		buildWebsocketErrorFunc(encoder)(res, req, 200, errors.New("blah"))
	})
}
