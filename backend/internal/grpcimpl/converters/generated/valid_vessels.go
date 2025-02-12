package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidVesselCreationRequestInputToValidVessel(input *messages.ValidVesselCreationRequestInput) *messages.ValidVessel {

output := &messages.ValidVessel{
    IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
    UsableForStorage: input.UsableForStorage,
    Shape: input.Shape,
    Description: input.Description,
    Slug: input.Slug,
    Capacity: input.Capacity,
    HeightInMillimeters: input.HeightInMillimeters,
    LengthInMillimeters: input.LengthInMillimeters,
    DisplayInSummaryLists: input.DisplayInSummaryLists,
    Name: input.Name,
    IconPath: input.IconPath,
    PluralName: input.PluralName,
    WidthInMillimeters: input.WidthInMillimeters,
}

return output
}
func ConvertValidVesselUpdateRequestInputToValidVessel(input *messages.ValidVesselUpdateRequestInput) *messages.ValidVessel {

output := &messages.ValidVessel{
    WidthInMillimeters: input.WidthInMillimeters,
    DisplayInSummaryLists: input.DisplayInSummaryLists,
    UsableForStorage: input.UsableForStorage,
    Shape: input.Shape,
    Description: input.Description,
    Slug: input.Slug,
    IconPath: input.IconPath,
    LengthInMillimeters: input.LengthInMillimeters,
    IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
    Name: input.Name,
    PluralName: input.PluralName,
    HeightInMillimeters: input.HeightInMillimeters,
    Capacity: input.Capacity,
}

return output
}
