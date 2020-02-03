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

func buildMockRowFromIterationMedia(x *models.IterationMedia) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(iterationMediasTableColumns).AddRow(
		x.ID,
		x.Path,
		x.Mimetype,
		x.RecipeIterationID,
		x.RecipeStepID,
		x.CreatedOn,
		x.UpdatedOn,
		x.ArchivedOn,
		x.BelongsTo,
	)

	return exampleRows
}

func buildErroneousMockRowFromIterationMedia(x *models.IterationMedia) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(iterationMediasTableColumns).AddRow(
		x.ArchivedOn,
		x.Path,
		x.Mimetype,
		x.RecipeIterationID,
		x.RecipeStepID,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsTo,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_buildGetIterationMediaQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleIterationMediaID := uint64(123)
		exampleUserID := uint64(321)

		expectedArgCount := 2
		expectedQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE belongs_to = $1 AND id = $2"
		actualQuery, args := p.buildGetIterationMediaQuery(exampleIterationMediaID, exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
		assert.Equal(t, exampleIterationMediaID, args[1].(uint64))
	})
}

func TestPostgres_GetIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE belongs_to = $1 AND id = $2"
		expected := &models.IterationMedia{
			ID: 123,
		}
		expectedUserID := uint64(321)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnRows(buildMockRowFromIterationMedia(expected))

		actual, err := p.GetIterationMedia(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE belongs_to = $1 AND id = $2"
		expected := &models.IterationMedia{
			ID: 123,
		}
		expectedUserID := uint64(321)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID, expected.ID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetIterationMedia(context.Background(), expected.ID, expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetIterationMediaCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT COUNT(id) FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		actualQuery, args := p.buildGetIterationMediaCountQuery(models.DefaultQueryFilter(), exampleUserID)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestPostgres_GetIterationMediaCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expectedQuery := "SELECT COUNT(id) FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetIterationMediaCount(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllIterationMediasCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expectedQuery := "SELECT COUNT(id) FROM iteration_medias WHERE archived_on IS NULL"

		actualQuery := p.buildGetAllIterationMediasCountQuery()
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllIterationMediasCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedQuery := "SELECT COUNT(id) FROM iteration_medias WHERE archived_on IS NULL"
		expectedCount := uint64(666)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllIterationMediasCount(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetIterationMediasQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		exampleUserID := uint64(321)

		expectedArgCount := 1
		expectedQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		actualQuery, args := p.buildGetIterationMediasQuery(models.DefaultQueryFilter(), exampleUserID)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, exampleUserID, args[0].(uint64))
	})
}

