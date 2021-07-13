package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeValidPreparationInstrument builds a faked valid preparation instrument.
func BuildFakeValidPreparationInstrument() *types.ValidPreparationInstrument {
	return &types.ValidPreparationInstrument{
		ID:            uint64(fake.Uint32()),
		ExternalID:    fake.UUID(),
		InstrumentID:  uint64(fake.Uint32()),
		PreparationID: uint64(fake.Uint32()),
		Notes:         fake.Word(),
		CreatedOn:     uint64(uint32(fake.Date().Unix())),
	}
}

// BuildFakeValidPreparationInstrumentList builds a faked ValidPreparationInstrumentList.
func BuildFakeValidPreparationInstrumentList() *types.ValidPreparationInstrumentList {
	var examples []*types.ValidPreparationInstrument
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidPreparationInstrument())
	}

	return &types.ValidPreparationInstrumentList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		ValidPreparationInstruments: examples,
	}
}

// BuildFakeValidPreparationInstrumentUpdateInput builds a faked ValidPreparationInstrumentUpdateInput from a valid preparation instrument.
func BuildFakeValidPreparationInstrumentUpdateInput() *types.ValidPreparationInstrumentUpdateInput {
	validPreparationInstrument := BuildFakeValidPreparationInstrument()
	return &types.ValidPreparationInstrumentUpdateInput{
		InstrumentID:  validPreparationInstrument.InstrumentID,
		PreparationID: validPreparationInstrument.PreparationID,
		Notes:         validPreparationInstrument.Notes,
	}
}

// BuildFakeValidPreparationInstrumentUpdateInputFromValidPreparationInstrument builds a faked ValidPreparationInstrumentUpdateInput from a valid preparation instrument.
func BuildFakeValidPreparationInstrumentUpdateInputFromValidPreparationInstrument(validPreparationInstrument *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentUpdateInput {
	return &types.ValidPreparationInstrumentUpdateInput{
		InstrumentID:  validPreparationInstrument.InstrumentID,
		PreparationID: validPreparationInstrument.PreparationID,
		Notes:         validPreparationInstrument.Notes,
	}
}

// BuildFakeValidPreparationInstrumentCreationInput builds a faked ValidPreparationInstrumentCreationInput.
func BuildFakeValidPreparationInstrumentCreationInput() *types.ValidPreparationInstrumentCreationInput {
	validPreparationInstrument := BuildFakeValidPreparationInstrument()
	return BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(validPreparationInstrument)
}

// BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument builds a faked ValidPreparationInstrumentCreationInput from a valid preparation instrument.
func BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(validPreparationInstrument *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentCreationInput {
	return &types.ValidPreparationInstrumentCreationInput{
		InstrumentID:  validPreparationInstrument.InstrumentID,
		PreparationID: validPreparationInstrument.PreparationID,
		Notes:         validPreparationInstrument.Notes,
	}
}
