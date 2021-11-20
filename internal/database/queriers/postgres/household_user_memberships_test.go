package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromHouseholdUserMemberships(memberships ...*types.HouseholdUserMembership) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(householdsUserMembershipTableColumns)

	for _, x := range memberships {
		rowValues := []driver.Value{
			&x.ID,
			&x.BelongsToUser,
			&x.BelongsToHousehold,
			strings.Join(x.HouseholdRoles, householdMemberRolesSeparator),
			&x.DefaultHousehold,
			&x.CreatedOn,
			&x.LastUpdatedOn,
			&x.ArchivedOn,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildInvalidMockRowsFromHouseholdUserMemberships(memberships ...*types.HouseholdUserMembership) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(householdsUserMembershipTableColumns)

	for _, x := range memberships {
		rowValues := []driver.Value{
			&x.DefaultHousehold,
			&x.BelongsToUser,
			&x.BelongsToHousehold,
			strings.Join(x.HouseholdRoles, householdMemberRolesSeparator),
			&x.CreatedOn,
			&x.LastUpdatedOn,
			&x.ArchivedOn,
			&x.ID,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanHouseholdUserMemberships(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := q.scanHouseholdUserMemberships(ctx, mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := q.scanHouseholdUserMemberships(ctx, mockRows)
		assert.Error(t, err)
	})
}

func TestQuerier_BuildSessionContextDataForUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.Members[0].DefaultHousehold = true

		examplePermsMap := map[string]*types.UserHouseholdMembershipInfo{}
		for _, membership := range exampleHousehold.Members {
			examplePermsMap[membership.BelongsToHousehold] = &types.UserHouseholdMembershipInfo{
				HouseholdName:  exampleHousehold.Name,
				HouseholdID:    membership.BelongsToHousehold,
				HouseholdRoles: membership.HouseholdRoles,
			}
		}

		exampleHouseholdPermissionsMap := map[string]authorization.HouseholdRolePermissionsChecker{}
		for _, membership := range exampleHousehold.Members {
			exampleHouseholdPermissionsMap[membership.BelongsToHousehold] = authorization.NewHouseholdRolePermissionChecker(membership.HouseholdRoles...)
		}

		c, db := buildTestClient(t)

		userRetrievalArgs := []interface{}{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getUserQuery)).
			WithArgs(interfaceToDriverValue(userRetrievalArgs)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		getHouseholdMembershipsForUserArgs := []interface{}{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdMembershipsForUserQuery)).
			WithArgs(interfaceToDriverValue(getHouseholdMembershipsForUserArgs)...).
			WillReturnRows(buildMockRowsFromHouseholdUserMemberships(exampleHousehold.Members...))

		expectedActiveHouseholdID := exampleHousehold.Members[0].BelongsToHousehold

		expected := &types.SessionContextData{
			Requester: types.RequesterInfo{
				UserID:                exampleUser.ID,
				Reputation:            exampleUser.ServiceHouseholdStatus,
				ReputationExplanation: exampleUser.ReputationExplanation,
				ServicePermissions:    authorization.NewServiceRolePermissionChecker(exampleUser.ServiceRoles...),
			},
			HouseholdPermissions: exampleHouseholdPermissionsMap,
			ActiveHouseholdID:    expectedActiveHouseholdID,
		}

		actual, err := c.BuildSessionContextDataForUser(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual, "expected and actual RequestContextData do not match")
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.BuildSessionContextDataForUser(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error retrieving user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		examplePermsMap := map[string]*types.UserHouseholdMembershipInfo{}
		for _, membership := range exampleHousehold.Members {
			examplePermsMap[membership.BelongsToHousehold] = &types.UserHouseholdMembershipInfo{
				HouseholdName:  exampleHousehold.Name,
				HouseholdID:    membership.BelongsToHousehold,
				HouseholdRoles: membership.HouseholdRoles,
			}
		}

		c, db := buildTestClient(t)

		userRetrievalArgs := []interface{}{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getUserQuery)).
			WithArgs(interfaceToDriverValue(userRetrievalArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.BuildSessionContextDataForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error retrieving household memberships", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		examplePermsMap := map[string]*types.UserHouseholdMembershipInfo{}
		for _, membership := range exampleHousehold.Members {
			examplePermsMap[membership.BelongsToHousehold] = &types.UserHouseholdMembershipInfo{
				HouseholdName:  exampleHousehold.Name,
				HouseholdID:    membership.BelongsToHousehold,
				HouseholdRoles: membership.HouseholdRoles,
			}
		}

		exampleHouseholdPermissionsMap := map[string]authorization.HouseholdRolePermissionsChecker{}
		for _, membership := range exampleHousehold.Members {
			exampleHouseholdPermissionsMap[membership.BelongsToHousehold] = authorization.NewHouseholdRolePermissionChecker(membership.HouseholdRoles...)
		}

		c, db := buildTestClient(t)

		userRetrievalArgs := []interface{}{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getUserQuery)).
			WithArgs(interfaceToDriverValue(userRetrievalArgs)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		getHouseholdMembershipsForUserArgs := []interface{}{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdMembershipsForUserQuery)).
			WithArgs(interfaceToDriverValue(getHouseholdMembershipsForUserArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.BuildSessionContextDataForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error scanning household user memberships", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		examplePermsMap := map[string]*types.UserHouseholdMembershipInfo{}
		for _, membership := range exampleHousehold.Members {
			examplePermsMap[membership.BelongsToHousehold] = &types.UserHouseholdMembershipInfo{
				HouseholdName:  exampleHousehold.Name,
				HouseholdID:    membership.BelongsToHousehold,
				HouseholdRoles: membership.HouseholdRoles,
			}
		}

		exampleHouseholdPermissionsMap := map[string]authorization.HouseholdRolePermissionsChecker{}
		for _, membership := range exampleHousehold.Members {
			exampleHouseholdPermissionsMap[membership.BelongsToHousehold] = authorization.NewHouseholdRolePermissionChecker(membership.HouseholdRoles...)
		}

		c, db := buildTestClient(t)

		userRetrievalArgs := []interface{}{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getUserQuery)).
			WithArgs(interfaceToDriverValue(userRetrievalArgs)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		getHouseholdMembershipsForUserArgs := []interface{}{exampleUser.ID}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdMembershipsForUserQuery)).
			WithArgs(interfaceToDriverValue(getHouseholdMembershipsForUserArgs)...).
			WillReturnRows(buildInvalidMockRowsFromHouseholdUserMemberships(exampleHousehold.Members...))

		actual, err := c.BuildSessionContextDataForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetDefaultHouseholdIDForUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()
		expected := exampleHouseholdID

		c, db := buildTestClient(t)

		args := []interface{}{exampleUserID, true}

		db.ExpectQuery(formatQueryForSQLMock(getDefaultHouseholdIDForUserQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(exampleHouseholdID))

		actual, err := c.GetDefaultHouseholdIDForUser(ctx, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, db.ExpectationsWereMet())
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetDefaultHouseholdIDForUser(ctx, "")
		assert.Error(t, err)
		assert.Zero(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		args := []interface{}{exampleUserID, true}

		db.ExpectQuery(formatQueryForSQLMock(getDefaultHouseholdIDForUserQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetDefaultHouseholdIDForUser(ctx, exampleUserID)
		assert.Error(t, err)
		assert.Zero(t, actual)

		assert.NoError(t, db.ExpectationsWereMet())
	})
}

func TestQuerier_MarkHouseholdAsUserDefault(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		markHouseholdAsUserDefaultArgs := []interface{}{
			exampleUserID,
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(markHouseholdAsUserDefaultQuery)).
			WithArgs(interfaceToDriverValue(markHouseholdAsUserDefaultArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.MarkHouseholdAsUserDefault(ctx, exampleUserID, exampleHouseholdID))
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkHouseholdAsUserDefault(ctx, "", exampleHousehold.ID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()

		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkHouseholdAsUserDefault(ctx, exampleUser.ID, ""))
	})

	T.Run("with error marking household as default", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		markHouseholdAsUserDefaultArgs := []interface{}{
			exampleUserID,
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(markHouseholdAsUserDefaultQuery)).
			WithArgs(interfaceToDriverValue(markHouseholdAsUserDefaultArgs)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.MarkHouseholdAsUserDefault(ctx, exampleUserID, exampleHouseholdID))
	})
}

func TestQuerier_UserIsMemberOfHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		userIsMemberOfHouseholdArgs := []interface{}{
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(userIsMemberOfHouseholdQuery)).
			WithArgs(interfaceToDriverValue(userIsMemberOfHouseholdArgs)...).
			WillReturnRows(sqlmock.NewRows([]string{"result"}).AddRow(true))

		actual, err := c.UserIsMemberOfHousehold(ctx, exampleUserID, exampleHouseholdID)
		assert.True(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.UserIsMemberOfHousehold(ctx, "", exampleHouseholdID)
		assert.False(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.UserIsMemberOfHousehold(ctx, exampleUserID, "")
		assert.False(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error performing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		userIsMemberOfHouseholdArgs := []interface{}{
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(userIsMemberOfHouseholdQuery)).
			WithArgs(interfaceToDriverValue(userIsMemberOfHouseholdArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.UserIsMemberOfHousehold(ctx, exampleUserID, exampleHouseholdID)
		assert.False(t, actual)
		assert.Error(t, err)
	})
}

func TestQuerier_ModifyUserPermissions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		c, db := buildTestClient(t)

		fakeArgs := []interface{}{
			strings.Join(exampleInput.NewRoles, householdMemberRolesSeparator),
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(modifyUserPermissionsQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ModifyUserPermissions(ctx, exampleHouseholdID, exampleUserID, exampleInput))
	})

	T.Run("with invalid household id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		c, _ := buildTestClient(t)

		assert.Error(t, c.ModifyUserPermissions(ctx, "", exampleUserID, exampleInput))
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		assert.Error(t, c.ModifyUserPermissions(ctx, exampleHouseholdID, exampleUserID, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		c, db := buildTestClient(t)

		fakeArgs := []interface{}{
			strings.Join(exampleInput.NewRoles, householdMemberRolesSeparator),
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(modifyUserPermissionsQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ModifyUserPermissions(ctx, exampleHouseholdID, exampleUserID, exampleInput))
	})
}

func TestQuerier_TransferHouseholdOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeHouseholdTransferArgs := []interface{}{
			exampleInput.NewOwner,
			exampleInput.CurrentOwner,
			exampleHouseholdID,
		}

		db.ExpectExec(formatQueryForSQLMock(transferHouseholdOwnershipQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdTransferArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		fakeHouseholdMembershipsTransferArgs := []interface{}{
			exampleInput.NewOwner,
			exampleHouseholdID,
			exampleInput.CurrentOwner,
		}

		db.ExpectExec(formatQueryForSQLMock(transferHouseholdMembershipQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdMembershipsTransferArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

		assert.NoError(t, c.TransferHouseholdOwnership(ctx, exampleHouseholdID, exampleInput))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, _ := buildTestClient(t)

		assert.Error(t, c.TransferHouseholdOwnership(ctx, "", exampleInput))
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, _ := buildTestClient(t)

		assert.Error(t, c.TransferHouseholdOwnership(ctx, exampleHousehold.ID, nil))
	})

	T.Run("with error starting transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.TransferHouseholdOwnership(ctx, exampleHousehold.ID, exampleInput))
	})

	T.Run("with error writing household transfer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeHouseholdTransferArgs := []interface{}{
			exampleInput.NewOwner,
			exampleInput.CurrentOwner,
			exampleHousehold.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(transferHouseholdOwnershipQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdTransferArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.TransferHouseholdOwnership(ctx, exampleHousehold.ID, exampleInput))
	})

	T.Run("with error writing membership transfers", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeHouseholdTransferArgs := []interface{}{
			exampleInput.NewOwner,
			exampleInput.CurrentOwner,
			exampleHouseholdID,
		}

		db.ExpectExec(formatQueryForSQLMock(transferHouseholdOwnershipQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdTransferArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		fakeHouseholdMembershipsTransferArgs := []interface{}{
			exampleInput.NewOwner,
			exampleHouseholdID,
			exampleInput.CurrentOwner,
		}

		db.ExpectExec(formatQueryForSQLMock(transferHouseholdMembershipQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdMembershipsTransferArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.TransferHouseholdOwnership(ctx, exampleHouseholdID, exampleInput))
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeHouseholdTransferArgs := []interface{}{
			exampleInput.NewOwner,
			exampleInput.CurrentOwner,
			exampleHouseholdID,
		}

		db.ExpectExec(formatQueryForSQLMock(transferHouseholdOwnershipQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdTransferArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		fakeHouseholdMembershipsTransferArgs := []interface{}{
			exampleInput.NewOwner,
			exampleHouseholdID,
			exampleInput.CurrentOwner,
		}

		db.ExpectExec(formatQueryForSQLMock(transferHouseholdMembershipQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdMembershipsTransferArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		assert.Error(t, c.TransferHouseholdOwnership(ctx, exampleHouseholdID, exampleInput))
	})
}

func TestSQLQuerier_addUserToHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleInput := fakes.BuildFakeHouseholdUserMembershipDatabaseCreationInput()

		c, db := buildTestClient(t)

		addUserToHouseholdArgs := []interface{}{
			exampleInput.ID,
			exampleInput.UserID,
			exampleInput.HouseholdID,
			strings.Join(exampleInput.HouseholdRoles, householdMemberRolesSeparator),
		}

		db.ExpectExec(formatQueryForSQLMock(addUserToHouseholdQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.addUserToHousehold(ctx, c.db, exampleInput)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		err := c.addUserToHousehold(ctx, c.db, nil)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_RemoveUserFromHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(removeUserFromHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.RemoveUserFromHousehold(ctx, exampleUserID, exampleHouseholdID))
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, _ := buildTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(ctx, "", exampleHousehold.ID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()

		c, _ := buildTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUser.ID, ""))
	})

	T.Run("with error writing removal to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(removeUserFromHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUserID, exampleHouseholdID))
	})
}
