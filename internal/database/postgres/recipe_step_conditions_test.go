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

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func buildMockRowsFromRecipeStepConditions(includeCounts bool, filteredCount uint64, recipeStepConditions ...*types.RecipeStepCondition) *sqlmock.Rows {
	columns := recipeStepConditionsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepConditions {
		rowValues := []driver.Value{
			x.ID,
			x.BelongsToRecipeStep,
			x.IngredientState.ID,
			x.Optional,
			x.Notes,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(recipeStepConditions))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanRecipeStepConditions(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipeStepConditions(ctx, mockRows, false)
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

		_, _, _, err := q.scanRecipeStepConditions(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_RecipeStepConditionExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		c, db := buildTestClient(t)
		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCondition.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepConditionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.RecipeStepConditionExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepConditionExists(ctx, "", exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepConditionExists(ctx, exampleRecipeID, "", exampleRecipeStepCondition.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step condition ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepConditionExists(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		c, db := buildTestClient(t)
		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCondition.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepConditionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.RecipeStepConditionExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		c, db := buildTestClient(t)
		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCondition.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepConditionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.RecipeStepConditionExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()
		exampleRecipeStepCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCondition.IngredientState.ID}
		exampleRecipeStepCondition.Ingredients = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCondition.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepConditions(false, 0, exampleRecipeStepCondition))

		actual, err := c.GetRecipeStepCondition(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepCondition, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCondition(ctx, "", exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCondition(ctx, exampleRecipeID, "", exampleRecipeStepCondition.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step condition ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCondition(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()
		exampleRecipeStepCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCondition.IngredientState.ID}
		exampleRecipeStepCondition.Ingredients = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCondition.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepCondition(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCondition.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepConditions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepConditionList := fakes.BuildFakeRecipeStepConditionList()

		for i := range exampleRecipeStepConditionList.Data {
			exampleRecipeStepConditionList.Data[i].IngredientState = types.ValidIngredientState{ID: exampleRecipeStepConditionList.Data[i].IngredientState.ID}
			exampleRecipeStepConditionList.Data[i].Ingredients = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_conditions", getRecipeStepConditionsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepConditionsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepConditions(true, exampleRecipeStepConditionList.FilteredCount, exampleRecipeStepConditionList.Data...))

		actual, err := c.GetRecipeStepConditions(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepConditionList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepConditions(ctx, "", exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepConditions(ctx, exampleRecipeID, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepConditionList := fakes.BuildFakeRecipeStepConditionList()
		exampleRecipeStepConditionList.Page = 0
		exampleRecipeStepConditionList.Limit = 0

		for i := range exampleRecipeStepConditionList.Data {
			exampleRecipeStepConditionList.Data[i].IngredientState = types.ValidIngredientState{ID: exampleRecipeStepConditionList.Data[i].IngredientState.ID}
			exampleRecipeStepConditionList.Data[i].Ingredients = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_conditions", getRecipeStepConditionsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepConditionsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepConditions(true, exampleRecipeStepConditionList.FilteredCount, exampleRecipeStepConditionList.Data...))

		actual, err := c.GetRecipeStepConditions(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepConditionList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_conditions", getRecipeStepConditionsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepConditionsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepConditions(ctx, exampleRecipeID, exampleRecipeStepID, filter)
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

		query, args := c.buildListQuery(ctx, "recipe_step_conditions", getRecipeStepConditionsJoins, []string{"valid_measurement_units.id", "valid_ingredients.id"}, nil, householdOwnershipColumn, recipeStepConditionsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepConditions(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipeStepCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()
		exampleRecipeStepCondition.ID = "1"
		exampleRecipeStepCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCondition.IngredientState.ID}
		exampleRecipeStepCondition.Ingredients = nil
		exampleInput := converters.ConvertRecipeStepConditionToRecipeStepConditionDatabaseCreationInput(exampleRecipeStepCondition)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.BelongsToRecipeStep,
			exampleInput.IngredientStateID,
			exampleInput.Optional,
			exampleInput.Notes,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepConditionCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleRecipeStepCondition.CreatedAt
		}

		actual, err := c.CreateRecipeStepCondition(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepCondition, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepCondition(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()
		exampleRecipeStepCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCondition.IngredientState.ID}
		exampleRecipeStepCondition.Ingredients = nil
		exampleInput := converters.ConvertRecipeStepConditionToRecipeStepConditionDatabaseCreationInput(exampleRecipeStepCondition)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.BelongsToRecipeStep,
			exampleInput.IngredientStateID,
			exampleInput.Optional,
			exampleInput.Notes,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepConditionCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleRecipeStepCondition.CreatedAt
		}

		actual, err := c.CreateRecipeStepCondition(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_createRecipeStepCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()
		exampleRecipeStepCondition.ID = "3"
		exampleRecipeStepCondition.BelongsToRecipeStep = "2"
		exampleRecipeStepCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCondition.IngredientState.ID}
		exampleRecipeStepCondition.Ingredients = nil

		exampleInput := converters.ConvertRecipeStepConditionToRecipeStepConditionDatabaseCreationInput(exampleRecipeStepCondition)

		ctx := context.Background()
		c, db := buildTestClient(t)

		recipeStepConditionCreationArgs := []any{
			exampleInput.ID,
			exampleInput.BelongsToRecipeStep,
			exampleInput.IngredientStateID,
			exampleInput.Optional,
			exampleInput.Notes,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepConditionCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeStepConditionCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleRecipeStepCondition.CreatedAt
		}

		actual, err := c.createRecipeStepCondition(ctx, c.db, exampleInput)
		assert.NoError(t, err)
		assert.NotNil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateRecipeStepCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()
		exampleRecipeStepCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCondition.IngredientState.ID}
		exampleRecipeStepCondition.Ingredients = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepCondition.Optional,
			exampleRecipeStepCondition.Notes,
			exampleRecipeStepCondition.BelongsToRecipeStep,
			exampleRecipeStepCondition.IngredientState.ID,
			exampleRecipeStepCondition.BelongsToRecipeStep,
			exampleRecipeStepCondition.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeStepConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateRecipeStepCondition(ctx, exampleRecipeStepCondition))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepCondition(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()
		exampleRecipeStepCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCondition.IngredientState.ID}
		exampleRecipeStepCondition.Ingredients = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepCondition.Optional,
			exampleRecipeStepCondition.Notes,
			exampleRecipeStepCondition.BelongsToRecipeStep,
			exampleRecipeStepCondition.IngredientState.ID,
			exampleRecipeStepCondition.BelongsToRecipeStep,
			exampleRecipeStepCondition.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeStepConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateRecipeStepCondition(ctx, exampleRecipeStepCondition))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveRecipeStepCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCondition.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeStepConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveRecipeStepCondition(ctx, exampleRecipeStepID, exampleRecipeStepCondition.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepCondition(ctx, "", exampleRecipeStepCondition.ID))
	})

	T.Run("with invalid recipe step condition ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepCondition(ctx, exampleRecipeStepID, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCondition := fakes.BuildFakeRecipeStepCondition()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCondition.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeStepConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveRecipeStepCondition(ctx, exampleRecipeStepID, exampleRecipeStepCondition.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
