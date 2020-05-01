package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkRecipeTagEquality(t *testing.T, expected, actual *models.RecipeTag) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.NotZero(t, actual.CreatedOn)
}

func TestRecipeTags(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe tag.
			exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
			exampleRecipeTag.BelongsToRecipe = createdRecipe.ID
			exampleRecipeTagInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)
			createdRecipeTag, err := prixfixeClient.CreateRecipeTag(ctx, exampleRecipeTagInput)
			checkValueAndError(t, createdRecipeTag, err)

			// Assert recipe tag equality.
			checkRecipeTagEquality(t, exampleRecipeTag, createdRecipeTag)

			// Clean up.
			err = prixfixeClient.ArchiveRecipeTag(ctx, createdRecipe.ID, createdRecipeTag.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetRecipeTag(ctx, createdRecipe.ID, createdRecipeTag.ID)
			checkValueAndError(t, actual, err)
			checkRecipeTagEquality(t, exampleRecipeTag, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("should fail to create for nonexistent recipe", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe tag.
			exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
			exampleRecipeTag.BelongsToRecipe = nonexistentID
			exampleRecipeTagInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)
			createdRecipeTag, err := prixfixeClient.CreateRecipeTag(ctx, exampleRecipeTagInput)

			assert.Nil(t, createdRecipeTag)
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

			// Create recipe tags.
			var expected []*models.RecipeTag
			for i := 0; i < 5; i++ {
				// Create recipe tag.
				exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
				exampleRecipeTag.BelongsToRecipe = createdRecipe.ID
				exampleRecipeTagInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)
				createdRecipeTag, recipeTagCreationErr := prixfixeClient.CreateRecipeTag(ctx, exampleRecipeTagInput)
				checkValueAndError(t, createdRecipeTag, recipeTagCreationErr)

				expected = append(expected, createdRecipeTag)
			}

			// Assert recipe tag list equality.
			actual, err := prixfixeClient.GetRecipeTags(ctx, createdRecipe.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeTags),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeTags),
			)

			// Clean up.
			for _, createdRecipeTag := range actual.RecipeTags {
				err = prixfixeClient.ArchiveRecipeTag(ctx, createdRecipe.ID, createdRecipeTag.ID)
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

			// Attempt to fetch nonexistent recipe tag.
			actual, err := prixfixeClient.RecipeTagExists(ctx, createdRecipe.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})

		T.Run("it should return true with no error when the relevant recipe tag exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakemodels.BuildFakeRecipe()
			exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := prixfixeClient.CreateRecipe(ctx, exampleRecipeInput)
			checkValueAndError(t, createdRecipe, err)

			// Create recipe tag.
			exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
			exampleRecipeTag.BelongsToRecipe = createdRecipe.ID
			exampleRecipeTagInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)
			createdRecipeTag, err := prixfixeClient.CreateRecipeTag(ctx, exampleRecipeTagInput)
			checkValueAndError(t, createdRecipeTag, err)

			// Fetch recipe tag.
			actual, err := prixfixeClient.RecipeTagExists(ctx, createdRecipe.ID, createdRecipeTag.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up recipe tag.
			assert.NoError(t, prixfixeClient.ArchiveRecipeTag(ctx, createdRecipe.ID, createdRecipeTag.ID))

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

			// Attempt to fetch nonexistent recipe tag.
			_, err = prixfixeClient.GetRecipeTag(ctx, createdRecipe.ID, nonexistentID)
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

			// Create recipe tag.
			exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
			exampleRecipeTag.BelongsToRecipe = createdRecipe.ID
			exampleRecipeTagInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)
			createdRecipeTag, err := prixfixeClient.CreateRecipeTag(ctx, exampleRecipeTagInput)
			checkValueAndError(t, createdRecipeTag, err)

			// Fetch recipe tag.
			actual, err := prixfixeClient.GetRecipeTag(ctx, createdRecipe.ID, createdRecipeTag.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe tag equality.
			checkRecipeTagEquality(t, exampleRecipeTag, actual)

			// Clean up recipe tag.
			assert.NoError(t, prixfixeClient.ArchiveRecipeTag(ctx, createdRecipe.ID, createdRecipeTag.ID))

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

			exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
			exampleRecipeTag.BelongsToRecipe = createdRecipe.ID
			exampleRecipeTag.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateRecipeTag(ctx, exampleRecipeTag))

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

			// Create recipe tag.
			exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
			exampleRecipeTag.BelongsToRecipe = createdRecipe.ID
			exampleRecipeTagInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)
			createdRecipeTag, err := prixfixeClient.CreateRecipeTag(ctx, exampleRecipeTagInput)
			checkValueAndError(t, createdRecipeTag, err)

			// Change recipe tag.
			createdRecipeTag.Update(exampleRecipeTag.ToUpdateInput())
			err = prixfixeClient.UpdateRecipeTag(ctx, createdRecipeTag)
			assert.NoError(t, err)

			// Fetch recipe tag.
			actual, err := prixfixeClient.GetRecipeTag(ctx, createdRecipe.ID, createdRecipeTag.ID)
			checkValueAndError(t, actual, err)

			// Assert recipe tag equality.
			checkRecipeTagEquality(t, exampleRecipeTag, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up recipe tag.
			assert.NoError(t, prixfixeClient.ArchiveRecipeTag(ctx, createdRecipe.ID, createdRecipeTag.ID))

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

			// Create recipe tag.
			exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
			exampleRecipeTag.BelongsToRecipe = createdRecipe.ID
			exampleRecipeTagInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)
			createdRecipeTag, err := prixfixeClient.CreateRecipeTag(ctx, exampleRecipeTagInput)
			checkValueAndError(t, createdRecipeTag, err)

			// Change recipe tag.
			createdRecipeTag.Update(exampleRecipeTag.ToUpdateInput())
			createdRecipeTag.BelongsToRecipe = nonexistentID
			err = prixfixeClient.UpdateRecipeTag(ctx, createdRecipeTag)
			assert.Error(t, err)

			// Clean up recipe tag.
			assert.NoError(t, prixfixeClient.ArchiveRecipeTag(ctx, createdRecipe.ID, createdRecipeTag.ID))

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

			assert.Error(t, prixfixeClient.ArchiveRecipeTag(ctx, createdRecipe.ID, nonexistentID))

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

			// Create recipe tag.
			exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
			exampleRecipeTag.BelongsToRecipe = createdRecipe.ID
			exampleRecipeTagInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)
			createdRecipeTag, err := prixfixeClient.CreateRecipeTag(ctx, exampleRecipeTagInput)
			checkValueAndError(t, createdRecipeTag, err)

			// Clean up recipe tag.
			assert.NoError(t, prixfixeClient.ArchiveRecipeTag(ctx, createdRecipe.ID, createdRecipeTag.ID))

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

			// Create recipe tag.
			exampleRecipeTag := fakemodels.BuildFakeRecipeTag()
			exampleRecipeTag.BelongsToRecipe = createdRecipe.ID
			exampleRecipeTagInput := fakemodels.BuildFakeRecipeTagCreationInputFromRecipeTag(exampleRecipeTag)
			createdRecipeTag, err := prixfixeClient.CreateRecipeTag(ctx, exampleRecipeTagInput)
			checkValueAndError(t, createdRecipeTag, err)

			assert.Error(t, prixfixeClient.ArchiveRecipeTag(ctx, nonexistentID, createdRecipeTag.ID))

			// Clean up recipe tag.
			assert.NoError(t, prixfixeClient.ArchiveRecipeTag(ctx, createdRecipe.ID, createdRecipeTag.ID))

			// Clean up recipe.
			assert.NoError(t, prixfixeClient.ArchiveRecipe(ctx, createdRecipe.ID))
		})
	})
}
