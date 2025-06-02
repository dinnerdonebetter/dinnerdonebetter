// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	UserIngredientPreference struct {
		ArchivedAt    string          `json:"archivedAt"`
		BelongsToUser string          `json:"belongsToUser"`
		CreatedAt     string          `json:"createdAt"`
		ID            string          `json:"id"`
		LastUpdatedAt string          `json:"lastUpdatedAt"`
		Notes         string          `json:"notes"`
		Ingredient    ValidIngredient `json:"ingredient"`
		Rating        uint64          `json:"rating"`
		Allergy       bool            `json:"allergy"`
	}
)
