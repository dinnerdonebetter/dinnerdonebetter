package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkRequiredPreparationInstrumentEquality(t *testing.T, expected, actual *models.RequiredPreparationInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.ValidInstrumentID, actual.ValidInstrumentID, "expected ValidInstrumentID for ID %d to be %v, but it was %v ", expected.ID, expected.ValidInstrumentID, actual.ValidInstrumentID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for ID %d to be %v, but it was %v ", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

func TestRequiredPreparationInstruments(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrument.BelongsToValidPreparation = createdValidPreparation.ID
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			// Assert required preparation instrument equality.
			checkRequiredPreparationInstrumentEquality(t, exampleRequiredPreparationInstrument, createdRequiredPreparationInstrument)

			// Clean up.
			err = prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetRequiredPreparationInstrument(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID)
			checkValueAndError(t, actual, err)
			checkRequiredPreparationInstrumentEquality(t, exampleRequiredPreparationInstrument, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})

		T.Run("should fail to create for nonexistent valid preparation", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrument.BelongsToValidPreparation = nonexistentID
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)

			assert.Nil(t, createdRequiredPreparationInstrument)
			assert.Error(t, err)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Create required preparation instruments.
			var expected []*models.RequiredPreparationInstrument
			for i := 0; i < 5; i++ {
				// Create required preparation instrument.
				exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
				exampleRequiredPreparationInstrument.BelongsToValidPreparation = createdValidPreparation.ID
				exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
				createdRequiredPreparationInstrument, requiredPreparationInstrumentCreationErr := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
				checkValueAndError(t, createdRequiredPreparationInstrument, requiredPreparationInstrumentCreationErr)

				expected = append(expected, createdRequiredPreparationInstrument)
			}

			// Assert required preparation instrument list equality.
			actual, err := prixfixeClient.GetRequiredPreparationInstruments(ctx, createdValidPreparation.ID, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RequiredPreparationInstruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RequiredPreparationInstruments),
			)

			// Clean up.
			for _, createdRequiredPreparationInstrument := range actual.RequiredPreparationInstruments {
				err = prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID)
				assert.NoError(t, err)
			}

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Attempt to fetch nonexistent required preparation instrument.
			actual, err := prixfixeClient.RequiredPreparationInstrumentExists(ctx, createdValidPreparation.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})

		T.Run("it should return true with no error when the relevant required preparation instrument exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrument.BelongsToValidPreparation = createdValidPreparation.ID
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			// Fetch required preparation instrument.
			actual, err := prixfixeClient.RequiredPreparationInstrumentExists(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up required preparation instrument.
			assert.NoError(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID))

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Attempt to fetch nonexistent required preparation instrument.
			_, err = prixfixeClient.GetRequiredPreparationInstrument(ctx, createdValidPreparation.ID, nonexistentID)
			assert.Error(t, err)

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrument.BelongsToValidPreparation = createdValidPreparation.ID
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			// Fetch required preparation instrument.
			actual, err := prixfixeClient.GetRequiredPreparationInstrument(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID)
			checkValueAndError(t, actual, err)

			// Assert required preparation instrument equality.
			checkRequiredPreparationInstrumentEquality(t, exampleRequiredPreparationInstrument, actual)

			// Clean up required preparation instrument.
			assert.NoError(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID))

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrument.BelongsToValidPreparation = createdValidPreparation.ID
			exampleRequiredPreparationInstrument.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument))

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrument.BelongsToValidPreparation = createdValidPreparation.ID
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			// Change required preparation instrument.
			createdRequiredPreparationInstrument.Update(exampleRequiredPreparationInstrument.ToUpdateInput())
			err = prixfixeClient.UpdateRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument)
			assert.NoError(t, err)

			// Fetch required preparation instrument.
			actual, err := prixfixeClient.GetRequiredPreparationInstrument(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID)
			checkValueAndError(t, actual, err)

			// Assert required preparation instrument equality.
			checkRequiredPreparationInstrumentEquality(t, exampleRequiredPreparationInstrument, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up required preparation instrument.
			assert.NoError(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID))

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})

		T.Run("it should return an error when trying to update something that belongs to a valid preparation that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrument.BelongsToValidPreparation = createdValidPreparation.ID
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			// Change required preparation instrument.
			createdRequiredPreparationInstrument.Update(exampleRequiredPreparationInstrument.ToUpdateInput())
			createdRequiredPreparationInstrument.BelongsToValidPreparation = nonexistentID
			err = prixfixeClient.UpdateRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument)
			assert.Error(t, err)

			// Clean up required preparation instrument.
			assert.NoError(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID))

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			assert.Error(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdValidPreparation.ID, nonexistentID))

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrument.BelongsToValidPreparation = createdValidPreparation.ID
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			// Clean up required preparation instrument.
			assert.NoError(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID))

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})

		T.Run("returns error when trying to archive post belonging to nonexistent valid preparation", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
			checkValueAndError(t, createdValidPreparation, err)

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrument.BelongsToValidPreparation = createdValidPreparation.ID
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			assert.Error(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, nonexistentID, createdRequiredPreparationInstrument.ID))

			// Clean up required preparation instrument.
			assert.NoError(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdValidPreparation.ID, createdRequiredPreparationInstrument.ID))

			// Clean up valid preparation.
			assert.NoError(t, prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID))
		})
	})
}
