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
	const count = 75

	// Store first instances for bridge relationships
	var firstValidIngredient *mealplanning.ValidIngredient
	var firstValidInstrument *mealplanning.ValidInstrument
	var firstValidPreparation *mealplanning.ValidPreparation
	var firstValidMeasurementUnitGram *mealplanning.ValidMeasurementUnit
	var firstValidMeasurementUnitKilogram *mealplanning.ValidMeasurementUnit
	var firstValidVessel *mealplanning.ValidVessel
	var firstValidIngredientState *mealplanning.ValidIngredientState

	// Create 75 ValidIngredients
	for i := 1; i <= count; i++ {
		validIngredient, err := repo.CreateValidIngredient(ctx, &mealplanning.ValidIngredientDatabaseCreationInput{
			ID:                     identifiers.New(),
			Name:                   fmt.Sprintf("garlic %d", i),
			Description:            "Fresh garlic cloves",
			PluralName:             "garlic cloves",
			StorageInstructions:    "Store in a cool, dry place",
			Slug:                   fmt.Sprintf("garlic-%d", i),
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
			return fmt.Errorf("failed to create valid ingredient %d: %w", i, err)
		}
		if i == 1 {
			firstValidIngredient = validIngredient
		}
		logger.Info("Created ValidIngredient: " + validIngredient.Name)
	}

	// Create 75 ValidInstruments
	for i := 1; i <= count; i++ {
		validInstrument, err := repo.CreateValidInstrument(ctx, &mealplanning.ValidInstrumentDatabaseCreationInput{
			ID:                             identifiers.New(),
			Name:                           fmt.Sprintf("chef's knife %d", i),
			Description:                    "A sharp chef's knife for cutting and chopping",
			PluralName:                     "chef's knives",
			Slug:                           fmt.Sprintf("chefs-knife-%d", i),
			DisplayInSummaryLists:          true,
			IncludeInGeneratedInstructions: true,
		})
		if err != nil {
			return fmt.Errorf("failed to create valid instrument %d: %w", i, err)
		}
		if i == 1 {
			firstValidInstrument = validInstrument
		}
		logger.Info("Created ValidInstrument: " + validInstrument.Name)
	}

	// Create 75 ValidPreparations
	for i := 1; i <= count; i++ {
		validPreparation, err := repo.CreateValidPreparation(ctx, &mealplanning.ValidPreparationDatabaseCreationInput{
			ID:                          identifiers.New(),
			Name:                        fmt.Sprintf("slicing %d", i),
			Description:                 "Cut into thin, flat pieces",
			Slug:                        fmt.Sprintf("slicing-%d", i),
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
			return fmt.Errorf("failed to create valid preparation %d: %w", i, err)
		}
		if i == 1 {
			firstValidPreparation = validPreparation
		}
		logger.Info("Created ValidPreparation: " + validPreparation.Name)
	}

	// Create 75 ValidMeasurementUnits (Gram)
	for i := 1; i <= count; i++ {
		validMeasurementUnitGram, err := repo.CreateValidMeasurementUnit(ctx, &mealplanning.ValidMeasurementUnitDatabaseCreationInput{
			ID:          identifiers.New(),
			Name:        fmt.Sprintf("gram %d", i),
			Description: "Metric unit of mass",
			PluralName:  "grams",
			Slug:        fmt.Sprintf("gram-%d", i),
			Volumetric:  false,
			Universal:   true,
			Metric:      true,
			Imperial:    false,
		})
		if err != nil {
			return fmt.Errorf("failed to create valid measurement unit (gram) %d: %w", i, err)
		}
		if i == 1 {
			firstValidMeasurementUnitGram = validMeasurementUnitGram
		}
		logger.Info("Created ValidMeasurementUnit: " + validMeasurementUnitGram.Name)
	}

	// Create 75 ValidMeasurementUnits (Kilogram)
	for i := 1; i <= count; i++ {
		validMeasurementUnitKilogram, err := repo.CreateValidMeasurementUnit(ctx, &mealplanning.ValidMeasurementUnitDatabaseCreationInput{
			ID:          identifiers.New(),
			Name:        fmt.Sprintf("kilogram %d", i),
			Description: "Metric unit of mass equal to 1000 grams",
			PluralName:  "kilograms",
			Slug:        fmt.Sprintf("kilogram-%d", i),
			Volumetric:  false,
			Universal:   true,
			Metric:      true,
			Imperial:    false,
		})
		if err != nil {
			return fmt.Errorf("failed to create valid measurement unit (kilogram) %d: %w", i, err)
		}
		if i == 1 {
			firstValidMeasurementUnitKilogram = validMeasurementUnitKilogram
		}
		logger.Info("Created ValidMeasurementUnit: " + validMeasurementUnitKilogram.Name)
	}

	// Create 75 ValidVessels
	for i := 1; i <= count; i++ {
		validVessel, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
			ID:                             identifiers.New(),
			Name:                           fmt.Sprintf("cutting board %d", i),
			Description:                    "A flat surface for cutting ingredients",
			PluralName:                     "cutting boards",
			Slug:                           fmt.Sprintf("cutting-board-%d", i),
			IncludeInGeneratedInstructions: true,
			DisplayInSummaryLists:          true,
			CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
			WidthInMillimeters:             300,
			LengthInMillimeters:            400,
			HeightInMillimeters:            20,
			Shape:                          mealplanning.VesselShapeRectangle,
			UsableForStorage:               true,
		})
		if err != nil {
			return fmt.Errorf("failed to create valid vessel %d: %w", i, err)
		}
		if i == 1 {
			firstValidVessel = validVessel
		}
		logger.Info("Created ValidVessel: " + validVessel.Name)
	}

	// Create 75 ValidIngredientStates
	for i := 1; i <= count; i++ {
		validIngredientState, err := repo.CreateValidIngredientState(ctx, &mealplanning.ValidIngredientStateDatabaseCreationInput{
			ID:            identifiers.New(),
			Name:          fmt.Sprintf("slice %d", i),
			Description:   "a sliced ingredient",
			AttributeType: mealplanning.ValidIngredientStateAttributeTypeOther,
			PastTense:     "sliced",
			Slug:          fmt.Sprintf("slice-%d", i),
		})
		if err != nil {
			return fmt.Errorf("failed to create valid ingredient state %d: %w", i, err)
		}
		if i == 1 {
			firstValidIngredientState = validIngredientState
		}
		logger.Info("Created ValidIngredientState: " + validIngredientState.Name)
	}

	// Create bridge types using first instances

	// ValidPreparationInstrument (Slicing requires Chef's Knife)
	_, err := repo.CreateValidPreparationInstrument(ctx, &mealplanning.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 identifiers.New(),
		ValidPreparationID: firstValidPreparation.ID,
		ValidInstrumentID:  firstValidInstrument.ID,
		Notes:              "A chef's knife is commonly used for slicing",
	})
	if err != nil {
		return fmt.Errorf("failed to create valid preparation instrument: %w", err)
	}
	logger.Info("Created ValidPreparationInstrument: slicing + chef's knife")

	// ValidIngredientMeasurementUnit (Garlic can be measured in Grams)
	_, err = repo.CreateValidIngredientMeasurementUnit(ctx, &mealplanning.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     identifiers.New(),
		ValidIngredientID:      firstValidIngredient.ID,
		ValidMeasurementUnitID: firstValidMeasurementUnitGram.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create valid ingredient measurement unit: %w", err)
	}
	logger.Info("Created ValidIngredientMeasurementUnit: garlic + gram")

	// ValidIngredientStateIngredient (Garlic can be in Whole state)
	_, err = repo.CreateValidIngredientStateIngredient(ctx, &mealplanning.ValidIngredientStateIngredientDatabaseCreationInput{
		ID:                     identifiers.New(),
		ValidIngredientID:      firstValidIngredient.ID,
		ValidIngredientStateID: firstValidIngredientState.ID,
		Notes:                  "Whole garlic cloves",
	})
	if err != nil {
		return fmt.Errorf("failed to create valid ingredient state ingredient: %w", err)
	}
	logger.Info("Created ValidIngredientStateIngredient: garlic + whole")

	// ValidPreparationVessel (Slicing can be done on a Cutting Board)
	_, err = repo.CreateValidPreparationVessel(ctx, &mealplanning.ValidPreparationVesselDatabaseCreationInput{
		ID:                 identifiers.New(),
		ValidPreparationID: firstValidPreparation.ID,
		ValidVesselID:      firstValidVessel.ID,
		Notes:              "Slicing is typically done on a cutting board",
	})
	if err != nil {
		return fmt.Errorf("failed to create valid preparation vessel: %w", err)
	}
	logger.Info("Created ValidPreparationVessel: slicing + cutting board")

	// ValidMeasurementUnitConversion (Gram to Kilogram)
	_, err = repo.CreateValidMeasurementUnitConversion(ctx, &mealplanning.ValidMeasurementUnitConversionDatabaseCreationInput{
		ID:       identifiers.New(),
		From:     firstValidMeasurementUnitGram.ID,
		To:       firstValidMeasurementUnitKilogram.ID,
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
		ValidPreparationID: firstValidPreparation.ID,
		ValidIngredientID:  firstValidIngredient.ID,
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
