package integration

import (
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	mpconverters "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	authgrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	identitygrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/identity"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/platform/types"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
	converters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"
	"github.com/dinnerdonebetter/backend/pkg/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMealPlans_WithRecipeOptionSelections(T *testing.T) {
	T.Parallel()

	T.Run("should create meal plan with alternative ingredients and verify grocery list", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		// Create the account admin user and their client
		_, accountAdminUserClient := createUserAndClientForTest(t)

		// Get the account ID
		currentStatus, statusErr := accountAdminUserClient.GetAuthStatus(ctx, &authgrpc.GetAuthStatusRequest{})
		require.NotNil(t, currentStatus)
		require.NoError(t, statusErr)
		relevantAccountID := currentStatus.ActiveAccount

		// Create 3 additional household members (total of 4 users including admin)
		householdClients := []client.Client{accountAdminUserClient}
		for i := 0; i < 3; i++ {
			u, c := createUserAndClientForTest(t)

			invitation, err := accountAdminUserClient.CreateAccountInvitation(ctx, &identitygrpc.CreateAccountInvitationRequest{
				Input: &identitygrpc.AccountInvitationCreationRequestInput{
					Note:    t.Name(),
					ToName:  t.Name(),
					ToEmail: u.EmailAddress,
				},
			})
			require.NoError(t, err)

			_, err = c.AcceptAccountInvitation(ctx, &identitygrpc.AcceptAccountInvitationRequest{
				AccountInvitationId: invitation.Created.Id,
				Input: &identitygrpc.AccountInvitationUpdateRequestInput{
					Token: invitation.Created.Token,
					Note:  t.Name(),
				},
			})
			require.NoError(t, err)

			_, err = c.SetDefaultAccount(ctx, &identitygrpc.SetDefaultAccountRequest{AccountId: relevantAccountID})
			require.NoError(t, err)

			tokenResponse, err := c.LoginForToken(ctx, &authgrpc.LoginForTokenRequest{Input: &authgrpc.UserLoginInput{
				Username:  u.Username,
				Password:  u.HashedPassword,
				TotpToken: generateTOTPCodeForUserForTest(t, u),
			}})
			require.NoError(t, err)
			assert.NotNil(t, tokenResponse)

			householdClients = append(householdClients, c)
		}
		require.Len(t, householdClients, 4)

		// Create Recipe 1 with alternative ingredients (option groups)
		// Step 1 will have: ingredient A (index=0, optionIndex=0) OR ingredient B (index=0, optionIndex=1)
		// Note: adminClient is used internally for recipe creation as it requires admin permissions
		recipe1ValidIngredients, recipe1ValidPreparation, recipe1 := createRecipeWithAlternativeIngredients(t, "Recipe1")
		require.NotNil(t, recipe1)
		require.NotEmpty(t, recipe1.Steps)
		require.GreaterOrEqual(t, len(recipe1ValidIngredients), 2)

		// Create Recipe 2 with alternative ingredients
		recipe2ValidIngredients, _, recipe2 := createRecipeWithAlternativeIngredients(t, "Recipe2")
		require.NotNil(t, recipe2)
		require.NotEmpty(t, recipe2.Steps)
		require.GreaterOrEqual(t, len(recipe2ValidIngredients), 2)

		// Create meals from the recipes using adminClient (meals can be created by admin)
		meal1 := createMealFromRecipe(t, recipe1, "Meal1")
		require.NotNil(t, meal1)

		meal2 := createMealFromRecipe(t, recipe2, "Meal2")
		require.NotNil(t, meal2)

		// Create a meal plan with both meals as options
		// We'll have all users vote for meal1
		now := time.Now()
		exampleMealPlan := &mealplanning.MealPlan{
			Notes:          t.Name(),
			Status:         string(mealplanning.MealPlanStatusAwaitingVotes),
			VotingDeadline: now.Add(10 * time.Second),
			ElectionMethod: mealplanning.MealPlanElectionMethodSchulze,
			Events: []*mealplanning.MealPlanEvent{
				{
					StartsAt: now.Add(24 * time.Hour),
					EndsAt:   now.Add(72 * time.Hour),
					MealName: mealplanning.BreakfastMealName,
					Options: []*mealplanning.MealPlanOption{
						{
							Meal:  mealplanning.Meal{ID: meal1.ID},
							Notes: "option A - meal with alternative ingredients",
						},
						{
							Meal:  mealplanning.Meal{ID: meal2.ID},
							Notes: "option B - another meal with alternative ingredients",
						},
					},
				},
			},
		}

		// Create the meal plan input with selections
		// We're selecting optionIndex=1 for ingredient index=0 in recipe1
		// This means we prefer ingredient B over ingredient A
		exampleMealPlanInput := mpconverters.ConvertMealPlanToMealPlanCreationRequestInput(exampleMealPlan)

		// Add selections - selecting the second alternative (optionIndex=1) for the first ingredient position (ingredientIndex=0)
		exampleMealPlanInput.Selections = []*mealplanning.MealPlanRecipeOptionSelectionCreationRequestInput{
			{
				RecipeID:            recipe1.ID,
				RecipeStepID:        recipe1.Steps[0].ID,
				IngredientIndex:     0,
				SelectedOptionIndex: 1, // Select the alternative ingredient (optionIndex=1)
				SelectionType:       mealplanning.MealPlanRecipeOptionSelectionTypeIngredient,
			},
		}

		// Create the meal plan
		createMealPlanRes, err := accountAdminUserClient.CreateMealPlan(ctx, &mealplanninggrpc.CreateMealPlanRequest{
			Input: converters.ConvertMealPlanCreationRequestInputToGRPCMealPlanCreationRequestInput(exampleMealPlanInput),
		})
		require.NoError(t, err)
		require.NotEmpty(t, createMealPlanRes.Created.Id)

		// Fetch the created meal plan
		createdMealPlanRes, err := accountAdminUserClient.GetMealPlan(ctx, &mealplanninggrpc.GetMealPlanRequest{
			MealPlanId: createMealPlanRes.Created.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, createdMealPlanRes)

		createdMealPlan := converters.ConvertGRPCMealPlanToMealPlan(createdMealPlanRes.Result)
		require.NotEmpty(t, createdMealPlan.Events)
		require.NotEmpty(t, createdMealPlan.Events[0].Options)

		createdMealPlanEvent := createdMealPlan.Events[0]

		// All 4 users vote for the same meal (option A - meal1)
		for _, userClient := range householdClients {
			userVotes := &mealplanning.MealPlanOptionVoteCreationRequestInput{
				Votes: []*mealplanning.MealPlanOptionVoteCreationInput{
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[0].ID,
						Rank:                    0, // First choice
					},
					{
						BelongsToMealPlanOption: createdMealPlanEvent.Options[1].ID,
						Rank:                    1, // Second choice
					},
				},
			}

			_, voteErr := userClient.CreateMealPlanOptionVote(ctx, &mealplanninggrpc.CreateMealPlanOptionVoteRequest{
				MealPlanId:      createdMealPlan.ID,
				MealPlanEventId: createdMealPlanEvent.ID,
				Input:           converters.ConvertMealPlanOptionVoteCreationRequestInputToGRPCMealPlanOptionVoteCreationRequestInput(userVotes),
			})
			require.NoError(t, voteErr)
		}

		// Force the voting deadline to have passed
		q := generated.New()
		rowsAffected, err := q.UpdateMealPlan(ctx, databaseClient.DB(), &generated.UpdateMealPlanParams{
			Notes:            createdMealPlan.Notes,
			Status:           generated.MealPlanStatus(createdMealPlan.Status),
			VotingDeadline:   time.Now().Add(-time.Minute),
			BelongsToAccount: relevantAccountID,
			ID:               createdMealPlan.ID,
		})
		require.NoError(t, err)
		require.NotZero(t, rowsAffected)

		// Run the finalization worker
		runFinalizeRes, err := adminClient.RunFinalizeMealPlanWorker(ctx, &mealplanninggrpc.RunFinalizeMealPlanWorkerRequest{})
		require.NoError(t, err)
		require.NotNil(t, runFinalizeRes)

		// Verify the meal plan is finalized
		finalizedMealPlanRes, err := accountAdminUserClient.GetMealPlan(ctx, &mealplanninggrpc.GetMealPlanRequest{
			MealPlanId: createdMealPlan.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, finalizedMealPlanRes)

		finalizedMealPlan := converters.ConvertGRPCMealPlanToMealPlan(finalizedMealPlanRes.Result)
		assert.Equal(t, string(mealplanning.MealPlanStatusFinalized), finalizedMealPlan.Status)

		// Verify the correct option was chosen (meal1 should be chosen since all users voted for it)
		var chosenOption *mealplanning.MealPlanOption
		for _, event := range finalizedMealPlan.Events {
			for _, opt := range event.Options {
				if opt.Chosen {
					chosenOption = opt
					break
				}
			}
		}
		require.NotNil(t, chosenOption, "expected one option to be chosen")
		assert.Equal(t, meal1.ID, chosenOption.Meal.ID, "expected meal1 to be chosen")

		// Run the grocery list initializer worker
		runGroceryListRes, err := adminClient.RunMealPlanGroceryListInitializerWorker(ctx, &mealplanninggrpc.RunMealPlanGroceryListInitializerWorkerRequest{})
		require.NoError(t, err)
		require.NotNil(t, runGroceryListRes)

		// Fetch the grocery list
		groceryListRes, err := accountAdminUserClient.GetMealPlanGroceryListItemsForMealPlan(ctx, &mealplanninggrpc.GetMealPlanGroceryListItemsForMealPlanRequest{
			MealPlanId: createdMealPlan.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, groceryListRes)
		require.NotEmpty(t, groceryListRes.Results, "expected grocery list to have items")

		// Verify that we have the selected option item (items with OptionIndex set)
		// Since we selected optionIndex=1, only Alternative B should be in the grocery list
		optionItemCount := 0
		for _, item := range groceryListRes.Results {
			if item.OptionIndex == nil {
				continue
			}

			optionItemCount++
			// Verify that the option items have the recipe context
			assert.NotNil(t, item.BelongsToMealPlanOption, "option items should have BelongsToMealPlanOption")
			assert.NotNil(t, item.RecipeId, "option items should have RecipeID")
			assert.NotNil(t, item.RecipeStepId, "option items should have RecipeStepID")
			assert.NotNil(t, item.IngredientIndex, "option items should have IngredientIndex")

			// Verify the selected optionIndex is correct
			// We selected optionIndex=1 for ingredientIndex=0 in the selection
			if item.IngredientIndex != nil && *item.IngredientIndex == 0 {
				assert.Equal(t, uint32(1), *item.OptionIndex, "expected the selected option (optionIndex=1) for ingredientIndex=0")
			}
		}
		// We expect 1 option item (only Alternative B since we selected optionIndex=1)
		// Alternative A (optionIndex=0) should NOT be included
		assert.Equal(t, 1, optionItemCount, "expected exactly 1 option item (only the selected alternative ingredient)")

		// Cleanup - use accountAdminUserClient for meal plan (owned by account)
		_, err = accountAdminUserClient.ArchiveMealPlan(ctx, &mealplanninggrpc.ArchiveMealPlanRequest{
			MealPlanId: createdMealPlan.ID,
		})
		assert.NoError(t, err)

		// Cleanup meals and recipes using adminClient (created by adminClient)
		_, _ = adminClient.ArchiveMeal(ctx, &mealplanninggrpc.ArchiveMealRequest{MealId: meal1.ID})
		_, _ = adminClient.ArchiveMeal(ctx, &mealplanninggrpc.ArchiveMealRequest{MealId: meal2.ID})

		_, _ = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: recipe1.ID})
		_, _ = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: recipe2.ID})

		// Cleanup valid ingredients and preparation
		for _, vi := range recipe1ValidIngredients {
			_, _ = adminClient.ArchiveValidIngredient(ctx, &mealplanninggrpc.ArchiveValidIngredientRequest{ValidIngredientId: vi.ID})
		}
		for _, vi := range recipe2ValidIngredients {
			_, _ = adminClient.ArchiveValidIngredient(ctx, &mealplanninggrpc.ArchiveValidIngredientRequest{ValidIngredientId: vi.ID})
		}
		_, _ = adminClient.ArchiveValidPreparation(ctx, &mealplanninggrpc.ArchiveValidPreparationRequest{ValidPreparationId: recipe1ValidPreparation.ID})
	})
}

