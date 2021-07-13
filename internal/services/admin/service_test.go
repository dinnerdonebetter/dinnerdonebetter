package admin

import (
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/authentication"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	mockrouting "gitlab.com/prixfixe/prixfixe/internal/routing/mock"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"

	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	logger := logging.NewNoopLogger()

	rpm := mockrouting.NewRouteParamManager()
	rpm.On(
		"BuildRouteParamIDFetcher",
		mock.IsType(logging.NewNoopLogger()), UserIDURIParamKey, "user").Return(func(*http.Request) uint64 { return 0 })

	s := ProvideService(
		logger,
		&authservice.Config{Cookies: authservice.CookieConfig{SigningKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!"}},
		&authentication.MockAuthenticator{},
		&mocktypes.AdminUserDataManager{},
		&mocktypes.AuditLogEntryDataManager{},
		scs.New(),
		encoding.ProvideServerEncoderDecoder(logger, encoding.ContentTypeJSON),
		rpm,
	)

	mock.AssertExpectationsForObjects(t, rpm)

	srv, ok := s.(*service)
	require.True(t, ok)

	return srv
}

func TestProvideAdminService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamIDFetcher",
			mock.IsType(logging.NewNoopLogger()), UserIDURIParamKey, "user").Return(func(*http.Request) uint64 { return 0 })

		s := ProvideService(
			logger,
			&authservice.Config{Cookies: authservice.CookieConfig{SigningKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!"}},
			&authentication.MockAuthenticator{},
			&mocktypes.AdminUserDataManager{},
			&mocktypes.AuditLogEntryDataManager{},
			scs.New(),
			encoding.ProvideServerEncoderDecoder(logger, encoding.ContentTypeJSON),
			rpm,
		)

		assert.NotNil(t, s)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}
