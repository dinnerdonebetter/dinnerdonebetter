package converters

import (
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// ConvertValidPreparationVesselCreationRequestInputToValidPreparationVesselDatabaseCreationInput creates a ValidPreparationVesselDatabaseCreationInput from a ValidPreparationVesselCreationRequestInput.
func ConvertValidPreparationVesselCreationRequestInputToValidPreparationVesselDatabaseCreationInput(x *types.ValidPreparationVesselCreationRequestInput) *types.ValidPreparationVesselDatabaseCreationInput {
	return &types.ValidPreparationVesselDatabaseCreationInput{
		ID:                 identifiers.New(),
		Notes:              x.Notes,
		ValidPreparationID: x.ValidPreparationID,
		ValidVesselID:      x.ValidVesselID,
	}
}

// ConvertValidPreparationVesselToValidPreparationVesselUpdateRequestInput builds a ValidPreparationVesselUpdateRequestInput from a ValidPreparationVessel.
func ConvertValidPreparationVesselToValidPreparationVesselUpdateRequestInput(x *types.ValidPreparationVessel) *types.ValidPreparationVesselUpdateRequestInput {
	return &types.ValidPreparationVesselUpdateRequestInput{
		Notes:              &x.Notes,
		ValidPreparationID: &x.Preparation.ID,
		ValidVesselID:      &x.Vessel.ID,
	}
}

// ConvertValidPreparationVesselToValidPreparationVesselCreationRequestInput builds a ValidPreparationVesselCreationRequestInput from a ValidPreparationVessel.
func ConvertValidPreparationVesselToValidPreparationVesselCreationRequestInput(x *types.ValidPreparationVessel) *types.ValidPreparationVesselCreationRequestInput {
	return &types.ValidPreparationVesselCreationRequestInput{
		Notes:              x.Notes,
		ValidPreparationID: x.Preparation.ID,
		ValidVesselID:      x.Vessel.ID,
	}
}

// ConvertValidPreparationVesselToValidPreparationVesselDatabaseCreationInput builds a ValidPreparationVesselDatabaseCreationInput from a ValidPreparationVessel.
func ConvertValidPreparationVesselToValidPreparationVesselDatabaseCreationInput(x *types.ValidPreparationVessel) *types.ValidPreparationVesselDatabaseCreationInput {
	return &types.ValidPreparationVesselDatabaseCreationInput{
		ID:                 x.ID,
		Notes:              x.Notes,
		ValidPreparationID: x.Preparation.ID,
		ValidVesselID:      x.Vessel.ID,
	}
}
