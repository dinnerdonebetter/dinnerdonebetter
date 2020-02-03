package mock

import (
	"context"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/stretchr/testify/mock"
)

var _ models.InvitationDataManager = (*InvitationDataManager)(nil)

// InvitationDataManager is a mocked models.InvitationDataManager for testing
type InvitationDataManager struct {
	mock.Mock
}

// GetInvitation is a mock function
func (m *InvitationDataManager) GetInvitation(ctx context.Context, invitationID, userID uint64) (*models.Invitation, error) {
	args := m.Called(ctx, invitationID, userID)
	return args.Get(0).(*models.Invitation), args.Error(1)
}

// GetInvitationCount is a mock function
func (m *InvitationDataManager) GetInvitationCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (uint64, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllInvitationsCount is a mock function
func (m *InvitationDataManager) GetAllInvitationsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetInvitations is a mock function
func (m *InvitationDataManager) GetInvitations(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.InvitationList, error) {
	args := m.Called(ctx, filter, userID)
	return args.Get(0).(*models.InvitationList), args.Error(1)
}

// GetAllInvitationsForUser is a mock function
func (m *InvitationDataManager) GetAllInvitationsForUser(ctx context.Context, userID uint64) ([]models.Invitation, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Invitation), args.Error(1)
}

// CreateInvitation is a mock function
func (m *InvitationDataManager) CreateInvitation(ctx context.Context, input *models.InvitationCreationInput) (*models.Invitation, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.Invitation), args.Error(1)
}

// UpdateInvitation is a mock function
func (m *InvitationDataManager) UpdateInvitation(ctx context.Context, updated *models.Invitation) error {
	return m.Called(ctx, updated).Error(0)
}

// ArchiveInvitation is a mock function
func (m *InvitationDataManager) ArchiveInvitation(ctx context.Context, id, userID uint64) error {
	return m.Called(ctx, id, userID).Error(0)
}
