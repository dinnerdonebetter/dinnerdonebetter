package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.InvitationDataManager = (*InvitationDataManager)(nil)

// InvitationDataManager is a mocked models.InvitationDataManager for testing.
type InvitationDataManager struct {
	mock.Mock
}

// InvitationExists is a mock function.
func (m *InvitationDataManager) InvitationExists(ctx context.Context, invitationID uint64) (bool, error) {
	args := m.Called(ctx, invitationID)
	return args.Bool(0), args.Error(1)
}

// GetInvitation is a mock function.
func (m *InvitationDataManager) GetInvitation(ctx context.Context, invitationID uint64) (*models.Invitation, error) {
	args := m.Called(ctx, invitationID)
	return args.Get(0).(*models.Invitation), args.Error(1)
}

// GetAllInvitationsCount is a mock function.
func (m *InvitationDataManager) GetAllInvitationsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllInvitations is a mock function.
func (m *InvitationDataManager) GetAllInvitations(ctx context.Context, results chan []models.Invitation) error {
	args := m.Called(ctx, results)
	return args.Error(0)
}

// GetInvitations is a mock function.
func (m *InvitationDataManager) GetInvitations(ctx context.Context, filter *models.QueryFilter) (*models.InvitationList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*models.InvitationList), args.Error(1)
}

// GetInvitationsWithIDs is a mock function.
func (m *InvitationDataManager) GetInvitationsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.Invitation, error) {
	args := m.Called(ctx, limit, ids)
	return args.Get(0).([]models.Invitation), args.Error(1)
}

// CreateInvitation is a mock function.
func (m *InvitationDataManager) CreateInvitation(ctx context.Context, input *models.InvitationCreationInput) (*models.Invitation, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.Invitation), args.Error(1)
}

// UpdateInvitation is a mock function.
func (m *InvitationDataManager) UpdateInvitation(ctx context.Context, updated *models.Invitation) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveInvitation is a mock function.
func (m *InvitationDataManager) ArchiveInvitation(ctx context.Context, invitationID, userID uint64) error {
	return m.Called(ctx, invitationID, userID).Error(0)
}
