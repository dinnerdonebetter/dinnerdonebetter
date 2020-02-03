package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomRequiredPreparationInstrumentCreationInput creates a random RequiredPreparationInstrumentInput
func RandomRequiredPreparationInstrumentCreationInput() *models.RequiredPreparationInstrumentCreationInput {
	x := &models.RequiredPreparationInstrumentCreationInput{
		InstrumentID:  fake.Uint64(),
		PreparationID: fake.Uint64(),
		Notes:         fake.Word(),
	}

	return x
}
