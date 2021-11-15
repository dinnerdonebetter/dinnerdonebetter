package tracing

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	useragent "github.com/mssola/user_agent"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/pkg/types"
)

func attachUint8ToSpan(span trace.Span, attachmentKey string, id uint8) {
	if span != nil {
		span.SetAttributes(attribute.Int64(attachmentKey, int64(id)))
	}
}

func attachUint64ToSpan(span trace.Span, attachmentKey string, id uint64) {
	if span != nil {
		span.SetAttributes(attribute.Int64(attachmentKey, int64(id)))
	}
}

func attachStringToSpan(span trace.Span, key, str string) {
	if span != nil {
		span.SetAttributes(attribute.String(key, str))
	}
}

func attachBooleanToSpan(span trace.Span, key string, b bool) {
	if span != nil {
		span.SetAttributes(attribute.Bool(key, b))
	}
}

func attachSliceToSpan(span trace.Span, key string, slice interface{}) {
	span.SetAttributes(attribute.Array(key, slice))
}

// AttachToSpan allows a user to attach any value to a span.
func AttachToSpan(span trace.Span, key string, val interface{}) {
	if span != nil {
		span.SetAttributes(attribute.Any(key, val))
	}
}

// AttachFilterToSpan provides a consistent way to attach a filter's info to a span.
func AttachFilterToSpan(span trace.Span, page uint64, limit uint8, sortBy string) {
	attachUint64ToSpan(span, keys.FilterPageKey, page)
	attachUint8ToSpan(span, keys.FilterLimitKey, limit)
	attachStringToSpan(span, keys.FilterSortByKey, sortBy)
}

// AttachEmailAddressToSpan provides a consistent way to attach a household's ID to a span.
func AttachEmailAddressToSpan(span trace.Span, emailAddress string) {
	attachStringToSpan(span, keys.UserEmailAddressKey, emailAddress)
}

// AttachHouseholdIDToSpan provides a consistent way to attach a household's ID to a span.
func AttachHouseholdIDToSpan(span trace.Span, householdID string) {
	attachStringToSpan(span, keys.HouseholdIDKey, householdID)
}

// AttachHouseholdInvitationIDToSpan provides a consistent way to attach a household's ID to a span.
func AttachHouseholdInvitationIDToSpan(span trace.Span, householdInvitationID string) {
	attachStringToSpan(span, keys.HouseholdInvitationIDKey, householdInvitationID)
}

// AttachHouseholdInvitationTokenToSpan provides a consistent way to attach a household's ID to a span.
func AttachHouseholdInvitationTokenToSpan(span trace.Span, householdInvitationTokenID string) {
	attachStringToSpan(span, keys.HouseholdInvitationTokenKey, householdInvitationTokenID)
}

// AttachActiveHouseholdIDToSpan provides a consistent way to attach a household's ID to a span.
func AttachActiveHouseholdIDToSpan(span trace.Span, householdID string) {
	attachStringToSpan(span, keys.ActiveHouseholdIDKey, householdID)
}

// AttachRequestingUserIDToSpan provides a consistent way to attach a user's ID to a span.
func AttachRequestingUserIDToSpan(span trace.Span, userID string) {
	attachStringToSpan(span, keys.RequesterIDKey, userID)
}

// AttachSessionContextDataToSpan provides a consistent way to attach a SessionContextData object to a span.
func AttachSessionContextDataToSpan(span trace.Span, sessionCtxData *types.SessionContextData) {
	if sessionCtxData != nil {
		AttachRequestingUserIDToSpan(span, sessionCtxData.Requester.UserID)
		AttachActiveHouseholdIDToSpan(span, sessionCtxData.ActiveHouseholdID)
		if sessionCtxData.Requester.ServicePermissions != nil {
			attachBooleanToSpan(span, keys.UserIsServiceAdminKey, sessionCtxData.Requester.ServicePermissions.IsServiceAdmin())
		}
	}
}

// AttachAPIClientDatabaseIDToSpan is a consistent way to attach an API client's database row ID to a span.
func AttachAPIClientDatabaseIDToSpan(span trace.Span, clientID string) {
	attachStringToSpan(span, keys.APIClientDatabaseIDKey, clientID)
}

// AttachAPIClientClientIDToSpan is a consistent way to attach an API client's ID to a span.
func AttachAPIClientClientIDToSpan(span trace.Span, clientID string) {
	attachStringToSpan(span, keys.APIClientClientIDKey, clientID)
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
	attachStringToSpan(span, keys.UserIDKey, userID)
}

// AttachUsernameToSpan provides a consistent way to attach a user's username to a span.
func AttachUsernameToSpan(span trace.Span, username string) {
	attachStringToSpan(span, keys.UsernameKey, username)
}

// AttachWebhookIDToSpan provides a consistent way to attach a webhook's ID to a span.
func AttachWebhookIDToSpan(span trace.Span, webhookID string) {
	attachStringToSpan(span, keys.WebhookIDKey, webhookID)
}

// AttachURLToSpan attaches a given URI to a span.
func AttachURLToSpan(span trace.Span, u *url.URL) {
	attachStringToSpan(span, keys.RequestURIKey, u.String())
}

// AttachRequestURIToSpan attaches a given URI to a span.
func AttachRequestURIToSpan(span trace.Span, uri string) {
	attachStringToSpan(span, keys.RequestURIKey, uri)
}

// AttachRequestToSpan attaches a given *http.Request to a span.
func AttachRequestToSpan(span trace.Span, req *http.Request) {
	if req != nil {
		attachStringToSpan(span, keys.RequestURIKey, req.URL.String())
		attachStringToSpan(span, keys.RequestMethodKey, req.Method)

		for k, v := range req.Header {
			attachSliceToSpan(span, fmt.Sprintf("%s.%s", keys.RequestHeadersKey, k), v)
		}
	}
}

