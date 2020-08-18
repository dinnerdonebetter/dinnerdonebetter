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

func buildMockRowsFromRecipeStepEvents(recipeStepEvents ...*models.RecipeStepEvent) *sqlmock.Rows {
	columns := recipeStepEventsTableColumns

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepEvents {
		rowValues := []driver.Value{
			x.ID,
			x.EventType,
			x.Done,
			x.RecipeIterationID,
			x.RecipeStepID,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipeStep,
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipeStepEvent(x *models.RecipeStepEvent) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepEventsTableColumns).AddRow(
		x.ArchivedOn,
		x.EventType,
		x.Done,
		x.RecipeIterationID,
		x.RecipeStepID,
		x.CreatedOn,
		x.LastUpdatedOn,
		x.BelongsToRecipeStep,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRecipeStepEvents(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, err := p.scanRecipeStepEvents(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanRecipeStepEvents(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRecipeStepEventExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "SELECT EXISTS ( SELECT recipe_step_events.id FROM recipe_step_events JOIN recipe_steps ON recipe_step_events.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_events.belongs_to_recipe_step = $1 AND recipe_step_events.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5 )"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepEvent.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildRecipeStepEventExistsQuery(exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RecipeStepEventExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT recipe_step_events.id FROM recipe_step_events JOIN recipe_steps ON recipe_step_events.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_events.belongs_to_recipe_step = $1 AND recipe_step_events.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepEvent.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RecipeStepEventExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)
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
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepEvent.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RecipeStepEventExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepEventQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM recipe_step_events JOIN recipe_steps ON recipe_step_events.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_events.belongs_to_recipe_step = $1 AND recipe_step_events.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepEvent.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepEventQuery(exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepEvent(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM recipe_step_events JOIN recipe_steps ON recipe_step_events.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_events.belongs_to_recipe_step = $1 AND recipe_step_events.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepEvent.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromRecipeStepEvents(exampleRecipeStepEvent))

		actual, err := p.GetRecipeStepEvent(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepEvent, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepEvent.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepEvent(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeStepEventsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(recipe_step_events.id) FROM recipe_step_events WHERE recipe_step_events.archived_on IS NULL"
		actualQuery := p.buildGetAllRecipeStepEventsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeStepEventsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_step_events.id) FROM recipe_step_events WHERE recipe_step_events.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeStepEventsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetBatchOfRecipeStepEventsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM recipe_step_events WHERE recipe_step_events.id > $1 AND recipe_step_events.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfRecipeStepEventsQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllRecipeStepEvents(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(recipe_step_events.id) FROM recipe_step_events WHERE recipe_step_events.archived_on IS NULL"
	expectedGetQuery := "SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM recipe_step_events WHERE recipe_step_events.id > $1 AND recipe_step_events.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeStepEventList := fakemodels.BuildFakeRecipeStepEventList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepEvents(
					&exampleRecipeStepEventList.RecipeStepEvents[0],
					&exampleRecipeStepEventList.RecipeStepEvents[1],
					&exampleRecipeStepEventList.RecipeStepEvents[2],
				),
			)

		out := make(chan []models.RecipeStepEvent)
		doneChan := make(chan bool, 1)

		err := p.GetAllRecipeStepEvents(ctx, out)
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

		out := make(chan []models.RecipeStepEvent)

		err := p.GetAllRecipeStepEvents(ctx, out)
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

		out := make(chan []models.RecipeStepEvent)

		err := p.GetAllRecipeStepEvents(ctx, out)
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

		out := make(chan []models.RecipeStepEvent)

		err := p.GetAllRecipeStepEvents(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepEvent(exampleRecipeStepEvent))

		out := make(chan []models.RecipeStepEvent)

		err := p.GetAllRecipeStepEvents(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepEventsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM recipe_step_events JOIN recipe_steps ON recipe_step_events.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_events.archived_on IS NULL AND recipe_step_events.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 AND recipe_step_events.created_on > $5 AND recipe_step_events.created_on < $6 AND recipe_step_events.last_updated_on > $7 AND recipe_step_events.last_updated_on < $8 ORDER BY recipe_step_events.id LIMIT 20 OFFSET 180"
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
		actualQuery, actualArgs := p.buildGetRecipeStepEventsQuery(exampleRecipe.ID, exampleRecipeStep.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepEvents(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM recipe_step_events JOIN recipe_steps ON recipe_step_events.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_events.archived_on IS NULL AND recipe_step_events.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 ORDER BY recipe_step_events.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepEventList := fakemodels.BuildFakeRecipeStepEventList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepEvents(
					&exampleRecipeStepEventList.RecipeStepEvents[0],
					&exampleRecipeStepEventList.RecipeStepEvents[1],
					&exampleRecipeStepEventList.RecipeStepEvents[2],
				),
			)

		actual, err := p.GetRecipeStepEvents(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepEventList, actual)

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

		actual, err := p.GetRecipeStepEvents(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
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

		actual, err := p.GetRecipeStepEvents(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step event", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepEvent(exampleRecipeStepEvent))

		actual, err := p.GetRecipeStepEvents(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepEventsWithIDsQuery(T *testing.T) {
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

		expectedQuery := fmt.Sprintf("SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM (SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM recipe_step_events JOIN unnest('{%s}'::int[]) WHERE recipe_step_events.archived_on IS NULL AND recipe_step_events.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_events WHERE recipe_step_events.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepEventsWithIDsQuery(exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepEventsWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleRecipeStepEventList := fakemodels.BuildFakeRecipeStepEventList()
		var exampleIDs []uint64
		for _, recipeStepEvent := range exampleRecipeStepEventList.RecipeStepEvents {
			exampleIDs = append(exampleIDs, recipeStepEvent.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM (SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM recipe_step_events JOIN unnest('{%s}'::int[]) WHERE recipe_step_events.archived_on IS NULL AND recipe_step_events.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_events WHERE recipe_step_events.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(
				buildMockRowsFromRecipeStepEvents(
					&exampleRecipeStepEventList.RecipeStepEvents[0],
					&exampleRecipeStepEventList.RecipeStepEvents[1],
					&exampleRecipeStepEventList.RecipeStepEvents[2],
				),
			)

		actual, err := p.GetRecipeStepEventsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepEventList.RecipeStepEvents, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM (SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM recipe_step_events JOIN unnest('{%s}'::int[]) WHERE recipe_step_events.archived_on IS NULL AND recipe_step_events.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_events WHERE recipe_step_events.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepEventsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

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
		expectedQuery := fmt.Sprintf("SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM (SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM recipe_step_events JOIN unnest('{%s}'::int[]) WHERE recipe_step_events.archived_on IS NULL AND recipe_step_events.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_events WHERE recipe_step_events.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeStepEventsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step event", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM (SELECT recipe_step_events.id, recipe_step_events.event_type, recipe_step_events.done, recipe_step_events.recipe_iteration_id, recipe_step_events.recipe_step_id, recipe_step_events.created_on, recipe_step_events.last_updated_on, recipe_step_events.archived_on, recipe_step_events.belongs_to_recipe_step FROM recipe_step_events JOIN unnest('{%s}'::int[]) WHERE recipe_step_events.archived_on IS NULL AND recipe_step_events.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS recipe_step_events WHERE recipe_step_events.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(buildErroneousMockRowFromRecipeStepEvent(exampleRecipeStepEvent))

		actual, err := p.GetRecipeStepEventsWithIDs(ctx, exampleRecipe.ID, exampleRecipeStep.ID, defaultLimit, exampleIDs)

		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeStepEventQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "INSERT INTO recipe_step_events (event_type,done,recipe_iteration_id,recipe_step_id,belongs_to_recipe_step) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipeStepEvent.EventType,
			exampleRecipeStepEvent.Done,
			exampleRecipeStepEvent.RecipeIterationID,
			exampleRecipeStepEvent.RecipeStepID,
			exampleRecipeStepEvent.BelongsToRecipeStep,
		}
		actualQuery, actualArgs := p.buildCreateRecipeStepEventQuery(exampleRecipeStepEvent)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRecipeStepEvent(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO recipe_step_events (event_type,done,recipe_iteration_id,recipe_step_id,belongs_to_recipe_step) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRecipeStepEvent.ID, exampleRecipeStepEvent.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStepEvent.EventType,
				exampleRecipeStepEvent.Done,
				exampleRecipeStepEvent.RecipeIterationID,
				exampleRecipeStepEvent.RecipeStepID,
				exampleRecipeStepEvent.BelongsToRecipeStep,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeStepEvent(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepEvent, actual)

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
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStepEvent.EventType,
				exampleRecipeStepEvent.Done,
				exampleRecipeStepEvent.RecipeIterationID,
				exampleRecipeStepEvent.RecipeStepID,
				exampleRecipeStepEvent.BelongsToRecipeStep,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipeStepEvent(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeStepEventQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "UPDATE recipe_step_events SET event_type = $1, done = $2, recipe_iteration_id = $3, recipe_step_id = $4, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_step = $5 AND id = $6 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleRecipeStepEvent.EventType,
			exampleRecipeStepEvent.Done,
			exampleRecipeStepEvent.RecipeIterationID,
			exampleRecipeStepEvent.RecipeStepID,
			exampleRecipeStepEvent.BelongsToRecipeStep,
			exampleRecipeStepEvent.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRecipeStepEventQuery(exampleRecipeStepEvent)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRecipeStepEvent(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_step_events SET event_type = $1, done = $2, recipe_iteration_id = $3, recipe_step_id = $4, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_step = $5 AND id = $6 RETURNING last_updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStepEvent.EventType,
				exampleRecipeStepEvent.Done,
				exampleRecipeStepEvent.RecipeIterationID,
				exampleRecipeStepEvent.RecipeStepID,
				exampleRecipeStepEvent.BelongsToRecipeStep,
				exampleRecipeStepEvent.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeStepEvent(ctx, exampleRecipeStepEvent)
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
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStepEvent.EventType,
				exampleRecipeStepEvent.Done,
				exampleRecipeStepEvent.RecipeIterationID,
				exampleRecipeStepEvent.RecipeStepID,
				exampleRecipeStepEvent.BelongsToRecipeStep,
				exampleRecipeStepEvent.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipeStepEvent(ctx, exampleRecipeStepEvent)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeStepEventQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "UPDATE recipe_step_events SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepEvent.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRecipeStepEventQuery(exampleRecipeStep.ID, exampleRecipeStepEvent.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRecipeStepEvent(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_step_events SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepEvent.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeStepEvent(ctx, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)
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
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepEvent.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRecipeStepEvent(ctx, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)
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
		exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
		exampleRecipeStepEvent.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepEvent.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeStepEvent(ctx, exampleRecipeStep.ID, exampleRecipeStepEvent.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
