// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	Recipe struct {
		Slug                string                      `json:"slug"`
		LastUpdatedAt       string                      `json:"lastUpdatedAt"`
		CreatedByUser       string                      `json:"createdByUser"`
		Description         string                      `json:"description"`
		YieldsComponentType string                      `json:"yieldsComponentType"`
		Source              string                      `json:"source"`
		ID                  string                      `json:"id"`
		InspiredByRecipeID  string                      `json:"inspiredByRecipeID"`
		ArchivedAt          string                      `json:"archivedAt"`
		CreatedAt           string                      `json:"createdAt"`
		PortionName         string                      `json:"portionName"`
		PluralPortionName   string                      `json:"pluralPortionName"`
		Name                string                      `json:"name"`
		PrepTasks           []RecipePrepTask            `json:"prepTasks"`
		Media               []RecipeMedia               `json:"media"`
		Steps               []RecipeStep                `json:"steps"`
		SupportingRecipes   []Recipe                    `json:"supportingRecipes"`
		EstimatedPortions   Float32RangeWithOptionalMax `json:"estimatedPortions"`
		SealOfApproval      bool                        `json:"sealOfApproval"`
		EligibleForMeals    bool                        `json:"eligibleForMeals"`
	}
)
