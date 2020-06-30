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

func buildMockRowsFromRecipeIteration(recipeIterations ...*models.RecipeIteration) *sqlmock.Rows {
	includeCount := len(recipeIterations) > 1
	columns := recipeIterationsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}
	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeIterations {
		rowValues := []driver.Value{
			x.ID,
			x.EndDifficultyRating,
			x.EndComplexityRating,
			x.EndTasteRating,
			x.EndOverallRating,
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipe,
		}

		if includeCount {
			rowValues = append(rowValues, len(recipeIterations))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipeIteration(x *models.RecipeIteration) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeIterationsTableColumns).AddRow(
		x.ArchivedOn,
		x.EndDifficultyRating,
		x.EndComplexityRating,
		x.EndTasteRating,
		x.EndOverallRating,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsToRecipe,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRecipeIterations(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanRecipeIterations(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanRecipeIterations(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRecipeIterationExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "SELECT EXISTS ( SELECT recipe_iterations.id FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE recipe_iterations.belongs_to_recipe = $1 AND recipe_iterations.id = $2 AND recipes.id = $3 )"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeIteration.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildRecipeIterationExistsQuery(exampleRecipe.ID, exampleRecipeIteration.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RecipeIterationExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT recipe_iterations.id FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE recipe_iterations.belongs_to_recipe = $1 AND recipe_iterations.id = $2 AND recipes.id = $3 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RecipeIterationExists(ctx, exampleRecipe.ID, exampleRecipeIteration.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RecipeIterationExists(ctx, exampleRecipe.ID, exampleRecipeIteration.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "SELECT recipe_iterations.id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE recipe_iterations.belongs_to_recipe = $1 AND recipe_iterations.id = $2 AND recipes.id = $3"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeIteration.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeIterationQuery(exampleRecipe.ID, exampleRecipeIteration.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeIteration(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT recipe_iterations.id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE recipe_iterations.belongs_to_recipe = $1 AND recipe_iterations.id = $2 AND recipes.id = $3"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromRecipeIteration(exampleRecipeIteration))

		actual, err := p.GetRecipeIteration(ctx, exampleRecipe.ID, exampleRecipeIteration.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIteration, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeIteration(ctx, exampleRecipe.ID, exampleRecipeIteration.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeIterationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(recipe_iterations.id) FROM recipe_iterations WHERE recipe_iterations.archived_on IS NULL"
		actualQuery := p.buildGetAllRecipeIterationsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeIterationsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_iterations.id) FROM recipe_iterations WHERE recipe_iterations.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeIterationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeIterationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_iterations.id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe, (SELECT COUNT(recipe_iterations.id) FROM recipe_iterations WHERE recipe_iterations.archived_on IS NULL) FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE recipe_iterations.archived_on IS NULL AND recipe_iterations.belongs_to_recipe = $1 AND recipes.id = $2 AND recipe_iterations.created_on > $3 AND recipe_iterations.created_on < $4 AND recipe_iterations.updated_on > $5 AND recipe_iterations.updated_on < $6 ORDER BY recipe_iterations.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRecipeIterationsQuery(exampleRecipe.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeIterations(T *testing.T) {
	T.Parallel()

	expectedListQuery := "SELECT recipe_iterations.id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe, (SELECT COUNT(recipe_iterations.id) FROM recipe_iterations WHERE recipe_iterations.archived_on IS NULL) FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE recipe_iterations.archived_on IS NULL AND recipe_iterations.belongs_to_recipe = $1 AND recipes.id = $2 ORDER BY recipe_iterations.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeIterationList := fakemodels.BuildFakeRecipeIterationList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeIteration(
					&exampleRecipeIterationList.RecipeIterations[0],
					&exampleRecipeIterationList.RecipeIterations[1],
					&exampleRecipeIterationList.RecipeIterations[2],
				),
			)

		actual, err := p.GetRecipeIterations(ctx, exampleRecipe.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeIterations(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeIterations(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe iteration", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeIteration(exampleRecipeIteration))

		actual, err := p.GetRecipeIterations(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "INSERT INTO recipe_iterations (end_difficulty_rating,end_complexity_rating,end_taste_rating,end_overall_rating,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipeIteration.EndDifficultyRating,
			exampleRecipeIteration.EndComplexityRating,
			exampleRecipeIteration.EndTasteRating,
			exampleRecipeIteration.EndOverallRating,
			exampleRecipeIteration.BelongsToRecipe,
		}
		actualQuery, actualArgs := p.buildCreateRecipeIterationQuery(exampleRecipeIteration)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRecipeIteration(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO recipe_iterations (end_difficulty_rating,end_complexity_rating,end_taste_rating,end_overall_rating,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRecipeIteration.ID, exampleRecipeIteration.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeIteration.EndDifficultyRating,
				exampleRecipeIteration.EndComplexityRating,
				exampleRecipeIteration.EndTasteRating,
				exampleRecipeIteration.EndOverallRating,
				exampleRecipeIteration.BelongsToRecipe,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeIteration(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIteration, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeIteration.EndDifficultyRating,
				exampleRecipeIteration.EndComplexityRating,
				exampleRecipeIteration.EndTasteRating,
				exampleRecipeIteration.EndOverallRating,
				exampleRecipeIteration.BelongsToRecipe,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipeIteration(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "UPDATE recipe_iterations SET end_difficulty_rating = $1, end_complexity_rating = $2, end_taste_rating = $3, end_overall_rating = $4, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $5 AND id = $6 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleRecipeIteration.EndDifficultyRating,
			exampleRecipeIteration.EndComplexityRating,
			exampleRecipeIteration.EndTasteRating,
			exampleRecipeIteration.EndOverallRating,
			exampleRecipeIteration.BelongsToRecipe,
			exampleRecipeIteration.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRecipeIterationQuery(exampleRecipeIteration)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRecipeIteration(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_iterations SET end_difficulty_rating = $1, end_complexity_rating = $2, end_taste_rating = $3, end_overall_rating = $4, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $5 AND id = $6 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIteration.EndDifficultyRating,
				exampleRecipeIteration.EndComplexityRating,
				exampleRecipeIteration.EndTasteRating,
				exampleRecipeIteration.EndOverallRating,
				exampleRecipeIteration.BelongsToRecipe,
				exampleRecipeIteration.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeIteration(ctx, exampleRecipeIteration)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIteration.EndDifficultyRating,
				exampleRecipeIteration.EndComplexityRating,
				exampleRecipeIteration.EndTasteRating,
				exampleRecipeIteration.EndOverallRating,
				exampleRecipeIteration.BelongsToRecipe,
				exampleRecipeIteration.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipeIteration(ctx, exampleRecipeIteration)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "UPDATE recipe_iterations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeIteration.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRecipeIterationQuery(exampleRecipe.ID, exampleRecipeIteration.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRecipeIteration(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_iterations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeIteration(ctx, exampleRecipe.ID, exampleRecipeIteration.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRecipeIteration(ctx, exampleRecipe.ID, exampleRecipeIteration.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeIteration(ctx, exampleRecipe.ID, exampleRecipeIteration.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