// createRecipeWithAlternativeIngredients creates a recipe where step 1 has alternative ingredients
// at the same index (index=0) with different optionIndex values (0 and 1).
// Uses adminClient internally since recipe creation requires admin permissions.
func createRecipeWithAlternativeIngredients(t *testing.T, nameSuffix string) ([]*mealplanning.ValidIngredient, *mealplanning.ValidPreparation, *mealplanning.Recipe) {
	t.Helper()
	ctx := t.Context()

	// Create valid entities (all use adminClient internally)
	createdValidPreparation := createValidPreparationForTest(t)
	createdValidMeasurementUnit := createValidMeasurementUnitForTest(t)
	createdValidInstrument := createValidInstrumentForTest(t)
	createdValidIngredientState := createValidIngredientStateForTest(t)
	createdValidVessel := createValidVesselForTest(t)

	// Create bridge table entries
	createdValidPreparationInstrument := createValidPreparationInstrumentWithEntitiesForTest(t, createdValidPreparation, createdValidInstrument)
	createdValidPreparationVessel := createValidPreparationVesselWithEntitiesForTest(t, createdValidPreparation, createdValidVessel)

	// Create valid ingredients - we need at least 2 for alternatives
	createdValidIngredients := []*mealplanning.ValidIngredient{}
	ingredientPreparationIDs := make(map[string]string)
	ingredientMeasurementUnitIDs := make(map[string]string)

	for i := 0; i < 3; i++ { // Create 3 ingredients: 2 alternatives + 1 regular
		createdValidIngredient := createValidIngredientForTest(t)
		createdValidIngredients = append(createdValidIngredients, createdValidIngredient)

		// Create bridge table entries for this ingredient
		createdVIP := createValidIngredientPreparationWithEntitiesForTest(t, createdValidIngredient, createdValidPreparation)
		createdVIMU := createValidIngredientMeasurementUnitWithEntitiesForTest(t, createdValidIngredient, createdValidMeasurementUnit)

		ingredientPreparationIDs[createdValidIngredient.ID] = createdVIP.ID
		ingredientMeasurementUnitIDs[createdValidIngredient.ID] = createdVIMU.ID
	}

	// Build a fake recipe and customize it
	exampleRecipe := fakes.BuildFakeRecipe()
	exampleRecipe.Name = "Test Recipe With Alternatives " + nameSuffix
	exampleRecipe.Media = []*mealplanning.RecipeMedia{}
	exampleRecipe.EligibleForMeals = true

	// Create exactly 2 steps (validation requires at least 2 steps)
	step1 := fakes.BuildFakeRecipeStep()
	step2 := fakes.BuildFakeRecipeStep()
	exampleRecipe.Steps = []*mealplanning.RecipeStep{step1, step2}

	// Create alternative ingredients at the same index (index=0) with different optionIndex values
	// Ingredient A: index=0, optionIndex=0
	// Ingredient B: index=0, optionIndex=1
	// Ingredient C: index=1, optionIndex=0 (regular, non-alternative)
	step1.Ingredients = []*mealplanning.RecipeStepIngredient{
		{
			ID:                  fakes.BuildFakeID(),
			Name:                "Alternative A " + nameSuffix,
			Ingredient:          createdValidIngredients[0],
			MeasurementUnit:     *createdValidMeasurementUnit,
			Quantity:            types.Float32RangeWithOptionalMax{Min: 1.0, Max: pointer.To(float32(2.0))},
			Index:               0, // Same index as Alternative B
			OptionIndex:         0, // First option
			BelongsToRecipeStep: step1.ID,
			Optional:            false,
		},
		{
			ID:                  fakes.BuildFakeID(),
			Name:                "Alternative B " + nameSuffix,
			Ingredient:          createdValidIngredients[1],
			MeasurementUnit:     *createdValidMeasurementUnit,
			Quantity:            types.Float32RangeWithOptionalMax{Min: 1.5, Max: pointer.To(float32(2.5))},
			Index:               0, // Same index as Alternative A
			OptionIndex:         1, // Second option (alternative)
			BelongsToRecipeStep: step1.ID,
			Optional:            false,
		},
		{
			ID:                  fakes.BuildFakeID(),
			Name:                "Regular Ingredient " + nameSuffix,
			Ingredient:          createdValidIngredients[2],
			MeasurementUnit:     *createdValidMeasurementUnit,
			Quantity:            types.Float32RangeWithOptionalMax{Min: 0.5, Max: pointer.To(float32(1.0))},
			Index:               1, // Different index
			OptionIndex:         0, // Only one option at this index
			BelongsToRecipeStep: step1.ID,
			Optional:            false,
		},
	}

	// Step 2 has a simple ingredient (using the product from step 1)
	step2.Ingredients = []*mealplanning.RecipeStepIngredient{
		{
			ID:                  fakes.BuildFakeID(),
			Name:                "Secondary Ingredient " + nameSuffix,
			Ingredient:          createdValidIngredients[0],
			MeasurementUnit:     *createdValidMeasurementUnit,
			Quantity:            types.Float32RangeWithOptionalMax{Min: 0.5, Max: pointer.To(float32(1.0))},
			Index:               0,
			OptionIndex:         0,
			BelongsToRecipeStep: step2.ID,
			Optional:            false,
		},
	}

	// Set up instruments and vessels for both steps
	for _, step := range exampleRecipe.Steps {
		if len(step.Instruments) == 0 {
			step.Instruments = []*mealplanning.RecipeStepInstrument{fakes.BuildFakeRecipeStepInstrument()}
		}
		for j := range step.Instruments {
			step.Instruments[j].Instrument = createdValidInstrument
			step.Instruments[j].Index = uint16(j)
			step.Instruments[j].OptionIndex = 0
		}

		if len(step.Vessels) == 0 {
			step.Vessels = []*mealplanning.RecipeStepVessel{fakes.BuildFakeRecipeStepVessel()}
		}
		for j := range step.Vessels {
			step.Vessels[j].Vessel = createdValidVessel
			step.Vessels[j].Index = uint16(j)
			step.Vessels[j].OptionIndex = 0
		}

		// Set up products
		if len(step.Products) == 0 {
			step.Products = []*mealplanning.RecipeStepProduct{fakes.BuildFakeRecipeStepProduct()}
		}
		for j := range step.Products {
			step.Products[j].MeasurementUnit = createdValidMeasurementUnit
			step.Products[j].Index = uint16(j)
		}

		// Set up completion conditions
		if len(step.CompletionConditions) > 0 && len(step.Ingredients) > 0 {
			for j := range step.CompletionConditions {
				step.CompletionConditions[j].IngredientState = *createdValidIngredientState
				for k := range step.CompletionConditions[j].Ingredients {
					step.CompletionConditions[j].Ingredients[k].RecipeStepIngredient = step.Ingredients[0].ID
				}
			}
		}
	}

	// Convert to creation input
	exampleRecipeInput := mpconverters.ConvertRecipeToRecipeCreationRequestInput(exampleRecipe)
	exampleRecipeInput.AlsoCreateMeal = false

	// Set up bridge table IDs for all steps
	for i := range exampleRecipeInput.Steps {
		exampleRecipeInput.Steps[i].PreparationID = createdValidPreparation.ID

		for j := range exampleRecipeInput.Steps[i].Ingredients {
			// For step 0 (with alternatives), map ingredients based on j
			// For other steps, use the first ingredient
			var ingredientID string
			if i == 0 && j < len(createdValidIngredients) {
				ingredientID = createdValidIngredients[j].ID
			} else {
				ingredientID = createdValidIngredients[0].ID
			}

			if vipID, ok := ingredientPreparationIDs[ingredientID]; ok {
				exampleRecipeInput.Steps[i].Ingredients[j].ValidIngredientPreparationID = &vipID
			}
			if vimuID, ok := ingredientMeasurementUnitIDs[ingredientID]; ok {
				exampleRecipeInput.Steps[i].Ingredients[j].ValidIngredientMeasurementUnitID = &vimuID
			}

			// Set the index explicitly for step 0
			if i == 0 {
				idx := uint16(0)
				if j == 2 { // Third ingredient has different index
					idx = 1
				}
				exampleRecipeInput.Steps[i].Ingredients[j].Index = &idx
			} else {
				idx := uint16(j)
				exampleRecipeInput.Steps[i].Ingredients[j].Index = &idx
			}
		}

		for j := range exampleRecipeInput.Steps[i].Instruments {
			exampleRecipeInput.Steps[i].Instruments[j].ValidPreparationInstrumentID = &createdValidPreparationInstrument.ID
		}

		for j := range exampleRecipeInput.Steps[i].Vessels {
			exampleRecipeInput.Steps[i].Vessels[j].ValidPreparationVesselID = &createdValidPreparationVessel.ID
		}
	}

	// Clear prep tasks since we're not testing them
	exampleRecipeInput.PrepTasks = nil

	// Create the recipe using adminClient (requires admin permissions)
	createdRecipeRes, err := adminClient.CreateRecipe(ctx, &mealplanninggrpc.CreateRecipeRequest{
		Input: converters.ConvertRecipeCreationRequestInputToGRPCRecipeCreationRequestInput(exampleRecipeInput),
	})
	require.NoError(t, err)
	require.NotEmpty(t, createdRecipeRes.Created.Id)

	// Fetch the created recipe
	fetchedRecipeRes, err := adminClient.GetRecipe(ctx, &mealplanninggrpc.GetRecipeRequest{
		RecipeId: createdRecipeRes.Created.Id,
	})
	require.NoError(t, err)

	createdRecipe := converters.ConvertGRPCRecipeToRecipe(fetchedRecipeRes.Result)
	require.NotEmpty(t, createdRecipe.Steps)
	require.GreaterOrEqual(t, len(createdRecipe.Steps[0].Ingredients), 2, "expected at least 2 ingredients")

	return createdValidIngredients, createdValidPreparation, createdRecipe
}

