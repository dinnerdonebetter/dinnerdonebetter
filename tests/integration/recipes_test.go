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

func checkRecipeEquality(t *testing.T, expected, actual *types.Recipe) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for recipe #%d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Source, actual.Source, "expected Source for recipe #%d to be %v, but it was %v ", expected.ID, expected.Source, actual.Source)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for recipe #%d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.InspiredByRecipeID, actual.InspiredByRecipeID, "expected InspiredByRecipeID for recipe #%d to be %v, but it was %v ", expected.ID, expected.InspiredByRecipeID, actual.InspiredByRecipeID)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestRecipes_Creating() {
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

			// assert recipe equality
			checkRecipeEquality(t, exampleRecipe, createdRecipe)

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipe(ctx, createdRecipe.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeCreationEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipe.ID, audit.RecipeAssignmentKey)

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipes_Listing() {
	s.runForEachClientExcept("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create recipes
			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				exampleRecipe := fakes.BuildFakeRecipe()
				exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)

				createdRecipe, recipeCreationErr := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
				requireNotNilAndNoProblems(t, createdRecipe, recipeCreationErr)

				expected = append(expected, createdRecipe)
			}

			// assert recipe list equality
			actual, err := testClients.main.GetRecipes(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.Recipes),
				"expected %d to be <= %d",
				len(expected),
				len(actual.Recipes),
			)

			// clean up
			for _, createdRecipe := range actual.Recipes {
				assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
			}
		}
	})
}

func (s *TestSuite) TestRecipes_ExistenceChecking_ReturnsFalseForNonexistentRecipe() {
	s.runForEachClientExcept("should not return an error for nonexistent recipe", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.main.RecipeExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		}
	})
}

func (s *TestSuite) TestRecipes_ExistenceChecking_ReturnsTrueForValidRecipe() {
	s.runForEachClientExcept("should not return an error for existent recipe", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create recipe
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// retrieve recipe
			actual, err := testClients.main.RecipeExists(ctx, createdRecipe.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// clean up recipe
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipes_Reading_Returns404ForNonexistentRecipe() {
	s.runForEachClientExcept("it should return an error when trying to read a recipe that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, err := testClients.main.GetRecipe(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestRecipes_Reading() {
	s.runForEachClientExcept("it should be readable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create recipe
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// retrieve recipe
			actual, err := testClients.main.GetRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			checkRecipeEquality(t, exampleRecipe, actual)

			// clean up recipe
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipes_Updating_Returns404ForNonexistentRecipe() {
	s.runForEachClientExcept("it should return an error when trying to update something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateRecipe(ctx, exampleRecipe))
		}
	})
}

// convertRecipeToRecipeUpdateInput creates an RecipeUpdateInput struct from a recipe.
func convertRecipeToRecipeUpdateInput(x *types.Recipe) *types.RecipeUpdateInput {
	return &types.RecipeUpdateInput{
		Name:               x.Name,
		Source:             x.Source,
		Description:        x.Description,
		InspiredByRecipeID: x.InspiredByRecipeID,
	}
}

func (s *TestSuite) TestRecipes_Updating() {
	s.runForEachClientExcept("it should be possible to update a recipe", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create recipe
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// change recipe
			createdRecipe.Update(convertRecipeToRecipeUpdateInput(exampleRecipe))
			assert.NoError(t, testClients.main.UpdateRecipe(ctx, createdRecipe))

			// retrieve changed recipe
			actual, err := testClients.main.GetRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			checkRecipeEquality(t, exampleRecipe, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipe(ctx, createdRecipe.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeCreationEvent},
				{EventType: audit.RecipeUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipe.ID, audit.RecipeAssignmentKey)

			// clean up recipe
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipes_Archiving_Returns404ForNonexistentRecipe() {
	s.runForEachClientExcept("it should return an error when trying to delete something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.main.ArchiveRecipe(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestRecipes_Archiving() {
	s.runForEachClientExcept("it should be possible to delete a recipe", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create recipe
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// clean up recipe
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipe(ctx, createdRecipe.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeCreationEvent},
				{EventType: audit.RecipeArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipe.ID, audit.RecipeAssignmentKey)
		}
	})
}

func (s *TestSuite) TestRecipes_Auditing_Returns404ForNonexistentRecipe() {
	s.runForEachClientExcept("it should return an error when trying to audit something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			x, err := testClients.admin.GetAuditLogForRecipe(ctx, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)
		}
	})
}
