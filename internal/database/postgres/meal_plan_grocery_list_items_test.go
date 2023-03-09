package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"
)

func buildMockRowsFromMealPlanGroceryListItems(includeCounts bool, filteredCount uint64, mealPlanGroceryListItems ...*types.MealPlanGroceryListItem) *sqlmock.Rows {
	columns := mealPlanGroceryListItemsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range mealPlanGroceryListItems {
		var purchasedMeasurementUnitID *string
		if x.PurchasedMeasurementUnit != nil {
			purchasedMeasurementUnitID = &x.PurchasedMeasurementUnit.ID
		}

		rowValues := []driver.Value{
			x.ID,
			x.BelongsToMealPlan,
			x.Ingredient.ID,
			x.MeasurementUnit.ID,
			x.MinimumQuantityNeeded,
			x.MaximumQuantityNeeded,
			x.QuantityPurchased,
			purchasedMeasurementUnitID,
			x.PurchasedUPC,
			x.PurchasePrice,
			x.StatusExplanation,
			x.Status,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(mealPlanGroceryListItems))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanMealPlanGroceryListItems(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := q.scanMealPlanGroceryListItems(ctx, mockRows)
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

		_, err := q.scanMealPlanGroceryListItems(ctx, mockRows)
		assert.Error(t, err)
	})
}

func TestQuerier_MealPlanGroceryListItemExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		c, db := buildTestClient(t)
		args := []any{
			exampleMealPlanGroceryListItem.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanGroceryListItemExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.MealPlanGroceryListItemExists(ctx, exampleMealPlan.ID, exampleMealPlanGroceryListItem.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan grocery list item ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlan := fakes.BuildFakeMealPlan()
		c, _ := buildTestClient(t)

		actual, err := c.MealPlanGroceryListItemExists(ctx, exampleMealPlan.ID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		c, db := buildTestClient(t)
		args := []any{
			exampleMealPlanGroceryListItem.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanGroceryListItemExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.MealPlanGroceryListItemExists(ctx, exampleMealPlan.ID, exampleMealPlanGroceryListItem.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		c, db := buildTestClient(t)
		args := []any{
			exampleMealPlanGroceryListItem.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(mealPlanGroceryListItemExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.MealPlanGroceryListItemExists(ctx, exampleMealPlan.ID, exampleMealPlanGroceryListItem.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_fleshOutMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		c, db := buildTestClient(t)

		getValidIngredientArgs := []any{
			exampleMealPlanGroceryListItem.Ingredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(getValidIngredientArgs)...).
			WillReturnRows(buildMockRowsFromValidIngredients(false, 0, &exampleMealPlanGroceryListItem.Ingredient))

		getValidMeasurementUnitArgs := []any{
			exampleMealPlanGroceryListItem.MeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(getValidMeasurementUnitArgs)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnits(false, 0, &exampleMealPlanGroceryListItem.MeasurementUnit))

		actual, err := c.fleshOutMealPlanGroceryListItem(ctx, exampleMealPlanGroceryListItem)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanGroceryListItem, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.fleshOutMealPlanGroceryListItem(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error fetching ingredient data", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		c, db := buildTestClient(t)

		getValidIngredientArgs := []any{
			exampleMealPlanGroceryListItem.Ingredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(getValidIngredientArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.fleshOutMealPlanGroceryListItem(ctx, exampleMealPlanGroceryListItem)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error fetching measurement unit data", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		c, db := buildTestClient(t)

		getValidIngredientArgs := []any{
			exampleMealPlanGroceryListItem.Ingredient.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientQuery)).
			WithArgs(interfaceToDriverValue(getValidIngredientArgs)...).
			WillReturnRows(buildMockRowsFromValidIngredients(false, 0, &exampleMealPlanGroceryListItem.Ingredient))

		getValidMeasurementUnitArgs := []any{
			exampleMealPlanGroceryListItem.MeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(getValidMeasurementUnitArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.fleshOutMealPlanGroceryListItem(ctx, exampleMealPlanGroceryListItem)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func prepareForMealPlanGroceryListItemDataHydration(t *testing.T, db sqlmock.Sqlmock, exampleMealPlanGroceryListItem *types.MealPlanGroceryListItem) {
	t.Helper()

	getValidIngredientArgs := []any{
		exampleMealPlanGroceryListItem.Ingredient.ID,
	}

	db.ExpectQuery(formatQueryForSQLMock(getValidIngredientQuery)).
		WithArgs(interfaceToDriverValue(getValidIngredientArgs)...).
		WillReturnRows(buildMockRowsFromValidIngredients(false, 0, &exampleMealPlanGroceryListItem.Ingredient))

	getValidMeasurementUnitArgs := []any{
		exampleMealPlanGroceryListItem.MeasurementUnit.ID,
	}

	db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementUnitQuery)).
		WithArgs(interfaceToDriverValue(getValidMeasurementUnitArgs)...).
		WillReturnRows(buildMockRowsFromValidMeasurementUnits(false, 0, &exampleMealPlanGroceryListItem.MeasurementUnit))
}

func TestQuerier_GetMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		c, db := buildTestClient(t)

		getMealPlanGroceryListItemArgs := []any{
			exampleMealPlan.ID,
			exampleMealPlanGroceryListItem.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanGroceryListItemQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanGroceryListItemArgs)...).
			WillReturnRows(buildMockRowsFromMealPlanGroceryListItems(false, 0, exampleMealPlanGroceryListItem))

		prepareForMealPlanGroceryListItemDataHydration(t, db, exampleMealPlanGroceryListItem)

		actual, err := c.GetMealPlanGroceryListItem(ctx, exampleMealPlan.ID, exampleMealPlanGroceryListItem.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanGroceryListItem, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan grocery list item ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanGroceryListItem(ctx, exampleMealPlan.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		c, db := buildTestClient(t)

		getMealPlanGroceryListItemArgs := []any{
			exampleMealPlan.ID,
			exampleMealPlanGroceryListItem.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getMealPlanGroceryListItemQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanGroceryListItemArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlanGroceryListItem(ctx, exampleMealPlan.ID, exampleMealPlanGroceryListItem.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetMealPlanGroceryListItemsForMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanGroceryListItemList := fakes.BuildFakeMealPlanGroceryListItemList().Data
		for i := range exampleMealPlanGroceryListItemList {
			exampleMealPlanGroceryListItemList[i].Ingredient = types.ValidIngredient{
				ID: exampleMealPlanGroceryListItemList[i].Ingredient.ID,
			}
			exampleMealPlanGroceryListItemList[i].MeasurementUnit = types.ValidMeasurementUnit{
				ID: exampleMealPlanGroceryListItemList[i].MeasurementUnit.ID,
			}
		}

		c, db := buildTestClient(t)

		getMealPlanGroceryListItemsForMealPlanArgs := []any{
			exampleMealPlan.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getMealPlanGroceryListItemsForMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanGroceryListItemsForMealPlanArgs)...).
			WillReturnRows(buildMockRowsFromMealPlanGroceryListItems(false, 0, exampleMealPlanGroceryListItemList...))

		for i := range exampleMealPlanGroceryListItemList {
			prepareForMealPlanGroceryListItemDataHydration(t, db, exampleMealPlanGroceryListItemList[i])
		}

		actual, err := c.GetMealPlanGroceryListItemsForMealPlan(ctx, exampleMealPlan.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanGroceryListItemList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		c, db := buildTestClient(t)

		getMealPlanGroceryListItemsForMealPlanArgs := []any{
			exampleMealPlan.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getMealPlanGroceryListItemsForMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanGroceryListItemsForMealPlanArgs)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetMealPlanGroceryListItemsForMealPlan(ctx, exampleMealPlan.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()

		c, db := buildTestClient(t)

		getMealPlanGroceryListItemsForMealPlanArgs := []any{
			exampleMealPlan.ID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getMealPlanGroceryListItemsForMealPlanQuery)).
			WithArgs(interfaceToDriverValue(getMealPlanGroceryListItemsForMealPlanArgs)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetMealPlanGroceryListItemsForMealPlan(ctx, exampleMealPlan.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_createMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleMealPlanGroceryListItem.ID = "1"
		exampleMealPlanGroceryListItem.PurchasedMeasurementUnit = &types.ValidMeasurementUnit{
			ID: fakes.BuildFakeID(),
		}
		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(exampleMealPlanGroceryListItem)

		exampleMealPlanGroceryListItem.Ingredient = types.ValidIngredient{
			ID: exampleMealPlanGroceryListItem.Ingredient.ID,
		}
		exampleMealPlanGroceryListItem.MeasurementUnit = types.ValidMeasurementUnit{
			ID: exampleMealPlanGroceryListItem.MeasurementUnit.ID,
		}
		c, db := buildTestClient(t)

		db.ExpectBegin()
		tx, err := c.db.BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, tx)

		args := []any{
			exampleInput.ID,
			exampleInput.BelongsToMealPlan,
			exampleInput.ValidIngredientID,
			exampleInput.ValidMeasurementUnitID,
			exampleInput.MinimumQuantityNeeded,
			exampleInput.MaximumQuantityNeeded,
			exampleInput.QuantityPurchased,
			exampleInput.PurchasedMeasurementUnitID,
			exampleInput.PurchasedUPC,
			exampleInput.PurchasePrice,
			exampleInput.StatusExplanation,
			exampleInput.Status,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanGroceryListItemCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleMealPlanGroceryListItem.CreatedAt
		}

		actual, err := c.createMealPlanGroceryListItem(ctx, tx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanGroceryListItem, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		tx, err := c.db.BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, tx)

		actual, err := c.createMealPlanGroceryListItem(ctx, tx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleMealPlanGroceryListItem.ID = "1"
		exampleMealPlanGroceryListItem.PurchasedMeasurementUnit = fakes.BuildFakeValidMeasurementUnit()
		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(exampleMealPlanGroceryListItem)

		exampleMealPlanGroceryListItem.Ingredient = types.ValidIngredient{
			ID: exampleMealPlanGroceryListItem.Ingredient.ID,
		}
		exampleMealPlanGroceryListItem.MeasurementUnit = types.ValidMeasurementUnit{
			ID: exampleMealPlanGroceryListItem.MeasurementUnit.ID,
		}
		c, db := buildTestClient(t)

		db.ExpectBegin()
		tx, err := c.db.BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, tx)

		args := []any{
			exampleInput.ID,
			exampleInput.BelongsToMealPlan,
			exampleInput.ValidIngredientID,
			exampleInput.ValidMeasurementUnitID,
			exampleInput.MinimumQuantityNeeded,
			exampleInput.MaximumQuantityNeeded,
			exampleInput.QuantityPurchased,
			exampleInput.PurchasedMeasurementUnitID,
			exampleInput.PurchasedUPC,
			exampleInput.PurchasePrice,
			exampleInput.StatusExplanation,
			exampleInput.Status,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanGroceryListItemCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		c.timeFunc = func() time.Time {
			return exampleMealPlanGroceryListItem.CreatedAt
		}

		actual, err := c.createMealPlanGroceryListItem(ctx, tx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleMealPlanGroceryListItem.ID = "1"
		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(exampleMealPlanGroceryListItem)

		exampleMealPlanGroceryListItem.Ingredient = types.ValidIngredient{
			ID: exampleMealPlanGroceryListItem.Ingredient.ID,
		}
		exampleMealPlanGroceryListItem.MeasurementUnit = types.ValidMeasurementUnit{
			ID: exampleMealPlanGroceryListItem.MeasurementUnit.ID,
		}

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.BelongsToMealPlan,
			exampleInput.ValidIngredientID,
			exampleInput.ValidMeasurementUnitID,
			exampleInput.MinimumQuantityNeeded,
			exampleInput.MaximumQuantityNeeded,
			exampleInput.QuantityPurchased,
			exampleInput.PurchasedMeasurementUnitID,
			exampleInput.PurchasedUPC,
			exampleInput.PurchasePrice,
			exampleInput.StatusExplanation,
			exampleInput.Status,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanGroceryListItemCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleMealPlanGroceryListItem.CreatedAt
		}

		actual, err := c.CreateMealPlanGroceryListItem(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleMealPlanGroceryListItem, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealPlanGroceryListItem(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expectedErr := errors.New(t.Name())
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(exampleMealPlanGroceryListItem)

		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.BelongsToMealPlan,
			exampleInput.ValidIngredientID,
			exampleInput.ValidMeasurementUnitID,
			exampleInput.MinimumQuantityNeeded,
			exampleInput.MaximumQuantityNeeded,
			exampleInput.QuantityPurchased,
			exampleInput.PurchasedMeasurementUnitID,
			exampleInput.PurchasedUPC,
			exampleInput.PurchasePrice,
			exampleInput.StatusExplanation,
			exampleInput.Status,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanGroceryListItemCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleMealPlanGroceryListItem.CreatedAt
		}

		actual, err := c.CreateMealPlanGroceryListItem(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlanGroceryListItemsForMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(exampleMealPlanGroceryListItem)
		inputs := []*types.MealPlanGroceryListItemDatabaseCreationInput{exampleInput}

		c, db := buildTestClient(t)

		db.ExpectBegin()

		for _, input := range inputs {
			args := []any{
				input.ID,
				input.BelongsToMealPlan,
				input.ValidIngredientID,
				input.ValidMeasurementUnitID,
				input.MinimumQuantityNeeded,
				input.MaximumQuantityNeeded,
				input.QuantityPurchased,
				input.PurchasedMeasurementUnitID,
				input.PurchasedUPC,
				input.PurchasePrice,
				input.StatusExplanation,
				input.Status,
			}

			db.ExpectExec(formatQueryForSQLMock(mealPlanGroceryListItemCreationQuery)).
				WithArgs(interfaceToDriverValue(args)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		markMealPlanOptionAsHavingStepsCreatedArgs := []any{
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(markMealPlanAsHavingGroceryListInitialized)).
			WithArgs(interfaceToDriverValue(markMealPlanOptionAsHavingStepsCreatedArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit()

		assert.NoError(t, c.CreateMealPlanGroceryListItemsForMealPlan(ctx, exampleMealPlan.ID, inputs))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(exampleMealPlanGroceryListItem)
		inputs := []*types.MealPlanGroceryListItemDatabaseCreationInput{exampleInput}

		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.CreateMealPlanGroceryListItemsForMealPlan(ctx, exampleMealPlan.ID, inputs))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing creation query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(exampleMealPlanGroceryListItem)
		inputs := []*types.MealPlanGroceryListItemDatabaseCreationInput{exampleInput}

		c, db := buildTestClient(t)

		db.ExpectBegin()

		args := []any{
			exampleInput.ID,
			exampleInput.BelongsToMealPlan,
			exampleInput.ValidIngredientID,
			exampleInput.ValidMeasurementUnitID,
			exampleInput.MinimumQuantityNeeded,
			exampleInput.MaximumQuantityNeeded,
			exampleInput.QuantityPurchased,
			exampleInput.PurchasedMeasurementUnitID,
			exampleInput.PurchasedUPC,
			exampleInput.PurchasePrice,
			exampleInput.StatusExplanation,
			exampleInput.Status,
		}

		db.ExpectExec(formatQueryForSQLMock(mealPlanGroceryListItemCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.CreateMealPlanGroceryListItemsForMealPlan(ctx, exampleMealPlan.ID, inputs))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error marking steps as created", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(exampleMealPlanGroceryListItem)
		inputs := []*types.MealPlanGroceryListItemDatabaseCreationInput{exampleInput}

		c, db := buildTestClient(t)

		db.ExpectBegin()

		for _, input := range inputs {
			args := []any{
				input.ID,
				input.BelongsToMealPlan,
				input.ValidIngredientID,
				input.ValidMeasurementUnitID,
				input.MinimumQuantityNeeded,
				input.MaximumQuantityNeeded,
				input.QuantityPurchased,
				input.PurchasedMeasurementUnitID,
				input.PurchasedUPC,
				input.PurchasePrice,
				input.StatusExplanation,
				input.Status,
			}

			db.ExpectExec(formatQueryForSQLMock(mealPlanGroceryListItemCreationQuery)).
				WithArgs(interfaceToDriverValue(args)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		markMealPlanOptionAsHavingStepsCreatedArgs := []any{
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(markMealPlanAsHavingGroceryListInitialized)).
			WithArgs(interfaceToDriverValue(markMealPlanOptionAsHavingStepsCreatedArgs)...).
			WillReturnError(errors.New("blah"))

		db.ExpectRollback()

		assert.Error(t, c.CreateMealPlanGroceryListItemsForMealPlan(ctx, exampleMealPlan.ID, inputs))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error committing transaction", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlan := fakes.BuildFakeMealPlan()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleInput := converters.ConvertMealPlanGroceryListItemToMealPlanGroceryListItemDatabaseCreationInput(exampleMealPlanGroceryListItem)
		inputs := []*types.MealPlanGroceryListItemDatabaseCreationInput{exampleInput}

		c, db := buildTestClient(t)

		db.ExpectBegin()

		for _, input := range inputs {
			args := []any{
				input.ID,
				input.BelongsToMealPlan,
				input.ValidIngredientID,
				input.ValidMeasurementUnitID,
				input.MinimumQuantityNeeded,
				input.MaximumQuantityNeeded,
				input.QuantityPurchased,
				input.PurchasedMeasurementUnitID,
				input.PurchasedUPC,
				input.PurchasePrice,
				input.StatusExplanation,
				input.Status,
			}

			db.ExpectExec(formatQueryForSQLMock(mealPlanGroceryListItemCreationQuery)).
				WithArgs(interfaceToDriverValue(args)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		markMealPlanOptionAsHavingStepsCreatedArgs := []any{
			exampleMealPlan.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(markMealPlanAsHavingGroceryListInitialized)).
			WithArgs(interfaceToDriverValue(markMealPlanOptionAsHavingStepsCreatedArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		db.ExpectCommit().WillReturnError(errors.New("blah"))

		assert.Error(t, c.CreateMealPlanGroceryListItemsForMealPlan(ctx, exampleMealPlan.ID, inputs))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()
		exampleMealPlanGroceryListItem.PurchasedMeasurementUnit = fakes.BuildFakeValidMeasurementUnit()

		c, db := buildTestClient(t)

		args := []any{
			exampleMealPlanGroceryListItem.BelongsToMealPlan,
			exampleMealPlanGroceryListItem.Ingredient.ID,
			exampleMealPlanGroceryListItem.MeasurementUnit.ID,
			exampleMealPlanGroceryListItem.MinimumQuantityNeeded,
			exampleMealPlanGroceryListItem.MaximumQuantityNeeded,
			exampleMealPlanGroceryListItem.QuantityPurchased,
			exampleMealPlanGroceryListItem.PurchasedMeasurementUnit.ID,
			exampleMealPlanGroceryListItem.PurchasedUPC,
			exampleMealPlanGroceryListItem.PurchasePrice,
			exampleMealPlanGroceryListItem.StatusExplanation,
			exampleMealPlanGroceryListItem.Status,
			exampleMealPlanGroceryListItem.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealPlanGroceryListItemQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateMealPlanGroceryListItem(ctx, exampleMealPlanGroceryListItem))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateMealPlanGroceryListItem(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		c, db := buildTestClient(t)

		var purchasedMeasurementUnitID *string
		if exampleMealPlanGroceryListItem.PurchasedMeasurementUnit != nil {
			purchasedMeasurementUnitID = &exampleMealPlanGroceryListItem.PurchasedMeasurementUnit.ID
		}

		args := []any{
			exampleMealPlanGroceryListItem.BelongsToMealPlan,
			exampleMealPlanGroceryListItem.Ingredient.ID,
			exampleMealPlanGroceryListItem.MeasurementUnit.ID,
			exampleMealPlanGroceryListItem.MinimumQuantityNeeded,
			exampleMealPlanGroceryListItem.MaximumQuantityNeeded,
			exampleMealPlanGroceryListItem.QuantityPurchased,
			purchasedMeasurementUnitID,
			exampleMealPlanGroceryListItem.PurchasedUPC,
			exampleMealPlanGroceryListItem.PurchasePrice,
			exampleMealPlanGroceryListItem.StatusExplanation,
			exampleMealPlanGroceryListItem.Status,
			exampleMealPlanGroceryListItem.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateMealPlanGroceryListItemQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateMealPlanGroceryListItem(ctx, exampleMealPlanGroceryListItem))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveMealPlanGroceryListItem(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		c, db := buildTestClient(t)

		args := []any{
			exampleMealPlanGroceryListItem.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealPlanGroceryListItemQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveMealPlanGroceryListItem(ctx, exampleMealPlanGroceryListItem.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan grocery list item ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanGroceryListItem(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleMealPlanGroceryListItem := fakes.BuildFakeMealPlanGroceryListItem()

		c, db := buildTestClient(t)

		args := []any{
			exampleMealPlanGroceryListItem.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveMealPlanGroceryListItemQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveMealPlanGroceryListItem(ctx, exampleMealPlanGroceryListItem.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
