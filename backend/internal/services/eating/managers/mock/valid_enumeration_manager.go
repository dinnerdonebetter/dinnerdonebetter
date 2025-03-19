package mockmanagers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

type MockValidEnumerationsManager struct {
	mock.Mock
}

func (m *MockValidEnumerationsManager) SearchValidIngredientGroups(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*types.ValidIngredientGroup), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidIngredientGroups(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientGroup, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidIngredientGroup), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredientGroup(ctx context.Context, input *types.ValidIngredientGroupCreationRequestInput) (*types.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidIngredientGroup), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredientGroup(ctx context.Context, validIngredientGroupID string) (*types.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, validIngredientGroupID)

	return returnValues.Get(0).(*types.ValidIngredientGroup), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredientGroup(ctx context.Context, validIngredientGroupID string, input *types.ValidIngredientGroupUpdateRequestInput) (*types.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, validIngredientGroupID, input)

	return returnValues.Get(0).(*types.ValidIngredientGroup), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredientGroup(ctx context.Context, validIngredientGroupID string) error {
	returnValues := m.Called(ctx, validIngredientGroupID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) ListValidIngredientMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientMeasurementUnit, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidIngredientMeasurementUnit), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredientMeasurementUnit(ctx context.Context, input *types.ValidIngredientMeasurementUnitCreationRequestInput) (*types.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*types.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*types.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string, input *types.ValidIngredientMeasurementUnitUpdateRequestInput) (*types.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*types.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error {
	returnValues := m.Called(ctx)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientMeasurementUnitsByIngredient(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).([]*types.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, validMeasurementUnitID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).([]*types.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidIngredientPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientPreparation, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidIngredientPreparation), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationCreationRequestInput) (*types.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidIngredientPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*types.ValidIngredientPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string, input *types.ValidIngredientPreparationUpdateRequestInput) (*types.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*types.ValidIngredientPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	returnValues := m.Called(ctx)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientPreparationsByIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).([]*types.ValidIngredientPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientPreparationsByPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).([]*types.ValidIngredientPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidIngredients(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredient, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*types.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredient, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidIngredient), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredient(ctx context.Context, input *types.ValidIngredientCreationRequestInput) (*types.ValidIngredient, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error) {
	returnValues := m.Called(ctx, validIngredientID)

	return returnValues.Get(0).(*types.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) RandomValidIngredient(ctx context.Context) (*types.ValidIngredient, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*types.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredient(ctx context.Context, validIngredientID string, input *types.ValidIngredientUpdateRequestInput) (*types.ValidIngredient, error) {
	returnValues := m.Called(ctx, validIngredientID, input)

	return returnValues.Get(0).(*types.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	returnValues := m.Called(ctx, validIngredientID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientsByPreparationAndIngredientName(ctx context.Context, preparationID, query string, filter *filtering.QueryFilter) ([]*types.ValidIngredient, error) {
	returnValues := m.Called(ctx, preparationID, query, filter)

	return returnValues.Get(0).([]*types.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidIngredientStateIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientStateIngredient, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidIngredientStateIngredient), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredientStateIngredient(ctx context.Context, input *types.ValidIngredientStateIngredientCreationRequestInput) (*types.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidIngredientStateIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*types.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, validIngredientStateIngredientID)

	return returnValues.Get(0).(*types.ValidIngredientStateIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string, input *types.ValidIngredientStateIngredientUpdateRequestInput) (*types.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, validIngredientStateIngredientID, input)

	return returnValues.Get(0).(*types.ValidIngredientStateIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error {
	returnValues := m.Called(ctx, validIngredientStateIngredientID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientStateIngredientsByIngredient(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, validIngredientID, filter)

	return returnValues.Get(0).([]*types.ValidIngredientStateIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientStateIngredientsByIngredientState(ctx context.Context, validIngredientStateID string, filter *filtering.QueryFilter) ([]*types.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, validIngredientStateID, filter)

	return returnValues.Get(0).([]*types.ValidIngredientStateIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientStates(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidIngredientState, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*types.ValidIngredientState), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidIngredientStates(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidIngredientState, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidIngredientState), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredientState(ctx context.Context, input *types.ValidIngredientStateCreationRequestInput) (*types.ValidIngredientState, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidIngredientState), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredientState(ctx context.Context, validIngredientStateID string) (*types.ValidIngredientState, error) {
	returnValues := m.Called(ctx, validIngredientStateID)

	return returnValues.Get(0).(*types.ValidIngredientState), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredientState(ctx context.Context, validIngredientStateID string, input *types.ValidIngredientStateUpdateRequestInput) (*types.ValidIngredientState, error) {
	returnValues := m.Called(ctx, validIngredientStateID, input)

	return returnValues.Get(0).(*types.ValidIngredientState), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error {
	returnValues := m.Called(ctx, validIngredientStateID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidMeasurementUnits(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*types.ValidMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidMeasurementUnitsByIngredientID(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, validIngredientID, filter)

	return returnValues.Get(0).([]*types.ValidMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidMeasurementUnit, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidMeasurementUnit), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitCreationRequestInput) (*types.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).(*types.ValidMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string, input *types.ValidMeasurementUnitUpdateRequestInput) (*types.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID, input)

	return returnValues.Get(0).(*types.ValidMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidInstruments(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidInstrument, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*types.ValidInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidInstrument, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidInstrument), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentCreationRequestInput) (*types.ValidInstrument, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	returnValues := m.Called(ctx, validInstrumentID)

	return returnValues.Get(0).(*types.ValidInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) RandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*types.ValidInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidInstrument(ctx context.Context, validInstrumentID string, input *types.ValidInstrumentUpdateRequestInput) (*types.ValidInstrument, error) {
	returnValues := m.Called(ctx, validInstrumentID, input)

	return returnValues.Get(0).(*types.ValidInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	returnValues := m.Called(ctx, validInstrumentID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) ValidMeasurementUnitConversionsFromMeasurementUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).([]*types.ValidMeasurementUnitConversion), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ValidMeasurementUnitConversionsToMeasurementUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).([]*types.ValidMeasurementUnitConversion), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) CreateValidMeasurementUnitConversion(ctx context.Context, input *types.ValidMeasurementUnitConversionCreationRequestInput) (*types.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidMeasurementUnitConversion), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) (*types.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitConversionID)

	return returnValues.Get(0).(*types.ValidMeasurementUnitConversion), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string, input *types.ValidMeasurementUnitConversionUpdateRequestInput) (*types.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitConversionID, input)

	return returnValues.Get(0).(*types.ValidMeasurementUnitConversion), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) error {
	returnValues := m.Called(ctx, validMeasurementUnitConversionID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) ListValidPreparationInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparationInstrument, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidPreparationInstrument), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidPreparationInstrument(ctx context.Context, input *types.ValidPreparationInstrumentCreationRequestInput) (*types.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidPreparationInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*types.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, validPreparationInstrumentID)

	return returnValues.Get(0).(*types.ValidPreparationInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string, input *types.ValidPreparationInstrumentUpdateRequestInput) (*types.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, validPreparationInstrumentID, input)

	return returnValues.Get(0).(*types.ValidPreparationInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	returnValues := m.Called(ctx, validPreparationInstrumentID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidPreparationInstrumentsByPreparation(ctx context.Context, validPreparationID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, validPreparationID, filter)

	return returnValues.Get(0).([]*types.ValidPreparationInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidPreparationInstrumentsByInstrument(ctx context.Context, validInstrumentID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, validInstrumentID, filter)

	return returnValues.Get(0).([]*types.ValidPreparationInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidPreparations(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidPreparation, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*types.ValidPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparation, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidPreparation), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidPreparation(ctx context.Context, input *types.ValidPreparationCreationRequestInput) (*types.ValidPreparation, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidPreparation(ctx context.Context, validPreparationID string) (*types.ValidPreparation, error) {
	returnValues := m.Called(ctx, validPreparationID)

	return returnValues.Get(0).(*types.ValidPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) RandomValidPreparation(ctx context.Context) (*types.ValidPreparation, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*types.ValidPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidPreparation(ctx context.Context, validPreparationID string, input *types.ValidPreparationUpdateRequestInput) (*types.ValidPreparation, error) {
	returnValues := m.Called(ctx, validPreparationID, input)

	return returnValues.Get(0).(*types.ValidPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	returnValues := m.Called(ctx, validPreparationID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) ListValidPreparationVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidPreparationVessel, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidPreparationVessel), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidPreparationVessel(ctx context.Context, input *types.ValidPreparationVesselCreationRequestInput) (*types.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidPreparationVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*types.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, validPreparationVesselID)

	return returnValues.Get(0).(*types.ValidPreparationVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidPreparationVessel(ctx context.Context, validPreparationVesselID string, input *types.ValidPreparationVesselUpdateRequestInput) (*types.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, validPreparationVesselID, input)

	return returnValues.Get(0).(*types.ValidPreparationVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error {
	returnValues := m.Called(ctx, validPreparationVesselID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidPreparationVesselsByPreparation(ctx context.Context, validPreparationID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, validPreparationID, filter)

	return returnValues.Get(0).([]*types.ValidPreparationVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidPreparationVesselsByVessel(ctx context.Context, validVesselID string, filter *filtering.QueryFilter) ([]*types.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, validVesselID, filter)

	return returnValues.Get(0).([]*types.ValidPreparationVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidVessels(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*types.ValidVessel, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*types.ValidVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*types.ValidVessel, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*types.ValidVessel), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidVessel(ctx context.Context, input *types.ValidVesselCreationRequestInput) (*types.ValidVessel, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*types.ValidVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidVessel(ctx context.Context, validVesselID string) (*types.ValidVessel, error) {
	returnValues := m.Called(ctx, validVesselID)

	return returnValues.Get(0).(*types.ValidVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) RandomValidVessel(ctx context.Context) (*types.ValidVessel, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*types.ValidVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidVessel(ctx context.Context, validVesselID string, input *types.ValidVesselUpdateRequestInput) (*types.ValidVessel, error) {
	returnValues := m.Called(ctx, validVesselID, input)

	return returnValues.Get(0).(*types.ValidVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidVessel(ctx context.Context, validVesselID string) error {
	returnValues := m.Called(ctx, validVesselID)

	return returnValues.Error(0)
}
