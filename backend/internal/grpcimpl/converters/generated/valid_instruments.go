package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidInstrumentCreationRequestInputToValidInstrument(input *messages.ValidInstrumentCreationRequestInput) *messages.ValidInstrument {

output := &messages.ValidInstrument{
    PluralName: input.PluralName,
    Description: input.Description,
    Slug: input.Slug,
    DisplayInSummaryLists: input.DisplayInSummaryLists,
    IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
    UsableForStorage: input.UsableForStorage,
    Name: input.Name,
    IconPath: input.IconPath,
}

return output
}
func ConvertValidInstrumentUpdateRequestInputToValidInstrument(input *messages.ValidInstrumentUpdateRequestInput) *messages.ValidInstrument {

output := &messages.ValidInstrument{
    Slug: input.Slug,
    DisplayInSummaryLists: input.DisplayInSummaryLists,
    IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
    UsableForStorage: input.UsableForStorage,
    Name: input.Name,
    IconPath: input.IconPath,
    PluralName: input.PluralName,
    Description: input.Description,
}

return output
}
