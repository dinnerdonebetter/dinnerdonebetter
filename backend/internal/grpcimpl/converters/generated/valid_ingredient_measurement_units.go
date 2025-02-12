package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidIngredientMeasurementUnitCreationRequestInputToValidIngredientMeasurementUnit(input *messages.ValidIngredientMeasurementUnitCreationRequestInput) *messages.ValidIngredientMeasurementUnit {

output := &messages.ValidIngredientMeasurementUnit{
    Notes: input.Notes,
    AllowableQuantity: input.AllowableQuantity,
}

return output
}
func ConvertValidIngredientMeasurementUnitUpdateRequestInputToValidIngredientMeasurementUnit(input *messages.ValidIngredientMeasurementUnitUpdateRequestInput) *messages.ValidIngredientMeasurementUnit {

output := &messages.ValidIngredientMeasurementUnit{
    Notes: input.Notes,
    AllowableQuantity: ConvertFloat32RangeWithOptionalMaxUpdateRequestInputToFloat32RangeWithOptionalMax(input.AllowableQuantity),
}

return output
}
