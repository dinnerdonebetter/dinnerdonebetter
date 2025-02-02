// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidIngredientPreparation struct {
		ArchivedAt    string           `json:"archivedAt"`
		CreatedAt     string           `json:"createdAt"`
		ID            string           `json:"id"`
		Ingredient    ValidIngredient  `json:"ingredient"`
		LastUpdatedAt string           `json:"lastUpdatedAt"`
		Notes         string           `json:"notes"`
		Preparation   ValidPreparation `json:"preparation"`
	}
)
