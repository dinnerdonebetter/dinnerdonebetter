package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkValidIngredientTagEquality(t *testing.T, expected, actual *models.ValidIngredientTag) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.NotZero(t, actual.CreatedOn)
}

func TestValidIngredientTags(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient tag.
			exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
			exampleValidIngredientTagInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)
			createdValidIngredientTag, err := prixfixeClient.CreateValidIngredientTag(ctx, exampleValidIngredientTagInput)
			checkValueAndError(t, createdValidIngredientTag, err)

			// Assert valid ingredient tag equality.
			checkValidIngredientTagEquality(t, exampleValidIngredientTag, createdValidIngredientTag)

			// Clean up.
			err = prixfixeClient.ArchiveValidIngredientTag(ctx, createdValidIngredientTag.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetValidIngredientTag(ctx, createdValidIngredientTag.ID)
			checkValueAndError(t, actual, err)
			checkValidIngredientTagEquality(t, exampleValidIngredientTag, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient tags.
			var expected []*models.ValidIngredientTag
			for i := 0; i < 5; i++ {
				// Create valid ingredient tag.
				exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
				exampleValidIngredientTagInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)
				createdValidIngredientTag, validIngredientTagCreationErr := prixfixeClient.CreateValidIngredientTag(ctx, exampleValidIngredientTagInput)
				checkValueAndError(t, createdValidIngredientTag, validIngredientTagCreationErr)

				expected = append(expected, createdValidIngredientTag)
			}

			// Assert valid ingredient tag list equality.
			actual, err := prixfixeClient.GetValidIngredientTags(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidIngredientTags),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidIngredientTags),
			)

			// Clean up.
			for _, createdValidIngredientTag := range actual.ValidIngredientTags {
				err = prixfixeClient.ArchiveValidIngredientTag(ctx, createdValidIngredientTag.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent valid ingredient tag.
			actual, err := prixfixeClient.ValidIngredientTagExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		})

		T.Run("it should return true with no error when the relevant valid ingredient tag exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient tag.
			exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
			exampleValidIngredientTagInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)
			createdValidIngredientTag, err := prixfixeClient.CreateValidIngredientTag(ctx, exampleValidIngredientTagInput)
			checkValueAndError(t, createdValidIngredientTag, err)

			// Fetch valid ingredient tag.
			actual, err := prixfixeClient.ValidIngredientTagExists(ctx, createdValidIngredientTag.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up valid ingredient tag.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredientTag(ctx, createdValidIngredientTag.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent valid ingredient tag.
			_, err := prixfixeClient.GetValidIngredientTag(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient tag.
			exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
			exampleValidIngredientTagInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)
			createdValidIngredientTag, err := prixfixeClient.CreateValidIngredientTag(ctx, exampleValidIngredientTagInput)
			checkValueAndError(t, createdValidIngredientTag, err)

			// Fetch valid ingredient tag.
			actual, err := prixfixeClient.GetValidIngredientTag(ctx, createdValidIngredientTag.ID)
			checkValueAndError(t, actual, err)

			// Assert valid ingredient tag equality.
			checkValidIngredientTagEquality(t, exampleValidIngredientTag, actual)

			// Clean up valid ingredient tag.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredientTag(ctx, createdValidIngredientTag.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
			exampleValidIngredientTag.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateValidIngredientTag(ctx, exampleValidIngredientTag))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient tag.
			exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
			exampleValidIngredientTagInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)
			createdValidIngredientTag, err := prixfixeClient.CreateValidIngredientTag(ctx, exampleValidIngredientTagInput)
			checkValueAndError(t, createdValidIngredientTag, err)

			// Change valid ingredient tag.
			createdValidIngredientTag.Update(exampleValidIngredientTag.ToUpdateInput())
			err = prixfixeClient.UpdateValidIngredientTag(ctx, createdValidIngredientTag)
			assert.NoError(t, err)

			// Fetch valid ingredient tag.
			actual, err := prixfixeClient.GetValidIngredientTag(ctx, createdValidIngredientTag.ID)
			checkValueAndError(t, actual, err)

			// Assert valid ingredient tag equality.
			checkValidIngredientTagEquality(t, exampleValidIngredientTag, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up valid ingredient tag.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredientTag(ctx, createdValidIngredientTag.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			assert.Error(t, prixfixeClient.ArchiveValidIngredientTag(ctx, nonexistentID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient tag.
			exampleValidIngredientTag := fakemodels.BuildFakeValidIngredientTag()
			exampleValidIngredientTagInput := fakemodels.BuildFakeValidIngredientTagCreationInputFromValidIngredientTag(exampleValidIngredientTag)
			createdValidIngredientTag, err := prixfixeClient.CreateValidIngredientTag(ctx, exampleValidIngredientTagInput)
			checkValueAndError(t, createdValidIngredientTag, err)

			// Clean up valid ingredient tag.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredientTag(ctx, createdValidIngredientTag.ID))
		})
	})
}
