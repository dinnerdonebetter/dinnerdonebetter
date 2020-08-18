package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowsFromRecipeStepInstruments(recipeStepInstruments ...*models.RecipeStepInstrument) *sqlmock.Rows {
	columns := recipeStepInstrumentsTableColumns

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

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipeStepInstrument(x *models.RecipeStepInstrument) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepInstrumentsTableColumns).AddRow(
		x.ArchivedOn,
		x.InstrumentID,
		x.RecipeStepID,
		x.Notes,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.BelongsToRecipeStep,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := p.scanRecipeStepInstruments(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanRecipeStepInstruments(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRecipeStepInstrumentExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "SELECT EXISTS ( SELECT recipe_step_instruments.id FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_step_instruments.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5 )"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepInstrument.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildRecipeStepInstrumentExistsQuery(exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RecipeStepInstrumentExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT recipe_step_instruments.id FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_step_instruments.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepInstrument.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RecipeStepInstrumentExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepInstrument.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RecipeStepInstrumentExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_step_instruments.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepInstrument.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepInstrumentQuery(exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_step_instruments.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepInstrument.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromRecipeStepInstruments(exampleRecipeStepInstrument))

		actual, err := p.GetRecipeStepInstrument(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrument, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepInstrument.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepInstrument(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeStepInstrumentsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(recipe_step_instruments.id) FROM recipe_step_instruments WHERE recipe_step_instruments.archived_on IS NULL"
		actualQuery := p.buildGetAllRecipeStepInstrumentsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeStepInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_step_instruments.id) FROM recipe_step_instruments WHERE recipe_step_instruments.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeStepInstrumentsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfRecipeStepInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM recipe_step_instruments WHERE recipe_step_instruments.id > $1 AND recipe_step_instruments.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfRecipeStepInstrumentsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(recipe_step_instruments.id) FROM recipe_step_instruments WHERE recipe_step_instruments.archived_on IS NULL"
	expectedGetQuery := "SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM recipe_step_instruments WHERE recipe_step_instruments.id > $1 AND recipe_step_instruments.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeStepInstrumentList := fakemodels.BuildFakeRecipeStepInstrumentList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepInstruments(
					&exampleRecipeStepInstrumentList.RecipeStepInstruments[0],
					&exampleRecipeStepInstrumentList.RecipeStepInstruments[1],
					&exampleRecipeStepInstrumentList.RecipeStepInstruments[2],
				),
			)

		out := make(chan []models.RecipeStepInstrument)
		doneChan := make(chan bool, 1)

		err := p.GetAllRecipeStepInstruments(ctx, out)
		assert.NoError(t, err)

		var stillQuerying = true
		for stillQuerying {
			select {
			case batch := <-out:
				assert.NotEmpty(t, batch)
				doneChan <- true
			case <-time.After(time.Second):
				t.FailNow()
			case <-doneChan:
				stillQuerying = false
			}
		}

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error fetching initial count", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		out := make(chan []models.RecipeStepInstrument)

		err := p.GetAllRecipeStepInstruments(ctx, out)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows returned", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(sql.ErrNoRows)

		out := make(chan []models.RecipeStepInstrument)

		err := p.GetAllRecipeStepInstruments(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnError(errors.New("blah"))

		out := make(chan []models.RecipeStepInstrument)

		err := p.GetAllRecipeStepInstruments(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepInstrument(exampleRecipeStepInstrument))

		out := make(chan []models.RecipeStepInstrument)

		err := p.GetAllRecipeStepInstruments(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.archived_on IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 AND recipe_step_instruments.created_on > $5 AND recipe_step_instruments.created_on < $6 AND recipe_step_instruments.last_updated_on > $7 AND recipe_step_instruments.last_updated_on < $8 ORDER BY recipe_step_instruments.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepInstrumentsQuery(exampleRecipe.ID, exampleRecipeStep.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.archived_on IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 ORDER BY recipe_step_instruments.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepInstrumentList := fakemodels.BuildFakeRecipeStepInstrumentList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepInstruments(
					&exampleRecipeStepInstrumentList.RecipeStepInstruments[0],
					&exampleRecipeStepInstrumentList.RecipeStepInstruments[1],
					&exampleRecipeStepInstrumentList.RecipeStepInstruments[2],
				),
			)

		actual, err := p.GetRecipeStepInstruments(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrumentList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepInstruments(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeStepInstruments(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step instrument", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepInstrument(exampleRecipeStepInstrument))

		actual, err := p.GetRecipeStepInstruments(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepInstrumentsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM (SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_instruments.archived_on IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_instruments WHERE recipe_step_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepInstrumentsWithIDsQuery(exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepInstrumentsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleRecipeStepInstrumentList := fakemodels.BuildFakeRecipeStepInstrumentList()
		var exampleIDs []uint64
		for _, recipeStepInstrument := range exampleRecipeStepInstrumentList.RecipeStepInstruments {
			exampleIDs = append(exampleIDs, recipeStepInstrument.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM (SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_instruments.archived_on IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_instruments WHERE recipe_step_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepInstruments(
					&exampleRecipeStepInstrumentList.RecipeStepInstruments[0],
					&exampleRecipeStepInstrumentList.RecipeStepInstruments[1],
					&exampleRecipeStepInstrumentList.RecipeStepInstruments[2],
				),
			)

		actual, err := p.GetRecipeStepInstrumentsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrumentList.RecipeStepInstruments, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM (SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_instruments.archived_on IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_instruments WHERE recipe_step_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepInstrumentsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM (SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_instruments.archived_on IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_instruments WHERE recipe_step_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeStepInstrumentsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step instrument", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM (SELECT recipe_step_instruments.id, recipe_step_instruments.instrument_id, recipe_step_instruments.recipe_step_id, recipe_step_instruments.notes, recipe_step_instruments.created_on, recipe_step_instruments.last_updated_on, recipe_step_instruments.archived_on, recipe_step_instruments.belongs_to_recipe_step FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_step_instruments.archived_on IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_instruments WHERE recipe_step_instruments.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepInstrument(exampleRecipeStepInstrument))

		actual, err := p.GetRecipeStepInstrumentsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeStepInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "INSERT INTO recipe_step_instruments (instrument_id,recipe_step_id,notes,belongs_to_recipe_step) VALUES ($1,$2,$3,$4) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipeStepInstrument.InstrumentID,
			exampleRecipeStepInstrument.RecipeStepID,
			exampleRecipeStepInstrument.Notes,
			exampleRecipeStepInstrument.BelongsToRecipeStep,
		}
		actualQuery, actualArgs := p.buildCreateRecipeStepInstrumentQuery(exampleRecipeStepInstrument)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO recipe_step_instruments (instrument_id,recipe_step_id,notes,belongs_to_recipe_step) VALUES ($1,$2,$3,$4) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRecipeStepInstrument.ID, exampleRecipeStepInstrument.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStepInstrument.InstrumentID,
				exampleRecipeStepInstrument.RecipeStepID,
				exampleRecipeStepInstrument.Notes,
				exampleRecipeStepInstrument.BelongsToRecipeStep,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeStepInstrument(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepInstrument, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStepInstrument.InstrumentID,
				exampleRecipeStepInstrument.RecipeStepID,
				exampleRecipeStepInstrument.Notes,
				exampleRecipeStepInstrument.BelongsToRecipeStep,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipeStepInstrument(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeStepInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "UPDATE recipe_step_instruments SET instrument_id = $1, recipe_step_id = $2, notes = $3, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_step = $4 AND id = $5 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleRecipeStepInstrument.InstrumentID,
			exampleRecipeStepInstrument.RecipeStepID,
			exampleRecipeStepInstrument.Notes,
			exampleRecipeStepInstrument.BelongsToRecipeStep,
			exampleRecipeStepInstrument.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRecipeStepInstrumentQuery(exampleRecipeStepInstrument)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_step_instruments SET instrument_id = $1, recipe_step_id = $2, notes = $3, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_step = $4 AND id = $5 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStepInstrument.InstrumentID,
				exampleRecipeStepInstrument.RecipeStepID,
				exampleRecipeStepInstrument.Notes,
				exampleRecipeStepInstrument.BelongsToRecipeStep,
				exampleRecipeStepInstrument.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeStepInstrument(ctx, exampleRecipeStepInstrument)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStepInstrument.InstrumentID,
				exampleRecipeStepInstrument.RecipeStepID,
				exampleRecipeStepInstrument.Notes,
				exampleRecipeStepInstrument.BelongsToRecipeStep,
				exampleRecipeStepInstrument.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipeStepInstrument(ctx, exampleRecipeStepInstrument)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeStepInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "UPDATE recipe_step_instruments SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepInstrument.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRecipeStepInstrumentQuery(exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_step_instruments SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepInstrument.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeStepInstrument(ctx, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepInstrument.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRecipeStepInstrument(ctx, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)
		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
		exampleRecipeStepInstrument.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepInstrument.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeStepInstrument(ctx, exampleRecipeStep.ID, exampleRecipeStepInstrument.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
