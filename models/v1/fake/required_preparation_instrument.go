package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRequiredPreparationInstrument builds a faked required preparation instrument.
func BuildFakeRequiredPreparationInstrument() *models.RequiredPreparationInstrument {
	return &models.RequiredPreparationInstrument{
		ID:                        fake.Uint64(),
		ValidInstrumentID:         uint64(fake.Uint32()),
		Notes:                     fake.Word(),
		CreatedOn:                 uint64(uint32(fake.Date().Unix())),
		BelongsToValidPreparation: fake.Uint64(),
	}
}

// BuildFakeRequiredPreparationInstrumentList builds a faked RequiredPreparationInstrumentList.
func BuildFakeRequiredPreparationInstrumentList() *models.RequiredPreparationInstrumentList {
	exampleRequiredPreparationInstrument1 := BuildFakeRequiredPreparationInstrument()
	exampleRequiredPreparationInstrument2 := BuildFakeRequiredPreparationInstrument()
	exampleRequiredPreparationInstrument3 := BuildFakeRequiredPreparationInstrument()

	return &models.RequiredPreparationInstrumentList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		RequiredPreparationInstruments: []models.RequiredPreparationInstrument{
			*exampleRequiredPreparationInstrument1,
			*exampleRequiredPreparationInstrument2,
			*exampleRequiredPreparationInstrument3,
		},
	}
}

// BuildFakeRequiredPreparationInstrumentUpdateInputFromRequiredPreparationInstrument builds a faked RequiredPreparationInstrumentUpdateInput from a required preparation instrument.
func BuildFakeRequiredPreparationInstrumentUpdateInputFromRequiredPreparationInstrument(requiredPreparationInstrument *models.RequiredPreparationInstrument) *models.RequiredPreparationInstrumentUpdateInput {
	return &models.RequiredPreparationInstrumentUpdateInput{
		ValidInstrumentID:         requiredPreparationInstrument.ValidInstrumentID,
		Notes:                     requiredPreparationInstrument.Notes,
		BelongsToValidPreparation: requiredPreparationInstrument.BelongsToValidPreparation,
	}
}

// BuildFakeRequiredPreparationInstrumentCreationInput builds a faked RequiredPreparationInstrumentCreationInput.
func BuildFakeRequiredPreparationInstrumentCreationInput() *models.RequiredPreparationInstrumentCreationInput {
	requiredPreparationInstrument := BuildFakeRequiredPreparationInstrument()
	return BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(requiredPreparationInstrument)
}

// BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument builds a faked RequiredPreparationInstrumentCreationInput from a required preparation instrument.
func BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(requiredPreparationInstrument *models.RequiredPreparationInstrument) *models.RequiredPreparationInstrumentCreationInput {
	return &models.RequiredPreparationInstrumentCreationInput{
		ValidInstrumentID:         requiredPreparationInstrument.ValidInstrumentID,
		Notes:                     requiredPreparationInstrument.Notes,
		BelongsToValidPreparation: requiredPreparationInstrument.BelongsToValidPreparation,
	}
}
