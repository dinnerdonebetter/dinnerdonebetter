package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidIngredientMeasurementUnitDataManager = (*ValidIngredientMeasurementUnitDataManagerMock)(nil)

// ValidIngredientMeasurementUnitDataManagerMock is a mocked types.ValidIngredientMeasurementUnitDataManager for testing.
type ValidIngredientMeasurementUnitDataManagerMock struct {
	mock.Mock
}

// ValidIngredientMeasurementUnitExists is a mock function.
func (m *ValidIngredientMeasurementUnitDataManagerMock) ValidIngredientMeasurementUnitExists(ctx context.Context, validIngredientMeasurementUnitID string) (bool, error) {
	returnValues := m.Called(ctx, validIngredientMeasurementUnitID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidIngredientMeasurementUnit is a mock function.
func (m *ValidIngredientMeasurementUnitDataManagerMock) GetValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*types.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx, validIngredientMeasurementUnitID)
	return returnValues.Get(0).(*types.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

// GetValidIngredientMeasurementUnits is a mock function.
func (m *ValidIngredientMeasurementUnitDataManagerMock) GetValidIngredientMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit]), returnValues.Error(1)
}

// GetValidIngredientMeasurementUnitsForIngredient is a mock function.
func (m *ValidIngredientMeasurementUnitDataManagerMock) GetValidIngredientMeasurementUnitsForIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit], error) {
	returnValues := m.Called(ctx, ingredientID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit]), returnValues.Error(1)
}

// GetValidIngredientMeasurementUnitsForMeasurementUnit is a mock function.
func (m *ValidIngredientMeasurementUnitDataManagerMock) GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx context.Context, measurementUnitID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit], error) {
	returnValues := m.Called(ctx, measurementUnitID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit]), returnValues.Error(1)
}

// CreateValidIngredientMeasurementUnit is a mock function.
func (m *ValidIngredientMeasurementUnitDataManagerMock) CreateValidIngredientMeasurementUnit(ctx context.Context, input *types.ValidIngredientMeasurementUnitDatabaseCreationInput) (*types.ValidIngredientMeasurementUnit, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.ValidIngredientMeasurementUnit), returnValues.Error(1)
}

// UpdateValidIngredientMeasurementUnit is a mock function.
func (m *ValidIngredientMeasurementUnitDataManagerMock) UpdateValidIngredientMeasurementUnit(ctx context.Context, updated *types.ValidIngredientMeasurementUnit) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientMeasurementUnit is a mock function.
func (m *ValidIngredientMeasurementUnitDataManagerMock) ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error {
	return m.Called(ctx, validIngredientMeasurementUnitID).Error(0)
}
