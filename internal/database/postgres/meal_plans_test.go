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

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		c, db := buildTestClient(t)
		args := []any{
			exampleMealPlan.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.MealPlanExists(ctx, exampleMealPlan.ID, exampleHouseholdID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		c, db := buildTestClient(t)
		args := []any{
			exampleMealPlan.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.MealPlanExists(ctx, exampleMealPlan.ID, exampleHouseholdID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		c, db := buildTestClient(t)
		args := []any{
			exampleMealPlan.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.MealPlanExists(ctx, exampleMealPlan.ID, exampleHouseholdID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func prepareMockToSuccessfullyGetMealPlan(t *testing.T, exampleMealPlan *types.MealPlan, householdID string, db *sqlmockExpecterWrapper, restrictToPastVotingDeadline bool) {
	t.Helper()

	if exampleMealPlan == nil {
		exampleMealPlan = fakes.BuildFakeMealPlan()
	}

	args := []any{
		exampleMealPlan.ID,
		householdID,
	}

	query := getMealPlanQuery
	if restrictToPastVotingDeadline {
		query = getMealPlanPastVotingDeadlineQuery
	}

	db.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(interfaceToDriverValue(args)...).
		WillReturnRows(buildMockRowsFromMealPlans(false, 0, exampleMealPlan))

	getMealPlanEventForMealPlanArgs := []any{
		exampleMealPlan.ID,
	}

	db.ExpectQuery(formatQueryForSQLMock(getMealPlanEventsForMealPlanQuery)).
		WithArgs(interfaceToDriverValue(getMealPlanEventForMealPlanArgs)...).
		WillReturnRows(buildMockRowsFromMealPlanEvents(false, 0, exampleMealPlan.Events...))

	for _, evt := range exampleMealPlan.Events {
		prepareMockToSuccessfullyGetMealPlanEvent(t, evt, exampleMealPlan.ID, db, nil)
	}
}

func prepareMockToSuccessfullyGetMealPlanEvent(t *testing.T, evt *types.MealPlanEvent, exampleMealPlanID string, db *sqlmockExpecterWrapper, err error) {
	t.Helper()

	getMealPlanOptionsForMealPlanEventsArgs := []any{
		evt.ID,
		exampleMealPlanID,
	}

	if err == nil {
		db.ExpectQuery(formatQueryForSQLMock(getMealPlanOptionsForMealPlanEventQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanOptionsForMealPlanEventsArgs)...).
			WillReturnRows(buildMockRowsFromMealPlanOptions(false, 0, evt.Options...))

		for _, opt := range evt.Options {
			prepareMockToSuccessfullyGetMealPlanOption(t, evt, exampleMealPlanID, opt, db)
		}
	} else {
		db.ExpectQuery(formatQueryForSQLMock(getMealPlanOptionsForMealPlanEventQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanOptionsForMealPlanEventsArgs)...).
			WillReturnError(err)
	}
}

func prepareMockToSuccessfullyGetMealPlanOption(t *testing.T, evt *types.MealPlanEvent, exampleMealPlanID string, opt *types.MealPlanOption, db *sqlmockExpecterWrapper) {
	t.Helper()

	getMealPlanOptionVotesForMealPlanOptionArgs := []any{
		exampleMealPlanID,
		evt.ID,
		opt.ID,
	}

	db.ExpectQuery(formatQueryForSQLMock(getMealPlanOptionVotesForMealPlanOptionQuery)).
		WithArgs(interfaceToDriverValue(getMealPlanOptionVotesForMealPlanOptionArgs)...).
		WillReturnRows(buildMockRowsFromMealPlanOptionVotes(false, 0, opt.Votes...))

	prepareMockToSuccessfullyGetMeal(t, &opt.Meal, db)
}

func TestQuerier_GetMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.Events = []*types.MealPlanEvent{exampleMealPlan.Events[0]}

		for i := range exampleMealPlan.Events[0].Options {
			exampleMealPlan.Events[0].Options[i].Meal = *fakes.BuildFakeMeal()
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		prepareMockToSuccessfullyGetMealPlan(t, exampleMealPlan, exampleHouseholdID, db, false)

		actual, err := c.GetMealPlan(ctx, exampleMealPlan.ID, exampleHouseholdID)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlan, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealPlan.ID,
			exampleHouseholdID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlan(ctx, exampleMealPlan.ID, exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealPlan.ID,
			exampleHouseholdID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealPlan(ctx, exampleMealPlan.ID, exampleHouseholdID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlans(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()
		exampleMealPlanList := fakes.BuildFakeMealPlanList()
		for i := range exampleMealPlanList.Data {
			exampleMealPlanList.Data[i].Events = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, exampleHouseholdID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlans(true, exampleMealPlanList.FilteredCount, exampleMealPlanList.Data...))

		for _, mp := range exampleMealPlanList.Data {
			prepareMockToSuccessfullyGetMealPlan(t, mp, exampleHouseholdID, db, false)
		}

		actual, err := c.GetMealPlans(ctx, exampleHouseholdID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		filter := (*types.QueryFilter)(nil)
		exampleMealPlanList := fakes.BuildFakeMealPlanList()
		exampleMealPlanList.Page = 0
		exampleMealPlanList.Limit = 0
		for i := range exampleMealPlanList.Data {
			exampleMealPlanList.Data[i].Events = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, exampleHouseholdID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlans(true, exampleMealPlanList.FilteredCount, exampleMealPlanList.Data...))

		for _, mp := range exampleMealPlanList.Data {
			prepareMockToSuccessfullyGetMealPlan(t, mp, exampleHouseholdID, db, false)
		}

		actual, err := c.GetMealPlans(ctx, exampleHouseholdID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, exampleHouseholdID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlans(ctx, exampleHouseholdID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, exampleHouseholdID, false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealPlans(ctx, exampleHouseholdID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.ID = "1"
		for i := range exampleMealPlan.Events {
			exampleMealPlan.Events[i].ID = "2"
			exampleMealPlan.Events[i].BelongsToMealPlan = exampleMealPlan.ID
			exampleMealPlan.Events[i].CreatedAt = exampleMealPlan.CreatedAt
			for j := range exampleMealPlan.Events[i].Options {
				exampleMealPlan.Events[i].Options[j].ID = "2"
				exampleMealPlan.Events[i].Options[j].Meal = types.Meal{ID: exampleMealPlan.Events[i].Options[j].Meal.ID}
				exampleMealPlan.Events[i].Options[j].Votes = []*types.MealPlanOptionVote{}
				exampleMealPlan.Events[i].Options[j].BelongsToMealPlanEvent = exampleMealPlan.Events[i].ID
				exampleMealPlan.Events[i].Options[j].CreatedAt = exampleMealPlan.CreatedAt
			}
		}
		exampleInput := converters.ConvertMealPlanToMealPlanDatabaseCreationInput(exampleMealPlan)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []any{
			exampleInput.ID,
			exampleInput.Notes,
			types.MealPlanStatusAwaitingVotes,
			exampleInput.VotingDeadline,
			exampleInput.BelongsToHousehold,
			exampleInput.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, event := range exampleInput.Events {
			mealPlanEventCreationArgs := []any{
				event.ID,
				event.Notes,
				event.StartsAt,
				event.EndsAt,
				event.MealName,
				event.BelongsToMealPlan,
			}

			db.ExpectExec(formatQueryForSQLMock(mealPlanEventCreationQuery)).
				WithArgs(interfaceToDriverValue(mealPlanEventCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())

			for _, option := range event.Options {
				mealPlanOptionCreationArgs := []any{
					option.ID,
					option.AssignedCook,
					option.AssignedDishwasher,
					option.MealID,
					option.Notes,
					option.MealScale,
					option.BelongsToMealPlanEvent,
					len(event.Options) == 1,
				}

				db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
					WithArgs(interfaceToDriverValue(mealPlanOptionCreationArgs)...).
					WillReturnResult(newArbitraryDatabaseResult())
			}
		}

		db.ExpectCommit()

		c.timeFunc = func() time.Time {
			return exampleMealPlan.CreatedAt
		}

		actual, err := c.CreateMealPlan(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlan, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealPlan(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.ID = "1"
		for i := range exampleMealPlan.Events {
			exampleMealPlan.Events[i].ID = "2"
			exampleMealPlan.Events[i].BelongsToMealPlan = "1"
			exampleMealPlan.Events[i].CreatedAt = exampleMealPlan.CreatedAt
		}
		exampleInput := converters.ConvertMealPlanToMealPlanDatabaseCreationInput(exampleMealPlan)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.CreateMealPlan(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleMealPlan := fakes.BuildFakeMealPlan()
		for i := range exampleMealPlan.Events {
			exampleMealPlan.Events[i].ID = "2"
			exampleMealPlan.Events[i].BelongsToMealPlan = "1"
			exampleMealPlan.Events[i].CreatedAt = exampleMealPlan.CreatedAt
		}
		exampleInput := converters.ConvertMealPlanToMealPlanDatabaseCreationInput(exampleMealPlan)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Notes,
			types.MealPlanStatusAwaitingVotes,
			exampleInput.VotingDeadline,
			exampleInput.BelongsToHousehold,
			exampleInput.CreatedByUser,
		}

		db.ExpectBegin()

		db.ExpectExec(formatQueryForSQLMock(mealPlanCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		db.ExpectRollback()

		c.timeFunc = func() time.Time {
			return exampleMealPlan.CreatedAt
		}

		actual, err := c.CreateMealPlan(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.ID = "1"
		for i := range exampleMealPlan.Events {
			exampleMealPlan.Events[i].ID = "2"
			exampleMealPlan.Events[i].BelongsToMealPlan = exampleMealPlan.ID
			exampleMealPlan.Events[i].CreatedAt = exampleMealPlan.CreatedAt
			for j := range exampleMealPlan.Events[i].Options {
				exampleMealPlan.Events[i].Options[j].ID = "2"
				exampleMealPlan.Events[i].Options[j].Meal = types.Meal{ID: exampleMealPlan.Events[i].Options[j].Meal.ID}
				exampleMealPlan.Events[i].Options[j].Votes = []*types.MealPlanOptionVote{}
				exampleMealPlan.Events[i].Options[j].BelongsToMealPlanEvent = exampleMealPlan.Events[i].ID
				exampleMealPlan.Events[i].Options[j].CreatedAt = exampleMealPlan.CreatedAt
			}
		}
		exampleInput := converters.ConvertMealPlanToMealPlanDatabaseCreationInput(exampleMealPlan)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []any{
			exampleInput.ID,
			exampleInput.Notes,
			types.MealPlanStatusAwaitingVotes,
			exampleInput.VotingDeadline,
			exampleInput.BelongsToHousehold,
			exampleInput.CreatedByUser,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, event := range exampleInput.Events {
			mealPlanEventCreationArgs := []any{
				event.ID,
				event.Notes,
				event.StartsAt,
				event.EndsAt,
				event.MealName,
				event.BelongsToMealPlan,
			}

			db.ExpectExec(formatQueryForSQLMock(mealPlanEventCreationQuery)).
				WithArgs(interfaceToDriverValue(mealPlanEventCreationArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())

			for _, option := range event.Options {
				mealPlanOptionCreationArgs := []any{
					option.ID,
					option.AssignedCook,
					option.AssignedDishwasher,
					option.MealID,
					option.Notes,
					option.MealScale,
					option.BelongsToMealPlanEvent,
					len(event.Options) == 1,
				}

				db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
					WithArgs(interfaceToDriverValue(mealPlanOptionCreationArgs)...).
					WillReturnResult(newArbitraryDatabaseResult())
			}
		}

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() time.Time {
			return exampleMealPlan.CreatedAt
		}

		actual, err := c.CreateMealPlan(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealPlan.Notes,
			exampleMealPlan.Status,
			exampleMealPlan.VotingDeadline,
			exampleMealPlan.BelongsToHousehold,
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealPlanQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateMealPlan(ctx, exampleMealPlan))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateMealPlan(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleMealPlan.Notes,
			exampleMealPlan.Status,
			exampleMealPlan.VotingDeadline,
			exampleMealPlan.BelongsToHousehold,
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealPlanQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateMealPlan(ctx, exampleMealPlan))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleAccountID := fakes.BuildFakeID()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleAccountID,
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealPlanQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveMealPlan(ctx, exampleMealPlan.ID, exampleAccountID))

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleAccountID := fakes.BuildFakeID()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleAccountID,
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealPlanQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveMealPlan(ctx, exampleMealPlan.ID, exampleAccountID))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_AttemptToFinalizeCompleteMealPlan(T *testing.T) {
	T.Parallel()

	optionA := "eggs benedict"
	optionB := "scrambled eggs"
	optionC := "buttered toast"
	userID1 := fakes.BuildFakeID()
	userID2 := fakes.BuildFakeID()
	userID3 := fakes.BuildFakeID()
	userID4 := fakes.BuildFakeID()

	T.Run("with all votes in", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.Members = []*types.HouseholdUserMembershipWithUser{
			{
				ID:                 fakes.BuildFakeID(),
				BelongsToUser:      &types.User{ID: userID1},
				BelongsToHousehold: exampleHousehold.ID,
			},
			{
				ID:                 fakes.BuildFakeID(),
				BelongsToUser:      &types.User{ID: userID2},
				BelongsToHousehold: exampleHousehold.ID,
			},
			{
				ID:                 fakes.BuildFakeID(),
				BelongsToUser:      &types.User{ID: userID3},
				BelongsToHousehold: exampleHousehold.ID,
			},
			{
				ID:                 fakes.BuildFakeID(),
				BelongsToUser:      &types.User{ID: userID4},
				BelongsToHousehold: exampleHousehold.ID,
			},
		}

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.BelongsToHousehold = exampleHousehold.ID
		exampleMealPlan.Events = []*types.MealPlanEvent{
			{
				ID:       fakes.BuildFakeID(),
				MealName: types.BreakfastMealName,
				Options: []*types.MealPlanOption{
					{
						ID:   optionA,
						Meal: *fakes.BuildFakeMeal(),
						Votes: []*types.MealPlanOptionVote{
							{
								BelongsToMealPlanOption: optionA,
								Rank:                    0,
								ByUser:                  userID1,
							},
							{
								BelongsToMealPlanOption: optionA,
								Rank:                    0,
								ByUser:                  userID2,
							},
							{
								BelongsToMealPlanOption: optionA,
								Rank:                    1,
								ByUser:                  userID3,
							},
							{
								BelongsToMealPlanOption: optionA,
								Rank:                    2,
								ByUser:                  userID4,
							},
						},
					},
					{
						ID:   optionB,
						Meal: *fakes.BuildFakeMeal(),
						Votes: []*types.MealPlanOptionVote{
							{
								BelongsToMealPlanOption: optionB,
								Rank:                    0,
								ByUser:                  userID3,
							},
							{
								BelongsToMealPlanOption: optionB,
								Rank:                    1,
								ByUser:                  userID2,
							},
							{
								BelongsToMealPlanOption: optionB,
								Rank:                    1,
								ByUser:                  userID4,
							},
							{
								BelongsToMealPlanOption: optionB,
								Rank:                    2,
								ByUser:                  userID1,
							},
						},
					},
					{
						ID:   optionC,
						Meal: *fakes.BuildFakeMeal(),
						Votes: []*types.MealPlanOptionVote{
							{
								BelongsToMealPlanOption: optionC,
								Rank:                    0,
								ByUser:                  userID4,
							},

							{
								BelongsToMealPlanOption: optionC,
								Rank:                    1,
								ByUser:                  userID1,
							},
							{
								BelongsToMealPlanOption: optionC,
								Rank:                    2,
								ByUser:                  userID2,
							},
							{
								BelongsToMealPlanOption: optionC,
								Rank:                    2,
								ByUser:                  userID3,
							},
						},
					},
				},
			},
		}

		c, db := buildTestClient(t)

		getHouseholdByIDArgs := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(getHouseholdByIDArgs)...).
			WillReturnRows(buildMockRowsFromHouseholds(false, 0, exampleHousehold))

		prepareMockToSuccessfullyGetMealPlan(t, exampleMealPlan, exampleHousehold.ID, db, false)

		db.ExpectBegin()

		for _, event := range exampleMealPlan.Events {
			winner, tiebroken, _ := c.decideOptionWinner(ctx, event.Options)

			finalizeMealPlanOptionsArgs := []any{
				event.ID,
				winner,
				tiebroken,
			}

			db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanOptionQuery)).
				WithArgs(interfaceToDriverValue(finalizeMealPlanOptionsArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		finalizeOptionsArgs := []any{
			types.MealPlanStatusFinalized,
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanQuery)).
			WithArgs(interfaceToDriverValue(finalizeOptionsArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

		actual, err := c.AttemptToFinalizeMealPlan(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with only some votes in, and voting deadline has passed", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.Members = []*types.HouseholdUserMembershipWithUser{
			{
				ID:                 fakes.BuildFakeID(),
				BelongsToUser:      &types.User{ID: userID1},
				BelongsToHousehold: exampleHousehold.ID,
			},
			{
				ID:                 fakes.BuildFakeID(),
				BelongsToUser:      &types.User{ID: userID2},
				BelongsToHousehold: exampleHousehold.ID,
			},
			{
				ID:                 fakes.BuildFakeID(),
				BelongsToUser:      &types.User{ID: userID3},
				BelongsToHousehold: exampleHousehold.ID,
			},
			{
				ID:                 fakes.BuildFakeID(),
				BelongsToUser:      &types.User{ID: userID4},
				BelongsToHousehold: exampleHousehold.ID,
			},
		}

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.VotingDeadline = time.Now().Add(-1 * time.Hour).Round(time.Second)
		exampleMealPlan.BelongsToHousehold = exampleHousehold.ID
		exampleMealPlan.Events = []*types.MealPlanEvent{
			{
				ID:       fakes.BuildFakeID(),
				MealName: types.BreakfastMealName,
				Options: []*types.MealPlanOption{
					{
						ID:   optionA,
						Meal: *fakes.BuildFakeMeal(),
						Votes: []*types.MealPlanOptionVote{
							{
								BelongsToMealPlanOption: optionA,
								Rank:                    0,
								ByUser:                  userID1,
							},
							{
								BelongsToMealPlanOption: optionA,
								Rank:                    0,
								ByUser:                  userID2,
							},
							{
								BelongsToMealPlanOption: optionA,
								Rank:                    2,
								ByUser:                  userID4,
							},
						},
					},
					{
						ID:   optionB,
						Meal: *fakes.BuildFakeMeal(),
						Votes: []*types.MealPlanOptionVote{
							{
								BelongsToMealPlanOption: optionB,
								Rank:                    1,
								ByUser:                  userID2,
							},
							{
								BelongsToMealPlanOption: optionB,
								Rank:                    1,
								ByUser:                  userID4,
							},
							{
								BelongsToMealPlanOption: optionB,
								Rank:                    2,
								ByUser:                  userID1,
							},
						},
					},
					{
						ID:   optionC,
						Meal: *fakes.BuildFakeMeal(),
						Votes: []*types.MealPlanOptionVote{
							{
								BelongsToMealPlanOption: optionC,
								Rank:                    0,
								ByUser:                  userID4,
							},

							{
								BelongsToMealPlanOption: optionC,
								Rank:                    1,
								ByUser:                  userID1,
							},
							{
								BelongsToMealPlanOption: optionC,
								Rank:                    2,
								ByUser:                  userID2,
							},
						},
					},
				},
			},
		}

		c, db := buildTestClient(t)

		getHouseholdByIDArgs := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(getHouseholdByIDArgs)...).
			WillReturnRows(buildMockRowsFromHouseholds(false, 0, exampleHousehold))

		prepareMockToSuccessfullyGetMealPlan(t, exampleMealPlan, exampleHousehold.ID, db, false)

		db.ExpectBegin()

		for _, event := range exampleMealPlan.Events {
			if len(event.Options) > 0 {
				winner, tiebroken, chosen := c.decideOptionWinner(ctx, event.Options)
				if chosen {
					finalizeMealPlanOptionsArgs := []any{
						event.ID,
						winner,
						tiebroken,
					}

					db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanOptionQuery)).
						WithArgs(interfaceToDriverValue(finalizeMealPlanOptionsArgs)...).
						WillReturnResult(newArbitraryDatabaseResult())
				}
			}
		}

		db.ExpectCommit()

		actual, err := c.AttemptToFinalizeMealPlan(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.False(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with no missing votes", func(t *testing.T) {
		t.Parallel()

		user1, user2 := fakes.BuildFakeUser(), fakes.BuildFakeUser()

		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.Members = []*types.HouseholdUserMembershipWithUser{
			{
				BelongsToUser:      user1,
				BelongsToHousehold: exampleHousehold.ID,
				HouseholdRole:      "user",
				DefaultHousehold:   true,
			},
			{
				BelongsToUser:      user2,
				BelongsToHousehold: exampleHousehold.ID,
				HouseholdRole:      "user",
				DefaultHousehold:   true,
			},
		}

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.Events = []*types.MealPlanEvent{exampleMealPlan.Events[0]}
		exampleMealPlan.Events[0].Options = []*types.MealPlanOption{exampleMealPlan.Events[0].Options[0], exampleMealPlan.Events[0].Options[1]}
		exampleMealPlan.Events[0].Options[0].Votes = []*types.MealPlanOptionVote{
			{
				ByUser: user1.ID,
				Rank:   0,
			},
			{
				ByUser: user2.ID,
				Rank:   1,
			},
		}
		exampleMealPlan.Events[0].Options[1].Votes = []*types.MealPlanOptionVote{
			{
				ByUser: user1.ID,
				Rank:   1,
			},
			{
				ByUser: user2.ID,
				Rank:   0,
			},
		}

		for i := range exampleMealPlan.Events[0].Options {
			exampleMealPlan.Events[0].Options[i].Meal = *fakes.BuildFakeMeal()
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(false, 0, exampleHousehold))

		prepareMockToSuccessfullyGetMealPlan(t, exampleMealPlan, exampleHousehold.ID, db, false)

		actual, err := c.FetchMissingVotesForMealPlan(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.NoError(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing votes", func(t *testing.T) {
		t.Parallel()

		user1, user2 := fakes.BuildFakeUser(), fakes.BuildFakeUser()

		exampleHousehold := fakes.BuildFakeHousehold()
		exampleHousehold.Members = []*types.HouseholdUserMembershipWithUser{
			{
				BelongsToUser:      user1,
				BelongsToHousehold: exampleHousehold.ID,
				HouseholdRole:      "user",
				DefaultHousehold:   true,
			},
			{
				BelongsToUser:      user2,
				BelongsToHousehold: exampleHousehold.ID,
				HouseholdRole:      "user",
				DefaultHousehold:   true,
			},
		}

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.Events = []*types.MealPlanEvent{exampleMealPlan.Events[0]}
		exampleMealPlan.Events[0].Options = []*types.MealPlanOption{exampleMealPlan.Events[0].Options[0], exampleMealPlan.Events[0].Options[1]}
		exampleMealPlan.Events[0].Options[0].Votes = []*types.MealPlanOptionVote{
			{
				ByUser: user1.ID,
				Rank:   0,
			},
			{
				ByUser: user2.ID,
				Rank:   1,
			},
		}
		exampleMealPlan.Events[0].Options[1].Votes = []*types.MealPlanOptionVote{
			{
				ByUser: user1.ID,
				Rank:   1,
			},
		}

		for i := range exampleMealPlan.Events[0].Options {
			exampleMealPlan.Events[0].Options[i].Meal = *fakes.BuildFakeMeal()
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromHouseholds(false, 0, exampleHousehold))

		prepareMockToSuccessfullyGetMealPlan(t, exampleMealPlan, exampleHousehold.ID, db, false)

		expected := []*types.MissingVote{
			{
				UserID:   user2.ID,
				OptionID: exampleMealPlan.Events[0].Options[1].ID,
				EventID:  exampleMealPlan.Events[0].ID,
			},
		}

		actual, err := c.FetchMissingVotesForMealPlan(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

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
