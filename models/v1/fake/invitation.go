package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeInvitation builds a faked invitation.
func BuildFakeInvitation() *models.Invitation {
	return &models.Invitation{
		ID:            fake.Uint64(),
		Code:          fake.Word(),
		Consumed:      fake.Bool(),
		CreatedOn:     uint64(uint32(fake.Date().Unix())),
		BelongsToUser: fake.Uint64(),
	}
}

// BuildFakeInvitationList builds a faked InvitationList.
func BuildFakeInvitationList() *models.InvitationList {
	exampleInvitation1 := BuildFakeInvitation()
	exampleInvitation2 := BuildFakeInvitation()
	exampleInvitation3 := BuildFakeInvitation()

	return &models.InvitationList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		Invitations: []models.Invitation{
			*exampleInvitation1,
			*exampleInvitation2,
			*exampleInvitation3,
		},
	}
}

// BuildFakeInvitationUpdateInputFromInvitation builds a faked InvitationUpdateInput from an invitation.
func BuildFakeInvitationUpdateInputFromInvitation(invitation *models.Invitation) *models.InvitationUpdateInput {
	return &models.InvitationUpdateInput{
		Code:          invitation.Code,
		Consumed:      invitation.Consumed,
		BelongsToUser: invitation.BelongsToUser,
	}
}

// BuildFakeInvitationCreationInput builds a faked InvitationCreationInput.
func BuildFakeInvitationCreationInput() *models.InvitationCreationInput {
	invitation := BuildFakeInvitation()
	return BuildFakeInvitationCreationInputFromInvitation(invitation)
}

// BuildFakeInvitationCreationInputFromInvitation builds a faked InvitationCreationInput from an invitation.
func BuildFakeInvitationCreationInputFromInvitation(invitation *models.Invitation) *models.InvitationCreationInput {
	return &models.InvitationCreationInput{
		Code:          invitation.Code,
		Consumed:      invitation.Consumed,
		BelongsToUser: invitation.BelongsToUser,
	}
}
