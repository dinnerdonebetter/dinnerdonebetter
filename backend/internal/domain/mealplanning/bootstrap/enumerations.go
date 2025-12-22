package bootstrap

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
)

// Enumerations holds all the created valid enumerations and their bridge types.
type Enumerations struct {
	Ingredients      map[string]*mealplanning.ValidIngredient
	Preparations     map[string]*mealplanning.ValidPreparation
	MeasurementUnits map[string]*mealplanning.ValidMeasurementUnit
	Instruments      map[string]*mealplanning.ValidInstrument
	Vessels          map[string]*mealplanning.ValidVessel

	// Bridge table lookups (keyed by [first entity ID][second entity ID])
	IngredientPreparations     map[string]map[string]*mealplanning.ValidIngredientPreparation     // [preparationID][ingredientID]
	IngredientMeasurementUnits map[string]map[string]*mealplanning.ValidIngredientMeasurementUnit // [ingredientID][unitID]
	PreparationInstruments     map[string]map[string]*mealplanning.ValidPreparationInstrument     // [preparationID][instrumentID]
	PreparationVessels         map[string]map[string]*mealplanning.ValidPreparationVessel         // [preparationID][vesselID]
}

// CreateEnumerations creates a comprehensive set of valid enumerations for local development.
func CreateEnumerations(ctx context.Context, repo mealplanning.Repository, logger logging.Logger) (*Enumerations, error) {
	const count = 75

	enums := &Enumerations{
		Ingredients:      make(map[string]*mealplanning.ValidIngredient),
		Preparations:     make(map[string]*mealplanning.ValidPreparation),
		MeasurementUnits: make(map[string]*mealplanning.ValidMeasurementUnit),
		Instruments:      make(map[string]*mealplanning.ValidInstrument),
		Vessels:          make(map[string]*mealplanning.ValidVessel),

		// Bridge table lookups
		IngredientPreparations:     make(map[string]map[string]*mealplanning.ValidIngredientPreparation),
		IngredientMeasurementUnits: make(map[string]map[string]*mealplanning.ValidIngredientMeasurementUnit),
		PreparationInstruments:     make(map[string]map[string]*mealplanning.ValidPreparationInstrument),
		PreparationVessels:         make(map[string]map[string]*mealplanning.ValidPreparationVessel),
	}

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
		// Additional ingredients for pan-seared steak recipe
		{ID: identifiers.New(), Name: "ribeye steak", Description: "Bone-in or boneless ribeye steak, thick-cut", PluralName: "ribeye steaks", StorageInstructions: "Keep refrigerated at or below 40°F, use within 3-5 days", Slug: "ribeye-steak", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "vegetable oil", Description: "Neutral-flavored vegetable or canola oil", PluralName: "vegetable oil", StorageInstructions: "Store in a cool, dark place", Slug: "vegetable-oil", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "shallot", Description: "Fresh shallots", PluralName: "shallots", StorageInstructions: "Store in a cool, dry, well-ventilated place", Slug: "shallot", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
	}

	for i, ing := range ingredients {
		validIngredient, err := repo.CreateValidIngredient(ctx, ing)
		if err != nil {
			return nil, fmt.Errorf("failed to create valid ingredient %d: %w", i+1, err)
		}
		if i == 0 {
			firstValidIngredient = validIngredient
		}
		// Store ingredients we'll need for recipes
		enums.Ingredients[validIngredient.Name] = validIngredient
		logger.Debug("Created ValidIngredient: " + validIngredient.Name)
	}

	// Create ValidInstruments
	instruments := []struct {
		name        string
		description string
		pluralName  string
		slug        string
		mapKey      string // key to store in enums.Instruments map (empty = don't store)
	}{
		{"chef's knife", "A sharp chef's knife for cutting and chopping", "chef's knives", "chefs-knife", "knife"},
		{"paper towels", "Absorbent paper towels for drying", "paper towels", "paper-towels", "paper towels"},
		{"tongs", "Kitchen tongs for flipping and handling food", "tongs", "tongs", "tongs"},
		{"spoon", "Large spoon for basting", "spoons", "spoon", "spoon"},
		{"instant-read thermometer", "Digital thermometer for checking internal temperature", "instant-read thermometers", "instant-read-thermometer", "instant-read thermometer"},
	}
	for i, inst := range instruments {
		validInstrument, err2 := repo.CreateValidInstrument(ctx, &mealplanning.ValidInstrumentDatabaseCreationInput{
			ID:                             identifiers.New(),
			Name:                           inst.name,
			Description:                    inst.description,
			PluralName:                     inst.pluralName,
			Slug:                           inst.slug,
			DisplayInSummaryLists:          true,
			IncludeInGeneratedInstructions: true,
		})
		if err2 != nil {
			return nil, fmt.Errorf("failed to create instrument %s: %w", inst.name, err2)
		}
		if i == 0 {
			firstValidInstrument = validInstrument
		}
		if inst.mapKey != "" {
			enums.Instruments[inst.mapKey] = validInstrument
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
			return nil, fmt.Errorf("failed to create valid preparation %d: %w", i, err)
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
			return nil, fmt.Errorf("failed to create valid measurement unit (gram) %d: %w", i, err)
		}
		if i == 1 {
			firstValidMeasurementUnitGram = validMeasurementUnitGram
			enums.MeasurementUnits["gram"] = validMeasurementUnitGram
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
			return nil, fmt.Errorf("failed to create valid measurement unit (kilogram) %d: %w", i, err)
		}
		if i == 1 {
			firstValidMeasurementUnitKilogram = validMeasurementUnitKilogram
		}
		logger.Debug("Created ValidMeasurementUnit: " + validMeasurementUnitKilogram.Name)
	}

	// Create additional measurement units
	measurementUnits := []struct {
		name        string
		description string
		pluralName  string
		slug        string
		volumetric  bool
		metric      bool
	}{
		{"unit", "A generic unit of measurement for recipe products", "units", "unit", false, false},
		{"milliliter", "Metric unit of volume", "milliliters", "milliliter", true, true},
		{"sprig", "A small stem with leaves, typically herbs", "sprigs", "sprig", false, false},
	}
	for _, unit := range measurementUnits {
		validUnit, err2 := repo.CreateValidMeasurementUnit(ctx, &mealplanning.ValidMeasurementUnitDatabaseCreationInput{
			ID:          identifiers.New(),
			Name:        unit.name,
			Description: unit.description,
			PluralName:  unit.pluralName,
			Slug:        unit.slug,
			Volumetric:  unit.volumetric,
			Universal:   true,
			Metric:      unit.metric,
			Imperial:    false,
		})
		if err2 != nil {
			return nil, fmt.Errorf("failed to create measurement unit %s: %w", unit.name, err2)
		}
		enums.MeasurementUnits[unit.name] = validUnit
		logger.Debug("Created ValidMeasurementUnit: " + validUnit.Name)
	}

	// Create 75 ValidVessels
	for i := 1; i <= count; i++ {
		validVessel, err2 := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
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
		if err2 != nil {
			return nil, fmt.Errorf("failed to create valid vessel %d: %w", i, err2)
		}
		if i == 1 {
			firstValidVessel = validVessel
			enums.Vessels["cutting board"] = validVessel
		}
		logger.Debug("Created ValidVessel: " + validVessel.Name)
	}

	// Create 75 ValidIngredientStates
	for i := 1; i <= count; i++ {
		validIngredientState, err2 := repo.CreateValidIngredientState(ctx, &mealplanning.ValidIngredientStateDatabaseCreationInput{
			ID:            identifiers.New(),
			Name:          fmt.Sprintf("slice %d", i),
			Description:   "a sliced ingredient",
			AttributeType: mealplanning.ValidIngredientStateAttributeTypeOther,
			PastTense:     "sliced",
			Slug:          fmt.Sprintf("slice-%d", i),
		})
		if err2 != nil {
			return nil, fmt.Errorf("failed to create valid ingredient state %d: %w", i, err2)
		}
		if i == 1 {
			firstValidIngredientState = validIngredientState
		}
		logger.Debug("Created ValidIngredientState: " + validIngredientState.Name)
	}

	// Create bridge types using first instances

	// ValidPreparationInstrument (Slicing requires Chef's Knife)
	createdVPI, err := repo.CreateValidPreparationInstrument(ctx, &mealplanning.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 identifiers.New(),
		ValidPreparationID: firstValidPreparation.ID,
		ValidInstrumentID:  firstValidInstrument.ID,
		Notes:              "A chef's knife is commonly used for slicing",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create valid preparation instrument: %w", err)
	}
	// Store in lookup map
	if enums.PreparationInstruments[firstValidPreparation.ID] == nil {
		enums.PreparationInstruments[firstValidPreparation.ID] = make(map[string]*mealplanning.ValidPreparationInstrument)
	}
	enums.PreparationInstruments[firstValidPreparation.ID][firstValidInstrument.ID] = createdVPI
	logger.Debug("Created ValidPreparationInstrument: slicing + chef's knife")

	// ValidIngredientMeasurementUnit (Garlic can be measured in Grams)
	createdVIMU, err := repo.CreateValidIngredientMeasurementUnit(ctx, &mealplanning.ValidIngredientMeasurementUnitDatabaseCreationInput{
		ID:                     identifiers.New(),
		ValidIngredientID:      firstValidIngredient.ID,
		ValidMeasurementUnitID: firstValidMeasurementUnitGram.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create valid ingredient measurement unit: %w", err)
	}
	// Store in lookup map
	if enums.IngredientMeasurementUnits[firstValidIngredient.ID] == nil {
		enums.IngredientMeasurementUnits[firstValidIngredient.ID] = make(map[string]*mealplanning.ValidIngredientMeasurementUnit)
	}
	enums.IngredientMeasurementUnits[firstValidIngredient.ID][firstValidMeasurementUnitGram.ID] = createdVIMU
	logger.Debug("Created ValidIngredientMeasurementUnit: garlic + gram")

	// ValidIngredientStateIngredient (Garlic can be in Whole state)
	_, err = repo.CreateValidIngredientStateIngredient(ctx, &mealplanning.ValidIngredientStateIngredientDatabaseCreationInput{
		ID:                     identifiers.New(),
		ValidIngredientID:      firstValidIngredient.ID,
		ValidIngredientStateID: firstValidIngredientState.ID,
		Notes:                  "Whole garlic cloves",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create valid ingredient state ingredient: %w", err)
	}
	logger.Debug("Created ValidIngredientStateIngredient: garlic + whole")

	// ValidPreparationVessel (Slicing can be done on a Cutting Board)
	createdVPV, err := repo.CreateValidPreparationVessel(ctx, &mealplanning.ValidPreparationVesselDatabaseCreationInput{
		ID:                 identifiers.New(),
		ValidPreparationID: firstValidPreparation.ID,
		ValidVesselID:      firstValidVessel.ID,
		Notes:              "Slicing is typically done on a cutting board",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create valid preparation vessel: %w", err)
	}
	// Store in lookup map
	if enums.PreparationVessels[firstValidPreparation.ID] == nil {
		enums.PreparationVessels[firstValidPreparation.ID] = make(map[string]*mealplanning.ValidPreparationVessel)
	}
	enums.PreparationVessels[firstValidPreparation.ID][firstValidVessel.ID] = createdVPV
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
		return nil, fmt.Errorf("failed to create valid measurement unit conversion: %w", err)
	}
	logger.Debug("Created ValidMeasurementUnitConversion: gram -> kilogram")

	createdVIP, err := repo.CreateValidIngredientPreparation(ctx, &mealplanning.ValidIngredientPreparationDatabaseCreationInput{
		ID:                 identifiers.New(),
		Notes:              "",
		ValidPreparationID: firstValidPreparation.ID,
		ValidIngredientID:  firstValidIngredient.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create valid ingredient preparation: %w", err)
	}
	// Store in lookup map
	if enums.IngredientPreparations[firstValidPreparation.ID] == nil {
		enums.IngredientPreparations[firstValidPreparation.ID] = make(map[string]*mealplanning.ValidIngredientPreparation)
	}
	enums.IngredientPreparations[firstValidPreparation.ID][firstValidIngredient.ID] = createdVIP
	logger.Debug("Created CreateValidIngredientPreparation: garlic -> slice")

	// Create additional vessels needed for recipes
	pan, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "pan",
		Description:                    "A frying pan for cooking",
		PluralName:                     "pans",
		Slug:                           "pan",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             200,
		LengthInMillimeters:            200,
		HeightInMillimeters:            50,
		Shape:                          mealplanning.VesselShapeCylinder,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create pan vessel: %w", err)
	}
	enums.Vessels["pan"] = pan

	pot, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "pot",
		Description:                    "A cooking pot",
		PluralName:                     "pots",
		Slug:                           "pot",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             200,
		LengthInMillimeters:            200,
		HeightInMillimeters:            150,
		Shape:                          mealplanning.VesselShapeCylinder,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create pot vessel: %w", err)
	}
	enums.Vessels["pot"] = pot

	bakingSheet, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "baking sheet",
		Description:                    "A flat baking sheet for the oven",
		PluralName:                     "baking sheets",
		Slug:                           "baking-sheet",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             300,
		LengthInMillimeters:            450,
		HeightInMillimeters:            5,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create baking sheet vessel: %w", err)
	}
	enums.Vessels["baking sheet"] = bakingSheet

	steamer, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "steamer",
		Description:                    "A steamer basket for steaming food",
		PluralName:                     "steamers",
		Slug:                           "steamer",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             200,
		LengthInMillimeters:            200,
		HeightInMillimeters:            80,
		Shape:                          mealplanning.VesselShapeCylinder,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create steamer vessel: %w", err)
	}
	enums.Vessels["steamer"] = steamer

	// Create additional vessels for steak recipe
	sheetPan, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "sheet pan",
		Description:                    "A flat rimmed baking sheet",
		PluralName:                     "sheet pans",
		Slug:                           "sheet-pan",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             330,
		LengthInMillimeters:            460,
		HeightInMillimeters:            25,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create sheet pan vessel: %w", err)
	}
	enums.Vessels["sheet pan"] = sheetPan

	castIronSkillet, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "cast iron skillet",
		Description:                    "A heavy-bottomed cast iron frying pan",
		PluralName:                     "cast iron skillets",
		Slug:                           "cast-iron-skillet",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             305,
		LengthInMillimeters:            305,
		HeightInMillimeters:            50,
		Shape:                          mealplanning.VesselShapeCylinder,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cast iron skillet vessel: %w", err)
	}
	enums.Vessels["cast iron skillet"] = castIronSkillet

	servingPlate, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "serving plate",
		Description:                    "A large plate for serving food",
		PluralName:                     "serving plates",
		Slug:                           "serving-plate",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             280,
		LengthInMillimeters:            280,
		HeightInMillimeters:            25,
		Shape:                          mealplanning.VesselShapeCylinder,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create serving plate vessel: %w", err)
	}
	enums.Vessels["serving plate"] = servingPlate

	// Create real preparations that we'll use for recipes
	prepInputs := []struct {
		name        string
		description string
		pastTense   string
		slug        string
		tempReq     bool
		timeReq     bool
	}{
		{"grill", "Cook over direct heat on a grill", "grilled", "grill", true, true},
		{"steam", "Cook with steam", "steamed", "steam", false, true},
		{"roast", "Cook in an oven with dry heat", "roasted", "roast", true, true},
		{"sauté", "Cook quickly in a small amount of fat", "sautéed", "saute", false, true},
		{"boil", "Cook in boiling liquid", "boiled", "boil", false, true},
		{"simmer", "Cook in liquid just below boiling", "simmered", "simmer", false, true},
		{"bake", "Cook in an oven", "baked", "bake", true, true},
		{"season", "Add salt, pepper, and other seasonings", "seasoned", "season", false, false},
		{"chop", "Cut into rough pieces", "chopped", "chop", false, false},
		{"dice", "Cut into small cubes", "diced", "dice", false, false},
		{"mince", "Cut into very small pieces", "minced", "mince", false, false},
		{"slice", "Cut into thin, flat pieces", "sliced", "slice", false, false},
		{"cook", "Prepare food by heating", "cooked", "cook", false, true},
		{"stir", "Mix ingredients together", "stirred", "stir", false, false},
		{"mix", "Combine ingredients together", "mixed", "mix", false, false},
		// Additional preparations for pan-seared steak recipe
		{"pan-sear", "Cook in a hot pan with oil to develop a crust", "pan-seared", "pan-sear", false, true},
		{"baste", "Spoon hot fat or liquid over food while cooking", "basted", "baste", true, true},
		{"rest", "Allow food to sit after cooking to redistribute juices", "rested", "rest", false, true},
	}

	for i := range prepInputs {
		prep := &prepInputs[i]
		validPrep, err2 := repo.CreateValidPreparation(ctx, &mealplanning.ValidPreparationDatabaseCreationInput{
			ID:                          identifiers.New(),
			Name:                        prep.name,
			Description:                 prep.description,
			Slug:                        prep.slug,
			PastTense:                   prep.pastTense,
			YieldsNothing:               false,
			RestrictToIngredients:       false,
			TemperatureRequired:         prep.tempReq,
			TimeEstimateRequired:        prep.timeReq,
			ConditionExpressionRequired: false,
			ConsumesVessel:              false,
			OnlyForVessels:              false,
		})
		if err2 != nil {
			return nil, fmt.Errorf("failed to create preparation %s: %w", prep.name, err2)
		}
		enums.Preparations[prep.name] = validPrep
		logger.Debug("Created ValidPreparation: " + validPrep.Name)
	}

	// Create bridge table entries for steak recipe
	if err = createSteakRecipeBridgeEntries(ctx, repo, logger, enums); err != nil {
		return nil, err
	}

	return enums, nil
}

