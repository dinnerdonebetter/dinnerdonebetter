package converters

import (
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

// ConvertInstrumentOwnershipToInstrumentOwnershipUpdateRequestInput creates a DatabaseCreationInput from a CreationInput.
func ConvertInstrumentOwnershipToInstrumentOwnershipUpdateRequestInput(x *types.InstrumentOwnership) *types.InstrumentOwnershipUpdateRequestInput {
	out := &types.InstrumentOwnershipUpdateRequestInput{
		Notes:             &x.Notes,
		Quantity:          &x.Quantity,
		ValidInstrumentID: &x.Instrument.ID,
	}

	return out
}

// ConvertInstrumentOwnershipCreationRequestInputToInstrumentOwnershipDatabaseCreationInput creates a DatabaseCreationInput from a CreationInput.
func ConvertInstrumentOwnershipCreationRequestInputToInstrumentOwnershipDatabaseCreationInput(x *types.InstrumentOwnershipCreationRequestInput) *types.InstrumentOwnershipDatabaseCreationInput {
	out := &types.InstrumentOwnershipDatabaseCreationInput{
		ID:                identifiers.New(),
		Notes:             x.Notes,
		Quantity:          x.Quantity,
		ValidInstrumentID: x.ValidInstrumentID,
	}

	return out
}

// ConvertInstrumentOwnershipToInstrumentOwnershipCreationRequestInput builds a InstrumentOwnershipCreationRequestInput from a Ingredient.
func ConvertInstrumentOwnershipToInstrumentOwnershipCreationRequestInput(x *types.InstrumentOwnership) *types.InstrumentOwnershipCreationRequestInput {
	return &types.InstrumentOwnershipCreationRequestInput{
		Notes:              x.Notes,
		Quantity:           x.Quantity,
		ValidInstrumentID:  x.Instrument.ID,
		BelongsToHousehold: x.BelongsToHousehold,
	}
}

// ConvertInstrumentOwnershipToInstrumentOwnershipDatabaseCreationInput builds a InstrumentOwnershipDatabaseCreationInput from a InstrumentOwnership.
func ConvertInstrumentOwnershipToInstrumentOwnershipDatabaseCreationInput(x *types.InstrumentOwnership) *types.InstrumentOwnershipDatabaseCreationInput {
	return &types.InstrumentOwnershipDatabaseCreationInput{
		ID:                 x.ID,
		Notes:              x.Notes,
		Quantity:           x.Quantity,
		ValidInstrumentID:  x.Instrument.ID,
		BelongsToHousehold: x.BelongsToHousehold,
	}
}
