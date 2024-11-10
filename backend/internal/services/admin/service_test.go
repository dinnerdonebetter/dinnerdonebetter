package admin

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	logger := logging.NewNoopLogger()

	s := ProvideService(
		logger,
		&mocktypes.AdminUserDataManagerMock{},
		encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON),
		tracing.NewNoopTracerProvider(),
		config.QueueSettings{},
		&mockpublishers.ProducerProvider{},
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
			&mocktypes.AdminUserDataManagerMock{},
			encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON),
			tracing.NewNoopTracerProvider(),
			config.QueueSettings{},
			&mockpublishers.ProducerProvider{},
		)

		assert.NotNil(t, s)
	})
}
