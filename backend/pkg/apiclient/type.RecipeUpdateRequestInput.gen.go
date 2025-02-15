// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeUpdateRequestInput struct {
		Description         string                                        `json:"description"`
		InspiredByRecipeID  string                                        `json:"inspiredByRecipeID"`
		Name                string                                        `json:"name"`
		PluralPortionName   string                                        `json:"pluralPortionName"`
		PortionName         string                                        `json:"portionName"`
		Slug                string                                        `json:"slug"`
		Source              string                                        `json:"source"`
		YieldsComponentType string                                        `json:"yieldsComponentType"`
		EstimatedPortions   Float32RangeWithOptionalMaxUpdateRequestInput `json:"estimatedPortions"`
		EligibleForMeals    bool                                          `json:"eligibleForMeals"`
		SealOfApproval      bool                                          `json:"sealOfApproval"`
	}
)
