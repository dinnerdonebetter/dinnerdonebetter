package mocktypes

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.HouseholdInvitationDataManager = (*HouseholdInvitationDataManagerMock)(nil)

// HouseholdInvitationDataManagerMock is a mocked types.HouseholdInvitationDataManager for testing.
type HouseholdInvitationDataManagerMock struct {
	mock.Mock
}

// HouseholdInvitationExists is a mock function.
func (m *HouseholdInvitationDataManagerMock) HouseholdInvitationExists(ctx context.Context, householdInvitationID string) (bool, error) {
	returnValues := m.Called(ctx, householdInvitationID)
	return returnValues.Bool(0), returnValues.Error(1)
}

// GetHouseholdInvitationByHouseholdAndID is a mock function.
func (m *HouseholdInvitationDataManagerMock) GetHouseholdInvitationByHouseholdAndID(ctx context.Context, householdID, householdInvitationID string) (*types.HouseholdInvitation, error) {
	returnValues := m.Called(ctx, householdID, householdInvitationID)
	return returnValues.Get(0).(*types.HouseholdInvitation), returnValues.Error(1)
}

// GetHouseholdInvitationByTokenAndID is a mock function.
func (m *HouseholdInvitationDataManagerMock) GetHouseholdInvitationByTokenAndID(ctx context.Context, householdInvitationID, token string) (*types.HouseholdInvitation, error) {
	returnValues := m.Called(ctx, householdInvitationID, token)
	return returnValues.Get(0).(*types.HouseholdInvitation), returnValues.Error(1)
}

// GetHouseholdInvitationByEmailAndToken is a mock function.
func (m *HouseholdInvitationDataManagerMock) GetHouseholdInvitationByEmailAndToken(ctx context.Context, emailAddress, token string) (*types.HouseholdInvitation, error) {
	returnValues := m.Called(ctx, emailAddress, token)
	return returnValues.Get(0).(*types.HouseholdInvitation), returnValues.Error(1)
}

// GetPendingHouseholdInvitationsFromUser is a mock function.
func (m *HouseholdInvitationDataManagerMock) GetPendingHouseholdInvitationsFromUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.HouseholdInvitation], error) {
	returnValues := m.Called(ctx, userID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.HouseholdInvitation]), returnValues.Error(1)
}

// GetPendingHouseholdInvitationsForUser is a mock function.
func (m *HouseholdInvitationDataManagerMock) GetPendingHouseholdInvitationsForUser(ctx context.Context, userID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.HouseholdInvitation], error) {
	returnValues := m.Called(ctx, userID, filter)
	return returnValues.Get(0).(*filtering.QueryFilteredResult[types.HouseholdInvitation]), returnValues.Error(1)
}

// CreateHouseholdInvitation is a mock function.
func (m *HouseholdInvitationDataManagerMock) CreateHouseholdInvitation(ctx context.Context, input *types.HouseholdInvitationDatabaseCreationInput) (*types.HouseholdInvitation, error) {
	returnValues := m.Called(ctx, input)
	return returnValues.Get(0).(*types.HouseholdInvitation), returnValues.Error(1)
}

// CancelHouseholdInvitation is a mock function.
func (m *HouseholdInvitationDataManagerMock) CancelHouseholdInvitation(ctx context.Context, householdInvitationID, note string) error {
	return m.Called(ctx, householdInvitationID, note).Error(0)
}

// AcceptHouseholdInvitation is a mock function.
func (m *HouseholdInvitationDataManagerMock) AcceptHouseholdInvitation(ctx context.Context, householdInvitationID, token, note string) error {
	return m.Called(ctx, householdInvitationID, token, note).Error(0)
}

// RejectHouseholdInvitation is a mock function.
func (m *HouseholdInvitationDataManagerMock) RejectHouseholdInvitation(ctx context.Context, householdInvitationID, note string) error {
	return m.Called(ctx, householdInvitationID, note).Error(0)
}

// ArchiveHouseholdInvitation is a mock function.
func (m *HouseholdInvitationDataManagerMock) ArchiveHouseholdInvitation(ctx context.Context, householdID, userID string) error {
	return m.Called(ctx, householdID, userID).Error(0)
}
