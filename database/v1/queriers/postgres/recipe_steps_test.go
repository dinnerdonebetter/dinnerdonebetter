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

func buildMockRowsFromRecipeStep(recipeSteps ...*models.RecipeStep) *sqlmock.Rows {
	includeCount := len(recipeSteps) > 1
	columns := recipeStepsTableColumns

	if includeCount {
		columns = append(columns, "count")
	}
	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipeSteps {
		rowValues := []driver.Value{
			x.ID,
			x.Index,
			x.ValidPreparationID,
			x.PrerequisiteStepID,
			x.MinEstimatedTimeInSeconds,
			x.MaxEstimatedTimeInSeconds,
			x.YieldsProductName,
			x.YieldsQuantity,
			x.Notes,
			x.CreatedOn,
			x.UpdatedOn,
			x.ArchivedOn,
			x.BelongsToRecipe,
		}

		if includeCount {
			rowValues = append(rowValues, len(recipeSteps))
		}

		exampleRows.AddRow(rowValues...)
	}

	return exampleRows
}

func buildErroneousMockRowFromRecipeStep(x *models.RecipeStep) *sqlmock.Rows {
	exampleRows := sqlmock.NewRows(recipeStepsTableColumns).AddRow(
		x.ArchivedOn,
		x.Index,
		x.ValidPreparationID,
		x.PrerequisiteStepID,
		x.MinEstimatedTimeInSeconds,
		x.MaxEstimatedTimeInSeconds,
		x.YieldsProductName,
		x.YieldsQuantity,
		x.Notes,
		x.CreatedOn,
		x.UpdatedOn,
		x.BelongsToRecipe,
		x.ID,
	)

	return exampleRows
}

func TestPostgres_ScanRecipeSteps(T *testing.T) {
	T.Parallel()

	T.Run("surfaces row errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(errors.New("blah"))

		_, _, err := p.scanRecipeSteps(mockRows)
		assert.Error(t, err)
	})

	T.Run("logs row closing errors", func(t *testing.T) {
		p, _ := buildTestService(t)
		mockRows := &database.MockResultIterator{}

		mockRows.On("Next").Return(false)
		mockRows.On("Err").Return(nil)
		mockRows.On("Close").Return(errors.New("blah"))

		_, _, err := p.scanRecipeSteps(mockRows)
		assert.NoError(t, err)
	})
}

