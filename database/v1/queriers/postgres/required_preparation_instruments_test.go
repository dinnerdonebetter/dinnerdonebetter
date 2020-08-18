package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowsFromRequiredPreparationInstruments(requiredPreparationInstruments ...*models.RequiredPreparationInstrument) *sqlmock.Rows {
	includeCount := len(requiredPreparationInstruments) > 1
	columns := requiredPreparationInstrumentsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range requiredPreparationInstruments {
		rowValues := []driver.Value{
			x.ID,
			x.InstrumentID,
			x.PreparationID,
			x.Notes,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
		}

		if includeCount {
			rowValues = append(rowValues, len(requiredPreparationInstruments))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRequiredPreparationInstrument(x *models.RequiredPreparationInstrument) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(requiredPreparationInstrumentsTableColumns).AddRow(
		x.ArchivedOn,
		x.InstrumentID,
		x.PreparationID,
		x.Notes,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRequiredPreparationInstruments(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanRequiredPreparationInstruments(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanRequiredPreparationInstruments(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRequiredPreparationInstrumentExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		expectedQuery := "SELECT EXISTS ( SELECT required_preparation_instruments.id FROM required_preparation_instruments WHERE required_preparation_instruments.id = $1 )"
		expectedArgs := []interface{}{
			exampleRequiredPreparationInstrument.ID,
		}
		actualQuery, actualArgs := p.buildRequiredPreparationInstrumentExistsQuery(exampleRequiredPreparationInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RequiredPreparationInstrumentExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT required_preparation_instruments.id FROM required_preparation_instruments WHERE required_preparation_instruments.id = $1 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RequiredPreparationInstrumentExists(ctx, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RequiredPreparationInstrumentExists(ctx, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		expectedQuery := "SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM required_preparation_instruments WHERE required_preparation_instruments.id = $1"
		expectedArgs := []interface{}{
			exampleRequiredPreparationInstrument.ID,
		}
		actualQuery, actualArgs := p.buildGetRequiredPreparationInstrumentQuery(exampleRequiredPreparationInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM required_preparation_instruments WHERE required_preparation_instruments.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.ID,
			).
			WillReturnRows(buildMockRowsFromRequiredPreparationInstruments(exampleRequiredPreparationInstrument))

		actual, err := p.GetRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrument, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRequiredPreparationInstrumentsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(required_preparation_instruments.id) FROM required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL"
		actualQuery := p.buildGetAllRequiredPreparationInstrumentsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRequiredPreparationInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(required_preparation_instruments.id) FROM required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRequiredPreparationInstrumentsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfRequiredPreparationInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM required_preparation_instruments WHERE required_preparation_instruments.id > $1 AND required_preparation_instruments.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfRequiredPreparationInstrumentsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllRequiredPreparationInstruments(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(required_preparation_instruments.id) FROM required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL"
	expectedGetQuery := "SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM required_preparation_instruments WHERE required_preparation_instruments.id > $1 AND required_preparation_instruments.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRequiredPreparationInstrumentList := fakemodels.BuildFakeRequiredPreparationInstrumentList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromRequiredPreparationInstruments(
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[0],
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[1],
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[2],
				),
			)

		out := make(chan []models.RequiredPreparationInstrument)
		doneChan := make(chan bool, 1)

		err := p.GetAllRequiredPreparationInstruments(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []models.RequiredPreparationInstrument)

		err := p.GetAllRequiredPreparationInstruments(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []models.RequiredPreparationInstrument)

		err := p.GetAllRequiredPreparationInstruments(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []models.RequiredPreparationInstrument)

		err := p.GetAllRequiredPreparationInstruments(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument))

		out := make(chan []models.RequiredPreparationInstrument)

		err := p.GetAllRequiredPreparationInstruments(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRequiredPreparationInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on, (SELECT COUNT(required_preparation_instruments.id) FROM required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL) FROM required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL AND required_preparation_instruments.created_on > $1 AND required_preparation_instruments.created_on < $2 AND required_preparation_instruments.last_updated_on > $3 AND required_preparation_instruments.last_updated_on < $4 ORDER BY required_preparation_instruments.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRequiredPreparationInstrumentsQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRequiredPreparationInstruments(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on, (SELECT COUNT(required_preparation_instruments.id) FROM required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL) FROM required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL ORDER BY required_preparation_instruments.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRequiredPreparationInstrumentList := fakemodels.BuildFakeRequiredPreparationInstrumentList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromRequiredPreparationInstruments(
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[0],
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[1],
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[2],
				),
			)

		actual, err := p.GetRequiredPreparationInstruments(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrumentList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRequiredPreparationInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRequiredPreparationInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning required preparation instrument", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument))

		actual, err := p.GetRequiredPreparationInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRequiredPreparationInstrumentsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM (SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM required_preparation_instruments JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}(nil)
		actualQuery, actualArgs := p.buildGetRequiredPreparationInstrumentsWithIDsQuery(defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRequiredPreparationInstrumentsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleRequiredPreparationInstrumentList := fakemodels.BuildFakeRequiredPreparationInstrumentList()
		var exampleIDs []uint64
		for _, requiredPreparationInstrument := range exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments {
			exampleIDs = append(exampleIDs, requiredPreparationInstrument.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM (SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM required_preparation_instruments JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(
				buildMockRowsFromRequiredPreparationInstruments(
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[0],
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[1],
					&exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments[2],
				),
			)

		actual, err := p.GetRequiredPreparationInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrumentList.RequiredPreparationInstruments, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM (SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM required_preparation_instruments JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRequiredPreparationInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM (SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM required_preparation_instruments JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRequiredPreparationInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning required preparation instrument", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM (SELECT required_preparation_instruments.id, required_preparation_instruments.instrument_id, required_preparation_instruments.preparation_id, required_preparation_instruments.notes, required_preparation_instruments.created_on, required_preparation_instruments.last_updated_on, required_preparation_instruments.archived_on FROM required_preparation_instruments JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS required_preparation_instruments WHERE required_preparation_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(buildErroneousMockRowFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument))

		actual, err := p.GetRequiredPreparationInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		expectedQuery := "INSERT INTO required_preparation_instruments (instrument_id,preparation_id,notes) VALUES ($1,$2,$3) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRequiredPreparationInstrument.InstrumentID,
			exampleRequiredPreparationInstrument.PreparationID,
			exampleRequiredPreparationInstrument.Notes,
		}
		actualQuery, actualArgs := p.buildCreateRequiredPreparationInstrumentQuery(exampleRequiredPreparationInstrument)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO required_preparation_instruments (instrument_id,preparation_id,notes) VALUES ($1,$2,$3) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRequiredPreparationInstrument.ID, exampleRequiredPreparationInstrument.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.InstrumentID,
				exampleRequiredPreparationInstrument.PreparationID,
				exampleRequiredPreparationInstrument.Notes,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRequiredPreparationInstrument(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRequiredPreparationInstrument, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
		exampleInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.InstrumentID,
				exampleRequiredPreparationInstrument.PreparationID,
				exampleRequiredPreparationInstrument.Notes,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRequiredPreparationInstrument(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		expectedQuery := "UPDATE required_preparation_instruments SET instrument_id = $1, preparation_id = $2, notes = $3, last_updated_on = extract(epoch FROM NOW()) WHERE id = $4 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleRequiredPreparationInstrument.InstrumentID,
			exampleRequiredPreparationInstrument.PreparationID,
			exampleRequiredPreparationInstrument.Notes,
			exampleRequiredPreparationInstrument.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRequiredPreparationInstrumentQuery(exampleRequiredPreparationInstrument)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE required_preparation_instruments SET instrument_id = $1, preparation_id = $2, notes = $3, last_updated_on = extract(epoch FROM NOW()) WHERE id = $4 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.InstrumentID,
				exampleRequiredPreparationInstrument.PreparationID,
				exampleRequiredPreparationInstrument.Notes,
				exampleRequiredPreparationInstrument.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.InstrumentID,
				exampleRequiredPreparationInstrument.PreparationID,
				exampleRequiredPreparationInstrument.Notes,
				exampleRequiredPreparationInstrument.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRequiredPreparationInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		expectedQuery := "UPDATE required_preparation_instruments SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRequiredPreparationInstrument.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRequiredPreparationInstrumentQuery(exampleRequiredPreparationInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRequiredPreparationInstrument(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE required_preparation_instruments SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRequiredPreparationInstrument.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
