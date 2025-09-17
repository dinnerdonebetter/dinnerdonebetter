package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"
	"github.com/dinnerdonebetter/backend/pkg/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipePrepTaskSliceEquality(t *testing.T, expected, actual []*mealplanning.RecipePrepTask) {
	t.Helper()

	require.Equal(t, len(expected), len(actual), "expected recipe prep tasks length")
	for i := range expected {
		checkRecipePrepTaskEquality(t, i, expected[i], actual[i])
	}
}

func checkRecipePrepTaskEquality(t *testing.T, taskIndex int, expected, actual *mealplanning.RecipePrepTask) {
	t.Helper()

	assert.NotEmpty(t, actual.ID, "expected prep task %d to have ID", taskIndex)
	assert.False(t, actual.CreatedAt.IsZero(), "expected prep task %d to have CreatedAt", taskIndex)
	assert.NotEmpty(t, actual.BelongsToRecipe, "expected prep task %d to have BelongsToRecipe", taskIndex)
	assert.Equal(t, expected.Name, actual.Name, "expected prep task %d Name", taskIndex)
	assert.Equal(t, expected.Description, actual.Description, "expected prep task %d Description", taskIndex)
	assert.Equal(t, expected.StorageType, actual.StorageType, "expected prep task %d StorageType", taskIndex)
	assert.Equal(t, expected.ExplicitStorageInstructions, actual.ExplicitStorageInstructions, "expected prep task %d ExplicitStorageInstructions", taskIndex)
	assert.Equal(t, expected.Notes, actual.Notes, "expected prep task %d Notes", taskIndex)
	assert.Equal(t, expected.Optional, actual.Optional, "expected prep task %d Optional", taskIndex)
	assert.Equal(t, expected.StorageTemperatureInCelsius, actual.StorageTemperatureInCelsius, "expected prep task %d StorageTemperatureInCelsius", taskIndex)
	assert.Equal(t, expected.TimeBufferBeforeRecipeInSeconds, actual.TimeBufferBeforeRecipeInSeconds, "expected prep task %d TimeBufferBeforeRecipeInSeconds", taskIndex)
	checkRecipePrepTaskStepSliceEquality(t, taskIndex, expected.TaskSteps, actual.TaskSteps)
}

func createRecipePrepTaskForTest(t *testing.T, userClient client.Client) (*mealplanning.Recipe, *mealplanning.RecipePrepTask) {
	ctx := t.Context()

	_, _, createdRecipe := createRecipeForTest(t, nil)

	var createdRecipeStep *mealplanning.RecipeStep
	for _, step := range createdRecipe.Steps {
		createdRecipeStep = step
		break
	}
	require.NotNil(t, createdRecipeStep)

	exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()
	exampleRecipePrepTask.BelongsToRecipe = createdRecipe.ID
	exampleRecipePrepTask.TaskSteps = []*mealplanning.RecipePrepTaskStep{
		{
			BelongsToRecipeStep:     createdRecipeStep.ID,
			BelongsToRecipePrepTask: exampleRecipePrepTask.ID,
			SatisfiesRecipeStep:     true,
		},
	}

	exampleInput := mpconverters.ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(exampleRecipePrepTask)

	createdRecipePrepTaskRes, err := userClient.CreateRecipePrepTask(ctx, &mealplanninggrpc.CreateRecipePrepTaskRequest{
		RecipeID: createdRecipe.ID,
		Input:    converters.ConvertRecipePrepTaskCreationRequestInputToGRPCRecipePrepTaskCreationRequestInput(exampleInput),
	})
	require.NoError(t, err)
	createdRecipePrepTask := converters.ConvertGRPCRecipePrepTaskToRecipePrepTask(createdRecipePrepTaskRes.Created)

	retrievedRes, err := userClient.GetRecipePrepTask(ctx, &mealplanninggrpc.GetRecipePrepTaskRequest{
		RecipeID:         createdRecipe.ID,
		RecipePrepTaskID: createdRecipePrepTask.ID,
	})
	require.NoError(t, err)

	actual := converters.ConvertGRPCRecipePrepTaskToRecipePrepTask(retrievedRes.Result)
	checkRecipePrepTaskEquality(t, -1, exampleRecipePrepTask, actual)

	return createdRecipe, actual
}

