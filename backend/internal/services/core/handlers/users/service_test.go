package users

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/analytics"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/featureflags"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/random"
	"github.com/dinnerdonebetter/backend/internal/lib/routing/chi"
	mockrouting "github.com/dinnerdonebetter/backend/internal/lib/routing/mock"
	authservice "github.com/dinnerdonebetter/backend/internal/services/core/handlers/authentication"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	msgCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

	pp := &mockpublishers.ProducerProvider{}
	pp.On("ProvidePublisher", msgCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	s, err := ProvideUsersService(
		&authservice.Config{},
		logging.NewNoopLogger(),
		&mocktypes.UserDataManagerMock{},
		&mocktypes.HouseholdInvitationDataManagerMock{},
		&mocktypes.HouseholdUserMembershipDataManagerMock{},
		&mockauthn.Authenticator{},
		mockencoding.NewMockEncoderDecoder(),
		chi.NewRouteParamManager(),
		tracing.NewNoopTracerProvider(),
		pp,
		random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
		&mocktypes.PasswordResetTokenDataManagerMock{},
		&featureflags.NoopFeatureFlagManager{},
		analytics.NewNoopEventReporter(),
		msgCfg,
	)

	require.NoError(t, err)

	return s.(*service)
}

func TestProvideUsersService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		msgCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			UserIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", msgCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideUsersService(
			&authservice.Config{},
			logging.NewNoopLogger(),
			&mocktypes.UserDataManagerMock{},
			&mocktypes.HouseholdInvitationDataManagerMock{},
			&mocktypes.HouseholdUserMembershipDataManagerMock{},
			&mockauthn.Authenticator{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			tracing.NewNoopTracerProvider(),
			pp,
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
			&mocktypes.PasswordResetTokenDataManagerMock{},
			&featureflags.NoopFeatureFlagManager{},
			analytics.NewNoopEventReporter(),
			msgCfg,
		)

		assert.NotNil(t, s)
		require.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}
