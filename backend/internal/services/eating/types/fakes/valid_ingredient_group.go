package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"
)

// BuildFakeValidIngredientGroup builds a faked valid ingredient group.
func BuildFakeValidIngredientGroup() *types.ValidIngredientGroup {
	groupID := BuildFakeID()

	var members []*types.ValidIngredientGroupMember
	for i := 0; i < exampleQuantity; i++ {
		newMember := BuildFakeValidIngredientGroupMember()
		newMember.BelongsToGroup = groupID
		members = append(members, newMember)
	}

	return &types.ValidIngredientGroup{
		ID:          groupID,
		Name:        buildUniqueString(),
		Description: buildUniqueString(),
		CreatedAt:   BuildFakeTime(),
		Slug:        buildUniqueString(),
		Members:     members,
	}
}

// BuildFakeValidIngredientGroupMember builds a faked valid ingredient group.
func BuildFakeValidIngredientGroupMember() *types.ValidIngredientGroupMember {
	return &types.ValidIngredientGroupMember{
		ID:              BuildFakeID(),
		ValidIngredient: *BuildFakeValidIngredient(),
		CreatedAt:       BuildFakeTime(),
		BelongsToGroup:  BuildFakeID(),
	}
}

// BuildFakeValidIngredientGroupsList builds a faked ValidIngredientGroupList.
func BuildFakeValidIngredientGroupsList() *filtering.QueryFilteredResult[types.ValidIngredientGroup] {
	var examples []*types.ValidIngredientGroup
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidIngredientGroup())
	}

	return &filtering.QueryFilteredResult[types.ValidIngredientGroup]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidIngredientGroupUpdateRequestInput builds a faked ValidIngredientGroupUpdateRequestInput from a valid ingredient group.
func BuildFakeValidIngredientGroupUpdateRequestInput() *types.ValidIngredientGroupUpdateRequestInput {
	validIngredient := BuildFakeValidIngredientGroup()
	return converters.ConvertValidIngredientGroupToValidIngredientGroupUpdateRequestInput(validIngredient)
}

// BuildFakeValidIngredientGroupCreationRequestInput builds a faked ValidIngredientGroupCreationRequestInput.
func BuildFakeValidIngredientGroupCreationRequestInput() *types.ValidIngredientGroupCreationRequestInput {
	validIngredient := BuildFakeValidIngredientGroup()
	return converters.ConvertValidIngredientGroupToValidIngredientGroupCreationRequestInput(validIngredient)
}
