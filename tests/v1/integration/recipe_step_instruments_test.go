package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkRecipeStepInstrumentEquality(t *testing.T, expected, actual *models.RecipeStepInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, *expected.InstrumentID, *actual.InstrumentID, "expected InstrumentID to be %v, but it was %v ", expected.InstrumentID, actual.InstrumentID)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for ID %d to be %v, but it was %v ", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

func TestRecipeStepInstruments(test *testing.T) {
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

			// Create recipe step instrument.
			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := prixfixeClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			checkValueAndError(t, createdRecipeStepInstrument, err)

			// Assert recipe step instrument equality.
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, createdRecipeStepInstrument)

			// Clean up.
			err = prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID)
			checkValueAndError(t, actual, err)
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, actual)
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

			// Create recipe step instrument.
			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = nonexistentID
			exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := prixfixeClient.CreateRecipeStepInstrument(ctx, nonexistentID, exampleRecipeStepInstrumentInput)

			assert.Nil(t, createdRecipeStepInstrument)
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

			// Create recipe step instrument.
			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = nonexistentID
			exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := prixfixeClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)

			assert.Nil(t, createdRecipeStepInstrument)
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

			// Create recipe step instruments.
			var expected []*models.RecipeStepInstrument
			for i := 0; i < 5; i++ {
				// Create recipe step instrument.
				exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
				exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
				exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
				createdRecipeStepInstrument, recipeStepInstrumentCreationErr := prixfixeClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
				checkValueAndError(t, createdRecipeStepInstrument, recipeStepInstrumentCreationErr)

				expected = append(expected, createdRecipeStepInstrument)
			}

			// Assert recipe step instrument list equality.
			actual, err := prixfixeClient.GetRecipeStepInstruments(ctx, createdRecipe.ID, createdRecipeStep.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepInstruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepInstruments),
			)

			// Clean up.
			for _, createdRecipeStepInstrument := range actual.RecipeStepInstruments {
				err = prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID)
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

			// Attempt to fetch nonexistent recipe step instrument.
			actual, err := prixfixeClient.RecipeStepInstrumentExists(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return true with no error when the relevant recipe step instrument exists", func(t *testing.T) {
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

			// Create recipe step instrument.
			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := prixfixeClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			checkValueAndError(t, createdRecipeStepInstrument, err)

			// Fetch recipe step instrument.
			actual, err := prixfixeClient.RecipeStepInstrumentExists(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up recipe step instrument.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID))

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

			// Attempt to fetch nonexistent recipe step instrument.
			_, err = prixfixeClient.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
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

			// Create recipe step instrument.
			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := prixfixeClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			checkValueAndError(t, createdRecipeStepInstrument, err)

			// Fetch recipe step instrument.
			actual, err := prixfixeClient.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step instrument equality.
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, actual)

			// Clean up recipe step instrument.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID))

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

			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrument.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrument))

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

			// Create recipe step instrument.
			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := prixfixeClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			checkValueAndError(t, createdRecipeStepInstrument, err)

			// Change recipe step instrument.
			createdRecipeStepInstrument.Update(exampleRecipeStepInstrument.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument)
			assert.NoError(t, err)

			// Fetch recipe step instrument.
			actual, err := prixfixeClient.GetRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step instrument equality.
			checkRecipeStepInstrumentEquality(t, exampleRecipeStepInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up recipe step instrument.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID))

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

			// Create recipe step instrument.
			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := prixfixeClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			checkValueAndError(t, createdRecipeStepInstrument, err)

			// Change recipe step instrument.
			createdRecipeStepInstrument.Update(exampleRecipeStepInstrument.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeStepInstrument(ctx, nonexistentID, createdRecipeStepInstrument)
			assert.Error(t, err)

			// Clean up recipe step instrument.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID))

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

			// Create recipe step instrument.
			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := prixfixeClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			checkValueAndError(t, createdRecipeStepInstrument, err)

			// Change recipe step instrument.
			createdRecipeStepInstrument.Update(exampleRecipeStepInstrument.ToUpdateInput())
			createdRecipeStepInstrument.BelongsToRecipeStep = nonexistentID
			err = prixfixeClient.UpdateRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStepInstrument)
			assert.Error(t, err)

			// Clean up recipe step instrument.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID))

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

			assert.Error(t, prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID))

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

			// Create recipe step instrument.
			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := prixfixeClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			checkValueAndError(t, createdRecipeStepInstrument, err)

			// Clean up recipe step instrument.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID))

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

			// Create recipe step instrument.
			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := prixfixeClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			checkValueAndError(t, createdRecipeStepInstrument, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeStepInstrument(ctx, nonexistentID, createdRecipeStep.ID, createdRecipeStepInstrument.ID))

			// Clean up recipe step instrument.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID))

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

			// Create recipe step instrument.
			exampleRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrument()
			exampleRecipeStepInstrument.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(exampleRecipeStepInstrument)
			createdRecipeStepInstrument, err := prixfixeClient.CreateRecipeStepInstrument(ctx, createdRecipe.ID, exampleRecipeStepInstrumentInput)
			checkValueAndError(t, createdRecipeStepInstrument, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, nonexistentID, createdRecipeStepInstrument.ID))

			// Clean up recipe step instrument.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepInstrument(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepInstrument.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})
}
