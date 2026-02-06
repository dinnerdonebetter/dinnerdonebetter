package dataprivacy

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/dataprivacy"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/backend/internal/domain/settings"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
)

// FetchUserDataCollection retrieves all user-associated data for GDPR/CCPA disclosure.
func (r *repository) FetchUserDataCollection(ctx context.Context, userID string) (*dataprivacy.UserDataCollection, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.WithValue("user_id", userID)
	logger.Info("fetching user data collection")

	collection := &dataprivacy.UserDataCollection{
		Identity:      identity.UserDataCollection{},
		MealPlanning:  mealplanning.UserDataCollection{},
		Webhooks:      webhooks.UserDataCollection{Data: make(map[string][]webhooks.Webhook)},
		Settings:      settings.UserDataCollection{},
		Notifications: notifications.UserDataCollection{},
	}

	// Fetch user profile
	user, err := r.identityRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching user")
	}
	collection.Identity.User = *user

	// Fetch user accounts
	accounts, err := r.identityRepo.GetAccounts(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching accounts")
	}
	for _, account := range accounts.Data {
		collection.Identity.Accounts = append(collection.Identity.Accounts, *account)
	}

	// Fetch account invitations (both sent and received)
	sentInvites, err := r.identityRepo.GetPendingAccountInvitationsFromUser(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching sent invitations")
	}
	receivedInvites, err := r.identityRepo.GetPendingAccountInvitationsForUser(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching received invitations")
	}
	// Combine and deduplicate invitations
	inviteMap := make(map[string]*identity.AccountInvitation)
	for _, invite := range sentInvites.Data {
		inviteMap[invite.ID] = invite
	}
	for _, invite := range receivedInvites.Data {
		inviteMap[invite.ID] = invite
	}
	for _, invite := range inviteMap {
		collection.Identity.AccountInvitations = append(collection.Identity.AccountInvitations, *invite)
	}

	// Fetch audit log entries
	auditLogs, err := r.auditLogRepo.GetAuditLogEntriesForUser(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching audit log entries")
	}
	for _, entry := range auditLogs.Data {
		collection.AuditLogEntries = append(collection.AuditLogEntries, *entry)
	}

	// Fetch meal planning data - user scoped
	recipes, err := r.mealPlanningRepo.GetRecipesCreatedByUser(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipes")
	}
	for _, recipe := range recipes.Data {
		collection.MealPlanning.Recipes = append(collection.MealPlanning.Recipes, *recipe)
	}

	meals, err := r.mealPlanningRepo.GetMealsCreatedByUser(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching meals")
	}
	for _, meal := range meals.Data {
		collection.MealPlanning.Meals = append(collection.MealPlanning.Meals, *meal)
	}

	preferences, err := r.mealPlanningRepo.GetUserIngredientPreferences(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching ingredient preferences")
	}
	for _, pref := range preferences.Data {
		collection.MealPlanning.UserIngredientPreferences = append(collection.MealPlanning.UserIngredientPreferences, *pref)
	}

	ratings, err := r.mealPlanningRepo.GetRecipeRatingsForUser(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching recipe ratings")
	}
	for _, rating := range ratings.Data {
		collection.MealPlanning.RecipeRatings = append(collection.MealPlanning.RecipeRatings, *rating)
	}

	// Fetch notifications
	notifs, err := r.notificationsRepo.GetUserNotifications(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching notifications")
	}
	for _, notif := range notifs.Data {
		collection.Notifications.Data = append(collection.Notifications.Data, *notif)
	}

	// Fetch user settings
	userSettings, err := r.settingsRepo.GetServiceSettingConfigurationsForUser(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching user settings")
	}
	for _, setting := range userSettings.Data {
		collection.Settings.UserSettings = append(collection.Settings.UserSettings, *setting)
	}

	// Fetch uploaded media
	media, err := r.uploadedMediaRepo.GetUploadedMediaForUser(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching uploaded media")
	}
	for _, m := range media.Data {
		collection.UploadedMedia = append(collection.UploadedMedia, *m)
	}

	// Fetch waitlist signups
	waitlistSignups, err := r.waitlistsRepo.GetWaitlistSignupsForUser(ctx, userID, nil)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching waitlist signups")
	}
	for _, signup := range waitlistSignups.Data {
		collection.WaitlistSignups = append(collection.WaitlistSignups, *signup)
	}

	// Fetch account-scoped data for each account
	for _, account := range accounts.Data {
		accountLogger := logger.WithValue("account_id", account.ID)

		// Meal plans
		mealPlans, mealPlanErr := r.mealPlanningRepo.GetMealPlansForAccount(ctx, account.ID, nil)
		if mealPlanErr != nil {
			return nil, observability.PrepareAndLogError(mealPlanErr, accountLogger, span, "fetching meal plans")
		}
		for _, mp := range mealPlans.Data {
			collection.MealPlanning.MealPlans = append(collection.MealPlanning.MealPlans, *mp)
		}

		// Account instrument ownerships
		ownerships, ownershipErr := r.mealPlanningRepo.GetAccountInstrumentOwnerships(ctx, account.ID, nil)
		if ownershipErr != nil {
			return nil, observability.PrepareAndLogError(ownershipErr, accountLogger, span, "fetching instrument ownerships")
		}
		for _, ownership := range ownerships.Data {
			collection.MealPlanning.AccountInstrumentOwnerships = append(collection.MealPlanning.AccountInstrumentOwnerships, *ownership)
		}

		// Webhooks
		hooks, webhookErr := r.webhooksRepo.GetWebhooks(ctx, account.ID, nil)
		if webhookErr != nil {
			return nil, observability.PrepareAndLogError(webhookErr, accountLogger, span, "fetching webhooks")
		}
		if len(hooks.Data) > 0 {
			var webhookList []webhooks.Webhook
			for _, hook := range hooks.Data {
				webhookList = append(webhookList, *hook)
			}
			collection.Webhooks.Data[account.ID] = webhookList
		}

		// Account settings
		accountSettings, settingsErr := r.settingsRepo.GetServiceSettingConfigurationsForAccount(ctx, account.ID, nil)
		if settingsErr != nil {
			return nil, observability.PrepareAndLogError(settingsErr, accountLogger, span, "fetching account settings")
		}
		for _, setting := range accountSettings.Data {
			collection.Settings.AccountSettings = append(collection.Settings.AccountSettings, *setting)
		}

		// Issue reports for account
		reports, reportsErr := r.issueReportsRepo.GetIssueReportsForAccount(ctx, account.ID, nil)
		if reportsErr != nil {
			return nil, observability.PrepareAndLogError(reportsErr, accountLogger, span, "fetching issue reports")
		}
		for _, report := range reports.Data {
			collection.IssueReports = append(collection.IssueReports, *report)
		}
	}

	logger.Info("user data collection complete")

	return collection, nil
}

// DeleteUser deletes a user and all associated data via ON DELETE CASCADE.
func (r *repository) DeleteUser(ctx context.Context, userID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.WithValue("user_id", userID)
	logger.Info("deleting user and all associated data")

	// The database schema uses ON DELETE CASCADE on all belongs_to_user foreign keys,
	// so deleting the user record will automatically delete all associated data.
	if err := r.identityRepo.DeleteUser(ctx, userID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "deleting user")
	}

	logger.Info("user deleted successfully")

	return nil
}
