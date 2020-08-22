// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"context"
	"database/sql"
	"gitlab.com/prixfixe/prixfixe/database/v1"
	"gitlab.com/prixfixe/prixfixe/internal/v1/auth"
	"gitlab.com/prixfixe/prixfixe/internal/v1/config"
	"gitlab.com/prixfixe/prixfixe/internal/v1/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/v1/search/bleve"
	"gitlab.com/prixfixe/prixfixe/server/v1"
	"gitlab.com/prixfixe/prixfixe/server/v1/http"
	auth2 "gitlab.com/prixfixe/prixfixe/services/v1/auth"
	"gitlab.com/prixfixe/prixfixe/services/v1/frontend"
	"gitlab.com/prixfixe/prixfixe/services/v1/invitations"
	"gitlab.com/prixfixe/prixfixe/services/v1/iterationmedias"
	"gitlab.com/prixfixe/prixfixe/services/v1/oauth2clients"
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
	"gitlab.com/prixfixe/prixfixe/services/v1/validingredientpreparations"
	"gitlab.com/prixfixe/prixfixe/services/v1/validingredients"
	"gitlab.com/prixfixe/prixfixe/services/v1/validinstruments"
	"gitlab.com/prixfixe/prixfixe/services/v1/validpreparations"
	"gitlab.com/prixfixe/prixfixe/services/v1/webhooks"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
	"gitlab.com/verygoodsoftwarenotvirus/newsman"
)

// Injectors from wire.go:

