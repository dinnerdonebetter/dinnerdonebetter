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

	require.Equal(t, len(expected.Steps), len(actual.Steps))
	for i := range expected.Steps {
		checkRecipeStepEquality(t, expected.Steps[i], actual.Steps[i])
	}

	assert.NotZero(t, actual.CreatedOn)
}

func up(u uint64) *uint64 {
	return &u
}

func fullRecipeToRecipe(s *types.FullRecipe) *types.Recipe {
	steps := []*types.RecipeStep{}
	for _, i := range s.Steps {
		steps = append(steps, fullRecipeStepToRecipeStep(i))
	}

	return &types.Recipe{
		LastUpdatedOn:      s.LastUpdatedOn,
		ArchivedOn:         s.ArchivedOn,
		InspiredByRecipeID: s.InspiredByRecipeID,
		Source:             s.Source,
		Description:        s.Description,
		ExternalID:         s.ExternalID,
		DisplayImageURL:    s.DisplayImageURL,
		Name:               s.Name,
		Steps:              steps,
		ID:                 s.ID,
		CreatedOn:          s.CreatedOn,
		BelongsToHousehold: s.BelongsToHousehold,
	}
}

func (s *TestSuite) TestRecipes_Creating() {
	s.runForEachClientExcept("should be creatable", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// Create recipe.
			exampleRecipe := fakes.BuildFakeRecipe()
			exampleRecipe.Steps = []*types.RecipeStep{}
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

			// create valid preparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// create valid ingredients
			var created []*types.ValidIngredient
			for i := 0; i < 5; i++ {
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

				createdValidIngredient, validIngredientCreationErr := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
				requireNotNilAndNoProblems(t, createdValidIngredient, validIngredientCreationErr)

				created = append(created, createdValidIngredient)
			}

			// create recipes
			var expected []*types.Recipe
			for i := 0; i < 5; i++ {
				exampleRecipe := fakes.BuildFakeRecipe()

				for j := range exampleRecipe.Steps {
					exampleRecipe.Steps[j].PreparationID = createdValidPreparation.ID
					for k := range exampleRecipe.Steps[j].Ingredients {
						exampleRecipe.Steps[j].Ingredients[k].IngredientID = up(created[k].ID)
					}
				}

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
			exampleRecipe.Steps = []*types.RecipeStep{}
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

			// create valid preparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// create valid ingredients
			var createdIngredients []*types.ValidIngredient
			for i := 0; i < 5; i++ {
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

				createdValidIngredient, validIngredientCreationErr := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
				requireNotNilAndNoProblems(t, createdValidIngredient, validIngredientCreationErr)

				createdIngredients = append(createdIngredients, createdValidIngredient)
			}

			// create recipe
			exampleRecipeInput := &types.RecipeCreationInput{
				Name:        t.Name(),
				Description: t.Name(),
				Steps: []*types.RecipeStepCreationInput{
					{
						Why:           t.Name(),
						PreparationID: createdValidPreparation.ID,
						Ingredients: []*types.RecipeStepIngredientCreationInput{
							{
								IngredientID:  up(createdIngredients[0].ID),
								QuantityType:  "grams",
								QuantityValue: 123,
							},
							{
								IngredientID:  up(createdIngredients[1].ID),
								QuantityType:  "grams",
								QuantityValue: 123,
							},
							{
								IngredientID:  up(createdIngredients[2].ID),
								QuantityType:  "grams",
								QuantityValue: 123,
							},
						},
					},
					{
						Why:           t.Name(),
						PreparationID: createdValidPreparation.ID,
						Ingredients: []*types.RecipeStepIngredientCreationInput{
							{
								IngredientID:  up(createdIngredients[0].ID),
								QuantityType:  "grams",
								QuantityValue: 123,
							},
							{
								IngredientID:  up(createdIngredients[1].ID),
								QuantityType:  "grams",
								QuantityValue: 123,
							},
							{
								IngredientID:  up(createdIngredients[2].ID),
								QuantityType:  "grams",
								QuantityValue: 123,
							},
						},
					},
					{
						Why:           t.Name(),
						PreparationID: createdValidPreparation.ID,
						Ingredients: []*types.RecipeStepIngredientCreationInput{
							{
								IngredientID:  up(createdIngredients[0].ID),
								QuantityType:  "grams",
								QuantityValue: 123,
							},
							{
								IngredientID:  up(createdIngredients[1].ID),
								QuantityType:  "grams",
								QuantityValue: 123,
							},
							{
								IngredientID:  up(createdIngredients[2].ID),
								QuantityType:  "grams",
								QuantityValue: 123,
							},
						},
					},
				},
			}

			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			// retrieve recipe
			actual, err := testClients.main.GetRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			assert.Equal(t, exampleRecipeInput.Name, actual.Name)
			assert.Equal(t, exampleRecipeInput.Description, actual.Description)
			assert.Equal(t, len(exampleRecipeInput.Steps), len(actual.Steps))
			for i, step := range actual.Steps {
				assert.Equal(t, exampleRecipeInput.Steps[i].Why, step.Why)
				assert.Equal(t, exampleRecipeInput.Steps[i].PreparationID, step.Preparation.ID)
				for j, ingredient := range step.Ingredients {
					assert.Equal(t, *exampleRecipeInput.Steps[i].Ingredients[j].IngredientID, ingredient.Ingredient.ID)
					assert.Equal(t, exampleRecipeInput.Steps[i].Ingredients[j].QuantityType, ingredient.QuantityType)
					assert.Equal(t, exampleRecipeInput.Steps[i].Ingredients[j].QuantityValue, ingredient.QuantityValue)
				}
			}
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
			exampleRecipe.Steps = []*types.RecipeStep{}
			exampleRecipe.ID = nonexistentID

			assert.Error(t, testClients.main.UpdateRecipe(ctx, exampleRecipe))
		}
	})
}

