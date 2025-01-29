// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
Meal struct {
   ArchivedAt string `json:"archivedAt"`
 Components []MealComponent `json:"components"`
 CreatedAt string `json:"createdAt"`
 CreatedByUser string `json:"createdByUser"`
 Description string `json:"description"`
 EligibleForMealPlans bool `json:"eligibleForMealPlans"`
 EstimatedPortions Float32RangeWithOptionalMax `json:"estimatedPortions"`
 ID string `json:"id"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Name string `json:"name"`

}
)
