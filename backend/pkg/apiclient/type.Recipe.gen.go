// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	Recipe struct {
		Slug                string                      `json:"slug"`
		CreatedAt           string                      `json:"createdAt"`
		CreatedByUser       string                      `json:"createdByUser"`
		Description         string                      `json:"description"`
		LastUpdatedAt       string                      `json:"lastUpdatedAt"`
		Source              string                      `json:"source"`
		ID                  string                      `json:"id"`
		Name                string                      `json:"name"`
		YieldsComponentType string                      `json:"yieldsComponentType"`
		ArchivedAt          string                      `json:"archivedAt"`
		InspiredByRecipeID  string                      `json:"inspiredByRecipeID"`
		PluralPortionName   string                      `json:"pluralPortionName"`
		PortionName         string                      `json:"portionName"`
		PrepTasks           []RecipePrepTask            `json:"prepTasks"`
		Steps               []RecipeStep                `json:"steps"`
		Media               []RecipeMedia               `json:"media"`
		SupportingRecipes   []Recipe                    `json:"supportingRecipes"`
		EstimatedPortions   Float32RangeWithOptionalMax `json:"estimatedPortions"`
		SealOfApproval      bool                        `json:"sealOfApproval"`
		EligibleForMeals    bool                        `json:"eligibleForMeals"`
	}
)
