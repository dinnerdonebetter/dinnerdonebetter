package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.ValidMeasurementConversionDataManager = (*ValidMeasurementConversionDataManager)(nil)

// ValidMeasurementConversionDataManager is a mocked types.ValidMeasurementConversionDataManager for testing.
type ValidMeasurementConversionDataManager struct {
	mock.Mock
}

// GetValidMeasurementConversionsFromUnit is a mock function.
func (m *ValidMeasurementConversionDataManager) GetValidMeasurementConversionsFromUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).([]*types.ValidMeasurementConversion), returnValues.Error(1)
}

// GetValidMeasurementConversionsToUnit is a mock function.
func (m *ValidMeasurementConversionDataManager) GetValidMeasurementConversionsToUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).([]*types.ValidMeasurementConversion), returnValues.Error(1)
}

// ValidMeasurementConversionExists is a mock function.
func (m *ValidMeasurementConversionDataManager) ValidMeasurementConversionExists(ctx context.Context, validPreparationID string) (bool, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Bool(0), args.Error(1)
}

// GetValidMeasurementConversion is a mock function.
func (m *ValidMeasurementConversionDataManager) GetValidMeasurementConversion(ctx context.Context, validPreparationID string) (*types.ValidMeasurementConversion, error) {
	args := m.Called(ctx, validPreparationID)
	return args.Get(0).(*types.ValidMeasurementConversion), args.Error(1)
}

// CreateValidMeasurementConversion is a mock function.
func (m *ValidMeasurementConversionDataManager) CreateValidMeasurementConversion(ctx context.Context, input *types.ValidMeasurementConversionDatabaseCreationInput) (*types.ValidMeasurementConversion, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.ValidMeasurementConversion), args.Error(1)
}

// UpdateValidMeasurementConversion is a mock function.
func (m *ValidMeasurementConversionDataManager) UpdateValidMeasurementConversion(ctx context.Context, updated *types.ValidMeasurementConversion) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidMeasurementConversion is a mock function.
func (m *ValidMeasurementConversionDataManager) ArchiveValidMeasurementConversion(ctx context.Context, validPreparationID string) error {
	return m.Called(ctx, validPreparationID).Error(0)
}
