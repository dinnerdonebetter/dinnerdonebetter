package reports

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
		logger:            logging.NewNoopLogger(),
		reportCounter:     &mockmetrics.UnitCounter{},
		reportDataManager: &mocktypes.ReportDataManager{},
		reportIDFetcher:   func(req *http.Request) uint64 { return 0 },
		encoderDecoder:    mockencoding.NewMockEncoderDecoder(),
		tracer:            tracing.NewTracer("test"),
	}
}

func TestProvideReportsService(T *testing.T) {
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
			ReportIDURIParamKey,
			"report",
		).Return(func(*http.Request) uint64 { return 0 })

		s, err := ProvideService(
			logging.NewNoopLogger(),
			Config{},
			&mocktypes.ReportDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			ucp,
			rpm,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}
