// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeRatingCreationRequestInput struct {
		ByUser       string  `json:"byUser"`
		Notes        string  `json:"notes"`
		RecipeID     string  `json:"recipeID"`
		Cleanup      float64 `json:"cleanup"`
		Difficulty   float64 `json:"difficulty"`
		Instructions float64 `json:"instructions"`
		Overall      float64 `json:"overall"`
		Taste        float64 `json:"taste"`
	}
)
