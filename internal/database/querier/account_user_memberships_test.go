package querier

import (
	"context"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/authorization"
	"gitlab.com/prixfixe/prixfixe/internal/database"
	"gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromHouseholdUserMemberships(memberships ...*types.HouseholdUserMembership) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(querybuilding.HouseholdsUserMembershipTableColumns)

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
	exampleRows := sqlmock.NewRows(querybuilding.HouseholdsUserMembershipTableColumns)

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

		examplePermsMap := map[uint64]*types.UserHouseholdMembershipInfo{}
		for _, membership := range exampleHousehold.Members {
			examplePermsMap[membership.BelongsToHousehold] = &types.UserHouseholdMembershipInfo{
				HouseholdName:  exampleHousehold.Name,
				HouseholdID:    membership.BelongsToHousehold,
				HouseholdRoles: membership.HouseholdRoles,
			}
		}

		exampleHouseholdPermissionsMap := map[uint64]authorization.HouseholdRolePermissionsChecker{}
		for _, membership := range exampleHousehold.Members {
			exampleHouseholdPermissionsMap[membership.BelongsToHousehold] = authorization.NewHouseholdRolePermissionChecker(membership.HouseholdRoles...)
		}

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeUserRetrievalQuery, fakeUserRetrievalArgs := fakes.BuildFakeSQLQuery()

		mockQueryBuilder.UserSQLQueryBuilder.On(
			"BuildGetUserQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
		).Return(fakeUserRetrievalQuery, fakeUserRetrievalArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeUserRetrievalQuery)).
			WithArgs(interfaceToDriverValue(fakeUserRetrievalArgs)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		fakeHouseholdMembershipsQuery, fakeHouseholdMembershipsArgs := fakes.BuildFakeSQLQuery()

		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildGetHouseholdMembershipsForUserQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
		).Return(fakeHouseholdMembershipsQuery, fakeHouseholdMembershipsArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeHouseholdMembershipsQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdMembershipsArgs)...).
			WillReturnRows(buildMockRowsFromHouseholdUserMemberships(exampleHousehold.Members...))

		c.sqlQueryBuilder = mockQueryBuilder

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

		actual, err := c.BuildSessionContextDataForUser(ctx, 0)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error retrieving user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		examplePermsMap := map[uint64]*types.UserHouseholdMembershipInfo{}
		for _, membership := range exampleHousehold.Members {
			examplePermsMap[membership.BelongsToHousehold] = &types.UserHouseholdMembershipInfo{
				HouseholdName:  exampleHousehold.Name,
				HouseholdID:    membership.BelongsToHousehold,
				HouseholdRoles: membership.HouseholdRoles,
			}
		}

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeUserRetrievalQuery, fakeUserRetrievalArgs := fakes.BuildFakeSQLQuery()

		mockQueryBuilder.UserSQLQueryBuilder.On(
			"BuildGetUserQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
		).Return(fakeUserRetrievalQuery, fakeUserRetrievalArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeUserRetrievalQuery)).
			WithArgs(interfaceToDriverValue(fakeUserRetrievalArgs)...).
			WillReturnError(errors.New("blah"))

		c.sqlQueryBuilder = mockQueryBuilder

		actual, err := c.BuildSessionContextDataForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error retrieving household memberships", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		examplePermsMap := map[uint64]*types.UserHouseholdMembershipInfo{}
		for _, membership := range exampleHousehold.Members {
			examplePermsMap[membership.BelongsToHousehold] = &types.UserHouseholdMembershipInfo{
				HouseholdName:  exampleHousehold.Name,
				HouseholdID:    membership.BelongsToHousehold,
				HouseholdRoles: membership.HouseholdRoles,
			}
		}

		exampleHouseholdPermissionsMap := map[uint64]authorization.HouseholdRolePermissionsChecker{}
		for _, membership := range exampleHousehold.Members {
			exampleHouseholdPermissionsMap[membership.BelongsToHousehold] = authorization.NewHouseholdRolePermissionChecker(membership.HouseholdRoles...)
		}

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeUserRetrievalQuery, fakeUserRetrievalArgs := fakes.BuildFakeSQLQuery()

		mockQueryBuilder.UserSQLQueryBuilder.On(
			"BuildGetUserQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
		).Return(fakeUserRetrievalQuery, fakeUserRetrievalArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeUserRetrievalQuery)).
			WithArgs(interfaceToDriverValue(fakeUserRetrievalArgs)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		fakeHouseholdMembershipsQuery, fakeHouseholdMembershipsArgs := fakes.BuildFakeSQLQuery()

		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildGetHouseholdMembershipsForUserQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
		).Return(fakeHouseholdMembershipsQuery, fakeHouseholdMembershipsArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeHouseholdMembershipsQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdMembershipsArgs)...).
			WillReturnError(errors.New("blah"))

		c.sqlQueryBuilder = mockQueryBuilder

		actual, err := c.BuildSessionContextDataForUser(ctx, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error scanning household user memberships", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		examplePermsMap := map[uint64]*types.UserHouseholdMembershipInfo{}
		for _, membership := range exampleHousehold.Members {
			examplePermsMap[membership.BelongsToHousehold] = &types.UserHouseholdMembershipInfo{
				HouseholdName:  exampleHousehold.Name,
				HouseholdID:    membership.BelongsToHousehold,
				HouseholdRoles: membership.HouseholdRoles,
			}
		}

		exampleHouseholdPermissionsMap := map[uint64]authorization.HouseholdRolePermissionsChecker{}
		for _, membership := range exampleHousehold.Members {
			exampleHouseholdPermissionsMap[membership.BelongsToHousehold] = authorization.NewHouseholdRolePermissionChecker(membership.HouseholdRoles...)
		}

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeUserRetrievalQuery, fakeUserRetrievalArgs := fakes.BuildFakeSQLQuery()

		mockQueryBuilder.UserSQLQueryBuilder.On(
			"BuildGetUserQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
		).Return(fakeUserRetrievalQuery, fakeUserRetrievalArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeUserRetrievalQuery)).
			WithArgs(interfaceToDriverValue(fakeUserRetrievalArgs)...).
			WillReturnRows(buildMockRowsFromUsers(false, 0, exampleUser))

		fakeHouseholdMembershipsQuery, fakeHouseholdMembershipsArgs := fakes.BuildFakeSQLQuery()

		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildGetHouseholdMembershipsForUserQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
		).Return(fakeHouseholdMembershipsQuery, fakeHouseholdMembershipsArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeHouseholdMembershipsQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdMembershipsArgs)...).
			WillReturnRows(buildInvalidMockRowsFromHouseholdUserMemberships(exampleHousehold.Members...))

		c.sqlQueryBuilder = mockQueryBuilder

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
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		expected := exampleHousehold.ID

		c, db := buildTestClient(t)

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildGetDefaultHouseholdIDForUserQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(exampleHousehold.ID))

		actual, err := c.GetDefaultHouseholdIDForUser(ctx, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, db.ExpectationsWereMet())
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetDefaultHouseholdIDForUser(ctx, 0)
		assert.Error(t, err)
		assert.Zero(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()

		c, db := buildTestClient(t)

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildGetDefaultHouseholdIDForUserQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetDefaultHouseholdIDForUser(ctx, exampleUser.ID)
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
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildMarkHouseholdAsUserDefaultQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		assert.NoError(t, c.MarkHouseholdAsUserDefault(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID))
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkHouseholdAsUserDefault(ctx, 0, exampleHousehold.ID, exampleUser.ID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()

		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkHouseholdAsUserDefault(ctx, exampleUser.ID, 0, exampleUser.ID))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.MarkHouseholdAsUserDefault(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID))
	})

	T.Run("with error marking household as default", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildMarkHouseholdAsUserDefaultQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.MarkHouseholdAsUserDefault(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID))
	})

	T.Run("with error writing audit log entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildMarkHouseholdAsUserDefaultQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(newDatabaseResultForID(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.MarkHouseholdAsUserDefault(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID))
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildMarkHouseholdAsUserDefaultQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(newDatabaseResultForID(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		assert.Error(t, c.MarkHouseholdAsUserDefault(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID))
	})
}

