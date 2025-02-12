package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertValidPreparationVesselCreationRequestInputToValidPreparationVessel(input *messages.ValidPreparationVesselCreationRequestInput) *messages.ValidPreparationVessel {

output := &messages.ValidPreparationVessel{
    Notes: input.Notes,
}

return output
}
func ConvertValidPreparationVesselUpdateRequestInputToValidPreparationVessel(input *messages.ValidPreparationVesselUpdateRequestInput) *messages.ValidPreparationVessel {

output := &messages.ValidPreparationVessel{
    Notes: input.Notes,
}

return output
}
