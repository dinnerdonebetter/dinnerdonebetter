package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v6"
)

// BuildFakeValidPreparation builds a faked valid preparation.
func BuildFakeValidPreparation() *types.ValidPreparation {
	minIngredients := BuildFakeNumber()
	minInstruments := BuildFakeNumber()
	minVessels := BuildFakeNumber()

	return &types.ValidPreparation{
		ID:                          BuildFakeID(),
		Name:                        buildUniqueString(),
		Description:                 buildUniqueString(),
		IconPath:                    buildUniqueString(),
		YieldsNothing:               fake.Bool(),
		RestrictToIngredients:       fake.Bool(),
		Slug:                        buildUniqueString(),
		PastTense:                   buildUniqueString(),
		CreatedAt:                   BuildFakeTime(),
		MinimumIngredientCount:      int32(minIngredients),
		MaximumIngredientCount:      pointers.Pointer(int32(minIngredients + 1)),
		MinimumInstrumentCount:      int32(minInstruments),
		MaximumInstrumentCount:      pointers.Pointer(int32(minInstruments + 1)),
		TemperatureRequired:         fake.Bool(),
		TimeEstimateRequired:        fake.Bool(),
		ConditionExpressionRequired: fake.Bool(),
		ConsumesVessel:              fake.Bool(),
		OnlyForVessels:              fake.Bool(),
		MinimumVesselCount:          int32(minVessels),
		MaximumVesselCount:          pointers.Pointer(int32(minVessels + 1)),
	}
}

// BuildFakeValidPreparationList builds a faked ValidPreparationList.
func BuildFakeValidPreparationList() *types.QueryFilteredResult[types.ValidPreparation] {
	var examples []*types.ValidPreparation
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidPreparation())
	}

	return &types.QueryFilteredResult[types.ValidPreparation]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidPreparationUpdateRequestInput builds a faked ValidPreparationUpdateRequestInput from a valid preparation.
func BuildFakeValidPreparationUpdateRequestInput() *types.ValidPreparationUpdateRequestInput {
	validPreparation := BuildFakeValidPreparation()
	return converters.ConvertValidPreparationToValidPreparationUpdateRequestInput(validPreparation)
}

// BuildFakeValidPreparationCreationRequestInput builds a faked ValidPreparationCreationRequestInput.
func BuildFakeValidPreparationCreationRequestInput() *types.ValidPreparationCreationRequestInput {
	validPreparation := BuildFakeValidPreparation()
	return converters.ConvertValidPreparationToValidPreparationCreationRequestInput(validPreparation)
}
