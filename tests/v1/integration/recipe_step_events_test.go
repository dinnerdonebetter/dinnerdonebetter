package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkRecipeStepEventEquality(t *testing.T, expected, actual *models.RecipeStepEvent) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.EventType, actual.EventType, "expected EventType for ID %d to be %v, but it was %v ", expected.ID, expected.EventType, actual.EventType)
	assert.Equal(t, expected.Done, actual.Done, "expected Done for ID %d to be %v, but it was %v ", expected.ID, expected.Done, actual.Done)
	assert.Equal(t, expected.RecipeIterationID, actual.RecipeIterationID, "expected RecipeIterationID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeIterationID, actual.RecipeIterationID)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.NotZero(t, actual.CreatedOn)
}

func TestRecipeStepEvents(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Create recipe step event.
			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
			createdRecipeStepEvent, err := prixfixeClient.CreateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEventInput)
			checkValueAndError(t, createdRecipeStepEvent, err)

			// Assert recipe step event equality.
			checkRecipeStepEventEquality(t, exampleRecipeStepEvent, createdRecipeStepEvent)

			// Clean up.
			err = prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID)
			checkValueAndError(t, actual, err)
			checkRecipeStepEventEquality(t, exampleRecipeStepEvent, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("should fail to create for nonexistent recipe", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe step event.
			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = nonexistentID
			exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
			createdRecipeStepEvent, err := prixfixeClient.CreateRecipeStepEvent(ctx, nonexistentID, exampleRecipeStepEventInput)

			assert.Nil(t, createdRecipeStepEvent)
			assert.Error(t, err)
		})

		T.Run("should fail to create for nonexistent recipe step", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe step event.
			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = nonexistentID
			exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
			createdRecipeStepEvent, err := prixfixeClient.CreateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEventInput)

			assert.Nil(t, createdRecipeStepEvent)
			assert.Error(t, err)

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
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

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Create recipe step events.
			var expected []*models.RecipeStepEvent
			for i := 0; i < 5; i++ {
				// Create recipe step event.
				exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
				exampleRecipeStepEvent.BelongsToRecipeStep = createdRecipeStep.ID
				exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
				createdRecipeStepEvent, recipeStepEventCreationErr := prixfixeClient.CreateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEventInput)
				checkValueAndError(t, createdRecipeStepEvent, recipeStepEventCreationErr)

				expected = append(expected, createdRecipeStepEvent)
			}

			// Assert recipe step event list equality.
			actual, err := prixfixeClient.GetRecipeStepEvents(ctx, createdRecipe.ID, createdRecipeStep.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepEvents),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepEvents),
			)

			// Clean up.
			for _, createdRecipeStepEvent := range actual.RecipeStepEvents {
				err = prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID)
				assert.NoError(t, err)
			}

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

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

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Attempt to fetch nonexistent recipe step event.
			actual, err := prixfixeClient.RecipeStepEventExists(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return true with no error when the relevant recipe step event exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Create recipe step event.
			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
			createdRecipeStepEvent, err := prixfixeClient.CreateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEventInput)
			checkValueAndError(t, createdRecipeStepEvent, err)

			// Fetch recipe step event.
			actual, err := prixfixeClient.RecipeStepEventExists(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up recipe step event.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

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

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Attempt to fetch nonexistent recipe step event.
			_, err = prixfixeClient.GetRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
			assert.Error(t, err)

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

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

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Create recipe step event.
			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
			createdRecipeStepEvent, err := prixfixeClient.CreateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEventInput)
			checkValueAndError(t, createdRecipeStepEvent, err)

			// Fetch recipe step event.
			actual, err := prixfixeClient.GetRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step event equality.
			checkRecipeStepEventEquality(t, exampleRecipeStepEvent, actual)

			// Clean up recipe step event.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

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

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepEvent.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEvent))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

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

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Create recipe step event.
			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
			createdRecipeStepEvent, err := prixfixeClient.CreateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEventInput)
			checkValueAndError(t, createdRecipeStepEvent, err)

			// Change recipe step event.
			createdRecipeStepEvent.Update(exampleRecipeStepEvent.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStepEvent)
			assert.NoError(t, err)

			// Fetch recipe step event.
			actual, err := prixfixeClient.GetRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step event equality.
			checkRecipeStepEventEquality(t, exampleRecipeStepEvent, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up recipe step event.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

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

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Create recipe step event.
			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
			createdRecipeStepEvent, err := prixfixeClient.CreateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEventInput)
			checkValueAndError(t, createdRecipeStepEvent, err)

			// Change recipe step event.
			createdRecipeStepEvent.Update(exampleRecipeStepEvent.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeStepEvent(ctx, nonexistentID, createdRecipeStepEvent)
			assert.Error(t, err)

			// Clean up recipe step event.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return an error when trying to update something that belongs to a recipe step that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Create recipe step event.
			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
			createdRecipeStepEvent, err := prixfixeClient.CreateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEventInput)
			checkValueAndError(t, createdRecipeStepEvent, err)

			// Change recipe step event.
			createdRecipeStepEvent.Update(exampleRecipeStepEvent.ToUpdateInput())
			createdRecipeStepEvent.BelongsToRecipeStep = nonexistentID
			err = prixfixeClient.UpdateRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStepEvent)
			assert.Error(t, err)

			// Clean up recipe step event.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

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

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

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

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Create recipe step event.
			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
			createdRecipeStepEvent, err := prixfixeClient.CreateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEventInput)
			checkValueAndError(t, createdRecipeStepEvent, err)

			// Clean up recipe step event.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

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

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Create recipe step event.
			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
			createdRecipeStepEvent, err := prixfixeClient.CreateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEventInput)
			checkValueAndError(t, createdRecipeStepEvent, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeStepEvent(ctx, nonexistentID, createdRecipeStep.ID, createdRecipeStepEvent.ID))

			// Clean up recipe step event.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("returns error when trying to archive post belonging to nonexistent recipe step", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := prixfixeClient.CreateRecipeStep(ctx, exampleRecipeStepInput)
			checkValueAndError(t, createdRecipeStep, err)

			// Create recipe step event.
			exampleRecipeStepEvent := fakemodels.BuildFakeRecipeStepEvent()
			exampleRecipeStepEvent.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(exampleRecipeStepEvent)
			createdRecipeStepEvent, err := prixfixeClient.CreateRecipeStepEvent(ctx, createdRecipe.ID, exampleRecipeStepEventInput)
			checkValueAndError(t, createdRecipeStepEvent, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, nonexistentID, createdRecipeStepEvent.ID))

			// Clean up recipe step event.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepEvent(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepEvent.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})
}
