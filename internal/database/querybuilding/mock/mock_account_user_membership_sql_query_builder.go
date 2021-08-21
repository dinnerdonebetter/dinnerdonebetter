package mock

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	"github.com/stretchr/testify/mock"
)

var _ querybuilding.HouseholdUserMembershipSQLQueryBuilder = (*HouseholdUserMembershipSQLQueryBuilder)(nil)

// HouseholdUserMembershipSQLQueryBuilder is a mocked types.HouseholdUserMembershipSQLQueryBuilder for testing.
type HouseholdUserMembershipSQLQueryBuilder struct {
	mock.Mock
}

// BuildGetDefaultHouseholdIDForUserQuery implements our interface.
func (m *HouseholdUserMembershipSQLQueryBuilder) BuildGetDefaultHouseholdIDForUserQuery(ctx context.Context, userID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, userID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildTransferHouseholdMembershipsQuery implements our interface.
func (m *HouseholdUserMembershipSQLQueryBuilder) BuildTransferHouseholdMembershipsQuery(ctx context.Context, currentOwnerID, newOwnerID, householdID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, currentOwnerID, newOwnerID, householdID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildArchiveHouseholdMembershipsForUserQuery implements our interface.
func (m *HouseholdUserMembershipSQLQueryBuilder) BuildArchiveHouseholdMembershipsForUserQuery(ctx context.Context, userID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, userID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildGetHouseholdMembershipsForUserQuery implements our interface.
func (m *HouseholdUserMembershipSQLQueryBuilder) BuildGetHouseholdMembershipsForUserQuery(ctx context.Context, userID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, userID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildMarkHouseholdAsUserDefaultQuery implements our interface.
func (m *HouseholdUserMembershipSQLQueryBuilder) BuildMarkHouseholdAsUserDefaultQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, userID, householdID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildCreateMembershipForNewUserQuery implements our interface.
func (m *HouseholdUserMembershipSQLQueryBuilder) BuildCreateMembershipForNewUserQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, userID, householdID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildUserIsMemberOfHouseholdQuery implements our interface.
func (m *HouseholdUserMembershipSQLQueryBuilder) BuildUserIsMemberOfHouseholdQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, userID, householdID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildAddUserToHouseholdQuery implements our interface.
func (m *HouseholdUserMembershipSQLQueryBuilder) BuildAddUserToHouseholdQuery(ctx context.Context, input *types.AddUserToHouseholdInput) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, input)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildRemoveUserFromHouseholdQuery implements our interface.
func (m *HouseholdUserMembershipSQLQueryBuilder) BuildRemoveUserFromHouseholdQuery(ctx context.Context, userID, householdID uint64) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, userID, householdID)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}

// BuildModifyUserPermissionsQuery implements our interface.
func (m *HouseholdUserMembershipSQLQueryBuilder) BuildModifyUserPermissionsQuery(ctx context.Context, userID, householdID uint64, newRoles []string) (query string, args []interface{}) {
	returnArgs := m.Called(ctx, userID, householdID, newRoles)

	return returnArgs.String(0), returnArgs.Get(1).([]interface{})
}
