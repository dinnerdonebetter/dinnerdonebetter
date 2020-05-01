package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkValidPreparationEquality(t *testing.T, expected, actual *models.ValidPreparation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for ID %d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Icon, actual.Icon, "expected Icon for ID %d to be %v, but it was %v ", expected.ID, expected.Icon, actual.Icon)
	assert.NotZero(t, actual.CreatedOn)
}

func TestValidPreparations(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Assert valid preparation equality.
			checkValidPreparationEquality(t, exampleValidPreparation, createdValidPreparation)

			// Clean up.
			err = prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetValidPreparation(ctx, createdValidPreparation.ID)
			checkValueAndError(t, actual, err)
			checkValidPreparationEquality(t, exampleValidPreparation, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparations.
			var expected []*models.ValidPreparation
			for i := 0; i < 5; i++ {
				// Create valid preparation.
				exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
				exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
				createdValidPreparation, validPreparationCreationErr := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
				checkValueAndError(t, createdValidPreparation, validPreparationCreationErr)

				expected = append(expected, createdValidPreparation)
			}

			// Assert valid preparation list equality.
			actual, err := prixfixeClient.GetValidPreparations(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidPreparations),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidPreparations),
			)

			// Clean up.
			for _, createdValidPreparation := range actual.ValidPreparations {
				err = prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent valid preparation.
			actual, err := prixfixeClient.ValidPreparationExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		})

		T.Run("it should return true with no error when the relevant valid preparation exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Fetch valid preparation.
			actual, err := prixfixeClient.ValidPreparationExists(ctx, createdValidPreparation.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent valid preparation.
			_, err := prixfixeClient.GetValidPreparation(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Fetch valid preparation.
			actual, err := prixfixeClient.GetValidPreparation(ctx, createdValidPreparation.ID)
			checkValueAndError(t, actual, err)

			// Assert valid preparation equality.
			checkValidPreparationEquality(t, exampleValidPreparation, actual)

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparation.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateValidPreparation(ctx, exampleValidPreparation))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Change valid preparation.
			createdValidPreparation.Update(exampleValidPreparation.ToUpdateInput())
			err = prixfixeClient.UpdateValidPreparation(ctx, createdValidPreparation)
			assert.NoError(t, err)

			// Fetch valid preparation.
			actual, err := prixfixeClient.GetValidPreparation(ctx, createdValidPreparation.ID)
			checkValueAndError(t, actual, err)

			// Assert valid preparation equality.
			checkValidPreparationEquality(t, exampleValidPreparation, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			assert.Error(t, prixfixeClient.ArchiveValidPreparation(ctx, nonexistentID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})
	})
}
