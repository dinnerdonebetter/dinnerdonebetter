package mariadb

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowFromRecipeStepEvent(x *models.RecipeStepEvent) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepEventsTableColumns).AddRow(
		x.ID,
		x.EventType,
		x.Done,
		x.RecipeIterationID,
		x.RecipeStepID,
		x.CreatedOn,
		x.UpdatedOn,
		x.ArchivedOn,
		x.BelongsTo,
	)

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
		x.UpdatedOn,
		x.BelongsTo,
		x.ID,
	)

	return exampleRows
}

func TestMariaDB_buildGetRecipeStepEventQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleRecipeStepEventID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE belongs_to = ? AND id = ?"
		actualQuery, args := m.buildGetRecipeStepEventQuery(exampleRecipeStepEventID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleRecipeStepEventID, args[1].(uint64))
	})
}

func TestMariaDB_GetRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeStepEvent{
			ID: 123,
		}
		expectedUserID := uint64(321)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromRecipeStepEvent(expected))

		actual, err := m.GetRecipeStepEvent(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeStepEvent{
			ID: 123,
		}
		expectedUserID := uint64(321)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetRecipeStepEvent(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetRecipeStepEventCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		actualQuery, args := m.buildGetRecipeStepEventCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestMariaDB_GetRecipeStepEventCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCount := uint64(666)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetRecipeStepEventCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetAllRecipeStepEventsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_events WHERE archived_on IS NULL"

		actualQuery := m.buildGetAllRecipeStepEventsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestMariaDB_GetAllRecipeStepEventsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_events WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetAllRecipeStepEventsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetRecipeStepEventsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		actualQuery, args := m.buildGetRecipeStepEventsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestMariaDB_GetRecipeStepEvents(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_step_events WHERE archived_on IS NULL"
		expectedRecipeStepEvent := &models.RecipeStepEvent{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.RecipeStepEventList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			RecipeStepEvents: []models.RecipeStepEvent{
				*expectedRecipeStepEvent,
			},
		}

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepEvent(expectedRecipeStepEvent))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := m.GetRecipeStepEvents(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetRecipeStepEvents(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetRecipeStepEvents(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step event", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStepEvent{
			ID: 321,
		}
		expectedListQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStepEvent(expected))

		actual, err := m.GetRecipeStepEvents(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStepEvent{
			ID: 321,
		}
		expectedListQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_step_events WHERE archived_on IS NULL"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepEvent(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetRecipeStepEvents(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_GetAllRecipeStepEventsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedRecipeStepEvent := &models.RecipeStepEvent{
			ID: 321,
		}
		expectedListQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepEvent(expectedRecipeStepEvent))

		expected := []models.RecipeStepEvent{*expectedRecipeStepEvent}
		actual, err := m.GetAllRecipeStepEventsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetAllRecipeStepEventsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetAllRecipeStepEventsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleRecipeStepEvent := &models.RecipeStepEvent{
			ID: 321,
		}
		expectedListQuery := "SELECT id, event_type, done, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM recipe_step_events WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStepEvent(exampleRecipeStepEvent))

		actual, err := m.GetAllRecipeStepEventsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildCreateRecipeStepEventQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepEvent{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 5
		expectedQuery := "INSERT INTO recipe_step_events (event_type,done,recipe_iteration_id,recipe_step_id,belongs_to,created_on) VALUES (?,?,?,?,?,UNIX_TIMESTAMP())"
		actualQuery, args := m.buildCreateRecipeStepEventQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.EventType, args[0].(string))
		assert.Equal(t, expected.Done, args[1].(bool))
		assert.Equal(t, expected.RecipeIterationID, args[2].(uint64))
		assert.Equal(t, expected.RecipeStepID, args[3].(uint64))
		assert.Equal(t, expected.BelongsTo, args[4].(uint64))
	})
}

func TestMariaDB_CreateRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepEvent{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeStepEventCreationInput{
			EventType:         expected.EventType,
			Done:              expected.Done,
			RecipeIterationID: expected.RecipeIterationID,
			RecipeStepID:      expected.RecipeStepID,
			BelongsTo:         expected.BelongsTo,
		}

		m, mockDB := buildTestService(t)

		expectedCreationQuery := "INSERT INTO recipe_step_events (event_type,done,recipe_iteration_id,recipe_step_id,belongs_to,created_on) VALUES (?,?,?,?,?,UNIX_TIMESTAMP())"
		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				expected.EventType,
				expected.Done,
				expected.RecipeIterationID,
				expected.RecipeStepID,
				expected.BelongsTo,
			).WillReturnResult(sqlmock.NewResult(int64(expected.ID), 1))

		expectedTimeQuery := "SELECT created_on FROM recipe_step_events WHERE id = ?"
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedTimeQuery)).
			WithArgs(expected.ID).
			WillReturnRows(sqlmock.NewRows([]string{"created_on"}).AddRow(expected.CreatedOn))

		actual, err := m.CreateRecipeStepEvent(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepEvent{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeStepEventCreationInput{
			EventType:         expected.EventType,
			Done:              expected.Done,
			RecipeIterationID: expected.RecipeIterationID,
			RecipeStepID:      expected.RecipeStepID,
			BelongsTo:         expected.BelongsTo,
		}
		expectedQuery := "INSERT INTO recipe_step_events (event_type,done,recipe_iteration_id,recipe_step_id,belongs_to,created_on) VALUES (?,?,?,?,?,UNIX_TIMESTAMP())"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.EventType,
				expected.Done,
				expected.RecipeIterationID,
				expected.RecipeStepID,
				expected.BelongsTo,
			).WillReturnError(errors.New("blah"))

		actual, err := m.CreateRecipeStepEvent(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildUpdateRecipeStepEventQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepEvent{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 6
		expectedQuery := "UPDATE recipe_step_events SET event_type = ?, done = ?, recipe_iteration_id = ?, recipe_step_id = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"
		actualQuery, args := m.buildUpdateRecipeStepEventQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.EventType, args[0].(string))
		assert.Equal(t, expected.Done, args[1].(bool))
		assert.Equal(t, expected.RecipeIterationID, args[2].(uint64))
		assert.Equal(t, expected.RecipeStepID, args[3].(uint64))
		assert.Equal(t, expected.BelongsTo, args[4].(uint64))
		assert.Equal(t, expected.ID, args[5].(uint64))
	})
}

func TestMariaDB_UpdateRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepEvent{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewResult(int64(expected.ID), 1)
		expectedQuery := "UPDATE recipe_step_events SET event_type = ?, done = ?, recipe_iteration_id = ?, recipe_step_id = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.EventType,
				expected.Done,
				expected.RecipeIterationID,
				expected.RecipeStepID,
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(exampleRows)

		err := m.UpdateRecipeStepEvent(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepEvent{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_events SET event_type = ?, done = ?, recipe_iteration_id = ?, recipe_step_id = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.EventType,
				expected.Done,
				expected.RecipeIterationID,
				expected.RecipeStepID,
				expected.BelongsTo,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := m.UpdateRecipeStepEvent(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildArchiveRecipeStepEventQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepEvent{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE recipe_step_events SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"
		actualQuery, args := m.buildArchiveRecipeStepEventQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestMariaDB_ArchiveRecipeStepEvent(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepEvent{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_events SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := m.ArchiveRecipeStepEvent(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.RecipeStepEvent{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_events SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := m.ArchiveRecipeStepEvent(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
