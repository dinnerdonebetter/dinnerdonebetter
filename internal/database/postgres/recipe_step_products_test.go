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

func buildMockRowsFromRecipeStepProducts(includeCounts bool, filteredCount uint64, recipeStepProducts ...*types.RecipeStepProduct) *sqlmock.Rows {
	columns := recipeStepProductsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepProducts {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Type,
			x.MeasurementUnit.ID,
			x.MeasurementUnit.Name,
			x.MeasurementUnit.Description,
			x.MeasurementUnit.Volumetric,
			x.MeasurementUnit.IconPath,
			x.MeasurementUnit.Universal,
			x.MeasurementUnit.Metric,
			x.MeasurementUnit.Imperial,
			x.MeasurementUnit.Slug,
			x.MeasurementUnit.PluralName,
			x.MeasurementUnit.CreatedAt,
			x.MeasurementUnit.LastUpdatedAt,
			x.MeasurementUnit.ArchivedAt,
			x.MinimumQuantity,
			x.MaximumQuantity,
			x.QuantityNotes,
			x.Compostable,
			x.MaximumStorageDurationInSeconds,
			x.MinimumStorageTemperatureInCelsius,
			x.MaximumStorageTemperatureInCelsius,
			x.StorageInstructions,
			x.IsLiquid,
			x.IsWaste,
			x.Index,
			x.ContainedInVesselIndex,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
			x.BelongsToRecipeStep,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(recipeStepProducts))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipeStepProducts(ctx, mockRows, false)
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

		_, _, _, err := q.scanRecipeStepProducts(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_RecipeStepProductExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		c, db := buildTestClient(t)
		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepProductExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.RecipeStepProductExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepProductExists(ctx, "", exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepProductExists(ctx, exampleRecipeID, "", exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step product ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepProductExists(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		c, db := buildTestClient(t)
		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepProductExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.RecipeStepProductExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		c, db := buildTestClient(t)
		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepProductExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.RecipeStepProductExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepProductQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepProducts(false, 0, exampleRecipeStepProduct))

		actual, err := c.GetRecipeStepProduct(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProduct, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepProduct(ctx, "", exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepProduct(ctx, exampleRecipeID, "", exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step product ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepProduct(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepProduct.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepProductQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepProduct(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepProduct.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_getRecipeStepProductsForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepProductList := fakes.BuildFakeRecipeStepProductList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepProductsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepProducts(false, exampleRecipeStepProductList.FilteredCount, exampleRecipeStepProductList.Data...))

		actual, err := c.getRecipeStepProductsForRecipe(ctx, exampleRecipeID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProductList.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.getRecipeStepProductsForRecipe(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepProductsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.getRecipeStepProductsForRecipe(ctx, exampleRecipeID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning query results", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeID,
		}
		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepProductsForRecipeQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildInvalidMockRowsFromListOfIDs([]string{"whatever"}))

		actual, err := c.getRecipeStepProductsForRecipe(ctx, exampleRecipeID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepProducts(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepProductList := fakes.BuildFakeRecipeStepProductList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_products", getRecipeStepProductsJoins, []string{"valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepProductsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepProducts(true, exampleRecipeStepProductList.FilteredCount, exampleRecipeStepProductList.Data...))

		actual, err := c.GetRecipeStepProducts(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProductList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepProducts(ctx, "", exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepProducts(ctx, exampleRecipeID, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_products", getRecipeStepProductsJoins, []string{"valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepProductsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepProducts(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_products", getRecipeStepProductsJoins, []string{"valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepProductsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepProducts(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleRecipeStepProduct.ID = "1"
		exampleInput := converters.ConvertRecipeStepProductToRecipeStepProductDatabaseCreationInput(exampleRecipeStepProduct)
		exampleRecipeStepProduct.MeasurementUnit = &types.ValidMeasurementUnit{ID: exampleRecipeStepProduct.MeasurementUnit.ID}

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Type,
			exampleInput.MeasurementUnitID,
			exampleInput.MinimumQuantity,
			exampleInput.MaximumQuantity,
			exampleInput.QuantityNotes,
			exampleInput.Compostable,
			exampleInput.MaximumStorageDurationInSeconds,
			exampleInput.MinimumStorageTemperatureInCelsius,
			exampleInput.MaximumStorageTemperatureInCelsius,
			exampleInput.StorageInstructions,
			exampleInput.BelongsToRecipeStep,
			exampleInput.IsLiquid,
			exampleInput.IsWaste,
			exampleRecipeStepProduct.Index,
			exampleRecipeStepProduct.ContainedInVesselIndex,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepProductCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleRecipeStepProduct.CreatedAt
		}

		actual, err := c.CreateRecipeStepProduct(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepProduct, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepProduct(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
		exampleInput := converters.ConvertRecipeStepProductToRecipeStepProductDatabaseCreationInput(exampleRecipeStepProduct)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Type,
			exampleInput.MeasurementUnitID,
			exampleInput.MinimumQuantity,
			exampleInput.MaximumQuantity,
			exampleInput.QuantityNotes,
			exampleInput.Compostable,
			exampleInput.MaximumStorageDurationInSeconds,
			exampleInput.MinimumStorageTemperatureInCelsius,
			exampleInput.MaximumStorageTemperatureInCelsius,
			exampleInput.StorageInstructions,
			exampleInput.BelongsToRecipeStep,
			exampleInput.IsLiquid,
			exampleInput.IsWaste,
			exampleRecipeStepProduct.Index,
			exampleRecipeStepProduct.ContainedInVesselIndex,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepProductCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleRecipeStepProduct.CreatedAt
		}

		actual, err := c.CreateRecipeStepProduct(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepProduct.Name,
			exampleRecipeStepProduct.Type,
			exampleRecipeStepProduct.MeasurementUnit.ID,
			exampleRecipeStepProduct.MinimumQuantity,
			exampleRecipeStepProduct.MaximumQuantity,
			exampleRecipeStepProduct.QuantityNotes,
			exampleRecipeStepProduct.Compostable,
			exampleRecipeStepProduct.MaximumStorageDurationInSeconds,
			exampleRecipeStepProduct.MinimumStorageTemperatureInCelsius,
			exampleRecipeStepProduct.MaximumStorageTemperatureInCelsius,
			exampleRecipeStepProduct.StorageInstructions,
			exampleRecipeStepProduct.IsLiquid,
			exampleRecipeStepProduct.IsWaste,
			exampleRecipeStepProduct.Index,
			exampleRecipeStepProduct.ContainedInVesselIndex,
			exampleRecipeStepProduct.BelongsToRecipeStep,
			exampleRecipeStepProduct.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeStepProductQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateRecipeStepProduct(ctx, exampleRecipeStepProduct))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepProduct(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepProduct.Name,
			exampleRecipeStepProduct.Type,
			exampleRecipeStepProduct.MeasurementUnit.ID,
			exampleRecipeStepProduct.MinimumQuantity,
			exampleRecipeStepProduct.MaximumQuantity,
			exampleRecipeStepProduct.QuantityNotes,
			exampleRecipeStepProduct.Compostable,
			exampleRecipeStepProduct.MaximumStorageDurationInSeconds,
			exampleRecipeStepProduct.MinimumStorageTemperatureInCelsius,
			exampleRecipeStepProduct.MaximumStorageTemperatureInCelsius,
			exampleRecipeStepProduct.StorageInstructions,
			exampleRecipeStepProduct.IsLiquid,
			exampleRecipeStepProduct.IsWaste,
			exampleRecipeStepProduct.Index,
			exampleRecipeStepProduct.ContainedInVesselIndex,
			exampleRecipeStepProduct.BelongsToRecipeStep,
			exampleRecipeStepProduct.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeStepProductQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateRecipeStepProduct(ctx, exampleRecipeStepProduct))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveRecipeStepProduct(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepProduct(ctx, "", exampleRecipeStepProduct.ID))
	})

	T.Run("with invalid recipe step product ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepProduct(ctx, exampleRecipeStepID, ""))
	})
}
