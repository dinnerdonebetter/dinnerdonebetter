package tracing

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/prixfixeco/backend/internal/observability/keys"
	"github.com/prixfixeco/backend/pkg/types"

	useragent "github.com/mssola/user_agent"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// AttachIntToSpan attaches an int to a span.
func AttachIntToSpan(span trace.Span, attachmentKey string, id int) {
	if span != nil {
		span.SetAttributes(attribute.Int64(attachmentKey, int64(id)))
	}
}

// AttachUint8ToSpan attaches a uint8 to a span.
func AttachUint8ToSpan(span trace.Span, attachmentKey string, id uint8) {
	if span != nil {
		span.SetAttributes(attribute.Int64(attachmentKey, int64(id)))
	}
}

// AttachUint16ToSpan attaches a uint16 to a span.
func AttachUint16ToSpan(span trace.Span, attachmentKey string, id uint16) {
	if span != nil {
		span.SetAttributes(attribute.Int64(attachmentKey, int64(id)))
	}
}

// AttachUint64ToSpan attaches a uint64 to a span.
func AttachUint64ToSpan(span trace.Span, attachmentKey string, id uint64) {
	if span != nil {
		span.SetAttributes(attribute.Int64(attachmentKey, int64(id)))
	}
}

// AttachStringToSpan attaches a string to a span.
func AttachStringToSpan(span trace.Span, key, str string) {
	if span != nil {
		span.SetAttributes(attribute.String(key, str))
	}
}

// AttachBooleanToSpan attaches a boolean to a span.
func AttachBooleanToSpan(span trace.Span, key string, b bool) {
	if span != nil {
		span.SetAttributes(attribute.Bool(key, b))
	}
}

// AttachSliceOfStringsToSpan attaches a slice of strings to a span.
func AttachSliceOfStringsToSpan(span trace.Span, key string, slice []string) {
	if span != nil {
		span.SetAttributes(attribute.StringSlice(key, slice))
	}
}

// AttachTimeToSpan attaches a uint64 to a span.
func AttachTimeToSpan(span trace.Span, attachmentKey string, t time.Time) {
	AttachStringToSpan(span, attachmentKey, t.Format(time.RFC3339Nano))
}

// AttachFilterDataToSpan provides a consistent way to attach a filter's info to a span.
func AttachFilterDataToSpan(span trace.Span, page *uint16, limit *uint8, sortBy *string) {
	if page != nil {
		AttachUint16ToSpan(span, keys.FilterPageKey, *page)
	}

	if limit != nil {
		AttachUint8ToSpan(span, keys.FilterLimitKey, *limit)
	}

	if sortBy != nil {
		AttachStringToSpan(span, keys.FilterSortByKey, *sortBy)
	}
}

// AttachEmailAddressToSpan provides a consistent way to attach a household's ID to a span.
func AttachEmailAddressToSpan(span trace.Span, emailAddress string) {
	AttachStringToSpan(span, keys.UserEmailAddressKey, emailAddress)
}

// AttachHouseholdIDToSpan provides a consistent way to attach a household's ID to a span.
func AttachHouseholdIDToSpan(span trace.Span, householdID string) {
	AttachStringToSpan(span, keys.HouseholdIDKey, householdID)
}

// AttachHouseholdInvitationIDToSpan provides a consistent way to attach a household's ID to a span.
func AttachHouseholdInvitationIDToSpan(span trace.Span, householdInvitationID string) {
	AttachStringToSpan(span, keys.HouseholdInvitationIDKey, householdInvitationID)
}

// AttachHouseholdInvitationTokenToSpan provides a consistent way to attach a household's ID to a span.
func AttachHouseholdInvitationTokenToSpan(span trace.Span, householdInvitationTokenID string) {
	AttachStringToSpan(span, keys.HouseholdInvitationTokenKey, householdInvitationTokenID)
}

// AttachActiveHouseholdIDToSpan provides a consistent way to attach a household's ID to a span.
func AttachActiveHouseholdIDToSpan(span trace.Span, householdID string) {
	AttachStringToSpan(span, keys.ActiveHouseholdIDKey, householdID)
}

// AttachRequestingUserIDToSpan provides a consistent way to attach a user's ID to a span.
func AttachRequestingUserIDToSpan(span trace.Span, userID string) {
	AttachStringToSpan(span, keys.RequesterIDKey, userID)
}

