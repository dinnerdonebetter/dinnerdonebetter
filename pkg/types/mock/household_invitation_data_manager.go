package mocktypes

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.HouseholdInvitationDataManager = (*HouseholdInvitationDataManager)(nil)

// HouseholdInvitationDataManager is a mocked types.HouseholdInvitationDataManager for testing.
type HouseholdInvitationDataManager struct {
	mock.Mock
}

// HouseholdInvitationExists is a mock function.
func (m *HouseholdInvitationDataManager) HouseholdInvitationExists(ctx context.Context, householdID, userID string) (bool, error) {
	args := m.Called(ctx, householdID, userID)
	return args.Bool(0), args.Error(1)
}

// GetHouseholdInvitation is a mock function.
func (m *HouseholdInvitationDataManager) GetHouseholdInvitation(ctx context.Context, householdID, userID string) (*types.HouseholdInvitation, error) {
	args := m.Called(ctx, householdID, userID)
	return args.Get(0).(*types.HouseholdInvitation), args.Error(1)
}

// GetAllHouseholdInvitationsCount is a mock function.
func (m *HouseholdInvitationDataManager) GetAllHouseholdInvitationsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllHouseholdInvitations is a mock function.
func (m *HouseholdInvitationDataManager) GetAllHouseholdInvitations(ctx context.Context, results chan []*types.HouseholdInvitation, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetHouseholdInvitations is a mock function.
func (m *HouseholdInvitationDataManager) GetHouseholdInvitations(ctx context.Context, userID string, filter *types.QueryFilter) (*types.HouseholdInvitationList, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).(*types.HouseholdInvitationList), args.Error(1)
}

// CreateHouseholdInvitation is a mock function.
func (m *HouseholdInvitationDataManager) CreateHouseholdInvitation(ctx context.Context, input *types.HouseholdInvitationCreationInput) (*types.HouseholdInvitation, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.HouseholdInvitation), args.Error(1)
}

// ArchiveHouseholdInvitation is a mock function.
func (m *HouseholdInvitationDataManager) ArchiveHouseholdInvitation(ctx context.Context, householdID, userID string) error {
	return m.Called(ctx, householdID, userID).Error(0)
}
