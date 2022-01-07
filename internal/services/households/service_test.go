package households

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/customerdata"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
	mockmetrics "github.com/prixfixeco/api_server/internal/observability/metrics/mock"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	mockrouting "github.com/prixfixeco/api_server/internal/routing/mock"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
)

func buildTestService() *service {
	return &service{
		logger:                         logging.NewNoopLogger(),
		householdCounter:               &mockmetrics.UnitCounter{},
		householdDataManager:           &mocktypes.HouseholdDataManager{},
		householdMembershipDataManager: &mocktypes.HouseholdUserMembershipDataManager{},
		householdIDFetcher:             func(req *http.Request) string { return "" },
		encoderDecoder:                 mockencoding.NewMockEncoderDecoder(),
		tracer:                         tracing.NewTracerForTest("test"),
	}
}

func TestProvideHouseholdsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		var ucp metrics.UnitCounterProvider = func(counterName, description string) metrics.UnitCounter {
			return &mockmetrics.UnitCounter{}
		}

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			HouseholdIDURIParamKey,
		).Return(func(*http.Request) string { return "" })
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			UserIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := Config{
			DataChangesTopicName: "data changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.HouseholdDataManager{},
			&mocktypes.HouseholdInvitationDataManager{},
			&mocktypes.HouseholdUserMembershipDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			ucp,
			rpm,
			pp,
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing publisher", func(t *testing.T) {
		t.Parallel()

		var ucp metrics.UnitCounterProvider = func(counterName, description string) metrics.UnitCounter {
			return &mockmetrics.UnitCounter{}
		}

		rpm := mockrouting.NewRouteParamManager()
		cfg := Config{
			DataChangesTopicName: "data changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, errors.New("blah"))

		s, err := ProvideService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.HouseholdDataManager{},
			&mocktypes.HouseholdInvitationDataManager{},
			&mocktypes.HouseholdUserMembershipDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			ucp,
			rpm,
			pp,
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})
}
