package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
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
}

func TestQuerier_getRecipeStepByID(T *testing.T) {
	T.Parallel()

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

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStep(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestSQLQuerier_createRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.createRecipeStep(ctx, c.db, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
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
