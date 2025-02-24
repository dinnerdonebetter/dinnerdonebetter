package types

type (
	EatingUserDataCollection struct {
		_ struct{} `json:"-"`

		InstrumentOwnerships  map[string][]InstrumentOwnership `json:"householdInstrumentOwnerships"`
		MealPlans             map[string][]MealPlan            `json:"mealPlans"`
		ReportID              string                           `json:"reportID"`
		RecipeRatings         []RecipeRating                   `json:"recipeRatings"`
		Recipes               []Recipe                         `json:"recipes"`
		Meals                 []Meal                           `json:"meals"`
		IngredientPreferences []IngredientPreference           `json:"userIngredientPreferences"`
	}
)
