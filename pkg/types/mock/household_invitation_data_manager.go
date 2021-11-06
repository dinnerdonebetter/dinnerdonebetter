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
func (m *HouseholdInvitationDataManager) HouseholdInvitationExists(ctx context.Context, householdID, householdInvitationID string) (bool, error) {
	args := m.Called(ctx, householdID, householdInvitationID)
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

func (m *HouseholdInvitationDataManager) GetSentPendingHouseholdInvitations(ctx context.Context, userID string, filter *types.QueryFilter) ([]*types.HouseholdInvitation, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).([]*types.HouseholdInvitation), args.Error(1)
}

func (m *HouseholdInvitationDataManager) GetReceivedPendingHouseholdInvitations(ctx context.Context, userID string, filter *types.QueryFilter) ([]*types.HouseholdInvitation, error) {
	args := m.Called(ctx, userID, filter)
	return args.Get(0).([]*types.HouseholdInvitation), args.Error(1)
}

// CreateHouseholdInvitation is a mock function.
func (m *HouseholdInvitationDataManager) CreateHouseholdInvitation(ctx context.Context, input *types.HouseholdInvitationDatabaseCreationInput) (*types.HouseholdInvitation, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*types.HouseholdInvitation), args.Error(1)
}

func (m *HouseholdInvitationDataManager) CancelHouseholdInvitation(ctx context.Context, invitationID string) error {
	args := m.Called(ctx, invitationID)
	return args.Error(0)
}

func (m *HouseholdInvitationDataManager) AcceptHouseholdInvitation(ctx context.Context, invitationID string) error {
	args := m.Called(ctx, invitationID)
	return args.Error(0)
}

func (m *HouseholdInvitationDataManager) RejectHouseholdInvitation(ctx context.Context, invitationID string) error {
	args := m.Called(ctx, invitationID)
	return args.Error(0)
}

// ArchiveHouseholdInvitation is a mock function.
func (m *HouseholdInvitationDataManager) ArchiveHouseholdInvitation(ctx context.Context, householdID, userID string) error {
	return m.Called(ctx, householdID, userID).Error(0)
}
