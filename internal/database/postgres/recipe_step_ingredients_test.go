package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			ingredient = converters.ConvertValidIngredientToNullableValidIngredient(x.Ingredient)
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
			ingredient.Slug,
			ingredient.ContainsAlcohol,
			ingredient.ShoppingSuggestions,
			ingredient.IsStarch,
			ingredient.IsProtein,
			ingredient.IsGrain,
			ingredient.IsFruit,
			ingredient.IsSalt,
			ingredient.IsFat,
			ingredient.IsAcid,
			ingredient.IsHeat,
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
			x.MeasurementUnit.Slug,
			x.MeasurementUnit.PluralName,
			x.MeasurementUnit.CreatedAt,
			x.MeasurementUnit.LastUpdatedAt,
			x.MeasurementUnit.ArchivedAt,
			x.MinimumQuantity,
			x.MaximumQuantity,
			x.QuantityNotes,
			x.RecipeStepProductID,
			x.IngredientNotes,
			x.OptionIndex,
			x.ToTaste,
			x.ProductPercentageToUse,
			x.VesselIndex,
			x.RecipeStepProductRecipeID,
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
}

func TestQuerier_GetRecipeStepIngredient(T *testing.T) {
	T.Parallel()

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
}

func TestQuerier_getRecipeStepIngredientsForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepIngredientList := fakes.BuildFakeRecipeStepIngredientList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		getRecipeStepIngredientsForRecipeArgs := []any{
			exampleRecipeID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepIngredientsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(getRecipeStepIngredientsForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(false, 0, exampleRecipeStepIngredientList.Data...))

		actual, err := c.getRecipeStepIngredientsForRecipe(ctx, exampleRecipeID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepIngredientList.Data, actual)

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

		getRecipeStepIngredientsForRecipeArgs := []any{
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

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, recipeStepIngredientsTableColumns, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepIngredients(true, exampleRecipeStepIngredientList.FilteredCount, exampleRecipeStepIngredientList.Data...))

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

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, recipeStepIngredientsTableColumns, filter)
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

		query, args := c.buildListQuery(ctx, "recipe_step_ingredients", getRecipeStepIngredientsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, recipeStepIngredientsTableColumns, filter)
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

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepIngredient(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestSQLQuerier_createRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.createRecipeStepIngredient(ctx, c.db, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipeStepIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepIngredient(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeStepIngredient(T *testing.T) {
	T.Parallel()

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
}
