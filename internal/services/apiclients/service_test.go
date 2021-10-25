package apiclients

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockauthn "gitlab.com/prixfixe/prixfixe/internal/authentication/mock"
	"gitlab.com/prixfixe/prixfixe/internal/database"
	mockencoding "gitlab.com/prixfixe/prixfixe/internal/encoding/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/observability/metrics/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	mockrandom "gitlab.com/prixfixe/prixfixe/internal/random/mock"
	mockrouting "gitlab.com/prixfixe/prixfixe/internal/routing/mock"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	return &service{
		apiClientDataManager:      database.BuildMockDatabase(),
		logger:                    logging.NewNoopLogger(),
		encoderDecoder:            mockencoding.NewMockEncoderDecoder(),
		authenticator:             &mockauthn.Authenticator{},
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		urlClientIDExtractor:      func(req *http.Request) string { return "" },
		apiClientCounter:          &mockmetrics.UnitCounter{},
		secretGenerator:           &mockrandom.Generator{},
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
			"BuildRouteParamStringIDFetcher",
			APIClientIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		s := ProvideAPIClientsService(
			logging.NewNoopLogger(),
			mockAPIClientDataManager,
			&mocktypes.UserDataManager{},
			&mockauthn.Authenticator{},
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
