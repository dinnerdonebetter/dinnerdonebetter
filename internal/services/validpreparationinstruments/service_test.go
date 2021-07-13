package validpreparationinstruments

import (
	"net/http"
	"testing"

	mockencoding "gitlab.com/prixfixe/prixfixe/internal/encoding/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/observability/metrics/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	mockrouting "gitlab.com/prixfixe/prixfixe/internal/routing/mock"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:                                logging.NewNoopLogger(),
		validPreparationInstrumentCounter:     &mockmetrics.UnitCounter{},
		validPreparationInstrumentDataManager: &mocktypes.ValidPreparationInstrumentDataManager{},
		validPreparationInstrumentIDFetcher:   func(req *http.Request) uint64 { return 0 },
		encoderDecoder:                        mockencoding.NewMockEncoderDecoder(),
		tracer:                                tracing.NewTracer("test"),
	}
}

func TestProvideValidPreparationInstrumentsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		var ucp metrics.UnitCounterProvider = func(counterName, description string) metrics.UnitCounter {
			return &mockmetrics.UnitCounter{}
		}

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamIDFetcher",
			mock.IsType(logging.NewNoopLogger()),
			ValidPreparationInstrumentIDURIParamKey,
			"valid_preparation_instrument",
		).Return(func(*http.Request) uint64 { return 0 })

		s, err := ProvideService(
			logging.NewNoopLogger(),
			Config{},
			&mocktypes.ValidPreparationInstrumentDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			ucp,
			rpm,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}
