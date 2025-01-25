package mealplanning

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding/mock"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	mockrouting "github.com/dinnerdonebetter/backend/internal/lib/routing/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:                  logging.NewNoopLogger(),
		mealPlanningDataManager: mocktypes.NewMealPlanningDataManagerMock(),
		mealIDFetcher:           func(req *http.Request) string { return "" },
		encoderDecoder:          encoding.ProvideServerEncoderDecoder(nil, nil, encoding.ContentTypeJSON),
		tracer:                  tracing.NewTracerForTest("test"),
	}
}

func TestProvideMealsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()

		for _, k := range allURIKeys {
			rpm.On("BuildRouteParamStringIDFetcher", k).Return(func(*http.Request) string { return "" })
		}

		cfg := &Config{}
		msgCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", msgCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			context.Background(),
			logging.NewNoopLogger(),
			cfg,
			&textsearchcfg.Config{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			tracing.NewNoopTracerProvider(),
			metrics.NewNoopMetricsProvider(),
			msgCfg,
			mocktypes.NewMealPlanningDataManagerMock(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing pre-writes producer", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}
		msgCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", msgCfg.DataChangesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			context.Background(),
			logging.NewNoopLogger(),
			cfg,
			&textsearchcfg.Config{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			tracing.NewNoopTracerProvider(),
			metrics.NewNoopMetricsProvider(),
			msgCfg,
			mocktypes.NewMealPlanningDataManagerMock(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
