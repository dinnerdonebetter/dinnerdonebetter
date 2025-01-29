// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
Recipe struct {
   ArchivedAt string `json:"archivedAt"`
 CreatedAt string `json:"createdAt"`
 CreatedByUser string `json:"createdByUser"`
 Description string `json:"description"`
 EligibleForMeals bool `json:"eligibleForMeals"`
 EstimatedPortions Float32RangeWithOptionalMax `json:"estimatedPortions"`
 ID string `json:"id"`
 InspiredByRecipeID string `json:"inspiredByRecipeID"`
 LastUpdatedAt string `json:"lastUpdatedAt"`
 Media []RecipeMedia `json:"media"`
 Name string `json:"name"`
 PluralPortionName string `json:"pluralPortionName"`
 PortionName string `json:"portionName"`
 PrepTasks []RecipePrepTask `json:"prepTasks"`
 SealOfApproval bool `json:"sealOfApproval"`
 Slug string `json:"slug"`
 Source string `json:"source"`
 Steps []RecipeStep `json:"steps"`
 SupportingRecipes []Recipe `json:"supportingRecipes"`
 YieldsComponentType string `json:"yieldsComponentType"`

}
)
