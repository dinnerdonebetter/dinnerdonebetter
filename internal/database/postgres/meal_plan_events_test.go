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
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
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

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		mealPlanEventExistenceArgs := []any{
			exampleMealPlanEventID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanEventExistenceQuery)).
			WithArgs(interfaceToDriverValue(mealPlanEventExistenceArgs)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.MealPlanEventExists(ctx, exampleMealPlanID, exampleMealPlanEventID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		mealPlanEventExistenceArgs := []any{
			exampleMealPlanEventID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanEventExistenceQuery)).
			WithArgs(interfaceToDriverValue(mealPlanEventExistenceArgs)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.MealPlanEventExists(ctx, exampleMealPlanID, exampleMealPlanEventID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		c, db := buildTestClient(t)
		mealPlanEventExistenceArgs := []any{
			exampleMealPlanEventID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanEventExistenceQuery)).
			WithArgs(interfaceToDriverValue(mealPlanEventExistenceArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.MealPlanEventExists(ctx, exampleMealPlanID, exampleMealPlanEventID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		exampleMealPlanEvent.Options = nil

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealPlanEventByIDArgs := []any{
			exampleMealPlanEvent.ID,
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanEventByIDQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanEventByIDArgs)...).
			WillReturnRows(buildMockRowsFromMealPlanEvents(false, 0, exampleMealPlanEvent))

		actual, err := c.GetMealPlanEvent(ctx, exampleMealPlanID, exampleMealPlanEvent.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanEvent, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()

		ctx := context.Background()
		c, db := buildTestClient(t)

		getMealPlanEventByIDArgs := []any{
			exampleMealPlanEvent.ID,
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanEventByIDQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanEventByIDArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealPlanEvent(ctx, exampleMealPlanID, exampleMealPlanEvent.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
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
			prepareMockToSuccessfullyGetMealPlanEvent(t, mealPlanEvent, exampleMealPlanID, db)
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
}

func TestQuerier_createMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		for i := range exampleMealPlanEvent.Options {
			exampleMealPlanEvent.Options[i].CreatedAt = exampleMealPlanEvent.CreatedAt
			exampleMealPlanEvent.Options[i].Meal = types.Meal{ID: exampleMealPlanEvent.Options[i].Meal.ID}
			exampleMealPlanEvent.Options[i].Votes = []*types.MealPlanOptionVote{}
		}
		exampleInput := converters.ConvertMealPlanEventToMealPlanEventDatabaseCreationInput(exampleMealPlanEvent)

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleMealPlanEvent.CreatedAt
		}

		db.ExpectBegin()
		tx, err := c.db.BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, tx)

		mealPlanEventCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.StartsAt,
			exampleInput.EndsAt,
			exampleInput.MealName,
			exampleInput.BelongsToMealPlan,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanEventCreationQuery)).
			WithArgs(interfaceToDriverValue(mealPlanEventCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, option := range exampleInput.Options {
			mealPlanOptionCreationArgs := []any{
				option.ID,
				option.AssignedCook,
				option.AssignedDishwasher,
				option.MealID,
				option.Notes,
				option.BelongsToMealPlanEvent,
				false,
			}

			db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
				WithArgs(interfaceToDriverValue(mealPlanOptionCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		actual, err := c.createMealPlanEvent(ctx, tx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanEvent, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with error creating option", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		for i := range exampleMealPlanEvent.Options {
			exampleMealPlanEvent.Options[i].CreatedAt = exampleMealPlanEvent.CreatedAt
			exampleMealPlanEvent.Options[i].Meal = types.Meal{ID: exampleMealPlanEvent.Options[i].Meal.ID}
			exampleMealPlanEvent.Options[i].Votes = []*types.MealPlanOptionVote{}
		}
		exampleInput := converters.ConvertMealPlanEventToMealPlanEventDatabaseCreationInput(exampleMealPlanEvent)

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleMealPlanEvent.CreatedAt
		}

		db.ExpectBegin()
		tx, err := c.db.BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, tx)

		mealPlanEventCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.StartsAt,
			exampleInput.EndsAt,
			exampleInput.MealName,
			exampleInput.BelongsToMealPlan,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanEventCreationQuery)).
			WithArgs(interfaceToDriverValue(mealPlanEventCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		mealPlanOptionCreationArgs := []any{
			exampleInput.Options[0].ID,
			exampleInput.Options[0].AssignedCook,
			exampleInput.Options[0].AssignedDishwasher,
			exampleInput.Options[0].MealID,
			exampleInput.Options[0].Notes,
			exampleInput.Options[0].BelongsToMealPlanEvent,
			false,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
			WithArgs(interfaceToDriverValue(mealPlanOptionCreationArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.createMealPlanEvent(ctx, tx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		for i := range exampleMealPlanEvent.Options {
			exampleMealPlanEvent.Options[i].CreatedAt = exampleMealPlanEvent.CreatedAt
			exampleMealPlanEvent.Options[i].Meal = types.Meal{ID: exampleMealPlanEvent.Options[i].Meal.ID}
			exampleMealPlanEvent.Options[i].Votes = []*types.MealPlanOptionVote{}
		}
		exampleInput := converters.ConvertMealPlanEventToMealPlanEventDatabaseCreationInput(exampleMealPlanEvent)

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleMealPlanEvent.CreatedAt
		}

		db.ExpectBegin()

		mealPlanEventCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.StartsAt,
			exampleInput.EndsAt,
			exampleInput.MealName,
			exampleInput.BelongsToMealPlan,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanEventCreationQuery)).
			WithArgs(interfaceToDriverValue(mealPlanEventCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, option := range exampleInput.Options {
			mealPlanOptionCreationArgs := []any{
				option.ID,
				option.AssignedCook,
				option.AssignedDishwasher,
				option.MealID,
				option.Notes,
				option.BelongsToMealPlanEvent,
				false,
			}

			db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
				WithArgs(interfaceToDriverValue(mealPlanOptionCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		db.ExpectCommit()

		actual, err := c.CreateMealPlanEvent(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanEvent, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealPlanEvent(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		for i := range exampleMealPlanEvent.Options {
			exampleMealPlanEvent.Options[i].CreatedAt = exampleMealPlanEvent.CreatedAt
			exampleMealPlanEvent.Options[i].Meal = types.Meal{ID: exampleMealPlanEvent.Options[i].Meal.ID}
			exampleMealPlanEvent.Options[i].Votes = []*types.MealPlanOptionVote{}
		}
		exampleInput := converters.ConvertMealPlanEventToMealPlanEventDatabaseCreationInput(exampleMealPlanEvent)

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleMealPlanEvent.CreatedAt
		}

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.CreateMealPlanEvent(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		for i := range exampleMealPlanEvent.Options {
			exampleMealPlanEvent.Options[i].CreatedAt = exampleMealPlanEvent.CreatedAt
			exampleMealPlanEvent.Options[i].Meal = types.Meal{ID: exampleMealPlanEvent.Options[i].Meal.ID}
			exampleMealPlanEvent.Options[i].Votes = []*types.MealPlanOptionVote{}
		}
		exampleInput := converters.ConvertMealPlanEventToMealPlanEventDatabaseCreationInput(exampleMealPlanEvent)

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleMealPlanEvent.CreatedAt
		}

		db.ExpectBegin()

		mealPlanEventCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.StartsAt,
			exampleInput.EndsAt,
			exampleInput.MealName,
			exampleInput.BelongsToMealPlan,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanEventCreationQuery)).
			WithArgs(interfaceToDriverValue(mealPlanEventCreationArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.CreateMealPlanEvent(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		for i := range exampleMealPlanEvent.Options {
			exampleMealPlanEvent.Options[i].CreatedAt = exampleMealPlanEvent.CreatedAt
			exampleMealPlanEvent.Options[i].Meal = types.Meal{ID: exampleMealPlanEvent.Options[i].Meal.ID}
			exampleMealPlanEvent.Options[i].Votes = []*types.MealPlanOptionVote{}
		}
		exampleInput := converters.ConvertMealPlanEventToMealPlanEventDatabaseCreationInput(exampleMealPlanEvent)

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleMealPlanEvent.CreatedAt
		}

		db.ExpectBegin()

		mealPlanEventCreationArgs := []any{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.StartsAt,
			exampleInput.EndsAt,
			exampleInput.MealName,
			exampleInput.BelongsToMealPlan,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanEventCreationQuery)).
			WithArgs(interfaceToDriverValue(mealPlanEventCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, option := range exampleInput.Options {
			mealPlanOptionCreationArgs := []any{
				option.ID,
				option.AssignedCook,
				option.AssignedDishwasher,
				option.MealID,
				option.Notes,
				option.BelongsToMealPlanEvent,
				false,
			}

			db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
				WithArgs(interfaceToDriverValue(mealPlanOptionCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		actual, err := c.CreateMealPlanEvent(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		for i := range exampleMealPlanEvent.Options {
			exampleMealPlanEvent.Options[i].CreatedAt = exampleMealPlanEvent.CreatedAt
			exampleMealPlanEvent.Options[i].Meal = types.Meal{ID: exampleMealPlanEvent.Options[i].Meal.ID}
			exampleMealPlanEvent.Options[i].Votes = []*types.MealPlanOptionVote{}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleMealPlanEvent.CreatedAt
		}

		updateMealPlanEventArgs := []any{
			exampleMealPlanEvent.Notes,
			exampleMealPlanEvent.StartsAt,
			exampleMealPlanEvent.EndsAt,
			exampleMealPlanEvent.MealName,
			exampleMealPlanEvent.BelongsToMealPlan,
			exampleMealPlanEvent.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealPlanEventQuery)).
			WithArgs(interfaceToDriverValue(updateMealPlanEventArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		err := c.UpdateMealPlanEvent(ctx, exampleMealPlanEvent)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.UpdateMealPlanEvent(ctx, nil)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
		for i := range exampleMealPlanEvent.Options {
			exampleMealPlanEvent.Options[i].CreatedAt = exampleMealPlanEvent.CreatedAt
			exampleMealPlanEvent.Options[i].Meal = types.Meal{ID: exampleMealPlanEvent.Options[i].Meal.ID}
			exampleMealPlanEvent.Options[i].Votes = []*types.MealPlanOptionVote{}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		c.timeFunc = func() time.Time {
			return exampleMealPlanEvent.CreatedAt
		}

		updateMealPlanEventArgs := []any{
			exampleMealPlanEvent.Notes,
			exampleMealPlanEvent.StartsAt,
			exampleMealPlanEvent.EndsAt,
			exampleMealPlanEvent.MealName,
			exampleMealPlanEvent.BelongsToMealPlan,
			exampleMealPlanEvent.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealPlanEventQuery)).
			WithArgs(interfaceToDriverValue(updateMealPlanEventArgs)...).
			WillReturnError(errors.New("blah"))

		err := c.UpdateMealPlanEvent(ctx, exampleMealPlanEvent)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleMealPlan.ID, exampleMealPlanEventID}

		db.ExpectExec(formatQueryForSQLMock(archiveMealPlanEventQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		actual := c.ArchiveMealPlanEvent(ctx, exampleMealPlanEventID, exampleMealPlan.ID)
		assert.NoError(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{exampleMealPlan.ID, exampleMealPlanEventID}

		db.ExpectExec(formatQueryForSQLMock(archiveMealPlanEventQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveMealPlanEvent(ctx, exampleMealPlanEventID, exampleMealPlan.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
