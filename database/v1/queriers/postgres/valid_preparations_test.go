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

func buildMockRowsFromValidPreparations(validPreparations ...*models.ValidPreparation) *sqlmock.Rows {
	includeCount := len(validPreparations) > 1
	columns := validPreparationsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validPreparations {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Description,
			x.Icon,
			x.ApplicableToAllIngredients,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
		}

		if includeCount {
			rowValues = append(rowValues, len(validPreparations))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromValidPreparation(x *models.ValidPreparation) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(validPreparationsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.Description,
		x.Icon,
		x.ApplicableToAllIngredients,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanValidPreparations(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanValidPreparations(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanValidPreparations(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildValidPreparationExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		expectedQuery := "SELECT EXISTS ( SELECT valid_preparations.id FROM valid_preparations WHERE valid_preparations.id = $1 )"
		expectedArgs := []interface{}{
			exampleValidPreparation.ID,
		}
		actualQuery, actualArgs := p.buildValidPreparationExistsQuery(exampleValidPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ValidPreparationExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT valid_preparations.id FROM valid_preparations WHERE valid_preparations.id = $1 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.ValidPreparationExists(ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.ValidPreparationExists(ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		expectedQuery := "SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations WHERE valid_preparations.id = $1"
		expectedArgs := []interface{}{
			exampleValidPreparation.ID,
		}
		actualQuery, actualArgs := p.buildGetValidPreparationQuery(exampleValidPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidPreparation(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations WHERE valid_preparations.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
			).
			WillReturnRows(buildMockRowsFromValidPreparations(exampleValidPreparation))

		actual, err := p.GetValidPreparation(ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparation, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidPreparation(ctx, exampleValidPreparation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllValidPreparationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(valid_preparations.id) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL"
		actualQuery := p.buildGetAllValidPreparationsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllValidPreparationsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(valid_preparations.id) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllValidPreparationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfValidPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations WHERE valid_preparations.id > $1 AND valid_preparations.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfValidPreparationsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllValidPreparations(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(valid_preparations.id) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL"
	expectedGetQuery := "SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations WHERE valid_preparations.id > $1 AND valid_preparations.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromValidPreparations(
					&exampleValidPreparationList.ValidPreparations[0],
					&exampleValidPreparationList.ValidPreparations[1],
					&exampleValidPreparationList.ValidPreparations[2],
				),
			)

		out := make(chan []models.ValidPreparation)
		doneChan := make(chan bool, 1)

		err := p.GetAllValidPreparations(ctx, out)
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

		out := make(chan []models.ValidPreparation)

		err := p.GetAllValidPreparations(ctx, out)
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

		out := make(chan []models.ValidPreparation)

		err := p.GetAllValidPreparations(ctx, out)
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

		out := make(chan []models.ValidPreparation)

		err := p.GetAllValidPreparations(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromValidPreparation(exampleValidPreparation))

		out := make(chan []models.ValidPreparation)

		err := p.GetAllValidPreparations(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on, (SELECT COUNT(valid_preparations.id) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.created_on > $1 AND valid_preparations.created_on < $2 AND valid_preparations.last_updated_on > $3 AND valid_preparations.last_updated_on < $4 ORDER BY valid_preparations.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetValidPreparationsQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidPreparations(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on, (SELECT COUNT(valid_preparations.id) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL ORDER BY valid_preparations.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromValidPreparations(
					&exampleValidPreparationList.ValidPreparations[0],
					&exampleValidPreparationList.ValidPreparations[1],
					&exampleValidPreparationList.ValidPreparations[2],
				),
			)

		actual, err := p.GetValidPreparations(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidPreparations(ctx, filter)
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

		actual, err := p.GetValidPreparations(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning valid preparation", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(buildErroneousMockRowFromValidPreparation(exampleValidPreparation))

		actual, err := p.GetValidPreparations(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidPreparationsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM (SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_preparations WHERE valid_preparations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}(nil)
		actualQuery, actualArgs := p.buildGetValidPreparationsWithIDsQuery(defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidPreparationsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList()
		var exampleIDs []uint64
		for _, validPreparation := range exampleValidPreparationList.ValidPreparations {
			exampleIDs = append(exampleIDs, validPreparation.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM (SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_preparations WHERE valid_preparations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(
				buildMockRowsFromValidPreparations(
					&exampleValidPreparationList.ValidPreparations[0],
					&exampleValidPreparationList.ValidPreparations[1],
					&exampleValidPreparationList.ValidPreparations[2],
				),
			)

		actual, err := p.GetValidPreparationsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationList.ValidPreparations, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM (SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_preparations WHERE valid_preparations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidPreparationsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM (SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_preparations WHERE valid_preparations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := p.GetValidPreparationsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning valid preparation", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM (SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.last_updated_on, valid_preparations.archived_on FROM valid_preparations JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS valid_preparations WHERE valid_preparations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(buildErroneousMockRowFromValidPreparation(exampleValidPreparation))

		actual, err := p.GetValidPreparationsWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateValidPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		expectedQuery := "INSERT INTO valid_preparations (name,description,icon,applicable_to_all_ingredients) VALUES ($1,$2,$3,$4) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleValidPreparation.Name,
			exampleValidPreparation.Description,
			exampleValidPreparation.Icon,
			exampleValidPreparation.ApplicableToAllIngredients,
		}
		actualQuery, actualArgs := p.buildCreateValidPreparationQuery(exampleValidPreparation)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateValidPreparation(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO valid_preparations (name,description,icon,applicable_to_all_ingredients) VALUES ($1,$2,$3,$4) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleValidPreparation.ID, exampleValidPreparation.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleValidPreparation.Name,
				exampleValidPreparation.Description,
				exampleValidPreparation.Icon,
				exampleValidPreparation.ApplicableToAllIngredients,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateValidPreparation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparation, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
		exampleInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleValidPreparation.Name,
				exampleValidPreparation.Description,
				exampleValidPreparation.Icon,
				exampleValidPreparation.ApplicableToAllIngredients,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateValidPreparation(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateValidPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		expectedQuery := "UPDATE valid_preparations SET name = $1, description = $2, icon = $3, applicable_to_all_ingredients = $4, last_updated_on = extract(epoch FROM NOW()) WHERE id = $5 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleValidPreparation.Name,
			exampleValidPreparation.Description,
			exampleValidPreparation.Icon,
			exampleValidPreparation.ApplicableToAllIngredients,
			exampleValidPreparation.ID,
		}
		actualQuery, actualArgs := p.buildUpdateValidPreparationQuery(exampleValidPreparation)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateValidPreparation(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE valid_preparations SET name = $1, description = $2, icon = $3, applicable_to_all_ingredients = $4, last_updated_on = extract(epoch FROM NOW()) WHERE id = $5 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.Name,
				exampleValidPreparation.Description,
				exampleValidPreparation.Icon,
				exampleValidPreparation.ApplicableToAllIngredients,
				exampleValidPreparation.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateValidPreparation(ctx, exampleValidPreparation)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.Name,
				exampleValidPreparation.Description,
				exampleValidPreparation.Icon,
				exampleValidPreparation.ApplicableToAllIngredients,
				exampleValidPreparation.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateValidPreparation(ctx, exampleValidPreparation)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveValidPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		expectedQuery := "UPDATE valid_preparations SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleValidPreparation.ID,
		}
		actualQuery, actualArgs := p.buildArchiveValidPreparationQuery(exampleValidPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveValidPreparation(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE valid_preparations SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveValidPreparation(ctx, exampleValidPreparation.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveValidPreparation(ctx, exampleValidPreparation.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveValidPreparation(ctx, exampleValidPreparation.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
