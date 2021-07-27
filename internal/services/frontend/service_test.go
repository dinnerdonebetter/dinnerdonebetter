package frontend

import (
	"net/http"
	"testing"

	capitalism "gitlab.com/prixfixe/prixfixe/internal/capitalism"
	database "gitlab.com/prixfixe/prixfixe/internal/database"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	mockrouting "gitlab.com/prixfixe/prixfixe/internal/routing/mock"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func dummyIDFetcher(*http.Request) uint64 {
	return 0
}

func TestProvideService(t *testing.T) {
	t.Parallel()

	cfg := &Config{}
	logger := logging.NewNoopLogger()
	authService := &mocktypes.AuthService{}
	usersService := &mocktypes.UsersService{}
	dataManager := database.BuildMockDatabase()

	rpm := mockrouting.NewRouteParamManager()
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), apiClientIDURLParamKey, "API client").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), accountIDURLParamKey, "account").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), webhookIDURLParamKey, "webhook").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), validInstrumentIDURLParamKey, "valid instrument").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), validPreparationIDURLParamKey, "valid preparation").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), validIngredientIDURLParamKey, "valid ingredient").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), validIngredientPreparationIDURLParamKey, "valid ingredient preparation").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), validPreparationInstrumentIDURLParamKey, "valid preparation instrument").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), recipeIDURLParamKey, "recipe").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), recipeStepIDURLParamKey, "recipe step").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), recipeStepIngredientIDURLParamKey, "recipe step ingredient").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), recipeStepProductIDURLParamKey, "recipe step product").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), invitationIDURLParamKey, "invitation").Return(dummyIDFetcher)
	rpm.On("BuildRouteParamIDFetcher", mock.IsType(logger), reportIDURLParamKey, "report").Return(dummyIDFetcher)

	s := ProvideService(
		cfg,
		logger,
		authService,
		usersService,
		dataManager,
		rpm,
		capitalism.NewMockPaymentManager(),
		&mocktypes.ValidIngredientDataManager{},
		&mocktypes.ValidInstrumentDataManager{},
		&mocktypes.ValidPreparationDataManager{},
	)

	mock.AssertExpectationsForObjects(t, authService, usersService, dataManager, rpm)
	assert.NotNil(t, s)
}
