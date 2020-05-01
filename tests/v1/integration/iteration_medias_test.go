package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkIterationMediaEquality(t *testing.T, expected, actual *models.IterationMedia) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Source, actual.Source, "expected Source for ID %d to be %v, but it was %v ", expected.ID, expected.Source, actual.Source)
	assert.Equal(t, expected.Mimetype, actual.Mimetype, "expected Mimetype for ID %d to be %v, but it was %v ", expected.ID, expected.Mimetype, actual.Mimetype)
	assert.NotZero(t, actual.CreatedOn)
}

func TestIterationMedias(test *testing.T) {
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

			// Create iteration media.
			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = createdRecipeIteration.ID
			exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
			createdIterationMedia, err := prixfixeClient.CreateIterationMedia(ctx, createdRecipe.ID, exampleIterationMediaInput)
			checkValueAndError(t, createdIterationMedia, err)

			// Assert iteration media equality.
			checkIterationMediaEquality(t, exampleIterationMedia, createdIterationMedia)

			// Clean up.
			err = prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID)
			checkValueAndError(t, actual, err)
			checkIterationMediaEquality(t, exampleIterationMedia, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("should fail to create for nonexistent recipe", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create iteration media.
			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = nonexistentID
			exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
			createdIterationMedia, err := prixfixeClient.CreateIterationMedia(ctx, nonexistentID, exampleIterationMediaInput)

			assert.Nil(t, createdIterationMedia)
			assert.Error(t, err)
		})

		T.Run("should fail to create for nonexistent recipe iteration", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create iteration media.
			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = nonexistentID
			exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
			createdIterationMedia, err := prixfixeClient.CreateIterationMedia(ctx, createdRecipe.ID, exampleIterationMediaInput)

			assert.Nil(t, createdIterationMedia)
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

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			// Create iteration medias.
			var expected []*models.IterationMedia
			for i := 0; i < 5; i++ {
				// Create iteration media.
				exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
				exampleIterationMedia.BelongsToRecipeIteration = createdRecipeIteration.ID
				exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
				createdIterationMedia, iterationMediaCreationErr := prixfixeClient.CreateIterationMedia(ctx, createdRecipe.ID, exampleIterationMediaInput)
				checkValueAndError(t, createdIterationMedia, iterationMediaCreationErr)

				expected = append(expected, createdIterationMedia)
			}

			// Assert iteration media list equality.
			actual, err := prixfixeClient.GetIterationMedias(ctx, createdRecipe.ID, createdRecipeIteration.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.IterationMedias),
				"expected %d to be <= %d",
				len(expected),
				len(actual.IterationMedias),
			)

			// Clean up.
			for _, createdIterationMedia := range actual.IterationMedias {
				err = prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID)
				assert.NoError(t, err)
			}

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

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

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			// Attempt to fetch nonexistent iteration media.
			actual, err := prixfixeClient.IterationMediaExists(ctx, createdRecipe.ID, createdRecipeIteration.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return true with no error when the relevant iteration media exists", func(t *testing.T) {
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

			// Create iteration media.
			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = createdRecipeIteration.ID
			exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
			createdIterationMedia, err := prixfixeClient.CreateIterationMedia(ctx, createdRecipe.ID, exampleIterationMediaInput)
			checkValueAndError(t, createdIterationMedia, err)

			// Fetch iteration media.
			actual, err := prixfixeClient.IterationMediaExists(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up iteration media.
			assert.NoError(t, prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID))

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

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			// Attempt to fetch nonexistent iteration media.
			_, err = prixfixeClient.GetIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, nonexistentID)
			assert.Error(t, err)

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

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

			// Create iteration media.
			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = createdRecipeIteration.ID
			exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
			createdIterationMedia, err := prixfixeClient.CreateIterationMedia(ctx, createdRecipe.ID, exampleIterationMediaInput)
			checkValueAndError(t, createdIterationMedia, err)

			// Fetch iteration media.
			actual, err := prixfixeClient.GetIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID)
			checkValueAndError(t, actual, err)

			// Assert iteration media equality.
			checkIterationMediaEquality(t, exampleIterationMedia, actual)

			// Clean up iteration media.
			assert.NoError(t, prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID))

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

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = createdRecipeIteration.ID
			exampleIterationMedia.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateIterationMedia(ctx, createdRecipe.ID, exampleIterationMedia))

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

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

			// Create iteration media.
			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = createdRecipeIteration.ID
			exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
			createdIterationMedia, err := prixfixeClient.CreateIterationMedia(ctx, createdRecipe.ID, exampleIterationMediaInput)
			checkValueAndError(t, createdIterationMedia, err)

			// Change iteration media.
			createdIterationMedia.Update(exampleIterationMedia.ToUpdateInput())
			err = prixfixeClient.UpdateIterationMedia(ctx, createdRecipe.ID, createdIterationMedia)
			assert.NoError(t, err)

			// Fetch iteration media.
			actual, err := prixfixeClient.GetIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID)
			checkValueAndError(t, actual, err)

			// Assert iteration media equality.
			checkIterationMediaEquality(t, exampleIterationMedia, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up iteration media.
			assert.NoError(t, prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID))

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

			// Create iteration media.
			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = createdRecipeIteration.ID
			exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
			createdIterationMedia, err := prixfixeClient.CreateIterationMedia(ctx, createdRecipe.ID, exampleIterationMediaInput)
			checkValueAndError(t, createdIterationMedia, err)

			// Change iteration media.
			createdIterationMedia.Update(exampleIterationMedia.ToUpdateInput())
			err = prixfixeClient.UpdateIterationMedia(ctx, nonexistentID, createdIterationMedia)
			assert.Error(t, err)

			// Clean up iteration media.
			assert.NoError(t, prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID))

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return an error when trying to update something that belongs to a recipe iteration that does not exist", func(t *testing.T) {
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

			// Create iteration media.
			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = createdRecipeIteration.ID
			exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
			createdIterationMedia, err := prixfixeClient.CreateIterationMedia(ctx, createdRecipe.ID, exampleIterationMediaInput)
			checkValueAndError(t, createdIterationMedia, err)

			// Change iteration media.
			createdIterationMedia.Update(exampleIterationMedia.ToUpdateInput())
			createdIterationMedia.BelongsToRecipeIteration = nonexistentID
			err = prixfixeClient.UpdateIterationMedia(ctx, createdRecipe.ID, createdIterationMedia)
			assert.Error(t, err)

			// Clean up iteration media.
			assert.NoError(t, prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID))

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

			// Create recipe iteration.
			exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
			exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
			exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
			createdRecipeIteration, err := prixfixeClient.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
			checkValueAndError(t, createdRecipeIteration, err)

			assert.Error(t, prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, nonexistentID))

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

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

			// Create iteration media.
			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = createdRecipeIteration.ID
			exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
			createdIterationMedia, err := prixfixeClient.CreateIterationMedia(ctx, createdRecipe.ID, exampleIterationMediaInput)
			checkValueAndError(t, createdIterationMedia, err)

			// Clean up iteration media.
			assert.NoError(t, prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID))

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

			// Create iteration media.
			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = createdRecipeIteration.ID
			exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
			createdIterationMedia, err := prixfixeClient.CreateIterationMedia(ctx, createdRecipe.ID, exampleIterationMediaInput)
			checkValueAndError(t, createdIterationMedia, err)

			assert.Error(t, prixfixeClient.ArchiveIterationMedia(ctx, nonexistentID, createdRecipeIteration.ID, createdIterationMedia.ID))

			// Clean up iteration media.
			assert.NoError(t, prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID))

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("returns error when trying to archive post belonging to nonexistent recipe iteration", func(t *testing.T) {
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

			// Create iteration media.
			exampleIterationMedia := fakemodels.BuildFakeIterationMedia()
			exampleIterationMedia.BelongsToRecipeIteration = createdRecipeIteration.ID
			exampleIterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInputFromIterationMedia(exampleIterationMedia)
			createdIterationMedia, err := prixfixeClient.CreateIterationMedia(ctx, createdRecipe.ID, exampleIterationMediaInput)
			checkValueAndError(t, createdIterationMedia, err)

			assert.Error(t, prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, nonexistentID, createdIterationMedia.ID))

			// Clean up iteration media.
			assert.NoError(t, prixfixeClient.ArchiveIterationMedia(ctx, createdRecipe.ID, createdRecipeIteration.ID, createdIterationMedia.ID))

			// Clean up recipe iteration.
			assert.NoError(t, prixfixeClient.ArchiveRecipeIteration(ctx, createdRecipe.ID, createdRecipeIteration.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})
}
