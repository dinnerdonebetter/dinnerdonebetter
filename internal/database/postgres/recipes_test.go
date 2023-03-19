package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/internal/pointers"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
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
	"valid_preparations.minimum_ingredient_count",
	"valid_preparations.maximum_ingredient_count",
	"valid_preparations.minimum_instrument_count",
	"valid_preparations.maximum_instrument_count",
	"valid_preparations.temperature_required",
	"valid_preparations.time_estimate_required",
	"valid_preparations.condition_expression_required",
	"valid_preparations.consumes_vessel",
	"valid_preparations.only_for_vessels",
	"valid_preparations.universal",
	"valid_preparations.minimum_vessel_count",
	"valid_preparations.maximum_vessel_count",
	"valid_preparations.slug",
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
	"recipe_steps.condition_expression",
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
			&step.Preparation.MinimumIngredientCount,
			&step.Preparation.MaximumIngredientCount,
			&step.Preparation.MinimumInstrumentCount,
			&step.Preparation.MaximumInstrumentCount,
			&step.Preparation.TemperatureRequired,
			&step.Preparation.TimeEstimateRequired,
			&step.Preparation.ConditionExpressionRequired,
			&step.Preparation.ConsumesVessel,
			&step.Preparation.OnlyForVessels,
			&step.Preparation.Universal,
			&step.Preparation.MinimumVesselCount,
			&step.Preparation.MaximumVesselCount,
			&step.Preparation.Slug,
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
			&step.ConditionExpression,
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
		args := []any{
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
		args := []any{
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
		args := []any{
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

func prepareMockToSuccessfullyGetRecipe(t *testing.T, recipe *types.Recipe, userID string, db *sqlmockExpecterWrapper) {
	t.Helper()

	var (
		exampleRecipe *types.Recipe
	)

	if recipe == nil {
		exampleRecipe = fakes.BuildFakeRecipe()
	} else {
		exampleRecipe = recipe
	}

	allIngredients := []*types.RecipeStepIngredient{}
	allInstruments := []*types.RecipeStepInstrument{}
	allVessels := []*types.RecipeStepVessel{}
	allProducts := []*types.RecipeStepProduct{}
	allCompletionConditions := []*types.RecipeStepCompletionCondition{}
	for _, step := range exampleRecipe.Steps {
		allIngredients = append(allIngredients, step.Ingredients...)
		allInstruments = append(allInstruments, step.Instruments...)
		allVessels = append(allVessels, step.Vessels...)
		allProducts = append(allProducts, step.Products...)
		allCompletionConditions = append(allCompletionConditions, step.CompletionConditions...)
	}

	args := []any{
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

	listRecipePrepTasksForRecipeArgs := []any{
		exampleRecipe.ID,
	}

	db.ExpectQuery(formatQueryForSQLMock(listRecipePrepTasksForRecipeQuery)).
		WithArgs(interfaceToDriverValue(listRecipePrepTasksForRecipeArgs)...).
		WillReturnRows(buildMockRowsFromRecipePrepTasks(exampleRecipe.PrepTasks...))

	recipeMediaForRecipeArgs := []any{
		exampleRecipe.ID,
	}

	db.ExpectQuery(formatQueryForSQLMock(recipeMediaForRecipeQuery)).
		WithArgs(interfaceToDriverValue(recipeMediaForRecipeArgs)...).
		WillReturnRows(buildMockRowsFromRecipeMedia(exampleRecipe.Media...))

	getRecipeStepIngredientsForRecipeArgs := []any{
		exampleRecipe.ID,
	}
	db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientsForRecipeQuery)).
		WithArgs(interfaceToDriverValue(getRecipeStepIngredientsForRecipeArgs)...).
		WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, allIngredients...))

	productsArgs := []any{
		exampleRecipe.ID,
	}
	db.ExpectQuery(formatQueryForSQLMock(getRecipeStepProductsForRecipeQuery)).
		WithArgs(interfaceToDriverValue(productsArgs)...).
		WillReturnRows(buildMockRowsFromRecipeStepProducts(false, 0, allProducts...))

	instrumentsArgs := []any{
		exampleRecipe.ID,
	}
	db.ExpectQuery(formatQueryForSQLMock(getRecipeStepInstrumentsForRecipeQuery)).
		WithArgs(interfaceToDriverValue(instrumentsArgs)...).
		WillReturnRows(buildMockRowsFromRecipeStepInstruments(false, 0, allInstruments...))

	vesselsArgs := []any{
		exampleRecipe.ID,
	}
	db.ExpectQuery(formatQueryForSQLMock(getRecipeStepVesselsForRecipeQuery)).
		WithArgs(interfaceToDriverValue(vesselsArgs)...).
		WillReturnRows(buildMockRowsFromRecipeStepVessels(false, 0, allVessels...))

	completionConditionsArgs := []any{
		exampleRecipe.ID,
	}
	db.ExpectQuery(formatQueryForSQLMock(getRecipeStepCompletionConditionsForRecipeQuery)).
		WithArgs(interfaceToDriverValue(completionConditionsArgs)...).
		WillReturnRows(buildMockRowsFromRecipeStepCompletionConditions(false, 0, allCompletionConditions...))

	for _, step := range exampleRecipe.Steps {
		recipeMediaForRecipeStepArgs := []any{
			exampleRecipe.ID,
			step.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeMediaForRecipeStepQuery)).
			WithArgs(interfaceToDriverValue(recipeMediaForRecipeStepArgs)...).
			WillReturnRows(buildMockRowsFromRecipeMedia(step.Media...))
	}
}

