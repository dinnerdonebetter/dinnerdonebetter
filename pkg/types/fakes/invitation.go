package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeInvitation builds a faked invitation.
func BuildFakeInvitation() *types.Invitation {
	return &types.Invitation{
		ID:                 uint64(fake.Uint32()),
		ExternalID:         fake.UUID(),
		Code:               fake.Word(),
		Consumed:           fake.Bool(),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
		BelongsToHousehold: fake.Uint64(),
	}
}

// BuildFakeInvitationList builds a faked InvitationList.
func BuildFakeInvitationList() *types.InvitationList {
	var examples []*types.Invitation
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeInvitation())
	}

	return &types.InvitationList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Invitations: examples,
	}
}

// BuildFakeInvitationUpdateInput builds a faked InvitationUpdateInput from an invitation.
func BuildFakeInvitationUpdateInput() *types.InvitationUpdateInput {
	invitation := BuildFakeInvitation()
	return &types.InvitationUpdateInput{
		Code:               invitation.Code,
		Consumed:           invitation.Consumed,
		BelongsToHousehold: invitation.BelongsToHousehold,
	}
}

// BuildFakeInvitationUpdateInputFromInvitation builds a faked InvitationUpdateInput from an invitation.
func BuildFakeInvitationUpdateInputFromInvitation(invitation *types.Invitation) *types.InvitationUpdateInput {
	return &types.InvitationUpdateInput{
		Code:               invitation.Code,
		Consumed:           invitation.Consumed,
		BelongsToHousehold: invitation.BelongsToHousehold,
	}
}

// BuildFakeInvitationCreationInput builds a faked InvitationCreationInput.
func BuildFakeInvitationCreationInput() *types.InvitationCreationInput {
	invitation := BuildFakeInvitation()
	return BuildFakeInvitationCreationInputFromInvitation(invitation)
}

// BuildFakeInvitationCreationInputFromInvitation builds a faked InvitationCreationInput from an invitation.
func BuildFakeInvitationCreationInputFromInvitation(invitation *types.Invitation) *types.InvitationCreationInput {
	return &types.InvitationCreationInput{
		Code:               invitation.Code,
		Consumed:           invitation.Consumed,
		BelongsToHousehold: invitation.BelongsToHousehold,
	}
}
