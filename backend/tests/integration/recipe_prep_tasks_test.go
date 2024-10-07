package integration

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipePrepTaskEquality(t *testing.T, expected, actual *types.RecipePrepTask) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe prep task %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for recipe prep task %s to be %v, but it was %v", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for recipe prep task %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Optional, actual.Optional, "expected Optional for recipe prep task %s to be %v, but it was %v", expected.ID, expected.Optional, actual.Optional)
	assert.Equal(t, expected.ExplicitStorageInstructions, actual.ExplicitStorageInstructions, "expected ExplicitStorageInstructions for recipe prep task %s to be %v, but it was %v", expected.ID, expected.ExplicitStorageInstructions, actual.ExplicitStorageInstructions)
	assert.Equal(t, expected.StorageType, actual.StorageType, "expected StorageType for recipe prep task %s to be %v, but it was %v", expected.ID, expected.StorageType, actual.StorageType)
	assert.Equal(t, expected.BelongsToRecipe, actual.BelongsToRecipe, "expected BelongsToRecipe for recipe prep task %s to be %v, but it was %v", expected.ID, expected.BelongsToRecipe, actual.BelongsToRecipe)
	assert.Equal(t, expected.TimeBufferBeforeRecipeInSeconds, actual.TimeBufferBeforeRecipeInSeconds, "expected TimeBufferBeforeRecipeInSeconds for recipe prep task %s to be %v, but it was %v", expected.ID, expected.TimeBufferBeforeRecipeInSeconds, actual.TimeBufferBeforeRecipeInSeconds)
	assert.Equal(t, expected.StorageTemperatureInCelsius, actual.StorageTemperatureInCelsius, "expected StorageTemperatureInCelsius for recipe prep task %s to be %v, but it was %v", expected.ID, expected.StorageTemperatureInCelsius, actual.StorageTemperatureInCelsius)
	assert.NotZero(t, actual.CreatedAt)
}

func createRecipePrepTaskForTest(ctx context.Context, t *testing.T, adminClient, client *apiclient.Client) (*types.Recipe, *types.RecipePrepTask) {
	_, _, createdRecipe := createRecipeForTest(ctx, t, adminClient, client, nil)

	var createdRecipeStep *types.RecipeStep
	for _, step := range createdRecipe.Steps {
		createdRecipeStep = step
		break
	}
	require.NotNil(t, createdRecipeStep)

	exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()
	exampleRecipePrepTask.BelongsToRecipe = createdRecipe.ID
	exampleRecipePrepTask.TaskSteps = []*types.RecipePrepTaskStep{
		{
			BelongsToRecipeStep:     createdRecipeStep.ID,
			BelongsToRecipePrepTask: exampleRecipePrepTask.ID,
			SatisfiesRecipeStep:     true,
		},
	}

	exampleInput := converters.ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(exampleRecipePrepTask)

	createdRecipePrepTask, err := client.CreateRecipePrepTask(ctx, exampleInput)
	requireNotNilAndNoProblems(t, createdRecipePrepTask, err)

	actual, err := client.GetRecipePrepTask(ctx, createdRecipe.ID, createdRecipePrepTask.ID)
	requireNotNilAndNoProblems(t, actual, err)

	checkRecipePrepTaskEquality(t, exampleRecipePrepTask, actual)

	return createdRecipe, actual
}

func (s *TestSuite) TestRecipePrepTasks_CompleteLifecycle() {
	s.runTest("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			createdRecipe, actual := createRecipePrepTaskForTest(ctx, t, testClients.adminClient, testClients.userClient)

			newRecipePrepTask := fakes.BuildFakeRecipePrepTask()
			newRecipePrepTask.ID = actual.ID
			newRecipePrepTask.BelongsToRecipe = createdRecipe.ID
			newRecipePrepTask.TaskSteps = actual.TaskSteps
			actual.Update(converters.ConvertRecipePrepTaskToRecipePrepTaskUpdateRequestInput(newRecipePrepTask))
			require.NoError(t, testClients.adminClient.UpdateRecipePrepTask(ctx, actual))

			actual, err := testClients.userClient.GetRecipePrepTask(ctx, createdRecipe.ID, actual.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe prep task equality
			checkRecipePrepTaskEquality(t, newRecipePrepTask, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.userClient.ArchiveRecipePrepTask(ctx, createdRecipe.ID, actual.ID))

			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipePrepTasks_Listing() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)

			var createdRecipeStep *types.RecipeStep
			for _, step := range createdRecipe.Steps {
				createdRecipeStep = step
				break
			}
			require.NotNil(t, createdRecipeStep)

			var expected []*types.RecipePrepTask
			for i := 0; i < 5; i++ {
				exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()
				exampleRecipePrepTask.BelongsToRecipe = createdRecipe.ID
				exampleRecipePrepTask.TaskSteps = []*types.RecipePrepTaskStep{
					{
						BelongsToRecipeStep:     createdRecipeStep.ID,
						BelongsToRecipePrepTask: exampleRecipePrepTask.ID,
						SatisfiesRecipeStep:     true,
					},
				}

				exampleInput := converters.ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(exampleRecipePrepTask)

				createdRecipePrepTask, err := testClients.adminClient.CreateRecipePrepTask(ctx, exampleInput)
				requireNotNilAndNoProblems(t, createdRecipePrepTask, err)

				exampleRecipePrepTaskInput := converters.ConvertRecipePrepTaskToRecipePrepTaskCreationRequestInput(exampleRecipePrepTask)

				createdRecipePrepTask, createdRecipePrepTaskErr := testClients.adminClient.CreateRecipePrepTask(ctx, exampleRecipePrepTaskInput)
				require.NoError(t, createdRecipePrepTaskErr)

				for j := range createdRecipePrepTask.TaskSteps {
					exampleRecipePrepTask.TaskSteps[j].ID = createdRecipePrepTask.TaskSteps[j].ID
					exampleRecipePrepTask.TaskSteps[j].BelongsToRecipePrepTask = createdRecipePrepTask.ID
				}

				checkRecipePrepTaskEquality(t, exampleRecipePrepTask, createdRecipePrepTask)

				createdRecipePrepTask, createdRecipePrepTaskErr = testClients.userClient.GetRecipePrepTask(ctx, createdRecipe.ID, createdRecipePrepTask.ID)
				requireNotNilAndNoProblems(t, createdRecipePrepTask, createdRecipePrepTaskErr)
				require.Equal(t, createdRecipe.ID, createdRecipePrepTask.BelongsToRecipe)

				expected = append(expected, createdRecipePrepTask)
			}

			// assert recipe prep task list equality
			actual, err := testClients.userClient.GetRecipePrepTasks(ctx, createdRecipe.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			for _, createdRecipePrepTask := range expected {
				assert.NoError(t, testClients.userClient.ArchiveRecipePrepTask(ctx, createdRecipe.ID, createdRecipePrepTask.ID))
			}

			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
