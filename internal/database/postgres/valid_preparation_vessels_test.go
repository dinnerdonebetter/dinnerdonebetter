package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromValidPreparationVessels(includeCounts bool, filteredCount uint64, validPreparationVessels ...*types.ValidPreparationVessel) *sqlmock.Rows {
	columns := fullValidPreparationVesselsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range validPreparationVessels {
		rowValues := []driver.Value{
			&x.ID,
			&x.Notes,
			&x.Preparation.ID,
			&x.Preparation.Name,
			&x.Preparation.Description,
			&x.Preparation.IconPath,
			&x.Preparation.YieldsNothing,
			&x.Preparation.RestrictToIngredients,
			&x.Preparation.MinimumIngredientCount,
			&x.Preparation.MaximumIngredientCount,
			&x.Preparation.MinimumInstrumentCount,
			&x.Preparation.MaximumInstrumentCount,
			&x.Preparation.TemperatureRequired,
			&x.Preparation.TimeEstimateRequired,
			&x.Preparation.ConditionExpressionRequired,
			&x.Preparation.ConsumesVessel,
			&x.Preparation.OnlyForVessels,
			&x.Preparation.MinimumVesselCount,
			&x.Preparation.MaximumVesselCount,
			&x.Preparation.Slug,
			&x.Preparation.PastTense,
			&x.Preparation.CreatedAt,
			&x.Preparation.LastUpdatedAt,
			&x.Preparation.ArchivedAt,
			&x.Vessel.ID,
			&x.Vessel.Name,
			&x.Vessel.PluralName,
			&x.Vessel.Description,
			&x.Vessel.IconPath,
			&x.Vessel.UsableForStorage,
			&x.Vessel.Slug,
			&x.Vessel.DisplayInSummaryLists,
			&x.Vessel.IncludeInGeneratedInstructions,
			&x.Vessel.Capacity,
			&x.Vessel.CapacityUnit.ID,
			&x.Vessel.CapacityUnit.Name,
			&x.Vessel.CapacityUnit.Description,
			&x.Vessel.CapacityUnit.Volumetric,
			&x.Vessel.CapacityUnit.IconPath,
			&x.Vessel.CapacityUnit.Universal,
			&x.Vessel.CapacityUnit.Metric,
			&x.Vessel.CapacityUnit.Imperial,
			&x.Vessel.CapacityUnit.Slug,
			&x.Vessel.CapacityUnit.PluralName,
			&x.Vessel.CapacityUnit.CreatedAt,
			&x.Vessel.CapacityUnit.LastUpdatedAt,
			&x.Vessel.CapacityUnit.ArchivedAt,
			&x.Vessel.WidthInMillimeters,
			&x.Vessel.LengthInMillimeters,
			&x.Vessel.HeightInMillimeters,
			&x.Vessel.Shape,
			&x.Vessel.CreatedAt,
			&x.Vessel.LastUpdatedAt,
			&x.Vessel.ArchivedAt,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(validPreparationVessels))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanValidPreparationVessels(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanValidPreparationVessels(ctx, mockRows, false)
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

		_, _, _, err := q.scanValidPreparationVessels(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_ValidPreparationVesselExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidPreparationVesselExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidPreparationVessel.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidPreparationVesselQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparationVessels(false, 0, exampleValidPreparationVessel))

		actual, err := c.GetValidPreparationVessel(ctx, exampleValidPreparationVessel.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationVessel, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidPreparationVessel(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidPreparationVessel := fakes.BuildFakeValidPreparationVessel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleValidPreparationVessel.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getValidPreparationVesselQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidPreparationVessel(ctx, exampleValidPreparationVessel.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidPreparationVessels(T *testing.T) {
	T.Parallel()

	joins := []string{
		validInstrumentsOnValidPreparationVesselsJoinClause,
		validPreparationsOnValidPreparationVesselsJoinClause,
		validMeasurementUnitsOnValidVesselsJoinClause,
	}

	groupBys := []string{
		"valid_preparations.id",
		"valid_vessels.id",
		"valid_preparation_vessels.id",
		"valid_measurement_units.id",
	}

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidPreparationVesselList := fakes.BuildFakeValidPreparationVesselList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparation_vessels", joins, groupBys, nil, householdOwnershipColumn, fullValidPreparationVesselsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparationVessels(true, exampleValidPreparationVesselList.FilteredCount, exampleValidPreparationVesselList.Data...))

		actual, err := c.GetValidPreparationVessels(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationVesselList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleValidPreparationVesselList := fakes.BuildFakeValidPreparationVesselList()
		exampleValidPreparationVesselList.Page = 1
		exampleValidPreparationVesselList.Limit = 20

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparation_vessels", joins, groupBys, nil, householdOwnershipColumn, fullValidPreparationVesselsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromValidPreparationVessels(true, exampleValidPreparationVesselList.FilteredCount, exampleValidPreparationVesselList.Data...))

		actual, err := c.GetValidPreparationVessels(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidPreparationVesselList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparation_vessels", joins, groupBys, nil, householdOwnershipColumn, fullValidPreparationVesselsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetValidPreparationVessels(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with erroneous response from database", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "valid_preparation_vessels", joins, groupBys, nil, householdOwnershipColumn, fullValidPreparationVesselsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetValidPreparationVessels(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidPreparationVessel(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidPreparationVessel(ctx, nil))
	})
}

func TestQuerier_ArchiveValidPreparationVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidPreparationVessel(ctx, ""))
	})
}

func TestQuerier_buildGetValidPreparationVesselsRestrictedByIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		expected := `SELECT valid_preparation_vessels.id, valid_preparation_vessels.notes, valid_preparations.id, valid_preparations.name, valid_preparations.description, valid_preparations.icon_path, valid_preparations.yields_nothing, valid_preparations.restrict_to_ingredients, valid_preparations.minimum_ingredient_count, valid_preparations.maximum_ingredient_count, valid_preparations.minimum_instrument_count, valid_preparations.maximum_instrument_count, valid_preparations.temperature_required, valid_preparations.time_estimate_required, valid_preparations.condition_expression_required, valid_preparations.consumes_vessel, valid_preparations.only_for_vessels, valid_preparations.minimum_vessel_count, valid_preparations.maximum_vessel_count, valid_preparations.slug, valid_preparations.past_tense, valid_preparations.created_at, valid_preparations.last_updated_at, valid_preparations.archived_at, valid_vessels.id, valid_vessels.name, valid_vessels.plural_name, valid_vessels.description, valid_vessels.icon_path, valid_vessels.usable_for_storage, valid_vessels.slug, valid_vessels.display_in_summary_lists, valid_vessels.include_in_generated_instructions, valid_vessels.capacity, valid_measurement_units.id, valid_measurement_units.name, valid_measurement_units.description, valid_measurement_units.volumetric, valid_measurement_units.icon_path, valid_measurement_units.universal, valid_measurement_units.metric, valid_measurement_units.imperial, valid_measurement_units.slug, valid_measurement_units.plural_name, valid_measurement_units.created_at, valid_measurement_units.last_updated_at, valid_measurement_units.archived_at, valid_vessels.width_in_millimeters, valid_vessels.length_in_millimeters, valid_vessels.height_in_millimeters, valid_vessels.shape, valid_vessels.created_at, valid_vessels.last_updated_at, valid_vessels.archived_at, valid_preparation_vessels.created_at, valid_preparation_vessels.last_updated_at, valid_preparation_vessels.archived_at FROM valid_preparation_vessels JOIN valid_vessels ON valid_preparation_vessels.valid_vessel_id = valid_vessels.id JOIN valid_preparations ON valid_preparation_vessels.valid_preparation_id = valid_preparations.id LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit = valid_measurement_units.id WHERE valid_preparation_vessels.archived_at IS NULL AND valid_preparation_vessels.valid_vessel_id IN ($1) LIMIT 20`
		actual, args := c.buildGetValidPreparationVesselsRestrictedByIDsQuery(ctx, "valid_vessel_id", 20, []string{"things"})

		assert.Len(t, args, 1)
		assert.Equal(t, expected, actual)
	})
}
