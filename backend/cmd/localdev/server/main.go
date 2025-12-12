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
			_, err = localdev.CreatePremadeAdminUser(ctx, logger, tracerProvider, repo, dbClient, premadeAdminUser)
			return err
		}),
		// Create OAuth2 client
		localdev.WithOAuth2Repository(func(ctx context.Context, repo oauth.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
			_, err = repo.CreateOAuth2Client(ctx, &oauth.OAuth2ClientDatabaseCreationInput{
				ID:           strings.Repeat("b", 20),
				Name:         "localdev_admin_client",
				Description:  "localdev admin client",
				ClientID:     strings.Repeat("A", oauth.ClientIDSize),
				ClientSecret: strings.Repeat("A", oauth.ClientSecretSize),
			})
			return err
		}),
		//// Create valid enumerations and bridge types
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
	ingredients := []*mealplanning.ValidIngredientDatabaseCreationInput{
		{ID: identifiers.New(), Name: "garlic", Description: "Fresh garlic cloves", PluralName: "garlic cloves", StorageInstructions: "Store in a cool, dry place", Slug: "garlic", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "onion", Description: "Yellow cooking onion", PluralName: "onions", StorageInstructions: "Store in a cool, dry, well-ventilated place", Slug: "onion", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "carrot", Description: "Fresh orange carrots", PluralName: "carrots", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "carrot", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "tomato", Description: "Ripe red tomatoes", PluralName: "tomatoes", StorageInstructions: "Store at room temperature until ripe, then refrigerate", Slug: "tomato", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "bell pepper", Description: "Red bell pepper", PluralName: "bell peppers", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "bell-pepper", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "broccoli", Description: "Fresh broccoli florets", PluralName: "broccoli", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "broccoli", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "chicken breast", Description: "Boneless, skinless chicken breast", PluralName: "chicken breasts", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "chicken-breast", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "ground beef", Description: "Lean ground beef", PluralName: "ground beef", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "ground-beef", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "salmon fillet", Description: "Fresh Atlantic salmon fillet", PluralName: "salmon fillets", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "salmon-fillet", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "milk", Description: "Whole milk", PluralName: "milk", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "milk", ContainsShellfish: false, ContainsDairy: true, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "butter", Description: "Unsalted butter", PluralName: "butter", StorageInstructions: "Keep refrigerated, can be kept at room temperature for short periods", Slug: "butter", ContainsShellfish: false, ContainsDairy: true, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "cheddar cheese", Description: "Sharp cheddar cheese", PluralName: "cheddar cheese", StorageInstructions: "Keep refrigerated, wrap tightly", Slug: "cheddar-cheese", ContainsShellfish: false, ContainsDairy: true, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "eggs", Description: "Large chicken eggs", PluralName: "eggs", StorageInstructions: "Keep refrigerated in original carton", Slug: "eggs", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: true, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "rice", Description: "Long-grain white rice", PluralName: "rice", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "rice", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "pasta", Description: "Dried spaghetti pasta", PluralName: "pasta", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "pasta", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: true, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "bread", Description: "White sandwich bread", PluralName: "bread", StorageInstructions: "Store at room temperature in a bread box or sealed bag", Slug: "bread", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: true, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "olive oil", Description: "Extra virgin olive oil", PluralName: "olive oil", StorageInstructions: "Store in a cool, dark place away from light", Slug: "olive-oil", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "salt", Description: "Fine sea salt", PluralName: "salt", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "salt", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "black pepper", Description: "Ground black pepper", PluralName: "black pepper", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "black-pepper", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "basil", Description: "Fresh basil leaves", PluralName: "basil", StorageInstructions: "Store in the refrigerator, stems in water", Slug: "basil", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "oregano", Description: "Dried oregano", PluralName: "oregano", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "oregano", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "thyme", Description: "Fresh thyme sprigs", PluralName: "thyme", StorageInstructions: "Store in the refrigerator, wrapped in damp paper towel", Slug: "thyme", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "parsley", Description: "Fresh flat-leaf parsley", PluralName: "parsley", StorageInstructions: "Store in the refrigerator, stems in water", Slug: "parsley", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "cilantro", Description: "Fresh cilantro leaves", PluralName: "cilantro", StorageInstructions: "Store in the refrigerator, stems in water", Slug: "cilantro", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "lemon", Description: "Fresh lemons", PluralName: "lemons", StorageInstructions: "Store at room temperature or in the refrigerator", Slug: "lemon", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "lime", Description: "Fresh limes", PluralName: "limes", StorageInstructions: "Store at room temperature or in the refrigerator", Slug: "lime", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "potato", Description: "Russet potatoes", PluralName: "potatoes", StorageInstructions: "Store in a cool, dark, well-ventilated place", Slug: "potato", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "sweet potato", Description: "Orange sweet potatoes", PluralName: "sweet potatoes", StorageInstructions: "Store in a cool, dark, well-ventilated place", Slug: "sweet-potato", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "spinach", Description: "Fresh baby spinach leaves", PluralName: "spinach", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "spinach", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "lettuce", Description: "Romaine lettuce", PluralName: "lettuce", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "lettuce", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "cucumber", Description: "English cucumber", PluralName: "cucumbers", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "cucumber", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "zucchini", Description: "Fresh zucchini", PluralName: "zucchini", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "zucchini", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "mushroom", Description: "White button mushrooms", PluralName: "mushrooms", StorageInstructions: "Store in the refrigerator in original packaging or paper bag", Slug: "mushroom", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "avocado", Description: "Hass avocado", PluralName: "avocados", StorageInstructions: "Store at room temperature until ripe, then refrigerate", Slug: "avocado", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "apple", Description: "Red delicious apples", PluralName: "apples", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "apple", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "banana", Description: "Yellow bananas", PluralName: "bananas", StorageInstructions: "Store at room temperature", Slug: "banana", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "strawberry", Description: "Fresh strawberries", PluralName: "strawberries", StorageInstructions: "Store in the refrigerator, do not wash until ready to use", Slug: "strawberry", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "blueberry", Description: "Fresh blueberries", PluralName: "blueberries", StorageInstructions: "Store in the refrigerator, do not wash until ready to use", Slug: "blueberry", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "almond", Description: "Raw almonds", PluralName: "almonds", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "almond", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: true, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "walnut", Description: "Raw walnut halves", PluralName: "walnuts", StorageInstructions: "Store in the refrigerator or freezer to prevent rancidity", Slug: "walnut", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: true, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "peanut", Description: "Raw peanuts", PluralName: "peanuts", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "peanut", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: true, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "tofu", Description: "Firm tofu", PluralName: "tofu", StorageInstructions: "Keep refrigerated, store in water and change daily", Slug: "tofu", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: true, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "black beans", Description: "Canned black beans", PluralName: "black beans", StorageInstructions: "Store in a cool, dry place, refrigerate after opening", Slug: "black-beans", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "chickpeas", Description: "Canned chickpeas", PluralName: "chickpeas", StorageInstructions: "Store in a cool, dry place, refrigerate after opening", Slug: "chickpeas", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "lentils", Description: "Dried brown lentils", PluralName: "lentils", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "lentils", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "quinoa", Description: "White quinoa", PluralName: "quinoa", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "quinoa", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "oats", Description: "Rolled oats", PluralName: "oats", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "oats", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "flour", Description: "All-purpose flour", PluralName: "flour", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "flour", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: true, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "sugar", Description: "Granulated white sugar", PluralName: "sugar", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "sugar", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "honey", Description: "Raw honey", PluralName: "honey", StorageInstructions: "Store at room temperature in a sealed container", Slug: "honey", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "vinegar", Description: "White wine vinegar", PluralName: "vinegar", StorageInstructions: "Store in a cool, dark place", Slug: "vinegar", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "soy sauce", Description: "Low-sodium soy sauce", PluralName: "soy sauce", StorageInstructions: "Store in a cool, dark place", Slug: "soy-sauce", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: true, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "ginger", Description: "Fresh ginger root", PluralName: "ginger", StorageInstructions: "Store in the refrigerator, wrapped in paper towel", Slug: "ginger", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "turmeric", Description: "Ground turmeric", PluralName: "turmeric", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "turmeric", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "cumin", Description: "Ground cumin", PluralName: "cumin", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "cumin", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "cinnamon", Description: "Ground cinnamon", PluralName: "cinnamon", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "cinnamon", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "paprika", Description: "Sweet paprika", PluralName: "paprika", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "paprika", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "chili powder", Description: "Mild chili powder", PluralName: "chili powder", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "chili-powder", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "cayenne pepper", Description: "Ground cayenne pepper", PluralName: "cayenne pepper", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "cayenne-pepper", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "bay leaf", Description: "Dried bay leaves", PluralName: "bay leaves", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "bay-leaf", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "rosemary", Description: "Fresh rosemary sprigs", PluralName: "rosemary", StorageInstructions: "Store in the refrigerator, wrapped in damp paper towel", Slug: "rosemary", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "sage", Description: "Fresh sage leaves", PluralName: "sage", StorageInstructions: "Store in the refrigerator, wrapped in damp paper towel", Slug: "sage", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "dill", Description: "Fresh dill weed", PluralName: "dill", StorageInstructions: "Store in the refrigerator, stems in water", Slug: "dill", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "mint", Description: "Fresh mint leaves", PluralName: "mint", StorageInstructions: "Store in the refrigerator, stems in water", Slug: "mint", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "corn", Description: "Fresh corn on the cob", PluralName: "corn", StorageInstructions: "Store in the refrigerator, keep husks on", Slug: "corn", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "peas", Description: "Frozen green peas", PluralName: "peas", StorageInstructions: "Keep frozen until ready to use", Slug: "peas", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "green beans", Description: "Fresh green beans", PluralName: "green beans", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "green-beans", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "asparagus", Description: "Fresh asparagus spears", PluralName: "asparagus", StorageInstructions: "Store in the refrigerator, stems in water", Slug: "asparagus", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "cauliflower", Description: "Fresh cauliflower head", PluralName: "cauliflower", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "cauliflower", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "cabbage", Description: "Green cabbage", PluralName: "cabbage", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "cabbage", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "celery", Description: "Fresh celery stalks", PluralName: "celery", StorageInstructions: "Store in the refrigerator crisper drawer", Slug: "celery", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "shrimp", Description: "Raw large shrimp", PluralName: "shrimp", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "shrimp", ContainsShellfish: true, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "cod", Description: "Fresh cod fillet", PluralName: "cod", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "cod", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "tuna", Description: "Fresh tuna steak", PluralName: "tuna", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "tuna", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "yogurt", Description: "Plain Greek yogurt", PluralName: "yogurt", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "yogurt", ContainsShellfish: false, ContainsDairy: true, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "sour cream", Description: "Full-fat sour cream", PluralName: "sour cream", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "sour-cream", ContainsShellfish: false, ContainsDairy: true, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "cream cheese", Description: "Plain cream cheese", PluralName: "cream cheese", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "cream-cheese", ContainsShellfish: false, ContainsDairy: true, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "parmesan cheese", Description: "Grated Parmesan cheese", PluralName: "Parmesan cheese", StorageInstructions: "Keep refrigerated, wrap tightly", Slug: "parmesan-cheese", ContainsShellfish: false, ContainsDairy: true, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "mozzarella cheese", Description: "Fresh mozzarella cheese", PluralName: "mozzarella cheese", StorageInstructions: "Keep refrigerated in brine or wrap tightly", Slug: "mozzarella-cheese", ContainsShellfish: false, ContainsDairy: true, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "chicken thigh", Description: "Bone-in chicken thighs", PluralName: "chicken thighs", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "chicken-thigh", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "pork chop", Description: "Bone-in pork chops", PluralName: "pork chops", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "pork-chop", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "ground turkey", Description: "Lean ground turkey", PluralName: "ground turkey", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "ground-turkey", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
	}

	for i, ing := range ingredients {
		validIngredient, err := repo.CreateValidIngredient(ctx, ing)
		if err != nil {
			return fmt.Errorf("failed to create valid ingredient %d: %w", i+1, err)
		}
		if i == 0 {
			firstValidIngredient = validIngredient
		}
		logger.Debug("Created ValidIngredient: " + validIngredient.Name)
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
		logger.Debug("Created ValidInstrument: " + validInstrument.Name)
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
		logger.Debug("Created ValidPreparation: " + validPreparation.Name)
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
		logger.Debug("Created ValidMeasurementUnit: " + validMeasurementUnitGram.Name)
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
		logger.Debug("Created ValidMeasurementUnit: " + validMeasurementUnitKilogram.Name)
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
		logger.Debug("Created ValidVessel: " + validVessel.Name)
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
		logger.Debug("Created ValidIngredientState: " + validIngredientState.Name)
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
	logger.Debug("Created ValidPreparationInstrument: slicing + chef's knife")

	// ValidIngredientMeasurementUnit (Garlic can be measured in Grams)
	_, err = repo.CreateValidIngredientMeasurementUnit(ctx, &mealplanning.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     identifiers.New(),
		ValidIngredientID:      firstValidIngredient.ID,
		ValidMeasurementUnitID: firstValidMeasurementUnitGram.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create valid ingredient measurement unit: %w", err)
	}
	logger.Debug("Created ValidIngredientMeasurementUnit: garlic + gram")

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
	logger.Debug("Created ValidIngredientStateIngredient: garlic + whole")

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
	logger.Debug("Created ValidPreparationVessel: slicing + cutting board")

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
	logger.Debug("Created ValidMeasurementUnitConversion: gram -> kilogram")

	_, err = repo.CreateValidIngredientPreparation(ctx, &mealplanning.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 identifiers.New(),
		Notes:              "",
		ValidPreparationID: firstValidPreparation.ID,
		ValidIngredientID:  firstValidIngredient.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create valid ingredient preparation: %w", err)
	}
	logger.Debug("Created CreateValidIngredientPreparation: garlic -> slice")

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
	logger.Debug("Created ServiceSetting: user_theme_preference (enumerated with default)")

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
	logger.Debug("Created ServiceSetting: membership_notification_frequency (enumerated with default)")

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
	logger.Debug("Created ServiceSetting: user_language (enumerated, non-admin)")

	return nil
}
