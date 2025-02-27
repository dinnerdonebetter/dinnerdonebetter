package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"
)

// BuildFakeInstrumentOwnership builds a faked valid ingredient.
func BuildFakeInstrumentOwnership() *types.InstrumentOwnership {
	return &types.InstrumentOwnership{
		CreatedAt:          BuildFakeTime(),
		ID:                 buildUniqueString(),
		Notes:              buildUniqueString(),
		BelongsToHousehold: buildUniqueString(),
		Instrument:         *BuildFakeValidInstrument(),
		Quantity:           uint16(buildFakeNumber()),
	}
}

// BuildFakeInstrumentOwnershipsList builds a faked InstrumentOwnershipList.
func BuildFakeInstrumentOwnershipsList() *filtering.QueryFilteredResult[types.InstrumentOwnership] {
	var examples []*types.InstrumentOwnership
	for range exampleQuantity {
		examples = append(examples, BuildFakeInstrumentOwnership())
	}

	return &filtering.QueryFilteredResult[types.InstrumentOwnership]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeInstrumentOwnershipUpdateRequestInput builds a faked InstrumentOwnershipUpdateRequestInput from a valid ingredient.
func BuildFakeInstrumentOwnershipUpdateRequestInput() *types.InstrumentOwnershipUpdateRequestInput {
	validIngredient := BuildFakeInstrumentOwnership()
	return converters.ConvertInstrumentOwnershipToInstrumentOwnershipUpdateRequestInput(validIngredient)
}

// BuildFakeInstrumentOwnershipCreationRequestInput builds a faked InstrumentOwnershipCreationRequestInput.
func BuildFakeInstrumentOwnershipCreationRequestInput() *types.InstrumentOwnershipCreationRequestInput {
	validIngredient := BuildFakeInstrumentOwnership()
	return converters.ConvertInstrumentOwnershipToInstrumentOwnershipCreationRequestInput(validIngredient)
}
