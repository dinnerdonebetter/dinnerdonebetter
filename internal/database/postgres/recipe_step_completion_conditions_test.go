package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromRecipeStepCompletionConditions(includeCounts bool, filteredCount uint64, recipeStepCompletionConditions ...*types.RecipeStepCompletionCondition) *sqlmock.Rows {
	columns := recipeStepCompletionConditionsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepCompletionConditions {
		for _, y := range x.Ingredients {
			rowValues := []driver.Value{
				&y.ID,
				&y.BelongsToRecipeStepCompletionCondition,
				&y.RecipeStepIngredient,
				&x.ID,
				&x.BelongsToRecipeStep,
				&x.IngredientState.ID,
				&x.IngredientState.Name,
				&x.IngredientState.Description,
				&x.IngredientState.IconPath,
				&x.IngredientState.Slug,
				&x.IngredientState.PastTense,
				&x.IngredientState.AttributeType,
				&x.IngredientState.CreatedAt,
				&x.IngredientState.LastUpdatedAt,
				&x.IngredientState.ArchivedAt,
				&x.Optional,
				&x.Notes,
				&x.CreatedAt,
				&x.LastUpdatedAt,
				&x.ArchivedAt,
			}

			if includeCounts {
				rowValues = append(rowValues, filteredCount, len(recipeStepCompletionConditions))
			}

			exampleRows.AddRow(rowValues...)
		}
	}

	return exampleRows
}

func TestQuerier_RecipeStepCompletionConditionExists(T *testing.T) {
	T.Parallel()

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
}

func TestQuerier_GetRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

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
			exampleRecipeStepCompletionConditionList.Data[i].BelongsToRecipeStep = exampleRecipeStepID
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			filter.Limit,
			filter.QueryOffset(),
		}

		returnedRows := buildMockRowsFromRecipeStepCompletionConditions(true, exampleRecipeStepCompletionConditionList.FilteredCount, exampleRecipeStepCompletionConditionList.Data...)
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepCompletionConditionsQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(returnedRows)

		actual, err := c.GetRecipeStepCompletionConditions(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)

		for i, actualEntry := range actual.Data {
			x, y := exampleRecipeStepCompletionConditionList.Data[i], actualEntry
			if !assert.Equal(t, x, y) {
				t.Log("it's happening")
			}
		}

		if !assert.Equal(t, exampleRecipeStepCompletionConditionList, actual) {
			t.Log("it's happening")
		}

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

		f := types.DefaultQueryFilter()
		exampleRecipeStepCompletionConditionList.Page = *f.Page
		exampleRecipeStepCompletionConditionList.Limit = *f.Limit

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			f.CreatedAfter,
			f.CreatedBefore,
			f.UpdatedAfter,
			f.UpdatedBefore,
			f.Limit,
			f.QueryOffset(),
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepCompletionConditionsQuery)).
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

		args := []any{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			filter.Limit,
			filter.QueryOffset(),
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepCompletionConditionsQuery)).
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

		args := []any{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			filter.Limit,
			filter.QueryOffset(),
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepCompletionConditionsQuery)).
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

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepCompletionCondition(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestSQLQuerier_createRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.createRecipeStepCompletionCondition(ctx, c.db, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepCompletionCondition(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeStepCompletionCondition(T *testing.T) {
	T.Parallel()

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
}
