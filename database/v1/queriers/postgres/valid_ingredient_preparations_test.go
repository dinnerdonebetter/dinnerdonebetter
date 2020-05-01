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

func buildMockRowsFromValidIngredientPreparation(validIngredientPreparations ...*models.ValidIngredientPreparation) *sqlmock.Rows {
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
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
			x.BelongsToValidIngredient,
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
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsToValidIngredient,
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

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		expectedQuery := "SELECT EXISTS ( SELECT valid_ingredient_preparations.id FROM valid_ingredient_preparations JOIN valid_ingredients ON valid_ingredient_preparations.belongs_to_valid_ingredient=valid_ingredients.id WHERE valid_ingredient_preparations.belongs_to_valid_ingredient = $1 AND valid_ingredient_preparations.id = $2 AND valid_ingredients.id = $3 )"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
			exampleValidIngredientPreparation.ID,
			exampleValidIngredient.ID,
		}
		actualQuery, actualArgs := p.buildValidIngredientPreparationExistsQuery(exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ValidIngredientPreparationExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT valid_ingredient_preparations.id FROM valid_ingredient_preparations JOIN valid_ingredients ON valid_ingredient_preparations.belongs_to_valid_ingredient=valid_ingredients.id WHERE valid_ingredient_preparations.belongs_to_valid_ingredient = $1 AND valid_ingredient_preparations.id = $2 AND valid_ingredients.id = $3 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredientPreparation.ID,
				exampleValidIngredient.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.ValidIngredientPreparationExists(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredientPreparation.ID,
				exampleValidIngredient.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.ValidIngredientPreparationExists(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetValidIngredientPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.created_on, valid_ingredient_preparations.updated_on, valid_ingredient_preparations.archived_on, valid_ingredient_preparations.belongs_to_valid_ingredient FROM valid_ingredient_preparations JOIN valid_ingredients ON valid_ingredient_preparations.belongs_to_valid_ingredient=valid_ingredients.id WHERE valid_ingredient_preparations.belongs_to_valid_ingredient = $1 AND valid_ingredient_preparations.id = $2 AND valid_ingredients.id = $3"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
			exampleValidIngredientPreparation.ID,
			exampleValidIngredient.ID,
		}
		actualQuery, actualArgs := p.buildGetValidIngredientPreparationQuery(exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.created_on, valid_ingredient_preparations.updated_on, valid_ingredient_preparations.archived_on, valid_ingredient_preparations.belongs_to_valid_ingredient FROM valid_ingredient_preparations JOIN valid_ingredients ON valid_ingredient_preparations.belongs_to_valid_ingredient=valid_ingredients.id WHERE valid_ingredient_preparations.belongs_to_valid_ingredient = $1 AND valid_ingredient_preparations.id = $2 AND valid_ingredients.id = $3"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredientPreparation.ID,
				exampleValidIngredient.ID,
			).
			WillReturnRows(buildMockRowsFromValidIngredientPreparation(exampleValidIngredientPreparation))

		actual, err := p.GetValidIngredientPreparation(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparation, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredientPreparation.ID,
				exampleValidIngredient.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidIngredientPreparation(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)
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

func TestPostgres_buildGetValidIngredientPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.created_on, valid_ingredient_preparations.updated_on, valid_ingredient_preparations.archived_on, valid_ingredient_preparations.belongs_to_valid_ingredient, COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations JOIN valid_ingredients ON valid_ingredient_preparations.belongs_to_valid_ingredient=valid_ingredients.id WHERE valid_ingredient_preparations.archived_on IS NULL AND valid_ingredient_preparations.belongs_to_valid_ingredient = $1 AND valid_ingredients.id = $2 AND valid_ingredient_preparations.created_on > $3 AND valid_ingredient_preparations.created_on < $4 AND valid_ingredient_preparations.updated_on > $5 AND valid_ingredient_preparations.updated_on < $6 GROUP BY valid_ingredient_preparations.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
			exampleValidIngredient.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetValidIngredientPreparationsQuery(exampleValidIngredient.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetValidIngredientPreparations(T *testing.T) {
	T.Parallel()

	expectedListQuery := "SELECT valid_ingredient_preparations.id, valid_ingredient_preparations.notes, valid_ingredient_preparations.created_on, valid_ingredient_preparations.updated_on, valid_ingredient_preparations.archived_on, valid_ingredient_preparations.belongs_to_valid_ingredient, COUNT(valid_ingredient_preparations.id) FROM valid_ingredient_preparations JOIN valid_ingredients ON valid_ingredient_preparations.belongs_to_valid_ingredient=valid_ingredients.id WHERE valid_ingredient_preparations.archived_on IS NULL AND valid_ingredient_preparations.belongs_to_valid_ingredient = $1 AND valid_ingredients.id = $2 GROUP BY valid_ingredient_preparations.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidIngredientPreparationList := fakemodels.BuildFakeValidIngredientPreparationList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredient.ID,
			).
			WillReturnRows(
				buildMockRowsFromValidIngredientPreparation(
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[0],
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[1],
					&exampleValidIngredientPreparationList.ValidIngredientPreparations[2],
				),
			)

		actual, err := p.GetValidIngredientPreparations(ctx, exampleValidIngredient.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparationList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredient.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetValidIngredientPreparations(ctx, exampleValidIngredient.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredient.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetValidIngredientPreparations(ctx, exampleValidIngredient.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning valid ingredient preparation", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredient.ID,
			).
			WillReturnRows(buildErroneousMockRowFromValidIngredientPreparation(exampleValidIngredientPreparation))

		actual, err := p.GetValidIngredientPreparations(ctx, exampleValidIngredient.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateValidIngredientPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		expectedQuery := "INSERT INTO valid_ingredient_preparations (notes,belongs_to_valid_ingredient) VALUES ($1,$2) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.Notes,
			exampleValidIngredientPreparation.BelongsToValidIngredient,
		}
		actualQuery, actualArgs := p.buildCreateValidIngredientPreparationQuery(exampleValidIngredientPreparation)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO valid_ingredient_preparations (notes,belongs_to_valid_ingredient) VALUES ($1,$2) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleValidIngredientPreparation.ID, exampleValidIngredientPreparation.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleValidIngredientPreparation.Notes,
				exampleValidIngredientPreparation.BelongsToValidIngredient,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateValidIngredientPreparation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientPreparation, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleValidIngredientPreparation.Notes,
				exampleValidIngredientPreparation.BelongsToValidIngredient,
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

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		expectedQuery := "UPDATE valid_ingredient_preparations SET notes = $1, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_valid_ingredient = $2 AND id = $3 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleValidIngredientPreparation.Notes,
			exampleValidIngredientPreparation.BelongsToValidIngredient,
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

	expectedQuery := "UPDATE valid_ingredient_preparations SET notes = $1, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_valid_ingredient = $2 AND id = $3 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientPreparation.Notes,
				exampleValidIngredientPreparation.BelongsToValidIngredient,
				exampleValidIngredientPreparation.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateValidIngredientPreparation(ctx, exampleValidIngredientPreparation)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredientPreparation.Notes,
				exampleValidIngredientPreparation.BelongsToValidIngredient,
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

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		expectedQuery := "UPDATE valid_ingredient_preparations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_valid_ingredient = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
			exampleValidIngredientPreparation.ID,
		}
		actualQuery, actualArgs := p.buildArchiveValidIngredientPreparationQuery(exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveValidIngredientPreparation(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE valid_ingredient_preparations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_valid_ingredient = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredientPreparation.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveValidIngredientPreparation(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredientPreparation.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveValidIngredientPreparation(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
		exampleValidIngredientPreparation.BelongsToValidIngredient = exampleValidIngredient.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredientPreparation.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveValidIngredientPreparation(ctx, exampleValidIngredient.ID, exampleValidIngredientPreparation.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
