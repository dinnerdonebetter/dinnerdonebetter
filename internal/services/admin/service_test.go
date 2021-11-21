package admin

import (
	"net/http"
	"testing"

	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockauthn "github.com/prixfixeco/api_server/internal/authentication/mock"
	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	mockrouting "github.com/prixfixeco/api_server/internal/routing/mock"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	logger := logging.NewNoopLogger()

	rpm := mockrouting.NewRouteParamManager()
	rpm.On(
		"BuildRouteParamStringIDFetcher",
		UserIDURIParamKey,
	).Return(func(*http.Request) string { return "" })

	s := ProvideService(
		logger,
		&authservice.Config{Cookies: authservice.CookieConfig{SigningKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!"}},
		&mockauthn.Authenticator{},
		&mocktypes.AdminUserDataManager{},
		scs.New(),
		encoding.ProvideServerEncoderDecoder(logger, encoding.ContentTypeJSON),
		rpm,
		&customerdata.MockCollector{},
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
			"BuildRouteParamStringIDFetcher",
			UserIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		s := ProvideService(
			logger,
			&authservice.Config{Cookies: authservice.CookieConfig{SigningKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!"}},
			&mockauthn.Authenticator{},
			&mocktypes.AdminUserDataManager{},
			scs.New(),
			encoding.ProvideServerEncoderDecoder(logger, encoding.ContentTypeJSON),
			rpm,
			&customerdata.MockCollector{},
		)

		assert.NotNil(t, s)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}