func TestPostgres_buildRecipeStepExistsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "SELECT EXISTS ( SELECT recipe_steps.id FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.id = $3 )"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildRecipeStepExistsQuery(exampleRecipe.ID, exampleRecipeStep.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_RecipeStepExists(T *testing.T) {
	T.Parallel()

	expectedQuery := "SELECT EXISTS ( SELECT recipe_steps.id FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.id = $3 )"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := p.RecipeStepExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
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

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.RecipeStepExists(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "SELECT recipe_steps.id, recipe_steps.index, recipe_steps.valid_preparation_id, recipe_steps.prerequisite_step_id, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.yields_product_name, recipe_steps.yields_quantity, recipe_steps.notes, recipe_steps.created_on, recipe_steps.updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.id = $3"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeStep.ID,
			exampleRecipe.ID,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepQuery(exampleRecipe.ID, exampleRecipeStep.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeStep(T *testing.T) {
	T.Parallel()

	exampleUser := fakemodels.BuildFakeUser()
	expectedQuery := "SELECT recipe_steps.id, recipe_steps.index, recipe_steps.valid_preparation_id, recipe_steps.prerequisite_step_id, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.yields_product_name, recipe_steps.yields_quantity, recipe_steps.notes, recipe_steps.created_on, recipe_steps.updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.belongs_to_recipe = $1 AND recipe_steps.id = $2 AND recipes.id = $3"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildMockRowsFromRecipeStep(exampleRecipeStep))

		actual, err := p.GetRecipeStep(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStep, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("surfaces sql.ErrNoRows", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
				exampleRecipe.ID,
			).
			WillReturnError(sql.ErrNoRows)

		actual, err := p.GetRecipeStep(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)
		assert.Equal(t, sql.ErrNoRows, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetAllRecipeStepsCountQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		expectedQuery := "SELECT COUNT(recipe_steps.id) FROM recipe_steps WHERE recipe_steps.archived_on IS NULL"
		actualQuery := p.buildGetAllRecipeStepsCountQuery()

		ensureArgCountMatchesQuery(t, actualQuery, []interface{}{})
		assert.Equal(t, expectedQuery, actualQuery)
	})
}

func TestPostgres_GetAllRecipeStepsCount(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedQuery := "SELECT COUNT(recipe_steps.id) FROM recipe_steps WHERE recipe_steps.archived_on IS NULL"
		expectedCount := uint64(123)

		p, mockDB := buildTestService(t)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(expectedCount))

		actualCount, err := p.GetAllRecipeStepsCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, actualCount)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildGetRecipeStepsQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleRecipe := fakemodels.BuildFakeRecipe()
		filter := fakemodels.BuildFleshedOutQueryFilter()

		expectedQuery := "SELECT recipe_steps.id, recipe_steps.index, recipe_steps.valid_preparation_id, recipe_steps.prerequisite_step_id, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.yields_product_name, recipe_steps.yields_quantity, recipe_steps.notes, recipe_steps.created_on, recipe_steps.updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe, COUNT(recipe_steps.id) FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipes.id = $2 AND recipe_steps.created_on > $3 AND recipe_steps.created_on < $4 AND recipe_steps.updated_on > $5 AND recipe_steps.updated_on < $6 GROUP BY recipe_steps.id LIMIT 20 OFFSET 180"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipe.ID,
			filter.CreatedAfter,
			filter.CreatedBefore,
			filter.UpdatedAfter,
			filter.UpdatedBefore,
		}
		actualQuery, actualArgs := p.buildGetRecipeStepsQuery(exampleRecipe.ID, filter)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_GetRecipeSteps(T *testing.T) {
	T.Parallel()

	expectedListQuery := "SELECT recipe_steps.id, recipe_steps.index, recipe_steps.valid_preparation_id, recipe_steps.prerequisite_step_id, recipe_steps.min_estimated_time_in_seconds, recipe_steps.max_estimated_time_in_seconds, recipe_steps.yields_product_name, recipe_steps.yields_quantity, recipe_steps.notes, recipe_steps.created_on, recipe_steps.updated_on, recipe_steps.archived_on, recipe_steps.belongs_to_recipe, COUNT(recipe_steps.id) FROM recipe_steps JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_steps.archived_on IS NULL AND recipe_steps.belongs_to_recipe = $1 AND recipes.id = $2 GROUP BY recipe_steps.id LIMIT 20"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStepList := fakemodels.BuildFakeRecipeStepList()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(
				buildMockRowsFromRecipeStep(
					&exampleRecipeStepList.RecipeSteps[0],
					&exampleRecipeStepList.RecipeSteps[1],
					&exampleRecipeStepList.RecipeSteps[2],
				),
			)

		actual, err := p.GetRecipeSteps(ctx, exampleRecipe.ID, filter)

		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStepList, actual)

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

		actual, err := p.GetRecipeSteps(ctx, exampleRecipe.ID, filter)
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

		actual, err := p.GetRecipeSteps(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})

	T.Run("with error scanning recipe step", func(t *testing.T) {
		ctx := context.Background()

		exampleRecipe := fakemodels.BuildFakeRecipe()

		p, mockDB := buildTestService(t)
		filter := models.DefaultQueryFilter()

		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedListQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipe.ID,
			).
			WillReturnRows(buildErroneousMockRowFromRecipeStep(exampleRecipeStep))

		actual, err := p.GetRecipeSteps(ctx, exampleRecipe.ID, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildCreateRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "INSERT INTO recipe_steps (index,valid_preparation_id,prerequisite_step_id,min_estimated_time_in_seconds,max_estimated_time_in_seconds,yields_product_name,yields_quantity,notes,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_on"
		expectedArgs := []interface{}{
			exampleRecipeStep.Index,
			exampleRecipeStep.ValidPreparationID,
			exampleRecipeStep.PrerequisiteStepID,
			exampleRecipeStep.MinEstimatedTimeInSeconds,
			exampleRecipeStep.MaxEstimatedTimeInSeconds,
			exampleRecipeStep.YieldsProductName,
			exampleRecipeStep.YieldsQuantity,
			exampleRecipeStep.Notes,
			exampleRecipeStep.BelongsToRecipe,
		}
		actualQuery, actualArgs := p.buildCreateRecipeStepQuery(exampleRecipeStep)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_CreateRecipeStep(T *testing.T) {
	T.Parallel()

	expectedCreationQuery := "INSERT INTO recipe_steps (index,valid_preparation_id,prerequisite_step_id,min_estimated_time_in_seconds,max_estimated_time_in_seconds,yields_product_name,yields_quantity,notes,belongs_to_recipe) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id, created_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID
		exampleInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)

		exampleRows := sqlmock.NewRows([]string{"id", "created_on"}).AddRow(exampleRecipeStep.ID, exampleRecipeStep.CreatedOn)
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStep.Index,
				exampleRecipeStep.ValidPreparationID,
				exampleRecipeStep.PrerequisiteStepID,
				exampleRecipeStep.MinEstimatedTimeInSeconds,
				exampleRecipeStep.MaxEstimatedTimeInSeconds,
				exampleRecipeStep.YieldsProductName,
				exampleRecipeStep.YieldsQuantity,
				exampleRecipeStep.Notes,
				exampleRecipeStep.BelongsToRecipe,
			).WillReturnRows(exampleRows)

		actual, err := p.CreateRecipeStep(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipeStep, actual)

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
		exampleInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedCreationQuery)).
			WithArgs(
				exampleRecipeStep.Index,
				exampleRecipeStep.ValidPreparationID,
				exampleRecipeStep.PrerequisiteStepID,
				exampleRecipeStep.MinEstimatedTimeInSeconds,
				exampleRecipeStep.MaxEstimatedTimeInSeconds,
				exampleRecipeStep.YieldsProductName,
				exampleRecipeStep.YieldsQuantity,
				exampleRecipeStep.Notes,
				exampleRecipeStep.BelongsToRecipe,
			).WillReturnError(errors.New("blah"))

		actual, err := p.CreateRecipeStep(ctx, exampleInput)
		assert.Error(t, err)
		assert.Nil(t, actual)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildUpdateRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "UPDATE recipe_steps SET index = $1, valid_preparation_id = $2, prerequisite_step_id = $3, min_estimated_time_in_seconds = $4, max_estimated_time_in_seconds = $5, yields_product_name = $6, yields_quantity = $7, notes = $8, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $9 AND id = $10 RETURNING updated_on"
		expectedArgs := []interface{}{
			exampleRecipeStep.Index,
			exampleRecipeStep.ValidPreparationID,
			exampleRecipeStep.PrerequisiteStepID,
			exampleRecipeStep.MinEstimatedTimeInSeconds,
			exampleRecipeStep.MaxEstimatedTimeInSeconds,
			exampleRecipeStep.YieldsProductName,
			exampleRecipeStep.YieldsQuantity,
			exampleRecipeStep.Notes,
			exampleRecipeStep.BelongsToRecipe,
			exampleRecipeStep.ID,
		}
		actualQuery, actualArgs := p.buildUpdateRecipeStepQuery(exampleRecipeStep)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_UpdateRecipeStep(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_steps SET index = $1, valid_preparation_id = $2, prerequisite_step_id = $3, min_estimated_time_in_seconds = $4, max_estimated_time_in_seconds = $5, yields_product_name = $6, yields_quantity = $7, notes = $8, updated_on = extract(epoch FROM NOW()) WHERE belongs_to_recipe = $9 AND id = $10 RETURNING updated_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		exampleRows := sqlmock.NewRows([]string{"updated_on"}).AddRow(uint64(time.Now().Unix()))
		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.Index,
				exampleRecipeStep.ValidPreparationID,
				exampleRecipeStep.PrerequisiteStepID,
				exampleRecipeStep.MinEstimatedTimeInSeconds,
				exampleRecipeStep.MaxEstimatedTimeInSeconds,
				exampleRecipeStep.YieldsProductName,
				exampleRecipeStep.YieldsQuantity,
				exampleRecipeStep.Notes,
				exampleRecipeStep.BelongsToRecipe,
				exampleRecipeStep.ID,
			).WillReturnRows(exampleRows)

		err := p.UpdateRecipeStep(ctx, exampleRecipeStep)
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

		mockDB.ExpectQuery(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipeStep.Index,
				exampleRecipeStep.ValidPreparationID,
				exampleRecipeStep.PrerequisiteStepID,
				exampleRecipeStep.MinEstimatedTimeInSeconds,
				exampleRecipeStep.MaxEstimatedTimeInSeconds,
				exampleRecipeStep.YieldsProductName,
				exampleRecipeStep.YieldsQuantity,
				exampleRecipeStep.Notes,
				exampleRecipeStep.BelongsToRecipe,
				exampleRecipeStep.ID,
			).WillReturnError(errors.New("blah"))

		err := p.UpdateRecipeStep(ctx, exampleRecipeStep)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}

