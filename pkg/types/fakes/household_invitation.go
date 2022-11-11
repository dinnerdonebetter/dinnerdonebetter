package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeHouseholdInvitation builds a faked HouseholdInvitation.
func BuildFakeHouseholdInvitation() *types.HouseholdInvitation {
	return &types.HouseholdInvitation{
		FromUser:             *BuildFakeUser(),
		ToEmail:              buildUniqueString(),
		ToUser:               func(s string) *string { return &s }(buildUniqueString()),
		Note:                 buildUniqueString(),
		StatusNote:           buildUniqueString(),
		Token:                fake.UUID(),
		DestinationHousehold: *BuildFakeHousehold(),
		ID:                   buildUniqueString(),
		ExpiresAt:            BuildFakeTime(),
		Status:               types.PendingHouseholdInvitationStatus,
		CreatedAt:            BuildFakeTime(),
	}
}

// BuildFakeHouseholdInvitationList builds a faked HouseholdInvitationList.
func BuildFakeHouseholdInvitationList() *types.HouseholdInvitationList {
	var examples []*types.HouseholdInvitation
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeHouseholdInvitation())
	}

	return &types.HouseholdInvitationList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		HouseholdInvitations: examples,
	}
}

// BuildFakeHouseholdInvitationCreationRequestInput builds a faked HouseholdInvitationCreationRequestInput from a webhook.
func BuildFakeHouseholdInvitationCreationRequestInput() *types.HouseholdInvitationCreationRequestInput {
	invitation := BuildFakeHouseholdInvitation()
	return converters.ConvertHouseholdInvitationToHouseholdInvitationCreationInput(invitation)
}