// BuildServer builds a server.
func BuildServer(ctx context.Context, cfg *config.ServerConfig, logger logging.Logger, database2 database.DataManager, db *sql.DB) (*server.Server, error) {
	authSettings := config.ProvideConfigAuthSettings(cfg)
	bcryptHashCost := auth.ProvideBcryptHashCost()
	authenticator := auth.ProvideBcryptAuthenticator(bcryptHashCost, logger)
	userDataManager := users.ProvideUserDataManager(database2)
	clientIDFetcher := httpserver.ProvideOAuth2ClientsServiceClientIDFetcher(logger)
	encoderDecoder := encoding.ProvideResponseEncoder()
	unitCounterProvider := metrics.ProvideUnitCounterProvider()
	service, err := oauth2clients.ProvideOAuth2ClientsService(logger, database2, authenticator, clientIDFetcher, encoderDecoder, unitCounterProvider)
	if err != nil {
		return nil, err
	}
	oAuth2ClientValidator := auth2.ProvideOAuth2ClientValidator(service)
	databaseSettings := config.ProvideConfigDatabaseSettings(cfg)
	sessionManager := config.ProvideSessionManager(authSettings, databaseSettings, db)
	authService, err := auth2.ProvideAuthService(logger, authSettings, authenticator, userDataManager, oAuth2ClientValidator, sessionManager, encoderDecoder)
	if err != nil {
		return nil, err
	}
	frontendSettings := config.ProvideConfigFrontendSettings(cfg)
	frontendService := frontend.ProvideFrontendService(logger, frontendSettings)
	validInstrumentDataManager := validinstruments.ProvideValidInstrumentDataManager(database2)
	validInstrumentIDFetcher := httpserver.ProvideValidInstrumentsServiceValidInstrumentIDFetcher(logger)
	websocketAuthFunc := auth2.ProvideWebsocketAuthFunc(authService)
	typeNameManipulationFunc := httpserver.ProvideNewsmanTypeNameManipulationFunc()
	newsmanNewsman := newsman.NewNewsman(websocketAuthFunc, typeNameManipulationFunc)
	reporter := ProvideReporter(newsmanNewsman)
	searchSettings := config.ProvideSearchSettings(cfg)
	indexManagerProvider := bleve.ProvideBleveIndexManagerProvider()
	searchIndex, err := validinstruments.ProvideValidInstrumentsServiceSearchIndex(searchSettings, indexManagerProvider, logger)
	if err != nil {
		return nil, err
	}
	validinstrumentsService, err := validinstruments.ProvideValidInstrumentsService(logger, validInstrumentDataManager, validInstrumentIDFetcher, encoderDecoder, unitCounterProvider, reporter, searchIndex)
	if err != nil {
		return nil, err
	}
	validInstrumentDataServer := validinstruments.ProvideValidInstrumentDataServer(validinstrumentsService)
	validIngredientDataManager := validingredients.ProvideValidIngredientDataManager(database2)
	validIngredientIDFetcher := httpserver.ProvideValidIngredientsServiceValidIngredientIDFetcher(logger)
	validingredientsSearchIndex, err := validingredients.ProvideValidIngredientsServiceSearchIndex(searchSettings, indexManagerProvider, logger)
	if err != nil {
		return nil, err
	}
	validingredientsService, err := validingredients.ProvideValidIngredientsService(logger, validIngredientDataManager, validIngredientIDFetcher, encoderDecoder, unitCounterProvider, reporter, validingredientsSearchIndex)
	if err != nil {
		return nil, err
	}
	validIngredientDataServer := validingredients.ProvideValidIngredientDataServer(validingredientsService)
	validPreparationDataManager := validpreparations.ProvideValidPreparationDataManager(database2)
	validPreparationIDFetcher := httpserver.ProvideValidPreparationsServiceValidPreparationIDFetcher(logger)
	validpreparationsSearchIndex, err := validpreparations.ProvideValidPreparationsServiceSearchIndex(searchSettings, indexManagerProvider, logger)
	if err != nil {
		return nil, err
	}
	validpreparationsService, err := validpreparations.ProvideValidPreparationsService(logger, validPreparationDataManager, validPreparationIDFetcher, encoderDecoder, unitCounterProvider, reporter, validpreparationsSearchIndex)
	if err != nil {
		return nil, err
	}
	validPreparationDataServer := validpreparations.ProvideValidPreparationDataServer(validpreparationsService)
	validIngredientPreparationDataManager := validingredientpreparations.ProvideValidIngredientPreparationDataManager(database2)
	validIngredientPreparationIDFetcher := httpserver.ProvideValidIngredientPreparationsServiceValidIngredientPreparationIDFetcher(logger)
	validingredientpreparationsService, err := validingredientpreparations.ProvideValidIngredientPreparationsService(logger, validIngredientPreparationDataManager, validIngredientPreparationIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	validIngredientPreparationDataServer := validingredientpreparations.ProvideValidIngredientPreparationDataServer(validingredientpreparationsService)
	requiredPreparationInstrumentDataManager := requiredpreparationinstruments.ProvideRequiredPreparationInstrumentDataManager(database2)
	requiredPreparationInstrumentIDFetcher := httpserver.ProvideRequiredPreparationInstrumentsServiceRequiredPreparationInstrumentIDFetcher(logger)
	requiredpreparationinstrumentsService, err := requiredpreparationinstruments.ProvideRequiredPreparationInstrumentsService(logger, requiredPreparationInstrumentDataManager, requiredPreparationInstrumentIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	requiredPreparationInstrumentDataServer := requiredpreparationinstruments.ProvideRequiredPreparationInstrumentDataServer(requiredpreparationinstrumentsService)
	recipeDataManager := recipes.ProvideRecipeDataManager(database2)
	recipeIDFetcher := httpserver.ProvideRecipesServiceRecipeIDFetcher(logger)
	userIDFetcher := httpserver.ProvideRecipesServiceUserIDFetcher()
	recipesService, err := recipes.ProvideRecipesService(logger, recipeDataManager, recipeIDFetcher, userIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	recipeDataServer := recipes.ProvideRecipeDataServer(recipesService)
	recipeStepDataManager := recipesteps.ProvideRecipeStepDataManager(database2)
	recipestepsRecipeIDFetcher := httpserver.ProvideRecipeStepsServiceRecipeIDFetcher(logger)
	recipeStepIDFetcher := httpserver.ProvideRecipeStepsServiceRecipeStepIDFetcher(logger)
	recipestepsUserIDFetcher := httpserver.ProvideRecipeStepsServiceUserIDFetcher()
	recipestepsService, err := recipesteps.ProvideRecipeStepsService(logger, recipeDataManager, recipeStepDataManager, recipestepsRecipeIDFetcher, recipeStepIDFetcher, recipestepsUserIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	recipeStepDataServer := recipesteps.ProvideRecipeStepDataServer(recipestepsService)
	recipeStepInstrumentDataManager := recipestepinstruments.ProvideRecipeStepInstrumentDataManager(database2)
	recipestepinstrumentsRecipeIDFetcher := httpserver.ProvideRecipeStepInstrumentsServiceRecipeIDFetcher(logger)
	recipestepinstrumentsRecipeStepIDFetcher := httpserver.ProvideRecipeStepInstrumentsServiceRecipeStepIDFetcher(logger)
	recipeStepInstrumentIDFetcher := httpserver.ProvideRecipeStepInstrumentsServiceRecipeStepInstrumentIDFetcher(logger)
	recipestepinstrumentsUserIDFetcher := httpserver.ProvideRecipeStepInstrumentsServiceUserIDFetcher()
	recipestepinstrumentsService, err := recipestepinstruments.ProvideRecipeStepInstrumentsService(logger, recipeDataManager, recipeStepDataManager, recipeStepInstrumentDataManager, recipestepinstrumentsRecipeIDFetcher, recipestepinstrumentsRecipeStepIDFetcher, recipeStepInstrumentIDFetcher, recipestepinstrumentsUserIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	recipeStepInstrumentDataServer := recipestepinstruments.ProvideRecipeStepInstrumentDataServer(recipestepinstrumentsService)
	recipeStepIngredientDataManager := recipestepingredients.ProvideRecipeStepIngredientDataManager(database2)
	recipestepingredientsRecipeIDFetcher := httpserver.ProvideRecipeStepIngredientsServiceRecipeIDFetcher(logger)
	recipestepingredientsRecipeStepIDFetcher := httpserver.ProvideRecipeStepIngredientsServiceRecipeStepIDFetcher(logger)
	recipeStepIngredientIDFetcher := httpserver.ProvideRecipeStepIngredientsServiceRecipeStepIngredientIDFetcher(logger)
	recipestepingredientsUserIDFetcher := httpserver.ProvideRecipeStepIngredientsServiceUserIDFetcher()
	recipestepingredientsService, err := recipestepingredients.ProvideRecipeStepIngredientsService(logger, recipeDataManager, recipeStepDataManager, recipeStepIngredientDataManager, recipestepingredientsRecipeIDFetcher, recipestepingredientsRecipeStepIDFetcher, recipeStepIngredientIDFetcher, recipestepingredientsUserIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	recipeStepIngredientDataServer := recipestepingredients.ProvideRecipeStepIngredientDataServer(recipestepingredientsService)
	recipeStepProductDataManager := recipestepproducts.ProvideRecipeStepProductDataManager(database2)
	recipestepproductsRecipeIDFetcher := httpserver.ProvideRecipeStepProductsServiceRecipeIDFetcher(logger)
	recipestepproductsRecipeStepIDFetcher := httpserver.ProvideRecipeStepProductsServiceRecipeStepIDFetcher(logger)
	recipeStepProductIDFetcher := httpserver.ProvideRecipeStepProductsServiceRecipeStepProductIDFetcher(logger)
	recipestepproductsUserIDFetcher := httpserver.ProvideRecipeStepProductsServiceUserIDFetcher()
	recipestepproductsService, err := recipestepproducts.ProvideRecipeStepProductsService(logger, recipeDataManager, recipeStepDataManager, recipeStepProductDataManager, recipestepproductsRecipeIDFetcher, recipestepproductsRecipeStepIDFetcher, recipeStepProductIDFetcher, recipestepproductsUserIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	recipeStepProductDataServer := recipestepproducts.ProvideRecipeStepProductDataServer(recipestepproductsService)
	recipeIterationDataManager := recipeiterations.ProvideRecipeIterationDataManager(database2)
	recipeiterationsRecipeIDFetcher := httpserver.ProvideRecipeIterationsServiceRecipeIDFetcher(logger)
	recipeIterationIDFetcher := httpserver.ProvideRecipeIterationsServiceRecipeIterationIDFetcher(logger)
	recipeiterationsUserIDFetcher := httpserver.ProvideRecipeIterationsServiceUserIDFetcher()
	recipeiterationsService, err := recipeiterations.ProvideRecipeIterationsService(logger, recipeDataManager, recipeIterationDataManager, recipeiterationsRecipeIDFetcher, recipeIterationIDFetcher, recipeiterationsUserIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	recipeIterationDataServer := recipeiterations.ProvideRecipeIterationDataServer(recipeiterationsService)
	recipeStepEventDataManager := recipestepevents.ProvideRecipeStepEventDataManager(database2)
	recipestepeventsRecipeIDFetcher := httpserver.ProvideRecipeStepEventsServiceRecipeIDFetcher(logger)
	recipestepeventsRecipeStepIDFetcher := httpserver.ProvideRecipeStepEventsServiceRecipeStepIDFetcher(logger)
	recipeStepEventIDFetcher := httpserver.ProvideRecipeStepEventsServiceRecipeStepEventIDFetcher(logger)
	recipestepeventsUserIDFetcher := httpserver.ProvideRecipeStepEventsServiceUserIDFetcher()
	recipestepeventsService, err := recipestepevents.ProvideRecipeStepEventsService(logger, recipeDataManager, recipeStepDataManager, recipeStepEventDataManager, recipestepeventsRecipeIDFetcher, recipestepeventsRecipeStepIDFetcher, recipeStepEventIDFetcher, recipestepeventsUserIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	recipeStepEventDataServer := recipestepevents.ProvideRecipeStepEventDataServer(recipestepeventsService)
	iterationMediaDataManager := iterationmedias.ProvideIterationMediaDataManager(database2)
	iterationmediasRecipeIDFetcher := httpserver.ProvideIterationMediasServiceRecipeIDFetcher(logger)
	iterationmediasRecipeIterationIDFetcher := httpserver.ProvideIterationMediasServiceRecipeIterationIDFetcher(logger)
	iterationMediaIDFetcher := httpserver.ProvideIterationMediasServiceIterationMediaIDFetcher(logger)
	iterationmediasUserIDFetcher := httpserver.ProvideIterationMediasServiceUserIDFetcher()
	iterationmediasService, err := iterationmedias.ProvideIterationMediasService(logger, recipeDataManager, recipeIterationDataManager, iterationMediaDataManager, iterationmediasRecipeIDFetcher, iterationmediasRecipeIterationIDFetcher, iterationMediaIDFetcher, iterationmediasUserIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	iterationMediaDataServer := iterationmedias.ProvideIterationMediaDataServer(iterationmediasService)
	invitationDataManager := invitations.ProvideInvitationDataManager(database2)
	invitationIDFetcher := httpserver.ProvideInvitationsServiceInvitationIDFetcher(logger)
	invitationsUserIDFetcher := httpserver.ProvideInvitationsServiceUserIDFetcher()
	invitationsService, err := invitations.ProvideInvitationsService(logger, invitationDataManager, invitationIDFetcher, invitationsUserIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	invitationDataServer := invitations.ProvideInvitationDataServer(invitationsService)
	reportDataManager := reports.ProvideReportDataManager(database2)
	reportIDFetcher := httpserver.ProvideReportsServiceReportIDFetcher(logger)
	reportsUserIDFetcher := httpserver.ProvideReportsServiceUserIDFetcher()
	reportsService, err := reports.ProvideReportsService(logger, reportDataManager, reportIDFetcher, reportsUserIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	reportDataServer := reports.ProvideReportDataServer(reportsService)
	usersUserIDFetcher := httpserver.ProvideUsersServiceUserIDFetcher(logger)
	usersService, err := users.ProvideUsersService(authSettings, logger, userDataManager, authenticator, usersUserIDFetcher, encoderDecoder, unitCounterProvider, reporter)
	if err != nil {
		return nil, err
	}
	userDataServer := users.ProvideUserDataServer(usersService)
	oAuth2ClientDataServer := oauth2clients.ProvideOAuth2ClientDataServer(service)
	webhookDataManager := webhooks.ProvideWebhookDataManager(database2)
	webhooksUserIDFetcher := httpserver.ProvideWebhooksServiceUserIDFetcher()
	webhookIDFetcher := httpserver.ProvideWebhooksServiceWebhookIDFetcher(logger)
	webhooksService, err := webhooks.ProvideWebhooksService(logger, webhookDataManager, webhooksUserIDFetcher, webhookIDFetcher, encoderDecoder, unitCounterProvider, newsmanNewsman)
	if err != nil {
		return nil, err
	}
	webhookDataServer := webhooks.ProvideWebhookDataServer(webhooksService)
	httpserverServer, err := httpserver.ProvideServer(ctx, cfg, authService, frontendService, validInstrumentDataServer, validIngredientDataServer, validPreparationDataServer, validIngredientPreparationDataServer, requiredPreparationInstrumentDataServer, recipeDataServer, recipeStepDataServer, recipeStepInstrumentDataServer, recipeStepIngredientDataServer, recipeStepProductDataServer, recipeIterationDataServer, recipeStepEventDataServer, iterationMediaDataServer, invitationDataServer, reportDataServer, userDataServer, oAuth2ClientDataServer, webhookDataServer, database2, logger, encoderDecoder, newsmanNewsman)
	if err != nil {
		return nil, err
	}
	serverServer, err := server.ProvideServer(cfg, httpserverServer)
	if err != nil {
		return nil, err
	}
	return serverServer, nil
}

// wire.go:

// ProvideReporter is an obligatory function that hopefully wire will eliminate for me one day.
func ProvideReporter(n *newsman.Newsman) newsman.Reporter {
	return n
}
