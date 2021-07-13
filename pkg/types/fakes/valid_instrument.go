package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeValidInstrument builds a faked valid instrument.
func BuildFakeValidInstrument() *types.ValidInstrument {
	return &types.ValidInstrument{
		ID:          uint64(fake.Uint32()),
		ExternalID:  fake.UUID(),
		Name:        fake.Word(),
		Variant:     fake.Word(),
		Description: fake.Word(),
		IconPath:    fake.Word(),
		CreatedOn:   uint64(uint32(fake.Date().Unix())),
	}
}

// BuildFakeValidInstrumentList builds a faked ValidInstrumentList.
func BuildFakeValidInstrumentList() *types.ValidInstrumentList {
	var examples []*types.ValidInstrument
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidInstrument())
	}

	return &types.ValidInstrumentList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		ValidInstruments: examples,
	}
}

// BuildFakeValidInstrumentUpdateInput builds a faked ValidInstrumentUpdateInput from a valid instrument.
func BuildFakeValidInstrumentUpdateInput() *types.ValidInstrumentUpdateInput {
	validInstrument := BuildFakeValidInstrument()
	return &types.ValidInstrumentUpdateInput{
		Name:        validInstrument.Name,
		Variant:     validInstrument.Variant,
		Description: validInstrument.Description,
		IconPath:    validInstrument.IconPath,
	}
}

// BuildFakeValidInstrumentUpdateInputFromValidInstrument builds a faked ValidInstrumentUpdateInput from a valid instrument.
func BuildFakeValidInstrumentUpdateInputFromValidInstrument(validInstrument *types.ValidInstrument) *types.ValidInstrumentUpdateInput {
	return &types.ValidInstrumentUpdateInput{
		Name:        validInstrument.Name,
		Variant:     validInstrument.Variant,
		Description: validInstrument.Description,
		IconPath:    validInstrument.IconPath,
	}
}

// BuildFakeValidInstrumentCreationInput builds a faked ValidInstrumentCreationInput.
func BuildFakeValidInstrumentCreationInput() *types.ValidInstrumentCreationInput {
	validInstrument := BuildFakeValidInstrument()
	return BuildFakeValidInstrumentCreationInputFromValidInstrument(validInstrument)
}

// BuildFakeValidInstrumentCreationInputFromValidInstrument builds a faked ValidInstrumentCreationInput from a valid instrument.
func BuildFakeValidInstrumentCreationInputFromValidInstrument(validInstrument *types.ValidInstrument) *types.ValidInstrumentCreationInput {
	return &types.ValidInstrumentCreationInput{
		Name:        validInstrument.Name,
		Variant:     validInstrument.Variant,
		Description: validInstrument.Description,
		IconPath:    validInstrument.IconPath,
	}
}
