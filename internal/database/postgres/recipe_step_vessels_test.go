package postgres

import (
	"context"
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

func buildMockRowsFromRecipeStepVessels(includeCounts bool, filteredCount uint64, recipeStepVessels ...*types.RecipeStepVessel) *sqlmock.Rows {
	columns := recipeStepVesselsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepVessels {
		rowValues := []driver.Value{
			&x.ID,
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
			&x.Name,
			&x.Notes,
			&x.BelongsToRecipeStep,
			&x.RecipeStepProductID,
			&x.VesselPreposition,
			&x.MinimumQuantity,
			&x.MaximumQuantity,
			&x.UnavailableAfterStep,
			&x.CreatedAt,
			&x.LastUpdatedAt,
			&x.ArchivedAt,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(recipeStepVessels))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanRecipeStepVessels(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipeStepVessels(ctx, mockRows, false)
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

		_, _, _, err := q.scanRecipeStepVessels(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_RecipeStepVesselExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepVesselExists(ctx, "", exampleRecipeStepID, exampleRecipeStepVessel.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepVesselExists(ctx, exampleRecipeID, "", exampleRecipeStepVessel.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepVesselExists(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepVessel.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepVesselQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepVessels(false, 0, exampleRecipeStepVessel))

		actual, err := c.GetRecipeStepVessel(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepVessel.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepVessel, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepVessel(ctx, "", exampleRecipeStepID, exampleRecipeStepVessel.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepVessel(ctx, exampleRecipeID, "", exampleRecipeStepVessel.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepVessel(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleRecipeStepID,
			exampleRecipeStepVessel.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepVesselQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepVessel(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepVessel.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepVessels(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepVesselList := fakes.BuildFakeRecipeStepVesselList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_vessels", getRecipeStepVesselsJoins, []string{"valid_vessels.id", "valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepVesselsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepVessels(true, exampleRecipeStepVesselList.FilteredCount, exampleRecipeStepVesselList.Data...))

		actual, err := c.GetRecipeStepVessels(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepVesselList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepVessels(ctx, "", exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepVessels(ctx, exampleRecipeID, "", filter)
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

		query, args := c.buildListQuery(ctx, "recipe_step_vessels", getRecipeStepVesselsJoins, []string{"valid_vessels.id", "valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepVesselsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepVessels(ctx, exampleRecipeID, exampleRecipeStepID, filter)
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

		query, args := c.buildListQuery(ctx, "recipe_step_vessels", getRecipeStepVesselsJoins, []string{"valid_vessels.id", "valid_measurement_units.id"}, nil, householdOwnershipColumn, recipeStepVesselsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepVessels(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()
		exampleRecipeStepVessel.ID = "1"
		exampleRecipeStepVessel.Vessel = &types.ValidVessel{ID: exampleRecipeStepVessel.ID}
		exampleInput := converters.ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput(exampleRecipeStepVessel)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Notes,
			exampleInput.BelongsToRecipeStep,
			exampleInput.RecipeStepProductID,
			exampleInput.VesselID,
			exampleInput.VesselPreposition,
			exampleInput.MinimumQuantity,
			exampleInput.MaximumQuantity,
			exampleInput.UnavailableAfterStep,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepVesselCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return exampleRecipeStepVessel.CreatedAt
		}

		actual, err := c.CreateRecipeStepVessel(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepVessel, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepVessel(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()
		exampleInput := converters.ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput(exampleRecipeStepVessel)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []any{
			exampleInput.ID,
			exampleInput.Name,
			exampleInput.Notes,
			exampleInput.BelongsToRecipeStep,
			exampleInput.RecipeStepProductID,
			exampleInput.VesselID,
			exampleInput.VesselPreposition,
			exampleInput.MinimumQuantity,
			exampleInput.MaximumQuantity,
			exampleInput.UnavailableAfterStep,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepVesselCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() time.Time {
			return exampleRecipeStepVessel.CreatedAt
		}

		actual, err := c.CreateRecipeStepVessel(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepVessel(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeStepVessel(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepVessel(ctx, "", exampleRecipeStepVessel.ID))
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepVessel(ctx, exampleRecipeStepID, ""))
	})
}
