package workers

import (
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/business/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/encoding/mock"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:         logging.NewNoopLogger(),
		encoderDecoder: encoding.ProvideServerEncoderDecoder(nil, nil, encoding.ContentTypeJSON),
		tracer:         tracing.NewTracerForTest("test"),
	}
}

func TestProvideService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		cfg := &msgconfig.QueuesConfig{
			DataChangesTopicName: "data_changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			logger,
			database.NewMockDatabase(),
			mockencoding.NewMockEncoderDecoder(),
			pp,
			tracing.NewNoopTracerProvider(),
			metrics.NewNoopMetricsProvider(),
			&recipeanalysis.MockRecipeAnalyzer{},
			cfg,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing data changes producer", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		cfg := &msgconfig.QueuesConfig{
			DataChangesTopicName: "data_changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			logger,
			database.NewMockDatabase(),
			mockencoding.NewMockEncoderDecoder(),
			pp,
			tracing.NewNoopTracerProvider(),
			metrics.NewNoopMetricsProvider(),
			&recipeanalysis.MockRecipeAnalyzer{},
			cfg,
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
