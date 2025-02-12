package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertHouseholdCreationRequestInputToHousehold(input *messages.HouseholdCreationRequestInput) *messages.Household {

output := &messages.Household{
    City: input.City,
    Country: input.Country,
    AddressLine2: input.AddressLine2,
    Latitude: input.Latitude,
    State: input.State,
    ContactPhone: input.ContactPhone,
    AddressLine1: input.AddressLine1,
    ZipCode: input.ZipCode,
    Name: input.Name,
    Longitude: input.Longitude,
}

return output
}
func ConvertHouseholdUpdateRequestInputToHousehold(input *messages.HouseholdUpdateRequestInput) *messages.Household {

output := &messages.Household{
    Longitude: input.Longitude,
    Latitude: input.Latitude,
    BelongsToUser: input.BelongsToUser,
    ContactPhone: input.ContactPhone,
    City: input.City,
    AddressLine1: input.AddressLine1,
    ZipCode: input.ZipCode,
    Country: input.Country,
    Name: input.Name,
    AddressLine2: input.AddressLine2,
    State: input.State,
}

return output
}
