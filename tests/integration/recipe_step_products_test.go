package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types/fakes"

	"github.com/prixfixeco/api_server/pkg/types"
)

func checkRecipeStepProductEquality(t *testing.T, expected, actual *types.RecipeStepProduct) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe step product %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Type, actual.Type, "expected Type for recipe step product %s to be %v, but it was %v", expected.ID, expected.Type, actual.Type)
	assert.Equal(t, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID, "expected MeasurementUnit.ID for recipe step product %s to be %v, but it was %v", expected.ID, expected.MeasurementUnit.ID, actual.MeasurementUnit.ID)
	assert.Equal(t, expected.MinimumQuantityValue, actual.MinimumQuantityValue, "expected MinimumQuantityValue for recipe step product %s to be %v, but it was %v", expected.ID, expected.MinimumQuantityValue, actual.MinimumQuantityValue)
	assert.Equal(t, expected.MaximumQuantityValue, actual.MaximumQuantityValue, "expected MaximumQuantityValue for recipe step product %s to be %v, but it was %v", expected.ID, expected.MaximumQuantityValue, actual.MaximumQuantityValue)
	assert.Equal(t, expected.QuantityNotes, actual.QuantityNotes, "expected QuantityNotes for recipe step product %s to be %v, but it was %v", expected.ID, expected.QuantityNotes, actual.QuantityNotes)
	assert.NotZero(t, actual.CreatedOn)
}

// convertRecipeStepProductToRecipeStepProductUpdateInput creates an RecipeStepProductUpdateRequestInput struct from a recipe step product.
func convertRecipeStepProductToRecipeStepProductUpdateInput(x *types.RecipeStepProduct) *types.RecipeStepProductUpdateRequestInput {
	return &types.RecipeStepProductUpdateRequestInput{
		Name:                 &x.Name,
		Type:                 &x.Type,
		MeasurementUnitID:    &x.MeasurementUnit.ID,
		MinimumQuantityValue: &x.MinimumQuantityValue,
		MaximumQuantityValue: &x.MaximumQuantityValue,
		QuantityNotes:        &x.QuantityNotes,
	}
}

func (s *TestSuite) TestRecipeStepProducts_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.main, nil)

			var (
				createdRecipeStepID string
			)
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating valid measurement unit")
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(exampleValidMeasurementUnit)
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
			exampleRecipeStepProduct.MeasurementUnit = createdValidMeasurementUnit
			exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			require.NoError(t, err)
			t.Logf("recipe step product %q created", createdRecipeStepProduct.ID)

			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

			createdRecipeStepProduct, err = testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID)
			requireNotNilAndNoProblems(t, createdRecipeStepProduct, err)
			require.Equal(t, createdRecipeStepID, createdRecipeStepProduct.BelongsToRecipeStep)

			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

			t.Log("changing recipe step product")
			newRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			newRecipeStepProduct.MeasurementUnit = createdValidMeasurementUnit
			createdRecipeStepProduct.Update(convertRecipeStepProductToRecipeStepProductUpdateInput(newRecipeStepProduct))

			require.NoError(t, testClients.main.UpdateRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepProduct))

			t.Log("fetching changed recipe step product")
			actual, err := testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step product equality
			checkRecipeStepProductEquality(t, newRecipeStepProduct, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step product")
			assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_Listing() {
	s.runForEachClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.main, nil)

			var (
				createdRecipeStepID string
			)
			for _, step := range createdRecipe.Steps {
				createdRecipeStepID = step.ID
				break
			}

			t.Log("creating valid measurement unit")
			exampleValidMeasurementUnit := fakes.BuildFakeValidMeasurementUnit()
			exampleValidMeasurementUnitInput := fakes.BuildFakeValidMeasurementUnitCreationRequestInputFromValidMeasurementUnit(exampleValidMeasurementUnit)
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
				exampleRecipeStepProduct.MeasurementUnit = createdValidMeasurementUnit
				exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(exampleRecipeStepProduct)
				createdRecipeStepProduct, createdRecipeStepProductErr := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
				require.NoError(t, createdRecipeStepProductErr)
				t.Logf("recipe step product %q created", createdRecipeStepProduct.ID)

				checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

				createdRecipeStepProduct, createdRecipeStepProductErr = testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID)
				requireNotNilAndNoProblems(t, createdRecipeStepProduct, createdRecipeStepProductErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepProduct.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepProduct)
			}

			// assert recipe step product list equality
			actual, err := testClients.main.GetRecipeStepProducts(ctx, createdRecipe.ID, createdRecipeStepID, nil)
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
				assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepProduct.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
