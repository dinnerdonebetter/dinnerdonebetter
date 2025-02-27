package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeValidPreparation builds a faked valid preparation.
func BuildFakeValidPreparation() *types.ValidPreparation {
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
		IngredientCount:             BuildFakeUint16RangeWithOptionalMax(),
		InstrumentCount:             BuildFakeUint16RangeWithOptionalMax(),
		TemperatureRequired:         fake.Bool(),
		TimeEstimateRequired:        fake.Bool(),
		ConditionExpressionRequired: fake.Bool(),
		ConsumesVessel:              fake.Bool(),
		OnlyForVessels:              fake.Bool(),
		VesselCount:                 BuildFakeUint16RangeWithOptionalMax(),
	}
}

// BuildFakeValidPreparationsList builds a faked ValidPreparationList.
func BuildFakeValidPreparationsList() *filtering.QueryFilteredResult[types.ValidPreparation] {
	var examples []*types.ValidPreparation
	for range exampleQuantity {
		examples = append(examples, BuildFakeValidPreparation())
	}

	return &filtering.QueryFilteredResult[types.ValidPreparation]{
		Pagination: filtering.Pagination{
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
