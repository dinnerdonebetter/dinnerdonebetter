package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			x.MinimumEstimatedPortions,
			x.MaximumEstimatedPortions,
			x.EligibleForMealPlans,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
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
	exampleRows := sqlmock.NewRows([]string{
		"meals.id",
		"meals.name",
		"meals.description",
		"meals.min_estimated_portions",
		"meals.max_estimated_portions",
		"meals.eligible_for_meal_plans",
		"meals.created_at",
		"meals.last_updated_at",
		"meals.archived_at",
		"meals.created_by_user",
		"meal_components.recipe_id",
		"meal_components.scale",
		"meal_components.meal_component_type",
	})

	for _, mealComponent := range meal.Components {
		exampleRows.AddRow(
			&meal.ID,
			&meal.Name,
			&meal.Description,
			&meal.MinimumEstimatedPortions,
			&meal.MaximumEstimatedPortions,
			&meal.EligibleForMealPlans,
			&meal.CreatedAt,
			&meal.LastUpdatedAt,
			&meal.ArchivedAt,
			&meal.CreatedByUser,
			&mealComponent.Recipe.ID,
			&mealComponent.RecipeScale,
			&mealComponent.ComponentType,
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
		args := []any{
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
		args := []any{
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
		args := []any{
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

func prepareMockToSuccessfullyGetMeal(t *testing.T, exampleMeal *types.Meal, db *sqlmockExpecterWrapper) {
	t.Helper()

	getMealArgs := []any{
		exampleMeal.ID,
	}

	db.ExpectQuery(formatQueryForSQLMock(getMealByIDQuery)).
		WithArgs(interfaceToDriverValue(getMealArgs)...).
		WillReturnRows(buildMockFullRowsFromMeal(exampleMeal))

	for _, component := range exampleMeal.Components {
		prepareMockToSuccessfullyGetRecipe(t, &component.Recipe, "", db)
	}
}

func TestQuerier_GetMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		for i := range exampleMeal.Components {
			for j := range exampleMeal.Components[i].Recipe.Steps {
				exampleMeal.Components[i].Recipe.Steps[j].Products = nil
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		prepareMockToSuccessfullyGetMeal(t, exampleMeal, db)

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

		getMealArgs := []any{
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

		getMealArgs := []any{
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

		getMealArgs := []any{
			exampleMeal.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealByIDQuery)).
			WithArgs(interfaceToDriverValue(getMealArgs)...).
			WillReturnRows(buildMockFullRowsFromMeal(exampleMeal))

		getRecipeArgs := []any{
			exampleMeal.Components[0].Recipe.ID,
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

func TestQuerier_GetMeals(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealList := fakes.BuildFakeMealList()
		for i := range exampleMealList.Data {
			exampleMealList.Data[i].Components = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meals", nil, nil, nil, householdOwnershipColumn, mealsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMeals(true, exampleMealList.FilteredCount, exampleMealList.Data...))

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
		for i := range exampleMealList.Data {
			exampleMealList.Data[i].Components = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meals", nil, nil, nil, householdOwnershipColumn, mealsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMeals(true, exampleMealList.FilteredCount, exampleMealList.Data...))

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

		query, args := c.buildListQuery(ctx, "meals", nil, nil, nil, householdOwnershipColumn, mealsTableColumns, "", false, filter)

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

		query, args := c.buildListQuery(ctx, "meals", nil, nil, nil, householdOwnershipColumn, mealsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMeals(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealThatNeedSearchIndexing(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealList := fakes.BuildFakeMealList()

		c, db := buildTestClient(t)

		exampleIDs := []string{}
		for _, exampleMeal := range exampleMealList.Data {
			exampleIDs = append(exampleIDs, exampleMeal.ID)
		}

		db.ExpectQuery(formatQueryForSQLMock(mealsNeedingIndexingQuery)).
			WithArgs(interfaceToDriverValue(nil)...).
			WillReturnRows(buildMockRowsFromIDs(exampleIDs...))

		actual, err := c.GetMealIDsThatNeedSearchIndexing(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleIDs, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForMeals(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealList := fakes.BuildFakeMealList()
		for i := range exampleMealList.Data {
			exampleMealList.Data[i].Components = nil
		}

		ctx := context.Background()
		recipeNameQuery := "example"
		c, db := buildTestClient(t)

		where := squirrel.ILike{"name": wrapQueryForILIKE(recipeNameQuery)}
		query, args := c.buildListQueryWithILike(ctx, "meals", nil, nil, where, "", mealsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMeals(true, exampleMealList.FilteredCount, exampleMealList.Data...))

		actual, err := c.SearchForMeals(ctx, recipeNameQuery, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealList := fakes.BuildFakeMealList()
		for i := range exampleMealList.Data {
			exampleMealList.Data[i].Components = nil
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
		for i := range exampleMealList.Data {
			exampleMealList.Data[i].Components = nil
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

func TestQuerier_CreateMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()
		exampleMeal.ID = "1"

		exampleInput := converters.ConvertMealToMealDatabaseCreationInput(exampleMeal)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		mealCreationArgs := []any{
			exampleMeal.ID,
			exampleMeal.Name,
			exampleMeal.Description,
			exampleMeal.MinimumEstimatedPortions,
			exampleMeal.MaximumEstimatedPortions,
			exampleMeal.EligibleForMealPlans,
			exampleMeal.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealCreationQuery)).
			WithArgs(interfaceToDriverValue(mealCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, component := range exampleInput.Components {
			mealRecipeCreationArgs := []any{
				&idMatcher{},
				exampleMeal.ID,
				component.RecipeID,
				component.ComponentType,
				component.RecipeScale,
			}

			db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
				WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		db.ExpectCommit()

		c.timeFunc = func() time.Time {
			return exampleMeal.CreatedAt
		}

		actual, err := c.CreateMeal(ctx, exampleInput)
		assert.NoError(t, err)
		exampleMeal.Components = nil

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
		exampleInput := converters.ConvertMealToMealDatabaseCreationInput(exampleMeal)

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

		exampleInput := converters.ConvertMealToMealDatabaseCreationInput(exampleMeal)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		mealCreationArgs := []any{
			exampleMeal.ID,
			exampleMeal.Name,
			exampleMeal.Description,
			exampleMeal.MinimumEstimatedPortions,
			exampleMeal.MaximumEstimatedPortions,
			exampleMeal.EligibleForMealPlans,
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

		exampleInput := converters.ConvertMealToMealDatabaseCreationInput(exampleMeal)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		mealCreationArgs := []any{
			exampleMeal.ID,
			exampleMeal.Name,
			exampleMeal.Description,
			exampleMeal.MinimumEstimatedPortions,
			exampleMeal.MaximumEstimatedPortions,
			exampleMeal.EligibleForMealPlans,
			exampleMeal.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealCreationQuery)).
			WithArgs(interfaceToDriverValue(mealCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		mealRecipeCreationArgs := []any{
			&idMatcher{},
			exampleMeal.ID,
			exampleInput.Components[0].RecipeID,
			exampleInput.Components[0].ComponentType,
			exampleInput.Components[0].RecipeScale,
		}

		db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
			WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() time.Time {
			return exampleMeal.CreatedAt
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

		exampleInput := converters.ConvertMealToMealDatabaseCreationInput(exampleMeal)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		mealCreationArgs := []any{
			exampleMeal.ID,
			exampleMeal.Name,
			exampleMeal.Description,
			exampleMeal.MinimumEstimatedPortions,
			exampleMeal.MaximumEstimatedPortions,
			exampleMeal.EligibleForMealPlans,
			exampleMeal.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealCreationQuery)).
			WithArgs(interfaceToDriverValue(mealCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, component := range exampleInput.Components {
			mealRecipeCreationArgs := []any{
				&idMatcher{},
				exampleMeal.ID,
				component.RecipeID,
				component.ComponentType,
				component.RecipeScale,
			}

			db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
				WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() time.Time {
			return exampleMeal.CreatedAt
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

		ctx := context.Background()
		c, db := buildTestClient(t)

		mealRecipeCreationArgs := []any{
			&idMatcher{},
			exampleMeal.ID,
			exampleMeal.Components[0].Recipe.ID,
			exampleMeal.Components[0].ComponentType,
			exampleMeal.Components[0].RecipeScale,
		}

		db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
			WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleMeal.CreatedAt
		}

		exampleInput := converters.ConvertMealComponentToMealComponentDatabaseCreationInput(exampleMeal.Components[0])

		err := c.CreateMealComponent(ctx, c.db, exampleMeal.ID, exampleInput)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing meal ID", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		exampleInput := converters.ConvertMealComponentToMealComponentDatabaseCreationInput(exampleMeal.Components[0])

		err := c.CreateMealComponent(ctx, c.db, "", exampleInput)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing input", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.CreateMealComponent(ctx, c.db, exampleMeal.ID, nil)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error performing query", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mealRecipeCreationArgs := []any{
			&idMatcher{},
			exampleMeal.ID,
			exampleMeal.Components[0].Recipe.ID,
			exampleMeal.Components[0].ComponentType,
			exampleMeal.Components[0].RecipeScale,
		}

		db.ExpectExec(formatQueryForSQLMock(mealRecipeCreationQuery)).
			WithArgs(interfaceToDriverValue(mealRecipeCreationArgs)...).
			WillReturnError(errors.New("blah"))

		c.timeFunc = func() time.Time {
			return exampleMeal.CreatedAt
		}

		exampleInput := converters.ConvertMealComponentToMealComponentDatabaseCreationInput(exampleMeal.Components[0])

		err := c.CreateMealComponent(ctx, c.db, exampleMeal.ID, exampleInput)
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

		args := []any{
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

		args := []any{
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

func TestQuerier_MarkMealAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMeal := fakes.BuildFakeMeal()

		c, db := buildTestClient(t)

		args := []any{
			exampleMeal.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealLastIndexedAtQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.MarkMealAsIndexed(ctx, exampleMeal.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkMealAsIndexed(ctx, ""))
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMeal := fakes.BuildFakeMeal()

		c, db := buildTestClient(t)

		args := []any{
			exampleMeal.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealLastIndexedAtQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.MarkMealAsIndexed(ctx, exampleMeal.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
