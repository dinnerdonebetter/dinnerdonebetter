package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomInstrumentCreationInput creates a random InstrumentInput
func RandomInstrumentCreationInput() *models.InstrumentCreationInput {
	x := &models.InstrumentCreationInput{
		Name:        fake.Word(),
		Variant:     fake.Word(),
		Description: fake.Word(),
		Icon:        fake.Word(),
	}

	return x
}
