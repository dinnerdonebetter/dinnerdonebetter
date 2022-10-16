package converters

import (
	"github.com/prixfixeco/api_server/pkg/types"
)

// ConvertValidPreparationInstrumentToValidPreparationInstrument creates a ValidPreparationInstrumentUpdateRequestInput from a ValidPreparationInstrument.
func ConvertValidPreparationInstrumentToValidPreparationInstrument(input *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentUpdateRequestInput {
	x := &types.ValidPreparationInstrumentUpdateRequestInput{
		Notes:              &input.Notes,
		ValidPreparationID: &input.Preparation.ID,
		ValidInstrumentID:  &input.Instrument.ID,
	}

	return x
}

// ConvertValidPreparationInstrumentCreationRequestInputToValidPreparationInstrumentDatabaseCreationInput creates a ValidPreparationInstrumentDatabaseCreationInput from a ValidPreparationInstrumentCreationRequestInput.
func ConvertValidPreparationInstrumentCreationRequestInputToValidPreparationInstrumentDatabaseCreationInput(input *types.ValidPreparationInstrumentCreationRequestInput) *types.ValidPreparationInstrumentDatabaseCreationInput {
	x := &types.ValidPreparationInstrumentDatabaseCreationInput{
		Notes:              input.Notes,
		ValidPreparationID: input.ValidPreparationID,
		ValidInstrumentID:  input.ValidInstrumentID,
	}

	return x
}

// ConvertValidPreparationInstrumentToValidPreparationInstrumentUpdateRequestInput builds a ValidPreparationInstrumentUpdateRequestInput from a ValidPreparationInstrument.
func ConvertValidPreparationInstrumentToValidPreparationInstrumentUpdateRequestInput(validPreparationInstrument *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentUpdateRequestInput {
	return &types.ValidPreparationInstrumentUpdateRequestInput{
		Notes:              &validPreparationInstrument.Notes,
		ValidPreparationID: &validPreparationInstrument.Preparation.ID,
		ValidInstrumentID:  &validPreparationInstrument.Instrument.ID,
	}
}

// ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput builds a ValidPreparationInstrumentCreationRequestInput from a ValidPreparationInstrument.
func ConvertValidPreparationInstrumentToValidPreparationInstrumentCreationRequestInput(validPreparationInstrument *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentCreationRequestInput {
	return &types.ValidPreparationInstrumentCreationRequestInput{
		ID:                 validPreparationInstrument.ID,
		Notes:              validPreparationInstrument.Notes,
		ValidPreparationID: validPreparationInstrument.Preparation.ID,
		ValidInstrumentID:  validPreparationInstrument.Instrument.ID,
	}
}

// ConvertValidPreparationInstrumentToValidPreparationInstrumentDatabaseCreationInput builds a ValidPreparationInstrumentDatabaseCreationInput from a ValidPreparationInstrument.
func ConvertValidPreparationInstrumentToValidPreparationInstrumentDatabaseCreationInput(validPreparationInstrument *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentDatabaseCreationInput {
	return &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 validPreparationInstrument.ID,
		Notes:              validPreparationInstrument.Notes,
		ValidPreparationID: validPreparationInstrument.Preparation.ID,
		ValidInstrumentID:  validPreparationInstrument.Instrument.ID,
	}
}