func TestRecipePrepTasks_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		createdRecipe, actual := createRecipePrepTaskForTest(t, userClient)

		newRecipePrepTask := fakes.BuildFakeRecipePrepTask()
		newRecipePrepTask.ID = actual.ID
		newRecipePrepTask.BelongsToRecipe = createdRecipe.ID
		newRecipePrepTask.TaskSteps = actual.TaskSteps

		updateInput := mpconverters.ConvertRecipePrepTaskToRecipePrepTaskUpdateRequestInput(newRecipePrepTask)
		actual.Update(updateInput)

		_, err := adminClient.UpdateRecipePrepTask(ctx, &mealplanninggrpc.UpdateRecipePrepTaskRequest{
			RecipeID:         createdRecipe.ID,
			RecipePrepTaskID: actual.ID,
			Input:            converters.ConvertRecipePrepTaskUpdateRequestInputToGRPCRecipePrepTaskUpdateRequestInput(updateInput),
		})
		require.NoError(t, err)

		retrievedRes, err := userClient.GetRecipePrepTask(ctx, &mealplanninggrpc.GetRecipePrepTaskRequest{
			RecipeID:         createdRecipe.ID,
			RecipePrepTaskID: actual.ID,
		})
		require.NoError(t, err)

		actual = converters.ConvertGRPCRecipePrepTaskToRecipePrepTask(retrievedRes.Result)

		// assert recipe prep task equality
		checkRecipePrepTaskEquality(t, -1, newRecipePrepTask, actual)
		assert.NotNil(t, actual.LastUpdatedAt)

		_, err = userClient.ArchiveRecipePrepTask(ctx, &mealplanninggrpc.ArchiveRecipePrepTaskRequest{
			RecipeID:         createdRecipe.ID,
			RecipePrepTaskID: actual.ID,
		})
		assert.NoError(t, err)

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
		assert.NoError(t, err)
	})
}

func TestRecipePrepTasks_Listing(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		_, _, createdRecipe := createRecipeForTest(t, nil)

		var createdRecipeStep *mealplanning.RecipeStep
		for _, step := range createdRecipe.Steps {
			createdRecipeStep = step
			break
		}
		require.NotNil(t, createdRecipeStep)

		var expected []*mealplanning.RecipePrepTask
		for i := 0; i < 5; i++ {
			exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()
			exampleRecipePrepTask.BelongsToRecipe = createdRecipe.ID
			exampleRecipePrepTask.TaskSteps = []*mealplanning.RecipePrepTaskStep{
				{
					BelongsToRecipeStep:     createdRecipeStep.ID,
					BelongsToRecipePrepTask: exampleRecipePrepTask.ID,
					SatisfiesRecipeStep:     true,
				},
			}

			exampleInput := mpconverters.ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(exampleRecipePrepTask)

			createdRecipePrepTaskRes, err := adminClient.CreateRecipePrepTask(ctx, &mealplanninggrpc.CreateRecipePrepTaskRequest{
				RecipeID: createdRecipe.ID,
				Input:    converters.ConvertRecipePrepTaskCreationRequestInputToGRPCRecipePrepTaskCreationRequestInput(exampleInput),
			})
			require.NoError(t, err)

			createdRecipePrepTask := converters.ConvertGRPCRecipePrepTaskToRecipePrepTask(createdRecipePrepTaskRes.Created)

			for j := range createdRecipePrepTask.TaskSteps {
				exampleRecipePrepTask.TaskSteps[j].ID = createdRecipePrepTask.TaskSteps[j].ID
				exampleRecipePrepTask.TaskSteps[j].BelongsToRecipePrepTask = createdRecipePrepTask.ID
			}

			checkRecipePrepTaskEquality(t, -1, exampleRecipePrepTask, createdRecipePrepTask)

			retrievedRecipePrepTaskRes, err := userClient.GetRecipePrepTask(ctx, &mealplanninggrpc.GetRecipePrepTaskRequest{
				RecipeID:         createdRecipe.ID,
				RecipePrepTaskID: createdRecipePrepTask.ID,
			})
			require.NoError(t, err)

			createdRecipePrepTask = converters.ConvertGRPCRecipePrepTaskToRecipePrepTask(retrievedRecipePrepTaskRes.Result)
			require.Equal(t, createdRecipe.ID, createdRecipePrepTask.BelongsToRecipe)

			expected = append(expected, createdRecipePrepTask)
		}

		// assert recipe prep task list equality
		actual, err := userClient.GetRecipePrepTasks(ctx, &mealplanninggrpc.GetRecipePrepTasksRequest{RecipeID: createdRecipe.ID})
		require.NoError(t, err)
		assert.True(
			t,
			len(expected) <= len(actual.Results),
			"expected %d to be <= %d",
			len(expected),
			len(actual.Results),
		)

		for _, createdRecipePrepTask := range expected {
			_, err = userClient.ArchiveRecipePrepTask(ctx, &mealplanninggrpc.ArchiveRecipePrepTaskRequest{
				RecipeID:         createdRecipe.ID,
				RecipePrepTaskID: createdRecipePrepTask.ID,
			})
			assert.NoError(t, err)
		}

		_, err = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeID: createdRecipe.ID})
		assert.NoError(t, err)
	})
}
