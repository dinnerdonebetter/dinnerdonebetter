package sqlite

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

func buildMockRowFromRecipeStep(x *models.RecipeStep) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepsTableColumns).AddRow(
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
		x.UpdatedOn,
		x.ArchivedOn,
		x.BelongsTo,
	)

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
		x.UpdatedOn,
		x.BelongsTo,
		x.ID,
	)

	return exampleRows
}

func TestSqlite_buildGetRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		exampleRecipeStepID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE belongs_to = ? AND id = ?"
		actualQuery, args := s.buildGetRecipeStepQuery(exampleRecipeStepID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleRecipeStepID, args[1].(uint64))
	})
}

func TestSqlite_GetRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeStep{
			ID: 123,
		}
		expectedUserID := uint64(321)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromRecipeStep(expected))

		actual, err := s.GetRecipeStep(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeStep{
			ID: 123,
		}
		expectedUserID := uint64(321)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetRecipeStep(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetRecipeStepCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		actualQuery, args := s.buildGetRecipeStepCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestSqlite_GetRecipeStepCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCount := uint64(666)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := s.GetRecipeStepCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetAllRecipeStepsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM recipe_steps WHERE archived_on IS NULL"

		actualQuery := s.buildGetAllRecipeStepsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestSqlite_GetAllRecipeStepsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM recipe_steps WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := s.GetAllRecipeStepsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildGetRecipeStepsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		actualQuery, args := s.buildGetRecipeStepsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestSqlite_GetRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_steps WHERE archived_on IS NULL"
		expectedRecipeStep := &models.RecipeStep{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.RecipeStepList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			RecipeSteps: []models.RecipeStep{
				*expectedRecipeStep,
			},
		}

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStep(expectedRecipeStep))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := s.GetRecipeSteps(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetRecipeSteps(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetRecipeSteps(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStep{
			ID: 321,
		}
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStep(expected))

		actual, err := s.GetRecipeSteps(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStep{
			ID: 321,
		}
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_steps WHERE archived_on IS NULL"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStep(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetRecipeSteps(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_GetAllRecipeStepsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedRecipeStep := &models.RecipeStep{
			ID: 321,
		}
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStep(expectedRecipeStep))

		expected := []models.RecipeStep{*expectedRecipeStep}
		actual, err := s.GetAllRecipeStepsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := s.GetAllRecipeStepsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := s.GetAllRecipeStepsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleRecipeStep := &models.RecipeStep{
			ID: 321,
		}
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStep(exampleRecipeStep))

		actual, err := s.GetAllRecipeStepsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildCreateRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.RecipeStep{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 9
		expectedQuery := "INSERT INTO recipe_steps (index,preparation_id,prerequisite_step,min_estimated_time_in_seconds,max_estimated_time_in_seconds,temperature_in_celsius,notes,recipe_id,belongs_to) VALUES (?,?,?,?,?,?,?,?,?)"
		actualQuery, args := s.buildCreateRecipeStepQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Index, args[0].(uint))
		assert.Equal(t, expected.PreparationID, args[1].(uint64))
		assert.Equal(t, expected.PrerequisiteStep, args[2].(uint64))
		assert.Equal(t, expected.MinEstimatedTimeInSeconds, args[3].(uint32))
		assert.Equal(t, expected.MaxEstimatedTimeInSeconds, args[4].(uint32))
		assert.Equal(t, expected.TemperatureInCelsius, args[5].(*uint16))
		assert.Equal(t, expected.Notes, args[6].(string))
		assert.Equal(t, expected.RecipeID, args[7].(uint64))
		assert.Equal(t, expected.BelongsTo, args[8].(uint64))
	})
}