// AttachSessionContextDataToSpan provides a consistent way to attach a SessionContextData object to a span.
func AttachSessionContextDataToSpan(span trace.Span, sessionCtxData *types.SessionContextData) {
	if sessionCtxData != nil {
		AttachRequestingUserIDToSpan(span, sessionCtxData.Requester.UserID)
		AttachActiveHouseholdIDToSpan(span, sessionCtxData.ActiveHouseholdID)
		if sessionCtxData.Requester.ServicePermissions != nil {
			AttachBooleanToSpan(span, keys.UserIsServiceAdminKey, sessionCtxData.Requester.ServicePermissions.IsServiceAdmin())
		}
	}
}

// AttachAPIClientDatabaseIDToSpan is a consistent way to attach an API client's database row ID to a span.
func AttachAPIClientDatabaseIDToSpan(span trace.Span, clientID string) {
	AttachStringToSpan(span, keys.APIClientDatabaseIDKey, clientID)
}

// AttachAPIClientClientIDToSpan is a consistent way to attach an API client's ID to a span.
func AttachAPIClientClientIDToSpan(span trace.Span, clientID string) {
	AttachStringToSpan(span, keys.APIClientClientIDKey, clientID)
}

// AttachUserToSpan provides a consistent way to attach a user to a span.
func AttachUserToSpan(span trace.Span, user *types.User) {
	if user != nil {
		AttachUserIDToSpan(span, user.ID)
		AttachUsernameToSpan(span, user.Username)
	}
}

// AttachUserIDToSpan provides a consistent way to attach a user's ID to a span.
func AttachUserIDToSpan(span trace.Span, userID string) {
	AttachStringToSpan(span, keys.UserIDKey, userID)
}

// AttachUsernameToSpan provides a consistent way to attach a user's username to a span.
func AttachUsernameToSpan(span trace.Span, username string) {
	AttachStringToSpan(span, keys.UsernameKey, username)
}

// AttachWebhookIDToSpan provides a consistent way to attach a webhook's ID to a span.
func AttachWebhookIDToSpan(span trace.Span, webhookID string) {
	AttachStringToSpan(span, keys.WebhookIDKey, webhookID)
}

// AttachURLToSpan attaches a given URI to a span.
func AttachURLToSpan(span trace.Span, u *url.URL) {
	AttachStringToSpan(span, keys.RequestURIKey, u.String())
}

// AttachRequestURIToSpan attaches a given URI to a span.
func AttachRequestURIToSpan(span trace.Span, uri string) {
	AttachStringToSpan(span, keys.RequestURIKey, uri)
}

// AttachRequestToSpan attaches a given *http.Request to a span.
func AttachRequestToSpan(span trace.Span, req *http.Request) {
	if req != nil {
		AttachStringToSpan(span, keys.RequestURIKey, req.URL.String())
		AttachStringToSpan(span, keys.RequestMethodKey, req.Method)

		for k, v := range req.Header {
			AttachSliceOfStringsToSpan(span, fmt.Sprintf("%s.%s", keys.RequestHeadersKey, k), v)
		}
	}
}

// AttachResponseToSpan attaches a given *http.Response to a span.
func AttachResponseToSpan(span trace.Span, res *http.Response) {
	if res != nil {
		AttachRequestToSpan(span, res.Request)

		span.SetAttributes(attribute.Int(keys.ResponseStatusKey, res.StatusCode))

		for k, v := range res.Header {
			AttachSliceOfStringsToSpan(span, fmt.Sprintf("%s.%s", keys.ResponseHeadersKey, k), v)
		}
	}
}

// AttachErrorToSpan attaches a given error to a span.
func AttachErrorToSpan(span trace.Span, description string, err error) {
	if err != nil {
		span.RecordError(
			err,
			trace.WithTimestamp(time.Now()),
			trace.WithAttributes(attribute.String("error.description", description)),
		)
	}
}

// AttachDatabaseQueryToSpan attaches a given search query to a span.
func AttachDatabaseQueryToSpan(span trace.Span, queryDescription, query string, args []any) {
	AttachStringToSpan(span, keys.DatabaseQueryKey, query)
	AttachStringToSpan(span, "query_description", queryDescription)

	for i, arg := range args {
		span.SetAttributes(attribute.String(fmt.Sprintf("query_args_%d", i), fmt.Sprintf("%+v", arg)))
	}
}

