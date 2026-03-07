package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	mealplangrocerylistinitializerbuild "github.com/dinnerdonebetter/backend/internal/build/jobs/meal_plan_grocery_list_initializer"
	mealplantaskcreatorbuild "github.com/dinnerdonebetter/backend/internal/build/jobs/meal_plan_task_creator"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identityconverters "github.com/dinnerdonebetter/backend/internal/domain/identity/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/bootstrap"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/domain/settings"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/platform/search/text/config"
	identitygenerated "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity/generated"
	mealplanningrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning"
)

const (
	apiConfigurationFilepath = "deploy/environments/testing/config_files/integration-tests-config.json"
	// createMealPlansAndVotes controls whether meal plans and votes are created during localdev startup.
	// Set to true via environment variable CREATE_MEAL_PLANS_AND_VOTES=true to enable.
)

var (
	createMealPlansAndVotes = true // os.Getenv("CREATE_MEAL_PLANS_AND_VOTES") == "true"
)

func cloneTime(t time.Time) time.Time {
	t, parseErr := time.Parse(time.RFC3339, t.Format(time.RFC3339))
	if parseErr != nil {
		panic(parseErr)
	}

	return t
}

// resolveEmptyRecipeIDs resolves empty RecipeStepProductRecipeID values in a recipe input
// by looking up recipes in the createdRecipes map. Prefers RecipeStepProductRecipeSlug when
// set (ensures correct recipe); otherwise falls back to matching by step index (ambiguous).
func resolveEmptyRecipeIDs(recipe *mealplanning.RecipeCreationRequestInput, createdRecipes map[string]*mealplanning.Recipe) {
	for _, step := range recipe.Steps {
		for _, ingredient := range step.Ingredients {
			if ingredient.RecipeStepProductRecipeID != nil && *ingredient.RecipeStepProductRecipeID == "" {
				// Prefer slug lookup when available - ensures we resolve to the correct recipe
				if ingredient.RecipeStepProductRecipeSlug != nil && *ingredient.RecipeStepProductRecipeSlug != "" {
					if refRecipe, ok := createdRecipes[*ingredient.RecipeStepProductRecipeSlug]; ok && refRecipe != nil {
						ingredient.RecipeStepProductRecipeID = &refRecipe.ID
						continue
					}
				}
				// Fallback: match by step index (ambiguous - many recipes may have a step at that index)
				stepIndex := ingredient.ProductOfRecipeStepIndex
				if stepIndex != nil {
					for _, refRecipe := range createdRecipes {
						if refRecipe != nil && int(*stepIndex) < len(refRecipe.Steps) {
							referencedStep := refRecipe.Steps[*stepIndex]
							if ingredient.ProductOfRecipeStepProductIndex != nil {
								productIndex := int(*ingredient.ProductOfRecipeStepProductIndex)
								if productIndex < len(referencedStep.Products) {
									ingredient.RecipeStepProductRecipeID = &refRecipe.ID
									break
								}
							}
						}
					}
				}
			}
		}
	}
}

