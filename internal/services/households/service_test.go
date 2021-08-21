package households

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
		logger:                         logging.NewNoopLogger(),
		householdCounter:               &mockmetrics.UnitCounter{},
		householdDataManager:           &mocktypes.HouseholdDataManager{},
		householdMembershipDataManager: &mocktypes.HouseholdUserMembershipDataManager{},
		householdIDFetcher:             func(req *http.Request) uint64 { return 0 },
		encoderDecoder:                 mockencoding.NewMockEncoderDecoder(),
		tracer:                         tracing.NewTracer("test"),
	}
}

func TestProvideHouseholdsService(t *testing.T) {
	t.Parallel()

	var ucp metrics.UnitCounterProvider = func(counterName, description string) metrics.UnitCounter {
		return &mockmetrics.UnitCounter{}
	}

	l := logging.NewNoopLogger()

	rpm := mockrouting.NewRouteParamManager()
	rpm.On(
		"BuildRouteParamIDFetcher",
		mock.IsType(l), HouseholdIDURIParamKey, "household").Return(func(*http.Request) uint64 { return 0 })
	rpm.On(
		"BuildRouteParamIDFetcher",
		mock.IsType(l), UserIDURIParamKey, "user").Return(func(*http.Request) uint64 { return 0 })

	s := ProvideService(
		logging.NewNoopLogger(),
		&mocktypes.HouseholdDataManager{},
		&mocktypes.HouseholdUserMembershipDataManager{},
		mockencoding.NewMockEncoderDecoder(),
		ucp,
		rpm,
	)

	assert.NotNil(t, s)

	mock.AssertExpectationsForObjects(t, rpm)
}
