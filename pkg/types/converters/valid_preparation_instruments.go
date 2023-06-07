package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertValidPreparationInstrumentCreationRequestInputToValidPreparationInstrumentDatabaseCreationInput creates a ValidPreparationInstrumentDatabaseCreationInput from a ValidPreparationInstrumentCreationRequestInput.
func ConvertValidPreparationInstrumentCreationRequestInputToValidPreparationInstrumentDatabaseCreationInput(input *types.ValidPreparationInstrumentCreationRequestInput) *types.ValidPreparationInstrumentDatabaseCreationInput {
	x := &types.ValidPreparationInstrumentDatabaseCreationInput{
		ID:                 identifiers.New(),
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

func ConvertValidPreparationInstrumentToValidPreparationInstrumentSearchSubset(x *types.ValidPreparationInstrument) *types.ValidPreparationInstrumentSearchSubset {
	y := &types.ValidPreparationInstrumentSearchSubset{
		ID:          x.ID,
		Notes:       x.Notes,
		Instrument:  types.NamedID{ID: x.Instrument.ID, Name: x.Instrument.Name},
		Preparation: types.NamedID{ID: x.Preparation.ID, Name: x.Preparation.Name},
	}

	return y
}

func ConvertValidIngredientPreparationToValidIngredientPreparationSearchSubset(x *types.ValidIngredientPreparation) *types.ValidIngredientPreparationSearchSubset {
	y := &types.ValidIngredientPreparationSearchSubset{
		ID:          x.ID,
		Notes:       x.Notes,
		Ingredient:  types.NamedID{ID: x.Ingredient.ID, Name: x.Ingredient.Name},
		Preparation: types.NamedID{ID: x.Preparation.ID, Name: x.Preparation.Name},
	}

	return y
}
