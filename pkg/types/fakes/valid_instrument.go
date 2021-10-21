package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/segmentio/ksuid"

	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// BuildFakeValidInstrument builds a faked valid instrument.
func BuildFakeValidInstrument() *types.ValidInstrument {
	return &types.ValidInstrument{
		ID:          ksuid.New().String(),
		Name:        fake.Word(),
		Variant:     fake.Word(),
		Description: fake.Word(),
		Icon:        fake.Word(),
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

// BuildFakeValidInstrumentUpdateRequestInput builds a faked ValidInstrumentUpdateRequestInput from a valid instrument.
func BuildFakeValidInstrumentUpdateRequestInput() *types.ValidInstrumentUpdateRequestInput {
	validInstrument := BuildFakeValidInstrument()
	return &types.ValidInstrumentUpdateRequestInput{
		Name:        validInstrument.Name,
		Variant:     validInstrument.Variant,
		Description: validInstrument.Description,
		Icon:        validInstrument.Icon,
	}
}

// BuildFakeValidInstrumentUpdateRequestInputFromValidInstrument builds a faked ValidInstrumentUpdateRequestInput from a valid instrument.
func BuildFakeValidInstrumentUpdateRequestInputFromValidInstrument(validInstrument *types.ValidInstrument) *types.ValidInstrumentUpdateRequestInput {
	return &types.ValidInstrumentUpdateRequestInput{
		Name:        validInstrument.Name,
		Variant:     validInstrument.Variant,
		Description: validInstrument.Description,
		Icon:        validInstrument.Icon,
	}
}

// BuildFakeValidInstrumentCreationRequestInput builds a faked ValidInstrumentCreationRequestInput.
func BuildFakeValidInstrumentCreationRequestInput() *types.ValidInstrumentCreationRequestInput {
	validInstrument := BuildFakeValidInstrument()
	return BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(validInstrument)
}

// BuildFakeValidInstrumentCreationRequestInputFromValidInstrument builds a faked ValidInstrumentCreationRequestInput from a valid instrument.
func BuildFakeValidInstrumentCreationRequestInputFromValidInstrument(validInstrument *types.ValidInstrument) *types.ValidInstrumentCreationRequestInput {
	return &types.ValidInstrumentCreationRequestInput{
		ID:          validInstrument.ID,
		Name:        validInstrument.Name,
		Variant:     validInstrument.Variant,
		Description: validInstrument.Description,
		Icon:        validInstrument.Icon,
	}
}

// BuildFakeValidInstrumentDatabaseCreationInput builds a faked ValidInstrumentDatabaseCreationInput.
func BuildFakeValidInstrumentDatabaseCreationInput() *types.ValidInstrumentDatabaseCreationInput {
	validInstrument := BuildFakeValidInstrument()
	return BuildFakeValidInstrumentDatabaseCreationInputFromValidInstrument(validInstrument)
}

// BuildFakeValidInstrumentDatabaseCreationInputFromValidInstrument builds a faked ValidInstrumentDatabaseCreationInput from a valid instrument.
func BuildFakeValidInstrumentDatabaseCreationInputFromValidInstrument(validInstrument *types.ValidInstrument) *types.ValidInstrumentDatabaseCreationInput {
	return &types.ValidInstrumentDatabaseCreationInput{
		ID:          validInstrument.ID,
		Name:        validInstrument.Name,
		Variant:     validInstrument.Variant,
		Description: validInstrument.Description,
		Icon:        validInstrument.Icon,
	}
}
