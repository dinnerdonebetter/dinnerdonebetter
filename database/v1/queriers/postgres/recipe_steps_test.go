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

func buildMockRowsFromRecipeSteps(recipeSteps ...*models.RecipeStep) *sqlmock.Rows {
	columns := recipeStepsTableColumns

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeSteps {
		rowValues := []driver.Value{
			x.ID,
			x.Index,
			x.PreparationID,
			x.PrerequisiteStep,
			x.MinEstimatedTimeInSeconds,
			x.MaxEstimatedTimeInSeconds,
			x.TemperatureInCelsius,
			x.Notes,
			x.RecipeID,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipe,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipeStep(x *models.RecipeStep) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepsTableColumns).AddRow(
		x.ArchivedOn,
		x.Index,
		x.PreparationID,
		x.PrerequisiteStep,
		x.MinEstimatedTimeInSeconds,
		x.MaxEstimatedTimeInSeconds,
		x.TemperatureInCelsius,
		x.Notes,
		x.RecipeID,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.BelongsToRecipe,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := p.scanRecipeSteps(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanRecipeSteps(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRecipeStepExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "SELECT EXISTS ( SELECT recipe_steps.id FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.id = $3 )"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildRecipeStepExistsQuery(exampleRecipe.ID, exampleRecipeStep.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RecipeStepExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT recipe_steps.id FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.id = $3 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RecipeStepExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
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

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RecipeStepExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.id = $3"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepQuery(exampleRecipe.ID, exampleRecipeStep.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStep(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.id = $3"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromRecipeSteps(exampleRecipeStep))

		actual, err := p.GetRecipeStep(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStep, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStep(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeStepsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(recipe_steps.id) FROM recipe_steps WHERE recipe_steps.archived_on IS NULL"
		actualQuery := p.buildGetAllRecipeStepsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeStepsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_steps.id) FROM recipe_steps WHERE recipe_steps.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeStepsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfRecipeStepsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps WHERE recipe_steps.id > $1 AND recipe_steps.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfRecipeStepsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllRecipeSteps(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(recipe_steps.id) FROM recipe_steps WHERE recipe_steps.archived_on IS NULL"
	expectedGetQuery := "SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps WHERE recipe_steps.id > $1 AND recipe_steps.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeStepList := fakemodels.BuildFakeRecipeStepList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromRecipeSteps(
					&exampleRecipeStepList.RecipeSteps[0],
					&exampleRecipeStepList.RecipeSteps[1],
					&exampleRecipeStepList.RecipeSteps[2],
				),
			)

		out := make(chan []models.RecipeStep)
		doneChan := make(chan bool, 1)

		err := p.GetAllRecipeSteps(ctx, out)
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

		out := make(chan []models.RecipeStep)

		err := p.GetAllRecipeSteps(ctx, out)
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

		out := make(chan []models.RecipeStep)

		err := p.GetAllRecipeSteps(ctx, out)
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

		out := make(chan []models.RecipeStep)

		err := p.GetAllRecipeSteps(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStep(exampleRecipeStep))

		out := make(chan []models.RecipeStep)

		err := p.GetAllRecipeSteps(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipes.id = $2 AND recipe_steps.created_on > $3 AND recipe_steps.created_on < $4 AND recipe_steps.last_updated_on > $5 AND recipe_steps.last_updated_on < $6 ORDER BY recipe_steps.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepsQuery(exampleRecipe.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeSteps(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipes.id = $2 ORDER BY recipe_steps.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepList := fakemodels.BuildFakeRecipeStepList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeSteps(
					&exampleRecipeStepList.RecipeSteps[0],
					&exampleRecipeStepList.RecipeSteps[1],
					&exampleRecipeStepList.RecipeSteps[2],
				),
			)

		actual, err := p.GetRecipeSteps(ctx, exampleRecipe.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeSteps(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeSteps(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStep(exampleRecipeStep))

		actual, err := p.GetRecipeSteps(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepsWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM (SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipes.id = $2 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_steps WHERE recipe_steps.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepsWithIDsQuery(exampleRecipe.ID, defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)

		exampleRecipeStepList := fakemodels.BuildFakeRecipeStepList()
		var exampleIDs []uint64
		for _, recipeStep := range exampleRecipeStepList.RecipeSteps {
			exampleIDs = append(exampleIDs, recipeStep.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM (SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipes.id = $2 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_steps WHERE recipe_steps.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeSteps(
					&exampleRecipeStepList.RecipeSteps[0],
					&exampleRecipeStepList.RecipeSteps[1],
					&exampleRecipeStepList.RecipeSteps[2],
				),
			)

		actual, err := p.GetRecipeStepsWithIDs(ctx, exampleRecipe.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepList.RecipeSteps, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM (SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipes.id = $2 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_steps WHERE recipe_steps.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepsWithIDs(ctx, exampleRecipe.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM (SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipes.id = $2 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_steps WHERE recipe_steps.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeStepsWithIDs(ctx, exampleRecipe.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM (SELECT recipe_steps.id, recipe_steps.index, recipe_steps.preparation_id, recipe_steps.prerequisite_step, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.temperature_in_celsius, recipe_steps.notes, recipe_steps.recipe_id, recipe_steps.created_on, recipe_steps.last_updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipes.id = $2 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_steps WHERE recipe_steps.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStep(exampleRecipeStep))

		actual, err := p.GetRecipeStepsWithIDs(ctx, exampleRecipe.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "INSERT INTO recipe_steps (index,preparation_id,prerequisite_step,min_estimated_time_in_seconds,max_estimated_time_in_seconds,temperature_in_celsius,notes,recipe_id,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipeStep.Index,
			exampleRecipeStep.PreparationID,
			exampleRecipeStep.PrerequisiteStep,
			exampleRecipeStep.MinEstimatedTimeInSeconds,
			exampleRecipeStep.MaxEstimatedTimeInSeconds,
			exampleRecipeStep.TemperatureInCelsius,
			exampleRecipeStep.Notes,
			exampleRecipeStep.RecipeID,
			exampleRecipeStep.BelongsToRecipe,
		}
		actualQuery, actualArgs := p.buildCreateRecipeStepQuery(exampleRecipeStep)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRecipeStep(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO recipe_steps (index,preparation_id,prerequisite_step,min_estimated_time_in_seconds,max_estimated_time_in_seconds,temperature_in_celsius,notes,recipe_id,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRecipeStep.ID, exampleRecipeStep.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStep.Index,
				exampleRecipeStep.PreparationID,
				exampleRecipeStep.PrerequisiteStep,
				exampleRecipeStep.MinEstimatedTimeInSeconds,
				exampleRecipeStep.MaxEstimatedTimeInSeconds,
				exampleRecipeStep.TemperatureInCelsius,
				exampleRecipeStep.Notes,
				exampleRecipeStep.RecipeID,
				exampleRecipeStep.BelongsToRecipe,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeStep(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStep, actual)

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
		exampleInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStep.Index,
				exampleRecipeStep.PreparationID,
				exampleRecipeStep.PrerequisiteStep,
				exampleRecipeStep.MinEstimatedTimeInSeconds,
				exampleRecipeStep.MaxEstimatedTimeInSeconds,
				exampleRecipeStep.TemperatureInCelsius,
				exampleRecipeStep.Notes,
				exampleRecipeStep.RecipeID,
				exampleRecipeStep.BelongsToRecipe,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipeStep(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "UPDATE recipe_steps SET index = $1, preparation_id = $2, prerequisite_step = $3, min_estimated_time_in_seconds = $4, max_estimated_time_in_seconds = $5, temperature_in_celsius = $6, notes = $7, recipe_id = $8, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $9 AND id = $10 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleRecipeStep.Index,
			exampleRecipeStep.PreparationID,
			exampleRecipeStep.PrerequisiteStep,
			exampleRecipeStep.MinEstimatedTimeInSeconds,
			exampleRecipeStep.MaxEstimatedTimeInSeconds,
			exampleRecipeStep.TemperatureInCelsius,
			exampleRecipeStep.Notes,
			exampleRecipeStep.RecipeID,
			exampleRecipeStep.BelongsToRecipe,
			exampleRecipeStep.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRecipeStepQuery(exampleRecipeStep)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRecipeStep(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_steps SET index = $1, preparation_id = $2, prerequisite_step = $3, min_estimated_time_in_seconds = $4, max_estimated_time_in_seconds = $5, temperature_in_celsius = $6, notes = $7, recipe_id = $8, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $9 AND id = $10 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.Index,
				exampleRecipeStep.PreparationID,
				exampleRecipeStep.PrerequisiteStep,
				exampleRecipeStep.MinEstimatedTimeInSeconds,
				exampleRecipeStep.MaxEstimatedTimeInSeconds,
				exampleRecipeStep.TemperatureInCelsius,
				exampleRecipeStep.Notes,
				exampleRecipeStep.RecipeID,
				exampleRecipeStep.BelongsToRecipe,
				exampleRecipeStep.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeStep(ctx, exampleRecipeStep)
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.Index,
				exampleRecipeStep.PreparationID,
				exampleRecipeStep.PrerequisiteStep,
				exampleRecipeStep.MinEstimatedTimeInSeconds,
				exampleRecipeStep.MaxEstimatedTimeInSeconds,
				exampleRecipeStep.TemperatureInCelsius,
				exampleRecipeStep.Notes,
				exampleRecipeStep.RecipeID,
				exampleRecipeStep.BelongsToRecipe,
				exampleRecipeStep.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipeStep(ctx, exampleRecipeStep)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "UPDATE recipe_steps SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeStep.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRecipeStepQuery(exampleRecipe.ID, exampleRecipeStep.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRecipeStep(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_steps SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeStep(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
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

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRecipeStep(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
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

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeStep(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
