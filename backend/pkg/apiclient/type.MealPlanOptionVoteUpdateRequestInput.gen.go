// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	MealPlanOptionVoteUpdateRequestInput struct {
		BelongsToMealPlanOption string `json:"belongsToMealPlanOption"`
		Notes                   string `json:"notes"`
		Rank                    uint64 `json:"rank"`
		Abstain                 bool   `json:"abstain"`
	}
)
