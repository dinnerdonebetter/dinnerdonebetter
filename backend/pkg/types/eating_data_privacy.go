package types

type (
	EatingUserDataCollection struct {
		_ struct{} `json:"-"`

		AccountInstrumentOwnerships map[string][]AccountInstrumentOwnership `json:"accountInstrumentOwnerships"`
		MealPlans                   map[string][]MealPlan                   `json:"mealPlans"`
		ReportID                    string                                  `json:"reportID"`
		RecipeRatings               []RecipeRating                          `json:"recipeRatings"`
		Recipes                     []Recipe                                `json:"recipes"`
		Meals                       []Meal                                  `json:"meals"`
		UserIngredientPreferences   []UserIngredientPreference              `json:"userIngredientPreferences"`
	}
)
