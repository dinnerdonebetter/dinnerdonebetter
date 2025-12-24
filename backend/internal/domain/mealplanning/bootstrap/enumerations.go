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
	IngredientStates map[string]*mealplanning.ValidIngredientState

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
		IngredientStates: make(map[string]*mealplanning.ValidIngredientState),

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
		{ID: identifiers.New(), Name: "bone-in skin-on chicken breast", Description: "Bone-in, skin-on chicken breast half", PluralName: "bone-in skin-on chicken breast halves", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "bone-in-skin-on-chicken-breast", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
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
		{ID: identifiers.New(), Name: "ribeye steak", Description: "Bone-in or boneless ribeye steak, thick-cut", PluralName: "ribeye steaks", StorageInstructions: "Keep refrigerated at or below 40°F, use within 3-5 days", Slug: "ribeye-steak", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "vegetable oil", Description: "Neutral-flavored vegetable or canola oil", PluralName: "vegetable oil", StorageInstructions: "Store in a cool, dark place", Slug: "vegetable-oil", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "shallot", Description: "Fresh shallots", PluralName: "shallots", StorageInstructions: "Store in a cool, dry, well-ventilated place", Slug: "shallot", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "water", Description: "Water", PluralName: "water", StorageInstructions: "Store at room temperature", Slug: "water", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "whole chicken", Description: "A whole chicken, about 4-5 pounds, with giblets removed and wing tips trimmed", PluralName: "whole chickens", StorageInstructions: "Keep refrigerated at or below 40°F", Slug: "whole-chicken", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "baking powder", Description: "Double-acting baking powder for leavening and crisping", PluralName: "baking powder", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "baking-powder", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		// Burger recipe ingredients
		{ID: identifiers.New(), Name: "beef sirloin", Description: "Lean beef from the sirloin section", PluralName: "beef sirloin", StorageInstructions: "Keep refrigerated at or below 40°F, use within 3-5 days", Slug: "beef-sirloin", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "beef brisket", Description: "Beef cut from the breast section, rich in connective tissue and fat", PluralName: "beef brisket", StorageInstructions: "Keep refrigerated at or below 40°F, use within 3-5 days", Slug: "beef-brisket", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "oxtail", Description: "Beef tail with bone, fat, and meat", PluralName: "oxtails", StorageInstructions: "Keep refrigerated at or below 40°F, use within 3-5 days", Slug: "oxtail", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "American cheese", Description: "Processed American cheese slices", PluralName: "American cheese slices", StorageInstructions: "Keep refrigerated", Slug: "american-cheese", ContainsShellfish: false, ContainsDairy: true, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "burger bun", Description: "Soft white burger bun", PluralName: "burger buns", StorageInstructions: "Store at room temperature in a sealed bag", Slug: "burger-bun", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: true, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "pickle", Description: "Pickled cucumber slices or chips", PluralName: "pickles", StorageInstructions: "Keep refrigerated after opening", Slug: "pickle", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		// Caesar roasted broccoli recipe ingredients
		{ID: identifiers.New(), Name: "anchovy paste", Description: "Concentrated anchovy paste for seasoning", PluralName: "anchovy paste", StorageInstructions: "Keep refrigerated after opening", Slug: "anchovy-paste", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
		// Haricots verts amandine recipe ingredients
		{ID: identifiers.New(), Name: "slivered almonds", Description: "Blanched almonds sliced into thin slivers", PluralName: "slivered almonds", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "slivered-almonds", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: true, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "breadcrumbs", Description: "Plain dry breadcrumbs", PluralName: "breadcrumbs", StorageInstructions: "Store in a cool, dry place in an airtight container", Slug: "breadcrumbs", ContainsShellfish: false, ContainsDairy: false, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: true, ContainsSoy: false, AnimalDerived: false, RestrictToPreparations: false},
		{ID: identifiers.New(), Name: "salted butter", Description: "Salted butter", PluralName: "salted butter", StorageInstructions: "Keep refrigerated, can be kept at room temperature for short periods", Slug: "salted-butter", ContainsShellfish: false, ContainsDairy: true, ContainsPeanut: false, ContainsTreeNut: false, ContainsEgg: false, ContainsWheat: false, ContainsSoy: false, AnimalDerived: true, RestrictToPreparations: false},
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
		{"meat pounder", "A heavy tool for pounding meat to even thickness", "meat pounders", "meat-pounder", "meat pounder"},
		{"rolling pin", "A cylindrical tool for rolling and flattening", "rolling pins", "rolling-pin", "rolling pin"},
		{"brush", "A brush for applying oil or sauces", "brushes", "brush", "brush"},
		{"sous vide cooker", "An immersion circulator for precision temperature cooking", "sous vide cookers", "sous-vide-cooker", "sous vide cooker"},
		{"spatula", "A flexible metal spatula for pressing food in pan", "spatulas", "spatula", "spatula"},
		// Additional instruments for roast chicken recipe
		{"butcher's twine", "Kitchen string for trussing meat and poultry", "butcher's twine", "butchers-twine", "butcher's twine"},
		{"bare hands", "Using clean bare hands to handle or apply ingredients", "bare hands", "bare-hands", "bare hands"},
		// Burger recipe instruments
		{"meat grinder", "A grinder for processing meat, with feed shaft, blade, and die", "meat grinders", "meat-grinder", "meat grinder"},
		{"wide spatula", "A wide, flexible spatula for flipping and pressing food", "wide spatulas", "wide-spatula", "wide spatula"},
		// Rice recipe instruments
		{"fork", "A standard eating fork for fluffing rice and other tasks", "forks", "fork", "fork"},
		{"wooden spoon", "A wooden spoon for stirring", "wooden spoons", "wooden-spoon", "wooden spoon"},
		// Mashed potatoes recipe instruments
		{"vegetable peeler", "A hand-held tool for peeling vegetables", "vegetable peelers", "vegetable-peeler", "vegetable peeler"},
		{"potato ricer", "A kitchen tool that processes potatoes by forcing them through small holes", "potato ricers", "potato-ricer", "potato ricer"},
		{"rubber spatula", "A flexible rubber spatula for folding and scraping", "rubber spatulas", "rubber-spatula", "rubber spatula"},
		// Caesar roasted broccoli recipe instruments
		{"aluminum foil", "Aluminum foil for lining pans and wrapping food", "aluminum foil", "aluminum-foil", "aluminum foil"},
		{"microplane", "A fine grater for zesting citrus and grating hard cheeses", "microplanes", "microplane", "microplane"},
		// Haricots verts amandine recipe instruments
		{"wire mesh spider", "A wide shallow wire-mesh strainer on a long handle for scooping food from hot liquids", "wire mesh spiders", "wire-mesh-spider", "wire mesh spider"},
		{"kitchen towels", "Absorbent towels for drying ingredients", "kitchen towels", "kitchen-towels", "kitchen towels"},
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
		{"liter", "Metric unit of volume equal to 1000 milliliters", "liters", "liter", true, true},
		{"cup", "A volumetric measurement equal to 240 milliliters", "cups", "cup", true, false},
		{"sprig", "A small stem with leaves, typically herbs", "sprigs", "sprig", false, false},
		{"tablespoon", "A volumetric measurement equal to 15 milliliters", "tablespoons", "tablespoon", true, false},
		{"teaspoon", "A volumetric measurement equal to 5 milliliters", "teaspoons", "teaspoon", true, false},
		{"ounce", "Imperial unit of weight equal to approximately 28 grams", "ounces", "ounce", false, false},
		{"slice", "A thin, flat piece cut from something", "slices", "slice", false, false},
		{"pinch", "A small amount picked up between thumb and forefinger", "pinches", "pinch", false, false},
		{"pound", "Imperial unit of weight equal to approximately 454 grams", "pounds", "pound", false, false},
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
	}

	// Create additional ingredient states for recipe completion conditions
	dryState, err := repo.CreateValidIngredientState(ctx, &mealplanning.ValidIngredientStateDatabaseCreationInput{
		ID:            identifiers.New(),
		Name:          "dry",
		Description:   "Ingredient has been dried, with moisture removed",
		AttributeType: mealplanning.ValidIngredientStateAttributeTypeTexture,
		PastTense:     "dried",
		Slug:          "dry",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create dry ingredient state: %w", err)
	}
	enums.IngredientStates["dry"] = dryState

	smokingState, err := repo.CreateValidIngredientState(ctx, &mealplanning.ValidIngredientStateDatabaseCreationInput{
		ID:            identifiers.New(),
		Name:          "smoking",
		Description:   "Oil or fat is hot enough to begin smoking",
		AttributeType: mealplanning.ValidIngredientStateAttributeTypeAppearance,
		PastTense:     "smoking",
		Slug:          "smoking",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create smoking ingredient state: %w", err)
	}
	enums.IngredientStates["smoking"] = smokingState

	atTemperatureState, err := repo.CreateValidIngredientState(ctx, &mealplanning.ValidIngredientStateDatabaseCreationInput{
		ID:            identifiers.New(),
		Name:          "at temperature",
		Description:   "Ingredient has reached the target internal temperature",
		AttributeType: mealplanning.ValidIngredientStateAttributeTypeTemperature,
		PastTense:     "at temperature",
		Slug:          "at-temperature",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create at temperature ingredient state: %w", err)
	}
	enums.IngredientStates["at temperature"] = atTemperatureState

	brownedState, err := repo.CreateValidIngredientState(ctx, &mealplanning.ValidIngredientStateDatabaseCreationInput{
		ID:            identifiers.New(),
		Name:          "browned",
		Description:   "Ingredient has developed a golden-brown color from cooking",
		AttributeType: mealplanning.ValidIngredientStateAttributeTypeAppearance,
		PastTense:     "browned",
		Slug:          "browned",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create browned ingredient state: %w", err)
	}
	enums.IngredientStates["browned"] = brownedState

	toastedState, err := repo.CreateValidIngredientState(ctx, &mealplanning.ValidIngredientStateDatabaseCreationInput{
		ID:            identifiers.New(),
		Name:          "toasted",
		Description:   "Ingredient has been dry-roasted until deeply browned and fragrant",
		AttributeType: mealplanning.ValidIngredientStateAttributeTypeAppearance,
		PastTense:     "toasted",
		Slug:          "toasted",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create toasted ingredient state: %w", err)
	}
	enums.IngredientStates["toasted"] = toastedState

	clearState, err := repo.CreateValidIngredientState(ctx, &mealplanning.ValidIngredientStateDatabaseCreationInput{
		ID:            identifiers.New(),
		Name:          "clear",
		Description:   "Liquid is transparent or translucent, free of cloudiness",
		AttributeType: mealplanning.ValidIngredientStateAttributeTypeAppearance,
		PastTense:     "cleared",
		Slug:          "clear",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create clear ingredient state: %w", err)
	}
	enums.IngredientStates["clear"] = clearState

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

	// Create additional vessels for grilled chicken recipe
	wireRack, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "wire rack",
		Description:                    "A wire rack for air circulation",
		PluralName:                     "wire racks",
		Slug:                           "wire-rack",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             300,
		LengthInMillimeters:            400,
		HeightInMillimeters:            30,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create wire rack vessel: %w", err)
	}
	enums.Vessels["wire rack"] = wireRack

	grill, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "grill",
		Description:                    "A grill for cooking over direct heat",
		PluralName:                     "grills",
		Slug:                           "grill",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             400,
		LengthInMillimeters:            600,
		HeightInMillimeters:            200,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create grill vessel: %w", err)
	}
	enums.Vessels["grill"] = grill

	plasticBag, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "plastic bag",
		Description:                    "A resealable plastic bag for containing food",
		PluralName:                     "plastic bags",
		Slug:                           "plastic-bag",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             200,
		LengthInMillimeters:            300,
		HeightInMillimeters:            5,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create plastic bag vessel: %w", err)
	}
	enums.Vessels["plastic bag"] = plasticBag

	// Create additional vessels for sous vide chicken recipe
	vacuumBag, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "vacuum bag",
		Description:                    "A vacuum-sealable bag for sous vide cooking",
		PluralName:                     "vacuum bags",
		Slug:                           "vacuum-bag",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             200,
		LengthInMillimeters:            300,
		HeightInMillimeters:            5,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create vacuum bag vessel: %w", err)
	}
	enums.Vessels["vacuum bag"] = vacuumBag

	waterBath, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "water bath",
		Description:                    "A container of water for sous vide cooking",
		PluralName:                     "water baths",
		Slug:                           "water-bath",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             300,
		LengthInMillimeters:            400,
		HeightInMillimeters:            200,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create water bath vessel: %w", err)
	}
	enums.Vessels["water bath"] = waterBath

	// Create additional vessels for roast chicken recipe
	smallBowl, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "small bowl",
		Description:                    "A small mixing bowl",
		PluralName:                     "small bowls",
		Slug:                           "small-bowl",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             150,
		LengthInMillimeters:            150,
		HeightInMillimeters:            80,
		Shape:                          mealplanning.VesselShapeHemisphere,
		UsableForStorage:               true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create small bowl vessel: %w", err)
	}
	enums.Vessels["small bowl"] = smallBowl

	stainlessSteelSkillet, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "stainless steel skillet",
		Description:                    "A 10 to 12-inch stainless steel frying pan",
		PluralName:                     "stainless steel skillets",
		Slug:                           "stainless-steel-skillet",
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
		return nil, fmt.Errorf("failed to create stainless steel skillet vessel: %w", err)
	}
	enums.Vessels["stainless steel skillet"] = stainlessSteelSkillet

	carvingBoard, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "carving board",
		Description:                    "A large cutting board for carving meat",
		PluralName:                     "carving boards",
		Slug:                           "carving-board",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             400,
		LengthInMillimeters:            500,
		HeightInMillimeters:            30,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create carving board vessel: %w", err)
	}
	enums.Vessels["carving board"] = carvingBoard

	// Create additional vessels for burger recipe
	largeBowl, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "large bowl",
		Description:                    "A large mixing bowl",
		PluralName:                     "large bowls",
		Slug:                           "large-bowl",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             300,
		LengthInMillimeters:            300,
		HeightInMillimeters:            150,
		Shape:                          mealplanning.VesselShapeHemisphere,
		UsableForStorage:               true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create large bowl vessel: %w", err)
	}
	enums.Vessels["large bowl"] = largeBowl

	sautePan, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "sauté pan",
		Description:                    "A heavy-bottomed sauté pan or skillet",
		PluralName:                     "sauté pans",
		Slug:                           "saute-pan",
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
		return nil, fmt.Errorf("failed to create sauté pan vessel: %w", err)
	}
	enums.Vessels["sauté pan"] = sautePan

	// Create saucepan for rice recipe
	saucepan, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "saucepan",
		Description:                    "A 2-quart saucepan with lid",
		PluralName:                     "saucepans",
		Slug:                           "saucepan",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             180,
		LengthInMillimeters:            180,
		HeightInMillimeters:            100,
		Shape:                          mealplanning.VesselShapeCylinder,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create saucepan vessel: %w", err)
	}
	enums.Vessels["saucepan"] = saucepan

	// Create fine-mesh strainer for rice recipe
	fineMeshStrainer, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "fine-mesh strainer",
		Description:                    "A strainer with fine mesh for draining rice and other small grains",
		PluralName:                     "fine-mesh strainers",
		Slug:                           "fine-mesh-strainer",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             150,
		LengthInMillimeters:            150,
		HeightInMillimeters:            80,
		Shape:                          mealplanning.VesselShapeHemisphere,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create fine-mesh strainer vessel: %w", err)
	}
	enums.Vessels["fine-mesh strainer"] = fineMeshStrainer

	// Create colander for mashed potatoes recipe
	colander, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "colander",
		Description:                    "A bowl-shaped strainer with holes for draining food",
		PluralName:                     "colanders",
		Slug:                           "colander",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             250,
		LengthInMillimeters:            250,
		HeightInMillimeters:            120,
		Shape:                          mealplanning.VesselShapeHemisphere,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create colander vessel: %w", err)
	}
	enums.Vessels["colander"] = colander

	// Caesar roasted broccoli recipe vessels
	smallNonstickSkillet, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "small nonstick skillet",
		Description:                    "A small nonstick frying pan for delicate cooking tasks",
		PluralName:                     "small nonstick skillets",
		Slug:                           "small-nonstick-skillet",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             200,
		LengthInMillimeters:            200,
		HeightInMillimeters:            40,
		Shape:                          mealplanning.VesselShapeCylinder,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create small nonstick skillet vessel: %w", err)
	}
	enums.Vessels["small nonstick skillet"] = smallNonstickSkillet

	servingPlatter, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "serving platter",
		Description:                    "A large flat platter for serving food",
		PluralName:                     "serving platters",
		Slug:                           "serving-platter",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             350,
		LengthInMillimeters:            450,
		HeightInMillimeters:            25,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create serving platter vessel: %w", err)
	}
	enums.Vessels["serving platter"] = servingPlatter

	oven, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "oven",
		Description:                    "A kitchen oven for baking and roasting",
		PluralName:                     "ovens",
		Slug:                           "oven",
		IncludeInGeneratedInstructions: false,
		DisplayInSummaryLists:          false,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             600,
		LengthInMillimeters:            600,
		HeightInMillimeters:            400,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create oven vessel: %w", err)
	}
	enums.Vessels["oven"] = oven

	freezer, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "freezer",
		Description:                    "A freezer compartment for chilling items",
		PluralName:                     "freezers",
		Slug:                           "freezer",
		IncludeInGeneratedInstructions: false,
		DisplayInSummaryLists:          false,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             500,
		LengthInMillimeters:            500,
		HeightInMillimeters:            400,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create freezer vessel: %w", err)
	}
	enums.Vessels["freezer"] = freezer

	// Haricots verts amandine recipe vessels
	mediumSkillet, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "medium skillet",
		Description:                    "A medium-sized frying pan, typically 10 inches in diameter",
		PluralName:                     "medium skillets",
		Slug:                           "medium-skillet",
		IncludeInGeneratedInstructions: true,
		DisplayInSummaryLists:          true,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             250,
		LengthInMillimeters:            250,
		HeightInMillimeters:            50,
		Shape:                          mealplanning.VesselShapeCylinder,
		UsableForStorage:               false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create medium skillet vessel: %w", err)
	}
	enums.Vessels["medium skillet"] = mediumSkillet

	refrigerator, err := repo.CreateValidVessel(ctx, &mealplanning.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           "refrigerator",
		Description:                    "A refrigerator for keeping items cold",
		PluralName:                     "refrigerators",
		Slug:                           "refrigerator",
		IncludeInGeneratedInstructions: false,
		DisplayInSummaryLists:          false,
		CapacityUnitID:                 &firstValidMeasurementUnitGram.ID,
		WidthInMillimeters:             600,
		LengthInMillimeters:            600,
		HeightInMillimeters:            1500,
		Shape:                          mealplanning.VesselShapeRectangle,
		UsableForStorage:               true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create refrigerator vessel: %w", err)
	}
	enums.Vessels["refrigerator"] = refrigerator

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
		{"dry", "Remove moisture using paper towels or cloth", "dried", "dry", false, false},
		{"heat", "Raise the temperature of an ingredient or vessel", "heated", "heat", true, true},
		{"pan-sear", "Cook in a hot pan with oil to develop a crust", "pan-seared", "pan-sear", false, true},
		{"baste", "Spoon hot fat or liquid over food while cooking", "basted", "baste", true, true},
		{"rest", "Allow food to sit after cooking to redistribute juices", "rested", "rest", false, true},
		// Additional preparations for sous vide chicken recipe
		{"pound", "Flatten meat to an even thickness using a mallet or heavy object", "pounded", "pound", false, false},
		{"wet-brine", "Soak in a saltwater solution to season and tenderize", "wet-brined", "wet-brine", false, true},
		{"dry-brine", "Salt and refrigerate uncovered to season and dry the surface", "dry-brined", "dry-brine", false, true},
		{"bag", "Place ingredients in a bag for cooking or storage", "bagged", "bag", false, false},
		{"sous-vide", "Cook in a temperature-controlled water bath while sealed in a bag", "sous-vided", "sous-vide", true, true},
		// Additional preparations for roast chicken recipe
		{"truss", "Tie meat or poultry with string to maintain shape during cooking", "trussed", "truss", false, false},
		{"rub", "Apply a seasoning or oil by rubbing onto the surface of food", "rubbed", "rub", false, false},
		// Burger recipe preparations
		{"chill", "Place in freezer or refrigerator to reduce temperature", "chilled", "chill", false, true},
		{"trim", "Remove unwanted parts such as gristle, silverskin, or fat", "trimmed", "trim", false, false},
		{"cube", "Cut into cube shapes, typically 1-inch or similar", "cubed", "cube", false, false},
		{"grind", "Process through a meat grinder to break down into smaller pieces", "ground", "grind", false, false},
		{"form", "Shape ingredients into a specific form such as patties or balls", "formed", "form", false, false},
		{"line", "Cover a surface with parchment, foil, or similar material", "lined", "line", false, false},
		{"flip", "Turn food over to expose the other side", "flipped", "flip", false, false},
		{"top", "Place ingredients on top of other ingredients", "topped", "top", false, false},
		{"assemble", "Put together components of a dish", "assembled", "assemble", false, false},
		{"toast", "Lightly brown with dry heat", "toasted", "toast", false, true},
		{"refrigerate", "Store in the refrigerator to keep cold", "refrigerated", "refrigerate", false, true},
		{"debone", "Remove bones from meat or poultry", "deboned", "debone", false, false},
		// Smash burger recipe preparations
		{"smash", "Press down firmly to flatten", "smashed", "smash", false, false},
		{"divide", "Separate into portions", "divided", "divide", false, false},
		// Rice recipe preparations
		{"rinse", "Wash with water to remove starch or impurities", "rinsed", "rinse", false, false},
		{"drain", "Remove liquid using a strainer or colander", "drained", "drain", false, false},
		{"cover", "Place a lid, wrap, or foil over a vessel", "covered", "cover", false, false},
		{"fluff", "Gently separate grains or fibers with a fork", "fluffed", "fluff", false, false},
		// Mashed potatoes recipe preparations
		{"peel", "Remove the outer skin from vegetables or fruits", "peeled", "peel", false, false},
		{"fold", "Gently combine ingredients by lifting from the bottom and folding over", "folded", "fold", false, false},
		{"rice", "Press cooked food through a ricer or food mill to create a smooth texture", "riced", "rice", false, false},
		{"submerge", "Cover completely with liquid", "submerged", "submerge", false, false},
		// Caesar roasted broccoli recipe preparations
		{"melt", "Heat a solid ingredient until it becomes liquid", "melted", "melt", false, false},
		{"preheat", "Heat an oven, vessel, or equipment to a specified temperature before use", "preheated", "preheat", true, true},
		{"toss", "Mix ingredients by lifting and turning", "tossed", "toss", false, false},
		{"zest", "Grate the outer peel of citrus fruit", "zested", "zest", false, false},
		{"transfer", "Move ingredients from one vessel to another", "transferred", "transfer", false, false},
		// Haricots verts amandine recipe preparations
		{"blanch", "Cook briefly in boiling water then shock in ice water", "blanched", "blanch", false, true},
		{"emulsify", "Combine fat and water-based liquids into a smooth, glossy sauce by rapid stirring or shaking", "emulsified", "emulsify", false, false},
		{"shock", "Immediately transfer hot food to ice water to stop cooking", "shocked", "shock", false, false},
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

	// Get preparations for new steps
	dryPrep := enums.Preparations["dry"]
	heatPrep := enums.Preparations["heat"]

	// === DRY PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(dryPrep, ribeye); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links (already created for ribeye)

	// Preparation-Instrument links
	if err := createVPI(dryPrep, paperTowels); err != nil {
		return err
	}

	// === HEAT PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(heatPrep, vegetableOil); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links (already created for vegetableOil)

	// Preparation-Vessel links
	if err := createVPV(heatPrep, castIronSkillet); err != nil {
		return err
	}

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
	bareHands := enums.Instruments["bare hands"]
	if err := createVPI(seasonPrep, bareHands); err != nil {
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

	// Preparation-Vessel links
	if err := createVPV(bastePrep, castIronSkillet); err != nil {
		return err
	}

	// === REST PREPARATION ===
	// Preparation-Instrument links
	if err := createVPI(restPrep, tongs); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(restPrep, sheetPan); err != nil {
		return err
	}
	if err := createVPV(restPrep, servingPlate); err != nil {
		return err
	}

	// === CHICKEN RECIPE BRIDGE ENTRIES ===
	// Get preparations for chicken recipe
	poundPrep := enums.Preparations["pound"]
	wetBrinePrep := enums.Preparations["wet-brine"]
	dryBrinePrep := enums.Preparations["dry-brine"]
	grillPrep := enums.Preparations["grill"]

	// Get ingredients for chicken recipe
	chickenBreast := enums.Ingredients["chicken breast"]
	water := enums.Ingredients["water"]
	oliveOil := enums.Ingredients["olive oil"]
	sugar := enums.Ingredients["sugar"]

	// Get instruments for chicken recipe
	meatPounder := enums.Instruments["meat pounder"]
	rollingPin := enums.Instruments["rolling pin"]
	brush := enums.Instruments["brush"]

	// Get vessels for chicken recipe
	wireRack := enums.Vessels["wire rack"]
	grillVessel := enums.Vessels["grill"]
	plasticBag := enums.Vessels["plastic bag"]

	// Get measurement units
	literMeasurement := enums.MeasurementUnits["liter"]

	// === SEASON PREPARATION (for chicken breast) ===
	// Ingredient-Preparation links
	if err := createVIP(seasonPrep, chickenBreast); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links
	if err := createVIMU(chickenBreast, unitMeasurement); err != nil {
		return err
	}

	// === POUND PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(poundPrep, chickenBreast); err != nil {
		return err
	}

	// Preparation-Instrument links
	if err := createVPI(poundPrep, meatPounder); err != nil {
		return err
	}
	if err := createVPI(poundPrep, rollingPin); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(poundPrep, plasticBag); err != nil {
		return err
	}

	// === WET-BRINE PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(wetBrinePrep, chickenBreast); err != nil {
		return err
	}
	if err := createVIP(wetBrinePrep, salt); err != nil {
		return err
	}
	if err := createVIP(wetBrinePrep, sugar); err != nil {
		return err
	}
	if err := createVIP(wetBrinePrep, water); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links
	if err := createVIMU(water, literMeasurement); err != nil {
		return err
	}
	if err := createVIMU(sugar, gramMeasurement); err != nil {
		return err
	}

	// === DRY-BRINE PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(dryBrinePrep, chickenBreast); err != nil {
		return err
	}
	if err := createVIP(dryBrinePrep, salt); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(dryBrinePrep, wireRack); err != nil {
		return err
	}
	if err := createVPV(dryBrinePrep, sheetPan); err != nil {
		return err
	}

	// === GRILL PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(grillPrep, chickenBreast); err != nil {
		return err
	}
	if err := createVIP(grillPrep, oliveOil); err != nil {
		return err
	}
	if err := createVIP(grillPrep, salt); err != nil {
		return err
	}
	if err := createVIP(grillPrep, blackPepper); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links
	if err := createVIMU(oliveOil, milliliterMeasurement); err != nil {
		return err
	}

	// Preparation-Instrument links
	if err := createVPI(grillPrep, brush); err != nil {
		return err
	}
	if err := createVPI(grillPrep, thermometer); err != nil {
		return err
	}
	if err := createVPI(grillPrep, tongs); err != nil {
		return err
	}
	if err := createVPI(grillPrep, paperTowels); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(grillPrep, grillVessel); err != nil {
		return err
	}

	// === SOUS VIDE CHICKEN RECIPE BRIDGE ENTRIES ===
	// Get preparations for sous vide recipe
	sousVidePrep := enums.Preparations["sous-vide"]
	bagPrep := enums.Preparations["bag"]
	panSearPrep = enums.Preparations["pan-sear"]

	// Get ingredients for sous vide recipe
	boneInSkinOnChickenBreast := enums.Ingredients["bone-in skin-on chicken breast"]
	thyme = enums.Ingredients["thyme"]
	rosemary = enums.Ingredients["rosemary"]

	// Get instruments for sous vide recipe
	sousVideCooker := enums.Instruments["sous vide cooker"]
	spatula := enums.Instruments["spatula"]

	// Get vessels for sous vide recipe
	vacuumBag := enums.Vessels["vacuum bag"]
	waterBath := enums.Vessels["water bath"]

	// === HEAT PREPARATION (for water bath) ===
	// Preparation-Vessel links
	if err := createVPV(heatPrep, waterBath); err != nil {
		return err
	}

	// === BAG PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(bagPrep, boneInSkinOnChickenBreast); err != nil {
		return err
	}
	if err := createVIP(bagPrep, thyme); err != nil {
		return err
	}
	if err := createVIP(bagPrep, rosemary); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links (already created for bone-in skin-on chicken)
	if err := createVIMU(thyme, sprigMeasurement); err != nil {
		return err
	}
	if err := createVIMU(rosemary, sprigMeasurement); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(bagPrep, plasticBag); err != nil {
		return err
	}
	if err := createVPV(bagPrep, vacuumBag); err != nil {
		return err
	}

	// === SOUS-VIDE PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(sousVidePrep, boneInSkinOnChickenBreast); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links
	if err := createVIMU(boneInSkinOnChickenBreast, unitMeasurement); err != nil {
		return err
	}

	// Preparation-Instrument links
	if err := createVPI(sousVidePrep, sousVideCooker); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(sousVidePrep, waterBath); err != nil {
		return err
	}

	// === PAN-SEAR PREPARATION (for finishing) ===
	// Ingredient-Preparation links for bone-in skin-on chicken
	if err := createVIP(panSearPrep, boneInSkinOnChickenBreast); err != nil {
		return err
	}
	// Already has oil bridges, but need to add paper towels and spatula
	if err := createVPI(panSearPrep, paperTowels); err != nil {
		return err
	}
	if err := createVPI(panSearPrep, spatula); err != nil {
		return err
	}
	// Need cast iron skillet vessel (may already exist, but ensure it)
	if err := createVPV(panSearPrep, castIronSkillet); err != nil {
		return err
	}

	// === GRILL PREPARATION (for finishing) ===
	// Ingredient-Preparation links for bone-in skin-on chicken
	if err := createVIP(grillPrep, boneInSkinOnChickenBreast); err != nil {
		return err
	}

	// === ROAST CHICKEN RECIPE BRIDGE ENTRIES ===
	// Get preparations for roast chicken recipe
	mixPrep := enums.Preparations["mix"]
	seasonPrep = enums.Preparations["season"]
	trussPrep := enums.Preparations["truss"]
	dryBrinePrep = enums.Preparations["dry-brine"]
	heatPrep = enums.Preparations["heat"]
	rubPrep := enums.Preparations["rub"]
	panSearPrep = enums.Preparations["pan-sear"]
	roastPrep := enums.Preparations["roast"]
	restPrep = enums.Preparations["rest"]

	// Get ingredients for roast chicken recipe
	wholeChicken := enums.Ingredients["whole chicken"]
	bakingPowder := enums.Ingredients["baking powder"]
	vegetableOil = enums.Ingredients["vegetable oil"]

	// Get instruments for roast chicken recipe
	butchersTwine := enums.Instruments["butcher's twine"]

	// Get vessels for roast chicken recipe
	smallBowl := enums.Vessels["small bowl"]
	stainlessSteelSkillet := enums.Vessels["stainless steel skillet"]
	carvingBoard := enums.Vessels["carving board"]
	bakingSheet := enums.Vessels["baking sheet"]
	// wireRack already defined above in chicken recipe section

	// Get measurement units for roast chicken
	tablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	teaspoonMeasurement := enums.MeasurementUnits["teaspoon"]

	// === MIX PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(mixPrep, salt); err != nil {
		return err
	}
	if err := createVIP(mixPrep, blackPepper); err != nil {
		return err
	}
	if err := createVIP(mixPrep, bakingPowder); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links
	if err := createVIMU(bakingPowder, teaspoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(salt, tablespoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(blackPepper, teaspoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(vegetableOil, tablespoonMeasurement); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(mixPrep, smallBowl); err != nil {
		return err
	}

	// === SEASON PREPARATION (for whole chicken) ===
	// Ingredient-Preparation links
	if err := createVIP(seasonPrep, wholeChicken); err != nil {
		return err
	}

	// Ingredient-MeasurementUnit links
	if err := createVIMU(wholeChicken, unitMeasurement); err != nil {
		return err
	}

	// === TRUSS PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(trussPrep, wholeChicken); err != nil {
		return err
	}

	// Preparation-Instrument links
	if err := createVPI(trussPrep, butchersTwine); err != nil {
		return err
	}

	// === DRY-BRINE PREPARATION (for whole chicken) ===
	// Ingredient-Preparation links
	if err := createVIP(dryBrinePrep, wholeChicken); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(dryBrinePrep, wireRack); err != nil {
		return err
	}
	if err := createVPV(dryBrinePrep, bakingSheet); err != nil {
		return err
	}

	// === HEAT PREPARATION (for stainless steel skillet) ===
	// Preparation-Vessel links
	if err := createVPV(heatPrep, stainlessSteelSkillet); err != nil {
		return err
	}

	// === RUB PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(rubPrep, wholeChicken); err != nil {
		return err
	}
	if err := createVIP(rubPrep, vegetableOil); err != nil {
		return err
	}

	// Preparation-Instrument links
	if err := createVPI(rubPrep, bareHands); err != nil {
		return err
	}

	// === PAN-SEAR PREPARATION (for whole chicken) ===
	// Ingredient-Preparation links
	if err := createVIP(panSearPrep, wholeChicken); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(panSearPrep, stainlessSteelSkillet); err != nil {
		return err
	}

	// === ROAST PREPARATION ===
	// Ingredient-Preparation links
	if err := createVIP(roastPrep, wholeChicken); err != nil {
		return err
	}

	// Preparation-Instrument links
	if err := createVPI(roastPrep, thermometer); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(roastPrep, stainlessSteelSkillet); err != nil {
		return err
	}

	// === REST PREPARATION (for whole chicken) ===
	// Ingredient-Preparation links
	if err := createVIP(restPrep, wholeChicken); err != nil {
		return err
	}

	// Preparation-Vessel links
	if err := createVPV(restPrep, carvingBoard); err != nil {
		return err
	}

	// === SOUS VIDE PORK CHOPS RECIPE BRIDGE ENTRIES ===
	// Get preparations for pork chops recipe
	porkSeasonPrep := enums.Preparations["season"]
	porkDryPrep := enums.Preparations["dry"]
	porkBagPrep := enums.Preparations["bag"]
	porkSousVidePrep := enums.Preparations["sous-vide"]
	porkPanSearPrep := enums.Preparations["pan-sear"]
	porkBastePrep := enums.Preparations["baste"]
	porkGrillPrep := enums.Preparations["grill"]
	porkRestPrep := enums.Preparations["rest"]
	porkHeatPrep := enums.Preparations["heat"]

	// Get ingredients for pork chops recipe
	porkChop := enums.Ingredients["pork chop"]
	porkGarlic := enums.Ingredients["garlic"]
	porkVegOil := enums.Ingredients["vegetable oil"]
	porkButter := enums.Ingredients["butter"]

	// Get measurement units for pork chops recipe
	porkUnitMeasurement := enums.MeasurementUnits["unit"]
	porkTablespoonMeasurement := enums.MeasurementUnits["tablespoon"]

	// Get instruments for pork chops recipe
	porkPaperTowels := enums.Instruments["paper towels"]
	porkTongs := enums.Instruments["tongs"]
	porkSpoon := enums.Instruments["spoon"]
	porkSousVideCooker := enums.Instruments["sous vide cooker"]
	porkBareHands := enums.Instruments["bare hands"]

	// Get vessels for pork chops recipe
	porkCastIronSkillet := enums.Vessels["cast iron skillet"]
	porkWaterBath := enums.Vessels["water bath"]
	porkPlasticBag := enums.Vessels["plastic bag"]
	porkVacuumBag := enums.Vessels["vacuum bag"]
	porkGrillVessel := enums.Vessels["grill"]
	porkWireRack := enums.Vessels["wire rack"]
	porkBakingSheet := enums.Vessels["baking sheet"]
	porkServingPlate := enums.Vessels["serving plate"]

	// === PORK CHOP INGREDIENT-PREPARATION LINKS ===
	if err := createVIP(porkSeasonPrep, porkChop); err != nil {
		return err
	}
	if err := createVIP(porkDryPrep, porkChop); err != nil {
		return err
	}
	if err := createVIP(porkBagPrep, porkChop); err != nil {
		return err
	}
	if err := createVIP(porkSousVidePrep, porkChop); err != nil {
		return err
	}
	if err := createVIP(porkPanSearPrep, porkChop); err != nil {
		return err
	}
	if err := createVIP(porkBastePrep, porkChop); err != nil {
		return err
	}
	if err := createVIP(porkGrillPrep, porkChop); err != nil {
		return err
	}
	if err := createVIP(porkRestPrep, porkChop); err != nil {
		return err
	}

	// Garlic for baste step
	if err := createVIP(porkBastePrep, porkGarlic); err != nil {
		return err
	}

	// === PORK CHOP INGREDIENT-MEASUREMENT UNIT LINKS ===
	if err := createVIMU(porkChop, porkUnitMeasurement); err != nil {
		return err
	}
	if err := createVIMU(porkGarlic, porkUnitMeasurement); err != nil {
		return err
	}
	if err := createVIMU(porkButter, porkTablespoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(porkVegOil, porkTablespoonMeasurement); err != nil {
		return err
	}

	// === PORK CHOP PREPARATION-INSTRUMENT LINKS ===
	// Season with bare hands
	if err := createVPI(porkSeasonPrep, porkBareHands); err != nil {
		return err
	}
	// Dry with paper towels
	if err := createVPI(porkDryPrep, porkPaperTowels); err != nil {
		return err
	}
	// Heat with sous vide cooker
	if err := createVPI(porkHeatPrep, porkSousVideCooker); err != nil {
		return err
	}
	// Sous vide with sous vide cooker
	if err := createVPI(porkSousVidePrep, porkSousVideCooker); err != nil {
		return err
	}
	// Pan-sear with tongs
	if err := createVPI(porkPanSearPrep, porkTongs); err != nil {
		return err
	}
	// Baste with tongs and spoon
	if err := createVPI(porkBastePrep, porkTongs); err != nil {
		return err
	}
	if err := createVPI(porkBastePrep, porkSpoon); err != nil {
		return err
	}
	// Grill with tongs and paper towels
	if err := createVPI(porkGrillPrep, porkTongs); err != nil {
		return err
	}
	if err := createVPI(porkGrillPrep, porkPaperTowels); err != nil {
		return err
	}
	// Rest with tongs
	if err := createVPI(porkRestPrep, porkTongs); err != nil {
		return err
	}

	// === PORK CHOP PREPARATION-VESSEL LINKS ===
	// Heat for water bath
	if err := createVPV(porkHeatPrep, porkWaterBath); err != nil {
		return err
	}
	// Bag with plastic bag or vacuum bag
	if err := createVPV(porkBagPrep, porkPlasticBag); err != nil {
		return err
	}
	if err := createVPV(porkBagPrep, porkVacuumBag); err != nil {
		return err
	}
	// Sous vide in water bath
	if err := createVPV(porkSousVidePrep, porkWaterBath); err != nil {
		return err
	}
	// Heat cast iron skillet
	if err := createVPV(porkHeatPrep, porkCastIronSkillet); err != nil {
		return err
	}
	// Pan-sear in cast iron skillet
	if err := createVPV(porkPanSearPrep, porkCastIronSkillet); err != nil {
		return err
	}
	// Baste in cast iron skillet
	if err := createVPV(porkBastePrep, porkCastIronSkillet); err != nil {
		return err
	}
	// Grill on grill
	if err := createVPV(porkGrillPrep, porkGrillVessel); err != nil {
		return err
	}
	// Rest on wire rack set over baking sheet
	if err := createVPV(porkRestPrep, porkWireRack); err != nil {
		return err
	}
	if err := createVPV(porkRestPrep, porkBakingSheet); err != nil {
		return err
	}
	// Also rest on serving plate (for final serve)
	if err := createVPV(porkRestPrep, porkServingPlate); err != nil {
		return err
	}

	// === CHEESEBURGER RECIPE BRIDGE ENTRIES ===
	// Get preparations for burger recipe
	chillPrep := enums.Preparations["chill"]
	trimPrep := enums.Preparations["trim"]
	cubePrep := enums.Preparations["cube"]
	grindPrep := enums.Preparations["grind"]
	formPrep := enums.Preparations["form"]
	linePrep := enums.Preparations["line"]
	burgerSeasonPrep := enums.Preparations["season"]
	flipPrep := enums.Preparations["flip"]
	refrigeratePrep := enums.Preparations["refrigerate"]
	burgerHeatPrep := enums.Preparations["heat"]
	burgerPanSearPrep := enums.Preparations["pan-sear"]
	topPrep := enums.Preparations["top"]
	assemblePrep := enums.Preparations["assemble"]
	toastPrep := enums.Preparations["toast"]
	debonePrep := enums.Preparations["debone"]

	// Get ingredients for burger recipe
	beefSirloin := enums.Ingredients["beef sirloin"]
	beefBrisket := enums.Ingredients["beef brisket"]
	oxtail := enums.Ingredients["oxtail"]
	americanCheese := enums.Ingredients["American cheese"]
	burgerBun := enums.Ingredients["burger bun"]
	pickle := enums.Ingredients["pickle"]
	burgerOnion := enums.Ingredients["onion"]
	burgerSalt := enums.Ingredients["salt"]
	burgerPepper := enums.Ingredients["black pepper"]
	burgerVegOil := enums.Ingredients["vegetable oil"]

	// Get instruments for burger recipe
	meatGrinder := enums.Instruments["meat grinder"]
	wideSpatula := enums.Instruments["wide spatula"]
	burgerBareHands := enums.Instruments["bare hands"]
	knife := enums.Instruments["knife"]

	// Get vessels for burger recipe
	largeBowl := enums.Vessels["large bowl"]
	sautePan := enums.Vessels["sauté pan"]
	freezer := enums.Vessels["freezer"]
	refrigerator := enums.Vessels["refrigerator"]
	burgerBakingSheet := enums.Vessels["baking sheet"]
	burgerCuttingBoard := enums.Vessels["cutting board"]
	burgerServingPlate := enums.Vessels["serving plate"]

	// Get measurement units for burger recipe
	ounceMeasurement := enums.MeasurementUnits["ounce"]
	sliceMeasurement := enums.MeasurementUnits["slice"]
	burgerUnitMeasurement := enums.MeasurementUnits["unit"]
	burgerTeaspoonMeasurement := enums.MeasurementUnits["teaspoon"]

	// === TRIM PREPARATION ===
	if err := createVIP(trimPrep, beefSirloin); err != nil {
		return err
	}
	if err := createVIP(trimPrep, beefBrisket); err != nil {
		return err
	}
	if err := createVIP(trimPrep, oxtail); err != nil {
		return err
	}
	if err := createVIMU(beefSirloin, ounceMeasurement); err != nil {
		return err
	}
	if err := createVIMU(beefBrisket, ounceMeasurement); err != nil {
		return err
	}
	if err := createVIMU(oxtail, ounceMeasurement); err != nil {
		return err
	}
	if err := createVPI(trimPrep, knife); err != nil {
		return err
	}
	if err := createVPV(trimPrep, burgerCuttingBoard); err != nil {
		return err
	}

	// === DEBONE PREPARATION ===
	if err := createVIP(debonePrep, oxtail); err != nil {
		return err
	}
	if err := createVPI(debonePrep, knife); err != nil {
		return err
	}
	if err := createVPV(debonePrep, burgerCuttingBoard); err != nil {
		return err
	}

	// === CUBE PREPARATION ===
	if err := createVIP(cubePrep, beefSirloin); err != nil {
		return err
	}
	if err := createVIP(cubePrep, beefBrisket); err != nil {
		return err
	}
	if err := createVIP(cubePrep, oxtail); err != nil {
		return err
	}
	if err := createVPI(cubePrep, knife); err != nil {
		return err
	}
	if err := createVPV(cubePrep, burgerCuttingBoard); err != nil {
		return err
	}

	// === CHILL PREPARATION ===
	if err := createVIP(chillPrep, beefSirloin); err != nil {
		return err
	}
	if err := createVIP(chillPrep, beefBrisket); err != nil {
		return err
	}
	if err := createVIP(chillPrep, oxtail); err != nil {
		return err
	}
	if err := createVPV(chillPrep, freezer); err != nil {
		return err
	}
	if err := createVPV(chillPrep, burgerBakingSheet); err != nil {
		return err
	}
	// Chill preparation for meat grinder instrument
	if err := createVPI(chillPrep, meatGrinder); err != nil {
		return err
	}

	// === MIX/COMBINE PREPARATION (using existing mix) ===
	burgerMixPrep := enums.Preparations["mix"]
	if err := createVIP(burgerMixPrep, beefSirloin); err != nil {
		return err
	}
	if err := createVIP(burgerMixPrep, beefBrisket); err != nil {
		return err
	}
	if err := createVIP(burgerMixPrep, oxtail); err != nil {
		return err
	}
	if err := createVPV(burgerMixPrep, largeBowl); err != nil {
		return err
	}
	if err := createVPI(burgerMixPrep, burgerBareHands); err != nil {
		return err
	}

	// === LINE PREPARATION ===
	if err := createVPV(linePrep, burgerBakingSheet); err != nil {
		return err
	}

	// === GRIND PREPARATION ===
	if err := createVIP(grindPrep, beefSirloin); err != nil {
		return err
	}
	if err := createVIP(grindPrep, beefBrisket); err != nil {
		return err
	}
	if err := createVIP(grindPrep, oxtail); err != nil {
		return err
	}
	if err := createVPI(grindPrep, meatGrinder); err != nil {
		return err
	}
	if err := createVPV(grindPrep, burgerBakingSheet); err != nil {
		return err
	}

	// === FORM PREPARATION ===
	if err := createVPI(formPrep, burgerBareHands); err != nil {
		return err
	}
	if err := createVPV(formPrep, burgerBakingSheet); err != nil {
		return err
	}

	// === SEASON PREPARATION for burger ingredients ===
	if err := createVIP(burgerSeasonPrep, burgerSalt); err != nil {
		return err
	}
	if err := createVIP(burgerSeasonPrep, burgerPepper); err != nil {
		return err
	}
	if err := createVPV(burgerSeasonPrep, burgerBakingSheet); err != nil {
		return err
	}

	// === FLIP PREPARATION ===
	if err := createVPI(flipPrep, wideSpatula); err != nil {
		return err
	}
	if err := createVPV(flipPrep, burgerBakingSheet); err != nil {
		return err
	}

	// === REFRIGERATE PREPARATION ===
	if err := createVPV(refrigeratePrep, refrigerator); err != nil {
		return err
	}
	if err := createVPV(refrigeratePrep, burgerBakingSheet); err != nil {
		return err
	}

	// === HEAT PREPARATION for sauté pan ===
	if err := createVIP(burgerHeatPrep, burgerVegOil); err != nil {
		return err
	}
	if err := createVIMU(burgerVegOil, burgerTeaspoonMeasurement); err != nil {
		return err
	}
	if err := createVPV(burgerHeatPrep, sautePan); err != nil {
		return err
	}

	// === PAN-SEAR PREPARATION ===
	if err := createVPI(burgerPanSearPrep, wideSpatula); err != nil {
		return err
	}
	if err := createVPV(burgerPanSearPrep, sautePan); err != nil {
		return err
	}

	// === TOP PREPARATION (adding cheese) ===
	if err := createVIP(topPrep, americanCheese); err != nil {
		return err
	}
	if err := createVIMU(americanCheese, sliceMeasurement); err != nil {
		return err
	}
	if err := createVPV(topPrep, sautePan); err != nil {
		return err
	}

	// === TOAST PREPARATION ===
	if err := createVIP(toastPrep, burgerBun); err != nil {
		return err
	}
	if err := createVIMU(burgerBun, burgerUnitMeasurement); err != nil {
		return err
	}

	// === ASSEMBLE PREPARATION ===
	if err := createVIP(assemblePrep, burgerBun); err != nil {
		return err
	}
	if err := createVIP(assemblePrep, pickle); err != nil {
		return err
	}
	if err := createVIP(assemblePrep, burgerOnion); err != nil {
		return err
	}
	if err := createVIMU(pickle, burgerUnitMeasurement); err != nil {
		return err
	}
	if err := createVIMU(burgerOnion, sliceMeasurement); err != nil {
		return err
	}
	if err := createVPV(assemblePrep, burgerServingPlate); err != nil {
		return err
	}

	// === SMASH BURGER RECIPE BRIDGE ENTRIES ===
	// Get new preparations
	smashPrep := enums.Preparations["smash"]
	dividePrep := enums.Preparations["divide"]

	// Get ground beef ingredient
	groundBeef := enums.Ingredients["ground beef"]

	// Get cast iron skillet (reuse from earlier)
	smashBurgerSkillet := enums.Vessels["cast iron skillet"]

	// === DIVIDE PREPARATION ===
	if err := createVIP(dividePrep, groundBeef); err != nil {
		return err
	}
	if err := createVPI(dividePrep, burgerBareHands); err != nil {
		return err
	}

	// === FORM PREPARATION for ground beef ===
	if err := createVIP(formPrep, groundBeef); err != nil {
		return err
	}

	// === SEASON PREPARATION for ground beef ===
	if err := createVIP(burgerSeasonPrep, groundBeef); err != nil {
		return err
	}

	// === SMASH PREPARATION ===
	if err := createVIP(smashPrep, groundBeef); err != nil {
		return err
	}
	if err := createVPI(smashPrep, wideSpatula); err != nil {
		return err
	}
	if err := createVPV(smashPrep, smashBurgerSkillet); err != nil {
		return err
	}

	// === HEAT PREPARATION for cast iron skillet ===
	if err := createVPV(burgerHeatPrep, smashBurgerSkillet); err != nil {
		return err
	}

	// === PAN-SEAR PREPARATION for cast iron skillet ===
	if err := createVPV(burgerPanSearPrep, smashBurgerSkillet); err != nil {
		return err
	}

	// === FLIP PREPARATION for cast iron skillet ===
	if err := createVPV(flipPrep, smashBurgerSkillet); err != nil {
		return err
	}

	// === TOP PREPARATION for cast iron skillet ===
	if err := createVPV(topPrep, smashBurgerSkillet); err != nil {
		return err
	}

	// Ground beef measurement unit (ounces)
	if err := createVIMU(groundBeef, ounceMeasurement); err != nil {
		return err
	}

	// === SIMPLE WHITE RICE RECIPE BRIDGE ENTRIES ===
	// Get preparations for rice recipe
	simmerPrep := enums.Preparations["simmer"]
	stirPrep := enums.Preparations["stir"]
	coverPrep := enums.Preparations["cover"]
	fluffPrep := enums.Preparations["fluff"]
	riceRestPrep := enums.Preparations["rest"]

	// Get ingredients for rice recipe
	rice := enums.Ingredients["rice"]
	riceWater := enums.Ingredients["water"]
	riceSalt := enums.Ingredients["salt"]
	riceOliveOil := enums.Ingredients["olive oil"]

	// Get vessels for rice recipe
	saucepan := enums.Vessels["saucepan"]

	// Get instruments for rice recipe
	fork := enums.Instruments["fork"]
	woodenSpoon := enums.Instruments["wooden spoon"]

	// Get measurement units for rice recipe
	cupMeasurement := enums.MeasurementUnits["cup"]
	pinchMeasurement := enums.MeasurementUnits["pinch"]
	tablespoonMeasurement = enums.MeasurementUnits["tablespoon"]

	// === SIMMER PREPARATION (combine ingredients and bring to simmer) ===
	if err := createVIP(simmerPrep, rice); err != nil {
		return err
	}
	if err := createVIP(simmerPrep, riceWater); err != nil {
		return err
	}
	if err := createVIP(simmerPrep, riceSalt); err != nil {
		return err
	}
	if err := createVIP(simmerPrep, riceOliveOil); err != nil {
		return err
	}
	if err := createVPV(simmerPrep, saucepan); err != nil {
		return err
	}

	// === STIR PREPARATION ===
	if err := createVIP(stirPrep, rice); err != nil {
		return err
	}
	if err := createVPI(stirPrep, woodenSpoon); err != nil {
		return err
	}
	if err := createVPV(stirPrep, saucepan); err != nil {
		return err
	}

	// === COVER PREPARATION ===
	if err := createVPV(coverPrep, saucepan); err != nil {
		return err
	}

	// === FLUFF PREPARATION ===
	if err := createVIP(fluffPrep, rice); err != nil {
		return err
	}
	if err := createVPI(fluffPrep, fork); err != nil {
		return err
	}
	if err := createVPV(fluffPrep, saucepan); err != nil {
		return err
	}

	// === REST PREPARATION for rice ===
	if err := createVIP(riceRestPrep, rice); err != nil {
		return err
	}
	if err := createVPV(riceRestPrep, saucepan); err != nil {
		return err
	}

	// Rice measurement units
	if err := createVIMU(rice, cupMeasurement); err != nil {
		return err
	}
	if err := createVIMU(riceWater, cupMeasurement); err != nil {
		return err
	}
	if err := createVIMU(riceSalt, pinchMeasurement); err != nil {
		return err
	}
	if err := createVIMU(riceOliveOil, tablespoonMeasurement); err != nil {
		return err
	}

	// === ULTRA-FLUFFY MASHED POTATOES RECIPE BRIDGE ENTRIES ===
	// Get preparations for mashed potatoes recipe
	peelPrep := enums.Preparations["peel"]
	cubePrep = enums.Preparations["cube"]
	rinsePrep := enums.Preparations["rinse"]
	submergePrep := enums.Preparations["submerge"]
	boilPrep := enums.Preparations["boil"]
	drainPrep := enums.Preparations["drain"]
	mashedRestPrep := enums.Preparations["rest"]
	ricePrep := enums.Preparations["rice"]
	foldPrep := enums.Preparations["fold"]
	mashedSeasonPrep := enums.Preparations["season"]
	mashedSimmerPrep := enums.Preparations["simmer"]

	// Get ingredients for mashed potatoes recipe
	potato := enums.Ingredients["potato"]
	mashedPepper := enums.Ingredients["black pepper"]
	mashedButter := enums.Ingredients["butter"]
	milk := enums.Ingredients["milk"]
	mashedWater := enums.Ingredients["water"]

	// Get instruments for mashed potatoes recipe
	vegetablePeeler := enums.Instruments["vegetable peeler"]
	potatoRicer := enums.Instruments["potato ricer"]
	rubberSpatula := enums.Instruments["rubber spatula"]
	mashedKnife := enums.Instruments["knife"]

	// Get vessels for mashed potatoes recipe
	mashedCuttingBoard := enums.Vessels["cutting board"]
	mashedPot := enums.Vessels["pot"]
	mashedColander := enums.Vessels["colander"]

	// Get measurement units for mashed potatoes recipe
	poundMeasurement := enums.MeasurementUnits["pound"]
	mashedCupMeasurement := enums.MeasurementUnits["cup"]
	mashedTablespoonMeasurement := enums.MeasurementUnits["tablespoon"]

	// === PEEL PREPARATION ===
	if err := createVIP(peelPrep, potato); err != nil {
		return err
	}
	if err := createVPI(peelPrep, vegetablePeeler); err != nil {
		return err
	}
	if err := createVPV(peelPrep, mashedCuttingBoard); err != nil {
		return err
	}

	// === CUBE PREPARATION for potato ===
	if err := createVIP(cubePrep, potato); err != nil {
		return err
	}
	if err := createVPI(cubePrep, mashedKnife); err != nil {
		return err
	}
	if err := createVPV(cubePrep, mashedCuttingBoard); err != nil {
		return err
	}

	// === RINSE PREPARATION ===
	if err := createVIP(rinsePrep, potato); err != nil {
		return err
	}
	if err := createVPV(rinsePrep, mashedPot); err != nil {
		return err
	}

	// === SUBMERGE PREPARATION ===
	if err := createVIP(submergePrep, potato); err != nil {
		return err
	}
	if err := createVIP(submergePrep, mashedWater); err != nil {
		return err
	}
	if err := createVPV(submergePrep, mashedPot); err != nil {
		return err
	}

	// === SEASON PREPARATION for pot (seasoning water) ===
	if err := createVPV(mashedSeasonPrep, mashedPot); err != nil {
		return err
	}

	// === BOIL PREPARATION ===
	if err := createVIP(boilPrep, potato); err != nil {
		return err
	}
	if err := createVPV(boilPrep, mashedPot); err != nil {
		return err
	}

	// === DRAIN PREPARATION ===
	if err := createVIP(drainPrep, potato); err != nil {
		return err
	}
	if err := createVPV(drainPrep, mashedColander); err != nil {
		return err
	}

	// === RINSE PREPARATION for colander ===
	if err := createVPV(rinsePrep, mashedColander); err != nil {
		return err
	}

	// === REST PREPARATION for potato ===
	if err := createVIP(mashedRestPrep, potato); err != nil {
		return err
	}
	if err := createVPV(mashedRestPrep, mashedColander); err != nil {
		return err
	}

	// === RICE PREPARATION ===
	if err := createVIP(ricePrep, potato); err != nil {
		return err
	}
	if err := createVPI(ricePrep, potatoRicer); err != nil {
		return err
	}
	if err := createVPV(ricePrep, mashedPot); err != nil {
		return err
	}

	// === FOLD PREPARATION ===
	if err := createVIP(foldPrep, potato); err != nil {
		return err
	}
	if err := createVIP(foldPrep, mashedButter); err != nil {
		return err
	}
	if err := createVIP(foldPrep, milk); err != nil {
		return err
	}
	if err := createVPI(foldPrep, rubberSpatula); err != nil {
		return err
	}
	if err := createVPV(foldPrep, mashedPot); err != nil {
		return err
	}

	// === SIMMER PREPARATION for milk ===
	if err := createVIP(mashedSimmerPrep, milk); err != nil {
		return err
	}
	if err := createVPV(mashedSimmerPrep, mashedPot); err != nil {
		return err
	}

	// === SEASON PREPARATION for mashed potatoes ===
	if err := createVIP(mashedSeasonPrep, potato); err != nil {
		return err
	}
	if err := createVIP(mashedSeasonPrep, mashedPepper); err != nil {
		return err
	}
	if err := createVPV(mashedSeasonPrep, mashedPot); err != nil {
		return err
	}

	// === POTATO MEASUREMENT UNITS ===
	if err := createVIMU(potato, poundMeasurement); err != nil {
		return err
	}
	if err := createVIMU(milk, mashedCupMeasurement); err != nil {
		return err
	}
	if err := createVIMU(mashedButter, mashedTablespoonMeasurement); err != nil {
		return err
	}

	// === CAESAR ROASTED BROCCOLI RECIPE BRIDGE ENTRIES ===
	// Get preparations for caesar roasted broccoli recipe
	caesarMeltPrep := enums.Preparations["melt"]
	caesarStirPrep := enums.Preparations["stir"]
	caesarCookPrep := enums.Preparations["cook"]
	caesarToastPrep := enums.Preparations["toast"]
	caesarSeasonPrep := enums.Preparations["season"]
	caesarTransferPrep := enums.Preparations["transfer"]
	caesarLinePrep := enums.Preparations["line"]
	caesarPreheatPrep := enums.Preparations["preheat"]
	caesarTossPrep := enums.Preparations["toss"]
	caesarRoastPrep := enums.Preparations["roast"]
	caesarTopPrep := enums.Preparations["top"]
	caesarZestPrep := enums.Preparations["zest"]

	// Get ingredients for caesar roasted broccoli recipe
	caesarSaltedButter := enums.Ingredients["salted butter"]
	caesarBreadcrumbs := enums.Ingredients["breadcrumbs"]
	caesarAnchovyPaste := enums.Ingredients["anchovy paste"]
	caesarGarlic := enums.Ingredients["garlic"]
	caesarLemon := enums.Ingredients["lemon"]
	caesarBroccoli := enums.Ingredients["broccoli"]
	caesarOliveOil := enums.Ingredients["olive oil"]
	caesarSalt := enums.Ingredients["salt"]
	caesarPepper := enums.Ingredients["black pepper"]
	caesarParmesan := enums.Ingredients["parmesan cheese"]

	// Get instruments for caesar roasted broccoli recipe
	caesarRubberSpatula := enums.Instruments["rubber spatula"]
	caesarAluminumFoil := enums.Instruments["aluminum foil"]
	caesarMicroplane := enums.Instruments["microplane"]

	// Get vessels for caesar roasted broccoli recipe
	caesarSmallNonstickSkillet := enums.Vessels["small nonstick skillet"]
	caesarSmallBowl := enums.Vessels["small bowl"]
	caesarBakingSheet := enums.Vessels["baking sheet"]
	caesarLargeBowl := enums.Vessels["large bowl"]
	caesarServingPlatter := enums.Vessels["serving platter"]
	caesarOven := enums.Vessels["oven"]

	// Get measurement units for caesar roasted broccoli recipe
	caesarTablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	caesarTeaspoonMeasurement := enums.MeasurementUnits["teaspoon"]
	caesarCupMeasurement := enums.MeasurementUnits["cup"]
	caesarPoundMeasurement := enums.MeasurementUnits["pound"]
	caesarGramMeasurement := enums.MeasurementUnits["gram"]
	caesarUnitMeasurement := enums.MeasurementUnits["unit"]

	// === MELT PREPARATION ===
	if err := createVIP(caesarMeltPrep, caesarSaltedButter); err != nil {
		return err
	}
	if err := createVPV(caesarMeltPrep, caesarSmallNonstickSkillet); err != nil {
		return err
	}

	// === STIR PREPARATION for anchovy paste and garlic ===
	if err := createVIP(caesarStirPrep, caesarAnchovyPaste); err != nil {
		return err
	}
	if err := createVIP(caesarStirPrep, caesarGarlic); err != nil {
		return err
	}
	if err := createVIP(caesarStirPrep, caesarBreadcrumbs); err != nil {
		return err
	}
	if err := createVIP(caesarStirPrep, caesarLemon); err != nil {
		return err
	}
	if err := createVPV(caesarStirPrep, caesarSmallNonstickSkillet); err != nil {
		return err
	}
	if err := createVPI(caesarStirPrep, caesarRubberSpatula); err != nil {
		return err
	}

	// === COOK PREPARATION for breadcrumbs ===
	if err := createVIP(caesarCookPrep, caesarAnchovyPaste); err != nil {
		return err
	}
	if err := createVIP(caesarCookPrep, caesarGarlic); err != nil {
		return err
	}
	if err := createVPV(caesarCookPrep, caesarSmallNonstickSkillet); err != nil {
		return err
	}

	// === TOAST PREPARATION for breadcrumbs ===
	if err := createVIP(caesarToastPrep, caesarBreadcrumbs); err != nil {
		return err
	}
	if err := createVPV(caesarToastPrep, caesarSmallNonstickSkillet); err != nil {
		return err
	}
	if err := createVPI(caesarToastPrep, caesarRubberSpatula); err != nil {
		return err
	}

	// === ZEST PREPARATION for lemon ===
	if err := createVIP(caesarZestPrep, caesarLemon); err != nil {
		return err
	}
	if err := createVPI(caesarZestPrep, caesarMicroplane); err != nil {
		return err
	}

	// === SEASON PREPARATION for breadcrumbs ===
	if err := createVIP(caesarSeasonPrep, caesarBreadcrumbs); err != nil {
		return err
	}

	// === TRANSFER PREPARATION for breadcrumbs ===
	if err := createVIP(caesarTransferPrep, caesarBreadcrumbs); err != nil {
		return err
	}
	if err := createVPV(caesarTransferPrep, caesarSmallBowl); err != nil {
		return err
	}
	if err := createVPV(caesarTransferPrep, caesarServingPlatter); err != nil {
		return err
	}

	// === LINE PREPARATION for baking sheet ===
	if err := createVPI(caesarLinePrep, caesarAluminumFoil); err != nil {
		return err
	}
	if err := createVPV(caesarLinePrep, caesarBakingSheet); err != nil {
		return err
	}

	// === PREHEAT PREPARATION for oven ===
	if err := createVPV(caesarPreheatPrep, caesarOven); err != nil {
		return err
	}
	if err := createVPV(caesarPreheatPrep, caesarBakingSheet); err != nil {
		return err
	}

	// === TOSS PREPARATION for broccoli ===
	if err := createVIP(caesarTossPrep, caesarBroccoli); err != nil {
		return err
	}
	if err := createVIP(caesarTossPrep, caesarOliveOil); err != nil {
		return err
	}
	if err := createVIP(caesarTossPrep, caesarSalt); err != nil {
		return err
	}
	if err := createVIP(caesarTossPrep, caesarPepper); err != nil {
		return err
	}
	if err := createVIP(caesarTossPrep, caesarLemon); err != nil {
		return err
	}
	if err := createVPV(caesarTossPrep, caesarLargeBowl); err != nil {
		return err
	}

	// === TRANSFER PREPARATION for broccoli ===
	if err := createVIP(caesarTransferPrep, caesarBroccoli); err != nil {
		return err
	}
	if err := createVPV(caesarTransferPrep, caesarBakingSheet); err != nil {
		return err
	}

	// === ROAST PREPARATION for broccoli ===
	if err := createVIP(caesarRoastPrep, caesarBroccoli); err != nil {
		return err
	}
	if err := createVPV(caesarRoastPrep, caesarBakingSheet); err != nil {
		return err
	}
	if err := createVPV(caesarRoastPrep, caesarOven); err != nil {
		return err
	}

	// === TOP PREPARATION for broccoli ===
	if err := createVIP(caesarTopPrep, caesarBroccoli); err != nil {
		return err
	}
	if err := createVIP(caesarTopPrep, caesarBreadcrumbs); err != nil {
		return err
	}
	if err := createVIP(caesarTopPrep, caesarParmesan); err != nil {
		return err
	}
	if err := createVPV(caesarTopPrep, caesarServingPlatter); err != nil {
		return err
	}

	// === CAESAR BROCCOLI INGREDIENT MEASUREMENT UNITS ===
	if err := createVIMU(caesarSaltedButter, caesarTablespoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(caesarBreadcrumbs, caesarCupMeasurement); err != nil {
		return err
	}
	if err := createVIMU(caesarAnchovyPaste, caesarTeaspoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(caesarGarlic, caesarUnitMeasurement); err != nil {
		return err
	}
	if err := createVIMU(caesarLemon, caesarTeaspoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(caesarBroccoli, caesarPoundMeasurement); err != nil {
		return err
	}
	if err := createVIMU(caesarParmesan, caesarTablespoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(caesarOliveOil, caesarTablespoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(caesarSalt, caesarTeaspoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(caesarPepper, caesarGramMeasurement); err != nil {
		return err
	}

	// === HARICOTS VERTS AMANDINE RECIPE BRIDGE ENTRIES ===
	// Get preparations for haricots verts amandine recipe
	hvaBoilPrep := enums.Preparations["boil"]
	hvaBlanch := enums.Preparations["blanch"]
	hvaShockPrep := enums.Preparations["shock"]
	hvaDrainPrep := enums.Preparations["drain"]
	hvaDryPrep := enums.Preparations["dry"]
	hvaHeatPrep := enums.Preparations["heat"]
	hvaToastPrep := enums.Preparations["toast"]
	hvaCookPrep := enums.Preparations["cook"]
	hvaStirPrep := enums.Preparations["stir"]
	hvaEmulsifyPrep := enums.Preparations["emulsify"]
	hvaSeasonPrep := enums.Preparations["season"]
	hvaTossPrep := enums.Preparations["toss"]
	hvaTransferPrep := enums.Preparations["transfer"]
	hvaTrimPrep := enums.Preparations["trim"]

	// Get ingredients for haricots verts amandine recipe
	hvaGreenBeans := enums.Ingredients["green beans"]
	hvaButter := enums.Ingredients["butter"]
	hvaSliveredAlmonds := enums.Ingredients["slivered almonds"]
	hvaGarlic := enums.Ingredients["garlic"]
	hvaShallot := enums.Ingredients["shallot"]
	hvaLemon := enums.Ingredients["lemon"]
	hvaSalt := enums.Ingredients["salt"]
	hvaPepper := enums.Ingredients["black pepper"]
	hvaWater := enums.Ingredients["water"]

	// Get instruments for haricots verts amandine recipe
	hvaWireMeshSpider := enums.Instruments["wire mesh spider"]
	hvaTongs := enums.Instruments["tongs"]
	hvaPaperTowels := enums.Instruments["paper towels"]
	hvaKitchenTowels := enums.Instruments["kitchen towels"]
	hvaRubberSpatula := enums.Instruments["rubber spatula"]

	// Get vessels for haricots verts amandine recipe
	hvaPot := enums.Vessels["pot"]
	hvaLargeBowl := enums.Vessels["large bowl"]
	hvaMediumSkillet := enums.Vessels["medium skillet"]
	hvaServingPlatter := enums.Vessels["serving platter"]
	hvaColander := enums.Vessels["colander"]

	// Get measurement units for haricots verts amandine recipe
	hvaPoundMeasurement := enums.MeasurementUnits["pound"]
	hvaTablespoonMeasurement := enums.MeasurementUnits["tablespoon"]
	hvaOunceMeasurement := enums.MeasurementUnits["ounce"]
	hvaUnitMeasurement := enums.MeasurementUnits["unit"]

	// === BOIL PREPARATION for water ===
	if err := createVIP(hvaBoilPrep, hvaWater); err != nil {
		return err
	}
	if err := createVIP(hvaBoilPrep, hvaSalt); err != nil {
		return err
	}
	if err := createVPV(hvaBoilPrep, hvaPot); err != nil {
		return err
	}

	// === BLANCH PREPARATION for green beans ===
	if err := createVIP(hvaBlanch, hvaGreenBeans); err != nil {
		return err
	}
	if err := createVPI(hvaBlanch, hvaWireMeshSpider); err != nil {
		return err
	}
	if err := createVPI(hvaBlanch, hvaTongs); err != nil {
		return err
	}
	if err := createVPV(hvaBlanch, hvaPot); err != nil {
		return err
	}

	// === SHOCK PREPARATION for green beans ===
	if err := createVIP(hvaShockPrep, hvaGreenBeans); err != nil {
		return err
	}
	if err := createVPI(hvaShockPrep, hvaWireMeshSpider); err != nil {
		return err
	}
	if err := createVPI(hvaShockPrep, hvaTongs); err != nil {
		return err
	}
	if err := createVPV(hvaShockPrep, hvaLargeBowl); err != nil {
		return err
	}

	// === DRAIN PREPARATION for green beans ===
	if err := createVIP(hvaDrainPrep, hvaGreenBeans); err != nil {
		return err
	}
	if err := createVPV(hvaDrainPrep, hvaColander); err != nil {
		return err
	}

	// === DRY PREPARATION for green beans ===
	if err := createVIP(hvaDryPrep, hvaGreenBeans); err != nil {
		return err
	}
	if err := createVPI(hvaDryPrep, hvaPaperTowels); err != nil {
		return err
	}
	if err := createVPI(hvaDryPrep, hvaKitchenTowels); err != nil {
		return err
	}

	// === HEAT PREPARATION for butter and skillet ===
	if err := createVIP(hvaHeatPrep, hvaButter); err != nil {
		return err
	}
	if err := createVIP(hvaHeatPrep, hvaSliveredAlmonds); err != nil {
		return err
	}
	if err := createVPV(hvaHeatPrep, hvaMediumSkillet); err != nil {
		return err
	}
	if err := createVPI(hvaHeatPrep, hvaRubberSpatula); err != nil {
		return err
	}

	// === TOAST PREPARATION for almonds ===
	if err := createVIP(hvaToastPrep, hvaSliveredAlmonds); err != nil {
		return err
	}
	if err := createVPI(hvaToastPrep, hvaRubberSpatula); err != nil {
		return err
	}
	if err := createVPV(hvaToastPrep, hvaMediumSkillet); err != nil {
		return err
	}

	// === COOK PREPARATION for garlic and shallot ===
	if err := createVIP(hvaCookPrep, hvaGarlic); err != nil {
		return err
	}
	if err := createVIP(hvaCookPrep, hvaShallot); err != nil {
		return err
	}
	if err := createVPI(hvaCookPrep, hvaRubberSpatula); err != nil {
		return err
	}
	if err := createVPV(hvaCookPrep, hvaMediumSkillet); err != nil {
		return err
	}

	// === STIR PREPARATION for lemon juice and water ===
	if err := createVIP(hvaStirPrep, hvaLemon); err != nil {
		return err
	}
	if err := createVIP(hvaStirPrep, hvaWater); err != nil {
		return err
	}
	if err := createVPV(hvaStirPrep, hvaMediumSkillet); err != nil {
		return err
	}

	// === EMULSIFY PREPARATION for sauce ===
	if err := createVIP(hvaEmulsifyPrep, hvaButter); err != nil {
		return err
	}
	if err := createVIP(hvaEmulsifyPrep, hvaLemon); err != nil {
		return err
	}
	if err := createVIP(hvaEmulsifyPrep, hvaWater); err != nil {
		return err
	}
	if err := createVPV(hvaEmulsifyPrep, hvaMediumSkillet); err != nil {
		return err
	}

	// === SEASON PREPARATION for sauce ===
	if err := createVIP(hvaSeasonPrep, hvaPepper); err != nil {
		return err
	}
	if err := createVPV(hvaSeasonPrep, hvaMediumSkillet); err != nil {
		return err
	}

	// === TOSS PREPARATION for green beans with sauce ===
	if err := createVIP(hvaTossPrep, hvaGreenBeans); err != nil {
		return err
	}
	if err := createVPV(hvaTossPrep, hvaMediumSkillet); err != nil {
		return err
	}

	// === TRANSFER PREPARATION for finished dish ===
	if err := createVIP(hvaTransferPrep, hvaGreenBeans); err != nil {
		return err
	}
	if err := createVPV(hvaTransferPrep, hvaServingPlatter); err != nil {
		return err
	}

	// === TRIM PREPARATION for green beans ===
	if err := createVIP(hvaTrimPrep, hvaGreenBeans); err != nil {
		return err
	}

	// === HARICOTS VERTS AMANDINE INGREDIENT MEASUREMENT UNITS ===
	if err := createVIMU(hvaGreenBeans, hvaPoundMeasurement); err != nil {
		return err
	}
	if err := createVIMU(hvaButter, hvaTablespoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(hvaSliveredAlmonds, hvaOunceMeasurement); err != nil {
		return err
	}
	if err := createVIMU(hvaGarlic, hvaUnitMeasurement); err != nil {
		return err
	}
	if err := createVIMU(hvaShallot, hvaUnitMeasurement); err != nil {
		return err
	}
	if err := createVIMU(hvaLemon, hvaTablespoonMeasurement); err != nil {
		return err
	}
	if err := createVIMU(hvaWater, hvaTablespoonMeasurement); err != nil {
		return err
	}

	return nil
}