// AttachQueryFilterToSpan attaches a given query filter to a span.
func AttachQueryFilterToSpan(span trace.Span, filter *types.QueryFilter) {
	if filter != nil {
		if filter.Limit != nil {
			AttachUint8ToSpan(span, keys.FilterLimitKey, *filter.Limit)
		}

		if filter.Page != nil {
			AttachUint16ToSpan(span, keys.FilterPageKey, *filter.Page)
		}

		if filter.CreatedAfter != nil {
			AttachTimeToSpan(span, keys.FilterCreatedAfterKey, *filter.CreatedAfter)
		}

		if filter.CreatedBefore != nil {
			AttachTimeToSpan(span, keys.FilterCreatedBeforeKey, *filter.CreatedBefore)
		}

		if filter.UpdatedAfter != nil {
			AttachTimeToSpan(span, keys.FilterUpdatedAfterKey, *filter.UpdatedAfter)
		}

		if filter.UpdatedBefore != nil {
			AttachTimeToSpan(span, keys.FilterUpdatedBeforeKey, *filter.UpdatedBefore)
		}

		if filter.SortBy != nil {
			AttachStringToSpan(span, keys.FilterSortByKey, *filter.SortBy)
		}
	} else {
		AttachBooleanToSpan(span, keys.FilterIsNilKey, true)
	}
}

// AttachSearchQueryToSpan attaches a given search query to a span.
func AttachSearchQueryToSpan(span trace.Span, query string) {
	AttachStringToSpan(span, keys.SearchQueryKey, query)
}

// AttachUserAgentDataToSpan attaches a given search query to a span.
func AttachUserAgentDataToSpan(span trace.Span, ua *useragent.UserAgent) {
	if ua != nil {
		AttachStringToSpan(span, keys.UserAgentOSKey, ua.OS())
		AttachBooleanToSpan(span, keys.UserAgentMobileKey, ua.Mobile())
		AttachBooleanToSpan(span, keys.UserAgentBotKey, ua.Bot())
	}
}

// AttachValidInstrumentIDToSpan attaches a valid instrument ID to a given span.
func AttachValidInstrumentIDToSpan(span trace.Span, validInstrumentID string) {
	AttachStringToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)
}

// AttachValidIngredientIDToSpan attaches a valid ingredient ID to a given span.
func AttachValidIngredientIDToSpan(span trace.Span, validIngredientID string) {
	AttachStringToSpan(span, keys.ValidIngredientIDKey, validIngredientID)
}

// AttachValidPreparationIDToSpan attaches a valid preparation ID to a given span.
func AttachValidPreparationIDToSpan(span trace.Span, validPreparationID string) {
	AttachStringToSpan(span, keys.ValidPreparationIDKey, validPreparationID)
}

// AttachValidIngredientPreparationIDToSpan attaches a valid ingredient preparation ID to a given span.
func AttachValidIngredientPreparationIDToSpan(span trace.Span, validIngredientPreparationID string) {
	AttachStringToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
}

// AttachValidIngredientStateIngredientIDToSpan attaches a valid ingredient state ingredient ID to a given span.
func AttachValidIngredientStateIngredientIDToSpan(span trace.Span, validIngredientStateIngredientID string) {
	AttachStringToSpan(span, keys.ValidIngredientStateIngredientIDKey, validIngredientStateIngredientID)
}

// AttachValidPreparationInstrumentIDToSpan attaches a valid preparation instrument ID to a given span.
func AttachValidPreparationInstrumentIDToSpan(span trace.Span, validPreparationInstrumentID string) {
	AttachStringToSpan(span, keys.ValidPreparationInstrumentIDKey, validPreparationInstrumentID)
}

// AttachValidIngredientMeasurementUnitIDToSpan attaches a valid ingredient measurement unit ID to a given span.
func AttachValidIngredientMeasurementUnitIDToSpan(span trace.Span, validIngredientMeasurementUnitID string) {
	AttachStringToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
}

// AttachMealIDToSpan attaches a recipe ID to a given span.
func AttachMealIDToSpan(span trace.Span, mealID string) {
	AttachStringToSpan(span, keys.MealIDKey, mealID)
}

// AttachRecipeIDToSpan attaches a recipe ID to a given span.
func AttachRecipeIDToSpan(span trace.Span, recipeID string) {
	AttachStringToSpan(span, keys.RecipeIDKey, recipeID)
}

// AttachRecipePrepTaskIDToSpan attaches a recipe prep task ID to a given span.
func AttachRecipePrepTaskIDToSpan(span trace.Span, recipePrepTaskID string) {
	AttachStringToSpan(span, keys.RecipePrepTaskIDKey, recipePrepTaskID)
}

// AttachValidIngredientStateIDToSpan attaches a valid ingredient state ID to a given span.
func AttachValidIngredientStateIDToSpan(span trace.Span, validIngredientStateID string) {
	AttachStringToSpan(span, keys.ValidIngredientStateIDKey, validIngredientStateID)
}

