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
	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromValidVessels(includeCounts bool, filteredCount uint64, validVessels ...*types.ValidVessel) *sqlmock.Rows {
	columns := validVesselsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validVessels {
		rowValues := []driver.Value{
			&x.ID,
			&x.Name,
			&x.PluralName,
			&x.Description,
			&x.IconPath,
			&x.UsableForStorage,
			&x.Slug,
			&x.DisplayInSummaryLists,
			&x.IncludeInGeneratedInstructions,
			&x.Capacity,
			&x.CapacityUnit.ID,
			&x.CapacityUnit.Name,
			&x.CapacityUnit.Description,
			&x.CapacityUnit.Volumetric,
			&x.CapacityUnit.IconPath,
			&x.CapacityUnit.Universal,
			&x.CapacityUnit.Metric,
			&x.CapacityUnit.Imperial,
			&x.CapacityUnit.Slug,
			&x.CapacityUnit.PluralName,
			&x.CapacityUnit.CreatedAt,
			&x.CapacityUnit.LastUpdatedAt,
			&x.CapacityUnit.ArchivedAt,
			&x.WidthInMillimeters,
			&x.LengthInMillimeters,
			&x.HeightInMillimeters,
			&x.Shape,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validVessels))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidVessels(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidVessels(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidVesselExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidVessel := fakes.BuildFakeValidVessel()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidVessel.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validVesselExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.ValidVesselExists(ctx, exampleValidVessel.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid vessel ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidVesselExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidVessel := fakes.BuildFakeValidVessel()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidVessel.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validVesselExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.ValidVesselExists(ctx, exampleValidVessel.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidVessel := fakes.BuildFakeValidVessel()

		c, db := buildTestClient(t)
		args := []any{
			exampleValidVessel.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(validVesselExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.ValidVesselExists(ctx, exampleValidVessel.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidVessel := fakes.BuildFakeValidVessel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidVessel.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidVesselQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidVessels(false, 0, exampleValidVessel))

		actual, err := c.GetValidVessel(ctx, exampleValidVessel.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidVessel, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid vessel ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidVessel(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidVessel := fakes.BuildFakeValidVessel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidVessel.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidVesselQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidVessel(ctx, exampleValidVessel.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRandomValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidVessel := fakes.BuildFakeValidVessel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{}

		db.ExpectQuery(formatQueryForSQLMock(getRandomValidVesselQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidVessels(false, 0, exampleValidVessel))

		actual, err := c.GetRandomValidVessel(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidVessel, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{}

		db.ExpectQuery(formatQueryForSQLMock(getRandomValidVesselQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRandomValidVessel(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForValidVessels(T *testing.T) {
	T.Parallel()

	exampleQuery := "blah"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidVessels := fakes.BuildFakeValidVesselList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			wrapQueryForILIKE(exampleQuery),
		}

		db.ExpectQuery(formatQueryForSQLMock(validVesselSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidVessels(false, 0, exampleValidVessels.Data...))

		actual, err := c.SearchForValidVessels(ctx, exampleQuery)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidVessels.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid vessel ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidVessels(ctx, "")
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

		db.ExpectQuery(formatQueryForSQLMock(validVesselSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.SearchForValidVessels(ctx, exampleQuery)
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

		db.ExpectQuery(formatQueryForSQLMock(validVesselSearchQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.SearchForValidVessels(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidVesselList := fakes.BuildFakeValidVesselList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			filter.CreatedBefore,
			filter.CreatedAfter,
			filter.UpdatedBefore,
			filter.UpdatedAfter,
			filter.QueryOffset(),
			filter.Limit,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidVesselsQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidVessels(true, exampleValidVesselList.FilteredCount, exampleValidVesselList.Data...))

		actual, err := c.GetValidVessels(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidVesselList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidVesselList := fakes.BuildFakeValidVesselList()
		exampleValidVesselList.Page = 0
		exampleValidVesselList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			filter.CreatedBefore,
			filter.CreatedAfter,
			filter.UpdatedBefore,
			filter.UpdatedAfter,
			filter.QueryOffset(),
			filter.Limit,
		}
		db.ExpectQuery(formatQueryForSQLMock(getValidVesselsQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidVessels(true, exampleValidVesselList.FilteredCount, exampleValidVesselList.Data...))

		actual, err := c.GetValidVessels(ctx, nil)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidVesselList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			filter.CreatedBefore,
			filter.CreatedAfter,
			filter.UpdatedBefore,
			filter.UpdatedAfter,
			filter.QueryOffset(),
			filter.Limit,
		}
		db.ExpectQuery(formatQueryForSQLMock(getValidVesselsQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidVessels(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			filter.CreatedBefore,
			filter.CreatedAfter,
			filter.UpdatedBefore,
			filter.UpdatedAfter,
			filter.QueryOffset(),
			filter.Limit,
		}
		db.ExpectQuery(formatQueryForSQLMock(getValidVesselsQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidVessels(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidVesselsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidVesselList := fakes.BuildFakeValidVesselList()

		exampleIDs := []string{}
		for _, exampleValidVessel := range exampleValidVesselList.Data {
			exampleIDs = append(exampleIDs, exampleValidVessel.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		where := squirrel.Eq{"valid_vessels.id": exampleIDs}
		joins := []string{
			"valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id",
		}
		groupBys := []string{
			"valid_vessels.id",
			"valid_measurement_units.id",
		}
		query, args := c.buildListQuery(ctx, validVesselsTable, joins, groupBys, where, householdOwnershipColumn, validVesselsTableColumns, "", false, nil)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidVessels(true, exampleValidVesselList.FilteredCount, exampleValidVesselList.Data...))

		actual, err := c.GetValidVesselsWithIDs(ctx, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidVesselList.Data, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidVesselThatNeedSearchIndexing(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidVesselList := fakes.BuildFakeValidVesselList()

		c, db := buildTestClient(t)

		exampleIDs := []string{}
		for _, exampleValidVessel := range exampleValidVesselList.Data {
			exampleIDs = append(exampleIDs, exampleValidVessel.ID)
		}

		db.ExpectQuery(formatQueryForSQLMock(validVesselsNeedingIndexingQuery)).
			WithArgs(interfaceToDriverValue(nil)...).
			WillReturnRows(buildMockRowsFromIDs(exampleIDs...))

		actual, err := c.GetValidVesselIDsThatNeedSearchIndexing(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleIDs, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidVessel := fakes.BuildFakeValidVessel()
		exampleValidVessel.CapacityUnit = &types.ValidMeasurementUnit{ID: exampleValidVessel.CapacityUnit.ID}
		exampleInput := converters.ConvertValidVesselToValidVesselDatabaseCreationInput(exampleValidVessel)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.PluralName,
			exampleInput.Description,
			exampleInput.IconPath,
			exampleInput.UsableForStorage,
			exampleInput.Slug,
			exampleInput.DisplayInSummaryLists,
			exampleInput.IncludeInGeneratedInstructions,
			exampleInput.Capacity,
			exampleInput.CapacityUnitID,
			exampleInput.WidthInMillimeters,
			exampleInput.LengthInMillimeters,
			exampleInput.HeightInMillimeters,
			exampleInput.Shape,
		}

		db.ExpectExec(formatQueryForSQLMock(validVesselCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleValidVessel.CreatedAt
		}

		actual, err := c.CreateValidVessel(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidVessel, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidVessel(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleValidVessel := fakes.BuildFakeValidVessel()
		exampleInput := converters.ConvertValidVesselToValidVesselDatabaseCreationInput(exampleValidVessel)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.PluralName,
			exampleInput.Description,
			exampleInput.IconPath,
			exampleInput.UsableForStorage,
			exampleInput.Slug,
			exampleInput.DisplayInSummaryLists,
			exampleInput.IncludeInGeneratedInstructions,
			exampleInput.Capacity,
			exampleInput.CapacityUnitID,
			exampleInput.WidthInMillimeters,
			exampleInput.LengthInMillimeters,
			exampleInput.HeightInMillimeters,
			exampleInput.Shape,
		}

		db.ExpectExec(formatQueryForSQLMock(validVesselCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleValidVessel.CreatedAt
		}

		actual, err := c.CreateValidVessel(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidVessel := fakes.BuildFakeValidVessel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidVessel.Name,
			exampleValidVessel.PluralName,
			exampleValidVessel.Description,
			exampleValidVessel.IconPath,
			exampleValidVessel.UsableForStorage,
			exampleValidVessel.Slug,
			exampleValidVessel.DisplayInSummaryLists,
			exampleValidVessel.IncludeInGeneratedInstructions,
			exampleValidVessel.Capacity,
			exampleValidVessel.CapacityUnit.ID,
			exampleValidVessel.WidthInMillimeters,
			exampleValidVessel.LengthInMillimeters,
			exampleValidVessel.HeightInMillimeters,
			exampleValidVessel.Shape,
			exampleValidVessel.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidVesselQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateValidVessel(ctx, exampleValidVessel))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidVessel(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidVessel := fakes.BuildFakeValidVessel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidVessel.Name,
			exampleValidVessel.PluralName,
			exampleValidVessel.Description,
			exampleValidVessel.IconPath,
			exampleValidVessel.UsableForStorage,
			exampleValidVessel.Slug,
			exampleValidVessel.DisplayInSummaryLists,
			exampleValidVessel.IncludeInGeneratedInstructions,
			exampleValidVessel.Capacity,
			exampleValidVessel.CapacityUnit.ID,
			exampleValidVessel.WidthInMillimeters,
			exampleValidVessel.LengthInMillimeters,
			exampleValidVessel.HeightInMillimeters,
			exampleValidVessel.Shape,
			exampleValidVessel.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidVesselQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateValidVessel(ctx, exampleValidVessel))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid vessel ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidVessel(ctx, ""))
	})
}

func TestQuerier_MarkValidVesselAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidVessel := fakes.BuildFakeValidVessel()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidVessel.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidVesselLastIndexedAtQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.MarkValidVesselAsIndexed(ctx, exampleValidVessel.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidVesselAsIndexed(ctx, ""))
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleValidVessel := fakes.BuildFakeValidVessel()

		c, db := buildTestClient(t)

		args := []any{
			exampleValidVessel.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateValidVesselLastIndexedAtQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.MarkValidVesselAsIndexed(ctx, exampleValidVessel.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
