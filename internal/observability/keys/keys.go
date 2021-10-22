package keys

const (
	// HouseholdSubscriptionPlanIDKey is the standard key for referring to an household subscription plan ID.
	HouseholdSubscriptionPlanIDKey = "household_subscription_plan.id"
	// PermissionsKey is the standard key for referring to an household user membership ID.
	PermissionsKey = "user.permissions"
	// RequesterIDKey is the standard key for referring to a requesting user's ID.
	RequesterIDKey = "request.made_by"
	// HouseholdIDKey is the standard key for referring to an household ID.
	HouseholdIDKey = "household.id"
	// ActiveHouseholdIDKey is the standard key for referring to an active household ID.
	ActiveHouseholdIDKey = "active_household_id"
	// UserIDKey is the standard key for referring to a user ID.
	UserIDKey = "user.id"
	// UserIsServiceAdminKey is the standard key for referring to a user's admin status.
	UserIsServiceAdminKey = "user.is_admin"
	// UsernameKey is the standard key for referring to a username.
	UsernameKey = "user.username"
	// ServiceRoleKey is the standard key for referring to a username.
	ServiceRoleKey = "user.service_role"
	// NameKey is the standard key for referring to a name.
	NameKey = "name"
	// FilterCreatedAfterKey is the standard key for referring to a types.QueryFilter's CreatedAfter field.
	FilterCreatedAfterKey = "query_filter.created_after"
	// FilterCreatedBeforeKey is the standard key for referring to a types.QueryFilter's CreatedBefore field.
	FilterCreatedBeforeKey = "query_filter.created_before"
	// FilterUpdatedAfterKey is the standard key for referring to a types.QueryFilter's UpdatedAfter field.
	FilterUpdatedAfterKey = "query_filter.updated_after"
	// FilterUpdatedBeforeKey is the standard key for referring to a types.QueryFilter's UpdatedAfter field.
	FilterUpdatedBeforeKey = "query_filter.updated_before"
	// FilterSortByKey is the standard key for referring to a types.QueryFilter's SortBy field.
	FilterSortByKey = "query_filter.sort_by"
	// FilterPageKey is the standard key for referring to a types.QueryFilter's page.
	FilterPageKey = "query_filter.page"
	// FilterLimitKey is the standard key for referring to a types.QueryFilter's limit.
	FilterLimitKey = "query_filter.limit"
	// FilterIsNilKey is the standard key for referring to a types.QueryFilter's null status.
	FilterIsNilKey = "query_filter.is_nil"
	// APIClientClientIDKey is the standard key for referring to an API client's database ID.
	APIClientClientIDKey = "api_client.client_id"
	// APIClientDatabaseIDKey is the standard key for referring to an API client's database ID.
	APIClientDatabaseIDKey = "api_client.id"
	// WebhookIDKey is the standard key for referring to a webhook's ID.
	WebhookIDKey = "webhook.id"
	// URLKey is the standard key for referring to a url.
	URLKey = "url"
	// RequestHeadersKey is the standard key for referring to an http.Request's Headers.
	RequestHeadersKey = "request.headers"
	// RequestMethodKey is the standard key for referring to an http.Request's Method.
	RequestMethodKey = "request.method"
	// RequestURIKey is the standard key for referring to an http.Request's URI.
	RequestURIKey = "request.uri"
	// ResponseStatusKey is the standard key for referring to an http.Request's URI.
	ResponseStatusKey = "response.status"
	// ResponseHeadersKey is the standard key for referring to an http.Response's Headers.
	ResponseHeadersKey = "response.headers"
	// ReasonKey is the standard key for referring to a reason.
	ReasonKey = "reason"
	// DatabaseQueryKey is the standard key for referring to a database query.
	DatabaseQueryKey = "database_query"
	// URLQueryKey is the standard key for referring to a url query.
	URLQueryKey = "url.query"
	// ConnectionDetailsKey is the standard key for referring to a database's URI.
	ConnectionDetailsKey = "database.connection_details"
	// SearchQueryKey is the standard key for referring to a search query parameter value.
	SearchQueryKey = "search_query"
	// UserAgentOSKey is the standard key for referring to a search query parameter value.
	UserAgentOSKey = "os"
	// UserAgentBotKey is the standard key for referring to a search query parameter value.
	UserAgentBotKey = "is_bot"
	// UserAgentMobileKey is the standard key for referring to a search query parameter value.
	UserAgentMobileKey = "is_mobile"
	// RollbackErrorKey is the standard key for referring to an error rolling back a transaction.
	RollbackErrorKey = "ROLLBACK_ERROR"
	// QueryErrorKey is the standard key for referring to an error building a query.
	QueryErrorKey = "QUERY_ERROR"
	// RowIDErrorKey is the standard key for referring to an error fetching a row ID.
	RowIDErrorKey = "ROW_ID_ERROR"
	// ValidationErrorKey is the standard key for referring to a struct validation error.
	ValidationErrorKey = "validation_error"

	// ValidInstrumentIDKey is the standard key for referring to a valid instrument ID.
	ValidInstrumentIDKey = "valid_instrument_id"

	// ValidIngredientIDKey is the standard key for referring to a valid ingredient ID.
	ValidIngredientIDKey = "valid_ingredient_id"

	// ValidPreparationIDKey is the standard key for referring to a valid preparation ID.
	ValidPreparationIDKey = "valid_preparation_id"

	// ValidIngredientPreparationIDKey is the standard key for referring to a valid ingredient preparation ID.
	ValidIngredientPreparationIDKey = "valid_ingredient_preparation_id"

	// RecipeIDKey is the standard key for referring to a recipe ID.
	RecipeIDKey = "recipe_id"

	// RecipeStepIDKey is the standard key for referring to a recipe step ID.
	RecipeStepIDKey = "recipe_step_id"

	// RecipeStepInstrumentIDKey is the standard key for referring to a recipe step instrument ID.
	RecipeStepInstrumentIDKey = "recipe_step_instrument_id"

	// RecipeStepIngredientIDKey is the standard key for referring to a recipe step ingredient ID.
	RecipeStepIngredientIDKey = "recipe_step_ingredient_id"

	// RecipeStepProductIDKey is the standard key for referring to a recipe step product ID.
	RecipeStepProductIDKey = "recipe_step_product_id"

	// MealPlanIDKey is the standard key for referring to a meal plan ID.
	MealPlanIDKey = "meal_plan_id"

	// MealPlanOptionIDKey is the standard key for referring to a meal plan option ID.
	MealPlanOptionIDKey = "meal_plan_option_id"

	// MealPlanOptionVoteIDKey is the standard key for referring to a meal plan option vote ID.
	MealPlanOptionVoteIDKey = "meal_plan_option_vote_id"
)
