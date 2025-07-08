package mockmanagers

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"

	"github.com/stretchr/testify/mock"
)

type MockValidEnumerationsManager struct {
	mock.Mock
}

func (m *MockValidEnumerationsManager) SearchValidIngredientGroups(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientGroup), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidIngredientGroups(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientGroup, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientGroup), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredientGroup(ctx context.Context, input *recipeenums.ValidIngredientGroupCreationRequestInput) (*recipeenums.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredientGroup), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredientGroup(ctx context.Context, validIngredientGroupID string) (*recipeenums.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, validIngredientGroupID)

	return returnValues.Get(0).(*recipeenums.ValidIngredientGroup), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredientGroup(ctx context.Context, validIngredientGroupID string, input *recipeenums.ValidIngredientGroupUpdateRequestInput) (*recipeenums.ValidIngredientGroup, error) {
	returnValues := m.Called(ctx, validIngredientGroupID, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredientGroup), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredientGroup(ctx context.Context, validIngredientGroupID string) error {
	returnValues := m.Called(ctx, validIngredientGroupID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) ListValidIngredientMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientMeasurementUnit, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientMeasurementUnit), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredientMeasurementUnit(ctx context.Context, input *recipeenums.ValidIngredientMeasurementUnitCreationRequestInput) (*recipeenums.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*recipeenums.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx, validIngredientMeasurementUnitID)

	return returnValues.Get(0).(*recipeenums.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string, input *recipeenums.ValidIngredientMeasurementUnitUpdateRequestInput) (*recipeenums.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx, validIngredientMeasurementUnitID, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error {
	returnValues := m.Called(ctx, validIngredientMeasurementUnitID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientMeasurementUnitsByIngredient(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx, validIngredientID, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientMeasurementUnitsByMeasurementUnit(ctx context.Context, validMeasurementUnitID string, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidIngredientPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientPreparation, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientPreparation), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredientPreparation(ctx context.Context, input *recipeenums.ValidIngredientPreparationCreationRequestInput) (*recipeenums.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredientPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*recipeenums.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx, validIngredientPreparationID)

	return returnValues.Get(0).(*recipeenums.ValidIngredientPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string, input *recipeenums.ValidIngredientPreparationUpdateRequestInput) (*recipeenums.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx, validIngredientPreparationID, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredientPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	returnValues := m.Called(ctx, validIngredientPreparationID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientPreparationsByIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx, ingredientID, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientPreparationsByPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientPreparation, error) {
	returnValues := m.Called(ctx, preparationID, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidIngredients(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredient, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredient, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredient), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredient(ctx context.Context, input *recipeenums.ValidIngredientCreationRequestInput) (*recipeenums.ValidIngredient, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredient(ctx context.Context, validIngredientID string) (*recipeenums.ValidIngredient, error) {
	returnValues := m.Called(ctx, validIngredientID)

	return returnValues.Get(0).(*recipeenums.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) RandomValidIngredient(ctx context.Context) (*recipeenums.ValidIngredient, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*recipeenums.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredient(ctx context.Context, validIngredientID string, input *recipeenums.ValidIngredientUpdateRequestInput) (*recipeenums.ValidIngredient, error) {
	returnValues := m.Called(ctx, validIngredientID, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	returnValues := m.Called(ctx, validIngredientID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientsByPreparationAndIngredientName(ctx context.Context, preparationID, query string, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredient, error) {
	returnValues := m.Called(ctx, preparationID, query, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidIngredientStateIngredients(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientStateIngredient, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientStateIngredient), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredientStateIngredient(ctx context.Context, input *recipeenums.ValidIngredientStateIngredientCreationRequestInput) (*recipeenums.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredientStateIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) (*recipeenums.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, validIngredientStateIngredientID)

	return returnValues.Get(0).(*recipeenums.ValidIngredientStateIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string, input *recipeenums.ValidIngredientStateIngredientUpdateRequestInput) (*recipeenums.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, validIngredientStateIngredientID, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredientStateIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredientStateIngredient(ctx context.Context, validIngredientStateIngredientID string) error {
	returnValues := m.Called(ctx, validIngredientStateIngredientID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientStateIngredientsByIngredient(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, validIngredientID, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientStateIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientStateIngredientsByIngredientState(ctx context.Context, validIngredientStateID string, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientStateIngredient, error) {
	returnValues := m.Called(ctx, validIngredientStateID, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientStateIngredient), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidIngredientStates(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientState, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientState), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidIngredientStates(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidIngredientState, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidIngredientState), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidIngredientState(ctx context.Context, input *recipeenums.ValidIngredientStateCreationRequestInput) (*recipeenums.ValidIngredientState, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredientState), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidIngredientState(ctx context.Context, validIngredientStateID string) (*recipeenums.ValidIngredientState, error) {
	returnValues := m.Called(ctx, validIngredientStateID)

	return returnValues.Get(0).(*recipeenums.ValidIngredientState), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidIngredientState(ctx context.Context, validIngredientStateID string, input *recipeenums.ValidIngredientStateUpdateRequestInput) (*recipeenums.ValidIngredientState, error) {
	returnValues := m.Called(ctx, validIngredientStateID, input)

	return returnValues.Get(0).(*recipeenums.ValidIngredientState), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidIngredientState(ctx context.Context, validIngredientStateID string) error {
	returnValues := m.Called(ctx, validIngredientStateID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidMeasurementUnits(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*recipeenums.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*recipeenums.ValidMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidMeasurementUnitsByIngredientID(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) ([]*recipeenums.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, validIngredientID, filter)

	return returnValues.Get(0).([]*recipeenums.ValidMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidMeasurementUnit, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidMeasurementUnit), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidMeasurementUnit(ctx context.Context, input *recipeenums.ValidMeasurementUnitCreationRequestInput) (*recipeenums.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*recipeenums.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).(*recipeenums.ValidMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string, input *recipeenums.ValidMeasurementUnitUpdateRequestInput) (*recipeenums.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID, input)

	return returnValues.Get(0).(*recipeenums.ValidMeasurementUnit), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidInstruments(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*recipeenums.ValidInstrument, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*recipeenums.ValidInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidInstrument, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidInstrument), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidInstrument(ctx context.Context, input *recipeenums.ValidInstrumentCreationRequestInput) (*recipeenums.ValidInstrument, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidInstrument(ctx context.Context, validInstrumentID string) (*recipeenums.ValidInstrument, error) {
	returnValues := m.Called(ctx, validInstrumentID)

	return returnValues.Get(0).(*recipeenums.ValidInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) RandomValidInstrument(ctx context.Context) (*recipeenums.ValidInstrument, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*recipeenums.ValidInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidInstrument(ctx context.Context, validInstrumentID string, input *recipeenums.ValidInstrumentUpdateRequestInput) (*recipeenums.ValidInstrument, error) {
	returnValues := m.Called(ctx, validInstrumentID, input)

	return returnValues.Get(0).(*recipeenums.ValidInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	returnValues := m.Called(ctx, validInstrumentID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) ValidMeasurementUnitConversionsFromMeasurementUnit(ctx context.Context, validMeasurementUnitID string) ([]*recipeenums.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).([]*recipeenums.ValidMeasurementUnitConversion), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ValidMeasurementUnitConversionsToMeasurementUnit(ctx context.Context, validMeasurementUnitID string) ([]*recipeenums.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).([]*recipeenums.ValidMeasurementUnitConversion), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) CreateValidMeasurementUnitConversion(ctx context.Context, input *recipeenums.ValidMeasurementUnitConversionCreationRequestInput) (*recipeenums.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidMeasurementUnitConversion), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) (*recipeenums.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitConversionID)

	return returnValues.Get(0).(*recipeenums.ValidMeasurementUnitConversion), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string, input *recipeenums.ValidMeasurementUnitConversionUpdateRequestInput) (*recipeenums.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitConversionID, input)

	return returnValues.Get(0).(*recipeenums.ValidMeasurementUnitConversion), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) error {
	returnValues := m.Called(ctx, validMeasurementUnitConversionID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) ListValidPreparationInstruments(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidPreparationInstrument, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidPreparationInstrument), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidPreparationInstrument(ctx context.Context, input *recipeenums.ValidPreparationInstrumentCreationRequestInput) (*recipeenums.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidPreparationInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) (*recipeenums.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, validPreparationInstrumentID)

	return returnValues.Get(0).(*recipeenums.ValidPreparationInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string, input *recipeenums.ValidPreparationInstrumentUpdateRequestInput) (*recipeenums.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, validPreparationInstrumentID, input)

	return returnValues.Get(0).(*recipeenums.ValidPreparationInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidPreparationInstrument(ctx context.Context, validPreparationInstrumentID string) error {
	returnValues := m.Called(ctx, validPreparationInstrumentID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidPreparationInstrumentsByPreparation(ctx context.Context, validPreparationID string, filter *filtering.QueryFilter) ([]*recipeenums.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, validPreparationID, filter)

	return returnValues.Get(0).([]*recipeenums.ValidPreparationInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidPreparationInstrumentsByInstrument(ctx context.Context, validInstrumentID string, filter *filtering.QueryFilter) ([]*recipeenums.ValidPreparationInstrument, error) {
	returnValues := m.Called(ctx, validInstrumentID, filter)

	return returnValues.Get(0).([]*recipeenums.ValidPreparationInstrument), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidPreparations(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*recipeenums.ValidPreparation, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*recipeenums.ValidPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidPreparations(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidPreparation, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidPreparation), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidPreparation(ctx context.Context, input *recipeenums.ValidPreparationCreationRequestInput) (*recipeenums.ValidPreparation, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidPreparation(ctx context.Context, validPreparationID string) (*recipeenums.ValidPreparation, error) {
	returnValues := m.Called(ctx, validPreparationID)

	return returnValues.Get(0).(*recipeenums.ValidPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) RandomValidPreparation(ctx context.Context) (*recipeenums.ValidPreparation, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*recipeenums.ValidPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidPreparation(ctx context.Context, validPreparationID string, input *recipeenums.ValidPreparationUpdateRequestInput) (*recipeenums.ValidPreparation, error) {
	returnValues := m.Called(ctx, validPreparationID, input)

	return returnValues.Get(0).(*recipeenums.ValidPreparation), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidPreparation(ctx context.Context, validPreparationID string) error {
	returnValues := m.Called(ctx, validPreparationID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) ListValidPreparationVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidPreparationVessel, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidPreparationVessel), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidPreparationVessel(ctx context.Context, input *recipeenums.ValidPreparationVesselCreationRequestInput) (*recipeenums.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidPreparationVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*recipeenums.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, validPreparationVesselID)

	return returnValues.Get(0).(*recipeenums.ValidPreparationVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidPreparationVessel(ctx context.Context, validPreparationVesselID string, input *recipeenums.ValidPreparationVesselUpdateRequestInput) (*recipeenums.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, validPreparationVesselID, input)

	return returnValues.Get(0).(*recipeenums.ValidPreparationVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error {
	returnValues := m.Called(ctx, validPreparationVesselID)

	return returnValues.Error(0)
}

func (m *MockValidEnumerationsManager) SearchValidPreparationVesselsByPreparation(ctx context.Context, validPreparationID string, filter *filtering.QueryFilter) ([]*recipeenums.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, validPreparationID, filter)

	return returnValues.Get(0).([]*recipeenums.ValidPreparationVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidPreparationVesselsByVessel(ctx context.Context, validVesselID string, filter *filtering.QueryFilter) ([]*recipeenums.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, validVesselID, filter)

	return returnValues.Get(0).([]*recipeenums.ValidPreparationVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) SearchValidVessels(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*recipeenums.ValidVessel, error) {
	returnValues := m.Called(ctx, query, useDatabase, filter)

	return returnValues.Get(0).([]*recipeenums.ValidVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ListValidVessels(ctx context.Context, filter *filtering.QueryFilter) ([]*recipeenums.ValidVessel, string, error) {
	returnValues := m.Called(ctx, filter)

	return returnValues.Get(0).([]*recipeenums.ValidVessel), returnValues.String(1), returnValues.Error(2)
}

func (m *MockValidEnumerationsManager) CreateValidVessel(ctx context.Context, input *recipeenums.ValidVesselCreationRequestInput) (*recipeenums.ValidVessel, error) {
	returnValues := m.Called(ctx, input)

	return returnValues.Get(0).(*recipeenums.ValidVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ReadValidVessel(ctx context.Context, validVesselID string) (*recipeenums.ValidVessel, error) {
	returnValues := m.Called(ctx, validVesselID)

	return returnValues.Get(0).(*recipeenums.ValidVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) RandomValidVessel(ctx context.Context) (*recipeenums.ValidVessel, error) {
	returnValues := m.Called(ctx)

	return returnValues.Get(0).(*recipeenums.ValidVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) UpdateValidVessel(ctx context.Context, validVesselID string, input *recipeenums.ValidVesselUpdateRequestInput) (*recipeenums.ValidVessel, error) {
	returnValues := m.Called(ctx, validVesselID, input)

	return returnValues.Get(0).(*recipeenums.ValidVessel), returnValues.Error(1)
}

func (m *MockValidEnumerationsManager) ArchiveValidVessel(ctx context.Context, validVesselID string) error {
	returnValues := m.Called(ctx, validVesselID)

	return returnValues.Error(0)
}
