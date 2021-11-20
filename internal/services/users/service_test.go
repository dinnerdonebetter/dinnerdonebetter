package users

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockauthn "github.com/prixfixeco/api_server/internal/authentication/mock"
	"github.com/prixfixeco/api_server/internal/database"
	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/metrics"
	mockmetrics "github.com/prixfixeco/api_server/internal/observability/metrics/mock"
	"github.com/prixfixeco/api_server/internal/routing/chi"
	mockrouting "github.com/prixfixeco/api_server/internal/routing/mock"
	authservice "github.com/prixfixeco/api_server/internal/services/authentication"
	"github.com/prixfixeco/api_server/internal/uploads/images"
	mockuploads "github.com/prixfixeco/api_server/internal/uploads/mock"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func buildTestService(t *testing.T) *service {
	t.Helper()

	expectedUserCount := uint64(123)

	uc := &mockmetrics.UnitCounter{}
	mockDB := database.NewMockDatabase()
	mockDB.UserDataManager.On(
		"GetAllUsersCount",
		testutils.ContextMatcher,
	).Return(expectedUserCount, nil)

	s := ProvideUsersService(
		&authservice.Config{},
		logging.NewNoopLogger(),
		&mocktypes.UserDataManager{},
		&mocktypes.HouseholdDataManager{},
		&mockauthn.Authenticator{},
		mockencoding.NewMockEncoderDecoder(),
		func(counterName, description string) metrics.UnitCounter {
			return uc
		},
		&images.MockImageUploadProcessor{},
		&mockuploads.UploadManager{},
		chi.NewRouteParamManager(),
	)

	mock.AssertExpectationsForObjects(t, mockDB, uc)

	return s.(*service)
}

func TestProvideUsersService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			UserIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		s := ProvideUsersService(
			&authservice.Config{},
			logging.NewNoopLogger(),
			&mocktypes.UserDataManager{},
			&mocktypes.HouseholdDataManager{},
			&mockauthn.Authenticator{},
			mockencoding.NewMockEncoderDecoder(),
			func(counterName, description string) metrics.UnitCounter {
				return &mockmetrics.UnitCounter{}
			},
			&images.MockImageUploadProcessor{},
			&mockuploads.UploadManager{},
			rpm,
		)

		assert.NotNil(t, s)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}
