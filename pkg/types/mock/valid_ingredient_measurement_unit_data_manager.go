package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.ValidIngredientMeasurementUnitDataManager = (*ValidIngredientMeasurementUnitDataManager)(nil)

// ValidIngredientMeasurementUnitDataManager is a mocked types.ValidIngredientMeasurementUnitDataManager for testing.
type ValidIngredientMeasurementUnitDataManager struct {
	mock.Mock
}

// ValidIngredientMeasurementUnitExists is a mock function.
func (m *ValidIngredientMeasurementUnitDataManager) ValidIngredientMeasurementUnitExists(ctx context.Context, validIngredientMeasurementUnitID string) (bool, error) {
	args := m.Called(ctx, validIngredientMeasurementUnitID)
	return args.Bool(0), args.Error(1)
}

// GetValidIngredientMeasurementUnit is a mock function.
func (m *ValidIngredientMeasurementUnitDataManager) GetValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*types.ValidIngredientMeasurementUnit, error) {
	args := m.Called(ctx, validIngredientMeasurementUnitID)
	return args.Get(0).(*types.ValidIngredientMeasurementUnit), args.Error(1)
}

// GetTotalValidIngredientMeasurementUnitCount is a mock function.
func (m *ValidIngredientMeasurementUnitDataManager) GetTotalValidIngredientMeasurementUnitCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetValidIngredientMeasurementUnits is a mock function.
func (m *ValidIngredientMeasurementUnitDataManager) GetValidIngredientMeasurementUnits(ctx context.Context, filter *types.QueryFilter) (*types.ValidIngredientMeasurementUnitList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ValidIngredientMeasurementUnitList), args.Error(1)
}

// CreateValidIngredientMeasurementUnit is a mock function.
func (m *ValidIngredientMeasurementUnitDataManager) CreateValidIngredientMeasurementUnit(ctx context.Context, input *types.ValidIngredientMeasurementUnitDatabaseCreationInput) (*types.ValidIngredientMeasurementUnit, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidIngredientMeasurementUnit), args.Error(1)
}

// UpdateValidIngredientMeasurementUnit is a mock function.
func (m *ValidIngredientMeasurementUnitDataManager) UpdateValidIngredientMeasurementUnit(ctx context.Context, updated *types.ValidIngredientMeasurementUnit) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidIngredientMeasurementUnit is a mock function.
func (m *ValidIngredientMeasurementUnitDataManager) ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error {
	return m.Called(ctx, validIngredientMeasurementUnitID).Error(0)
}
