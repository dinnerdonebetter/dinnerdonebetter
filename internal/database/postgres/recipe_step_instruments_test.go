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

func buildMockRowsFromRecipeStepInstruments(includeCounts bool, filteredCount uint64, recipeStepInstruments ...*types.RecipeStepInstrument) *sqlmock.Rows {
	columns := recipeStepInstrumentsTableColumns

	if includeCounts {
		columns = append(columns, "filtered_count", "total_count")
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepInstruments {
		rowValues := []driver.Value{
			x.ID,
			x.Instrument.ID,
			x.Instrument.Name,
			x.Instrument.PluralName,
			x.Instrument.Description,
			x.Instrument.IconPath,
			x.Instrument.UsableForStorage,
			x.Instrument.DisplayInSummaryLists,
			x.Instrument.IncludeInGeneratedInstructions,
			x.Instrument.Slug,
			x.Instrument.CreatedAt,
			x.Instrument.LastUpdatedAt,
			x.Instrument.ArchivedAt,
			x.RecipeStepProductID,
			x.Name,
			x.Notes,
			x.PreferenceRank,
			&x.Optional,
			&x.MinimumQuantity,
			&x.MaximumQuantity,
			&x.OptionIndex,
			x.CreatedAt,
			x.LastUpdatedAt,
			x.ArchivedAt,
			x.BelongsToRecipeStep,
		}

		if includeCounts {
			rowValues = append(rowValues, filteredCount, len(recipeStepInstruments))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func TestQuerier_ScanRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errs", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		q, _ := buildTestClient(t)

		mockRows := &database.MockResultIterator{}
		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, _, err := q.scanRecipeStepInstruments(ctx, mockRows, false)
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

		_, _, _, err := q.scanRecipeStepInstruments(ctx, mockRows, false)
		assert.Error(t, err)
	})
}

func TestQuerier_RecipeStepInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepInstrumentExists(ctx, "", exampleRecipeStepID, exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepInstrumentExists(ctx, exampleRecipeID, "", exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.RecipeStepInstrumentExists(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstrument(ctx, "", exampleRecipeStepID, exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstrument(ctx, exampleRecipeID, "", exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstrument(ctx, exampleRecipeID, exampleRecipeStepID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrumentList := fakes.BuildFakeRecipeStepInstrumentList()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(ctx, "recipe_step_instruments", getRecipeStepInstrumentsJoins, []string{"valid_instruments.id", "recipe_step_instruments.id"}, nil, householdOwnershipColumn, recipeStepInstrumentsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepInstruments(true, exampleRecipeStepInstrumentList.FilteredCount, exampleRecipeStepInstrumentList.Data...))

		actual, err := c.GetRecipeStepInstruments(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstruments(ctx, "", exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstruments(ctx, exampleRecipeID, "", filter)
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

		query, args := c.buildListQuery(ctx, "recipe_step_instruments", getRecipeStepInstrumentsJoins, []string{"valid_instruments.id", "recipe_step_instruments.id"}, nil, householdOwnershipColumn, recipeStepInstrumentsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepInstruments(ctx, exampleRecipeID, exampleRecipeStepID, filter)
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

		query, args := c.buildListQuery(ctx, "recipe_step_instruments", getRecipeStepInstrumentsJoins, []string{"valid_instruments.id", "recipe_step_instruments.id"}, nil, householdOwnershipColumn, recipeStepInstrumentsTableColumns, "", false, filter)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepInstruments(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepInstrument(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepInstrument(ctx, nil))
	})
}

func TestQuerier_ArchiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepInstrument(ctx, "", exampleRecipeStepInstrument.ID))
	})

	T.Run("with invalid recipe step instrument ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipeStepInstrument(ctx, exampleRecipeStepID, ""))
	})
}