func TestSqlite_CreateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStep{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeStepCreationInput{
			Index:                     expected.Index,
			PreparationID:             expected.PreparationID,
			PrerequisiteStep:          expected.PrerequisiteStep,
			MinEstimatedTimeInSeconds: expected.MinEstimatedTimeInSeconds,
			MaxEstimatedTimeInSeconds: expected.MaxEstimatedTimeInSeconds,
			TemperatureInCelsius:      expected.TemperatureInCelsius,
			Notes:                     expected.Notes,
			RecipeID:                  expected.RecipeID,
			BelongsTo:                 expected.BelongsTo,
		}

		s, mockDB := buildTestService(t)

		expectedCreationQuery := "INSERT INTO recipe_steps (index,preparation_id,prerequisite_step,min_estimated_time_in_seconds,max_estimated_time_in_seconds,temperature_in_celsius,notes,recipe_id,belongs_to) VALUES (?,?,?,?,?,?,?,?,?)"
		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				expected.Index,
				expected.PreparationID,
				expected.PrerequisiteStep,
				expected.MinEstimatedTimeInSeconds,
				expected.MaxEstimatedTimeInSeconds,
				expected.TemperatureInCelsius,
				expected.Notes,
				expected.RecipeID,
				expected.BelongsTo,
			).WillReturnResult(sqlmock.NewResult(int64(expected.ID), 1))

		expectedTimeQuery := "SELECT created_on FROM recipe_steps WHERE id = ?"
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedTimeQuery)).
			WithArgs(expected.ID).
			WillReturnRows(sqlmock.NewRows([]string{"created_on"}).AddRow(expected.CreatedOn))

		actual, err := s.CreateRecipeStep(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStep{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeStepCreationInput{
			Index:                     expected.Index,
			PreparationID:             expected.PreparationID,
			PrerequisiteStep:          expected.PrerequisiteStep,
			MinEstimatedTimeInSeconds: expected.MinEstimatedTimeInSeconds,
			MaxEstimatedTimeInSeconds: expected.MaxEstimatedTimeInSeconds,
			TemperatureInCelsius:      expected.TemperatureInCelsius,
			Notes:                     expected.Notes,
			RecipeID:                  expected.RecipeID,
			BelongsTo:                 expected.BelongsTo,
		}
		expectedQuery := "INSERT INTO recipe_steps (index,preparation_id,prerequisite_step,min_estimated_time_in_seconds,max_estimated_time_in_seconds,temperature_in_celsius,notes,recipe_id,belongs_to) VALUES (?,?,?,?,?,?,?,?,?)"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Index,
				expected.PreparationID,
				expected.PrerequisiteStep,
				expected.MinEstimatedTimeInSeconds,
				expected.MaxEstimatedTimeInSeconds,
				expected.TemperatureInCelsius,
				expected.Notes,
				expected.RecipeID,
				expected.BelongsTo,
			).WillReturnError(errors.New("blah"))

		actual, err := s.CreateRecipeStep(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildUpdateRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.RecipeStep{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 10
		expectedQuery := "UPDATE recipe_steps SET index = ?, preparation_id = ?, prerequisite_step = ?, min_estimated_time_in_seconds = ?, max_estimated_time_in_seconds = ?, temperature_in_celsius = ?, notes = ?, recipe_id = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"
		actualQuery, args := s.buildUpdateRecipeStepQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Index, args[0].(uint))
		assert.Equal(t, expected.PreparationID, args[1].(uint64))
		assert.Equal(t, expected.PrerequisiteStep, args[2].(uint64))
		assert.Equal(t, expected.MinEstimatedTimeInSeconds, args[3].(uint32))
		assert.Equal(t, expected.MaxEstimatedTimeInSeconds, args[4].(uint32))
		assert.Equal(t, expected.TemperatureInCelsius, args[5].(*uint16))
		assert.Equal(t, expected.Notes, args[6].(string))
		assert.Equal(t, expected.RecipeID, args[7].(uint64))
		assert.Equal(t, expected.BelongsTo, args[8].(uint64))
		assert.Equal(t, expected.ID, args[9].(uint64))
	})
}

func TestSqlite_UpdateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStep{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewResult(int64(expected.ID), 1)
		expectedQuery := "UPDATE recipe_steps SET index = ?, preparation_id = ?, prerequisite_step = ?, min_estimated_time_in_seconds = ?, max_estimated_time_in_seconds = ?, temperature_in_celsius = ?, notes = ?, recipe_id = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Index,
				expected.PreparationID,
				expected.PrerequisiteStep,
				expected.MinEstimatedTimeInSeconds,
				expected.MaxEstimatedTimeInSeconds,
				expected.TemperatureInCelsius,
				expected.Notes,
				expected.RecipeID,
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(exampleRows)

		err := s.UpdateRecipeStep(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStep{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_steps SET index = ?, preparation_id = ?, prerequisite_step = ?, min_estimated_time_in_seconds = ?, max_estimated_time_in_seconds = ?, temperature_in_celsius = ?, notes = ?, recipe_id = ?, updated_on = (strftime('%s','now')) WHERE belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Index,
				expected.PreparationID,
				expected.PrerequisiteStep,
				expected.MinEstimatedTimeInSeconds,
				expected.MaxEstimatedTimeInSeconds,
				expected.TemperatureInCelsius,
				expected.Notes,
				expected.RecipeID,
				expected.BelongsTo,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := s.UpdateRecipeStep(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestSqlite_buildArchiveRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		s, _ := buildTestService(t)
		expected := &models.RecipeStep{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE recipe_steps SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"
		actualQuery, args := s.buildArchiveRecipeStepQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestSqlite_ArchiveRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStep{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_steps SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.ArchiveRecipeStep(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.RecipeStep{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_steps SET updated_on = (strftime('%s','now')), archived_on = (strftime('%s','now')) WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		s, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := s.ArchiveRecipeStep(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
