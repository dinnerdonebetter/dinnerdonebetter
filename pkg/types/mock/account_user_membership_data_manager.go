package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.HouseholdUserMembershipDataManager = (*HouseholdUserMembershipDataManager)(nil)

// HouseholdUserMembershipDataManager is a mocked types.HouseholdUserMembershipDataManager for testing.
type HouseholdUserMembershipDataManager struct {
	mock.Mock
}

// BuildSessionContextDataForUser satisfies our interface contract.
func (m *HouseholdUserMembershipDataManager) BuildSessionContextDataForUser(ctx context.Context, userID uint64) (*types.SessionContextData, error) {
	returnValues := m.Called(ctx, userID)

	return returnValues.Get(0).(*types.SessionContextData), returnValues.Error(1)
}

// GetDefaultHouseholdIDForUser satisfies our interface contract.
func (m *HouseholdUserMembershipDataManager) GetDefaultHouseholdIDForUser(ctx context.Context, userID uint64) (uint64, error) {
	returnValues := m.Called(ctx, userID)

	return returnValues.Get(0).(uint64), returnValues.Error(1)
}

// MarkHouseholdAsUserDefault implements the interface.
func (m *HouseholdUserMembershipDataManager) MarkHouseholdAsUserDefault(ctx context.Context, userID, householdID, changedByUser uint64) error {
	return m.Called(ctx, userID, householdID, changedByUser).Error(0)
}

// UserIsMemberOfHousehold implements the interface.
func (m *HouseholdUserMembershipDataManager) UserIsMemberOfHousehold(ctx context.Context, userID, householdID uint64) (bool, error) {
	returnValues := m.Called(ctx, userID, householdID)

	return returnValues.Bool(0), returnValues.Error(1)
}

// AddUserToHousehold implements the interface.
func (m *HouseholdUserMembershipDataManager) AddUserToHousehold(ctx context.Context, input *types.AddUserToHouseholdInput, addedByUser uint64) error {
	return m.Called(ctx, input, addedByUser).Error(0)
}

// RemoveUserFromHousehold implements the interface.
func (m *HouseholdUserMembershipDataManager) RemoveUserFromHousehold(ctx context.Context, userID, householdID, removedByUser uint64, reason string) error {
	return m.Called(ctx, userID, householdID, removedByUser, reason).Error(0)
}

// ModifyUserPermissions implements the interface.
func (m *HouseholdUserMembershipDataManager) ModifyUserPermissions(ctx context.Context, userID, householdID, changedByUser uint64, input *types.ModifyUserPermissionsInput) error {
	return m.Called(ctx, userID, householdID, changedByUser, input).Error(0)
}

// TransferHouseholdOwnership implements the interface.
func (m *HouseholdUserMembershipDataManager) TransferHouseholdOwnership(ctx context.Context, householdID, transferredBy uint64, input *types.HouseholdOwnershipTransferInput) error {
	return m.Called(ctx, householdID, transferredBy, input).Error(0)
}
