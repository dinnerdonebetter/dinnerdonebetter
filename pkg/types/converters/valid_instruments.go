package converters

import "github.com/prixfixeco/backend/pkg/types"

// ConvertValidInstrumentToValidInstrumentUpdateRequestInput creates a ValidInstrumentUpdateRequestInput from a ValidInstrument.
func ConvertValidInstrumentToValidInstrumentUpdateRequestInput(input *types.ValidInstrument) *types.ValidInstrumentUpdateRequestInput {
	x := &types.ValidInstrumentUpdateRequestInput{
		Name:             &input.Name,
		PluralName:       &input.PluralName,
		Description:      &input.Description,
		IconPath:         &input.IconPath,
		UsableForStorage: &input.UsableForStorage,
	}

	return x
}

// ConvertValidInstrumentCreationRequestInputToValidInstrumentDatabaseCreationInput creates a ValidInstrumentDatabaseCreationInput from a ValidInstrumentCreationRequestInput.
func ConvertValidInstrumentCreationRequestInputToValidInstrumentDatabaseCreationInput(input *types.ValidInstrumentCreationRequestInput) *types.ValidInstrumentDatabaseCreationInput {
	x := &types.ValidInstrumentDatabaseCreationInput{
		Name:             input.Name,
		PluralName:       input.PluralName,
		Description:      input.Description,
		IconPath:         input.IconPath,
		UsableForStorage: input.UsableForStorage,
	}

	return x
}

// ConvertNullableValidInstrumentToValidInstrument produces a ValidInstrument from a NullableValidInstrument.
func ConvertNullableValidInstrumentToValidInstrument(x types.NullableValidInstrument) *types.ValidInstrument {
	return &types.ValidInstrument{
		LastUpdatedAt:    x.LastUpdatedAt,
		ArchivedAt:       x.ArchivedAt,
		Description:      *x.Description,
		IconPath:         *x.IconPath,
		ID:               *x.ID,
		Name:             *x.Name,
		PluralName:       *x.PluralName,
		CreatedAt:        *x.CreatedAt,
		UsableForStorage: *x.UsableForStorage,
	}
}

// ConvertValidInstrumentToValidInstrumentCreationRequestInput builds a ValidInstrumentCreationRequestInput from a ValidInstrument.
func ConvertValidInstrumentToValidInstrumentCreationRequestInput(validInstrument *types.ValidInstrument) *types.ValidInstrumentCreationRequestInput {
	return &types.ValidInstrumentCreationRequestInput{
		ID:               validInstrument.ID,
		Name:             validInstrument.Name,
		PluralName:       validInstrument.PluralName,
		Description:      validInstrument.Description,
		IconPath:         validInstrument.IconPath,
		UsableForStorage: validInstrument.UsableForStorage,
	}
}

// ConvertValidInstrumentToValidInstrumentDatabaseCreationInput builds a ValidInstrumentDatabaseCreationInput from a ValidInstrument.
func ConvertValidInstrumentToValidInstrumentDatabaseCreationInput(validInstrument *types.ValidInstrument) *types.ValidInstrumentDatabaseCreationInput {
	return &types.ValidInstrumentDatabaseCreationInput{
		ID:               validInstrument.ID,
		Name:             validInstrument.Name,
		PluralName:       validInstrument.PluralName,
		Description:      validInstrument.Description,
		IconPath:         validInstrument.IconPath,
		UsableForStorage: validInstrument.UsableForStorage,
	}
}
