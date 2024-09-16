package main

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
)

type routeDetails struct {
	ResponseType  any
	InputType     any
	ListRoute     bool
	InputTypeName string
}

var routeInfoMap = map[string]routeDetails{
	"GET /_meta_/live":                       {},
	"GET /_meta_/ready":                      {},
	"POST /api/v1/admin/cycle_cookie_secret": {},
	"POST /api/v1/admin/users/status": {
		ResponseType: &types.UserStatusResponse{},
	},
	"GET /api/v1/audit_log_entries/for_household": {
		ResponseType: &types.AuditLogEntry{},
	},
	"GET /api/v1/audit_log_entries/for_user": {
		ResponseType: &types.AuditLogEntry{},
	},
	"GET /api/v1/audit_log_entries/{auditLogEntryID}": {
		ResponseType: &types.AuditLogEntry{},
	},
	"GET /api/v1/household_invitations/received": {
		ResponseType: &types.HouseholdInvitation{},
	},
	"GET /api/v1/household_invitations/sent": {
		ResponseType: &types.HouseholdInvitation{},
	},
	"GET /api/v1/household_invitations/{householdInvitationID}/": {
		ResponseType: &types.HouseholdInvitation{},
	},
	"PUT /api/v1/household_invitations/{householdInvitationID}/accept": {
		ResponseType: &types.HouseholdInvitation{},
		InputType:    &types.HouseholdInvitationUpdateRequestInput{},
	},
	"PUT /api/v1/household_invitations/{householdInvitationID}/cancel": {
		ResponseType: &types.HouseholdInvitation{},
		InputType:    &types.HouseholdInvitationUpdateRequestInput{},
	},
	"PUT /api/v1/household_invitations/{householdInvitationID}/reject": {
		ResponseType: &types.HouseholdInvitation{},
		InputType:    &types.HouseholdInvitationUpdateRequestInput{},
	},
	"GET /api/v1/households/": {
		ResponseType: &types.Household{},
		ListRoute:    true,
	},
	"POST /api/v1/households/": {
		ResponseType: &types.Household{},
		InputType:    &types.HouseholdCreationRequestInput{},
	},
	"GET /api/v1/households/current": {
		ResponseType: &types.Household{},
	},
	"POST /api/v1/households/instruments/": {
		ResponseType: &types.HouseholdInstrumentOwnership{},
		InputType:    &types.HouseholdInstrumentOwnershipCreationRequestInput{},
	},
	"GET /api/v1/households/instruments/": {
		ResponseType: &types.HouseholdInstrumentOwnership{},
		ListRoute:    true,
	},
	"GET /api/v1/households/instruments/{householdInstrumentOwnershipID}/": {
		ResponseType: &types.HouseholdInstrumentOwnership{},
	},
	"PUT /api/v1/households/instruments/{householdInstrumentOwnershipID}/": {
		ResponseType: &types.HouseholdInstrumentOwnership{},
		InputType:    &types.HouseholdInstrumentOwnershipUpdateRequestInput{},
	},
	"DELETE /api/v1/households/instruments/{householdInstrumentOwnershipID}/": {
		ResponseType: &types.HouseholdInstrumentOwnership{},
	},
	"PUT /api/v1/households/{householdID}/": {
		ResponseType: &types.Household{},
		InputType:    &types.HouseholdCreationRequestInput{},
	},
	"DELETE /api/v1/households/{householdID}/": {
		ResponseType: &types.Household{},
	},
	"GET /api/v1/households/{householdID}/": {
		ResponseType: &types.Household{},
	},
	"POST /api/v1/households/{householdID}/default": {
		ResponseType: &types.Household{},
	},
	"GET /api/v1/households/{householdID}/invitations/{householdInvitationID}/": {
		ResponseType: &types.HouseholdInvitation{},
	},
	// these are duplicate routes lol
	"POST /api/v1/households/{householdID}/invitations/": {
		ResponseType: &types.HouseholdInvitation{},
		InputType:    &types.HouseholdInvitationCreationRequestInput{},
	},
	"POST /api/v1/households/{householdID}/invite": {
		ResponseType: &types.HouseholdInvitation{},
		InputType:    &types.HouseholdInvitationCreationRequestInput{},
	},
	"DELETE /api/v1/households/{householdID}/members/{userID}": {
		ResponseType: &types.HouseholdUserMembership{},
	},
	"PATCH /api/v1/households/{householdID}/members/{userID}/permissions": {
		ResponseType: &types.UserPermissionsResponse{},
	},
	"POST /api/v1/households/{householdID}/transfer": {
		ResponseType: &types.Household{},
		InputType:    &types.HouseholdOwnershipTransferInput{},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/": {
		ResponseType: &types.MealPlanOptionVote{},
		ListRoute:    true,
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/": {
		ResponseType: &types.MealPlanOptionVote{},
		InputType:    &types.MealPlanOptionVoteUpdateRequestInput{},
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/": {
		ResponseType: &types.MealPlanOptionVote{},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/": {
		ResponseType: &types.MealPlanOptionVote{},
	},
	"POST /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/": {
		ResponseType: &types.MealPlanOption{},
		InputType:    &types.MealPlanOptionCreationRequestInput{},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/": {
		ResponseType: &types.MealPlanOption{},
		ListRoute:    true,
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/": {
		ResponseType: &types.MealPlanOption{},
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/": {
		ResponseType: &types.MealPlanOption{},
		InputType:    &types.MealPlanOptionUpdateRequestInput{},
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/": {
		ResponseType: &types.MealPlanOption{},
	},
	"POST /api/v1/meal_plans/{mealPlanID}/events/": {
		ResponseType: &types.MealPlanEvent{},
		InputType:    &types.MealPlanEventCreationRequestInput{},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/": {
		ResponseType: &types.MealPlanEvent{},
		ListRoute:    true,
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/": {
		ResponseType: &types.MealPlanEvent{},
		InputType:    &types.MealPlanEventUpdateRequestInput{},
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/": {
		ResponseType: &types.MealPlanEvent{},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/": {
		ResponseType: &types.MealPlanEvent{},
	},
	"POST /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/vote": {
		ResponseType: &types.MealPlanOptionVote{},
		InputType:    &types.MealPlanOptionVoteCreationRequestInput{},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/grocery_list_items/": {
		ResponseType: &types.MealPlanGroceryListItem{},
		ListRoute:    true,
	},
	"POST /api/v1/meal_plans/{mealPlanID}/grocery_list_items/": {
		ResponseType: &types.MealPlanGroceryListItem{},
		InputType:    &types.MealPlanGroceryListItemCreationRequestInput{},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/": {
		ResponseType: &types.MealPlanGroceryListItem{},
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/": {
		ResponseType: &types.MealPlanGroceryListItem{},
		InputType:    &types.MealPlanGroceryListItemUpdateRequestInput{},
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/": {
		ResponseType: &types.MealPlanGroceryListItem{},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/tasks/": {
		ResponseType: &types.MealPlanTask{},
		ListRoute:    true,
	},
	"POST /api/v1/meal_plans/{mealPlanID}/tasks/": {
		ResponseType: &types.MealPlanTask{},
		InputType:    &types.MealPlanTaskCreationRequestInput{},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}/": {
		ResponseType: &types.MealPlanTask{},
	},
	"PATCH /api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}/": {
		ResponseType: &types.MealPlanTask{},
	},
	"POST /api/v1/meal_plans/": {
		ResponseType: &types.MealPlan{},
		InputType:    &types.MealPlanCreationRequestInput{},
	},
	"GET /api/v1/meal_plans/": {
		ResponseType: &types.MealPlan{},
		ListRoute:    true,
	},
	"GET /api/v1/meal_plans/{mealPlanID}/": {
		ResponseType: &types.MealPlan{},
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/": {
		ResponseType: &types.MealPlan{},
		InputType:    &types.MealPlanUpdateRequestInput{},
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/": {
		ResponseType: &types.MealPlan{},
	},
	"POST /api/v1/meal_plans/{mealPlanID}/finalize": {
		ResponseType: &types.MealPlan{},
		// No input type for this route
	},
	"POST /api/v1/meals/": {
		ResponseType: &types.Meal{},
		InputType:    &types.MealCreationRequestInput{},
	},
	"GET /api/v1/meals/": {
		ResponseType: &types.Meal{},
		ListRoute:    true,
	},
	"GET /api/v1/meals/search": {
		ResponseType: &types.Meal{},
	},
	"DELETE /api/v1/meals/{mealID}/": {
		ResponseType: &types.Meal{},
	},
	"GET /api/v1/meals/{mealID}/": {
		ResponseType: &types.Meal{},
	},
	"GET /api/v1/oauth2_clients/": {
		ResponseType: &types.OAuth2Client{},
		ListRoute:    true,
	},
	"POST /api/v1/oauth2_clients/": {
		ResponseType: &types.OAuth2Client{},
		InputType:    &types.OAuth2ClientCreationRequestInput{},
	},
	"GET /api/v1/oauth2_clients/{oauth2ClientID}/": {
		ResponseType: &types.OAuth2Client{},
	},
	"DELETE /api/v1/oauth2_clients/{oauth2ClientID}/": {
		ResponseType: &types.OAuth2Client{},
	},
	"POST /api/v1/recipes/{recipeID}/prep_tasks/": {
		ResponseType: &types.RecipePrepTask{},
		InputType:    &types.RecipePrepTaskCreationRequestInput{},
	},
	"GET /api/v1/recipes/{recipeID}/prep_tasks/": {
		ResponseType: &types.RecipePrepTask{},
		ListRoute:    true,
	},
	"DELETE /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/": {
		ResponseType: &types.RecipePrepTask{},
	},
	"GET /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/": {
		ResponseType: &types.RecipePrepTask{},
	},
	"PUT /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/": {
		ResponseType: &types.RecipePrepTask{},
		InputType:    &types.RecipePrepTaskUpdateRequestInput{},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/": {
		ResponseType: &types.RecipeStepCompletionCondition{},
		ListRoute:    true,
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/": {
		ResponseType: &types.RecipeStepCompletionCondition{},
		InputType:    &types.RecipeStepCompletionConditionCreationRequestInput{},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/": {
		ResponseType: &types.RecipeStepCompletionCondition{},
		InputType:    &types.RecipeStepCompletionConditionUpdateRequestInput{},
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/": {
		ResponseType: &types.RecipeStepCompletionCondition{},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/": {
		ResponseType: &types.RecipeStepCompletionCondition{},
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/": {
		ResponseType: &types.RecipeStepIngredient{},
		InputType:    &types.RecipeStepIngredientCreationRequestInput{},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/": {
		ResponseType: &types.RecipeStepIngredient{},
		ListRoute:    true,
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/": {
		ResponseType: &types.RecipeStepIngredient{},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/": {
		ResponseType: &types.RecipeStepIngredient{},
		InputType:    &types.RecipeStepIngredientUpdateRequestInput{},
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/": {
		ResponseType: &types.RecipeStepIngredient{},
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/": {
		ResponseType: &types.RecipeStepInstrument{},
		InputType:    &types.RecipeStepInstrumentCreationRequestInput{},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/": {
		ResponseType: &types.RecipeStepInstrument{},
		ListRoute:    true,
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/": {
		ResponseType: &types.RecipeStepInstrument{},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/": {
		ResponseType: &types.RecipeStepInstrument{},
		InputType:    &types.RecipeStepInstrumentUpdateRequestInput{},
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/": {
		ResponseType: &types.RecipeStepInstrument{},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/": {
		ResponseType: &types.RecipeStepProduct{},
		ListRoute:    true,
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/": {
		ResponseType: &types.RecipeStepProduct{},
		InputType:    &types.RecipeStepProductCreationRequestInput{},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/": {
		ResponseType: &types.RecipeStepProduct{},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/": {
		ResponseType: &types.RecipeStepProduct{},
		InputType:    &types.RecipeStepProductUpdateRequestInput{},
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/": {
		ResponseType: &types.RecipeStepProduct{},
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/": {
		ResponseType: &types.RecipeStepVessel{},
		InputType:    &types.RecipeStepVesselCreationRequestInput{},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/": {
		ResponseType: &types.RecipeStepVessel{},
		ListRoute:    true,
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/": {
		ResponseType: &types.RecipeStepVessel{},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/": {
		ResponseType: &types.RecipeStepVessel{},
		InputType:    &types.RecipeStepVesselUpdateRequestInput{},
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/": {
		ResponseType: &types.RecipeStepVessel{},
	},
	"POST /api/v1/recipes/{recipeID}/steps/": {
		ResponseType: &types.RecipeStep{},
		InputType:    &types.RecipeStepCreationRequestInput{},
	},
	"GET /api/v1/recipes/{recipeID}/steps/": {
		ResponseType: &types.RecipeStep{},
		ListRoute:    true,
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/": {
		ResponseType: &types.RecipeStep{},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/": {
		ResponseType: &types.RecipeStep{},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/": {
		ResponseType: &types.RecipeStep{},
		InputType:    &types.RecipeStepUpdateRequestInput{},
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/images": {},
	"POST /api/v1/recipes/": {
		ResponseType: &types.Recipe{},
		InputType:    &types.RecipeCreationRequestInput{},
	},
	"GET /api/v1/recipes/": {
		ResponseType: &types.Recipe{},
		ListRoute:    true,
	},
	"GET /api/v1/recipes/search": {
		ResponseType: &types.Recipe{},
	},
	"DELETE /api/v1/recipes/{recipeID}/": {
		ResponseType: &types.Recipe{},
	},
	"GET /api/v1/recipes/{recipeID}/": {
		ResponseType: &types.Recipe{},
	},
	"PUT /api/v1/recipes/{recipeID}/": {
		ResponseType: &types.Recipe{},
		InputType:    &types.RecipeUpdateRequestInput{},
	},
	"POST /api/v1/recipes/{recipeID}/clone": {
		ResponseType: &types.Recipe{},
		// No input type for this route
	},
	"GET /api/v1/recipes/{recipeID}/dag": {
		ResponseType: &types.APIError{},
	},
	"POST /api/v1/recipes/{recipeID}/images": {},
	"GET /api/v1/recipes/{recipeID}/mermaid": {},
	"GET /api/v1/recipes/{recipeID}/prep_steps": {
		ResponseType: &types.RecipePrepTaskStep{},
	},
	"POST /api/v1/recipes/{recipeID}/ratings/": {
		ResponseType: &types.RecipeRating{},
		InputType:    &types.RecipeRatingCreationRequestInput{},
	},
	"GET /api/v1/recipes/{recipeID}/ratings/": {
		ResponseType: &types.RecipeRating{},
		ListRoute:    true,
	},
	"PUT /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/": {
		ResponseType: &types.RecipeRating{},
		InputType:    &types.RecipeRatingUpdateRequestInput{},
	},
	"DELETE /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/": {
		ResponseType: &types.RecipeRating{},
	},
	"GET /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/": {
		ResponseType: &types.RecipeRating{},
	},
	"POST /api/v1/settings/": {
		ResponseType: &types.ServiceSetting{},
		InputType:    &types.ServiceSettingCreationRequestInput{},
	},
	"GET /api/v1/settings/": {
		ResponseType: &types.ServiceSetting{},
		ListRoute:    true,
	},
	"POST /api/v1/settings/configurations/": {
		ResponseType: &types.ServiceSettingConfiguration{},
		InputType:    &types.ServiceSettingConfigurationCreationRequestInput{},
	},
	"GET /api/v1/settings/configurations/household": {
		ResponseType: &types.ServiceSettingConfiguration{},
	},
	"GET /api/v1/settings/configurations/user": {
		ResponseType: &types.ServiceSettingConfiguration{},
	},
	"GET /api/v1/settings/configurations/user/{serviceSettingConfigurationName}": {
		ResponseType: &types.ServiceSettingConfiguration{},
	},
	"DELETE /api/v1/settings/configurations/{serviceSettingConfigurationID}": {
		ResponseType: &types.ServiceSettingConfiguration{},
	},
	"PUT /api/v1/settings/configurations/{serviceSettingConfigurationID}": {
		ResponseType: &types.ServiceSettingConfiguration{},
		InputType:    &types.ServiceSettingConfigurationUpdateRequestInput{},
	},
	"GET /api/v1/settings/search": {
		ResponseType: &types.ServiceSetting{},
	},
	"GET /api/v1/settings/{serviceSettingID}/": {
		ResponseType: &types.ServiceSetting{},
	},
	"DELETE /api/v1/settings/{serviceSettingID}/": {
		ResponseType: &types.ServiceSetting{},
	},
	"POST /api/v1/user_ingredient_preferences/": {
		ResponseType: &types.UserIngredientPreference{},
		InputType:    &types.UserIngredientPreferenceCreationRequestInput{},
	},
	"GET /api/v1/user_ingredient_preferences/": {
		ResponseType: &types.UserIngredientPreference{},
		ListRoute:    true,
	},
	"PUT /api/v1/user_ingredient_preferences/{userIngredientPreferenceID}/": {
		ResponseType: &types.UserIngredientPreference{},
		InputType:    &types.UserIngredientPreferenceUpdateRequestInput{},
	},
	"DELETE /api/v1/user_ingredient_preferences/{userIngredientPreferenceID}/": {
		ResponseType: &types.UserIngredientPreference{},
	},
	"POST /api/v1/user_notifications/": {
		ResponseType: &types.UserNotification{},
		InputType:    &types.UserNotificationCreationRequestInput{},
	},
	"GET /api/v1/user_notifications/": {
		ResponseType: &types.UserNotification{},
		ListRoute:    true,
	},
	"GET /api/v1/user_notifications/{userNotificationID}/": {
		ResponseType: &types.UserNotification{},
	},
	"PATCH /api/v1/user_notifications/{userNotificationID}/": {
		ResponseType: &types.UserNotification{},
	},
	"GET /api/v1/users/": {
		ResponseType: &types.User{},
		ListRoute:    true,
	},
	"POST /api/v1/users/avatar/upload": {},
	"PUT /api/v1/users/details": {
		ResponseType: &types.User{},
		InputType:    &types.UserDetailsUpdateRequestInput{},
	},
	"PUT /api/v1/users/email_address": {
		ResponseType: &types.User{},
		InputType:    &types.UserEmailAddressUpdateInput{},
	},
	"POST /api/v1/users/email_address_verification": {
		ResponseType: &types.User{},
		InputType:    &types.EmailAddressVerificationRequestInput{},
	},
	"POST /api/v1/users/household/select": {
		ResponseType: &types.Household{},
		InputType:    &types.ChangeActiveHouseholdInput{},
	},
	"PUT /api/v1/users/password/new": {
		// No output type for this route
	},
	"POST /api/v1/users/permissions/check": {
		ResponseType: &types.UserPermissionsResponse{},
		InputType:    &types.UserPermissionsRequestInput{},
	},
	"GET /api/v1/users/search": {
		ResponseType: &types.User{},
		ListRoute:    true,
	},
	"GET /api/v1/users/self": {
		ResponseType: &types.User{},
	},
	"POST /api/v1/users/totp_secret/new": {
		ResponseType: &types.APIError{},
		InputType:    &types.TOTPSecretRefreshInput{},
	},
	"PUT /api/v1/users/username": {
		ResponseType: &types.User{},
		InputType:    &types.UsernameUpdateInput{},
	},
	"GET /api/v1/users/{userID}/": {
		ResponseType: &types.User{},
	},
	"DELETE /api/v1/users/{userID}/": {
		ResponseType: &types.User{},
	},
	"POST /api/v1/valid_ingredient_groups/": {
		ResponseType: &types.ValidIngredientGroup{},
		InputType:    &types.ValidIngredientGroupCreationRequestInput{},
	},
	"GET /api/v1/valid_ingredient_groups/": {
		ResponseType: &types.ValidIngredientGroup{},
		ListRoute:    true,
	},
	"GET /api/v1/valid_ingredient_groups/search": {
		ResponseType: &types.ValidIngredientGroup{},
	},
	"DELETE /api/v1/valid_ingredient_groups/{validIngredientGroupID}/": {
		ResponseType: &types.ValidIngredientGroup{},
	},
	"GET /api/v1/valid_ingredient_groups/{validIngredientGroupID}/": {
		ResponseType: &types.ValidIngredientGroup{},
	},
	"PUT /api/v1/valid_ingredient_groups/{validIngredientGroupID}/": {
		ResponseType: &types.ValidIngredientGroup{},
		InputType:    &types.ValidIngredientGroupUpdateRequestInput{},
	},
	"POST /api/v1/valid_ingredient_measurement_units/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
		InputType:    &types.ValidIngredientMeasurementUnitCreationRequestInput{},
	},
	"GET /api/v1/valid_ingredient_measurement_units/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
		ListRoute:    true,
	},
	"GET /api/v1/valid_ingredient_measurement_units/by_ingredient/{validIngredientID}/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
	},
	"GET /api/v1/valid_ingredient_measurement_units/by_measurement_unit/{validMeasurementUnitID}/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
	},
	"GET /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
	},
	"PUT /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
		InputType:    &types.ValidIngredientMeasurementUnitUpdateRequestInput{},
	},
	"DELETE /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
	},
	"GET /api/v1/valid_ingredient_preparations/": {
		ResponseType: &types.ValidIngredientPreparation{},
		ListRoute:    true,
	},
	"POST /api/v1/valid_ingredient_preparations/": {
		ResponseType: &types.ValidIngredientPreparation{},
		InputType:    &types.ValidIngredientPreparationCreationRequestInput{},
	},
	"GET /api/v1/valid_ingredient_preparations/by_ingredient/{validIngredientID}/": {
		ResponseType: &types.ValidIngredientPreparation{},
	},
	"GET /api/v1/valid_ingredient_preparations/by_preparation/{validPreparationID}/": {
		ResponseType: &types.ValidIngredientPreparation{},
	},
	"GET /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/": {
		ResponseType: &types.ValidIngredientPreparation{},
	},
	"PUT /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/": {
		ResponseType: &types.ValidIngredientPreparation{},
		InputType:    &types.ValidIngredientPreparationUpdateRequestInput{},
	},
	"DELETE /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/": {
		ResponseType: &types.ValidIngredientPreparation{},
	},
	"POST /api/v1/valid_ingredient_state_ingredients/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
		InputType:    &types.ValidIngredientStateIngredientCreationRequestInput{},
	},
	"GET /api/v1/valid_ingredient_state_ingredients/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
		ListRoute:    true,
	},
	"GET /api/v1/valid_ingredient_state_ingredients/by_ingredient/{validIngredientID}/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
	},
	"GET /api/v1/valid_ingredient_state_ingredients/by_ingredient_state/{validIngredientStateID}/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
	},
	"GET /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
	},
	"PUT /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
		InputType:    &types.ValidIngredientStateIngredientUpdateRequestInput{},
	},
	"DELETE /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
	},
	"POST /api/v1/valid_ingredient_states/": {
		ResponseType: &types.ValidIngredientState{},
		InputType:    &types.ValidIngredientStateCreationRequestInput{},
	},
	"GET /api/v1/valid_ingredient_states/": {
		ResponseType: &types.ValidIngredientState{},
		ListRoute:    true,
	},
	"GET /api/v1/valid_ingredient_states/search": {
		ResponseType: &types.ValidIngredientState{},
	},
	"PUT /api/v1/valid_ingredient_states/{validIngredientStateID}/": {
		ResponseType: &types.ValidIngredientState{},
		InputType:    &types.ValidIngredientStateUpdateRequestInput{},
	},
	"DELETE /api/v1/valid_ingredient_states/{validIngredientStateID}/": {
		ResponseType: &types.ValidIngredientState{},
	},
	"GET /api/v1/valid_ingredient_states/{validIngredientStateID}/": {
		ResponseType: &types.ValidIngredientState{},
	},
	"POST /api/v1/valid_ingredients/": {
		ResponseType: &types.ValidIngredient{},
		InputType:    &types.ValidIngredientCreationRequestInput{},
	},
	"GET /api/v1/valid_ingredients/": {
		ResponseType: &types.ValidIngredient{},
		ListRoute:    true,
	},
	"GET /api/v1/valid_ingredients/by_preparation/{validPreparationID}/": {
		ResponseType: &types.ValidIngredient{},
	},
	"GET /api/v1/valid_ingredients/random": {
		ResponseType: &types.ValidIngredient{},
	},
	"GET /api/v1/valid_ingredients/search": {
		ResponseType: &types.ValidIngredient{},
	},
	"PUT /api/v1/valid_ingredients/{validIngredientID}/": {
		ResponseType: &types.ValidIngredient{},
		InputType:    &types.ValidIngredientUpdateRequestInput{},
	},
	"DELETE /api/v1/valid_ingredients/{validIngredientID}/": {
		ResponseType: &types.ValidIngredient{},
	},
	"GET /api/v1/valid_ingredients/{validIngredientID}/": {
		ResponseType: &types.ValidIngredient{},
	},
	"GET /api/v1/valid_instruments/": {
		ResponseType: &types.ValidInstrument{},
		ListRoute:    true,
	},
	"POST /api/v1/valid_instruments/": {
		ResponseType: &types.ValidInstrument{},
		InputType:    &types.ValidInstrumentCreationRequestInput{},
	},
	"GET /api/v1/valid_instruments/random": {
		ResponseType: &types.ValidInstrument{},
	},
	"GET /api/v1/valid_instruments/search": {
		ResponseType: &types.ValidInstrument{},
	},
	"DELETE /api/v1/valid_instruments/{validInstrumentID}/": {
		ResponseType: &types.ValidInstrument{},
	},
	"GET /api/v1/valid_instruments/{validInstrumentID}/": {
		ResponseType: &types.ValidInstrument{},
	},
	"PUT /api/v1/valid_instruments/{validInstrumentID}/": {
		ResponseType: &types.ValidInstrument{},
		InputType:    &types.ValidInstrumentUpdateRequestInput{},
	},
	"POST /api/v1/valid_measurement_conversions/": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
		InputType:    &types.ValidMeasurementUnitConversionCreationRequestInput{},
	},
	"GET /api/v1/valid_measurement_conversions/from_unit/{validMeasurementUnitID}": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
	},
	"GET /api/v1/valid_measurement_conversions/to_unit/{validMeasurementUnitID}": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
	},
	"PUT /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
		InputType:    &types.ValidMeasurementUnitConversionUpdateRequestInput{},
	},
	"DELETE /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
	},
	"GET /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
	},
	"POST /api/v1/valid_measurement_units/": {
		ResponseType: &types.ValidMeasurementUnit{},
		InputType:    &types.ValidMeasurementUnitCreationRequestInput{},
	},
	"GET /api/v1/valid_measurement_units/": {
		ResponseType: &types.ValidMeasurementUnit{},
		ListRoute:    true,
	},
	"GET /api/v1/valid_measurement_units/by_ingredient/{validIngredientID}": {
		ResponseType: &types.ValidMeasurementUnit{},
	},
	"GET /api/v1/valid_measurement_units/search": {
		ResponseType: &types.ValidMeasurementUnit{},
	},
	"GET /api/v1/valid_measurement_units/{validMeasurementUnitID}/": {
		ResponseType: &types.ValidMeasurementUnit{},
	},
	"PUT /api/v1/valid_measurement_units/{validMeasurementUnitID}/": {
		ResponseType: &types.ValidMeasurementUnit{},
		InputType:    &types.ValidMeasurementUnitUpdateRequestInput{},
	},
	"DELETE /api/v1/valid_measurement_units/{validMeasurementUnitID}/": {
		ResponseType: &types.ValidMeasurementUnit{},
	},
	"GET /api/v1/valid_preparation_instruments/": {
		ResponseType: &types.ValidPreparationInstrument{},
		ListRoute:    true,
	},
	"POST /api/v1/valid_preparation_instruments/": {
		ResponseType: &types.ValidPreparationInstrument{},
		InputType:    &types.ValidPreparationInstrumentCreationRequestInput{},
	},
	"GET /api/v1/valid_preparation_instruments/by_instrument/{validInstrumentID}/": {
		ResponseType: &types.ValidPreparationInstrument{},
	},
	"GET /api/v1/valid_preparation_instruments/by_preparation/{validPreparationID}/": {
		ResponseType: &types.ValidPreparationInstrument{},
	},
	"DELETE /api/v1/valid_preparation_instruments/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationInstrument{},
	},
	"GET /api/v1/valid_preparation_instruments/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationInstrument{},
	},
	"PUT /api/v1/valid_preparation_instruments/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationInstrument{},
		InputType:    &types.ValidPreparationInstrumentUpdateRequestInput{},
	},
	"POST /api/v1/valid_preparation_vessels/": {
		ResponseType: &types.ValidPreparationVessel{},
		InputType:    &types.ValidPreparationVesselCreationRequestInput{},
	},
	"GET /api/v1/valid_preparation_vessels/": {
		ResponseType: &types.ValidPreparationVessel{},
		ListRoute:    true,
	},
	"GET /api/v1/valid_preparation_vessels/by_preparation/{validPreparationID}/": {
		ResponseType: &types.ValidPreparationVessel{},
	},
	"GET /api/v1/valid_preparation_vessels/by_vessel/{ValidVesselID}/": {
		ResponseType: &types.ValidPreparationVessel{},
	},
	"PUT /api/v1/valid_preparation_vessels/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationVessel{},
		InputType:    &types.ValidPreparationVesselUpdateRequestInput{},
	},
	"DELETE /api/v1/valid_preparation_vessels/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationVessel{},
	},
	"GET /api/v1/valid_preparation_vessels/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationVessel{},
	},
	"GET /api/v1/valid_preparations/": {
		ResponseType: &types.ValidPreparation{},
		ListRoute:    true,
	},
	"POST /api/v1/valid_preparations/": {
		ResponseType: &types.ValidPreparation{},
		InputType:    &types.ValidPreparationCreationRequestInput{},
	},
	"GET /api/v1/valid_preparations/random": {
		ResponseType: &types.ValidPreparation{},
	},
	"GET /api/v1/valid_preparations/search": {
		ResponseType: &types.ValidPreparation{},
	},
	"PUT /api/v1/valid_preparations/{validPreparationID}/": {
		ResponseType: &types.ValidPreparation{},
		InputType:    &types.ValidPreparationUpdateRequestInput{},
	},
	"DELETE /api/v1/valid_preparations/{validPreparationID}/": {
		ResponseType: &types.ValidPreparation{},
	},
	"GET /api/v1/valid_preparations/{validPreparationID}/": {
		ResponseType: &types.ValidPreparation{},
	},
	"POST /api/v1/valid_vessels/": {
		ResponseType: &types.ValidVessel{},
		InputType:    &types.ValidVesselCreationRequestInput{},
	},
	"GET /api/v1/valid_vessels/": {
		ResponseType: &types.ValidVessel{},
		ListRoute:    true,
	},
	"GET /api/v1/valid_vessels/random": {
		ResponseType: &types.ValidVessel{},
	},
	"GET /api/v1/valid_vessels/search": {
		ResponseType: &types.ValidVessel{},
	},
	"GET /api/v1/valid_vessels/{validVesselID}/": {
		ResponseType: &types.ValidVessel{},
	},
	"PUT /api/v1/valid_vessels/{validVesselID}/": {
		ResponseType: &types.ValidVessel{},
		InputType:    &types.ValidVesselUpdateRequestInput{},
	},
	"DELETE /api/v1/valid_vessels/{validVesselID}/": {
		ResponseType: &types.ValidVessel{},
	},
	"GET /api/v1/webhooks/": {
		ResponseType: &types.Webhook{},
		ListRoute:    true,
	},
	"POST /api/v1/webhooks/": {
		ResponseType: &types.Webhook{},
		InputType:    &types.WebhookCreationRequestInput{},
	},
	"GET /api/v1/webhooks/{webhookID}/": {
		ResponseType: &types.Webhook{},
	},
	"DELETE /api/v1/webhooks/{webhookID}/": {
		ResponseType: &types.Webhook{},
	},
	"POST /api/v1/webhooks/{webhookID}/trigger_events": {
		ResponseType: &types.WebhookTriggerEvent{},
		InputType:    &types.WebhookTriggerEventCreationRequestInput{},
	},
	"DELETE /api/v1/webhooks/{webhookID}/trigger_events/{webhookTriggerEventID}/": {
		ResponseType: &types.WebhookTriggerEvent{},
	},
	"POST /api/v1/workers/finalize_meal_plans": {
		ResponseType: &types.FinalizeMealPlansRequest{},
		InputType:    &types.FinalizeMealPlansRequest{},
	},
	"POST /api/v1/workers/meal_plan_grocery_list_init": {
		// no input or output types for this route
	},
	"POST /api/v1/workers/meal_plan_tasks": {
		// no input or output types for this route
	},
	"GET /auth/status": {
		ResponseType: &types.UserStatusResponse{},
		// no input type for this route
	},
	"GET /auth/{auth_provider}": {
		// we don't really have control over this route
	},
	"GET /auth/{auth_provider}/callback": {
		// we don't really have control over this route
	},
	"GET /oauth2/authorize": {
		// we don't really have control over this route
	},
	"POST /oauth2/token": {
		// we don't really have control over this route
	},
	"POST /users/": {
		ResponseType: &types.User{},
		InputType:    &types.UserRegistrationInput{},
	},
	"POST /users/email_address/verify": {
		ResponseType: &types.User{},
		InputType:    &types.EmailAddressVerificationRequestInput{},
	},
	"POST /users/login": {
		ResponseType: &types.UserStatusResponse{},
		InputType:    &types.UserLoginInput{},
	},
	"POST /users/login/admin": {
		ResponseType: &types.UserStatusResponse{},
		InputType:    &types.UserLoginInput{},
	},
	"POST /users/logout": {
		ResponseType: &types.UserStatusResponse{},
	},
	"POST /users/password/reset": {
		ResponseType: &types.PasswordResetToken{},
		InputType:    &types.PasswordResetTokenCreationRequestInput{},
	},
	"POST /users/password/reset/redeem": {
		ResponseType: &types.User{},
		InputType:    &types.PasswordResetTokenRedemptionRequestInput{},
	},
	"POST /users/totp_secret/verify": {
		ResponseType: &types.User{},
		InputType:    &types.EmailAddressVerificationRequestInput{},
	},
	"POST /users/username/reminder": {
		ResponseType: &types.User{},
		InputType:    &types.UsernameReminderRequestInput{},
	},
}
