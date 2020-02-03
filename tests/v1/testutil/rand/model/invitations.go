package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomInvitationCreationInput creates a random InvitationInput
func RandomInvitationCreationInput() *models.InvitationCreationInput {
	x := &models.InvitationCreationInput{
		Code:     fake.Word(),
		Consumed: fake.Bool(),
	}

	return x
}
