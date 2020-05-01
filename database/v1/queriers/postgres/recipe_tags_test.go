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

func buildMockRowsFromRecipeTag(recipeTags ...*models.RecipeTag) *sqlmock.Rows {
	includeCount := len(recipeTags) > 1
	columns := recipeTagsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}
	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeTags {
		rowValues := []driver.Value{
			x.ID,
			x.Name,
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipe,
		}

		if includeCount {
			rowValues = append(rowValues, len(recipeTags))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipeTag(x *models.RecipeTag) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeTagsTableColumns).AddRow(
		x.ArchivedOn,
		x.Name,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsToRecipe,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRecipeTags(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanRecipeTags(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanRecipeTags(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRecipeTagExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "SELECT EXISTS ( SELECT recipe_tags.id FROM recipe_tags JOIN recipes ON recipe_tags.belongs_to_recipe=recipes.id WHERE recipe_tags.belongs_to_recipe = $1 AND recipe_tags.id = $2 AND recipes.id = $3 )"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeTag.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildRecipeTagExistsQuery(exampleRecipe.ID, exampleRecipeTag.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RecipeTagExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT recipe_tags.id FROM recipe_tags JOIN recipes ON recipe_tags.belongs_to_recipe=recipes.id WHERE recipe_tags.belongs_to_recipe = $1 AND recipe_tags.id = $2 AND recipes.id = $3 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeTag.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RecipeTagExists(ctx, exampleRecipe.ID, exampleRecipeTag.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeTag.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RecipeTagExists(ctx, exampleRecipe.ID, exampleRecipeTag.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeTagQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "SELECT recipe_tags.id, recipe_tags.name, recipe_tags.created_on, recipe_tags.updated_on, recipe_tags.archived_on, recipe_tags.belongs_to_recipe FROM recipe_tags JOIN recipes ON recipe_tags.belongs_to_recipe=recipes.id WHERE recipe_tags.belongs_to_recipe = $1 AND recipe_tags.id = $2 AND recipes.id = $3"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeTag.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeTagQuery(exampleRecipe.ID, exampleRecipeTag.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeTag(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT recipe_tags.id, recipe_tags.name, recipe_tags.created_on, recipe_tags.updated_on, recipe_tags.archived_on, recipe_tags.belongs_to_recipe FROM recipe_tags JOIN recipes ON recipe_tags.belongs_to_recipe=recipes.id WHERE recipe_tags.belongs_to_recipe = $1 AND recipe_tags.id = $2 AND recipes.id = $3"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeTag.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromRecipeTag(exampleRecipeTag))

		actual, err := p.GetRecipeTag(ctx, exampleRecipe.ID, exampleRecipeTag.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeTag, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeTag.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeTag(ctx, exampleRecipe.ID, exampleRecipeTag.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeTagsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(recipe_tags.id) FROM recipe_tags WHERE recipe_tags.archived_on IS NULL"
		actualQuery := p.buildGetAllRecipeTagsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeTagsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_tags.id) FROM recipe_tags WHERE recipe_tags.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeTagsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeTagsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_tags.id, recipe_tags.name, recipe_tags.created_on, recipe_tags.updated_on, recipe_tags.archived_on, recipe_tags.belongs_to_recipe, COUNT(recipe_tags.id) FROM recipe_tags JOIN recipes ON recipe_tags.belongs_to_recipe=recipes.id WHERE recipe_tags.archived_on IS NULL AND recipe_tags.belongs_to_recipe = $1 AND recipes.id = $2 AND recipe_tags.created_on > $3 AND recipe_tags.created_on < $4 AND recipe_tags.updated_on > $5 AND recipe_tags.updated_on < $6 GROUP BY recipe_tags.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRecipeTagsQuery(exampleRecipe.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeTags(T *testing.T) {
	T.Parallel()

	expectedListQuery := "SELECT recipe_tags.id, recipe_tags.name, recipe_tags.created_on, recipe_tags.updated_on, recipe_tags.archived_on, recipe_tags.belongs_to_recipe, COUNT(recipe_tags.id) FROM recipe_tags JOIN recipes ON recipe_tags.belongs_to_recipe=recipes.id WHERE recipe_tags.archived_on IS NULL AND recipe_tags.belongs_to_recipe = $1 AND recipes.id = $2 GROUP BY recipe_tags.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeTagList := fakemodels.BuildFakeRecipeTagList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeTag(
					&exampleRecipeTagList.RecipeTags[0],
					&exampleRecipeTagList.RecipeTags[1],
					&exampleRecipeTagList.RecipeTags[2],
				),
			)

		actual, err := p.GetRecipeTags(ctx, exampleRecipe.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeTagList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeTags(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeTags(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe tag", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeTag(exampleRecipeTag))

		actual, err := p.GetRecipeTags(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeTagQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "INSERT INTO recipe_tags (name,belongs_to_recipe) VALUES ($1,$2) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipeTag.Name,
			exampleRecipeTag.BelongsToRecipe,
		}
		actualQuery, actualArgs := p.buildCreateRecipeTagQuery(exampleRecipeTag)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRecipeTag(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO recipe_tags (name,belongs_to_recipe) VALUES ($1,$2) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRecipeTag.ID, exampleRecipeTag.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeTag.Name,
				exampleRecipeTag.BelongsToRecipe,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeTag(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeTag, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeTag.Name,
				exampleRecipeTag.BelongsToRecipe,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipeTag(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeTagQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "UPDATE recipe_tags SET name = $1, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $2 AND id = $3 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleRecipeTag.Name,
			exampleRecipeTag.BelongsToRecipe,
			exampleRecipeTag.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRecipeTagQuery(exampleRecipeTag)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRecipeTag(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_tags SET name = $1, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $2 AND id = $3 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeTag.Name,
				exampleRecipeTag.BelongsToRecipe,
				exampleRecipeTag.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeTag(ctx, exampleRecipeTag)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeTag.Name,
				exampleRecipeTag.BelongsToRecipe,
				exampleRecipeTag.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipeTag(ctx, exampleRecipeTag)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeTagQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "UPDATE recipe_tags SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeTag.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRecipeTagQuery(exampleRecipe.ID, exampleRecipeTag.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRecipeTag(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_tags SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeTag.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeTag(ctx, exampleRecipe.ID, exampleRecipeTag.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeTag.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRecipeTag(ctx, exampleRecipe.ID, exampleRecipeTag.ID)
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
		exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
		exampleRecipeTag.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeTag.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeTag(ctx, exampleRecipe.ID, exampleRecipeTag.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
