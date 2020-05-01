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

func buildMockRowsFromIngredientTagMapping(ingredientTagMappings ...*models.IngredientTagMapping) *sqlmock.Rows {
	includeCount := len(ingredientTagMappings) > 1
	columns := ingredientTagMappingsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}
	exampleRows := sqlmock.NewRows(columns)

	for _, x := range ingredientTagMappings {
		rowValues := []driver.Value{
			x.ID,
			x.ValidIngredientTagID,
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
			x.BelongsToValidIngredient,
		}

		if includeCount {
			rowValues = append(rowValues, len(ingredientTagMappings))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromIngredientTagMapping(x *models.IngredientTagMapping) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(ingredientTagMappingsTableColumns).AddRow(
		x.ArchivedOn,
		x.ValidIngredientTagID,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsToValidIngredient,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanIngredientTagMappings(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanIngredientTagMappings(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanIngredientTagMappings(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildIngredientTagMappingExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		expectedQuery := "SELECT EXISTS ( SELECT ingredient_tag_mappings.id FROM ingredient_tag_mappings JOIN valid_ingredients ON ingredient_tag_mappings.belongs_to_valid_ingredient=valid_ingredients.id WHERE ingredient_tag_mappings.belongs_to_valid_ingredient = $1 AND ingredient_tag_mappings.id = $2 AND valid_ingredients.id = $3 )"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
			exampleIngredientTagMapping.ID,
			exampleValidIngredient.ID,
		}
		actualQuery, actualArgs := p.buildIngredientTagMappingExistsQuery(exampleValidIngredient.ID, exampleIngredientTagMapping.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_IngredientTagMappingExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT ingredient_tag_mappings.id FROM ingredient_tag_mappings JOIN valid_ingredients ON ingredient_tag_mappings.belongs_to_valid_ingredient=valid_ingredients.id WHERE ingredient_tag_mappings.belongs_to_valid_ingredient = $1 AND ingredient_tag_mappings.id = $2 AND valid_ingredients.id = $3 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleIngredientTagMapping.ID,
				exampleValidIngredient.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.IngredientTagMappingExists(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleIngredientTagMapping.ID,
				exampleValidIngredient.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.IngredientTagMappingExists(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetIngredientTagMappingQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		expectedQuery := "SELECT ingredient_tag_mappings.id, ingredient_tag_mappings.valid_ingredient_tag_id, ingredient_tag_mappings.created_on, ingredient_tag_mappings.updated_on, ingredient_tag_mappings.archived_on, ingredient_tag_mappings.belongs_to_valid_ingredient FROM ingredient_tag_mappings JOIN valid_ingredients ON ingredient_tag_mappings.belongs_to_valid_ingredient=valid_ingredients.id WHERE ingredient_tag_mappings.belongs_to_valid_ingredient = $1 AND ingredient_tag_mappings.id = $2 AND valid_ingredients.id = $3"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
			exampleIngredientTagMapping.ID,
			exampleValidIngredient.ID,
		}
		actualQuery, actualArgs := p.buildGetIngredientTagMappingQuery(exampleValidIngredient.ID, exampleIngredientTagMapping.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetIngredientTagMapping(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT ingredient_tag_mappings.id, ingredient_tag_mappings.valid_ingredient_tag_id, ingredient_tag_mappings.created_on, ingredient_tag_mappings.updated_on, ingredient_tag_mappings.archived_on, ingredient_tag_mappings.belongs_to_valid_ingredient FROM ingredient_tag_mappings JOIN valid_ingredients ON ingredient_tag_mappings.belongs_to_valid_ingredient=valid_ingredients.id WHERE ingredient_tag_mappings.belongs_to_valid_ingredient = $1 AND ingredient_tag_mappings.id = $2 AND valid_ingredients.id = $3"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleIngredientTagMapping.ID,
				exampleValidIngredient.ID,
			).
			WillReturnRows(buildMockRowsFromIngredientTagMapping(exampleIngredientTagMapping))

		actual, err := p.GetIngredientTagMapping(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleIngredientTagMapping, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleIngredientTagMapping.ID,
				exampleValidIngredient.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetIngredientTagMapping(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllIngredientTagMappingsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(ingredient_tag_mappings.id) FROM ingredient_tag_mappings WHERE ingredient_tag_mappings.archived_on IS NULL"
		actualQuery := p.buildGetAllIngredientTagMappingsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllIngredientTagMappingsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(ingredient_tag_mappings.id) FROM ingredient_tag_mappings WHERE ingredient_tag_mappings.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllIngredientTagMappingsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetIngredientTagMappingsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT ingredient_tag_mappings.id, ingredient_tag_mappings.valid_ingredient_tag_id, ingredient_tag_mappings.created_on, ingredient_tag_mappings.updated_on, ingredient_tag_mappings.archived_on, ingredient_tag_mappings.belongs_to_valid_ingredient, COUNT(ingredient_tag_mappings.id) FROM ingredient_tag_mappings JOIN valid_ingredients ON ingredient_tag_mappings.belongs_to_valid_ingredient=valid_ingredients.id WHERE ingredient_tag_mappings.archived_on IS NULL AND ingredient_tag_mappings.belongs_to_valid_ingredient = $1 AND valid_ingredients.id = $2 AND ingredient_tag_mappings.created_on > $3 AND ingredient_tag_mappings.created_on < $4 AND ingredient_tag_mappings.updated_on > $5 AND ingredient_tag_mappings.updated_on < $6 GROUP BY ingredient_tag_mappings.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
			exampleValidIngredient.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetIngredientTagMappingsQuery(exampleValidIngredient.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetIngredientTagMappings(T *testing.T) {
	T.Parallel()

	expectedListQuery := "SELECT ingredient_tag_mappings.id, ingredient_tag_mappings.valid_ingredient_tag_id, ingredient_tag_mappings.created_on, ingredient_tag_mappings.updated_on, ingredient_tag_mappings.archived_on, ingredient_tag_mappings.belongs_to_valid_ingredient, COUNT(ingredient_tag_mappings.id) FROM ingredient_tag_mappings JOIN valid_ingredients ON ingredient_tag_mappings.belongs_to_valid_ingredient=valid_ingredients.id WHERE ingredient_tag_mappings.archived_on IS NULL AND ingredient_tag_mappings.belongs_to_valid_ingredient = $1 AND valid_ingredients.id = $2 GROUP BY ingredient_tag_mappings.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleIngredientTagMappingList := fakemodels.BuildFakeIngredientTagMappingList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredient.ID,
			).
			WillReturnRows(
				buildMockRowsFromIngredientTagMapping(
					&exampleIngredientTagMappingList.IngredientTagMappings[0],
					&exampleIngredientTagMappingList.IngredientTagMappings[1],
					&exampleIngredientTagMappingList.IngredientTagMappings[2],
				),
			)

		actual, err := p.GetIngredientTagMappings(ctx, exampleValidIngredient.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleIngredientTagMappingList, actual)

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

		actual, err := p.GetIngredientTagMappings(ctx, exampleValidIngredient.ID, filter)
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

		actual, err := p.GetIngredientTagMappings(ctx, exampleValidIngredient.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning ingredient tag mapping", func(t *testing.T) {
		ctx := context.Background()

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleValidIngredient.ID,
			).
			WillReturnRows(buildErroneousMockRowFromIngredientTagMapping(exampleIngredientTagMapping))

		actual, err := p.GetIngredientTagMappings(ctx, exampleValidIngredient.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateIngredientTagMappingQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		expectedQuery := "INSERT INTO ingredient_tag_mappings (valid_ingredient_tag_id,belongs_to_valid_ingredient) VALUES ($1,$2) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleIngredientTagMapping.ValidIngredientTagID,
			exampleIngredientTagMapping.BelongsToValidIngredient,
		}
		actualQuery, actualArgs := p.buildCreateIngredientTagMappingQuery(exampleIngredientTagMapping)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateIngredientTagMapping(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO ingredient_tag_mappings (valid_ingredient_tag_id,belongs_to_valid_ingredient) VALUES ($1,$2) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleIngredientTagMapping.ID, exampleIngredientTagMapping.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleIngredientTagMapping.ValidIngredientTagID,
				exampleIngredientTagMapping.BelongsToValidIngredient,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateIngredientTagMapping(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleIngredientTagMapping, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID
		exampleInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleIngredientTagMapping.ValidIngredientTagID,
				exampleIngredientTagMapping.BelongsToValidIngredient,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateIngredientTagMapping(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateIngredientTagMappingQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		expectedQuery := "UPDATE ingredient_tag_mappings SET valid_ingredient_tag_id = $1, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_valid_ingredient = $2 AND id = $3 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleIngredientTagMapping.ValidIngredientTagID,
			exampleIngredientTagMapping.BelongsToValidIngredient,
			exampleIngredientTagMapping.ID,
		}
		actualQuery, actualArgs := p.buildUpdateIngredientTagMappingQuery(exampleIngredientTagMapping)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateIngredientTagMapping(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE ingredient_tag_mappings SET valid_ingredient_tag_id = $1, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_valid_ingredient = $2 AND id = $3 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleIngredientTagMapping.ValidIngredientTagID,
				exampleIngredientTagMapping.BelongsToValidIngredient,
				exampleIngredientTagMapping.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateIngredientTagMapping(ctx, exampleIngredientTagMapping)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleIngredientTagMapping.ValidIngredientTagID,
				exampleIngredientTagMapping.BelongsToValidIngredient,
				exampleIngredientTagMapping.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateIngredientTagMapping(ctx, exampleIngredientTagMapping)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveIngredientTagMappingQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		expectedQuery := "UPDATE ingredient_tag_mappings SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_valid_ingredient = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleValidIngredient.ID,
			exampleIngredientTagMapping.ID,
		}
		actualQuery, actualArgs := p.buildArchiveIngredientTagMappingQuery(exampleValidIngredient.ID, exampleIngredientTagMapping.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveIngredientTagMapping(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE ingredient_tag_mappings SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_valid_ingredient = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleIngredientTagMapping.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveIngredientTagMapping(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleIngredientTagMapping.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveIngredientTagMapping(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
		exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
		exampleIngredientTagMapping.BelongsToValidIngredient = exampleValidIngredient.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleValidIngredient.ID,
				exampleIngredientTagMapping.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveIngredientTagMapping(ctx, exampleValidIngredient.ID, exampleIngredientTagMapping.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
