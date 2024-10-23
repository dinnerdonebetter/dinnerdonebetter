package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
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

	allHouseholds, err := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.Household], error) {
		return q.getHouseholdsForUser(ctx, q.db, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting households: %w", err)
	}

	allUserAuditLogEntries, err := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
		return q.GetAuditLogEntriesForUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting user audit log entries: %w", err)
	}

	allSettingConfigs, err := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
		return q.GetServiceSettingConfigurationsForUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting service setting configurations: %w", err)
	}

	userIngredientPreferences, err := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.UserIngredientPreference], error) {
		return q.GetUserIngredientPreferences(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting user ingredient preferences: %w", err)
	}

	receivedInvites, err := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.HouseholdInvitation], error) {
		return q.GetPendingHouseholdInvitationsForUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting received invites: %w", err)
	}

	sentInvites, err := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.HouseholdInvitation], error) {
		return q.GetPendingHouseholdInvitationsFromUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting sent invites: %w", err)
	}

	recipes, err := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.Recipe], error) {
		return q.GetRecipesCreatedByUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting recipes: %w", err)
	}

	recipeRatings, err := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeRating], error) {
		return q.GetRecipeRatingsForUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting recipe ratings: %w", err)
	}

	meals, err := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.Meal], error) {
		return q.GetMealsCreatedByUser(ctx, userID, filter)
	})
	if err != nil {
		return nil, fmt.Errorf("getting meals: %w", err)
	}

	output := &types.UserDataCollection{
		User:                             *user,
		Households:                       allHouseholds,
		UserAuditLogEntries:              allUserAuditLogEntries,
		UserServiceSettingConfigurations: allSettingConfigs,
		UserIngredientPreferences:        userIngredientPreferences,
		ReceivedInvites:                  receivedInvites,
		SentInvites:                      sentInvites,
		RecipeRatings:                    recipeRatings,
		Recipes:                          recipes,
		Meals:                            meals,
		AuditLogEntries:                  map[string][]types.AuditLogEntry{},
		MealPlans:                        map[string][]types.MealPlan{},
		HouseholdInstrumentOwnerships:    map[string][]types.HouseholdInstrumentOwnership{},
		ServiceSettingConfigurations:     map[string][]types.ServiceSettingConfiguration{},
		Webhooks:                         map[string][]types.Webhook{},
	}

	// set up data collections for all households
	for i := range allHouseholds {
		household := allHouseholds[i]
		auditLogEntries, fetchErr := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.AuditLogEntry], error) {
			return q.GetAuditLogEntriesForHousehold(ctx, household.ID, filter)
		})
		if fetchErr != nil {
			return nil, fmt.Errorf("fetching audit log entries for household %s", household.ID)
		}
		output.AuditLogEntries[household.ID] = auditLogEntries

		serviceSettingConfigs, fetchErr := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.ServiceSettingConfiguration], error) {
			return q.GetServiceSettingConfigurationsForHousehold(ctx, household.ID, filter)
		})
		if fetchErr != nil {
			return nil, fmt.Errorf("fetching audit log entries for household %s", household.ID)
		}
		output.ServiceSettingConfigurations[household.ID] = serviceSettingConfigs

		webhooks, fetchErr := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.Webhook], error) {
			return q.GetWebhooks(ctx, household.ID, filter)
		})
		if fetchErr != nil {
			return nil, fmt.Errorf("fetching audit log entries for household %s", household.ID)
		}
		output.Webhooks[household.ID] = webhooks

		householdInstrumentOwnerships, fetchErr := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.HouseholdInstrumentOwnership], error) {
			return q.GetHouseholdInstrumentOwnerships(ctx, household.ID, filter)
		})
		if fetchErr != nil {
			return nil, fmt.Errorf("fetching audit log entries for household %s", household.ID)
		}
		output.HouseholdInstrumentOwnerships[household.ID] = householdInstrumentOwnerships

		mealPlans, fetchErr := fetchAllRows(func(filter *types.QueryFilter) (*types.QueryFilteredResult[types.MealPlan], error) {
			return q.GetMealPlansForHousehold(ctx, household.ID, filter)
		})
		if fetchErr != nil {
			return nil, fmt.Errorf("fetching audit log entries for household %s", household.ID)
		}
		output.MealPlans[household.ID] = mealPlans
	}

	logger.Info("user data collected")

	return output, nil
}