package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/domain/settings"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	apiConfigurationFilepath = "deploy/environments/testing/config_files/integration-tests-config.json"
)

func main() {
	ctx := context.Background()

	// create premade admin user
	premadeAdminUser := &identity.User{
		ID:              strings.Repeat("a", 20),
		TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		EmailAddress:    "integration_tests@example.email",
		Username:        "admin_user",
		HashedPassword:  "admin_pass",
	}

	apiConfig, err := config.LoadConfigFromPath[config.APIServiceConfig](ctx, apiConfigurationFilepath)
	if err != nil {
		log.Fatal(err)
	}

	server, err := localdev.AllInOne(
		ctx,
		apiConfig,
		// Create admin user
		localdev.WithIdentityRepository(func(ctx context.Context, repo identity.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider, dbClient database.Client) error {
			_, err := localdev.CreatePremadeAdminUser(ctx, logger, tracerProvider, repo, dbClient, premadeAdminUser)
			return err
		}),
		// Create OAuth2 client
		localdev.WithOAuth2Repository(func(ctx context.Context, repo oauth.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
			_, err := repo.CreateOAuth2Client(ctx, &oauth.OAuth2ClientDatabaseCreationInput{
				ID:           strings.Repeat("b", 20),
				Name:         "localdev_admin_client",
				Description:  "localdev admin client",
				ClientID:     strings.Repeat("A", oauth.ClientIDSize),
				ClientSecret: strings.Repeat("A", oauth.ClientSecretSize),
			})
			return err
		}),
		// Create valid enumerations and bridge types
		localdev.WithMealPlanningRepository(func(ctx context.Context, repo mealplanning.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
			return createTestEnumerations(ctx, repo, logger)
		}),
		// Create example service settings
		localdev.WithSettingsRepository(func(ctx context.Context, repo settings.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
			return createExampleServiceSettings(ctx, repo, logger)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting server")
	server.Run()
}

func createTestEnumerations(ctx context.Context, repo mealplanning.Repository, logger logging.Logger) error {
	// Create ValidIngredient (Garlic)
	validIngredient, err := repo.CreateValidIngredient(ctx, &mealplanning.ValidIngredientDatabaseCreationInput{
		ID:                     identifiers.New(),
		Name:                   "garlic",
		Description:            "Fresh garlic cloves",
		PluralName:             "garlic cloves",
		StorageInstructions:    "Store in a cool, dry place",
		Slug:                   "garlic",
		ContainsShellfish:      false,
		ContainsDairy:          false,
		ContainsPeanut:         false,
		ContainsTreeNut:        false,
		ContainsEgg:            false,
		ContainsWheat:          false,
		ContainsSoy:            false,
		AnimalDerived:          false,
		RestrictToPreparations: false,
	})
	if err != nil {
		return fmt.Errorf("failed to create valid ingredient: %w", err)
	}
	logger.Info("Created ValidIngredient: " + validIngredient.Name)

	// Create ValidInstrument (Chef's Knife)
	validInstrument, err := repo.CreateValidInstrument(ctx, &mealplanning.ValidInstrumentDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "chef's knife",
		Description:                    "A sharp chef's knife for cutting and chopping",
		PluralName:                     "chef's knives",
		Slug:                           "chefs-knife",
		DisplayInSummaryLists:          true,
		IncludeInGeneratedInstructions: true,
	})
	if err != nil {
		return fmt.Errorf("failed to create valid instrument: %w", err)
	}
	logger.Info("Created ValidInstrument: " + validInstrument.Name)

	// Create ValidPreparation (Slicing)
	validPreparation, err := repo.CreateValidPreparation(ctx, &mealplanning.ValidPreparationDatabaseCreationInput{
		ID:                          identifiers.New(),
		Name:                        "slicing",
		Description:                 "Cut into thin, flat pieces",
		Slug:                        "slicing",
		PastTense:                   "sliced",
		YieldsNothing:               false,
		RestrictToIngredients:       false,
		TemperatureRequired:         false,
		TimeEstimateRequired:        false,
		ConditionExpressionRequired: false,
		ConsumesVessel:              false,
		OnlyForVessels:              false,
	})
	if err != nil {
		return fmt.Errorf("failed to create valid preparation: %w", err)
	}
	logger.Info("Created ValidPreparation: " + validPreparation.Name)

	// Create ValidMeasurementUnits (Gram and Kilogram for conversion)
	validMeasurementUnitGram, err := repo.CreateValidMeasurementUnit(ctx, &mealplanning.ValidMeasurementUnitDatabaseCreationInput{
		ID:          identifiers.New(),
		Name:        "gram",
		Description: "Metric unit of mass",
		PluralName:  "grams",
		Slug:        "gram",
		Volumetric:  false,
		Universal:   true,
		Metric:      true,
		Imperial:    false,
	})
	if err != nil {
		return fmt.Errorf("failed to create valid measurement unit (gram): %w", err)
	}
	logger.Info("Created ValidMeasurementUnit: " + validMeasurementUnitGram.Name)

	validMeasurementUnitKilogram, err := repo.CreateValidMeasurementUnit(ctx, &mealplanning.ValidMeasurementUnitDatabaseCreationInput{
		ID:          identifiers.New(),
		Name:        "kilogram",
		Description: "Metric unit of mass equal to 1000 grams",
		PluralName:  "kilograms",
		Slug:        "kilogram",
		Volumetric:  false,
		Universal:   true,
		Metric:      true,
		Imperial:    false,
	})
	if err != nil {
		return fmt.Errorf("failed to create valid measurement unit (kilogram): %w", err)
	}
	logger.Info("Created ValidMeasurementUnit: " + validMeasurementUnitKilogram.Name)

	// Create ValidVessel (Cutting Board)
	validVessel, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "cutting board",
		Description:                    "A flat surface for cutting ingredients",
		PluralName:                     "cutting boards",
		Slug:                           "cutting-board",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &validMeasurementUnitGram.ID,
		WidthInMillimeters:             300,
		LengthInMillimeters:            400,
		HeightInMillimeters:            20,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               true,
	})
	if err != nil {
		return fmt.Errorf("failed to create valid vessel: %w", err)
	}
	logger.Info("Created ValidVessel: " + validVessel.Name)

	// Create ValidIngredientState (Whole)
	validIngredientState, err := repo.CreateValidIngredientState(ctx, &mealplanning.ValidIngredientStateDatabaseCreationInput{
		ID:            identifiers.New(),
		Name:          "slice",
		Description:   "a sliced ingredient",
		AttributeType: mealplanning.ValidIngredientStateAttributeTypeOther,
		PastTense:     "sliced",
		Slug:          "slice",
	})
	if err != nil {
		return fmt.Errorf("failed to create valid ingredient state: %w", err)
	}
	logger.Info("Created ValidIngredientState: " + validIngredientState.Name)

	// Create bridge types

	// ValidPreparationInstrument (Slicing requires Chef's Knife)
	_, err = repo.CreateValidPreparationInstrument(ctx, &mealplanning.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 identifiers.New(),
		ValidPreparationID: validPreparation.ID,
		ValidInstrumentID:  validInstrument.ID,
		Notes:              "A chef's knife is commonly used for slicing",
	})
	if err != nil {
		return fmt.Errorf("failed to create valid preparation instrument: %w", err)
	}
	logger.Info("Created ValidPreparationInstrument: slicing + chef's knife")

	// ValidIngredientMeasurementUnit (Garlic can be measured in Grams)
	_, err = repo.CreateValidIngredientMeasurementUnit(ctx, &mealplanning.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     identifiers.New(),
		ValidIngredientID:      validIngredient.ID,
		ValidMeasurementUnitID: validMeasurementUnitGram.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create valid ingredient measurement unit: %w", err)
	}
	logger.Info("Created ValidIngredientMeasurementUnit: garlic + gram")

	// ValidIngredientStateIngredient (Garlic can be in Whole state)
	_, err = repo.CreateValidIngredientStateIngredient(ctx, &mealplanning.ValidIngredientStateIngredientDatabaseCreationInput{
		ID:                     identifiers.New(),
		ValidIngredientID:      validIngredient.ID,
		ValidIngredientStateID: validIngredientState.ID,
		Notes:                  "Whole garlic cloves",
	})
	if err != nil {
		return fmt.Errorf("failed to create valid ingredient state ingredient: %w", err)
	}
	logger.Info("Created ValidIngredientStateIngredient: garlic + whole")

	// ValidPreparationVessel (Slicing can be done on a Cutting Board)
	_, err = repo.CreateValidPreparationVessel(ctx, &mealplanning.ValidPreparationVesselDatabaseCreationInput{
		ID:                 identifiers.New(),
		ValidPreparationID: validPreparation.ID,
		ValidVesselID:      validVessel.ID,
		Notes:              "Slicing is typically done on a cutting board",
	})
	if err != nil {
		return fmt.Errorf("failed to create valid preparation vessel: %w", err)
	}
	logger.Info("Created ValidPreparationVessel: slicing + cutting board")

	// ValidMeasurementUnitConversion (Gram to Kilogram)
	_, err = repo.CreateValidMeasurementUnitConversion(ctx, &mealplanning.ValidMeasurementUnitConversionDatabaseCreationInput{
		ID:       identifiers.New(),
		From:     validMeasurementUnitGram.ID,
		To:       validMeasurementUnitKilogram.ID,
		Notes:    "conversion from grams to kilograms",
		Modifier: 0.001, // 1 gram = 0.001 kilograms
	})
	if err != nil {
		return fmt.Errorf("failed to create valid measurement unit conversion: %w", err)
	}
	logger.Info("Created ValidMeasurementUnitConversion: gram -> kilogram")

	_, err = repo.CreateValidIngredientPreparation(ctx, &mealplanning.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 identifiers.New(),
		Notes:              "",
		ValidPreparationID: validPreparation.ID,
		ValidIngredientID:  validIngredient.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create valid ingredient preparation: %w", err)
	}
	logger.Info("Created CreateValidIngredientPreparation: garlic -> slice")

	return nil
}