// AttachResponseToSpan attaches a given *http.Response to a span.
func AttachResponseToSpan(span trace.Span, res *http.Response) {
	if res != nil {
		AttachRequestToSpan(span, res.Request)

		span.SetAttributes(attribute.Int(keys.ResponseStatusKey, res.StatusCode))

		for k, v := range res.Header {
			attachSliceToSpan(span, fmt.Sprintf("%s.%s", keys.ResponseHeadersKey, k), v)
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
func AttachDatabaseQueryToSpan(span trace.Span, queryDescription, query string, args []interface{}) {
	attachStringToSpan(span, keys.DatabaseQueryKey, query)
	attachStringToSpan(span, "query_description", queryDescription)

	for i, arg := range args {
		span.SetAttributes(attribute.Any(fmt.Sprintf("query_args_%d", i), arg))
	}
}

// AttachQueryFilterToSpan attaches a given query filter to a span.
func AttachQueryFilterToSpan(span trace.Span, filter *types.QueryFilter) {
	if filter != nil {
		attachUint8ToSpan(span, keys.FilterLimitKey, filter.Limit)
		attachUint64ToSpan(span, keys.FilterPageKey, filter.Page)
		attachUint64ToSpan(span, keys.FilterCreatedAfterKey, filter.CreatedAfter)
		attachUint64ToSpan(span, keys.FilterCreatedBeforeKey, filter.CreatedBefore)
		attachUint64ToSpan(span, keys.FilterUpdatedAfterKey, filter.UpdatedAfter)
		attachUint64ToSpan(span, keys.FilterUpdatedBeforeKey, filter.UpdatedBefore)
		attachStringToSpan(span, keys.FilterSortByKey, string(filter.SortBy))
	} else {
		attachBooleanToSpan(span, keys.FilterIsNilKey, true)
	}
}

// AttachSearchQueryToSpan attaches a given search query to a span.
func AttachSearchQueryToSpan(span trace.Span, query string) {
	attachStringToSpan(span, keys.SearchQueryKey, query)
}

// AttachUserAgentDataToSpan attaches a given search query to a span.
func AttachUserAgentDataToSpan(span trace.Span, ua *useragent.UserAgent) {
	if ua != nil {
		attachStringToSpan(span, keys.UserAgentOSKey, ua.OS())
		attachBooleanToSpan(span, keys.UserAgentMobileKey, ua.Mobile())
		attachBooleanToSpan(span, keys.UserAgentBotKey, ua.Bot())
	}
}

// AttachValidInstrumentIDToSpan attaches a valid instrument ID to a given span.
func AttachValidInstrumentIDToSpan(span trace.Span, validInstrumentID string) {
	attachStringToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)
}

// AttachValidIngredientIDToSpan attaches a valid ingredient ID to a given span.
func AttachValidIngredientIDToSpan(span trace.Span, validIngredientID string) {
	attachStringToSpan(span, keys.ValidIngredientIDKey, validIngredientID)
}

// AttachValidPreparationIDToSpan attaches a valid preparation ID to a given span.
func AttachValidPreparationIDToSpan(span trace.Span, validPreparationID string) {
	attachStringToSpan(span, keys.ValidPreparationIDKey, validPreparationID)
}

// AttachValidIngredientPreparationIDToSpan attaches a valid ingredient preparation ID to a given span.
func AttachValidIngredientPreparationIDToSpan(span trace.Span, validIngredientPreparationID string) {
	attachStringToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
}

// AttachRecipeIDToSpan attaches a recipe ID to a given span.
func AttachRecipeIDToSpan(span trace.Span, recipeID string) {
	attachStringToSpan(span, keys.RecipeIDKey, recipeID)
}

// AttachRecipeStepIDToSpan attaches a recipe step ID to a given span.
func AttachRecipeStepIDToSpan(span trace.Span, recipeStepID string) {
	attachStringToSpan(span, keys.RecipeStepIDKey, recipeStepID)
}

// AttachRecipeStepInstrumentIDToSpan attaches a recipe step instrument ID to a given span.
func AttachRecipeStepInstrumentIDToSpan(span trace.Span, recipeStepInstrumentID string) {
	attachStringToSpan(span, keys.RecipeStepInstrumentIDKey, recipeStepInstrumentID)
}

// AttachRecipeStepIngredientIDToSpan attaches a recipe step ingredient ID to a given span.
func AttachRecipeStepIngredientIDToSpan(span trace.Span, recipeStepIngredientID string) {
	attachStringToSpan(span, keys.RecipeStepIngredientIDKey, recipeStepIngredientID)
}

// AttachRecipeStepProductIDToSpan attaches a recipe step product ID to a given span.
func AttachRecipeStepProductIDToSpan(span trace.Span, recipeStepProductID string) {
	attachStringToSpan(span, keys.RecipeStepProductIDKey, recipeStepProductID)
}

// AttachMealPlanIDToSpan attaches a meal plan ID to a given span.
func AttachMealPlanIDToSpan(span trace.Span, mealPlanID string) {
	attachStringToSpan(span, keys.MealPlanIDKey, mealPlanID)
}

// AttachMealPlanOptionIDToSpan attaches a meal plan option ID to a given span.
func AttachMealPlanOptionIDToSpan(span trace.Span, mealPlanOptionID string) {
	attachStringToSpan(span, keys.MealPlanOptionIDKey, mealPlanOptionID)
}

// AttachMealPlanOptionVoteIDToSpan attaches a meal plan option vote ID to a given span.
func AttachMealPlanOptionVoteIDToSpan(span trace.Span, mealPlanOptionVoteID string) {
	attachStringToSpan(span, keys.MealPlanOptionVoteIDKey, mealPlanOptionVoteID)
}
