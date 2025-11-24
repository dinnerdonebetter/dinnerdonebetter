package fakes

import (
	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

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
		Shape:                          types.VesselShapeOther,
		CreatedAt:                      BuildFakeTime(),
	}
}

// BuildFakeValidVesselsList builds a faked ValidVesselList.
func BuildFakeValidVesselsList() *filtering.QueryFilteredResult[types.ValidVessel] {
	var examples []*types.ValidVessel
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidVessel())
	}

	return &filtering.QueryFilteredResult[types.ValidVessel]{
		Pagination: filtering.Pagination{
			Cursor:        BuildFakeID(),
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
