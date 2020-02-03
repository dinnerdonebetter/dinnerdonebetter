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

func TestPostgres_buildGetRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleRecipeIterationID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE belongs_to = $1 AND id = $2"
		actualQuery, args := p.buildGetRecipeIterationQuery(exampleRecipeIterationID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleRecipeIterationID, args[1].(uint64))
	})
}

func TestPostgres_GetRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE belongs_to = $1 AND id = $2"
		expected := &models.RecipeIteration{
			ID: 123,
		}
		expectedUserID := uint64(321)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromRecipeIteration(expected))

		actual, err := p.GetRecipeIteration(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE belongs_to = $1 AND id = $2"
		expected := &models.RecipeIteration{
			ID: 123,
		}
		expectedUserID := uint64(321)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeIteration(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeIterationCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		actualQuery, args := p.buildGetRecipeIterationCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestPostgres_GetRecipeIterationCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetRecipeIterationCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeIterationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM recipe_iterations WHERE archived_on IS NULL"

		actualQuery := p.buildGetAllRecipeIterationsCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeIterationsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM recipe_iterations WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeIterationsCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeIterationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		actualQuery, args := p.buildGetRecipeIterationsQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestPostgres_GetRecipeIterations(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
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

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeIteration(expectedRecipeIteration))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := p.GetRecipeIterations(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeIterations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeIterations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe iteration", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeIteration{
			ID: 321,
		}
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeIteration(expected))

		actual, err := p.GetRecipeIterations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.RecipeIteration{
			ID: 321,
		}
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM recipe_iterations WHERE archived_on IS NULL"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeIteration(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeIterations(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_GetAllRecipeIterationsForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedRecipeIteration := &models.RecipeIteration{
			ID: 321,
		}
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromRecipeIteration(expectedRecipeIteration))

		expected := []models.RecipeIteration{*expectedRecipeIteration}
		actual, err := p.GetAllRecipeIterationsForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetAllRecipeIterationsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetAllRecipeIterationsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleRecipeIteration := &models.RecipeIteration{
			ID: 321,
		}
		expectedListQuery := "SELECT id, recipe_id, end_difficulty_rating, end_complexity_rating, end_taste_rating, end_overall_rating, created_on, updated_on, archived_on, belongs_to FROM recipe_iterations WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromRecipeIteration(exampleRecipeIteration))

		actual, err := p.GetAllRecipeIterationsForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RecipeIteration{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 6
		expectedQuery := "INSERT INTO recipe_iterations (recipe_id,end_difficulty_rating,end_complexity_rating,end_taste_rating,end_overall_rating,belongs_to) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_on"
		actualQuery, args := p.buildCreateRecipeIterationQuery(expected)

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

func TestPostgres_CreateRecipeIteration(T *testing.T) {
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
		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(expected.ID, uint64(time.Now().Unix()))
		expectedQuery := "INSERT INTO recipe_iterations (recipe_id,end_difficulty_rating,end_complexity_rating,end_taste_rating,end_overall_rating,belongs_to) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.RecipeID,
				expected.EndDifficultyRating,
				expected.EndComplexityRating,
				expected.EndTasteRating,
				expected.EndOverallRating,
				expected.BelongsTo,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeIteration(context.Background(), expectedInput)
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
		expectedQuery := "INSERT INTO recipe_iterations (recipe_id,end_difficulty_rating,end_complexity_rating,end_taste_rating,end_overall_rating,belongs_to) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.RecipeID,
				expected.EndDifficultyRating,
				expected.EndComplexityRating,
				expected.EndTasteRating,
				expected.EndOverallRating,
				expected.BelongsTo,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipeIteration(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RecipeIteration{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 7
		expectedQuery := "UPDATE recipe_iterations SET recipe_id = $1, end_difficulty_rating = $2, end_complexity_rating = $3, end_taste_rating = $4, end_overall_rating = $5, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $6 AND id = $7 RETURNING updated_on"
		actualQuery, args := p.buildUpdateRecipeIterationQuery(expected)

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

func TestPostgres_UpdateRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeIteration{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		expectedQuery := "UPDATE recipe_iterations SET recipe_id = $1, end_difficulty_rating = $2, end_complexity_rating = $3, end_taste_rating = $4, end_overall_rating = $5, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $6 AND id = $7 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.RecipeID,
				expected.EndDifficultyRating,
				expected.EndComplexityRating,
				expected.EndTasteRating,
				expected.EndOverallRating,
				expected.BelongsTo,
				expected.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeIteration(context.Background(), expected)
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
		expectedQuery := "UPDATE recipe_iterations SET recipe_id = $1, end_difficulty_rating = $2, end_complexity_rating = $3, end_taste_rating = $4, end_overall_rating = $5, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $6 AND id = $7 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.RecipeID,
				expected.EndDifficultyRating,
				expected.EndComplexityRating,
				expected.EndTasteRating,
				expected.EndOverallRating,
				expected.BelongsTo,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipeIteration(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeIterationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.RecipeIteration{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE recipe_iterations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"
		actualQuery, args := p.buildArchiveRecipeIterationQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestPostgres_ArchiveRecipeIteration(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.RecipeIteration{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE recipe_iterations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeIteration(context.Background(), expected.ID, expectedUserID)
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
		expectedQuery := "UPDATE recipe_iterations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeIteration(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
