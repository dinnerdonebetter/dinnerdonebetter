package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidMeasurementUnitConversionDataManager = (*ValidMeasurementUnitConversionDataManagerMock)(nil)

// ValidMeasurementUnitConversionDataManagerMock is a mocked types.ValidMeasurementUnitConversionDataManager for testing.
type ValidMeasurementUnitConversionDataManagerMock struct {
	mock.Mock
}

// GetValidMeasurementUnitConversionsFromUnit is a mock function.
func (m *ValidMeasurementUnitConversionDataManagerMock) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).([]*types.ValidMeasurementUnitConversion), returnValues.Error(1)
}

// GetValidMeasurementUnitConversionsToUnit is a mock function.
func (m *ValidMeasurementUnitConversionDataManagerMock) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validMeasurementUnitID)

	return returnValues.Get(0).([]*types.ValidMeasurementUnitConversion), returnValues.Error(1)
}

// ValidMeasurementUnitConversionExists is a mock function.
func (m *ValidMeasurementUnitConversionDataManagerMock) ValidMeasurementUnitConversionExists(ctx context.Context, validPreparationID string) (bool, error) {
	returnValues := m.Called(ctx, validPreparationID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidMeasurementUnitConversion is a mock function.
func (m *ValidMeasurementUnitConversionDataManagerMock) GetValidMeasurementUnitConversion(ctx context.Context, validPreparationID string) (*types.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, validPreparationID)
	return returnValues.Get(0).(*types.ValidMeasurementUnitConversion), returnValues.Error(1)
}

// CreateValidMeasurementUnitConversion is a mock function.
func (m *ValidMeasurementUnitConversionDataManagerMock) CreateValidMeasurementUnitConversion(ctx context.Context, input *types.ValidMeasurementUnitConversionDatabaseCreationInput) (*types.ValidMeasurementUnitConversion, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.ValidMeasurementUnitConversion), returnValues.Error(1)
}

// UpdateValidMeasurementUnitConversion is a mock function.
func (m *ValidMeasurementUnitConversionDataManagerMock) UpdateValidMeasurementUnitConversion(ctx context.Context, updated *types.ValidMeasurementUnitConversion) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidMeasurementUnitConversion is a mock function.
func (m *ValidMeasurementUnitConversionDataManagerMock) ArchiveValidMeasurementUnitConversion(ctx context.Context, validPreparationID string) error {
	return m.Called(ctx, validPreparationID).Error(0)
}
