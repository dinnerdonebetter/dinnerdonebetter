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

func buildMockRowFromRecipeStepInstrument(x *models.RecipeStepInstrument) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepInstrumentsTableColumns).AddRow(
		x.ID,
		x.InstrumentID,
		x.RecipeStepID,
		x.Notes,
		x.CreatedOn,
		x.UpdatedOn,
		x.ArchivedOn,
		x.BelongsTo,
	)

	return exampleRows
}

func buildErroneousMockRowFromRecipeStepInstrument(x *models.RecipeStepInstrument) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepInstrumentsTableColumns).AddRow(
		x.ArchivedOn,
		x.InstrumentID,
		x.RecipeStepID,
		x.Notes,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsTo,
		x.ID,
	)

	return exampleRows
}

func TestMariaDB_buildGetRecipeStepInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleRecipeStepInstrumentID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE belongs_to = ? AND id = ?"
		actualQuery, args := m.buildGetRecipeStepInstrumentQuery(exampleRecipeStepInstrumentID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleRecipeStepInstrumentID, args[1].(uint64))
	})
}

func TestMariaDB_GetRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}
		expectedUserID := uint64(321)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromRecipeStepInstrument(expected))

		actual, err := m.GetRecipeStepInstrument(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeStepInstrument{
			ID: 123,
		}
		expectedUserID := uint64(321)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetRecipeStepInstrument(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetRecipeStepInstrumentCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		actualQuery, args := m.buildGetRecipeStepInstrumentCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestMariaDB_GetRecipeStepInstrumentCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCount := uint64(666)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetRecipeStepInstrumentCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetAllRecipeStepInstrumentsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_instruments WHERE archived_on IS NULL"

		actualQuery := m.buildGetAllRecipeStepInstrumentsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestMariaDB_GetAllRecipeStepInstrumentsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM recipe_step_instruments WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetAllRecipeStepInstrumentsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetRecipeStepInstrumentsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		actualQuery, args := m.buildGetRecipeStepInstrumentsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestMariaDB_GetRecipeStepInstruments(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_step_instruments WHERE archived_on IS NULL"
		expectedRecipeStepInstrument := &models.RecipeStepInstrument{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.RecipeStepInstrumentList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			RecipeStepInstruments: []models.RecipeStepInstrument{
				*expectedRecipeStepInstrument,
			},
		}

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepInstrument(expectedRecipeStepInstrument))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := m.GetRecipeStepInstruments(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetRecipeStepInstruments(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetRecipeStepInstruments(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step instrument", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStepInstrument{
			ID: 321,
		}
		expectedListQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStepInstrument(expected))

		actual, err := m.GetRecipeStepInstruments(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStepInstrument{
			ID: 321,
		}
		expectedListQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_step_instruments WHERE archived_on IS NULL"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepInstrument(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetRecipeStepInstruments(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_GetAllRecipeStepInstrumentsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedRecipeStepInstrument := &models.RecipeStepInstrument{
			ID: 321,
		}
		expectedListQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStepInstrument(expectedRecipeStepInstrument))

		expected := []models.RecipeStepInstrument{*expectedRecipeStepInstrument}
		actual, err := m.GetAllRecipeStepInstrumentsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetAllRecipeStepInstrumentsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetAllRecipeStepInstrumentsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleRecipeStepInstrument := &models.RecipeStepInstrument{
			ID: 321,
		}
		expectedListQuery := "SELECT id, instrument_id, recipe_step_id, notes, created_on, updated_on, archived_on, belongs_to FROM recipe_step_instruments WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStepInstrument(exampleRecipeStepInstrument))

		actual, err := m.GetAllRecipeStepInstrumentsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildCreateRecipeStepInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepInstrument{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 4
		expectedQuery := "INSERT INTO recipe_step_instruments (instrument_id,recipe_step_id,notes,belongs_to,created_on) VALUES (?,?,?,?,UNIX_TIMESTAMP())"
		actualQuery, args := m.buildCreateRecipeStepInstrumentQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.InstrumentID, args[0].(*uint64))
		assert.Equal(t, expected.RecipeStepID, args[1].(uint64))
		assert.Equal(t, expected.Notes, args[2].(string))
		assert.Equal(t, expected.BelongsTo, args[3].(uint64))
	})
}

func TestMariaDB_CreateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepInstrument{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeStepInstrumentCreationInput{
			InstrumentID: expected.InstrumentID,
			RecipeStepID: expected.RecipeStepID,
			Notes:        expected.Notes,
			BelongsTo:    expected.BelongsTo,
		}

		m, mockDB := buildTestService(t)

		expectedCreationQuery := "INSERT INTO recipe_step_instruments (instrument_id,recipe_step_id,notes,belongs_to,created_on) VALUES (?,?,?,?,UNIX_TIMESTAMP())"
		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				expected.InstrumentID,
				expected.RecipeStepID,
				expected.Notes,
				expected.BelongsTo,
			).WillReturnResult(sqlmock.NewResult(int64(expected.ID), 1))

		expectedTimeQuery := "SELECT created_on FROM recipe_step_instruments WHERE id = ?"
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedTimeQuery)).
			WithArgs(expected.ID).
			WillReturnRows(sqlmock.NewRows([]string{"created_on"}).AddRow(expected.CreatedOn))

		actual, err := m.CreateRecipeStepInstrument(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepInstrument{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeStepInstrumentCreationInput{
			InstrumentID: expected.InstrumentID,
			RecipeStepID: expected.RecipeStepID,
			Notes:        expected.Notes,
			BelongsTo:    expected.BelongsTo,
		}
		expectedQuery := "INSERT INTO recipe_step_instruments (instrument_id,recipe_step_id,notes,belongs_to,created_on) VALUES (?,?,?,?,UNIX_TIMESTAMP())"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.InstrumentID,
				expected.RecipeStepID,
				expected.Notes,
				expected.BelongsTo,
			).WillReturnError(errors.New("blah"))

		actual, err := m.CreateRecipeStepInstrument(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildUpdateRecipeStepInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepInstrument{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 5
		expectedQuery := "UPDATE recipe_step_instruments SET instrument_id = ?, recipe_step_id = ?, notes = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"
		actualQuery, args := m.buildUpdateRecipeStepInstrumentQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.InstrumentID, args[0].(*uint64))
		assert.Equal(t, expected.RecipeStepID, args[1].(uint64))
		assert.Equal(t, expected.Notes, args[2].(string))
		assert.Equal(t, expected.BelongsTo, args[3].(uint64))
		assert.Equal(t, expected.ID, args[4].(uint64))
	})
}

func TestMariaDB_UpdateRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepInstrument{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewResult(int64(expected.ID), 1)
		expectedQuery := "UPDATE recipe_step_instruments SET instrument_id = ?, recipe_step_id = ?, notes = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.InstrumentID,
				expected.RecipeStepID,
				expected.Notes,
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(exampleRows)

		err := m.UpdateRecipeStepInstrument(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepInstrument{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_instruments SET instrument_id = ?, recipe_step_id = ?, notes = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.InstrumentID,
				expected.RecipeStepID,
				expected.Notes,
				expected.BelongsTo,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := m.UpdateRecipeStepInstrument(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildArchiveRecipeStepInstrumentQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeStepInstrument{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE recipe_step_instruments SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"
		actualQuery, args := m.buildArchiveRecipeStepInstrumentQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestMariaDB_ArchiveRecipeStepInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStepInstrument{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_instruments SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := m.ArchiveRecipeStepInstrument(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.RecipeStepInstrument{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_step_instruments SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := m.ArchiveRecipeStepInstrument(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
