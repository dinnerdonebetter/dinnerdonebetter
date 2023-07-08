package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertValidVesselToValidVesselUpdateRequestInput creates a ValidVesselUpdateRequestInput from a ValidVessel.
func ConvertValidVesselToValidVesselUpdateRequestInput(x *types.ValidVessel) *types.ValidVesselUpdateRequestInput {
	return &types.ValidVesselUpdateRequestInput{
		Name:                           &x.Name,
		PluralName:                     &x.PluralName,
		Description:                    &x.Description,
		IconPath:                       &x.IconPath,
		UsableForStorage:               &x.UsableForStorage,
		Slug:                           &x.Slug,
		DisplayInSummaryLists:          &x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: &x.IncludeInGeneratedInstructions,
		Capacity:                       &x.Capacity,
		CapacityUnitID:                 &x.CapacityUnit.ID,
		WidthInMillimeters:             &x.WidthInMillimeters,
		LengthInMillimeters:            &x.LengthInMillimeters,
		HeightInMillimeters:            &x.HeightInMillimeters,
		Shape:                          &x.Shape,
	}
}

// ConvertValidVesselCreationRequestInputToValidVesselDatabaseCreationInput creates a ValidVesselDatabaseCreationInput from a ValidVesselCreationRequestInput.
func ConvertValidVesselCreationRequestInputToValidVesselDatabaseCreationInput(x *types.ValidVesselCreationRequestInput) *types.ValidVesselDatabaseCreationInput {
	return &types.ValidVesselDatabaseCreationInput{
		ID:                             identifiers.New(),
		Name:                           x.Name,
		PluralName:                     x.PluralName,
		Description:                    x.Description,
		IconPath:                       x.IconPath,
		UsableForStorage:               x.UsableForStorage,
		Slug:                           x.Slug,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
		Capacity:                       x.Capacity,
		CapacityUnitID:                 x.CapacityUnitID,
		WidthInMillimeters:             x.WidthInMillimeters,
		LengthInMillimeters:            x.LengthInMillimeters,
		HeightInMillimeters:            x.HeightInMillimeters,
		Shape:                          x.Shape,
	}
}

// ConvertNullableValidVesselToValidVessel produces a ValidVessel from a NullableValidVessel.
func ConvertNullableValidVesselToValidVessel(x *types.NullableValidVessel) *types.ValidVessel {
	return &types.ValidVessel{
		ID:                             *x.ID,
		Name:                           *x.Name,
		PluralName:                     *x.PluralName,
		Description:                    *x.Description,
		IconPath:                       *x.IconPath,
		UsableForStorage:               *x.UsableForStorage,
		Slug:                           *x.Slug,
		DisplayInSummaryLists:          *x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: *x.IncludeInGeneratedInstructions,
		Capacity:                       *x.Capacity,
		CapacityUnit:                   *x.CapacityUnit,
		WidthInMillimeters:             *x.WidthInMillimeters,
		LengthInMillimeters:            *x.LengthInMillimeters,
		HeightInMillimeters:            *x.HeightInMillimeters,
		Shape:                          *x.Shape,
	}
}

// ConvertValidVesselToValidVesselCreationRequestInput builds a ValidVesselCreationRequestInput from a ValidVessel.
func ConvertValidVesselToValidVesselCreationRequestInput(x *types.ValidVessel) *types.ValidVesselCreationRequestInput {
	return &types.ValidVesselCreationRequestInput{
		Name:                           x.Name,
		PluralName:                     x.PluralName,
		Description:                    x.Description,
		IconPath:                       x.IconPath,
		UsableForStorage:               x.UsableForStorage,
		Slug:                           x.Slug,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
		Capacity:                       x.Capacity,
		CapacityUnitID:                 x.CapacityUnit.ID,
		WidthInMillimeters:             x.WidthInMillimeters,
		LengthInMillimeters:            x.LengthInMillimeters,
		HeightInMillimeters:            x.HeightInMillimeters,
		Shape:                          x.Shape,
	}
}

// ConvertValidVesselToValidVesselDatabaseCreationInput builds a ValidVesselDatabaseCreationInput from a ValidVessel.
func ConvertValidVesselToValidVesselDatabaseCreationInput(x *types.ValidVessel) *types.ValidVesselDatabaseCreationInput {
	return &types.ValidVesselDatabaseCreationInput{
		ID:                             x.ID,
		Name:                           x.Name,
		PluralName:                     x.PluralName,
		Description:                    x.Description,
		IconPath:                       x.IconPath,
		UsableForStorage:               x.UsableForStorage,
		Slug:                           x.Slug,
		DisplayInSummaryLists:          x.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: x.IncludeInGeneratedInstructions,
		Capacity:                       x.Capacity,
		CapacityUnitID:                 x.CapacityUnit.ID,
		WidthInMillimeters:             x.WidthInMillimeters,
		LengthInMillimeters:            x.LengthInMillimeters,
		HeightInMillimeters:            x.HeightInMillimeters,
		Shape:                          x.Shape,
	}
}

// ConvertValidVesselToValidVesselSearchSubset converts a ValidVessel to a ValidVesselSearchSubset.
func ConvertValidVesselToValidVesselSearchSubset(x *types.ValidVessel) *types.ValidVesselSearchSubset {
	return &types.ValidVesselSearchSubset{
		ID:          x.ID,
		Name:        x.Name,
		PluralName:  x.PluralName,
		Description: x.Description,
	}
}
