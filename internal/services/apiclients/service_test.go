package apiclients

import (
	"net/http"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/authentication"
	"gitlab.com/prixfixe/prixfixe/internal/database"
	mockencoding "gitlab.com/prixfixe/prixfixe/internal/encoding/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/observability/metrics/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/random"
	mockrouting "gitlab.com/prixfixe/prixfixe/internal/routing/mock"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	return &service{
		apiClientDataManager:      database.BuildMockDatabase(),
		logger:                    logging.NewNoopLogger(),
		encoderDecoder:            mockencoding.NewMockEncoderDecoder(),
		authenticator:             &authentication.MockAuthenticator{},
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		urlClientIDExtractor:      func(req *http.Request) uint64 { return 0 },
		apiClientCounter:          &mockmetrics.UnitCounter{},
		secretGenerator:           &random.MockGenerator{},
		tracer:                    tracing.NewTracer(serviceName),
		cfg:                       &config{},
	}
}

func TestProvideAPIClientsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		mockAPIClientDataManager := &mocktypes.APIClientDataManager{}

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamIDFetcher",
			mock.IsType(logging.NewNoopLogger()),
			APIClientIDURIParamKey,
			"api client",
		).Return(func(*http.Request) uint64 { return 0 })

		s := ProvideAPIClientsService(
			logging.NewNoopLogger(),
			mockAPIClientDataManager,
			&mocktypes.UserDataManager{},
			&authentication.MockAuthenticator{},
			mockencoding.NewMockEncoderDecoder(),
			func(counterName, description string) metrics.UnitCounter {
				return nil
			},
			rpm,
			&config{},
		)
		assert.NotNil(t, s)

		mock.AssertExpectationsForObjects(t, mockAPIClientDataManager, rpm)
	})
}
