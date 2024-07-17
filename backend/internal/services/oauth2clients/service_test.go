package oauth2clients

import (
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/encoding/mock"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"
	randommock "github.com/dinnerdonebetter/backend/internal/pkg/random/mock"
	mockrouting "github.com/dinnerdonebetter/backend/internal/routing/mock"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	return &service{
		oauth2ClientDataManager:   database.NewMockDatabase(),
		logger:                    logging.NewNoopLogger(),
		encoderDecoder:            encoding.ProvideServerEncoderDecoder(nil, nil, encoding.ContentTypeJSON),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		urlClientIDExtractor:      func(req *http.Request) string { return "" },
		secretGenerator:           &randommock.Generator{},
		tracer:                    tracing.NewTracerForTest(serviceName),
	}
}

func TestProvideOAuth2ClientsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		mockOAuth2ClientDataManager := &mocktypes.OAuth2ClientDataManagerMock{}

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			OAuth2ClientIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
			DataChangesTopicName: t.Name(),
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideOAuth2ClientsService(
			logging.NewNoopLogger(),
			mockOAuth2ClientDataManager,
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			cfg,
			tracing.NewNoopTracerProvider(),
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
			pp,
		)
		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockOAuth2ClientDataManager, rpm)
	})
}
