package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

// BuildFakeValidIngredientGroup builds a faked valid ingredient.
func BuildFakeValidIngredientGroup() *types.ValidIngredientGroup {
	return &types.ValidIngredientGroup{
		ID:          BuildFakeID(),
		Name:        buildUniqueString(),
		Description: buildUniqueString(),
		CreatedAt:   BuildFakeTime(),
		Slug:        buildUniqueString(),
	}
}

// BuildFakeValidIngredientGroupList builds a faked ValidIngredientGroupList.
func BuildFakeValidIngredientGroupList() *types.QueryFilteredResult[types.ValidIngredientGroup] {
	var examples []*types.ValidIngredientGroup
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidIngredientGroup())
	}

	return &types.QueryFilteredResult[types.ValidIngredientGroup]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidIngredientGroupUpdateRequestInput builds a faked ValidIngredientGroupUpdateRequestInput from a valid ingredient.
func BuildFakeValidIngredientGroupUpdateRequestInput() *types.ValidIngredientGroupUpdateRequestInput {
	validIngredient := BuildFakeValidIngredientGroup()
	return converters.ConvertValidIngredientGroupToValidIngredientGroupUpdateRequestInput(validIngredient)
}

// BuildFakeValidIngredientGroupCreationRequestInput builds a faked ValidIngredientGroupCreationRequestInput.
func BuildFakeValidIngredientGroupCreationRequestInput() *types.ValidIngredientGroupCreationRequestInput {
	validIngredient := BuildFakeValidIngredientGroup()
	return converters.ConvertValidIngredientGroupToValidIngredientGroupCreationRequestInput(validIngredient)
}
