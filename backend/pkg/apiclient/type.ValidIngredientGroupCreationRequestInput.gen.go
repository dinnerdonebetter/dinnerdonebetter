// GENERATED CODE, DO NOT EDIT MANUALLY

package apiclient

type (
	ValidIngredientGroupCreationRequestInput struct {
		Description string                                           `json:"description"`
		Name        string                                           `json:"name"`
		Slug        string                                           `json:"slug"`
		Members     []ValidIngredientGroupMemberCreationRequestInput `json:"members"`
	}
)
