package websockets

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockconsumers "github.com/prixfixeco/backend/internal/messagequeue/mock"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/internal/encoding"
	mockencoding "github.com/prixfixeco/backend/internal/encoding/mock"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	testutils "github.com/prixfixeco/backend/tests/utils"
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
		logger := logging.NewNoopLogger()
		encoder := encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		consumer := &mockconsumers.Consumer{}
		consumer.On("Consume", chan bool(nil), chan error(nil))

		consumerProvider := &mockconsumers.ConsumerProvider{}
		consumerProvider.On(
			"ProvideConsumer",
			testutils.ContextMatcher,
			topicName,
			mock.Anything,
		).Return(consumer, nil)

		actual, err := ProvideService(
			ctx,
			authCfg,
			logger,
			encoder,
			consumerProvider,
			tracing.NewNoopTracerProvider(),
		)

		require.NoError(t, err)
		require.NotNil(t, actual)

		mock.AssertExpectationsForObjects(t, consumerProvider)
	})

	T.Run("with consumer provider error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		authCfg := &authservice.Config{}
		logger := logging.NewNoopLogger()
		encoder := encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		consumerProvider := &mockconsumers.ConsumerProvider{}
		consumerProvider.On(
			"ProvideConsumer",
			testutils.ContextMatcher,
			topicName,
			mock.Anything,
		).Return(&mockconsumers.Consumer{}, errors.New("blah"))

		actual, err := ProvideService(
			ctx,
			authCfg,
			logger,
			encoder,
			consumerProvider,
			tracing.NewNoopTracerProvider(),
		)

		require.Error(t, err)
		require.Nil(t, actual)
	})
}

func Test_buildWebsocketErrorFunc(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		encoder := encoding.ProvideServerEncoderDecoder(nil, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		buildWebsocketErrorFunc(encoder)(res, req, 200, errors.New("blah"))
	})
}
