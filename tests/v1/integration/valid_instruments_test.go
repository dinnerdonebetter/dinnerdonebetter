package integration

import (
	"context"
	"fmt"
	"testing"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"
	"gitlab.com/prixfixe/prixfixe/tests/v1/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
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

	test.Run("Searching", func(T *testing.T) {
		T.Run("should be able to be search for valid instruments", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid instruments.
			exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
			var expected []*models.ValidInstrument
			for i := 0; i < 5; i++ {
				// Create valid instrument.
				exampleValidInstrumentInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)
				exampleValidInstrumentInput.Name = fmt.Sprintf("%s %d", exampleValidInstrumentInput.Name, i)
				createdValidInstrument, validInstrumentCreationErr := prixfixeClient.CreateValidInstrument(ctx, exampleValidInstrumentInput)
				checkValueAndError(t, createdValidInstrument, validInstrumentCreationErr)

				expected = append(expected, createdValidInstrument)
			}

			exampleLimit := uint8(20)

			// Assert valid instrument list equality.
			actual, err := prixfixeClient.SearchValidInstruments(ctx, exampleValidInstrument.Name, exampleLimit)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// Clean up.
			for _, createdValidInstrument := range expected {
				err = prixfixeClient.ArchiveValidInstrument(ctx, createdValidInstrument.ID)
				assert.NoError(t, err)
			}
		})

		T.Run("should only receive your own valid instruments", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// create user and oauth2 client A.
			userA, err := testutil.CreateObligatoryUser(urlToUse, debug)
			require.NoError(t, err)

			ca, err := testutil.CreateObligatoryClient(urlToUse, userA)
			require.NoError(t, err)

			clientA, err := client.NewClient(
				ctx,
				ca.ClientID,
				ca.ClientSecret,
				prixfixeClient.URL,
				noop.ProvideNoopLogger(),
				buildHTTPClient(),
				ca.Scopes,
				true,
			)
			checkValueAndError(test, clientA, err)

			// Create valid instruments for user A.
			exampleValidInstrumentA := fakemodels.BuildFakeValidInstrument()
			var createdForA []*models.ValidInstrument
			for i := 0; i < 5; i++ {
				// Create valid instrument.
				exampleValidInstrumentInputA := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrumentA)
				exampleValidInstrumentInputA.Name = fmt.Sprintf("%s %d", exampleValidInstrumentInputA.Name, i)

				createdValidInstrument, validInstrumentCreationErr := clientA.CreateValidInstrument(ctx, exampleValidInstrumentInputA)
				checkValueAndError(t, createdValidInstrument, validInstrumentCreationErr)

				createdForA = append(createdForA, createdValidInstrument)
			}

			exampleLimit := uint8(20)
			query := exampleValidInstrumentA.Name

			// create user and oauth2 client B.
			userB, err := testutil.CreateObligatoryUser(urlToUse, debug)
			require.NoError(t, err)

			cb, err := testutil.CreateObligatoryClient(urlToUse, userB)
			require.NoError(t, err)

			clientB, err := client.NewClient(
				ctx,
				cb.ClientID,
				cb.ClientSecret,
				prixfixeClient.URL,
				noop.ProvideNoopLogger(),
				buildHTTPClient(),
				cb.Scopes,
				true,
			)
			checkValueAndError(test, clientB, err)

			// Create valid instruments for user B.
			exampleValidInstrumentB := fakemodels.BuildFakeValidInstrument()
			exampleValidInstrumentB.Name = reverse(exampleValidInstrumentA.Name)
			var createdForB []*models.ValidInstrument
			for i := 0; i < 5; i++ {
				// Create valid instrument.
				exampleValidInstrumentInputB := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrumentB)
				exampleValidInstrumentInputB.Name = fmt.Sprintf("%s %d", exampleValidInstrumentInputB.Name, i)

				createdValidInstrument, validInstrumentCreationErr := clientB.CreateValidInstrument(ctx, exampleValidInstrumentInputB)
				checkValueAndError(t, createdValidInstrument, validInstrumentCreationErr)

				createdForB = append(createdForB, createdValidInstrument)
			}

			expected := createdForA

			// Assert valid instrument list equality.
			actual, err := clientA.SearchValidInstruments(ctx, query, exampleLimit)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// Clean up.
			for _, createdValidInstrument := range createdForA {
				err = clientA.ArchiveValidInstrument(ctx, createdValidInstrument.ID)
				assert.NoError(t, err)
			}

			for _, createdValidInstrument := range createdForB {
				err = clientB.ArchiveValidInstrument(ctx, createdValidInstrument.ID)
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
			assert.NotNil(t, actual.LastUpdatedOn)

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
