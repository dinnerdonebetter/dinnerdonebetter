package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkRecipeStepIngredientEquality(t *testing.T, expected, actual *models.RecipeStepIngredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.ValidIngredientID, actual.ValidIngredientID, "expected ValidIngredientID for ID %d to be %v, but it was %v ", expected.ID, expected.ValidIngredientID, actual.ValidIngredientID)
	assert.Equal(t, expected.IngredientNotes, actual.IngredientNotes, "expected IngredientNotes for ID %d to be %v, but it was %v ", expected.ID, expected.IngredientNotes, actual.IngredientNotes)
	assert.Equal(t, expected.QuantityType, actual.QuantityType, "expected QuantityType for ID %d to be %v, but it was %v ", expected.ID, expected.QuantityType, actual.QuantityType)
	assert.Equal(t, expected.QuantityValue, actual.QuantityValue, "expected QuantityValue for ID %d to be %v, but it was %v ", expected.ID, expected.QuantityValue, actual.QuantityValue)
	assert.Equal(t, expected.QuantityNotes, actual.QuantityNotes, "expected QuantityNotes for ID %d to be %v, but it was %v ", expected.ID, expected.QuantityNotes, actual.QuantityNotes)
	assert.Equal(t, *expected.ProductOfRecipeStepID, *actual.ProductOfRecipeStepID, "expected ProductOfRecipeStepID to be %v, but it was %v ", expected.ProductOfRecipeStepID, actual.ProductOfRecipeStepID)
	assert.NotZero(t, actual.CreatedOn)
}

func TestRecipeStepIngredients(test *testing.T) {
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

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := prixfixeClient.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			checkValueAndError(t, createdRecipeStepIngredient, err)

			// Assert recipe step ingredient equality.
			checkRecipeStepIngredientEquality(t, exampleRecipeStepIngredient, createdRecipeStepIngredient)

			// Clean up.
			err = prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
			checkValueAndError(t, actual, err)
			checkRecipeStepIngredientEquality(t, exampleRecipeStepIngredient, actual)
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

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = nonexistentID
			exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := prixfixeClient.CreateRecipeStepIngredient(ctx, nonexistentID, exampleRecipeStepIngredientInput)

			assert.Nil(t, createdRecipeStepIngredient)
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

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = nonexistentID
			exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := prixfixeClient.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)

			assert.Nil(t, createdRecipeStepIngredient)
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

			// Create recipe step ingredients.
			var expected []*models.RecipeStepIngredient
			for i := 0; i < 5; i++ {
				// Create recipe step ingredient.
				exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
				exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
				exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
				createdRecipeStepIngredient, recipeStepIngredientCreationErr := prixfixeClient.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
				checkValueAndError(t, createdRecipeStepIngredient, recipeStepIngredientCreationErr)

				expected = append(expected, createdRecipeStepIngredient)
			}

			// Assert recipe step ingredient list equality.
			actual, err := prixfixeClient.GetRecipeStepIngredients(ctx, createdRecipe.ID, createdRecipeStep.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepIngredients),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepIngredients),
			)

			// Clean up.
			for _, createdRecipeStepIngredient := range actual.RecipeStepIngredients {
				err = prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
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

			// Attempt to fetch nonexistent recipe step ingredient.
			actual, err := prixfixeClient.RecipeStepIngredientExists(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return true with no error when the relevant recipe step ingredient exists", func(t *testing.T) {
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

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := prixfixeClient.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			checkValueAndError(t, createdRecipeStepIngredient, err)

			// Fetch recipe step ingredient.
			actual, err := prixfixeClient.RecipeStepIngredientExists(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up recipe step ingredient.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

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

			// Attempt to fetch nonexistent recipe step ingredient.
			_, err = prixfixeClient.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
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

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := prixfixeClient.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			checkValueAndError(t, createdRecipeStepIngredient, err)

			// Fetch recipe step ingredient.
			actual, err := prixfixeClient.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step ingredient equality.
			checkRecipeStepIngredientEquality(t, exampleRecipeStepIngredient, actual)

			// Clean up recipe step ingredient.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

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

			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredient.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredient))

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

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := prixfixeClient.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			checkValueAndError(t, createdRecipeStepIngredient, err)

			// Change recipe step ingredient.
			createdRecipeStepIngredient.Update(exampleRecipeStepIngredient.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepIngredient)
			assert.NoError(t, err)

			// Fetch recipe step ingredient.
			actual, err := prixfixeClient.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step ingredient equality.
			checkRecipeStepIngredientEquality(t, exampleRecipeStepIngredient, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up recipe step ingredient.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

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

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := prixfixeClient.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			checkValueAndError(t, createdRecipeStepIngredient, err)

			// Change recipe step ingredient.
			createdRecipeStepIngredient.Update(exampleRecipeStepIngredient.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeStepIngredient(ctx, nonexistentID, createdRecipeStepIngredient)
			assert.Error(t, err)

			// Clean up recipe step ingredient.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

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

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := prixfixeClient.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			checkValueAndError(t, createdRecipeStepIngredient, err)

			// Change recipe step ingredient.
			createdRecipeStepIngredient.Update(exampleRecipeStepIngredient.ToUpdateInput())
			createdRecipeStepIngredient.BelongsToRecipeStep = nonexistentID
			err = prixfixeClient.UpdateRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepIngredient)
			assert.Error(t, err)

			// Clean up recipe step ingredient.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

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

			assert.Error(t, prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID))

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

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := prixfixeClient.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			checkValueAndError(t, createdRecipeStepIngredient, err)

			// Clean up recipe step ingredient.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

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

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := prixfixeClient.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			checkValueAndError(t, createdRecipeStepIngredient, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeStepIngredient(ctx, nonexistentID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

			// Clean up recipe step ingredient.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

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

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := prixfixeClient.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			checkValueAndError(t, createdRecipeStepIngredient, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, nonexistentID, createdRecipeStepIngredient.ID))

			// Clean up recipe step ingredient.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})
}
