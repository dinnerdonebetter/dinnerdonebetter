package integration

import (
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkRecipeStepIngredientEquality(t *testing.T, expected, actual *types.RecipeStepIngredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.IngredientID, actual.IngredientID, "expected IngredientID for recipe step ingredient #%d to be %v, but it was %v ", expected.ID, expected.IngredientID, actual.IngredientID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe step ingredient #%d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.QuantityType, actual.QuantityType, "expected QuantityType for recipe step ingredient #%d to be %v, but it was %v ", expected.ID, expected.QuantityType, actual.QuantityType)
	assert.Equal(t, expected.QuantityValue, actual.QuantityValue, "expected QuantityValue for recipe step ingredient #%d to be %v, but it was %v ", expected.ID, expected.QuantityValue, actual.QuantityValue)
	assert.Equal(t, expected.QuantityNotes, actual.QuantityNotes, "expected QuantityNotes for recipe step ingredient #%d to be %v, but it was %v ", expected.ID, expected.QuantityNotes, actual.QuantityNotes)
	assert.Equal(t, expected.ProductOfRecipeStep, actual.ProductOfRecipeStep, "expected ProductOfRecipeStep for recipe step ingredient #%d to be %v, but it was %v ", expected.ID, expected.ProductOfRecipeStep, actual.ProductOfRecipeStep)
	assert.Equal(t, expected.IngredientNotes, actual.IngredientNotes, "expected IngredientNotes for recipe step ingredient #%d to be %v, but it was %v ", expected.ID, expected.IngredientNotes, actual.IngredientNotes)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestRecipeStepIngredients_Creating() {
	s.runForEachClientExcept("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// Create valid ingredient.
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.PreparationID = createdValidPreparation.ID
			exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{}
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// Create recipe step ingredient.
			exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredient.IngredientID = &createdValidIngredient.ID
			exampleRecipeStepIngredientInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := testClients.main.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			requireNotNilAndNoProblems(t, createdRecipeStepIngredient, err)

			// assert recipe step ingredient equality
			checkRecipeStepIngredientEquality(t, exampleRecipeStepIngredient, createdRecipeStepIngredient)

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeStepIngredientCreationEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipeStepIngredient.ID, audit.RecipeStepIngredientAssignmentKey)

			// Clean up recipe step ingredient.
			assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_Listing() {
	s.runForEachClientExcept("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// Create valid ingredient.
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.PreparationID = createdValidPreparation.ID
			exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{}
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// create recipe step ingredients
			var expected []*types.RecipeStepIngredient
			for i := 0; i < 5; i++ {
				exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
				exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
				exampleRecipeStepIngredient.IngredientID = &createdValidIngredient.ID
				exampleRecipeStepIngredientInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)

				createdRecipeStepIngredient, recipeStepIngredientCreationErr := testClients.main.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
				requireNotNilAndNoProblems(t, createdRecipeStepIngredient, recipeStepIngredientCreationErr)

				expected = append(expected, createdRecipeStepIngredient)
			}

			// assert recipe step ingredient list equality
			actual, err := testClients.main.GetRecipeStepIngredients(ctx, createdRecipe.ID, createdRecipeStep.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeStepIngredients),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeStepIngredients),
			)

			// clean up
			for _, createdRecipeStepIngredient := range actual.RecipeStepIngredients {
				assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))
			}

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_ExistenceChecking_ReturnsFalseForNonexistentRecipeStepIngredient() {
	s.runForEachClientExcept("should not return an error for nonexistent recipe step ingredient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.PreparationID = createdValidPreparation.ID
			exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{}
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			actual, err := testClients.main.RecipeStepIngredientExists(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_ExistenceChecking_ReturnsTrueForValidRecipeStepIngredient() {
	s.runForEachClientExcept("should not return an error for existent recipe step ingredient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// Create valid ingredient.
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.PreparationID = createdValidPreparation.ID
			exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{}
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// create recipe step ingredient
			exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredient.IngredientID = &createdValidIngredient.ID
			exampleRecipeStepIngredientInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := testClients.main.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			requireNotNilAndNoProblems(t, createdRecipeStepIngredient, err)

			// retrieve recipe step ingredient
			actual, err := testClients.main.RecipeStepIngredientExists(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// clean up recipe step ingredient
			assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_Reading_Returns404ForNonexistentRecipeStepIngredient() {
	s.runForEachClientExcept("it should return an error when trying to read a recipe step ingredient that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.PreparationID = createdValidPreparation.ID
			exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{}
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			_, err = testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
			assert.Error(t, err)

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_Reading() {
	s.runForEachClientExcept("it should be readable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// Create valid ingredient.
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.PreparationID = createdValidPreparation.ID
			exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{}
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// create recipe step ingredient
			exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredient.IngredientID = &createdValidIngredient.ID
			exampleRecipeStepIngredientInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := testClients.main.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			requireNotNilAndNoProblems(t, createdRecipeStepIngredient, err)

			// retrieve recipe step ingredient
			actual, err := testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step ingredient equality
			checkRecipeStepIngredientEquality(t, exampleRecipeStepIngredient, actual)

			// clean up recipe step ingredient
			assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_Updating_Returns404ForNonexistentRecipeStepIngredient() {
	s.runForEachClientExcept("it should return an error when trying to update something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.PreparationID = createdValidPreparation.ID
			exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{}
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredient))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

// convertRecipeStepIngredientToRecipeStepIngredientUpdateInput creates an RecipeStepIngredientUpdateInput struct from a recipe step ingredient.
func convertRecipeStepIngredientToRecipeStepIngredientUpdateInput(x *types.RecipeStepIngredient) *types.RecipeStepIngredientUpdateInput {
	return &types.RecipeStepIngredientUpdateInput{
		IngredientID:        x.IngredientID,
		Name:                x.Name,
		QuantityType:        x.QuantityType,
		QuantityValue:       x.QuantityValue,
		QuantityNotes:       x.QuantityNotes,
		ProductOfRecipeStep: x.ProductOfRecipeStep,
		IngredientNotes:     x.IngredientNotes,
	}
}

func (s *TestSuite) TestRecipeStepIngredients_Updating() {
	s.runForEachClientExcept("it should be possible to update a recipe step ingredient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// Create valid ingredient.
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.PreparationID = createdValidPreparation.ID
			exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{}
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// create recipe step ingredient
			exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredient.IngredientID = &createdValidIngredient.ID
			exampleRecipeStepIngredientInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := testClients.main.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			requireNotNilAndNoProblems(t, createdRecipeStepIngredient, err)

			// change recipe step ingredient
			createdRecipeStepIngredient.Update(convertRecipeStepIngredientToRecipeStepIngredientUpdateInput(exampleRecipeStepIngredient))
			assert.NoError(t, testClients.main.UpdateRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStepIngredient))

			// retrieve changed recipe step ingredient
			actual, err := testClients.main.GetRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step ingredient equality
			checkRecipeStepIngredientEquality(t, exampleRecipeStepIngredient, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeStepIngredientCreationEvent},
				{EventType: audit.RecipeStepIngredientUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipeStepIngredient.ID, audit.RecipeStepIngredientAssignmentKey)

			// clean up recipe step ingredient
			assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_Archiving_Returns404ForNonexistentRecipeStepIngredient() {
	s.runForEachClientExcept("it should return an error when trying to delete something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.PreparationID = createdValidPreparation.ID
			exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{}
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			assert.Error(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_Archiving() {
	s.runForEachClientExcept("it should be possible to delete a recipe step ingredient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// Create valid ingredient.
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.PreparationID = createdValidPreparation.ID
			exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{}
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// create recipe step ingredient
			exampleRecipeStepIngredient := fakes.BuildFakeRecipeStepIngredient()
			exampleRecipeStepIngredient.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepIngredient.IngredientID = &createdValidIngredient.ID
			exampleRecipeStepIngredientInput := fakes.BuildFakeRecipeStepIngredientCreationInputFromRecipeStepIngredient(exampleRecipeStepIngredient)
			createdRecipeStepIngredient, err := testClients.main.CreateRecipeStepIngredient(ctx, createdRecipe.ID, exampleRecipeStepIngredientInput)
			requireNotNilAndNoProblems(t, createdRecipeStepIngredient, err)

			// clean up recipe step ingredient
			assert.NoError(t, testClients.main.ArchiveRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepIngredient.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeStepIngredientCreationEvent},
				{EventType: audit.RecipeStepIngredientArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipeStepIngredient.ID, audit.RecipeStepIngredientAssignmentKey)

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepIngredients_Auditing_Returns404ForNonexistentRecipeStepIngredient() {
	s.runForEachClientExcept("it should return an error when trying to audit something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.PreparationID = createdValidPreparation.ID
			exampleRecipeStep.Ingredients = []*types.RecipeStepIngredient{}
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			x, err := testClients.admin.GetAuditLogForRecipeStepIngredient(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
