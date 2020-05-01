package integration

import (
	"context"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
)

func checkValidInstrumentEquality(t *testing.T, expected, actual *models.ValidInstrument) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Variant, actual.Variant, "expected Variant for ID %d to be %v, but it was %v ", expected.ID, expected.Variant, actual.Variant)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for ID %d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Icon, actual.Icon, "expected Icon for ID %d to be %v, but it was %v ", expected.ID, expected.Icon, actual.Icon)
	assert.NotZero(t, actual.CreatedOn)
}

func TestValidInstruments(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid instrument.
			exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := prixfixeClient.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			checkValueAndError(t, createdValidInstrument, err)

			// Assert valid instrument equality.
			checkValidInstrumentEquality(t, exampleValidInstrument, createdValidInstrument)

			// Clean up.
			err = prixfixeClient.ArchiveValidInstrument(ctx, createdValidInstrument.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetValidInstrument(ctx, createdValidInstrument.ID)
			checkValueAndError(t, actual, err)
			checkValidInstrumentEquality(t, exampleValidInstrument, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid instruments.
			var expected []*models.ValidInstrument
			for i := 0; i < 5; i++ {
				// Create valid instrument.
				exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
				exampleValidInstrumentInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
				createdValidInstrument, validInstrumentCreationErr := prixfixeClient.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				checkValueAndError(t, createdValidInstrument, validInstrumentCreationErr)

				expected = append(expected, createdValidInstrument)
			}

			// Assert valid instrument list equality.
			actual, err := prixfixeClient.GetValidInstruments(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidInstruments),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidInstruments),
			)

			// Clean up.
			for _, createdValidInstrument := range actual.ValidInstruments {
				err = prixfixeClient.ArchiveValidInstrument(ctx, createdValidInstrument.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent valid instrument.
			actual, err := prixfixeClient.ValidInstrumentExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		})

		T.Run("it should return true with no error when the relevant valid instrument exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid instrument.
			exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := prixfixeClient.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			checkValueAndError(t, createdValidInstrument, err)

			// Fetch valid instrument.
			actual, err := prixfixeClient.ValidInstrumentExists(ctx, createdValidInstrument.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up valid instrument.
			assert.NoError(t, prixfixeClient.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent valid instrument.
			_, err := prixfixeClient.GetValidInstrument(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid instrument.
			exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := prixfixeClient.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			checkValueAndError(t, createdValidInstrument, err)

			// Fetch valid instrument.
			actual, err := prixfixeClient.GetValidInstrument(ctx, createdValidInstrument.ID)
			checkValueAndError(t, actual, err)

			// Assert valid instrument equality.
			checkValidInstrumentEquality(t, exampleValidInstrument, actual)

			// Clean up valid instrument.
			assert.NoError(t, prixfixeClient.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
			exampleValidInstrument.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateValidInstrument(ctx, exampleValidInstrument))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid instrument.
			exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := prixfixeClient.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			checkValueAndError(t, createdValidInstrument, err)

			// Change valid instrument.
			createdValidInstrument.Update(exampleValidInstrument.ToUpdateInput())
			err = prixfixeClient.UpdateValidInstrument(ctx, createdValidInstrument)
			assert.NoError(t, err)

			// Fetch valid instrument.
			actual, err := prixfixeClient.GetValidInstrument(ctx, createdValidInstrument.ID)
			checkValueAndError(t, actual, err)

			// Assert valid instrument equality.
			checkValidInstrumentEquality(t, exampleValidInstrument, actual)
			assert.NotNil(t, actual.UpdatedOn)

			// Clean up valid instrument.
			assert.NoError(t, prixfixeClient.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			assert.Error(t, prixfixeClient.ArchiveValidInstrument(ctx, nonexistentID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid instrument.
			exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
			exampleValidInstrumentInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
			createdValidInstrument, err := prixfixeClient.CreateValidInstrument(ctx, exampleValidInstrumentInput)
			checkValueAndError(t, createdValidInstrument, err)

			// Clean up valid instrument.
			assert.NoError(t, prixfixeClient.ArchiveValidInstrument(ctx, createdValidInstrument.ID))
		})
	})
}
