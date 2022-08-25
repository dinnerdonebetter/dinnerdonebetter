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

	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
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
			x.PluralName,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
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
		args := []interface{}{
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
		args := []interface{}{
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
		args := []interface{}{
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

		args := []interface{}{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
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

		args := []interface{}{
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

		args := []interface{}{}

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

		args := []interface{}{}

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

		args := []interface{}{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnits(false, 0, exampleValidMeasurementUnits.ValidMeasurementUnits...))

		actual, err := c.SearchForValidMeasurementUnits(ctx, exampleQuery)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnits.ValidMeasurementUnits, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid ingredient ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidMeasurementUnits(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForValidMeasurementUnits(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning response", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validMeasurementUnitSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForValidMeasurementUnits(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetTotalValidMeasurementUnitCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalValidMeasurementUnitsCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetTotalValidMeasurementUnitCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalValidMeasurementUnitsCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetTotalValidMeasurementUnitCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

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

		query, args := c.buildListQuery(ctx, "valid_measurement_units", nil, nil, nil, householdOwnershipColumn, validMeasurementUnitsTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnits(true, exampleValidMeasurementUnitList.FilteredCount, exampleValidMeasurementUnitList.ValidMeasurementUnits...))

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

		query, args := c.buildListQuery(ctx, "valid_measurement_units", nil, nil, nil, householdOwnershipColumn, validMeasurementUnitsTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnits(true, exampleValidMeasurementUnitList.FilteredCount, exampleValidMeasurementUnitList.ValidMeasurementUnits...))

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

		query, args := c.buildListQuery(ctx, "valid_measurement_units", nil, nil, nil, householdOwnershipColumn, validMeasurementUnitsTableColumns, "", false, filter, true)

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

		query, args := c.buildListQuery(ctx, "valid_measurement_units", nil, nil, nil, householdOwnershipColumn, validMeasurementUnitsTableColumns, "", false, filter, true)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidMeasurementUnits(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidMeasurementUnitsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnitList := fakes.BuildFakeValidMeasurementUnitList()

		var exampleIDs []string
		for _, x := range exampleValidMeasurementUnitList.ValidMeasurementUnits {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetValidMeasurementUnitsWithIDsQuery(ctx, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidMeasurementUnits(false, 0, exampleValidMeasurementUnitList.ValidMeasurementUnits...))

		actual, err := c.GetValidMeasurementUnitsWithIDs(ctx, 0, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidMeasurementUnitList.ValidMeasurementUnits, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid IDs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidMeasurementUnitsWithIDs(ctx, defaultLimit, nil)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnitList := fakes.BuildFakeValidMeasurementUnitList()

		var exampleIDs []string
		for _, x := range exampleValidMeasurementUnitList.ValidMeasurementUnits {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetValidMeasurementUnitsWithIDsQuery(ctx, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidMeasurementUnitsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning query results", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnitList := fakes.BuildFakeValidMeasurementUnitList()

		var exampleIDs []string
		for _, x := range exampleValidMeasurementUnitList.ValidMeasurementUnits {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetValidMeasurementUnitsWithIDsQuery(ctx, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidMeasurementUnitsWithIDs(ctx, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidMeasurementUnit(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
		exampleValidMeasurementUnit.ID = "1"
		exampleInput := fakes.BuildFakeValidMeasurementUnitDatabaseCreationInputFromValidMeasurementUnit(exampleValidMeasurementUnit)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.Volumetric,
			exampleInput.IconPath,
			exampleInput.Universal,
			exampleInput.Metric,
			exampleInput.Imperial,
			exampleInput.PluralName,
		}

		db.ExpectExec(formatQueryForSQLMock(validMeasurementUnitCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() uint64 {
			return exampleValidMeasurementUnit.CreatedOn
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
		exampleInput := fakes.BuildFakeValidMeasurementUnitDatabaseCreationInputFromValidMeasurementUnit(exampleValidMeasurementUnit)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Description,
			exampleInput.Volumetric,
			exampleInput.IconPath,
			exampleInput.Universal,
			exampleInput.Metric,
			exampleInput.Imperial,
			exampleInput.PluralName,
		}

		db.ExpectExec(formatQueryForSQLMock(validMeasurementUnitCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() uint64 {
			return exampleValidMeasurementUnit.CreatedOn
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

		args := []interface{}{
			exampleValidMeasurementUnit.Name,
			exampleValidMeasurementUnit.Description,
			exampleValidMeasurementUnit.Volumetric,
			exampleValidMeasurementUnit.IconPath,
			exampleValidMeasurementUnit.Universal,
			exampleValidMeasurementUnit.Metric,
			exampleValidMeasurementUnit.Imperial,
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

		args := []interface{}{
			exampleValidMeasurementUnit.Name,
			exampleValidMeasurementUnit.Description,
			exampleValidMeasurementUnit.Volumetric,
			exampleValidMeasurementUnit.IconPath,
			exampleValidMeasurementUnit.Universal,
			exampleValidMeasurementUnit.Metric,
			exampleValidMeasurementUnit.Imperial,
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

		args := []interface{}{
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

		args := []interface{}{
			exampleValidMeasurementUnit.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveValidMeasurementUnitQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveValidMeasurementUnit(ctx, exampleValidMeasurementUnit.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
