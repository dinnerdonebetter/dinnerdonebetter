package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.HouseholdDataManager = (*HouseholdDataManager)(nil)

// HouseholdDataManager is a mocked types.HouseholdDataManager for testing.
type HouseholdDataManager struct {
	mock.Mock
}

// HouseholdExists is a mock function.
func (m *HouseholdDataManager) HouseholdExists(ctx context.Context, householdID, userID uint64) (bool, error) {
	args := m.Called(ctx, householdID, userID)
	return args.Bool(0), args.Error(1)
}

// GetHousehold is a mock function.
func (m *HouseholdDataManager) GetHousehold(ctx context.Context, householdID, userID uint64) (*types.Household, error) {
	args := m.Called(ctx, householdID, userID)
	return args.Get(0).(*types.Household), args.Error(1)
}

// GetAllHouseholdsCount is a mock function.
func (m *HouseholdDataManager) GetAllHouseholdsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllHouseholds is a mock function.
func (m *HouseholdDataManager) GetAllHouseholds(ctx context.Context, results chan []*types.Household, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetHouseholds is a mock function.
func (m *HouseholdDataManager) GetHouseholds(ctx context.Context, userID uint64, filter *types.QueryFilter) (*types.HouseholdList, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*types.HouseholdList), args.Error(1)
}

// GetHouseholdsForAdmin is a mock function.
func (m *HouseholdDataManager) GetHouseholdsForAdmin(ctx context.Context, filter *types.QueryFilter) (*types.HouseholdList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.HouseholdList), args.Error(1)
}

// CreateHousehold is a mock function.
func (m *HouseholdDataManager) CreateHousehold(ctx context.Context, input *types.HouseholdCreationInput, createdByUser uint64) (*types.Household, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.Household), args.Error(1)
}

// UpdateHousehold is a mock function.
func (m *HouseholdDataManager) UpdateHousehold(ctx context.Context, updated *types.Household, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveHousehold is a mock function.
func (m *HouseholdDataManager) ArchiveHousehold(ctx context.Context, householdID, userID, archivedByUser uint64) error {
	return m.Called(ctx, householdID, userID, archivedByUser).Error(0)
}

// GetAuditLogEntriesForHousehold is a mock function.
func (m *HouseholdDataManager) GetAuditLogEntriesForHousehold(ctx context.Context, householdID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, householdID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
