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
	ingredienttagmappingsservice "gitlab.com/prixfixe/prixfixe/services/v1/ingredienttagmappings"
	invitationsservice "gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	iterationmediasservice "gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	oauth2clientsservice "gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
	recipeiterationsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipeiterations"
	recipeiterationstepsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipeiterationsteps"
	recipesservice "gitlab.com/prixfixe/prixfixe/services/v1/recipes"
	recipestepingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipestepingredients"
	recipesteppreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipesteppreparations"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipesteps"
	recipetagsservice "gitlab.com/prixfixe/prixfixe/services/v1/recipetags"
	reportsservice "gitlab.com/prixfixe/prixfixe/services/v1/reports"
	requiredpreparationinstrumentsservice "gitlab.com/prixfixe/prixfixe/services/v1/requiredpreparationinstruments"
	usersservice "gitlab.com/prixfixe/prixfixe/services/v1/users"
	validingredientpreparationsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredientpreparations"
	validingredientsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredients"
	validingredienttagsservice "gitlab.com/prixfixe/prixfixe/services/v1/validingredienttags"
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
		validIngredientTagsService:            &mockmodels.ValidIngredientTagDataServer{},
		ingredientTagMappingsService:          &mockmodels.IngredientTagMappingDataServer{},
		validPreparationsService:              &mockmodels.ValidPreparationDataServer{},
		requiredPreparationInstrumentsService: &mockmodels.RequiredPreparationInstrumentDataServer{},
		validIngredientPreparationsService:    &mockmodels.ValidIngredientPreparationDataServer{},
		recipesService:                        &mockmodels.RecipeDataServer{},
		recipeTagsService:                     &mockmodels.RecipeTagDataServer{},
		recipeStepsService:                    &mockmodels.RecipeStepDataServer{},
		recipeStepPreparationsService:         &mockmodels.RecipeStepPreparationDataServer{},
		recipeStepIngredientsService:          &mockmodels.RecipeStepIngredientDataServer{},
		recipeIterationsService:               &mockmodels.RecipeIterationDataServer{},
		recipeIterationStepsService:           &mockmodels.RecipeIterationStepDataServer{},
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
			&validingredienttagsservice.Service{},
			&ingredienttagmappingsservice.Service{},
			&validpreparationsservice.Service{},
			&requiredpreparationinstrumentsservice.Service{},
			&validingredientpreparationsservice.Service{},
			&recipesservice.Service{},
			&recipetagsservice.Service{},
			&recipestepsservice.Service{},
			&recipesteppreparationsservice.Service{},
			&recipestepingredientsservice.Service{},
			&recipeiterationsservice.Service{},
			&recipeiterationstepsservice.Service{},
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
			&validingredienttagsservice.Service{},
			&ingredienttagmappingsservice.Service{},
			&validpreparationsservice.Service{},
			&requiredpreparationinstrumentsservice.Service{},
			&validingredientpreparationsservice.Service{},
			&recipesservice.Service{},
			&recipetagsservice.Service{},
			&recipestepsservice.Service{},
			&recipesteppreparationsservice.Service{},
			&recipestepingredientsservice.Service{},
			&recipeiterationsservice.Service{},
			&recipeiterationstepsservice.Service{},
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
			&validingredienttagsservice.Service{},
			&ingredienttagmappingsservice.Service{},
			&validpreparationsservice.Service{},
			&requiredpreparationinstrumentsservice.Service{},
			&validingredientpreparationsservice.Service{},
			&recipesservice.Service{},
			&recipetagsservice.Service{},
			&recipestepsservice.Service{},
			&recipesteppreparationsservice.Service{},
			&recipestepingredientsservice.Service{},
			&recipeiterationsservice.Service{},
			&recipeiterationstepsservice.Service{},
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
