package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidPreparationCreationRequestInputToValidPreparation(input *messages.ValidPreparationCreationRequestInput) *messages.ValidPreparation {

output := &messages.ValidPreparation{
    IngredientCount: input.IngredientCount,
    VesselCount: input.VesselCount,
    Name: input.Name,
    IconPath: input.IconPath,
    Description: input.Description,
    TemperatureRequired: input.TemperatureRequired,
    ConditionExpressionRequired: input.ConditionExpressionRequired,
    InstrumentCount: input.InstrumentCount,
    OnlyForVessels: input.OnlyForVessels,
    YieldsNothing: input.YieldsNothing,
    TimeEstimateRequired: input.TimeEstimateRequired,
    RestrictToIngredients: input.RestrictToIngredients,
    Slug: input.Slug,
    PastTense: input.PastTense,
    ConsumesVessel: input.ConsumesVessel,
}

return output
}
func ConvertValidPreparationUpdateRequestInputToValidPreparation(input *messages.ValidPreparationUpdateRequestInput) *messages.ValidPreparation {

output := &messages.ValidPreparation{
    InstrumentCount: ConvertUint16RangeWithOptionalMaxUpdateRequestInputToUint16RangeWithOptionalMax(input.InstrumentCount),
    IconPath: input.IconPath,
    Description: input.Description,
    VesselCount: ConvertUint16RangeWithOptionalMaxUpdateRequestInputToUint16RangeWithOptionalMax(input.VesselCount),
    Name: input.Name,
    ConsumesVessel: input.ConsumesVessel,
    IngredientCount: ConvertUint16RangeWithOptionalMaxUpdateRequestInputToUint16RangeWithOptionalMax(input.IngredientCount),
    OnlyForVessels: input.OnlyForVessels,
    Slug: input.Slug,
    PastTense: input.PastTense,
    TemperatureRequired: input.TemperatureRequired,
    ConditionExpressionRequired: input.ConditionExpressionRequired,
    YieldsNothing: input.YieldsNothing,
    TimeEstimateRequired: input.TimeEstimateRequired,
    RestrictToIngredients: input.RestrictToIngredients,
}

return output
}
