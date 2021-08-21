package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.InvitationDataManager = (*InvitationDataManager)(nil)

// InvitationDataManager is a mocked types.InvitationDataManager for testing.
type InvitationDataManager struct {
	mock.Mock
}

// InvitationExists is a mock function.
func (m *InvitationDataManager) InvitationExists(ctx context.Context, invitationID uint64) (bool, error) {
	args := m.Called(ctx, invitationID)
	return args.Bool(0), args.Error(1)
}

// GetInvitation is a mock function.
func (m *InvitationDataManager) GetInvitation(ctx context.Context, invitationID uint64) (*types.Invitation, error) {
	args := m.Called(ctx, invitationID)
	return args.Get(0).(*types.Invitation), args.Error(1)
}

// GetAllInvitationsCount is a mock function.
func (m *InvitationDataManager) GetAllInvitationsCount(ctx context.Context) (uint64, error) {
	args := m.Called(ctx)
	return args.Get(0).(uint64), args.Error(1)
}

// GetAllInvitations is a mock function.
func (m *InvitationDataManager) GetAllInvitations(ctx context.Context, results chan []*types.Invitation, bucketSize uint16) error {
	args := m.Called(ctx, results, bucketSize)
	return args.Error(0)
}

// GetInvitations is a mock function.
func (m *InvitationDataManager) GetInvitations(ctx context.Context, filter *types.QueryFilter) (*types.InvitationList, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(*types.InvitationList), args.Error(1)
}

// GetInvitationsWithIDs is a mock function.
func (m *InvitationDataManager) GetInvitationsWithIDs(ctx context.Context, householdID uint64, limit uint8, ids []uint64) ([]*types.Invitation, error) {
	args := m.Called(ctx, householdID, limit, ids)
	return args.Get(0).([]*types.Invitation), args.Error(1)
}

// CreateInvitation is a mock function.
func (m *InvitationDataManager) CreateInvitation(ctx context.Context, input *types.InvitationCreationInput, createdByUser uint64) (*types.Invitation, error) {
	args := m.Called(ctx, input, createdByUser)
	return args.Get(0).(*types.Invitation), args.Error(1)
}

// UpdateInvitation is a mock function.
func (m *InvitationDataManager) UpdateInvitation(ctx context.Context, updated *types.Invitation, changedByUser uint64, changes []*types.FieldChangeSummary) error {
	return m.Called(ctx, updated, changedByUser, changes).Error(0)
}

// ArchiveInvitation is a mock function.
func (m *InvitationDataManager) ArchiveInvitation(ctx context.Context, invitationID, householdID, archivedBy uint64) error {
	return m.Called(ctx, invitationID, householdID, archivedBy).Error(0)
}

// GetAuditLogEntriesForInvitation is a mock function.
func (m *InvitationDataManager) GetAuditLogEntriesForInvitation(ctx context.Context, invitationID uint64) ([]*types.AuditLogEntry, error) {
	args := m.Called(ctx, invitationID)
	return args.Get(0).([]*types.AuditLogEntry), args.Error(1)
}
