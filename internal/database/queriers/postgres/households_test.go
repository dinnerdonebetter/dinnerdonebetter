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

func buildMockRowsFromHouseholds(includeCounts bool, filteredCount uint64, households ...*types.Household) *sqlmock.Rows {
	columns := append(householdsTableColumns, householdsUserMembershipTableColumns...)

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range households {
		for _, y := range x.Members {
			rowValues := []driver.Value{
				x.ID,
				x.Name,
				x.BillingStatus,
				x.ContactEmail,
				x.ContactPhone,
				x.PaymentProcessorCustomerID,
				x.SubscriptionPlanID,
				x.CreatedOn,
				x.LastUpdatedOn,
				x.ArchivedOn,
				x.BelongsToUser,
				y.ID,
				y.BelongsToUser,
				y.BelongsToHousehold,
				strings.Join(y.HouseholdRoles, householdMemberRolesSeparator),
				y.DefaultHousehold,
				y.CreatedOn,
				x.LastUpdatedOn,
				y.ArchivedOn,
			}

			if includeCounts {
				rowValues = append(rowValues, filteredCount, len(households))
			}

			exampleRows.AddRow(rowValues...)
		}
	}

	return exampleRows
}

func TestQuerier_ScanHouseholds(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanHouseholds(ctx, mockRows, false)
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

		_, _, _, err := q.scanHouseholds(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_GetHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleUserID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(false, 0, exampleHousehold))

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, exampleHousehold, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		c, _ := buildTestClient(t)

		actual, err := c.GetHousehold(ctx, "", exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		c, _ := buildTestClient(t)

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleUserID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID, exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleUserID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID, exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with no returned households", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		c, db := buildTestClient(t)

		columns := append(householdsTableColumns, householdsUserMembershipTableColumns...)

		args := []interface{}{
			exampleUserID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows(columns))

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID, exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetAllHouseholdsCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getAllHouseholdsCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(exampleCount))

		actual, err := c.GetAllHouseholdsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getAllHouseholdsCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetAllHouseholdsCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetHouseholds(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdList := fakes.BuildFakeHouseholdList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Households...))

		actual, err := c.GetHouseholds(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholds(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdList := fakes.BuildFakeHouseholdList()
		exampleHouseholdList.Page, exampleHouseholdList.Limit = 0, 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Households...))

		actual, err := c.GetHouseholds(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleUserID := fakes.BuildFakeID()
		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHouseholds(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleUserID := fakes.BuildFakeID()
		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholds(ctx, exampleUserID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetHouseholdsForAdmin(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleHouseholdList := fakes.BuildFakeHouseholdList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdsQuery(ctx, "", true, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Households...))

		actual, err := c.GetHouseholdsForAdmin(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleHouseholdList := fakes.BuildFakeHouseholdList()
		exampleHouseholdList.Page, exampleHouseholdList.Limit = 0, 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdsQuery(ctx, "", true, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Households...))

		actual, err := c.GetHouseholdsForAdmin(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdsQuery(ctx, "", true, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHouseholdsForAdmin(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdsQuery(ctx, "", true, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholdsForAdmin(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BillingStatus = types.UnpaidHouseholdBillingStatus
		exampleHousehold.PaymentProcessorCustomerID = ""
		exampleHousehold.ID = ""
		exampleHousehold.BelongsToUser = exampleUserID
		exampleHousehold.Members = []*types.HouseholdUserMembership(nil)
		exampleInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		householdCreationArgs := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			types.UnpaidHouseholdBillingStatus,
			exampleInput.ContactEmail,
			exampleInput.ContactPhone,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHousehold.ID))

		addUserToHouseholdArgs := []interface{}{
			&idMatcher{},
			exampleInput.BelongsToUser,
			&idMatcher{},
			strings.Join([]string{authorization.HouseholdAdminRole.String()}, householdMemberRolesSeparator),
		}

		db.ExpectExec(formatQueryForSQLMock(addUserToHouseholdDuringCreationQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHousehold.ID))

		db.ExpectCommit()

		c.timeFunc = func() uint64 {
			return exampleHousehold.CreatedOn
		}

		actual, err := c.CreateHousehold(ctx, exampleInput)
		assert.NoError(t, err)
		assert.NotEmpty(t, actual.ID)
		actual.ID = ""

		assert.Equal(t, exampleHousehold, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateHousehold(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.ID = ""
		exampleHousehold.BelongsToUser = exampleUserID
		exampleHousehold.Members = []*types.HouseholdUserMembership(nil)
		exampleCreationInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.CreateHousehold(ctx, exampleCreationInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.ID = ""
		exampleHousehold.BelongsToUser = exampleUserID
		exampleHousehold.Members = []*types.HouseholdUserMembership{}
		exampleInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		householdCreationArgs := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			types.UnpaidHouseholdBillingStatus,
			exampleInput.ContactEmail,
			exampleInput.ContactPhone,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleHousehold.CreatedOn
		}

		actual, err := c.CreateHousehold(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error writing user addition to database", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.ID = ""
		exampleHousehold.BelongsToUser = exampleUserID
		exampleHousehold.Members = []*types.HouseholdUserMembership(nil)
		exampleInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		householdCreationArgs := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			types.UnpaidHouseholdBillingStatus,
			exampleInput.ContactEmail,
			exampleInput.ContactPhone,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHousehold.ID))

		addUserToHouseholdArgs := []interface{}{
			&idMatcher{},
			exampleInput.BelongsToUser,
			&idMatcher{},
			strings.Join([]string{authorization.HouseholdAdminRole.String()}, householdMemberRolesSeparator),
		}

		db.ExpectExec(formatQueryForSQLMock(addUserToHouseholdDuringCreationQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleHousehold.CreatedOn
		}

		actual, err := c.CreateHousehold(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.ID = ""
		exampleHousehold.BelongsToUser = exampleUserID
		exampleHousehold.Members = []*types.HouseholdUserMembership(nil)
		exampleInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		householdCreationArgs := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			types.UnpaidHouseholdBillingStatus,
			exampleInput.ContactEmail,
			exampleInput.ContactPhone,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHousehold.ID))

		addUserToHouseholdArgs := []interface{}{
			&idMatcher{},
			exampleInput.BelongsToUser,
			&idMatcher{},
			strings.Join([]string{authorization.HouseholdAdminRole.String()}, householdMemberRolesSeparator),
		}

		db.ExpectExec(formatQueryForSQLMock(addUserToHouseholdDuringCreationQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHousehold.ID))

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() uint64 {
			return exampleHousehold.CreatedOn
		}

		actual, err := c.CreateHousehold(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHousehold.Name,
			exampleHousehold.ContactEmail,
			exampleHousehold.ContactPhone,
			exampleHousehold.BelongsToUser,
			exampleHousehold.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHousehold.ID))

		assert.NoError(t, c.UpdateHousehold(ctx, exampleHousehold))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateHousehold(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUserID

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHousehold.Name,
			exampleHousehold.ContactEmail,
			exampleHousehold.ContactPhone,
			exampleHousehold.BelongsToUser,
			exampleHousehold.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateHousehold(ctx, exampleHousehold))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleUserID,
			exampleHouseholdID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleHouseholdID))

		assert.NoError(t, c.ArchiveHousehold(ctx, exampleHouseholdID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHousehold(ctx, "", exampleUserID))
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHousehold(ctx, exampleHouseholdID, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleUserID,
			exampleHouseholdID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveHousehold(ctx, exampleHouseholdID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
