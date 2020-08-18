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
	assert.Equal(t, expected.InstrumentID, actual.InstrumentID, "expected InstrumentID for ID %d to be %v, but it was %v ", expected.ID, expected.InstrumentID, actual.InstrumentID)
	assert.Equal(t, expected.PreparationID, actual.PreparationID, "expected PreparationID for ID %d to be %v, but it was %v ", expected.ID, expected.PreparationID, actual.PreparationID)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for ID %d to be %v, but it was %v ", expected.ID, expected.Notes, actual.Notes)
	assert.NotZero(t, actual.CreatedOn)
}

func TestRequiredPreparationInstruments(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			// Assert required preparation instrument equality.
			checkRequiredPreparationInstrumentEquality(t, exampleRequiredPreparationInstrument, createdRequiredPreparationInstrument)

			// Clean up.
			err = prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument.ID)
			checkValueAndError(t, actual, err)
			checkRequiredPreparationInstrumentEquality(t, exampleRequiredPreparationInstrument, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create required preparation instruments.
			var expected []*models.RequiredPreparationInstrument
			for i := 0; i < 5; i++ {
				// Create required preparation instrument.
				exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
				exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
				createdRequiredPreparationInstrument, requiredPreparationInstrumentCreationErr := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
				checkValueAndError(t, createdRequiredPreparationInstrument, requiredPreparationInstrumentCreationErr)

				expected = append(expected, createdRequiredPreparationInstrument)
			}

			// Assert required preparation instrument list equality.
			actual, err := prixfixeClient.GetRequiredPreparationInstruments(ctx, nil)
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
				err = prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent required preparation instrument.
			actual, err := prixfixeClient.RequiredPreparationInstrumentExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		})

		T.Run("it should return true with no error when the relevant required preparation instrument exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			// Fetch required preparation instrument.
			actual, err := prixfixeClient.RequiredPreparationInstrumentExists(ctx, createdRequiredPreparationInstrument.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up required preparation instrument.
			assert.NoError(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent required preparation instrument.
			_, err := prixfixeClient.GetRequiredPreparationInstrument(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			// Fetch required preparation instrument.
			actual, err := prixfixeClient.GetRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument.ID)
			checkValueAndError(t, actual, err)

			// Assert required preparation instrument equality.
			checkRequiredPreparationInstrumentEquality(t, exampleRequiredPreparationInstrument, actual)

			// Clean up required preparation instrument.
			assert.NoError(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrument.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrument))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			// Change required preparation instrument.
			createdRequiredPreparationInstrument.Update(exampleRequiredPreparationInstrument.ToUpdateInput())
			err = prixfixeClient.UpdateRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument)
			assert.NoError(t, err)

			// Fetch required preparation instrument.
			actual, err := prixfixeClient.GetRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument.ID)
			checkValueAndError(t, actual, err)

			// Assert required preparation instrument equality.
			checkRequiredPreparationInstrumentEquality(t, exampleRequiredPreparationInstrument, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up required preparation instrument.
			assert.NoError(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			assert.Error(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, nonexistentID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create required preparation instrument.
			exampleRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrument()
			exampleRequiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInputFromRequiredPreparationInstrument(exampleRequiredPreparationInstrument)
			createdRequiredPreparationInstrument, err := prixfixeClient.CreateRequiredPreparationInstrument(ctx, exampleRequiredPreparationInstrumentInput)
			checkValueAndError(t, createdRequiredPreparationInstrument, err)

			// Clean up required preparation instrument.
			assert.NoError(t, prixfixeClient.ArchiveRequiredPreparationInstrument(ctx, createdRequiredPreparationInstrument.ID))
		})
	})
}
