package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertHouseholdInstrumentOwnershipCreationRequestInputToHouseholdInstrumentOwnership(input *messages.HouseholdInstrumentOwnershipCreationRequestInput) *messages.HouseholdInstrumentOwnership {

output := &messages.HouseholdInstrumentOwnership{
    Notes: input.Notes,
    BelongsToHousehold: input.BelongsToHousehold,
    Quantity: input.Quantity,
}

return output
}
func ConvertHouseholdInstrumentOwnershipUpdateRequestInputToHouseholdInstrumentOwnership(input *messages.HouseholdInstrumentOwnershipUpdateRequestInput) *messages.HouseholdInstrumentOwnership {

output := &messages.HouseholdInstrumentOwnership{
    Notes: input.Notes,
    Quantity: input.Quantity,
}

return output
}
