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

func buildMockRowsFromRecipeStepIngredients(recipeStepIngredients ...*models.RecipeStepIngredient) *sqlmock.Rows {
	columns := recipeStepIngredientsTableColumns

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepIngredients {
		rowValues := []driver.Value{
			x.ID,
			x.IngredientID,
			x.QuantityType,
			x.QuantityValue,
			x.QuantityNotes,
			x.ProductOfRecipe,
			x.IngredientNotes,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipeStep,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipeStepIngredient(x *models.RecipeStepIngredient) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepIngredientsTableColumns).AddRow(
		x.ArchivedOn,
		x.IngredientID,
		x.QuantityType,
		x.QuantityValue,
		x.QuantityNotes,
		x.ProductOfRecipe,
		x.IngredientNotes,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.BelongsToRecipeStep,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := p.scanRecipeStepIngredients(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanRecipeStepIngredients(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRecipeStepIngredientExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "SELECT EXISTS ( SELECT recipe_step_ingredients.id FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_step_ingredients.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5 )"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepIngredient.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildRecipeStepIngredientExistsQuery(exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RecipeStepIngredientExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT recipe_step_ingredients.id FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_step_ingredients.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepIngredient.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RecipeStepIngredientExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepIngredient.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RecipeStepIngredientExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_step_ingredients.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepIngredient.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepIngredientQuery(exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_step_ingredients.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepIngredient.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(exampleRecipeStepIngredient))

		actual, err := p.GetRecipeStepIngredient(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredient, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepIngredient.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepIngredient(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeStepIngredientsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(recipe_step_ingredients.id) FROM recipe_step_ingredients WHERE recipe_step_ingredients.archived_on IS NULL"
		actualQuery := p.buildGetAllRecipeStepIngredientsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeStepIngredientsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_step_ingredients.id) FROM recipe_step_ingredients WHERE recipe_step_ingredients.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeStepIngredientsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfRecipeStepIngredientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients WHERE recipe_step_ingredients.id > $1 AND recipe_step_ingredients.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfRecipeStepIngredientsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(recipe_step_ingredients.id) FROM recipe_step_ingredients WHERE recipe_step_ingredients.archived_on IS NULL"
	expectedGetQuery := "SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients WHERE recipe_step_ingredients.id > $1 AND recipe_step_ingredients.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeStepIngredientList := fakemodels.BuildFakeRecipeStepIngredientList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepIngredients(
					&exampleRecipeStepIngredientList.RecipeStepIngredients[0],
					&exampleRecipeStepIngredientList.RecipeStepIngredients[1],
					&exampleRecipeStepIngredientList.RecipeStepIngredients[2],
				),
			)

		out := make(chan []models.RecipeStepIngredient)
		doneChan := make(chan bool, 1)

		err := p.GetAllRecipeStepIngredients(ctx, out)
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

		out := make(chan []models.RecipeStepIngredient)

		err := p.GetAllRecipeStepIngredients(ctx, out)
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

		out := make(chan []models.RecipeStepIngredient)

		err := p.GetAllRecipeStepIngredients(ctx, out)
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

		out := make(chan []models.RecipeStepIngredient)

		err := p.GetAllRecipeStepIngredients(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepIngredient(exampleRecipeStepIngredient))

		out := make(chan []models.RecipeStepIngredient)

		err := p.GetAllRecipeStepIngredients(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepIngredientsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 AND recipe_step_ingredients.created_on > $5 AND recipe_step_ingredients.created_on < $6 AND recipe_step_ingredients.last_updated_on > $7 AND recipe_step_ingredients.last_updated_on < $8 ORDER BY recipe_step_ingredients.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepIngredientsQuery(exampleRecipe.ID, exampleRecipeStep.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 ORDER BY recipe_step_ingredients.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepIngredientList := fakemodels.BuildFakeRecipeStepIngredientList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepIngredients(
					&exampleRecipeStepIngredientList.RecipeStepIngredients[0],
					&exampleRecipeStepIngredientList.RecipeStepIngredients[1],
					&exampleRecipeStepIngredientList.RecipeStepIngredients[2],
				),
			)

		actual, err := p.GetRecipeStepIngredients(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepIngredients(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeStepIngredients(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step ingredient", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepIngredient(exampleRecipeStepIngredient))

		actual, err := p.GetRecipeStepIngredients(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepIngredientsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM (SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_ingredients WHERE recipe_step_ingredients.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepIngredientsWithIDsQuery(exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepIngredientsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleRecipeStepIngredientList := fakemodels.BuildFakeRecipeStepIngredientList()
		var exampleIDs []uint64
		for _, recipeStepIngredient := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleIDs = append(exampleIDs, recipeStepIngredient.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM (SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_ingredients WHERE recipe_step_ingredients.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepIngredients(
					&exampleRecipeStepIngredientList.RecipeStepIngredients[0],
					&exampleRecipeStepIngredientList.RecipeStepIngredients[1],
					&exampleRecipeStepIngredientList.RecipeStepIngredients[2],
				),
			)

		actual, err := p.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList.RecipeStepIngredients, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM (SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_ingredients WHERE recipe_step_ingredients.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM (SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_ingredients WHERE recipe_step_ingredients.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step ingredient", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM (SELECT recipe_step_ingredients.id, recipe_step_ingredients.ingredient_id, recipe_step_ingredients.quantity_type, recipe_step_ingredients.quantity_value, recipe_step_ingredients.quantity_notes, recipe_step_ingredients.product_of_recipe, recipe_step_ingredients.ingredient_notes, recipe_step_ingredients.created_on, recipe_step_ingredients.last_updated_on, recipe_step_ingredients.archived_on, recipe_step_ingredients.belongs_to_recipe_step FROM recipe_step_ingredients JOIN recipe_steps ON recipe_step_ingredients.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_ingredients.archived_on IS NULL AND recipe_step_ingredients.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_ingredients WHERE recipe_step_ingredients.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepIngredient(exampleRecipeStepIngredient))

		actual, err := p.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "INSERT INTO recipe_step_ingredients (ingredient_id,quantity_type,quantity_value,quantity_notes,product_of_recipe,ingredient_notes,belongs_to_recipe_step) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipeStepIngredient.IngredientID,
			exampleRecipeStepIngredient.QuantityType,
			exampleRecipeStepIngredient.QuantityValue,
			exampleRecipeStepIngredient.QuantityNotes,
			exampleRecipeStepIngredient.ProductOfRecipe,
			exampleRecipeStepIngredient.IngredientNotes,
			exampleRecipeStepIngredient.BelongsToRecipeStep,
		}
		actualQuery, actualArgs := p.buildCreateRecipeStepIngredientQuery(exampleRecipeStepIngredient)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO recipe_step_ingredients (ingredient_id,quantity_type,quantity_value,quantity_notes,product_of_recipe,ingredient_notes,belongs_to_recipe_step) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRecipeStepIngredient.ID, exampleRecipeStepIngredient.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStepIngredient.IngredientID,
				exampleRecipeStepIngredient.QuantityType,
				exampleRecipeStepIngredient.QuantityValue,
				exampleRecipeStepIngredient.QuantityNotes,
				exampleRecipeStepIngredient.ProductOfRecipe,
				exampleRecipeStepIngredient.IngredientNotes,
				exampleRecipeStepIngredient.BelongsToRecipeStep,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeStepIngredient(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredient, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStepIngredient.IngredientID,
				exampleRecipeStepIngredient.QuantityType,
				exampleRecipeStepIngredient.QuantityValue,
				exampleRecipeStepIngredient.QuantityNotes,
				exampleRecipeStepIngredient.ProductOfRecipe,
				exampleRecipeStepIngredient.IngredientNotes,
				exampleRecipeStepIngredient.BelongsToRecipeStep,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipeStepIngredient(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "UPDATE recipe_step_ingredients SET ingredient_id = $1, quantity_type = $2, quantity_value = $3, quantity_notes = $4, product_of_recipe = $5, ingredient_notes = $6, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_step = $7 AND id = $8 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleRecipeStepIngredient.IngredientID,
			exampleRecipeStepIngredient.QuantityType,
			exampleRecipeStepIngredient.QuantityValue,
			exampleRecipeStepIngredient.QuantityNotes,
			exampleRecipeStepIngredient.ProductOfRecipe,
			exampleRecipeStepIngredient.IngredientNotes,
			exampleRecipeStepIngredient.BelongsToRecipeStep,
			exampleRecipeStepIngredient.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRecipeStepIngredientQuery(exampleRecipeStepIngredient)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_step_ingredients SET ingredient_id = $1, quantity_type = $2, quantity_value = $3, quantity_notes = $4, product_of_recipe = $5, ingredient_notes = $6, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_step = $7 AND id = $8 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStepIngredient.IngredientID,
				exampleRecipeStepIngredient.QuantityType,
				exampleRecipeStepIngredient.QuantityValue,
				exampleRecipeStepIngredient.QuantityNotes,
				exampleRecipeStepIngredient.ProductOfRecipe,
				exampleRecipeStepIngredient.IngredientNotes,
				exampleRecipeStepIngredient.BelongsToRecipeStep,
				exampleRecipeStepIngredient.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeStepIngredient(ctx, exampleRecipeStepIngredient)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStepIngredient.IngredientID,
				exampleRecipeStepIngredient.QuantityType,
				exampleRecipeStepIngredient.QuantityValue,
				exampleRecipeStepIngredient.QuantityNotes,
				exampleRecipeStepIngredient.ProductOfRecipe,
				exampleRecipeStepIngredient.IngredientNotes,
				exampleRecipeStepIngredient.BelongsToRecipeStep,
				exampleRecipeStepIngredient.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipeStepIngredient(ctx, exampleRecipeStepIngredient)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeStepIngredientQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "UPDATE recipe_step_ingredients SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepIngredient.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRecipeStepIngredientQuery(exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_step_ingredients SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepIngredient.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeStepIngredient(ctx, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepIngredient.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRecipeStepIngredient(ctx, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)
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
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepIngredient.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeStepIngredient(ctx, exampleRecipeStep.ID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
