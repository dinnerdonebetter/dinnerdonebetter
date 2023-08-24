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

func buildMockRowsFromMealPlans(includeCounts bool, filteredCount uint64, mealPlans ...*types.MealPlan) *sqlmock.Rows {
	columns := mealPlansTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range mealPlans {
		rowValues := []driver.Value{
			x.ID,
			x.Notes,
			x.Status,
			x.VotingDeadline,
			x.GroceryListInitialized,
			x.TasksCreated,
			x.ElectionMethod,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
			x.BelongsToHousehold,
			x.CreatedByUser,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(mealPlans))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanMealPlans(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanMealPlans(ctx, mockRows, false)
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

		_, _, _, err := q.scanMealPlans(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_MealPlanExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		actual, err := c.MealPlanExists(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		actual, err := c.MealPlanExists(ctx, exampleMealPlanID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlan(ctx, "", exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlan(ctx, exampleMealPlanID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealPlan(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateMealPlan(ctx, nil))
	})
}

func TestQuerier_ArchiveMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleAccountID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlan(ctx, "", exampleAccountID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlan(ctx, exampleMealPlan.ID, ""))
	})
}

func TestQuerier_AttemptToFinalizeCompleteMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleHousehold := fakes.BuildFakeHousehold()
		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.AttemptToFinalizeMealPlan(ctx, "", exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.AttemptToFinalizeMealPlan(ctx, exampleMealPlan.ID, "")
		assert.False(t, actual)
		assert.Error(t, err)
	})
}

func TestQuerier_FetchExpiredAndUnresolvedMealPlanIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := []*types.MealPlan{}
		exampleMealPlanList := fakes.BuildFakeMealPlanList()
		for _, mp := range exampleMealPlanList.Data {
			mp.Events = nil
			expected = append(expected, mp)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getExpiredAndUnresolvedMealPlansQuery)).
			WithArgs().
			WillReturnRows(buildMockRowsFromMealPlans(false, exampleMealPlanList.FilteredCount, exampleMealPlanList.Data...))

		actual, err := c.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error performing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getExpiredAndUnresolvedMealPlansQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanList := fakes.BuildFakeMealPlanList()
		for _, mp := range exampleMealPlanList.Data {
			mp.Events = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getExpiredAndUnresolvedMealPlansQuery)).
			WithArgs().
			WillReturnRows(buildInvalidMockRowsFromListOfIDs([]string{"things", "and", "stuff"}))

		actual, err := c.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error closing rows", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanList := fakes.BuildFakeMealPlanList()
		for _, mp := range exampleMealPlanList.Data {
			mp.Events = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getExpiredAndUnresolvedMealPlansQuery)).
			WithArgs().
			WillReturnRows(buildMockRowsFromMealPlans(false, exampleMealPlanList.FilteredCount, exampleMealPlanList.Data...).RowError(0, errors.New("blah")))

		actual, err := c.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_FetchMissingVotesForMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("with missing meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)
		exampleHousehold := fakes.BuildFakeHousehold()

		actual, err := c.FetchMissingVotesForMealPlan(ctx, "", exampleHousehold.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)
		exampleMealPlan := fakes.BuildFakeMealPlan()

		actual, err := c.FetchMissingVotesForMealPlan(ctx, exampleMealPlan.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}
