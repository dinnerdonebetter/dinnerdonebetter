package integration

import (
	"fmt"
	"testing"

	audit "gitlab.com/prixfixe/prixfixe/internal/audit"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	"gitlab.com/prixfixe/prixfixe/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func checkValidIngredientEquality(t *testing.T, expected, actual *types.ValidIngredient) {
	t.Helper()

	assert.NotZero(t, actual.ID)
	assert.Equal(t, expected.Name, actual.Name, "expected Name for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.Name, actual.Name)
	assert.Equal(t, expected.Variant, actual.Variant, "expected Variant for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.Variant, actual.Variant)
	assert.Equal(t, expected.Description, actual.Description, "expected Description for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.Description, actual.Description)
	assert.Equal(t, expected.Warning, actual.Warning, "expected Warning for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.Warning, actual.Warning)
	assert.Equal(t, expected.ContainsEgg, actual.ContainsEgg, "expected ContainsEgg for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.ContainsEgg, actual.ContainsEgg)
	assert.Equal(t, expected.ContainsDairy, actual.ContainsDairy, "expected ContainsDairy for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.ContainsDairy, actual.ContainsDairy)
	assert.Equal(t, expected.ContainsPeanut, actual.ContainsPeanut, "expected ContainsPeanut for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.ContainsPeanut, actual.ContainsPeanut)
	assert.Equal(t, expected.ContainsTreeNut, actual.ContainsTreeNut, "expected ContainsTreeNut for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.ContainsTreeNut, actual.ContainsTreeNut)
	assert.Equal(t, expected.ContainsSoy, actual.ContainsSoy, "expected ContainsSoy for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.ContainsSoy, actual.ContainsSoy)
	assert.Equal(t, expected.ContainsWheat, actual.ContainsWheat, "expected ContainsWheat for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.ContainsWheat, actual.ContainsWheat)
	assert.Equal(t, expected.ContainsShellfish, actual.ContainsShellfish, "expected ContainsShellfish for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.ContainsShellfish, actual.ContainsShellfish)
	assert.Equal(t, expected.ContainsSesame, actual.ContainsSesame, "expected ContainsSesame for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.ContainsSesame, actual.ContainsSesame)
	assert.Equal(t, expected.ContainsFish, actual.ContainsFish, "expected ContainsFish for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.ContainsFish, actual.ContainsFish)
	assert.Equal(t, expected.ContainsGluten, actual.ContainsGluten, "expected ContainsGluten for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.ContainsGluten, actual.ContainsGluten)
	assert.Equal(t, expected.AnimalFlesh, actual.AnimalFlesh, "expected AnimalFlesh for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.AnimalFlesh, actual.AnimalFlesh)
	assert.Equal(t, expected.AnimalDerived, actual.AnimalDerived, "expected AnimalDerived for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.AnimalDerived, actual.AnimalDerived)
	assert.Equal(t, expected.Volumetric, actual.Volumetric, "expected Volumetric for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.Volumetric, actual.Volumetric)
	assert.Equal(t, expected.IconPath, actual.IconPath, "expected IconPath for valid ingredient #%d to be %v, but it was %v ", expected.ID, expected.IconPath, actual.IconPath)
	assert.NotZero(t, actual.CreatedOn)
}

func (s *TestSuite) TestValidIngredients_Creating() {
	s.runForEachClientExcept("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid ingredient.
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			// assert valid ingredient equality
			checkValidIngredientEquality(t, exampleValidIngredient, createdValidIngredient)

			auditLogEntries, err := testClients.admin.GetAuditLogForValidIngredient(ctx, createdValidIngredient.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidIngredientCreationEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidIngredient.ID, audit.ValidIngredientAssignmentKey)

			// Clean up valid ingredient.
			assert.NoError(t, testClients.main.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredients_Listing() {
	s.runForEachClientExcept("should be readable in paginated form", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid ingredients
			var expected []*types.ValidIngredient
			for i := 0; i < 5; i++ {
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

				createdValidIngredient, validIngredientCreationErr := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
				requireNotNilAndNoProblems(t, createdValidIngredient, validIngredientCreationErr)

				expected = append(expected, createdValidIngredient)
			}

			// assert valid ingredient list equality
			actual, err := testClients.main.GetValidIngredients(ctx, nil)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual.ValidIngredients),
				"expected %d to be <= %d",
				len(expected),
				len(actual.ValidIngredients),
			)

			// clean up
			for _, createdValidIngredient := range actual.ValidIngredients {
				assert.NoError(t, testClients.main.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredients_Searching() {
	s.runForEachClientExcept("should be able to be search for valid ingredients", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create valid preparation.
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// create valid ingredients
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			var expected []*types.ValidIngredient
			for i := 0; i < 5; i++ {
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
				exampleValidIngredientInput.Name = fmt.Sprintf("%s %d", exampleValidIngredientInput.Name, i)

				createdValidIngredient, validIngredientCreationErr := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
				requireNotNilAndNoProblems(t, createdValidIngredient, validIngredientCreationErr)

				expected = append(expected, createdValidIngredient)
			}

			exampleLimit := uint8(20)

			// assert valid ingredient list equality
			actual, err := testClients.main.SearchValidIngredients(ctx, exampleValidIngredient.Name, createdValidPreparation.ID, exampleLimit)
			requireNotNilAndNoProblems(t, actual, err)
			assert.True(
				t,
				len(expected) <= len(actual),
				"expected results length %d to be <= %d",
				len(expected),
				len(actual),
			)

			// Clean up valid preparation.
			assert.NoError(t, testClients.main.ArchiveValidPreparation(ctx, createdValidPreparation.ID))

			// clean up
			for _, createdValidIngredient := range expected {
				assert.NoError(t, testClients.main.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
			}
		}
	})
}

func (s *TestSuite) TestValidIngredients_ExistenceChecking_ReturnsFalseForNonexistentValidIngredient() {
	s.runForEachClientExcept("should not return an error for nonexistent valid ingredient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			actual, err := testClients.main.ValidIngredientExists(ctx, nonexistentID)
			assert.NoError(t, err)
			assert.False(t, actual)
		}
	})
}

func (s *TestSuite) TestValidIngredients_ExistenceChecking_ReturnsTrueForValidValidIngredient() {
	s.runForEachClientExcept("should not return an error for existent valid ingredient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid ingredient
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			// retrieve valid ingredient
			actual, err := testClients.main.ValidIngredientExists(ctx, createdValidIngredient.ID)
			assert.NoError(t, err)
			assert.True(t, actual)

			// clean up valid ingredient
			assert.NoError(t, testClients.main.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredients_Reading_Returns404ForNonexistentValidIngredient() {
	s.runForEachClientExcept("it should return an error when trying to read a valid ingredient that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			_, err := testClients.main.GetValidIngredient(ctx, nonexistentID)
			assert.Error(t, err)
		}
	})
}

func (s *TestSuite) TestValidIngredients_Reading() {
	s.runForEachClientExcept("it should be readable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid ingredient
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			// retrieve valid ingredient
			actual, err := testClients.main.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient equality
			checkValidIngredientEquality(t, exampleValidIngredient, actual)

			// clean up valid ingredient
			assert.NoError(t, testClients.main.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredients_Updating_Returns404ForNonexistentValidIngredient() {
	s.runForEachClientExcept("it should return an error when trying to update something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredient.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateValidIngredient(ctx, exampleValidIngredient))
		}
	})
}

// convertValidIngredientToValidIngredientUpdateInput creates an ValidIngredientUpdateInput struct from a valid ingredient.
func convertValidIngredientToValidIngredientUpdateInput(x *types.ValidIngredient) *types.ValidIngredientUpdateInput {
	return &types.ValidIngredientUpdateInput{
		Name:              x.Name,
		Variant:           x.Variant,
		Description:       x.Description,
		Warning:           x.Warning,
		ContainsEgg:       x.ContainsEgg,
		ContainsDairy:     x.ContainsDairy,
		ContainsPeanut:    x.ContainsPeanut,
		ContainsTreeNut:   x.ContainsTreeNut,
		ContainsSoy:       x.ContainsSoy,
		ContainsWheat:     x.ContainsWheat,
		ContainsShellfish: x.ContainsShellfish,
		ContainsSesame:    x.ContainsSesame,
		ContainsFish:      x.ContainsFish,
		ContainsGluten:    x.ContainsGluten,
		AnimalFlesh:       x.AnimalFlesh,
		AnimalDerived:     x.AnimalDerived,
		Volumetric:        x.Volumetric,
		IconPath:          x.IconPath,
	}
}

func (s *TestSuite) TestValidIngredients_Updating() {
	s.runForEachClientExcept("it should be possible to update a valid ingredient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid ingredient
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			// change valid ingredient
			createdValidIngredient.Update(convertValidIngredientToValidIngredientUpdateInput(exampleValidIngredient))
			assert.NoError(t, testClients.main.UpdateValidIngredient(ctx, createdValidIngredient))

			// retrieve changed valid ingredient
			actual, err := testClients.main.GetValidIngredient(ctx, createdValidIngredient.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert valid ingredient equality
			checkValidIngredientEquality(t, exampleValidIngredient, actual)
			assert.NotNil(t, actual.LastUpdatedOn)

			auditLogEntries, err := testClients.admin.GetAuditLogForValidIngredient(ctx, createdValidIngredient.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidIngredientCreationEvent},
				{EventType: audit.ValidIngredientUpdateEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidIngredient.ID, audit.ValidIngredientAssignmentKey)

			// clean up valid ingredient
			assert.NoError(t, testClients.main.ArchiveValidIngredient(ctx, createdValidIngredient.ID))
		}
	})
}

func (s *TestSuite) TestValidIngredients_Archiving_Returns404ForNonexistentValidIngredient() {
	s.runForEachClientExcept("it should return an error when trying to delete something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			assert.Error(t, testClients.main.ArchiveValidIngredient(ctx, nonexistentID))
		}
	})
}

func (s *TestSuite) TestValidIngredients_Archiving() {
	s.runForEachClientExcept("it should be possible to delete a valid ingredient", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid ingredient
			exampleValidIngredient := fakes.BuildFakeValidIngredient()
			exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
			createdValidIngredient, err := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
			requireNotNilAndNoProblems(t, createdValidIngredient, err)

			// clean up valid ingredient
			assert.NoError(t, testClients.main.ArchiveValidIngredient(ctx, createdValidIngredient.ID))

			auditLogEntries, err := testClients.admin.GetAuditLogForValidIngredient(ctx, createdValidIngredient.ID)
			require.NoError(t, err)

			expectedAuditLogEntries := []*types.AuditLogEntry{
				{EventType: audit.ValidIngredientCreationEvent},
				{EventType: audit.ValidIngredientArchiveEvent},
			}
			validateAuditLogEntries(t, expectedAuditLogEntries, auditLogEntries, createdValidIngredient.ID, audit.ValidIngredientAssignmentKey)
		}
	})
}

func (s *TestSuite) TestValidIngredients_Auditing_Returns404ForNonexistentValidIngredient() {
	s.runForEachClientExcept("it should return an error when trying to audit something that does not exist", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			x, err := testClients.admin.GetAuditLogForValidIngredient(ctx, nonexistentID)

			assert.NoError(t, err)
			assert.Empty(t, x)
		}
	})
}
