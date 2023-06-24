// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package build

import (
	"context"
	config8 "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	config4 "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	config6 "github.com/dinnerdonebetter/backend/internal/email/config"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	config7 "github.com/dinnerdonebetter/backend/internal/featureflags/config"
	"github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	config5 "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	config2 "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	config3 "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	"github.com/dinnerdonebetter/backend/internal/services/admin"
	"github.com/dinnerdonebetter/backend/internal/services/apiclients"
	authentication2 "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/internal/services/householdinstrumentownerships"
	"github.com/dinnerdonebetter/backend/internal/services/householdinvitations"
	"github.com/dinnerdonebetter/backend/internal/services/households"
	"github.com/dinnerdonebetter/backend/internal/services/mealplanevents"
	"github.com/dinnerdonebetter/backend/internal/services/mealplangrocerylistitems"
	"github.com/dinnerdonebetter/backend/internal/services/mealplanoptions"
	"github.com/dinnerdonebetter/backend/internal/services/mealplanoptionvotes"
	"github.com/dinnerdonebetter/backend/internal/services/mealplans"
	"github.com/dinnerdonebetter/backend/internal/services/mealplantasks"
	"github.com/dinnerdonebetter/backend/internal/services/meals"
	"github.com/dinnerdonebetter/backend/internal/services/oauth2clients"
	"github.com/dinnerdonebetter/backend/internal/services/recipepreptasks"
	"github.com/dinnerdonebetter/backend/internal/services/reciperatings"
	"github.com/dinnerdonebetter/backend/internal/services/recipes"
	recipestepingredients2 "github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	"github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	"github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	"github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	"github.com/dinnerdonebetter/backend/internal/services/recipestepvessels"
	"github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	"github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	"github.com/dinnerdonebetter/backend/internal/services/useringredientpreferences"
	"github.com/dinnerdonebetter/backend/internal/services/users"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientgroups"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	"github.com/dinnerdonebetter/backend/internal/services/validingredients"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientstateingredients"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	"github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/validmeasurementconversions"
	"github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	"github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	"github.com/dinnerdonebetter/backend/internal/services/vendorproxy"
	"github.com/dinnerdonebetter/backend/internal/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/uploads/images"
)

// Injectors from build.go:

