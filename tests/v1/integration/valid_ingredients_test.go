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

func checkValidIngredientEquality(t *testing.T, expected, actual *models.ValidIngredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for ID %d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Variant, actual.Variant, "expected Variant for ID %d to be %v, but it was %v ", expected.ID, expected.Variant, actual.Variant)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for ID %d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Warning, actual.Warning, "expected Warning for ID %d to be %v, but it was %v ", expected.ID, expected.Warning, actual.Warning)
	assert.Equal(t, expected.ContainsEgg, actual.ContainsEgg, "expected ContainsEgg for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsEgg, actual.ContainsEgg)
	assert.Equal(t, expected.ContainsDairy, actual.ContainsDairy, "expected ContainsDairy for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsDairy, actual.ContainsDairy)
	assert.Equal(t, expected.ContainsPeanut, actual.ContainsPeanut, "expected ContainsPeanut for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsPeanut, actual.ContainsPeanut)
	assert.Equal(t, expected.ContainsTreeNut, actual.ContainsTreeNut, "expected ContainsTreeNut for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsTreeNut, actual.ContainsTreeNut)
	assert.Equal(t, expected.ContainsSoy, actual.ContainsSoy, "expected ContainsSoy for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsSoy, actual.ContainsSoy)
	assert.Equal(t, expected.ContainsWheat, actual.ContainsWheat, "expected ContainsWheat for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsWheat, actual.ContainsWheat)
	assert.Equal(t, expected.ContainsShellfish, actual.ContainsShellfish, "expected ContainsShellfish for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsShellfish, actual.ContainsShellfish)
	assert.Equal(t, expected.ContainsSesame, actual.ContainsSesame, "expected ContainsSesame for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsSesame, actual.ContainsSesame)
	assert.Equal(t, expected.ContainsFish, actual.ContainsFish, "expected ContainsFish for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsFish, actual.ContainsFish)
	assert.Equal(t, expected.ContainsGluten, actual.ContainsGluten, "expected ContainsGluten for ID %d to be %v, but it was %v ", expected.ID, expected.ContainsGluten, actual.ContainsGluten)
	assert.Equal(t, expected.AnimalFlesh, actual.AnimalFlesh, "expected AnimalFlesh for ID %d to be %v, but it was %v ", expected.ID, expected.AnimalFlesh, actual.AnimalFlesh)
	assert.Equal(t, expected.AnimalDerived, actual.AnimalDerived, "expected AnimalDerived for ID %d to be %v, but it was %v ", expected.ID, expected.AnimalDerived, actual.AnimalDerived)
	assert.Equal(t, expected.MeasurableByVolume, actual.MeasurableByVolume, "expected MeasurableByVolume for ID %d to be %v, but it was %v ", expected.ID, expected.MeasurableByVolume, actual.MeasurableByVolume)
	assert.Equal(t, expected.Icon, actual.Icon, "expected Icon for ID %d to be %v, but it was %v ", expected.ID, expected.Icon, actual.Icon)
	assert.NotZero(t, actual.CreatedOn)
}

func TestValidIngredients(test *testing.T) {
	test.Run("Creating", func(T *testing.T) {
		T.Run("should be createable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Assert valid ingredient equality.
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			// Clean up.
			err = prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID)
			assert.NoError(t, err)

			actual, err := prixfixeClient.GetValidIngredient(ctx, createdValidIngredient.ID)
			checkValueAndError(t, actual, err)
			checkValidIngredientEquality(t, exampleValidIngredient, actual)
			assert.NotNil(t, actual.ArchivedOn)
			assert.NotZero(t, actual.ArchivedOn)
		})
	})

	test.Run("Listing", func(T *testing.T) {
		T.Run("should be able to be read in a list", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredients.
			var expected []*models.ValidIngredient
			for i := 0; i < 5; i++ {
				// Create valid ingredient.
				exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
				createdValidIngredient, validIngredientCreationErr := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
				checkValueAndError(t, createdValidIngredient, validIngredientCreationErr)

				expected = append(expected, createdValidIngredient)
			}

			// Assert valid ingredient list equality.
			actual, err := prixfixeClient.GetValidIngredients(ctx, nil)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidIngredients),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidIngredients),
			)

			// Clean up.
			for _, createdValidIngredient := range actual.ValidIngredients {
				err = prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("Searching", func(T *testing.T) {
		T.Run("should be able to be search for valid ingredients", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredients.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			var expected []*models.ValidIngredient
			for i := 0; i < 5; i++ {
				// Create valid ingredient.
				exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
				exampleValidIngredientInput.Name = fmt.Sprintf("%s %d", exampleValidIngredientInput.Name, i)
				createdValidIngredient, validIngredientCreationErr := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
				checkValueAndError(t, createdValidIngredient, validIngredientCreationErr)

				expected = append(expected, createdValidIngredient)
			}

			exampleLimit := uint8(20)

			// Assert valid ingredient list equality.
			actual, err := prixfixeClient.SearchValidIngredients(ctx, exampleValidIngredient.Name, exampleLimit)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// Clean up.
			for _, createdValidIngredient := range expected {
				err = prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID)
				assert.NoError(t, err)
			}
		})

		T.Run("should only receive your own valid ingredients", func(t *testing.T) {
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

			// Create valid ingredients for user A.
			exampleValidIngredientA := fakemodels.BuildFakeValidIngredient()
			var createdForA []*models.ValidIngredient
			for i := 0; i < 5; i++ {
				// Create valid ingredient.
				exampleValidIngredientInputA := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredientA)
				exampleValidIngredientInputA.Name = fmt.Sprintf("%s %d", exampleValidIngredientInputA.Name, i)

				createdValidIngredient, validIngredientCreationErr := clientA.CreateValidIngredient(ctx, exampleValidIngredientInputA)
				checkValueAndError(t, createdValidIngredient, validIngredientCreationErr)

				createdForA = append(createdForA, createdValidIngredient)
			}

			exampleLimit := uint8(20)
			query := exampleValidIngredientA.Name

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

			// Create valid ingredients for user B.
			exampleValidIngredientB := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientB.Name = reverse(exampleValidIngredientA.Name)
			var createdForB []*models.ValidIngredient
			for i := 0; i < 5; i++ {
				// Create valid ingredient.
				exampleValidIngredientInputB := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredientB)
				exampleValidIngredientInputB.Name = fmt.Sprintf("%s %d", exampleValidIngredientInputB.Name, i)

				createdValidIngredient, validIngredientCreationErr := clientB.CreateValidIngredient(ctx, exampleValidIngredientInputB)
				checkValueAndError(t, createdValidIngredient, validIngredientCreationErr)

				createdForB = append(createdForB, createdValidIngredient)
			}

			expected := createdForA

			// Assert valid ingredient list equality.
			actual, err := clientA.SearchValidIngredients(ctx, query, exampleLimit)
			checkValueAndError(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// Clean up.
			for _, createdValidIngredient := range createdForA {
				err = clientA.ArchiveValidIngredient(ctx, createdValidIngredient.ID)
				assert.NoError(t, err)
			}

			for _, createdValidIngredient := range createdForB {
				err = clientB.ArchiveValidIngredient(ctx, createdValidIngredient.ID)
				assert.NoError(t, err)
			}
		})
	})

	test.Run("ExistenceChecking", func(T *testing.T) {
		T.Run("it should return false with no error when checking something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent valid ingredient.
			actual, err := prixfixeClient.ValidIngredientExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		})

		T.Run("it should return true with no error when the relevant valid ingredient exists", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Fetch valid ingredient.
			actual, err := prixfixeClient.ValidIngredientExists(ctx, createdValidIngredient.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})
	})

	test.Run("Reading", func(T *testing.T) {
		T.Run("it should return an error when trying to read something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Attempt to fetch nonexistent valid ingredient.
			_, err := prixfixeClient.GetValidIngredient(ctx, nonexistentID)
			assert.Error(t, err)
		})

		T.Run("it should be readable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Fetch valid ingredient.
			actual, err := prixfixeClient.GetValidIngredient(ctx, createdValidIngredient.ID)
			checkValueAndError(t, actual, err)

			// Assert valid ingredient equality.
			checkValidIngredientEquality(t, exampleValidIngredient, actual)

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})
	})

	test.Run("Updating", func(T *testing.T) {
		T.Run("it should return an error when trying to update something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredient.ID = nonexistentID

			assert.Error(t, prixfixeClient.UpdateValidIngredient(ctx, exampleValidIngredient))
		})

		T.Run("it should be updatable", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Change valid ingredient.
			createdValidIngredient.Update(exampleValidIngredient.ToUpdateInput())
			err = prixfixeClient.UpdateValidIngredient(ctx, createdValidIngredient)
			assert.NoError(t, err)

			// Fetch valid ingredient.
			actual, err := prixfixeClient.GetValidIngredient(ctx, createdValidIngredient.ID)
			checkValueAndError(t, actual, err)

			// Assert valid ingredient equality.
			checkValidIngredientEquality(t, exampleValidIngredient, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})
	})

	test.Run("Deleting", func(T *testing.T) {
		T.Run("it should return an error when trying to delete something that does not exist", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			assert.Error(t, prixfixeClient.ArchiveValidIngredient(ctx, nonexistentID))
		})

		T.Run("should be able to be deleted", func(t *testing.T) {
			ctx, span := tracing.StartSpan(context.Background(), t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := prixfixeClient.CreateValidIngredient(ctx, exampleValidIngredientInput)
			checkValueAndError(t, createdValidIngredient, err)

			// Clean up valid ingredient.
			assert.NoError(t, prixfixeClient.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		})
	})
}
