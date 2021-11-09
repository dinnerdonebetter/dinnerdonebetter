package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	database "github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromRecipes(includeCounts bool, filteredCount uint64, recipes ...*types.Recipe) *sqlmock.Rows {
	columns := recipesTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

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
			x.CreatedByUser,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(recipes))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildMockFullRowsFromRecipe(recipe *types.Recipe) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(completeRecipeColumns)

	for _, step := range recipe.Steps {
		for _, ingredient := range step.Ingredients {
			exampleRows.AddRow(
				&recipe.ID,
				&recipe.Name,
				&recipe.Source,
				&recipe.Description,
				&recipe.InspiredByRecipeID,
				&recipe.CreatedOn,
				&recipe.LastUpdatedOn,
				&recipe.ArchivedOn,
				&recipe.CreatedByUser,
				&step.ID,
				&step.Index,
				&step.Preparation.ID,
				&step.Preparation.Name,
				&step.Preparation.Description,
				&step.Preparation.IconPath,
				&step.Preparation.CreatedOn,
				&step.Preparation.LastUpdatedOn,
				&step.Preparation.ArchivedOn,
				&step.PrerequisiteStep,
				&step.MinEstimatedTimeInSeconds,
				&step.MaxEstimatedTimeInSeconds,
				&step.TemperatureInCelsius,
				&step.Notes,
				&step.Why,
				&step.CreatedOn,
				&step.LastUpdatedOn,
				&step.ArchivedOn,
				&step.BelongsToRecipe,
				&ingredient.ID,
				&ingredient.Ingredient.ID,
				&ingredient.Ingredient.Name,
				&ingredient.Ingredient.Variant,
				&ingredient.Ingredient.Description,
				&ingredient.Ingredient.Warning,
				&ingredient.Ingredient.ContainsEgg,
				&ingredient.Ingredient.ContainsDairy,
				&ingredient.Ingredient.ContainsPeanut,
				&ingredient.Ingredient.ContainsTreeNut,
				&ingredient.Ingredient.ContainsSoy,
				&ingredient.Ingredient.ContainsWheat,
				&ingredient.Ingredient.ContainsShellfish,
				&ingredient.Ingredient.ContainsSesame,
				&ingredient.Ingredient.ContainsFish,
				&ingredient.Ingredient.ContainsGluten,
				&ingredient.Ingredient.AnimalFlesh,
				&ingredient.Ingredient.AnimalDerived,
				&ingredient.Ingredient.Volumetric,
				&ingredient.Ingredient.IconPath,
				&ingredient.Ingredient.CreatedOn,
				&ingredient.Ingredient.LastUpdatedOn,
				&ingredient.Ingredient.ArchivedOn,
				&ingredient.QuantityType,
				&ingredient.QuantityValue,
				&ingredient.QuantityNotes,
				&ingredient.ProductOfRecipeStep,
				&ingredient.IngredientNotes,
				&ingredient.CreatedOn,
				&ingredient.LastUpdatedOn,
				&ingredient.ArchivedOn,
				&ingredient.BelongsToRecipeStep,
			)
		}
	}

	return exampleRows
}

func buildInvalidMockFullRowsFromRecipe(recipe *types.Recipe) *sqlmock.Rows {
	columns := completeRecipeColumns
	exampleRows := sqlmock.NewRows(columns)

	for _, step := range recipe.Steps {
		for range step.Ingredients {
			exampleRows.AddRow(
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
				driver.Value(nil),
			)
		}
	}

	return exampleRows
}

