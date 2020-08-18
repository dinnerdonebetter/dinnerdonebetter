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

func buildMockRowsFromRecipes(recipes ...*models.Recipe) *sqlmock.Rows {
	columns := recipesTableColumns

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipes {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Source,
			x.Description,
			x.InspiredByRecipeID,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToUser,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipe(x *models.Recipe) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipesTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.Source,
		x.Description,
		x.InspiredByRecipeID,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.BelongsToUser,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRecipes(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := p.scanRecipes(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanRecipes(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRecipeExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT EXISTS ( SELECT recipes.id FROM recipes WHERE recipes.id = $1 )"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildRecipeExistsQuery(exampleRecipe.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RecipeExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT recipes.id FROM recipes WHERE recipes.id = $1 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RecipeExists(ctx, exampleRecipe.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RecipeExists(ctx, exampleRecipe.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		expectedQuery := "SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM recipes WHERE recipes.id = $1"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeQuery(exampleRecipe.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipe(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM recipes WHERE recipes.id = $1"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromRecipes(exampleRecipe))

		actual, err := p.GetRecipe(ctx, exampleRecipe.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipe, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipe(ctx, exampleRecipe.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipesCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(recipes.id) FROM recipes WHERE recipes.archived_on IS NULL"
		actualQuery := p.buildGetAllRecipesCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipesCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipes.id) FROM recipes WHERE recipes.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipesCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfRecipesQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM recipes WHERE recipes.id > $1 AND recipes.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfRecipesQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllRecipes(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(recipes.id) FROM recipes WHERE recipes.archived_on IS NULL"
	expectedGetQuery := "SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM recipes WHERE recipes.id > $1 AND recipes.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeList := fakemodels.BuildFakeRecipeList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromRecipes(
					&exampleRecipeList.Recipes[0],
					&exampleRecipeList.Recipes[1],
					&exampleRecipeList.Recipes[2],
				),
			)

		out := make(chan []models.Recipe)
		doneChan := make(chan bool, 1)

		err := p.GetAllRecipes(ctx, out)
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

		out := make(chan []models.Recipe)

		err := p.GetAllRecipes(ctx, out)
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

		out := make(chan []models.Recipe)

		err := p.GetAllRecipes(ctx, out)
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

		out := make(chan []models.Recipe)

		err := p.GetAllRecipes(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipe := fakemodels.BuildFakeRecipe()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromRecipe(exampleRecipe))

		out := make(chan []models.Recipe)

		err := p.GetAllRecipes(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipesQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM recipes WHERE recipes.archived_on IS NULL AND recipes.created_on > $1 AND recipes.created_on < $2 AND recipes.last_updated_on > $3 AND recipes.last_updated_on < $4 ORDER BY recipes.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRecipesQuery(filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipes(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM recipes WHERE recipes.archived_on IS NULL ORDER BY recipes.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeList := fakemodels.BuildFakeRecipeList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildMockRowsFromRecipes(
					&exampleRecipeList.Recipes[0],
					&exampleRecipeList.Recipes[1],
					&exampleRecipeList.Recipes[2],
				),
			)

		actual, err := p.GetRecipes(ctx, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		actual, err := p.GetRecipes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(
				buildErroneousMockRowFromRecipe(exampleRecipe),
			)

		actual, err := p.GetRecipes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipesWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM (SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM recipes JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipes WHERE recipes.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}(nil)
		actualQuery, actualArgs := p.buildGetRecipesWithIDsQuery(defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipesWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleRecipeList := fakemodels.BuildFakeRecipeList()
		var exampleIDs []uint64
		for _, recipe := range exampleRecipeList.Recipes {
			exampleIDs = append(exampleIDs, recipe.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM (SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM recipes JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipes WHERE recipes.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(
				buildMockRowsFromRecipes(
					&exampleRecipeList.Recipes[0],
					&exampleRecipeList.Recipes[1],
					&exampleRecipeList.Recipes[2],
				),
			)

		actual, err := p.GetRecipesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeList.Recipes, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM (SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM recipes JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipes WHERE recipes.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM (SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM recipes JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipes WHERE recipes.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM (SELECT recipes.id, recipes.name, recipes.source, recipes.description, recipes.inspired_by_recipe_id, recipes.created_on, recipes.last_updated_on, recipes.archived_on, recipes.belongs_to_user FROM recipes JOIN unnest('{%s}'::int[]) WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipes WHERE recipes.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleRecipe := fakemodels.BuildFakeRecipe()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(buildErroneousMockRowFromRecipe(exampleRecipe))

		actual, err := p.GetRecipesWithIDs(ctx, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		expectedQuery := "INSERT INTO recipes (name,source,description,inspired_by_recipe_id,belongs_to_user) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipe.Name,
			exampleRecipe.Source,
			exampleRecipe.Description,
			exampleRecipe.InspiredByRecipeID,
			exampleRecipe.BelongsToUser,
		}
		actualQuery, actualArgs := p.buildCreateRecipeQuery(exampleRecipe)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRecipe(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO recipes (name,source,description,inspired_by_recipe_id,belongs_to_user) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRecipe.ID, exampleRecipe.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipe.Name,
				exampleRecipe.Source,
				exampleRecipe.Description,
				exampleRecipe.InspiredByRecipeID,
				exampleRecipe.BelongsToUser,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipe(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipe, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipe.Name,
				exampleRecipe.Source,
				exampleRecipe.Description,
				exampleRecipe.InspiredByRecipeID,
				exampleRecipe.BelongsToUser,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipe(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE recipes SET name = $1, source = $2, description = $3, inspired_by_recipe_id = $4, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $5 AND id = $6 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleRecipe.Name,
			exampleRecipe.Source,
			exampleRecipe.Description,
			exampleRecipe.InspiredByRecipeID,
			exampleRecipe.BelongsToUser,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRecipeQuery(exampleRecipe)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRecipe(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipes SET name = $1, source = $2, description = $3, inspired_by_recipe_id = $4, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_user = $5 AND id = $6 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.Name,
				exampleRecipe.Source,
				exampleRecipe.Description,
				exampleRecipe.InspiredByRecipeID,
				exampleRecipe.BelongsToUser,
				exampleRecipe.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipe(ctx, exampleRecipe)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.Name,
				exampleRecipe.Source,
				exampleRecipe.Description,
				exampleRecipe.InspiredByRecipeID,
				exampleRecipe.BelongsToUser,
				exampleRecipe.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipe(ctx, exampleRecipe)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		expectedQuery := "UPDATE recipes SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleUser.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRecipeQuery(exampleRecipe.ID, exampleUser.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRecipe(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipes SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_user = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleRecipe.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipe(ctx, exampleRecipe.ID, exampleUser.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleRecipe.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRecipe(ctx, exampleRecipe.ID, exampleUser.ID)
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

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleUser.ID,
				exampleRecipe.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipe(ctx, exampleRecipe.ID, exampleUser.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
