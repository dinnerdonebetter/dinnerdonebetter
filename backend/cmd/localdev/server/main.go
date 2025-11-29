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

	// Create 75 ValidIngredients - diverse list of real ingredients
	ingredients := []struct {
		name                   string
		description            string
		pluralName             string
		storageInstructions    string
		slug                   string
		containsShellfish      bool
		containsDairy          bool
		containsPeanut         bool
		containsTreeNut        bool
		containsEgg            bool
		containsWheat          bool
		containsSoy            bool
		animalDerived          bool
		restrictToPreparations bool
	}{
		{"garlic", "Fresh garlic cloves", "garlic cloves", "Store in a cool, dry place", "garlic", false, false, false, false, false, false, false, false, false},
		{"onion", "Yellow cooking onion", "onions", "Store in a cool, dry, well-ventilated place", "onion", false, false, false, false, false, false, false, false, false},
		{"carrot", "Fresh orange carrots", "carrots", "Store in the refrigerator crisper drawer", "carrot", false, false, false, false, false, false, false, false, false},
		{"tomato", "Ripe red tomatoes", "tomatoes", "Store at room temperature until ripe, then refrigerate", "tomato", false, false, false, false, false, false, false, false, false},
		{"bell pepper", "Red bell pepper", "bell peppers", "Store in the refrigerator crisper drawer", "bell-pepper", false, false, false, false, false, false, false, false, false},
		{"broccoli", "Fresh broccoli florets", "broccoli", "Store in the refrigerator crisper drawer", "broccoli", false, false, false, false, false, false, false, false, false},
		{"chicken breast", "Boneless, skinless chicken breast", "chicken breasts", "Keep refrigerated at or below 40°F", "chicken-breast", false, false, false, false, false, false, false, true, false},
		{"ground beef", "Lean ground beef", "ground beef", "Keep refrigerated at or below 40°F", "ground-beef", false, false, false, false, false, false, false, true, false},
		{"salmon fillet", "Fresh Atlantic salmon fillet", "salmon fillets", "Keep refrigerated at or below 40°F", "salmon-fillet", false, false, false, false, false, false, false, true, false},
		{"milk", "Whole milk", "milk", "Keep refrigerated at or below 40°F", "milk", false, true, false, false, false, false, false, true, false},
		{"butter", "Unsalted butter", "butter", "Keep refrigerated, can be kept at room temperature for short periods", "butter", false, true, false, false, false, false, false, true, false},
		{"cheddar cheese", "Sharp cheddar cheese", "cheddar cheese", "Keep refrigerated, wrap tightly", "cheddar-cheese", false, true, false, false, false, false, false, true, false},
		{"eggs", "Large chicken eggs", "eggs", "Keep refrigerated in original carton", "eggs", false, false, false, false, true, false, false, true, false},
		{"rice", "Long-grain white rice", "rice", "Store in a cool, dry place in an airtight container", "rice", false, false, false, false, false, false, false, false, false},
		{"pasta", "Dried spaghetti pasta", "pasta", "Store in a cool, dry place in an airtight container", "pasta", false, false, false, false, false, true, false, false, false},
		{"bread", "White sandwich bread", "bread", "Store at room temperature in a bread box or sealed bag", "bread", false, false, false, false, false, true, false, false, false},
		{"olive oil", "Extra virgin olive oil", "olive oil", "Store in a cool, dark place away from light", "olive-oil", false, false, false, false, false, false, false, false, false},
		{"salt", "Fine sea salt", "salt", "Store in a cool, dry place in an airtight container", "salt", false, false, false, false, false, false, false, false, false},
		{"black pepper", "Ground black pepper", "black pepper", "Store in a cool, dry place in an airtight container", "black-pepper", false, false, false, false, false, false, false, false, false},
		{"basil", "Fresh basil leaves", "basil", "Store in the refrigerator, stems in water", "basil", false, false, false, false, false, false, false, false, false},
		{"oregano", "Dried oregano", "oregano", "Store in a cool, dry place in an airtight container", "oregano", false, false, false, false, false, false, false, false, false},
		{"thyme", "Fresh thyme sprigs", "thyme", "Store in the refrigerator, wrapped in damp paper towel", "thyme", false, false, false, false, false, false, false, false, false},
		{"parsley", "Fresh flat-leaf parsley", "parsley", "Store in the refrigerator, stems in water", "parsley", false, false, false, false, false, false, false, false, false},
		{"cilantro", "Fresh cilantro leaves", "cilantro", "Store in the refrigerator, stems in water", "cilantro", false, false, false, false, false, false, false, false, false},
		{"lemon", "Fresh lemons", "lemons", "Store at room temperature or in the refrigerator", "lemon", false, false, false, false, false, false, false, false, false},
		{"lime", "Fresh limes", "limes", "Store at room temperature or in the refrigerator", "lime", false, false, false, false, false, false, false, false, false},
		{"potato", "Russet potatoes", "potatoes", "Store in a cool, dark, well-ventilated place", "potato", false, false, false, false, false, false, false, false, false},
		{"sweet potato", "Orange sweet potatoes", "sweet potatoes", "Store in a cool, dark, well-ventilated place", "sweet-potato", false, false, false, false, false, false, false, false, false},
		{"spinach", "Fresh baby spinach leaves", "spinach", "Store in the refrigerator crisper drawer", "spinach", false, false, false, false, false, false, false, false, false},
		{"lettuce", "Romaine lettuce", "lettuce", "Store in the refrigerator crisper drawer", "lettuce", false, false, false, false, false, false, false, false, false},
		{"cucumber", "English cucumber", "cucumbers", "Store in the refrigerator crisper drawer", "cucumber", false, false, false, false, false, false, false, false, false},
		{"zucchini", "Fresh zucchini", "zucchini", "Store in the refrigerator crisper drawer", "zucchini", false, false, false, false, false, false, false, false, false},
		{"mushroom", "White button mushrooms", "mushrooms", "Store in the refrigerator in original packaging or paper bag", "mushroom", false, false, false, false, false, false, false, false, false},
		{"avocado", "Hass avocado", "avocados", "Store at room temperature until ripe, then refrigerate", "avocado", false, false, false, false, false, false, false, false, false},
		{"apple", "Red delicious apples", "apples", "Store in the refrigerator crisper drawer", "apple", false, false, false, false, false, false, false, false, false},
		{"banana", "Yellow bananas", "bananas", "Store at room temperature", "banana", false, false, false, false, false, false, false, false, false},
		{"strawberry", "Fresh strawberries", "strawberries", "Store in the refrigerator, do not wash until ready to use", "strawberry", false, false, false, false, false, false, false, false, false},
		{"blueberry", "Fresh blueberries", "blueberries", "Store in the refrigerator, do not wash until ready to use", "blueberry", false, false, false, false, false, false, false, false, false},
		{"almond", "Raw almonds", "almonds", "Store in a cool, dry place in an airtight container", "almond", false, false, false, true, false, false, false, false, false},
		{"walnut", "Raw walnut halves", "walnuts", "Store in the refrigerator or freezer to prevent rancidity", "walnut", false, false, false, true, false, false, false, false, false},
		{"peanut", "Raw peanuts", "peanuts", "Store in a cool, dry place in an airtight container", "peanut", false, false, true, false, false, false, false, false, false},
		{"tofu", "Firm tofu", "tofu", "Keep refrigerated, store in water and change daily", "tofu", false, false, false, false, false, false, true, false, false},
		{"black beans", "Canned black beans", "black beans", "Store in a cool, dry place, refrigerate after opening", "black-beans", false, false, false, false, false, false, false, false, false},
		{"chickpeas", "Canned chickpeas", "chickpeas", "Store in a cool, dry place, refrigerate after opening", "chickpeas", false, false, false, false, false, false, false, false, false},
		{"lentils", "Dried brown lentils", "lentils", "Store in a cool, dry place in an airtight container", "lentils", false, false, false, false, false, false, false, false, false},
		{"quinoa", "White quinoa", "quinoa", "Store in a cool, dry place in an airtight container", "quinoa", false, false, false, false, false, false, false, false, false},
		{"oats", "Rolled oats", "oats", "Store in a cool, dry place in an airtight container", "oats", false, false, false, false, false, false, false, false, false},
		{"flour", "All-purpose flour", "flour", "Store in a cool, dry place in an airtight container", "flour", false, false, false, false, false, true, false, false, false},
		{"sugar", "Granulated white sugar", "sugar", "Store in a cool, dry place in an airtight container", "sugar", false, false, false, false, false, false, false, false, false},
		{"honey", "Raw honey", "honey", "Store at room temperature in a sealed container", "honey", false, false, false, false, false, false, false, false, false},
		{"vinegar", "White wine vinegar", "vinegar", "Store in a cool, dark place", "vinegar", false, false, false, false, false, false, false, false, false},
		{"soy sauce", "Low-sodium soy sauce", "soy sauce", "Store in a cool, dark place", "soy-sauce", false, false, false, false, false, false, true, false, false},
		{"ginger", "Fresh ginger root", "ginger", "Store in the refrigerator, wrapped in paper towel", "ginger", false, false, false, false, false, false, false, false, false},
		{"turmeric", "Ground turmeric", "turmeric", "Store in a cool, dry place in an airtight container", "turmeric", false, false, false, false, false, false, false, false, false},
		{"cumin", "Ground cumin", "cumin", "Store in a cool, dry place in an airtight container", "cumin", false, false, false, false, false, false, false, false, false},
		{"cinnamon", "Ground cinnamon", "cinnamon", "Store in a cool, dry place in an airtight container", "cinnamon", false, false, false, false, false, false, false, false, false},
		{"paprika", "Sweet paprika", "paprika", "Store in a cool, dry place in an airtight container", "paprika", false, false, false, false, false, false, false, false, false},
		{"chili powder", "Mild chili powder", "chili powder", "Store in a cool, dry place in an airtight container", "chili-powder", false, false, false, false, false, false, false, false, false},
		{"cayenne pepper", "Ground cayenne pepper", "cayenne pepper", "Store in a cool, dry place in an airtight container", "cayenne-pepper", false, false, false, false, false, false, false, false, false},
		{"bay leaf", "Dried bay leaves", "bay leaves", "Store in a cool, dry place in an airtight container", "bay-leaf", false, false, false, false, false, false, false, false, false},
		{"rosemary", "Fresh rosemary sprigs", "rosemary", "Store in the refrigerator, wrapped in damp paper towel", "rosemary", false, false, false, false, false, false, false, false, false},
		{"sage", "Fresh sage leaves", "sage", "Store in the refrigerator, wrapped in damp paper towel", "sage", false, false, false, false, false, false, false, false, false},
		{"dill", "Fresh dill weed", "dill", "Store in the refrigerator, stems in water", "dill", false, false, false, false, false, false, false, false, false},
		{"mint", "Fresh mint leaves", "mint", "Store in the refrigerator, stems in water", "mint", false, false, false, false, false, false, false, false, false},
		{"corn", "Fresh corn on the cob", "corn", "Store in the refrigerator, keep husks on", "corn", false, false, false, false, false, false, false, false, false},
		{"peas", "Frozen green peas", "peas", "Keep frozen until ready to use", "peas", false, false, false, false, false, false, false, false, false},
		{"green beans", "Fresh green beans", "green beans", "Store in the refrigerator crisper drawer", "green-beans", false, false, false, false, false, false, false, false, false},
		{"asparagus", "Fresh asparagus spears", "asparagus", "Store in the refrigerator, stems in water", "asparagus", false, false, false, false, false, false, false, false, false},
		{"cauliflower", "Fresh cauliflower head", "cauliflower", "Store in the refrigerator crisper drawer", "cauliflower", false, false, false, false, false, false, false, false, false},
		{"cabbage", "Green cabbage", "cabbage", "Store in the refrigerator crisper drawer", "cabbage", false, false, false, false, false, false, false, false, false},
		{"celery", "Fresh celery stalks", "celery", "Store in the refrigerator crisper drawer", "celery", false, false, false, false, false, false, false, false, false},
		{"shrimp", "Raw large shrimp", "shrimp", "Keep refrigerated at or below 40°F", "shrimp", true, false, false, false, false, false, false, true, false},
		{"cod", "Fresh cod fillet", "cod", "Keep refrigerated at or below 40°F", "cod", false, false, false, false, false, false, false, true, false},
		{"tuna", "Fresh tuna steak", "tuna", "Keep refrigerated at or below 40°F", "tuna", false, false, false, false, false, false, false, true, false},
		{"yogurt", "Plain Greek yogurt", "yogurt", "Keep refrigerated at or below 40°F", "yogurt", false, true, false, false, false, false, false, true, false},
		{"sour cream", "Full-fat sour cream", "sour cream", "Keep refrigerated at or below 40°F", "sour-cream", false, true, false, false, false, false, false, true, false},
		{"cream cheese", "Plain cream cheese", "cream cheese", "Keep refrigerated at or below 40°F", "cream-cheese", false, true, false, false, false, false, false, true, false},
		{"parmesan cheese", "Grated Parmesan cheese", "Parmesan cheese", "Keep refrigerated, wrap tightly", "parmesan-cheese", false, true, false, false, false, false, false, true, false},
		{"mozzarella cheese", "Fresh mozzarella cheese", "mozzarella cheese", "Keep refrigerated in brine or wrap tightly", "mozzarella-cheese", false, true, false, false, false, false, false, true, false},
		{"chicken thigh", "Bone-in chicken thighs", "chicken thighs", "Keep refrigerated at or below 40°F", "chicken-thigh", false, false, false, false, false, false, false, true, false},
		{"pork chop", "Bone-in pork chops", "pork chops", "Keep refrigerated at or below 40°F", "pork-chop", false, false, false, false, false, false, false, true, false},
		{"ground turkey", "Lean ground turkey", "ground turkey", "Keep refrigerated at or below 40°F", "ground-turkey", false, false, false, false, false, false, false, true, false},
	}

	for i, ing := range ingredients {
		validIngredient, err := repo.CreateValidIngredient(ctx, &mealplanning.ValidIngredientDatabaseCreationInput{
			ID:                     identifiers.New(),
			Name:                   ing.name,
			Description:            ing.description,
			PluralName:             ing.pluralName,
			StorageInstructions:    ing.storageInstructions,
			Slug:                   ing.slug,
			ContainsShellfish:      ing.containsShellfish,
			ContainsDairy:          ing.containsDairy,
			ContainsPeanut:         ing.containsPeanut,
			ContainsTreeNut:        ing.containsTreeNut,
			ContainsEgg:            ing.containsEgg,
			ContainsWheat:          ing.containsWheat,
			ContainsSoy:            ing.containsSoy,
			AnimalDerived:          ing.animalDerived,
			RestrictToPreparations: ing.restrictToPreparations,
		})
		if err != nil {
			return fmt.Errorf("failed to create valid ingredient %d: %w", i+1, err)
		}
		if i == 0 {
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
