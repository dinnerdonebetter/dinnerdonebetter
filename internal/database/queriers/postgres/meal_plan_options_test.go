package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	database "github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
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
			x.Day,
			x.RecipeID,
			x.Notes,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToMealPlan,
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
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanOption.ID,
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanOptionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.MealPlanOptionExists(ctx, exampleMealPlanID, exampleMealPlanOption.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionExists(ctx, "", exampleMealPlanOption.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanOptionExists(ctx, exampleMealPlanID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanOption.ID,
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanOptionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.MealPlanOptionExists(ctx, exampleMealPlanID, exampleMealPlanOption.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanOption.ID,
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanOptionExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.MealPlanOptionExists(ctx, exampleMealPlanID, exampleMealPlanOption.ID)
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
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanOption.ID,
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlanOptions(false, 0, exampleMealPlanOption))

		actual, err := c.GetMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanOption.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOption, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOption(ctx, "", exampleMealPlanOption.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOption(ctx, exampleMealPlanID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanOption.ID,
			exampleMealPlanID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanOption.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetTotalMealPlanOptionCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalMealPlanOptionsCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetTotalMealPlanOptionCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalMealPlanOptionsCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetTotalMealPlanOptionCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlanOptions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionList := fakes.BuildFakeMealPlanOptionList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plan_options", getMealPlanOptionsJoins, nil, nil, householdOwnershipColumn, mealPlanOptionsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlanOptions(true, exampleMealPlanOptionList.FilteredCount, exampleMealPlanOptionList.MealPlanOptions...))

		actual, err := c.GetMealPlanOptions(ctx, exampleMealPlanID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOptionList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptions(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionList := fakes.BuildFakeMealPlanOptionList()
		exampleMealPlanOptionList.Page = 0
		exampleMealPlanOptionList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plan_options", getMealPlanOptionsJoins, nil, nil, householdOwnershipColumn, mealPlanOptionsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlanOptions(true, exampleMealPlanOptionList.FilteredCount, exampleMealPlanOptionList.MealPlanOptions...))

		actual, err := c.GetMealPlanOptions(ctx, exampleMealPlanID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOptionList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plan_options", getMealPlanOptionsJoins, nil, nil, householdOwnershipColumn, mealPlanOptionsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlanOptions(ctx, exampleMealPlanID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleMealPlanID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "meal_plan_options", getMealPlanOptionsJoins, nil, nil, householdOwnershipColumn, mealPlanOptionsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealPlanOptions(ctx, exampleMealPlanID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlanOptionsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionList := fakes.BuildFakeMealPlanOptionList()

		var exampleIDs []string
		for _, x := range exampleMealPlanOptionList.MealPlanOptions {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetMealPlanOptionsWithIDsQuery(ctx, exampleMealPlanID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromMealPlanOptions(false, 0, exampleMealPlanOptionList.MealPlanOptions...))

		actual, err := c.GetMealPlanOptionsWithIDs(ctx, exampleMealPlanID, 0, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanOptionList.MealPlanOptions, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid IDs", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionsWithIDs(ctx, exampleMealPlanID, defaultLimit, nil)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOptionList := fakes.BuildFakeMealPlanOptionList()

		var exampleIDs []string
		for _, x := range exampleMealPlanOptionList.MealPlanOptions {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanOptionsWithIDs(ctx, "", defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionList := fakes.BuildFakeMealPlanOptionList()

		var exampleIDs []string
		for _, x := range exampleMealPlanOptionList.MealPlanOptions {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetMealPlanOptionsWithIDsQuery(ctx, exampleMealPlanID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlanOptionsWithIDs(ctx, exampleMealPlanID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning query results", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOptionList := fakes.BuildFakeMealPlanOptionList()

		var exampleIDs []string
		for _, x := range exampleMealPlanOptionList.MealPlanOptions {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetMealPlanOptionsWithIDsQuery(ctx, exampleMealPlanID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealPlanOptionsWithIDs(ctx, exampleMealPlanID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlanOption(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
		exampleMealPlanOption.ID = "1"
		exampleInput := fakes.BuildFakeMealPlanOptionDatabaseCreationInputFromMealPlanOption(exampleMealPlanOption)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Day,
			exampleInput.RecipeID,
			exampleInput.Notes,
			exampleInput.BelongsToMealPlan,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleMealPlanOption.ID))

		c.timeFunc = func() uint64 {
			return exampleMealPlanOption.CreatedOn
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
		exampleInput := fakes.BuildFakeMealPlanOptionDatabaseCreationInputFromMealPlanOption(exampleMealPlanOption)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Day,
			exampleInput.RecipeID,
			exampleInput.Notes,
			exampleInput.BelongsToMealPlan,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanOptionCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() uint64 {
			return exampleMealPlanOption.CreatedOn
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
			exampleMealPlanOption.Day,
			exampleMealPlanOption.RecipeID,
			exampleMealPlanOption.Notes,
			exampleMealPlanOption.BelongsToMealPlan,
			exampleMealPlanOption.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleMealPlanOption.ID))

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
			exampleMealPlanOption.Day,
			exampleMealPlanOption.RecipeID,
			exampleMealPlanOption.Notes,
			exampleMealPlanOption.BelongsToMealPlan,
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
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanOption.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleMealPlanOption.ID))

		assert.NoError(t, c.ArchiveMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanOption.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanOption(ctx, "", exampleMealPlanOption.ID))
	})

	T.Run("with invalid meal plan option ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanOption(ctx, exampleMealPlanID, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanID := fakes.BuildFakeID()
		exampleMealPlanOption := fakes.BuildFakeMealPlanOption()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleMealPlanID,
			exampleMealPlanOption.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealPlanOptionQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveMealPlanOption(ctx, exampleMealPlanID, exampleMealPlanOption.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
