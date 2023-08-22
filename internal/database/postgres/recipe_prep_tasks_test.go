package postgres

import (
	"context"
	"database/sql/driver"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildMockRowsFromRecipePrepTasks(recipePrepTasks ...*types.RecipePrepTask) *sqlmock.Rows {
	columns := []string{
		"recipe_prep_tasks.id",
		"recipe_prep_tasks.name",
		"recipe_prep_tasks.description",
		"recipe_prep_tasks.notes",
		"recipe_prep_tasks.optional",
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
				x.Name,
				x.Description,
				x.Notes,
				x.Optional,
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

func TestQuerier_RecipePrepTaskExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		c, _ := buildTestClient(t)

		actual, err := c.RecipePrepTaskExists(ctx, "", exampleRecipePrepTask.ID)
		assert.Error(t, err)
		assert.False(t, actual)
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

		args := []any{
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

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipePrepTask(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetRecipePrepTasksForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipe := fakes.BuildFakeRecipe()
		expected := fakes.BuildFakeRecipePrepTaskList().Data

		c, db := buildTestClient(t)

		listRecipePrepTasksForRecipeArgs := []any{
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

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		assert.Error(t, c.UpdateRecipePrepTask(ctx, nil))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("with missing recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		c, db := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipePrepTask(ctx, "", exampleRecipePrepTask.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
