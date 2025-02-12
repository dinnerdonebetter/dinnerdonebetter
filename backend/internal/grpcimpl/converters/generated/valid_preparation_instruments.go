package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidPreparationInstrumentCreationRequestInputToValidPreparationInstrument(input *messages.ValidPreparationInstrumentCreationRequestInput) *messages.ValidPreparationInstrument {

output := &messages.ValidPreparationInstrument{
    Notes: input.Notes,
}

return output
}
func ConvertValidPreparationInstrumentUpdateRequestInputToValidPreparationInstrument(input *messages.ValidPreparationInstrumentUpdateRequestInput) *messages.ValidPreparationInstrument {

output := &messages.ValidPreparationInstrument{
    Notes: input.Notes,
}

return output
}