func TestPostgres_GetIterationMedias(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM iteration_medias WHERE archived_on IS NULL"
		expectedIterationMedia := &models.IterationMedia{
			ID: 321,
		}
		expectedCount := uint64(666)
		expected := &models.IterationMediaList{
			Pagination: models.Pagination{
				Page:       1,
				Limit:      20,
				TotalCount: expectedCount,
			},
			IterationMedias: []models.IterationMedia{
				*expectedIterationMedia,
			},
		}

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromIterationMedia(expectedIterationMedia))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actual, err := p.GetIterationMedias(context.Background(), models.DefaultQueryFilter(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetIterationMedias(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetIterationMedias(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning iteration media", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.IterationMedia{
			ID: 321,
		}
		expectedListQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromIterationMedia(expected))

		actual, err := p.GetIterationMedias(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying for count", func(t *testing.T) {
		expectedUserID := uint64(123)
		expected := &models.IterationMedia{
			ID: 321,
		}
		expectedListQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1 LIMIT 20"
		expectedCountQuery := "SELECT COUNT(id) FROM iteration_medias WHERE archived_on IS NULL"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromIterationMedia(expected))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetIterationMedias(context.Background(), models.DefaultQueryFilter(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_GetAllIterationMediasForUser(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedIterationMedia := &models.IterationMedia{
			ID: 321,
		}
		expectedListQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildMockRowFromIterationMedia(expectedIterationMedia))

		expected := []models.IterationMedia{*expectedIterationMedia}
		actual, err := p.GetAllIterationMediasForUser(context.Background(), expectedUserID)

		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetAllIterationMediasForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error querying database", func(t *testing.T) {
		expectedUserID := uint64(123)
		expectedListQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetAllIterationMediasForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with unscannable response", func(t *testing.T) {
		expectedUserID := uint64(123)
		exampleIterationMedia := &models.IterationMedia{
			ID: 321,
		}
		expectedListQuery := "SELECT id, path, mimetype, recipe_iteration_id, recipe_step_id, created_on, updated_on, archived_on, belongs_to FROM iteration_medias WHERE archived_on IS NULL AND belongs_to = $1"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(expectedUserID).
			WillReturnRows(buildErroneousMockRowFromIterationMedia(exampleIterationMedia))

		actual, err := p.GetAllIterationMediasForUser(context.Background(), expectedUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateIterationMediaQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.IterationMedia{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 5
		expectedQuery := "INSERT INTO iteration_medias (path,mimetype,recipe_iteration_id,recipe_step_id,belongs_to) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"
		actualQuery, args := p.buildCreateIterationMediaQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Path, args[0].(string))
		assert.Equal(t, expected.Mimetype, args[1].(string))
		assert.Equal(t, expected.RecipeIterationID, args[2].(uint64))
		assert.Equal(t, expected.RecipeStepID, args[3].(*uint64))
		assert.Equal(t, expected.BelongsTo, args[4].(uint64))
	})
}

func TestPostgres_CreateIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.IterationMedia{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.IterationMediaCreationInput{
			Path:              expected.Path,
			Mimetype:          expected.Mimetype,
			RecipeIterationID: expected.RecipeIterationID,
			RecipeStepID:      expected.RecipeStepID,
			BelongsTo:         expected.BelongsTo,
		}
		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(expected.ID, uint64(time.Now().Unix()))
		expectedQuery := "INSERT INTO iteration_medias (path,mimetype,recipe_iteration_id,recipe_step_id,belongs_to) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Path,
				expected.Mimetype,
				expected.RecipeIterationID,
				expected.RecipeStepID,
				expected.BelongsTo,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateIterationMedia(context.Background(), expectedInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.IterationMedia{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedInput := &models.IterationMediaCreationInput{
			Path:              expected.Path,
			Mimetype:          expected.Mimetype,
			RecipeIterationID: expected.RecipeIterationID,
			RecipeStepID:      expected.RecipeStepID,
			BelongsTo:         expected.BelongsTo,
		}
		expectedQuery := "INSERT INTO iteration_medias (path,mimetype,recipe_iteration_id,recipe_step_id,belongs_to) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Path,
				expected.Mimetype,
				expected.RecipeIterationID,
				expected.RecipeStepID,
				expected.BelongsTo,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateIterationMedia(context.Background(), expectedInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateIterationMediaQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.IterationMedia{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 6
		expectedQuery := "UPDATE iteration_medias SET path = $1, mimetype = $2, recipe_iteration_id = $3, recipe_step_id = $4, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $5 AND id = $6 RETURNING updated_on"
		actualQuery, args := p.buildUpdateIterationMediaQuery(expected)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.Path, args[0].(string))
		assert.Equal(t, expected.Mimetype, args[1].(string))
		assert.Equal(t, expected.RecipeIterationID, args[2].(uint64))
		assert.Equal(t, expected.RecipeStepID, args[3].(*uint64))
		assert.Equal(t, expected.BelongsTo, args[4].(uint64))
		assert.Equal(t, expected.ID, args[5].(uint64))
	})
}

func TestPostgres_UpdateIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.IterationMedia{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		expectedQuery := "UPDATE iteration_medias SET path = $1, mimetype = $2, recipe_iteration_id = $3, recipe_step_id = $4, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $5 AND id = $6 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Path,
				expected.Mimetype,
				expected.RecipeIterationID,
				expected.RecipeStepID,
				expected.BelongsTo,
				expected.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateIterationMedia(context.Background(), expected)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.IterationMedia{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE iteration_medias SET path = $1, mimetype = $2, recipe_iteration_id = $3, recipe_step_id = $4, updated_on = extract(epoch FROM NOW()) WHERE belongs_to = $5 AND id = $6 RETURNING updated_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.Path,
				expected.Mimetype,
				expected.RecipeIterationID,
				expected.RecipeStepID,
				expected.BelongsTo,
				expected.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateIterationMedia(context.Background(), expected)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveIterationMediaQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)
		expected := &models.IterationMedia{
			ID:        321,
			BelongsTo: 123,
		}
		expectedArgCount := 2
		expectedQuery := "UPDATE iteration_medias SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"
		actualQuery, args := p.buildArchiveIterationMediaQuery(expected.ID, expected.BelongsTo)

		assert.Equal(t, expectedQuery, actualQuery)
		assert.Len(t, args, expectedArgCount)
		assert.Equal(t, expected.BelongsTo, args[0].(uint64))
		assert.Equal(t, expected.ID, args[1].(uint64))
	})
}

func TestPostgres_ArchiveIterationMedia(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		expectedUserID := uint64(321)
		expected := &models.IterationMedia{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE iteration_medias SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				expected.BelongsTo,
				expected.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveIterationMedia(context.Background(), expected.ID, expectedUserID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		expectedUserID := uint64(321)
		example := &models.IterationMedia{
			ID:        123,
			BelongsTo: expectedUserID,
			CreatedOn: uint64(time.Now().Unix()),
		}
		expectedQuery := "UPDATE iteration_medias SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to = $1 AND id = $2 RETURNING archived_on"

		p, mockDB := buildTestService(t)
		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				example.BelongsTo,
				example.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveIterationMedia(context.Background(), example.ID, expectedUserID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
