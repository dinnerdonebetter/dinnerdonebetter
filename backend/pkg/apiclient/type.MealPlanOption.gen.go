// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	MealPlanOption struct {
		LastUpdatedAt          string               `json:"lastUpdatedAt"`
		ID                     string               `json:"id"`
		AssignedDishwasher     string               `json:"assignedDishwasher"`
		BelongsToMealPlanEvent string               `json:"belongsToMealPlanEvent"`
		CreatedAt              string               `json:"createdAt"`
		ArchivedAt             string               `json:"archivedAt"`
		AssignedCook           string               `json:"assignedCook"`
		Notes                  string               `json:"notes"`
		Votes                  []MealPlanOptionVote `json:"votes"`
		Meal                   Meal                 `json:"meal"`
		MealScale              float64              `json:"mealScale"`
		Chosen                 bool                 `json:"chosen"`
		TieBroken              bool                 `json:"tieBroken"`
	}
)
