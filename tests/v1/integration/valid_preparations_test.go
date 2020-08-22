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

func checkValidPreparationEquality(t *testing.T, expected, actual *models.ValidPreparation) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for ID %d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Icon, actual.Icon, "expected Icon for ID %d to be %v, but it was %v ", expected.ID, expected.Icon, actual.Icon)
	assert.Equal(t, expected.ApplicableToAllIngredients, actual.ApplicableToAllIngredients, "expected ApplicableToAllIngredients for ID %d to be %v, but it was %v ", expected.ID, expected.ApplicableToAllIngredients, actual.ApplicableToAllIngredients)
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

	test.Run("Searching", func(T *testing.T) {
		T.Run("should be able to be search for valid preparations", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid preparations.
			exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
			var expected []*models.ValidPreparation
			for i := 0; i < 5; i++ {
				// Create valid preparation.
				exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
				exampleValidPreparationInput.Name = fmt.Sprintf("%s %d", exampleValidPreparationInput.Name, i)
				createdValidPreparation, validPreparationCreationErr := prixfixeClient.CreateValidPreparation(ctx, exampleValidPreparationInput)
				checkValueAndError(t, createdValidPreparation, validPreparationCreationErr)

				expected = append(expected, createdValidPreparation)
			}

			exampleLimit := uint8(20)

			// Assert valid preparation list equality.
			actual, err := prixfixeClient.SearchValidPreparations(ctx, exampleValidPreparation.Name, exampleLimit)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// Clean up.
			for _, createdValidPreparation := range expected {
				err = prixfixeClient.ArchiveValidPreparation(ctx, createdValidPreparation.ID)
				assert.NoError(t, err)
			}
		})

		T.Run("should only receive your own valid preparations", func(t *testing.T) {
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

			// Create valid preparations for user A.
			exampleValidPreparationA := fakemodels.BuildFakeValidPreparation()
			var createdForA []*models.ValidPreparation
			for i := 0; i < 5; i++ {
				// Create valid preparation.
				exampleValidPreparationInputA := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparationA)
				exampleValidPreparationInputA.Name = fmt.Sprintf("%s %d", exampleValidPreparationInputA.Name, i)

				createdValidPreparation, validPreparationCreationErr := clientA.CreateValidPreparation(ctx, exampleValidPreparationInputA)
				checkValueAndError(t, createdValidPreparation, validPreparationCreationErr)

				createdForA = append(createdForA, createdValidPreparation)
			}

			exampleLimit := uint8(20)
			query := exampleValidPreparationA.Name

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

			// Create valid preparations for user B.
			exampleValidPreparationB := fakemodels.BuildFakeValidPreparation()
			exampleValidPreparationB.Name = reverse(exampleValidPreparationA.Name)
			var createdForB []*models.ValidPreparation
			for i := 0; i < 5; i++ {
				// Create valid preparation.
				exampleValidPreparationInputB := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparationB)
				exampleValidPreparationInputB.Name = fmt.Sprintf("%s %d", exampleValidPreparationInputB.Name, i)

				createdValidPreparation, validPreparationCreationErr := clientB.CreateValidPreparation(ctx, exampleValidPreparationInputB)
				checkValueAndError(t, createdValidPreparation, validPreparationCreationErr)

				createdForB = append(createdForB, createdValidPreparation)
			}

			expected := createdForA

			// Assert valid preparation list equality.
			actual, err := clientA.SearchValidPreparations(ctx, query, exampleLimit)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// Clean up.
			for _, createdValidPreparation := range createdForA {
				err = clientA.ArchiveValidPreparation(ctx, createdValidPreparation.ID)
				assert.NoError(t, err)
			}

			for _, createdValidPreparation := range createdForB {
				err = clientB.ArchiveValidPreparation(ctx, createdValidPreparation.ID)
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
			assert.NotNil(t, actual.LastUpdatedOn)

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
