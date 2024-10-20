// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package build

import (
	"context"

	config6 "github.com/dinnerdonebetter/backend/internal/analytics/config"
	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	config5 "github.com/dinnerdonebetter/backend/internal/featureflags/config"
	"github.com/dinnerdonebetter/backend/internal/features/recipeanalysis"
	config4 "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	config2 "github.com/dinnerdonebetter/backend/internal/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	config3 "github.com/dinnerdonebetter/backend/internal/observability/tracing/config"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"
	"github.com/dinnerdonebetter/backend/internal/routing/chi"
	"github.com/dinnerdonebetter/backend/internal/server/http"
	"github.com/dinnerdonebetter/backend/internal/services/admin"
	"github.com/dinnerdonebetter/backend/internal/services/auditlogentries"
	authentication2 "github.com/dinnerdonebetter/backend/internal/services/authentication"
	workers2 "github.com/dinnerdonebetter/backend/internal/services/dataprivacy"
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
	"github.com/dinnerdonebetter/backend/internal/services/recipestepcompletionconditions"
	"github.com/dinnerdonebetter/backend/internal/services/recipestepingredients"
	"github.com/dinnerdonebetter/backend/internal/services/recipestepinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/recipestepproducts"
	"github.com/dinnerdonebetter/backend/internal/services/recipesteps"
	"github.com/dinnerdonebetter/backend/internal/services/recipestepvessels"
	"github.com/dinnerdonebetter/backend/internal/services/servicesettingconfigurations"
	"github.com/dinnerdonebetter/backend/internal/services/servicesettings"
	"github.com/dinnerdonebetter/backend/internal/services/useringredientpreferences"
	"github.com/dinnerdonebetter/backend/internal/services/usernotifications"
	"github.com/dinnerdonebetter/backend/internal/services/users"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientgroups"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientmeasurementunits"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientpreparations"
	"github.com/dinnerdonebetter/backend/internal/services/validingredients"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientstateingredients"
	"github.com/dinnerdonebetter/backend/internal/services/validingredientstates"
	"github.com/dinnerdonebetter/backend/internal/services/validinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/validmeasurementunitconversions"
	"github.com/dinnerdonebetter/backend/internal/services/validmeasurementunits"
	"github.com/dinnerdonebetter/backend/internal/services/validpreparationinstruments"
	"github.com/dinnerdonebetter/backend/internal/services/validpreparations"
	"github.com/dinnerdonebetter/backend/internal/services/validpreparationvessels"
	"github.com/dinnerdonebetter/backend/internal/services/validvessels"
	"github.com/dinnerdonebetter/backend/internal/services/webhooks"
	"github.com/dinnerdonebetter/backend/internal/services/workers"
	"github.com/dinnerdonebetter/backend/internal/uploads/images"
)

// Injectors from build.go:

// Build builds a server.
func Build(ctx context.Context, cfg *config.InstanceConfig) (http.Server, error) {
	httpConfig := cfg.Server
	observabilityConfig := &cfg.Observability
	configConfig := &observabilityConfig.Logging
	logger := config2.ProvideLogger(configConfig)
	config7 := &observabilityConfig.Tracing
	tracerProvider, err := config3.ProvideTracerProvider(ctx, config7, logger)
	if err != nil {
		return nil, err
	}
	config8 := &cfg.Database
	dataManager, err := postgres.ProvideDatabaseClient(ctx, logger, tracerProvider, config8)
	if err != nil {
		return nil, err
	}
	routingConfig := &cfg.Routing
	router := chi.NewRouter(logger, tracerProvider, routingConfig)
	servicesConfig := &cfg.Services
	authenticationConfig := &servicesConfig.Auth
	authenticator := authentication.ProvideArgon2Authenticator(logger, tracerProvider)
	householdUserMembershipDataManager := database.ProvideHouseholdUserMembershipDataManager(dataManager)
	encodingConfig := cfg.Encoding
	contentType := encoding.ProvideContentType(encodingConfig)
	serverEncoderDecoder := encoding.ProvideServerEncoderDecoder(logger, tracerProvider, contentType)
	config9 := &cfg.Events
	publisherProvider, err := config4.ProvidePublisherProvider(ctx, logger, tracerProvider, config9)
	if err != nil {
		return nil, err
	}
	config10 := &cfg.FeatureFlags
	client := tracing.BuildTracedHTTPClient()
	featureFlagManager, err := config5.ProvideFeatureFlagManager(config10, logger, tracerProvider, client)
	if err != nil {
		return nil, err
	}
	config11 := &cfg.Analytics
	eventReporter, err := config6.ProvideEventReporter(config11, logger, tracerProvider)
	if err != nil {
		return nil, err
	}
	routeParamManager := chi.NewRouteParamManager()
	authService, err := authentication2.ProvideService(ctx, logger, authenticationConfig, authenticator, dataManager, householdUserMembershipDataManager, serverEncoderDecoder, tracerProvider, publisherProvider, featureFlagManager, eventReporter, routeParamManager)
	if err != nil {
		return nil, err
	}
	usersConfig := &servicesConfig.Users
	userDataManager := database.ProvideUserDataManager(dataManager)
	householdInvitationDataManager := database.ProvideHouseholdInvitationDataManager(dataManager)
	generator := random.NewGenerator(logger, tracerProvider)
	passwordResetTokenDataManager := database.ProvidePasswordResetTokenDataManager(dataManager)
	userDataService, err := users.ProvideUsersService(ctx, usersConfig, authenticationConfig, logger, userDataManager, householdInvitationDataManager, householdUserMembershipDataManager, authenticator, serverEncoderDecoder, routeParamManager, tracerProvider, publisherProvider, generator, passwordResetTokenDataManager, featureFlagManager, eventReporter)
	if err != nil {
		return nil, err
	}
	householdsConfig := servicesConfig.Households
	householdDataManager := database.ProvideHouseholdDataManager(dataManager)
	householdDataService, err := households.ProvideService(logger, householdsConfig, householdDataManager, householdUserMembershipDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider, generator)
	if err != nil {
		return nil, err
	}
	householdinvitationsConfig := &servicesConfig.HouseholdInvitations
	householdInvitationDataService, err := householdinvitations.ProvideHouseholdInvitationsService(logger, householdinvitationsConfig, userDataManager, householdInvitationDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider, generator)
	if err != nil {
		return nil, err
	}
	validinstrumentsConfig := &servicesConfig.ValidInstruments
	config12 := &cfg.Search
	validInstrumentDataManager := database.ProvideValidInstrumentDataManager(dataManager)
	validInstrumentDataService, err := validinstruments.ProvideService(ctx, logger, validinstrumentsConfig, config12, validInstrumentDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientsConfig := &servicesConfig.ValidIngredients
	validIngredientDataManager := database.ProvideValidIngredientDataManager(dataManager)
	validIngredientDataService, err := validingredients.ProvideService(ctx, logger, validingredientsConfig, config12, validIngredientDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientgroupsConfig := &servicesConfig.ValidIngredientGroups
	validIngredientGroupDataManager := database.ProvideValidIngredientGroupDataManager(dataManager)
	validIngredientGroupDataService, err := validingredientgroups.ProvideService(ctx, logger, validingredientgroupsConfig, validIngredientGroupDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validpreparationsConfig := &servicesConfig.ValidPreparations
	validPreparationDataManager := database.ProvideValidPreparationDataManager(dataManager)
	validPreparationDataService, err := validpreparations.ProvideService(ctx, logger, validpreparationsConfig, config12, validPreparationDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
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
	mealDataService, err := meals.ProvideService(ctx, logger, mealsConfig, config12, mealDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipesConfig := &servicesConfig.Recipes
	recipeDataManager := database.ProvideRecipeDataManager(dataManager)
	recipeMediaDataManager := database.ProvideRecipeMediaDataManager(dataManager)
	recipeAnalyzer := recipeanalysis.NewRecipeAnalyzer(logger, tracerProvider)
	mediaUploadProcessor := images.NewImageUploadProcessor(logger, tracerProvider)
	recipeDataService, err := recipes.ProvideService(ctx, logger, recipesConfig, config12, recipeDataManager, recipeMediaDataManager, recipeAnalyzer, serverEncoderDecoder, routeParamManager, publisherProvider, mediaUploadProcessor, tracerProvider)
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
	validMeasurementUnitDataService, err := validmeasurementunits.ProvideService(ctx, logger, validmeasurementunitsConfig, config12, validMeasurementUnitDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validingredientstatesConfig := &servicesConfig.ValidIngredientStates
	validIngredientStateDataManager := database.ProvideValidIngredientStateDataManager(dataManager)
	validIngredientStateDataService, err := validingredientstates.ProvideService(ctx, logger, validingredientstatesConfig, config12, validIngredientStateDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
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
	validmeasurementunitconversionsConfig := &servicesConfig.ValidMeasurementUnitConversions
	validMeasurementUnitConversionDataManager := database.ProvideValidMeasurementUnitConversionDataManager(dataManager)
	validMeasurementUnitConversionDataService, err := validmeasurementunitconversions.ProvideService(ctx, logger, validmeasurementunitconversionsConfig, validMeasurementUnitConversionDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	recipestepcompletionconditionsConfig := &servicesConfig.RecipeStepCompletionConditions
	recipeStepCompletionConditionDataManager := database.ProvideRecipeStepCompletionConditionDataManager(dataManager)
	recipeStepCompletionConditionDataService, err := recipestepcompletionconditions.ProvideService(logger, recipestepcompletionconditionsConfig, recipeStepCompletionConditionDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
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
	adminService := admin.ProvideService(logger, adminUserDataManager, serverEncoderDecoder, tracerProvider)
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
	oAuth2ClientDataService, err := oauth2clients.ProvideOAuth2ClientsService(logger, oAuth2ClientDataManager, serverEncoderDecoder, routeParamManager, oauth2clientsConfig, tracerProvider, generator, publisherProvider)
	if err != nil {
		return nil, err
	}
	validvesselsConfig := &servicesConfig.ValidVessels
	validVesselDataManager := database.ProvideValidVesselDataManager(dataManager)
	validVesselDataService, err := validvessels.ProvideService(ctx, logger, validvesselsConfig, config12, validVesselDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	validpreparationvesselsConfig := &servicesConfig.ValidPreparationVessels
	validPreparationVesselDataManager := database.ProvideValidPreparationVesselDataManager(dataManager)
	validPreparationVesselDataService, err := validpreparationvessels.ProvideService(logger, validpreparationvesselsConfig, validPreparationVesselDataManager, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	workersConfig := &servicesConfig.Workers
	workerService, err := workers.ProvideService(ctx, logger, workersConfig, dataManager, serverEncoderDecoder, publisherProvider, tracerProvider, recipeAnalyzer)
	if err != nil {
		return nil, err
	}
	usernotificationsConfig := &servicesConfig.UserNotifications
	userNotificationDataManager := database.ProvideUserNotificationDataManager(dataManager)
	userNotificationDataService, err := usernotifications.ProvideService(ctx, logger, usernotificationsConfig, serverEncoderDecoder, routeParamManager, publisherProvider, tracerProvider, userNotificationDataManager)
	if err != nil {
		return nil, err
	}
	auditLogEntryDataManager := database.ProvideAuditLogEntryDataManager(dataManager)
	auditLogEntryDataService, err := auditlogentries.ProvideService(ctx, logger, auditLogEntryDataManager, serverEncoderDecoder, routeParamManager, tracerProvider)
	if err != nil {
		return nil, err
	}
	config13 := &servicesConfig.DataPrivacy
	dataPrivacyDataManager := database.ProvideDataPrivacyDataManager(dataManager)
	dataPrivacyService, err := workers2.ProvideService(ctx, logger, config13, dataPrivacyDataManager, serverEncoderDecoder, publisherProvider, tracerProvider)
	if err != nil {
		return nil, err
	}
	server, err := http.ProvideHTTPServer(ctx, httpConfig, dataManager, logger, router, tracerProvider, authService, userDataService, householdDataService, householdInvitationDataService, validInstrumentDataService, validIngredientDataService, validIngredientGroupDataService, validPreparationDataService, validIngredientPreparationDataService, mealDataService, recipeDataService, recipeStepDataService, recipeStepProductDataService, recipeStepInstrumentDataService, recipeStepIngredientDataService, mealPlanDataService, mealPlanOptionDataService, mealPlanOptionVoteDataService, validMeasurementUnitDataService, validIngredientStateDataService, validPreparationInstrumentDataService, validIngredientMeasurementUnitDataService, mealPlanEventDataService, mealPlanTaskDataService, recipePrepTaskDataService, mealPlanGroceryListItemDataService, validMeasurementUnitConversionDataService, recipeStepCompletionConditionDataService, validIngredientStateIngredientDataService, recipeStepVesselDataService, webhookDataService, adminService, serviceSettingDataService, serviceSettingConfigurationDataService, userIngredientPreferenceDataService, recipeRatingDataService, householdInstrumentOwnershipDataService, oAuth2ClientDataService, validVesselDataService, validPreparationVesselDataService, workerService, userNotificationDataService, auditLogEntryDataService, dataPrivacyService)
	if err != nil {
		return nil, err
	}
	return server, nil
}
