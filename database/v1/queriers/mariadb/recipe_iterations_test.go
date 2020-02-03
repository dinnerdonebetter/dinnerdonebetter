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

func buildMockRowFromRecipeIteration(x *models.RecipeIteration) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeIterationsTableColumns).AddRow(
		x.ID,
		x.RecipeID,
		x.EndDifficultyRating,
		x.EndComplexityRating,
		x.EndTasteRating,
		x.EndOverallRating,
		x.CreatedOn,
		x.UpdatedOn,
		x.ArchivedOn,
		x.BelongsTo,
	)

	return exampleRows
}

func buildErroneousMockRowFromRecipeIteration(x *models.RecipeIteration) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeIterationsTableColumns).AddRow(
		x.ArchivedOn,
		x.RecipeID,
		x.EndDifficultyRating,
		x.EndComplexityRating,
		x.EndTasteRating,
		x.EndOverallRating,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsTo,
		x.ID,
	)

	return exampleRows
}

func TestMariaDB_buildGetRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleRecipeIterationID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE belongs_to = ? AND id = ?"
		actualQuery, args := m.buildGetRecipeIterationQuery(exampleRecipeIterationID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleRecipeIterationID, args[1].(uint64))
	})
}

func TestMariaDB_GetRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeIteration{
			ID: 123,
		}
		expectedUserID := uint64(321)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromRecipeIteration(expected))

		actual, err := m.GetRecipeIteration(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE belongs_to = ? AND id = ?"
		expected := &models.RecipeIteration{
			ID: 123,
		}
		expectedUserID := uint64(321)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetRecipeIteration(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetRecipeIterationCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		actualQuery, args := m.buildGetRecipeIterationCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestMariaDB_GetRecipeIterationCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCount := uint64(666)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetRecipeIterationCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetAllRecipeIterationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM recipe_iterations WHERE archived_on IS NULL"

		actualQuery := m.buildGetAllRecipeIterationsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestMariaDB_GetAllRecipeIterationsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM recipe_iterations WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := m.GetAllRecipeIterationsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildGetRecipeIterationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		actualQuery, args := m.buildGetRecipeIterationsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestMariaDB_GetRecipeIterations(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_iterations WHERE archived_on IS NULL"
		expectedRecipeIteration := &models.RecipeIteration{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.RecipeIterationList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			RecipeIterations: []models.RecipeIteration{
				*expectedRecipeIteration,
			},
		}

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeIteration(expectedRecipeIteration))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := m.GetRecipeIterations(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetRecipeIterations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetRecipeIterations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe iteration", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeIteration{
			ID: 321,
		}
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeIteration(expected))

		actual, err := m.GetRecipeIterations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeIteration{
			ID: 321,
		}
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ? LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_iterations WHERE archived_on IS NULL"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeIteration(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetRecipeIterations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_GetAllRecipeIterationsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedRecipeIteration := &models.RecipeIteration{
			ID: 321,
		}
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeIteration(expectedRecipeIteration))

		expected := []models.RecipeIteration{*expectedRecipeIteration}
		actual, err := m.GetAllRecipeIterationsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := m.GetAllRecipeIterationsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := m.GetAllRecipeIterationsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleRecipeIteration := &models.RecipeIteration{
			ID: 321,
		}
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeIteration(exampleRecipeIteration))

		actual, err := m.GetAllRecipeIterationsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildCreateRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeIteration{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 6
		expectedQuery := "INSERT INTO recipe_iterations (recipe_id,end_difficulty_rating,end_complexity_rating,end_taste_rating,end_overall_rating,belongs_to,created_on) VALUES (?,?,?,?,?,?,UNIX_TIMESTAMP())"
		actualQuery, args := m.buildCreateRecipeIterationQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.RecipeID, args[0].(uint64))
		assert.Equal(t, expected.EndDifficultyRating, args[1].(float32))
		assert.Equal(t, expected.EndComplexityRating, args[2].(float32))
		assert.Equal(t, expected.EndTasteRating, args[3].(float32))
		assert.Equal(t, expected.EndOverallRating, args[4].(float32))
		assert.Equal(t, expected.BelongsTo, args[5].(uint64))
	})
}