func TestQuerier_UserIsMemberOfHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildUserIsMemberOfHouseholdQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(sqlmock.NewRows([]string{"result"}).AddRow(true))

		actual, err := c.UserIsMemberOfHousehold(ctx, exampleUser.ID, exampleHousehold.ID)
		assert.True(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, _ := buildTestClient(t)

		actual, err := c.UserIsMemberOfHousehold(ctx, 0, exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()

		c, _ := buildTestClient(t)

		actual, err := c.UserIsMemberOfHousehold(ctx, exampleUser.ID, 0)
		assert.False(t, actual)
		assert.Error(t, err)
	})

	T.Run("with error performing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildUserIsMemberOfHouseholdQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.UserIsMemberOfHousehold(ctx, exampleUser.ID, exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)
	})
}

func TestQuerier_ModifyUserPermissions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildModifyUserPermissionsQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
			exampleInput.NewRoles,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		assert.NoError(t, c.ModifyUserPermissions(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		c, _ := buildTestClient(t)

		assert.Error(t, c.ModifyUserPermissions(ctx, 0, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})

	T.Run("with invalid household id", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		c, _ := buildTestClient(t)

		assert.Error(t, c.ModifyUserPermissions(ctx, exampleUser.ID, 0, exampleUser.ID, exampleInput))
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, _ := buildTestClient(t)

		assert.Error(t, c.ModifyUserPermissions(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, nil))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.ModifyUserPermissions(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildModifyUserPermissionsQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
			exampleInput.NewRoles,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.ModifyUserPermissions(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})

	T.Run("with error writing audit log entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildModifyUserPermissionsQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
			exampleInput.NewRoles,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.ModifyUserPermissions(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeUserPermissionModificationInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildModifyUserPermissionsQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
			exampleInput.NewRoles,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		assert.Error(t, c.ModifyUserPermissions(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})
}

func TestQuerier_TransferHouseholdOwnership(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeHouseholdTransferQuery, fakeHouseholdTransferArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildTransferHouseholdOwnershipQuery",
			testutils.ContextMatcher,
			exampleInput.CurrentOwner,
			exampleInput.NewOwner,
			exampleHousehold.ID,
		).Return(fakeHouseholdTransferQuery, fakeHouseholdTransferArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdTransferQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdTransferArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		fakeHouseholdMembershipsTransferQuery, fakeHouseholdMembershipsTransferArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildTransferHouseholdMembershipsQuery",
			testutils.ContextMatcher,
			exampleInput.CurrentOwner,
			exampleInput.NewOwner,
			exampleHousehold.ID,
		).Return(fakeHouseholdMembershipsTransferQuery, fakeHouseholdMembershipsTransferArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdMembershipsTransferQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdMembershipsTransferArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.NoError(t, c.TransferHouseholdOwnership(ctx, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, _ := buildTestClient(t)

		assert.Error(t, c.TransferHouseholdOwnership(ctx, 0, exampleUser.ID, exampleInput))
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, _ := buildTestClient(t)

		assert.Error(t, c.TransferHouseholdOwnership(ctx, exampleHousehold.ID, exampleUser.ID, nil))
	})

	T.Run("with error starting transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.TransferHouseholdOwnership(ctx, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})

	T.Run("with error writing household transfer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeHouseholdTransferQuery, fakeHouseholdTransferArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildTransferHouseholdOwnershipQuery",
			testutils.ContextMatcher,
			exampleInput.CurrentOwner,
			exampleInput.NewOwner,
			exampleHousehold.ID,
		).Return(fakeHouseholdTransferQuery, fakeHouseholdTransferArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdTransferQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdTransferArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.TransferHouseholdOwnership(ctx, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})

	T.Run("with error writing membership transfers", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeHouseholdTransferQuery, fakeHouseholdTransferArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildTransferHouseholdOwnershipQuery",
			testutils.ContextMatcher,
			exampleInput.CurrentOwner,
			exampleInput.NewOwner,
			exampleHousehold.ID,
		).Return(fakeHouseholdTransferQuery, fakeHouseholdTransferArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdTransferQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdTransferArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		fakeHouseholdMembershipsTransferQuery, fakeHouseholdMembershipsTransferArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildTransferHouseholdMembershipsQuery",
			testutils.ContextMatcher,
			exampleInput.CurrentOwner,
			exampleInput.NewOwner,
			exampleHousehold.ID,
		).Return(fakeHouseholdMembershipsTransferQuery, fakeHouseholdMembershipsTransferArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdMembershipsTransferQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdMembershipsTransferArgs)...).
			WillReturnError(errors.New("blah"))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.TransferHouseholdOwnership(ctx, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})

	T.Run("with error writing membership transfers audit log entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeHouseholdTransferQuery, fakeHouseholdTransferArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildTransferHouseholdOwnershipQuery",
			testutils.ContextMatcher,
			exampleInput.CurrentOwner,
			exampleInput.NewOwner,
			exampleHousehold.ID,
		).Return(fakeHouseholdTransferQuery, fakeHouseholdTransferArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdTransferQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdTransferArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		fakeHouseholdMembershipsTransferQuery, fakeHouseholdMembershipsTransferArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildTransferHouseholdMembershipsQuery",
			testutils.ContextMatcher,
			exampleInput.CurrentOwner,
			exampleInput.NewOwner,
			exampleHousehold.ID,
		).Return(fakeHouseholdMembershipsTransferQuery, fakeHouseholdMembershipsTransferArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdMembershipsTransferQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdMembershipsTransferArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.TransferHouseholdOwnership(ctx, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleInput := fakes.BuildFakeTransferHouseholdOwnershipInput()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeHouseholdTransferQuery, fakeHouseholdTransferArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildTransferHouseholdOwnershipQuery",
			testutils.ContextMatcher,
			exampleInput.CurrentOwner,
			exampleInput.NewOwner,
			exampleHousehold.ID,
		).Return(fakeHouseholdTransferQuery, fakeHouseholdTransferArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdTransferQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdTransferArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		fakeHouseholdMembershipsTransferQuery, fakeHouseholdMembershipsTransferArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildTransferHouseholdMembershipsQuery",
			testutils.ContextMatcher,
			exampleInput.CurrentOwner,
			exampleInput.NewOwner,
			exampleHousehold.ID,
		).Return(fakeHouseholdMembershipsTransferQuery, fakeHouseholdMembershipsTransferArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdMembershipsTransferQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdMembershipsTransferArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.TransferHouseholdOwnership(ctx, exampleHousehold.ID, exampleUser.ID, exampleInput))
	})
}

func TestQuerier_AddUserToHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHouseholdUserMembership := fakes.BuildFakeHouseholdUserMembership()
		exampleHouseholdUserMembership.BelongsToHousehold = exampleHousehold.ID

		exampleInput := &types.AddUserToHouseholdInput{
			Reason:         t.Name(),
			HouseholdID:    exampleHousehold.ID,
			UserID:         exampleHousehold.BelongsToUser,
			HouseholdRoles: []string{householdMemberRolesSeparator},
		}

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildAddUserToHouseholdQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHouseholdUserMembership.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.NoError(t, c.AddUserToHousehold(ctx, exampleInput, exampleUser.ID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid actor ID", func(t *testing.T) {
		t.Parallel()

		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHouseholdUserMembership := fakes.BuildFakeHouseholdUserMembership()
		exampleHouseholdUserMembership.BelongsToHousehold = exampleHousehold.ID

		exampleInput := &types.AddUserToHouseholdInput{
			Reason:         t.Name(),
			HouseholdID:    exampleHousehold.ID,
			UserID:         exampleHousehold.BelongsToUser,
			HouseholdRoles: []string{householdMemberRolesSeparator},
		}

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.AddUserToHousehold(ctx, exampleInput, 0))
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHouseholdUserMembership := fakes.BuildFakeHouseholdUserMembership()
		exampleHouseholdUserMembership.BelongsToHousehold = exampleHousehold.ID

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.AddUserToHousehold(ctx, nil, exampleUser.ID))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHouseholdUserMembership := fakes.BuildFakeHouseholdUserMembership()
		exampleHouseholdUserMembership.BelongsToHousehold = exampleHousehold.ID

		exampleInput := &types.AddUserToHouseholdInput{
			Reason:         t.Name(),
			HouseholdID:    exampleHousehold.ID,
			UserID:         exampleHousehold.BelongsToUser,
			HouseholdRoles: []string{householdMemberRolesSeparator},
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.AddUserToHousehold(ctx, exampleInput, exampleUser.ID))
	})

	T.Run("with error writing add query", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHouseholdUserMembership := fakes.BuildFakeHouseholdUserMembership()
		exampleHouseholdUserMembership.BelongsToHousehold = exampleHousehold.ID

		exampleInput := &types.AddUserToHouseholdInput{
			Reason:         t.Name(),
			HouseholdID:    exampleHousehold.ID,
			UserID:         exampleHousehold.BelongsToUser,
			HouseholdRoles: []string{householdMemberRolesSeparator},
		}

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildAddUserToHouseholdQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.AddUserToHousehold(ctx, exampleInput, exampleUser.ID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error writing audit log entry", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHouseholdUserMembership := fakes.BuildFakeHouseholdUserMembership()
		exampleHouseholdUserMembership.BelongsToHousehold = exampleHousehold.ID

		exampleInput := &types.AddUserToHouseholdInput{
			Reason:         t.Name(),
			HouseholdID:    exampleHousehold.ID,
			UserID:         exampleHousehold.BelongsToUser,
			HouseholdRoles: []string{householdMemberRolesSeparator},
		}

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildAddUserToHouseholdQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHouseholdUserMembership.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.AddUserToHousehold(ctx, exampleInput, exampleUser.ID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHouseholdUserMembership := fakes.BuildFakeHouseholdUserMembership()
		exampleHouseholdUserMembership.BelongsToHousehold = exampleHousehold.ID

		exampleInput := &types.AddUserToHouseholdInput{
			Reason:         t.Name(),
			HouseholdID:    exampleHousehold.ID,
			UserID:         exampleHousehold.BelongsToUser,
			HouseholdRoles: []string{householdMemberRolesSeparator},
		}

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildAddUserToHouseholdQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHouseholdUserMembership.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.AddUserToHousehold(ctx, exampleInput, exampleUser.ID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_RemoveUserFromHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildRemoveUserFromHouseholdQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		assert.NoError(t, c.RemoveUserFromHousehold(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, t.Name()))
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, _ := buildTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(ctx, 0, exampleHousehold.ID, exampleUser.ID, t.Name()))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()

		c, _ := buildTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUser.ID, 0, exampleUser.ID, t.Name()))
	})

	T.Run("with invalid actor ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, _ := buildTestClient(t)

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUser.ID, exampleHousehold.ID, 0, t.Name()))
	})

	T.Run("with empty reason", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildRemoveUserFromHouseholdQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, ""))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, t.Name()))
	})

	T.Run("with error writing removal to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildRemoveUserFromHouseholdQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, t.Name()))
	})

	T.Run("with error writing audit log entry", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildRemoveUserFromHouseholdQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, t.Name()))
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildRemoveUserFromHouseholdQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		assert.Error(t, c.RemoveUserFromHousehold(ctx, exampleUser.ID, exampleHousehold.ID, exampleUser.ID, t.Name()))
	})
}
