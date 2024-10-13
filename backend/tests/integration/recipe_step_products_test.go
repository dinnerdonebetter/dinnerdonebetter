package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipeStepProductEquality(t *testing.T, expected, actual *types.RecipeStepProduct) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe step product %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Type, actual.Type, "expected IndexType for recipe step product %s to be %v, but it was %v", expected.ID, expected.Type, actual.Type)
	assert.Equal(t, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID, "expected MeasurementUnit.ID for recipe step product %s to be %v, but it was %v", expected.ID, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID)
	assert.Equal(t, expected.Quantity, actual.Quantity, "expected Quantity for recipe step product %s to be %v, but it was %v", expected.ID, expected.Quantity, actual.Quantity)
	assert.Equal(t, expected.QuantityNotes, actual.QuantityNotes, "expected QuantityNotes for recipe step product %s to be %v, but it was %v", expected.ID, expected.QuantityNotes, actual.QuantityNotes)
	assert.Equal(t, expected.Compostable, actual.Compostable, "expected Compostable for recipe step product %s to be %v, but was %v", expected.ID, expected.Compostable, actual.Compostable)
	assert.Equal(t, expected.StorageDurationInSeconds, actual.StorageDurationInSeconds, "expected StorageDurationInSeconds for recipe step product %s to be %v, but was %v", expected.ID, expected.StorageDurationInSeconds, actual.StorageDurationInSeconds)
	assert.Equal(t, expected.StorageTemperatureInCelsius, actual.StorageTemperatureInCelsius, "expected StorageTemperatureInCelsius for recipe step product %s to be %v, but was %v", expected.ID, expected.StorageTemperatureInCelsius, actual.StorageTemperatureInCelsius)
	assert.Equal(t, expected.StorageInstructions, actual.StorageInstructions, "expected StorageInstructions for recipe step product %s to be %v, but was %v", expected.ID, expected.StorageInstructions, actual.StorageInstructions)
	assert.Equal(t, expected.IsLiquid, actual.IsLiquid, "expected IsLiquid for recipe step product %s to be %v, but was %v", expected.ID, expected.IsLiquid, actual.IsLiquid)
	assert.Equal(t, expected.IsWaste, actual.IsWaste, "expected IsWaste for recipe step product %s to be %v, but was %v", expected.ID, expected.IsWaste, actual.IsWaste)
	assert.Equal(t, expected.Index, actual.Index, "expected Index for recipe step product %s to be %v, but was %v", expected.ID, expected.Index, actual.Index)
	assert.Equal(t, expected.ContainedInVesselIndex, actual.ContainedInVesselIndex, "expected ContainedInVesselIndex for recipe step product %s to be %v, but was %v", expected.ID, expected.ContainedInVesselIndex, actual.ContainedInVesselIndex)
	assert.NotZero(t, actual.CreatedAt)
}

func (s *TestSuite) TestRecipeStepProducts_CompleteLifecycle() {
	s.runTest("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)

			var (
				createdRecipeStepID string
			)
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.adminClient.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.adminClient.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStepID
			exampleRecipeStepProduct.MeasurementUnit = createdValidMeasurementUnit
			exampleRecipeStepProductInput := converters.ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := testClients.adminClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, exampleRecipeStepProductInput)
			require.NoError(t, err)

			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

			createdRecipeStepProduct, err = testClients.userClient.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID)
			requireNotNilAndNoProblems(t, createdRecipeStepProduct, err)
			require.Equal(t, createdRecipeStepID, createdRecipeStepProduct.BelongsToRecipeStep)

			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

			newRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			newRecipeStepProduct.MeasurementUnit = createdValidMeasurementUnit
			updateInput := converters.ConvertRecipeStepProductToRecipeStepProductUpdateRequestInput(newRecipeStepProduct)
			createdRecipeStepProduct.Update(updateInput)

			require.NoError(t, testClients.adminClient.UpdateRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID, updateInput))

			actual, err := testClients.userClient.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step product equality
			checkRecipeStepProductEquality(t, newRecipeStepProduct, actual)
			assert.NotNil(t, actual.LastUpdatedAt)

			assert.NoError(t, testClients.userClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID))

			assert.NoError(t, testClients.userClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_Listing() {
	s.runTest("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.adminClient, testClients.userClient, nil)

			var (
				createdRecipeStepID string
			)
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := converters.ConvertValidMeasurementUnitToValidMeasurementUnitCreationRequestInput(exampleValidMeasurementUnit)
			createdValidMeasurementUnit, err := testClients.adminClient.CreateValidMeasurementUnit(ctx, exampleValidMeasurementUnitInput)
			require.NoError(t, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			createdValidMeasurementUnit, err = testClients.adminClient.GetValidMeasurementUnit(ctx, createdValidMeasurementUnit.ID)
			requireNotNilAndNoProblems(t, createdValidMeasurementUnit, err)
			checkValidMeasurementUnitEquality(t, exampleValidMeasurementUnit, createdValidMeasurementUnit)

			var expected []*types.RecipeStepProduct
			for i := 0; i < 5; i++ {
				exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
				exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepProduct.MeasurementUnit = createdValidMeasurementUnit
				exampleRecipeStepProductInput := converters.ConvertRecipeStepProductToRecipeStepProductCreationRequestInput(exampleRecipeStepProduct)
				createdRecipeStepProduct, createdRecipeStepProductErr := testClients.adminClient.CreateRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, exampleRecipeStepProductInput)
				require.NoError(t, createdRecipeStepProductErr)

				checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

				createdRecipeStepProduct, createdRecipeStepProductErr = testClients.userClient.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID)
				requireNotNilAndNoProblems(t, createdRecipeStepProduct, createdRecipeStepProductErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepProduct.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepProduct)
			}

			// assert recipe step product list equality
			actual, err := testClients.userClient.GetRecipeStepProducts(ctx, createdRecipe.ID, createdRecipeStepID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Data),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Data),
			)

			for _, createdRecipeStepProduct := range expected {
				assert.NoError(t, testClients.userClient.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID))
			}

			assert.NoError(t, testClients.userClient.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			assert.NoError(t, testClients.adminClient.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
