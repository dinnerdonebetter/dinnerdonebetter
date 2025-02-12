package converters

import (
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
)

func ConvertHouseholdInvitationCreationRequestInputToHouseholdInvitation(input *messages.HouseholdInvitationCreationRequestInput) *messages.HouseholdInvitation {

output := &messages.HouseholdInvitation{
    ToEmail: input.ToEmail,
    Note: input.Note,
    ToName: input.ToName,
    ExpiresAt: input.ExpiresAt,
}

return output
}
func ConvertHouseholdInvitationUpdateRequestInputToHouseholdInvitation(input *messages.HouseholdInvitationUpdateRequestInput) *messages.HouseholdInvitation {

output := &messages.HouseholdInvitation{
    Token: input.Token,
    Note: input.Note,
}

return output
}
