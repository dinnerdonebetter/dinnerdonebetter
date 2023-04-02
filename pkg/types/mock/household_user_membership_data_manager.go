package mocktypes

import (
	"context"

	"github.com/prixfixeco/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ types.HouseholdUserMembershipDataManager = (*HouseholdUserMembershipDataManager)(nil)

// HouseholdUserMembershipDataManager is a mocked types.HouseholdUserMembershipDataManager for testing.
type HouseholdUserMembershipDataManager struct {
	mock.Mock
}

// BuildSessionContextDataForUser satisfies our interface contract.
func (m *HouseholdUserMembershipDataManager) BuildSessionContextDataForUser(ctx context.Context, userID string) (*types.SessionContextData, error) {
	returnValues := m.Called(ctx, userID)

	return returnValues.Get(0).(*types.SessionContextData), returnValues.Error(1)
}

// GetDefaultHouseholdIDForUser satisfies our interface contract.
func (m *HouseholdUserMembershipDataManager) GetDefaultHouseholdIDForUser(ctx context.Context, userID string) (string, error) {
	returnValues := m.Called(ctx, userID)

	return returnValues.Get(0).(string), returnValues.Error(1)
}

// MarkHouseholdAsUserDefault implements the interface.
func (m *HouseholdUserMembershipDataManager) MarkHouseholdAsUserDefault(ctx context.Context, userID, householdID string) error {
	return m.Called(ctx, userID, householdID).Error(0)
}

// UserIsMemberOfHousehold implements the interface.
func (m *HouseholdUserMembershipDataManager) UserIsMemberOfHousehold(ctx context.Context, userID, householdID string) (bool, error) {
	returnValues := m.Called(ctx, userID, householdID)

	return returnValues.Bool(0), returnValues.Error(1)
}

// AddUserToHousehold implements the interface.
func (m *HouseholdUserMembershipDataManager) AddUserToHousehold(ctx context.Context, input *types.HouseholdUserMembershipDatabaseCreationInput) error {
	return m.Called(ctx, input).Error(0)
}

// RemoveUserFromHousehold implements the interface.
func (m *HouseholdUserMembershipDataManager) RemoveUserFromHousehold(ctx context.Context, userID, householdID string) error {
	return m.Called(ctx, userID, householdID).Error(0)
}

// ModifyUserPermissions implements the interface.
func (m *HouseholdUserMembershipDataManager) ModifyUserPermissions(ctx context.Context, householdID, userID string, input *types.ModifyUserPermissionsInput) error {
	return m.Called(ctx, userID, householdID, input).Error(0)
}

// TransferHouseholdOwnership implements the interface.
func (m *HouseholdUserMembershipDataManager) TransferHouseholdOwnership(ctx context.Context, householdID string, input *types.HouseholdOwnershipTransferInput) error {
	return m.Called(ctx, householdID, input).Error(0)
}
