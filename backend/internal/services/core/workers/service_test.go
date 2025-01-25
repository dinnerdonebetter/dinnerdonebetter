package workers

import (
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/businesslogic/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding/mock"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

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
