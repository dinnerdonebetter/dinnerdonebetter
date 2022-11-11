package fakes

import (
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeValidPreparationInstrument builds a faked valid preparation instrument.
func BuildFakeValidPreparationInstrument() *types.ValidPreparationInstrument {
	return &types.ValidPreparationInstrument{
		ID:          BuildFakeID(),
		Notes:       buildUniqueString(),
		Preparation: *BuildFakeValidPreparation(),
		Instrument:  *BuildFakeValidInstrument(),
		CreatedAt:   BuildFakeTime(),
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

// BuildFakeValidPreparationInstrumentUpdateRequestInput builds a faked ValidPreparationInstrumentUpdateRequestInput from a valid preparation instrument.
func BuildFakeValidPreparationInstrumentUpdateRequestInput() *types.ValidPreparationInstrumentUpdateRequestInput {
	validPreparationInstrument := BuildFakeValidPreparationInstrument()
	return &types.ValidPreparationInstrumentUpdateRequestInput{
		Notes:              &validPreparationInstrument.Notes,
		ValidPreparationID: &validPreparationInstrument.Preparation.ID,
		ValidInstrumentID:  &validPreparationInstrument.Instrument.ID,
	}
}

// BuildFakeValidPreparationInstrumentCreationRequestInput builds a faked ValidPreparationInstrumentCreationRequestInput.
func BuildFakeValidPreparationInstrumentCreationRequestInput() *types.ValidPreparationInstrumentCreationRequestInput {
	validPreparationInstrument := BuildFakeValidPreparationInstrument()
	return converters.ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput(validPreparationInstrument)
}
