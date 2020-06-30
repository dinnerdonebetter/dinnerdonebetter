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

func buildMockRowsFromRecipeStepPreparation(recipeStepPreparations ...*models.RecipeStepPreparation) *sqlmock.Rows {
	includeCount := len(recipeStepPreparations) > 1
	columns := recipeStepPreparationsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}
	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeStepPreparations {
		rowValues := []driver.Value{
			x.ID,
			x.ValidPreparationID,
			x.Notes,
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipeStep,
		}

		if includeCount {
			rowValues = append(rowValues, len(recipeStepPreparations))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipeStepPreparation(x *models.RecipeStepPreparation) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepPreparationsTableColumns).AddRow(
		x.ArchivedOn,
		x.ValidPreparationID,
		x.Notes,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsToRecipeStep,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRecipeStepPreparations(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanRecipeStepPreparations(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanRecipeStepPreparations(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRecipeStepPreparationExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "SELECT EXISTS ( SELECT recipe_step_preparations.id FROM recipe_step_preparations JOIN recipe_steps ON recipe_step_preparations.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_preparations.belongs_to_recipe_step = $1 AND recipe_step_preparations.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5 )"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepPreparation.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildRecipeStepPreparationExistsQuery(exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RecipeStepPreparationExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT recipe_step_preparations.id FROM recipe_step_preparations JOIN recipe_steps ON recipe_step_preparations.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_preparations.belongs_to_recipe_step = $1 AND recipe_step_preparations.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepPreparation.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RecipeStepPreparationExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepPreparation.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RecipeStepPreparationExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "SELECT recipe_step_preparations.id, recipe_step_preparations.valid_preparation_id, recipe_step_preparations.notes, recipe_step_preparations.created_on, recipe_step_preparations.updated_on, recipe_step_preparations.archived_on, recipe_step_preparations.belongs_to_recipe_step FROM recipe_step_preparations JOIN recipe_steps ON recipe_step_preparations.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_preparations.belongs_to_recipe_step = $1 AND recipe_step_preparations.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepPreparation.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepPreparationQuery(exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepPreparation(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT recipe_step_preparations.id, recipe_step_preparations.valid_preparation_id, recipe_step_preparations.notes, recipe_step_preparations.created_on, recipe_step_preparations.updated_on, recipe_step_preparations.archived_on, recipe_step_preparations.belongs_to_recipe_step FROM recipe_step_preparations JOIN recipe_steps ON recipe_step_preparations.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_preparations.belongs_to_recipe_step = $1 AND recipe_step_preparations.id = $2 AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.id = $5"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepPreparation.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromRecipeStepPreparation(exampleRecipeStepPreparation))

		actual, err := p.GetRecipeStepPreparation(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepPreparation, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepPreparation.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepPreparation(ctx, exampleRecipe.ID, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeStepPreparationsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(recipe_step_preparations.id) FROM recipe_step_preparations WHERE recipe_step_preparations.archived_on IS NULL"
		actualQuery := p.buildGetAllRecipeStepPreparationsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeStepPreparationsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_step_preparations.id) FROM recipe_step_preparations WHERE recipe_step_preparations.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeStepPreparationsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepPreparationsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_step_preparations.id, recipe_step_preparations.valid_preparation_id, recipe_step_preparations.notes, recipe_step_preparations.created_on, recipe_step_preparations.updated_on, recipe_step_preparations.archived_on, recipe_step_preparations.belongs_to_recipe_step, (SELECT COUNT(recipe_step_preparations.id) FROM recipe_step_preparations WHERE recipe_step_preparations.archived_on IS NULL) FROM recipe_step_preparations JOIN recipe_steps ON recipe_step_preparations.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_preparations.archived_on IS NULL AND recipe_step_preparations.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 AND recipe_step_preparations.created_on > $5 AND recipe_step_preparations.created_on < $6 AND recipe_step_preparations.updated_on > $7 AND recipe_step_preparations.updated_on < $8 ORDER BY recipe_step_preparations.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepPreparationsQuery(exampleRecipe.ID, exampleRecipeStep.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStepPreparations(T *testing.T) {
	T.Parallel()

	expectedListQuery := "SELECT recipe_step_preparations.id, recipe_step_preparations.valid_preparation_id, recipe_step_preparations.notes, recipe_step_preparations.created_on, recipe_step_preparations.updated_on, recipe_step_preparations.archived_on, recipe_step_preparations.belongs_to_recipe_step, (SELECT COUNT(recipe_step_preparations.id) FROM recipe_step_preparations WHERE recipe_step_preparations.archived_on IS NULL) FROM recipe_step_preparations JOIN recipe_steps ON recipe_step_preparations.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_preparations.archived_on IS NULL AND recipe_step_preparations.belongs_to_recipe_step = $1 AND recipe_steps.belongs_to_recipe = $2 AND recipe_steps.id = $3 AND recipes.id = $4 ORDER BY recipe_step_preparations.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepPreparationList := fakemodels.BuildFakeRecipeStepPreparationList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeStepPreparation(
					&exampleRecipeStepPreparationList.RecipeStepPreparations[0],
					&exampleRecipeStepPreparationList.RecipeStepPreparations[1],
					&exampleRecipeStepPreparationList.RecipeStepPreparations[2],
				),
			)

		actual, err := p.GetRecipeStepPreparations(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepPreparationList, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStepPreparations(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error executing read query", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(errors.New("blah"))

		actual, err := p.GetRecipeStepPreparations(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step preparation", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStepPreparation(exampleRecipeStepPreparation))

		actual, err := p.GetRecipeStepPreparations(ctx, exampleRecipe.ID, exampleRecipeStep.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeStepPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "INSERT INTO recipe_step_preparations (valid_preparation_id,notes,belongs_to_recipe_step) VALUES ($1,$2,$3) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipeStepPreparation.ValidPreparationID,
			exampleRecipeStepPreparation.Notes,
			exampleRecipeStepPreparation.BelongsToRecipeStep,
		}
		actualQuery, actualArgs := p.buildCreateRecipeStepPreparationQuery(exampleRecipeStepPreparation)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRecipeStepPreparation(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO recipe_step_preparations (valid_preparation_id,notes,belongs_to_recipe_step) VALUES ($1,$2,$3) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRecipeStepPreparation.ID, exampleRecipeStepPreparation.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStepPreparation.ValidPreparationID,
				exampleRecipeStepPreparation.Notes,
				exampleRecipeStepPreparation.BelongsToRecipeStep,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeStepPreparation(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepPreparation, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID
		exampleInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStepPreparation.ValidPreparationID,
				exampleRecipeStepPreparation.Notes,
				exampleRecipeStepPreparation.BelongsToRecipeStep,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipeStepPreparation(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeStepPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "UPDATE recipe_step_preparations SET valid_preparation_id = $1, notes = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_step = $3 AND id = $4 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleRecipeStepPreparation.ValidPreparationID,
			exampleRecipeStepPreparation.Notes,
			exampleRecipeStepPreparation.BelongsToRecipeStep,
			exampleRecipeStepPreparation.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRecipeStepPreparationQuery(exampleRecipeStepPreparation)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRecipeStepPreparation(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_step_preparations SET valid_preparation_id = $1, notes = $2, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe_step = $3 AND id = $4 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStepPreparation.ValidPreparationID,
				exampleRecipeStepPreparation.Notes,
				exampleRecipeStepPreparation.BelongsToRecipeStep,
				exampleRecipeStepPreparation.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeStepPreparation(ctx, exampleRecipeStepPreparation)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStepPreparation.ValidPreparationID,
				exampleRecipeStepPreparation.Notes,
				exampleRecipeStepPreparation.BelongsToRecipeStep,
				exampleRecipeStepPreparation.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipeStepPreparation(ctx, exampleRecipeStepPreparation)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeStepPreparationQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		expectedQuery := "UPDATE recipe_step_preparations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRecipeStep.ID,
			exampleRecipeStepPreparation.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRecipeStepPreparationQuery(exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRecipeStepPreparation(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_step_preparations SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe_step = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepPreparation.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeStepPreparation(ctx, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepPreparation.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRecipeStepPreparation(ctx, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)
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
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
		exampleRecipeStepPreparation.BelongsToRecipeStep = exampleRecipeStep.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.ID,
				exampleRecipeStepPreparation.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeStepPreparation(ctx, exampleRecipeStep.ID, exampleRecipeStepPreparation.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