// Build builds a server.
func Build(ctx context.Context, cfg *config.InstanceConfig) (*http.Server, error) {
	httpConfig := cfg.Server
	observabilityConfig := &cfg.Observability
	configConfig := &observabilityConfig.Logging
	logger, err := config2.ProvideLogger(configConfig)
	if err != nil {
		return nil, err
	}
	config9 := &cfg.Database
	config10 := &observabilityConfig.Tracing
	tracerProvider, err := config3.ProvideTracerProvider(ctx, config10, logger)
	if err != nil {
		return nil, err
	}
	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, config9, tracerProvider)
	if err != nil {
		return nil, err
	}
	encodingConfig := cfg.Encoding
	contentType := encoding.ProvideContentType(encodingConfig)
	serverEncoderDecoder := encoding.ProvideServerEncoderDecoder(logger, tracerProvider, contentType)
	routingConfig := &cfg.Routing
	router := chi.NewRouter(logger, tracerProvider, routingConfig)
	servicesConfig := &cfg.Services
	authenticationConfig := &servicesConfig.Auth
	authenticator := authentication.ProvideArgon2Authenticator(logger, tracerProvider)
	householdUserMembershipDataManager := database.ProvideHouseholdUserMembershipDataManager(dataManager)
	cookieConfig := authenticationConfig.Cookies
	sessionManager, err := config4.ProvideSessionManager(cookieConfig, dataManager)
	if err != nil {
		return nil, err
	}
	config11 := &cfg.Events
	publisherProvider, err := config5.ProvidePublisherProvider(logger, tracerProvider, config11)
	if err != nil {
		return nil, err
	}
	generator := random.NewGenerator(logger, tracerProvider)
	config12 := &cfg.Email
	client := tracing.BuildTracedHTTPClient()
	emailer, err := config6.ProvideEmailer(config12, logger, tracerProvider, client)
	if err != nil {
		return nil, err
	}
	config13 := &cfg.FeatureFlags
	featureFlagManager, err := config7.ProvideFeatureFlagManager(config13, logger, tracerProvider, client)
	if err != nil {
		return nil, err
	}
	authService, err := authentication2.ProvideService(ctx, logger, authenticationConfig, authenticator, dataManager, householdUserMembershipDataManager, sessionManager, serverEncoderDecoder, tracerProvider, publisherProvider, generator, emailer, featureFlagManager)
	if err != nil {
		return nil, err
	}
	usersConfig := &servicesConfig.Users
	userDataManager := database.ProvideUserDataManager(dataManager)
	householdDataManager := database.ProvideHouseholdDataManager(dataManager)
	householdInvitationDataManager := database.ProvideHouseholdInvitationDataManager(dataManager)
	mediaUploadProcessor := images.NewImageUploadProcessor(logger, tracerProvider)
	routeParamManager := chi.NewRouteParamManager()
	passwordResetTokenDataManager := database.ProvidePasswordResetTokenDataManager(dataManager)
	userDataService, err := users.ProvideUsersService(ctx, usersConfig, authenticationConfig, logger, userDataManager, householdDataManager, householdInvitationDataManager, householdUserMembershipDataManager, authenticator, serverEncoderDecoder, mediaUploadProcessor, routeParamManager, tracerProvider, publisherProvider, generator, passwordResetTokenDataManager, featureFlagManager)
	if err != nil {
		return nil, err
	}
	householdsConfig := servicesConfig.Households
	householdDataService, err := households.ProvideService(logger, householdsConfig, householdDataManager, householdInvitationDataManager, householdUserMembershipDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	householdinvitationsConfig := &servicesConfig.HouseholdInvitations
	householdInvitationDataService, err := householdinvitations.ProvideHouseholdInvitationsService(logger, householdinvitationsConfig, userDataManager, householdInvitationDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider, emailer, generator)
	if err != nil {
		return nil, err
	}
	apiClientDataManager := database.ProvideAPIClientDataManager(dataManager)
	apiclientsConfig := apiclients.ProvideConfig(authenticationConfig)
	apiClientDataService, err := apiclients.ProvideAPIClientsService(logger, apiClientDataManager, userDataManager, authenticator, serverEncoderDecoder, routeParamManager, apiclientsConfig, tracerProvider, generator, publisherProvider)
	if err != nil {
		return nil, err
	}
	validinstrumentsConfig := &servicesConfig.ValidInstruments
	config14 := &cfg.Search
	validInstrumentDataManager := database.ProvideValidInstrumentDataManager(dataManager)
	validInstrumentDataService, err := validinstruments.ProvideService(ctx, logger, validinstrumentsConfig, config14, validInstrumentDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientsConfig := &servicesConfig.ValidIngredients
	validIngredientDataManager := database.ProvideValidIngredientDataManager(dataManager)
	validIngredientDataService, err := validingredients.ProvideService(ctx, logger, validingredientsConfig, config14, validIngredientDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientgroupsConfig := &servicesConfig.ValidIngredientGroups
	validIngredientGroupDataManager := database.ProvideValidIngredientGroupDataManager(dataManager)
	validIngredientGroupDataService, err := validingredientgroups.ProvideService(logger, validingredientgroupsConfig, validIngredientGroupDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validpreparationsConfig := &servicesConfig.ValidPreparations
	validPreparationDataManager := database.ProvideValidPreparationDataManager(dataManager)
	validPreparationDataService, err := validpreparations.ProvideService(ctx, logger, validpreparationsConfig, config14, validPreparationDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientpreparationsConfig := &servicesConfig.ValidIngredientPreparations
	validIngredientPreparationDataManager := database.ProvideValidIngredientPreparationDataManager(dataManager)
	validIngredientPreparationDataService, err := validingredientpreparations.ProvideService(logger, validingredientpreparationsConfig, validIngredientPreparationDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealsConfig := &servicesConfig.Meals
	mealDataManager := database.ProvideMealDataManager(dataManager)
	mealDataService, err := meals.ProvideService(ctx, logger, mealsConfig, config14, mealDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipesConfig := &servicesConfig.Recipes
	recipeDataManager := database.ProvideRecipeDataManager(dataManager)
	recipeMediaDataManager := database.ProvideRecipeMediaDataManager(dataManager)
	recipeAnalyzer := recipeanalysis.NewRecipeAnalyzer(logger, tracerProvider)
	recipeDataService, err := recipes.ProvideService(ctx, logger, recipesConfig, config14, recipeDataManager, recipeMediaDataManager, recipeAnalyzer, serverEncoderDecoder, routeParamManager, publisherProvider, mediaUploadProcessor, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepsConfig := &servicesConfig.RecipeSteps
	recipeStepDataManager := database.ProvideRecipeStepDataManager(dataManager)
	recipeStepDataService, err := recipesteps.ProvideService(ctx, logger, recipestepsConfig, recipeStepDataManager, recipeMediaDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider, mediaUploadProcessor)
	if err != nil {
		return nil, err
	}
	recipestepproductsConfig := &servicesConfig.RecipeStepProducts
	recipeStepProductDataManager := database.ProvideRecipeStepProductDataManager(dataManager)
	recipeStepProductDataService, err := recipestepproducts.ProvideService(logger, recipestepproductsConfig, recipeStepProductDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepinstrumentsConfig := &servicesConfig.RecipeStepInstruments
	recipeStepInstrumentDataManager := database.ProvideRecipeStepInstrumentDataManager(dataManager)
	recipeStepInstrumentDataService, err := recipestepinstruments.ProvideService(logger, recipestepinstrumentsConfig, recipeStepInstrumentDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepingredientsConfig := &servicesConfig.RecipeStepIngredients
	recipeStepIngredientDataManager := database.ProvideRecipeStepIngredientDataManager(dataManager)
	recipeStepIngredientDataService, err := recipestepingredients.ProvideService(logger, recipestepingredientsConfig, recipeStepIngredientDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplansConfig := &servicesConfig.MealPlans
	mealPlanDataManager := database.ProvideMealPlanDataManager(dataManager)
	mealPlanDataService, err := mealplans.ProvideService(logger, mealplansConfig, mealPlanDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplanoptionsConfig := &servicesConfig.MealPlanOptions
	mealPlanOptionDataManager := database.ProvideMealPlanOptionDataManager(dataManager)
	mealPlanOptionDataService, err := mealplanoptions.ProvideService(logger, mealplanoptionsConfig, mealPlanOptionDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplanoptionvotesConfig := &servicesConfig.MealPlanOptionVotes
	mealPlanOptionVoteDataService, err := mealplanoptionvotes.ProvideService(logger, mealplanoptionvotesConfig, dataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validmeasurementunitsConfig := &servicesConfig.ValidMeasurementUnits
	validMeasurementUnitDataManager := database.ProvideValidMeasurementUnitDataManager(dataManager)
	validMeasurementUnitDataService, err := validmeasurementunits.ProvideService(ctx, logger, validmeasurementunitsConfig, config14, validMeasurementUnitDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientstatesConfig := &servicesConfig.ValidIngredientStates
	validIngredientStateDataManager := database.ProvideValidIngredientStateDataManager(dataManager)
	validIngredientStateDataService, err := validingredientstates.ProvideService(ctx, logger, validingredientstatesConfig, config14, validIngredientStateDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validpreparationinstrumentsConfig := &servicesConfig.ValidPreparationInstruments
	validPreparationInstrumentDataManager := database.ProvideValidPreparationInstrumentDataManager(dataManager)
	validPreparationInstrumentDataService, err := validpreparationinstruments.ProvideService(logger, validpreparationinstrumentsConfig, validPreparationInstrumentDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientmeasurementunitsConfig := &servicesConfig.ValidInstrumentMeasurementUnits
	validIngredientMeasurementUnitDataManager := database.ProvideValidIngredientMeasurementUnitDataManager(dataManager)
	validIngredientMeasurementUnitDataService, err := validingredientmeasurementunits.ProvideService(logger, validingredientmeasurementunitsConfig, validIngredientMeasurementUnitDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplaneventsConfig := &servicesConfig.MealPlanEvents
	mealPlanEventDataManager := database.ProvideMealPlanEventDataManager(dataManager)
	mealPlanEventDataService, err := mealplanevents.ProvideService(logger, mealplaneventsConfig, mealPlanEventDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplantasksConfig := &servicesConfig.MealPlanTasks
	mealPlanTaskDataManager := database.ProvideMealPlanTaskDataManager(dataManager)
	mealPlanTaskDataService, err := mealplantasks.ProvideService(logger, mealplantasksConfig, mealPlanTaskDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipepreptasksConfig := &servicesConfig.RecipePrepTasks
	recipePrepTaskDataManager := database.ProvideRecipePrepTaskDataManager(dataManager)
	recipePrepTaskDataService, err := recipepreptasks.ProvideService(logger, recipepreptasksConfig, recipePrepTaskDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	mealplangrocerylistitemsConfig := &servicesConfig.MealPlanGroceryListItems
	mealPlanGroceryListItemDataManager := database.ProvideMealPlanGroceryListItemDataManager(dataManager)
	mealPlanGroceryListItemDataService, err := mealplangrocerylistitems.ProvideService(logger, mealplangrocerylistitemsConfig, mealPlanGroceryListItemDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validmeasurementconversionsConfig := &servicesConfig.ValidMeasurementConversions
	validMeasurementConversionDataManager := database.ProvideValidMeasurementConversionDataManager(dataManager)
	validMeasurementConversionDataService, err := validmeasurementconversions.ProvideService(ctx, logger, validmeasurementconversionsConfig, validMeasurementConversionDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	config15 := &servicesConfig.RecipeStepCompletionConditions
	recipeStepCompletionConditionDataManager := database.ProvideRecipeStepCompletionConditionDataManager(dataManager)
	recipeStepCompletionConditionDataService, err := recipestepingredients2.ProvideService(logger, config15, recipeStepCompletionConditionDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientstateingredientsConfig := &servicesConfig.ValidIngredientStateIngredients
	validIngredientStateIngredientDataManager := database.ProvideValidIngredientStateIngredientDataManager(dataManager)
	validIngredientStateIngredientDataService, err := validingredientstateingredients.ProvideService(logger, validingredientstateingredientsConfig, validIngredientStateIngredientDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepvesselsConfig := &servicesConfig.RecipeStepVessels
	recipeStepVesselDataManager := database.ProvideRecipeStepVesselDataManager(dataManager)
	recipeStepVesselDataService, err := recipestepvessels.ProvideService(logger, recipestepvesselsConfig, recipeStepVesselDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	webhooksConfig := &servicesConfig.Webhooks
	webhookDataManager := database.ProvideWebhookDataManager(dataManager)
	webhookDataService, err := webhooks.ProvideWebhooksService(logger, webhooksConfig, webhookDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	adminUserDataManager := database.ProvideAdminUserDataManager(dataManager)
	adminService := admin.ProvideService(logger, authenticationConfig, authenticator, adminUserDataManager, sessionManager, serverEncoderDecoder, routeParamManager, tracerProvider)
	vendorproxyConfig := &servicesConfig.VendorProxy
	config16 := &cfg.Analytics
	eventReporter, err := config8.ProvideEventReporter(config16, logger, tracerProvider)
	if err != nil {
		return nil, err
	}
	service, err := vendorproxy.ProvideService(logger, vendorproxyConfig, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider, featureFlagManager, eventReporter)
	if err != nil {
		return nil, err
	}
	servicesettingsConfig := &servicesConfig.ServiceSettings
	serviceSettingDataManager := database.ProvideServiceSettingDataManager(dataManager)
	serviceSettingDataService, err := servicesettings.ProvideService(logger, servicesettingsConfig, serviceSettingDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	servicesettingconfigurationsConfig := &servicesConfig.ServiceSettingConfigurations
	serviceSettingConfigurationDataManager := database.ProvideServiceSettingConfigurationDataManager(dataManager)
	serviceSettingConfigurationDataService, err := servicesettingconfigurations.ProvideService(logger, servicesettingconfigurationsConfig, serviceSettingConfigurationDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	useringredientpreferencesConfig := &servicesConfig.UserIngredientPreferences
	userIngredientPreferenceDataManager := database.ProvideUserIngredientPreferenceDataManager(dataManager)
	userIngredientPreferenceDataService, err := useringredientpreferences.ProvideService(ctx, logger, useringredientpreferencesConfig, userIngredientPreferenceDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	reciperatingsConfig := &servicesConfig.RecipeRatings
	recipeRatingDataManager := database.ProvideRecipeRatingDataManager(dataManager)
	recipeRatingDataService, err := reciperatings.ProvideService(logger, reciperatingsConfig, recipeRatingDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	householdinstrumentownershipsConfig := &servicesConfig.HouseholdInstrumentOwnerships
	householdInstrumentOwnershipDataManager := database.ProvideHouseholdInstrumentOwnershipDataManager(dataManager)
	householdInstrumentOwnershipDataService, err := householdinstrumentownerships.ProvideService(logger, householdinstrumentownershipsConfig, householdInstrumentOwnershipDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	oAuth2ClientDataManager := database.ProvideOAuth2ClientDataManager(dataManager)
	oauth2clientsConfig := oauth2clients.ProvideConfig(authenticationConfig)
	oAuth2ClientDataService, err := oauth2clients.ProvideOAuth2ClientsService(logger, oAuth2ClientDataManager, userDataManager, authenticator, serverEncoderDecoder, routeParamManager, oauth2clientsConfig, tracerProvider, generator, publisherProvider)
	if err != nil {
		return nil, err
	}
	server, err := http.ProvideHTTPServer(ctx, httpConfig, dataManager, logger, serverEncoderDecoder, router, tracerProvider, authService, userDataService, householdDataService, householdInvitationDataService, apiClientDataService, validInstrumentDataService, validIngredientDataService, validIngredientGroupDataService, validPreparationDataService, validIngredientPreparationDataService, mealDataService, recipeDataService, recipeStepDataService, recipeStepProductDataService, recipeStepInstrumentDataService, recipeStepIngredientDataService, mealPlanDataService, mealPlanOptionDataService, mealPlanOptionVoteDataService, validMeasurementUnitDataService, validIngredientStateDataService, validPreparationInstrumentDataService, validIngredientMeasurementUnitDataService, mealPlanEventDataService, mealPlanTaskDataService, recipePrepTaskDataService, mealPlanGroceryListItemDataService, validMeasurementConversionDataService, recipeStepCompletionConditionDataService, validIngredientStateIngredientDataService, recipeStepVesselDataService, webhookDataService, adminService, service, serviceSettingDataService, serviceSettingConfigurationDataService, userIngredientPreferenceDataService, recipeRatingDataService, householdInstrumentOwnershipDataService, oAuth2ClientDataService)
	if err != nil {
		return nil, err
	}
	return server, nil
}
