package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeValidVessel builds a faked valid instrument.
func BuildFakeValidVessel() *types.ValidVessel {
	return &types.ValidVessel{
		ID:                             BuildFakeID(),
		Name:                           buildUniqueString(),
		PluralName:                     buildUniqueString(),
		Description:                    buildUniqueString(),
		IconPath:                       buildUniqueString(),
		Slug:                           buildUniqueString(),
		DisplayInSummaryLists:          fake.Bool(),
		IncludeInGeneratedInstructions: fake.Bool(),
		UsableForStorage:               fake.Bool(),
		Capacity:                       float32(buildFakeNumber()),
		CapacityUnit:                   BuildFakeValidMeasurementUnit(),
		WidthInMillimeters:             float32(buildFakeNumber()),
		LengthInMillimeters:            float32(buildFakeNumber()),
		HeightInMillimeters:            float32(buildFakeNumber()),
		Shape:                          "other",
		CreatedAt:                      BuildFakeTime(),
	}
}

// BuildFakeValidVesselsList builds a faked ValidVesselList.
func BuildFakeValidVesselsList() *filtering.QueryFilteredResult[types.ValidVessel] {
	var examples []*types.ValidVessel
	for range exampleQuantity {
		examples = append(examples, BuildFakeValidVessel())
	}

	return &filtering.QueryFilteredResult[types.ValidVessel]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidVesselUpdateRequestInput builds a faked ValidVesselUpdateRequestInput from a valid instrument.
func BuildFakeValidVesselUpdateRequestInput() *types.ValidVesselUpdateRequestInput {
	validVessel := BuildFakeValidVessel()
	return converters.ConvertValidVesselToValidVesselUpdateRequestInput(validVessel)
}

// BuildFakeValidVesselCreationRequestInput builds a faked ValidVesselCreationRequestInput.
func BuildFakeValidVesselCreationRequestInput() *types.ValidVesselCreationRequestInput {
	validVessel := BuildFakeValidVessel()
	return converters.ConvertValidVesselToValidVesselCreationRequestInput(validVessel)
}
