package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeValidInstrument builds a faked valid instrument.
func BuildFakeValidInstrument() *types.ValidInstrument {
	return &types.ValidInstrument{
		ID:                             BuildFakeID(),
		Name:                           buildUniqueString(),
		PluralName:                     buildUniqueString(),
		Description:                    buildUniqueString(),
		IconPath:                       buildUniqueString(),
		Slug:                           buildUniqueString(),
		DisplayInSummaryLists:          fake.Bool(),
		IncludeInGeneratedInstructions: fake.Bool(),
		UsableForStorage:               fake.Bool(),
		CreatedAt:                      BuildFakeTime(),
	}
}

// BuildFakeValidInstrumentsList builds a faked ValidInstrumentList.
func BuildFakeValidInstrumentsList() *filtering.QueryFilteredResult[types.ValidInstrument] {
	var examples []*types.ValidInstrument
	for range exampleQuantity {
		examples = append(examples, BuildFakeValidInstrument())
	}

	return &filtering.QueryFilteredResult[types.ValidInstrument]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidInstrumentUpdateRequestInput builds a faked ValidInstrumentUpdateRequestInput from a valid instrument.
func BuildFakeValidInstrumentUpdateRequestInput() *types.ValidInstrumentUpdateRequestInput {
	validInstrument := BuildFakeValidInstrument()
	return converters.ConvertValidInstrumentToValidInstrumentUpdateRequestInput(validInstrument)
}

// BuildFakeValidInstrumentCreationRequestInput builds a faked ValidInstrumentCreationRequestInput.
func BuildFakeValidInstrumentCreationRequestInput() *types.ValidInstrumentCreationRequestInput {
	validInstrument := BuildFakeValidInstrument()
	return converters.ConvertValidInstrumentToValidInstrumentCreationRequestInput(validInstrument)
}
