package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromHouseholdUserMembershipsWithUsers(memberships ...*types.HouseholdUserMembershipWithUser) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(householdsUserMembershipTableColumns)

	for _, x := range memberships {
		rowValues := []driver.Value{
			&x.ID,
			&x.BelongsToUser.ID,
			&x.BelongsToHousehold,
			x.HouseholdRole,
			&x.DefaultHousehold,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildInvalidRowsFromHouseholdUserMembershipsWithUsers(memberships ...*types.HouseholdUserMembershipWithUser) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(householdsUserMembershipTableColumns)

	for _, x := range memberships {
		rowValues := []driver.Value{
			&x.ArchivedAt,
			&x.ID,
			&x.BelongsToUser.ID,
			&x.BelongsToHousehold,
			x.HouseholdRole,
			&x.DefaultHousehold,
			&x.CreatedAt,
			&x.LastUpdatedAt,
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

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.BuildSessionContextDataForUser(ctx, "")
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

		args := []any{exampleUserID, true}

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

		args := []any{exampleUserID, true}

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

		markHouseholdAsUserDefaultArgs := []any{
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

		markHouseholdAsUserDefaultArgs := []any{
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

		userIsMemberOfHouseholdArgs := []any{
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

		userIsMemberOfHouseholdArgs := []any{
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

		fakeArgs := []any{
			exampleInput.NewRole,
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

		fakeArgs := []any{
			exampleInput.NewRole,
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(modifyUserPermissionsQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ModifyUserPermissions(ctx, exampleHouseholdID, exampleUserID, exampleInput))
	})
}

func TestSQLQuerier_addUserToHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleInput := fakes.BuildFakeHouseholdUserMembershipDatabaseCreationInput()

		c, db := buildTestClient(t)

		addUserToHouseholdArgs := []any{
			exampleInput.ID,
			exampleInput.UserID,
			exampleInput.HouseholdID,
			exampleInput.HouseholdRole,
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
		exampleHouseholdList := fakes.BuildFakeHouseholdList()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []any{
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(removeUserFromHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		query, args := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, nil)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, 0, exampleHouseholdList.Data...))

		markHouseholdAsUserDefaultArgs := []any{
			exampleUserID,
			exampleHouseholdList.Data[0].ID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(markHouseholdAsUserDefaultQuery)).
			WithArgs(interfaceToDriverValue(markHouseholdAsUserDefaultArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

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

	T.Run("with error creating transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUserID, exampleHouseholdID))
	})

	T.Run("with error fetching households", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []any{
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(removeUserFromHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		query, args := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, nil)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUserID, exampleHouseholdID))
	})

	T.Run("with error writing removal to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []any{
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(removeUserFromHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUserID, exampleHouseholdID))
	})

	T.Run("creates new household when none are left", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdList := fakes.BuildFakeHouseholdList()
		exampleHouseholdList.Data = []*types.Household{}

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []any{
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(removeUserFromHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		query, args := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, nil)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, 0, exampleHouseholdList.Data...))

		// create household user membership for created user
		householdCreationInput := &types.HouseholdDatabaseCreationInput{
			Name:          fmt.Sprintf("%s_default", exampleUserID),
			BelongsToUser: exampleUserID,
		}
		createHouseholdForNewUserArgs := []any{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactPhone,
			householdCreationInput.AddressLine1,
			householdCreationInput.AddressLine2,
			householdCreationInput.City,
			householdCreationInput.State,
			householdCreationInput.ZipCode,
			householdCreationInput.Country,
			householdCreationInput.Latitude,
			householdCreationInput.Longitude,
			householdCreationInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		// create household user membership for created user
		createHouseholdMembershipForNewUserArgs := []any{
			&idMatcher{},
			&idMatcher{},
			&idMatcher{},
			true,
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(createHouseholdMembershipForNewUserQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdMembershipForNewUserArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

		assert.NoError(t, c.RemoveUserFromHousehold(ctx, exampleUserID, exampleHouseholdID))
	})

	T.Run("with error creating new household when none are left", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdList := fakes.BuildFakeHouseholdList()
		exampleHouseholdList.Data = []*types.Household{}

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []any{
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(removeUserFromHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		query, args := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, nil)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, 0, exampleHouseholdList.Data...))

		// create household user membership for created user
		householdCreationInput := &types.HouseholdDatabaseCreationInput{
			Name:          fmt.Sprintf("%s_default", exampleUserID),
			BelongsToUser: exampleUserID,
		}
		createHouseholdForNewUserArgs := []any{
			&idMatcher{},
			householdCreationInput.Name,
			types.UnpaidHouseholdBillingStatus,
			householdCreationInput.ContactPhone,
			householdCreationInput.AddressLine1,
			householdCreationInput.AddressLine2,
			householdCreationInput.City,
			householdCreationInput.State,
			householdCreationInput.ZipCode,
			householdCreationInput.Country,
			householdCreationInput.Latitude,
			householdCreationInput.Longitude,
			householdCreationInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(createHouseholdForNewUserArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUserID, exampleHouseholdID))
	})

	T.Run("with error marking new household as user default", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdList := fakes.BuildFakeHouseholdList()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []any{
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(removeUserFromHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		query, args := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, nil)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, 0, exampleHouseholdList.Data...))

		markHouseholdAsUserDefaultArgs := []any{
			exampleUserID,
			exampleHouseholdList.Data[0].ID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(markHouseholdAsUserDefaultQuery)).
			WithArgs(interfaceToDriverValue(markHouseholdAsUserDefaultArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUserID, exampleHouseholdID))
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()
		exampleHouseholdList := fakes.BuildFakeHouseholdList()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []any{
			exampleHouseholdID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(removeUserFromHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		query, args := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, nil)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, 0, exampleHouseholdList.Data...))

		markHouseholdAsUserDefaultArgs := []any{
			exampleUserID,
			exampleHouseholdList.Data[0].ID,
			exampleUserID,
		}

		db.ExpectExec(formatQueryForSQLMock(markHouseholdAsUserDefaultQuery)).
			WithArgs(interfaceToDriverValue(markHouseholdAsUserDefaultArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUserID, exampleHouseholdID))
	})
}
