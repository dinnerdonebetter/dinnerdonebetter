// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	EatingUserDataCollection struct {
		HouseholdInstrumentOwnerships map[string]any             `json:"householdInstrumentOwnerships"`
		MealPlans                     map[string]any             `json:"mealPlans"`
		Meals                         []Meal                     `json:"meals"`
		RecipeRatings                 []RecipeRating             `json:"recipeRatings"`
		Recipes                       []Recipe                   `json:"recipes"`
		ReportID                      string                     `json:"reportID"`
		UserIngredientPreferences     []UserIngredientPreference `json:"userIngredientPreferences"`
	}
)
