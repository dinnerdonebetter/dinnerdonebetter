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

func buildMockRowsFromRecipeIterationStep(recipeIterationSteps ...*models.RecipeIterationStep) *sqlmock.Rows {
	includeCount := len(recipeIterationSteps) > 1
	columns := recipeIterationStepsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}
	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeIterationSteps {
		rowValues := []driver.Value{
			x.ID,
			x.StartedOn,
			x.EndedOn,
			x.State,
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipe,
		}

		if includeCount {
			rowValues = append(rowValues, len(recipeIterationSteps))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipeIterationStep(x *models.RecipeIterationStep) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeIterationStepsTableColumns).AddRow(
		x.ArchivedOn,
		x.StartedOn,
		x.EndedOn,
		x.State,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsToRecipe,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRecipeIterationSteps(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanRecipeIterationSteps(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanRecipeIterationSteps(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRecipeIterationStepExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "SELECT EXISTS ( SELECT recipe_iteration_steps.id FROM recipe_iteration_steps JOIN recipes ON recipe_iteration_steps.belongs_to_recipe=recipes.id WHERE recipe_iteration_steps.belongs_to_recipe = $1 AND recipe_iteration_steps.id = $2 AND recipes.id = $3 )"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeIterationStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildRecipeIterationStepExistsQuery(exampleRecipe.ID, exampleRecipeIterationStep.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RecipeIterationStepExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT recipe_iteration_steps.id FROM recipe_iteration_steps JOIN recipes ON recipe_iteration_steps.belongs_to_recipe=recipes.id WHERE recipe_iteration_steps.belongs_to_recipe = $1 AND recipe_iteration_steps.id = $2 AND recipes.id = $3 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIterationStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RecipeIterationStepExists(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with no rows", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIterationStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RecipeIterationStepExists(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeIterationStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "SELECT recipe_iteration_steps.id, recipe_iteration_steps.started_on, recipe_iteration_steps.ended_on, recipe_iteration_steps.state, recipe_iteration_steps.created_on, recipe_iteration_steps.updated_on, recipe_iteration_steps.archived_on, recipe_iteration_steps.belongs_to_recipe FROM recipe_iteration_steps JOIN recipes ON recipe_iteration_steps.belongs_to_recipe=recipes.id WHERE recipe_iteration_steps.belongs_to_recipe = $1 AND recipe_iteration_steps.id = $2 AND recipes.id = $3"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeIterationStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeIterationStepQuery(exampleRecipe.ID, exampleRecipeIterationStep.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeIterationStep(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT recipe_iteration_steps.id, recipe_iteration_steps.started_on, recipe_iteration_steps.ended_on, recipe_iteration_steps.state, recipe_iteration_steps.created_on, recipe_iteration_steps.updated_on, recipe_iteration_steps.archived_on, recipe_iteration_steps.belongs_to_recipe FROM recipe_iteration_steps JOIN recipes ON recipe_iteration_steps.belongs_to_recipe=recipes.id WHERE recipe_iteration_steps.belongs_to_recipe = $1 AND recipe_iteration_steps.id = $2 AND recipes.id = $3"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIterationStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromRecipeIterationStep(exampleRecipeIterationStep))

		actual, err := p.GetRecipeIterationStep(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationStep, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIterationStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeIterationStep(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeIterationStepsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(recipe_iteration_steps.id) FROM recipe_iteration_steps WHERE recipe_iteration_steps.archived_on IS NULL"
		actualQuery := p.buildGetAllRecipeIterationStepsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeIterationStepsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_iteration_steps.id) FROM recipe_iteration_steps WHERE recipe_iteration_steps.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeIterationStepsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeIterationStepsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_iteration_steps.id, recipe_iteration_steps.started_on, recipe_iteration_steps.ended_on, recipe_iteration_steps.state, recipe_iteration_steps.created_on, recipe_iteration_steps.updated_on, recipe_iteration_steps.archived_on, recipe_iteration_steps.belongs_to_recipe, COUNT(recipe_iteration_steps.id) FROM recipe_iteration_steps JOIN recipes ON recipe_iteration_steps.belongs_to_recipe=recipes.id WHERE recipe_iteration_steps.archived_on IS NULL AND recipe_iteration_steps.belongs_to_recipe = $1 AND recipes.id = $2 AND recipe_iteration_steps.created_on > $3 AND recipe_iteration_steps.created_on < $4 AND recipe_iteration_steps.updated_on > $5 AND recipe_iteration_steps.updated_on < $6 GROUP BY recipe_iteration_steps.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRecipeIterationStepsQuery(exampleRecipe.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeIterationSteps(T *testing.T) {
	T.Parallel()

	expectedListQuery := "SELECT recipe_iteration_steps.id, recipe_iteration_steps.started_on, recipe_iteration_steps.ended_on, recipe_iteration_steps.state, recipe_iteration_steps.created_on, recipe_iteration_steps.updated_on, recipe_iteration_steps.archived_on, recipe_iteration_steps.belongs_to_recipe, COUNT(recipe_iteration_steps.id) FROM recipe_iteration_steps JOIN recipes ON recipe_iteration_steps.belongs_to_recipe=recipes.id WHERE recipe_iteration_steps.archived_on IS NULL AND recipe_iteration_steps.belongs_to_recipe = $1 AND recipes.id = $2 GROUP BY recipe_iteration_steps.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeIterationStepList := fakemodels.BuildFakeRecipeIterationStepList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeIterationStep(
					&exampleRecipeIterationStepList.RecipeIterationSteps[0],
					&exampleRecipeIterationStepList.RecipeIterationSteps[1],
					&exampleRecipeIterationStepList.RecipeIterationSteps[2],
				),
			)

		actual, err := p.GetRecipeIterationSteps(ctx, exampleRecipe.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationStepList, actual)

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

		actual, err := p.GetRecipeIterationSteps(ctx, exampleRecipe.ID, filter)
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

		actual, err := p.GetRecipeIterationSteps(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe iteration step", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeIterationStep(exampleRecipeIterationStep))

		actual, err := p.GetRecipeIterationSteps(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeIterationStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "INSERT INTO recipe_iteration_steps (started_on,ended_on,state,belongs_to_recipe) VALUES ($1,$2,$3,$4) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipeIterationStep.StartedOn,
			exampleRecipeIterationStep.EndedOn,
			exampleRecipeIterationStep.State,
			exampleRecipeIterationStep.BelongsToRecipe,
		}
		actualQuery, actualArgs := p.buildCreateRecipeIterationStepQuery(exampleRecipeIterationStep)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRecipeIterationStep(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO recipe_iteration_steps (started_on,ended_on,state,belongs_to_recipe) VALUES ($1,$2,$3,$4) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRecipeIterationStep.ID, exampleRecipeIterationStep.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeIterationStep.StartedOn,
				exampleRecipeIterationStep.EndedOn,
				exampleRecipeIterationStep.State,
				exampleRecipeIterationStep.BelongsToRecipe,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeIterationStep(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeIterationStep, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeIterationStep.StartedOn,
				exampleRecipeIterationStep.EndedOn,
				exampleRecipeIterationStep.State,
				exampleRecipeIterationStep.BelongsToRecipe,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipeIterationStep(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeIterationStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "UPDATE recipe_iteration_steps SET started_on = $1, ended_on = $2, state = $3, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $4 AND id = $5 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleRecipeIterationStep.StartedOn,
			exampleRecipeIterationStep.EndedOn,
			exampleRecipeIterationStep.State,
			exampleRecipeIterationStep.BelongsToRecipe,
			exampleRecipeIterationStep.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRecipeIterationStepQuery(exampleRecipeIterationStep)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRecipeIterationStep(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_iteration_steps SET started_on = $1, ended_on = $2, state = $3, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $4 AND id = $5 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIterationStep.StartedOn,
				exampleRecipeIterationStep.EndedOn,
				exampleRecipeIterationStep.State,
				exampleRecipeIterationStep.BelongsToRecipe,
				exampleRecipeIterationStep.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeIterationStep(ctx, exampleRecipeIterationStep)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error writing to database", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeIterationStep.StartedOn,
				exampleRecipeIterationStep.EndedOn,
				exampleRecipeIterationStep.State,
				exampleRecipeIterationStep.BelongsToRecipe,
				exampleRecipeIterationStep.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipeIterationStep(ctx, exampleRecipeIterationStep)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeIterationStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "UPDATE recipe_iteration_steps SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeIterationStep.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRecipeIterationStepQuery(exampleRecipe.ID, exampleRecipeIterationStep.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRecipeIterationStep(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_iteration_steps SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIterationStep.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeIterationStep(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)
		assert.NoError(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("returns sql.ErrNoRows with no rows affected", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIterationStep.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRecipeIterationStep(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)
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
		exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
		exampleRecipeIterationStep.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeIterationStep.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeIterationStep(ctx, exampleRecipe.ID, exampleRecipeIterationStep.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
