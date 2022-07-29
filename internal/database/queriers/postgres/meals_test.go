package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromMeals(includeCounts bool, filteredCount uint64, meals ...*types.Meal) *sqlmock.Rows {
	columns := mealsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range meals {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Description,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.CreatedByUser,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(meals))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildMockFullRowsFromMeal(meal *types.Meal) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(mealsTableWithRecipeIDColumns)

	for _, recipe := range meal.Recipes {
		exampleRows.AddRow(
			&meal.ID,
			&meal.Name,
			&meal.Description,
			&meal.CreatedOn,
			&meal.LastUpdatedOn,
			&meal.ArchivedOn,
			&meal.CreatedByUser,
			&recipe.ID,
		)
	}

	return exampleRows
}

func TestQuerier_ScanMeals(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanMeals(ctx, mockRows, false)
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

		_, _, _, err := q.scanMeals(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_MealExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMeal := fakes.BuildFakeMeal()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleMeal.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.MealExists(ctx, exampleMeal.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.MealExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMeal := fakes.BuildFakeMeal()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleMeal.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.MealExists(ctx, exampleMeal.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMeal := fakes.BuildFakeMeal()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleMeal.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.MealExists(ctx, exampleMeal.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		for i := range exampleMeal.Recipes {
			for j := range exampleMeal.Recipes[i].Steps {
				exampleMeal.Recipes[i].Steps[j].Products = nil
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealArgs := []interface{}{
			exampleMeal.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealByIDQuery)).
			WithArgs(interfaceToDriverValue(getMealArgs)...).
			WillReturnRows(buildMockFullRowsFromMeal(exampleMeal))

		for _, recipe := range exampleMeal.Recipes {
			allIngredients := []*types.RecipeStepIngredient{}
			allProducts := []*types.RecipeStepProduct{}
			for _, step := range recipe.Steps {
				allIngredients = append(allIngredients, step.Ingredients...)
				allProducts = append(allProducts, step.Products...)
			}

			getRecipeArgs := []interface{}{
				recipe.ID,
			}

			db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDQuery)).
				WithArgs(interfaceToDriverValue(getRecipeArgs)...).
				WillReturnRows(buildMockFullRowsFromRecipe(recipe))

			query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, nil, false)
			db.ExpectQuery(formatQueryForSQLMock(query)).
				WithArgs(interfaceToDriverValue(args)...).
				WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, allIngredients...))

			productsArgs := []interface{}{
				recipe.ID,
				recipe.ID,
			}
			db.ExpectQuery(formatQueryForSQLMock(getRecipeStepProductsForRecipeQuery)).
				WithArgs(interfaceToDriverValue(productsArgs)...).
				WillReturnRows(buildMockRowsFromRecipeStepProducts(false, 0, allProducts...))
		}

		actual, err := c.GetMeal(ctx, exampleMeal.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleMeal, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMeal(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealArgs := []interface{}{
			exampleMeal.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealByIDQuery)).
			WithArgs(interfaceToDriverValue(getMealArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMeal(ctx, exampleMeal.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response from database", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealArgs := []interface{}{
			exampleMeal.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealByIDQuery)).
			WithArgs(interfaceToDriverValue(getMealArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMeal(ctx, exampleMeal.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error querying for recipes", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealArgs := []interface{}{
			exampleMeal.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealByIDQuery)).
			WithArgs(interfaceToDriverValue(getMealArgs)...).
			WillReturnRows(buildMockFullRowsFromMeal(exampleMeal))

		getRecipeArgs := []interface{}{
			exampleMeal.Recipes[0].ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeByIDQuery)).
			WithArgs(interfaceToDriverValue(getRecipeArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMeal(ctx, exampleMeal.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetTotalMealCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalMealsCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetTotalMealCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalMealsCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetTotalMealCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMeals(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealList := fakes.BuildFakeMealList()
		for i := range exampleMealList.Meals {
			exampleMealList.Meals[i].Recipes = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meals", nil, nil, nil, householdOwnershipColumn, mealsTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMeals(true, exampleMealList.FilteredCount, exampleMealList.Meals...))

		actual, err := c.GetMeals(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleMealList := fakes.BuildFakeMealList()
		exampleMealList.Page = 0
		exampleMealList.Limit = 0
		for i := range exampleMealList.Meals {
			exampleMealList.Meals[i].Recipes = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meals", nil, nil, nil, householdOwnershipColumn, mealsTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMeals(true, exampleMealList.FilteredCount, exampleMealList.Meals...))

		actual, err := c.GetMeals(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meals", nil, nil, nil, householdOwnershipColumn, mealsTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMeals(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meals", nil, nil, nil, householdOwnershipColumn, mealsTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMeals(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForMeals(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealList := fakes.BuildFakeMealList()
		for i := range exampleMealList.Meals {
			exampleMealList.Meals[i].Recipes = nil
		}

		ctx := context.Background()
		recipeNameQuery := "example"
		c, db := buildTestClient(t)

		where := squirrel.ILike{"name": wrapQueryForILIKE(recipeNameQuery)}
		query, args := c.buildListQueryWithILike(ctx, "meals", nil, nil, where, "", mealsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMeals(true, exampleMealList.FilteredCount, exampleMealList.Meals...))

		actual, err := c.SearchForMeals(ctx, recipeNameQuery, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealList := fakes.BuildFakeMealList()
		for i := range exampleMealList.Meals {
			exampleMealList.Meals[i].Recipes = nil
		}

		ctx := context.Background()
		recipeNameQuery := "example"
		c, db := buildTestClient(t)

		where := squirrel.ILike{"name": wrapQueryForILIKE(recipeNameQuery)}
		query, args := c.buildListQueryWithILike(ctx, "meals", nil, nil, where, "", mealsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForMeals(ctx, recipeNameQuery, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealList := fakes.BuildFakeMealList()
		for i := range exampleMealList.Meals {
			exampleMealList.Meals[i].Recipes = nil
		}

		ctx := context.Background()
		recipeNameQuery := "example"
		c, db := buildTestClient(t)

		where := squirrel.ILike{"name": wrapQueryForILIKE(recipeNameQuery)}
		query, args := c.buildListQueryWithILike(ctx, "meals", nil, nil, where, "", mealsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForMeals(ctx, recipeNameQuery, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealList := fakes.BuildFakeMealList()
		for i := range exampleMealList.Meals {
			exampleMealList.Meals[i].Recipes = nil
		}

		var exampleIDs []string
		for _, x := range exampleMealList.Meals {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetMealsWithIDsQuery(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMeals(false, 0, exampleMealList.Meals...))

		actual, err := c.GetMealsWithIDs(ctx, exampleHouseholdID, 0, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealList.Meals, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid IDs", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealsWithIDs(ctx, exampleHouseholdID, defaultLimit, nil)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealList := fakes.BuildFakeMealList()

		var exampleIDs []string
		for _, x := range exampleMealList.Meals {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetMealsWithIDsQuery(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealsWithIDs(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning query results", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealList := fakes.BuildFakeMealList()

		var exampleIDs []string
		for _, x := range exampleMealList.Meals {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetMealsWithIDsQuery(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealsWithIDs(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()
		exampleMeal.ID = "1"

		exampleInput := fakes.BuildFakeMealDatabaseCreationInputFromMeal(exampleMeal)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		mealCreationArgs := []interface{}{
			exampleMeal.ID,
			exampleMeal.Name,
			exampleMeal.Description,
			exampleMeal.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealCreationQuery)).
			WithArgs(interfaceToDriverValue(mealCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, recipeID := range exampleInput.Recipes {
			mealRecipeCreationArgs := []interface{}{
				&idMatcher{},
				exampleMeal.ID,
				recipeID,
			}

			db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
				WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		db.ExpectCommit()

		c.timeFunc = func() uint64 {
			return exampleMeal.CreatedOn
		}

		actual, err := c.CreateMeal(ctx, exampleInput)
		assert.NoError(t, err)
		exampleMeal.Recipes = nil

		assert.Equal(t, exampleMeal, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.CreateMeal(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error starting transaction", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()
		exampleInput := fakes.BuildFakeMealDatabaseCreationInputFromMeal(exampleMeal)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.CreateMeal(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error creating meal", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()
		exampleMeal.ID = "1"

		exampleInput := fakes.BuildFakeMealDatabaseCreationInputFromMeal(exampleMeal)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		mealCreationArgs := []interface{}{
			exampleMeal.ID,
			exampleMeal.Name,
			exampleMeal.Description,
			exampleMeal.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealCreationQuery)).
			WithArgs(interfaceToDriverValue(mealCreationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.CreateMeal(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error creating meal recipe", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()
		exampleMeal.ID = "1"

		exampleInput := fakes.BuildFakeMealDatabaseCreationInputFromMeal(exampleMeal)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		mealCreationArgs := []interface{}{
			exampleMeal.ID,
			exampleMeal.Name,
			exampleMeal.Description,
			exampleMeal.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealCreationQuery)).
			WithArgs(interfaceToDriverValue(mealCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		mealRecipeCreationArgs := []interface{}{
			&idMatcher{},
			exampleMeal.ID,
			exampleInput.Recipes[0],
		}

		db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
			WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleMeal.CreatedOn
		}

		actual, err := c.CreateMeal(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()
		exampleMeal.ID = "1"

		exampleInput := fakes.BuildFakeMealDatabaseCreationInputFromMeal(exampleMeal)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		mealCreationArgs := []interface{}{
			exampleMeal.ID,
			exampleMeal.Name,
			exampleMeal.Description,
			exampleMeal.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealCreationQuery)).
			WithArgs(interfaceToDriverValue(mealCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, recipeID := range exampleInput.Recipes {
			mealRecipeCreationArgs := []interface{}{
				&idMatcher{},
				exampleMeal.ID,
				recipeID,
			}

			db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
				WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() uint64 {
			return exampleMeal.CreatedOn
		}

		actual, err := c.CreateMeal(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()
		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mealRecipeCreationArgs := []interface{}{
			&idMatcher{},
			exampleMeal.ID,
			exampleRecipe.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
			WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() uint64 {
			return exampleMeal.CreatedOn
		}

		err := c.CreateMealRecipe(ctx, c.db, exampleMeal.ID, exampleRecipe.ID)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing meal ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.CreateMealRecipe(ctx, c.db, "", exampleRecipe.ID)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.CreateMealRecipe(ctx, c.db, exampleMeal.ID, "")
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error performing query", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()
		exampleRecipe := fakes.BuildFakeRecipe()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mealRecipeCreationArgs := []interface{}{
			&idMatcher{},
			exampleMeal.ID,
			exampleRecipe.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
			WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
			WillReturnError(errors.New("blah"))

		c.timeFunc = func() uint64 {
			return exampleMeal.CreatedOn
		}

		err := c.CreateMealRecipe(ctx, c.db, exampleMeal.ID, exampleRecipe.ID)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHouseholdID,
			exampleMeal.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveMeal(ctx, exampleMeal.ID, exampleHouseholdID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMeal(ctx, "", exampleHouseholdID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMeal(ctx, exampleMeal.ID, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleHouseholdID,
			exampleMeal.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveMeal(ctx, exampleMeal.ID, exampleHouseholdID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
