package types

type (
	EatingUserDataCollection struct {
		_ struct{} `json:"-"`

		HouseholdInstrumentOwnerships map[string][]HouseholdInstrumentOwnership `json:"householdInstrumentOwnerships"`
		MealPlans                     map[string][]MealPlan                     `json:"mealPlans"`
		ReportID                      string                                    `json:"reportID"`
		RecipeRatings                 []RecipeRating                            `json:"recipeRatings"`
		Recipes                       []Recipe                                  `json:"recipes"`
		Meals                         []Meal                                    `json:"meals"`
		UserIngredientPreferences     []UserIngredientPreference                `json:"userIngredientPreferences"`
	}
)
