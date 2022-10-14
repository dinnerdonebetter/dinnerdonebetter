package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types/fakes"

	"github.com/prixfixeco/api_server/pkg/types"
)

func checkRecipePrepTaskEquality(t *testing.T, expected, actual *types.RecipePrepTask) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for recipe prep task %s to be %v, but it was %v", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.ExplicitStorageInstructions, actual.ExplicitStorageInstructions, "expected ExplicitStorageInstructions for recipe prep task %s to be %v, but it was %v", expected.ID, expected.ExplicitStorageInstructions, actual.ExplicitStorageInstructions)
	assert.Equal(t, expected.StorageType, actual.StorageType, "expected StorageType for recipe prep task %s to be %v, but it was %v", expected.ID, expected.StorageType, actual.StorageType)
	assert.Equal(t, expected.BelongsToRecipe, actual.BelongsToRecipe, "expected BelongsToRecipe for recipe prep task %s to be %v, but it was %v", expected.ID, expected.BelongsToRecipe, actual.BelongsToRecipe)
	assert.Equal(t, expected.TaskSteps, actual.TaskSteps, "expected TaskSteps for recipe prep task %s to be %v, but it was %v", expected.ID, expected.TaskSteps, actual.TaskSteps)
	assert.Equal(t, expected.MinimumTimeBufferBeforeRecipeInSeconds, actual.MinimumTimeBufferBeforeRecipeInSeconds, "expected MinimumTimeBufferBeforeRecipeInSeconds for recipe prep task %s to be %v, but it was %v", expected.ID, expected.MinimumTimeBufferBeforeRecipeInSeconds, actual.MinimumTimeBufferBeforeRecipeInSeconds)
	assert.Equal(t, expected.MaximumStorageTemperatureInCelsius, actual.MaximumStorageTemperatureInCelsius, "expected MaximumStorageTemperatureInCelsius for recipe prep task %s to be %v, but it was %v", expected.ID, expected.MaximumStorageTemperatureInCelsius, actual.MaximumStorageTemperatureInCelsius)
	assert.Equal(t, expected.MaximumTimeBufferBeforeRecipeInSeconds, actual.MaximumTimeBufferBeforeRecipeInSeconds, "expected MaximumTimeBufferBeforeRecipeInSeconds for recipe prep task %s to be %v, but it was %v", expected.ID, expected.MaximumTimeBufferBeforeRecipeInSeconds, actual.MaximumTimeBufferBeforeRecipeInSeconds)
	assert.Equal(t, expected.MinimumStorageTemperatureInCelsius, actual.MinimumStorageTemperatureInCelsius, "expected MinimumStorageTemperatureInCelsius for recipe prep task %s to be %v, but it was %v", expected.ID, expected.MinimumStorageTemperatureInCelsius, actual.MinimumStorageTemperatureInCelsius)
	assert.NotZero(t, actual.CreatedAt)
}

// convertRecipePrepTaskToRecipePrepTaskUpdateInput creates an RecipePrepTaskUpdateRequestInput struct from a recipe prep task.
func convertRecipePrepTaskToRecipePrepTaskUpdateInput(x *types.RecipePrepTask) *types.RecipePrepTaskUpdateRequestInput {
	updateSteps := []*types.RecipePrepTaskStepUpdateRequestInput{}
	for _, taskStep := range x.TaskSteps {
		updateSteps = append(updateSteps, &types.RecipePrepTaskStepUpdateRequestInput{
			SatisfiesRecipeStep:     &taskStep.SatisfiesRecipeStep,
			BelongsToRecipeStep:     &taskStep.BelongsToRecipeStep,
			BelongsToRecipePrepTask: &taskStep.BelongsToRecipePrepTask,
			ID:                      taskStep.ID,
		})
	}

	return &types.RecipePrepTaskUpdateRequestInput{
		Notes:                                  &x.Notes,
		ExplicitStorageInstructions:            &x.ExplicitStorageInstructions,
		MinimumTimeBufferBeforeRecipeInSeconds: &x.MinimumTimeBufferBeforeRecipeInSeconds,
		MaximumTimeBufferBeforeRecipeInSeconds: &x.MaximumTimeBufferBeforeRecipeInSeconds,
		StorageType:                            &x.StorageType,
		MinimumStorageTemperatureInCelsius:     &x.MinimumStorageTemperatureInCelsius,
		MaximumStorageTemperatureInCelsius:     &x.MaximumStorageTemperatureInCelsius,
		BelongsToRecipe:                        &x.BelongsToRecipe,
		TaskSteps:                              updateSteps,
	}
}

