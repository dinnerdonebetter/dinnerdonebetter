package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeHouseholdInvitation builds a faked household.
func BuildFakeHouseholdInvitation() *types.HouseholdInvitation {
	return &types.HouseholdInvitation{
		ID:                   ksuid.New().String(),
		Status:               types.PendingHouseholdBillingStatus,
		FromUser:             fake.LoremIpsumSentence(exampleQuantity),
		ToUser:               fake.LoremIpsumSentence(exampleQuantity),
		DestinationHousehold: fake.LoremIpsumSentence(exampleQuantity),
		CreatedOn:            uint64(uint32(fake.Date().Unix())),
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

// BuildFakeHouseholdInvitationCreationInput builds a faked HouseholdInvitationCreationInput.
func BuildFakeHouseholdInvitationCreationInput() *types.HouseholdInvitationCreationInput {
	householdInvitation := BuildFakeHouseholdInvitation()
	return BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation(householdInvitation)
}

// BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation builds a faked HouseholdInvitationCreationInput from a household.
func BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation(householdInvitation *types.HouseholdInvitation) *types.HouseholdInvitationCreationInput {
	return &types.HouseholdInvitationCreationInput{
		ID:                   ksuid.New().String(),
		Status:               householdInvitation.Status,
		ToUser:               householdInvitation.ToUser,
		FromUser:             householdInvitation.FromUser,
		DestinationHousehold: householdInvitation.DestinationHousehold,
	}
}
