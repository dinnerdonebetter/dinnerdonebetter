package serversentevents

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding/mock"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	testutils "github.com/dinnerdonebetter/backend/tests/utils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:         logging.NewNoopLogger(),
		dataManager:    database.NewMockDatabase(),
		encoderDecoder: mockencoding.NewMockEncoderDecoder(),
		tracer:         tracing.NewTracerForTest("test"),
		cfg:            &Config{},
	}
}

func TestProvideValidVesselsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := &Config{
			DataChangesTopicName: "data_changes",
		}

		mockConsumer := &mockpublishers.Consumer{}
		mockConsumer.On("Consume", mock.Anything, mock.Anything).Return(nil)

		pp := &mockpublishers.ConsumerProvider{}
		pp.On("ProvideConsumer", testutils.ContextMatcher, cfg.DataChangesTopicName, mock.Anything).Return(mockConsumer, nil)

		s, err := ProvideService(
			ctx,
			cfg,
			logger,
			database.NewMockDatabase(),
			mockencoding.NewMockEncoderDecoder(),
			pp,
			tracing.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing data changes consumer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := &Config{
			DataChangesTopicName: "data_changes",
		}

		pp := &mockpublishers.ConsumerProvider{}
		pp.On("ProvideConsumer", testutils.ContextMatcher, cfg.DataChangesTopicName, mock.Anything).Return((*mockpublishers.Consumer)(nil), errors.New("blah"))

		s, err := ProvideService(
			ctx,
			cfg,
			logger,
			database.NewMockDatabase(),
			mockencoding.NewMockEncoderDecoder(),
			pp,
			tracing.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