func TestMariaDB_CreateRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeIteration{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeIterationCreationInput{
			RecipeID:            expected.RecipeID,
			EndDifficultyRating: expected.EndDifficultyRating,
			EndComplexityRating: expected.EndComplexityRating,
			EndTasteRating:      expected.EndTasteRating,
			EndOverallRating:    expected.EndOverallRating,
			BelongsTo:           expected.BelongsTo,
		}

		m, mockDB := buildTestService(t)

		expectedCreationQuery := "INSERT INTO recipe_iterations (recipe_id,end_difficulty_rating,end_complexity_rating,end_taste_rating,end_overall_rating,belongs_to,created_on) VALUES (?,?,?,?,?,?,UNIX_TIMESTAMP())"
		mockDB.ExpectExec(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				expected.RecipeID,
				expected.EndDifficultyRating,
				expected.EndComplexityRating,
				expected.EndTasteRating,
				expected.EndOverallRating,
				expected.BelongsTo,
			).WillReturnResult(sqlmock.NewResult(int64(expected.ID), 1))

		expectedTimeQuery := "SELECT created_on FROM recipe_iterations WHERE id = ?"
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedTimeQuery)).
			WithArgs(expected.ID).
			WillReturnRows(sqlmock.NewRows([]string{"created_on"}).AddRow(expected.CreatedOn))

		actual, err := m.CreateRecipeIteration(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeIteration{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.RecipeIterationCreationInput{
			RecipeID:            expected.RecipeID,
			EndDifficultyRating: expected.EndDifficultyRating,
			EndComplexityRating: expected.EndComplexityRating,
			EndTasteRating:      expected.EndTasteRating,
			EndOverallRating:    expected.EndOverallRating,
			BelongsTo:           expected.BelongsTo,
		}
		expectedQuery := "INSERT INTO recipe_iterations (recipe_id,end_difficulty_rating,end_complexity_rating,end_taste_rating,end_overall_rating,belongs_to,created_on) VALUES (?,?,?,?,?,?,UNIX_TIMESTAMP())"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.RecipeID,
				expected.EndDifficultyRating,
				expected.EndComplexityRating,
				expected.EndTasteRating,
				expected.EndOverallRating,
				expected.BelongsTo,
			).WillReturnError(errors.New("blah"))

		actual, err := m.CreateRecipeIteration(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildUpdateRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeIteration{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 7
		expectedQuery := "UPDATE recipe_iterations SET recipe_id = ?, end_difficulty_rating = ?, end_complexity_rating = ?, end_taste_rating = ?, end_overall_rating = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"
		actualQuery, args := m.buildUpdateRecipeIterationQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.RecipeID, args[0].(uint64))
		assert.Equal(t, expected.EndDifficultyRating, args[1].(float32))
		assert.Equal(t, expected.EndComplexityRating, args[2].(float32))
		assert.Equal(t, expected.EndTasteRating, args[3].(float32))
		assert.Equal(t, expected.EndOverallRating, args[4].(float32))
		assert.Equal(t, expected.BelongsTo, args[5].(uint64))
		assert.Equal(t, expected.ID, args[6].(uint64))
	})
}

func TestMariaDB_UpdateRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeIteration{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewResult(int64(expected.ID), 1)
		expectedQuery := "UPDATE recipe_iterations SET recipe_id = ?, end_difficulty_rating = ?, end_complexity_rating = ?, end_taste_rating = ?, end_overall_rating = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.RecipeID,
				expected.EndDifficultyRating,
				expected.EndComplexityRating,
				expected.EndTasteRating,
				expected.EndOverallRating,
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(exampleRows)

		err := m.UpdateRecipeIteration(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeIteration{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_iterations SET recipe_id = ?, end_difficulty_rating = ?, end_complexity_rating = ?, end_taste_rating = ?, end_overall_rating = ?, updated_on = UNIX_TIMESTAMP() WHERE belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.RecipeID,
				expected.EndDifficultyRating,
				expected.EndComplexityRating,
				expected.EndTasteRating,
				expected.EndOverallRating,
				expected.BelongsTo,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := m.UpdateRecipeIteration(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestMariaDB_buildArchiveRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		m, _ := buildTestService(t)
		expected := &models.RecipeIteration{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE recipe_iterations SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"
		actualQuery, args := m.buildArchiveRecipeIterationQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestMariaDB_ArchiveRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeIteration{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_iterations SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := m.ArchiveRecipeIteration(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.RecipeIteration{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_iterations SET updated_on = UNIX_TIMESTAMP(), archived_on = UNIX_TIMESTAMP() WHERE archived_on IS NULL AND belongs_to = ? AND id = ?"

		m, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := m.ArchiveRecipeIteration(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