func TestPostgres_buildArchiveRecipeStepQuery(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		p, _ := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		expectedQuery := "UPDATE recipe_steps SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"
		expectedArgs := []interface{}{
			exampleRecipe.ID,
			exampleRecipeStep.ID,
		}
		actualQuery, actualArgs := p.buildArchiveRecipeStepQuery(exampleRecipe.ID, exampleRecipeStep.ID)

		ensureArgCountMatchesQuery(t, actualQuery, actualArgs)
		assert.Equal(t, expectedQuery, actualQuery)
		assert.Equal(t, expectedArgs, actualArgs)
	})
}

func TestPostgres_ArchiveRecipeStep(T *testing.T) {
	T.Parallel()

	expectedQuery := "UPDATE recipe_steps SET updated_on = extract(epoch FROM NOW()), archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_recipe = $1 AND id = $2 RETURNING archived_on"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		p, mockDB := buildTestService(t)

		exampleUser := fakemodels.BuildFakeUser()
		exampleRecipe := fakemodels.BuildFakeRecipe()
		exampleRecipe.BelongsToUser = exampleUser.ID
		exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
		exampleRecipeStep.BelongsToRecipe = exampleRecipe.ID

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
			).WillReturnResult(sqlmock.NewResult(1, 1))

		err := p.ArchiveRecipeStep(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
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

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
			).WillReturnResult(sqlmock.NewResult(0, 0))

		err := p.ArchiveRecipeStep(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
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

		mockDB.ExpectExec(formatQueryForSQLMock(expectedQuery)).
			WithArgs(
				exampleRecipe.ID,
				exampleRecipeStep.ID,
			).WillReturnError(errors.New("blah"))

		err := p.ArchiveRecipeStep(ctx, exampleRecipe.ID, exampleRecipeStep.ID)
		assert.Error(t, err)

		assert.NoError(t, mockDB.ExpectationsWereMet(), "not all database expectations were met")
	})
}
