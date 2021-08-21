package querier

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"testing"
	"time"

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

func buildMockRowsFromHouseholds(includeCounts bool, filteredCount uint64, households ...*types.Household) *sqlmock.Rows {
	columns := append(querybuilding.HouseholdsTableColumns, querybuilding.HouseholdsUserMembershipTableColumns...)

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range households {
		for _, y := range x.Members {
			rowValues := []driver.Value{
				x.ID,
				x.ExternalID,
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

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold.ID, exampleUser.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromHouseholds(false, 0, exampleHousehold))

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleHousehold, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		c, _ := buildTestClient(t)

		actual, err := c.GetHousehold(ctx, 0, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		c, _ := buildTestClient(t)

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID, 0)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error reading from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold.ID, exampleUser.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold.ID, exampleUser.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with no returned households", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold.ID, exampleUser.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		columns := append(querybuilding.HouseholdsTableColumns, querybuilding.HouseholdsUserMembershipTableColumns...)

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(sqlmock.NewRows(columns))

		actual, err := c.GetHousehold(ctx, exampleHousehold.ID, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetAllHouseholdsCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetAllHouseholdsCountQuery", testutils.ContextMatcher).Return(fakeQuery)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(exampleCount))

		actual, err := c.GetAllHouseholdsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetAllHouseholds(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		results := make(chan []*types.Household)
		expectedCount := uint64(20)
		doneChan := make(chan bool, 1)
		exampleHouseholdList := fakes.BuildFakeHouseholdList()

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetAllHouseholdsCountQuery", testutils.ContextMatcher).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetBatchOfHouseholdsQuery",
			testutils.ContextMatcher,
			uint64(1), uint64(defaultBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnRows(buildMockRowsFromHouseholds(false, 0, exampleHouseholdList.Households...))

		err := c.GetAllHouseholds(ctx, results, 0)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-results:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with nil results channel", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		assert.Error(t, c.GetAllHouseholds(ctx, nil, 0))
	})

	T.Run("with now rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		results := make(chan []*types.Household)
		expectedCount := uint64(20)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetAllHouseholdsCountQuery", testutils.ContextMatcher).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetBatchOfHouseholdsQuery",
			testutils.ContextMatcher,
			uint64(1), uint64(exampleBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnError(sql.ErrNoRows)

		err := c.GetAllHouseholds(ctx, results, exampleBatchSize)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		results := make(chan []*types.Household)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetAllHouseholdsCountQuery", testutils.ContextMatcher).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		c.sqlQueryBuilder = mockQueryBuilder

		err := c.GetAllHouseholds(ctx, results, exampleBatchSize)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		results := make(chan []*types.Household)
		expectedCount := uint64(20)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetAllHouseholdsCountQuery", testutils.ContextMatcher).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetBatchOfHouseholdsQuery",
			testutils.ContextMatcher,
			uint64(1), uint64(exampleBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnError(errors.New("blah"))

		err := c.GetAllHouseholds(ctx, results, exampleBatchSize)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		results := make(chan []*types.Household)
		expectedCount := uint64(20)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetAllHouseholdsCountQuery", testutils.ContextMatcher).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetBatchOfHouseholdsQuery",
			testutils.ContextMatcher,
			uint64(1), uint64(exampleBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		err := c.GetAllHouseholds(ctx, results, exampleBatchSize)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetHouseholds(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleUser := fakes.BuildFakeUser()
		exampleHouseholdList := fakes.BuildFakeHouseholdList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdsQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Households...))

		actual, err := c.GetHouseholds(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdList, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetHouseholds(ctx, 0, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleUser := fakes.BuildFakeUser()
		exampleHouseholdList := fakes.BuildFakeHouseholdList()
		exampleHouseholdList.Page, exampleHouseholdList.Limit = 0, 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdsQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Households...))

		actual, err := c.GetHouseholds(ctx, exampleUser.ID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdList, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleUser := fakes.BuildFakeUser()
		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdsQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHouseholds(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleUser := fakes.BuildFakeUser()
		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdsQuery",
			testutils.ContextMatcher,
			exampleUser.ID,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholds(ctx, exampleUser.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
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

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdsQuery",
			testutils.ContextMatcher,
			uint64(0),
			true,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Households...))

		actual, err := c.GetHouseholdsForAdmin(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdList, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleHouseholdList := fakes.BuildFakeHouseholdList()
		exampleHouseholdList.Page, exampleHouseholdList.Limit = 0, 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdsQuery",
			testutils.ContextMatcher,
			uint64(0),
			true,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromHouseholds(true, exampleHouseholdList.FilteredCount, exampleHouseholdList.Households...))

		actual, err := c.GetHouseholdsForAdmin(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleHouseholdList, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdsQuery",
			testutils.ContextMatcher,
			uint64(0),
			true,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetHouseholdsForAdmin(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetHouseholdsQuery",
			testutils.ContextMatcher,
			uint64(0),
			true,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetHouseholdsForAdmin(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_CreateHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.ExternalID = ""
		exampleHousehold.BillingStatus = types.UnpaidHouseholdBillingStatus
		exampleHousehold.PaymentProcessorCustomerID = ""
		exampleHousehold.BelongsToUser = exampleUser.ID
		exampleHousehold.Members = []*types.HouseholdUserMembership(nil)
		exampleCreationInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)
		exampleHouseholdAdditionInput := &types.AddUserToHouseholdInput{
			Reason:         "household creation",
			UserID:         exampleUser.ID,
			HouseholdID:    exampleHousehold.ID,
			HouseholdRoles: []string{authorization.HouseholdAdminRole.String()},
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeCreationQuery, fakeCreationArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildHouseholdCreationQuery",
			testutils.ContextMatcher,
			exampleCreationInput,
		).Return(fakeCreationQuery, fakeCreationArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeCreationQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationArgs)...).
			WillReturnRows(newDatabaseResultForID(exampleHousehold.ID))

		fakeCreationAuditLogEntryQuery, fakeCreationAuditLogEntryArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.AuditLogEntrySQLQueryBuilder.On(
			"BuildCreateAuditLogEntryQuery",
			testutils.ContextMatcher,
			mock.IsType(&types.AuditLogEntryCreationInput{}),
		).Return(fakeCreationAuditLogEntryQuery, fakeCreationAuditLogEntryArgs).Once()

		db.ExpectExec(formatQueryForSQLMock(fakeCreationAuditLogEntryQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationAuditLogEntryArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(123))

		fakeHouseholdAdditionQuery, fakeHouseholdAdditionArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildAddUserToHouseholdQuery",
			testutils.ContextMatcher,
			exampleHouseholdAdditionInput,
		).Return(fakeHouseholdAdditionQuery, fakeHouseholdAdditionArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdAdditionQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdAdditionArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		fakeHouseholdAdditionAuditLogEntryQuery, fakeHouseholdAdditionAuditLogEntryArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.AuditLogEntrySQLQueryBuilder.On(
			"BuildCreateAuditLogEntryQuery",
			testutils.ContextMatcher,
			mock.IsType(&types.AuditLogEntryCreationInput{}),
		).Return(fakeHouseholdAdditionAuditLogEntryQuery, fakeHouseholdAdditionAuditLogEntryArgs).Once()

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdAdditionAuditLogEntryQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdAdditionAuditLogEntryArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(123))

		db.ExpectCommit()

		c.timeFunc = func() uint64 {
			return exampleHousehold.CreatedOn
		}
		c.sqlQueryBuilder = mockQueryBuilder

		actual, err := c.CreateHousehold(ctx, exampleCreationInput, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleHousehold, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid actor ID", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.ExternalID = ""
		exampleHousehold.BelongsToUser = exampleUser.ID
		exampleHousehold.Members = []*types.HouseholdUserMembership(nil)
		exampleCreationInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateHousehold(ctx, exampleCreationInput, 0)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateHousehold(ctx, nil, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.ExternalID = ""
		exampleHousehold.BelongsToUser = exampleUser.ID
		exampleHousehold.Members = []*types.HouseholdUserMembership(nil)
		exampleCreationInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.CreateHousehold(ctx, exampleCreationInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID
		exampleHousehold.Members = []*types.HouseholdUserMembership{}
		exampleInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildHouseholdCreationQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleHousehold.CreatedOn
		}

		actual, err := c.CreateHousehold(ctx, exampleInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error writing household creation audit log entry", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.ExternalID = ""
		exampleHousehold.BelongsToUser = exampleUser.ID
		exampleHousehold.Members = []*types.HouseholdUserMembership(nil)
		exampleCreationInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeCreationQuery, fakeCreationArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildHouseholdCreationQuery",
			testutils.ContextMatcher,
			exampleCreationInput,
		).Return(fakeCreationQuery, fakeCreationArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeCreationQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationArgs)...).
			WillReturnRows(newDatabaseResultForID(exampleHousehold.ID))

		fakeCreationAuditLogEntryQuery, fakeCreationAuditLogEntryArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.AuditLogEntrySQLQueryBuilder.On(
			"BuildCreateAuditLogEntryQuery",
			testutils.ContextMatcher,
			mock.IsType(&types.AuditLogEntryCreationInput{}),
		).Return(fakeCreationAuditLogEntryQuery, fakeCreationAuditLogEntryArgs).Once()

		db.ExpectExec(formatQueryForSQLMock(fakeCreationAuditLogEntryQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationAuditLogEntryArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleHousehold.CreatedOn
		}
		c.sqlQueryBuilder = mockQueryBuilder

		actual, err := c.CreateHousehold(ctx, exampleCreationInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error writing user addition to database", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.ExternalID = ""
		exampleHousehold.BelongsToUser = exampleUser.ID
		exampleHousehold.Members = []*types.HouseholdUserMembership(nil)
		exampleCreationInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)
		exampleHouseholdAdditionInput := &types.AddUserToHouseholdInput{
			Reason:         "household creation",
			UserID:         exampleUser.ID,
			HouseholdID:    exampleHousehold.ID,
			HouseholdRoles: []string{authorization.HouseholdAdminRole.String()},
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeCreationQuery, fakeCreationArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildHouseholdCreationQuery",
			testutils.ContextMatcher,
			exampleCreationInput,
		).Return(fakeCreationQuery, fakeCreationArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeCreationQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationArgs)...).
			WillReturnRows(newDatabaseResultForID(exampleHousehold.ID))

		fakeCreationAuditLogEntryQuery, fakeCreationAuditLogEntryArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.AuditLogEntrySQLQueryBuilder.On(
			"BuildCreateAuditLogEntryQuery",
			testutils.ContextMatcher,
			mock.IsType(&types.AuditLogEntryCreationInput{}),
		).Return(fakeCreationAuditLogEntryQuery, fakeCreationAuditLogEntryArgs).Once()

		db.ExpectExec(formatQueryForSQLMock(fakeCreationAuditLogEntryQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationAuditLogEntryArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(123))

		fakeHouseholdAdditionQuery, fakeHouseholdAdditionArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildAddUserToHouseholdQuery",
			testutils.ContextMatcher,
			exampleHouseholdAdditionInput,
		).Return(fakeHouseholdAdditionQuery, fakeHouseholdAdditionArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdAdditionQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdAdditionArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleHousehold.CreatedOn
		}
		c.sqlQueryBuilder = mockQueryBuilder

		actual, err := c.CreateHousehold(ctx, exampleCreationInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error writing user membership addition audit log entry", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.ExternalID = ""
		exampleHousehold.BelongsToUser = exampleUser.ID
		exampleHousehold.Members = []*types.HouseholdUserMembership(nil)
		exampleCreationInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)
		exampleHouseholdAdditionInput := &types.AddUserToHouseholdInput{
			Reason:         "household creation",
			UserID:         exampleUser.ID,
			HouseholdID:    exampleHousehold.ID,
			HouseholdRoles: []string{authorization.HouseholdAdminRole.String()},
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeCreationQuery, fakeCreationArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildHouseholdCreationQuery",
			testutils.ContextMatcher,
			exampleCreationInput,
		).Return(fakeCreationQuery, fakeCreationArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeCreationQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationArgs)...).
			WillReturnRows(newDatabaseResultForID(exampleHousehold.ID))

		fakeCreationAuditLogEntryQuery, fakeCreationAuditLogEntryArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.AuditLogEntrySQLQueryBuilder.On(
			"BuildCreateAuditLogEntryQuery",
			testutils.ContextMatcher,
			mock.IsType(&types.AuditLogEntryCreationInput{}),
		).Return(fakeCreationAuditLogEntryQuery, fakeCreationAuditLogEntryArgs).Once()

		db.ExpectExec(formatQueryForSQLMock(fakeCreationAuditLogEntryQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationAuditLogEntryArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(123))

		fakeHouseholdAdditionQuery, fakeHouseholdAdditionArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildAddUserToHouseholdQuery",
			testutils.ContextMatcher,
			exampleHouseholdAdditionInput,
		).Return(fakeHouseholdAdditionQuery, fakeHouseholdAdditionArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdAdditionQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdAdditionArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		fakeHouseholdAdditionAuditLogEntryQuery, fakeHouseholdAdditionAuditLogEntryArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.AuditLogEntrySQLQueryBuilder.On(
			"BuildCreateAuditLogEntryQuery",
			testutils.ContextMatcher,
			mock.IsType(&types.AuditLogEntryCreationInput{}),
		).Return(fakeHouseholdAdditionAuditLogEntryQuery, fakeHouseholdAdditionAuditLogEntryArgs).Once()

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdAdditionAuditLogEntryQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdAdditionAuditLogEntryArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleHousehold.CreatedOn
		}
		c.sqlQueryBuilder = mockQueryBuilder

		actual, err := c.CreateHousehold(ctx, exampleCreationInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.ExternalID = ""
		exampleHousehold.BelongsToUser = exampleUser.ID
		exampleHousehold.Members = []*types.HouseholdUserMembership(nil)
		exampleCreationInput := fakes.BuildFakeHouseholdCreationInputFromHousehold(exampleHousehold)
		exampleHouseholdAdditionInput := &types.AddUserToHouseholdInput{
			Reason:         "household creation",
			UserID:         exampleUser.ID,
			HouseholdID:    exampleHousehold.ID,
			HouseholdRoles: []string{authorization.HouseholdAdminRole.String()},
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeCreationQuery, fakeCreationArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildHouseholdCreationQuery",
			testutils.ContextMatcher,
			exampleCreationInput,
		).Return(fakeCreationQuery, fakeCreationArgs)

		db.ExpectQuery(formatQueryForSQLMock(fakeCreationQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationArgs)...).
			WillReturnRows(newDatabaseResultForID(exampleHousehold.ID))

		fakeCreationAuditLogEntryQuery, fakeCreationAuditLogEntryArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.AuditLogEntrySQLQueryBuilder.On(
			"BuildCreateAuditLogEntryQuery",
			testutils.ContextMatcher,
			mock.IsType(&types.AuditLogEntryCreationInput{}),
		).Return(fakeCreationAuditLogEntryQuery, fakeCreationAuditLogEntryArgs).Once()

		db.ExpectExec(formatQueryForSQLMock(fakeCreationAuditLogEntryQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationAuditLogEntryArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(123))

		fakeHouseholdAdditionQuery, fakeHouseholdAdditionArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdUserMembershipSQLQueryBuilder.On(
			"BuildAddUserToHouseholdQuery",
			testutils.ContextMatcher,
			exampleHouseholdAdditionInput,
		).Return(fakeHouseholdAdditionQuery, fakeHouseholdAdditionArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdAdditionQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdAdditionArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		fakeHouseholdAdditionAuditLogEntryQuery, fakeHouseholdAdditionAuditLogEntryArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.AuditLogEntrySQLQueryBuilder.On(
			"BuildCreateAuditLogEntryQuery",
			testutils.ContextMatcher,
			mock.IsType(&types.AuditLogEntryCreationInput{}),
		).Return(fakeHouseholdAdditionAuditLogEntryQuery, fakeHouseholdAdditionAuditLogEntryArgs).Once()

		db.ExpectExec(formatQueryForSQLMock(fakeHouseholdAdditionAuditLogEntryQuery)).
			WithArgs(interfaceToDriverValue(fakeHouseholdAdditionAuditLogEntryArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(123))

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() uint64 {
			return exampleHousehold.CreatedOn
		}
		c.sqlQueryBuilder = mockQueryBuilder

		actual, err := c.CreateHousehold(ctx, exampleCreationInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_UpdateHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildUpdateHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		assert.NoError(t, c.UpdateHousehold(ctx, exampleHousehold, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateHousehold(ctx, nil, exampleUser.ID, nil))
	})

	T.Run("with invalid actor ID", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateHousehold(ctx, exampleHousehold, 0, nil))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateHousehold(ctx, exampleHousehold, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildUpdateHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.UpdateHousehold(ctx, exampleHousehold, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error writing audit log entry", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildUpdateHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.UpdateHousehold(ctx, exampleHousehold, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildUpdateHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateHousehold(ctx, exampleHousehold, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_ArchiveHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildArchiveHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold.ID, exampleUser.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		assert.NoError(t, c.ArchiveHousehold(ctx, exampleHousehold.ID, exampleUser.ID, exampleUser.ID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHousehold(ctx, 0, exampleUser.ID, exampleUser.ID))
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveHousehold(ctx, exampleHousehold.ID, 0, exampleUser.ID))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveHousehold(ctx, exampleHousehold.ID, exampleUser.ID, exampleUser.ID))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildArchiveHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold.ID, exampleUser.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.ArchiveHousehold(ctx, exampleHousehold.ID, exampleUser.ID, exampleUser.ID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error writing audit log entry", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildArchiveHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold.ID, exampleUser.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.ArchiveHousehold(ctx, exampleHousehold.ID, exampleUser.ID, exampleUser.ID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildArchiveHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold.ID, exampleUser.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleHousehold.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveHousehold(ctx, exampleHousehold.ID, exampleUser.ID, exampleUser.ID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetAuditLogEntriesForHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHousehold := fakes.BuildFakeHousehold()
		exampleAuditLogEntryList := fakes.BuildFakeAuditLogEntryList()
		c, db := buildTestClient(t)

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetAuditLogEntriesForHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromAuditLogEntries(false, exampleAuditLogEntryList.Entries...))

		actual, err := c.GetAuditLogEntriesForHousehold(ctx, exampleHousehold.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntryList.Entries, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.GetAuditLogEntriesForHousehold(ctx, 0)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHousehold := fakes.BuildFakeHousehold()
		c, db := buildTestClient(t)

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetAuditLogEntriesForHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetAuditLogEntriesForHousehold(ctx, exampleHousehold.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHousehold := fakes.BuildFakeHousehold()
		c, db := buildTestClient(t)

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.HouseholdSQLQueryBuilder.On(
			"BuildGetAuditLogEntriesForHouseholdQuery",
			testutils.ContextMatcher,
			exampleHousehold.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetAuditLogEntriesForHousehold(ctx, exampleHousehold.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}
