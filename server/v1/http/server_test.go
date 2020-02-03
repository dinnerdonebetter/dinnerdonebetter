package httpserver

import (
	"context"
	"testing"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	"gitlab.com/prixfixe/prixfixe/internal/v1/config"
	mockencoding "gitlab.com/prixfixe/prixfixe/internal/v1/encoding/mock"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	mockmodels "gitlab.com/prixfixe/prixfixe/models/v1/mock"
	"gitlab.com/prixfixe/prixfixe/services/v1/auth"
	"gitlab.com/prixfixe/prixfixe/services/v1/frontend"
	"gitlab.com/prixfixe/prixfixe/services/v1/ingredients"
	"gitlab.com/prixfixe/prixfixe/services/v1/instruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	"gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	"gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
	"gitlab.com/prixfixe/prixfixe/services/v1/preparations"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipeiterations"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipes"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepevents"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepingredients"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepinstruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipestepproducts"
	"gitlab.com/prixfixe/prixfixe/services/v1/recipesteps"
	"gitlab.com/prixfixe/prixfixe/services/v1/reports"
	"gitlab.com/prixfixe/prixfixe/services/v1/requiredpreparationinstruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/users"
	"gitlab.com/prixfixe/prixfixe/services/v1/webhooks"

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
		frontendService: frontend.ProvideFrontendService(
			noop.ProvideNoopLogger(),
			config.FrontendSettings{},
		),
		webhooksService:                       &mockmodels.WebhookDataServer{},
		usersService:                          &mockmodels.UserDataServer{},
		authService:                           &auth.Service{},
		instrumentsService:                    &mockmodels.InstrumentDataServer{},
		ingredientsService:                    &mockmodels.IngredientDataServer{},
		preparationsService:                   &mockmodels.PreparationDataServer{},
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
		mockDB := database.BuildMockDatabase()
		mockDB.WebhookDataManager.On("GetAllWebhooks", mock.Anything).Return(&models.WebhookList{}, nil)

		actual, err := ProvideServer(
			context.Background(),
			&config.ServerConfig{
				Auth: config.AuthSettings{
					CookieSecret: "THISISAVERYLONGSTRINGFORTESTPURPOSES",
				},
			},
			&auth.Service{},
			&frontend.Service{},
			&instruments.Service{},
			&ingredients.Service{},
			&preparations.Service{},
			&requiredpreparationinstruments.Service{},
			&recipes.Service{},
			&recipesteps.Service{},
			&recipestepinstruments.Service{},
			&recipestepingredients.Service{},
			&recipestepproducts.Service{},
			&recipeiterations.Service{},
			&recipestepevents.Service{},
			&iterationmedias.Service{},
			&invitations.Service{},
			&reports.Service{},
			&users.Service{},
			&oauth2clients.Service{},
			&webhooks.Service{},
			mockDB,
			noop.ProvideNoopLogger(),
			&mockencoding.EncoderDecoder{},
			newsman.NewNewsman(nil, nil),
		)

		assert.NotNil(t, actual)
		assert.NoError(t, err)
	})
}
