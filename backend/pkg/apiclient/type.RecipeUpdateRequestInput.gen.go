// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient



type (
RecipeUpdateRequestInput struct {
   Description string `json:"description"`
 EligibleForMeals bool `json:"eligibleForMeals"`
 EstimatedPortions Float32RangeWithOptionalMaxUpdateRequestInput `json:"estimatedPortions"`
 InspiredByRecipeID string `json:"inspiredByRecipeID"`
 Name string `json:"name"`
 PluralPortionName string `json:"pluralPortionName"`
 PortionName string `json:"portionName"`
 SealOfApproval bool `json:"sealOfApproval"`
 Slug string `json:"slug"`
 Source string `json:"source"`
 YieldsComponentType string `json:"yieldsComponentType"`

}
)
