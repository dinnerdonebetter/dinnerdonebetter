package authentication

import (
	"encoding/base64"
	"net/http"
	"testing"

	mockauthn "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication/mock"
	tokenscfg "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/authentication/tokens/config"
	identitymanagermock "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/identity/manager/mock"
	oauthmock "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/oauth/mock"

	noopanalytics "github.com/primandproper/platform/analytics/noop"
	"github.com/primandproper/platform/encoding"
	msgconfig "github.com/primandproper/platform/messagequeue/config"
	mockpublishers "github.com/primandproper/platform/messagequeue/mock"
	"github.com/primandproper/platform/observability/logging"
	"github.com/primandproper/platform/observability/tracing"
	"github.com/primandproper/platform/reflection"
	mockrouting "github.com/primandproper/platform/routing/mock"
	"github.com/primandproper/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	ctx := t.Context()
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
	pp.On(reflection.GetMethodName(pp.ProvidePublisher), queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

	rpm := mockrouting.NewRouteParamManager()
	rpm.On(
		"BuildRouteParamStringIDFetcher",
		AuthProviderParamKey,
	).Return(func(*http.Request) string { return "" })

	s, err := ProvideService(
		ctx,
		logger,
		cfg,
		&mockauthn.Authenticator{},
		&oauthmock.RepositoryMock{},
		&identitymanagermock.IdentityDataManager{},
		encoderDecoder,
		tracing.NewNoopTracerProvider(),
		pp,
		noopanalytics.NewEventReporter(),
		rpm,
		queueCfg,
	)
	require.NoError(t, err)

	mock.AssertExpectationsForObjects(t, pp, rpm)

	return s.(*service)
}

func TestProvideService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
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
		pp.On(reflection.GetMethodName(pp.ProvidePublisher), queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			AuthProviderParamKey,
		).Return(func(*http.Request) string { return "" })

		s, err := ProvideService(
			ctx,
			logger,
			cfg,
			&mockauthn.Authenticator{},
			&oauthmock.RepositoryMock{},
			&identitymanagermock.IdentityDataManager{},
			encoderDecoder,
			tracing.NewNoopTracerProvider(),
			pp,
			noopanalytics.NewEventReporter(),
			rpm,
			queueCfg,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)
	})
}
