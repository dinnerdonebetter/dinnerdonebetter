package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/pkg/types"
)

// BuildFakeValidPreparationInstrument builds a faked valid ingredient preparation.
func BuildFakeValidPreparationInstrument() *types.ValidPreparationInstrument {
	return &types.ValidPreparationInstrument{
		ID:          ksuid.New().String(),
		Notes:       buildUniqueString(),
		Preparation: *BuildFakeValidPreparation(),
		Instrument:  *BuildFakeValidInstrument(),
		CreatedAt:   fake.Date(),
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

// BuildFakeValidPreparationInstrumentUpdateRequestInput builds a faked ValidPreparationInstrumentUpdateRequestInput from a valid ingredient preparation.
func BuildFakeValidPreparationInstrumentUpdateRequestInput() *types.ValidPreparationInstrumentUpdateRequestInput {
	validPreparationInstrument := BuildFakeValidPreparationInstrument()
	return &types.ValidPreparationInstrumentUpdateRequestInput{
		Notes:              &validPreparationInstrument.Notes,
		ValidPreparationID: &validPreparationInstrument.Preparation.ID,
		ValidInstrumentID:  &validPreparationInstrument.Instrument.ID,
	}
}

// BuildFakeValidPreparationInstrumentUpdateRequestInputFromValidPreparationInstrument builds a faked ValidPreparationInstrumentUpdateRequestInput from a valid ingredient preparation.
func BuildFakeValidPreparationInstrumentUpdateRequestInputFromValidPreparationInstrument(validPreparationInstrument *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentUpdateRequestInput {
	return &types.ValidPreparationInstrumentUpdateRequestInput{
		Notes:              &validPreparationInstrument.Notes,
		ValidPreparationID: &validPreparationInstrument.Preparation.ID,
		ValidInstrumentID:  &validPreparationInstrument.Instrument.ID,
	}
}

// BuildFakeValidPreparationInstrumentCreationRequestInput builds a faked ValidPreparationInstrumentCreationRequestInput.
func BuildFakeValidPreparationInstrumentCreationRequestInput() *types.ValidPreparationInstrumentCreationRequestInput {
	validPreparationInstrument := BuildFakeValidPreparationInstrument()
	return BuildFakeValidPreparationInstrumentCreationRequestInputFromValidPreparationInstrument(validPreparationInstrument)
}

// BuildFakeValidPreparationInstrumentCreationRequestInputFromValidPreparationInstrument builds a faked ValidPreparationInstrumentCreationRequestInput from a valid ingredient preparation.
func BuildFakeValidPreparationInstrumentCreationRequestInputFromValidPreparationInstrument(validPreparationInstrument *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentCreationRequestInput {
	return &types.ValidPreparationInstrumentCreationRequestInput{
		ID:                 validPreparationInstrument.ID,
		Notes:              validPreparationInstrument.Notes,
		ValidPreparationID: validPreparationInstrument.Preparation.ID,
		ValidInstrumentID:  validPreparationInstrument.Instrument.ID,
	}
}

// BuildFakeValidPreparationInstrumentDatabaseCreationInput builds a faked ValidPreparationInstrumentDatabaseCreationInput.
func BuildFakeValidPreparationInstrumentDatabaseCreationInput() *types.ValidPreparationInstrumentDatabaseCreationInput {
	validPreparationInstrument := BuildFakeValidPreparationInstrument()
	return BuildFakeValidPreparationInstrumentDatabaseCreationInputFromValidPreparationInstrument(validPreparationInstrument)
}

// BuildFakeValidPreparationInstrumentDatabaseCreationInputFromValidPreparationInstrument builds a faked ValidPreparationInstrumentDatabaseCreationInput from a valid ingredient preparation.
func BuildFakeValidPreparationInstrumentDatabaseCreationInputFromValidPreparationInstrument(validPreparationInstrument *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentDatabaseCreationInput {
	return &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 validPreparationInstrument.ID,
		Notes:              validPreparationInstrument.Notes,
		ValidPreparationID: validPreparationInstrument.Preparation.ID,
		ValidInstrumentID:  validPreparationInstrument.Instrument.ID,
	}
}
