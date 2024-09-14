package main

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
)

type routeDetails struct {
	ResponseType       any
	ResponseTypeSchema *OpenAPISchema
	InputType          any
	InputTypeSchema    *OpenAPISchema
	ListRoute          bool
	InputTypeName      string
}

var routeInfoMap = map[string]routeDetails{
	"GET /_meta_/live":                       {},
	"GET /_meta_/ready":                      {},
	"POST /api/v1/admin/cycle_cookie_secret": {},
	"POST /api/v1/admin/users/status": {
		ResponseType:       &types.UserStatusResponse{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserStatusResponse{}),
	},
	"GET /api/v1/audit_log_entries/for_household": {
		ResponseType:       &types.AuditLogEntry{},
		ResponseTypeSchema: SchemaFromInstance(&types.AuditLogEntry{}),
	},
	"GET /api/v1/audit_log_entries/for_user": {
		ResponseType:       &types.AuditLogEntry{},
		ResponseTypeSchema: SchemaFromInstance(&types.AuditLogEntry{}),
	},
	"GET /api/v1/audit_log_entries/{auditLogEntryID}": {
		ResponseType:       &types.AuditLogEntry{},
		ResponseTypeSchema: SchemaFromInstance(&types.AuditLogEntry{}),
	},
	"GET /api/v1/household_invitations/received": {
		ResponseType:       &types.HouseholdInvitation{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInvitation{}),
	},
	"GET /api/v1/household_invitations/sent": {
		ResponseType:       &types.HouseholdInvitation{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInvitation{}),
	},
	"GET /api/v1/household_invitations/{householdInvitationID}/": {
		ResponseType:       &types.HouseholdInvitation{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInvitation{}),
	},
	"PUT /api/v1/household_invitations/{householdInvitationID}/accept": {
		ResponseType:       &types.HouseholdInvitation{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInvitation{}),
		InputType:          &types.HouseholdInvitationUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.HouseholdInvitationUpdateRequestInput{}),
	},
	"PUT /api/v1/household_invitations/{householdInvitationID}/cancel": {
		ResponseType:       &types.HouseholdInvitation{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInvitation{}),
		InputType:          &types.HouseholdInvitationUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.HouseholdInvitationUpdateRequestInput{}),
	},
	"PUT /api/v1/household_invitations/{householdInvitationID}/reject": {
		ResponseType:       &types.HouseholdInvitation{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInvitation{}),
		InputType:          &types.HouseholdInvitationUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.HouseholdInvitationUpdateRequestInput{}),
	},
	"GET /api/v1/households/": {
		ResponseType:       &types.Household{},
		ResponseTypeSchema: SchemaFromInstance(&types.Household{}),
		ListRoute:          true,
	},
	"POST /api/v1/households/": {
		ResponseType:       &types.Household{},
		ResponseTypeSchema: SchemaFromInstance(&types.Household{}),
		InputType:          &types.HouseholdCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.HouseholdCreationRequestInput{}),
	},
	"GET /api/v1/households/current": {
		ResponseType:       &types.Household{},
		ResponseTypeSchema: SchemaFromInstance(&types.Household{}),
	},
	"POST /api/v1/households/instruments/": {
		ResponseType:       &types.HouseholdInstrumentOwnership{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInstrumentOwnership{}),
		InputType:          &types.HouseholdInstrumentOwnershipCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.HouseholdInstrumentOwnershipCreationRequestInput{}),
	},
	"GET /api/v1/households/instruments/": {
		ResponseType:       &types.HouseholdInstrumentOwnership{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInstrumentOwnership{}),
		ListRoute:          true,
	},
	"GET /api/v1/households/instruments/{householdInstrumentOwnershipID}/": {
		ResponseType:       &types.HouseholdInstrumentOwnership{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInstrumentOwnership{}),
	},
	"PUT /api/v1/households/instruments/{householdInstrumentOwnershipID}/": {
		ResponseType:       &types.HouseholdInstrumentOwnership{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInstrumentOwnership{}),
		InputType:          &types.HouseholdInstrumentOwnershipUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.HouseholdInstrumentOwnershipUpdateRequestInput{}),
	},
	"DELETE /api/v1/households/instruments/{householdInstrumentOwnershipID}/": {
		ResponseType:       &types.HouseholdInstrumentOwnership{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInstrumentOwnership{}),
	},
	"PUT /api/v1/households/{householdID}/": {
		ResponseType:       &types.Household{},
		ResponseTypeSchema: SchemaFromInstance(&types.Household{}),
		InputType:          &types.HouseholdCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.HouseholdCreationRequestInput{}),
	},
	"DELETE /api/v1/households/{householdID}/": {
		ResponseType:       &types.Household{},
		ResponseTypeSchema: SchemaFromInstance(&types.Household{}),
	},
	"GET /api/v1/households/{householdID}/": {
		ResponseType:       &types.Household{},
		ResponseTypeSchema: SchemaFromInstance(&types.Household{}),
	},
	"POST /api/v1/households/{householdID}/default": {
		ResponseType:       &types.Household{},
		ResponseTypeSchema: SchemaFromInstance(&types.Household{}),
	},
	"GET /api/v1/households/{householdID}/invitations/{householdInvitationID}/": {
		ResponseType:       &types.HouseholdInvitation{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInvitation{}),
	},
	// these are duplicate routes lol
	"POST /api/v1/households/{householdID}/invitations/": {
		ResponseType:       &types.HouseholdInvitation{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInvitation{}),
		InputType:          &types.HouseholdInvitationCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.HouseholdInvitationCreationRequestInput{}),
	},
	"POST /api/v1/households/{householdID}/invite": {
		ResponseType:       &types.HouseholdInvitation{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdInvitation{}),
		InputType:          &types.HouseholdInvitationCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.HouseholdInvitationCreationRequestInput{}),
	},
	"DELETE /api/v1/households/{householdID}/members/{userID}": {
		ResponseType:       &types.HouseholdUserMembership{},
		ResponseTypeSchema: SchemaFromInstance(&types.HouseholdUserMembership{}),
	},
	"PATCH /api/v1/households/{householdID}/members/{userID}/permissions": {
		ResponseType:       &types.UserPermissionsResponse{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserPermissionsResponse{}),
	},
	"POST /api/v1/households/{householdID}/transfer": {
		ResponseType:       &types.Household{},
		ResponseTypeSchema: SchemaFromInstance(&types.Household{}),
		InputType:          &types.HouseholdOwnershipTransferInput{},
		InputTypeSchema:    SchemaFromInstance(&types.HouseholdOwnershipTransferInput{}),
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/": {
		ResponseType:       &types.MealPlanOptionVote{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanOptionVote{}),
		ListRoute:          true,
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/": {
		ResponseType:       &types.MealPlanOptionVote{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanOptionVote{}),
		InputType:          &types.MealPlanOptionVoteUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealPlanOptionVoteUpdateRequestInput{}),
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/": {
		ResponseType:       &types.MealPlanOptionVote{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanOptionVote{}),
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/": {
		ResponseType:       &types.MealPlanOptionVote{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanOptionVote{}),
	},
	"POST /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/": {
		ResponseType:       &types.MealPlanOption{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanOption{}),
		InputType:          &types.MealPlanOptionCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealPlanOptionCreationRequestInput{}),
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/": {
		ResponseType:       &types.MealPlanOption{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanOption{}),
		ListRoute:          true,
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/": {
		ResponseType:       &types.MealPlanOption{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanOption{}),
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/": {
		ResponseType:       &types.MealPlanOption{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanOption{}),
		InputType:          &types.MealPlanOptionUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealPlanOptionUpdateRequestInput{}),
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/": {
		ResponseType:       &types.MealPlanOption{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanOption{}),
	},
	"POST /api/v1/meal_plans/{mealPlanID}/events/": {
		ResponseType:       &types.MealPlanEvent{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanEvent{}),
		InputType:          &types.MealPlanEventCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealPlanEventCreationRequestInput{}),
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/": {
		ResponseType:       &types.MealPlanEvent{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanEvent{}),
		ListRoute:          true,
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/": {
		ResponseType:       &types.MealPlanEvent{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanEvent{}),
		InputType:          &types.MealPlanEventUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealPlanEventUpdateRequestInput{}),
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/": {
		ResponseType:       &types.MealPlanEvent{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanEvent{}),
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/": {
		ResponseType:       &types.MealPlanEvent{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanEvent{}),
	},
	"POST /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/vote": {
		ResponseType:       &types.MealPlanOptionVote{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanOptionVote{}),
		InputType:          &types.MealPlanOptionVoteCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealPlanOptionVoteCreationRequestInput{}),
	},
	"GET /api/v1/meal_plans/{mealPlanID}/grocery_list_items/": {
		ResponseType:       &types.MealPlanGroceryListItem{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanGroceryListItem{}),
		ListRoute:          true,
	},
	"POST /api/v1/meal_plans/{mealPlanID}/grocery_list_items/": {
		ResponseType:       &types.MealPlanGroceryListItem{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanGroceryListItem{}),
		InputType:          &types.MealPlanGroceryListItemCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealPlanGroceryListItemCreationRequestInput{}),
	},
	"GET /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/": {
		ResponseType:       &types.MealPlanGroceryListItem{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanGroceryListItem{}),
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/": {
		ResponseType:       &types.MealPlanGroceryListItem{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanGroceryListItem{}),
		InputType:          &types.MealPlanGroceryListItemUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealPlanGroceryListItemUpdateRequestInput{}),
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/": {
		ResponseType:       &types.MealPlanGroceryListItem{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanGroceryListItem{}),
	},
	"GET /api/v1/meal_plans/{mealPlanID}/tasks/": {
		ResponseType:       &types.MealPlanTask{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanTask{}),
		ListRoute:          true,
	},
	"POST /api/v1/meal_plans/{mealPlanID}/tasks/": {
		ResponseType:       &types.MealPlanTask{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanTask{}),
		InputType:          &types.MealPlanTaskCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealPlanTaskCreationRequestInput{}),
	},
	"GET /api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}/": {
		ResponseType:       &types.MealPlanTask{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanTask{}),
	},
	"PATCH /api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}/": {
		ResponseType:       &types.MealPlanTask{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlanTask{}),
	},
	"POST /api/v1/meal_plans/": {
		ResponseType:       &types.MealPlan{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlan{}),
		InputType:          &types.MealPlanCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealPlanCreationRequestInput{}),
	},
	"GET /api/v1/meal_plans/": {
		ResponseType:       &types.MealPlan{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlan{}),
		ListRoute:          true,
	},
	"GET /api/v1/meal_plans/{mealPlanID}/": {
		ResponseType:       &types.MealPlan{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlan{}),
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/": {
		ResponseType:       &types.MealPlan{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlan{}),
		InputType:          &types.MealPlanUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealPlanUpdateRequestInput{}),
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/": {
		ResponseType:       &types.MealPlan{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlan{}),
	},
	"POST /api/v1/meal_plans/{mealPlanID}/finalize": {
		ResponseType:       &types.MealPlan{},
		ResponseTypeSchema: SchemaFromInstance(&types.MealPlan{}),
		// No input type for this route
	},
	"POST /api/v1/meals/": {
		ResponseType:       &types.Meal{},
		ResponseTypeSchema: SchemaFromInstance(&types.Meal{}),
		InputType:          &types.MealCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.MealCreationRequestInput{}),
	},
	"GET /api/v1/meals/": {
		ResponseType:       &types.Meal{},
		ResponseTypeSchema: SchemaFromInstance(&types.Meal{}),
		ListRoute:          true,
	},
	"GET /api/v1/meals/search": {
		ResponseType:       &types.Meal{},
		ResponseTypeSchema: SchemaFromInstance(&types.Meal{}),
	},
	"DELETE /api/v1/meals/{mealID}/": {
		ResponseType:       &types.Meal{},
		ResponseTypeSchema: SchemaFromInstance(&types.Meal{}),
	},
	"GET /api/v1/meals/{mealID}/": {
		ResponseType:       &types.Meal{},
		ResponseTypeSchema: SchemaFromInstance(&types.Meal{}),
	},
	"GET /api/v1/oauth2_clients/": {
		ResponseType:       &types.OAuth2Client{},
		ResponseTypeSchema: SchemaFromInstance(&types.OAuth2Client{}),
		ListRoute:          true,
	},
	"POST /api/v1/oauth2_clients/": {
		ResponseType:       &types.OAuth2Client{},
		ResponseTypeSchema: SchemaFromInstance(&types.OAuth2Client{}),
		InputType:          &types.OAuth2ClientCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.OAuth2ClientCreationRequestInput{}),
	},
	"GET /api/v1/oauth2_clients/{oauth2ClientID}/": {
		ResponseType:       &types.OAuth2Client{},
		ResponseTypeSchema: SchemaFromInstance(&types.OAuth2Client{}),
	},
	"DELETE /api/v1/oauth2_clients/{oauth2ClientID}/": {
		ResponseType:       &types.OAuth2Client{},
		ResponseTypeSchema: SchemaFromInstance(&types.OAuth2Client{}),
	},
	"POST /api/v1/recipes/{recipeID}/prep_tasks/": {
		ResponseType:       &types.RecipePrepTask{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipePrepTask{}),
		InputType:          &types.RecipePrepTaskCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipePrepTaskCreationRequestInput{}),
	},
	"GET /api/v1/recipes/{recipeID}/prep_tasks/": {
		ResponseType:       &types.RecipePrepTask{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipePrepTask{}),
		ListRoute:          true,
	},
	"DELETE /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/": {
		ResponseType:       &types.RecipePrepTask{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipePrepTask{}),
	},
	"GET /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/": {
		ResponseType:       &types.RecipePrepTask{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipePrepTask{}),
	},
	"PUT /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/": {
		ResponseType:       &types.RecipePrepTask{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipePrepTask{}),
		InputType:          &types.RecipePrepTaskUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipePrepTaskUpdateRequestInput{}),
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/": {
		ResponseType:       &types.RecipeStepCompletionCondition{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepCompletionCondition{}),
		ListRoute:          true,
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/": {
		ResponseType:       &types.RecipeStepCompletionCondition{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepCompletionCondition{}),
		InputType:          &types.RecipeStepCompletionConditionCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepCompletionConditionCreationRequestInput{}),
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/": {
		ResponseType:       &types.RecipeStepCompletionCondition{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepCompletionCondition{}),
		InputType:          &types.RecipeStepCompletionConditionUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepCompletionConditionUpdateRequestInput{}),
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/": {
		ResponseType:       &types.RecipeStepCompletionCondition{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepCompletionCondition{}),
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/": {
		ResponseType:       &types.RecipeStepCompletionCondition{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepCompletionCondition{}),
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/": {
		ResponseType:       &types.RecipeStepIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepIngredient{}),
		InputType:          &types.RecipeStepIngredientCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepIngredientCreationRequestInput{}),
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/": {
		ResponseType:       &types.RecipeStepIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepIngredient{}),
		ListRoute:          true,
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/": {
		ResponseType:       &types.RecipeStepIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepIngredient{}),
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/": {
		ResponseType:       &types.RecipeStepIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepIngredient{}),
		InputType:          &types.RecipeStepIngredientUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepIngredientUpdateRequestInput{}),
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/": {
		ResponseType:       &types.RecipeStepIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepIngredient{}),
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/": {
		ResponseType:       &types.RecipeStepInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepInstrument{}),
		InputType:          &types.RecipeStepInstrumentCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepInstrumentCreationRequestInput{}),
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/": {
		ResponseType:       &types.RecipeStepInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepInstrument{}),
		ListRoute:          true,
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/": {
		ResponseType:       &types.RecipeStepInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepInstrument{}),
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/": {
		ResponseType:       &types.RecipeStepInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepInstrument{}),
		InputType:          &types.RecipeStepInstrumentUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepInstrumentUpdateRequestInput{}),
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/": {
		ResponseType:       &types.RecipeStepInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepInstrument{}),
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/": {
		ResponseType:       &types.RecipeStepProduct{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepProduct{}),
		ListRoute:          true,
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/": {
		ResponseType:       &types.RecipeStepProduct{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepProduct{}),
		InputType:          &types.RecipeStepProductCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepProductCreationRequestInput{}),
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/": {
		ResponseType:       &types.RecipeStepProduct{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepProduct{}),
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/": {
		ResponseType:       &types.RecipeStepProduct{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepProduct{}),
		InputType:          &types.RecipeStepProductUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepProductUpdateRequestInput{}),
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/": {
		ResponseType:       &types.RecipeStepProduct{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepProduct{}),
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/": {
		ResponseType:       &types.RecipeStepVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepVessel{}),
		InputType:          &types.RecipeStepVesselCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepVesselCreationRequestInput{}),
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/": {
		ResponseType:       &types.RecipeStepVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepVessel{}),
		ListRoute:          true,
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/": {
		ResponseType:       &types.RecipeStepVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepVessel{}),
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/": {
		ResponseType:       &types.RecipeStepVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepVessel{}),
		InputType:          &types.RecipeStepVesselUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepVesselUpdateRequestInput{}),
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/": {
		ResponseType:       &types.RecipeStepVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStepVessel{}),
	},
	"POST /api/v1/recipes/{recipeID}/steps/": {
		ResponseType:       &types.RecipeStep{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStep{}),
		InputType:          &types.RecipeStepCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepCreationRequestInput{}),
	},
	"GET /api/v1/recipes/{recipeID}/steps/": {
		ResponseType:       &types.RecipeStep{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStep{}),
		ListRoute:          true,
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/": {
		ResponseType:       &types.RecipeStep{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStep{}),
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/": {
		ResponseType:       &types.RecipeStep{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStep{}),
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/": {
		ResponseType:       &types.RecipeStep{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeStep{}),
		InputType:          &types.RecipeStepUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeStepUpdateRequestInput{}),
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/images": {},
	"POST /api/v1/recipes/": {
		ResponseType:       &types.Recipe{},
		ResponseTypeSchema: SchemaFromInstance(&types.Recipe{}),
		InputType:          &types.RecipeCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeCreationRequestInput{}),
	},
	"GET /api/v1/recipes/": {
		ResponseType:       &types.Recipe{},
		ResponseTypeSchema: SchemaFromInstance(&types.Recipe{}),
		ListRoute:          true,
	},
	"GET /api/v1/recipes/search": {
		ResponseType:       &types.Recipe{},
		ResponseTypeSchema: SchemaFromInstance(&types.Recipe{}),
	},
	"DELETE /api/v1/recipes/{recipeID}/": {
		ResponseType:       &types.Recipe{},
		ResponseTypeSchema: SchemaFromInstance(&types.Recipe{}),
	},
	"GET /api/v1/recipes/{recipeID}/": {
		ResponseType:       &types.Recipe{},
		ResponseTypeSchema: SchemaFromInstance(&types.Recipe{}),
	},
	"PUT /api/v1/recipes/{recipeID}/": {
		ResponseType:       &types.Recipe{},
		ResponseTypeSchema: SchemaFromInstance(&types.Recipe{}),
		InputType:          &types.RecipeUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeUpdateRequestInput{}),
	},
	"POST /api/v1/recipes/{recipeID}/clone": {
		ResponseType:       &types.Recipe{},
		ResponseTypeSchema: SchemaFromInstance(&types.Recipe{}),
		// No input type for this route
	},
	"GET /api/v1/recipes/{recipeID}/dag": {
		ResponseType:       &types.APIError{},
		ResponseTypeSchema: SchemaFromInstance(&types.APIError{}),
	},
	"POST /api/v1/recipes/{recipeID}/images": {},
	"GET /api/v1/recipes/{recipeID}/mermaid": {},
	"GET /api/v1/recipes/{recipeID}/prep_steps": {
		ResponseType:       &types.RecipePrepTaskStep{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipePrepTaskStep{}),
	},
	"POST /api/v1/recipes/{recipeID}/ratings/": {
		ResponseType:       &types.RecipeRating{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeRating{}),
		InputType:          &types.RecipeRatingCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeRatingCreationRequestInput{}),
	},
	"GET /api/v1/recipes/{recipeID}/ratings/": {
		ResponseType:       &types.RecipeRating{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeRating{}),
		ListRoute:          true,
	},
	"PUT /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/": {
		ResponseType:       &types.RecipeRating{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeRating{}),
		InputType:          &types.RecipeRatingUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.RecipeRatingUpdateRequestInput{}),
	},
	"DELETE /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/": {
		ResponseType:       &types.RecipeRating{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeRating{}),
	},
	"GET /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/": {
		ResponseType:       &types.RecipeRating{},
		ResponseTypeSchema: SchemaFromInstance(&types.RecipeRating{}),
	},
	"POST /api/v1/settings/": {
		ResponseType:       &types.ServiceSetting{},
		ResponseTypeSchema: SchemaFromInstance(&types.ServiceSetting{}),
		InputType:          &types.ServiceSettingCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ServiceSettingCreationRequestInput{}),
	},
	"GET /api/v1/settings/": {
		ResponseType:       &types.ServiceSetting{},
		ResponseTypeSchema: SchemaFromInstance(&types.ServiceSetting{}),
		ListRoute:          true,
	},
	"POST /api/v1/settings/configurations/": {
		ResponseType:       &types.ServiceSettingConfiguration{},
		ResponseTypeSchema: SchemaFromInstance(&types.ServiceSettingConfiguration{}),
		InputType:          &types.ServiceSettingConfigurationCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ServiceSettingConfigurationCreationRequestInput{}),
	},
	"GET /api/v1/settings/configurations/household": {
		ResponseType:       &types.ServiceSettingConfiguration{},
		ResponseTypeSchema: SchemaFromInstance(&types.ServiceSettingConfiguration{}),
	},
	"GET /api/v1/settings/configurations/user": {
		ResponseType:       &types.ServiceSettingConfiguration{},
		ResponseTypeSchema: SchemaFromInstance(&types.ServiceSettingConfiguration{}),
	},
	"GET /api/v1/settings/configurations/user/{serviceSettingConfigurationName}": {
		ResponseType:       &types.ServiceSettingConfiguration{},
		ResponseTypeSchema: SchemaFromInstance(&types.ServiceSettingConfiguration{}),
	},
	"DELETE /api/v1/settings/configurations/{serviceSettingConfigurationID}": {
		ResponseType:       &types.ServiceSettingConfiguration{},
		ResponseTypeSchema: SchemaFromInstance(&types.ServiceSettingConfiguration{}),
	},
	"PUT /api/v1/settings/configurations/{serviceSettingConfigurationID}": {
		ResponseType:       &types.ServiceSettingConfiguration{},
		ResponseTypeSchema: SchemaFromInstance(&types.ServiceSettingConfiguration{}),
		InputType:          &types.ServiceSettingConfigurationUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ServiceSettingConfigurationUpdateRequestInput{}),
	},
	"GET /api/v1/settings/search": {
		ResponseType:       &types.ServiceSetting{},
		ResponseTypeSchema: SchemaFromInstance(&types.ServiceSetting{}),
	},
	"GET /api/v1/settings/{serviceSettingID}/": {
		ResponseType:       &types.ServiceSetting{},
		ResponseTypeSchema: SchemaFromInstance(&types.ServiceSetting{}),
	},
	"DELETE /api/v1/settings/{serviceSettingID}/": {
		ResponseType:       &types.ServiceSetting{},
		ResponseTypeSchema: SchemaFromInstance(&types.ServiceSetting{}),
	},
	"POST /api/v1/user_ingredient_preferences/": {
		ResponseType:       &types.UserIngredientPreference{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserIngredientPreference{}),
		InputType:          &types.UserIngredientPreferenceCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.UserIngredientPreferenceCreationRequestInput{}),
	},
	"GET /api/v1/user_ingredient_preferences/": {
		ResponseType:       &types.UserIngredientPreference{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserIngredientPreference{}),
		ListRoute:          true,
	},
	"PUT /api/v1/user_ingredient_preferences/{userIngredientPreferenceID}/": {
		ResponseType:       &types.UserIngredientPreference{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserIngredientPreference{}),
		InputType:          &types.UserIngredientPreferenceUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.UserIngredientPreferenceUpdateRequestInput{}),
	},
	"DELETE /api/v1/user_ingredient_preferences/{userIngredientPreferenceID}/": {
		ResponseType:       &types.UserIngredientPreference{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserIngredientPreference{}),
	},
	"POST /api/v1/user_notifications/": {
		ResponseType:       &types.UserNotification{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserNotification{}),
		InputType:          &types.UserNotificationCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.UserNotificationCreationRequestInput{}),
	},
	"GET /api/v1/user_notifications/": {
		ResponseType:       &types.UserNotification{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserNotification{}),
		ListRoute:          true,
	},
	"GET /api/v1/user_notifications/{userNotificationID}/": {
		ResponseType:       &types.UserNotification{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserNotification{}),
	},
	"PATCH /api/v1/user_notifications/{userNotificationID}/": {
		ResponseType:       &types.UserNotification{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserNotification{}),
	},
	"GET /api/v1/users/": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
		ListRoute:          true,
	},
	"POST /api/v1/users/avatar/upload": {},
	"PUT /api/v1/users/details": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
		InputType:          &types.UserDetailsUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.UserDetailsUpdateRequestInput{}),
	},
	"PUT /api/v1/users/email_address": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
		InputType:          &types.UserEmailAddressUpdateInput{},
		InputTypeSchema:    SchemaFromInstance(&types.UserEmailAddressUpdateInput{}),
	},
	"POST /api/v1/users/email_address_verification": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
		InputType:          &types.EmailAddressVerificationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.EmailAddressVerificationRequestInput{}),
	},
	"POST /api/v1/users/household/select": {
		ResponseType:       &types.Household{},
		ResponseTypeSchema: SchemaFromInstance(&types.Household{}),
		InputType:          &types.ChangeActiveHouseholdInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ChangeActiveHouseholdInput{}),
	},
	"PUT /api/v1/users/password/new": {
		// No output type for this route
	},
	"POST /api/v1/users/permissions/check": {
		ResponseType:       &types.UserPermissionsResponse{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserPermissionsResponse{}),
		InputType:          &types.UserPermissionsRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.UserPermissionsRequestInput{}),
	},
	"GET /api/v1/users/search": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
		ListRoute:          true,
	},
	"GET /api/v1/users/self": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
	},
	"POST /api/v1/users/totp_secret/new": {
		ResponseType:       &types.APIError{},
		ResponseTypeSchema: SchemaFromInstance(&types.APIError{}),
		InputType:          &types.TOTPSecretRefreshInput{},
		InputTypeSchema:    SchemaFromInstance(&types.TOTPSecretRefreshInput{}),
	},
	"PUT /api/v1/users/username": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
		InputType:          &types.UsernameUpdateInput{},
		InputTypeSchema:    SchemaFromInstance(&types.UsernameUpdateInput{}),
	},
	"GET /api/v1/users/{userID}/": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
	},
	"DELETE /api/v1/users/{userID}/": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
	},
	"POST /api/v1/valid_ingredient_groups/": {
		ResponseType:       &types.ValidIngredientGroup{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientGroup{}),
		InputType:          &types.ValidIngredientGroupCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientGroupCreationRequestInput{}),
	},
	"GET /api/v1/valid_ingredient_groups/": {
		ResponseType:       &types.ValidIngredientGroup{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientGroup{}),
		ListRoute:          true,
	},
	"GET /api/v1/valid_ingredient_groups/search": {
		ResponseType:       &types.ValidIngredientGroup{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientGroup{}),
	},
	"DELETE /api/v1/valid_ingredient_groups/{validIngredientGroupID}/": {
		ResponseType:       &types.ValidIngredientGroup{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientGroup{}),
	},
	"GET /api/v1/valid_ingredient_groups/{validIngredientGroupID}/": {
		ResponseType:       &types.ValidIngredientGroup{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientGroup{}),
	},
	"PUT /api/v1/valid_ingredient_groups/{validIngredientGroupID}/": {
		ResponseType:       &types.ValidIngredientGroup{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientGroup{}),
		InputType:          &types.ValidIngredientGroupUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientGroupUpdateRequestInput{}),
	},
	"POST /api/v1/valid_ingredient_measurement_units/": {
		ResponseType:       &types.ValidIngredientMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientMeasurementUnit{}),
		InputType:          &types.ValidIngredientMeasurementUnitCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientMeasurementUnitCreationRequestInput{}),
	},
	"GET /api/v1/valid_ingredient_measurement_units/": {
		ResponseType:       &types.ValidIngredientMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientMeasurementUnit{}),
		ListRoute:          true,
	},
	"GET /api/v1/valid_ingredient_measurement_units/by_ingredient/{validIngredientID}/": {
		ResponseType:       &types.ValidIngredientMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientMeasurementUnit{}),
	},
	"GET /api/v1/valid_ingredient_measurement_units/by_measurement_unit/{validMeasurementUnitID}/": {
		ResponseType:       &types.ValidIngredientMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientMeasurementUnit{}),
	},
	"GET /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/": {
		ResponseType:       &types.ValidIngredientMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientMeasurementUnit{}),
	},
	"PUT /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/": {
		ResponseType:       &types.ValidIngredientMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientMeasurementUnit{}),
		InputType:          &types.ValidIngredientMeasurementUnitUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientMeasurementUnitUpdateRequestInput{}),
	},
	"DELETE /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/": {
		ResponseType:       &types.ValidIngredientMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientMeasurementUnit{}),
	},
	"GET /api/v1/valid_ingredient_preparations/": {
		ResponseType:       &types.ValidIngredientPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientPreparation{}),
		ListRoute:          true,
	},
	"POST /api/v1/valid_ingredient_preparations/": {
		ResponseType:       &types.ValidIngredientPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientPreparation{}),
		InputType:          &types.ValidIngredientPreparationCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientPreparationCreationRequestInput{}),
	},
	"GET /api/v1/valid_ingredient_preparations/by_ingredient/{validIngredientID}/": {
		ResponseType:       &types.ValidIngredientPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientPreparation{}),
	},
	"GET /api/v1/valid_ingredient_preparations/by_preparation/{validPreparationID}/": {
		ResponseType:       &types.ValidIngredientPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientPreparation{}),
	},
	"GET /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/": {
		ResponseType:       &types.ValidIngredientPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientPreparation{}),
	},
	"PUT /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/": {
		ResponseType:       &types.ValidIngredientPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientPreparation{}),
		InputType:          &types.ValidIngredientPreparationUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientPreparationUpdateRequestInput{}),
	},
	"DELETE /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/": {
		ResponseType:       &types.ValidIngredientPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientPreparation{}),
	},
	"POST /api/v1/valid_ingredient_state_ingredients/": {
		ResponseType:       &types.ValidIngredientStateIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientStateIngredient{}),
		InputType:          &types.ValidIngredientStateIngredientCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientStateIngredientCreationRequestInput{}),
	},
	"GET /api/v1/valid_ingredient_state_ingredients/": {
		ResponseType:       &types.ValidIngredientStateIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientStateIngredient{}),
		ListRoute:          true,
	},
	"GET /api/v1/valid_ingredient_state_ingredients/by_ingredient/{validIngredientID}/": {
		ResponseType:       &types.ValidIngredientStateIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientStateIngredient{}),
	},
	"GET /api/v1/valid_ingredient_state_ingredients/by_ingredient_state/{validIngredientStateID}/": {
		ResponseType:       &types.ValidIngredientStateIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientStateIngredient{}),
	},
	"GET /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/": {
		ResponseType:       &types.ValidIngredientStateIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientStateIngredient{}),
	},
	"PUT /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/": {
		ResponseType:       &types.ValidIngredientStateIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientStateIngredient{}),
		InputType:          &types.ValidIngredientStateIngredientUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientStateIngredientUpdateRequestInput{}),
	},
	"DELETE /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/": {
		ResponseType:       &types.ValidIngredientStateIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientStateIngredient{}),
	},
	"POST /api/v1/valid_ingredient_states/": {
		ResponseType:       &types.ValidIngredientState{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientState{}),
		InputType:          &types.ValidIngredientStateCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientStateCreationRequestInput{}),
	},
	"GET /api/v1/valid_ingredient_states/": {
		ResponseType:       &types.ValidIngredientState{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientState{}),
		ListRoute:          true,
	},
	"GET /api/v1/valid_ingredient_states/search": {
		ResponseType:       &types.ValidIngredientState{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientState{}),
	},
	"PUT /api/v1/valid_ingredient_states/{validIngredientStateID}/": {
		ResponseType:       &types.ValidIngredientState{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientState{}),
		InputType:          &types.ValidIngredientStateUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientStateUpdateRequestInput{}),
	},
	"DELETE /api/v1/valid_ingredient_states/{validIngredientStateID}/": {
		ResponseType:       &types.ValidIngredientState{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientState{}),
	},
	"GET /api/v1/valid_ingredient_states/{validIngredientStateID}/": {
		ResponseType:       &types.ValidIngredientState{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredientState{}),
	},
	"POST /api/v1/valid_ingredients/": {
		ResponseType:       &types.ValidIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredient{}),
		InputType:          &types.ValidIngredientCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientCreationRequestInput{}),
	},
	"GET /api/v1/valid_ingredients/": {
		ResponseType:       &types.ValidIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredient{}),
		ListRoute:          true,
	},
	"GET /api/v1/valid_ingredients/by_preparation/{validPreparationID}/": {
		ResponseType:       &types.ValidIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredient{}),
	},
	"GET /api/v1/valid_ingredients/random": {
		ResponseType:       &types.ValidIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredient{}),
	},
	"GET /api/v1/valid_ingredients/search": {
		ResponseType:       &types.ValidIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredient{}),
	},
	"PUT /api/v1/valid_ingredients/{validIngredientID}/": {
		ResponseType:       &types.ValidIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredient{}),
		InputType:          &types.ValidIngredientUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidIngredientUpdateRequestInput{}),
	},
	"DELETE /api/v1/valid_ingredients/{validIngredientID}/": {
		ResponseType:       &types.ValidIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredient{}),
	},
	"GET /api/v1/valid_ingredients/{validIngredientID}/": {
		ResponseType:       &types.ValidIngredient{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidIngredient{}),
	},
	"GET /api/v1/valid_instruments/": {
		ResponseType:       &types.ValidInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidInstrument{}),
		ListRoute:          true,
	},
	"POST /api/v1/valid_instruments/": {
		ResponseType:       &types.ValidInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidInstrument{}),
		InputType:          &types.ValidInstrumentCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidInstrumentCreationRequestInput{}),
	},
	"GET /api/v1/valid_instruments/random": {
		ResponseType:       &types.ValidInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidInstrument{}),
	},
	"GET /api/v1/valid_instruments/search": {
		ResponseType:       &types.ValidInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidInstrument{}),
	},
	"DELETE /api/v1/valid_instruments/{validInstrumentID}/": {
		ResponseType:       &types.ValidInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidInstrument{}),
	},
	"GET /api/v1/valid_instruments/{validInstrumentID}/": {
		ResponseType:       &types.ValidInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidInstrument{}),
	},
	"PUT /api/v1/valid_instruments/{validInstrumentID}/": {
		ResponseType:       &types.ValidInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidInstrument{}),
		InputType:          &types.ValidInstrumentUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidInstrumentUpdateRequestInput{}),
	},
	"POST /api/v1/valid_measurement_conversions/": {
		ResponseType:       &types.ValidMeasurementUnitConversion{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnitConversion{}),
		InputType:          &types.ValidMeasurementUnitConversionCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidMeasurementUnitConversionCreationRequestInput{}),
	},
	"GET /api/v1/valid_measurement_conversions/from_unit/{validMeasurementUnitID}": {
		ResponseType:       &types.ValidMeasurementUnitConversion{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnitConversion{}),
	},
	"GET /api/v1/valid_measurement_conversions/to_unit/{validMeasurementUnitID}": {
		ResponseType:       &types.ValidMeasurementUnitConversion{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnitConversion{}),
	},
	"PUT /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/": {
		ResponseType:       &types.ValidMeasurementUnitConversion{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnitConversion{}),
		InputType:          &types.ValidMeasurementUnitConversionUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidMeasurementUnitConversionUpdateRequestInput{}),
	},
	"DELETE /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/": {
		ResponseType:       &types.ValidMeasurementUnitConversion{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnitConversion{}),
	},
	"GET /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/": {
		ResponseType:       &types.ValidMeasurementUnitConversion{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnitConversion{}),
	},
	"POST /api/v1/valid_measurement_units/": {
		ResponseType:       &types.ValidMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnit{}),
		InputType:          &types.ValidMeasurementUnitCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidMeasurementUnitCreationRequestInput{}),
	},
	"GET /api/v1/valid_measurement_units/": {
		ResponseType:       &types.ValidMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnit{}),
		ListRoute:          true,
	},
	"GET /api/v1/valid_measurement_units/by_ingredient/{validIngredientID}": {
		ResponseType:       &types.ValidMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnit{}),
	},
	"GET /api/v1/valid_measurement_units/search": {
		ResponseType:       &types.ValidMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnit{}),
	},
	"GET /api/v1/valid_measurement_units/{validMeasurementUnitID}/": {
		ResponseType:       &types.ValidMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnit{}),
	},
	"PUT /api/v1/valid_measurement_units/{validMeasurementUnitID}/": {
		ResponseType:       &types.ValidMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnit{}),
		InputType:          &types.ValidMeasurementUnitUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidMeasurementUnitUpdateRequestInput{}),
	},
	"DELETE /api/v1/valid_measurement_units/{validMeasurementUnitID}/": {
		ResponseType:       &types.ValidMeasurementUnit{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidMeasurementUnit{}),
	},
	"GET /api/v1/valid_preparation_instruments/": {
		ResponseType:       &types.ValidPreparationInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationInstrument{}),
		ListRoute:          true,
	},
	"POST /api/v1/valid_preparation_instruments/": {
		ResponseType:       &types.ValidPreparationInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationInstrument{}),
		InputType:          &types.ValidPreparationInstrumentCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidPreparationInstrumentCreationRequestInput{}),
	},
	"GET /api/v1/valid_preparation_instruments/by_instrument/{validInstrumentID}/": {
		ResponseType:       &types.ValidPreparationInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationInstrument{}),
	},
	"GET /api/v1/valid_preparation_instruments/by_preparation/{validPreparationID}/": {
		ResponseType:       &types.ValidPreparationInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationInstrument{}),
	},
	"DELETE /api/v1/valid_preparation_instruments/{validPreparationVesselID}/": {
		ResponseType:       &types.ValidPreparationInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationInstrument{}),
	},
	"GET /api/v1/valid_preparation_instruments/{validPreparationVesselID}/": {
		ResponseType:       &types.ValidPreparationInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationInstrument{}),
	},
	"PUT /api/v1/valid_preparation_instruments/{validPreparationVesselID}/": {
		ResponseType:       &types.ValidPreparationInstrument{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationInstrument{}),
		InputType:          &types.ValidPreparationInstrumentUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidPreparationInstrumentUpdateRequestInput{}),
	},
	"POST /api/v1/valid_preparation_vessels/": {
		ResponseType:       &types.ValidPreparationVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationVessel{}),
		InputType:          &types.ValidPreparationVesselCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidPreparationVesselCreationRequestInput{}),
	},
	"GET /api/v1/valid_preparation_vessels/": {
		ResponseType:       &types.ValidPreparationVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationVessel{}),
		ListRoute:          true,
	},
	"GET /api/v1/valid_preparation_vessels/by_preparation/{validPreparationID}/": {
		ResponseType:       &types.ValidPreparationVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationVessel{}),
	},
	"GET /api/v1/valid_preparation_vessels/by_vessel/{ValidVesselID}/": {
		ResponseType:       &types.ValidPreparationVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationVessel{}),
	},
	"PUT /api/v1/valid_preparation_vessels/{validPreparationVesselID}/": {
		ResponseType:       &types.ValidPreparationVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationVessel{}),
		InputType:          &types.ValidPreparationVesselUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidPreparationVesselUpdateRequestInput{}),
	},
	"DELETE /api/v1/valid_preparation_vessels/{validPreparationVesselID}/": {
		ResponseType:       &types.ValidPreparationVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationVessel{}),
	},
	"GET /api/v1/valid_preparation_vessels/{validPreparationVesselID}/": {
		ResponseType:       &types.ValidPreparationVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparationVessel{}),
	},
	"GET /api/v1/valid_preparations/": {
		ResponseType:       &types.ValidPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparation{}),
		ListRoute:          true,
	},
	"POST /api/v1/valid_preparations/": {
		ResponseType:       &types.ValidPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparation{}),
		InputType:          &types.ValidPreparationCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidPreparationCreationRequestInput{}),
	},
	"GET /api/v1/valid_preparations/random": {
		ResponseType:       &types.ValidPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparation{}),
	},
	"GET /api/v1/valid_preparations/search": {
		ResponseType:       &types.ValidPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparation{}),
	},
	"PUT /api/v1/valid_preparations/{validPreparationID}/": {
		ResponseType:       &types.ValidPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparation{}),
		InputType:          &types.ValidPreparationUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidPreparationUpdateRequestInput{}),
	},
	"DELETE /api/v1/valid_preparations/{validPreparationID}/": {
		ResponseType:       &types.ValidPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparation{}),
	},
	"GET /api/v1/valid_preparations/{validPreparationID}/": {
		ResponseType:       &types.ValidPreparation{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidPreparation{}),
	},
	"POST /api/v1/valid_vessels/": {
		ResponseType:       &types.ValidVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidVessel{}),
		InputType:          &types.ValidVesselCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidVesselCreationRequestInput{}),
	},
	"GET /api/v1/valid_vessels/": {
		ResponseType:       &types.ValidVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidVessel{}),
		ListRoute:          true,
	},
	"GET /api/v1/valid_vessels/random": {
		ResponseType:       &types.ValidVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidVessel{}),
	},
	"GET /api/v1/valid_vessels/search": {
		ResponseType:       &types.ValidVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidVessel{}),
	},
	"GET /api/v1/valid_vessels/{validVesselID}/": {
		ResponseType:       &types.ValidVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidVessel{}),
	},
	"PUT /api/v1/valid_vessels/{validVesselID}/": {
		ResponseType:       &types.ValidVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidVessel{}),
		InputType:          &types.ValidVesselUpdateRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.ValidVesselUpdateRequestInput{}),
	},
	"DELETE /api/v1/valid_vessels/{validVesselID}/": {
		ResponseType:       &types.ValidVessel{},
		ResponseTypeSchema: SchemaFromInstance(&types.ValidVessel{}),
	},
	"GET /api/v1/webhooks/": {
		ResponseType:       &types.Webhook{},
		ResponseTypeSchema: SchemaFromInstance(&types.Webhook{}),
		ListRoute:          true,
	},
	"POST /api/v1/webhooks/": {
		ResponseType:       &types.Webhook{},
		ResponseTypeSchema: SchemaFromInstance(&types.Webhook{}),
		InputType:          &types.WebhookCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.WebhookCreationRequestInput{}),
	},
	"GET /api/v1/webhooks/{webhookID}/": {
		ResponseType:       &types.Webhook{},
		ResponseTypeSchema: SchemaFromInstance(&types.Webhook{}),
	},
	"DELETE /api/v1/webhooks/{webhookID}/": {
		ResponseType:       &types.Webhook{},
		ResponseTypeSchema: SchemaFromInstance(&types.Webhook{}),
	},
	"POST /api/v1/webhooks/{webhookID}/trigger_events": {
		ResponseType:       &types.WebhookTriggerEvent{},
		ResponseTypeSchema: SchemaFromInstance(&types.WebhookTriggerEvent{}),
		InputType:          &types.WebhookTriggerEventCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.WebhookTriggerEventCreationRequestInput{}),
	},
	"DELETE /api/v1/webhooks/{webhookID}/trigger_events/{webhookTriggerEventID}/": {
		ResponseType:       &types.WebhookTriggerEvent{},
		ResponseTypeSchema: SchemaFromInstance(&types.WebhookTriggerEvent{}),
	},
	"POST /api/v1/workers/finalize_meal_plans": {
		ResponseType:       &types.FinalizeMealPlansRequest{},
		ResponseTypeSchema: SchemaFromInstance(&types.FinalizeMealPlansRequest{}),
		InputType:          &types.FinalizeMealPlansRequest{},
		InputTypeSchema:    SchemaFromInstance(&types.FinalizeMealPlansRequest{}),
	},
	"POST /api/v1/workers/meal_plan_grocery_list_init": {
		// no input or output types for this route
	},
	"POST /api/v1/workers/meal_plan_tasks": {
		// no input or output types for this route
	},
	"GET /auth/status": {
		ResponseType:       &types.UserStatusResponse{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserStatusResponse{}),
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
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
		InputType:          &types.UserRegistrationInput{},
		InputTypeSchema:    SchemaFromInstance(&types.UserRegistrationInput{}),
	},
	"POST /users/email_address/verify": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
		InputType:          &types.EmailAddressVerificationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.EmailAddressVerificationRequestInput{}),
	},
	"POST /users/login": {
		ResponseType:       &types.UserStatusResponse{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserStatusResponse{}),
		InputType:          &types.UserLoginInput{},
		InputTypeSchema:    SchemaFromInstance(&types.UserLoginInput{}),
	},
	"POST /users/login/admin": {
		ResponseType:       &types.UserStatusResponse{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserStatusResponse{}),
		InputType:          &types.UserLoginInput{},
		InputTypeSchema:    SchemaFromInstance(&types.UserLoginInput{}),
	},
	"POST /users/logout": {
		ResponseType:       &types.UserStatusResponse{},
		ResponseTypeSchema: SchemaFromInstance(&types.UserStatusResponse{}),
	},
	"POST /users/password/reset": {
		ResponseType:       &types.PasswordResetToken{},
		ResponseTypeSchema: SchemaFromInstance(&types.PasswordResetToken{}),
		InputType:          &types.PasswordResetTokenCreationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.PasswordResetTokenCreationRequestInput{}),
	},
	"POST /users/password/reset/redeem": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
		InputType:          &types.PasswordResetTokenRedemptionRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.PasswordResetTokenRedemptionRequestInput{}),
	},
	"POST /users/totp_secret/verify": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
		InputType:          &types.EmailAddressVerificationRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.EmailAddressVerificationRequestInput{}),
	},
	"POST /users/username/reminder": {
		ResponseType:       &types.User{},
		ResponseTypeSchema: SchemaFromInstance(&types.User{}),
		InputType:          &types.UsernameReminderRequestInput{},
		InputTypeSchema:    SchemaFromInstance(&types.UsernameReminderRequestInput{}),
	},
}
