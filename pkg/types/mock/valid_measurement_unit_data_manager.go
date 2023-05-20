package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidMeasurementUnitDataManager = (*ValidMeasurementUnitDataManager)(nil)

// ValidMeasurementUnitDataManager is a mocked types.ValidMeasurementUnitDataManager for testing.
type ValidMeasurementUnitDataManager struct {
	mock.Mock
}

// ValidMeasurementUnitsForIngredientID is a mock function.
func (m *ValidMeasurementUnitDataManager) ValidMeasurementUnitsForIngredientID(ctx context.Context, validIngredientID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidMeasurementUnit], error) {
	args := m.Called(ctx, validIngredientID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidMeasurementUnit]), args.Error(1)
}

// ValidMeasurementUnitExists is a mock function.
func (m *ValidMeasurementUnitDataManager) ValidMeasurementUnitExists(ctx context.Context, validMeasurementUnitID string) (bool, error) {
	args := m.Called(ctx, validMeasurementUnitID)
	return args.Bool(0), args.Error(1)
}

// GetValidMeasurementUnit is a mock function.
func (m *ValidMeasurementUnitDataManager) GetValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error) {
	args := m.Called(ctx, validMeasurementUnitID)
	return args.Get(0).(*types.ValidMeasurementUnit), args.Error(1)
}

// SearchForValidMeasurementUnitsByName is a mock function.
func (m *ValidMeasurementUnitDataManager) SearchForValidMeasurementUnitsByName(ctx context.Context, query string) ([]*types.ValidMeasurementUnit, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*types.ValidMeasurementUnit), args.Error(1)
}

// GetValidMeasurementUnits is a mock function.
func (m *ValidMeasurementUnitDataManager) GetValidMeasurementUnits(ctx context.Context, filter *types.QueryFilter) (*types.QueryFilteredResult[types.ValidMeasurementUnit], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.ValidMeasurementUnit]), args.Error(1)
}

// CreateValidMeasurementUnit is a mock function.
func (m *ValidMeasurementUnitDataManager) CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitDatabaseCreationInput) (*types.ValidMeasurementUnit, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidMeasurementUnit), args.Error(1)
}

// UpdateValidMeasurementUnit is a mock function.
func (m *ValidMeasurementUnitDataManager) UpdateValidMeasurementUnit(ctx context.Context, updated *types.ValidMeasurementUnit) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidMeasurementUnit is a mock function.
func (m *ValidMeasurementUnitDataManager) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	return m.Called(ctx, validMeasurementUnitID).Error(0)
}
