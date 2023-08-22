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
	"github.com/stretchr/testify/require"
)

func buildMockRowsFromMealPlanEvents(includeCounts bool, filteredCount uint64, mealPlans ...*types.MealPlanEvent) *sqlmock.Rows {
	columns := mealPlanEventsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range mealPlans {
		rowValues := []driver.Value{
			x.ID,
			x.Notes,
			x.StartsAt,
			x.EndsAt,
			x.MealName,
			x.BelongsToMealPlan,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(mealPlans))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_MealPlanEventExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanEventID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		actual, err := c.MealPlanEventExists(ctx, "", exampleMealPlanEventID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan event ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanEventExists(ctx, exampleMealPlanID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanEvent(ctx, "", exampleMealPlanEventID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid meal plan event ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanEvent(ctx, exampleMealPlan.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_getMealPlanEventsForMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEvents := fakes.BuildFakeMealPlanEventList()
		for i, mealPlanEvent := range exampleMealPlanEvents.Data {
			for j := range mealPlanEvent.Options {
				exampleMealPlanEvents.Data[i].Options[j].Meal = *fakes.BuildFakeMeal()
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealPlanEventsForMealPlanArgs := []any{
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanEventsForMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanEventsForMealPlanArgs)...).
			WillReturnRows(buildMockRowsFromMealPlanEvents(false, 0, exampleMealPlanEvents.Data...))

		for _, mealPlanEvent := range exampleMealPlanEvents.Data {
			prepareMockToSuccessfullyGetMealPlanEvent(t, mealPlanEvent, exampleMealPlanID, db, nil)
		}

		actual, err := c.getMealPlanEventsForMealPlan(ctx, exampleMealPlanID)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanEvents.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.getMealPlanEventsForMealPlan(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error making initial query", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEvents := fakes.BuildFakeMealPlanEventList()
		for i, mealPlanEvent := range exampleMealPlanEvents.Data {
			for j := range mealPlanEvent.Options {
				exampleMealPlanEvents.Data[i].Options[j].Meal = *fakes.BuildFakeMeal()
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealPlanEventsForMealPlanArgs := []any{
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanEventsForMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanEventsForMealPlanArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.getMealPlanEventsForMealPlan(ctx, exampleMealPlanID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEvents := fakes.BuildFakeMealPlanEventList()
		for i, mealPlanEvent := range exampleMealPlanEvents.Data {
			for j := range mealPlanEvent.Options {
				exampleMealPlanEvents.Data[i].Options[j].Meal = *fakes.BuildFakeMeal()
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealPlanEventsForMealPlanArgs := []any{
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanEventsForMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanEventsForMealPlanArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.getMealPlanEventsForMealPlan(ctx, exampleMealPlanID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error fetching meal plan options", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEvents := fakes.BuildFakeMealPlanEventList()
		for i, mealPlanEvent := range exampleMealPlanEvents.Data {
			for j := range mealPlanEvent.Options {
				exampleMealPlanEvents.Data[i].Options[j].Meal = *fakes.BuildFakeMeal()
			}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealPlanEventsForMealPlanArgs := []any{
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanEventsForMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanEventsForMealPlanArgs)...).
			WillReturnRows(buildMockRowsFromMealPlanEvents(false, 0, exampleMealPlanEvents.Data...))

		prepareMockToSuccessfullyGetMealPlanEvent(t, exampleMealPlanEvents.Data[0], exampleMealPlanID, db, errors.New("blah"))

		actual, err := c.getMealPlanEventsForMealPlan(ctx, exampleMealPlanID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlanEvents(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEvents := fakes.BuildFakeMealPlanEventList()
		for i := range exampleMealPlanEvents.Data {
			exampleMealPlanEvents.Data[i].Options = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealPlanEventsQuery, getMealPlanEventsArgs := c.buildListQuery(ctx, "meal_plan_events", nil, nil, nil, "", mealPlanEventsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanEventsQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanEventsArgs)...).
			WillReturnRows(buildMockRowsFromMealPlanEvents(true, exampleMealPlanEvents.FilteredCount, exampleMealPlanEvents.Data...))

		actual, err := c.GetMealPlanEvents(ctx, exampleMealPlanID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanEvents, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEvents := fakes.BuildFakeMealPlanEventList()
		for i := range exampleMealPlanEvents.Data {
			exampleMealPlanEvents.Data[i].Options = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealPlanEventsQuery, getMealPlanEventsArgs := c.buildListQuery(ctx, "meal_plan_events", nil, nil, nil, "", mealPlanEventsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanEventsQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanEventsArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlanEvents(ctx, exampleMealPlanID, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEvents := fakes.BuildFakeMealPlanEventList()
		for i := range exampleMealPlanEvents.Data {
			exampleMealPlanEvents.Data[i].Options = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealPlanEventsQuery, getMealPlanEventsArgs := c.buildListQuery(ctx, "meal_plan_events", nil, nil, nil, "", mealPlanEventsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanEventsQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanEventsArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealPlanEvents(ctx, exampleMealPlanID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_createMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		tx, err := c.db.BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, tx)

		actual, err := c.createMealPlanEvent(ctx, tx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealPlanEvent(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.UpdateMealPlanEvent(ctx, nil)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan event ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanEvent(ctx, "", exampleMealPlan.ID))
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanEvent(ctx, exampleMealPlanEventID, ""))
	})
}