// AttachRecipeStepIDToSpan attaches a recipe step ID to a given span.
func AttachRecipeStepIDToSpan(span trace.Span, recipeStepID string) {
	AttachStringToSpan(span, keys.RecipeStepIDKey, recipeStepID)
}

// AttachRecipeStepInstrumentIDToSpan attaches a recipe step instrument ID to a given span.
func AttachRecipeStepInstrumentIDToSpan(span trace.Span, recipeStepInstrumentID string) {
	AttachStringToSpan(span, keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
}

// AttachRecipeStepVesselIDToSpan attaches a recipe step vessel ID to a given span.
func AttachRecipeStepVesselIDToSpan(span trace.Span, recipeStepVesselID string) {
	AttachStringToSpan(span, keys.RecipeStepVesselIDKey, recipeStepVesselID)
}

// AttachRecipeStepIngredientIDToSpan attaches a recipe step ingredient ID to a given span.
func AttachRecipeStepIngredientIDToSpan(span trace.Span, recipeStepIngredientID string) {
	AttachStringToSpan(span, keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
}

// AttachRecipeStepCompletionConditionIDToSpan attaches a recipe step completion condition ID to a given span.
func AttachRecipeStepCompletionConditionIDToSpan(span trace.Span, recipeStepCompletionConditionID string) {
	AttachStringToSpan(span, keys.RecipeStepCompletionConditionIDKey, recipeStepCompletionConditionID)
}

// AttachRecipeStepProductIDToSpan attaches a recipe step product ID to a given span.
func AttachRecipeStepProductIDToSpan(span trace.Span, recipeStepProductID string) {
	AttachStringToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)
}

// AttachMealPlanIDToSpan attaches a meal plan ID to a given span.
func AttachMealPlanIDToSpan(span trace.Span, mealPlanID string) {
	AttachStringToSpan(span, keys.MealPlanIDKey, mealPlanID)
}

// AttachMealPlanEventIDToSpan attaches a meal plan ID to a given span.
func AttachMealPlanEventIDToSpan(span trace.Span, mealPlanEventID string) {
	AttachStringToSpan(span, keys.MealPlanEventIDKey, mealPlanEventID)
}

// AttachMealPlanOptionIDToSpan attaches a meal plan option ID to a given span.
func AttachMealPlanOptionIDToSpan(span trace.Span, mealPlanOptionID string) {
	AttachStringToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)
}

// AttachMealPlanOptionVoteIDToSpan attaches a meal plan option vote ID to a given span.
func AttachMealPlanOptionVoteIDToSpan(span trace.Span, mealPlanOptionVoteID string) {
	AttachStringToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
}

// AttachPasswordResetTokenIDToSpan attaches a password reset token ID to a given span.
func AttachPasswordResetTokenIDToSpan(span trace.Span, passwordResetTokenID string) {
	AttachStringToSpan(span, keys.PasswordResetTokenIDKey, passwordResetTokenID)
}

// AttachPasswordResetTokenToSpan attaches a password reset token to a given span.
func AttachPasswordResetTokenToSpan(span trace.Span, passwordResetToken string) {
	AttachStringToSpan(span, keys.PasswordResetTokenIDKey, passwordResetToken)
}

// AttachValidMeasurementUnitIDToSpan attaches a valid measurement unit ID to a given span.
func AttachValidMeasurementUnitIDToSpan(span trace.Span, validMeasurementUnitID string) {
	AttachStringToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
}

// AttachMealPlanTaskIDToSpan attaches a meal plan task ID to a given span.
func AttachMealPlanTaskIDToSpan(span trace.Span, mealPlanTaskID string) {
	AttachStringToSpan(span, keys.MealPlanTaskIDKey, mealPlanTaskID)
}

// AttachMealPlanGroceryListItemIDToSpan attaches a meal plan task ID to a given span.
func AttachMealPlanGroceryListItemIDToSpan(span trace.Span, mealPlanGroceryListItemID string) {
	AttachStringToSpan(span, keys.MealPlanGroceryListItemIDKey, mealPlanGroceryListItemID)
}

// AttachValidMeasurementConversionIDToSpan attaches a valid measurement conversion ID to a given span.
func AttachValidMeasurementConversionIDToSpan(span trace.Span, validMeasurementConversionID string) {
	AttachStringToSpan(span, keys.ValidMeasurementConversionIDKey, validMeasurementConversionID)
}

// AttachRecipeMediaIDToSpan attaches a recipe media ID to a given span.
func AttachRecipeMediaIDToSpan(span trace.Span, recipeMediaID string) {
	AttachStringToSpan(span, keys.RecipeMediaIDKey, recipeMediaID)
}
