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

func buildMockRowsFromValidInstruments(validInstruments ...*models.ValidInstrument) *sqlmock.Rows {
	includeCount := len(validInstruments) > 1
	columns := validInstrumentsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validInstruments {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Variant,
			x.Description,
			x.Icon,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
		}

		if includeCount {
			rowValues = append(rowValues, len(validInstruments))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromValidInstrument(x *models.ValidInstrument) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(validInstrumentsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.Variant,
		x.Description,
		x.Icon,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanValidInstruments(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanValidInstruments(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildValidInstrumentExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		expectedQuery := "SELECT EXISTS ( SELECT valid_instruments.id FROM valid_instruments WHERE valid_instruments.id = $1 )"
		expectedArgs := []interface{}{
			exampleValidInstrument.ID,
		}
		actualQuery, actualArgs := p.buildValidInstrumentExistsQuery(exampleValidInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ValidInstrumentExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT valid_instruments.id FROM valid_instruments WHERE valid_instruments.id = $1 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidInstrument.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.ValidInstrumentExists(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidInstrument.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.ValidInstrumentExists(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		expectedQuery := "SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments WHERE valid_instruments.id = $1"
		expectedArgs := []interface{}{
			exampleValidInstrument.ID,
		}
		actualQuery, actualArgs := p.buildGetValidInstrumentQuery(exampleValidInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidInstrument(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments WHERE valid_instruments.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidInstrument.ID,
			).
			WillReturnRows(buildMockRowsFromValidInstruments(exampleValidInstrument))

		actual, err := p.GetValidInstrument(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrument, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidInstrument.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidInstrument(ctx, exampleValidInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllValidInstrumentsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(valid_instruments.id) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL"
		actualQuery := p.buildGetAllValidInstrumentsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllValidInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(valid_instruments.id) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllValidInstrumentsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfValidInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments WHERE valid_instruments.id > $1 AND valid_instruments.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfValidInstrumentsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllValidInstruments(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(valid_instruments.id) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL"
	expectedGetQuery := "SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments WHERE valid_instruments.id > $1 AND valid_instruments.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleValidInstrumentList := fakemodels.BuildFakeValidInstrumentList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromValidInstruments(
					&exampleValidInstrumentList.ValidInstruments[0],
					&exampleValidInstrumentList.ValidInstruments[1],
					&exampleValidInstrumentList.ValidInstruments[2],
				),
			)

		out := make(chan []models.ValidInstrument)
		doneChan := make(chan bool, 1)

		err := p.GetAllValidInstruments(ctx, out)
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

		out := make(chan []models.ValidInstrument)

		err := p.GetAllValidInstruments(ctx, out)
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

		out := make(chan []models.ValidInstrument)

		err := p.GetAllValidInstruments(ctx, out)
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

		out := make(chan []models.ValidInstrument)

		err := p.GetAllValidInstruments(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromValidInstrument(exampleValidInstrument))

		out := make(chan []models.ValidInstrument)

		err := p.GetAllValidInstruments(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on, (SELECT COUNT(valid_instruments.id) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL AND valid_instruments.created_on > $1 AND valid_instruments.created_on < $2 AND valid_instruments.last_updated_on > $3 AND valid_instruments.last_updated_on < $4 ORDER BY valid_instruments.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetValidInstrumentsQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidInstruments(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on, (SELECT COUNT(valid_instruments.id) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL) FROM valid_instruments WHERE valid_instruments.archived_on IS NULL ORDER BY valid_instruments.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidInstrumentList := fakemodels.BuildFakeValidInstrumentList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromValidInstruments(
					&exampleValidInstrumentList.ValidInstruments[0],
					&exampleValidInstrumentList.ValidInstruments[1],
					&exampleValidInstrumentList.ValidInstruments[2],
				),
			)

		actual, err := p.GetValidInstruments(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrumentList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidInstruments(ctx, filter)
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

		actual, err := p.GetValidInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning valid instrument", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromValidInstrument(exampleValidInstrument))

		actual, err := p.GetValidInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidInstrumentsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM (SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_instruments WHERE valid_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}(nil)
		actualQuery, actualArgs := p.buildGetValidInstrumentsWithIDsQuery(defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidInstrumentsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidInstrumentList := fakemodels.BuildFakeValidInstrumentList()
		var exampleIDs []uint64
		for _, validInstrument := range exampleValidInstrumentList.ValidInstruments {
			exampleIDs = append(exampleIDs, validInstrument.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM (SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_instruments WHERE valid_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(
				buildMockRowsFromValidInstruments(
					&exampleValidInstrumentList.ValidInstruments[0],
					&exampleValidInstrumentList.ValidInstruments[1],
					&exampleValidInstrumentList.ValidInstruments[2],
				),
			)

		actual, err := p.GetValidInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrumentList.ValidInstruments, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM (SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_instruments WHERE valid_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM (SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_instruments WHERE valid_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := p.GetValidInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning valid instrument", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM (SELECT valid_instruments.id, valid_instruments.name, valid_instruments.variant, valid_instruments.description, valid_instruments.icon, valid_instruments.created_on, valid_instruments.last_updated_on, valid_instruments.archived_on FROM valid_instruments JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_instruments WHERE valid_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(buildErroneousMockRowFromValidInstrument(exampleValidInstrument))

		actual, err := p.GetValidInstrumentsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateValidInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		expectedQuery := "INSERT INTO valid_instruments (name,variant,description,icon) VALUES ($1,$2,$3,$4) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleValidInstrument.Name,
			exampleValidInstrument.Variant,
			exampleValidInstrument.Description,
			exampleValidInstrument.Icon,
		}
		actualQuery, actualArgs := p.buildCreateValidInstrumentQuery(exampleValidInstrument)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateValidInstrument(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO valid_instruments (name,variant,description,icon) VALUES ($1,$2,$3,$4) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
		exampleInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleValidInstrument.ID, exampleValidInstrument.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleValidInstrument.Name,
				exampleValidInstrument.Variant,
				exampleValidInstrument.Description,
				exampleValidInstrument.Icon,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateValidInstrument(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrument, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
		exampleInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleValidInstrument.Name,
				exampleValidInstrument.Variant,
				exampleValidInstrument.Description,
				exampleValidInstrument.Icon,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateValidInstrument(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateValidInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		expectedQuery := "UPDATE valid_instruments SET name = $1, variant = $2, description = $3, icon = $4, last_updated_on = extract(epoch FROM NOW()) WHERE id = $5 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleValidInstrument.Name,
			exampleValidInstrument.Variant,
			exampleValidInstrument.Description,
			exampleValidInstrument.Icon,
			exampleValidInstrument.ID,
		}
		actualQuery, actualArgs := p.buildUpdateValidInstrumentQuery(exampleValidInstrument)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE valid_instruments SET name = $1, variant = $2, description = $3, icon = $4, last_updated_on = extract(epoch FROM NOW()) WHERE id = $5 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidInstrument.Name,
				exampleValidInstrument.Variant,
				exampleValidInstrument.Description,
				exampleValidInstrument.Icon,
				exampleValidInstrument.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateValidInstrument(ctx, exampleValidInstrument)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidInstrument.Name,
				exampleValidInstrument.Variant,
				exampleValidInstrument.Description,
				exampleValidInstrument.Icon,
				exampleValidInstrument.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateValidInstrument(ctx, exampleValidInstrument)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveValidInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		expectedQuery := "UPDATE valid_instruments SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleValidInstrument.ID,
		}
		actualQuery, actualArgs := p.buildArchiveValidInstrumentQuery(exampleValidInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveValidInstrument(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE valid_instruments SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidInstrument.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveValidInstrument(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidInstrument.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveValidInstrument(ctx, exampleValidInstrument.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidInstrument.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveValidInstrument(ctx, exampleValidInstrument.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
