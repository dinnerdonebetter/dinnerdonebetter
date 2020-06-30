package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"testing"
	"time"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func buildMockRowsFromIterationMedia(iterationMedias ...*models.IterationMedia) *sqlmock.Rows {
	includeCount := len(iterationMedias) > 1
	columns := iterationMediasTableColumns

	if includeCount {
		columns = append(columns, "count")
	}
	exampleRows := sqlmock.NewRows(columns)

	for _, x := range iterationMedias {
		rowValues := []driver.Value{
			x.ID,
			x.Source,
			x.Mimetype,
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipeIteration,
		}

		if includeCount {
			rowValues = append(rowValues, len(iterationMedias))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromIterationMedia(x *models.IterationMedia) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(iterationMediasTableColumns).AddRow(
		x.ArchivedOn,
		x.Source,
		x.Mimetype,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsToRecipeIteration,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanIterationMedias(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanIterationMedias(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanIterationMedias(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildIterationMediaExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		expectedQuery := "SELECT EXISTS ( SELECT iteration_medias.id FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE iteration_medias.belongs_to_recipe_iteration = $1 AND iteration_medias.id = $2 AND recipe_iterations.belongs_to_recipe = $3 AND recipe_iterations.id = $4 AND recipes.id = $5 )"
		expectedArgs := []interface{}{
			exampleRecipeIteration.ID,
			exampleIterationMedia.ID,
			exampleRecipe.ID,
			exampleRecipeIteration.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildIterationMediaExistsQuery(exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_IterationMediaExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT iteration_medias.id FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE iteration_medias.belongs_to_recipe_iteration = $1 AND iteration_medias.id = $2 AND recipe_iterations.belongs_to_recipe = $3 AND recipe_iterations.id = $4 AND recipes.id = $5 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleIterationMedia.ID,
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.IterationMediaExists(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleIterationMedia.ID,
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.IterationMediaExists(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetIterationMediaQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		expectedQuery := "SELECT iteration_medias.id, iteration_medias.source, iteration_medias.mimetype, iteration_medias.created_on, iteration_medias.updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE iteration_medias.belongs_to_recipe_iteration = $1 AND iteration_medias.id = $2 AND recipe_iterations.belongs_to_recipe = $3 AND recipe_iterations.id = $4 AND recipes.id = $5"
		expectedArgs := []interface{}{
			exampleRecipeIteration.ID,
			exampleIterationMedia.ID,
			exampleRecipe.ID,
			exampleRecipeIteration.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetIterationMediaQuery(exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetIterationMedia(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT iteration_medias.id, iteration_medias.source, iteration_medias.mimetype, iteration_medias.created_on, iteration_medias.updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE iteration_medias.belongs_to_recipe_iteration = $1 AND iteration_medias.id = $2 AND recipe_iterations.belongs_to_recipe = $3 AND recipe_iterations.id = $4 AND recipes.id = $5"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleIterationMedia.ID,
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromIterationMedia(exampleIterationMedia))

		actual, err := p.GetIterationMedia(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleIterationMedia, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleIterationMedia.ID,
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetIterationMedia(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, exampleIterationMedia.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllIterationMediasCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(iteration_medias.id) FROM iteration_medias WHERE iteration_medias.archived_on IS NULL"
		actualQuery := p.buildGetAllIterationMediasCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllIterationMediasCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(iteration_medias.id) FROM iteration_medias WHERE iteration_medias.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllIterationMediasCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetIterationMediasQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT iteration_medias.id, iteration_medias.source, iteration_medias.mimetype, iteration_medias.created_on, iteration_medias.updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration, (SELECT COUNT(iteration_medias.id) FROM iteration_medias WHERE iteration_medias.archived_on IS NULL) FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE iteration_medias.archived_on IS NULL AND iteration_medias.belongs_to_recipe_iteration = $1 AND recipe_iterations.belongs_to_recipe = $2 AND recipe_iterations.id = $3 AND recipes.id = $4 AND iteration_medias.created_on > $5 AND iteration_medias.created_on < $6 AND iteration_medias.updated_on > $7 AND iteration_medias.updated_on < $8 ORDER BY iteration_medias.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipeIteration.ID,
			exampleRecipe.ID,
			exampleRecipeIteration.ID,
			exampleRecipe.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetIterationMediasQuery(exampleRecipe.ID, exampleRecipeIteration.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetIterationMedias(T *testing.T) {
	T.Parallel()

	expectedListQuery := "SELECT iteration_medias.id, iteration_medias.source, iteration_medias.mimetype, iteration_medias.created_on, iteration_medias.updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration, (SELECT COUNT(iteration_medias.id) FROM iteration_medias WHERE iteration_medias.archived_on IS NULL) FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE iteration_medias.archived_on IS NULL AND iteration_medias.belongs_to_recipe_iteration = $1 AND recipe_iterations.belongs_to_recipe = $2 AND recipe_iterations.id = $3 AND recipes.id = $4 ORDER BY iteration_medias.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleIterationMediaList := fakemodels.BuildFakeIterationMediaList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromIterationMedia(
					&exampleIterationMediaList.IterationMedia[0],
					&exampleIterationMediaList.IterationMedia[1],
					&exampleIterationMediaList.IterationMedia[2],
				),
			)

		actual, err := p.GetIterationMedias(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleIterationMediaList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetIterationMedias(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetIterationMedias(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning iteration media", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromIterationMedia(exampleIterationMedia))

		actual, err := p.GetIterationMedias(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateIterationMediaQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		expectedQuery := "INSERT INTO iteration_medias (source,mimetype,belongs_to_recipe_iteration) VALUES ($1,$2,$3) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleIterationMedia.Source,
			exampleIterationMedia.Mimetype,
			exampleIterationMedia.BelongsToRecipeIteration,
		}
		actualQuery, actualArgs := p.buildCreateIterationMediaQuery(exampleIterationMedia)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateIterationMedia(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO iteration_medias (source,mimetype,belongs_to_recipe_iteration) VALUES ($1,$2,$3) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleIterationMedia.ID, exampleIterationMedia.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleIterationMedia.Source,
				exampleIterationMedia.Mimetype,
				exampleIterationMedia.BelongsToRecipeIteration,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateIterationMedia(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleIterationMedia, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID
		exampleInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleIterationMedia.Source,
				exampleIterationMedia.Mimetype,
				exampleIterationMedia.BelongsToRecipeIteration,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateIterationMedia(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateIterationMediaQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		expectedQuery := "UPDATE iteration_medias SET source = $1, mimetype = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_iteration = $3 AND id = $4 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleIterationMedia.Source,
			exampleIterationMedia.Mimetype,
			exampleIterationMedia.BelongsToRecipeIteration,
			exampleIterationMedia.ID,
		}
		actualQuery, actualArgs := p.buildUpdateIterationMediaQuery(exampleIterationMedia)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateIterationMedia(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE iteration_medias SET source = $1, mimetype = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_iteration = $3 AND id = $4 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleIterationMedia.Source,
				exampleIterationMedia.Mimetype,
				exampleIterationMedia.BelongsToRecipeIteration,
				exampleIterationMedia.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateIterationMedia(ctx, exampleIterationMedia)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleIterationMedia.Source,
				exampleIterationMedia.Mimetype,
				exampleIterationMedia.BelongsToRecipeIteration,
				exampleIterationMedia.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateIterationMedia(ctx, exampleIterationMedia)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveIterationMediaQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		expectedQuery := "UPDATE iteration_medias SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_iteration = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRecipeIteration.ID,
			exampleIterationMedia.ID,
		}
		actualQuery, actualArgs := p.buildArchiveIterationMediaQuery(exampleRecipeIteration.ID, exampleIterationMedia.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveIterationMedia(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE iteration_medias SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_iteration = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleIterationMedia.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveIterationMedia(ctx, exampleRecipeIteration.ID, exampleIterationMedia.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleIterationMedia.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveIterationMedia(ctx, exampleRecipeIteration.ID, exampleIterationMedia.ID)
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
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		exampleIterationMedia.BelongsToRecipeIteration = exampleRecipeIteration.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleIterationMedia.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveIterationMedia(ctx, exampleRecipeIteration.ID, exampleIterationMedia.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
