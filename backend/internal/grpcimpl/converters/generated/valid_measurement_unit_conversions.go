package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidMeasurementUnitConversionCreationRequestInputToValidMeasurementUnitConversion(input *messages.ValidMeasurementUnitConversionCreationRequestInput) *messages.ValidMeasurementUnitConversion {

output := &messages.ValidMeasurementUnitConversion{
    To: input.To,
    Notes: input.Notes,
    Modifier: input.Modifier,
    OnlyForIngredient: input.OnlyForIngredient,
    From: input.From,
}

return output
}
func ConvertValidMeasurementUnitConversionUpdateRequestInputToValidMeasurementUnitConversion(input *messages.ValidMeasurementUnitConversionUpdateRequestInput) *messages.ValidMeasurementUnitConversion {

output := &messages.ValidMeasurementUnitConversion{
    OnlyForIngredient: input.OnlyForIngredient,
    From: input.From,
    To: input.To,
    Notes: input.Notes,
    Modifier: input.Modifier,
}

return output
}
