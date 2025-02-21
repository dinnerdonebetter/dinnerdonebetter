package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidPreparationVesselDataManager = (*ValidPreparationVesselDataManagerMock)(nil)

// ValidPreparationVesselDataManagerMock is a mocked types.ValidPreparationVesselDataManager for testing.
type ValidPreparationVesselDataManagerMock struct {
	mock.Mock
}

// ValidPreparationVesselExists is a mock function.
func (m *ValidPreparationVesselDataManagerMock) ValidPreparationVesselExists(ctx context.Context, validPreparationVesselID string) (bool, error) {
	returnValues := m.Called(ctx, validPreparationVesselID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidPreparationVessel is a mock function.
func (m *ValidPreparationVesselDataManagerMock) GetValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*types.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, validPreparationVesselID)
	return returnValues.Get(0).(*types.ValidPreparationVessel), returnValues.Error(1)
}

// GetValidPreparationVessels is a mock function.
func (m *ValidPreparationVesselDataManagerMock) GetValidPreparationVessels(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidPreparationVessel], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidPreparationVessel]), returnValues.Error(1)
}

// GetValidPreparationVesselsForPreparation is a mock function.
func (m *ValidPreparationVesselDataManagerMock) GetValidPreparationVesselsForPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidPreparationVessel], error) {
	returnValues := m.Called(ctx, preparationID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidPreparationVessel]), returnValues.Error(1)
}

// GetValidPreparationVesselsForVessel is a mock function.
func (m *ValidPreparationVesselDataManagerMock) GetValidPreparationVesselsForVessel(ctx context.Context, vesselID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidPreparationVessel], error) {
	returnValues := m.Called(ctx, vesselID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidPreparationVessel]), returnValues.Error(1)
}

// CreateValidPreparationVessel is a mock function.
func (m *ValidPreparationVesselDataManagerMock) CreateValidPreparationVessel(ctx context.Context, input *types.ValidPreparationVesselDatabaseCreationInput) (*types.ValidPreparationVessel, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.ValidPreparationVessel), returnValues.Error(1)
}

// UpdateValidPreparationVessel is a mock function.
func (m *ValidPreparationVesselDataManagerMock) UpdateValidPreparationVessel(ctx context.Context, updated *types.ValidPreparationVessel) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidPreparationVessel is a mock function.
func (m *ValidPreparationVesselDataManagerMock) ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error {
	return m.Called(ctx, validPreparationVesselID).Error(0)
}