func main() {
	ctx := context.Background()

	// create premade admin user
	premadeAdminUser := &identity.User{
		ID:              strings.Repeat("a", 20),
		TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		EmailAddress:    "integration_tests@example.email",
		Username:        "admin_user",
		HashedPassword:  "admin_pass",
	}

	apiConfig, err := config.LoadConfigFromPath[config.APIServiceConfig](ctx, apiConfigurationFilepath)
	if err != nil {
		log.Fatal(err)
	}

	var adminUserID string
	var adminAccountID string
	var currentMealPlanID string
	var memberUserIDs []string

	server, err := localdev.AllInOne(
		ctx,
		apiConfig,
		// Create admin user and get account
		localdev.WithIdentityRepository(func(ctx context.Context, repo identity.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider, dbClient database.Client) error {
			user, userErr := localdev.CreatePremadeAdminUser(ctx, logger, tracerProvider, repo, dbClient, premadeAdminUser)
			if userErr != nil {
				return userErr
			}
			adminUserID = user.ID

			// Get or create account for admin user
			accounts, accountsErr := repo.GetAccounts(ctx, adminUserID, nil)
			if accountsErr != nil {
				return fmt.Errorf("failed to get accounts for admin user: %w", accountsErr)
			}

			if len(accounts.Data) > 0 {
				// Use first account
				adminAccountID = accounts.Data[0].ID
			} else {
				// Create a new account for the admin user
				accountInput := &identity.AccountCreationRequestInput{
					Name:          "Admin Household",
					BelongsToUser: adminUserID,
				}
				account, accountErr := repo.CreateAccount(ctx, identityconverters.ConvertAccountCreationInputToAccountDatabaseCreationInput(accountInput))
				if accountErr != nil {
					return fmt.Errorf("failed to create account for admin user: %w", accountErr)
				}
				adminAccountID = account.ID
			}

			if adminAccountID == "" {
				return fmt.Errorf("admin account ID not set")
			}

			hasher := authentication.ProvideArgon2Authenticator(logger, tracerProvider)
			generatedQuerier := identitygenerated.New()

			// Create two member users
			memberUsers := []*struct {
				username  string
				email     string
				password  string
				firstName string
				lastName  string
				userID    string
			}{
				{
					username:  "member_user_1",
					email:     "member1@example.email",
					password:  "member_pass_1",
					firstName: "Member",
					lastName:  "One",
					userID:    strings.Repeat("c", 20),
				},
				{
					username:  "member_user_2",
					email:     "member2@example.email",
					password:  "member_pass_2",
					firstName: "Member",
					lastName:  "Two",
					userID:    strings.Repeat("d", 20),
				},
			}

			for _, memberUser := range memberUsers {
				// Check if user already exists
				existingUser, userExistsErr := repo.GetUserByUsername(ctx, memberUser.username)
				if userExistsErr == nil && existingUser != nil {
					logger.Info(fmt.Sprintf("User %s already exists, skipping creation", memberUser.username))
					// Still add to account if not already a member
					isMember, memberErr := repo.UserIsMemberOfAccount(ctx, existingUser.ID, adminAccountID)
					if memberErr == nil && !isMember {
						membershipID := identifiers.New()
						if err = generatedQuerier.AddUserToAccount(ctx, dbClient.WriteDB(), &identitygenerated.AddUserToAccountParams{
							ID:               membershipID,
							BelongsToUser:    existingUser.ID,
							BelongsToAccount: adminAccountID,
							AccountRole:      authorization.AccountMemberRole.String(),
						}); err != nil {
							return fmt.Errorf("failed to add existing user %s to account: %w", memberUser.username, err)
						}
						logger.Info(fmt.Sprintf("Added existing user %s to admin account", memberUser.username))
					}
					memberUserIDs = append(memberUserIDs, existingUser.ID)
					continue
				}

				// Hash password
				hashedPassword, hashErr := hasher.HashPassword(ctx, memberUser.password)
				if hashErr != nil {
					return fmt.Errorf("failed to hash password for user %s: %w", memberUser.username, hashErr)
				}

				// Create user
				userInput := &identity.User{
					ID:              memberUser.userID,
					Username:        memberUser.username,
					EmailAddress:    memberUser.email,
					HashedPassword:  hashedPassword,
					TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
					FirstName:       memberUser.firstName,
					LastName:        memberUser.lastName,
				}

				user, userErr = repo.CreateUser(ctx, identityconverters.ConvertUserToUserDatabaseCreationInput(userInput))
				if userErr != nil {
					return fmt.Errorf("failed to create user %s: %w", memberUser.username, userErr)
				}

				// Mark two-factor secret as verified
				if err = repo.MarkUserTwoFactorSecretAsVerified(ctx, user.ID); err != nil {
					return fmt.Errorf("failed to mark user %s as verified: %w", memberUser.username, err)
				}

				// Add user to admin account as a member
				membershipID := identifiers.New()
				if err = generatedQuerier.AddUserToAccount(ctx, dbClient.WriteDB(), &identitygenerated.AddUserToAccountParams{
					ID:               membershipID,
					BelongsToUser:    user.ID,
					BelongsToAccount: adminAccountID,
					AccountRole:      authorization.AccountMemberRole.String(),
				}); err != nil {
					return fmt.Errorf("failed to add user %s to account: %w", memberUser.username, err)
				}

				if err = generatedQuerier.MarkAccountUserMembershipAsUserDefault(ctx, dbClient.WriteDB(), &identitygenerated.MarkAccountUserMembershipAsUserDefaultParams{
					BelongsToUser:    user.ID,
					BelongsToAccount: adminAccountID,
				}); err != nil {
					return fmt.Errorf("failed to mark user %s account as default: %w", memberUser.username, err)
				}

				memberUserIDs = append(memberUserIDs, user.ID)
				logger.Info(fmt.Sprintf("Created user %s and added to admin account", memberUser.username))
			}

			// Also add admin user to member list
			memberUserIDs = append(memberUserIDs, adminUserID)

			return nil
		}),
		// Create OAuth2 client
		localdev.WithOAuth2Repository(func(ctx context.Context, repo oauth.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
			_, err = repo.CreateOAuth2Client(ctx, &oauth.OAuth2ClientDatabaseCreationInput{
				ID:           strings.Repeat("b", 20),
				Name:         "localdev_admin_client",
				Description:  "localdev admin client",
				ClientID:     strings.Repeat("A", oauth.ClientIDSize),
				ClientSecret: strings.Repeat("A", oauth.ClientSecretSize),
			})
			return err
		}),
		// Create valid enumerations and bridge types, then create all bootstrap recipes
		localdev.WithMealPlanningRepository(func(ctx context.Context, repo mealplanning.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
			logger.Info("Creating enumerations...")
			enums, enumsErr := bootstrap.CreateEnumerations(ctx, repo, logger)
			if enumsErr != nil {
				return fmt.Errorf("failed to create enumerations: %w", enumsErr)
			}
			logger.Info("Enumerations created successfully!")

			// Create RecipeManager to create the first recipe
			logger.Info("Creating RecipeManager...")
			queueCfg := &msgconfig.QueuesConfig{
				DataChangesTopicName: "data_changes",
			}
			publisherProvider := messagequeue.NewNoopPublisherProvider()
			recipeAnalyzer := recipeanalysis.NewRecipeAnalyzer(logger, tracerProvider)
			searchConfig := &textsearchcfg.Config{}
			metricsProvider := metrics.NewNoopMetricsProvider()

			recipeManager, recipeManagerErr := managers.NewRecipeManager(
				ctx,
				logger,
				tracerProvider,
				repo,
				queueCfg,
				publisherProvider,
				recipeAnalyzer,
				searchConfig,
				metricsProvider,
			)
			if recipeManagerErr != nil {
				return fmt.Errorf("failed to create recipe manager: %w", recipeManagerErr)
			}
			logger.Info("RecipeManager created successfully!")

			logger.Info("Creating remaining bootstrap recipes...")

			// Phase 1: Create recipes without prerequisites
			allRecipes := bootstrap.AllRecipes(enums)
			logger.Info(fmt.Sprintf("Found %d recipes without prerequisites to create", len(allRecipes)))

			createdRecipes := make(map[string]*mealplanning.Recipe)
			// Create recipes without prerequisites
			for i, recipe := range allRecipes {
				logger.Info(fmt.Sprintf("Creating recipe %d: %s (%d steps)", i+1, recipe.Name, len(recipe.Steps)))
				r, createErr := recipeManager.CreateRecipe(ctx, adminUserID, recipe)
				if createErr != nil {
					return fmt.Errorf("failed to create recipe #%d %s: %w", i, recipe.Name, createErr)
				}

				createdRecipes[r.Slug] = r
			}
			logger.Info("All recipes without prerequisites created successfully!")

			recipes := slices.Collect(maps.Values(createdRecipes))

			// Phase 2: Create recipes with prerequisites
			recipesWithPrerequisites := bootstrap.AllRecipesWithPrerequisites(enums, createdRecipes)
			logger.Info(fmt.Sprintf("Found %d recipes with prerequisites to create", len(recipesWithPrerequisites)))

			for i, recipe := range recipesWithPrerequisites {
				// Resolve empty recipe IDs in cross-recipe references before creating
				// This is needed because getRecipeIDBySlug may return empty strings when
				// called during recipe input construction (before prerequisite recipes exist)
				resolveEmptyRecipeIDs(recipe, createdRecipes)

				logger.Info(fmt.Sprintf("Creating recipe with prerequisites %d: %s (%d steps)", i+1, recipe.Name, len(recipe.Steps)))
				r, createErr := recipeManager.CreateRecipe(ctx, adminUserID, recipe)
				if createErr != nil {
					return fmt.Errorf("failed to create recipe with prerequisites #%d %s: %w", i, recipe.Name, createErr)
				}

				recipes = append(recipes, r)
				// Update lookup map so subsequent recipes in phase 2 can reference this one
				createdRecipes[r.Slug] = r
			}
			logger.Info("All bootstrap recipes created successfully!")

			// Approve all bootstrap recipes
			for _, r := range recipes {
				if err := recipeManager.UpdateRecipeStatus(ctx, r.ID, mealplanning.RecipeStatusApproved); err != nil {
					return fmt.Errorf("failed to approve recipe %s: %w", r.Name, err)
				}
			}
			logger.Info("All bootstrap recipes approved!")

			// Always create meals
			logger.Info("Creating bootstrap meals...")
			meals := bootstrap.AllMeals(adminUserID, recipes)
			logger.Info(fmt.Sprintf("Found %d meals to create", len(meals)))

			for i, meal := range meals {
				logger.Info(fmt.Sprintf("Creating meal %d: %s (%d components)", i+1, meal.Name, len(meal.Components)))
				_, err = repo.CreateMeal(ctx, meal)
				if err != nil {
					return fmt.Errorf("failed to create meal %s: %w", meal.Name, err)
				}
			}
			logger.Info("All bootstrap meals created successfully!")

			return nil
		}),
		// Create meal plan with 3 chicken dishes
		localdev.WithMealPlanningRepository(func(ctx context.Context, repo mealplanning.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
			// Check if meal plan creation is enabled (via constant or environment variable)
			shouldCreate := createMealPlansAndVotes
			if !shouldCreate {
				logger.Info("Skipping meal plan creation (CREATE_MEAL_PLANS_AND_VOTES=false)")
				return nil
			}
			if adminUserID == "" || adminAccountID == "" {
				return fmt.Errorf("admin user ID or account ID not set")
			}

			logger.Info("Creating meal plan with chicken dishes...")

			// Get all meals created by admin user
			mealsResult, mealsErr := repo.GetMealsCreatedByUser(ctx, adminUserID, nil)
			if mealsErr != nil {
				return fmt.Errorf("failed to get meals: %w", mealsErr)
			}

			// Find the 3 chicken dishes
			chickenMealNames := []string{
				"Sous Vide Chicken Breast with Rice",
				"Roast Chicken with Caesar Broccoli",
				"Soy Sauce Braised Chicken Thighs with Rice",
			}

			otherMealNames := []string{
				"Pan-Seared Steak with Mashed Potatoes",
				"Classic Burgers with Mixed Green Salad",
				"Grilled Pork Tenderloin with Brussels Sprouts",
			}

			var (
				chickenMeals []*mealplanning.Meal
				otherMeals   []*mealplanning.Meal
			)
			for _, meal := range mealsResult.Data {
				if slices.Contains(chickenMealNames, meal.Name) {
					chickenMeals = append(chickenMeals, meal)
				}

				for _, name := range otherMealNames {
					if meal.Name == name {
						otherMeals = append(otherMeals, meal)
					}
				}
			}

			now := time.Now()

			// Voting deadline is Friday before the event (midnight)
			votingDeadline := cloneTime(now).Add(24 * time.Hour * 3)

			eventStart := cloneTime(votingDeadline).Add(24 * time.Hour)
			eventEnd := cloneTime(eventStart).Add(2 * time.Hour) // 2 hour duration

			// Create options for all three chicken meals
			var chickenOptions []*mealplanning.MealPlanOptionDatabaseCreationInput
			for _, chickenMeal := range chickenMeals {
				chickenOptions = append(chickenOptions, &mealplanning.MealPlanOptionDatabaseCreationInput{
					ID:        identifiers.New(),
					MealID:    chickenMeal.ID,
					MealScale: 1.0,
				})
			}

			var otherOptions []*mealplanning.MealPlanOptionDatabaseCreationInput
			for _, otherMeal := range otherMeals {
				otherOptions = append(otherOptions, &mealplanning.MealPlanOptionDatabaseCreationInput{
					ID:        identifiers.New(),
					MealID:    otherMeal.ID,
					MealScale: 1.0,
				})
			}

			// Create a single event with all three chickenOptions
			events := []*mealplanning.MealPlanEventDatabaseCreationInput{
				{
					ID:       identifiers.New(),
					StartsAt: eventStart,
					EndsAt:   eventEnd,
					MealName: mealplanning.DinnerMealName,
					Options:  chickenOptions,
				},
				{
					ID:       identifiers.New(),
					StartsAt: cloneTime(eventStart).Add(24 * time.Hour),
					EndsAt:   cloneTime(eventEnd).Add(24 * time.Hour),
					MealName: mealplanning.SupperMealName,
					Options:  otherOptions,
				},
			}

			// Create meal plan
			mealPlanInput := &mealplanning.MealPlanDatabaseCreationInput{
				ID:               identifiers.New(),
				Notes:            "Example 	Meal Plan",
				VotingDeadline:   votingDeadline,
				ElectionMethod:   mealplanning.MealPlanElectionMethodSchulze,
				BelongsToAccount: adminAccountID,
				CreatedByUser:    adminUserID,
				Events:           events,
			}

			createdMealPlan, mealPlanErr := repo.CreateMealPlan(ctx, mealPlanInput)
			if mealPlanErr != nil {
				return fmt.Errorf("failed to create meal plan: %w", mealPlanErr)
			}

			currentMealPlanID = createdMealPlan.ID
			logger.Info(fmt.Sprintf("Created meal plan %s with %d events", createdMealPlan.ID, len(events)))
			return nil
		}),
		// Create finalized meal plan with votes and extend current meal plan deadline
		localdev.WithMealPlanningRepository(func(ctx context.Context, repo mealplanning.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
			// Check if meal plan creation is enabled (via constant or environment variable)
			shouldCreate := createMealPlansAndVotes
			if !shouldCreate {
				logger.Info("Skipping finalized meal plan and vote creation (CREATE_MEAL_PLANS_AND_VOTES=false)")
				return nil
			}
			if adminUserID == "" || adminAccountID == "" {
				return fmt.Errorf("admin user ID or account ID not set")
			}

			if currentMealPlanID == "" {
				return fmt.Errorf("current meal plan ID not set")
			}

			if len(memberUserIDs) == 0 {
				return fmt.Errorf("member user IDs not set")
			}

			// Get the current meal plan to use its timing
			currentMealPlan, mealPlanErr := repo.GetMealPlan(ctx, currentMealPlanID, adminAccountID)
			if mealPlanErr != nil {
				return fmt.Errorf("failed to get current meal plan: %w", mealPlanErr)
			}

			logger.Info("Creating finalized meal plan with votes...")

			// Get all meals created by admin user
			mealsResult, mealsErr := repo.GetMealsCreatedByUser(ctx, adminUserID, nil)
			if mealsErr != nil {
				return fmt.Errorf("failed to get meals: %w", mealsErr)
			}

			// Find the 3 chicken dishes
			chickenMealNames := []string{
				"Sous Vide Chicken Breast with Rice",
				"Roast Chicken with Caesar Broccoli",
				"Soy Sauce Braised Chicken Thighs with Rice",
			}

			otherMealNames := []string{
				"Pan-Seared Steak with Mashed Potatoes",
				"Classic Burgers with Mixed Green Salad",
				"Grilled Pork Tenderloin with Brussels Sprouts",
			}

			var (
				chickenMeals []*mealplanning.Meal
				otherMeals   []*mealplanning.Meal
			)
			for _, meal := range mealsResult.Data {
				if slices.Contains(chickenMealNames, meal.Name) {
					chickenMeals = append(chickenMeals, meal)
				}

				for _, name := range otherMealNames {
					if meal.Name == name {
						otherMeals = append(otherMeals, meal)
					}
				}
			}

			// Use the same timing as the current meal plan
			finalizedVotingDeadline := currentMealPlan.VotingDeadline
			finalizedEventStart := cloneTime(finalizedVotingDeadline).Add(24 * time.Hour)
			finalizedEventEnd := cloneTime(finalizedEventStart).Add(2 * time.Hour)

			// Create options for all three chicken meals
			var chickenOptions []*mealplanning.MealPlanOptionDatabaseCreationInput
			for _, chickenMeal := range chickenMeals {
				chickenOptions = append(chickenOptions, &mealplanning.MealPlanOptionDatabaseCreationInput{
					ID:        identifiers.New(),
					MealID:    chickenMeal.ID,
					MealScale: 1.0,
				})
			}

			var otherOptions []*mealplanning.MealPlanOptionDatabaseCreationInput
			for _, otherMeal := range otherMeals {
				otherOptions = append(otherOptions, &mealplanning.MealPlanOptionDatabaseCreationInput{
					ID:        identifiers.New(),
					MealID:    otherMeal.ID,
					MealScale: 1.0,
				})
			}

			// Create events for finalized meal plan
			finalizedEvents := []*mealplanning.MealPlanEventDatabaseCreationInput{
				{
					ID:       identifiers.New(),
					StartsAt: finalizedEventStart,
					EndsAt:   finalizedEventEnd,
					MealName: mealplanning.DinnerMealName,
					Options:  chickenOptions,
				},
				{
					ID:       identifiers.New(),
					StartsAt: cloneTime(finalizedEventStart).Add(24 * time.Hour),
					EndsAt:   cloneTime(finalizedEventEnd).Add(24 * time.Hour),
					MealName: mealplanning.SupperMealName,
					Options:  otherOptions,
				},
			}

			// Create finalized meal plan
			finalizedMealPlanInput := &mealplanning.MealPlanDatabaseCreationInput{
				ID:               identifiers.New(),
				Notes:            "Finalized Example Meal Plan",
				VotingDeadline:   finalizedVotingDeadline,
				ElectionMethod:   mealplanning.MealPlanElectionMethodSchulze,
				BelongsToAccount: adminAccountID,
				CreatedByUser:    adminUserID,
				Events:           finalizedEvents,
			}

			finalizedMealPlan, finalizedErr := repo.CreateMealPlan(ctx, finalizedMealPlanInput)
			if finalizedErr != nil {
				return fmt.Errorf("failed to create finalized meal plan: %w", finalizedErr)
			}

			logger.Info(fmt.Sprintf("Created finalized meal plan %s with %d events", finalizedMealPlan.ID, len(finalizedEvents)))

			// Create votes from all members for all options in all events
			// We need to reload the meal plan to get the created options
			finalizedMealPlanWithEvents, finalizeErr := repo.GetMealPlan(ctx, finalizedMealPlan.ID, adminAccountID)
			if finalizeErr != nil {
				return fmt.Errorf("failed to get finalized meal plan with events: %w", finalizeErr)
			}

			for _, event := range finalizedMealPlanWithEvents.Events {
				for _, memberUserID := range memberUserIDs {
					// Create votes for this user for all options in this event
					var votes []*mealplanning.MealPlanOptionVoteDatabaseCreationInput
					for rank, option := range event.Options {
						votes = append(votes, &mealplanning.MealPlanOptionVoteDatabaseCreationInput{
							ID:                      identifiers.New(),
							ByUser:                  memberUserID,
							BelongsToMealPlanOption: option.ID,
							Rank:                    uint8(rank),
							Abstain:                 false,
							Notes:                   "",
						})
					}

					// Create votes for this user
					voteInput := &mealplanning.MealPlanOptionVotesDatabaseCreationInput{
						Votes: votes,
					}
					_, voteErr := repo.CreateMealPlanOptionVote(ctx, voteInput)
					if voteErr != nil {
						return fmt.Errorf("failed to create votes for user %s: %w", memberUserID, voteErr)
					}
				}
			}

			logger.Info("Created votes from all members for finalized meal plan")

			// Finalize the meal plan (idempotent: already-finalized is OK on re-run)
			finalized, finalizeErr := repo.AttemptToFinalizeMealPlan(ctx, finalizedMealPlan.ID, adminAccountID)
			switch {
			case finalizeErr != nil:
				if errors.Is(finalizeErr, mealplanningrepo.ErrAlreadyFinalized) {
					logger.Info("Meal plan already finalized (idempotent re-run), continuing")
				} else {
					return fmt.Errorf("failed to finalize meal plan: %w", finalizeErr)
				}
			case finalized:
				logger.Info("Finalized meal plan successfully")
			default:
				return fmt.Errorf("meal plan was not finalized")
			}

			// Extend the current meal plan's voting deadline by one week
			updatedMealPlan := *currentMealPlan
			updatedMealPlan.VotingDeadline = cloneTime(currentMealPlan.VotingDeadline).Add(7 * 24 * time.Hour)
			updateErr := repo.UpdateMealPlan(ctx, &updatedMealPlan)
			if updateErr != nil {
				return fmt.Errorf("failed to update current meal plan voting deadline: %w", updateErr)
			}

			logger.Info(fmt.Sprintf("Extended current meal plan %s voting deadline by one week", currentMealPlanID))
			return nil
		}),
		// Run grocery list initializer and task creator workers for finalized meal plans
		localdev.WithMealPlanningRepository(func(ctx context.Context, repo mealplanning.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
			// Check if meal plan creation is enabled (via constant or environment variable)
			shouldCreate := createMealPlansAndVotes
			if !shouldCreate {
				logger.Info("Skipping grocery list and task creator workers (CREATE_MEAL_PLANS_AND_VOTES=false)")
				return nil
			}
			logger.Info("Running grocery list initializer and task creator workers...")

			// Build grocery list initializer worker
			groceryListConfig := &config.MealPlanGroceryListInitializerConfig{
				Database:      apiConfig.Database,
				Observability: apiConfig.Observability,
				Events:        apiConfig.Events,
				Queues:        apiConfig.Queues,
				Analytics:     apiConfig.Analytics,
			}
			groceryListWorker, workerErr := mealplangrocerylistinitializerbuild.Build(ctx, groceryListConfig)
			if workerErr != nil {
				return fmt.Errorf("failed to build grocery list initializer worker: %w", workerErr)
			}

			// Run grocery list initializer
			if err = groceryListWorker.Work(ctx); err != nil {
				return fmt.Errorf("failed to run grocery list initializer worker: %w", err)
			}
			logger.Info("Grocery list initializer worker completed successfully")

			// Build task creator worker
			taskCreatorConfig := &config.MealPlanTaskCreatorConfig{
				Database:      apiConfig.Database,
				Observability: apiConfig.Observability,
				Events:        apiConfig.Events,
				Queues:        apiConfig.Queues,
				Analytics:     apiConfig.Analytics,
			}
			taskCreatorWorker, workerErr := mealplantaskcreatorbuild.Build(ctx, taskCreatorConfig)
			if workerErr != nil {
				return fmt.Errorf("failed to build task creator worker: %w", workerErr)
			}

			// Run task creator
			if err = taskCreatorWorker.Work(ctx); err != nil {
				return fmt.Errorf("failed to run task creator worker: %w", err)
			}
			logger.Info("Task creator worker completed successfully")

			return nil
		}),
		// Create example service settings
		localdev.WithSettingsRepository(func(ctx context.Context, repo settings.Repository, logger logging.Logger, tracerProvider tracing.TracerProvider) error {
			return createExampleServiceSettings(ctx, repo, logger)
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting server")

	if os.Getenv("DRY_RUN") == "true" {
		log.Println("dry run is enabled, skipping server run")
		return
	}
	server.Run()
}

func createExampleServiceSettings(ctx context.Context, repo settings.Repository, logger logging.Logger) error {
	defaultTheme := "light"
	_, err := repo.CreateServiceSetting(ctx, &settings.ServiceSettingDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         "user_theme_preference",
		Type:         "user",
		Description:  "User's preferred theme for the application interface",
		Enumeration:  []string{"light", "dark", "auto"},
		DefaultValue: &defaultTheme,
		AdminsOnly:   true,
	})
	if err != nil {
		return fmt.Errorf("failed to create theme preference setting: %w", err)
	}
	logger.Debug("Created ServiceSetting: user_theme_preference (enumerated with default)")

	defaultNotificationFreq := "daily"
	_, err = repo.CreateServiceSetting(ctx, &settings.ServiceSettingDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         "membership_notification_frequency",
		Type:         "membership",
		Description:  "How often to send notifications to membership members",
		Enumeration:  []string{"immediate", "daily", "weekly", "never"},
		DefaultValue: &defaultNotificationFreq,
		AdminsOnly:   true,
	})
	if err != nil {
		return fmt.Errorf("failed to create notification frequency setting: %w", err)
	}
	logger.Debug("Created ServiceSetting: membership_notification_frequency (enumerated with default)")

	defaultLanguage := "en"
	_, err = repo.CreateServiceSetting(ctx, &settings.ServiceSettingDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         "user_language",
		Type:         "user",
		Description:  "User's preferred language for the application",
		Enumeration:  []string{"en", "es", "fr", "de", "it"},
		DefaultValue: &defaultLanguage,
		AdminsOnly:   false,
	})
	if err != nil {
		return fmt.Errorf("failed to create language setting: %w", err)
	}
	logger.Debug("Created ServiceSetting: user_language (enumerated, non-admin)")

	return nil
}
