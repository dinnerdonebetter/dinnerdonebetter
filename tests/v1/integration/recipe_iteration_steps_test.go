package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkRecipeIterationStepEquality(t *testing.T, expected, actual *models.RecipeIterationStep) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, *expected.StartedOn, *actual.StartedOn, "expected StartedOn to be %v, but it was %v ", expected.StartedOn, actual.StartedOn)
	assert.Equal(t, *expected.EndedOn, *actual.EndedOn, "expected EndedOn to be %v, but it was %v ", expected.EndedOn, actual.EndedOn)
	assert.Equal(t, expected.State, actual.State, "expected State for ID %d to be %v, but it was %v ", expected.ID, expected.State, actual.State)
	assert.NotZero(t, actual.CreatedOn)
}

func TestRecipeIterationSteps(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe iteration step.
			exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
			exampleRecipeIterationStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationStepInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)
			createdRecipeIterationStep, err := prixfixeClient.CreateRecipeIterationStep(ctx, exampleRecipeIterationStepInput)
			checkValueAndError(t, createdRecipeIterationStep, err)

			// Assert recipe iteration step equality.
			checkRecipeIterationStepEquality(t, exampleRecipeIterationStep, createdRecipeIterationStep)

			// Clean up.
			err = prixfixeClient.ArchiveRecipeIterationStep(ctx, createdRecipe.ID, createdRecipeIterationStep.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetRecipeIterationStep(ctx, createdRecipe.ID, createdRecipeIterationStep.ID)
			checkValueAndError(t, actual, err)
			checkRecipeIterationStepEquality(t, exampleRecipeIterationStep, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("should fail to create for nonexistent recipe", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe iteration step.
			exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
			exampleRecipeIterationStep.BelongsToRecipe = nonexistentID
			exampleRecipeIterationStepInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)
			createdRecipeIterationStep, err := prixfixeClient.CreateRecipeIterationStep(ctx, exampleRecipeIterationStepInput)

			assert.Nil(t, createdRecipeIterationStep)
			assert.Error(t, err)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe iteration steps.
			var expected []*models.RecipeIterationStep
			for i := 0; i < 5; i++ {
				// Create recipe iteration step.
				exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
				exampleRecipeIterationStep.BelongsToRecipe = createdRecipe.ID
				exampleRecipeIterationStepInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)
				createdRecipeIterationStep, recipeIterationStepCreationErr := prixfixeClient.CreateRecipeIterationStep(ctx, exampleRecipeIterationStepInput)
				checkValueAndError(t, createdRecipeIterationStep, recipeIterationStepCreationErr)

				expected = append(expected, createdRecipeIterationStep)
			}

			// Assert recipe iteration step list equality.
			actual, err := prixfixeClient.GetRecipeIterationSteps(ctx, createdRecipe.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeIterationSteps),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeIterationSteps),
			)

			// Clean up.
			for _, createdRecipeIterationStep := range actual.RecipeIterationSteps {
				err = prixfixeClient.ArchiveRecipeIterationStep(ctx, createdRecipe.ID, createdRecipeIterationStep.ID)
				assert.NoError(t, err)
			}

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Attempt to fetch nonexistent recipe iteration step.
			actual, err := prixfixeClient.RecipeIterationStepExists(ctx, createdRecipe.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return true with no error when the relevant recipe iteration step exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe iteration step.
			exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
			exampleRecipeIterationStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationStepInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)
			createdRecipeIterationStep, err := prixfixeClient.CreateRecipeIterationStep(ctx, exampleRecipeIterationStepInput)
			checkValueAndError(t, createdRecipeIterationStep, err)

			// Fetch recipe iteration step.
			actual, err := prixfixeClient.RecipeIterationStepExists(ctx, createdRecipe.ID, createdRecipeIterationStep.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up recipe iteration step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIterationStep(ctx, createdRecipe.ID, createdRecipeIterationStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Attempt to fetch nonexistent recipe iteration step.
			_, err = prixfixeClient.GetRecipeIterationStep(ctx, createdRecipe.ID, nonexistentID)
			assert.Error(t, err)

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe iteration step.
			exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
			exampleRecipeIterationStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationStepInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)
			createdRecipeIterationStep, err := prixfixeClient.CreateRecipeIterationStep(ctx, exampleRecipeIterationStepInput)
			checkValueAndError(t, createdRecipeIterationStep, err)

			// Fetch recipe iteration step.
			actual, err := prixfixeClient.GetRecipeIterationStep(ctx, createdRecipe.ID, createdRecipeIterationStep.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe iteration step equality.
			checkRecipeIterationStepEquality(t, exampleRecipeIterationStep, actual)

			// Clean up recipe iteration step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIterationStep(ctx, createdRecipe.ID, createdRecipeIterationStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
			exampleRecipeIterationStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationStep.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateRecipeIterationStep(ctx, exampleRecipeIterationStep))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe iteration step.
			exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
			exampleRecipeIterationStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationStepInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)
			createdRecipeIterationStep, err := prixfixeClient.CreateRecipeIterationStep(ctx, exampleRecipeIterationStepInput)
			checkValueAndError(t, createdRecipeIterationStep, err)

			// Change recipe iteration step.
			createdRecipeIterationStep.Update(exampleRecipeIterationStep.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeIterationStep(ctx, createdRecipeIterationStep)
			assert.NoError(t, err)

			// Fetch recipe iteration step.
			actual, err := prixfixeClient.GetRecipeIterationStep(ctx, createdRecipe.ID, createdRecipeIterationStep.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe iteration step equality.
			checkRecipeIterationStepEquality(t, exampleRecipeIterationStep, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up recipe iteration step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIterationStep(ctx, createdRecipe.ID, createdRecipeIterationStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return an error when trying to update something that belongs to a recipe that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe iteration step.
			exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
			exampleRecipeIterationStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationStepInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)
			createdRecipeIterationStep, err := prixfixeClient.CreateRecipeIterationStep(ctx, exampleRecipeIterationStepInput)
			checkValueAndError(t, createdRecipeIterationStep, err)

			// Change recipe iteration step.
			createdRecipeIterationStep.Update(exampleRecipeIterationStep.ToUpdateInput())
			createdRecipeIterationStep.BelongsToRecipe = nonexistentID
			err = prixfixeClient.UpdateRecipeIterationStep(ctx, createdRecipeIterationStep)
			assert.Error(t, err)

			// Clean up recipe iteration step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIterationStep(ctx, createdRecipe.ID, createdRecipeIterationStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeIterationStep(ctx, createdRecipe.ID, nonexistentID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe iteration step.
			exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
			exampleRecipeIterationStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationStepInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)
			createdRecipeIterationStep, err := prixfixeClient.CreateRecipeIterationStep(ctx, exampleRecipeIterationStepInput)
			checkValueAndError(t, createdRecipeIterationStep, err)

			// Clean up recipe iteration step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIterationStep(ctx, createdRecipe.ID, createdRecipeIterationStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("returns error when trying to archive post belonging to nonexistent recipe", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe iteration step.
			exampleRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStep()
			exampleRecipeIterationStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationStepInput := fakemodels.BuildFakeRecipeIterationStepCreationInputFromRecipeIterationStep(exampleRecipeIterationStep)
			createdRecipeIterationStep, err := prixfixeClient.CreateRecipeIterationStep(ctx, exampleRecipeIterationStepInput)
			checkValueAndError(t, createdRecipeIterationStep, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeIterationStep(ctx, nonexistentID, createdRecipeIterationStep.ID))

			// Clean up recipe iteration step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIterationStep(ctx, createdRecipe.ID, createdRecipeIterationStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})
}
