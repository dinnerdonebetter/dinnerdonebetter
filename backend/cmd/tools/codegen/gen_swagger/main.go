package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/server/http/build"
)

type RouteDefinition struct {
	Method        string
	Summary       string
	Path          string
	PathArguments []string
	ListRoute     bool
	Tags          []string
	ResponseType  string
}

func (d *RouteDefinition) Render() string {
	return fmt.Sprintf(`  {{ .Path }}:
		{{ .Method }}:
      tags:
        - recipes
      {{- if .Summary }} summary: {{ .Summary }} {{ end }}
      security:
        - oAuth2:
        - cookieAuth:
      {{- if or .ListRoute .PathArguments }} 
      parameters:
      	{{- if .ListRoute }}
        - in: query
          name: page
          schema:
            type: integer
        - in: query
          name: createdBefore
          schema:
            type: string
        - in: query
          name: createdAfter
          schema:
            type: string
        - in: query
          name: updatedBefore
          schema:
            type: string
        - in: query
          name: updatedAfter
          schema:
            type: string
        - in: query
          name: includeArchived
          schema:
            type: string
            enum: [ '1', 't', 'T', 'true', 'TRUE', 'True', '0', 'f', 'F', 'false', 'FALSE', 'False' ]
        - in: query
          name: sortBy
          schema:
            type: string
            enum: [ 'asc', 'desc' ]
        {{ end -}}
        {{ range $i, $x := .PathArguments }}
        - in: path
          type: string
          name: {{ $x }}
        {{ end }}
      {{- end }}
      responses:
        200:
          description: successful response
          content:
            application/json:
              schema:
                allOf:
                  - $ref: '#/components/schemas/APIResponse'
                  - type: object
                    properties:
                      data:
                      {{- if .ListRoute }}
                        type: array
                        items:
                          $ref: '#/components/schemas/{{ .ResponseType }}'
                      {{- else }}
                      data:
                        type: object
                        $ref: '#/components/schemas/{{ .ResponseType }}'
                      {{ end }}                     	
        400:
          description: bad request
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/APIResponseWithError'
        401:
          description: unauthorized
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/APIResponseWithError'
        500:
          description: internal error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/APIResponseWithError'`)
}

var routeParamRegex = regexp.MustCompile("\\{[a-zA-Z]+\\}")

var tagReplacements = map[string]string{
	"steps":                 "recipe_steps",
	"prep_tasks":            "recipe_prep_tasks",
	"completion_conditions": "recipe_step_completion_conditions",
}

