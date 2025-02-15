// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidIngredientGroup struct {
		ArchivedAt    string                       `json:"archivedAt"`
		CreatedAt     string                       `json:"createdAt"`
		Description   string                       `json:"description"`
		ID            string                       `json:"id"`
		LastUpdatedAt string                       `json:"lastUpdatedAt"`
		Name          string                       `json:"name"`
		Slug          string                       `json:"slug"`
		Members       []ValidIngredientGroupMember `json:"members"`
	}
)
