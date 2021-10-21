package integration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"
)

func checkRecipeStepProductEquality(t *testing.T, expected, actual *types.RecipeStepProduct) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe step product %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for recipe step product %s to be %v, but it was %v", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.NotZero(t, actual.CreatedOn)
}

// convertRecipeStepProductToRecipeStepProductUpdateInput creates an RecipeStepProductUpdateRequestInput struct from a recipe step product.
func convertRecipeStepProductToRecipeStepProductUpdateInput(x *types.RecipeStepProduct) *types.RecipeStepProductUpdateRequestInput {
	return &types.RecipeStepProductUpdateRequestInput{
		Name:         x.Name,
		RecipeStepID: x.RecipeStepID,
	}
}

func (s *TestSuite) TestRecipeStepProducts_CompleteLifecycle() {
	s.runForCookieClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", createdRecipeID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeDataType)
			require.NotNil(t, n.Recipe)
			checkRecipeEquality(t, exampleRecipe, n.Recipe)

			createdRecipe, err := testClients.main.GetRecipe(ctx, createdRecipeID)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			t.Log("creating prerequisite recipe step")
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStepID, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			require.NoError(t, err)
			t.Logf("recipe step %q created", createdRecipeStepID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepDataType)
			require.NotNil(t, n.RecipeStep)
			checkRecipeStepEquality(t, exampleRecipeStep, n.RecipeStep)

			createdRecipeStep, err := testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)
			require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)

			t.Log("creating recipe step product")
			exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProductID, err := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			require.NoError(t, err)
			t.Logf("recipe step product %q created", createdRecipeStepProductID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepProductDataType)
			require.NotNil(t, n.RecipeStepProduct)
			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, n.RecipeStepProduct)

			createdRecipeStepProduct, err := testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProductID)
			requireNotNilAndNoProblems(t, createdRecipeStepProduct, err)
			require.Equal(t, createdRecipeStep.ID, createdRecipeStepProduct.BelongsToRecipeStep)

			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

			t.Log("changing recipe step product")
			newRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			createdRecipeStepProduct.Update(convertRecipeStepProductToRecipeStepProductUpdateInput(newRecipeStepProduct))
			assert.NoError(t, testClients.main.UpdateRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepProduct))

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepProductDataType)

			t.Log("fetching changed recipe step product")
			actual, err := testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProductID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step product equality
			checkRecipeStepProductEquality(t, newRecipeStepProduct, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step product")
			assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProductID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})

	s.runForPASETOClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", createdRecipeID)

			var createdRecipe *types.Recipe
			checkFunc = func() bool {
				createdRecipe, err = testClients.main.GetRecipe(ctx, createdRecipeID)
				return assert.NotNil(t, createdRecipe) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkRecipeEquality(t, exampleRecipe, createdRecipe)

			t.Log("creating prerequisite recipe step")
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStepID, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			require.NoError(t, err)
			t.Logf("recipe step %q created", createdRecipeStepID)

			var createdRecipeStep *types.RecipeStep
			checkFunc = func() bool {
				createdRecipeStep, err = testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID)
				return assert.NotNil(t, createdRecipeStep) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)
			checkRecipeStepEquality(t, exampleRecipeStep, createdRecipeStep)

			t.Log("creating recipe step product")
			exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProductID, err := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			require.NoError(t, err)
			t.Logf("recipe step product %q created", createdRecipeStepProductID)

			var createdRecipeStepProduct *types.RecipeStepProduct
			checkFunc = func() bool {
				createdRecipeStepProduct, err = testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProductID)
				return assert.NotNil(t, createdRecipeStepProduct) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			require.Equal(t, createdRecipeStep.ID, createdRecipeStepProduct.BelongsToRecipeStep)
			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

			// assert recipe step product equality
			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

			// change recipe step product
			newRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			createdRecipeStepProduct.Update(convertRecipeStepProductToRecipeStepProductUpdateInput(newRecipeStepProduct))
			assert.NoError(t, testClients.main.UpdateRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepProduct))

			time.Sleep(time.Second)

			// retrieve changed recipe step product
			var actual *types.RecipeStepProduct
			checkFunc = func() bool {
				actual, err = testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProductID)
				return assert.NotNil(t, createdRecipeStepProduct) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)

			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step product equality
			checkRecipeStepProductEquality(t, newRecipeStepProduct, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step product")
			assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProductID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_Listing() {
	s.runForCookieClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			stopChan := make(chan bool, 1)
			notificationsChan, err := testClients.main.SubscribeToDataChangeNotifications(ctx, stopChan)
			require.NotNil(t, notificationsChan)
			require.NoError(t, err)

			var n *types.DataChangeMessage

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", createdRecipeID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeDataType)
			require.NotNil(t, n.Recipe)
			checkRecipeEquality(t, exampleRecipe, n.Recipe)

			createdRecipe, err := testClients.main.GetRecipe(ctx, createdRecipeID)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			t.Log("creating prerequisite recipe step")
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStepID, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			require.NoError(t, err)
			t.Logf("recipe step %q created", createdRecipeStepID)

			n = <-notificationsChan
			assert.Equal(t, n.DataType, types.RecipeStepDataType)
			require.NotNil(t, n.RecipeStep)
			checkRecipeStepEquality(t, exampleRecipeStep, n.RecipeStep)

			createdRecipeStep, err := testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)
			require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)

			t.Log("creating recipe step products")
			var expected []*types.RecipeStepProduct
			for i := 0; i < 5; i++ {
				exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
				exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
				exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(exampleRecipeStepProduct)
				createdRecipeStepProductID, createdRecipeStepProductErr := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
				require.NoError(t, createdRecipeStepProductErr)
				t.Logf("recipe step product %q created", createdRecipeStepProductID)

				n = <-notificationsChan
				assert.Equal(t, n.DataType, types.RecipeStepProductDataType)
				require.NotNil(t, n.RecipeStepProduct)
				checkRecipeStepProductEquality(t, exampleRecipeStepProduct, n.RecipeStepProduct)

				createdRecipeStepProduct, createdRecipeStepProductErr := testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProductID)
				requireNotNilAndNoProblems(t, createdRecipeStepProduct, createdRecipeStepProductErr)
				require.Equal(t, createdRecipeStep.ID, createdRecipeStepProduct.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepProduct)
			}

			// assert recipe step product list equality
			actual, err := testClients.main.GetRecipeStepProducts(ctx, createdRecipe.ID, createdRecipeStep.ID, nil)
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
				assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})

	s.runForPASETOClient("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			var checkFunc func() bool
			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			t.Log("creating prerequisite recipe")
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationRequestInputFromRecipe(exampleRecipe)
			createdRecipeID, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			require.NoError(t, err)
			t.Logf("recipe %q created", createdRecipeID)

			var createdRecipe *types.Recipe
			checkFunc = func() bool {
				createdRecipe, err = testClients.main.GetRecipe(ctx, createdRecipeID)
				return assert.NotNil(t, createdRecipe) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			checkRecipeEquality(t, exampleRecipe, createdRecipe)

			t.Log("creating prerequisite recipe step")
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationRequestInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStepID, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			require.NoError(t, err)
			t.Logf("recipe step %q created", createdRecipeStepID)

			var createdRecipeStep *types.RecipeStep
			checkFunc = func() bool {
				createdRecipeStep, err = testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID)
				return assert.NotNil(t, createdRecipeStep) && assert.NoError(t, err)
			}
			assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
			require.Equal(t, createdRecipe.ID, createdRecipeStep.BelongsToRecipe)
			checkRecipeStepEquality(t, exampleRecipeStep, createdRecipeStep)

			t.Log("creating recipe step products")
			var expected []*types.RecipeStepProduct
			for i := 0; i < 5; i++ {
				exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
				exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
				exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationRequestInputFromRecipeStepProduct(exampleRecipeStepProduct)
				createdRecipeStepProductID, createdRecipeStepProductErr := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
				require.NoError(t, createdRecipeStepProductErr)

				var createdRecipeStepProduct *types.RecipeStepProduct
				checkFunc = func() bool {
					createdRecipeStepProduct, createdRecipeStepProductErr = testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProductID)
					return assert.NotNil(t, createdRecipeStepProduct) && assert.NoError(t, createdRecipeStepProductErr)
				}
				assert.Eventually(t, checkFunc, creationTimeout, waitPeriod)
				checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

				expected = append(expected, createdRecipeStepProduct)
			}

			// assert recipe step product list equality
			actual, err := testClients.main.GetRecipeStepProducts(ctx, createdRecipe.ID, createdRecipeStep.ID, nil)
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
				assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipeID))
		}
	})
}
