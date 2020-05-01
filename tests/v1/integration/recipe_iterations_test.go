package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkRecipeIterationEquality(t *testing.T, expected, actual *models.RecipeIteration) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.EndDifficultyRating, actual.EndDifficultyRating, "expected EndDifficultyRating for ID %d to be %v, but it was %v ", expected.ID, expected.EndDifficultyRating, actual.EndDifficultyRating)
	assert.Equal(t, expected.EndComplexityRating, actual.EndComplexityRating, "expected EndComplexityRating for ID %d to be %v, but it was %v ", expected.ID, expected.EndComplexityRating, actual.EndComplexityRating)
	assert.Equal(t, expected.EndTasteRating, actual.EndTasteRating, "expected EndTasteRating for ID %d to be %v, but it was %v ", expected.ID, expected.EndTasteRating, actual.EndTasteRating)
	assert.Equal(t, expected.EndOverallRating, actual.EndOverallRating, "expected EndOverallRating for ID %d to be %v, but it was %v ", expected.ID, expected.EndOverallRating, actual.EndOverallRating)
	assert.NotZero(t, actual.CreatedOn)
}

func TestRecipeIterations(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			// Assert recipe iteration equality.
			checkRecipeIterationEquality(t, exampleRecipeIteration, createdRecipeIteration)

			// Clean up.
			err = prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID)
			checkValueAndError(t, actual, err)
			checkRecipeIterationEquality(t, exampleRecipeIteration, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("should fail to create for nonexistent recipe", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = nonexistentID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)

			assert.Nil(t, createdRecipeIteration)
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

			// Create recipe iterations.
			var expected []*models.RecipeIteration
			for i := 0; i < 5; i++ {
				// Create recipe iteration.
				exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
				exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
				exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
				createdRecipeIteration, recipeIterationCreationErr := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
				checkValueAndError(t, createdRecipeIteration, recipeIterationCreationErr)

				expected = append(expected, createdRecipeIteration)
			}

			// Assert recipe iteration list equality.
			actual, err := prixfixeClient.GetRecipeIterations(ctx, createdRecipe.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeIterations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeIterations),
			)

			// Clean up.
			for _, createdRecipeIteration := range actual.RecipeIterations {
				err = prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID)
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

			// Attempt to fetch nonexistent recipe iteration.
			actual, err := prixfixeClient.RecipeIterationExists(ctx, createdRecipe.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return true with no error when the relevant recipe iteration exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			// Fetch recipe iteration.
			actual, err := prixfixeClient.RecipeIterationExists(ctx, createdRecipe.ID, createdRecipeIteration.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

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

			// Attempt to fetch nonexistent recipe iteration.
			_, err = prixfixeClient.GetRecipeIteration(ctx, createdRecipe.ID, nonexistentID)
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

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			// Fetch recipe iteration.
			actual, err := prixfixeClient.GetRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe iteration equality.
			checkRecipeIterationEquality(t, exampleRecipeIteration, actual)

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

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

			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIteration.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateRecipeIteration(ctx, exampleRecipeIteration))

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

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			// Change recipe iteration.
			createdRecipeIteration.Update(exampleRecipeIteration.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeIteration(ctx, createdRecipeIteration)
			assert.NoError(t, err)

			// Fetch recipe iteration.
			actual, err := prixfixeClient.GetRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe iteration equality.
			checkRecipeIterationEquality(t, exampleRecipeIteration, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

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

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			// Change recipe iteration.
			createdRecipeIteration.Update(exampleRecipeIteration.ToUpdateInput())
			createdRecipeIteration.BelongsToRecipe = nonexistentID
			err = prixfixeClient.UpdateRecipeIteration(ctx, createdRecipeIteration)
			assert.Error(t, err)

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

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

			assert.Error(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, nonexistentID))

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

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

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

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeIteration(ctx, nonexistentID, createdRecipeIteration.ID))

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})
}
