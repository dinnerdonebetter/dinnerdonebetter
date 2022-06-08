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

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
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
			x.StartsAt,
			x.EndsAt,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToHousehold,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(mealPlans))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildMockRowsFromFullMealPlans(includeCounts bool, filteredCount uint64, mealPlans ...*types.MealPlan) *sqlmock.Rows {
	columns := []string{
		"meal_plans.id",
		"meal_plans.notes",
		"meal_plans.status",
		"meal_plans.voting_deadline",
		"meal_plans.starts_at",
		"meal_plans.ends_at",
		"meal_plans.created_on",
		"meal_plans.last_updated_on",
		"meal_plans.archived_on",
		"meal_plans.belongs_to_household",
		"meal_plan_options.id",
		"meal_plan_options.day",
		"meal_plan_options.meal_name",
		"meal_plan_options.chosen",
		"meal_plan_options.tiebroken",
		"meal_plan_options.meal_id",
		"meal_plan_options.notes",
		"meal_plan_options.created_on",
		"meal_plan_options.last_updated_on",
		"meal_plan_options.archived_on",
		"meal_plan_options.belongs_to_meal_plan",
		"meal_plan_option_votes.id",
		"meal_plan_option_votes.rank",
		"meal_plan_option_votes.abstain",
		"meal_plan_option_votes.notes",
		"meal_plan_option_votes.by_user",
		"meal_plan_option_votes.created_on",
		"meal_plan_option_votes.last_updated_on",
		"meal_plan_option_votes.archived_on",
		"meal_plan_option_votes.belongs_to_meal_plan_option",
		"meals.id",
		"meals.name",
		"meals.description",
		"meals.created_on",
		"meals.last_updated_on",
		"meals.archived_on",
		"meals.created_by_user",
	}

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range mealPlans {
		for _, opt := range x.Options {
			for _, vote := range opt.Votes {
				rowValues := []driver.Value{
					x.ID,
					x.Notes,
					x.Status,
					x.VotingDeadline,
					x.StartsAt,
					x.EndsAt,
					x.CreatedOn,
					x.LastUpdatedOn,
					x.ArchivedOn,
					x.BelongsToHousehold,
					opt.ID,
					opt.Day,
					opt.MealName,
					opt.Chosen,
					opt.TieBroken,
					opt.Meal.ID,
					opt.Notes,
					opt.CreatedOn,
					opt.LastUpdatedOn,
					opt.ArchivedOn,
					opt.BelongsToMealPlan,
					vote.ID,
					vote.Rank,
					vote.Abstain,
					vote.Notes,
					vote.ByUser,
					vote.CreatedOn,
					vote.LastUpdatedOn,
					vote.ArchivedOn,
					vote.BelongsToMealPlanOption,
					opt.Meal.ID,
					opt.Meal.Name,
					opt.Meal.Description,
					opt.Meal.CreatedOn,
					opt.Meal.LastUpdatedOn,
					opt.Meal.ArchivedOn,
					opt.Meal.CreatedByUser,
				}

				if includeCounts {
					rowValues = append(rowValues, filteredCount, len(mealPlans))
				}

				exampleRows.AddRow(rowValues...)
			}
		}
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
		args := []interface{}{
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
		args := []interface{}{
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
		args := []interface{}{
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

func TestQuerier_GetMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMealPlan.ID,
			exampleHouseholdID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromFullMealPlans(false, 0, exampleMealPlan))

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

		args := []interface{}{
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

		args := []interface{}{
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

func TestQuerier_GetTotalMealPlanCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalMealPlansCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetTotalMealPlanCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalMealPlansCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetTotalMealPlanCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlans(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanList := fakes.BuildFakeMealPlanList()
		for i := range exampleMealPlanList.MealPlans {
			exampleMealPlanList.MealPlans[i].Options = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlans(true, exampleMealPlanList.FilteredCount, exampleMealPlanList.MealPlans...))

		actual, err := c.GetMealPlans(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleMealPlanList := fakes.BuildFakeMealPlanList()
		exampleMealPlanList.Page = 0
		exampleMealPlanList.Limit = 0
		for i := range exampleMealPlanList.MealPlans {
			exampleMealPlanList.MealPlans[i].Options = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlans(true, exampleMealPlanList.FilteredCount, exampleMealPlanList.MealPlans...))

		actual, err := c.GetMealPlans(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlans(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plans", nil, nil, nil, householdOwnershipColumn, mealPlansTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealPlans(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlansWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealPlanList := fakes.BuildFakeMealPlanList()
		for i := range exampleMealPlanList.MealPlans {
			exampleMealPlanList.MealPlans[i].Options = nil
		}

		var exampleIDs []string
		for _, x := range exampleMealPlanList.MealPlans {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetMealPlansWithIDsQuery(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlans(false, 0, exampleMealPlanList.MealPlans...))

		actual, err := c.GetMealPlansWithIDs(ctx, exampleHouseholdID, 0, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanList.MealPlans, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid IDs", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlansWithIDs(ctx, exampleHouseholdID, defaultLimit, nil)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealPlanList := fakes.BuildFakeMealPlanList()

		var exampleIDs []string
		for _, x := range exampleMealPlanList.MealPlans {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetMealPlansWithIDsQuery(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlansWithIDs(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning query results", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()
		exampleMealPlanList := fakes.BuildFakeMealPlanList()

		var exampleIDs []string
		for _, x := range exampleMealPlanList.MealPlans {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetMealPlansWithIDsQuery(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealPlansWithIDs(ctx, exampleHouseholdID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.ID = "1"
		for i := range exampleMealPlan.Options {
			exampleMealPlan.Options[i].ID = "2"
			exampleMealPlan.Options[i].BelongsToMealPlan = "1"
			exampleMealPlan.Options[i].Meal = types.Meal{ID: exampleMealPlan.Options[i].Meal.ID}
			exampleMealPlan.Options[i].CreatedOn = exampleMealPlan.CreatedOn
			exampleMealPlan.Options[i].Votes = []*types.MealPlanOptionVote{}
		}
		exampleInput := fakes.BuildFakeMealPlanDatabaseCreationInputFromMealPlan(exampleMealPlan)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.Status,
			exampleInput.VotingDeadline,
			exampleInput.StartsAt,
			exampleInput.EndsAt,
			exampleInput.BelongsToHousehold,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, option := range exampleInput.Options {
			optionArgs := []interface{}{
				option.ID,
				option.Day,
				option.MealName,
				option.MealID,
				option.Notes,
				option.BelongsToMealPlan,
			}

			db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
				WithArgs(interfaceToDriverValue(optionArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		db.ExpectCommit()

		c.timeFunc = func() uint64 {
			return exampleMealPlan.CreatedOn
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
		for i := range exampleMealPlan.Options {
			exampleMealPlan.Options[i].ID = "2"
			exampleMealPlan.Options[i].BelongsToMealPlan = "1"
			exampleMealPlan.Options[i].CreatedOn = exampleMealPlan.CreatedOn
			exampleMealPlan.Options[i].Votes = []*types.MealPlanOptionVote{}
		}
		exampleInput := fakes.BuildFakeMealPlanDatabaseCreationInputFromMealPlan(exampleMealPlan)

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
		for i := range exampleMealPlan.Options {
			exampleMealPlan.Options[i].ID = "2"
			exampleMealPlan.Options[i].BelongsToMealPlan = "1"
			exampleMealPlan.Options[i].CreatedOn = exampleMealPlan.CreatedOn
		}
		exampleInput := fakes.BuildFakeMealPlanDatabaseCreationInputFromMealPlan(exampleMealPlan)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.Status,
			exampleInput.VotingDeadline,
			exampleInput.StartsAt,
			exampleInput.EndsAt,
			exampleInput.BelongsToHousehold,
		}

		db.ExpectBegin()

		db.ExpectExec(formatQueryForSQLMock(mealPlanCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleMealPlan.CreatedOn
		}

		actual, err := c.CreateMealPlan(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error creating meal plan option", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.ID = "1"
		for i := range exampleMealPlan.Options {
			exampleMealPlan.Options[i].ID = "2"
			exampleMealPlan.Options[i].BelongsToMealPlan = "1"
			exampleMealPlan.Options[i].CreatedOn = exampleMealPlan.CreatedOn
			exampleMealPlan.Options[i].Votes = []*types.MealPlanOptionVote{}
		}
		exampleInput := fakes.BuildFakeMealPlanDatabaseCreationInputFromMealPlan(exampleMealPlan)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.Status,
			exampleInput.VotingDeadline,
			exampleInput.StartsAt,
			exampleInput.EndsAt,
			exampleInput.BelongsToHousehold,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		optionArgs := []interface{}{
			exampleInput.Options[0].ID,
			exampleInput.Options[0].Day,
			exampleInput.Options[0].MealName,
			exampleInput.Options[0].MealID,
			exampleInput.Options[0].Notes,
			exampleInput.Options[0].BelongsToMealPlan,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
			WithArgs(interfaceToDriverValue(optionArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		c.timeFunc = func() uint64 {
			return exampleMealPlan.CreatedOn
		}

		actual, err := c.CreateMealPlan(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.ID = "1"
		for i := range exampleMealPlan.Options {
			exampleMealPlan.Options[i].ID = "2"
			exampleMealPlan.Options[i].BelongsToMealPlan = "1"
			exampleMealPlan.Options[i].CreatedOn = exampleMealPlan.CreatedOn
			exampleMealPlan.Options[i].Votes = []*types.MealPlanOptionVote{}
		}
		exampleInput := fakes.BuildFakeMealPlanDatabaseCreationInputFromMealPlan(exampleMealPlan)

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.Status,
			exampleInput.VotingDeadline,
			exampleInput.StartsAt,
			exampleInput.EndsAt,
			exampleInput.BelongsToHousehold,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		for _, option := range exampleInput.Options {
			optionArgs := []interface{}{
				option.ID,
				option.Day,
				option.MealName,
				option.MealID,
				option.Notes,
				option.BelongsToMealPlan,
			}

			db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
				WithArgs(interfaceToDriverValue(optionArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		c.timeFunc = func() uint64 {
			return exampleMealPlan.CreatedOn
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

		args := []interface{}{
			exampleMealPlan.Notes,
			exampleMealPlan.Status,
			exampleMealPlan.VotingDeadline,
			exampleMealPlan.StartsAt,
			exampleMealPlan.EndsAt,
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

		args := []interface{}{
			exampleMealPlan.Notes,
			exampleMealPlan.Status,
			exampleMealPlan.VotingDeadline,
			exampleMealPlan.StartsAt,
			exampleMealPlan.EndsAt,
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

		args := []interface{}{
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

		args := []interface{}{
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

func Test_byDayAndMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := []*types.MealPlanOption{
			{
				Day:      time.Wednesday,
				MealName: types.SecondBreakfastMealName,
			},
		}
		options := []*types.MealPlanOption{
			{
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
			},
			{
				Day:      time.Tuesday,
				MealName: types.SecondBreakfastMealName,
			},
			expected[0],
			{
				Day:      time.Thursday,
				MealName: types.BrunchMealName,
			},
			{
				Day:      time.Friday,
				MealName: types.LunchMealName,
			},
			{
				Day:      time.Saturday,
				MealName: types.SupperMealName,
			},
			{
				Day:      time.Sunday,
				MealName: types.DinnerMealName,
			},
		}

		actual := byDayAndMeal(options, time.Wednesday, types.SecondBreakfastMealName)

		assert.Equal(t, expected, actual)
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

	T.Run("standard", func(t *testing.T) {
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
		exampleMealPlan.Options = []*types.MealPlanOption{
			{
				ID:       optionA,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				Chosen: true,
			},
			{
				ID:       optionB,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				ID:       optionC,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
		}

		c, db := buildTestClient(t)

		getMealPlanArgs := []interface{}{
			exampleMealPlan.ID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanArgs)...).
			WillReturnRows(buildMockRowsFromFullMealPlans(false, 0, exampleMealPlan))

		finalizeOptionsArgs := []interface{}{
			types.FinalizedMealPlanStatus,
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanQuery)).
			WithArgs(interfaceToDriverValue(finalizeOptionsArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		actual, err := c.AttemptToFinalizeCompleteMealPlan(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.True(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleHousehold := fakes.BuildFakeHousehold()
		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.AttemptToFinalizeCompleteMealPlan(ctx, "", exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.AttemptToFinalizeCompleteMealPlan(ctx, exampleMealPlan.ID, "")
		assert.False(t, actual)
		assert.Error(t, err)
	})

	T.Run("without selection", func(t *testing.T) {
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
		exampleMealPlan.Options = []*types.MealPlanOption{
			{
				ID:       optionA,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				ID:       optionB,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				ID:       optionC,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
		}

		c, db := buildTestClient(t)

		getMealPlanArgs := []interface{}{
			exampleMealPlan.ID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanArgs)...).
			WillReturnRows(buildMockRowsFromFullMealPlans(false, 0, exampleMealPlan))

		actual, err := c.AttemptToFinalizeCompleteMealPlan(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.False(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error fetching meal plan", func(t *testing.T) {
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
		exampleMealPlan.Options = []*types.MealPlanOption{
			{
				ID:       optionA,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				Chosen: true,
			},
			{
				ID:       optionB,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				ID:       optionC,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
		}

		c, db := buildTestClient(t)

		getMealPlanArgs := []interface{}{
			exampleMealPlan.ID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.AttemptToFinalizeCompleteMealPlan(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error marking meal plan as finalized", func(t *testing.T) {
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
		exampleMealPlan.Options = []*types.MealPlanOption{
			{
				ID:       optionA,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				Chosen: true,
			},
			{
				ID:       optionB,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				ID:       optionC,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
		}

		c, db := buildTestClient(t)

		getMealPlanArgs := []interface{}{
			exampleMealPlan.ID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanArgs)...).
			WillReturnRows(buildMockRowsFromFullMealPlans(false, 0, exampleMealPlan))

		finalizeOptionsArgs := []interface{}{
			types.FinalizedMealPlanStatus,
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanQuery)).
			WithArgs(interfaceToDriverValue(finalizeOptionsArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.AttemptToFinalizeCompleteMealPlan(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_FetchExpiredAndUnresolvedMealPlanIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := []*types.MealPlan{}
		exampleMealPlanList := fakes.BuildFakeMealPlanList()
		for _, mp := range exampleMealPlanList.MealPlans {
			mp.Options = nil
			expected = append(expected, mp)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getExpiredAndUnresolvedMealPlanIDsQuery)).
			WithArgs().
			WillReturnRows(buildMockRowsFromMealPlans(false, exampleMealPlanList.FilteredCount, exampleMealPlanList.MealPlans...))

		actual, err := c.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error performing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getExpiredAndUnresolvedMealPlanIDsQuery)).
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
		for _, mp := range exampleMealPlanList.MealPlans {
			mp.Options = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getExpiredAndUnresolvedMealPlanIDsQuery)).
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
		for _, mp := range exampleMealPlanList.MealPlans {
			mp.Options = nil
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getExpiredAndUnresolvedMealPlanIDsQuery)).
			WithArgs().
			WillReturnRows(buildMockRowsFromMealPlans(false, exampleMealPlanList.FilteredCount, exampleMealPlanList.MealPlans...).RowError(0, errors.New("blah")))

		actual, err := c.GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_FinalizeMealPlanWithExpiredVotingPeriod(T *testing.T) {
	T.Parallel()

	optionA := "eggs benedict"
	optionB := "scrambled eggs"
	optionC := "buttered toast"
	userID1 := fakes.BuildFakeID()
	userID2 := fakes.BuildFakeID()
	userID3 := fakes.BuildFakeID()
	userID4 := fakes.BuildFakeID()

	T.Run("standard", func(t *testing.T) {
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
		exampleMealPlan.Options = []*types.MealPlanOption{
			{
				ID:       optionA,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				Chosen: true,
			},
			{
				ID:       optionB,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				ID:       optionC,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
		}

		c, db := buildTestClient(t)

		getMealPlanArgs := []interface{}{
			exampleMealPlan.ID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanArgs)...).
			WillReturnRows(buildMockRowsFromFullMealPlans(false, 0, exampleMealPlan))

		db.ExpectBegin()

		for _, day := range allDays {
			for _, mealName := range allMealNames {
				options := byDayAndMeal(exampleMealPlan.Options, day, mealName)
				if len(options) > 0 {
					winner, tiebroken, _ := c.decideOptionWinner(options)

					finalizeMealPlanOptionsArgs := []interface{}{
						exampleMealPlan.ID,
						winner,
						tiebroken,
					}

					db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanOptionQuery)).
						WithArgs(interfaceToDriverValue(finalizeMealPlanOptionsArgs)...).
						WillReturnResult(newArbitraryDatabaseResult())
				}
			}
		}

		finalizeOptionsArgs := []interface{}{
			types.FinalizedMealPlanStatus,
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanQuery)).
			WithArgs(interfaceToDriverValue(finalizeOptionsArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

		actual, err := c.FinalizeMealPlanWithExpiredVotingPeriod(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.True(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()

		c, db := buildTestClient(t)

		actual, err := c.FinalizeMealPlanWithExpiredVotingPeriod(ctx, "", exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		c, db := buildTestClient(t)

		actual, err := c.FinalizeMealPlanWithExpiredVotingPeriod(ctx, exampleMealPlan.ID, "")
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error fetching meal plan", func(t *testing.T) {
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
		exampleMealPlan.Options = []*types.MealPlanOption{
			{
				ID:       optionA,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				Chosen: true,
			},
			{
				ID:       optionB,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				ID:       optionC,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
		}

		c, db := buildTestClient(t)

		getMealPlanArgs := []interface{}{
			exampleMealPlan.ID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.FinalizeMealPlanWithExpiredVotingPeriod(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleHousehold := fakes.BuildFakeHousehold()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlan.BelongsToHousehold = exampleHousehold.ID

		c, db := buildTestClient(t)

		getMealPlanArgs := []interface{}{
			exampleMealPlan.ID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanArgs)...).
			WillReturnRows(buildMockRowsFromFullMealPlans(false, 0, exampleMealPlan))

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		actual, err := c.FinalizeMealPlanWithExpiredVotingPeriod(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error finalizing meal plan option", func(t *testing.T) {
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
		exampleMealPlan.Options = []*types.MealPlanOption{
			{
				ID:       optionA,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				Chosen: true,
			},
			{
				ID:       optionB,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				ID:       optionC,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
		}

		c, db := buildTestClient(t)

		getMealPlanArgs := []interface{}{
			exampleMealPlan.ID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanArgs)...).
			WillReturnRows(buildMockRowsFromFullMealPlans(false, 0, exampleMealPlan))

		db.ExpectBegin()

		options := byDayAndMeal(exampleMealPlan.Options, time.Monday, types.BreakfastMealName)
		winner, tiebroken, _ := c.decideOptionWinner(options)

		finalizeMealPlanOptionsArgs := []interface{}{
			exampleMealPlan.ID,
			winner,
			tiebroken,
		}

		db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(finalizeMealPlanOptionsArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.FinalizeMealPlanWithExpiredVotingPeriod(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error finalizing meal plan", func(t *testing.T) {
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
		exampleMealPlan.Options = []*types.MealPlanOption{
			{
				ID:       optionA,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				Chosen: true,
			},
			{
				ID:       optionB,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				ID:       optionC,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
		}

		c, db := buildTestClient(t)

		getMealPlanArgs := []interface{}{
			exampleMealPlan.ID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanArgs)...).
			WillReturnRows(buildMockRowsFromFullMealPlans(false, 0, exampleMealPlan))

		db.ExpectBegin()

		for _, day := range allDays {
			for _, mealName := range allMealNames {
				options := byDayAndMeal(exampleMealPlan.Options, day, mealName)
				if len(options) > 0 {
					winner, tiebroken, _ := c.decideOptionWinner(options)

					finalizeMealPlanOptionsArgs := []interface{}{
						exampleMealPlan.ID,
						winner,
						tiebroken,
					}

					db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanOptionQuery)).
						WithArgs(interfaceToDriverValue(finalizeMealPlanOptionsArgs)...).
						WillReturnResult(newArbitraryDatabaseResult())
				}
			}
		}

		finalizeOptionsArgs := []interface{}{
			types.FinalizedMealPlanStatus,
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanQuery)).
			WithArgs(interfaceToDriverValue(finalizeOptionsArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		actual, err := c.FinalizeMealPlanWithExpiredVotingPeriod(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
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
		exampleMealPlan.Options = []*types.MealPlanOption{
			{
				ID:       optionA,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				Chosen: true,
			},
			{
				ID:       optionB,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
				ID:       optionC,
				Day:      time.Monday,
				MealName: types.BreakfastMealName,
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
		}

		c, db := buildTestClient(t)

		getMealPlanArgs := []interface{}{
			exampleMealPlan.ID,
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanArgs)...).
			WillReturnRows(buildMockRowsFromFullMealPlans(false, 0, exampleMealPlan))

		db.ExpectBegin()

		for _, day := range allDays {
			for _, mealName := range allMealNames {
				options := byDayAndMeal(exampleMealPlan.Options, day, mealName)
				if len(options) > 0 {
					winner, tiebroken, _ := c.decideOptionWinner(options)

					finalizeMealPlanOptionsArgs := []interface{}{
						exampleMealPlan.ID,
						winner,
						tiebroken,
					}

					db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanOptionQuery)).
						WithArgs(interfaceToDriverValue(finalizeMealPlanOptionsArgs)...).
						WillReturnResult(newArbitraryDatabaseResult())
				}
			}
		}

		finalizeOptionsArgs := []interface{}{
			types.FinalizedMealPlanStatus,
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanQuery)).
			WithArgs(interfaceToDriverValue(finalizeOptionsArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		actual, err := c.FinalizeMealPlanWithExpiredVotingPeriod(ctx, exampleMealPlan.ID, exampleHousehold.ID)
		assert.False(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}
