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

func buildMockRowsFromMealPlanOptionVotes(includeCounts bool, filteredCount uint64, mealPlanOptionVotes ...*types.MealPlanOptionVote) *sqlmock.Rows {
	columns := mealPlanOptionVotesTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range mealPlanOptionVotes {
		rowValues := []driver.Value{
			x.ID,
			x.Rank,
			x.Abstain,
			x.Notes,
			x.ByUser,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
			x.BelongsToMealPlanOption,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(mealPlanOptionVotes))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanMealPlanOptionVotes(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanMealPlanOptionVotes(ctx, mockRows, false)
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

		_, _, _, err := q.scanMealPlanOptionVotes(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_MealPlanOptionVoteExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionVoteExists(ctx, "", exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionVoteExists(ctx, exampleMealPlanID, exampleMealPlanEventID, "", exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionVoteExists(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealPlanOptionID,
			exampleMealPlanOptionVote.ID,
			exampleMealPlanEventID,
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanOptionVoteQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlanOptionVotes(false, 0, exampleMealPlanOptionVote))

		actual, err := c.GetMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOptionVote, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionVote(ctx, "", exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, "", exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealPlanOptionID,
			exampleMealPlanOptionVote.ID,
			exampleMealPlanEventID,
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanOptionVoteQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, exampleMealPlanOptionVote.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlanOptionVotes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()
		exampleMealPlanOptionVoteList := fakes.BuildFakeMealPlanOptionVoteList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plan_option_votes", getMealPlanOptionVotesJoins, nil, nil, householdOwnershipColumn, mealPlanOptionVotesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlanOptionVotes(true, exampleMealPlanOptionVoteList.FilteredCount, exampleMealPlanOptionVoteList.Data...))

		actual, err := c.GetMealPlanOptionVotes(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOptionVoteList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionVotes(ctx, "", exampleMealPlanEventID, exampleMealPlanOptionID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionVotes(ctx, exampleMealPlanID, exampleMealPlanEventID, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plan_option_votes", getMealPlanOptionVotesJoins, nil, nil, householdOwnershipColumn, mealPlanOptionVotesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlanOptionVotes(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plan_option_votes", getMealPlanOptionVotesJoins, nil, nil, householdOwnershipColumn, mealPlanOptionVotesTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealPlanOptionVotes(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealPlanOptionVote(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error creating transaction", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()
		exampleMealPlanOptionVote.ID = "1"
		exampleInput := converters.ConvertMealPlanOptionVoteToMealPlanOptionVoteDatabaseCreationInput(exampleMealPlanOptionVote)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.CreateMealPlanOptionVote(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateMealPlanOptionVote(ctx, nil))
	})
}

func TestQuerier_ArchiveMealPlanOptionVote(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionVote := fakes.BuildFakeMealPlanOptionVote()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, "", exampleMealPlanOptionVote.ID))
	})

	T.Run("with invalid meal plan option vote ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanOptionVote(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOptionID, ""))
	})
}
