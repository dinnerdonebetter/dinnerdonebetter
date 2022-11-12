package integration

import (
	"github.com/prixfixeco/backend/pkg/types/converters"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/prixfixeco/backend/pkg/types"
)

func checkRecipeStepProductEquality(t *testing.T, expected, actual *types.RecipeStepProduct) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe step product %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Type, actual.Type, "expected Type for recipe step product %s to be %v, but it was %v", expected.ID, expected.Type, actual.Type)
	assert.Equal(t, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID, "expected MeasurementUnit.ID for recipe step product %s to be %v, but it was %v", expected.ID, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID)
	assert.Equal(t, expected.MinimumQuantity, actual.MinimumQuantity, "expected MinimumQuantity for recipe step product %s to be %v, but it was %v", expected.ID, expected.MinimumQuantity, actual.MinimumQuantity)
	assert.Equal(t, expected.MaximumQuantity, actual.MaximumQuantity, "expected MaximumQuantity for recipe step product %s to be %v, but it was %v", expected.ID, expected.MaximumQuantity, actual.MaximumQuantity)
	assert.Equal(t, expected.QuantityNotes, actual.QuantityNotes, "expected QuantityNotes for recipe step product %s to be %v, but it was %v", expected.ID, expected.QuantityNotes, actual.QuantityNotes)
	assert.Equal(t, expected.Compostable, actual.Compostable, "expected Compostable for recipe step product %s to be %v, but was %v", expected.ID, expected.Compostable, actual.Compostable)
	assert.Equal(t, expected.MaximumStorageDurationInSeconds, actual.MaximumStorageDurationInSeconds, "expected MaximumStorageDurationInSeconds for recipe step product %s to be %v, but was %v", expected.ID, expected.MaximumStorageDurationInSeconds, actual.MaximumStorageDurationInSeconds)
	assert.Equal(t, expected.MinimumStorageTemperatureInCelsius, actual.MinimumStorageTemperatureInCelsius, "expected MinimumStorageTemperatureInCelsius for recipe step product %s to be %v, but was %v", expected.ID, expected.MinimumStorageTemperatureInCelsius, actual.MinimumStorageTemperatureInCelsius)
	assert.Equal(t, expected.MaximumStorageTemperatureInCelsius, actual.MaximumStorageTemperatureInCelsius, "expected MaximumStorageTemperatureInCelsius for recipe step product %s to be %v, but was %v", expected.ID, expected.MaximumStorageTemperatureInCelsius, actual.MaximumStorageTemperatureInCelsius)
	assert.Equal(t, expected.StorageInstructions, actual.StorageInstructions, "expected StorageInstructions for recipe step product %s to be %v, but was %v", expected.ID, expected.StorageInstructions, actual.StorageInstructions)
	assert.Equal(t, expected.IsLiquid, actual.IsLiquid, "expected IsLiquid for recipe step product %s to be %v, but was %v", expected.ID, expected.IsLiquid, actual.IsLiquid)
	assert.Equal(t, expected.IsWaste, actual.IsWaste, "expected IsWaste for recipe step product %s to be %v, but was %v", expected.ID, expected.IsWaste, actual.IsWaste)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestRecipeStepProducts_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var (
				createdRecipeStepID string
			)
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating valid measurement unit")
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.admin.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			t.Log("creating recipe step product")
			exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepProduct.MeasurementUnit = *createdValidMeasurementUnit
			exampleRecipeStepProductInput := converters.ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := testClients.user.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			require.NoError(t, err)
			t.Logf("recipe step product %q created", createdRecipeStepProduct.ID)

			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

			createdRecipeStepProduct, err = testClients.user.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID)
			requireNotNilAndNoProblems(t, createdRecipeStepProduct, err)
			require.Equal(t, createdRecipeStepID, createdRecipeStepProduct.BelongsToRecipeStep)

			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

			t.Log("changing recipe step product")
			newRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			newRecipeStepProduct.MeasurementUnit = *createdValidMeasurementUnit
			createdRecipeStepProduct.Update(converters.ConvertRecipeStepProductToRecipeStepProductUpdateRequestInput(newRecipeStepProduct))

			require.NoError(t, testClients.user.UpdateRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepProduct))

			t.Log("fetching changed recipe step product")
			actual, err := testClients.user.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step product equality
			checkRecipeStepProductEquality(t, newRecipeStepProduct, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			t.Log("cleaning up recipe step product")
			assert.NoError(t, testClients.user.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.user, nil)

			var (
				createdRecipeStepID string
			)
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating valid measurement unit")
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.admin.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			t.Logf("valid measurement unit %q created", createdValidMeasurementUnit.ID)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.admin.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			t.Log("creating recipe step products")
			var expected []*types.RecipeStepProduct
			for i := 0; i < 5; i++ {
				exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
				exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepProduct.MeasurementUnit = *createdValidMeasurementUnit
				exampleRecipeStepProductInput := converters.ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(exampleRecipeStepProduct)
				createdRecipeStepProduct, createdRecipeStepProductErr := testClients.user.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
				require.NoError(t, createdRecipeStepProductErr)
				t.Logf("recipe step product %q created", createdRecipeStepProduct.ID)

				checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

				createdRecipeStepProduct, createdRecipeStepProductErr = testClients.user.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID)
				requireNotNilAndNoProblems(t, createdRecipeStepProduct, createdRecipeStepProductErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepProduct.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepProduct)
			}

			// assert recipe step product list equality
			actual, err := testClients.user.GetRecipeStepProducts(ctx, createdRecipe.ID, createdRecipeStepID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepProducts),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepProducts),
			)

			t.Log("cleaning up")
			for _, createdRecipeStepProduct := range expected {
				assert.NoError(t, testClients.user.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.user.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.user.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
