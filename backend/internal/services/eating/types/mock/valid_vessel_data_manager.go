package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

var _ types.ValidVesselDataManager = (*ValidVesselDataManagerMock)(nil)

// ValidVesselDataManagerMock is a mocked types.ValidVesselDataManager for testing.
type ValidVesselDataManagerMock struct {
	mock.Mock
}

// ValidVesselExists is a mock function.
func (m *ValidVesselDataManagerMock) ValidVesselExists(ctx context.Context, validVesselID string) (bool, error) {
	returnValues := m.Called(ctx, validVesselID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetValidVessel is a mock function.
func (m *ValidVesselDataManagerMock) GetValidVessel(ctx context.Context, validVesselID string) (*types.ValidVessel, error) {
	returnValues := m.Called(ctx, validVesselID)
	return returnValues.Get(0).(*types.ValidVessel), returnValues.Error(1)
}

// GetRandomValidVessel is a mock function.
func (m *ValidVesselDataManagerMock) GetRandomValidVessel(ctx context.Context) (*types.ValidVessel, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).(*types.ValidVessel), returnValues.Error(1)
}

// SearchForValidVessels is a mock function.
func (m *ValidVesselDataManagerMock) SearchForValidVessels(ctx context.Context, query string) ([]*types.ValidVessel, error) {
	returnValues := m.Called(ctx, query)
	return returnValues.Get(0).([]*types.ValidVessel), returnValues.Error(1)
}

// GetValidVessels is a mock function.
func (m *ValidVesselDataManagerMock) GetValidVessels(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidVessel], error) {
	returnValues := m.Called(ctx, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.ValidVessel]), returnValues.Error(1)
}

// CreateValidVessel is a mock function.
func (m *ValidVesselDataManagerMock) CreateValidVessel(ctx context.Context, input *types.ValidVesselDatabaseCreationInput) (*types.ValidVessel, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.ValidVessel), returnValues.Error(1)
}

// UpdateValidVessel is a mock function.
func (m *ValidVesselDataManagerMock) UpdateValidVessel(ctx context.Context, updated *types.ValidVessel) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveValidVessel is a mock function.
func (m *ValidVesselDataManagerMock) ArchiveValidVessel(ctx context.Context, validVesselID string) error {
	return m.Called(ctx, validVesselID).Error(0)
}

// MarkValidVesselAsIndexed is a mock function.
func (m *ValidVesselDataManagerMock) MarkValidVesselAsIndexed(ctx context.Context, validVesselID string) error {
	return m.Called(ctx, validVesselID).Error(0)
}

// GetValidVesselIDsThatNeedSearchIndexing is a mock function.
func (m *ValidVesselDataManagerMock) GetValidVesselIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	returnValues := m.Called(ctx)
	return returnValues.Get(0).([]string), returnValues.Error(1)
}

// GetValidVesselsWithIDs is a mock function.
func (m *ValidVesselDataManagerMock) GetValidVesselsWithIDs(ctx context.Context, ids []string) ([]*types.ValidVessel, error) {
	returnValues := m.Called(ctx, ids)
	return returnValues.Get(0).([]*types.ValidVessel), returnValues.Error(1)
}
