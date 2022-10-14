package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
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
			x.YieldsPortions,
			x.SealOfApproval,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
			x.CreatedByUser,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(recipes))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

// fullRecipesColumns are the columns for the recipes table.
var fullRecipesColumns = []string{
	"recipes.id",
	"recipes.name",
	"recipes.source",
	"recipes.description",
	"recipes.inspired_by_recipe_id",
	"recipes.yields_portions",
	"recipes.seal_of_approval",
	"recipes.created_at",
	"recipes.last_updated_at",
	"recipes.archived_at",
	"recipes.created_by_user",
	"recipe_steps.id",
	"recipe_steps.index",
	"valid_preparations.id",
	"valid_preparations.name",
	"valid_preparations.description",
	"valid_preparations.icon_path",
	"valid_preparations.yields_nothing",
	"valid_preparations.restrict_to_ingredients",
	"valid_preparations.zero_ingredients_allowable",
	"valid_preparations.past_tense",
	"valid_preparations.created_at",
	"valid_preparations.last_updated_at",
	"valid_preparations.archived_at",
	"recipe_steps.minimum_estimated_time_in_seconds",
	"recipe_steps.maximum_estimated_time_in_seconds",
	"recipe_steps.minimum_temperature_in_celsius",
	"recipe_steps.maximum_temperature_in_celsius",
	"recipe_steps.notes",
	"recipe_steps.explicit_instructions",
	"recipe_steps.optional",
	"recipe_steps.created_at",
	"recipe_steps.last_updated_at",
	"recipe_steps.archived_at",
	"recipe_steps.belongs_to_recipe",
}

func buildMockFullRowsFromRecipe(recipe *types.Recipe) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(fullRecipesColumns)

	for _, step := range recipe.Steps {
		exampleRows.AddRow(
			&recipe.ID,
			&recipe.Name,
			&recipe.Source,
			&recipe.Description,
			&recipe.InspiredByRecipeID,
			&recipe.YieldsPortions,
			&recipe.SealOfApproval,
			&recipe.CreatedAt,
			&recipe.LastUpdatedAt,
			&recipe.ArchivedAt,
			&recipe.CreatedByUser,
			&step.ID,
			&step.Index,
			&step.Preparation.ID,
			&step.Preparation.Name,
			&step.Preparation.Description,
			&step.Preparation.IconPath,
			&step.Preparation.YieldsNothing,
			&step.Preparation.RestrictToIngredients,
			&step.Preparation.ZeroIngredientsAllowable,
			&step.Preparation.PastTense,
			&step.Preparation.CreatedAt,
			&step.Preparation.LastUpdatedAt,
			&step.Preparation.ArchivedAt,
			&step.MinimumEstimatedTimeInSeconds,
			&step.MaximumEstimatedTimeInSeconds,
			&step.MinimumTemperatureInCelsius,
			&step.MaximumTemperatureInCelsius,
			&step.Notes,
			&step.ExplicitInstructions,
			&step.Optional,
			&step.CreatedAt,
			&step.LastUpdatedAt,
			&step.ArchivedAt,
			&step.BelongsToRecipe,
		)
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

func prepareMockToSuccessfullyGetRecipe(ctx context.Context, t *testing.T, recipe *types.Recipe, userID string, c *Querier, db *sqlmockExpecterWrapper) {
	t.Helper()

	var (
		exampleRecipe *types.Recipe
	)

	if recipe == nil {
		exampleRecipe = fakes.BuildFakeRecipe()
		exampleRecipe.Steps = []*types.RecipeStep{
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
		}
	} else {
		exampleRecipe = recipe
	}

	allIngredients := []*types.RecipeStepIngredient{}
	allInstruments := []*types.RecipeStepInstrument{}
	allProducts := []*types.RecipeStepProduct{}
	for _, step := range exampleRecipe.Steps {
		allIngredients = append(allIngredients, step.Ingredients...)
		allInstruments = append(allInstruments, step.Instruments...)
		allProducts = append(allProducts, step.Products...)
	}

	args := []interface{}{
		exampleRecipe.ID,
	}

	query := getRecipeByIDQuery
	if userID != "" {
		query = getRecipeByIDAndAuthorIDQuery
		args = append(args, userID)
	}

	db.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(interfaceToDriverValue(args)...).
		WillReturnRows(buildMockFullRowsFromRecipe(exampleRecipe))

	listRecipePrepTasksForRecipeArgs := []interface{}{
		exampleRecipe.ID,
	}

	db.ExpectQuery(formatQueryForSQLMock(listRecipePrepTasksForRecipeQuery)).
		WithArgs(interfaceToDriverValue(listRecipePrepTasksForRecipeArgs)...).
		WillReturnRows(buildMockRowsFromRecipePrepTasks(exampleRecipe.PrepTasks...))

	getRecipeStepIngredientsForRecipeArgs := []interface{}{
		exampleRecipe.ID,
	}
	db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientsForRecipeQuery)).
		WithArgs(interfaceToDriverValue(getRecipeStepIngredientsForRecipeArgs)...).
		WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, allIngredients...))

	productsArgs := []interface{}{
		exampleRecipe.ID,
	}
	db.ExpectQuery(formatQueryForSQLMock(getRecipeStepProductsForRecipeQuery)).
		WithArgs(interfaceToDriverValue(productsArgs)...).
		WillReturnRows(buildMockRowsFromRecipeStepProducts(false, 0, allProducts...))

	instrumentsArgs := []interface{}{
		exampleRecipe.ID,
	}
	db.ExpectQuery(formatQueryForSQLMock(getRecipeStepInstrumentsForRecipeQuery)).
		WithArgs(interfaceToDriverValue(instrumentsArgs)...).
		WillReturnRows(buildMockRowsFromRecipeStepInstruments(false, 0, allInstruments...))
}

