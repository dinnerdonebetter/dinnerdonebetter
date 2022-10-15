package postgres

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func buildMockRowsFromRecipePrepTasks(recipePrepTasks ...*types.RecipePrepTask) *sqlmock.Rows {
	columns := []string{
		"recipe_prep_tasks.id",
		"recipe_prep_tasks.notes",
		"recipe_prep_tasks.explicit_storage_instructions",
		"recipe_prep_tasks.minimum_time_buffer_before_recipe_in_seconds",
		"recipe_prep_tasks.maximum_time_buffer_before_recipe_in_seconds",
		"recipe_prep_tasks.storage_type",
		"recipe_prep_tasks.minimum_storage_temperature_in_celsius",
		"recipe_prep_tasks.maximum_storage_temperature_in_celsius",
		"recipe_prep_tasks.belongs_to_recipe",
		"recipe_prep_tasks.created_at",
		"recipe_prep_tasks.last_updated_at",
		"recipe_prep_tasks.archived_at",
		"recipe_prep_task_steps.id",
		"recipe_prep_task_steps.belongs_to_recipe_step",
		"recipe_prep_task_steps.belongs_to_recipe_prep_task",
		"recipe_prep_task_steps.satisfies_recipe_step",
	}

	exampleRows := sqlmock.NewRows(columns)

	for _, x := range recipePrepTasks {
		for _, y := range x.TaskSteps {
			rowValues := []driver.Value{
				x.ID,
				x.Notes,
				x.ExplicitStorageInstructions,
				x.MinimumTimeBufferBeforeRecipeInSeconds,
				x.MaximumTimeBufferBeforeRecipeInSeconds,
				x.StorageType,
				x.MinimumStorageTemperatureInCelsius,
				x.MaximumStorageTemperatureInCelsius,
				x.BelongsToRecipe,
				x.CreatedAt,
				x.LastUpdatedAt,
				x.ArchivedAt,
				y.ID,
				y.BelongsToRecipeStep,
				y.BelongsToRecipePrepTask,
				y.SatisfiesRecipeStep,
			}

			exampleRows.AddRow(rowValues...)
		}
	}

	return exampleRows
}

func TestQuerier_scanRecipePrepTaskWithRecipes(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		//
	})
}

func TestQuerier_scanRecipePrepTasksWithSteps(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		//
	})
}

func TestQuerier_RecipePrepTaskExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		c, db := buildTestClient(t)
		args := []interface{}{
			exampleRecipe.ID,
			exampleRecipePrepTask.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(recipePrepTasksExistsQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		actual, err := c.RecipePrepTaskExists(ctx, exampleRecipe.ID, exampleRecipePrepTask.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipePrepTask.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(getRecipePrepTasksQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnRows(buildMockRowsFromRecipePrepTasks(exampleRecipePrepTask))

		actual, err := c.GetRecipePrepTask(ctx, exampleRecipe.ID, exampleRecipePrepTask.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleRecipePrepTask, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expected := fakes.BuildFakeRecipePrepTask()
		exampleInput := fakes.BuildFakeRecipePrepTaskDatabaseCreationInputFromRecipePrepTask(expected)

		c, db := buildTestClient(t)

		createRecipePrepTaskQueryArgs := []interface{}{
			exampleInput.ID,
			exampleInput.Notes,
			exampleInput.ExplicitStorageInstructions,
			exampleInput.MinimumTimeBufferBeforeRecipeInSeconds,
			exampleInput.MaximumTimeBufferBeforeRecipeInSeconds,
			exampleInput.StorageType,
			exampleInput.MinimumStorageTemperatureInCelsius,
			exampleInput.MaximumStorageTemperatureInCelsius,
			exampleInput.BelongsToRecipe,
		}

		db.ExpectBegin()

		db.ExpectExec(formatQueryForSQLMock(createRecipePrepTaskQuery)).
			WithArgs(interfaceToDriverValue(createRecipePrepTaskQueryArgs)...).
			WillReturnResult(newArbitraryDatabaseResult())

		c.timeFunc = func() time.Time {
			return expected.CreatedAt
		}

		for _, taskStep := range exampleInput.TaskSteps {
			createRecipePrepTaskStepArgs := []interface{}{
				taskStep.ID,
				taskStep.BelongsToRecipePrepTask,
				taskStep.BelongsToRecipeStep,
				taskStep.SatisfiesRecipeStep,
			}

			db.ExpectExec(formatQueryForSQLMock(createRecipePrepTaskStepQuery)).
				WithArgs(interfaceToDriverValue(createRecipePrepTaskStepArgs)...).
				WillReturnResult(newArbitraryDatabaseResult())
		}

		db.ExpectCommit()

		actual, err := c.CreateRecipePrepTask(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_createRecipePrepTaskStep(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		expected := fakes.BuildFakeRecipePrepTaskStep()
		exampleInput := fakes.BuildFakeRecipePrepTaskStepDatabaseCreationInputFromRecipePrepTaskStep(expected)

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleInput.ID,
			exampleInput.BelongsToRecipePrepTask,
			exampleInput.BelongsToRecipeStep,
			exampleInput.SatisfiesRecipeStep,
		}

		db.ExpectBegin()
		tx, err := c.DB().Begin()
		require.NoError(t, err)

		db.ExpectExec(formatQueryForSQLMock(createRecipePrepTaskStepQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		actual, err := c.createRecipePrepTaskStep(ctx, tx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetRecipePrepTasksForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipe := fakes.BuildFakeRecipe()
		expected := fakes.BuildFakeRecipePrepTaskList().RecipePrepTasks

		c, db := buildTestClient(t)

		listRecipePrepTasksForRecipeArgs := []interface{}{
			exampleRecipe.ID,
		}

		db.ExpectQuery(formatQueryForSQLMock(listRecipePrepTasksForRecipeQuery)).
			WithArgs(interfaceToDriverValue(listRecipePrepTasksForRecipeArgs)...).
			WillReturnRows(buildMockRowsFromRecipePrepTasks(expected...))

		actual, err := c.GetRecipePrepTasksForRecipe(ctx, exampleRecipe.ID)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipePrepTask.Notes,
			exampleRecipePrepTask.ExplicitStorageInstructions,
			exampleRecipePrepTask.MinimumTimeBufferBeforeRecipeInSeconds,
			exampleRecipePrepTask.MaximumTimeBufferBeforeRecipeInSeconds,
			exampleRecipePrepTask.StorageType,
			exampleRecipePrepTask.MinimumStorageTemperatureInCelsius,
			exampleRecipePrepTask.MaximumStorageTemperatureInCelsius,
			exampleRecipePrepTask.BelongsToRecipe,
			exampleRecipePrepTask.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(updateRecipePrepStepTaskQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.UpdateRecipePrepTask(ctx, exampleRecipePrepTask))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipe := fakes.BuildFakeRecipe()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		c, db := buildTestClient(t)

		args := []interface{}{
			exampleRecipePrepTask.ID,
		}

		db.ExpectExec(formatQueryForSQLMock(archiveRecipePrepStepTaskQuery)).
			WithArgs(interfaceToDriverValue(args)...).
			WillReturnResult(newArbitraryDatabaseResult())

		assert.NoError(t, c.ArchiveRecipePrepTask(ctx, exampleRecipe.ID, exampleRecipePrepTask.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
