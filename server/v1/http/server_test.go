package httpserver

import (
	"context"
	"errors"
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	"gitlab.com/prixfixe/prixfixe/internal/v1/config"
	mockencoding "gitlab.com/prixfixe/prixfixe/internal/v1/encoding/mock"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"
	mockmodels "gitlab.com/prixfixe/prixfixe/models/v1/mock"
	authservice "gitlab.com/prixfixe/prixfixe/services/v1/auth"
	frontendservice "gitlab.com/prixfixe/prixfixe/services/v1/frontend"
	invitationsservice "gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	iterationmediasservice "gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	oauth2clientsservice "gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
	recipeiterationsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipeiterations"
	recipesservice "gitlab.com/prixfixe/prixfixe/services/v1/recipes"
	recipestepeventsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepevents"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepingredients"
	recipestepinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepinstruments"
	recipestepproductsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepproducts"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipesteps"
	reportsservice "gitlab.com/prixfixe/prixfixe/services/v1/reports"
	requiredpreparationinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/requiredpreparationinstruments"
	usersservice "gitlab.com/prixfixe/prixfixe/services/v1/users"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredients"
	validinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/validinstruments"
	validpreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/validpreparations"
	webhooksservice "gitlab.com/prixfixe/prixfixe/services/v1/webhooks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

func buildTestServer() *Server {
	s := &Server{
		DebugMode:  true,
		db:         database.BuildMockDatabase(),
		config:     &config.ServerConfig{},
		encoder:    &mockencoding.EncoderDecoder{},
		httpServer: provideHTTPServer(),
		logger:     noop.ProvideNoopLogger(),
		frontendService: frontendservice.ProvideFrontendService(
			noop.ProvideNoopLogger(),
			config.FrontendSettings{},
		),
		webhooksService:                       &mockmodels.WebhookDataServer{},
		usersService:                          &mockmodels.UserDataServer{},
		authService:                           &authservice.Service{},
		validInstrumentsService:               &mockmodels.ValidInstrumentDataServer{},
		validIngredientsService:               &mockmodels.ValidIngredientDataServer{},
		validPreparationsService:              &mockmodels.ValidPreparationDataServer{},
		validIngredientPreparationsService:    &mockmodels.ValidIngredientPreparationDataServer{},
		requiredPreparationInstrumentsService: &mockmodels.RequiredPreparationInstrumentDataServer{},
		recipesService:                        &mockmodels.RecipeDataServer{},
		recipeStepsService:                    &mockmodels.RecipeStepDataServer{},
		recipeStepInstrumentsService:          &mockmodels.RecipeStepInstrumentDataServer{},
		recipeStepIngredientsService:          &mockmodels.RecipeStepIngredientDataServer{},
		recipeStepProductsService:             &mockmodels.RecipeStepProductDataServer{},
		recipeIterationsService:               &mockmodels.RecipeIterationDataServer{},
		recipeStepEventsService:               &mockmodels.RecipeStepEventDataServer{},
		iterationMediasService:                &mockmodels.IterationMediaDataServer{},
		invitationsService:                    &mockmodels.InvitationDataServer{},
		reportsService:                        &mockmodels.ReportDataServer{},
		oauth2ClientsService:                  &mockmodels.OAuth2ClientDataServer{},
	}

	return s
}

func TestProvideServer(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fakemodels.BuildFakeWebhookList()

		mockDB := database.BuildMockDatabase()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock.Anything).Return(exampleWebhookList, nil)

		actual, err := ProvideServer(
			ctx,
			&config.ServerConfig{
				Auth: config.AuthSettings{
					CookieSecret: "THISISAVERYLONGSTRINGFORTESTPURPOSES",
				},
			},
			&authservice.Service{},
			&frontendservice.Service{},
			&validinstrumentsservice.Service{},
			&validingredientsservice.Service{},
			&validpreparationsservice.Service{},
			&validingredientpreparationsservice.Service{},
			&requiredpreparationinstrumentsservice.Service{},
			&recipesservice.Service{},
			&recipestepsservice.Service{},
			&recipestepinstrumentsservice.Service{},
			&recipestepingredientsservice.Service{},
			&recipestepproductsservice.Service{},
			&recipeiterationsservice.Service{},
			&recipestepeventsservice.Service{},
			&iterationmediasservice.Service{},
			&invitationsservice.Service{},
			&reportsservice.Service{},
			&usersservice.Service{},
			&oauth2clientsservice.Service{},
			&webhooksservice.Service{},
			mockDB,
			noop.ProvideNoopLogger(),
			&mockencoding.EncoderDecoder{},
			newsman.NewNewsman(nil, nil),
		)

		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with invalid cookie secret", func(t *testing.T) {
		ctx := context.Background()

		exampleWebhookList := fakemodels.BuildFakeWebhookList()

		mockDB := database.BuildMockDatabase()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock.Anything).Return(exampleWebhookList, nil)

		actual, err := ProvideServer(
			ctx,
			&config.ServerConfig{
				Auth: config.AuthSettings{
					CookieSecret: "THISSTRINGISNTLONGENOUGH:(",
				},
			},
			&authservice.Service{},
			&frontendservice.Service{},
			&validinstrumentsservice.Service{},
			&validingredientsservice.Service{},
			&validpreparationsservice.Service{},
			&validingredientpreparationsservice.Service{},
			&requiredpreparationinstrumentsservice.Service{},
			&recipesservice.Service{},
			&recipestepsservice.Service{},
			&recipestepinstrumentsservice.Service{},
			&recipestepingredientsservice.Service{},
			&recipestepproductsservice.Service{},
			&recipeiterationsservice.Service{},
			&recipestepeventsservice.Service{},
			&iterationmediasservice.Service{},
			&invitationsservice.Service{},
			&reportsservice.Service{},
			&usersservice.Service{},
			&oauth2clientsservice.Service{},
			&webhooksservice.Service{},
			mockDB,
			noop.ProvideNoopLogger(),
			&mockencoding.EncoderDecoder{},
			newsman.NewNewsman(nil, nil),
		)

		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})

	T.Run("with error fetching webhooks", func(t *testing.T) {
		ctx := context.Background()

		mockDB := database.BuildMockDatabase()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock.Anything).Return((*models.WebhookList)(nil), errors.New("blah"))

		actual, err := ProvideServer(
			ctx,
			&config.ServerConfig{
				Auth: config.AuthSettings{
					CookieSecret: "THISISAVERYLONGSTRINGFORTESTPURPOSES",
				},
			},
			&authservice.Service{},
			&frontendservice.Service{},
			&validinstrumentsservice.Service{},
			&validingredientsservice.Service{},
			&validpreparationsservice.Service{},
			&validingredientpreparationsservice.Service{},
			&requiredpreparationinstrumentsservice.Service{},
			&recipesservice.Service{},
			&recipestepsservice.Service{},
			&recipestepinstrumentsservice.Service{},
			&recipestepingredientsservice.Service{},
			&recipestepproductsservice.Service{},
			&recipeiterationsservice.Service{},
			&recipestepeventsservice.Service{},
			&iterationmediasservice.Service{},
			&invitationsservice.Service{},
			&reportsservice.Service{},
			&usersservice.Service{},
			&oauth2clientsservice.Service{},
			&webhooksservice.Service{},
			mockDB,
			noop.ProvideNoopLogger(),
			&mockencoding.EncoderDecoder{},
			newsman.NewNewsman(nil, nil),
		)

		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, mockDB)
	})
}
