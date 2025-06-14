package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.DataPrivacyDataManager = (*Querier)(nil)
)

// DeleteUser archives a user.
func (q *Querier) DeleteUser(ctx context.Context, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	changed, err := q.generatedQuerier.DeleteUser(ctx, q.db, userID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving user")
	}

	if changed == 0 {
		return sql.ErrNoRows
	}

	logger.Info("user deleted")

	return nil
}

// AggregateUserData collects all of a user's data.
func (q *Querier) AggregateUserData(ctx context.Context, userID string) (*types.UserDataCollection, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	user, err := q.GetUser(ctx, userID)
	if err != nil || user == nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("getting user: %w", err)
	}

	// TODO: var outputWG sync.WG; var outputLock sync.Mutex; go func() {}()

	allAccounts, err := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Account], error) {
		return q.getAccountsForUser(ctx, q.db, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting accounts: %w", err)
	}

	allUserAuditLogEntries, err := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.AuditLogEntry], error) {
		return q.GetAuditLogEntriesForUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting user audit log entries: %w", err)
	}

	allSettingConfigs, err := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
		return q.GetServiceSettingConfigurationsForUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting service setting configurations: %w", err)
	}

	userIngredientPreferences, err := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.UserIngredientPreference], error) {
		return q.GetUserIngredientPreferences(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting user ingredient preferences: %w", err)
	}

	receivedInvites, err := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.AccountInvitation], error) {
		return q.GetPendingAccountInvitationsForUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting received invites: %w", err)
	}

	sentInvites, err := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.AccountInvitation], error) {
		return q.GetPendingAccountInvitationsFromUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting sent invites: %w", err)
	}

	recipes, err := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Recipe], error) {
		return q.GetRecipesCreatedByUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting recipes: %w", err)
	}

	recipeRatings, err := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.RecipeRating], error) {
		return q.GetRecipeRatingsForUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting recipe ratings: %w", err)
	}

	meals, err := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Meal], error) {
		return q.GetMealsCreatedByUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting meals: %w", err)
	}

	output := &types.UserDataCollection{
		User:     *user,
		ReportID: identifiers.New(),
		Core: types.CoreUserDataCollection{
			Accounts:                         allAccounts,
			UserAuditLogEntries:              allUserAuditLogEntries,
			UserServiceSettingConfigurations: allSettingConfigs,
			ReceivedInvites:                  receivedInvites,
			SentInvites:                      sentInvites,
			AuditLogEntries:                  map[string][]types.AuditLogEntry{},
			ServiceSettingConfigurations:     map[string][]types.ServiceSettingConfiguration{},
			Webhooks:                         map[string][]types.Webhook{},
		},
		Eating: types.EatingUserDataCollection{
			UserIngredientPreferences:   userIngredientPreferences,
			RecipeRatings:               recipeRatings,
			Recipes:                     recipes,
			Meals:                       meals,
			MealPlans:                   map[string][]types.MealPlan{},
			AccountInstrumentOwnerships: map[string][]types.AccountInstrumentOwnership{},
		},
	}

	// set up data collections for all accounts
	for i := range allAccounts {
		account := allAccounts[i]
		auditLogEntries, fetchErr := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.AuditLogEntry], error) {
			return q.GetAuditLogEntriesForAccount(ctx, account.ID, filter)
		})
		if fetchErr != nil {
			return nil, fmt.Errorf("fetching audit log entries for account %s", account.ID)
		}
		output.Core.AuditLogEntries[account.ID] = auditLogEntries

		serviceSettingConfigs, fetchErr := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
			return q.GetServiceSettingConfigurationsForAccount(ctx, account.ID, filter)
		})
		if fetchErr != nil {
			return nil, fmt.Errorf("fetching audit log entries for account %s", account.ID)
		}
		output.Core.ServiceSettingConfigurations[account.ID] = serviceSettingConfigs

		webhooks, fetchErr := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.Webhook], error) {
			return q.GetWebhooks(ctx, account.ID, filter)
		})
		if fetchErr != nil {
			return nil, fmt.Errorf("fetching audit log entries for account %s", account.ID)
		}
		output.Core.Webhooks[account.ID] = webhooks

		accountInstrumentOwnerships, fetchErr := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.AccountInstrumentOwnership], error) {
			return q.GetAccountInstrumentOwnerships(ctx, account.ID, filter)
		})
		if fetchErr != nil {
			return nil, fmt.Errorf("fetching audit log entries for account %s", account.ID)
		}
		output.Eating.AccountInstrumentOwnerships[account.ID] = accountInstrumentOwnerships

		mealPlans, fetchErr := fetchAllRows(func(filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.MealPlan], error) {
			return q.GetMealPlansForAccount(ctx, account.ID, filter)
		})
		if fetchErr != nil {
			return nil, fmt.Errorf("fetching audit log entries for account %s", account.ID)
		}
		output.Eating.MealPlans[account.ID] = mealPlans
	}

	logger.Info("user data collected")

	return output, nil
}
