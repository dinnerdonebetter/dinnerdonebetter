package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identityconverters "github.com/dinnerdonebetter/backend/internal/domain/identity/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/bootstrap"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/domain/settings"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	identitygenerated "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity/generated"
)

const (
	apiConfigurationFilepath = "deploy/environments/testing/config_files/integration-tests-config.json"
)

func cloneTime(t time.Time) time.Time {
	t, parseErr := time.Parse(time.RFC3339, t.Format(time.RFC3339))
	if parseErr != nil {
		panic(parseErr)
	}

	return t
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
						if err = generatedQuerier.AddUserToAccount(ctx, dbClient.DB(), &identitygenerated.AddUserToAccountParams{
							ID:               membershipID,
							BelongsToUser:    existingUser.ID,
							BelongsToAccount: adminAccountID,
							AccountRole:      authorization.AccountMemberRole.String(),
						}); err != nil {
							return fmt.Errorf("failed to add existing user %s to account: %w", memberUser.username, err)
						}
						logger.Info(fmt.Sprintf("Added existing user %s to admin account", memberUser.username))
					}
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
				if err = generatedQuerier.AddUserToAccount(ctx, dbClient.DB(), &identitygenerated.AddUserToAccountParams{
					ID:               membershipID,
					BelongsToUser:    user.ID,
					BelongsToAccount: adminAccountID,
					AccountRole:      authorization.AccountMemberRole.String(),
				}); err != nil {
					return fmt.Errorf("failed to add user %s to account: %w", memberUser.username, err)
				}

				if err = generatedQuerier.MarkAccountUserMembershipAsUserDefault(ctx, dbClient.DB(), &identitygenerated.MarkAccountUserMembershipAsUserDefaultParams{
					BelongsToUser:    user.ID,
					BelongsToAccount: adminAccountID,
				}); err != nil {
					return fmt.Errorf("failed to mark user %s account as default: %w", memberUser.username, err)
				}

				logger.Info(fmt.Sprintf("Created user %s and added to admin account", memberUser.username))
			}

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

			logger.Info("Creating bootstrap recipes...")
			recipes := bootstrap.AllRecipes(adminUserID, enums)
			logger.Info(fmt.Sprintf("Found %d recipes to create", len(recipes)))

			for i, recipe := range recipes {
				logger.Info(fmt.Sprintf("Creating recipe %d: %s (%d steps)", i+1, recipe.Name, len(recipe.Steps)))
				_, err = repo.CreateRecipe(ctx, recipe)
				if err != nil {
					return fmt.Errorf("failed to create recipe %s: %w", recipe.Name, err)
				}
			}

			logger.Info("All bootstrap recipes created successfully!")

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
				for _, name := range chickenMealNames {
					if meal.Name == name {
						chickenMeals = append(chickenMeals, meal)
						break
					}
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
					ID:                     identifiers.New(),
					MealID:                 chickenMeal.ID,
					MealScale:              1.0,
					Notes:                  "",
					AssignedCook:           nil,
					AssignedDishwasher:     nil,
					BelongsToMealPlanEvent: "", // Will be set by converter
				})
			}

			var otherOptions []*mealplanning.MealPlanOptionDatabaseCreationInput
			for _, otherMeal := range otherMeals {
				otherOptions = append(otherOptions, &mealplanning.MealPlanOptionDatabaseCreationInput{
					ID:                     identifiers.New(),
					MealID:                 otherMeal.ID,
					MealScale:              1.0,
					Notes:                  "",
					AssignedCook:           nil,
					AssignedDishwasher:     nil,
					BelongsToMealPlanEvent: "", // Will be set by converter
				})
			}

			// Create a single event with all three chickenOptions
			events := []*mealplanning.MealPlanEventDatabaseCreationInput{
				{
					ID:                identifiers.New(),
					StartsAt:          eventStart,
					EndsAt:            eventEnd,
					MealName:          mealplanning.DinnerMealName,
					Notes:             "",
					BelongsToMealPlan: "", // Will be set by converter
					Options:           chickenOptions,
				},
				{
					ID:                identifiers.New(),
					StartsAt:          cloneTime(eventStart).Add(24 * time.Hour),
					EndsAt:            cloneTime(eventEnd).Add(24 * time.Hour),
					MealName:          mealplanning.SupperMealName,
					Notes:             "",
					BelongsToMealPlan: "", // Will be set by converter
					Options:           otherOptions,
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

			logger.Info(fmt.Sprintf("Created meal plan %s with %d events", createdMealPlan.ID, len(events)))
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
