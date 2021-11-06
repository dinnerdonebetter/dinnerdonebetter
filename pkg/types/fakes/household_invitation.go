package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeHouseholdInvitation builds a faked HouseholdInvitation.
func BuildFakeHouseholdInvitation() *types.HouseholdInvitation {
	return &types.HouseholdInvitation{
		FromUser:             fake.LoremIpsumSentence(exampleQuantity),
		ToEmail:              fake.LoremIpsumSentence(exampleQuantity),
		ToUser:               func(s string) *string { return &s }(fake.LoremIpsumSentence(exampleQuantity)),
		Note:                 fake.LoremIpsumSentence(exampleQuantity),
		Token:                fake.UUID(),
		DestinationHousehold: fake.LoremIpsumSentence(exampleQuantity),
		ID:                   fake.LoremIpsumSentence(exampleQuantity),
		Status:               types.PendingHouseholdInvitationStatus,
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
			FilteredCount: exampleQuantity,
			TotalCount:    exampleQuantity * 2,
		},
		HouseholdInvitations: examples,
	}
}

// BuildFakeHouseholdInvitationCreationInput builds a faked HouseholdInvitationCreationRequestInput from a webhook.
func BuildFakeHouseholdInvitationCreationInput() *types.HouseholdInvitationCreationRequestInput {
	webhook := BuildFakeHouseholdInvitation()
	return BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation(webhook)
}

// BuildFakeHouseholdInvitationDatabaseCreationInput builds a faked HouseholdInvitationCreationRequestInput from a webhook.
func BuildFakeHouseholdInvitationDatabaseCreationInput() *types.HouseholdInvitationDatabaseCreationInput {
	webhook := BuildFakeHouseholdInvitation()
	return BuildFakeHouseholdInvitationDatabaseCreationInputFromHouseholdInvitation(webhook)
}

// BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation builds a faked HouseholdInvitationCreationRequestInput.
func BuildFakeHouseholdInvitationCreationInputFromHouseholdInvitation(householdInvitation *types.HouseholdInvitation) *types.HouseholdInvitationCreationRequestInput {
	return &types.HouseholdInvitationCreationRequestInput{
		ID:                   householdInvitation.ID,
		FromUser:             householdInvitation.FromUser,
		Note:                 householdInvitation.Note,
		ToEmail:              householdInvitation.ToEmail,
		DestinationHousehold: householdInvitation.DestinationHousehold,
	}
}

// BuildFakeHouseholdInvitationDatabaseCreationInputFromHouseholdInvitation builds a faked HouseholdInvitationCreationRequestInput.
func BuildFakeHouseholdInvitationDatabaseCreationInputFromHouseholdInvitation(householdInvitation *types.HouseholdInvitation) *types.HouseholdInvitationDatabaseCreationInput {
	return &types.HouseholdInvitationDatabaseCreationInput{
		ID:                   householdInvitation.ID,
		FromUser:             householdInvitation.FromUser,
		ToUser:               householdInvitation.ToUser,
		Note:                 householdInvitation.Note,
		ToEmail:              householdInvitation.ToEmail,
		Token:                householdInvitation.Token,
		DestinationHousehold: householdInvitation.DestinationHousehold,
	}
}