func TestQuerier_getRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleUserID := fakes.BuildFakeID()
		exampleRecipe.Steps = []*types.RecipeStep{
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
		}

		prepareMockToSuccessfullyGetRecipe(ctx, t, exampleRecipe, exampleUserID, c, db)

		actual, err := c.getRecipe(ctx, exampleRecipe.ID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipe, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error fetching recipe step ingredients", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleUserID := fakes.BuildFakeID()
		exampleRecipe.Steps = []*types.RecipeStep{
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDAndAuthorIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockFullRowsFromRecipe(exampleRecipe))

		listRecipePrepTasksForRecipeArgs := []interface{}{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(listRecipePrepTasksForRecipeQuery)).
			WithArgs(interfaceToDriverValue(listRecipePrepTasksForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipePrepTasks(exampleRecipe.PrepTasks...))

		getRecipeStepIngredientsForRecipeArgs := []interface{}{
			exampleRecipe.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(getRecipeStepIngredientsForRecipeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.getRecipe(ctx, exampleRecipe.ID, exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error retrieving recipe step products", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleUserID := fakes.BuildFakeID()
		exampleRecipe.Steps = []*types.RecipeStep{
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
			fakes.BuildFakeRecipeStep(),
		}

		allIngredients := []*types.RecipeStepIngredient{}
		for _, step := range exampleRecipe.Steps {
			allIngredients = append(allIngredients, step.Ingredients...)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDAndAuthorIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockFullRowsFromRecipe(exampleRecipe))

		listRecipePrepTasksForRecipeArgs := []interface{}{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(listRecipePrepTasksForRecipeQuery)).
			WithArgs(interfaceToDriverValue(listRecipePrepTasksForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipePrepTasks(exampleRecipe.PrepTasks...))

		getRecipeStepIngredientsForRecipeArgs := []interface{}{
			exampleRecipe.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(getRecipeStepIngredientsForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, allIngredients...))

		productsArgs := []interface{}{
			exampleRecipe.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepProductsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(productsArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.getRecipe(ctx, exampleRecipe.ID, exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

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

		allIngredients := []*types.RecipeStepIngredient{}
		allInstruments := []*types.RecipeStepInstrument{}
		allProducts := []*types.RecipeStepProduct{}
		for _, step := range exampleRecipe.Steps {
			allIngredients = append(allIngredients, step.Ingredients...)
			allInstruments = append(allInstruments, step.Instruments...)
			allProducts = append(allProducts, step.Products...)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockFullRowsFromRecipe(exampleRecipe))

		listRecipePrepTasksForRecipeArgs := []interface{}{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(listRecipePrepTasksForRecipeQuery)).
			WithArgs(interfaceToDriverValue(listRecipePrepTasksForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipePrepTasks(exampleRecipe.PrepTasks...))

		getRecipeStepIngredientsForRecipeArgs := []interface{}{
			exampleRecipe.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(getRecipeStepIngredientsForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, allIngredients...))

		productsArgs := []interface{}{
			exampleRecipe.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepProductsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(productsArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepProducts(false, 0, allProducts...))

		instrumentsArgs := []interface{}{
			exampleRecipe.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepInstrumentsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(instrumentsArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepInstruments(false, 0, allInstruments...))

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
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

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
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

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

		allIngredients := []*types.RecipeStepIngredient{}
		allInstruments := []*types.RecipeStepInstrument{}
		allProducts := []*types.RecipeStepProduct{}
		for _, step := range exampleRecipe.Steps {
			allIngredients = append(allIngredients, step.Ingredients...)
			allInstruments = append(allInstruments, step.Instruments...)
			allProducts = append(allProducts, step.Products...)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDAndAuthorIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockFullRowsFromRecipe(exampleRecipe))

		listRecipePrepTasksForRecipeArgs := []interface{}{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(listRecipePrepTasksForRecipeQuery)).
			WithArgs(interfaceToDriverValue(listRecipePrepTasksForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipePrepTasks(exampleRecipe.PrepTasks...))

		getRecipeStepIngredientsForRecipeArgs := []interface{}{
			exampleRecipe.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(getRecipeStepIngredientsForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, allIngredients...))

		productsArgs := []interface{}{
			exampleRecipe.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepProductsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(productsArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepProducts(false, 0, allProducts...))

		instrumentsArgs := []interface{}{
			exampleRecipe.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepInstrumentsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(instrumentsArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepInstruments(false, 0, allInstruments...))

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
			exampleRecipe.CreatedByUser,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDAndAuthorIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

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
			exampleRecipe.CreatedByUser,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDAndAuthorIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.GetRecipeByIDAndUser(ctx, exampleRecipe.ID, exampleRecipe.CreatedByUser)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.True(t, errors.Is(err, sql.ErrNoRows))

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
			exampleRecipeList.Recipes[i].PrepTasks = nil
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
			exampleRecipeList.Recipes[i].PrepTasks = nil
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

func TestQuerier_getRecipeIDsForMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()
		exampleRecipeList := fakes.BuildFakeRecipeList()
		exampleRecipeIDs := []string{}
		for i := range exampleRecipeList.Recipes {
			exampleRecipeList.Recipes[i].Steps = nil
			exampleRecipeList.Recipes[i].PrepTasks = nil
			exampleRecipeIDs = append(exampleRecipeIDs, exampleRecipeList.Recipes[i].ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMeal.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipesForMealQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromIDs(exampleRecipeIDs...))

		actual, err := c.getRecipeIDsForMeal(ctx, exampleMeal.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIDs, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForRecipes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeList := fakes.BuildFakeRecipeList()
		for i := range exampleRecipeList.Recipes {
			exampleRecipeList.Recipes[i].Steps = nil
			exampleRecipeList.Recipes[i].PrepTasks = nil
		}

		ctx := context.Background()
		recipeNameQuery := "example"
		c, db := buildTestClient(t)

		where := squirrel.ILike{"name": wrapQueryForILIKE(recipeNameQuery)}
		query, args := c.buildListQueryWithILike(ctx, "recipes", nil, nil, where, "", recipesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipes(true, exampleRecipeList.FilteredCount, exampleRecipeList.Recipes...))

		actual, err := c.SearchForRecipes(ctx, recipeNameQuery, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeList := fakes.BuildFakeRecipeList()
		for i := range exampleRecipeList.Recipes {
			exampleRecipeList.Recipes[i].Steps = nil
			exampleRecipeList.Recipes[i].PrepTasks = nil
		}

		ctx := context.Background()
		recipeNameQuery := "example"
		c, db := buildTestClient(t)

		where := squirrel.ILike{"name": wrapQueryForILIKE(recipeNameQuery)}
		query, args := c.buildListQueryWithILike(ctx, "recipes", nil, nil, where, "", recipesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForRecipes(ctx, recipeNameQuery, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeList := fakes.BuildFakeRecipeList()
		for i := range exampleRecipeList.Recipes {
			exampleRecipeList.Recipes[i].Steps = nil
			exampleRecipeList.Recipes[i].PrepTasks = nil
		}

		ctx := context.Background()
		recipeNameQuery := "example"
		c, db := buildTestClient(t)

		where := squirrel.ILike{"name": wrapQueryForILIKE(recipeNameQuery)}
		query, args := c.buildListQueryWithILike(ctx, "recipes", nil, nil, where, "", recipesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForRecipes(ctx, recipeNameQuery, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

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
			exampleRecipe.Steps[i].Preparation = types.ValidPreparation{ID: exampleRecipe.Steps[i].Preparation.ID}

			for j := range step.Ingredients {
				exampleRecipe.Steps[i].Ingredients[j].ID = "3"
				exampleRecipe.Steps[i].Ingredients[j].Ingredient = &types.ValidIngredient{ID: exampleRecipe.Steps[i].Ingredients[j].Ingredient.ID}
				exampleRecipe.Steps[i].Ingredients[j].BelongsToRecipeStep = "2"
				exampleRecipe.Steps[i].Ingredients[j].MeasurementUnit = types.ValidMeasurementUnit{ID: exampleRecipe.Steps[i].Ingredients[j].MeasurementUnit.ID}
			}

			for j := range step.Instruments {
				exampleRecipe.Steps[i].Instruments[j].ID = "3"
				exampleRecipe.Steps[i].Instruments[j].Instrument = &types.ValidInstrument{ID: exampleRecipe.Steps[i].Instruments[j].Instrument.ID}
				exampleRecipe.Steps[i].Instruments[j].BelongsToRecipeStep = "2"
			}

			step.Products = nil
		}

		for i := range exampleRecipe.PrepTasks {
			exampleRecipe.PrepTasks[i].BelongsToRecipe = exampleRecipe.ID
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
			exampleRecipe.YieldsPortions,
			exampleRecipe.SealOfApproval,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, step := range exampleInput.Steps {
			recipeStepCreationArgs := []interface{}{
				step.ID,
				step.Index,
				step.PreparationID,
				step.MinimumEstimatedTimeInSeconds,
				step.MaximumEstimatedTimeInSeconds,
				step.MinimumTemperatureInCelsius,
				step.MaximumTemperatureInCelsius,
				step.Notes,
				step.ExplicitInstructions,
				step.Optional,
				step.BelongsToRecipe,
			}

			db.ExpectExec(formatQueryForSQLMock(recipeStepCreationQuery)).
				WithArgs(interfaceToDriverValue(recipeStepCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())

			for _, ingredient := range step.Ingredients {
				recipeStepIngredientCreationArgs := []interface{}{
					ingredient.ID,
					ingredient.Name,
					ingredient.Optional,
					ingredient.IngredientID,
					ingredient.MeasurementUnitID,
					ingredient.MinimumQuantity,
					ingredient.MaximumQuantity,
					ingredient.QuantityNotes,
					ingredient.ProductOfRecipeStep,
					ingredient.RecipeStepProductID,
					ingredient.IngredientNotes,
					ingredient.BelongsToRecipeStep,
				}

				db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
					WithArgs(interfaceToDriverValue(recipeStepIngredientCreationArgs)...).
					WillReturnResult(newArbitraryDatabaseResult())
			}

			for _, instrument := range step.Instruments {
				recipeStepInstrumentCreationArgs := []interface{}{
					instrument.ID,
					instrument.InstrumentID,
					instrument.RecipeStepProductID,
					instrument.Name,
					instrument.ProductOfRecipeStep,
					instrument.Notes,
					instrument.PreferenceRank,
					instrument.Optional,
					instrument.MinimumQuantity,
					instrument.MaximumQuantity,
					instrument.BelongsToRecipeStep,
				}

				db.ExpectExec(formatQueryForSQLMock(recipeStepInstrumentCreationQuery)).
					WithArgs(interfaceToDriverValue(recipeStepInstrumentCreationArgs)...).
					WillReturnResult(newArbitraryDatabaseResult())
			}
		}

		for _, prepTask := range exampleInput.PrepTasks {
			createRecipePrepTaskQueryArgs := []interface{}{
				prepTask.ID,
				prepTask.Notes,
				prepTask.ExplicitStorageInstructions,
				prepTask.MinimumTimeBufferBeforeRecipeInSeconds,
				prepTask.MaximumTimeBufferBeforeRecipeInSeconds,
				prepTask.StorageType,
				prepTask.MinimumStorageTemperatureInCelsius,
				prepTask.MaximumStorageTemperatureInCelsius,
				prepTask.BelongsToRecipe,
			}

			db.ExpectExec(formatQueryForSQLMock(createRecipePrepTaskQuery)).
				WithArgs(interfaceToDriverValue(createRecipePrepTaskQueryArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())

			for _, taskStep := range prepTask.TaskSteps {
				createRecipePrepTaskStepArgs := []interface{}{
					taskStep.ID,
					taskStep.BelongsToRecipePrepTask,
					taskStep.BelongsToRecipeStep,
					taskStep.SatisfiesRecipeStep,
				}

				db.ExpectExec(formatQueryForSQLMock(createRecipePrepTaskStepQuery)).
					WithArgs(interfaceToDriverValue(createRecipePrepTaskStepArgs)...).
					WillReturnResult(newArbitraryDatabaseResult())
			}
		}

		db.ExpectCommit()

		c.timeFunc = func() time.Time {
			return exampleRecipe.CreatedAt
		}

		actual, err := c.CreateRecipe(ctx, exampleInput)
		require.NoError(t, err)
		require.NotNil(t, actual)
		require.Equal(t, len(exampleRecipe.Steps), len(actual.Steps))

		for i, step := range exampleRecipe.Steps {
			step.BelongsToRecipe = actual.ID
			step.CreatedAt = actual.Steps[i].CreatedAt

			for j, ingredient := range step.Ingredients {
				ingredient.BelongsToRecipeStep = step.ID
				ingredient.CreatedAt = actual.Steps[i].Ingredients[j].CreatedAt
			}

			for j, instrument := range step.Instruments {
				instrument.BelongsToRecipeStep = step.ID
				instrument.CreatedAt = actual.Steps[i].Instruments[j].CreatedAt
			}
		}

		for i, prepTask := range exampleRecipe.PrepTasks {
			prepTask.BelongsToRecipe = actual.ID
			prepTask.CreatedAt = actual.PrepTasks[i].CreatedAt
		}

		assert.Equal(t, exampleRecipe, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("while also creating meal", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.Steps = nil
		exampleRecipe.PrepTasks = nil
		exampleRecipe.ID = "1"

		exampleInput := fakes.BuildFakeRecipeDatabaseCreationInputFromRecipe(exampleRecipe)
		exampleInput.AlsoCreateMeal = true
		exampleInput.Steps = []*types.RecipeStepDatabaseCreationInput{}

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		recipeCreationArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.Name,
			exampleRecipe.Source,
			exampleRecipe.Description,
			exampleRecipe.InspiredByRecipeID,
			exampleRecipe.YieldsPortions,
			exampleRecipe.SealOfApproval,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		mealCreationArgs := []interface{}{
			&idMatcher{},
			exampleRecipe.Name,
			exampleRecipe.Description,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealCreationQuery)).
			WithArgs(interfaceToDriverValue(mealCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		mealRecipeCreationArgs := []interface{}{
			&idMatcher{},
			&idMatcher{},
			exampleRecipe.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
			WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

		c.timeFunc = func() time.Time {
			return exampleRecipe.CreatedAt
		}

		actual, err := c.CreateRecipe(ctx, exampleInput)
		require.NotNil(t, actual)
		require.NoError(t, err)
		require.Equal(t, len(exampleRecipe.Steps), len(actual.Steps))

		for i, step := range exampleRecipe.Steps {
			step.BelongsToRecipe = actual.ID
			step.CreatedAt = actual.Steps[i].CreatedAt

			for j, ingredient := range step.Ingredients {
				ingredient.BelongsToRecipeStep = step.ID
				ingredient.CreatedAt = actual.Steps[i].Ingredients[j].CreatedAt
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

		c.timeFunc = func() time.Time {
			return exampleRecipe.CreatedAt
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
			exampleRecipe.YieldsPortions,
			exampleRecipe.SealOfApproval,
			exampleInput.CreatedByUser,
		}

		db.ExpectBegin()

		db.ExpectExec(formatQueryForSQLMock(recipeCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		db.ExpectRollback()

		c.timeFunc = func() time.Time {
			return exampleRecipe.CreatedAt
		}

		actual, err := c.CreateRecipe(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error creating recipe step", func(t *testing.T) {
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
			exampleRecipe.YieldsPortions,
			exampleRecipe.SealOfApproval,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		recipeStepCreationArgs := []interface{}{
			exampleInput.Steps[0].ID,
			0,
			exampleInput.Steps[0].PreparationID,
			exampleInput.Steps[0].MinimumEstimatedTimeInSeconds,
			exampleInput.Steps[0].MaximumEstimatedTimeInSeconds,
			exampleInput.Steps[0].MinimumTemperatureInCelsius,
			exampleInput.Steps[0].MaximumTemperatureInCelsius,
			exampleInput.Steps[0].Notes,
			exampleInput.Steps[0].ExplicitInstructions,
			exampleInput.Steps[0].Optional,
			exampleInput.Steps[0].BelongsToRecipe,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeStepCreationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() time.Time {
			return exampleRecipe.CreatedAt
		}

		actual, err := c.CreateRecipe(ctx, exampleInput)
		assert.Error(t, err)
		require.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.ID = "1"
		exampleRecipe.Steps = nil
		exampleRecipe.PrepTasks = nil
		exampleInput := fakes.BuildFakeRecipeDatabaseCreationInputFromRecipe(exampleRecipe)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Source,
			exampleInput.Description,
			exampleInput.InspiredByRecipeID,
			exampleInput.YieldsPortions,
			exampleInput.SealOfApproval,
			exampleInput.CreatedByUser,
		}

		db.ExpectBegin()

		db.ExpectExec(formatQueryForSQLMock(recipeCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() time.Time {
			return exampleRecipe.CreatedAt
		}

		actual, err := c.CreateRecipe(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error while also creating meal", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.Steps = nil
		exampleRecipe.PrepTasks = nil
		exampleRecipe.ID = "1"

		exampleInput := fakes.BuildFakeRecipeDatabaseCreationInputFromRecipe(exampleRecipe)
		exampleInput.AlsoCreateMeal = true
		exampleInput.Steps = []*types.RecipeStepDatabaseCreationInput{}

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		recipeCreationArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.Name,
			exampleRecipe.Source,
			exampleRecipe.Description,
			exampleRecipe.InspiredByRecipeID,
			exampleRecipe.YieldsPortions,
			exampleRecipe.SealOfApproval,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		mealCreationArgs := []interface{}{
			&idMatcher{},
			exampleRecipe.Name,
			exampleRecipe.Description,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealCreationQuery)).
			WithArgs(interfaceToDriverValue(mealCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		mealRecipeCreationArgs := []interface{}{
			&idMatcher{},
			&idMatcher{},
			exampleRecipe.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
			WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
			WillReturnError(errors.New("fart"))

		db.ExpectRollback()

		actual, err := c.CreateRecipe(ctx, exampleInput)
		require.Nil(t, actual)
		require.Error(t, err)

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
			exampleRecipe.YieldsPortions,
			exampleRecipe.SealOfApproval,
			exampleRecipe.CreatedByUser,
			exampleRecipe.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

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
			exampleRecipe.YieldsPortions,
			exampleRecipe.SealOfApproval,
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
			WillReturnResult(newArbitraryDatabaseResult())

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

func Test_findCreatedRecipeStepProductsForIngredients(T *testing.T) {
	T.Parallel()

	T.Run("sopa de frijol", func(t *testing.T) {
		t.Parallel()

		soak := fakes.BuildFakeValidPreparation()
		water := fakes.BuildFakeValidIngredient()
		pintoBeans := fakes.BuildFakeValidIngredient()
		garlicPaste := fakes.BuildFakeValidIngredient()
		productName := "soaked pinto beans"

		expected := &types.Recipe{
			Name:        "sopa de frijol",
			Description: "",
			Steps: []*types.RecipeStep{
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProduct{
						{
							ID:              fakes.BuildFakeID(),
							Name:            productName,
							MeasurementUnit: *fakes.BuildFakeValidMeasurementUnit(),
							Type:            types.RecipeStepProductIngredientType,
						},
					},
					Notes:       "first step",
					Preparation: *soak,
					Ingredients: []*types.RecipeStepIngredient{
						{
							RecipeStepProductID: nil,
							Ingredient:          pintoBeans,
							Name:                "pinto beans",
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     500,
							ProductOfRecipeStep: false,
						},
						{
							RecipeStepProductID: nil,
							Ingredient:          water,
							Name:                "water",
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     500,
							ProductOfRecipeStep: false,
						},
					},
					Index: 0,
				},
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProduct{
						{
							Name:            "final output",
							MeasurementUnit: *fakes.BuildFakeValidMeasurementUnit(),
							Type:            types.RecipeStepProductIngredientType,
						},
					},
					Notes:       "second step",
					Preparation: *soak,
					Ingredients: []*types.RecipeStepIngredient{
						{
							Ingredient:          nil,
							RecipeStepProductID: nil,
							Name:                productName,
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     1000,
							ProductOfRecipeStep: true,
						},
						{
							RecipeStepProductID: nil,
							Ingredient:          garlicPaste,
							Name:                "garlic paste",
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     10,
							ProductOfRecipeStep: false,
						},
					},
					Index: 1,
				},
			},
		}

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        expected.Name,
			Description: expected.Description,
		}

		for _, step := range expected.Steps {
			newStep := &types.RecipeStepDatabaseCreationInput{
				MinimumTemperatureInCelsius:   step.MinimumTemperatureInCelsius,
				MaximumTemperatureInCelsius:   step.MaximumTemperatureInCelsius,
				Notes:                         step.Notes,
				ExplicitInstructions:          step.ExplicitInstructions,
				PreparationID:                 step.Preparation.ID,
				BelongsToRecipe:               step.BelongsToRecipe,
				ID:                            step.ID,
				Index:                         step.Index,
				MinimumEstimatedTimeInSeconds: step.MinimumEstimatedTimeInSeconds,
				MaximumEstimatedTimeInSeconds: step.MaximumEstimatedTimeInSeconds,
				Optional:                      step.Optional,
			}

			for _, ingredient := range step.Ingredients {
				newIngredient := &types.RecipeStepIngredientDatabaseCreationInput{
					ID:                  ingredient.ID,
					BelongsToRecipeStep: ingredient.BelongsToRecipeStep,
					Name:                ingredient.Name,
					Optional:            ingredient.Optional,
					RecipeStepProductID: ingredient.RecipeStepProductID,
					MeasurementUnitID:   ingredient.MeasurementUnit.ID,
					QuantityNotes:       ingredient.QuantityNotes,
					IngredientNotes:     ingredient.IngredientNotes,
					MinimumQuantity:     ingredient.MinimumQuantity,
					ProductOfRecipeStep: ingredient.ProductOfRecipeStep,
				}
				if ingredient.Ingredient != nil {
					newIngredient.IngredientID = &ingredient.Ingredient.ID
				}
				newStep.Ingredients = append(newStep.Ingredients, newIngredient)
			}

			for _, product := range step.Products {
				newProduct := &types.RecipeStepProductDatabaseCreationInput{
					ID:                                 product.ID,
					Name:                               product.Name,
					Type:                               product.Type,
					MeasurementUnitID:                  product.MeasurementUnit.ID,
					QuantityNotes:                      product.QuantityNotes,
					Compostable:                        product.Compostable,
					MaximumStorageDurationInSeconds:    product.MaximumStorageDurationInSeconds,
					MinimumStorageTemperatureInCelsius: product.MinimumStorageTemperatureInCelsius,
					MaximumStorageTemperatureInCelsius: product.MaximumStorageTemperatureInCelsius,
					StorageInstructions:                product.StorageInstructions,
					BelongsToRecipeStep:                product.BelongsToRecipeStep,
					MinimumQuantity:                    product.MinimumQuantity,
					MaximumQuantity:                    product.MaximumQuantity,
				}
				newStep.Products = append(newStep.Products, newProduct)
			}

			exampleRecipeInput.Steps = append(exampleRecipeInput.Steps, newStep)
		}

		findCreatedRecipeStepProductsForIngredients(exampleRecipeInput, len(exampleRecipeInput.Steps)-1)

		require.NotNil(t, exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[0].Products[0].ID, *exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
	})

	T.Run("slightly more complicated recipe", func(t *testing.T) {
		t.Parallel()

		soak := fakes.BuildFakeValidPreparation()
		water := fakes.BuildFakeValidIngredient()
		pintoBeans := fakes.BuildFakeValidIngredient()
		garlicPaste := fakes.BuildFakeValidIngredient()
		productName := "soaked pinto beans"

		expected := &types.Recipe{
			Name:        "sopa de frijol",
			Description: "",
			Steps: []*types.RecipeStep{
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProduct{
						{
							ID:              fakes.BuildFakeID(),
							Name:            productName,
							MeasurementUnit: *fakes.BuildFakeValidMeasurementUnit(),
							Type:            types.RecipeStepProductIngredientType,
						},
					},
					Notes:       "first step",
					Preparation: *soak,
					Ingredients: []*types.RecipeStepIngredient{
						{
							RecipeStepProductID: nil,
							Ingredient:          pintoBeans,
							Name:                "pinto beans",
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     500,
							ProductOfRecipeStep: false,
						},
						{
							RecipeStepProductID: nil,
							Ingredient:          water,
							Name:                "water",
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     5,
							ProductOfRecipeStep: false,
						},
					},
					Index: 0,
				},
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProduct{
						{
							Name:            "pressure cooked beans",
							MeasurementUnit: *fakes.BuildFakeValidMeasurementUnit(),
							Type:            types.RecipeStepProductIngredientType,
						},
					},
					Notes:       "second step",
					Preparation: *soak,
					Ingredients: []*types.RecipeStepIngredient{
						{
							Ingredient:          nil,
							RecipeStepProductID: nil,
							Name:                productName,
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     1000,
							ProductOfRecipeStep: true,
						},
						{
							RecipeStepProductID: nil,
							Ingredient:          garlicPaste,
							Name:                "garlic paste",
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     10,
							ProductOfRecipeStep: false,
						},
					},
					Index: 1,
				},
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProduct{
						{
							ID:              fakes.BuildFakeID(),
							Name:            productName,
							MeasurementUnit: *fakes.BuildFakeValidMeasurementUnit(),
							Type:            types.RecipeStepProductIngredientType,
						},
					},
					Notes:       "third step",
					Preparation: *soak,
					Ingredients: []*types.RecipeStepIngredient{
						{
							RecipeStepProductID: nil,
							Ingredient:          pintoBeans,
							Name:                "pinto beans",
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     500,
							ProductOfRecipeStep: false,
						},
						{
							RecipeStepProductID: nil,
							Ingredient:          water,
							Name:                "water",
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     5,
							ProductOfRecipeStep: false,
						},
					},
					Index: 2,
				},
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProduct{
						{
							Name:            "final output",
							MeasurementUnit: *fakes.BuildFakeValidMeasurementUnit(),
							Type:            types.RecipeStepProductIngredientType,
						},
					},
					Notes:       "fourth step",
					Preparation: *soak,
					Ingredients: []*types.RecipeStepIngredient{
						{
							Ingredient:          nil,
							RecipeStepProductID: nil,
							Name:                productName,
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     1000,
							ProductOfRecipeStep: true,
						},
						{
							RecipeStepProductID: nil,
							Ingredient:          nil,
							Name:                "pressure cooked beans",
							MeasurementUnit:     *fakes.BuildFakeValidMeasurementUnit(),
							MinimumQuantity:     10,
							ProductOfRecipeStep: true,
						},
					},
					Index: 3,
				},
			},
		}

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        expected.Name,
			Description: expected.Description,
		}

		for _, step := range expected.Steps {
			newStep := &types.RecipeStepDatabaseCreationInput{
				MinimumTemperatureInCelsius:   step.MinimumTemperatureInCelsius,
				MaximumTemperatureInCelsius:   step.MaximumTemperatureInCelsius,
				Notes:                         step.Notes,
				ExplicitInstructions:          step.ExplicitInstructions,
				PreparationID:                 step.Preparation.ID,
				BelongsToRecipe:               step.BelongsToRecipe,
				ID:                            step.ID,
				Index:                         step.Index,
				MinimumEstimatedTimeInSeconds: step.MinimumEstimatedTimeInSeconds,
				MaximumEstimatedTimeInSeconds: step.MaximumEstimatedTimeInSeconds,
				Optional:                      step.Optional,
			}

			for _, ingredient := range step.Ingredients {
				newIngredient := &types.RecipeStepIngredientDatabaseCreationInput{
					ID:                  ingredient.ID,
					BelongsToRecipeStep: ingredient.BelongsToRecipeStep,
					Name:                ingredient.Name,
					RecipeStepProductID: ingredient.RecipeStepProductID,
					MeasurementUnitID:   ingredient.MeasurementUnit.ID,
					QuantityNotes:       ingredient.QuantityNotes,
					IngredientNotes:     ingredient.IngredientNotes,
					MinimumQuantity:     ingredient.MinimumQuantity,
					ProductOfRecipeStep: ingredient.ProductOfRecipeStep,
				}
				if ingredient.Ingredient != nil {
					newIngredient.IngredientID = &ingredient.Ingredient.ID
				}
				newStep.Ingredients = append(newStep.Ingredients, newIngredient)
			}

			for _, product := range step.Products {
				newProduct := &types.RecipeStepProductDatabaseCreationInput{
					ID:                                 product.ID,
					Name:                               product.Name,
					Type:                               product.Type,
					MeasurementUnitID:                  product.MeasurementUnit.ID,
					QuantityNotes:                      product.QuantityNotes,
					Compostable:                        product.Compostable,
					MaximumStorageDurationInSeconds:    product.MaximumStorageDurationInSeconds,
					MinimumStorageTemperatureInCelsius: product.MinimumStorageTemperatureInCelsius,
					MaximumStorageTemperatureInCelsius: product.MaximumStorageTemperatureInCelsius,
					StorageInstructions:                product.StorageInstructions,
					BelongsToRecipeStep:                product.BelongsToRecipeStep,
					MinimumQuantity:                    product.MinimumQuantity,
					MaximumQuantity:                    product.MaximumQuantity,
				}
				newStep.Products = append(newStep.Products, newProduct)
			}

			exampleRecipeInput.Steps = append(exampleRecipeInput.Steps, newStep)
		}

		for stepIndex := range exampleRecipeInput.Steps {
			findCreatedRecipeStepProductsForIngredients(exampleRecipeInput, stepIndex)
		}

		require.NotNil(t, exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[0].Products[0].ID, *exampleRecipeInput.Steps[1].Ingredients[0].RecipeStepProductID)
		require.NotNil(t, exampleRecipeInput.Steps[3].Ingredients[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[2].Products[0].ID, *exampleRecipeInput.Steps[3].Ingredients[0].RecipeStepProductID)
	})
}

func Test_findCreatedRecipeStepProductsForInstruments(T *testing.T) {
	T.Parallel()

	T.Run("example", func(t *testing.T) {
		t.Parallel()

		bake := fakes.BuildFakeValidPreparation()
		line := fakes.BuildFakeValidPreparation()
		bakingSheet := fakes.BuildFakeValidInstrument()
		aluminumFoil := fakes.BuildFakeValidIngredient()
		asparagus := fakes.BuildFakeValidIngredient()
		grams := fakes.BuildFakeValidMeasurementUnit()
		sheet := fakes.BuildFakeValidMeasurementUnit()

		productName := "lined baking sheet"

		expected := &types.Recipe{
			Name:        "example",
			Description: "",
			Steps: []*types.RecipeStep{
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProduct{
						{
							ID:   fakes.BuildFakeID(),
							Name: productName,
							Type: types.RecipeStepProductInstrumentType,
						},
					},
					Instruments: []*types.RecipeStepInstrument{
						{
							Instrument:          bakingSheet,
							RecipeStepProductID: nil,
							Name:                "baking sheet",
							ProductOfRecipeStep: false,
						},
					},
					Notes:       "first step",
					Preparation: *line,
					Ingredients: []*types.RecipeStepIngredient{
						{
							RecipeStepProductID: nil,
							Ingredient:          aluminumFoil,
							Name:                "aluminum foil",
							MeasurementUnit:     *sheet,
							MinimumQuantity:     1,
							ProductOfRecipeStep: false,
						},
					},
					Index: 0,
				},
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProduct{
						{
							ID:   fakes.BuildFakeID(),
							Name: "roasted asparagus",
							Type: types.RecipeStepProductInstrumentType,
						},
					},
					Instruments: []*types.RecipeStepInstrument{
						{
							Instrument:          bakingSheet,
							RecipeStepProductID: nil,
							Name:                productName,
							ProductOfRecipeStep: true,
						},
					},
					Notes:       "second step",
					Preparation: *bake,
					Ingredients: []*types.RecipeStepIngredient{
						{
							RecipeStepProductID: nil,
							Ingredient:          asparagus,
							Name:                "asparagus",
							MeasurementUnit:     *grams,
							MinimumQuantity:     1000,
							ProductOfRecipeStep: false,
						},
					},
					Index: 1,
				},
			},
		}

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        expected.Name,
			Description: expected.Description,
		}

		for _, step := range expected.Steps {
			newStep := &types.RecipeStepDatabaseCreationInput{
				MinimumTemperatureInCelsius:   step.MinimumTemperatureInCelsius,
				MaximumTemperatureInCelsius:   step.MaximumTemperatureInCelsius,
				Notes:                         step.Notes,
				ExplicitInstructions:          step.ExplicitInstructions,
				PreparationID:                 step.Preparation.ID,
				BelongsToRecipe:               step.BelongsToRecipe,
				ID:                            step.ID,
				Index:                         step.Index,
				MinimumEstimatedTimeInSeconds: step.MinimumEstimatedTimeInSeconds,
				MaximumEstimatedTimeInSeconds: step.MaximumEstimatedTimeInSeconds,
				Optional:                      step.Optional,
			}

			for _, ingredient := range step.Ingredients {
				newIngredient := &types.RecipeStepIngredientDatabaseCreationInput{
					IngredientID:        &ingredient.Ingredient.ID,
					ID:                  ingredient.ID,
					BelongsToRecipeStep: ingredient.BelongsToRecipeStep,
					Name:                ingredient.Name,
					RecipeStepProductID: ingredient.RecipeStepProductID,
					MeasurementUnitID:   ingredient.MeasurementUnit.ID,
					QuantityNotes:       ingredient.QuantityNotes,
					IngredientNotes:     ingredient.IngredientNotes,
					MinimumQuantity:     ingredient.MinimumQuantity,
					ProductOfRecipeStep: ingredient.ProductOfRecipeStep,
				}
				newStep.Ingredients = append(newStep.Ingredients, newIngredient)
			}

			for _, instrument := range step.Instruments {
				var instrumentID *string
				if instrument.Instrument != nil {
					instrumentID = &instrument.Instrument.ID
				}

				newInstrument := &types.RecipeStepInstrumentDatabaseCreationInput{
					InstrumentID:        instrumentID,
					RecipeStepProductID: instrument.RecipeStepProductID,
					ID:                  instrument.ID,
					Name:                instrument.Name,
					Notes:               instrument.Notes,
					BelongsToRecipeStep: instrument.BelongsToRecipeStep,
					ProductOfRecipeStep: instrument.ProductOfRecipeStep,
					PreferenceRank:      instrument.PreferenceRank,
				}
				newStep.Instruments = append(newStep.Instruments, newInstrument)
			}

			for _, product := range step.Products {
				measurementUnitID := product.MeasurementUnit.ID

				newProduct := &types.RecipeStepProductDatabaseCreationInput{
					ID:                                 product.ID,
					Name:                               product.Name,
					Type:                               product.Type,
					MeasurementUnitID:                  measurementUnitID,
					QuantityNotes:                      product.QuantityNotes,
					Compostable:                        product.Compostable,
					MaximumStorageDurationInSeconds:    product.MaximumStorageDurationInSeconds,
					MinimumStorageTemperatureInCelsius: product.MinimumStorageTemperatureInCelsius,
					MaximumStorageTemperatureInCelsius: product.MaximumStorageTemperatureInCelsius,
					StorageInstructions:                product.StorageInstructions,
					BelongsToRecipeStep:                product.BelongsToRecipeStep,
					MinimumQuantity:                    product.MinimumQuantity,
					MaximumQuantity:                    product.MaximumQuantity,
				}
				newStep.Products = append(newStep.Products, newProduct)
			}

			exampleRecipeInput.Steps = append(exampleRecipeInput.Steps, newStep)
		}

		findCreatedRecipeStepProductsForInstruments(exampleRecipeInput, len(exampleRecipeInput.Steps)-1)

		require.NotNil(t, exampleRecipeInput.Steps[1].Instruments[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[0].Products[0].ID, *exampleRecipeInput.Steps[1].Instruments[0].RecipeStepProductID)
	})
}
