package validmeasurementunits

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockencoding "github.com/prixfixeco/backend/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/backend/internal/messagequeue/mock"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	mockrouting "github.com/prixfixeco/backend/internal/routing/mock"
	mocktypes "github.com/prixfixeco/backend/pkg/types/mock"
)

func buildTestService() *service {
	return &service{
		logger:                          logging.NewNoopLogger(),
		validMeasurementUnitDataManager: &mocktypes.ValidMeasurementUnitDataManager{},
		validMeasurementUnitIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:                  mockencoding.NewMockEncoderDecoder(),
		tracer:                          tracing.NewTracerForTest("test"),
	}
}

func TestProvideValidMeasurementUnitsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			ValidMeasurementUnitIDURIParamKey,
		).Return(func(*http.Request) string { return "" })
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			ValidIngredientIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
			DataChangesTopicName: "data_changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			logger,
			cfg,
			&mocktypes.ValidMeasurementUnitDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			tracing.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing data changes producer", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		cfg := &Config{
			DataChangesTopicName: "data_changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			logger,
			cfg,
			&mocktypes.ValidMeasurementUnitDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			tracing.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
