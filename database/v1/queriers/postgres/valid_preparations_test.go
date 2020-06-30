package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowsFromValidPreparation(validPreparations ...*models.ValidPreparation) *sqlmock.Rows {
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
			x.UpdatedOn,
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
		x.UpdatedOn,
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

		expectedQuery := "SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.updated_on, valid_preparations.archived_on FROM valid_preparations WHERE valid_preparations.id = $1"
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

	expectedQuery := "SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.updated_on, valid_preparations.archived_on FROM valid_preparations WHERE valid_preparations.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.ID,
			).
			WillReturnRows(buildMockRowsFromValidPreparation(exampleValidPreparation))

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

func TestPostgres_buildGetValidPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.updated_on, valid_preparations.archived_on, (SELECT COUNT(valid_preparations.id) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL AND valid_preparations.created_on > $1 AND valid_preparations.created_on < $2 AND valid_preparations.updated_on > $3 AND valid_preparations.updated_on < $4 ORDER BY valid_preparations.id LIMIT 20 OFFSET 180"
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

	expectedListQuery := "SELECT valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon, valid_preparations.applicable_to_all_ingredients, valid_preparations.created_on, valid_preparations.updated_on, valid_preparations.archived_on, (SELECT COUNT(valid_preparations.id) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL) FROM valid_preparations WHERE valid_preparations.archived_on IS NULL ORDER BY valid_preparations.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidPreparationList := fakemodels.BuildFakeValidPreparationList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(
				buildMockRowsFromValidPreparation(
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WillReturnRows(buildErroneousMockRowFromValidPreparation(exampleValidPreparation))

		actual, err := p.GetValidPreparations(ctx, filter)
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

		expectedQuery := "INSERT INTO valid_preparations (name,description,icon) VALUES ($1,$2,$3) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleValidPreparation.Name,
			exampleValidPreparation.Description,
			exampleValidPreparation.Icon,
		}
		actualQuery, actualArgs := p.buildCreateValidPreparationQuery(exampleValidPreparation)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateValidPreparation(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO valid_preparations (name,description,icon) VALUES ($1,$2,$3) RETURNING id, created_on"

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

		expectedQuery := "UPDATE valid_preparations SET name = $1, description = $2, icon = $3, updated_on = extract(epoch FROM NOW()) WHERE id = $4 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleValidPreparation.Name,
			exampleValidPreparation.Description,
			exampleValidPreparation.Icon,
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

	expectedQuery := "UPDATE valid_preparations SET name = $1, description = $2, icon = $3, updated_on = extract(epoch FROM NOW()) WHERE id = $4 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidPreparation := fakemodels.BuildFakeValidPreparation()

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidPreparation.Name,
				exampleValidPreparation.Description,
				exampleValidPreparation.Icon,
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

		expectedQuery := "UPDATE valid_preparations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"
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

	expectedQuery := "UPDATE valid_preparations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1 RETURNING archived_on"

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
