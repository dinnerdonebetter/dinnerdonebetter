package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkRecipeStepProductEquality(t *testing.T, expected, actual *models.RecipeStepProduct) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for ID %d to be %v, but it was %v ", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.NotZero(t, actual.CreatedOn)
}

func TestRecipeStepProducts(test *testing.T) {
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

			// Create recipe step product.
			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := prixfixeClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			checkValueAndError(t, createdRecipeStepProduct, err)

			// Assert recipe step product equality.
			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

			// Clean up.
			err = prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
			checkValueAndError(t, actual, err)
			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, actual)
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

			// Create recipe step product.
			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = nonexistentID
			exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := prixfixeClient.CreateRecipeStepProduct(ctx, nonexistentID, exampleRecipeStepProductInput)

			assert.Nil(t, createdRecipeStepProduct)
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

			// Create recipe step product.
			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = nonexistentID
			exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := prixfixeClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)

			assert.Nil(t, createdRecipeStepProduct)
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

			// Create recipe step products.
			var expected []*models.RecipeStepProduct
			for i := 0; i < 5; i++ {
				// Create recipe step product.
				exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
				exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
				exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
				createdRecipeStepProduct, recipeStepProductCreationErr := prixfixeClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
				checkValueAndError(t, createdRecipeStepProduct, recipeStepProductCreationErr)

				expected = append(expected, createdRecipeStepProduct)
			}

			// Assert recipe step product list equality.
			actual, err := prixfixeClient.GetRecipeStepProducts(ctx, createdRecipe.ID, createdRecipeStep.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepProducts),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepProducts),
			)

			// Clean up.
			for _, createdRecipeStepProduct := range actual.RecipeStepProducts {
				err = prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
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

			// Attempt to fetch nonexistent recipe step product.
			actual, err := prixfixeClient.RecipeStepProductExists(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return true with no error when the relevant recipe step product exists", func(t *testing.T) {
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

			// Create recipe step product.
			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := prixfixeClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			checkValueAndError(t, createdRecipeStepProduct, err)

			// Fetch recipe step product.
			actual, err := prixfixeClient.RecipeStepProductExists(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up recipe step product.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

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

			// Attempt to fetch nonexistent recipe step product.
			_, err = prixfixeClient.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
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

			// Create recipe step product.
			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := prixfixeClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			checkValueAndError(t, createdRecipeStepProduct, err)

			// Fetch recipe step product.
			actual, err := prixfixeClient.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step product equality.
			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, actual)

			// Clean up recipe step product.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

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

			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProduct.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProduct))

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

			// Create recipe step product.
			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := prixfixeClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			checkValueAndError(t, createdRecipeStepProduct, err)

			// Change recipe step product.
			createdRecipeStepProduct.Update(exampleRecipeStepProduct.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepProduct)
			assert.NoError(t, err)

			// Fetch recipe step product.
			actual, err := prixfixeClient.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe step product equality.
			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up recipe step product.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

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

			// Create recipe step product.
			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := prixfixeClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			checkValueAndError(t, createdRecipeStepProduct, err)

			// Change recipe step product.
			createdRecipeStepProduct.Update(exampleRecipeStepProduct.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeStepProduct(ctx, nonexistentID, createdRecipeStepProduct)
			assert.Error(t, err)

			// Clean up recipe step product.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

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

			// Create recipe step product.
			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := prixfixeClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			checkValueAndError(t, createdRecipeStepProduct, err)

			// Change recipe step product.
			createdRecipeStepProduct.Update(exampleRecipeStepProduct.ToUpdateInput())
			createdRecipeStepProduct.BelongsToRecipeStep = nonexistentID
			err = prixfixeClient.UpdateRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepProduct)
			assert.Error(t, err)

			// Clean up recipe step product.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

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

			assert.Error(t, prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID))

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

			// Create recipe step product.
			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := prixfixeClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			checkValueAndError(t, createdRecipeStepProduct, err)

			// Clean up recipe step product.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

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

			// Create recipe step product.
			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := prixfixeClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			checkValueAndError(t, createdRecipeStepProduct, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeStepProduct(ctx, nonexistentID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

			// Clean up recipe step product.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

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

			// Create recipe step product.
			exampleRecipeStepProduct := fakemodels.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := prixfixeClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			checkValueAndError(t, createdRecipeStepProduct, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, nonexistentID, createdRecipeStepProduct.ID))

			// Clean up recipe step product.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

			// Clean up recipe step.
			assert.NoError(t, prixfixeClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})
}
