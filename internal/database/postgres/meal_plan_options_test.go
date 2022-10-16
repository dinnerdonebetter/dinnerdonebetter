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
	"resenje.org/schulze"

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/converters"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromMealPlanOptions(includeCounts bool, filteredCount uint64, mealPlanOptions ...*types.MealPlanOption) *sqlmock.Rows {
	columns := mealPlanOptionsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range mealPlanOptions {
		rowValues := []driver.Value{
			x.ID,
			x.AssignedCook,
			x.AssignedDishwasher,
			x.Chosen,
			x.TieBroken,
			x.Meal.ID,
			x.Notes,
			x.PrepStepsCreated,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
			x.BelongsToMealPlanEvent,
			x.Meal.ID,
			x.Meal.Name,
			x.Meal.Description,
			x.Meal.CreatedAt,
			x.Meal.LastUpdatedAt,
			x.Meal.ArchivedAt,
			x.Meal.CreatedByUser,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(mealPlanOptions))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanMealPlanOptions(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanMealPlanOptions(ctx, mockRows, false)
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

		_, _, _, err := q.scanMealPlanOptions(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_MealPlanOptionExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanEventID,
			exampleMealPlanOption.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanOptionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.MealPlanOptionExists(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionExists(ctx, "", "", exampleMealPlanOption.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionExists(ctx, exampleMealPlanID, exampleMealPlanEventID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanEventID,
			exampleMealPlanOption.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanOptionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.MealPlanOptionExists(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanEventID,
			exampleMealPlanOption.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanOptionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.MealPlanOptionExists(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
		exampleMealPlanOption.Votes = []*types.MealPlanOptionVote{}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanEventID,
			exampleMealPlanOption.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlanOptions(false, 0, exampleMealPlanOption))

		actual, err := c.GetMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOption, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOption(ctx, "", "", exampleMealPlanOption.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanEventID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanEventID,
			exampleMealPlanOption.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlanOptions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionList := fakes.BuildFakeMealPlanOptionList()

		for i := range exampleMealPlanOptionList.MealPlanOptions {
			exampleMealPlanOptionList.MealPlanOptions[i].Votes = []*types.MealPlanOptionVote{}
			exampleMealPlanOptionList.MealPlanOptions[i].Meal = types.Meal{ID: exampleMealPlanOptionList.MealPlanOptions[i].Meal.ID}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		groupBys := []string{"meal_plan_options.id", "meals.id"}
		query, args := c.buildListQuery(ctx, "meal_plan_options", getMealPlanOptionsJoins, groupBys, nil, householdOwnershipColumn, mealPlanOptionsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlanOptions(true, exampleMealPlanOptionList.FilteredCount, exampleMealPlanOptionList.MealPlanOptions...))

		actual, err := c.GetMealPlanOptions(ctx, exampleMealPlanID, exampleMealPlanEventID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOptionList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptions(ctx, "", "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOptionList := fakes.BuildFakeMealPlanOptionList()
		exampleMealPlanOptionList.Page = 0
		exampleMealPlanOptionList.Limit = 0

		for i := range exampleMealPlanOptionList.MealPlanOptions {
			exampleMealPlanOptionList.MealPlanOptions[i].Votes = []*types.MealPlanOptionVote{}
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		groupBys := []string{"meal_plan_options.id", "meals.id"}
		query, args := c.buildListQuery(ctx, "meal_plan_options", getMealPlanOptionsJoins, groupBys, nil, householdOwnershipColumn, mealPlanOptionsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlanOptions(true, exampleMealPlanOptionList.FilteredCount, exampleMealPlanOptionList.MealPlanOptions...))

		actual, err := c.GetMealPlanOptions(ctx, exampleMealPlanID, exampleMealPlanEventID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOptionList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		groupBys := []string{"meal_plan_options.id", "meals.id"}
		query, args := c.buildListQuery(ctx, "meal_plan_options", getMealPlanOptionsJoins, groupBys, nil, householdOwnershipColumn, mealPlanOptionsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlanOptions(ctx, exampleMealPlanID, exampleMealPlanEventID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		groupBys := []string{"meal_plan_options.id", "meals.id"}
		query, args := c.buildListQuery(ctx, "meal_plan_options", getMealPlanOptionsJoins, groupBys, nil, householdOwnershipColumn, mealPlanOptionsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealPlanOptions(ctx, exampleMealPlanID, exampleMealPlanEventID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
		exampleMealPlanOption.ID = "1"
		exampleMealPlanOption.Votes = []*types.MealPlanOptionVote{}
		exampleMealPlanOption.Meal = types.Meal{ID: exampleMealPlanOption.Meal.ID}
		exampleInput := converters.ConvertMealPlanOptionToMealPlanOptionDatabaseCreationInput(exampleMealPlanOption)

		ctx := context.Background()
		c, db := buildTestClient(t)

		mealPlanOptionCreationArgs := []interface{}{
			exampleInput.ID,
			exampleInput.AssignedCook,
			exampleInput.AssignedDishwasher,
			exampleInput.MealID,
			exampleInput.Notes,
			exampleInput.BelongsToMealPlanEvent,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
			WithArgs(interfaceToDriverValue(mealPlanOptionCreationArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleMealPlanOption.CreatedAt
		}

		actual, err := c.CreateMealPlanOption(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOption, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealPlanOption(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
		exampleMealPlanOption.Meal = types.Meal{ID: exampleMealPlanOption.Meal.ID}
		exampleInput := converters.ConvertMealPlanOptionToMealPlanOptionDatabaseCreationInput(exampleMealPlanOption)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.AssignedCook,
			exampleInput.AssignedDishwasher,
			exampleInput.MealID,
			exampleInput.Notes,
			exampleInput.BelongsToMealPlanEvent,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleMealPlanOption.CreatedAt
		}

		actual, err := c.CreateMealPlanOption(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMealPlanOption.AssignedCook,
			exampleMealPlanOption.AssignedDishwasher,
			exampleMealPlanOption.Meal.ID,
			exampleMealPlanOption.Notes,
			exampleMealPlanOption.PrepStepsCreated,
			exampleMealPlanOption.BelongsToMealPlanEvent,
			exampleMealPlanOption.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateMealPlanOption(ctx, exampleMealPlanOption))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateMealPlanOption(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMealPlanOption.AssignedCook,
			exampleMealPlanOption.AssignedDishwasher,
			exampleMealPlanOption.Meal.ID,
			exampleMealPlanOption.Notes,
			exampleMealPlanOption.PrepStepsCreated,
			exampleMealPlanOption.BelongsToMealPlanEvent,
			exampleMealPlanOption.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateMealPlanOption(ctx, exampleMealPlanOption))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMealPlanEventID,
			exampleMealPlanOption.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanOption(ctx, "", "", exampleMealPlanOption.ID))
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanEventID, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanEventID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMealPlanEventID,
			exampleMealPlanOption.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanEventID, exampleMealPlanOption.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func Test_determineWinner(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		c, _ := buildTestClient(t)

		expected := "blah blah blah"
		exampleWinners := []schulze.Score{
			{
				Choice: t.Name(),
				Wins:   1,
			},
			{
				Choice: "",
				Wins:   2,
			},
			{
				Choice: expected,
				Wins:   3,
			},
		}

		actual := c.determineWinner(exampleWinners)

		assert.Equal(t, expected, actual)
	})

	T.Run("with tie", func(t *testing.T) {
		t.Parallel()

		c, _ := buildTestClient(t)

		expectedA := "blah blah blah"
		expectedB := "beeble beeble"
		exampleWinners := []schulze.Score{
			{
				Choice: expectedA,
				Wins:   3,
			},
			{
				Choice: "",
				Wins:   1,
			},
			{
				Choice: expectedB,
				Wins:   3,
			},
		}

		actual := c.determineWinner(exampleWinners)

		assert.True(t, expectedA == actual || expectedB == actual)
	})
}

func Test_decideOptionWinner(T *testing.T) {
	T.Parallel()

	optionA := "eggs benedict"
	optionB := "scrambled eggs"
	optionC := "buttered toast"
	userID1 := fakes.BuildFakeID()
	userID2 := fakes.BuildFakeID()
	userID3 := fakes.BuildFakeID()
	userID4 := fakes.BuildFakeID()

	T.Run("with clear winner", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		expected := optionA
		exampleOptions := []*types.MealPlanOption{
			{
				ID: optionA,
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
				ID: optionB,
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
				ID: optionC,
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

		actual, tiebroken, chosen := c.decideOptionWinner(ctx, exampleOptions)
		assert.Equal(t, expected, actual)
		assert.False(t, tiebroken)
		assert.True(t, chosen)
	})

	T.Run("with tie", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		exampleOptions := []*types.MealPlanOption{
			{
				ID: optionA,
				Votes: []*types.MealPlanOptionVote{
					{
						BelongsToMealPlanOption: optionA,
						Rank:                    0,
						ByUser:                  userID1,
					},
				},
			},
			{
				ID: optionB,
				Votes: []*types.MealPlanOptionVote{
					{
						BelongsToMealPlanOption: optionB,
						Rank:                    0,
						ByUser:                  userID2,
					},
				},
			},
		}

		actual, tiebroken, chosen := c.decideOptionWinner(ctx, exampleOptions)
		assert.NotEmpty(t, actual)
		assert.True(t, tiebroken)
		assert.True(t, chosen)
	})

	T.Run("without enough votes", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		exampleOptions := []*types.MealPlanOption{
			{
				ID:    optionA,
				Votes: nil,
			},
		}

		actual, tiebroken, chosen := c.decideOptionWinner(ctx, exampleOptions)
		assert.Empty(t, actual)
		assert.False(t, tiebroken)
		assert.False(t, chosen)
	})
}

func TestQuerier_MealPlanOptionCanBeFinalized(T *testing.T) {
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
		exampleMealPlan.Events = []*types.MealPlanEvent{
			{
				ID: fakes.BuildFakeID(),
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

		for i := range exampleMealPlan.Events {
			for j := range exampleMealPlan.Events[i].Options {
				for k := range exampleMealPlan.Events[i].Options[j].Meal.Recipes {
					exampleMealPlan.Events[i].Options[j].Meal.Recipes[k].PrepTasks = nil
				}
			}
		}

		c, db := buildTestClient(t)

		prepareMockToSuccessfullyGetMealPlan(t, exampleMealPlan, exampleHousehold.ID, db, false)

		getHouseholdArgs := []interface{}{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(getHouseholdArgs)...).
			WillReturnRows(buildMockRowsFromHouseholds(false, 0, exampleHousehold))

		finalizeMealPlanOptionsArgs := []interface{}{
			exampleMealPlan.Events[0].ID,
			optionA,
			false,
		}

		db.ExpectExec(formatQueryForSQLMock(finalizeMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(finalizeMealPlanOptionsArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		actual, err := c.FinalizeMealPlanOption(ctx, exampleMealPlan.ID, exampleMealPlan.Events[0].ID, exampleMealPlan.Events[0].Options[0].ID, exampleHousehold.ID)
		assert.True(t, actual)
		assert.NoError(t, err)
	})

	T.Run("with missing votes", func(t *testing.T) {
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
				ID: fakes.BuildFakeID(),
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

		prepareMockToSuccessfullyGetMealPlan(t, exampleMealPlan, exampleHousehold.ID, db, false)

		getHouseholdArgs := []interface{}{
			exampleHousehold.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getHouseholdAndMembershipsByIDQuery)).
			WithArgs(interfaceToDriverValue(getHouseholdArgs)...).
			WillReturnRows(buildMockRowsFromHouseholds(false, 0, exampleHousehold))

		actual, err := c.FinalizeMealPlanOption(ctx, exampleMealPlan.ID, exampleMealPlan.Events[0].ID, exampleMealPlan.Events[0].Options[0].ID, exampleHousehold.ID)
		assert.False(t, actual)
		assert.NoError(t, err)
	})
}
