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

func buildMockRowsFromRecipeIterations(recipeIterations ...*models.RecipeIteration) *sqlmock.Rows {
	columns := recipeIterationsTableColumns

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeIterations {
		rowValues := []driver.Value{
			x.ID,
			x.RecipeID,
			x.EndDifficultyRating,
			x.EndComplexityRating,
			x.EndTasteRating,
			x.EndOverallRating,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipe,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipeIteration(x *models.RecipeIteration) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeIterationsTableColumns).AddRow(
		x.ArchivedOn,
		x.RecipeID,
		x.EndDifficultyRating,
		x.EndComplexityRating,
		x.EndTasteRating,
		x.EndOverallRating,
		x.CreatedOn,
		x.LastUpdatedOn,
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

		_, err := p.scanRecipeIterations(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanRecipeIterations(mockRows)
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

		expectedQuery := "SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE recipe_iterations.belongs_to_recipe = $1 AND recipe_iterations.id = $2 AND recipes.id = $3"
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
	expectedQuery := "SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE recipe_iterations.belongs_to_recipe = $1 AND recipe_iterations.id = $2 AND recipes.id = $3"

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
			WillReturnRows(buildMockRowsFromRecipeIterations(exampleRecipeIteration))

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

func TestPostgres_buildGetBatchOfRecipeIterationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations WHERE recipe_iterations.id > $1 AND recipe_iterations.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfRecipeIterationsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllRecipeIterations(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(recipe_iterations.id) FROM recipe_iterations WHERE recipe_iterations.archived_on IS NULL"
	expectedGetQuery := "SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations WHERE recipe_iterations.id > $1 AND recipe_iterations.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeIterationList := fakemodels.BuildFakeRecipeIterationList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromRecipeIterations(
					&exampleRecipeIterationList.RecipeIterations[0],
					&exampleRecipeIterationList.RecipeIterations[1],
					&exampleRecipeIterationList.RecipeIterations[2],
				),
			)

		out := make(chan []models.RecipeIteration)
		doneChan := make(chan bool, 1)

		err := p.GetAllRecipeIterations(ctx, out)
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

		out := make(chan []models.RecipeIteration)

		err := p.GetAllRecipeIterations(ctx, out)
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

		out := make(chan []models.RecipeIteration)

		err := p.GetAllRecipeIterations(ctx, out)
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

		out := make(chan []models.RecipeIteration)

		err := p.GetAllRecipeIterations(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromRecipeIteration(exampleRecipeIteration))

		out := make(chan []models.RecipeIteration)

		err := p.GetAllRecipeIterations(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeIterationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE recipe_iterations.archived_on IS NULL AND recipe_iterations.belongs_to_recipe = $1 AND recipes.id = $2 AND recipe_iterations.created_on > $3 AND recipe_iterations.created_on < $4 AND recipe_iterations.last_updated_on > $5 AND recipe_iterations.last_updated_on < $6 ORDER BY recipe_iterations.id LIMIT 20 OFFSET 180"
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

	expectedQuery := "SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE recipe_iterations.archived_on IS NULL AND recipe_iterations.belongs_to_recipe = $1 AND recipes.id = $2 ORDER BY recipe_iterations.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeIterationList := fakemodels.BuildFakeRecipeIterationList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeIterations(
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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

func TestPostgres_buildGetRecipeIterationsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM (SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_iterations.archived_on IS NULL AND recipe_iterations.belongs_to_recipe = $1 AND recipes.id = $2 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_iterations WHERE recipe_iterations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeIterationsWithIDsQuery(exampleRecipe.ID, defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeIterationsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)

		exampleRecipeIterationList := fakemodels.BuildFakeRecipeIterationList()
		var exampleIDs []uint64
		for _, recipeIteration := range exampleRecipeIterationList.RecipeIterations {
			exampleIDs = append(exampleIDs, recipeIteration.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM (SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_iterations.archived_on IS NULL AND recipe_iterations.belongs_to_recipe = $1 AND recipes.id = $2 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_iterations WHERE recipe_iterations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeIterations(
					&exampleRecipeIterationList.RecipeIterations[0],
					&exampleRecipeIterationList.RecipeIterations[1],
					&exampleRecipeIterationList.RecipeIterations[2],
				),
			)

		actual, err := p.GetRecipeIterationsWithIDs(ctx, exampleRecipe.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationList.RecipeIterations, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM (SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_iterations.archived_on IS NULL AND recipe_iterations.belongs_to_recipe = $1 AND recipes.id = $2 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_iterations WHERE recipe_iterations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeIterationsWithIDs(ctx, exampleRecipe.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM (SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_iterations.archived_on IS NULL AND recipe_iterations.belongs_to_recipe = $1 AND recipes.id = $2 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_iterations WHERE recipe_iterations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeIterationsWithIDs(ctx, exampleRecipe.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe iteration", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM (SELECT recipe_iterations.id, recipe_iterations.recipe_id, recipe_iterations.end_difficulty_rating, recipe_iterations.end_complexity_rating, recipe_iterations.end_taste_rating, recipe_iterations.end_overall_rating, recipe_iterations.created_on, recipe_iterations.last_updated_on, recipe_iterations.archived_on, recipe_iterations.belongs_to_recipe FROM recipe_iterations JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_iterations.archived_on IS NULL AND recipe_iterations.belongs_to_recipe = $1 AND recipes.id = $2 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_iterations WHERE recipe_iterations.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeIteration(exampleRecipeIteration))

		actual, err := p.GetRecipeIterationsWithIDs(ctx, exampleRecipe.ID, defaultLimit, exampleIDs)

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

		expectedQuery := "INSERT INTO recipe_iterations (recipe_id,end_difficulty_rating,end_complexity_rating,end_taste_rating,end_overall_rating,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipeIteration.RecipeID,
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

	expectedCreationQuery := "INSERT INTO recipe_iterations (recipe_id,end_difficulty_rating,end_complexity_rating,end_taste_rating,end_overall_rating,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_on"

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
				exampleRecipeIteration.RecipeID,
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
				exampleRecipeIteration.RecipeID,
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

		expectedQuery := "UPDATE recipe_iterations SET recipe_id = $1, end_difficulty_rating = $2, end_complexity_rating = $3, end_taste_rating = $4, end_overall_rating = $5, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $6 AND id = $7 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleRecipeIteration.RecipeID,
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

	expectedQuery := "UPDATE recipe_iterations SET recipe_id = $1, end_difficulty_rating = $2, end_complexity_rating = $3, end_taste_rating = $4, end_overall_rating = $5, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $6 AND id = $7 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIteration.RecipeID,
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
				exampleRecipeIteration.RecipeID,
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

		expectedQuery := "UPDATE recipe_iterations SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"
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

	expectedQuery := "UPDATE recipe_iterations SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"

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
