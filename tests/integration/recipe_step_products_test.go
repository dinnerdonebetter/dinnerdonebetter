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

func checkRecipeStepProductEquality(t *testing.T, expected, actual *types.RecipeStepProduct) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe step product #%d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.QuantityType, actual.QuantityType, "expected QuantityType for recipe step product #%d to be %v, but it was %v ", expected.ID, expected.QuantityType, actual.QuantityType)
	assert.Equal(t, expected.QuantityValue, actual.QuantityValue, "expected QuantityValue for recipe step product #%d to be %v, but it was %v ", expected.ID, expected.QuantityValue, actual.QuantityValue)
	assert.Equal(t, expected.QuantityNotes, actual.QuantityNotes, "expected QuantityNotes for recipe step product #%d to be %v, but it was %v ", expected.ID, expected.QuantityNotes, actual.QuantityNotes)
	assert.Equal(t, expected.RecipeStepID, actual.RecipeStepID, "expected RecipeStepID for recipe step product #%d to be %v, but it was %v ", expected.ID, expected.RecipeStepID, actual.RecipeStepID)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestRecipeStepProducts_Creating() {
	s.runForEachClientExcept("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// Create recipe step product.
			exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			requireNotNilAndNoProblems(t, createdRecipeStepProduct, err)

			// assert recipe step product equality
			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, createdRecipeStepProduct)

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeStepProductCreationEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipeStepProduct.ID, audit.RecipeStepProductAssignmentKey)

			// Clean up recipe step product.
			assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_Listing() {
	s.runForEachClientExcept("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// create recipe step products
			var expected []*types.RecipeStepProduct
			for i := 0; i < 5; i++ {
				exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
				exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
				exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)

				createdRecipeStepProduct, recipeStepProductCreationErr := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
				requireNotNilAndNoProblems(t, createdRecipeStepProduct, recipeStepProductCreationErr)

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

			// clean up
			for _, createdRecipeStepProduct := range actual.RecipeStepProducts {
				assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))
			}

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_ExistenceChecking_ReturnsFalseForNonexistentRecipeStepProduct() {
	s.runForEachClientExcept("should not return an error for nonexistent recipe step product", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			actual, err := testClients.main.RecipeStepProductExists(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_ExistenceChecking_ReturnsTrueForValidRecipeStepProduct() {
	s.runForEachClientExcept("should not return an error for existent recipe step product", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// create recipe step product
			exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			requireNotNilAndNoProblems(t, createdRecipeStepProduct, err)

			// retrieve recipe step product
			actual, err := testClients.main.RecipeStepProductExists(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// clean up recipe step product
			assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_Reading_Returns404ForNonexistentRecipeStepProduct() {
	s.runForEachClientExcept("it should return an error when trying to read a recipe step product that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			_, err = testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)
			assert.Error(t, err)

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_Reading() {
	s.runForEachClientExcept("it should be readable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// create recipe step product
			exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			requireNotNilAndNoProblems(t, createdRecipeStepProduct, err)

			// retrieve recipe step product
			actual, err := testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step product equality
			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, actual)

			// clean up recipe step product
			assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_Updating_Returns404ForNonexistentRecipeStepProduct() {
	s.runForEachClientExcept("it should return an error when trying to update something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProduct))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

// convertRecipeStepProductToRecipeStepProductUpdateInput creates an RecipeStepProductUpdateInput struct from a recipe step product.
func convertRecipeStepProductToRecipeStepProductUpdateInput(x *types.RecipeStepProduct) *types.RecipeStepProductUpdateInput {
	return &types.RecipeStepProductUpdateInput{
		Name:          x.Name,
		QuantityType:  x.QuantityType,
		QuantityValue: x.QuantityValue,
		QuantityNotes: x.QuantityNotes,
		RecipeStepID:  x.RecipeStepID,
	}
}

func (s *TestSuite) TestRecipeStepProducts_Updating() {
	s.runForEachClientExcept("it should be possible to update a recipe step product", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// create recipe step product
			exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			requireNotNilAndNoProblems(t, createdRecipeStepProduct, err)

			// change recipe step product
			createdRecipeStepProduct.Update(convertRecipeStepProductToRecipeStepProductUpdateInput(exampleRecipeStepProduct))
			assert.NoError(t, testClients.main.UpdateRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStepProduct))

			// retrieve changed recipe step product
			actual, err := testClients.main.GetRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step product equality
			checkRecipeStepProductEquality(t, exampleRecipeStepProduct, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeStepProductCreationEvent},
				{EventType: audit.RecipeStepProductUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipeStepProduct.ID, audit.RecipeStepProductAssignmentKey)

			// clean up recipe step product
			assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_Archiving_Returns404ForNonexistentRecipeStepProduct() {
	s.runForEachClientExcept("it should return an error when trying to delete something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			assert.Error(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID))

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_Archiving() {
	s.runForEachClientExcept("it should be possible to delete a recipe step product", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// create recipe step product
			exampleRecipeStepProduct := fakes.BuildFakeRecipeStepProduct()
			exampleRecipeStepProduct.BelongsToRecipeStep = createdRecipeStep.ID
			exampleRecipeStepProductInput := fakes.BuildFakeRecipeStepProductCreationInputFromRecipeStepProduct(exampleRecipeStepProduct)
			createdRecipeStepProduct, err := testClients.main.CreateRecipeStepProduct(ctx, createdRecipe.ID, exampleRecipeStepProductInput)
			requireNotNilAndNoProblems(t, createdRecipeStepProduct, err)

			// clean up recipe step product
			assert.NoError(t, testClients.main.ArchiveRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, createdRecipeStepProduct.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeStepProductCreationEvent},
				{EventType: audit.RecipeStepProductArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipeStepProduct.ID, audit.RecipeStepProductAssignmentKey)

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeStepProducts_Auditing_Returns404ForNonexistentRecipeStepProduct() {
	s.runForEachClientExcept("it should return an error when trying to audit something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// Create recipe step.
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			x, err := testClients.admin.GetAuditLogForRecipeStepProduct(ctx, createdRecipe.ID, createdRecipeStep.ID, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
