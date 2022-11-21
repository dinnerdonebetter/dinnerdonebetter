package converters

import (
	"github.com/prixfixeco/backend/pkg/types"
)

// ConvertHouseholdCreationInputToHouseholdDatabaseCreationInput creates a HouseholdDatabaseCreationInput from a HouseholdCreationRequestInput.
func ConvertHouseholdCreationInputToHouseholdDatabaseCreationInput(input *types.HouseholdCreationRequestInput) *types.HouseholdDatabaseCreationInput {
	x := &types.HouseholdDatabaseCreationInput{
		Name:         input.Name,
		ContactEmail: input.ContactEmail,
		ContactPhone: input.ContactPhone,
		TimeZone:     input.TimeZone,
	}

	return x
}

// ConvertHouseholdToHouseholdUpdateRequestInput creates a HouseholdUpdateRequestInput from a Household.
func ConvertHouseholdToHouseholdUpdateRequestInput(input *types.Household) *types.HouseholdUpdateRequestInput {
	x := &types.HouseholdUpdateRequestInput{
		Name:          &input.Name,
		ContactEmail:  &input.ContactEmail,
		ContactPhone:  &input.ContactPhone,
		BelongsToUser: input.BelongsToUser,
		TimeZone:      &input.TimeZone,
	}

	return x
}

// ConvertHouseholdToHouseholdCreationRequestInput builds a faked HouseholdCreationRequestInput from a household.
func ConvertHouseholdToHouseholdCreationRequestInput(household *types.Household) *types.HouseholdCreationRequestInput {
	return &types.HouseholdCreationRequestInput{
		Name:         household.Name,
		ContactEmail: household.ContactEmail,
		ContactPhone: household.ContactPhone,
		TimeZone:     household.TimeZone,
	}
}

// ConvertHouseholdToHouseholdDatabaseCreationInput builds a faked HouseholdCreationRequestInput.
func ConvertHouseholdToHouseholdDatabaseCreationInput(household *types.Household) *types.HouseholdDatabaseCreationInput {
	return &types.HouseholdDatabaseCreationInput{
		ID:            household.ID,
		Name:          household.Name,
		ContactEmail:  household.ContactEmail,
		ContactPhone:  household.ContactPhone,
		BelongsToUser: household.BelongsToUser,
		TimeZone:      household.TimeZone,
	}
}

// ConvertHouseholdUserMembershipCreationRequestInputToHouseholdUserMembershipDatabaseCreationInput builds a HouseholdUserMembershipDatabaseCreationInput from a HouseholdUserMembershipCreationRequestInput.
func ConvertHouseholdUserMembershipCreationRequestInputToHouseholdUserMembershipDatabaseCreationInput(input *types.HouseholdUserMembershipCreationRequestInput) *types.HouseholdUserMembershipDatabaseCreationInput {
	return &types.HouseholdUserMembershipDatabaseCreationInput{
		Reason: input.Reason,
		UserID: input.UserID,
	}
}
