package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/pkg/types"
)

var _ types.ValidMeasurementUnitDataManager = (*ValidMeasurementUnitDataManager)(nil)

// ValidMeasurementUnitDataManager is a mocked types.ValidMeasurementUnitDataManager for testing.
type ValidMeasurementUnitDataManager struct {
	mock.Mock
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

// SearchForValidMeasurementUnits is a mock function.
func (m *ValidMeasurementUnitDataManager) SearchForValidMeasurementUnits(ctx context.Context, query string) ([]*types.ValidMeasurementUnit, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*types.ValidMeasurementUnit), args.Error(1)
}

// GetValidMeasurementUnits is a mock function.
func (m *ValidMeasurementUnitDataManager) GetValidMeasurementUnits(ctx context.Context, filter *types.QueryFilter) (*types.ValidMeasurementUnitList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.ValidMeasurementUnitList), args.Error(1)
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
