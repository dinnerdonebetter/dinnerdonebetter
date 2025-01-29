// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
MealCreationRequestInput struct {
   Components []MealComponentCreationRequestInput `json:"components"`
 Description string `json:"description"`
 EligibleForMealPlans bool `json:"eligibleForMealPlans"`
 EstimatedPortions Float32RangeWithOptionalMax `json:"estimatedPortions"`
 Name string `json:"name"`

}
)
