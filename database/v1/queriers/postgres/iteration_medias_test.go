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

func buildMockRowsFromIterationMedias(iterationMedias ...*models.IterationMedia) *sqlmock.Rows {
	columns := iterationMediasTableColumns

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range iterationMedias {
		rowValues := []driver.Value{
			x.ID,
			x.Path,
			x.Mimetype,
			x.RecipeIterationID,
			x.RecipeStepID,
			x.CreatedOn,
			x.LastUpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipeIteration,
		}

		exampleRows.AddRow(rowValues...)
	}

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
		x.LastUpdatedOn,
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

		_, err := p.scanIterationMedias(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, err := p.scanIterationMedias(mockRows)
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

		expectedQuery := "SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE iteration_medias.belongs_to_recipe_iteration = $1 AND iteration_medias.id = $2 AND recipe_iterations.belongs_to_recipe = $3 AND recipe_iterations.id = $4 AND recipes.id = $5"
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
	expectedQuery := "SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE iteration_medias.belongs_to_recipe_iteration = $1 AND iteration_medias.id = $2 AND recipe_iterations.belongs_to_recipe = $3 AND recipe_iterations.id = $4 AND recipes.id = $5"

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
			WillReturnRows(buildMockRowsFromIterationMedias(exampleIterationMedia))

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

func TestPostgres_buildGetBatchOfIterationMediasQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		beginID, endID := uint64(1), uint64(1000)

		expectedQuery := "SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias WHERE iteration_medias.id > $1 AND iteration_medias.id < $2"
		expectedArgs := []interface{}{
			beginID,
			endID,
		}
		actualQuery, actualArgs := p.buildGetBatchOfIterationMediasQuery(beginID, endID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetAllIterationMedias(T *testing.T) {
	T.Parallel()

	expectedCountQuery := "SELECT COUNT(iteration_medias.id) FROM iteration_medias WHERE iteration_medias.archived_on IS NULL"
	expectedGetQuery := "SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias WHERE iteration_medias.id > $1 AND iteration_medias.id < $2"

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleIterationMediaList := fakemodels.BuildFakeIterationMediaList()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(
				buildMockRowsFromIterationMedias(
					&exampleIterationMediaList.IterationMedias[0],
					&exampleIterationMediaList.IterationMedias[1],
					&exampleIterationMediaList.IterationMedias[2],
				),
			)

		out := make(chan []models.IterationMedia)
		doneChan := make(chan bool, 1)

		err := p.GetAllIterationMedias(ctx, out)
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

		out := make(chan []models.IterationMedia)

		err := p.GetAllIterationMedias(ctx, out)
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

		out := make(chan []models.IterationMedia)

		err := p.GetAllIterationMedias(ctx, out)
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

		out := make(chan []models.IterationMedia)

		err := p.GetAllIterationMedias(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with invalid response from database", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		p, mockDB := buildTestService(t)
		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
		expectedCount := uint64(20)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCountQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedGetQuery)).
			WithArgs(
				uint64(1),
				uint64(1001),
			).
			WillReturnRows(buildErroneousMockRowFromIterationMedia(exampleIterationMedia))

		out := make(chan []models.IterationMedia)

		err := p.GetAllIterationMedias(ctx, out)
		assert.NoError(t, err)

		time.Sleep(time.Second)

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

		expectedQuery := "SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE iteration_medias.archived_on IS NULL AND iteration_medias.belongs_to_recipe_iteration = $1 AND recipe_iterations.belongs_to_recipe = $2 AND recipe_iterations.id = $3 AND recipes.id = $4 AND iteration_medias.created_on > $5 AND iteration_medias.created_on < $6 AND iteration_medias.last_updated_on > $7 AND iteration_medias.last_updated_on < $8 ORDER BY iteration_medias.id LIMIT 20 OFFSET 180"
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

	expectedQuery := "SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id WHERE iteration_medias.archived_on IS NULL AND iteration_medias.belongs_to_recipe_iteration = $1 AND recipe_iterations.belongs_to_recipe = $2 AND recipe_iterations.id = $3 AND recipes.id = $4 ORDER BY iteration_medias.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleIterationMediaList := fakemodels.BuildFakeIterationMediaList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
				exampleRecipeIteration.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromIterationMedias(
					&exampleIterationMediaList.IterationMedias[0],
					&exampleIterationMediaList.IterationMedias[1],
					&exampleIterationMediaList.IterationMedias[2],
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
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

func TestPostgres_buildGetIterationMediasWithIDsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID
		exampleIDs := []uint64{
			789,
			123,
			456,
		}

		expectedQuery := fmt.Sprintf("SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM (SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE iteration_medias.archived_on IS NULL AND iteration_medias.belongs_to_recipe_iteration = $1 AND recipe_iterations.belongs_to_recipe = $2 AND recipe_iterations.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS iteration_medias WHERE iteration_medias.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		expectedArgs := []interface{}{
			exampleRecipeIteration.ID,
			exampleRecipe.ID,
			exampleRecipeIteration.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetIterationMediasWithIDsQuery(exampleRecipe.ID, exampleRecipeIteration.ID, defaultLimit, exampleIDs)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetIterationMediasWithIDs(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIterationMediaList := fakemodels.BuildFakeIterationMediaList()
		var exampleIDs []uint64
		for _, iterationMedia := range exampleIterationMediaList.IterationMedias {
			exampleIDs = append(exampleIDs, iterationMedia.ID)
		}

		expectedQuery := fmt.Sprintf("SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM (SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE iteration_medias.archived_on IS NULL AND iteration_medias.belongs_to_recipe_iteration = $1 AND recipe_iterations.belongs_to_recipe = $2 AND recipe_iterations.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS iteration_medias WHERE iteration_medias.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(
				buildMockRowsFromIterationMedias(
					&exampleIterationMediaList.IterationMedias[0],
					&exampleIterationMediaList.IterationMedias[1],
					&exampleIterationMediaList.IterationMedias[2],
				),
			)

		actual, err := p.GetIterationMediasWithIDs(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, defaultLimit, exampleIDs)

		assert.NoError(t, err)
		assert.Equal(t, exampleIterationMediaList.IterationMedias, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
		exampleRecipeIteration.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM (SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE iteration_medias.archived_on IS NULL AND iteration_medias.belongs_to_recipe_iteration = $1 AND recipe_iterations.belongs_to_recipe = $2 AND recipe_iterations.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS iteration_medias WHERE iteration_medias.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetIterationMediasWithIDs(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, defaultLimit, exampleIDs)

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

		exampleIDs := []uint64{123, 456, 789}

		expectedQuery := fmt.Sprintf("SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM (SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE iteration_medias.archived_on IS NULL AND iteration_medias.belongs_to_recipe_iteration = $1 AND recipe_iterations.belongs_to_recipe = $2 AND recipe_iterations.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS iteration_medias WHERE iteration_medias.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnError(errors.New("blah"))

		actual, err := p.GetIterationMediasWithIDs(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, defaultLimit, exampleIDs)

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

		exampleIDs := []uint64{123, 456, 789}
		expectedQuery := fmt.Sprintf("SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM (SELECT iteration_medias.id, iteration_medias.path, iteration_medias.mimetype, iteration_medias.recipe_iteration_id, iteration_medias.recipe_step_id, iteration_medias.created_on, iteration_medias.last_updated_on, iteration_medias.archived_on, iteration_medias.belongs_to_recipe_iteration FROM iteration_medias JOIN recipe_iterations ON iteration_medias.belongs_to_recipe_iteration=recipe_iterations.id JOIN recipes ON recipe_iterations.belongs_to_recipe=recipes.id JOIN unnest('{%s}'::int[]) WHERE iteration_medias.archived_on IS NULL AND iteration_medias.belongs_to_recipe_iteration = $1 AND recipe_iterations.belongs_to_recipe = $2 AND recipe_iterations.id = $3 AND recipes.id = $4 WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d) AS iteration_medias WHERE iteration_medias.archived_on IS NULL", joinUint64s(exampleIDs), defaultLimit)

		exampleIterationMedia := fakemodels.BuildFakeIterationMedia()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs().
			WillReturnRows(buildErroneousMockRowFromIterationMedia(exampleIterationMedia))

		actual, err := p.GetIterationMediasWithIDs(ctx, exampleRecipe.ID, exampleRecipeIteration.ID, defaultLimit, exampleIDs)

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

		expectedQuery := "INSERT INTO iteration_medias (path,mimetype,recipe_iteration_id,recipe_step_id,belongs_to_recipe_iteration) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleIterationMedia.Path,
			exampleIterationMedia.Mimetype,
			exampleIterationMedia.RecipeIterationID,
			exampleIterationMedia.RecipeStepID,
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

	expectedCreationQuery := "INSERT INTO iteration_medias (path,mimetype,recipe_iteration_id,recipe_step_id,belongs_to_recipe_iteration) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_on"

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
				exampleIterationMedia.Path,
				exampleIterationMedia.Mimetype,
				exampleIterationMedia.RecipeIterationID,
				exampleIterationMedia.RecipeStepID,
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
				exampleIterationMedia.Path,
				exampleIterationMedia.Mimetype,
				exampleIterationMedia.RecipeIterationID,
				exampleIterationMedia.RecipeStepID,
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

		expectedQuery := "UPDATE iteration_medias SET path = $1, mimetype = $2, recipe_iteration_id = $3, recipe_step_id = $4, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_iteration = $5 AND id = $6 RETURNING last_updated_on"
		expectedArgs := []interface{}{
			exampleIterationMedia.Path,
			exampleIterationMedia.Mimetype,
			exampleIterationMedia.RecipeIterationID,
			exampleIterationMedia.RecipeStepID,
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

	expectedQuery := "UPDATE iteration_medias SET path = $1, mimetype = $2, recipe_iteration_id = $3, recipe_step_id = $4, last_updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_iteration = $5 AND id = $6 RETURNING last_updated_on"

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

		exampleRows := sqlmock.NewRows([]string{"last_updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleIterationMedia.Path,
				exampleIterationMedia.Mimetype,
				exampleIterationMedia.RecipeIterationID,
				exampleIterationMedia.RecipeStepID,
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
				exampleIterationMedia.Path,
				exampleIterationMedia.Mimetype,
				exampleIterationMedia.RecipeIterationID,
				exampleIterationMedia.RecipeStepID,
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

		expectedQuery := "UPDATE iteration_medias SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_iteration = $1 AND id = $2 RETURNING archived_on"
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

	expectedQuery := "UPDATE iteration_medias SET last_updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_iteration = $1 AND id = $2 RETURNING archived_on"

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
