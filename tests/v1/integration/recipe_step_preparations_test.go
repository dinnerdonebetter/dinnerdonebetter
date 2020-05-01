package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkRecipeStepPreparationEquality(t *testing.T, expected, actual *models.RecipeStepPreparation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.ValidPreparationID, actual.ValidPreparationID, "expected ValidPreparationID for ID %d to be %v, but it was %v ", expected.ID, expected.ValidPreparationID, actual.ValidPreparationID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for ID %d to be %v, but it was %v ", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

func TestRecipeStepPreparations(test *testing.T) {
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

			// Create recipe step preparation.
			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
			createdRecipeStepPreparation, err := prixfixeClient.CreateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparationInput)
			checkValueAndError(t, createdRecipeStepPreparation, err)

			// Assert recipe step preparation equality.
			checkRecipeStepPreparationEquality(t, exampleRecipeStepPreparation, createdRecipeStepPreparation)

			// Clean up.
			err = prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID)
			checkValueAndError(t, actual, err)
			checkRecipeStepPreparationEquality(t, exampleRecipeStepPreparation, actual)
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

			// Create recipe step preparation.
			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = nonexistentID
			exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
			createdRecipeStepPreparation, err := prixfixeClient.CreateRecipeStepPreparation(ctx, nonexistentID, exampleRecipeStepPreparationInput)

			assert.Nil(t, createdRecipeStepPreparation)
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

			// Create recipe step preparation.
			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = nonexistentID
			exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
			createdRecipeStepPreparation, err := prixfixeClient.CreateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparationInput)

			assert.Nil(t, createdRecipeStepPreparation)
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

			// Create recipe step preparations.
			var expected []*models.RecipeStepPreparation
			for i := 0; i < 5; i++ {
				// Create recipe step preparation.
				exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
				exampleRecipeStepPreparation.BelongsToRecipeStep = createdRecipeStep.ID
				exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
				createdRecipeStepPreparation, recipeStepPreparationCreationErr := prixfixeClient.CreateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparationInput)
				checkValueAndError(t, createdRecipeStepPreparation, recipeStepPreparationCreationErr)

				expected = append(expected, createdRecipeStepPreparation)
			}

			// Assert recipe step preparation list equality.
			actual, err := prixfixeClient.GetRecipeStepPreparations(ctx, createdRecipe.ID, createdRecipeStep.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepPreparations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepPreparations),
			)

			// Clean up.
			for _, createdRecipeStepPreparation := range actual.RecipeStepPreparations {
				err = prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID)
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

			// Attempt to fetch nonexistent recipe step preparation.
			actual, err := prixfixeClient.RecipeStepPreparationExists(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return true with no error when the relevant recipe step preparation exists", func(t *testing.T) {
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

			// Create recipe step preparation.
			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
			createdRecipeStepPreparation, err := prixfixeClient.CreateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparationInput)
			checkValueAndError(t, createdRecipeStepPreparation, err)

			// Fetch recipe step preparation.
			actual, err := prixfixeClient.RecipeStepPreparationExists(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up recipe step preparation.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID))

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

			// Attempt to fetch nonexistent recipe step preparation.
			_, err = prixfixeClient.GetRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
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

			// Create recipe step preparation.
			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
			createdRecipeStepPreparation, err := prixfixeClient.CreateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparationInput)
			checkValueAndError(t, createdRecipeStepPreparation, err)

			// Fetch recipe step preparation.
			actual, err := prixfixeClient.GetRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step preparation equality.
			checkRecipeStepPreparationEquality(t, exampleRecipeStepPreparation, actual)

			// Clean up recipe step preparation.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID))

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

			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepPreparation.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparation))

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

			// Create recipe step preparation.
			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
			createdRecipeStepPreparation, err := prixfixeClient.CreateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparationInput)
			checkValueAndError(t, createdRecipeStepPreparation, err)

			// Change recipe step preparation.
			createdRecipeStepPreparation.Update(exampleRecipeStepPreparation.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStepPreparation)
			assert.NoError(t, err)

			// Fetch recipe step preparation.
			actual, err := prixfixeClient.GetRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step preparation equality.
			checkRecipeStepPreparationEquality(t, exampleRecipeStepPreparation, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up recipe step preparation.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID))

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

			// Create recipe step preparation.
			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
			createdRecipeStepPreparation, err := prixfixeClient.CreateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparationInput)
			checkValueAndError(t, createdRecipeStepPreparation, err)

			// Change recipe step preparation.
			createdRecipeStepPreparation.Update(exampleRecipeStepPreparation.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeStepPreparation(ctx, nonexistentID, createdRecipeStepPreparation)
			assert.Error(t, err)

			// Clean up recipe step preparation.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID))

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

			// Create recipe step preparation.
			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
			createdRecipeStepPreparation, err := prixfixeClient.CreateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparationInput)
			checkValueAndError(t, createdRecipeStepPreparation, err)

			// Change recipe step preparation.
			createdRecipeStepPreparation.Update(exampleRecipeStepPreparation.ToUpdateInput())
			createdRecipeStepPreparation.BelongsToRecipeStep = nonexistentID
			err = prixfixeClient.UpdateRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStepPreparation)
			assert.Error(t, err)

			// Clean up recipe step preparation.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID))

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

			assert.Error(t, prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID))

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

			// Create recipe step preparation.
			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
			createdRecipeStepPreparation, err := prixfixeClient.CreateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparationInput)
			checkValueAndError(t, createdRecipeStepPreparation, err)

			// Clean up recipe step preparation.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID))

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

			// Create recipe step preparation.
			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
			createdRecipeStepPreparation, err := prixfixeClient.CreateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparationInput)
			checkValueAndError(t, createdRecipeStepPreparation, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeStepPreparation(ctx, nonexistentID, createdRecipeStep.ID, createdRecipeStepPreparation.ID))

			// Clean up recipe step preparation.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID))

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

			// Create recipe step preparation.
			exampleRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparation()
			exampleRecipeStepPreparation.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInputFromRecipeStepPreparation(exampleRecipeStepPreparation)
			createdRecipeStepPreparation, err := prixfixeClient.CreateRecipeStepPreparation(ctx, createdRecipe.ID, exampleRecipeStepPreparationInput)
			checkValueAndError(t, createdRecipeStepPreparation, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, nonexistentID, createdRecipeStepPreparation.ID))

			// Clean up recipe step preparation.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepPreparation(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepPreparation.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})
}
