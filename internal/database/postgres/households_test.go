package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/prixfixeco/backend/internal/authorization"
	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/pkg/pointers"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
				x.ContactPhone,
				x.PaymentProcessorCustomerID,
				x.SubscriptionPlanID,
				x.TimeZone,
				x.CreatedAt,
				x.LastUpdatedAt,
				x.ArchivedAt,
				x.BelongsToUser,
				&y.BelongsToUser.ID,
				&y.BelongsToUser.FirstName,
				&y.BelongsToUser.LastName,
				&y.BelongsToUser.Username,
				&y.BelongsToUser.EmailAddress,
				&y.BelongsToUser.EmailAddressVerifiedAt,
				&y.BelongsToUser.AvatarSrc,
				&y.BelongsToUser.RequiresPasswordChange,
				&y.BelongsToUser.PasswordLastChangedAt,
				&y.BelongsToUser.TwoFactorSecretVerifiedAt,
				y.BelongsToUser.ServiceRole,
				&y.BelongsToUser.AccountStatus,
				&y.BelongsToUser.AccountStatusExplanation,
				&y.BelongsToUser.Birthday,
				&y.BelongsToUser.CreatedAt,
				&y.BelongsToUser.LastUpdatedAt,
				&y.BelongsToUser.ArchivedAt,
				y.ID,
				y.BelongsToUser.ID,
				y.BelongsToHousehold,
				y.HouseholdRole,
				y.DefaultHousehold,
				y.CreatedAt,
				x.LastUpdatedAt,
				y.ArchivedAt,
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

		args := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(false, 0, exampleHousehold))

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID)
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

		actual, err := c.GetHousehold(ctx, "")
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

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID)
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

		args := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID)
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

		args := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID)
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

		args := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows(columns))

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetHouseholdByID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		args := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(false, 0, exampleHousehold))

		actual, err := c.GetHouseholdByID(ctx, exampleHousehold.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleHousehold, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholdByID(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		args := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholdByID(ctx, exampleHousehold.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		args := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHouseholdByID(ctx, exampleHousehold.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with no returned households", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		columns := append(householdsTableColumns, householdsUserMembershipTableColumns...)

		args := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows(columns))

		actual, err := c.GetHouseholdByID(ctx, exampleHousehold.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

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
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Data...))

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
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Data...))

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
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdList := fakes.BuildFakeHouseholdList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdsQuery(ctx, "", true, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Data...))

		actual, err := c.GetHouseholdsForAdmin(ctx, exampleUserID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleUserID := fakes.BuildFakeID()
		exampleHouseholdList := fakes.BuildFakeHouseholdList()
		exampleHouseholdList.Page, exampleHouseholdList.Limit = 0, 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetHouseholdsQuery(ctx, "", true, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Data...))

		actual, err := c.GetHouseholdsForAdmin(ctx, exampleUserID, filter)
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

		query, args := c.buildGetHouseholdsQuery(ctx, "", true, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHouseholdsForAdmin(ctx, exampleUserID, filter)
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

		query, args := c.buildGetHouseholdsQuery(ctx, "", true, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholdsForAdmin(ctx, exampleUserID, filter)
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
		exampleHousehold.BillingStatus = string(types.UnpaidHouseholdBillingStatus)
		exampleHousehold.PaymentProcessorCustomerID = ""
		exampleHousehold.BelongsToUser = exampleUserID
		exampleHousehold.Members = []*types.HouseholdUserMembershipWithUser(nil)
		exampleInput := converters.ConvertHouseholdToHouseholdDatabaseCreationInput(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		householdCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Name,
			types.UnpaidHouseholdBillingStatus,
			exampleInput.ContactPhone,
			exampleHousehold.TimeZone,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		addUserToHouseholdArgs := []any{
			&idMatcher{},
			exampleInput.BelongsToUser,
			&idMatcher{},
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(addUserToHouseholdQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

		c.timeFunc = func() time.Time {
			return exampleHousehold.CreatedAt
		}

		actual, err := c.CreateHousehold(ctx, exampleInput)
		assert.NoError(t, err)
		assert.NotEmpty(t, actual.ID)

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
		exampleHousehold.Members = []*types.HouseholdUserMembershipWithUser(nil)
		exampleCreationInput := converters.ConvertHouseholdToHouseholdDatabaseCreationInput(exampleHousehold)

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
		exampleHousehold.Members = []*types.HouseholdUserMembershipWithUser{}
		exampleInput := converters.ConvertHouseholdToHouseholdDatabaseCreationInput(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		householdCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Name,
			types.UnpaidHouseholdBillingStatus,
			exampleInput.ContactPhone,
			exampleHousehold.TimeZone,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() time.Time {
			return exampleHousehold.CreatedAt
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
		exampleHousehold.BelongsToUser = exampleUserID
		exampleHousehold.Members = []*types.HouseholdUserMembershipWithUser(nil)
		exampleInput := converters.ConvertHouseholdToHouseholdDatabaseCreationInput(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		householdCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Name,
			types.UnpaidHouseholdBillingStatus,
			exampleInput.ContactPhone,
			exampleHousehold.TimeZone,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		addUserToHouseholdArgs := []any{
			&idMatcher{},
			exampleInput.BelongsToUser,
			&idMatcher{},
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(addUserToHouseholdQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() time.Time {
			return exampleHousehold.CreatedAt
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
		exampleHousehold.BelongsToUser = exampleUserID
		exampleHousehold.Members = []*types.HouseholdUserMembershipWithUser(nil)
		exampleInput := converters.ConvertHouseholdToHouseholdDatabaseCreationInput(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		householdCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Name,
			types.UnpaidHouseholdBillingStatus,
			exampleInput.ContactPhone,
			exampleHousehold.TimeZone,
			exampleInput.BelongsToUser,
		}

		db.ExpectExec(formatQueryForSQLMock(householdCreationQuery)).
			WithArgs(interfaceToDriverValue(householdCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		addUserToHouseholdArgs := []any{
			&idMatcher{},
			exampleInput.BelongsToUser,
			&idMatcher{},
			authorization.HouseholdAdminRole.String(),
		}

		db.ExpectExec(formatQueryForSQLMock(addUserToHouseholdQuery)).
			WithArgs(interfaceToDriverValue(addUserToHouseholdArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() time.Time {
			return exampleHousehold.CreatedAt
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

		args := []any{
			exampleHousehold.Name,
			exampleHousehold.ContactPhone,
			exampleHousehold.TimeZone,
			exampleHousehold.BelongsToUser,
			exampleHousehold.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

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

		args := []any{
			exampleHousehold.Name,
			exampleHousehold.ContactPhone,
			exampleHousehold.TimeZone,
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

		args := []any{
			exampleUserID,
			exampleHouseholdID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveHouseholdQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

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

		args := []any{
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

func TestQuerier_buildGetHouseholdsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		expectedQuery := "SELECT households.id, households.name, households.billing_status, households.contact_phone, households.payment_processor_customer_id, households.subscription_plan_id, households.time_zone, households.created_at, households.last_updated_at, households.archived_at, households.belongs_to_user, users.id, users.first_name, users.last_name, users.username, users.email_address, users.email_address_verified_at, users.avatar_src, users.requires_password_change, users.password_last_changed_at, users.two_factor_secret_verified_at, users.service_role, users.user_account_status, users.user_account_status_explanation, users.birthday, users.created_at, users.last_updated_at, users.archived_at, household_user_memberships.id, household_user_memberships.belongs_to_user, household_user_memberships.belongs_to_household, household_user_memberships.household_role, household_user_memberships.default_household, household_user_memberships.created_at, household_user_memberships.last_updated_at, household_user_memberships.archived_at, (SELECT COUNT(households.id) FROM households WHERE households.archived_at IS NULL AND households.belongs_to_user = $1) as filtered_count, (SELECT COUNT(households.id) FROM households WHERE households.archived_at IS NULL AND households.belongs_to_user = $2) as total_count FROM households JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id JOIN users ON household_user_memberships.belongs_to_user = users.id WHERE household_user_memberships.archived_at IS NULL AND household_user_memberships.belongs_to_user = $3 AND households.archived_at IS NULL GROUP BY households.id, users.id, household_user_memberships.id LIMIT 20"
		expectedArgs := []any{
			exampleUserID,
			exampleUserID,
			exampleUserID,
		}

		actualQuery, actualArgs := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, types.DefaultQueryFilter())

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})

	T.Run("include archived", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		expectedQuery := "SELECT households.id, households.name, households.billing_status, households.contact_phone, households.payment_processor_customer_id, households.subscription_plan_id, households.time_zone, households.created_at, households.last_updated_at, households.archived_at, households.belongs_to_user, users.id, users.first_name, users.last_name, users.username, users.email_address, users.email_address_verified_at, users.avatar_src, users.requires_password_change, users.password_last_changed_at, users.two_factor_secret_verified_at, users.service_role, users.user_account_status, users.user_account_status_explanation, users.birthday, users.created_at, users.last_updated_at, users.archived_at, household_user_memberships.id, household_user_memberships.belongs_to_user, household_user_memberships.belongs_to_household, household_user_memberships.household_role, household_user_memberships.default_household, household_user_memberships.created_at, household_user_memberships.last_updated_at, household_user_memberships.archived_at, (SELECT COUNT(households.id) FROM households WHERE households.archived_at IS NULL AND households.belongs_to_user = $1) as filtered_count, (SELECT COUNT(households.id) FROM households WHERE households.archived_at IS NULL AND households.belongs_to_user = $2) as total_count FROM households JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id JOIN users ON household_user_memberships.belongs_to_user = users.id WHERE household_user_memberships.archived_at IS NULL AND household_user_memberships.belongs_to_user = $3 AND households.archived_at IS NULL GROUP BY households.id, users.id, household_user_memberships.id LIMIT 20"
		expectedArgs := []any{
			exampleUserID,
			exampleUserID,
			exampleUserID,
		}

		filter := types.DefaultQueryFilter()
		filter.IncludeArchived = pointers.Pointer(true)

		actualQuery, actualArgs := c.buildGetHouseholdsQuery(ctx, exampleUserID, false, filter)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}
