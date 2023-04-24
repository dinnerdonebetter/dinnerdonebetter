package keys

const (
	// RequesterIDKey is the standard key for referring to a requesting user's ID.
	RequesterIDKey = "request.made_by"
	// HouseholdIDKey is the standard key for referring to a household ID.
	HouseholdIDKey = "household.id"
	// HouseholdInvitationIDKey is the standard key for referring to a household ID.
	HouseholdInvitationIDKey = "household_invitation.id"
	// HouseholdInvitationTokenKey is the standard key for referring to a household invitation token.
	HouseholdInvitationTokenKey = "household_invitation.token"
	// ActiveHouseholdIDKey is the standard key for referring to an active household ID.
	ActiveHouseholdIDKey = "active_household.id"
	// UserIDKey is the standard key for referring to a user ID.
	UserIDKey = "user.id"
	// UserEmailAddressKey is the standard key for referring to a user's email address.
	UserEmailAddressKey = "user.email_address"
	// UserIsServiceAdminKey is the standard key for referring to a user's admin status.
	UserIsServiceAdminKey = "user.is_admin"
	// UsernameKey is the standard key for referring to a username.
	UsernameKey = "user.username"
	// ServiceRoleKey is the standard key for referring to a user's service role.
	ServiceRoleKey = "user.service_role"
	// NameKey is the standard key for referring to a name.
	NameKey = "name"
	// SpanIDKey is the standard key for referring to a span ID.
	SpanIDKey = "span.id"
	// TraceIDKey is the standard key for referring to a trace ID.
	TraceIDKey = "trace.id"
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
	// APIClientClientIDKey is the standard key for referring to an API client's client ID.
	APIClientClientIDKey = "api_client.client.id"
	// APIClientDatabaseIDKey is the standard key for referring to an API client's database ID.
	APIClientDatabaseIDKey = "api_client.id"
	// WebhookIDKey is the standard key for referring to a webhook's ID.
	WebhookIDKey = "webhook.id"
	// URLKey is the standard key for referring to a URL.
	URLKey = "url"
	// PasswordResetTokenIDKey is the standard key for referring to a password reset token's ID.
	PasswordResetTokenIDKey = "password_reset_token.id"
	// RequestHeadersKey is the standard key for referring to an http.Request's Headers.
	RequestHeadersKey = "request.headers"
	// RequestMethodKey is the standard key for referring to an http.Request's Method.
	RequestMethodKey = "request.method"
	// RequestURIKey is the standard key for referring to an http.Request's URI.
	RequestURIKey = "request.uri"
	// RequestURIPathKey is the standard key for referring to an http.Request's URI.
	RequestURIPathKey = "request.uri.path"
	// RequestURIQueryKey is the standard key for referring to an http.Request's URI.
	RequestURIQueryKey = "request.uri.query"
	// ResponseStatusKey is the standard key for referring to an http.Request's URI.
	ResponseStatusKey = "response.status"
	// ResponseHeadersKey is the standard key for referring to an http.Response's Headers.
	ResponseHeadersKey = "response.headers"
	// ReasonKey is the standard key for referring to a reason for a change.
	ReasonKey = "reason"
	// DatabaseQueryKey is the standard key for referring to a database query.
	DatabaseQueryKey = "database_query"
	// URLQueryKey is the standard key for referring to a URL query.
	URLQueryKey = "url.query"
	// SearchQueryKey is the standard key for referring to a search query parameter value.
	SearchQueryKey = "search_query"
	// UserAgentOSKey is the standard key for referring to a user agent's OS.
	UserAgentOSKey = "os"
	// UserAgentBotKey is the standard key for referring to a user agent's bot status.
	UserAgentBotKey = "is_bot"
	// UserAgentMobileKey is the standard key for referring to user agent's mobile status.
	UserAgentMobileKey = "is_mobile"
	// QueryErrorKey is the standard key for referring to an error building a query.
	QueryErrorKey = "QUERY_ERROR"
	// ValidationErrorKey is the standard key for referring to a struct validation error.
	ValidationErrorKey = "validation_error"

	// ValidInstrumentIDKey is the standard key for referring to a valid instrument's ID.
	ValidInstrumentIDKey = "valid_instrument.id"

	// ValidIngredientIDKey is the standard key for referring to a valid ingredient's ID.
	ValidIngredientIDKey = "valid_ingredient.id"

	// ValidPreparationIDKey is the standard key for referring to a valid preparation's ID.
	ValidPreparationIDKey = "valid_preparation.id"

	// ValidIngredientStateIDKey is the standard key for referring to a valid ingredient state's ID.
	ValidIngredientStateIDKey = "valid_ingredient_state.id"

	// ValidIngredientStateIngredientIDKey is the standard key for referring to a valid ingredient state ingredient's ID.
	ValidIngredientStateIngredientIDKey = "valid_ingredient_state_ingredient.id"

	// ValidIngredientPreparationIDKey is the standard key for referring to a valid ingredient preparation's ID.
	ValidIngredientPreparationIDKey = "valid_ingredient_preparation.id"

	// ValidPreparationInstrumentIDKey is the standard key for referring to a valid instrument preparation's ID.
	ValidPreparationInstrumentIDKey = "valid_preparation_instrument.id"

	// ValidIngredientMeasurementUnitIDKey is the standard key for referring to a valid instrument preparation's ID.
	ValidIngredientMeasurementUnitIDKey = "valid_ingredient_measurement_unit.id"

	// MealIDKey is the standard key for referring to a meal's ID.
	MealIDKey = "meal.id"

	// RecipeIDKey is the standard key for referring to a recipe's ID.
	RecipeIDKey = "recipe.id"

	// RecipeStepIDKey is the standard key for referring to a recipe step's ID.
	RecipeStepIDKey = "recipe_step.id"

	// RecipePrepTaskIDKey is the standard key for referring to a recipe prep task's ID.
	RecipePrepTaskIDKey = "recipe_prep_task.id"

	// RecipeStepInstrumentIDKey is the standard key for referring to a recipe step instrument's ID.
	RecipeStepInstrumentIDKey = "recipe_step_instrument.id"

	// RecipeStepVesselIDKey is the standard key for referring to a recipe step vessel's ID.
	RecipeStepVesselIDKey = "recipe_step_vessel.id"

	// RecipeStepIngredientIDKey is the standard key for referring to a recipe step ingredient's ID.
	RecipeStepIngredientIDKey = "recipe_step_ingredient.id"

	// RecipeStepCompletionConditionIDKey is the standard key for referring to a recipe step completion condition's ID.
	RecipeStepCompletionConditionIDKey = "recipe_step_completion_condition.id"

	// RecipeStepProductIDKey is the standard key for referring to a recipe step product's ID.
	RecipeStepProductIDKey = "recipe_step_product.id"

	// MealPlanIDKey is the standard key for referring to a meal plan's ID.
	MealPlanIDKey = "meal_plan.id"

	// MealPlanEventIDKey is the standard key for referring to a meal plan event's ID.
	MealPlanEventIDKey = "meal_plan_event.id"

	// MealPlanOptionIDKey is the standard key for referring to a meal plan option's ID.
	MealPlanOptionIDKey = "meal_plan_option.id"

	// MealPlanOptionVoteIDKey is the standard key for referring to a meal plan option vote's ID.
	MealPlanOptionVoteIDKey = "meal_plan_option_vote.id"

	// ValidMeasurementUnitIDKey is the standard key for referring to a valid measurement unit's ID.
	ValidMeasurementUnitIDKey = "valid_measurement_unit.id"

	// MealPlanTaskIDKey is the standard key for referring to a meal plan task's ID.
	MealPlanTaskIDKey = "meal_plan_task.id"

	// MealPlanGroceryListItemIDKey is the standard key for referring to a meal plan grocery list item's ID.
	MealPlanGroceryListItemIDKey = "meal_plan_grocery_list_item.id"

	// ValidMeasurementConversionIDKey is the standard key for referring to a valid measurement conversion's ID.
	ValidMeasurementConversionIDKey = "valid_measurement_conversion.id"

	// RecipeMediaIDKey is the standard key for referring to a recipe media's ID.
	RecipeMediaIDKey = "recipe_media.id"

	// ServiceSettingIDKey is the standard key for referring to a service setting's ID.
	ServiceSettingIDKey = "service_setting.id"
)
