package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkIngredientTagMappingEquality(t *testing.T, expected, actual *models.IngredientTagMapping) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.ValidIngredientTagID, actual.ValidIngredientTagID, "expected ValidIngredientTagID for ID %d to be %v, but it was %v ", expected.ID, expected.ValidIngredientTagID, actual.ValidIngredientTagID)
	assert.NotZero(t, actual.CreatedOn)
}

func TestIngredientTagMappings(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Create ingredient tag mapping.
			exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
			exampleIngredientTagMapping.BelongsToValidIngredient = createdValidIngredient.ID
			exampleIngredientTagMappingInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)
			createdIngredientTagMapping, err := prixfixeClient.CreateIngredientTagMapping(ctx, exampleIngredientTagMappingInput)
			checkValueAndError(t, createdIngredientTagMapping, err)

			// Assert ingredient tag mapping equality.
			checkIngredientTagMappingEquality(t, exampleIngredientTagMapping, createdIngredientTagMapping)

			// Clean up.
			err = prixfixeClient.ArchiveIngredientTagMapping(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetIngredientTagMapping(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID)
			checkValueAndError(t, actual, err)
			checkIngredientTagMappingEquality(t, exampleIngredientTagMapping, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})

		T.Run("should fail to create for nonexistent valid ingredient", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create ingredient tag mapping.
			exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
			exampleIngredientTagMapping.BelongsToValidIngredient = nonexistentID
			exampleIngredientTagMappingInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)
			createdIngredientTagMapping, err := prixfixeClient.CreateIngredientTagMapping(ctx, exampleIngredientTagMappingInput)

			assert.Nil(t, createdIngredientTagMapping)
			assert.Error(t, err)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Create ingredient tag mappings.
			var expected []*models.IngredientTagMapping
			for i := 0; i < 5; i++ {
				// Create ingredient tag mapping.
				exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
				exampleIngredientTagMapping.BelongsToValidIngredient = createdValidIngredient.ID
				exampleIngredientTagMappingInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)
				createdIngredientTagMapping, ingredientTagMappingCreationErr := prixfixeClient.CreateIngredientTagMapping(ctx, exampleIngredientTagMappingInput)
				checkValueAndError(t, createdIngredientTagMapping, ingredientTagMappingCreationErr)

				expected = append(expected, createdIngredientTagMapping)
			}

			// Assert ingredient tag mapping list equality.
			actual, err := prixfixeClient.GetIngredientTagMappings(ctx, createdValidIngredient.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.IngredientTagMappings),
				"expected %d to be <= %d",
				len(expected),
				len(actual.IngredientTagMappings),
			)

			// Clean up.
			for _, createdIngredientTagMapping := range actual.IngredientTagMappings {
				err = prixfixeClient.ArchiveIngredientTagMapping(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID)
				assert.NoError(t, err)
			}

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Attempt to fetch nonexistent ingredient tag mapping.
			actual, err := prixfixeClient.IngredientTagMappingExists(ctx, createdValidIngredient.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})

		T.Run("it should return true with no error when the relevant ingredient tag mapping exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Create ingredient tag mapping.
			exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
			exampleIngredientTagMapping.BelongsToValidIngredient = createdValidIngredient.ID
			exampleIngredientTagMappingInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)
			createdIngredientTagMapping, err := prixfixeClient.CreateIngredientTagMapping(ctx, exampleIngredientTagMappingInput)
			checkValueAndError(t, createdIngredientTagMapping, err)

			// Fetch ingredient tag mapping.
			actual, err := prixfixeClient.IngredientTagMappingExists(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up ingredient tag mapping.
			assert.NoError(t, prixfixeClient.ArchiveIngredientTagMapping(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID))

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Attempt to fetch nonexistent ingredient tag mapping.
			_, err = prixfixeClient.GetIngredientTagMapping(ctx, createdValidIngredient.ID, nonexistentID)
			assert.Error(t, err)

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Create ingredient tag mapping.
			exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
			exampleIngredientTagMapping.BelongsToValidIngredient = createdValidIngredient.ID
			exampleIngredientTagMappingInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)
			createdIngredientTagMapping, err := prixfixeClient.CreateIngredientTagMapping(ctx, exampleIngredientTagMappingInput)
			checkValueAndError(t, createdIngredientTagMapping, err)

			// Fetch ingredient tag mapping.
			actual, err := prixfixeClient.GetIngredientTagMapping(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID)
			checkValueAndError(t, actual, err)

			// Assert ingredient tag mapping equality.
			checkIngredientTagMappingEquality(t, exampleIngredientTagMapping, actual)

			// Clean up ingredient tag mapping.
			assert.NoError(t, prixfixeClient.ArchiveIngredientTagMapping(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID))

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
			exampleIngredientTagMapping.BelongsToValidIngredient = createdValidIngredient.ID
			exampleIngredientTagMapping.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateIngredientTagMapping(ctx, exampleIngredientTagMapping))

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Create ingredient tag mapping.
			exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
			exampleIngredientTagMapping.BelongsToValidIngredient = createdValidIngredient.ID
			exampleIngredientTagMappingInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)
			createdIngredientTagMapping, err := prixfixeClient.CreateIngredientTagMapping(ctx, exampleIngredientTagMappingInput)
			checkValueAndError(t, createdIngredientTagMapping, err)

			// Change ingredient tag mapping.
			createdIngredientTagMapping.Update(exampleIngredientTagMapping.ToUpdateInput())
			err = prixfixeClient.UpdateIngredientTagMapping(ctx, createdIngredientTagMapping)
			assert.NoError(t, err)

			// Fetch ingredient tag mapping.
			actual, err := prixfixeClient.GetIngredientTagMapping(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID)
			checkValueAndError(t, actual, err)

			// Assert ingredient tag mapping equality.
			checkIngredientTagMappingEquality(t, exampleIngredientTagMapping, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up ingredient tag mapping.
			assert.NoError(t, prixfixeClient.ArchiveIngredientTagMapping(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID))

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})

		T.Run("it should return an error when trying to update something that belongs to a valid ingredient that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Create ingredient tag mapping.
			exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
			exampleIngredientTagMapping.BelongsToValidIngredient = createdValidIngredient.ID
			exampleIngredientTagMappingInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)
			createdIngredientTagMapping, err := prixfixeClient.CreateIngredientTagMapping(ctx, exampleIngredientTagMappingInput)
			checkValueAndError(t, createdIngredientTagMapping, err)

			// Change ingredient tag mapping.
			createdIngredientTagMapping.Update(exampleIngredientTagMapping.ToUpdateInput())
			createdIngredientTagMapping.BelongsToValidIngredient = nonexistentID
			err = prixfixeClient.UpdateIngredientTagMapping(ctx, createdIngredientTagMapping)
			assert.Error(t, err)

			// Clean up ingredient tag mapping.
			assert.NoError(t, prixfixeClient.ArchiveIngredientTagMapping(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID))

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			assert.Error(t, prixfixeClient.ArchiveIngredientTagMapping(ctx, createdValidIngredient.ID, nonexistentID))

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Create ingredient tag mapping.
			exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
			exampleIngredientTagMapping.BelongsToValidIngredient = createdValidIngredient.ID
			exampleIngredientTagMappingInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)
			createdIngredientTagMapping, err := prixfixeClient.CreateIngredientTagMapping(ctx, exampleIngredientTagMappingInput)
			checkValueAndError(t, createdIngredientTagMapping, err)

			// Clean up ingredient tag mapping.
			assert.NoError(t, prixfixeClient.ArchiveIngredientTagMapping(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID))

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})

		T.Run("returns error when trying to archive post belonging to nonexistent valid ingredient", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Create ingredient tag mapping.
			exampleIngredientTagMapping := fakemodels.BuildFakeIngredientTagMapping()
			exampleIngredientTagMapping.BelongsToValidIngredient = createdValidIngredient.ID
			exampleIngredientTagMappingInput := fakemodels.BuildFakeIngredientTagMappingCreationInputFromIngredientTagMapping(exampleIngredientTagMapping)
			createdIngredientTagMapping, err := prixfixeClient.CreateIngredientTagMapping(ctx, exampleIngredientTagMappingInput)
			checkValueAndError(t, createdIngredientTagMapping, err)

			assert.Error(t, prixfixeClient.ArchiveIngredientTagMapping(ctx, nonexistentID, createdIngredientTagMapping.ID))

			// Clean up ingredient tag mapping.
			assert.NoError(t, prixfixeClient.ArchiveIngredientTagMapping(ctx, createdValidIngredient.ID, createdIngredientTagMapping.ID))

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})
	})
}
