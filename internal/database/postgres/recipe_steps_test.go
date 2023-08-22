package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromRecipeSteps(includeCounts bool, filteredCount uint64, recipeSteps ...*types.RecipeStep) *sqlmock.Rows {
	columns := recipeStepsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeSteps {
		rowValues := []driver.Value{
			x.ID,
			x.Index,
			x.Preparation.ID,
			x.Preparation.Name,
			x.Preparation.Description,
			x.Preparation.IconPath,
			x.Preparation.YieldsNothing,
			x.Preparation.RestrictToIngredients,
			x.Preparation.MinimumIngredientCount,
			x.Preparation.MaximumIngredientCount,
			x.Preparation.MinimumInstrumentCount,
			x.Preparation.MaximumInstrumentCount,
			x.Preparation.TemperatureRequired,
			x.Preparation.TimeEstimateRequired,
			x.Preparation.ConditionExpressionRequired,
			x.Preparation.ConsumesVessel,
			x.Preparation.OnlyForVessels,
			x.Preparation.MinimumVesselCount,
			x.Preparation.MaximumVesselCount,
			x.Preparation.Slug,
			x.Preparation.PastTense,
			x.Preparation.CreatedAt,
			x.Preparation.LastUpdatedAt,
			x.Preparation.ArchivedAt,
			x.MinimumEstimatedTimeInSeconds,
			x.MaximumEstimatedTimeInSeconds,
			x.MinimumTemperatureInCelsius,
			x.MaximumTemperatureInCelsius,
			x.Notes,
			x.ExplicitInstructions,
			x.ConditionExpression,
			x.Optional,
			x.StartTimerAutomatically,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
			x.BelongsToRecipe,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(recipeSteps))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipeSteps(ctx, mockRows, false)
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

		_, _, _, err := q.scanRecipeSteps(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_RecipeStepExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepExists(ctx, "", exampleRecipeStep.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepExists(ctx, exampleRecipeID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.Instruments = nil
		exampleRecipeStep.Vessels = nil
		exampleRecipeStep.Ingredients = nil
		exampleRecipeStep.Products = nil
		exampleRecipeStep.CompletionConditions = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeID,
			exampleRecipeStep.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeSteps(false, 0, exampleRecipeStep))

		actual, err := c.GetRecipeStep(ctx, exampleRecipeID, exampleRecipeStep.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStep, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStep(ctx, "", exampleRecipeStep.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStep(ctx, exampleRecipeID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeID,
			exampleRecipeStep.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStep(ctx, exampleRecipeID, exampleRecipeStep.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_getRecipeStepByID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.Instruments = nil
		exampleRecipeStep.Vessels = nil
		exampleRecipeStep.Ingredients = nil
		exampleRecipeStep.Products = nil
		exampleRecipeStep.CompletionConditions = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStep.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeSteps(false, 0, exampleRecipeStep))

		actual, err := c.getRecipeStepByID(ctx, c.db, exampleRecipeStep.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStep, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.getRecipeStepByID(ctx, c.db, exampleRecipeStep.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.getRecipeStepByID(ctx, c.db, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStep.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.getRecipeStepByID(ctx, c.db, exampleRecipeStep.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepList := fakes.BuildFakeRecipeStepList()

		for i := range exampleRecipeStepList.Data {
			exampleRecipeStepList.Data[i].Instruments = nil
			exampleRecipeStepList.Data[i].Vessels = nil
			exampleRecipeStepList.Data[i].Ingredients = nil
			exampleRecipeStepList.Data[i].Products = nil
			exampleRecipeStepList.Data[i].CompletionConditions = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_steps", getRecipeStepsJoins, []string{"valid_preparations.id"}, nil, householdOwnershipColumn, recipeStepsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeSteps(true, exampleRecipeStepList.FilteredCount, exampleRecipeStepList.Data...))

		actual, err := c.GetRecipeSteps(ctx, exampleRecipeID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeSteps(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_steps", getRecipeStepsJoins, []string{"valid_preparations.id"}, nil, householdOwnershipColumn, recipeStepsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeSteps(ctx, exampleRecipeID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_steps", getRecipeStepsJoins, []string{"valid_preparations.id"}, nil, householdOwnershipColumn, recipeStepsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeSteps(ctx, exampleRecipeID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.ID = "1"
		exampleRecipeStep.Ingredients = nil
		exampleRecipeStep.Products = nil
		exampleRecipeStep.Instruments = nil
		exampleRecipeStep.Vessels = nil
		exampleRecipeStep.CompletionConditions = nil
		exampleRecipeStep.Preparation = types.ValidPreparation{}
		exampleInput := converters.ConvertRecipeStepToRecipeStepDatabaseCreationInput(exampleRecipeStep)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Index,
			exampleInput.PreparationID,
			exampleInput.MinimumEstimatedTimeInSeconds,
			exampleInput.MaximumEstimatedTimeInSeconds,
			exampleInput.MinimumTemperatureInCelsius,
			exampleInput.MaximumTemperatureInCelsius,
			exampleInput.Notes,
			exampleInput.ExplicitInstructions,
			exampleInput.ConditionExpression,
			exampleInput.Optional,
			exampleInput.StartTimerAutomatically,
			exampleInput.BelongsToRecipe,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleRecipeStep.CreatedAt
		}

		actual, err := c.CreateRecipeStep(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStep, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStep(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleInput := converters.ConvertRecipeStepToRecipeStepDatabaseCreationInput(exampleRecipeStep)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Index,
			exampleInput.PreparationID,
			exampleInput.MinimumEstimatedTimeInSeconds,
			exampleInput.MaximumEstimatedTimeInSeconds,
			exampleInput.MinimumTemperatureInCelsius,
			exampleInput.MaximumTemperatureInCelsius,
			exampleInput.Notes,
			exampleInput.ExplicitInstructions,
			exampleInput.ConditionExpression,
			exampleInput.Optional,
			exampleInput.StartTimerAutomatically,
			exampleInput.BelongsToRecipe,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleRecipeStep.CreatedAt
		}

		actual, err := c.CreateRecipeStep(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_createRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.ID = "2"
		exampleRecipeStep.BelongsToRecipe = "1"
		exampleRecipeStep.Preparation = types.ValidPreparation{}
		for i := range exampleRecipeStep.Ingredients {
			exampleRecipeStep.Ingredients[i].ID = "3"
			exampleRecipeStep.Ingredients[i].BelongsToRecipeStep = exampleRecipeStep.ID
		}

		for i := range exampleRecipeStep.Products {
			exampleRecipeStep.Products[i].ID = "4"
			exampleRecipeStep.Products[i].BelongsToRecipeStep = exampleRecipeStep.ID
		}

		for i := range exampleRecipeStep.Instruments {
			exampleRecipeStep.Instruments[i].ID = "5"
			exampleRecipeStep.Instruments[i].BelongsToRecipeStep = exampleRecipeStep.ID
		}

		for i := range exampleRecipeStep.Vessels {
			exampleRecipeStep.Vessels[i].ID = "6"
			exampleRecipeStep.Vessels[i].BelongsToRecipeStep = exampleRecipeStep.ID
		}

		for i := range exampleRecipeStep.CompletionConditions {
			exampleRecipeStep.CompletionConditions[i].ID = "7"
			exampleRecipeStep.CompletionConditions[i].BelongsToRecipeStep = exampleRecipeStep.ID
		}

		exampleInput := converters.ConvertRecipeStepToRecipeStepDatabaseCreationInput(exampleRecipeStep)

		ctx := context.Background()
		c, db := buildTestClient(t)

		recipeStepCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Index,
			exampleInput.PreparationID,
			exampleInput.MinimumEstimatedTimeInSeconds,
			exampleInput.MaximumEstimatedTimeInSeconds,
			exampleInput.MinimumTemperatureInCelsius,
			exampleInput.MaximumTemperatureInCelsius,
			exampleInput.Notes,
			exampleInput.ExplicitInstructions,
			exampleInput.ConditionExpression,
			exampleInput.Optional,
			exampleInput.StartTimerAutomatically,
			exampleInput.BelongsToRecipe,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeStepCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, ingredient := range exampleInput.Ingredients {
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
				ingredient.ToTaste,
				ingredient.ProductPercentageToUse,
				ingredient.VesselIndex,
				ingredient.RecipeStepProductRecipeID,
				ingredient.BelongsToRecipeStep,
			}

			db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
				WithArgs(interfaceToDriverValue(recipeStepIngredientCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		for _, product := range exampleInput.Products {
			args := []any{
				product.ID,
				product.Name,
				product.Type,
				product.MeasurementUnitID,
				product.MinimumQuantity,
				product.MaximumQuantity,
				product.QuantityNotes,
				product.Compostable,
				product.MaximumStorageDurationInSeconds,
				product.MinimumStorageTemperatureInCelsius,
				product.MaximumStorageTemperatureInCelsius,
				product.StorageInstructions,
				product.BelongsToRecipeStep,
				product.IsLiquid,
				product.IsWaste,
				product.Index,
				product.ContainedInVesselIndex,
			}

			db.ExpectExec(formatQueryForSQLMock(recipeStepProductCreationQuery)).
				WithArgs(interfaceToDriverValue(args)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		for _, instrument := range exampleInput.Instruments {
			args := []any{
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
				WithArgs(interfaceToDriverValue(args)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		for _, vessel := range exampleInput.Vessels {
			args := []any{
				vessel.ID,
				vessel.Name,
				vessel.Notes,
				vessel.BelongsToRecipeStep,
				vessel.RecipeStepProductID,
				vessel.VesselID,
				vessel.VesselPreposition,
				vessel.MinimumQuantity,
				vessel.MaximumQuantity,
				vessel.UnavailableAfterStep,
			}

			db.ExpectExec(formatQueryForSQLMock(recipeStepVesselCreationQuery)).
				WithArgs(interfaceToDriverValue(args)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		for _, completionCondition := range exampleInput.CompletionConditions {
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

		c.timeFunc = func() time.Time {
			return exampleRecipeStep.CreatedAt
		}

		actual, err := c.createRecipeStep(ctx, c.db, exampleInput)
		assert.NoError(t, err)
		assert.NotNil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error creating recipe step ingredient", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.ID = "2"
		exampleRecipeStep.BelongsToRecipe = "1"
		exampleRecipeStep.Preparation = types.ValidPreparation{}
		for i := range exampleRecipeStep.Ingredients {
			exampleRecipeStep.Ingredients[i].ID = "3"
			exampleRecipeStep.Ingredients[i].BelongsToRecipeStep = "2"
		}

		exampleInput := converters.ConvertRecipeStepToRecipeStepDatabaseCreationInput(exampleRecipeStep)

		ctx := context.Background()
		c, db := buildTestClient(t)

		recipeStepCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Index,
			exampleInput.PreparationID,
			exampleInput.MinimumEstimatedTimeInSeconds,
			exampleInput.MaximumEstimatedTimeInSeconds,
			exampleInput.MinimumTemperatureInCelsius,
			exampleInput.MaximumTemperatureInCelsius,
			exampleInput.Notes,
			exampleInput.ExplicitInstructions,
			exampleInput.ConditionExpression,
			exampleInput.Optional,
			exampleInput.StartTimerAutomatically,
			exampleInput.BelongsToRecipe,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeStepCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		ingredient := exampleInput.Ingredients[0]
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
			ingredient.ToTaste,
			ingredient.ProductPercentageToUse,
			ingredient.VesselIndex,
			ingredient.RecipeStepProductRecipeID,
			ingredient.BelongsToRecipeStep,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeStepIngredientCreationArgs)...).
			WillReturnError(errors.New("blah"))

		c.timeFunc = func() time.Time {
			return exampleRecipeStep.CreatedAt
		}

		actual, err := c.createRecipeStep(ctx, c.db, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error creating recipe step product", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()
		exampleRecipeStep.ID = "2"
		exampleRecipeStep.BelongsToRecipe = "1"
		exampleRecipeStep.Preparation = types.ValidPreparation{}
		for i := range exampleRecipeStep.Ingredients {
			exampleRecipeStep.Ingredients[i].ID = "3"
			exampleRecipeStep.Ingredients[i].BelongsToRecipeStep = exampleRecipeStep.ID
		}

		for i := range exampleRecipeStep.Products {
			exampleRecipeStep.Products[i].ID = "3"
			exampleRecipeStep.Products[i].BelongsToRecipeStep = exampleRecipeStep.ID
		}

		exampleInput := converters.ConvertRecipeStepToRecipeStepDatabaseCreationInput(exampleRecipeStep)

		ctx := context.Background()
		c, db := buildTestClient(t)

		recipeStepCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Index,
			exampleInput.PreparationID,
			exampleInput.MinimumEstimatedTimeInSeconds,
			exampleInput.MaximumEstimatedTimeInSeconds,
			exampleInput.MinimumTemperatureInCelsius,
			exampleInput.MaximumTemperatureInCelsius,
			exampleInput.Notes,
			exampleInput.ExplicitInstructions,
			exampleInput.ConditionExpression,
			exampleInput.Optional,
			exampleInput.StartTimerAutomatically,
			exampleInput.BelongsToRecipe,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeStepCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, ingredient := range exampleInput.Ingredients {
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
				ingredient.ToTaste,
				ingredient.ProductPercentageToUse,
				ingredient.VesselIndex,
				ingredient.RecipeStepProductRecipeID,
				ingredient.BelongsToRecipeStep,
			}

			db.ExpectExec(formatQueryForSQLMock(recipeStepIngredientCreationQuery)).
				WithArgs(interfaceToDriverValue(recipeStepIngredientCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		args := []any{
			exampleInput.Products[0].ID,
			exampleInput.Products[0].Name,
			exampleInput.Products[0].Type,
			exampleInput.Products[0].MeasurementUnitID,
			exampleInput.Products[0].MinimumQuantity,
			exampleInput.Products[0].MaximumQuantity,
			exampleInput.Products[0].QuantityNotes,
			exampleInput.Products[0].Compostable,
			exampleInput.Products[0].MaximumStorageDurationInSeconds,
			exampleInput.Products[0].MinimumStorageTemperatureInCelsius,
			exampleInput.Products[0].MaximumStorageTemperatureInCelsius,
			exampleInput.Products[0].StorageInstructions,
			exampleInput.Products[0].BelongsToRecipeStep,
			exampleInput.Products[0].IsLiquid,
			exampleInput.Products[0].IsWaste,
			exampleInput.Products[0].Index,
			exampleInput.Products[0].ContainedInVesselIndex,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepProductCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		c.timeFunc = func() time.Time {
			return exampleRecipeStep.CreatedAt
		}

		actual, err := c.createRecipeStep(ctx, c.db, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStep(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStep := fakes.BuildFakeRecipeStep()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStep(ctx, "", exampleRecipeStep.ID))
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStep(ctx, exampleRecipeID, ""))
	})
}