func TestQuerier_getRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleUserID := fakes.BuildFakeID()

		prepareMockToSuccessfullyGetRecipe(t, exampleRecipe, exampleUserID, db)

		actual, err := c.getRecipe(ctx, exampleRecipe.ID, exampleUserID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipe, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error fetching recipe step ingredients", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipe.ID,
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDAndAuthorIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockFullRowsFromRecipe(exampleRecipe))

		listRecipePrepTasksForRecipeArgs := []any{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(listRecipePrepTasksForRecipeQuery)).
			WithArgs(interfaceToDriverValue(listRecipePrepTasksForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipePrepTasks(exampleRecipe.PrepTasks...))

		recipeMediaForRecipeArgs := []any{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeMediaForRecipeQuery)).
			WithArgs(interfaceToDriverValue(recipeMediaForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeMedia(exampleRecipe.Media...))

		getRecipeStepIngredientsForRecipeArgs := []any{
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

		allIngredients := []*types.RecipeStepIngredient{}
		for _, step := range exampleRecipe.Steps {
			allIngredients = append(allIngredients, step.Ingredients...)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipe.ID,
			exampleUserID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDAndAuthorIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockFullRowsFromRecipe(exampleRecipe))

		listRecipePrepTasksForRecipeArgs := []any{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(listRecipePrepTasksForRecipeQuery)).
			WithArgs(interfaceToDriverValue(listRecipePrepTasksForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipePrepTasks(exampleRecipe.PrepTasks...))

		recipeMediaForRecipeArgs := []any{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeMediaForRecipeQuery)).
			WithArgs(interfaceToDriverValue(recipeMediaForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeMedia(exampleRecipe.Media...))

		getRecipeStepIngredientsForRecipeArgs := []any{
			exampleRecipe.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(getRecipeStepIngredientsForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, allIngredients...))

		productsArgs := []any{
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

		ctx := context.Background()
		c, db := buildTestClient(t)

		prepareMockToSuccessfullyGetRecipe(t, exampleRecipe, "", db)

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

		args := []any{
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

		for _, step := range exampleRecipe.Steps {
			step.Ingredients = []*types.RecipeStepIngredient{
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
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

		ctx := context.Background()
		c, db := buildTestClient(t)

		prepareMockToSuccessfullyGetRecipe(t, exampleRecipe, exampleRecipe.CreatedByUser, db)

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

		args := []any{
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

		for _, step := range exampleRecipe.Steps {
			step.Ingredients = []*types.RecipeStepIngredient{
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
				fakes.BuildFakeRecipeStepIngredient(),
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
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
		for i := range exampleRecipeList.Data {
			exampleRecipeList.Data[i].Steps = nil
			exampleRecipeList.Data[i].PrepTasks = nil
			exampleRecipeList.Data[i].Media = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipes", nil, nil, nil, householdOwnershipColumn, recipesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipes(true, exampleRecipeList.FilteredCount, exampleRecipeList.Data...))

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
		for i := range exampleRecipeList.Data {
			exampleRecipeList.Data[i].Steps = nil
			exampleRecipeList.Data[i].PrepTasks = nil
			exampleRecipeList.Data[i].Media = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipes", nil, nil, nil, householdOwnershipColumn, recipesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipes(true, exampleRecipeList.FilteredCount, exampleRecipeList.Data...))

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
		for i := range exampleRecipeList.Data {
			exampleRecipeList.Data[i].Steps = nil
			exampleRecipeList.Data[i].PrepTasks = nil
			exampleRecipeIDs = append(exampleRecipeIDs, exampleRecipeList.Data[i].ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
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
		for i := range exampleRecipeList.Data {
			exampleRecipeList.Data[i].Steps = nil
			exampleRecipeList.Data[i].PrepTasks = nil
			exampleRecipeList.Data[i].Media = nil
		}

		ctx := context.Background()
		recipeNameQuery := "example"
		c, db := buildTestClient(t)

		where := squirrel.ILike{"name": wrapQueryForILIKE(recipeNameQuery)}
		query, args := c.buildListQueryWithILike(ctx, "recipes", nil, nil, where, "", recipesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipes(true, exampleRecipeList.FilteredCount, exampleRecipeList.Data...))

		actual, err := c.SearchForRecipes(ctx, recipeNameQuery, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeList := fakes.BuildFakeRecipeList()
		for i := range exampleRecipeList.Data {
			exampleRecipeList.Data[i].Steps = nil
			exampleRecipeList.Data[i].PrepTasks = nil
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
		for i := range exampleRecipeList.Data {
			exampleRecipeList.Data[i].Steps = nil
			exampleRecipeList.Data[i].PrepTasks = nil
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
		exampleRecipe.Media = nil
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
				exampleRecipe.Steps[i].Instruments[j].ID = "4"
				exampleRecipe.Steps[i].Instruments[j].Instrument = &types.ValidInstrument{ID: exampleRecipe.Steps[i].Instruments[j].Instrument.ID}
				exampleRecipe.Steps[i].Instruments[j].BelongsToRecipeStep = "2"
			}

			for j := range step.Vessels {
				exampleRecipe.Steps[i].Vessels[j].ID = "5"
				exampleRecipe.Steps[i].Vessels[j].Instrument = &types.ValidInstrument{ID: exampleRecipe.Steps[i].Vessels[j].Instrument.ID}
				exampleRecipe.Steps[i].Vessels[j].BelongsToRecipeStep = "2"
			}

			for j := range step.CompletionConditions {
				exampleRecipe.Steps[i].CompletionConditions[j].ID = "6"
				exampleRecipe.Steps[i].CompletionConditions[j].BelongsToRecipeStep = "2"
				exampleRecipe.Steps[i].CompletionConditions[j].Ingredients = nil
			}

			step.Products = nil
			step.Media = nil
		}

		for i := range exampleRecipe.PrepTasks {
			exampleRecipe.PrepTasks[i].BelongsToRecipe = exampleRecipe.ID
		}

		exampleInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(exampleRecipe)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		recipeCreationArgs := []any{
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
			recipeStepCreationArgs := []any{
				step.ID,
				step.Index,
				step.PreparationID,
				step.MinimumEstimatedTimeInSeconds,
				step.MaximumEstimatedTimeInSeconds,
				step.MinimumTemperatureInCelsius,
				step.MaximumTemperatureInCelsius,
				step.Notes,
				step.ExplicitInstructions,
				step.ConditionExpression,
				step.Optional,
				step.BelongsToRecipe,
			}

			db.ExpectExec(formatQueryForSQLMock(recipeStepCreationQuery)).
				WithArgs(interfaceToDriverValue(recipeStepCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())

			for _, ingredient := range step.Ingredients {
				recipeStepIngredientCreationArgs := []any{
					ingredient.ID,
					ingredient.Name,
					ingredient.Optional,
					ingredient.IngredientID,
					ingredient.MeasurementUnitID,
					ingredient.MinimumQuantity,
					ingredient.MaximumQuantity,
					ingredient.QuantityNotes,
					ingredient.RecipeStepProductID,
					ingredient.IngredientNotes,
					ingredient.OptionIndex,
					ingredient.RequiresDefrost,
					ingredient.ToTaste,
					ingredient.ProductPercentageToUse,
					ingredient.VesselIndex,
					ingredient.BelongsToRecipeStep,
				}

				db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
					WithArgs(interfaceToDriverValue(recipeStepIngredientCreationArgs)...).
					WillReturnResult(newArbitraryDatabaseResult())
			}

			for _, instrument := range step.Instruments {
				recipeStepInstrumentCreationArgs := []any{
					instrument.ID,
					instrument.InstrumentID,
					instrument.RecipeStepProductID,
					instrument.Name,
					instrument.Notes,
					instrument.PreferenceRank,
					instrument.Optional,
					instrument.OptionIndex,
					instrument.MinimumQuantity,
					instrument.MaximumQuantity,
					instrument.BelongsToRecipeStep,
				}

				db.ExpectExec(formatQueryForSQLMock(recipeStepInstrumentCreationQuery)).
					WithArgs(interfaceToDriverValue(recipeStepInstrumentCreationArgs)...).
					WillReturnResult(newArbitraryDatabaseResult())
			}

			for _, vessel := range step.Vessels {
				recipeStepVesselCreationArgs := []any{
					vessel.ID,
					vessel.Name,
					vessel.Notes,
					vessel.BelongsToRecipeStep,
					vessel.RecipeStepProductID,
					vessel.InstrumentID,
					vessel.VesselPredicate,
					vessel.MinimumQuantity,
					vessel.MaximumQuantity,
					vessel.UnavailableAfterStep,
				}

				db.ExpectExec(formatQueryForSQLMock(recipeStepVesselCreationQuery)).
					WithArgs(interfaceToDriverValue(recipeStepVesselCreationArgs)...).
					WillReturnResult(newArbitraryDatabaseResult())
			}

			for _, completionCondition := range step.CompletionConditions {
				recipeStepCompletionConditionCreationArgs := []any{
					completionCondition.ID,
					completionCondition.BelongsToRecipeStep,
					completionCondition.IngredientStateID,
					completionCondition.Optional,
					completionCondition.Notes,
				}

				db.ExpectExec(formatQueryForSQLMock(recipeStepCompletionConditionCreationQuery)).
					WithArgs(interfaceToDriverValue(recipeStepCompletionConditionCreationArgs)...).
					WillReturnResult(newArbitraryDatabaseResult())

				for _, completionConditionIngredient := range completionCondition.Ingredients {
					recipeStepCompletionConditionIngredientCreationArgs := []any{
						completionConditionIngredient.ID,
						completionConditionIngredient.BelongsToRecipeStepCompletionCondition,
						completionConditionIngredient.RecipeStepIngredient,
					}

					db.ExpectExec(formatQueryForSQLMock(recipeStepCompletionConditionIngredientCreationQuery)).
						WithArgs(interfaceToDriverValue(recipeStepCompletionConditionIngredientCreationArgs)...).
						WillReturnResult(newArbitraryDatabaseResult())
				}
			}
		}

		for _, prepTask := range exampleInput.PrepTasks {
			createRecipePrepTaskQueryArgs := []any{
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
				createRecipePrepTaskStepArgs := []any{
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

			for j, vessel := range step.Vessels {
				vessel.BelongsToRecipeStep = step.ID
				vessel.CreatedAt = actual.Steps[i].Vessels[j].CreatedAt
			}

			for j, completionCondition := range step.CompletionConditions {
				completionCondition.BelongsToRecipeStep = step.ID
				completionCondition.CreatedAt = actual.Steps[i].CompletionConditions[j].CreatedAt
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
		exampleRecipe.Media = nil
		exampleRecipe.ID = "1"

		exampleInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(exampleRecipe)
		exampleInput.AlsoCreateMeal = true
		exampleInput.Steps = []*types.RecipeStepDatabaseCreationInput{}

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		recipeCreationArgs := []any{
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

		mealCreationArgs := []any{
			&idMatcher{},
			exampleRecipe.Name,
			exampleRecipe.Description,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealCreationQuery)).
			WithArgs(interfaceToDriverValue(mealCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		mealRecipeCreationArgs := []any{
			&idMatcher{},
			&idMatcher{},
			exampleRecipe.ID,
			types.MealComponentTypesMain,
		}

		db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
			WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

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
		}

		assert.Equal(t, exampleRecipe, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipe.ID = "1"
		exampleInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(exampleRecipe)

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
		exampleInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(exampleRecipe)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
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

		exampleInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(exampleRecipe)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		recipeCreationArgs := []any{
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

		recipeStepCreationArgs := []any{
			exampleInput.Steps[0].ID,
			0,
			exampleInput.Steps[0].PreparationID,
			exampleInput.Steps[0].MinimumEstimatedTimeInSeconds,
			exampleInput.Steps[0].MaximumEstimatedTimeInSeconds,
			exampleInput.Steps[0].MinimumTemperatureInCelsius,
			exampleInput.Steps[0].MaximumTemperatureInCelsius,
			exampleInput.Steps[0].Notes,
			exampleInput.Steps[0].ExplicitInstructions,
			exampleInput.Steps[0].ConditionExpression,
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
		exampleInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(exampleRecipe)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
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

		exampleInput := converters.ConvertRecipeToRecipeDatabaseCreationInput(exampleRecipe)
		exampleInput.AlsoCreateMeal = true
		exampleInput.Steps = []*types.RecipeStepDatabaseCreationInput{}

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		recipeCreationArgs := []any{
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

		mealCreationArgs := []any{
			&idMatcher{},
			exampleRecipe.Name,
			exampleRecipe.Description,
			exampleRecipe.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealCreationQuery)).
			WithArgs(interfaceToDriverValue(mealCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		mealRecipeCreationArgs := []any{
			&idMatcher{},
			&idMatcher{},
			exampleRecipe.ID,
			types.MealComponentTypesMain,
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

		args := []any{
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

		args := []any{
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

		args := []any{
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

		args := []any{
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

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        "sopa de frijol",
			Description: "",
			Steps: []*types.RecipeStepDatabaseCreationInput{
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:                fakes.BuildFakeID(),
							Name:              productName,
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "first step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							IngredientID:      &pintoBeans.ID,
							Name:              "pinto beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
						{
							IngredientID:      &water.ID,
							Name:              "water",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
					},
					Index: 0,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							Name:              "final output",
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "second step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							Name:                            productName,
							MeasurementUnitID:               fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:                 1000,
							ProductOfRecipeStepProductIndex: pointers.Pointer(uint64(0)),
							ProductOfRecipeStepIndex:        pointers.Pointer(uint64(0)),
						},
						{
							IngredientID:      &garlicPaste.ID,
							Name:              "garlic paste",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   10,
						},
					},
					Index: 1,
				},
			},
		}

		findCreatedRecipeStepProductsForIngredients(exampleRecipeInput)

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

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        "sopa de frijol",
			Description: "",
			Steps: []*types.RecipeStepDatabaseCreationInput{
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:                fakes.BuildFakeID(),
							Name:              productName,
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "first step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							IngredientID:      &pintoBeans.ID,
							Name:              "pinto beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
						{
							IngredientID:      &water.ID,
							Name:              "water",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   5,
						},
					},
					Index: 0,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							Name:              "pressure cooked beans",
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "second step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							Name:                            productName,
							MeasurementUnitID:               fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:                 1000,
							ProductOfRecipeStepIndex:        pointers.Pointer(uint64(0)),
							ProductOfRecipeStepProductIndex: pointers.Pointer(uint64(0)),
						},
						{
							IngredientID:      &garlicPaste.ID,
							Name:              "garlic paste",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   10,
						},
					},
					Index: 1,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:                fakes.BuildFakeID(),
							Name:              productName,
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "third step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							IngredientID:      &pintoBeans.ID,
							Name:              "pinto beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   500,
						},
						{
							IngredientID:      &water.ID,
							Name:              "water",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   5,
						},
					},
					Index: 2,
				},
				{
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							Name:              "final output",
							MeasurementUnitID: &fakes.BuildFakeValidMeasurementUnit().ID,
							Type:              types.RecipeStepProductIngredientType,
						},
					},
					Notes:         "fourth step",
					PreparationID: soak.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							Name:                            productName,
							MeasurementUnitID:               fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:                 1000,
							ProductOfRecipeStepIndex:        pointers.Pointer(uint64(2)),
							ProductOfRecipeStepProductIndex: pointers.Pointer(uint64(0)),
						},
						{
							Name:              "pressure cooked beans",
							MeasurementUnitID: fakes.BuildFakeValidMeasurementUnit().ID,
							MinimumQuantity:   10,
						},
					},
					Index: 3,
				},
			},
		}

		findCreatedRecipeStepProductsForIngredients(exampleRecipeInput)

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

		exampleRecipeInput := &types.RecipeDatabaseCreationInput{
			Name:        "example",
			Description: "",
			Steps: []*types.RecipeStepDatabaseCreationInput{
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:   fakes.BuildFakeID(),
							Name: productName,
							Type: types.RecipeStepProductInstrumentType,
						},
					},
					Instruments: []*types.RecipeStepInstrumentDatabaseCreationInput{
						{
							InstrumentID:        &bakingSheet.ID,
							RecipeStepProductID: nil,
							Name:                "baking sheet",
						},
					},
					Notes:         "first step",
					PreparationID: line.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							RecipeStepProductID: nil,
							IngredientID:        &aluminumFoil.ID,
							Name:                "aluminum foil",
							MeasurementUnitID:   sheet.ID,
							MinimumQuantity:     1,
						},
					},
					Index: 0,
				},
				{
					MinimumTemperatureInCelsius: nil,
					MaximumTemperatureInCelsius: nil,
					Products: []*types.RecipeStepProductDatabaseCreationInput{
						{
							ID:   fakes.BuildFakeID(),
							Name: "roasted asparagus",
							Type: types.RecipeStepProductInstrumentType,
						},
					},
					Instruments: []*types.RecipeStepInstrumentDatabaseCreationInput{
						{
							InstrumentID:                    &bakingSheet.ID,
							RecipeStepProductID:             nil,
							Name:                            productName,
							ProductOfRecipeStepIndex:        pointers.Pointer(uint64(0)),
							ProductOfRecipeStepProductIndex: pointers.Pointer(uint64(0)),
						},
					},
					Notes:         "second step",
					PreparationID: bake.ID,
					Ingredients: []*types.RecipeStepIngredientDatabaseCreationInput{
						{
							RecipeStepProductID: nil,
							IngredientID:        &asparagus.ID,
							Name:                "asparagus",
							MeasurementUnitID:   grams.ID,
							MinimumQuantity:     1000,
						},
					},
					Index: 1,
				},
			},
		}

		findCreatedRecipeStepProductsForInstruments(exampleRecipeInput)

		require.NotNil(t, exampleRecipeInput.Steps[1].Instruments[0].RecipeStepProductID)
		assert.Equal(t, exampleRecipeInput.Steps[0].Products[0].ID, *exampleRecipeInput.Steps[1].Instruments[0].RecipeStepProductID)
	})
}
