package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidMeasurementUnitCreationRequestInputToValidMeasurementUnit(input *messages.ValidMeasurementUnitCreationRequestInput) *messages.ValidMeasurementUnit {

output := &messages.ValidMeasurementUnit{
    Name: input.Name,
    Slug: input.Slug,
    Universal: input.Universal,
    IconPath: input.IconPath,
    Description: input.Description,
    Volumetric: input.Volumetric,
    Metric: input.Metric,
    Imperial: input.Imperial,
    PluralName: input.PluralName,
}

return output
}
func ConvertValidMeasurementUnitUpdateRequestInputToValidMeasurementUnit(input *messages.ValidMeasurementUnitUpdateRequestInput) *messages.ValidMeasurementUnit {

output := &messages.ValidMeasurementUnit{
    PluralName: input.PluralName,
    Description: input.Description,
    Metric: input.Metric,
    Imperial: input.Imperial,
    IconPath: input.IconPath,
    Name: input.Name,
    Slug: input.Slug,
    Volumetric: input.Volumetric,
    Universal: input.Universal,
}

return output
}
