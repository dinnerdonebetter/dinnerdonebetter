package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
)

func checkRecipeStepIngredientEquality(t *testing.T, expected, actual *types.RecipeStepIngredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, *expected.IngredientID, *actual.IngredientID, "expected IngredientID for recipe step ingredient %s to be %v, but it was %v", expected.ID, *expected.IngredientID, *actual.IngredientID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.QuantityType, actual.QuantityType, "expected ValidMeasurementID for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.QuantityType, actual.QuantityType)
	assert.Equal(t, expected.MinimumQuantityValue, actual.MinimumQuantityValue, "expected MinimumQuantityValue for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.MinimumQuantityValue, actual.MinimumQuantityValue)
	assert.Equal(t, expected.QuantityNotes, actual.QuantityNotes, "expected QuantityNotes for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.QuantityNotes, actual.QuantityNotes)
	assert.Equal(t, expected.ProductOfRecipeStep, actual.ProductOfRecipeStep, "expected ProductOfRecipeStep for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.ProductOfRecipeStep, actual.ProductOfRecipeStep)
	assert.Equal(t, expected.IngredientNotes, actual.IngredientNotes, "expected IngredientNotes for recipe step ingredient %s to be %v, but it was %v", expected.ID, expected.IngredientNotes, actual.IngredientNotes)
	assert.NotZero(t, actual.CreatedOn)
}

// convertRecipeStepIngredientToRecipeStepIngredientUpdateInput creates an RecipeStepIngredientUpdateRequestInput struct from a recipe step ingredient.
func convertRecipeStepIngredientToRecipeStepIngredientUpdateInput(x *types.RecipeStepIngredient) *types.RecipeStepIngredientUpdateRequestInput {
	return &types.RecipeStepIngredientUpdateRequestInput{
		IngredientID:         x.IngredientID,
		Name:                 &x.Name,
		ValidMeasurementID:   &x.QuantityType,
		MinimumQuantityValue: &x.MinimumQuantityValue,
		QuantityNotes:        &x.QuantityNotes,
		ProductOfRecipeStep:  &x.ProductOfRecipeStep,
		IngredientNotes:      &x.IngredientNotes,
		RecipeStepProductID:  x.RecipeStepProductID,
		BelongsToRecipeStep:  &x.BelongsToRecipeStep,
	}
}

func (s *TestSuite) TestRecipeStepIngredients_CompleteLifecycle() {
	s.runForEachClient("should be creatable and readable and updatable and deletable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, _, createdRecipe := createRecipeForTest(ctx, t, testClients.admin, testClients.main, nil)

			firstStep := createdRecipe.Steps[0]
			createdRecipeStepID := firstStep.ID
			createdRecipeStepIngredientID := firstStep.Ingredients[0].ID

			require.NotEmpty(t, createdRecipeStepID, "created recipe step ID must not be empty")
			require.NotEmpty(t, createdRecipeStepIngredientID, "created recipe step ingredient ID must not be empty")

			t.Log("fetching recipe step ingredient")
			createdRecipeStepIngredient, err := testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID)
			requireNotNilAndNoProblems(t, createdRecipeStepIngredient, err)

			t.Log("creating new valid ingredient for recipe step ingredient")
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationRequestInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.admin.CreateValidIngredient(ctx, exampleValidIngredientInput)
			require.NoError(t, err)
			t.Logf("valid ingredient %q created", createdValidIngredient.ID)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			createdValidIngredient, err = testClients.admin.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			t.Log("changing recipe step ingredient")
			newRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
			newRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStepID
			newRecipeStepIngredient.IngredientID = &createdValidIngredient.ID
			newRecipeStepIngredient.ID = createdRecipeStepIngredientID

			createdRecipeStepIngredient.Update(convertRecipeStepIngredientToRecipeStepIngredientUpdateInput(newRecipeStepIngredient))

			t.Logf("updating recipe step ingredient: %+v", createdRecipeStepIngredient)

			require.NoError(t, testClients.main.UpdateRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepIngredient))

			t.Log("fetching changed recipe step ingredient")
			actual, err := testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step ingredient equality
			checkRecipeStepIngredientEquality(t, newRecipeStepIngredient, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			t.Log("cleaning up recipe step ingredient")
			assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredientID))

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_Listing() {
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

			t.Log("creating recipe step ingredients")
			var expected []*types.RecipeStepIngredient
			for i := 0; i < 5; i++ {
				x, _, _ := createRecipeForTest(ctx, t, testClients.admin, testClients.main, nil)

				exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
				exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStepID
				exampleRecipeStepIngredient.IngredientID = &x[0].ID

				exampleRecipeStepIngredientInput := fakes.BuildFakeRecipeStepIngredientCreationRequestInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
				createdRecipeStepIngredient, createdRecipeStepIngredientErr := testClients.main.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
				require.NoError(t, createdRecipeStepIngredientErr)

				t.Logf("recipe step ingredient %q created", createdRecipeStepIngredient.ID)
				checkRecipeStepIngredientEquality(t, exampleRecipeStepIngredient, createdRecipeStepIngredient)

				createdRecipeStepIngredient, createdRecipeStepIngredientErr = testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredient.ID)
				requireNotNilAndNoProblems(t, createdRecipeStepIngredient, createdRecipeStepIngredientErr)
				require.Equal(t, createdRecipeStepID, createdRecipeStepIngredient.BelongsToRecipeStep)

				expected = append(expected, createdRecipeStepIngredient)
			}

			// assert recipe step ingredient list equality
			actual, err := testClients.main.GetRecipeStepIngredients(ctx, createdRecipe.ID, createdRecipeStepID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepIngredients),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepIngredients),
			)

			t.Log("cleaning up")
			for _, createdRecipeStepIngredient := range expected {
				assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepID, createdRecipeStepIngredient.ID))
			}

			t.Log("cleaning up recipe step")
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStepID))

			t.Log("cleaning up recipe")
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