func (s *TestSuite) TestRecipes_Updating() {
	s.runForEachClientExcept("it should be possible to update a recipe", func(testClients *testClientWrapper) func() {
		return func() {
			t := s.T()

			ctx, span := tracing.StartCustomSpan(s.ctx, t.Name())
			defer span.End()

			// create valid preparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// create valid ingredients
			var createdIngredients []*types.ValidIngredient
			for i := 0; i < 5; i++ {
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

				createdValidIngredient, validIngredientCreationErr := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
				requireNotNilAndNoProblems(t, createdValidIngredient, validIngredientCreationErr)

				createdIngredients = append(createdIngredients, createdValidIngredient)
			}

			// create recipe
			exampleRecipe := fakes.BuildFakeRecipe()

			for i := range exampleRecipe.Steps {
				exampleRecipe.Steps[i].PreparationID = createdValidPreparation.ID
				for j := range exampleRecipe.Steps[i].Ingredients {
					exampleRecipe.Steps[i].Ingredients[j].IngredientID = up(createdIngredients[i].ID)
				}
			}

			exampleRecipeInput := fakes.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
			createdRecipe, err := testClients.main.CreateRecipe(ctx, exampleRecipeInput)
			requireNotNilAndNoProblems(t, createdRecipe, err)

			exampleRecipe.Description = t.Name()

			// change recipe
			createdRecipe.Update(&types.RecipeUpdateInput{
				Name:               exampleRecipe.Name,
				Source:             exampleRecipe.Source,
				Description:        exampleRecipe.Description,
				InspiredByRecipeID: exampleRecipe.InspiredByRecipeID,
			})
			assert.NoError(t, testClients.main.UpdateRecipe(ctx, createdRecipe))

			// retrieve changed recipe
			actual, err := testClients.main.GetRecipe(ctx, createdRecipe.ID)
			requireNotNilAndNoProblems(t, actual, err)

			// assert recipe equality
			checkRecipeEquality(t, exampleRecipe, fullRecipeToRecipe(actual))
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

			// create valid preparation
			exampleValidPreparation := fakes.BuildFakeValidPreparation()
			exampleValidPreparationInput := fakes.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
			createdValidPreparation, err := testClients.main.CreateValidPreparation(ctx, exampleValidPreparationInput)
			requireNotNilAndNoProblems(t, createdValidPreparation, err)

			// create valid ingredients
			var created []*types.ValidIngredient
			for i := 0; i < 5; i++ {
				exampleValidIngredient := fakes.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakes.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)

				createdValidIngredient, validIngredientCreationErr := testClients.main.CreateValidIngredient(ctx, exampleValidIngredientInput)
				requireNotNilAndNoProblems(t, createdValidIngredient, validIngredientCreationErr)

				created = append(created, createdValidIngredient)
			}

			// create recipe
			exampleRecipe := fakes.BuildFakeRecipe()

			for i := range exampleRecipe.Steps {
				exampleRecipe.Steps[i].PreparationID = createdValidPreparation.ID
				for j := range exampleRecipe.Steps[i].Ingredients {
					exampleRecipe.Steps[i].Ingredients[j].IngredientID = up(created[i].ID)
				}
			}

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
