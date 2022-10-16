package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromRecipeStepIngredients(includeCounts bool, filteredCount uint64, recipeStepIngredients ...*types.RecipeStepIngredient) *sqlmock.Rows {
	columns := recipeStepIngredientsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepIngredients {
		var ingredient *types.NullableValidIngredient
		if x.Ingredient != nil {
			ingredient = &types.NullableValidIngredient{
				CreatedAt:                               &x.Ingredient.CreatedAt,
				LastUpdatedAt:                           x.Ingredient.LastUpdatedAt,
				ArchivedAt:                              x.Ingredient.ArchivedAt,
				ID:                                      &x.Ingredient.ID,
				Warning:                                 &x.Ingredient.Warning,
				Description:                             &x.Ingredient.Description,
				IconPath:                                &x.Ingredient.IconPath,
				PluralName:                              &x.Ingredient.PluralName,
				StorageInstructions:                     &x.Ingredient.StorageInstructions,
				Name:                                    &x.Ingredient.Name,
				MaximumIdealStorageTemperatureInCelsius: x.Ingredient.MaximumIdealStorageTemperatureInCelsius,
				MinimumIdealStorageTemperatureInCelsius: x.Ingredient.MinimumIdealStorageTemperatureInCelsius,
				ContainsShellfish:                       &x.Ingredient.ContainsShellfish,
				ContainsDairy:                           &x.Ingredient.ContainsDairy,
				AnimalFlesh:                             &x.Ingredient.AnimalFlesh,
				IsMeasuredVolumetrically:                &x.Ingredient.IsMeasuredVolumetrically,
				IsLiquid:                                &x.Ingredient.IsLiquid,
				ContainsPeanut:                          &x.Ingredient.ContainsPeanut,
				ContainsTreeNut:                         &x.Ingredient.ContainsTreeNut,
				ContainsEgg:                             &x.Ingredient.ContainsEgg,
				ContainsWheat:                           &x.Ingredient.ContainsWheat,
				ContainsSoy:                             &x.Ingredient.ContainsSoy,
				AnimalDerived:                           &x.Ingredient.AnimalDerived,
				RestrictToPreparations:                  &x.Ingredient.RestrictToPreparations,
				ContainsSesame:                          &x.Ingredient.ContainsSesame,
				ContainsFish:                            &x.Ingredient.ContainsFish,
				ContainsGluten:                          &x.Ingredient.ContainsGluten,
			}
		}

		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Optional,
			ingredient.ID,
			ingredient.Name,
			ingredient.Description,
			ingredient.Warning,
			ingredient.ContainsEgg,
			ingredient.ContainsDairy,
			ingredient.ContainsPeanut,
			ingredient.ContainsTreeNut,
			ingredient.ContainsSoy,
			ingredient.ContainsWheat,
			ingredient.ContainsShellfish,
			ingredient.ContainsSesame,
			ingredient.ContainsFish,
			ingredient.ContainsGluten,
			ingredient.AnimalFlesh,
			ingredient.IsMeasuredVolumetrically,
			ingredient.IsLiquid,
			ingredient.IconPath,
			ingredient.AnimalDerived,
			ingredient.PluralName,
			ingredient.RestrictToPreparations,
			ingredient.MinimumIdealStorageTemperatureInCelsius,
			ingredient.MaximumIdealStorageTemperatureInCelsius,
			ingredient.StorageInstructions,
			ingredient.CreatedAt,
			ingredient.LastUpdatedAt,
			ingredient.ArchivedAt,
			x.MeasurementUnit.ID,
			x.MeasurementUnit.Name,
			x.MeasurementUnit.Description,
			x.MeasurementUnit.Volumetric,
			x.MeasurementUnit.IconPath,
			x.MeasurementUnit.Universal,
			x.MeasurementUnit.Metric,
			x.MeasurementUnit.Imperial,
			x.MeasurementUnit.PluralName,
			x.MeasurementUnit.CreatedAt,
			x.MeasurementUnit.LastUpdatedAt,
			x.MeasurementUnit.ArchivedAt,
			x.MinimumQuantity,
			x.MaximumQuantity,
			x.QuantityNotes,
			x.ProductOfRecipeStep,
			x.RecipeStepProductID,
			x.IngredientNotes,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
			x.BelongsToRecipeStep,
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

func TestQuerier_getRecipeStepIngredientsForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		getRecipeStepIngredientsForRecipeArgs := []interface{}{
			exampleRecipeID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(getRecipeStepIngredientsForRecipeArgs)...).
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

	T.Run("with error scanning results", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		getRecipeStepIngredientsForRecipeArgs := []interface{}{
			exampleRecipeID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(getRecipeStepIngredientsForRecipeArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.getRecipeStepIngredientsForRecipe(ctx, exampleRecipeID)
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

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, filter)
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

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, filter)
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

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, filter)
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

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepIngredientsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepIngredients(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
		exampleRecipeStepIngredient.ID = "1"
		exampleRecipeStepIngredient.MeasurementUnit = types.ValidMeasurementUnit{ID: exampleRecipeStepIngredient.MeasurementUnit.ID}
		exampleRecipeStepIngredient.Ingredient = &types.ValidIngredient{ID: exampleRecipeStepIngredient.Ingredient.ID}
		exampleInput := converters.ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Optional,
			exampleInput.IngredientID,
			exampleInput.MeasurementUnitID,
			exampleInput.MinimumQuantity,
			exampleInput.MaximumQuantity,
			exampleInput.QuantityNotes,
			exampleInput.ProductOfRecipeStep,
			exampleInput.RecipeStepProductID,
			exampleInput.IngredientNotes,
			exampleInput.BelongsToRecipeStep,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleRecipeStepIngredient.CreatedAt
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
		exampleInput := converters.ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Optional,
			exampleInput.IngredientID,
			exampleInput.MeasurementUnitID,
			exampleInput.MinimumQuantity,
			exampleInput.MaximumQuantity,
			exampleInput.QuantityNotes,
			exampleInput.ProductOfRecipeStep,
			exampleInput.RecipeStepProductID,
			exampleInput.IngredientNotes,
			exampleInput.BelongsToRecipeStep,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleRecipeStepIngredient.CreatedAt
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

		exampleInput := converters.ConvertRecipeStepIngredientToRecipeStepIngredientDatabaseCreationInput(exampleRecipeStepIngredient)

		ctx := context.Background()
		c, db := buildTestClient(t)

		recipeStepIngredientCreationArgs := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Optional,
			exampleInput.IngredientID,
			exampleInput.MeasurementUnitID,
			exampleInput.MinimumQuantity,
			exampleInput.MaximumQuantity,
			exampleInput.QuantityNotes,
			exampleInput.ProductOfRecipeStep,
			exampleInput.RecipeStepProductID,
			exampleInput.IngredientNotes,
			exampleInput.BelongsToRecipeStep,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeStepIngredientCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleRecipeStepIngredient.CreatedAt
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
		exampleRecipeStepIngredient.MeasurementUnit = types.ValidMeasurementUnit{ID: exampleRecipeStepIngredient.MeasurementUnit.ID}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepIngredient.Ingredient.ID,
			exampleRecipeStepIngredient.Name,
			exampleRecipeStepIngredient.Optional,
			exampleRecipeStepIngredient.MeasurementUnit.ID,
			exampleRecipeStepIngredient.MinimumQuantity,
			exampleRecipeStepIngredient.MaximumQuantity,
			exampleRecipeStepIngredient.QuantityNotes,
			exampleRecipeStepIngredient.ProductOfRecipeStep,
			exampleRecipeStepIngredient.RecipeStepProductID,
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
		exampleRecipeStepIngredient.MeasurementUnit = types.ValidMeasurementUnit{ID: exampleRecipeStepIngredient.MeasurementUnit.ID}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepIngredient.Ingredient.ID,
			exampleRecipeStepIngredient.Name,
			exampleRecipeStepIngredient.Optional,
			exampleRecipeStepIngredient.MeasurementUnit.ID,
			exampleRecipeStepIngredient.MinimumQuantity,
			exampleRecipeStepIngredient.MaximumQuantity,
			exampleRecipeStepIngredient.QuantityNotes,
			exampleRecipeStepIngredient.ProductOfRecipeStep,
			exampleRecipeStepIngredient.RecipeStepProductID,
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
