package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertValidIngredientGroupToValidIngredientGroupUpdateRequestInput creates a DatabaseCreationInput from a CreationInput.
func ConvertValidIngredientGroupToValidIngredientGroupUpdateRequestInput(input *types.ValidIngredientGroup) *types.ValidIngredientGroupUpdateRequestInput {
	x := &types.ValidIngredientGroupUpdateRequestInput{
		Name:        &input.Name,
		Description: &input.Description,
		Slug:        &input.Slug,
	}

	return x
}

// ConvertValidIngredientGroupCreationRequestInputToValidIngredientGroupDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertValidIngredientGroupCreationRequestInputToValidIngredientGroupDatabaseCreationInput(input *types.ValidIngredientGroupCreationRequestInput) *types.ValidIngredientGroupDatabaseCreationInput {
	x := &types.ValidIngredientGroupDatabaseCreationInput{
		ID:          identifiers.New(),
		Name:        input.Name,
		Description: input.Description,
		Slug:        input.Slug,
	}

	return x
}

// ConvertValidIngredientGroupToValidIngredientGroupCreationRequestInput builds a ValidIngredientGroupCreationRequestInput from a Ingredient.
func ConvertValidIngredientGroupToValidIngredientGroupCreationRequestInput(validIngredient *types.ValidIngredientGroup) *types.ValidIngredientGroupCreationRequestInput {
	return &types.ValidIngredientGroupCreationRequestInput{
		Name:        validIngredient.Name,
		Description: validIngredient.Description,
		Slug:        validIngredient.Slug,
	}
}

// ConvertValidIngredientGroupToValidIngredientGroupDatabaseCreationInput builds a ValidIngredientGroupDatabaseCreationInput from a ValidIngredientGroup.
func ConvertValidIngredientGroupToValidIngredientGroupDatabaseCreationInput(validIngredient *types.ValidIngredientGroup) *types.ValidIngredientGroupDatabaseCreationInput {
	members := make([]*types.ValidIngredientGroupMemberDatabaseCreationInput, len(validIngredient.Members))
	for i, member := range validIngredient.Members {
		members[i] = &types.ValidIngredientGroupMemberDatabaseCreationInput{
			ID:                member.ID,
			ValidIngredientID: member.ValidIngredientID,
		}
	}

	return &types.ValidIngredientGroupDatabaseCreationInput{
		ID:          validIngredient.ID,
		Name:        validIngredient.Name,
		Description: validIngredient.Description,
		Slug:        validIngredient.Slug,
		Members:     members,
	}
}
