package users

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockauthn "github.com/prixfixeco/backend/internal/authentication/mock"
	"github.com/prixfixeco/backend/internal/email"
	mockencoding "github.com/prixfixeco/backend/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/backend/internal/messagequeue/mock"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/internal/random"
	"github.com/prixfixeco/backend/internal/routing/chi"
	mockrouting "github.com/prixfixeco/backend/internal/routing/mock"
	authservice "github.com/prixfixeco/backend/internal/services/authentication"
	"github.com/prixfixeco/backend/internal/storage"
	"github.com/prixfixeco/backend/internal/uploads"
	"github.com/prixfixeco/backend/internal/uploads/images"
	mocktypes "github.com/prixfixeco/backend/pkg/types/mock"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	cfg := &Config{
		Uploads: uploads.Config{
			Storage: storage.Config{
				FilesystemConfig: &storage.FilesystemConfig{RootDirectory: t.Name()},
				BucketName:       t.Name(),
				Provider:         storage.FilesystemProvider,
			},
			Debug: false,
		},
	}

	pp := &mockpublishers.ProducerProvider{}
	pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	s, err := ProvideUsersService(
		context.Background(),
		cfg,
		&authservice.Config{},
		logging.NewNoopLogger(),
		&mocktypes.UserDataManager{},
		&mocktypes.HouseholdDataManager{},
		&mocktypes.HouseholdInvitationDataManager{},
		&mockauthn.Authenticator{},
		mockencoding.NewMockEncoderDecoder(),
		&images.MockImageUploadProcessor{},
		chi.NewRouteParamManager(),
		tracing.NewNoopTracerProvider(),
		pp,
		random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
		&mocktypes.PasswordResetTokenDataManager{},
		&email.MockEmailer{},
	)

	require.NoError(t, err)

	return s.(*service)
}

func TestProvideUsersService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			UserIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
			Uploads: uploads.Config{
				Storage: storage.Config{
					FilesystemConfig: &storage.FilesystemConfig{RootDirectory: t.Name()},
					BucketName:       t.Name(),
					Provider:         storage.FilesystemProvider,
				},
				Debug: false,
			},
		}

		rpm.On(
			"BuildRouteParamStringIDFetcher",
			cfg.Uploads.Storage.UploadFilenameKey,
		).Return(func(*http.Request) string { return "" })

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideUsersService(
			context.Background(),
			cfg,
			&authservice.Config{},
			logging.NewNoopLogger(),
			&mocktypes.UserDataManager{},
			&mocktypes.HouseholdDataManager{},
			&mocktypes.HouseholdInvitationDataManager{},
			&mockauthn.Authenticator{},
			mockencoding.NewMockEncoderDecoder(),
			&images.MockImageUploadProcessor{},
			rpm,
			tracing.NewNoopTracerProvider(),
			pp,
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
			&mocktypes.PasswordResetTokenDataManager{},
			&email.MockEmailer{},
		)

		assert.NotNil(t, s)
		require.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}
