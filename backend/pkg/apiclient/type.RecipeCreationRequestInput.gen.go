// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeCreationRequestInput struct {
   AlsoCreateMeal bool `json:"alsoCreateMeal"`
 Description string `json:"description"`
 EligibleForMeals bool `json:"eligibleForMeals"`
 EstimatedPortions Float32RangeWithOptionalMax `json:"estimatedPortions"`
 InspiredByRecipeID string `json:"inspiredByRecipeID"`
 Name string `json:"name"`
 PluralPortionName string `json:"pluralPortionName"`
 PortionName string `json:"portionName"`
 PrepTasks []RecipePrepTaskWithinRecipeCreationRequestInput `json:"prepTasks"`
 SealOfApproval bool `json:"sealOfApproval"`
 Slug string `json:"slug"`
 Source string `json:"source"`
 Steps []RecipeStepCreationRequestInput `json:"steps"`
 YieldsComponentType string `json:"yieldsComponentType"`

}
)
