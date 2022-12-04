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

func buildMockRowsFromRecipeStepCompletionConditions(includeCounts bool, filteredCount uint64, recipeStepCompletionConditions ...*types.RecipeStepCompletionCondition) *sqlmock.Rows {
	columns := recipeStepCompletionConditionsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepCompletionConditions {
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
			rowValues = append(rowValues, filteredCount, len(recipeStepCompletionConditions))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanRecipeStepCompletionConditions(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipeStepCompletionConditions(ctx, mockRows, false)
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

		_, _, _, err := q.scanRecipeStepCompletionConditions(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_RecipeStepCompletionConditionExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		c, db := buildTestClient(t)
		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCompletionCondition.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepCompletionConditionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.RecipeStepCompletionConditionExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepCompletionConditionExists(ctx, "", exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepCompletionConditionExists(ctx, exampleRecipeID, "", exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step completion condition ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepCompletionConditionExists(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		c, db := buildTestClient(t)
		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCompletionCondition.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepCompletionConditionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.RecipeStepCompletionConditionExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		c, db := buildTestClient(t)
		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCompletionCondition.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepCompletionConditionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.RecipeStepCompletionConditionExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()
		exampleRecipeStepCompletionCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCompletionCondition.IngredientState.ID}
		exampleRecipeStepCompletionCondition.Ingredients = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepCompletionCondition.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepCompletionConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepCompletionConditions(false, 0, exampleRecipeStepCompletionCondition))

		actual, err := c.GetRecipeStepCompletionCondition(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepCompletionCondition, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCompletionCondition(ctx, "", exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCompletionCondition(ctx, exampleRecipeID, "", exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step completion condition ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCompletionCondition(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()
		exampleRecipeStepCompletionCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCompletionCondition.IngredientState.ID}
		exampleRecipeStepCompletionCondition.Ingredients = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeStepCompletionCondition.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepCompletionConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepCompletionCondition(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepCompletionConditions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionConditionList := fakes.BuildFakeRecipeStepCompletionConditionList()

		for i := range exampleRecipeStepCompletionConditionList.Data {
			exampleRecipeStepCompletionConditionList.Data[i].IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCompletionConditionList.Data[i].IngredientState.ID}
			exampleRecipeStepCompletionConditionList.Data[i].Ingredients = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_completion_conditions", getRecipeStepCompletionConditionsJoins, []string{}, nil, "", recipeStepCompletionConditionsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepCompletionConditions(true, exampleRecipeStepCompletionConditionList.FilteredCount, exampleRecipeStepCompletionConditionList.Data...))

		actual, err := c.GetRecipeStepCompletionConditions(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepCompletionConditionList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCompletionConditions(ctx, "", exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepCompletionConditions(ctx, exampleRecipeID, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionConditionList := fakes.BuildFakeRecipeStepCompletionConditionList()
		exampleRecipeStepCompletionConditionList.Page = 0
		exampleRecipeStepCompletionConditionList.Limit = 0

		for i := range exampleRecipeStepCompletionConditionList.Data {
			exampleRecipeStepCompletionConditionList.Data[i].IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCompletionConditionList.Data[i].IngredientState.ID}
			exampleRecipeStepCompletionConditionList.Data[i].Ingredients = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_completion_conditions", getRecipeStepCompletionConditionsJoins, []string{}, nil, "", recipeStepCompletionConditionsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepCompletionConditions(true, exampleRecipeStepCompletionConditionList.FilteredCount, exampleRecipeStepCompletionConditionList.Data...))

		actual, err := c.GetRecipeStepCompletionConditions(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepCompletionConditionList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_completion_conditions", getRecipeStepCompletionConditionsJoins, []string{}, nil, "", recipeStepCompletionConditionsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepCompletionConditions(ctx, exampleRecipeID, exampleRecipeStepID, filter)
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

		query, args := c.buildListQuery(ctx, "recipe_step_completion_conditions", getRecipeStepCompletionConditionsJoins, []string{}, nil, "", recipeStepCompletionConditionsTableColumns, "", false, filter)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepCompletionConditions(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()
		exampleRecipeStepCompletionCondition.ID = "1"
		exampleRecipeStepCompletionCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCompletionCondition.IngredientState.ID}
		exampleRecipeStepCompletionCondition.Ingredients = nil
		exampleInput := converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionDatabaseCreationInput(exampleRecipeStepCompletionCondition)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.BelongsToRecipeStep,
			exampleInput.IngredientStateID,
			exampleInput.Optional,
			exampleInput.Notes,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepCompletionConditionCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleRecipeStepCompletionCondition.CreatedAt
		}

		actual, err := c.CreateRecipeStepCompletionCondition(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepCompletionCondition, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepCompletionCondition(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()
		exampleRecipeStepCompletionCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCompletionCondition.IngredientState.ID}
		exampleRecipeStepCompletionCondition.Ingredients = nil
		exampleInput := converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionDatabaseCreationInput(exampleRecipeStepCompletionCondition)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.BelongsToRecipeStep,
			exampleInput.IngredientStateID,
			exampleInput.Optional,
			exampleInput.Notes,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepCompletionConditionCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleRecipeStepCompletionCondition.CreatedAt
		}

		actual, err := c.CreateRecipeStepCompletionCondition(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestSQLQuerier_createRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()
		exampleRecipeStepCompletionCondition.ID = "3"
		exampleRecipeStepCompletionCondition.BelongsToRecipeStep = "2"
		exampleRecipeStepCompletionCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCompletionCondition.IngredientState.ID}
		exampleRecipeStepCompletionCondition.Ingredients = nil

		exampleInput := converters.ConvertRecipeStepCompletionConditionToRecipeStepCompletionConditionDatabaseCreationInput(exampleRecipeStepCompletionCondition)

		ctx := context.Background()
		c, db := buildTestClient(t)

		recipeStepCompletionConditionCreationArgs := []any{
			exampleInput.ID,
			exampleInput.BelongsToRecipeStep,
			exampleInput.IngredientStateID,
			exampleInput.Optional,
			exampleInput.Notes,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepCompletionConditionCreationQuery)).
			WithArgs(interfaceToDriverValue(recipeStepCompletionConditionCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleRecipeStepCompletionCondition.CreatedAt
		}

		actual, err := c.createRecipeStepCompletionCondition(ctx, c.db, exampleInput)
		assert.NoError(t, err)
		assert.NotNil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()
		exampleRecipeStepCompletionCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCompletionCondition.IngredientState.ID}
		exampleRecipeStepCompletionCondition.Ingredients = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepCompletionCondition.Optional,
			exampleRecipeStepCompletionCondition.Notes,
			exampleRecipeStepCompletionCondition.BelongsToRecipeStep,
			exampleRecipeStepCompletionCondition.IngredientState.ID,
			exampleRecipeStepCompletionCondition.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeStepCompletionConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateRecipeStepCompletionCondition(ctx, exampleRecipeStepCompletionCondition))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepCompletionCondition(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()
		exampleRecipeStepCompletionCondition.IngredientState = types.ValidIngredientState{ID: exampleRecipeStepCompletionCondition.IngredientState.ID}
		exampleRecipeStepCompletionCondition.Ingredients = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepCompletionCondition.Optional,
			exampleRecipeStepCompletionCondition.Notes,
			exampleRecipeStepCompletionCondition.BelongsToRecipeStep,
			exampleRecipeStepCompletionCondition.IngredientState.ID,
			exampleRecipeStepCompletionCondition.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeStepCompletionConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateRecipeStepCompletionCondition(ctx, exampleRecipeStepCompletionCondition))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCompletionCondition.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeStepCompletionConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveRecipeStepCompletionCondition(ctx, exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepCompletionCondition(ctx, "", exampleRecipeStepCompletionCondition.ID))
	})

	T.Run("with invalid recipe step completion condition ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepCompletionCondition(ctx, exampleRecipeStepID, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepCompletionCondition := fakes.BuildFakeRecipeStepCompletionCondition()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepCompletionCondition.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeStepCompletionConditionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveRecipeStepCompletionCondition(ctx, exampleRecipeStepID, exampleRecipeStepCompletionCondition.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