func getTypeName(input interface{}) string {
	t := reflect.TypeOf(input)

	// Dereference pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

var responseTypeMap = map[string]string{
	"GET /_meta_/live":                                                                               "",
	"GET /_meta_/ready":                                                                              "",
	"POST /api/v1/admin/cycle_cookie_secret":                                                         "",
	"POST /api/v1/admin/users/status":                                                                getTypeName(&types.UserStatusResponse{}),
	"GET /api/v1/audit_log_entries/for_household":                                                    getTypeName(&types.AuditLogEntry{}),
	"GET /api/v1/audit_log_entries/for_user":                                                         getTypeName(&types.AuditLogEntry{}),
	"GET /api/v1/audit_log_entries/{auditLogEntryID}":                                                getTypeName(&types.AuditLogEntry{}),
	"GET /api/v1/household_invitations/received":                                                     getTypeName(&types.HouseholdInvitation{}),
	"GET /api/v1/household_invitations/sent":                                                         getTypeName(&types.HouseholdInvitation{}),
	"GET /api/v1/household_invitations/{householdInvitationID}/":                                     getTypeName(&types.HouseholdInvitation{}),
	"PUT /api/v1/household_invitations/{householdInvitationID}/accept":                               getTypeName(&types.HouseholdInvitation{}),
	"PUT /api/v1/household_invitations/{householdInvitationID}/cancel":                               getTypeName(&types.HouseholdInvitation{}),
	"PUT /api/v1/household_invitations/{householdInvitationID}/reject":                               getTypeName(&types.HouseholdInvitation{}),
	"GET /api/v1/households/":                                                                        getTypeName(&types.Household{}),
	"POST /api/v1/households/":                                                                       getTypeName(&types.Household{}),
	"GET /api/v1/households/current":                                                                 getTypeName(&types.Household{}),
	"POST /api/v1/households/instruments/":                                                           getTypeName(&types.HouseholdInstrumentOwnership{}),
	"GET /api/v1/households/instruments/":                                                            getTypeName(&types.HouseholdInstrumentOwnership{}),
	"GET /api/v1/households/instruments/{householdInstrumentOwnershipID}/":                           getTypeName(&types.HouseholdInstrumentOwnership{}),
	"PUT /api/v1/households/instruments/{householdInstrumentOwnershipID}/":                           getTypeName(&types.HouseholdInstrumentOwnership{}),
	"DELETE /api/v1/households/instruments/{householdInstrumentOwnershipID}/":                        getTypeName(&types.HouseholdInstrumentOwnership{}),
	"PUT /api/v1/households/{householdID}/":                                                          getTypeName(&types.Household{}),
	"DELETE /api/v1/households/{householdID}/":                                                       getTypeName(&types.Household{}),
	"GET /api/v1/households/{householdID}/":                                                          getTypeName(&types.Household{}),
	"POST /api/v1/households/{householdID}/default":                                                  getTypeName(&types.Household{}),
	"POST /api/v1/households/{householdID}/invitations/":                                             getTypeName(&types.HouseholdInvitation{}),
	"GET /api/v1/households/{householdID}/invitations/{householdInvitationID}/":                      getTypeName(&types.HouseholdInvitation{}),
	"POST /api/v1/households/{householdID}/invite":                                                   getTypeName(&types.HouseholdInvitation{}),
	"DELETE /api/v1/households/{householdID}/members/{userID}":                                       getTypeName(&types.HouseholdUserMembership{}),
	"PATCH /api/v1/households/{householdID}/members/{userID}/permissions":                            getTypeName(&types.UserPermissionsResponse{}),
	"POST /api/v1/households/{householdID}/transfer":                                                 getTypeName(&types.Household{}),
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/": getTypeName(&types.MealPlanOptionVote{}),
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/":    getTypeName(&types.MealPlanOptionVote{}),
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/": getTypeName(&types.MealPlanOptionVote{}),
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/":    getTypeName(&types.MealPlanOptionVote{}),
	"POST /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/":                                                   getTypeName(&types.MealPlanOption{}),
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/":                                                    getTypeName(&types.MealPlanOption{}),
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/":                                 getTypeName(&types.MealPlanOption{}),
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/":                                 getTypeName(&types.MealPlanOption{}),
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/":                              getTypeName(&types.MealPlanOption{}),
	"POST /api/v1/meal_plans/{mealPlanID}/events/":                                                                             getTypeName(&types.MealPlanEvent{}),
	"GET /api/v1/meal_plans/{mealPlanID}/events/":                                                                              getTypeName(&types.MealPlanEvent{}),
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/":                                                            getTypeName(&types.MealPlanEvent{}),
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/":                                                         getTypeName(&types.MealPlanEvent{}),
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/":                                                            getTypeName(&types.MealPlanEvent{}),
	"POST /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/vote":                                                       getTypeName(&types.MealPlanOptionVote{}),
	"GET /api/v1/meal_plans/{mealPlanID}/grocery_list_items/":                                                                  getTypeName(&types.MealPlanGroceryListItem{}),
	"POST /api/v1/meal_plans/{mealPlanID}/grocery_list_items/":                                                                 getTypeName(&types.MealPlanGroceryListItem{}),
	"GET /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/":                                      getTypeName(&types.MealPlanGroceryListItem{}),
	"PUT /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/":                                      getTypeName(&types.MealPlanGroceryListItem{}),
	"DELETE /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/":                                   getTypeName(&types.MealPlanGroceryListItem{}),
	"GET /api/v1/meal_plans/{mealPlanID}/tasks/":                                                                               getTypeName(&types.MealPlanTask{}),
	"POST /api/v1/meal_plans/{mealPlanID}/tasks/":                                                                              getTypeName(&types.MealPlanTask{}),
	"GET /api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}/":                                                              getTypeName(&types.MealPlanTask{}),
	"PATCH /api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}/":                                                            getTypeName(&types.MealPlanTask{}),
	"POST /api/v1/meal_plans/":                                                    getTypeName(&types.MealPlan{}),
	"GET /api/v1/meal_plans/":                                                     getTypeName(&types.MealPlan{}),
	"GET /api/v1/meal_plans/{mealPlanID}/":                                        getTypeName(&types.MealPlan{}),
	"PUT /api/v1/meal_plans/{mealPlanID}/":                                        getTypeName(&types.MealPlan{}),
	"DELETE /api/v1/meal_plans/{mealPlanID}/":                                     getTypeName(&types.MealPlan{}),
	"POST /api/v1/meal_plans/{mealPlanID}/finalize":                               getTypeName(&types.MealPlan{}),
	"POST /api/v1/meals/":                                                         getTypeName(&types.Meal{}),
	"GET /api/v1/meals/":                                                          getTypeName(&types.Meal{}),
	"GET /api/v1/meals/search":                                                    getTypeName(&types.Meal{}),
	"DELETE /api/v1/meals/{mealID}/":                                              getTypeName(&types.Meal{}),
	"GET /api/v1/meals/{mealID}/":                                                 getTypeName(&types.Meal{}),
	"GET /api/v1/oauth2_clients/":                                                 getTypeName(&types.OAuth2Client{}),
	"POST /api/v1/oauth2_clients/":                                                getTypeName(&types.OAuth2Client{}),
	"GET /api/v1/oauth2_clients/{oauth2ClientID}/":                                getTypeName(&types.OAuth2Client{}),
	"DELETE /api/v1/oauth2_clients/{oauth2ClientID}/":                             getTypeName(&types.OAuth2Client{}),
	"POST /api/v1/recipes/{recipeID}/prep_tasks/":                                 getTypeName(&types.RecipePrepTask{}),
	"GET /api/v1/recipes/{recipeID}/prep_tasks/":                                  getTypeName(&types.RecipePrepTask{}),
	"DELETE /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/":            getTypeName(&types.RecipePrepTask{}),
	"GET /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/":               getTypeName(&types.RecipePrepTask{}),
	"PUT /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/":               getTypeName(&types.RecipePrepTask{}),
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/":  getTypeName(&types.RecipeStepCompletionCondition{}),
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/": getTypeName(&types.RecipeStepCompletionCondition{}),
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/":    getTypeName(&types.RecipeStepCompletionCondition{}),
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/": getTypeName(&types.RecipeStepCompletionCondition{}),
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/":    getTypeName(&types.RecipeStepCompletionCondition{}),
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/":                                               getTypeName(&types.RecipeStepIngredient{}),
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/":                                                getTypeName(&types.RecipeStepIngredient{}),
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/":                       getTypeName(&types.RecipeStepIngredient{}),
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/":                       getTypeName(&types.RecipeStepIngredient{}),
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/":                    getTypeName(&types.RecipeStepIngredient{}),
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/":                                               getTypeName(&types.RecipeStepInstrument{}),
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/":                                                getTypeName(&types.RecipeStepInstrument{}),
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/":                       getTypeName(&types.RecipeStepInstrument{}),
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/":                       getTypeName(&types.RecipeStepInstrument{}),
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/":                    getTypeName(&types.RecipeStepInstrument{}),
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/":                                                   getTypeName(&types.RecipeStepProduct{}),
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/":                                                  getTypeName(&types.RecipeStepProduct{}),
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/":                             getTypeName(&types.RecipeStepProduct{}),
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/":                             getTypeName(&types.RecipeStepProduct{}),
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/":                          getTypeName(&types.RecipeStepProduct{}),
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/":                                                   getTypeName(&types.RecipeStepVessel{}),
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/":                                                    getTypeName(&types.RecipeStepVessel{}),
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/":                               getTypeName(&types.RecipeStepVessel{}),
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/":                               getTypeName(&types.RecipeStepVessel{}),
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/":                            getTypeName(&types.RecipeStepVessel{}),
	"POST /api/v1/recipes/{recipeID}/steps/":                                                                          getTypeName(&types.RecipeStep{}),
	"GET /api/v1/recipes/{recipeID}/steps/":                                                                           getTypeName(&types.RecipeStep{}),
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/":                                                         getTypeName(&types.RecipeStep{}),
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/":                                                            getTypeName(&types.RecipeStep{}),
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/":                                                            getTypeName(&types.RecipeStep{}),
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/images":                                                     "",
	"POST /api/v1/recipes/":                                                                        getTypeName(&types.Recipe{}),
	"GET /api/v1/recipes/":                                                                         getTypeName(&types.Recipe{}),
	"GET /api/v1/recipes/search":                                                                   getTypeName(&types.Recipe{}),
	"DELETE /api/v1/recipes/{recipeID}/":                                                           getTypeName(&types.Recipe{}),
	"GET /api/v1/recipes/{recipeID}/":                                                              getTypeName(&types.Recipe{}),
	"PUT /api/v1/recipes/{recipeID}/":                                                              getTypeName(&types.Recipe{}),
	"POST /api/v1/recipes/{recipeID}/clone":                                                        getTypeName(&types.Recipe{}),
	"GET /api/v1/recipes/{recipeID}/dag":                                                           getTypeName(&types.APIError{}),
	"POST /api/v1/recipes/{recipeID}/images":                                                       "",
	"GET /api/v1/recipes/{recipeID}/mermaid":                                                       "",
	"GET /api/v1/recipes/{recipeID}/prep_steps":                                                    getTypeName(&types.RecipePrepTaskStep{}),
	"POST /api/v1/recipes/{recipeID}/ratings/":                                                     getTypeName(&types.RecipeRating{}),
	"GET /api/v1/recipes/{recipeID}/ratings/":                                                      getTypeName(&types.RecipeRating{}),
	"PUT /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/":                                     getTypeName(&types.RecipeRating{}),
	"DELETE /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/":                                  getTypeName(&types.RecipeRating{}),
	"GET /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/":                                     getTypeName(&types.RecipeRating{}),
	"POST /api/v1/settings/":                                                                       getTypeName(&types.ServiceSetting{}),
	"GET /api/v1/settings/":                                                                        getTypeName(&types.ServiceSetting{}),
	"POST /api/v1/settings/configurations/":                                                        getTypeName(&types.ServiceSettingConfiguration{}),
	"GET /api/v1/settings/configurations/household":                                                getTypeName(&types.ServiceSettingConfiguration{}),
	"GET /api/v1/settings/configurations/user":                                                     getTypeName(&types.ServiceSettingConfiguration{}),
	"GET /api/v1/settings/configurations/user/{serviceSettingConfigurationName}":                   getTypeName(&types.ServiceSettingConfiguration{}),
	"DELETE /api/v1/settings/configurations/{serviceSettingConfigurationID}":                       getTypeName(&types.ServiceSettingConfiguration{}),
	"PUT /api/v1/settings/configurations/{serviceSettingConfigurationID}":                          getTypeName(&types.ServiceSettingConfiguration{}),
	"GET /api/v1/settings/search":                                                                  getTypeName(&types.ServiceSetting{}),
	"GET /api/v1/settings/{serviceSettingID}/":                                                     getTypeName(&types.ServiceSetting{}),
	"DELETE /api/v1/settings/{serviceSettingID}/":                                                  getTypeName(&types.ServiceSetting{}),
	"POST /api/v1/user_ingredient_preferences/":                                                    getTypeName(&types.UserIngredientPreference{}),
	"GET /api/v1/user_ingredient_preferences/":                                                     getTypeName(&types.UserIngredientPreference{}),
	"PUT /api/v1/user_ingredient_preferences/{userIngredientPreferenceID}/":                        getTypeName(&types.UserIngredientPreference{}),
	"DELETE /api/v1/user_ingredient_preferences/{userIngredientPreferenceID}/":                     getTypeName(&types.UserIngredientPreference{}),
	"POST /api/v1/user_notifications/":                                                             getTypeName(&types.UserNotification{}),
	"GET /api/v1/user_notifications/":                                                              getTypeName(&types.UserNotification{}),
	"GET /api/v1/user_notifications/{userNotificationID}/":                                         getTypeName(&types.UserNotification{}),
	"PATCH /api/v1/user_notifications/{userNotificationID}/":                                       getTypeName(&types.UserNotification{}),
	"GET /api/v1/users/":                                                                           getTypeName(&types.User{}),
	"POST /api/v1/users/avatar/upload":                                                             "",
	"PUT /api/v1/users/details":                                                                    getTypeName(&types.User{}),
	"PUT /api/v1/users/email_address":                                                              getTypeName(&types.User{}),
	"POST /api/v1/users/email_address_verification":                                                getTypeName(&types.User{}),
	"POST /api/v1/users/household/select":                                                          getTypeName(&types.Household{}),
	"PUT /api/v1/users/password/new":                                                               getTypeName(&types.User{}),
	"POST /api/v1/users/permissions/check":                                                         "",
	"GET /api/v1/users/search":                                                                     getTypeName(&types.User{}),
	"GET /api/v1/users/self":                                                                       getTypeName(&types.User{}),
	"POST /api/v1/users/totp_secret/new":                                                           getTypeName(&types.APIError{}),
	"PUT /api/v1/users/username":                                                                   getTypeName(&types.User{}),
	"GET /api/v1/users/{userID}/":                                                                  getTypeName(&types.User{}),
	"DELETE /api/v1/users/{userID}/":                                                               getTypeName(&types.User{}),
	"POST /api/v1/valid_ingredient_groups/":                                                        getTypeName(&types.ValidIngredientGroup{}),
	"GET /api/v1/valid_ingredient_groups/":                                                         getTypeName(&types.ValidIngredientGroup{}),
	"GET /api/v1/valid_ingredient_groups/search":                                                   getTypeName(&types.ValidIngredientGroup{}),
	"DELETE /api/v1/valid_ingredient_groups/{validIngredientGroupID}/":                             getTypeName(&types.ValidIngredientGroup{}),
	"GET /api/v1/valid_ingredient_groups/{validIngredientGroupID}/":                                getTypeName(&types.ValidIngredientGroup{}),
	"PUT /api/v1/valid_ingredient_groups/{validIngredientGroupID}/":                                getTypeName(&types.ValidIngredientGroup{}),
	"POST /api/v1/valid_ingredient_measurement_units/":                                             getTypeName(&types.ValidIngredientMeasurementUnit{}),
	"GET /api/v1/valid_ingredient_measurement_units/":                                              getTypeName(&types.ValidIngredientMeasurementUnit{}),
	"GET /api/v1/valid_ingredient_measurement_units/by_ingredient/{validIngredientID}/":            getTypeName(&types.ValidIngredientMeasurementUnit{}),
	"GET /api/v1/valid_ingredient_measurement_units/by_measurement_unit/{validMeasurementUnitID}/": getTypeName(&types.ValidIngredientMeasurementUnit{}),
	"GET /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/":           getTypeName(&types.ValidIngredientMeasurementUnit{}),
	"PUT /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/":           getTypeName(&types.ValidIngredientMeasurementUnit{}),
	"DELETE /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/":        getTypeName(&types.ValidIngredientMeasurementUnit{}),
	"GET /api/v1/valid_ingredient_preparations/":                                                   getTypeName(&types.ValidIngredientPreparation{}),
	"POST /api/v1/valid_ingredient_preparations/":                                                  getTypeName(&types.ValidIngredientPreparation{}),
	"GET /api/v1/valid_ingredient_preparations/by_ingredient/{validIngredientID}/":                 getTypeName(&types.ValidIngredientPreparation{}),
	"GET /api/v1/valid_ingredient_preparations/by_preparation/{validPreparationID}/":               getTypeName(&types.ValidIngredientPreparation{}),
	"GET /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/":                    getTypeName(&types.ValidIngredientPreparation{}),
	"PUT /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/":                    getTypeName(&types.ValidIngredientPreparation{}),
	"DELETE /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/":                 getTypeName(&types.ValidIngredientPreparation{}),
	"POST /api/v1/valid_ingredient_state_ingredients/":                                             getTypeName(&types.ValidIngredientStateIngredient{}),
	"GET /api/v1/valid_ingredient_state_ingredients/":                                              getTypeName(&types.ValidIngredientStateIngredient{}),
	"GET /api/v1/valid_ingredient_state_ingredients/by_ingredient/{validIngredientID}/":            getTypeName(&types.ValidIngredientStateIngredient{}),
	"GET /api/v1/valid_ingredient_state_ingredients/by_ingredient_state/{validIngredientStateID}/": getTypeName(&types.ValidIngredientStateIngredient{}),
	"GET /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/":           getTypeName(&types.ValidIngredientStateIngredient{}),
	"PUT /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/":           getTypeName(&types.ValidIngredientStateIngredient{}),
	"DELETE /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/":        getTypeName(&types.ValidIngredientStateIngredient{}),
	"POST /api/v1/valid_ingredient_states/":                                                        getTypeName(&types.ValidIngredientState{}),
	"GET /api/v1/valid_ingredient_states/":                                                         getTypeName(&types.ValidIngredientState{}),
	"GET /api/v1/valid_ingredient_states/search":                                                   getTypeName(&types.ValidIngredientState{}),
	"PUT /api/v1/valid_ingredient_states/{validIngredientStateID}/":                                getTypeName(&types.ValidIngredientState{}),
	"DELETE /api/v1/valid_ingredient_states/{validIngredientStateID}/":                             getTypeName(&types.ValidIngredientState{}),
	"GET /api/v1/valid_ingredient_states/{validIngredientStateID}/":                                getTypeName(&types.ValidIngredientState{}),
	"POST /api/v1/valid_ingredients/":                                                              getTypeName(&types.ValidIngredient{}),
	"GET /api/v1/valid_ingredients/":                                                               getTypeName(&types.ValidIngredient{}),
	"GET /api/v1/valid_ingredients/by_preparation/{validPreparationID}/":                           getTypeName(&types.ValidIngredient{}),
	"GET /api/v1/valid_ingredients/random":                                                         getTypeName(&types.ValidIngredient{}),
	"GET /api/v1/valid_ingredients/search":                                                         getTypeName(&types.ValidIngredient{}),
	"PUT /api/v1/valid_ingredients/{validIngredientID}/":                                           getTypeName(&types.ValidIngredient{}),
	"DELETE /api/v1/valid_ingredients/{validIngredientID}/":                                        getTypeName(&types.ValidIngredient{}),
	"GET /api/v1/valid_ingredients/{validIngredientID}/":                                           getTypeName(&types.ValidIngredient{}),
	"GET /api/v1/valid_instruments/":                                                               getTypeName(&types.ValidInstrument{}),
	"POST /api/v1/valid_instruments/":                                                              getTypeName(&types.ValidInstrument{}),
	"GET /api/v1/valid_instruments/random":                                                         getTypeName(&types.ValidInstrument{}),
	"GET /api/v1/valid_instruments/search":                                                         getTypeName(&types.ValidInstrument{}),
	"DELETE /api/v1/valid_instruments/{validInstrumentID}/":                                        getTypeName(&types.ValidInstrument{}),
	"GET /api/v1/valid_instruments/{validInstrumentID}/":                                           getTypeName(&types.ValidInstrument{}),
	"PUT /api/v1/valid_instruments/{validInstrumentID}/":                                           getTypeName(&types.ValidInstrument{}),
	"POST /api/v1/valid_measurement_conversions/":                                                  getTypeName(&types.ValidMeasurementUnitConversion{}),
	"GET /api/v1/valid_measurement_conversions/from_unit/{validMeasurementUnitID}":                 getTypeName(&types.ValidMeasurementUnitConversion{}),
	"GET /api/v1/valid_measurement_conversions/to_unit/{validMeasurementUnitID}":                   getTypeName(&types.ValidMeasurementUnitConversion{}),
	"PUT /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/":                getTypeName(&types.ValidMeasurementUnitConversion{}),
	"DELETE /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/":             getTypeName(&types.ValidMeasurementUnitConversion{}),
	"GET /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/":                getTypeName(&types.ValidMeasurementUnitConversion{}),
	"POST /api/v1/valid_measurement_units/":                                                        getTypeName(&types.ValidMeasurementUnit{}),
	"GET /api/v1/valid_measurement_units/":                                                         getTypeName(&types.ValidMeasurementUnit{}),
	"GET /api/v1/valid_measurement_units/by_ingredient/{validIngredientID}":                        getTypeName(&types.ValidMeasurementUnit{}),
	"GET /api/v1/valid_measurement_units/search":                                                   getTypeName(&types.ValidMeasurementUnit{}),
	"GET /api/v1/valid_measurement_units/{validMeasurementUnitID}/":                                getTypeName(&types.ValidMeasurementUnit{}),
	"PUT /api/v1/valid_measurement_units/{validMeasurementUnitID}/":                                getTypeName(&types.ValidMeasurementUnit{}),
	"DELETE /api/v1/valid_measurement_units/{validMeasurementUnitID}/":                             getTypeName(&types.ValidMeasurementUnit{}),
	"GET /api/v1/valid_preparation_instruments/":                                                   getTypeName(&types.ValidPreparationInstrument{}),
	"POST /api/v1/valid_preparation_instruments/":                                                  getTypeName(&types.ValidPreparationInstrument{}),
	"GET /api/v1/valid_preparation_instruments/by_instrument/{validInstrumentID}/":                 getTypeName(&types.ValidPreparationInstrument{}),
	"GET /api/v1/valid_preparation_instruments/by_preparation/{validPreparationID}/":               getTypeName(&types.ValidPreparationInstrument{}),
	"DELETE /api/v1/valid_preparation_instruments/{validPreparationVesselID}/":                     getTypeName(&types.ValidPreparationInstrument{}),
	"GET /api/v1/valid_preparation_instruments/{validPreparationVesselID}/":                        getTypeName(&types.ValidPreparationInstrument{}),
	"PUT /api/v1/valid_preparation_instruments/{validPreparationVesselID}/":                        getTypeName(&types.ValidPreparationInstrument{}),
	"POST /api/v1/valid_preparation_vessels/":                                                      getTypeName(&types.ValidPreparationVessel{}),
	"GET /api/v1/valid_preparation_vessels/":                                                       getTypeName(&types.ValidPreparationVessel{}),
	"GET /api/v1/valid_preparation_vessels/by_preparation/{validPreparationID}/":                   getTypeName(&types.ValidPreparationVessel{}),
	"GET /api/v1/valid_preparation_vessels/by_vessel/{ValidVesselID}/":                             getTypeName(&types.ValidPreparationVessel{}),
	"PUT /api/v1/valid_preparation_vessels/{validPreparationVesselID}/":                            getTypeName(&types.ValidPreparationVessel{}),
	"DELETE /api/v1/valid_preparation_vessels/{validPreparationVesselID}/":                         getTypeName(&types.ValidPreparationVessel{}),
	"GET /api/v1/valid_preparation_vessels/{validPreparationVesselID}/":                            getTypeName(&types.ValidPreparationVessel{}),
	"GET /api/v1/valid_preparations/":                                                              getTypeName(&types.ValidPreparation{}),
	"POST /api/v1/valid_preparations/":                                                             getTypeName(&types.ValidPreparation{}),
	"GET /api/v1/valid_preparations/random":                                                        getTypeName(&types.ValidPreparation{}),
	"GET /api/v1/valid_preparations/search":                                                        getTypeName(&types.ValidPreparation{}),
	"PUT /api/v1/valid_preparations/{validPreparationID}/":                                         getTypeName(&types.ValidPreparation{}),
	"DELETE /api/v1/valid_preparations/{validPreparationID}/":                                      getTypeName(&types.ValidPreparation{}),
	"GET /api/v1/valid_preparations/{validPreparationID}/":                                         getTypeName(&types.ValidPreparation{}),
	"POST /api/v1/valid_vessels/":                                                                  getTypeName(&types.ValidVessel{}),
	"GET /api/v1/valid_vessels/":                                                                   getTypeName(&types.ValidVessel{}),
	"GET /api/v1/valid_vessels/random":                                                             getTypeName(&types.ValidVessel{}),
	"GET /api/v1/valid_vessels/search":                                                             getTypeName(&types.ValidVessel{}),
	"GET /api/v1/valid_vessels/{validVesselID}/":                                                   getTypeName(&types.ValidVessel{}),
	"PUT /api/v1/valid_vessels/{validVesselID}/":                                                   getTypeName(&types.ValidVessel{}),
	"DELETE /api/v1/valid_vessels/{validVesselID}/":                                                getTypeName(&types.ValidVessel{}),
	"GET /api/v1/webhooks/":                                                                        getTypeName(&types.Webhook{}),
	"POST /api/v1/webhooks/":                                                                       getTypeName(&types.Webhook{}),
	"GET /api/v1/webhooks/{webhookID}/":                                                            getTypeName(&types.Webhook{}),
	"DELETE /api/v1/webhooks/{webhookID}/":                                                         getTypeName(&types.Webhook{}),
	"POST /api/v1/webhooks/{webhookID}/trigger_events":                                             getTypeName(&types.WebhookTriggerEvent{}),
	"DELETE /api/v1/webhooks/{webhookID}/trigger_events/{webhookTriggerEventID}/":                  getTypeName(&types.WebhookTriggerEvent{}),
	"POST /api/v1/workers/finalize_meal_plans":                                                     "",
	"POST /api/v1/workers/meal_plan_grocery_list_init":                                             "",
	"POST /api/v1/workers/meal_plan_tasks":                                                         "",
	"GET /auth/status":                                                                             getTypeName(&types.UserStatusResponse{}),
	"GET /auth/{auth_provider}":                                                                    "",
	"GET /auth/{auth_provider}/callback":                                                           "",
	"GET /oauth2/authorize":                                                                        "",
	"POST /oauth2/token":                                                                           "",
	"POST /users/":                                                                                 getTypeName(&types.User{}),
	"POST /users/email_address/verify":                                                             "",
	"POST /users/login":                                                                            "",
	"POST /users/login/admin":                                                                      "",
	"POST /users/logout":                                                                           "",
	"POST /users/password/reset":                                                                   getTypeName(&types.PasswordResetToken{}),
	"POST /users/password/reset/redeem":                                                            "",
	"POST /users/totp_secret/verify":                                                               "",
	"POST /users/username/reminder":                                                                "",
}

func main() {
	ctx := context.Background()

	rawCfg, err := os.ReadFile("environments/dev/config_files/service-config.json")
	if err != nil {
		log.Fatal(err)
	}

	var cfg *config.InstanceConfig
	if err = json.Unmarshal(rawCfg, &cfg); err != nil {
		log.Fatal(err)
	}

	cfg.Neutralize()

	// build our server struct.
	srv, err := build.Build(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	routeDefinitions := []RouteDefinition{}
	for _, route := range srv.Router().Routes() {
		pathArgs := []string{}
		for _, pathArg := range routeParamRegex.FindAllString(route.Path, -1) {
			pathArgs = append(pathArgs, strings.TrimPrefix(strings.TrimSuffix(pathArg, "}"), "{"))
		}

		routeDef := RouteDefinition{
			Method:        route.Method,
			Path:          route.Path,
			PathArguments: pathArgs,
			ListRoute:     strings.HasSuffix(route.Path, "s/") && route.Method == http.MethodGet,
			ResponseType:  responseTypeMap[fmt.Sprintf("%s %s", route.Method, route.Path)],
		}

		for _, part := range strings.Split(route.Path, "/") {
			if strings.TrimSpace(part) != "" && !strings.HasPrefix(part, "{") && part != "api" && part != "v1" {
				if part == "steps" {
					print("")
				}

				if rep, ok := tagReplacements[part]; ok {
					routeDef.Tags = append(routeDef.Tags, rep)
				} else {
					routeDef.Tags = append(routeDef.Tags, part)
				}
			}
		}

		routeDefinitions = append(routeDefinitions, routeDef)
	}

	for _, rd := range routeDefinitions {
		fmt.Printf("%q \n", rd.ResponseType)
	}
}
