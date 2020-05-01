package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkValidIngredientPreparationEquality(t *testing.T, expected, actual *models.ValidIngredientPreparation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for ID %d to be %v, but it was %v ", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

func TestValidIngredientPreparations(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Create valid ingredient preparation.
			exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.BelongsToValidIngredient = createdValidIngredient.ID
			exampleValidIngredientPreparationInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := prixfixeClient.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			checkValueAndError(t, createdValidIngredientPreparation, err)

			// Assert valid ingredient preparation equality.
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, createdValidIngredientPreparation)

			// Clean up.
			err = prixfixeClient.ArchiveValidIngredientPreparation(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetValidIngredientPreparation(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID)
			checkValueAndError(t, actual, err)
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})

		T.Run("should fail to create for nonexistent valid ingredient", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient preparation.
			exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.BelongsToValidIngredient = nonexistentID
			exampleValidIngredientPreparationInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := prixfixeClient.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)

			assert.Nil(t, createdValidIngredientPreparation)
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

			// Create valid ingredient preparations.
			var expected []*models.ValidIngredientPreparation
			for i := 0; i < 5; i++ {
				// Create valid ingredient preparation.
				exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
				exampleValidIngredientPreparation.BelongsToValidIngredient = createdValidIngredient.ID
				exampleValidIngredientPreparationInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
				createdValidIngredientPreparation, validIngredientPreparationCreationErr := prixfixeClient.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
				checkValueAndError(t, createdValidIngredientPreparation, validIngredientPreparationCreationErr)

				expected = append(expected, createdValidIngredientPreparation)
			}

			// Assert valid ingredient preparation list equality.
			actual, err := prixfixeClient.GetValidIngredientPreparations(ctx, createdValidIngredient.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidIngredientPreparations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidIngredientPreparations),
			)

			// Clean up.
			for _, createdValidIngredientPreparation := range actual.ValidIngredientPreparations {
				err = prixfixeClient.ArchiveValidIngredientPreparation(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID)
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

			// Attempt to fetch nonexistent valid ingredient preparation.
			actual, err := prixfixeClient.ValidIngredientPreparationExists(ctx, createdValidIngredient.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})

		T.Run("it should return true with no error when the relevant valid ingredient preparation exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Create valid ingredient preparation.
			exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.BelongsToValidIngredient = createdValidIngredient.ID
			exampleValidIngredientPreparationInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := prixfixeClient.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			checkValueAndError(t, createdValidIngredientPreparation, err)

			// Fetch valid ingredient preparation.
			actual, err := prixfixeClient.ValidIngredientPreparationExists(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up valid ingredient preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredientPreparation(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID))

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

			// Attempt to fetch nonexistent valid ingredient preparation.
			_, err = prixfixeClient.GetValidIngredientPreparation(ctx, createdValidIngredient.ID, nonexistentID)
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

			// Create valid ingredient preparation.
			exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.BelongsToValidIngredient = createdValidIngredient.ID
			exampleValidIngredientPreparationInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := prixfixeClient.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			checkValueAndError(t, createdValidIngredientPreparation, err)

			// Fetch valid ingredient preparation.
			actual, err := prixfixeClient.GetValidIngredientPreparation(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID)
			checkValueAndError(t, actual, err)

			// Assert valid ingredient preparation equality.
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, actual)

			// Clean up valid ingredient preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredientPreparation(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID))

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

			exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.BelongsToValidIngredient = createdValidIngredient.ID
			exampleValidIngredientPreparation.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateValidIngredientPreparation(ctx, exampleValidIngredientPreparation))

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

			// Create valid ingredient preparation.
			exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.BelongsToValidIngredient = createdValidIngredient.ID
			exampleValidIngredientPreparationInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := prixfixeClient.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			checkValueAndError(t, createdValidIngredientPreparation, err)

			// Change valid ingredient preparation.
			createdValidIngredientPreparation.Update(exampleValidIngredientPreparation.ToUpdateInput())
			err = prixfixeClient.UpdateValidIngredientPreparation(ctx, createdValidIngredientPreparation)
			assert.NoError(t, err)

			// Fetch valid ingredient preparation.
			actual, err := prixfixeClient.GetValidIngredientPreparation(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID)
			checkValueAndError(t, actual, err)

			// Assert valid ingredient preparation equality.
			checkValidIngredientPreparationEquality(t, exampleValidIngredientPreparation, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up valid ingredient preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredientPreparation(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID))

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

			// Create valid ingredient preparation.
			exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.BelongsToValidIngredient = createdValidIngredient.ID
			exampleValidIngredientPreparationInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := prixfixeClient.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			checkValueAndError(t, createdValidIngredientPreparation, err)

			// Change valid ingredient preparation.
			createdValidIngredientPreparation.Update(exampleValidIngredientPreparation.ToUpdateInput())
			createdValidIngredientPreparation.BelongsToValidIngredient = nonexistentID
			err = prixfixeClient.UpdateValidIngredientPreparation(ctx, createdValidIngredientPreparation)
			assert.Error(t, err)

			// Clean up valid ingredient preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredientPreparation(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID))

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

			assert.Error(t, prixfixeClient.ArchiveValidIngredientPreparation(ctx, createdValidIngredient.ID, nonexistentID))

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

			// Create valid ingredient preparation.
			exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.BelongsToValidIngredient = createdValidIngredient.ID
			exampleValidIngredientPreparationInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := prixfixeClient.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			checkValueAndError(t, createdValidIngredientPreparation, err)

			// Clean up valid ingredient preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredientPreparation(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID))

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

			// Create valid ingredient preparation.
			exampleValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparation()
			exampleValidIngredientPreparation.BelongsToValidIngredient = createdValidIngredient.ID
			exampleValidIngredientPreparationInput := fakemodels.BuildFakeValidIngredientPreparationCreationInputFromValidIngredientPreparation(exampleValidIngredientPreparation)
			createdValidIngredientPreparation, err := prixfixeClient.CreateValidIngredientPreparation(ctx, exampleValidIngredientPreparationInput)
			checkValueAndError(t, createdValidIngredientPreparation, err)

			assert.Error(t, prixfixeClient.ArchiveValidIngredientPreparation(ctx, nonexistentID, createdValidIngredientPreparation.ID))

			// Clean up valid ingredient preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredientPreparation(ctx, createdValidIngredient.ID, createdValidIngredientPreparation.ID))

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})
	})
}