func TestQuerier_ScanRecipes(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipes(ctx, mockRows, false)
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

		_, _, _, err := q.scanRecipes(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_RecipeExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipe := fakes.BuildFakeRecipe()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.RecipeExists(ctx, exampleRecipe.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipe := fakes.BuildFakeRecipe()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.RecipeExists(ctx, exampleRecipe.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipe := fakes.BuildFakeRecipe()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.RecipeExists(ctx, exampleRecipe.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		exampleRecipe.Steps = []*types.RecipeStep{
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
		}

		for i, step := range exampleRecipe.Steps {
			exampleRecipe.Steps[i].Ingredients = []*types.RecipeStepIngredient{}
			for j := 0; j < 3; j++ {
				ingredient := fakes.BuildFakeRecipeStepIngredient()
				ingredient.IngredientID = nil

				exampleRecipe.Steps[i].Ingredients = append(step.Ingredients, ingredient)
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getCompleteRecipeByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockFullRowsFromRecipe(exampleRecipe))

		actual, err := c.GetRecipe(ctx, exampleRecipe.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipe, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipe(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getCompleteRecipeByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipe(ctx, exampleRecipe.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response from database", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		exampleRecipe.Steps = []*types.RecipeStep{
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
		}

		for _, step := range exampleRecipe.Steps {
			step.Ingredients = []*types.RecipeStepIngredient{
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getCompleteRecipeByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildInvalidMockFullRowsFromRecipe(exampleRecipe))

		actual, err := c.GetRecipe(ctx, exampleRecipe.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with no results returned", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		exampleRecipe.Steps = []*types.RecipeStep{
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
		}

		for _, step := range exampleRecipe.Steps {
			step.Ingredients = []*types.RecipeStepIngredient{
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getCompleteRecipeByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"things"}))

		actual, err := c.GetRecipe(ctx, exampleRecipe.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.True(t, errors.Is(err, sql.ErrNoRows))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeByUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		exampleRecipe.Steps = []*types.RecipeStep{
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
		}

		for i, step := range exampleRecipe.Steps {
			exampleRecipe.Steps[i].Ingredients = []*types.RecipeStepIngredient{}
			for j := 0; j < 3; j++ {
				ingredient := fakes.BuildFakeRecipeStepIngredient()
				ingredient.IngredientID = nil

				exampleRecipe.Steps[i].Ingredients = append(step.Ingredients, ingredient)
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectQuery(formatQueryForSQLMock(getCompleteRecipeByIDAndAuthorIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockFullRowsFromRecipe(exampleRecipe))

		actual, err := c.GetRecipeByIDAndUser(ctx, exampleRecipe.ID, exampleRecipe.CreatedByUser)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipe, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)
		exampleUserID := fakes.BuildFakeID()

		actual, err := c.GetRecipeByIDAndUser(ctx, "", exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)
		exampleRecipeID := fakes.BuildFakeID()

		actual, err := c.GetRecipeByIDAndUser(ctx, exampleRecipeID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectQuery(formatQueryForSQLMock(getCompleteRecipeByIDAndAuthorIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeByIDAndUser(ctx, exampleRecipe.ID, exampleRecipe.CreatedByUser)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response from database", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		exampleRecipe.Steps = []*types.RecipeStep{
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
		}

		for _, step := range exampleRecipe.Steps {
			step.Ingredients = []*types.RecipeStepIngredient{
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectQuery(formatQueryForSQLMock(getCompleteRecipeByIDAndAuthorIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildInvalidMockFullRowsFromRecipe(exampleRecipe))

		actual, err := c.GetRecipeByIDAndUser(ctx, exampleRecipe.ID, exampleRecipe.CreatedByUser)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with no results returned", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		exampleRecipe.Steps = []*types.RecipeStep{
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
		}

		for _, step := range exampleRecipe.Steps {
			step.Ingredients = []*types.RecipeStepIngredient{
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectQuery(formatQueryForSQLMock(getCompleteRecipeByIDAndAuthorIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"things"}))

		actual, err := c.GetRecipeByIDAndUser(ctx, exampleRecipe.ID, exampleRecipe.CreatedByUser)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.True(t, errors.Is(err, sql.ErrNoRows))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetTotalRecipeCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalRecipesCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetTotalRecipeCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalRecipesCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetTotalRecipeCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeList := fakes.BuildFakeRecipeList()
		for i := range exampleRecipeList.Recipes {
			exampleRecipeList.Recipes[i].Steps = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipes", nil, nil, nil, householdOwnershipColumn, recipesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipes(true, exampleRecipeList.FilteredCount, exampleRecipeList.Recipes...))

		actual, err := c.GetRecipes(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleRecipeList := fakes.BuildFakeRecipeList()
		exampleRecipeList.Page = 0
		exampleRecipeList.Limit = 0
		for i := range exampleRecipeList.Recipes {
			exampleRecipeList.Recipes[i].Steps = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipes", nil, nil, nil, householdOwnershipColumn, recipesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipes(true, exampleRecipeList.FilteredCount, exampleRecipeList.Recipes...))

		actual, err := c.GetRecipes(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipes", nil, nil, nil, householdOwnershipColumn, recipesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipes", nil, nil, nil, householdOwnershipColumn, recipesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipes(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipesWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleRecipeList := fakes.BuildFakeRecipeList()
		for i := range exampleRecipeList.Recipes {
			exampleRecipeList.Recipes[i].Steps = nil
		}

		var exampleIDs []string
		for _, x := range exampleRecipeList.Recipes {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetRecipesWithIDsQuery(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipes(false, 0, exampleRecipeList.Recipes...))

		actual, err := c.GetRecipesWithIDs(ctx, exampleHouseholdID, 0, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeList.Recipes, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid IDs", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipesWithIDs(ctx, exampleHouseholdID, defaultLimit, nil)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleRecipeList := fakes.BuildFakeRecipeList()

		var exampleIDs []string
		for _, x := range exampleRecipeList.Recipes {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetRecipesWithIDsQuery(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipesWithIDs(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning query results", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleRecipeList := fakes.BuildFakeRecipeList()

		var exampleIDs []string
		for _, x := range exampleRecipeList.Recipes {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetRecipesWithIDsQuery(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipesWithIDs(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.ID = "1"
		for i, step := range exampleRecipe.Steps {
			exampleRecipe.Steps[i].ID = "2"
			exampleRecipe.Steps[i].BelongsToRecipe = "1"
			exampleRecipe.Steps[i].Preparation = types.ValidPreparation{}
			for j := range step.Ingredients {
				exampleRecipe.Steps[i].Ingredients[j].ID = "3"
				exampleRecipe.Steps[i].Ingredients[j].BelongsToRecipeStep = "2"
				exampleRecipe.Steps[i].Ingredients[j].Ingredient = types.ValidIngredient{}
			}
		}

		exampleInput := fakes.BuildFakeRecipeDatabaseCreationInputFromRecipe(exampleRecipe)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		recipeCreationArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.Name,
			exampleRecipe.Source,
			exampleRecipe.Description,
			exampleRecipe.InspiredByRecipeID,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleRecipe.ID))

		for _, step := range exampleInput.Steps {
			recipeStepCreationArgs := []interface{}{
				step.ID,
				step.Index,
				step.PreparationID,
				step.PrerequisiteStep,
				step.MinEstimatedTimeInSeconds,
				step.MaxEstimatedTimeInSeconds,
				step.TemperatureInCelsius,
				step.Notes,
				step.Why,
				step.BelongsToRecipe,
			}

			db.ExpectExec(formatQueryForSQLMock(recipeStepCreationQuery)).
				WithArgs(interfaceToDriverValue(recipeStepCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult(step.ID))

			for _, ingredient := range step.Ingredients {
				recipeStepIngredientCreationArgs := []interface{}{
					ingredient.ID,
					ingredient.IngredientID,
					ingredient.QuantityType,
					ingredient.QuantityValue,
					ingredient.QuantityNotes,
					ingredient.ProductOfRecipe,
					ingredient.IngredientNotes,
					ingredient.BelongsToRecipeStep,
				}

				db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
					WithArgs(interfaceToDriverValue(recipeStepIngredientCreationArgs)...).
					WillReturnResult(newArbitraryDatabaseResult(ingredient.ID))
			}
		}

		db.ExpectCommit()

		c.timeFunc = func() uint64 {
			return exampleRecipe.CreatedOn
		}

		actual, err := c.CreateRecipe(ctx, exampleInput)
		assert.NoError(t, err)
		require.Equal(t, len(exampleRecipe.Steps), len(actual.Steps))

		for i, step := range exampleRecipe.Steps {
			step.BelongsToRecipe = actual.ID
			step.CreatedOn = actual.Steps[i].CreatedOn

			for j, ingredient := range step.Ingredients {
				ingredient.BelongsToRecipeStep = step.ID
				ingredient.CreatedOn = actual.Steps[i].Ingredients[j].CreatedOn
			}
		}

		assert.Equal(t, exampleRecipe, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.ID = "1"
		exampleInput := fakes.BuildFakeRecipeDatabaseCreationInputFromRecipe(exampleRecipe)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		c.timeFunc = func() uint64 {
			return exampleRecipe.CreatedOn
		}

		actual, err := c.CreateRecipe(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipe(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.ID = "1"
		for i, step := range exampleRecipe.Steps {
			exampleRecipe.Steps[i].ID = "2"
			exampleRecipe.Steps[i].BelongsToRecipe = "1"
			for j := range step.Ingredients {
				exampleRecipe.Steps[i].Ingredients[j].ID = "3"
				exampleRecipe.Steps[i].Ingredients[j].BelongsToRecipeStep = "2"
			}
		}
		exampleInput := fakes.BuildFakeRecipeDatabaseCreationInputFromRecipe(exampleRecipe)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Source,
			exampleInput.Description,
			exampleInput.InspiredByRecipeID,
			exampleInput.CreatedByUser,
		}

		db.ExpectBegin()

		db.ExpectExec(formatQueryForSQLMock(recipeCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleRecipe.CreatedOn
		}

		actual, err := c.CreateRecipe(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.ID = "1"
		exampleRecipe.Steps = nil
		exampleInput := fakes.BuildFakeRecipeDatabaseCreationInputFromRecipe(exampleRecipe)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Source,
			exampleInput.Description,
			exampleInput.InspiredByRecipeID,
			exampleInput.CreatedByUser,
		}

		db.ExpectBegin()

		db.ExpectExec(formatQueryForSQLMock(recipeCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleRecipe.ID))

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() uint64 {
			return exampleRecipe.CreatedOn
		}

		actual, err := c.CreateRecipe(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.Name,
			exampleRecipe.Source,
			exampleRecipe.Description,
			exampleRecipe.InspiredByRecipeID,
			exampleRecipe.CreatedByUser,
			exampleRecipe.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleRecipe.ID))

		assert.NoError(t, c.UpdateRecipe(ctx, exampleRecipe))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipe(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.Name,
			exampleRecipe.Source,
			exampleRecipe.Description,
			exampleRecipe.InspiredByRecipeID,
			exampleRecipe.CreatedByUser,
			exampleRecipe.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateRecipe(ctx, exampleRecipe))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHouseholdID,
			exampleRecipe.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleRecipe.ID))

		assert.NoError(t, c.ArchiveRecipe(ctx, exampleRecipe.ID, exampleHouseholdID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipe(ctx, "", exampleHouseholdID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipe(ctx, exampleRecipe.ID, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHouseholdID,
			exampleRecipe.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveRecipe(ctx, exampleRecipe.ID, exampleHouseholdID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
