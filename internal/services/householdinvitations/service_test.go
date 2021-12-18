package householdinvitations

import (
	"errors"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/customerdata"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	mockrouting "github.com/prixfixeco/api_server/internal/routing/mock"
	householdsservice "github.com/prixfixeco/api_server/internal/services/households"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
)

func buildTestService() *service {
	return &service{
		logger:                         logging.NewNoopLogger(),
		householdInvitationDataManager: &mocktypes.HouseholdInvitationDataManager{},
		householdInvitationIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:                 mockencoding.NewMockEncoderDecoder(),
		tracer:                         tracing.NewTracerForTest("test"),
	}
}

func TestProvideHouseholdInvitationsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			householdsservice.HouseholdIDURIParamKey,
		).Return(func(*http.Request) string { return "" })
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			HouseholdInvitationIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
			DataChangesTopicName: "data_changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		actual, err := ProvideHouseholdInvitationsService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.UserDataManager{},
			&mocktypes.HouseholdInvitationDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)

		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing data changes publisher", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			DataChangesTopicName: "data_changes",
		}
		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		actual, err := ProvideHouseholdInvitationsService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.UserDataManager{},
			&mocktypes.HouseholdInvitationDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)

		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
