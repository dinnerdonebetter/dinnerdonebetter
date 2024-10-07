package converters

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertValidPreparationToValidPreparationUpdateRequestInput creates a ValidPreparationUpdateRequestInput from a ValidPreparation.
func ConvertValidPreparationToValidPreparationUpdateRequestInput(input *types.ValidPreparation) *types.ValidPreparationUpdateRequestInput {
	x := &types.ValidPreparationUpdateRequestInput{
		Name:                  &input.Name,
		Description:           &input.Description,
		IconPath:              &input.IconPath,
		YieldsNothing:         &input.YieldsNothing,
		RestrictToIngredients: &input.RestrictToIngredients,
		Slug:                  &input.Slug,
		PastTense:             &input.PastTense,
		IngredientCount: types.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Max: input.IngredientCount.Max,
			Min: &input.IngredientCount.Min,
		},
		InstrumentCount: types.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Max: input.InstrumentCount.Max,
			Min: &input.InstrumentCount.Min,
		},
		VesselCount: types.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Max: input.VesselCount.Max,
			Min: &input.VesselCount.Min,
		},
		TemperatureRequired:         &input.TemperatureRequired,
		TimeEstimateRequired:        &input.TimeEstimateRequired,
		ConditionExpressionRequired: &input.ConditionExpressionRequired,
		ConsumesVessel:              &input.ConsumesVessel,
		OnlyForVessels:              &input.OnlyForVessels,
	}

	return x
}

// ConvertValidPreparationCreationRequestInputToValidPreparationDatabaseCreationInput creates a ValidPreparationDatabaseCreationInput from a ValidPreparationCreationRequestInput.
func ConvertValidPreparationCreationRequestInputToValidPreparationDatabaseCreationInput(input *types.ValidPreparationCreationRequestInput) *types.ValidPreparationDatabaseCreationInput {
	x := &types.ValidPreparationDatabaseCreationInput{
		ID:                    identifiers.New(),
		Name:                  input.Name,
		Description:           input.Description,
		IconPath:              input.IconPath,
		YieldsNothing:         input.YieldsNothing,
		RestrictToIngredients: input.RestrictToIngredients,
		Slug:                  input.Slug,
		PastTense:             input.PastTense,
		IngredientCount: types.Uint16RangeWithOptionalMax{
			Max: input.IngredientCount.Max,
			Min: input.IngredientCount.Min,
		},
		InstrumentCount: types.Uint16RangeWithOptionalMax{
			Max: input.InstrumentCount.Max,
			Min: input.InstrumentCount.Min,
		},
		VesselCount: types.Uint16RangeWithOptionalMax{
			Max: input.VesselCount.Max,
			Min: input.VesselCount.Min,
		},
		TemperatureRequired:         input.TemperatureRequired,
		TimeEstimateRequired:        input.TimeEstimateRequired,
		ConditionExpressionRequired: input.ConditionExpressionRequired,
		ConsumesVessel:              input.ConsumesVessel,
		OnlyForVessels:              input.OnlyForVessels,
	}

	return x
}

// ConvertValidPreparationToValidPreparationCreationRequestInput builds a ValidPreparationCreationRequestInput from a ValidPreparation.
func ConvertValidPreparationToValidPreparationCreationRequestInput(validPreparation *types.ValidPreparation) *types.ValidPreparationCreationRequestInput {
	return &types.ValidPreparationCreationRequestInput{
		Name:                  validPreparation.Name,
		Description:           validPreparation.Description,
		IconPath:              validPreparation.IconPath,
		YieldsNothing:         validPreparation.YieldsNothing,
		RestrictToIngredients: validPreparation.RestrictToIngredients,
		Slug:                  validPreparation.Slug,
		PastTense:             validPreparation.PastTense,
		IngredientCount: types.Uint16RangeWithOptionalMax{
			Max: validPreparation.IngredientCount.Max,
			Min: validPreparation.IngredientCount.Min,
		},
		InstrumentCount: types.Uint16RangeWithOptionalMax{
			Max: validPreparation.InstrumentCount.Max,
			Min: validPreparation.InstrumentCount.Min,
		},
		VesselCount: types.Uint16RangeWithOptionalMax{
			Max: validPreparation.VesselCount.Max,
			Min: validPreparation.VesselCount.Min,
		},
		TemperatureRequired:         validPreparation.TemperatureRequired,
		TimeEstimateRequired:        validPreparation.TimeEstimateRequired,
		ConditionExpressionRequired: validPreparation.ConditionExpressionRequired,
		ConsumesVessel:              validPreparation.ConsumesVessel,
		OnlyForVessels:              validPreparation.OnlyForVessels,
	}
}

// ConvertValidPreparationToValidPreparationDatabaseCreationInput builds a ValidPreparationDatabaseCreationInput from a ValidPreparation.
func ConvertValidPreparationToValidPreparationDatabaseCreationInput(validPreparation *types.ValidPreparation) *types.ValidPreparationDatabaseCreationInput {
	return &types.ValidPreparationDatabaseCreationInput{
		ID:                    validPreparation.ID,
		Name:                  validPreparation.Name,
		Description:           validPreparation.Description,
		IconPath:              validPreparation.IconPath,
		YieldsNothing:         validPreparation.YieldsNothing,
		RestrictToIngredients: validPreparation.RestrictToIngredients,
		Slug:                  validPreparation.Slug,
		PastTense:             validPreparation.PastTense,
		IngredientCount: types.Uint16RangeWithOptionalMax{
			Max: validPreparation.IngredientCount.Max,
			Min: validPreparation.IngredientCount.Min,
		},
		InstrumentCount: types.Uint16RangeWithOptionalMax{
			Max: validPreparation.InstrumentCount.Max,
			Min: validPreparation.InstrumentCount.Min,
		},
		VesselCount: types.Uint16RangeWithOptionalMax{
			Max: validPreparation.VesselCount.Max,
			Min: validPreparation.VesselCount.Min,
		},
		TemperatureRequired:         validPreparation.TemperatureRequired,
		TimeEstimateRequired:        validPreparation.TimeEstimateRequired,
		ConditionExpressionRequired: validPreparation.ConditionExpressionRequired,
		ConsumesVessel:              validPreparation.ConsumesVessel,
		OnlyForVessels:              validPreparation.OnlyForVessels,
	}
}

// ConvertValidPreparationToValidPreparationSearchSubset converts a ValidPreparation to a ValidPreparationSearchSubset.
func ConvertValidPreparationToValidPreparationSearchSubset(x *types.ValidPreparation) *types.ValidPreparationSearchSubset {
	return &types.ValidPreparationSearchSubset{
		ID:          x.ID,
		Name:        x.Name,
		PastTense:   x.PastTense,
		Description: x.Description,
	}
}
