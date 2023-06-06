package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.HouseholdInstrumentOwnershipDataManager = (*HouseholdInstrumentOwnershipDataManager)(nil)

// HouseholdInstrumentOwnershipDataManager is a mocked types.HouseholdInstrumentOwnershipDataManager for testing.
type HouseholdInstrumentOwnershipDataManager struct {
	mock.Mock
}

// HouseholdInstrumentOwnershipExists is a mock function.
func (m *HouseholdInstrumentOwnershipDataManager) HouseholdInstrumentOwnershipExists(ctx context.Context, householdInstrumentOwnershipID, householdID string) (bool, error) {
	args := m.Called(ctx, householdInstrumentOwnershipID, householdID)
	return args.Bool(0), args.Error(1)
}

// GetHouseholdInstrumentOwnership is a mock function.
func (m *HouseholdInstrumentOwnershipDataManager) GetHouseholdInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) (*types.HouseholdInstrumentOwnership, error) {
	args := m.Called(ctx, householdInstrumentOwnershipID, householdID)
	return args.Get(0).(*types.HouseholdInstrumentOwnership), args.Error(1)
}

// GetHouseholdInstrumentOwnerships is a mock function.
func (m *HouseholdInstrumentOwnershipDataManager) GetHouseholdInstrumentOwnerships(ctx context.Context, householdID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.HouseholdInstrumentOwnership], error) {
	args := m.Called(ctx, householdID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.HouseholdInstrumentOwnership]), args.Error(1)
}

// CreateHouseholdInstrumentOwnership is a mock function.
func (m *HouseholdInstrumentOwnershipDataManager) CreateHouseholdInstrumentOwnership(ctx context.Context, input *types.HouseholdInstrumentOwnershipDatabaseCreationInput) (*types.HouseholdInstrumentOwnership, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.HouseholdInstrumentOwnership), args.Error(1)
}

// UpdateHouseholdInstrumentOwnership is a mock function.
func (m *HouseholdInstrumentOwnershipDataManager) UpdateHouseholdInstrumentOwnership(ctx context.Context, updated *types.HouseholdInstrumentOwnership) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveHouseholdInstrumentOwnership is a mock function.
func (m *HouseholdInstrumentOwnershipDataManager) ArchiveHouseholdInstrumentOwnership(ctx context.Context, householdInstrumentOwnershipID, householdID string) error {
	return m.Called(ctx, householdInstrumentOwnershipID, householdID).Error(0)
}
