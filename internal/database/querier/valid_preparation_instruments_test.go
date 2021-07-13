package querier

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/internal/database"
	querybuilding "gitlab.com/prixfixe/prixfixe/internal/database/querybuilding"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromValidPreparationInstruments(includeCounts bool, filteredCount uint64, validPreparationInstruments ...*types.ValidPreparationInstrument) *sqlmock.Rows {
	columns := querybuilding.ValidPreparationInstrumentsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validPreparationInstruments {
		rowValues := []driver.Value{
			x.ID,
			x.ExternalID,
			x.InstrumentID,
			x.PreparationID,
			x.Notes,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validPreparationInstruments))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidPreparationInstruments(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidPreparationInstruments(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidPreparationInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildValidPreparationInstrumentExistsQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidPreparationInstrumentExists(ctx, exampleValidPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidPreparationInstrumentExists(ctx, 0)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildValidPreparationInstrumentExistsQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidPreparationInstrumentExists(ctx, exampleValidPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildValidPreparationInstrumentExistsQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidPreparationInstrumentExists(ctx, exampleValidPreparationInstrument.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromValidPreparationInstruments(false, 0, exampleValidPreparationInstrument))

		actual, err := c.GetValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationInstrument, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidPreparationInstrument(ctx, 0)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetAllValidPreparationInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		exampleCount := uint64(123)

		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetAllValidPreparationInstrumentsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetAllValidPreparationInstrumentsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetAllValidPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		results := make(chan []*types.ValidPreparationInstrument)
		doneChan := make(chan bool, 1)
		expectedCount := uint64(20)
		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()
		exampleBatchSize := uint16(1000)

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetAllValidPreparationInstrumentsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetBatchOfValidPreparationInstrumentsQuery",
			testutils.ContextMatcher,
			uint64(1),
			uint64(exampleBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnRows(buildMockRowsFromValidPreparationInstruments(false, 0, exampleValidPreparationInstrumentList.ValidPreparationInstruments...))

		assert.NoError(t, c.GetAllValidPreparationInstruments(ctx, results, exampleBatchSize))

		stillQuerying := true
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

		exampleBatchSize := uint16(1000)
		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.GetAllValidPreparationInstruments(ctx, nil, exampleBatchSize))
	})

	T.Run("with now rows returned", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		results := make(chan []*types.ValidPreparationInstrument)
		expectedCount := uint64(20)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetAllValidPreparationInstrumentsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetBatchOfValidPreparationInstrumentsQuery",
			testutils.ContextMatcher,
			uint64(1),
			uint64(exampleBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnError(sql.ErrNoRows)

		assert.NoError(t, c.GetAllValidPreparationInstruments(ctx, results, exampleBatchSize))

		time.Sleep(time.Second)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		results := make(chan []*types.ValidPreparationInstrument)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetAllValidPreparationInstrumentsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		c.sqlQueryBuilder = mockQueryBuilder

		err := c.GetAllValidPreparationInstruments(ctx, results, exampleBatchSize)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		results := make(chan []*types.ValidPreparationInstrument)
		expectedCount := uint64(20)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetAllValidPreparationInstrumentsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetBatchOfValidPreparationInstrumentsQuery",
			testutils.ContextMatcher,
			uint64(1),
			uint64(exampleBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnError(errors.New("blah"))

		assert.NoError(t, c.GetAllValidPreparationInstruments(ctx, results, exampleBatchSize))

		time.Sleep(time.Second)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()

		results := make(chan []*types.ValidPreparationInstrument)
		expectedCount := uint64(20)
		exampleBatchSize := uint16(1000)

		c, db := buildTestClient(t)

		fakeQuery, _ := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetAllValidPreparationInstrumentsCountQuery",
			testutils.ContextMatcher,
		).Return(fakeQuery, []interface{}{})

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(expectedCount))

		secondFakeQuery, secondFakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetBatchOfValidPreparationInstrumentsQuery",
			testutils.ContextMatcher,
			uint64(1),
			uint64(exampleBatchSize+1),
		).Return(secondFakeQuery, secondFakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(secondFakeQuery)).
			WithArgs(interfaceToDriverValue(secondFakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		assert.NoError(t, c.GetAllValidPreparationInstruments(ctx, results, exampleBatchSize))

		time.Sleep(time.Second)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetValidPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetValidPreparationInstrumentsQuery",
			testutils.ContextMatcher,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromValidPreparationInstruments(true, exampleValidPreparationInstrumentList.FilteredCount, exampleValidPreparationInstrumentList.ValidPreparationInstruments...))

		actual, err := c.GetValidPreparationInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()
		exampleValidPreparationInstrumentList.Page = 0
		exampleValidPreparationInstrumentList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetValidPreparationInstrumentsQuery",
			testutils.ContextMatcher,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromValidPreparationInstruments(true, exampleValidPreparationInstrumentList.FilteredCount, exampleValidPreparationInstrumentList.ValidPreparationInstruments...))

		actual, err := c.GetValidPreparationInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetValidPreparationInstrumentsQuery",
			testutils.ContextMatcher,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidPreparationInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetValidPreparationInstrumentsQuery",
			testutils.ContextMatcher,
			false,
			filter,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidPreparationInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetValidPreparationInstrumentsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()

		var exampleIDs []uint64
		for _, x := range exampleValidPreparationInstrumentList.ValidPreparationInstruments {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetValidPreparationInstrumentsWithIDsQuery",
			testutils.ContextMatcher,
			defaultLimit,
			exampleIDs,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromValidPreparationInstruments(false, 0, exampleValidPreparationInstrumentList.ValidPreparationInstruments...))

		actual, err := c.GetValidPreparationInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationInstrumentList.ValidPreparationInstruments, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("sets limit if not present", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()
		var exampleIDs []uint64
		for _, x := range exampleValidPreparationInstrumentList.ValidPreparationInstruments {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetValidPreparationInstrumentsWithIDsQuery",
			testutils.ContextMatcher,
			defaultLimit,
			exampleIDs,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromValidPreparationInstruments(false, 0, exampleValidPreparationInstrumentList.ValidPreparationInstruments...))

		actual, err := c.GetValidPreparationInstrumentsWithIDs(ctx, 0, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationInstrumentList.ValidPreparationInstruments, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()
		var exampleIDs []uint64
		for _, x := range exampleValidPreparationInstrumentList.ValidPreparationInstruments {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetValidPreparationInstrumentsWithIDsQuery",
			testutils.ContextMatcher,
			defaultLimit,
			exampleIDs,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidPreparationInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrumentList := fakes.BuildFakeValidPreparationInstrumentList()
		var exampleIDs []uint64
		for _, x := range exampleValidPreparationInstrumentList.ValidPreparationInstruments {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetValidPreparationInstrumentsWithIDsQuery",
			testutils.ContextMatcher,
			defaultLimit,
			exampleIDs,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidPreparationInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_CreateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleValidPreparationInstrument.ExternalID = ""
		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeCreationQuery, fakeCreationArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildCreateValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeCreationQuery, fakeCreationArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeCreationQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleValidPreparationInstrument.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		c.timeFunc = func() uint64 {
			return exampleValidPreparationInstrument.CreatedOn
		}
		c.sqlQueryBuilder = mockQueryBuilder

		actual, err := c.CreateValidPreparationInstrument(ctx, exampleInput, exampleUser.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationInstrument, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleValidPreparationInstrument.ExternalID = ""

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidPreparationInstrument(ctx, nil, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid actor ID", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleValidPreparationInstrument.ExternalID = ""
		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidPreparationInstrument(ctx, exampleInput, 0)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleValidPreparationInstrument.ExternalID = ""
		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.CreateValidPreparationInstrument(ctx, exampleInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleUser := fakes.BuildFakeUser()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildCreateValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() uint64 {
			return exampleValidPreparationInstrument.CreatedOn
		}

		actual, err := c.CreateValidPreparationInstrument(ctx, exampleInput, exampleUser.ID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error creating audit log entry", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleValidPreparationInstrument.ExternalID = ""
		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeCreationQuery, fakeCreationArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildCreateValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeCreationQuery, fakeCreationArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeCreationQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleValidPreparationInstrument.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.CreateValidPreparationInstrument(ctx, exampleInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleValidPreparationInstrument.ExternalID = ""
		exampleInput := fakes.BuildFakeValidPreparationInstrumentCreationInputFromValidPreparationInstrument(exampleValidPreparationInstrument)

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeCreationQuery, fakeCreationArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildCreateValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleInput,
		).Return(fakeCreationQuery, fakeCreationArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectExec(formatQueryForSQLMock(fakeCreationQuery)).
			WithArgs(interfaceToDriverValue(fakeCreationArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleValidPreparationInstrument.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() uint64 {
			return exampleValidPreparationInstrument.CreatedOn
		}

		actual, err := c.CreateValidPreparationInstrument(ctx, exampleInput, exampleUser.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_UpdateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildUpdateValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleValidPreparationInstrument.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectCommit()

		assert.NoError(t, c.UpdateValidPreparationInstrument(ctx, exampleValidPreparationInstrument, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidPreparationInstrument(ctx, nil, exampleUser.ID, nil))
	})

	T.Run("with invalid actor ID", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidPreparationInstrument(ctx, exampleValidPreparationInstrument, 0, nil))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidPreparationInstrument(ctx, exampleValidPreparationInstrument, exampleUser.ID, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildUpdateValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.UpdateValidPreparationInstrument(ctx, exampleValidPreparationInstrument, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error writing audit log entry to database", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildUpdateValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleValidPreparationInstrument.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.UpdateValidPreparationInstrument(ctx, exampleValidPreparationInstrument, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)
		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeUpdateQuery, fakeUpdateArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildUpdateValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument,
		).Return(fakeUpdateQuery, fakeUpdateArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeUpdateQuery)).
			WithArgs(interfaceToDriverValue(fakeUpdateArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleValidPreparationInstrument.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.UpdateValidPreparationInstrument(ctx, exampleValidPreparationInstrument, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_ArchiveValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildArchiveValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleValidPreparationInstrument.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.NoError(t, c.ArchiveValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidPreparationInstrument(ctx, 0, exampleUserID))
	})

	T.Run("with invalid actor ID", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID, 0))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID, exampleUserID))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildArchiveValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.ArchiveValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error writing audit log entry", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildArchiveValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleValidPreparationInstrument.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, errors.New("blah"))

		db.ExpectRollback()

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.ArchiveValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()
		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		db.ExpectBegin()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildArchiveValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)

		db.ExpectExec(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnResult(newSuccessfulDatabaseResult(exampleValidPreparationInstrument.ID))

		expectAuditLogEntryInTransaction(mockQueryBuilder, db, nil)

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.sqlQueryBuilder = mockQueryBuilder

		assert.Error(t, c.ArchiveValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID, exampleUserID))

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}

func TestQuerier_GetAuditLogEntriesForValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()
		exampleAuditLogEntriesList := fakes.BuildFakeAuditLogEntryList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetAuditLogEntriesForValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildMockRowsFromAuditLogEntries(false, exampleAuditLogEntriesList.Entries...))

		actual, err := c.GetAuditLogEntriesForValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleAuditLogEntriesList.Entries, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetAuditLogEntriesForValidPreparationInstrument(ctx, 0)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetAuditLogEntriesForValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetAuditLogEntriesForValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationInstrument := fakes.BuildFakeValidPreparationInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockQueryBuilder := database.BuildMockSQLQueryBuilder()

		fakeQuery, fakeArgs := fakes.BuildFakeSQLQuery()
		mockQueryBuilder.ValidPreparationInstrumentSQLQueryBuilder.On(
			"BuildGetAuditLogEntriesForValidPreparationInstrumentQuery",
			testutils.ContextMatcher,
			exampleValidPreparationInstrument.ID,
		).Return(fakeQuery, fakeArgs)
		c.sqlQueryBuilder = mockQueryBuilder

		db.ExpectQuery(formatQueryForSQLMock(fakeQuery)).
			WithArgs(interfaceToDriverValue(fakeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetAuditLogEntriesForValidPreparationInstrument(ctx, exampleValidPreparationInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockQueryBuilder)
	})
}
