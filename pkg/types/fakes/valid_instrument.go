package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeValidInstrument builds a faked valid instrument.
func BuildFakeValidInstrument() *types.ValidInstrument {
	return &types.ValidInstrument{
		ID:               BuildFakeID(),
		Name:             buildUniqueString(),
		PluralName:       buildUniqueString(),
		Description:      buildUniqueString(),
		IconPath:         buildUniqueString(),
		UsableForStorage: fake.Bool(),
		CreatedAt:        fake.Date(),
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

// BuildFakeValidInstrumentUpdateRequestInput builds a faked ValidInstrumentUpdateRequestInput from a valid instrument.
func BuildFakeValidInstrumentUpdateRequestInput() *types.ValidInstrumentUpdateRequestInput {
	validInstrument := BuildFakeValidInstrument()
	return &types.ValidInstrumentUpdateRequestInput{
		Name:             &validInstrument.Name,
		PluralName:       &validInstrument.PluralName,
		Description:      &validInstrument.Description,
		IconPath:         &validInstrument.IconPath,
		UsableForStorage: &validInstrument.UsableForStorage,
	}
}

// BuildFakeValidInstrumentCreationRequestInput builds a faked ValidInstrumentCreationRequestInput.
func BuildFakeValidInstrumentCreationRequestInput() *types.ValidInstrumentCreationRequestInput {
	validInstrument := BuildFakeValidInstrument()
	return converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(validInstrument)
}
