// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
MealPlanOption struct {
   ArchivedAt string `json:"archivedAt"`
 AssignedCook string `json:"assignedCook"`
 AssignedDishwasher string `json:"assignedDishwasher"`
 BelongsToMealPlanEvent string `json:"belongsToMealPlanEvent"`
 Chosen bool `json:"chosen"`
 CreatedAt string `json:"createdAt"`
 ID string `json:"id"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Meal Meal `json:"meal"`
 MealScale float64 `json:"mealScale"`
 Notes string `json:"notes"`
 TieBroken bool `json:"tieBroken"`
 Votes []MealPlanOptionVote `json:"votes"`

}
)
