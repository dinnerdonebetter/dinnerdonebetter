package fakes

import (
	recipeenumfakes "github.com/dinnerdonebetter/backend/internal/domain/recipeenums/fakes"
	types "github.com/dinnerdonebetter/backend/internal/domain/recipes"
	"github.com/dinnerdonebetter/backend/internal/domain/recipes/converters"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
)

// BuildFakeAccountInstrumentOwnership builds a faked valid ingredient.
func BuildFakeAccountInstrumentOwnership() *types.AccountInstrumentOwnership {
	return &types.AccountInstrumentOwnership{
		CreatedAt:        BuildFakeTime(),
		ID:               buildUniqueString(),
		Notes:            buildUniqueString(),
		BelongsToAccount: buildUniqueString(),
		Instrument:       *recipeenumfakes.BuildFakeValidInstrument(),
		Quantity:         uint16(buildFakeNumber()),
	}
}

// BuildFakeAccountInstrumentOwnershipsList builds a faked AccountInstrumentOwnershipList.
func BuildFakeAccountInstrumentOwnershipsList() *filtering.QueryFilteredResult[types.AccountInstrumentOwnership] {
	var examples []*types.AccountInstrumentOwnership
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeAccountInstrumentOwnership())
	}

	return &filtering.QueryFilteredResult[types.AccountInstrumentOwnership]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeAccountInstrumentOwnershipUpdateRequestInput builds a faked AccountInstrumentOwnershipUpdateRequestInput from a valid ingredient.
func BuildFakeAccountInstrumentOwnershipUpdateRequestInput() *types.AccountInstrumentOwnershipUpdateRequestInput {
	validIngredient := BuildFakeAccountInstrumentOwnership()
	return converters.ConvertAccountInstrumentOwnershipToAccountInstrumentOwnershipUpdateRequestInput(validIngredient)
}

// BuildFakeAccountInstrumentOwnershipCreationRequestInput builds a faked AccountInstrumentOwnershipCreationRequestInput.
func BuildFakeAccountInstrumentOwnershipCreationRequestInput() *types.AccountInstrumentOwnershipCreationRequestInput {
	validIngredient := BuildFakeAccountInstrumentOwnership()
	return converters.ConvertAccountInstrumentOwnershipToAccountInstrumentOwnershipCreationRequestInput(validIngredient)
}
