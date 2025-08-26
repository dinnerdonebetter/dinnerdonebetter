package http

import (
	"encoding/base64"
	"net/http"
	"testing"

	mockauthn "github.com/dinnerdonebetter/backend/internal/authentication/mock"
	tokenscfg "github.com/dinnerdonebetter/backend/internal/authentication/tokens/config"
	identitymock "github.com/dinnerdonebetter/backend/internal/domain/identity/mock"
	oauthmock "github.com/dinnerdonebetter/backend/internal/domain/oauth/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/analytics"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/featureflags"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	mockrouting "github.com/dinnerdonebetter/backend/internal/platform/routing/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	logger := logging.NewNoopLogger()
	encoderDecoder := encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

	cfg := &Config{
		Tokens: tokenscfg.Config{
			Provider:                tokenscfg.ProviderJWT,
			Audience:                "",
			Base64EncodedSigningKey: base64.URLEncoding.EncodeToString([]byte(testutils.Example32ByteKey)),
		},
	}
	queueCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

	pp := &mockpublishers.PublisherProvider{}
	pp.On("ProvidePublisher", queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	rpm := mockrouting.NewRouteParamManager()
	rpm.On(
		"BuildRouteParamStringIDFetcher",
		AuthProviderParamKey,
	).Return(func(*http.Request) string { return "" })

	mmp := &metrics.MockProvider{}
	mmp.On("NewInt64Counter", rejectedRequestCounterName, []metric.Int64CounterOption(nil)).Return(
		metrics.Int64CounterForTest(t.Name()), nil,
	)

	s, err := ProvideService(
		logger,
		cfg,
		&mockauthn.Authenticator{},
		&oauthmock.RepositoryMock{},
		&identitymock.RepositoryMock{},
		encoderDecoder,
		tracing.NewNoopTracerProvider(),
		pp,
		&featureflags.NoopFeatureFlagManager{},
		analytics.NewNoopEventReporter(),
		rpm,
		mmp,
		queueCfg,
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, pp, rpm, mmp)

	return s.(*service)
}

func TestProvideService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		logger := logging.NewNoopLogger()
		encoderDecoder := encoding.ProvideServerEncoderDecoder(logger, tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

		cfg := &Config{
			Tokens: tokenscfg.Config{
				Provider:                tokenscfg.ProviderJWT,
				Audience:                "",
				Base64EncodedSigningKey: base64.URLEncoding.EncodeToString([]byte(testutils.Example32ByteKey)),
			},
		}
		queueCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

		pp := &mockpublishers.PublisherProvider{}
		pp.On("ProvidePublisher", queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			AuthProviderParamKey,
		).Return(func(*http.Request) string { return "" })

		s, err := ProvideService(
			logger,
			cfg,
			&mockauthn.Authenticator{},
			&oauthmock.RepositoryMock{},
			&identitymock.RepositoryMock{},
			encoderDecoder,
			tracing.NewNoopTracerProvider(),
			pp,
			&featureflags.NoopFeatureFlagManager{},
			analytics.NewNoopEventReporter(),
			rpm,
			metrics.NewNoopMetricsProvider(),
			queueCfg,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)
	})
}