func (s *TestSuite) TestRecipePrepTasks_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var createdRecipeStep *types.RecipeStep
			for _, step := range createdRecipe.Steps {
				createdRecipeStep = step
				break
			}
			require.NotNil(t, createdRecipeStep)

			t.Log("creating recipe prep task")
			exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()
			exampleRecipePrepTask.BelongsToRecipe = createdRecipe.ID
			exampleRecipePrepTask.TaskSteps = []*types.RecipePrepTaskStep{
				{
					BelongsToRecipeStep:     createdRecipeStep.ID,
					BelongsToRecipePrepTask: exampleRecipePrepTask.ID,
					SatisfiesRecipeStep:     true,
				},
			}

			exampleInput := fakes.BuildFakeRecipePrepTaskCreationRequestInputFromRecipePrepTask(exampleRecipePrepTask)

			createdRecipePrepTask, err := testClients.user.CreateRecipePrepTask(ctx, exampleInput)
			requireNotNilAndNoProblems(t, createdRecipePrepTask, err)

			t.Log("fetching changed recipe prep task")
			actual, err := testClients.user.GetRecipePrepTask(ctx, createdRecipe.ID, createdRecipePrepTask.ID)
			requireNotNilAndNoProblems(t, actual, err)

			t.Log("changing recipe prep task")
			newRecipePrepTask := fakes.BuildFakeRecipePrepTask()
			newRecipePrepTask.ID = createdRecipePrepTask.ID
			newRecipePrepTask.BelongsToRecipe = createdRecipe.ID
			newRecipePrepTask.TaskSteps = createdRecipePrepTask.TaskSteps
			createdRecipePrepTask.Update(convertRecipePrepTaskToRecipePrepTaskUpdateInput(newRecipePrepTask))
			require.NoError(t, testClients.user.UpdateRecipePrepTask(ctx, createdRecipePrepTask))

			t.Log("fetching changed recipe prep task")
			actual, err = testClients.user.GetRecipePrepTask(ctx, createdRecipe.ID, createdRecipePrepTask.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe prep task equality
			checkRecipePrepTaskEquality(t, newRecipePrepTask, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up recipe prep task")
			assert.NoError(t, testClients.user.ArchiveRecipePrepTask(ctx, createdRecipe.ID, createdRecipePrepTask.ID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipePrepTasks_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var createdRecipeStep *types.RecipeStep
			for _, step := range createdRecipe.Steps {
				createdRecipeStep = step
				break
			}
			require.NotNil(t, createdRecipeStep)

			t.Log("creating recipe prep tasks")
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

				exampleInput := fakes.BuildFakeRecipePrepTaskCreationRequestInputFromRecipePrepTask(exampleRecipePrepTask)

				createdRecipePrepTask, err := testClients.user.CreateRecipePrepTask(ctx, exampleInput)
				requireNotNilAndNoProblems(t, createdRecipePrepTask, err)

				exampleRecipePrepTaskInput := fakes.BuildFakeRecipePrepTaskCreationRequestInputFromRecipePrepTask(exampleRecipePrepTask)

				createdRecipePrepTask, createdRecipePrepTaskErr := testClients.user.CreateRecipePrepTask(ctx, exampleRecipePrepTaskInput)
				require.NoError(t, createdRecipePrepTaskErr)
				t.Logf("recipe prep task %q created", createdRecipePrepTask.ID)

				for j := range createdRecipePrepTask.TaskSteps {
					exampleRecipePrepTask.TaskSteps[j].ID = createdRecipePrepTask.TaskSteps[j].ID
					exampleRecipePrepTask.TaskSteps[j].BelongsToRecipePrepTask = createdRecipePrepTask.ID
				}

				checkRecipePrepTaskEquality(t, exampleRecipePrepTask, createdRecipePrepTask)

				createdRecipePrepTask, createdRecipePrepTaskErr = testClients.user.GetRecipePrepTask(ctx, createdRecipe.ID, createdRecipePrepTask.ID)
				requireNotNilAndNoProblems(t, createdRecipePrepTask, createdRecipePrepTaskErr)
				require.Equal(t, createdRecipe.ID, createdRecipePrepTask.BelongsToRecipe)

				expected = append(expected, createdRecipePrepTask)
			}

			// assert recipe prep task list equality
			actual, err := testClients.user.GetRecipePrepTasks(ctx, createdRecipe.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected %d to be <= %d",
				len(expected),
				len(actual),
			)

			t.Log("cleaning up")
			for _, createdRecipePrepTask := range expected {
				assert.NoError(t, testClients.user.ArchiveRecipePrepTask(ctx, createdRecipe.ID, createdRecipePrepTask.ID))
			}

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}