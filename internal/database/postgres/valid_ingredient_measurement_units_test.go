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

func buildMockRowsFromValidIngredientMeasurementUnits(includeCounts bool, filteredCount uint64, validIngredientMeasurementUnits ...*types.ValidIngredientMeasurementUnit) *sqlmock.Rows {
	columns := validIngredientMeasurementUnitsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validIngredientMeasurementUnits {
		rowValues := []driver.Value{
			&x.ID,
			&x.Notes,
			&x.MeasurementUnit.ID,
			&x.MeasurementUnit.Name,
			&x.MeasurementUnit.Description,
			&x.MeasurementUnit.Volumetric,
			&x.MeasurementUnit.IconPath,
			&x.MeasurementUnit.Universal,
			&x.MeasurementUnit.Metric,
			&x.MeasurementUnit.Imperial,
			&x.MeasurementUnit.PluralName,
			&x.MeasurementUnit.CreatedAt,
			&x.MeasurementUnit.LastUpdatedAt,
			&x.MeasurementUnit.ArchivedAt,
			&x.Ingredient.ID,
			&x.Ingredient.Name,
			&x.Ingredient.Description,
			&x.Ingredient.Warning,
			&x.Ingredient.ContainsEgg,
			&x.Ingredient.ContainsDairy,
			&x.Ingredient.ContainsPeanut,
			&x.Ingredient.ContainsTreeNut,
			&x.Ingredient.ContainsSoy,
			&x.Ingredient.ContainsWheat,
			&x.Ingredient.ContainsShellfish,
			&x.Ingredient.ContainsSesame,
			&x.Ingredient.ContainsFish,
			&x.Ingredient.ContainsGluten,
			&x.Ingredient.AnimalFlesh,
			&x.Ingredient.IsMeasuredVolumetrically,
			&x.Ingredient.IsLiquid,
			&x.Ingredient.IconPath,
			&x.Ingredient.AnimalDerived,
			&x.Ingredient.PluralName,
			&x.Ingredient.RestrictToPreparations,
			&x.Ingredient.MinimumIdealStorageTemperatureInCelsius,
			&x.Ingredient.MaximumIdealStorageTemperatureInCelsius,
			&x.Ingredient.StorageInstructions,
			&x.Ingredient.CreatedAt,
			&x.Ingredient.LastUpdatedAt,
			&x.Ingredient.ArchivedAt,
			&x.MinimumAllowableQuantity,
			&x.MaximumAllowableQuantity,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validIngredientMeasurementUnits))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidIngredientMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidIngredientMeasurementUnits(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidIngredientMeasurementUnits(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidIngredientMeasurementUnitExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidIngredientMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientMeasurementUnitExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidIngredientMeasurementUnitExists(ctx, exampleValidIngredientMeasurementUnit.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient measurement unit ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientMeasurementUnitExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidIngredientMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientMeasurementUnitExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidIngredientMeasurementUnitExists(ctx, exampleValidIngredientMeasurementUnit.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleValidIngredientMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validIngredientMeasurementUnitExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidIngredientMeasurementUnitExists(ctx, exampleValidIngredientMeasurementUnit.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientMeasurementUnits(false, 0, exampleValidIngredientMeasurementUnit))

		actual, err := c.GetValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnit.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientMeasurementUnit, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient measurement unit ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientMeasurementUnit(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidIngredientMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnit.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidIngredientMeasurementUnits(T *testing.T) {
	T.Parallel()

	joins := []string{
		validMeasurementUnitsOnValidIngredientMeasurementUnitsJoinClause,
		validIngredientsOnValidIngredientMeasurementUnitsJoinClause,
	}

	groupBys := []string{
		"valid_ingredients.id",
		"valid_measurement_units.id",
		"valid_ingredient_measurement_units.id",
	}

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidIngredientMeasurementUnitList := fakes.BuildFakeValidIngredientMeasurementUnitList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_measurement_units", joins, groupBys, nil, householdOwnershipColumn, validIngredientMeasurementUnitsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientMeasurementUnits(true, exampleValidIngredientMeasurementUnitList.FilteredCount, exampleValidIngredientMeasurementUnitList.ValidIngredientMeasurementUnits...))

		actual, err := c.GetValidIngredientMeasurementUnits(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientMeasurementUnitList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleValidIngredientMeasurementUnitList := fakes.BuildFakeValidIngredientMeasurementUnitList()
		exampleValidIngredientMeasurementUnitList.Page = 0
		exampleValidIngredientMeasurementUnitList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_measurement_units", joins, groupBys, nil, householdOwnershipColumn, validIngredientMeasurementUnitsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidIngredientMeasurementUnits(true, exampleValidIngredientMeasurementUnitList.FilteredCount, exampleValidIngredientMeasurementUnitList.ValidIngredientMeasurementUnits...))

		actual, err := c.GetValidIngredientMeasurementUnits(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientMeasurementUnitList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_measurement_units", joins, groupBys, nil, householdOwnershipColumn, validIngredientMeasurementUnitsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidIngredientMeasurementUnits(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_ingredient_measurement_units", joins, groupBys, nil, householdOwnershipColumn, validIngredientMeasurementUnitsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidIngredientMeasurementUnits(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
		exampleValidIngredientMeasurementUnit.ID = "1"
		exampleValidIngredientMeasurementUnit.Ingredient = types.ValidIngredient{ID: exampleValidIngredientMeasurementUnit.Ingredient.ID}
		exampleValidIngredientMeasurementUnit.MeasurementUnit = types.ValidMeasurementUnit{ID: exampleValidIngredientMeasurementUnit.MeasurementUnit.ID}

		exampleInput := fakes.BuildFakeValidIngredientMeasurementUnitDatabaseCreationInputFromValidIngredientMeasurementUnit(exampleValidIngredientMeasurementUnit)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.ValidMeasurementUnitID,
			exampleInput.ValidIngredientID,
			exampleInput.MinimumAllowableQuantity,
			exampleInput.MaximumAllowableQuantity,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientMeasurementUnitCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleValidIngredientMeasurementUnit.CreatedAt
		}

		actual, err := c.CreateValidIngredientMeasurementUnit(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidIngredientMeasurementUnit, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredientMeasurementUnit(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()
		exampleValidIngredientMeasurementUnit.Ingredient = types.ValidIngredient{ID: exampleValidIngredientMeasurementUnit.Ingredient.ID}
		exampleValidIngredientMeasurementUnit.MeasurementUnit = types.ValidMeasurementUnit{ID: exampleValidIngredientMeasurementUnit.MeasurementUnit.ID}

		exampleInput := fakes.BuildFakeValidIngredientMeasurementUnitDatabaseCreationInputFromValidIngredientMeasurementUnit(exampleValidIngredientMeasurementUnit)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.ValidMeasurementUnitID,
			exampleInput.ValidIngredientID,
			exampleInput.MinimumAllowableQuantity,
			exampleInput.MaximumAllowableQuantity,
		}

		db.ExpectExec(formatQueryForSQLMock(validIngredientMeasurementUnitCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleValidIngredientMeasurementUnit.CreatedAt
		}

		actual, err := c.CreateValidIngredientMeasurementUnit(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientMeasurementUnit.Notes,
			exampleValidIngredientMeasurementUnit.MeasurementUnit.ID,
			exampleValidIngredientMeasurementUnit.Ingredient.ID,
			exampleValidIngredientMeasurementUnit.MinimumAllowableQuantity,
			exampleValidIngredientMeasurementUnit.MaximumAllowableQuantity,
			exampleValidIngredientMeasurementUnit.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnit))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredientMeasurementUnit(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientMeasurementUnit.Notes,
			exampleValidIngredientMeasurementUnit.MeasurementUnit.ID,
			exampleValidIngredientMeasurementUnit.Ingredient.ID,
			exampleValidIngredientMeasurementUnit.MinimumAllowableQuantity,
			exampleValidIngredientMeasurementUnit.MaximumAllowableQuantity,
			exampleValidIngredientMeasurementUnit.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidIngredientMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnit))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidIngredientMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientMeasurementUnit.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnit.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient measurement unit ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredientMeasurementUnit(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidIngredientMeasurementUnit := fakes.BuildFakeValidIngredientMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleValidIngredientMeasurementUnit.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidIngredientMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidIngredientMeasurementUnit(ctx, exampleValidIngredientMeasurementUnit.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
