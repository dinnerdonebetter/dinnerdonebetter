package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	"github.com/prixfixeco/backend/internal/database"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromValidMeasurementUnits(includeCounts bool, filteredCount uint64, validMeasurementUnits ...*types.ValidMeasurementUnit) *sqlmock.Rows {
	columns := validMeasurementUnitsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validMeasurementUnits {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.Description,
			x.Volumetric,
			x.IconPath,
			x.Universal,
			x.Metric,
			x.Imperial,
			x.Slug,
			x.PluralName,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validMeasurementUnits))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidMeasurementUnits(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidMeasurementUnits(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidMeasurementUnitExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidMeasurementUnitExists(ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidMeasurementUnitExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidMeasurementUnitExists(ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidMeasurementUnitExists(ctx, exampleValidMeasurementUnit.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		getValidMeasurementUnitArgs := []any{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(getValidMeasurementUnitArgs)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnits(false, 0, exampleValidMeasurementUnit))

		actual, err := c.GetValidMeasurementUnit(ctx, exampleValidMeasurementUnit.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnit, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidMeasurementUnit(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidMeasurementUnit(ctx, exampleValidMeasurementUnit.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRandomValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{}

		db.ExpectQuery(formatQueryForSQLMock(getRandomValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnits(false, 0, exampleValidMeasurementUnit))

		actual, err := c.GetRandomValidMeasurementUnit(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnit, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{}

		db.ExpectQuery(formatQueryForSQLMock(getRandomValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRandomValidMeasurementUnit(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForValidMeasurementUnits(T *testing.T) {
	T.Parallel()

	exampleQuery := "blah"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnits := fakes.BuildFakeValidMeasurementUnitList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnits(false, 0, exampleValidMeasurementUnits.Data...))

		actual, err := c.SearchForValidMeasurementUnitsByName(ctx, exampleQuery)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnits.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidMeasurementUnitsByName(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForValidMeasurementUnitsByName(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForValidMeasurementUnitsByName(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ValidMeasurementUnitsForIngredientID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		exampleValidIngredientID := fakes.BuildFakeID()
		exampleValidMeasurementUnits := fakes.BuildFakeValidMeasurementUnitList()

		c, db := buildTestClient(t)

		args := []any{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			exampleValidIngredientID,
			filter.Limit,
			filter.QueryOffset(),
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitSearchByIngredientIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnits(true, exampleValidMeasurementUnits.FilteredCount, exampleValidMeasurementUnits.Data...))

		actual, err := c.ValidMeasurementUnitsForIngredientID(ctx, exampleValidIngredientID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnits, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		c, _ := buildTestClient(t)

		actual, err := c.ValidMeasurementUnitsForIngredientID(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		exampleValidIngredientID := fakes.BuildFakeID()
		c, db := buildTestClient(t)

		args := []any{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			exampleValidIngredientID,
			filter.Limit,
			filter.QueryOffset(),
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitSearchByIngredientIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidMeasurementUnitsForIngredientID(ctx, exampleValidIngredientID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		filter := types.DefaultQueryFilter()
		exampleValidIngredientID := fakes.BuildFakeID()
		c, db := buildTestClient(t)

		args := []any{
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
			exampleValidIngredientID,
			filter.Limit,
			filter.QueryOffset(),
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitSearchByIngredientIDQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.ValidMeasurementUnitsForIngredientID(ctx, exampleValidIngredientID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidMeasurementUnits(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidMeasurementUnitList := fakes.BuildFakeValidMeasurementUnitList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_measurement_units", nil, nil, nil, householdOwnershipColumn, validMeasurementUnitsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnits(true, exampleValidMeasurementUnitList.FilteredCount, exampleValidMeasurementUnitList.Data...))

		actual, err := c.GetValidMeasurementUnits(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnitList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleValidMeasurementUnitList := fakes.BuildFakeValidMeasurementUnitList()
		exampleValidMeasurementUnitList.Page = 0
		exampleValidMeasurementUnitList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_measurement_units", nil, nil, nil, householdOwnershipColumn, validMeasurementUnitsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnits(true, exampleValidMeasurementUnitList.FilteredCount, exampleValidMeasurementUnitList.Data...))

		actual, err := c.GetValidMeasurementUnits(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnitList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_measurement_units", nil, nil, nil, householdOwnershipColumn, validMeasurementUnitsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidMeasurementUnits(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_measurement_units", nil, nil, nil, householdOwnershipColumn, validMeasurementUnitsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidMeasurementUnits(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
		exampleValidMeasurementUnit.ID = "1"
		exampleInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitDatabaseCreationInput(exampleValidMeasurementUnit)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.Volumetric,
			exampleInput.IconPath,
			exampleInput.Universal,
			exampleInput.Metric,
			exampleInput.Imperial,
			exampleInput.PluralName,
			exampleInput.Slug,
		}

		db.ExpectExec(formatQueryForSQLMock(validMeasurementUnitCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleValidMeasurementUnit.CreatedAt
		}

		actual, err := c.CreateValidMeasurementUnit(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnit, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidMeasurementUnit(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
		exampleInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitDatabaseCreationInput(exampleValidMeasurementUnit)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.Volumetric,
			exampleInput.IconPath,
			exampleInput.Universal,
			exampleInput.Metric,
			exampleInput.Imperial,
			exampleInput.PluralName,
			exampleInput.Slug,
		}

		db.ExpectExec(formatQueryForSQLMock(validMeasurementUnitCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleValidMeasurementUnit.CreatedAt
		}

		actual, err := c.CreateValidMeasurementUnit(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidMeasurementUnit.Name,
			exampleValidMeasurementUnit.Description,
			exampleValidMeasurementUnit.Volumetric,
			exampleValidMeasurementUnit.IconPath,
			exampleValidMeasurementUnit.Universal,
			exampleValidMeasurementUnit.Metric,
			exampleValidMeasurementUnit.Imperial,
			exampleValidMeasurementUnit.Slug,
			exampleValidMeasurementUnit.PluralName,
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidMeasurementUnit(ctx, exampleValidMeasurementUnit))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidMeasurementUnit(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidMeasurementUnit.Name,
			exampleValidMeasurementUnit.Description,
			exampleValidMeasurementUnit.Volumetric,
			exampleValidMeasurementUnit.IconPath,
			exampleValidMeasurementUnit.Universal,
			exampleValidMeasurementUnit.Metric,
			exampleValidMeasurementUnit.Imperial,
			exampleValidMeasurementUnit.Slug,
			exampleValidMeasurementUnit.PluralName,
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidMeasurementUnit(ctx, exampleValidMeasurementUnit))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveValidMeasurementUnit(ctx, exampleValidMeasurementUnit.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidMeasurementUnit(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidMeasurementUnit(ctx, exampleValidMeasurementUnit.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