// createMealFromRecipe creates a meal from a single recipe.
// Uses adminClient internally since meal creation requires admin permissions.
func createMealFromRecipe(t *testing.T, recipe *mealplanning.Recipe, nameSuffix string) *mealplanning.Meal {
	t.Helper()
	ctx := t.Context()

	exampleMeal := fakes.BuildFakeMeal()
	exampleMeal.Name = "Test Meal " + nameSuffix
	exampleMeal.EligibleForMealPlans = true

	exampleMealInput := mpconverters.ConvertMealToMealCreationRequestInput(exampleMeal)
	exampleMealInput.Components = []*mealplanning.MealComponentCreationRequestInput{
		{
			RecipeID:      recipe.ID,
			ComponentType: mealplanning.MealComponentTypesMain,
			RecipeScale:   1.0,
		},
	}

	// Use adminClient for meal creation
	createdMealRes, err := adminClient.CreateMeal(ctx, &mealplanninggrpc.CreateMealRequest{
		Input: converters.ConvertMealCreationRequestInputToGRPCMealCreationRequestInput(exampleMealInput),
	})
	require.NoError(t, err)

	fetchedMealRes, err := adminClient.GetMeal(ctx, &mealplanninggrpc.GetMealRequest{
		MealId: createdMealRes.Created.Id,
	})
	require.NoError(t, err)

	createdMeal := converters.ConvertGRPCMealToMeal(fetchedMealRes.Result)
	return createdMeal
}
