package householdinvitations

import (
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/email"
	mockencoding "github.com/dinnerdonebetter/backend/internal/encoding/mock"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/random"
	mockrouting "github.com/dinnerdonebetter/backend/internal/routing/mock"
	householdsservice "github.com/dinnerdonebetter/backend/internal/services/households"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:                       logging.NewNoopLogger(),
		householdInvitationIDFetcher: func(req *http.Request) string { return "" },
		encoderDecoder:               mockencoding.NewMockEncoderDecoder(),
		tracer:                       tracing.NewTracerForTest("test"),
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
		pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		actual, err := ProvideHouseholdInvitationsService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.UserDataManager{},
			&mocktypes.HouseholdInvitationDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			tracing.NewNoopTracerProvider(),
			&email.MockEmailer{},
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
		)

		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing data changes publisher", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			DataChangesTopicName: "pre-writes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		actual, err := ProvideHouseholdInvitationsService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.UserDataManager{},
			&mocktypes.HouseholdInvitationDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			tracing.NewNoopTracerProvider(),
			&email.MockEmailer{},
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
		)

		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
