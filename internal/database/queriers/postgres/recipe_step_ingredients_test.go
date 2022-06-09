package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromRecipeStepIngredients(includeCounts bool, filteredCount uint64, recipeStepIngredients ...*types.RecipeStepIngredient) *sqlmock.Rows {
	columns := recipeStepIngredientsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepIngredients {
		rowValues := []driver.Value{
			x.ID,
			x.IngredientID,
			x.QuantityType,
			x.QuantityValue,
			x.QuantityNotes,
			x.ProductOfRecipeStep,
			x.IngredientNotes,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipeStep,
			x.Ingredient.ID,
			x.Ingredient.Name,
			x.Ingredient.Variant,
			x.Ingredient.Description,
			x.Ingredient.Warning,
			x.Ingredient.ContainsEgg,
			x.Ingredient.ContainsDairy,
			x.Ingredient.ContainsPeanut,
			x.Ingredient.ContainsTreeNut,
			x.Ingredient.ContainsSoy,
			x.Ingredient.ContainsWheat,
			x.Ingredient.ContainsShellfish,
			x.Ingredient.ContainsSesame,
			x.Ingredient.ContainsFish,
			x.Ingredient.ContainsGluten,
			x.Ingredient.AnimalFlesh,
			x.Ingredient.AnimalDerived,
			x.Ingredient.Volumetric,
			x.Ingredient.IconPath,
			x.Ingredient.CreatedOn,
			x.Ingredient.LastUpdatedOn,
			x.Ingredient.ArchivedOn,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(recipeStepIngredients))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipeStepIngredients(ctx, mockRows, false)
		assert.Error(t, err)
	})

	T.Run("logs row closing errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipeStepIngredients(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_RecipeStepIngredientExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepIngredientExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepIngredientExists(ctx, "", exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, "", exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepIngredientExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepIngredientExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.RecipeStepIngredientExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.Ingredient = types.ValidIngredient{}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, exampleRecipeStepIngredient))

		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredient, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredient(ctx, "", exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipeID, "", exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepIngredient(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepIngredient.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetTotalRecipeStepIngredientCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalRecipeStepIngredientsCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetTotalRecipeStepIngredientCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalRecipeStepIngredientsCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetTotalRecipeStepIngredientCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_getRecipeStepIngredientsForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		for i := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleRecipeStepIngredientList.RecipeStepIngredients[i].Ingredient = types.ValidIngredient{}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"recipe_step_ingredients.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, nil, false)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, exampleRecipeStepIngredientList.RecipeStepIngredients...))

		actual, err := c.getRecipeStepIngredientsForRecipe(ctx, exampleRecipeID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList.RecipeStepIngredients, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.getRecipeStepIngredientsForRecipe(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepIngredients(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		for i := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleRecipeStepIngredientList.RecipeStepIngredients[i].Ingredient = types.ValidIngredient{}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"recipe_step_ingredients.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, filter, true)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(true, exampleRecipeStepIngredientList.FilteredCount, exampleRecipeStepIngredientList.RecipeStepIngredients...))

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredients(ctx, "", exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()
		exampleRecipeStepIngredientList.Page = 0
		exampleRecipeStepIngredientList.Limit = 0
		for i := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleRecipeStepIngredientList.RecipeStepIngredients[i].Ingredient = types.ValidIngredient{}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"recipe_step_ingredients.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, filter, true)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(true, exampleRecipeStepIngredientList.FilteredCount, exampleRecipeStepIngredientList.RecipeStepIngredients...))

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"recipe_step_ingredients.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, filter, true)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"recipe_step_ingredients.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, filter, true)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepIngredientsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		var exampleIDs []string
		for i, x := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleIDs = append(exampleIDs, x.ID)
			exampleRecipeStepIngredientList.RecipeStepIngredients[i].Ingredient = types.ValidIngredient{}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetRecipeStepIngredientsWithIDsQuery(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, exampleRecipeStepIngredientList.RecipeStepIngredients...))

		actual, err := c.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipeStepID, 0, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList.RecipeStepIngredients, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid IDs", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipeStepID, defaultLimit, nil)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		var exampleIDs []string
		for _, x := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepIngredientsWithIDs(ctx, "", defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		var exampleIDs []string
		for _, x := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetRecipeStepIngredientsWithIDsQuery(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning query results", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		var exampleIDs []string
		for _, x := range exampleRecipeStepIngredientList.RecipeStepIngredients {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetRecipeStepIngredientsWithIDsQuery(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepIngredientsWithIDs(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.ID = "1"
		exampleRecipeStepIngredient.Ingredient = types.ValidIngredient{}
		exampleInput := fakes.BuildFakeRecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.IngredientID,
			exampleInput.QuantityType,
			exampleInput.QuantityValue,
			exampleInput.QuantityNotes,
			exampleInput.ProductOfRecipe,
			exampleInput.IngredientNotes,
			exampleInput.BelongsToRecipeStep,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() uint64 {
			return exampleRecipeStepIngredient.CreatedOn
		}

		actual, err := c.CreateRecipeStepIngredient(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredient, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepIngredient(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleInput := fakes.BuildFakeRecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.IngredientID,
			exampleInput.QuantityType,
			exampleInput.QuantityValue,
			exampleInput.QuantityNotes,
			exampleInput.ProductOfRecipe,
			exampleInput.IngredientNotes,
			exampleInput.BelongsToRecipeStep,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() uint64 {
			return exampleRecipeStepIngredient.CreatedOn
		}

		actual, err := c.CreateRecipeStepIngredient(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_createRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.ID = "3"
		exampleRecipeStepIngredient.BelongsToRecipeStep = "2"
		exampleRecipeStepIngredient.Ingredient = types.ValidIngredient{}

		exampleInput := fakes.BuildFakeRecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		recipeStepIngredientCreationArgs := []interface{}{
			exampleInput.ID,
			exampleInput.IngredientID,
			exampleInput.QuantityType,
			exampleInput.QuantityValue,
			exampleInput.QuantityNotes,
			exampleInput.ProductOfRecipe,
			exampleInput.IngredientNotes,
			exampleInput.BelongsToRecipeStep,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeStepIngredientCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() uint64 {
			return exampleRecipeStepIngredient.CreatedOn
		}

		actual, err := c.createRecipeStepIngredient(ctx, c.db, exampleInput)
		assert.NoError(t, err)
		assert.NotNil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepIngredient.IngredientID,
			exampleRecipeStepIngredient.QuantityType,
			exampleRecipeStepIngredient.QuantityValue,
			exampleRecipeStepIngredient.QuantityNotes,
			exampleRecipeStepIngredient.ProductOfRecipeStep,
			exampleRecipeStepIngredient.IngredientNotes,
			exampleRecipeStepIngredient.BelongsToRecipeStep,
			exampleRecipeStepIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeStepIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateRecipeStepIngredient(ctx, exampleRecipeStepIngredient))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepIngredient(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepIngredient.IngredientID,
			exampleRecipeStepIngredient.QuantityType,
			exampleRecipeStepIngredient.QuantityValue,
			exampleRecipeStepIngredient.QuantityNotes,
			exampleRecipeStepIngredient.ProductOfRecipeStep,
			exampleRecipeStepIngredient.IngredientNotes,
			exampleRecipeStepIngredient.BelongsToRecipeStep,
			exampleRecipeStepIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeStepIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateRecipeStepIngredient(ctx, exampleRecipeStepIngredient))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeStepIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepID, exampleRecipeStepIngredient.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, "", exampleRecipeStepIngredient.ID))
	})

	T.Run("with invalid recipe step ingredient ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepID, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepIngredient.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeStepIngredientQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveRecipeStepIngredient(ctx, exampleRecipeStepID, exampleRecipeStepIngredient.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
