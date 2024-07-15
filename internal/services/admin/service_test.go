package admin

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/alexedwards/scs/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	logger := logging.NewNoopLogger()

	s := ProvideService(
		logger,
		&authservice.Config{Cookies: authservice.CookieConfig{BlockKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!"}},
		&mocktypes.AdminUserDataManagerMock{},
		scs.New(),
		encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON),
		tracing.NewNoopTracerProvider(),
	)

	srv, ok := s.(*service)
	require.True(t, ok)

	return srv
}

func TestProvideAdminService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		s := ProvideService(
			logger,
			&authservice.Config{Cookies: authservice.CookieConfig{BlockKey: "BLAHBLAHBLAHPRETENDTHISISSECRET!"}},
			&mocktypes.AdminUserDataManagerMock{},
			scs.New(),
			encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON),
			tracing.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
	})
}