// createSteakRecipeBridgeEntries creates all the bridge table entries needed for the steak recipe.
func createSteakRecipeBridgeEntries(ctx context.Context, repo mealplanning.Repository, logger logging.Logger, enums *Enumerations) error {
	// Helper to create ValidIngredientPreparation and store in map
	createVIP := func(prep *mealplanning.ValidPreparation, ing *mealplanning.ValidIngredient) error {
		vip, err := repo.CreateValidIngredientPreparation(ctx, &mealplanning.ValidIngredientPreparationDatabaseCreationInput{
			ID:                 identifiers.New(),
			ValidPreparationID: prep.ID,
			ValidIngredientID:  ing.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to create VIP %s+%s: %w", prep.Name, ing.Name, err)
		}
		if enums.IngredientPreparations[prep.ID] == nil {
			enums.IngredientPreparations[prep.ID] = make(map[string]*mealplanning.ValidIngredientPreparation)
		}
		enums.IngredientPreparations[prep.ID][ing.ID] = vip
		logger.Debug(fmt.Sprintf("Created ValidIngredientPreparation: %s + %s", prep.Name, ing.Name))
		return nil
	}

	// Helper to create ValidIngredientMeasurementUnit and store in map
	createVIMU := func(ing *mealplanning.ValidIngredient, unit *mealplanning.ValidMeasurementUnit) error {
		vimu, err := repo.CreateValidIngredientMeasurementUnit(ctx, &mealplanning.ValidIngredientMeasurementUnitDatabaseCreationInput{
			ID:                     identifiers.New(),
			ValidIngredientID:      ing.ID,
			ValidMeasurementUnitID: unit.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to create VIMU %s+%s: %w", ing.Name, unit.Name, err)
		}
		if enums.IngredientMeasurementUnits[ing.ID] == nil {
			enums.IngredientMeasurementUnits[ing.ID] = make(map[string]*mealplanning.ValidIngredientMeasurementUnit)
		}
		enums.IngredientMeasurementUnits[ing.ID][unit.ID] = vimu
		logger.Debug(fmt.Sprintf("Created ValidIngredientMeasurementUnit: %s + %s", ing.Name, unit.Name))
		return nil
	}

	// Helper to create ValidPreparationInstrument and store in map
	createVPI := func(prep *mealplanning.ValidPreparation, inst *mealplanning.ValidInstrument) error {
		vpi, err := repo.CreateValidPreparationInstrument(ctx, &mealplanning.ValidPreparationInstrumentDatabaseCreationInput{
			ID:                 identifiers.New(),
			ValidPreparationID: prep.ID,
			ValidInstrumentID:  inst.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to create VPI %s+%s: %w", prep.Name, inst.Name, err)
		}
		if enums.PreparationInstruments[prep.ID] == nil {
			enums.PreparationInstruments[prep.ID] = make(map[string]*mealplanning.ValidPreparationInstrument)
		}
		enums.PreparationInstruments[prep.ID][inst.ID] = vpi
		logger.Debug(fmt.Sprintf("Created ValidPreparationInstrument: %s + %s", prep.Name, inst.Name))
		return nil
	}

	// Helper to create ValidPreparationVessel and store in map
	createVPV := func(prep *mealplanning.ValidPreparation, vessel *mealplanning.ValidVessel) error {
		vpv, err := repo.CreateValidPreparationVessel(ctx, &mealplanning.ValidPreparationVesselDatabaseCreationInput{
			ID:                 identifiers.New(),
			ValidPreparationID: prep.ID,
			ValidVesselID:      vessel.ID,
		})
		if err != nil {
			return fmt.Errorf("failed to create VPV %s+%s: %w", prep.Name, vessel.Name, err)
		}
		if enums.PreparationVessels[prep.ID] == nil {
			enums.PreparationVessels[prep.ID] = make(map[string]*mealplanning.ValidPreparationVessel)
		}
		enums.PreparationVessels[prep.ID][vessel.ID] = vpv
		logger.Debug(fmt.Sprintf("Created ValidPreparationVessel: %s + %s", prep.Name, vessel.Name))
		return nil
	}

	// Get preparations
	seasonPrep := enums.Preparations["season"]
	panSearPrep := enums.Preparations["pan-sear"]
	bastePrep := enums.Preparations["baste"]
	restPrep := enums.Preparations["rest"]

	// Get ingredients
	ribeye := enums.Ingredients["ribeye steak"]
	salt := enums.Ingredients["salt"]
	blackPepper := enums.Ingredients["black pepper"]
	vegetableOil := enums.Ingredients["vegetable oil"]
	butter := enums.Ingredients["butter"]
	thyme := enums.Ingredients["thyme"]
	rosemary := enums.Ingredients["rosemary"]
	shallot := enums.Ingredients["shallot"]

	// Get measurement units
	unitMeasurement := enums.MeasurementUnits["unit"]
	gramMeasurement := enums.MeasurementUnits["gram"]
	milliliterMeasurement := enums.MeasurementUnits["milliliter"]
	sprigMeasurement := enums.MeasurementUnits["sprig"]

	// Get instruments
	paperTowels := enums.Instruments["paper towels"]
	tongs := enums.Instruments["tongs"]
	spoon := enums.Instruments["spoon"]
	thermometer := enums.Instruments["instant-read thermometer"]

	// Get vessels
	sheetPan := enums.Vessels["sheet pan"]
	castIronSkillet := enums.Vessels["cast iron skillet"]
	servingPlate := enums.Vessels["serving plate"]

	// === SEASON PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(seasonPrep, ribeye); err != nil {
		return err
	}
	if err := createVIP(seasonPrep, salt); err != nil {
		return err
	}
	if err := createVIP(seasonPrep, blackPepper); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links
	if err := createVIMU(ribeye, unitMeasurement); err != nil {
		return err
	}
	if err := createVIMU(salt, gramMeasurement); err != nil {
		return err
	}
	if err := createVIMU(blackPepper, gramMeasurement); err != nil {
		return err
	}

	// Preparation-Instrument links
	if err := createVPI(seasonPrep, paperTowels); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(seasonPrep, sheetPan); err != nil {
		return err
	}

	// === PAN-SEAR PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(panSearPrep, vegetableOil); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links
	if err := createVIMU(vegetableOil, milliliterMeasurement); err != nil {
		return err
	}

	// Preparation-Instrument links
	if err := createVPI(panSearPrep, tongs); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(panSearPrep, castIronSkillet); err != nil {
		return err
	}

	// === BASTE PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(bastePrep, butter); err != nil {
		return err
	}
	if err := createVIP(bastePrep, thyme); err != nil {
		return err
	}
	if err := createVIP(bastePrep, rosemary); err != nil {
		return err
	}
	if err := createVIP(bastePrep, shallot); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links
	if err := createVIMU(butter, gramMeasurement); err != nil {
		return err
	}
	if err := createVIMU(thyme, sprigMeasurement); err != nil {
		return err
	}
	if err := createVIMU(rosemary, sprigMeasurement); err != nil {
		return err
	}
	if err := createVIMU(shallot, gramMeasurement); err != nil {
		return err
	}

	// Preparation-Instrument links
	if err := createVPI(bastePrep, spoon); err != nil {
		return err
	}
	if err := createVPI(bastePrep, thermometer); err != nil {
		return err
	}
	if err := createVPI(bastePrep, tongs); err != nil {
		return err
	}

	// === REST PREPARATION ===
	// Preparation-Instrument links
	if err := createVPI(restPrep, tongs); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(restPrep, servingPlate); err != nil {
		return err
	}

	return nil
}
