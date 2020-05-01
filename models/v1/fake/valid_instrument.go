package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeValidInstrument builds a faked valid instrument.
func BuildFakeValidInstrument() *models.ValidInstrument {
	return &models.ValidInstrument{
		ID:          fake.Uint64(),
		Name:        fake.Word(),
		Variant:     fake.Word(),
		Description: fake.Word(),
		Icon:        fake.Word(),
		CreatedOn:   uint64(uint32(fake.Date().Unix())),
	}
}

// BuildFakeValidInstrumentList builds a faked ValidInstrumentList.
func BuildFakeValidInstrumentList() *models.ValidInstrumentList {
	exampleValidInstrument1 := BuildFakeValidInstrument()
	exampleValidInstrument2 := BuildFakeValidInstrument()
	exampleValidInstrument3 := BuildFakeValidInstrument()

	return &models.ValidInstrumentList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		ValidInstruments: []models.ValidInstrument{
			*exampleValidInstrument1,
			*exampleValidInstrument2,
			*exampleValidInstrument3,
		},
	}
}

// BuildFakeValidInstrumentUpdateInputFromValidInstrument builds a faked ValidInstrumentUpdateInput from a valid instrument.
func BuildFakeValidInstrumentUpdateInputFromValidInstrument(validInstrument *models.ValidInstrument) *models.ValidInstrumentUpdateInput {
	return &models.ValidInstrumentUpdateInput{
		Name:        validInstrument.Name,
		Variant:     validInstrument.Variant,
		Description: validInstrument.Description,
		Icon:        validInstrument.Icon,
	}
}

// BuildFakeValidInstrumentCreationInput builds a faked ValidInstrumentCreationInput.
func BuildFakeValidInstrumentCreationInput() *models.ValidInstrumentCreationInput {
	validInstrument := BuildFakeValidInstrument()
	return BuildFakeValidInstrumentCreationInputFromValidInstrument(validInstrument)
}

// BuildFakeValidInstrumentCreationInputFromValidInstrument builds a faked ValidInstrumentCreationInput from a valid instrument.
func BuildFakeValidInstrumentCreationInputFromValidInstrument(validInstrument *models.ValidInstrument) *models.ValidInstrumentCreationInput {
	return &models.ValidInstrumentCreationInput{
		Name:        validInstrument.Name,
		Variant:     validInstrument.Variant,
		Description: validInstrument.Description,
		Icon:        validInstrument.Icon,
	}
}
