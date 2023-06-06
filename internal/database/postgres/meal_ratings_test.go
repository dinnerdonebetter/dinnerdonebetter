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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromMealRating(includeCounts bool, filteredCount uint64, mealRatings ...*types.MealRating) *sqlmock.Rows {
	columns := mealRatingsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range mealRatings {
		rowValues := []driver.Value{
			&x.ID,
			&x.MealID,
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
			rowValues = append(rowValues, filteredCount, len(mealRatings))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanMealRatings(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanMealRatings(ctx, mockRows, false)
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

		_, _, _, err := q.scanMealRatings(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_MealRatingExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealRating := fakes.BuildFakeMealRating()

		c, db := buildTestClient(t)
		args := []any{
			exampleMealRating.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealRatingExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.MealRatingExists(ctx, exampleMealRating.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.MealRatingExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealRating := fakes.BuildFakeMealRating()

		c, db := buildTestClient(t)
		args := []any{
			exampleMealRating.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealRatingExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.MealRatingExists(ctx, exampleMealRating.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealRating := fakes.BuildFakeMealRating()

		c, db := buildTestClient(t)
		args := []any{
			exampleMealRating.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealRatingExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.MealRatingExists(ctx, exampleMealRating.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealRating := fakes.BuildFakeMealRating()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealRating.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealRatingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealRating(false, 0, exampleMealRating))

		actual, err := c.GetMealRating(ctx, exampleMealRating.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealRating, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealRating(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleMealRating := fakes.BuildFakeMealRating()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealRating.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealRatingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealRating(ctx, exampleMealRating.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealRatings(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealRatingList := fakes.BuildFakeMealRatingList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_ratings", nil, nil, nil, "", mealRatingsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealRating(true, exampleMealRatingList.FilteredCount, exampleMealRatingList.Data...))

		actual, err := c.GetMealRatings(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealRatingList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealRatingList := fakes.BuildFakeMealRatingList()
		exampleMealRatingList.Page = 0
		exampleMealRatingList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_ratings", nil, nil, nil, "", mealRatingsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealRating(true, exampleMealRatingList.FilteredCount, exampleMealRatingList.Data...))

		actual, err := c.GetMealRatings(ctx, nil)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealRatingList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_ratings", nil, nil, nil, "", mealRatingsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealRatings(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_ratings", nil, nil, nil, "", mealRatingsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealRatings(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealRating := fakes.BuildFakeMealRating()
		exampleMealRating.ID = "1"
		exampleInput := converters.ConvertMealRatingToMealRatingDatabaseCreationInput(exampleMealRating)

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

		db.ExpectExec(formatQueryForSQLMock(mealRatingCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleMealRating.CreatedAt
		}

		actual, err := c.CreateMealRating(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealRating, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealRating(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleMealRating := fakes.BuildFakeMealRating()
		exampleInput := converters.ConvertMealRatingToMealRatingDatabaseCreationInput(exampleMealRating)

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

		db.ExpectExec(formatQueryForSQLMock(mealRatingCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleMealRating.CreatedAt
		}

		actual, err := c.CreateMealRating(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateMealRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealRating := fakes.BuildFakeMealRating()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealRating.MealID,
			exampleMealRating.Taste,
			exampleMealRating.Difficulty,
			exampleMealRating.Cleanup,
			exampleMealRating.Instructions,
			exampleMealRating.Overall,
			exampleMealRating.Notes,
			exampleMealRating.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealRatingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateMealRating(ctx, exampleMealRating))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateMealRating(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleMealRating := fakes.BuildFakeMealRating()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealRating.MealID,
			exampleMealRating.Taste,
			exampleMealRating.Difficulty,
			exampleMealRating.Cleanup,
			exampleMealRating.Instructions,
			exampleMealRating.Overall,
			exampleMealRating.Notes,
			exampleMealRating.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealRatingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateMealRating(ctx, exampleMealRating))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveMealRating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealRating := fakes.BuildFakeMealRating()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealRating.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealRatingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveMealRating(ctx, exampleMealRating.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealRating(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleMealRating := fakes.BuildFakeMealRating()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealRating.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealRatingQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveMealRating(ctx, exampleMealRating.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
