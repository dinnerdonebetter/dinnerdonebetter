package oauth2clients

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	mocktypes "github.com/dinnerdonebetter/backend/internal/domain/oauth/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	mockencoding "github.com/dinnerdonebetter/backend/internal/platform/encoding/mock"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	randommock "github.com/dinnerdonebetter/backend/internal/platform/random/mock"
	mockrouting "github.com/dinnerdonebetter/backend/internal/platform/routing/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	return &service{
		oauth2ClientDataManager:   mocktypes.NewRepositoryMock(),
		logger:                    logging.NewNoopLogger(),
		encoderDecoder:            encoding.ProvideServerEncoderDecoder(nil, nil, encoding.ContentTypeJSON),
		sessionContextDataFetcher: sessions.FetchContextDataFromRequest,
		urlClientIDExtractor:      func(req *http.Request) string { return "" },
		secretGenerator:           &randommock.Generator{},
		tracer:                    tracing.NewTracerForTest(serviceName),
	}
}

func TestProvideOAuth2ClientsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		mockOAuth2ClientDataManager := mocktypes.NewRepositoryMock()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			OAuth2ClientIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{OAuth2ClientCreationDisabled: true}
		msgCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

		pp := &mockpublishers.PublisherProvider{}
		pp.On("ProvidePublisher", msgCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideOAuth2ClientsService(
			logging.NewNoopLogger(),
			cfg,
			mockOAuth2ClientDataManager,
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			tracing.NewNoopTracerProvider(),
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
			pp,
			msgCfg,
		)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockOAuth2ClientDataManager, rpm)
	})
}
