package converters

import (
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertHouseholdInvitationCreationInputToHouseholdInvitationDatabaseCreationInput creates a HouseholdInvitationDatabaseCreationInput from a HouseholdInvitationCreationRequestInput.
func ConvertHouseholdInvitationCreationInputToHouseholdInvitationDatabaseCreationInput(input *types.HouseholdInvitationCreationRequestInput) *types.HouseholdInvitationDatabaseCreationInput {
	x := &types.HouseholdInvitationDatabaseCreationInput{
		ID:                     input.ID,
		FromUser:               input.FromUser,
		ToEmail:                input.ToEmail,
		DestinationHouseholdID: input.DestinationHouseholdID,
	}

	return x
}

// ConvertHouseholdInvitationToHouseholdInvitationCreationInput builds a faked HouseholdInvitationCreationRequestInput.
func ConvertHouseholdInvitationToHouseholdInvitationCreationInput(householdInvitation *types.HouseholdInvitation) *types.HouseholdInvitationCreationRequestInput {
	return &types.HouseholdInvitationCreationRequestInput{
		ID:                     householdInvitation.ID,
		FromUser:               householdInvitation.FromUser.ID,
		Note:                   householdInvitation.Note,
		ToEmail:                householdInvitation.ToEmail,
		DestinationHouseholdID: householdInvitation.DestinationHousehold.ID,
	}
}

// ConvertHouseholdInvitationToHouseholdInvitationUpdateInput builds a faked HouseholdInvitationUpdateRequestInput.
func ConvertHouseholdInvitationToHouseholdInvitationUpdateInput(householdInvitation *types.HouseholdInvitation) *types.HouseholdInvitationUpdateRequestInput {
	return &types.HouseholdInvitationUpdateRequestInput{
		Token: householdInvitation.Token,
		Note:  householdInvitation.Note,
	}
}

// ConvertHouseholdInvitationToHouseholdInvitationDatabaseCreationInput builds a faked HouseholdInvitationCreationRequestInput.
func ConvertHouseholdInvitationToHouseholdInvitationDatabaseCreationInput(householdInvitation *types.HouseholdInvitation) *types.HouseholdInvitationDatabaseCreationInput {
	return &types.HouseholdInvitationDatabaseCreationInput{
		ID:                     householdInvitation.ID,
		FromUser:               householdInvitation.FromUser.ID,
		ToUser:                 householdInvitation.ToUser,
		Note:                   householdInvitation.Note,
		ToEmail:                householdInvitation.ToEmail,
		Token:                  householdInvitation.Token,
		DestinationHouseholdID: householdInvitation.DestinationHousehold.ID,
	}
}
