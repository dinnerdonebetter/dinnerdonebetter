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

func buildMockRowsFromRecipeRating(includeCounts bool, filteredCount uint64, recipeRatings ...*types.RecipeRating) *sqlmock.Rows {
	columns := recipeRatingsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeRatings {
		rowValues := []driver.Value{
			&x.ID,
			&x.RecipeID,
			&x.Taste,
			&x.Difficulty,
			&x.Cleanup,
			&x.Instructions,
			&x.Overall,
			&x.Notes,
			&x.ByUser,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(recipeRatings))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanRecipeRatings(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipeRatings(ctx, mockRows, false)
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

		_, _, _, err := q.scanRecipeRatings(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_RecipeRatingExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.RecipeRatingExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeRating := fakes.BuildFakeRecipeRating()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeRating.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeRatingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeRating(false, 0, exampleRecipeRating))

		actual, err := c.GetRecipeRating(ctx, exampleRecipeRating.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeRating, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeRating(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeRating := fakes.BuildFakeRecipeRating()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeRating.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeRatingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeRating(ctx, exampleRecipeRating.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeRatings(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeRatingList := fakes.BuildFakeRecipeRatingList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_ratings", nil, nil, nil, "", recipeRatingsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeRating(true, exampleRecipeRatingList.FilteredCount, exampleRecipeRatingList.Data...))

		actual, err := c.GetRecipeRatings(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeRatingList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_ratings", nil, nil, nil, "", recipeRatingsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeRatings(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_ratings", nil, nil, nil, "", recipeRatingsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeRatings(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeRating := fakes.BuildFakeRecipeRating()
		exampleRecipeRating.ID = "1"
		exampleInput := converters.ConvertRecipeRatingToRecipeRatingDatabaseCreationInput(exampleRecipeRating)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.MealID,
			exampleInput.Taste,
			exampleInput.Difficulty,
			exampleInput.Cleanup,
			exampleInput.Instructions,
			exampleInput.Overall,
			exampleInput.Notes,
			exampleInput.ByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeRatingCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleRecipeRating.CreatedAt
		}

		actual, err := c.CreateRecipeRating(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeRating, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeRating(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleRecipeRating := fakes.BuildFakeRecipeRating()
		exampleInput := converters.ConvertRecipeRatingToRecipeRatingDatabaseCreationInput(exampleRecipeRating)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.MealID,
			exampleInput.Taste,
			exampleInput.Difficulty,
			exampleInput.Cleanup,
			exampleInput.Instructions,
			exampleInput.Overall,
			exampleInput.Notes,
			exampleInput.ByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeRatingCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleRecipeRating.CreatedAt
		}

		actual, err := c.CreateRecipeRating(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeRating(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeRating(T *testing.T) {
	T.Parallel()

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeRating(ctx, ""))
	})
}
