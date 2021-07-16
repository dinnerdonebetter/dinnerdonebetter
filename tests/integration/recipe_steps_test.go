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

func checkRecipeStepEquality(t *testing.T, expected, actual *types.RecipeStep) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Index, actual.Index, "expected Index for recipe step #%d to be %v, but it was %v ", expected.ID, expected.Index, actual.Index)
	assert.Equal(t, expected.PreparationID, actual.PreparationID, "expected PreparationID for recipe step #%d to be %v, but it was %v ", expected.ID, expected.PreparationID, actual.PreparationID)
	assert.Equal(t, expected.PrerequisiteStep, actual.PrerequisiteStep, "expected PrerequisiteStep for recipe step #%d to be %v, but it was %v ", expected.ID, expected.PrerequisiteStep, actual.PrerequisiteStep)
	assert.Equal(t, expected.MinEstimatedTimeInSeconds, actual.MinEstimatedTimeInSeconds, "expected MinEstimatedTimeInSeconds for recipe step #%d to be %v, but it was %v ", expected.ID, expected.MinEstimatedTimeInSeconds, actual.MinEstimatedTimeInSeconds)
	assert.Equal(t, expected.MaxEstimatedTimeInSeconds, actual.MaxEstimatedTimeInSeconds, "expected MaxEstimatedTimeInSeconds for recipe step #%d to be %v, but it was %v ", expected.ID, expected.MaxEstimatedTimeInSeconds, actual.MaxEstimatedTimeInSeconds)
	assert.Equal(t, expected.TemperatureInCelsius, actual.TemperatureInCelsius, "expected TemperatureInCelsius for recipe step #%d to be %v, but it was %v ", expected.ID, expected.TemperatureInCelsius, actual.TemperatureInCelsius)
	assert.Equal(t, expected.Notes, actual.Notes, "expected Notes for recipe step #%d to be %v, but it was %v ", expected.ID, expected.Notes, actual.Notes)
	assert.Equal(t, expected.Why, actual.Why, "expected Why for recipe step #%d to be %v, but it was %v ", expected.ID, expected.Why, actual.Why)
	assert.Equal(t, expected.RecipeID, actual.RecipeID, "expected RecipeID for recipe step #%d to be %v, but it was %v ", expected.ID, expected.RecipeID, actual.RecipeID)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestRecipeSteps_Creating() {
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

			// assert recipe step equality
			checkRecipeStepEquality(t, exampleRecipeStep, createdRecipeStep)

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeStepCreationEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipeStep.ID, audit.RecipeStepAssignmentKey)

			// Clean up recipe step.
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_Listing() {
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

			// create recipe steps
			var expected []*types.RecipeStep
			for i := 0; i < 5; i++ {
				exampleRecipeStep := fakes.BuildFakeRecipeStep()
				exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
				exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)

				createdRecipeStep, recipeStepCreationErr := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
				requireNotNilAndNoProblems(t, createdRecipeStep, recipeStepCreationErr)

				expected = append(expected, createdRecipeStep)
			}

			// assert recipe step list equality
			actual, err := testClients.main.GetRecipeSteps(ctx, createdRecipe.ID, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.RecipeSteps),
				"expected %d to be <= %d",
				len(expected),
				len(actual.RecipeSteps),
			)

			// clean up
			for _, createdRecipeStep := range actual.RecipeSteps {
				assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))
			}

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_ExistenceChecking_ReturnsFalseForNonexistentRecipeStep() {
	s.runForEachClientExcept("should not return an error for nonexistent recipe step", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			actual, err := testClients.main.RecipeStepExists(ctx, createdRecipe.ID, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_ExistenceChecking_ReturnsTrueForValidRecipeStep() {
	s.runForEachClientExcept("should not return an error for existent recipe step", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// create recipe step
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// retrieve recipe step
			actual, err := testClients.main.RecipeStepExists(ctx, createdRecipe.ID, createdRecipeStep.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// clean up recipe step
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_Reading_Returns404ForNonexistentRecipeStep() {
	s.runForEachClientExcept("it should return an error when trying to read a recipe step that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			_, err = testClients.main.GetRecipeStep(ctx, createdRecipe.ID, nonexistentID)
			assert.Error(t, err)

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_Reading() {
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

			// create recipe step
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// retrieve recipe step
			actual, err := testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step equality
			checkRecipeStepEquality(t, exampleRecipeStep, actual)

			// clean up recipe step
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_Updating_Returns404ForNonexistentRecipeStep() {
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

			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateRecipeStep(ctx, exampleRecipeStep))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

// convertRecipeStepToRecipeStepUpdateInput creates an RecipeStepUpdateInput struct from a recipe step.
func convertRecipeStepToRecipeStepUpdateInput(x *types.RecipeStep) *types.RecipeStepUpdateInput {
	return &types.RecipeStepUpdateInput{
		Index:                     x.Index,
		PreparationID:             x.PreparationID,
		PrerequisiteStep:          x.PrerequisiteStep,
		MinEstimatedTimeInSeconds: x.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: x.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      x.TemperatureInCelsius,
		Notes:                     x.Notes,
		Why:                       x.Why,
		RecipeID:                  x.RecipeID,
	}
}

func (s *TestSuite) TestRecipeSteps_Updating() {
	s.runForEachClientExcept("it should be possible to update a recipe step", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// create recipe step
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// change recipe step
			createdRecipeStep.Update(convertRecipeStepToRecipeStepUpdateInput(exampleRecipeStep))
			assert.NoError(t, testClients.main.UpdateRecipeStep(ctx, createdRecipeStep))

			// retrieve changed recipe step
			actual, err := testClients.main.GetRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe step equality
			checkRecipeStepEquality(t, exampleRecipeStep, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeStepCreationEvent},
				{EventType: audit.RecipeStepUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipeStep.ID, audit.RecipeStepAssignmentKey)

			// clean up recipe step
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_Archiving_Returns404ForNonexistentRecipeStep() {
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

			assert.Error(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, nonexistentID))

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_Archiving() {
	s.runForEachClientExcept("it should be possible to delete a recipe step", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// create recipe step
			exampleRecipeStep := fakes.BuildFakeRecipeStep()
			exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
			exampleRecipeStepInput := fakes.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
			createdRecipeStep, err := testClients.main.CreateRecipeStep(ctx, exampleRecipeStepInput)
			requireNotNilAndNoProblems(t, createdRecipeStep, err)

			// clean up recipe step
			assert.NoError(t, testClients.main.ArchiveRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForRecipeStep(ctx, createdRecipe.ID, createdRecipeStep.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.RecipeStepCreationEvent},
				{EventType: audit.RecipeStepArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdRecipeStep.ID, audit.RecipeStepAssignmentKey)

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}

func (s *TestSuite) TestRecipeSteps_Auditing_Returns404ForNonexistentRecipeStep() {
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

			x, err := testClients.admin.GetAuditLogForRecipeStep(ctx, createdRecipe.ID, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)

			// Clean up recipe.
			assert.NoError(t, testClients.main.ArchiveRecipe(ctx, createdRecipe.ID))
		}
	})
}
