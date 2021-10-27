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

	database "github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
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
			x.InstrumentID,
			x.RecipeStepID,
			x.Notes,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
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

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepInstrument.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepInstrumentExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.RecipeStepInstrumentExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepInstrument.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepInstrumentExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(sql.ErrNoRows)

		actual, err := c.RecipeStepInstrumentExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepInstrument.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipeStepInstrumentExistenceQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.RecipeStepInstrumentExists(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepInstrument.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepInstruments(false, 0, exampleRecipeStepInstrument))

		actual, err := c.GetRecipeStepInstrument(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrument, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepInstrument.ID,
			exampleRecipeID,
			exampleRecipeStepID,
			exampleRecipeID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipeStepInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepInstrument(ctx, exampleRecipeID, exampleRecipeStepID, exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetTotalRecipeStepInstrumentCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleCount := uint64(123)

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalRecipeStepInstrumentsCountQuery)).
			WithArgs().
			WillReturnRows(newCountDBRowResponse(uint64(123)))

		actual, err := c.GetTotalRecipeStepInstrumentCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		db.ExpectQuery(formatQueryForSQLMock(getTotalRecipeStepInstrumentsCountQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := c.GetTotalRecipeStepInstrumentCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

		mock.AssertExpectationsForObjects(t, db)
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

		query, args := c.buildListQuery(
			ctx,
			"recipe_step_instruments",
			getRecipeStepInstrumentsJoins,
			nil,
			householdOwnershipColumn,
			recipeStepInstrumentsTableColumns,
			"",
			false,
			filter,
		)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepInstruments(true, exampleRecipeStepInstrumentList.FilteredCount, exampleRecipeStepInstrumentList.RecipeStepInstruments...))

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

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrumentList := fakes.BuildFakeRecipeStepInstrumentList()
		exampleRecipeStepInstrumentList.Page = 0
		exampleRecipeStepInstrumentList.Limit = 0

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(
			ctx,
			"recipe_step_instruments",
			getRecipeStepInstrumentsJoins,
			nil,
			householdOwnershipColumn,
			recipeStepInstrumentsTableColumns,
			"",
			false,
			filter,
		)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepInstruments(true, exampleRecipeStepInstrumentList.FilteredCount, exampleRecipeStepInstrumentList.RecipeStepInstruments...))

		actual, err := c.GetRecipeStepInstruments(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleRecipeID := fakes.BuildFakeID()
		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildListQuery(
			ctx,
			"recipe_step_instruments",
			getRecipeStepInstrumentsJoins,
			nil,
			householdOwnershipColumn,
			recipeStepInstrumentsTableColumns,
			"",
			false,
			filter,
		)

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

		query, args := c.buildListQuery(
			ctx,
			"recipe_step_instruments",
			getRecipeStepInstrumentsJoins,
			nil,
			householdOwnershipColumn,
			recipeStepInstrumentsTableColumns,
			"",
			false,
			filter,
		)

		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepInstruments(ctx, exampleRecipeID, exampleRecipeStepID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipeStepInstrumentsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrumentList := fakes.BuildFakeRecipeStepInstrumentList()

		var exampleIDs []string
		for _, x := range exampleRecipeStepInstrumentList.RecipeStepInstruments {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetRecipeStepInstrumentsWithIDsQuery(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipeStepInstruments(false, 0, exampleRecipeStepInstrumentList.RecipeStepInstruments...))

		actual, err := c.GetRecipeStepInstrumentsWithIDs(ctx, exampleRecipeStepID, 0, exampleIDs)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrumentList.RecipeStepInstruments, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid IDs", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstrumentsWithIDs(ctx, exampleRecipeStepID, defaultLimit, nil)
		assert.Error(t, err)
		assert.Empty(t, actual)
	})

	T.Run("with invalid recipe step ID", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepInstrumentList := fakes.BuildFakeRecipeStepInstrumentList()

		var exampleIDs []string
		for _, x := range exampleRecipeStepInstrumentList.RecipeStepInstruments {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipeStepInstrumentsWithIDs(ctx, "", defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrumentList := fakes.BuildFakeRecipeStepInstrumentList()

		var exampleIDs []string
		for _, x := range exampleRecipeStepInstrumentList.RecipeStepInstruments {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetRecipeStepInstrumentsWithIDsQuery(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		actual, err := c.GetRecipeStepInstrumentsWithIDs(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with error scanning query results", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrumentList := fakes.BuildFakeRecipeStepInstrumentList()

		var exampleIDs []string
		for _, x := range exampleRecipeStepInstrumentList.RecipeStepInstruments {
			exampleIDs = append(exampleIDs, x.ID)
		}

		ctx := context.Background()
		c, db := buildTestClient(t)

		query, args := c.buildGetRecipeStepInstrumentsWithIDsQuery(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		db.ExpectQuery(formatQueryForSQLMock(query)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildErroneousMockRow())

		actual, err := c.GetRecipeStepInstrumentsWithIDs(ctx, exampleRecipeStepID, defaultLimit, exampleIDs)
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.ID = "1"
		exampleInput := fakes.BuildFakeRecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.InstrumentID,
			exampleInput.RecipeStepID,
			exampleInput.Notes,
			exampleInput.BelongsToRecipeStep,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepInstrumentCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleRecipeStepInstrument.ID))

		c.timeFunc = func() uint64 {
			return exampleRecipeStepInstrument.CreatedOn
		}

		actual, err := c.CreateRecipeStepInstrument(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrument, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipeStepInstrument(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()
		exampleInput := fakes.BuildFakeRecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.InstrumentID,
			exampleInput.RecipeStepID,
			exampleInput.Notes,
			exampleInput.BelongsToRecipeStep,
		}

		db.ExpectExec(formatQueryForSQLMock(recipeStepInstrumentCreationQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(expectedErr)

		c.timeFunc = func() uint64 {
			return exampleRecipeStepInstrument.CreatedOn
		}

		actual, err := c.CreateRecipeStepInstrument(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepInstrument.InstrumentID,
			exampleRecipeStepInstrument.RecipeStepID,
			exampleRecipeStepInstrument.Notes,
			exampleRecipeStepInstrument.BelongsToRecipeStep,
			exampleRecipeStepInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeStepInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleRecipeStepInstrument.ID))

		assert.NoError(t, c.UpdateRecipeStepInstrument(ctx, exampleRecipeStepInstrument))

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateRecipeStepInstrument(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepInstrument.InstrumentID,
			exampleRecipeStepInstrument.RecipeStepID,
			exampleRecipeStepInstrument.Notes,
			exampleRecipeStepInstrument.BelongsToRecipeStep,
			exampleRecipeStepInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipeStepInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.UpdateRecipeStepInstrument(ctx, exampleRecipeStepInstrument))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeStepInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult(exampleRecipeStepInstrument.ID))

		assert.NoError(t, c.ArchiveRecipeStepInstrument(ctx, exampleRecipeStepID, exampleRecipeStepInstrument.ID))

		mock.AssertExpectationsForObjects(t, db)
	})

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

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleRecipeStepID := fakes.BuildFakeID()
		exampleRecipeStepInstrument := fakes.BuildFakeRecipeStepInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipeStepID,
			exampleRecipeStepInstrument.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipeStepInstrumentQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveRecipeStepInstrument(ctx, exampleRecipeStepID, exampleRecipeStepInstrument.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
