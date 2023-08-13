package postgres

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
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

func TestQuerier_ValidVesselExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid vessel ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidVesselExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid vessel ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidVessel(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidVessels(T *testing.T) {
	T.Parallel()

	T.Run("with invalid query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidVessels(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
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

		_, db := buildTestClient(t)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidVessel(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidVessel(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidVessel(ctx, nil))
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

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidVesselAsIndexed(ctx, ""))
	})
}
