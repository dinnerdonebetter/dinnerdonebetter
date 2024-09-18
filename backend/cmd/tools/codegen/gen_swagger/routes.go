package main

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
)

type routeDetails struct {
	ResponseType  any
	InputType     any
	InputTypeName string
	OAuth2Scopes  []string
	ListRoute     bool
}

const (
	serviceAdmin    = "service_admin"
	householdAdmin  = "household_admin"
	householdMember = "household_member"
)

var routeInfoMap = map[string]routeDetails{
	"GET /_meta_/live":  {},
	"GET /_meta_/ready": {},
	"POST /api/v1/admin/cycle_cookie_secret": {
		OAuth2Scopes: []string{serviceAdmin},
	},
	"POST /api/v1/admin/users/status": {
		ResponseType: &types.UserStatusResponse{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/audit_log_entries/for_household": {
		ResponseType: &types.AuditLogEntry{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"GET /api/v1/audit_log_entries/for_user": {
		ResponseType: &types.AuditLogEntry{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/audit_log_entries/{auditLogEntryID}": {
		ResponseType: &types.AuditLogEntry{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/household_invitations/received": {
		ResponseType: &types.HouseholdInvitation{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/household_invitations/sent": {
		ResponseType: &types.HouseholdInvitation{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/household_invitations/{householdInvitationID}/": {
		ResponseType: &types.HouseholdInvitation{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/household_invitations/{householdInvitationID}/accept": {
		ResponseType: &types.HouseholdInvitation{},
		InputType:    &types.HouseholdInvitationUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/household_invitations/{householdInvitationID}/cancel": {
		ResponseType: &types.HouseholdInvitation{},
		InputType:    &types.HouseholdInvitationUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/household_invitations/{householdInvitationID}/reject": {
		ResponseType: &types.HouseholdInvitation{},
		InputType:    &types.HouseholdInvitationUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/households/": {
		ResponseType: &types.Household{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/households/": {
		ResponseType: &types.Household{},
		InputType:    &types.HouseholdCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/households/current": {
		ResponseType: &types.Household{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/households/instruments/": {
		ResponseType: &types.HouseholdInstrumentOwnership{},
		InputType:    &types.HouseholdInstrumentOwnershipCreationRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"GET /api/v1/households/instruments/": {
		ResponseType: &types.HouseholdInstrumentOwnership{},
		OAuth2Scopes: []string{householdMember},
		ListRoute:    true,
	},
	"GET /api/v1/households/instruments/{householdInstrumentOwnershipID}/": {
		ResponseType: &types.HouseholdInstrumentOwnership{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/households/instruments/{householdInstrumentOwnershipID}/": {
		ResponseType: &types.HouseholdInstrumentOwnership{},
		InputType:    &types.HouseholdInstrumentOwnershipUpdateRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"DELETE /api/v1/households/instruments/{householdInstrumentOwnershipID}/": {
		ResponseType: &types.HouseholdInstrumentOwnership{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"PUT /api/v1/households/{householdID}/": {
		ResponseType: &types.Household{},
		InputType:    &types.HouseholdCreationRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"DELETE /api/v1/households/{householdID}/": {
		ResponseType: &types.Household{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"GET /api/v1/households/{householdID}/": {
		ResponseType: &types.Household{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/households/{householdID}/default": {
		ResponseType: &types.Household{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/households/{householdID}/invitations/{householdInvitationID}/": {
		ResponseType: &types.HouseholdInvitation{},
		OAuth2Scopes: []string{householdMember},
	},
	// these are duplicate routes lol
	"POST /api/v1/households/{householdID}/invitations/": {
		ResponseType: &types.HouseholdInvitation{},
		InputType:    &types.HouseholdInvitationCreationRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"POST /api/v1/households/{householdID}/invite": {
		ResponseType: &types.HouseholdInvitation{},
		InputType:    &types.HouseholdInvitationCreationRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"DELETE /api/v1/households/{householdID}/members/{userID}": {
		ResponseType: &types.HouseholdUserMembership{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"PATCH /api/v1/households/{householdID}/members/{userID}/permissions": {
		ResponseType: &types.UserPermissionsResponse{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"POST /api/v1/households/{householdID}/transfer": {
		ResponseType: &types.Household{},
		InputType:    &types.HouseholdOwnershipTransferInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/": {
		ResponseType: &types.MealPlanOptionVote{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/": {
		ResponseType: &types.MealPlanOptionVote{},
		InputType:    &types.MealPlanOptionVoteUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/": {
		ResponseType: &types.MealPlanOptionVote{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/": {
		ResponseType: &types.MealPlanOptionVote{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/": {
		ResponseType: &types.MealPlanOption{},
		InputType:    &types.MealPlanOptionCreationRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/": {
		ResponseType: &types.MealPlanOption{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/": {
		ResponseType: &types.MealPlanOption{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/": {
		ResponseType: &types.MealPlanOption{},
		InputType:    &types.MealPlanOptionUpdateRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/": {
		ResponseType: &types.MealPlanOption{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"POST /api/v1/meal_plans/{mealPlanID}/events/": {
		ResponseType: &types.MealPlanEvent{},
		InputType:    &types.MealPlanEventCreationRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/": {
		ResponseType: &types.MealPlanEvent{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/": {
		ResponseType: &types.MealPlanEvent{},
		InputType:    &types.MealPlanEventUpdateRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/": {
		ResponseType: &types.MealPlanEvent{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/": {
		ResponseType: &types.MealPlanEvent{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/vote": {
		ResponseType: &types.MealPlanOptionVote{},
		InputType:    &types.MealPlanOptionVoteCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/grocery_list_items/": {
		ResponseType: &types.MealPlanGroceryListItem{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/meal_plans/{mealPlanID}/grocery_list_items/": {
		ResponseType: &types.MealPlanGroceryListItem{},
		InputType:    &types.MealPlanGroceryListItemCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/": {
		ResponseType: &types.MealPlanGroceryListItem{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/": {
		ResponseType: &types.MealPlanGroceryListItem{},
		InputType:    &types.MealPlanGroceryListItemUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/": {
		ResponseType: &types.MealPlanGroceryListItem{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/tasks/": {
		ResponseType: &types.MealPlanTask{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/meal_plans/{mealPlanID}/tasks/": {
		ResponseType: &types.MealPlanTask{},
		InputType:    &types.MealPlanTaskCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}/": {
		ResponseType: &types.MealPlanTask{},
		OAuth2Scopes: []string{householdMember},
	},
	"PATCH /api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}/": {
		ResponseType: &types.MealPlanTask{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/meal_plans/": {
		ResponseType: &types.MealPlan{},
		InputType:    &types.MealPlanCreationRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"GET /api/v1/meal_plans/": {
		ResponseType: &types.MealPlan{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/meal_plans/{mealPlanID}/": {
		ResponseType: &types.MealPlan{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/meal_plans/{mealPlanID}/": {
		ResponseType: &types.MealPlan{},
		InputType:    &types.MealPlanUpdateRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"DELETE /api/v1/meal_plans/{mealPlanID}/": {
		ResponseType: &types.MealPlan{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"POST /api/v1/meal_plans/{mealPlanID}/finalize": {
		ResponseType: &types.MealPlan{},
		OAuth2Scopes: []string{householdAdmin},
		// No input type for this route
	},
	"POST /api/v1/meals/": {
		ResponseType: &types.Meal{},
		InputType:    &types.MealCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/meals/": {
		ResponseType: &types.Meal{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/meals/search": {
		ResponseType: &types.Meal{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/meals/{mealID}/": {
		ResponseType: &types.Meal{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/meals/{mealID}/": {
		ResponseType: &types.Meal{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/oauth2_clients/": {
		ResponseType: &types.OAuth2Client{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/oauth2_clients/": {
		ResponseType: &types.OAuth2Client{},
		InputType:    &types.OAuth2ClientCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/oauth2_clients/{oauth2ClientID}/": {
		ResponseType: &types.OAuth2Client{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/oauth2_clients/{oauth2ClientID}/": {
		ResponseType: &types.OAuth2Client{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	// TODO: done up to here
	"POST /api/v1/recipes/{recipeID}/prep_tasks/": {
		ResponseType: &types.RecipePrepTask{},
		InputType:    &types.RecipePrepTaskCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/prep_tasks/": {
		ResponseType: &types.RecipePrepTask{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/": {
		ResponseType: &types.RecipePrepTask{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/": {
		ResponseType: &types.RecipePrepTask{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/": {
		ResponseType: &types.RecipePrepTask{},
		InputType:    &types.RecipePrepTaskUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/": {
		ResponseType: &types.RecipeStepCompletionCondition{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/": {
		ResponseType: &types.RecipeStepCompletionCondition{},
		InputType:    &types.RecipeStepCompletionConditionCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/": {
		ResponseType: &types.RecipeStepCompletionCondition{},
		InputType:    &types.RecipeStepCompletionConditionUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/": {
		ResponseType: &types.RecipeStepCompletionCondition{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/": {
		ResponseType: &types.RecipeStepCompletionCondition{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/": {
		ResponseType: &types.RecipeStepIngredient{},
		InputType:    &types.RecipeStepIngredientCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/": {
		ResponseType: &types.RecipeStepIngredient{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/": {
		ResponseType: &types.RecipeStepIngredient{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/": {
		ResponseType: &types.RecipeStepIngredient{},
		InputType:    &types.RecipeStepIngredientUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/": {
		ResponseType: &types.RecipeStepIngredient{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/": {
		ResponseType: &types.RecipeStepInstrument{},
		InputType:    &types.RecipeStepInstrumentCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/": {
		ResponseType: &types.RecipeStepInstrument{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/": {
		ResponseType: &types.RecipeStepInstrument{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/": {
		ResponseType: &types.RecipeStepInstrument{},
		InputType:    &types.RecipeStepInstrumentUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/": {
		ResponseType: &types.RecipeStepInstrument{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/": {
		ResponseType: &types.RecipeStepProduct{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/": {
		ResponseType: &types.RecipeStepProduct{},
		InputType:    &types.RecipeStepProductCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/": {
		ResponseType: &types.RecipeStepProduct{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/": {
		ResponseType: &types.RecipeStepProduct{},
		InputType:    &types.RecipeStepProductUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/": {
		ResponseType: &types.RecipeStepProduct{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/": {
		ResponseType: &types.RecipeStepVessel{},
		InputType:    &types.RecipeStepVesselCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/": {
		ResponseType: &types.RecipeStepVessel{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/": {
		ResponseType: &types.RecipeStepVessel{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/": {
		ResponseType: &types.RecipeStepVessel{},
		InputType:    &types.RecipeStepVesselUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/": {
		ResponseType: &types.RecipeStepVessel{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/recipes/{recipeID}/steps/": {
		ResponseType: &types.RecipeStep{},
		InputType:    &types.RecipeStepCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/": {
		ResponseType: &types.RecipeStep{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/": {
		ResponseType: &types.RecipeStep{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/": {
		ResponseType: &types.RecipeStep{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/": {
		ResponseType: &types.RecipeStep{},
		InputType:    &types.RecipeStepUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/images": {},
	"POST /api/v1/recipes/": {
		ResponseType: &types.Recipe{},
		InputType:    &types.RecipeCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/": {
		ResponseType: &types.Recipe{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/search": {
		ResponseType: &types.Recipe{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/recipes/{recipeID}/": {
		ResponseType: &types.Recipe{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/": {
		ResponseType: &types.Recipe{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/recipes/{recipeID}/": {
		ResponseType: &types.Recipe{},
		InputType:    &types.RecipeUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/recipes/{recipeID}/clone": {
		ResponseType: &types.Recipe{},
		// No input type for this route
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/dag": {
		ResponseType: &types.APIError{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/recipes/{recipeID}/images": {},
	"GET /api/v1/recipes/{recipeID}/mermaid": {},
	"GET /api/v1/recipes/{recipeID}/prep_steps": {
		ResponseType: &types.RecipePrepTaskStep{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/recipes/{recipeID}/ratings/": {
		ResponseType: &types.RecipeRating{},
		InputType:    &types.RecipeRatingCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/ratings/": {
		ResponseType: &types.RecipeRating{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/": {
		ResponseType: &types.RecipeRating{},
		InputType:    &types.RecipeRatingUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/": {
		ResponseType: &types.RecipeRating{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/": {
		ResponseType: &types.RecipeRating{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/settings/": {
		ResponseType: &types.ServiceSetting{},
		InputType:    &types.ServiceSettingCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/settings/": {
		ResponseType: &types.ServiceSetting{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/settings/configurations/": {
		ResponseType: &types.ServiceSettingConfiguration{},
		InputType:    &types.ServiceSettingConfigurationCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/settings/configurations/household": {
		ResponseType: &types.ServiceSettingConfiguration{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/settings/configurations/user": {
		ResponseType: &types.ServiceSettingConfiguration{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/settings/configurations/user/{serviceSettingConfigurationName}": {
		ResponseType: &types.ServiceSettingConfiguration{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/settings/configurations/{serviceSettingConfigurationID}": {
		ResponseType: &types.ServiceSettingConfiguration{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/settings/configurations/{serviceSettingConfigurationID}": {
		ResponseType: &types.ServiceSettingConfiguration{},
		InputType:    &types.ServiceSettingConfigurationUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/settings/search": {
		ResponseType: &types.ServiceSetting{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/settings/{serviceSettingID}/": {
		ResponseType: &types.ServiceSetting{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/settings/{serviceSettingID}/": {
		ResponseType: &types.ServiceSetting{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"POST /api/v1/user_ingredient_preferences/": {
		ResponseType: &types.UserIngredientPreference{},
		InputType:    &types.UserIngredientPreferenceCreationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/user_ingredient_preferences/": {
		ResponseType: &types.UserIngredientPreference{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/user_ingredient_preferences/{userIngredientPreferenceID}/": {
		ResponseType: &types.UserIngredientPreference{},
		InputType:    &types.UserIngredientPreferenceUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/user_ingredient_preferences/{userIngredientPreferenceID}/": {
		ResponseType: &types.UserIngredientPreference{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/user_notifications/": {
		ResponseType: &types.UserNotification{},
		InputType:    &types.UserNotificationCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/user_notifications/": {
		ResponseType: &types.UserNotification{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/user_notifications/{userNotificationID}/": {
		ResponseType: &types.UserNotification{},
		OAuth2Scopes: []string{householdMember},
	},
	"PATCH /api/v1/user_notifications/{userNotificationID}/": {
		ResponseType: &types.UserNotification{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/users/": {
		ResponseType: &types.User{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/users/avatar/upload": {
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/users/details": {
		ResponseType: &types.User{},
		InputType:    &types.UserDetailsUpdateRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/users/email_address": {
		ResponseType: &types.User{},
		InputType:    &types.UserEmailAddressUpdateInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/users/email_address_verification": {
		ResponseType: &types.User{},
		InputType:    &types.EmailAddressVerificationRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/users/household/select": {
		ResponseType: &types.Household{},
		InputType:    &types.ChangeActiveHouseholdInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/users/password/new": {
		// No output type for this route
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/users/permissions/check": {
		ResponseType: &types.UserPermissionsResponse{},
		InputType:    &types.UserPermissionsRequestInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/users/search": {
		ResponseType: &types.User{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/users/self": {
		ResponseType: &types.User{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/users/totp_secret/new": {
		ResponseType: &types.APIError{},
		InputType:    &types.TOTPSecretRefreshInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/users/username": {
		ResponseType: &types.User{},
		InputType:    &types.UsernameUpdateInput{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/users/{userID}/": {
		ResponseType: &types.User{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/users/{userID}/": {
		ResponseType: &types.User{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/valid_ingredient_groups/": {
		ResponseType: &types.ValidIngredientGroup{},
		InputType:    &types.ValidIngredientGroupCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_ingredient_groups/": {
		ResponseType: &types.ValidIngredientGroup{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredient_groups/search": {
		ResponseType: &types.ValidIngredientGroup{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/valid_ingredient_groups/{validIngredientGroupID}/": {
		ResponseType: &types.ValidIngredientGroup{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_ingredient_groups/{validIngredientGroupID}/": {
		ResponseType: &types.ValidIngredientGroup{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_ingredient_groups/{validIngredientGroupID}/": {
		ResponseType: &types.ValidIngredientGroup{},
		InputType:    &types.ValidIngredientGroupUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"POST /api/v1/valid_ingredient_measurement_units/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
		InputType:    &types.ValidIngredientMeasurementUnitCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_ingredient_measurement_units/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredient_measurement_units/by_ingredient/{validIngredientID}/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredient_measurement_units/by_measurement_unit/{validMeasurementUnitID}/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
		InputType:    &types.ValidIngredientMeasurementUnitUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"DELETE /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/": {
		ResponseType: &types.ValidIngredientMeasurementUnit{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_ingredient_preparations/": {
		ResponseType: &types.ValidIngredientPreparation{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/valid_ingredient_preparations/": {
		ResponseType: &types.ValidIngredientPreparation{},
		InputType:    &types.ValidIngredientPreparationCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_ingredient_preparations/by_ingredient/{validIngredientID}/": {
		ResponseType: &types.ValidIngredientPreparation{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredient_preparations/by_preparation/{validPreparationID}/": {
		ResponseType: &types.ValidIngredientPreparation{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/": {
		ResponseType: &types.ValidIngredientPreparation{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/": {
		ResponseType: &types.ValidIngredientPreparation{},
		InputType:    &types.ValidIngredientPreparationUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"DELETE /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/": {
		ResponseType: &types.ValidIngredientPreparation{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"POST /api/v1/valid_ingredient_state_ingredients/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
		InputType:    &types.ValidIngredientStateIngredientCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_ingredient_state_ingredients/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredient_state_ingredients/by_ingredient/{validIngredientID}/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredient_state_ingredients/by_ingredient_state/{validIngredientStateID}/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
		InputType:    &types.ValidIngredientStateIngredientUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"DELETE /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/": {
		ResponseType: &types.ValidIngredientStateIngredient{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"POST /api/v1/valid_ingredient_states/": {
		ResponseType: &types.ValidIngredientState{},
		InputType:    &types.ValidIngredientStateCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_ingredient_states/": {
		ResponseType: &types.ValidIngredientState{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredient_states/search": {
		ResponseType: &types.ValidIngredientState{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_ingredient_states/{validIngredientStateID}/": {
		ResponseType: &types.ValidIngredientState{},
		InputType:    &types.ValidIngredientStateUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"DELETE /api/v1/valid_ingredient_states/{validIngredientStateID}/": {
		ResponseType: &types.ValidIngredientState{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_ingredient_states/{validIngredientStateID}/": {
		ResponseType: &types.ValidIngredientState{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/valid_ingredients/": {
		ResponseType: &types.ValidIngredient{},
		InputType:    &types.ValidIngredientCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_ingredients/": {
		ResponseType: &types.ValidIngredient{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredients/by_preparation/{validPreparationID}/": {
		ResponseType: &types.ValidIngredient{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredients/random": {
		ResponseType: &types.ValidIngredient{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_ingredients/search": {
		ResponseType: &types.ValidIngredient{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_ingredients/{validIngredientID}/": {
		ResponseType: &types.ValidIngredient{},
		InputType:    &types.ValidIngredientUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"DELETE /api/v1/valid_ingredients/{validIngredientID}/": {
		ResponseType: &types.ValidIngredient{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_ingredients/{validIngredientID}/": {
		ResponseType: &types.ValidIngredient{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_instruments/": {
		ResponseType: &types.ValidInstrument{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/valid_instruments/": {
		ResponseType: &types.ValidInstrument{},
		InputType:    &types.ValidInstrumentCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_instruments/random": {
		ResponseType: &types.ValidInstrument{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_instruments/search": {
		ResponseType: &types.ValidInstrument{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/valid_instruments/{validInstrumentID}/": {
		ResponseType: &types.ValidInstrument{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_instruments/{validInstrumentID}/": {
		ResponseType: &types.ValidInstrument{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_instruments/{validInstrumentID}/": {
		ResponseType: &types.ValidInstrument{},
		InputType:    &types.ValidInstrumentUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"POST /api/v1/valid_measurement_conversions/": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
		InputType:    &types.ValidMeasurementUnitConversionCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_measurement_conversions/from_unit/{validMeasurementUnitID}": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_measurement_conversions/to_unit/{validMeasurementUnitID}": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
		InputType:    &types.ValidMeasurementUnitConversionUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"DELETE /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/": {
		ResponseType: &types.ValidMeasurementUnitConversion{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/valid_measurement_units/": {
		ResponseType: &types.ValidMeasurementUnit{},
		InputType:    &types.ValidMeasurementUnitCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_measurement_units/": {
		ResponseType: &types.ValidMeasurementUnit{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_measurement_units/by_ingredient/{validIngredientID}": {
		ResponseType: &types.ValidMeasurementUnit{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_measurement_units/search": {
		ResponseType: &types.ValidMeasurementUnit{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_measurement_units/{validMeasurementUnitID}/": {
		ResponseType: &types.ValidMeasurementUnit{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_measurement_units/{validMeasurementUnitID}/": {
		ResponseType: &types.ValidMeasurementUnit{},
		InputType:    &types.ValidMeasurementUnitUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"DELETE /api/v1/valid_measurement_units/{validMeasurementUnitID}/": {
		ResponseType: &types.ValidMeasurementUnit{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_preparation_instruments/": {
		ResponseType: &types.ValidPreparationInstrument{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/valid_preparation_instruments/": {
		ResponseType: &types.ValidPreparationInstrument{},
		InputType:    &types.ValidPreparationInstrumentCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_preparation_instruments/by_instrument/{validInstrumentID}/": {
		ResponseType: &types.ValidPreparationInstrument{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_preparation_instruments/by_preparation/{validPreparationID}/": {
		ResponseType: &types.ValidPreparationInstrument{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/valid_preparation_instruments/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationInstrument{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_preparation_instruments/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationInstrument{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_preparation_instruments/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationInstrument{},
		InputType:    &types.ValidPreparationInstrumentUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"POST /api/v1/valid_preparation_vessels/": {
		ResponseType: &types.ValidPreparationVessel{},
		InputType:    &types.ValidPreparationVesselCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_preparation_vessels/": {
		ResponseType: &types.ValidPreparationVessel{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_preparation_vessels/by_preparation/{validPreparationID}/": {
		ResponseType: &types.ValidPreparationVessel{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_preparation_vessels/by_vessel/{ValidVesselID}/": {
		ResponseType: &types.ValidPreparationVessel{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_preparation_vessels/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationVessel{},
		InputType:    &types.ValidPreparationVesselUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"DELETE /api/v1/valid_preparation_vessels/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationVessel{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_preparation_vessels/{validPreparationVesselID}/": {
		ResponseType: &types.ValidPreparationVessel{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_preparations/": {
		ResponseType: &types.ValidPreparation{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/valid_preparations/": {
		ResponseType: &types.ValidPreparation{},
		InputType:    &types.ValidPreparationCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_preparations/random": {
		ResponseType: &types.ValidPreparation{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_preparations/search": {
		ResponseType: &types.ValidPreparation{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_preparations/{validPreparationID}/": {
		ResponseType: &types.ValidPreparation{},
		InputType:    &types.ValidPreparationUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"DELETE /api/v1/valid_preparations/{validPreparationID}/": {
		ResponseType: &types.ValidPreparation{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_preparations/{validPreparationID}/": {
		ResponseType: &types.ValidPreparation{},
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/valid_vessels/": {
		ResponseType: &types.ValidVessel{},
		InputType:    &types.ValidVesselCreationRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/valid_vessels/": {
		ResponseType: &types.ValidVessel{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_vessels/random": {
		ResponseType: &types.ValidVessel{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_vessels/search": {
		ResponseType: &types.ValidVessel{},
		OAuth2Scopes: []string{householdMember},
	},
	"GET /api/v1/valid_vessels/{validVesselID}/": {
		ResponseType: &types.ValidVessel{},
		OAuth2Scopes: []string{householdMember},
	},
	"PUT /api/v1/valid_vessels/{validVesselID}/": {
		ResponseType: &types.ValidVessel{},
		InputType:    &types.ValidVesselUpdateRequestInput{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"DELETE /api/v1/valid_vessels/{validVesselID}/": {
		ResponseType: &types.ValidVessel{},
		OAuth2Scopes: []string{serviceAdmin},
	},
	"GET /api/v1/webhooks/": {
		ResponseType: &types.Webhook{},
		ListRoute:    true,
		OAuth2Scopes: []string{householdMember},
	},
	"POST /api/v1/webhooks/": {
		ResponseType: &types.Webhook{},
		InputType:    &types.WebhookCreationRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"GET /api/v1/webhooks/{webhookID}/": {
		ResponseType: &types.Webhook{},
		OAuth2Scopes: []string{householdMember},
	},
	"DELETE /api/v1/webhooks/{webhookID}/": {
		ResponseType: &types.Webhook{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"POST /api/v1/webhooks/{webhookID}/trigger_events": {
		ResponseType: &types.WebhookTriggerEvent{},
		InputType:    &types.WebhookTriggerEventCreationRequestInput{},
		OAuth2Scopes: []string{householdAdmin},
	},
	"DELETE /api/v1/webhooks/{webhookID}/trigger_events/{webhookTriggerEventID}/": {
		ResponseType: &types.WebhookTriggerEvent{},
		OAuth2Scopes: []string{householdAdmin},
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
