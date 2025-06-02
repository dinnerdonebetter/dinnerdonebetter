// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	RecipeStepInstrument struct {
		Name                string                     `json:"name"`
		BelongsToRecipeStep string                     `json:"belongsToRecipeStep"`
		CreatedAt           string                     `json:"createdAt"`
		ID                  string                     `json:"id"`
		LastUpdatedAt       string                     `json:"lastUpdatedAt"`
		ArchivedAt          string                     `json:"archivedAt"`
		Notes               string                     `json:"notes"`
		RecipeStepProductID string                     `json:"recipeStepProductID"`
		Instrument          ValidInstrument            `json:"instrument"`
		Quantity            Uint32RangeWithOptionalMax `json:"quantity"`
		OptionIndex         uint64                     `json:"optionIndex"`
		PreferenceRank      uint64                     `json:"preferenceRank"`
		Optional            bool                       `json:"optional"`
	}
)