func createExampleServiceSettings(ctx context.Context, repo settings.Repository, logger logging.Logger) error {
	defaultTheme := "light"
	_, err := repo.CreateServiceSetting(ctx, &settings.ServiceSettingDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         "user_theme_preference",
		Type:         "user",
		Description:  "User's preferred theme for the application interface",
		Enumeration:  []string{"light", "dark", "auto"},
		DefaultValue: &defaultTheme,
		AdminsOnly:   true,
	})
	if err != nil {
		return fmt.Errorf("failed to create theme preference setting: %w", err)
	}
	logger.Info("Created ServiceSetting: user_theme_preference (enumerated with default)")

	defaultNotificationFreq := "daily"
	_, err = repo.CreateServiceSetting(ctx, &settings.ServiceSettingDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         "membership_notification_frequency",
		Type:         "membership",
		Description:  "How often to send notifications to membership members",
		Enumeration:  []string{"immediate", "daily", "weekly", "never"},
		DefaultValue: &defaultNotificationFreq,
		AdminsOnly:   true,
	})
	if err != nil {
		return fmt.Errorf("failed to create notification frequency setting: %w", err)
	}
	logger.Info("Created ServiceSetting: membership_notification_frequency (enumerated with default)")

	defaultLanguage := "en"
	_, err = repo.CreateServiceSetting(ctx, &settings.ServiceSettingDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         "user_language",
		Type:         "user",
		Description:  "User's preferred language for the application",
		Enumeration:  []string{"en", "es", "fr", "de", "it"},
		DefaultValue: &defaultLanguage,
		AdminsOnly:   false,
	})
	if err != nil {
		return fmt.Errorf("failed to create language setting: %w", err)
	}
	logger.Info("Created ServiceSetting: user_language (enumerated, non-admin)")

	return nil
}
