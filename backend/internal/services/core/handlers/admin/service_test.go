package admin

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	logger := logging.NewNoopLogger()

	s, err := ProvideService(
		logger,
		&mocktypes.AdminUserDataManagerMock{},
		encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON),
		tracing.NewNoopTracerProvider(),
		&msgconfig.QueuesConfig{
			DataChangesTopicName:              "DataChangesTopicName",
			OutboundEmailsTopicName:           "OutboundEmailsTopicName",
			SearchIndexRequestsTopicName:      "SearchIndexRequestsTopicName",
			UserDataAggregationTopicName:      "UserDataAggregationTopicName",
			WebhookExecutionRequestsTopicName: "WebhookExecutionRequestsTopicName",
		},
		&mockpublishers.PublisherProvider{},
	)
	assert.NoError(t, err)

	srv, ok := s.(*service)
	require.True(t, ok)

	return srv
}

func TestProvideAdminService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		s, err := ProvideService(
			logger,
			&mocktypes.AdminUserDataManagerMock{},
			encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON),
			tracing.NewNoopTracerProvider(),
			&msgconfig.QueuesConfig{},
			&mockpublishers.PublisherProvider{},
		)
		assert.NoError(t, err)

		assert.NotNil(t, s)
	})
}
