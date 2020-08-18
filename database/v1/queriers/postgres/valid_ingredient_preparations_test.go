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

func buildMockRowsFromValidIngredientPreparations(validIngredientPreparations ...*models.ValidIngredientPreparation) *sqlmock.Rows {
	includeCount := len(validIngredientPreparations) > 1
	columns := validIngredientPreparationsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validIngredientPreparations {
		rowValues := []driver.Value{
			x.ID,
			x.Notes,
			x.ValidPreparationID,
			x.ValidIngredientID,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
		}

		if includeCount {
			rowValues = append(rowValues, len(validIngredientPreparations))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromValidIngredientPreparation(x *models.ValidIngredientPreparation) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(validIngredientPreparationsTableColumns).AddRow(
		x.ArchivedOn,
		x.Notes,
		x.ValidPreparationID,
		x.ValidIngredientID,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanValidIngredientPreparations(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanValidIngredientPreparations(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildValidIngredientPreparationExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		expectedQuery := "SELECT EXISTS ( SELECT valid_ingredient_preparations.id FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.id = $1 )"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.ID,
		}
		actualQuery, actualArgs := p.buildValidIngredientPreparationExistsQuery(exampleValidIngredientPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ValidIngredientPreparationExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT valid_ingredient_preparations.id FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.id = $1 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientPreparation.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.ValidIngredientPreparationExists(ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientPreparation.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.ValidIngredientPreparationExists(ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidIngredientPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.id = $1"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.ID,
		}
		actualQuery, actualArgs := p.buildGetValidIngredientPreparationQuery(exampleValidIngredientPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientPreparation.ID,
			).
			WillReturnRows(buildMockRowsFromValidIngredientPreparations(exampleValidIngredientPreparation))

		actual, err := p.GetValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparation, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientPreparation.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllValidIngredientPreparationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL"
		actualQuery := p.buildGetAllValidIngredientPreparationsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllValidIngredientPreparationsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllValidIngredientPreparationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfValidIngredientPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.id > $1 AND valid_ingredient_preparations.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfValidIngredientPreparationsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL"
	expectedGetQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.id > $1 AND valid_ingredient_preparations.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromValidIngredientPreparations(
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[0],
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[1],
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[2],
				),
			)

		out := make(chan []models.ValidIngredientPreparation)
		doneChan := make(chan bool, 1)

		err := p.GetAllValidIngredientPreparations(ctx, out)
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

		out := make(chan []models.ValidIngredientPreparation)

		err := p.GetAllValidIngredientPreparations(ctx, out)
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

		out := make(chan []models.ValidIngredientPreparation)

		err := p.GetAllValidIngredientPreparations(ctx, out)
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

		out := make(chan []models.ValidIngredientPreparation)

		err := p.GetAllValidIngredientPreparations(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromValidIngredientPreparation(exampleValidIngredientPreparation))

		out := make(chan []models.ValidIngredientPreparation)

		err := p.GetAllValidIngredientPreparations(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidIngredientPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on, (SELECT COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL AND valid_ingredient_preparations.created_on > $1 AND valid_ingredient_preparations.created_on < $2 AND valid_ingredient_preparations.last_updated_on > $3 AND valid_ingredient_preparations.last_updated_on < $4 ORDER BY valid_ingredient_preparations.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetValidIngredientPreparationsQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on, (SELECT COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL) FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL ORDER BY valid_ingredient_preparations.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromValidIngredientPreparations(
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[0],
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[1],
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[2],
				),
			)

		actual, err := p.GetValidIngredientPreparations(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparationList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidIngredientPreparations(ctx, filter)
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

		actual, err := p.GetValidIngredientPreparations(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning valid ingredient preparation", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromValidIngredientPreparation(exampleValidIngredientPreparation))

		actual, err := p.GetValidIngredientPreparations(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidIngredientPreparationsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM (SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}(nil)
		actualQuery, actualArgs := p.buildGetValidIngredientPreparationsWithIDsQuery(defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidIngredientPreparationsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList()
		var exampleIDs []uint64
		for _, validIngredientPreparation := range exampleValidIngredientPreparationList.ValidIngredientPreparations {
			exampleIDs = append(exampleIDs, validIngredientPreparation.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM (SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(
				buildMockRowsFromValidIngredientPreparations(
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[0],
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[1],
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[2],
				),
			)

		actual, err := p.GetValidIngredientPreparationsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparationList.ValidIngredientPreparations, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM (SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidIngredientPreparationsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM (SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := p.GetValidIngredientPreparationsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning valid ingredient preparation", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM (SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.valid_preparation_id, valid_ingredient_preparations.valid_ingredient_id, valid_ingredient_preparations.created_on, valid_ingredient_preparations.last_updated_on, valid_ingredient_preparations.archived_on FROM valid_ingredient_preparations JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(buildErroneousMockRowFromValidIngredientPreparation(exampleValidIngredientPreparation))

		actual, err := p.GetValidIngredientPreparationsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateValidIngredientPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		expectedQuery := "INSERT INTO valid_ingredient_preparations (notes,valid_preparation_id,valid_ingredient_id) VALUES ($1,$2,$3) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.Notes,
			exampleValidIngredientPreparation.ValidPreparationID,
			exampleValidIngredientPreparation.ValidIngredientID,
		}
		actualQuery, actualArgs := p.buildCreateValidIngredientPreparationQuery(exampleValidIngredientPreparation)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO valid_ingredient_preparations (notes,valid_preparation_id,valid_ingredient_id) VALUES ($1,$2,$3) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleValidIngredientPreparation.ID, exampleValidIngredientPreparation.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleValidIngredientPreparation.Notes,
				exampleValidIngredientPreparation.ValidPreparationID,
				exampleValidIngredientPreparation.ValidIngredientID,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateValidIngredientPreparation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparation, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleValidIngredientPreparation.Notes,
				exampleValidIngredientPreparation.ValidPreparationID,
				exampleValidIngredientPreparation.ValidIngredientID,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateValidIngredientPreparation(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateValidIngredientPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		expectedQuery := "UPDATE valid_ingredient_preparations SET notes = $1, valid_preparation_id = $2, valid_ingredient_id = $3, last_updated_on = extract(epoch FROM NOW()) WHERE id = $4 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.Notes,
			exampleValidIngredientPreparation.ValidPreparationID,
			exampleValidIngredientPreparation.ValidIngredientID,
			exampleValidIngredientPreparation.ID,
		}
		actualQuery, actualArgs := p.buildUpdateValidIngredientPreparationQuery(exampleValidIngredientPreparation)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE valid_ingredient_preparations SET notes = $1, valid_preparation_id = $2, valid_ingredient_id = $3, last_updated_on = extract(epoch FROM NOW()) WHERE id = $4 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientPreparation.Notes,
				exampleValidIngredientPreparation.ValidPreparationID,
				exampleValidIngredientPreparation.ValidIngredientID,
				exampleValidIngredientPreparation.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateValidIngredientPreparation(ctx, exampleValidIngredientPreparation)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientPreparation.Notes,
				exampleValidIngredientPreparation.ValidPreparationID,
				exampleValidIngredientPreparation.ValidIngredientID,
				exampleValidIngredientPreparation.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateValidIngredientPreparation(ctx, exampleValidIngredientPreparation)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveValidIngredientPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		expectedQuery := "UPDATE valid_ingredient_preparations SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.ID,
		}
		actualQuery, actualArgs := p.buildArchiveValidIngredientPreparationQuery(exampleValidIngredientPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE valid_ingredient_preparations SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientPreparation.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientPreparation.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientPreparation.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveValidIngredientPreparation(ctx, exampleValidIngredientPreparation.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
