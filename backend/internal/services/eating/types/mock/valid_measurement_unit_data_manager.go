package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidMeasurementUnitDataManager = (*ValidMeasurementUnitDataManagerMock)(nil)

// ValidMeasurementUnitDataManagerMock is a mocked types.ValidMeasurementUnitDataManager for testing.
type ValidMeasurementUnitDataManagerMock struct {
	mock.Mock
}

// ValidMeasurementUnitsForIngredientID is a mock function.
func (m *ValidMeasurementUnitDataManagerMock) ValidMeasurementUnitsForIngredientID(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidMeasurementUnit], error) {
	returnValues := m.Called(ctx, validIngredientID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidMeasurementUnit]), returnValues.Error(1)
}

// ValidMeasurementUnitExists is a mock function.
func (m *ValidMeasurementUnitDataManagerMock) ValidMeasurementUnitExists(ctx context.Context, validMeasurementUnitID string) (bool, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidMeasurementUnit is a mock function.
func (m *ValidMeasurementUnitDataManagerMock) GetValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)
	return returnValues.Get(0).(*types.ValidMeasurementUnit), returnValues.Error(1)
}

// SearchForValidMeasurementUnitsByName is a mock function.
func (m *ValidMeasurementUnitDataManagerMock) SearchForValidMeasurementUnits(ctx context.Context, query string) ([]*types.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, query)
	return returnValues.Get(0).([]*types.ValidMeasurementUnit), returnValues.Error(1)
}

// GetValidMeasurementUnits is a mock function.
func (m *ValidMeasurementUnitDataManagerMock) GetValidMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidMeasurementUnit], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidMeasurementUnit]), returnValues.Error(1)
}

// CreateValidMeasurementUnit is a mock function.
func (m *ValidMeasurementUnitDataManagerMock) CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitDatabaseCreationInput) (*types.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.ValidMeasurementUnit), returnValues.Error(1)
}

// UpdateValidMeasurementUnit is a mock function.
func (m *ValidMeasurementUnitDataManagerMock) UpdateValidMeasurementUnit(ctx context.Context, updated *types.ValidMeasurementUnit) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidMeasurementUnit is a mock function.
func (m *ValidMeasurementUnitDataManagerMock) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	return m.Called(ctx, validMeasurementUnitID).Error(0)
}

// MarkValidMeasurementUnitAsIndexed is a mock function.
func (m *ValidMeasurementUnitDataManagerMock) MarkValidMeasurementUnitAsIndexed(ctx context.Context, validMeasurementUnitID string) error {
	return m.Called(ctx, validMeasurementUnitID).Error(0)
}

// GetValidMeasurementUnitIDsThatNeedSearchIndexing is a mock function.
func (m *ValidMeasurementUnitDataManagerMock) GetValidMeasurementUnitIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetValidMeasurementUnitsWithIDs is a mock function.
func (m *ValidMeasurementUnitDataManagerMock) GetValidMeasurementUnitsWithIDs(ctx context.Context, ids []string) ([]*types.ValidMeasurementUnit, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*types.ValidMeasurementUnit), returnValues.Error(1)
}
