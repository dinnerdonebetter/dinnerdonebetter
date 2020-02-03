package postgres

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

func TestPostgres_buildGetRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleRecipeStepID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE belongs_to = $1 AND id = $2"
		actualQuery, args := p.buildGetRecipeStepQuery(exampleRecipeStepID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleRecipeStepID, args[1].(uint64))
	})
}

func TestPostgres_GetRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE belongs_to = $1 AND id = $2"
		expected := &models.RecipeStep{
			ID: 123,
		}
		expectedUserID := uint64(321)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromRecipeStep(expected))

		actual, err := p.GetRecipeStep(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE belongs_to = $1 AND id = $2"
		expected := &models.RecipeStep{
			ID: 123,
		}
		expectedUserID := uint64(321)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStep(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		actualQuery, args := p.buildGetRecipeStepCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestPostgres_GetRecipeStepCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetRecipeStepCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeStepsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM recipe_steps WHERE archived_on IS NULL"

		actualQuery := p.buildGetAllRecipeStepsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeStepsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM recipe_steps WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeStepsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		actualQuery, args := p.buildGetRecipeStepsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestPostgres_GetRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
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

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStep(expectedRecipeStep))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := p.GetRecipeSteps(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeSteps(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeSteps(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStep{
			ID: 321,
		}
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStep(expected))

		actual, err := p.GetRecipeSteps(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeStep{
			ID: 321,
		}
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_steps WHERE archived_on IS NULL"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStep(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeSteps(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_GetAllRecipeStepsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedRecipeStep := &models.RecipeStep{
			ID: 321,
		}
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeStep(expectedRecipeStep))

		expected := []models.RecipeStep{*expectedRecipeStep}
		actual, err := p.GetAllRecipeStepsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetAllRecipeStepsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetAllRecipeStepsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleRecipeStep := &models.RecipeStep{
			ID: 321,
		}
		expectedListQuery := "SELECT id, index, preparation_id, prerequisite_step, min_estimated_time_in_seconds, max_estimated_time_in_seconds, temperature_in_celsius, notes, recipe_id, created_on, updated_on, archived_on, belongs_to FROM recipe_steps WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeStep(exampleRecipeStep))

		actual, err := p.GetAllRecipeStepsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RecipeStep{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 9
		expectedQuery := "INSERT INTO recipe_steps (index,preparation_id,prerequisite_step,min_estimated_time_in_seconds,max_estimated_time_in_seconds,temperature_in_celsius,notes,recipe_id,belongs_to) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_on"
		actualQuery, args := p.buildCreateRecipeStepQuery(expected)

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

func TestPostgres_CreateRecipeStep(T *testing.T) {
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
		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(expected.ID, uint64(time.Now().Unix()))
		expectedQuery := "INSERT INTO recipe_steps (index,preparation_id,prerequisite_step,min_estimated_time_in_seconds,max_estimated_time_in_seconds,temperature_in_celsius,notes,recipe_id,belongs_to) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeStep(context.Background(), expectedInput)
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
		expectedQuery := "INSERT INTO recipe_steps (index,preparation_id,prerequisite_step,min_estimated_time_in_seconds,max_estimated_time_in_seconds,temperature_in_celsius,notes,recipe_id,belongs_to) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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

		actual, err := p.CreateRecipeStep(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RecipeStep{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 10
		expectedQuery := "UPDATE recipe_steps SET index = $1, preparation_id = $2, prerequisite_step = $3, min_estimated_time_in_seconds = $4, max_estimated_time_in_seconds = $5, temperature_in_celsius = $6, notes = $7, recipe_id = $8, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $9 AND id = $10 RETURNING updated_on"
		actualQuery, args := p.buildUpdateRecipeStepQuery(expected)

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

func TestPostgres_UpdateRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStep{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		expectedQuery := "UPDATE recipe_steps SET index = $1, preparation_id = $2, prerequisite_step = $3, min_estimated_time_in_seconds = $4, max_estimated_time_in_seconds = $5, temperature_in_celsius = $6, notes = $7, recipe_id = $8, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $9 AND id = $10 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeStep(context.Background(), expected)
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
		expectedQuery := "UPDATE recipe_steps SET index = $1, preparation_id = $2, prerequisite_step = $3, min_estimated_time_in_seconds = $4, max_estimated_time_in_seconds = $5, temperature_in_celsius = $6, notes = $7, recipe_id = $8, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $9 AND id = $10 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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

		err := p.UpdateRecipeStep(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RecipeStep{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE recipe_steps SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"
		actualQuery, args := p.buildArchiveRecipeStepQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestPostgres_ArchiveRecipeStep(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeStep{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_steps SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeStep(context.Background(), expected.ID, expectedUserID)
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
		expectedQuery := "UPDATE recipe_steps SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeStep(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
