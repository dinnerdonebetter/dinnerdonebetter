package mocktypes

import (
	"context"

	"github.com/prixfixeco/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.HouseholdDataManager = (*HouseholdDataManager)(nil)

// HouseholdDataManager is a mocked types.HouseholdDataManager for testing.
type HouseholdDataManager struct {
	mock.Mock
}

// HouseholdExists is a mock function.
func (m *HouseholdDataManager) HouseholdExists(ctx context.Context, householdID, userID string) (bool, error) {
	args := m.Called(ctx, householdID, userID)
	return args.Bool(0), args.Error(1)
}

// GetHousehold is a mock function.
func (m *HouseholdDataManager) GetHousehold(ctx context.Context, householdID, userID string) (*types.Household, error) {
	args := m.Called(ctx, householdID, userID)
	return args.Get(0).(*types.Household), args.Error(1)
}

// GetHouseholdByID is a mock function.
func (m *HouseholdDataManager) GetHouseholdByID(ctx context.Context, householdID string) (*types.Household, error) {
	args := m.Called(ctx, householdID)
	return args.Get(0).(*types.Household), args.Error(1)
}

// GetAllHouseholds is a mock function.
func (m *HouseholdDataManager) GetAllHouseholds(ctx context.Context, results chan []*types.Household, bucketSize uint16) error {
	return m.Called(ctx, results, bucketSize).Error(0)
}

// GetHouseholds is a mock function.
func (m *HouseholdDataManager) GetHouseholds(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Household], error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.Household]), args.Error(1)
}

// GetHouseholdsForAdmin is a mock function.
func (m *HouseholdDataManager) GetHouseholdsForAdmin(ctx context.Context, userID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.Household], error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*types.QueryFilteredResult[types.Household]), args.Error(1)
}

// CreateHousehold is a mock function.
func (m *HouseholdDataManager) CreateHousehold(ctx context.Context, input *types.HouseholdDatabaseCreationInput) (*types.Household, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.Household), args.Error(1)
}

// UpdateHousehold is a mock function.
func (m *HouseholdDataManager) UpdateHousehold(ctx context.Context, updated *types.Household) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveHousehold is a mock function.
func (m *HouseholdDataManager) ArchiveHousehold(ctx context.Context, householdID, userID string) error {
	return m.Called(ctx, householdID, userID).Error(0)
}
