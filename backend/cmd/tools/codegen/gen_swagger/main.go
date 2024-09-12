package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/server/http/build"
	"github.com/dinnerdonebetter/backend/pkg/types"
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

type routeDetails struct {
	ResponseTypeName string
	InputTypeName    string
}

var routeInfoMap = map[string]routeDetails{
	"GET /_meta_/live":                                                                               {ResponseTypeName: "", InputTypeName: ""},
	"GET /_meta_/ready":                                                                              {ResponseTypeName: "", InputTypeName: ""},
	"POST /api/v1/admin/cycle_cookie_secret":                                                         {ResponseTypeName: "", InputTypeName: ""},
	"POST /api/v1/admin/users/status":                                                                {ResponseTypeName: getTypeName(&types.UserStatusResponse{}), InputTypeName: ""},
	"GET /api/v1/audit_log_entries/for_household":                                                    {ResponseTypeName: getTypeName(&types.AuditLogEntry{}), InputTypeName: ""},
	"GET /api/v1/audit_log_entries/for_user":                                                         {ResponseTypeName: getTypeName(&types.AuditLogEntry{}), InputTypeName: ""},
	"GET /api/v1/audit_log_entries/{auditLogEntryID}":                                                {ResponseTypeName: getTypeName(&types.AuditLogEntry{}), InputTypeName: ""},
	"GET /api/v1/household_invitations/received":                                                     {ResponseTypeName: getTypeName(&types.HouseholdInvitation{}), InputTypeName: ""},
	"GET /api/v1/household_invitations/sent":                                                         {ResponseTypeName: getTypeName(&types.HouseholdInvitation{}), InputTypeName: ""},
	"GET /api/v1/household_invitations/{householdInvitationID}/":                                     {ResponseTypeName: getTypeName(&types.HouseholdInvitation{}), InputTypeName: ""},
	"PUT /api/v1/household_invitations/{householdInvitationID}/accept":                               {ResponseTypeName: getTypeName(&types.HouseholdInvitation{}), InputTypeName: ""},
	"PUT /api/v1/household_invitations/{householdInvitationID}/cancel":                               {ResponseTypeName: getTypeName(&types.HouseholdInvitation{}), InputTypeName: ""},
	"PUT /api/v1/household_invitations/{householdInvitationID}/reject":                               {ResponseTypeName: getTypeName(&types.HouseholdInvitation{}), InputTypeName: ""},
	"GET /api/v1/households/":                                                                        {ResponseTypeName: getTypeName(&types.Household{}), InputTypeName: ""},
	"POST /api/v1/households/":                                                                       {ResponseTypeName: getTypeName(&types.Household{}), InputTypeName: ""},
	"GET /api/v1/households/current":                                                                 {ResponseTypeName: getTypeName(&types.Household{}), InputTypeName: ""},
	"POST /api/v1/households/instruments/":                                                           {ResponseTypeName: getTypeName(&types.HouseholdInstrumentOwnership{}), InputTypeName: ""},
	"GET /api/v1/households/instruments/":                                                            {ResponseTypeName: getTypeName(&types.HouseholdInstrumentOwnership{}), InputTypeName: ""},
	"GET /api/v1/households/instruments/{householdInstrumentOwnershipID}/":                           {ResponseTypeName: getTypeName(&types.HouseholdInstrumentOwnership{}), InputTypeName: ""},
	"PUT /api/v1/households/instruments/{householdInstrumentOwnershipID}/":                           {ResponseTypeName: getTypeName(&types.HouseholdInstrumentOwnership{}), InputTypeName: ""},
	"DELETE /api/v1/households/instruments/{householdInstrumentOwnershipID}/":                        {ResponseTypeName: getTypeName(&types.HouseholdInstrumentOwnership{}), InputTypeName: ""},
	"PUT /api/v1/households/{householdID}/":                                                          {ResponseTypeName: getTypeName(&types.Household{}), InputTypeName: ""},
	"DELETE /api/v1/households/{householdID}/":                                                       {ResponseTypeName: getTypeName(&types.Household{}), InputTypeName: ""},
	"GET /api/v1/households/{householdID}/":                                                          {ResponseTypeName: getTypeName(&types.Household{}), InputTypeName: ""},
	"POST /api/v1/households/{householdID}/default":                                                  {ResponseTypeName: getTypeName(&types.Household{}), InputTypeName: ""},
	"POST /api/v1/households/{householdID}/invitations/":                                             {ResponseTypeName: getTypeName(&types.HouseholdInvitation{}), InputTypeName: ""},
	"GET /api/v1/households/{householdID}/invitations/{householdInvitationID}/":                      {ResponseTypeName: getTypeName(&types.HouseholdInvitation{}), InputTypeName: ""},
	"POST /api/v1/households/{householdID}/invite":                                                   {ResponseTypeName: getTypeName(&types.HouseholdInvitation{}), InputTypeName: ""},
	"DELETE /api/v1/households/{householdID}/members/{userID}":                                       {ResponseTypeName: getTypeName(&types.HouseholdUserMembership{}), InputTypeName: ""},
	"PATCH /api/v1/households/{householdID}/members/{userID}/permissions":                            {ResponseTypeName: getTypeName(&types.UserPermissionsResponse{}), InputTypeName: ""},
	"POST /api/v1/households/{householdID}/transfer":                                                 {ResponseTypeName: getTypeName(&types.Household{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/": {ResponseTypeName: getTypeName(&types.MealPlanOptionVote{}), InputTypeName: ""},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/":    {ResponseTypeName: getTypeName(&types.MealPlanOptionVote{}), InputTypeName: ""},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/": {ResponseTypeName: getTypeName(&types.MealPlanOptionVote{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/votes/{mealPlanOptionVoteID}/":    {ResponseTypeName: getTypeName(&types.MealPlanOptionVote{}), InputTypeName: ""},
	"POST /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/":                                                   {ResponseTypeName: getTypeName(&types.MealPlanOption{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/":                                                    {ResponseTypeName: getTypeName(&types.MealPlanOption{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/":                                 {ResponseTypeName: getTypeName(&types.MealPlanOption{}), InputTypeName: ""},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/":                                 {ResponseTypeName: getTypeName(&types.MealPlanOption{}), InputTypeName: ""},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/options/{mealPlanOptionID}/":                              {ResponseTypeName: getTypeName(&types.MealPlanOption{}), InputTypeName: ""},
	"POST /api/v1/meal_plans/{mealPlanID}/events/":                                                                             {ResponseTypeName: getTypeName(&types.MealPlanEvent{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/{mealPlanID}/events/":                                                                              {ResponseTypeName: getTypeName(&types.MealPlanEvent{}), InputTypeName: ""},
	"PUT /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/":                                                            {ResponseTypeName: getTypeName(&types.MealPlanEvent{}), InputTypeName: ""},
	"DELETE /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/":                                                         {ResponseTypeName: getTypeName(&types.MealPlanEvent{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/":                                                            {ResponseTypeName: getTypeName(&types.MealPlanEvent{}), InputTypeName: ""},
	"POST /api/v1/meal_plans/{mealPlanID}/events/{mealPlanEventID}/vote":                                                       {ResponseTypeName: getTypeName(&types.MealPlanOptionVote{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/{mealPlanID}/grocery_list_items/":                                                                  {ResponseTypeName: getTypeName(&types.MealPlanGroceryListItem{}), InputTypeName: ""},
	"POST /api/v1/meal_plans/{mealPlanID}/grocery_list_items/":                                                                 {ResponseTypeName: getTypeName(&types.MealPlanGroceryListItem{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/":                                      {ResponseTypeName: getTypeName(&types.MealPlanGroceryListItem{}), InputTypeName: ""},
	"PUT /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/":                                      {ResponseTypeName: getTypeName(&types.MealPlanGroceryListItem{}), InputTypeName: ""},
	"DELETE /api/v1/meal_plans/{mealPlanID}/grocery_list_items/{mealPlanGroceryListItemID}/":                                   {ResponseTypeName: getTypeName(&types.MealPlanGroceryListItem{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/{mealPlanID}/tasks/":                                                                               {ResponseTypeName: getTypeName(&types.MealPlanTask{}), InputTypeName: ""},
	"POST /api/v1/meal_plans/{mealPlanID}/tasks/":                                                                              {ResponseTypeName: getTypeName(&types.MealPlanTask{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}/":                                                              {ResponseTypeName: getTypeName(&types.MealPlanTask{}), InputTypeName: ""},
	"PATCH /api/v1/meal_plans/{mealPlanID}/tasks/{mealPlanTaskID}/":                                                            {ResponseTypeName: getTypeName(&types.MealPlanTask{}), InputTypeName: ""},
	"POST /api/v1/meal_plans/":                                                    {ResponseTypeName: getTypeName(&types.MealPlan{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/":                                                     {ResponseTypeName: getTypeName(&types.MealPlan{}), InputTypeName: ""},
	"GET /api/v1/meal_plans/{mealPlanID}/":                                        {ResponseTypeName: getTypeName(&types.MealPlan{}), InputTypeName: ""},
	"PUT /api/v1/meal_plans/{mealPlanID}/":                                        {ResponseTypeName: getTypeName(&types.MealPlan{}), InputTypeName: ""},
	"DELETE /api/v1/meal_plans/{mealPlanID}/":                                     {ResponseTypeName: getTypeName(&types.MealPlan{}), InputTypeName: ""},
	"POST /api/v1/meal_plans/{mealPlanID}/finalize":                               {ResponseTypeName: getTypeName(&types.MealPlan{}), InputTypeName: ""},
	"POST /api/v1/meals/":                                                         {ResponseTypeName: getTypeName(&types.Meal{}), InputTypeName: ""},
	"GET /api/v1/meals/":                                                          {ResponseTypeName: getTypeName(&types.Meal{}), InputTypeName: ""},
	"GET /api/v1/meals/search":                                                    {ResponseTypeName: getTypeName(&types.Meal{}), InputTypeName: ""},
	"DELETE /api/v1/meals/{mealID}/":                                              {ResponseTypeName: getTypeName(&types.Meal{}), InputTypeName: ""},
	"GET /api/v1/meals/{mealID}/":                                                 {ResponseTypeName: getTypeName(&types.Meal{}), InputTypeName: ""},
	"GET /api/v1/oauth2_clients/":                                                 {ResponseTypeName: getTypeName(&types.OAuth2Client{}), InputTypeName: ""},
	"POST /api/v1/oauth2_clients/":                                                {ResponseTypeName: getTypeName(&types.OAuth2Client{}), InputTypeName: ""},
	"GET /api/v1/oauth2_clients/{oauth2ClientID}/":                                {ResponseTypeName: getTypeName(&types.OAuth2Client{}), InputTypeName: ""},
	"DELETE /api/v1/oauth2_clients/{oauth2ClientID}/":                             {ResponseTypeName: getTypeName(&types.OAuth2Client{}), InputTypeName: ""},
	"POST /api/v1/recipes/{recipeID}/prep_tasks/":                                 {ResponseTypeName: getTypeName(&types.RecipePrepTask{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/prep_tasks/":                                  {ResponseTypeName: getTypeName(&types.RecipePrepTask{}), InputTypeName: ""},
	"DELETE /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/":            {ResponseTypeName: getTypeName(&types.RecipePrepTask{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/":               {ResponseTypeName: getTypeName(&types.RecipePrepTask{}), InputTypeName: ""},
	"PUT /api/v1/recipes/{recipeID}/prep_tasks/{recipePrepTaskID}/":               {ResponseTypeName: getTypeName(&types.RecipePrepTask{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/":  {ResponseTypeName: getTypeName(&types.RecipeStepCompletionCondition{}), InputTypeName: ""},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/": {ResponseTypeName: getTypeName(&types.RecipeStepCompletionCondition{}), InputTypeName: ""},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/":    {ResponseTypeName: getTypeName(&types.RecipeStepCompletionCondition{}), InputTypeName: ""},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/": {ResponseTypeName: getTypeName(&types.RecipeStepCompletionCondition{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/completion_conditions/{recipeStepCompletionConditionID}/":    {ResponseTypeName: getTypeName(&types.RecipeStepCompletionCondition{}), InputTypeName: ""},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/":                                               {ResponseTypeName: getTypeName(&types.RecipeStepIngredient{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/":                                                {ResponseTypeName: getTypeName(&types.RecipeStepIngredient{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/":                       {ResponseTypeName: getTypeName(&types.RecipeStepIngredient{}), InputTypeName: ""},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/":                       {ResponseTypeName: getTypeName(&types.RecipeStepIngredient{}), InputTypeName: ""},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/ingredients/{recipeStepIngredientID}/":                    {ResponseTypeName: getTypeName(&types.RecipeStepIngredient{}), InputTypeName: ""},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/":                                               {ResponseTypeName: getTypeName(&types.RecipeStepInstrument{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/":                                                {ResponseTypeName: getTypeName(&types.RecipeStepInstrument{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/":                       {ResponseTypeName: getTypeName(&types.RecipeStepInstrument{}), InputTypeName: ""},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/":                       {ResponseTypeName: getTypeName(&types.RecipeStepInstrument{}), InputTypeName: ""},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/instruments/{recipeStepInstrumentID}/":                    {ResponseTypeName: getTypeName(&types.RecipeStepInstrument{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/":                                                   {ResponseTypeName: getTypeName(&types.RecipeStepProduct{}), InputTypeName: ""},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/":                                                  {ResponseTypeName: getTypeName(&types.RecipeStepProduct{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/":                             {ResponseTypeName: getTypeName(&types.RecipeStepProduct{}), InputTypeName: ""},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/":                             {ResponseTypeName: getTypeName(&types.RecipeStepProduct{}), InputTypeName: ""},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/products/{recipeStepProductID}/":                          {ResponseTypeName: getTypeName(&types.RecipeStepProduct{}), InputTypeName: ""},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/":                                                   {ResponseTypeName: getTypeName(&types.RecipeStepVessel{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/":                                                    {ResponseTypeName: getTypeName(&types.RecipeStepVessel{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/":                               {ResponseTypeName: getTypeName(&types.RecipeStepVessel{}), InputTypeName: ""},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/":                               {ResponseTypeName: getTypeName(&types.RecipeStepVessel{}), InputTypeName: ""},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/vessels/{recipeStepVesselID}/":                            {ResponseTypeName: getTypeName(&types.RecipeStepVessel{}), InputTypeName: ""},
	"POST /api/v1/recipes/{recipeID}/steps/":                                                                          {ResponseTypeName: getTypeName(&types.RecipeStep{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/":                                                                           {ResponseTypeName: getTypeName(&types.RecipeStep{}), InputTypeName: ""},
	"DELETE /api/v1/recipes/{recipeID}/steps/{recipeStepID}/":                                                         {ResponseTypeName: getTypeName(&types.RecipeStep{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/steps/{recipeStepID}/":                                                            {ResponseTypeName: getTypeName(&types.RecipeStep{}), InputTypeName: ""},
	"PUT /api/v1/recipes/{recipeID}/steps/{recipeStepID}/":                                                            {ResponseTypeName: getTypeName(&types.RecipeStep{}), InputTypeName: ""},
	"POST /api/v1/recipes/{recipeID}/steps/{recipeStepID}/images":                                                     {ResponseTypeName: "", InputTypeName: ""},
	"POST /api/v1/recipes/":                                                                        {ResponseTypeName: getTypeName(&types.Recipe{}), InputTypeName: ""},
	"GET /api/v1/recipes/":                                                                         {ResponseTypeName: getTypeName(&types.Recipe{}), InputTypeName: ""},
	"GET /api/v1/recipes/search":                                                                   {ResponseTypeName: getTypeName(&types.Recipe{}), InputTypeName: ""},
	"DELETE /api/v1/recipes/{recipeID}/":                                                           {ResponseTypeName: getTypeName(&types.Recipe{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/":                                                              {ResponseTypeName: getTypeName(&types.Recipe{}), InputTypeName: ""},
	"PUT /api/v1/recipes/{recipeID}/":                                                              {ResponseTypeName: getTypeName(&types.Recipe{}), InputTypeName: ""},
	"POST /api/v1/recipes/{recipeID}/clone":                                                        {ResponseTypeName: getTypeName(&types.Recipe{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/dag":                                                           {ResponseTypeName: getTypeName(&types.APIError{}), InputTypeName: ""},
	"POST /api/v1/recipes/{recipeID}/images":                                                       {ResponseTypeName: "", InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/mermaid":                                                       {ResponseTypeName: "", InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/prep_steps":                                                    {ResponseTypeName: getTypeName(&types.RecipePrepTaskStep{}), InputTypeName: ""},
	"POST /api/v1/recipes/{recipeID}/ratings/":                                                     {ResponseTypeName: getTypeName(&types.RecipeRating{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/ratings/":                                                      {ResponseTypeName: getTypeName(&types.RecipeRating{}), InputTypeName: ""},
	"PUT /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/":                                     {ResponseTypeName: getTypeName(&types.RecipeRating{}), InputTypeName: ""},
	"DELETE /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/":                                  {ResponseTypeName: getTypeName(&types.RecipeRating{}), InputTypeName: ""},
	"GET /api/v1/recipes/{recipeID}/ratings/{recipeRatingID}/":                                     {ResponseTypeName: getTypeName(&types.RecipeRating{}), InputTypeName: ""},
	"POST /api/v1/settings/":                                                                       {ResponseTypeName: getTypeName(&types.ServiceSetting{}), InputTypeName: ""},
	"GET /api/v1/settings/":                                                                        {ResponseTypeName: getTypeName(&types.ServiceSetting{}), InputTypeName: ""},
	"POST /api/v1/settings/configurations/":                                                        {ResponseTypeName: getTypeName(&types.ServiceSettingConfiguration{}), InputTypeName: ""},
	"GET /api/v1/settings/configurations/household":                                                {ResponseTypeName: getTypeName(&types.ServiceSettingConfiguration{}), InputTypeName: ""},
	"GET /api/v1/settings/configurations/user":                                                     {ResponseTypeName: getTypeName(&types.ServiceSettingConfiguration{}), InputTypeName: ""},
	"GET /api/v1/settings/configurations/user/{serviceSettingConfigurationName}":                   {ResponseTypeName: getTypeName(&types.ServiceSettingConfiguration{}), InputTypeName: ""},
	"DELETE /api/v1/settings/configurations/{serviceSettingConfigurationID}":                       {ResponseTypeName: getTypeName(&types.ServiceSettingConfiguration{}), InputTypeName: ""},
	"PUT /api/v1/settings/configurations/{serviceSettingConfigurationID}":                          {ResponseTypeName: getTypeName(&types.ServiceSettingConfiguration{}), InputTypeName: ""},
	"GET /api/v1/settings/search":                                                                  {ResponseTypeName: getTypeName(&types.ServiceSetting{}), InputTypeName: ""},
	"GET /api/v1/settings/{serviceSettingID}/":                                                     {ResponseTypeName: getTypeName(&types.ServiceSetting{}), InputTypeName: ""},
	"DELETE /api/v1/settings/{serviceSettingID}/":                                                  {ResponseTypeName: getTypeName(&types.ServiceSetting{}), InputTypeName: ""},
	"POST /api/v1/user_ingredient_preferences/":                                                    {ResponseTypeName: getTypeName(&types.UserIngredientPreference{}), InputTypeName: ""},
	"GET /api/v1/user_ingredient_preferences/":                                                     {ResponseTypeName: getTypeName(&types.UserIngredientPreference{}), InputTypeName: ""},
	"PUT /api/v1/user_ingredient_preferences/{userIngredientPreferenceID}/":                        {ResponseTypeName: getTypeName(&types.UserIngredientPreference{}), InputTypeName: ""},
	"DELETE /api/v1/user_ingredient_preferences/{userIngredientPreferenceID}/":                     {ResponseTypeName: getTypeName(&types.UserIngredientPreference{}), InputTypeName: ""},
	"POST /api/v1/user_notifications/":                                                             {ResponseTypeName: getTypeName(&types.UserNotification{}), InputTypeName: ""},
	"GET /api/v1/user_notifications/":                                                              {ResponseTypeName: getTypeName(&types.UserNotification{}), InputTypeName: ""},
	"GET /api/v1/user_notifications/{userNotificationID}/":                                         {ResponseTypeName: getTypeName(&types.UserNotification{}), InputTypeName: ""},
	"PATCH /api/v1/user_notifications/{userNotificationID}/":                                       {ResponseTypeName: getTypeName(&types.UserNotification{}), InputTypeName: ""},
	"GET /api/v1/users/":                                                                           {ResponseTypeName: getTypeName(&types.User{}), InputTypeName: ""},
	"POST /api/v1/users/avatar/upload":                                                             {ResponseTypeName: "", InputTypeName: ""},
	"PUT /api/v1/users/details":                                                                    {ResponseTypeName: getTypeName(&types.User{}), InputTypeName: ""},
	"PUT /api/v1/users/email_address":                                                              {ResponseTypeName: getTypeName(&types.User{}), InputTypeName: ""},
	"POST /api/v1/users/email_address_verification":                                                {ResponseTypeName: getTypeName(&types.User{}), InputTypeName: ""},
	"POST /api/v1/users/household/select":                                                          {ResponseTypeName: getTypeName(&types.Household{}), InputTypeName: ""},
	"PUT /api/v1/users/password/new":                                                               {ResponseTypeName: getTypeName(&types.User{}), InputTypeName: ""},
	"POST /api/v1/users/permissions/check":                                                         {ResponseTypeName: "", InputTypeName: ""},
	"GET /api/v1/users/search":                                                                     {ResponseTypeName: getTypeName(&types.User{}), InputTypeName: ""},
	"GET /api/v1/users/self":                                                                       {ResponseTypeName: getTypeName(&types.User{}), InputTypeName: ""},
	"POST /api/v1/users/totp_secret/new":                                                           {ResponseTypeName: getTypeName(&types.APIError{}), InputTypeName: ""},
	"PUT /api/v1/users/username":                                                                   {ResponseTypeName: getTypeName(&types.User{}), InputTypeName: ""},
	"GET /api/v1/users/{userID}/":                                                                  {ResponseTypeName: getTypeName(&types.User{}), InputTypeName: ""},
	"DELETE /api/v1/users/{userID}/":                                                               {ResponseTypeName: getTypeName(&types.User{}), InputTypeName: ""},
	"POST /api/v1/valid_ingredient_groups/":                                                        {ResponseTypeName: getTypeName(&types.ValidIngredientGroup{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_groups/":                                                         {ResponseTypeName: getTypeName(&types.ValidIngredientGroup{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_groups/search":                                                   {ResponseTypeName: getTypeName(&types.ValidIngredientGroup{}), InputTypeName: ""},
	"DELETE /api/v1/valid_ingredient_groups/{validIngredientGroupID}/":                             {ResponseTypeName: getTypeName(&types.ValidIngredientGroup{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_groups/{validIngredientGroupID}/":                                {ResponseTypeName: getTypeName(&types.ValidIngredientGroup{}), InputTypeName: ""},
	"PUT /api/v1/valid_ingredient_groups/{validIngredientGroupID}/":                                {ResponseTypeName: getTypeName(&types.ValidIngredientGroup{}), InputTypeName: ""},
	"POST /api/v1/valid_ingredient_measurement_units/":                                             {ResponseTypeName: getTypeName(&types.ValidIngredientMeasurementUnit{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_measurement_units/":                                              {ResponseTypeName: getTypeName(&types.ValidIngredientMeasurementUnit{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_measurement_units/by_ingredient/{validIngredientID}/":            {ResponseTypeName: getTypeName(&types.ValidIngredientMeasurementUnit{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_measurement_units/by_measurement_unit/{validMeasurementUnitID}/": {ResponseTypeName: getTypeName(&types.ValidIngredientMeasurementUnit{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/":           {ResponseTypeName: getTypeName(&types.ValidIngredientMeasurementUnit{}), InputTypeName: ""},
	"PUT /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/":           {ResponseTypeName: getTypeName(&types.ValidIngredientMeasurementUnit{}), InputTypeName: ""},
	"DELETE /api/v1/valid_ingredient_measurement_units/{validIngredientMeasurementUnitID}/":        {ResponseTypeName: getTypeName(&types.ValidIngredientMeasurementUnit{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_preparations/":                                                   {ResponseTypeName: getTypeName(&types.ValidIngredientPreparation{}), InputTypeName: ""},
	"POST /api/v1/valid_ingredient_preparations/":                                                  {ResponseTypeName: getTypeName(&types.ValidIngredientPreparation{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_preparations/by_ingredient/{validIngredientID}/":                 {ResponseTypeName: getTypeName(&types.ValidIngredientPreparation{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_preparations/by_preparation/{validPreparationID}/":               {ResponseTypeName: getTypeName(&types.ValidIngredientPreparation{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/":                    {ResponseTypeName: getTypeName(&types.ValidIngredientPreparation{}), InputTypeName: ""},
	"PUT /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/":                    {ResponseTypeName: getTypeName(&types.ValidIngredientPreparation{}), InputTypeName: ""},
	"DELETE /api/v1/valid_ingredient_preparations/{validIngredientPreparationID}/":                 {ResponseTypeName: getTypeName(&types.ValidIngredientPreparation{}), InputTypeName: ""},
	"POST /api/v1/valid_ingredient_state_ingredients/":                                             {ResponseTypeName: getTypeName(&types.ValidIngredientStateIngredient{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_state_ingredients/":                                              {ResponseTypeName: getTypeName(&types.ValidIngredientStateIngredient{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_state_ingredients/by_ingredient/{validIngredientID}/":            {ResponseTypeName: getTypeName(&types.ValidIngredientStateIngredient{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_state_ingredients/by_ingredient_state/{validIngredientStateID}/": {ResponseTypeName: getTypeName(&types.ValidIngredientStateIngredient{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/":           {ResponseTypeName: getTypeName(&types.ValidIngredientStateIngredient{}), InputTypeName: ""},
	"PUT /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/":           {ResponseTypeName: getTypeName(&types.ValidIngredientStateIngredient{}), InputTypeName: ""},
	"DELETE /api/v1/valid_ingredient_state_ingredients/{validIngredientStateIngredientID}/":        {ResponseTypeName: getTypeName(&types.ValidIngredientStateIngredient{}), InputTypeName: ""},
	"POST /api/v1/valid_ingredient_states/":                                                        {ResponseTypeName: getTypeName(&types.ValidIngredientState{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_states/":                                                         {ResponseTypeName: getTypeName(&types.ValidIngredientState{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_states/search":                                                   {ResponseTypeName: getTypeName(&types.ValidIngredientState{}), InputTypeName: ""},
	"PUT /api/v1/valid_ingredient_states/{validIngredientStateID}/":                                {ResponseTypeName: getTypeName(&types.ValidIngredientState{}), InputTypeName: ""},
	"DELETE /api/v1/valid_ingredient_states/{validIngredientStateID}/":                             {ResponseTypeName: getTypeName(&types.ValidIngredientState{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredient_states/{validIngredientStateID}/":                                {ResponseTypeName: getTypeName(&types.ValidIngredientState{}), InputTypeName: ""},
	"POST /api/v1/valid_ingredients/":                                                              {ResponseTypeName: getTypeName(&types.ValidIngredient{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredients/":                                                               {ResponseTypeName: getTypeName(&types.ValidIngredient{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredients/by_preparation/{validPreparationID}/":                           {ResponseTypeName: getTypeName(&types.ValidIngredient{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredients/random":                                                         {ResponseTypeName: getTypeName(&types.ValidIngredient{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredients/search":                                                         {ResponseTypeName: getTypeName(&types.ValidIngredient{}), InputTypeName: ""},
	"PUT /api/v1/valid_ingredients/{validIngredientID}/":                                           {ResponseTypeName: getTypeName(&types.ValidIngredient{}), InputTypeName: ""},
	"DELETE /api/v1/valid_ingredients/{validIngredientID}/":                                        {ResponseTypeName: getTypeName(&types.ValidIngredient{}), InputTypeName: ""},
	"GET /api/v1/valid_ingredients/{validIngredientID}/":                                           {ResponseTypeName: getTypeName(&types.ValidIngredient{}), InputTypeName: ""},
	"GET /api/v1/valid_instruments/":                                                               {ResponseTypeName: getTypeName(&types.ValidInstrument{}), InputTypeName: ""},
	"POST /api/v1/valid_instruments/":                                                              {ResponseTypeName: getTypeName(&types.ValidInstrument{}), InputTypeName: ""},
	"GET /api/v1/valid_instruments/random":                                                         {ResponseTypeName: getTypeName(&types.ValidInstrument{}), InputTypeName: ""},
	"GET /api/v1/valid_instruments/search":                                                         {ResponseTypeName: getTypeName(&types.ValidInstrument{}), InputTypeName: ""},
	"DELETE /api/v1/valid_instruments/{validInstrumentID}/":                                        {ResponseTypeName: getTypeName(&types.ValidInstrument{}), InputTypeName: ""},
	"GET /api/v1/valid_instruments/{validInstrumentID}/":                                           {ResponseTypeName: getTypeName(&types.ValidInstrument{}), InputTypeName: ""},
	"PUT /api/v1/valid_instruments/{validInstrumentID}/":                                           {ResponseTypeName: getTypeName(&types.ValidInstrument{}), InputTypeName: ""},
	"POST /api/v1/valid_measurement_conversions/":                                                  {ResponseTypeName: getTypeName(&types.ValidMeasurementUnitConversion{}), InputTypeName: ""},
	"GET /api/v1/valid_measurement_conversions/from_unit/{validMeasurementUnitID}":                 {ResponseTypeName: getTypeName(&types.ValidMeasurementUnitConversion{}), InputTypeName: ""},
	"GET /api/v1/valid_measurement_conversions/to_unit/{validMeasurementUnitID}":                   {ResponseTypeName: getTypeName(&types.ValidMeasurementUnitConversion{}), InputTypeName: ""},
	"PUT /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/":                {ResponseTypeName: getTypeName(&types.ValidMeasurementUnitConversion{}), InputTypeName: ""},
	"DELETE /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/":             {ResponseTypeName: getTypeName(&types.ValidMeasurementUnitConversion{}), InputTypeName: ""},
	"GET /api/v1/valid_measurement_conversions/{validMeasurementUnitConversionID}/":                {ResponseTypeName: getTypeName(&types.ValidMeasurementUnitConversion{}), InputTypeName: ""},
	"POST /api/v1/valid_measurement_units/":                                                        {ResponseTypeName: getTypeName(&types.ValidMeasurementUnit{}), InputTypeName: ""},
	"GET /api/v1/valid_measurement_units/":                                                         {ResponseTypeName: getTypeName(&types.ValidMeasurementUnit{}), InputTypeName: ""},
	"GET /api/v1/valid_measurement_units/by_ingredient/{validIngredientID}":                        {ResponseTypeName: getTypeName(&types.ValidMeasurementUnit{}), InputTypeName: ""},
	"GET /api/v1/valid_measurement_units/search":                                                   {ResponseTypeName: getTypeName(&types.ValidMeasurementUnit{}), InputTypeName: ""},
	"GET /api/v1/valid_measurement_units/{validMeasurementUnitID}/":                                {ResponseTypeName: getTypeName(&types.ValidMeasurementUnit{}), InputTypeName: ""},
	"PUT /api/v1/valid_measurement_units/{validMeasurementUnitID}/":                                {ResponseTypeName: getTypeName(&types.ValidMeasurementUnit{}), InputTypeName: ""},
	"DELETE /api/v1/valid_measurement_units/{validMeasurementUnitID}/":                             {ResponseTypeName: getTypeName(&types.ValidMeasurementUnit{}), InputTypeName: ""},
	"GET /api/v1/valid_preparation_instruments/":                                                   {ResponseTypeName: getTypeName(&types.ValidPreparationInstrument{}), InputTypeName: ""},
	"POST /api/v1/valid_preparation_instruments/":                                                  {ResponseTypeName: getTypeName(&types.ValidPreparationInstrument{}), InputTypeName: ""},
	"GET /api/v1/valid_preparation_instruments/by_instrument/{validInstrumentID}/":                 {ResponseTypeName: getTypeName(&types.ValidPreparationInstrument{}), InputTypeName: ""},
	"GET /api/v1/valid_preparation_instruments/by_preparation/{validPreparationID}/":               {ResponseTypeName: getTypeName(&types.ValidPreparationInstrument{}), InputTypeName: ""},
	"DELETE /api/v1/valid_preparation_instruments/{validPreparationVesselID}/":                     {ResponseTypeName: getTypeName(&types.ValidPreparationInstrument{}), InputTypeName: ""},
	"GET /api/v1/valid_preparation_instruments/{validPreparationVesselID}/":                        {ResponseTypeName: getTypeName(&types.ValidPreparationInstrument{}), InputTypeName: ""},
	"PUT /api/v1/valid_preparation_instruments/{validPreparationVesselID}/":                        {ResponseTypeName: getTypeName(&types.ValidPreparationInstrument{}), InputTypeName: ""},
	"POST /api/v1/valid_preparation_vessels/":                                                      {ResponseTypeName: getTypeName(&types.ValidPreparationVessel{}), InputTypeName: ""},
	"GET /api/v1/valid_preparation_vessels/":                                                       {ResponseTypeName: getTypeName(&types.ValidPreparationVessel{}), InputTypeName: ""},
	"GET /api/v1/valid_preparation_vessels/by_preparation/{validPreparationID}/":                   {ResponseTypeName: getTypeName(&types.ValidPreparationVessel{}), InputTypeName: ""},
	"GET /api/v1/valid_preparation_vessels/by_vessel/{ValidVesselID}/":                             {ResponseTypeName: getTypeName(&types.ValidPreparationVessel{}), InputTypeName: ""},
	"PUT /api/v1/valid_preparation_vessels/{validPreparationVesselID}/":                            {ResponseTypeName: getTypeName(&types.ValidPreparationVessel{}), InputTypeName: ""},
	"DELETE /api/v1/valid_preparation_vessels/{validPreparationVesselID}/":                         {ResponseTypeName: getTypeName(&types.ValidPreparationVessel{}), InputTypeName: ""},
	"GET /api/v1/valid_preparation_vessels/{validPreparationVesselID}/":                            {ResponseTypeName: getTypeName(&types.ValidPreparationVessel{}), InputTypeName: ""},
	"GET /api/v1/valid_preparations/":                                                              {ResponseTypeName: getTypeName(&types.ValidPreparation{}), InputTypeName: ""},
	"POST /api/v1/valid_preparations/":                                                             {ResponseTypeName: getTypeName(&types.ValidPreparation{}), InputTypeName: ""},
	"GET /api/v1/valid_preparations/random":                                                        {ResponseTypeName: getTypeName(&types.ValidPreparation{}), InputTypeName: ""},
	"GET /api/v1/valid_preparations/search":                                                        {ResponseTypeName: getTypeName(&types.ValidPreparation{}), InputTypeName: ""},
	"PUT /api/v1/valid_preparations/{validPreparationID}/":                                         {ResponseTypeName: getTypeName(&types.ValidPreparation{}), InputTypeName: ""},
	"DELETE /api/v1/valid_preparations/{validPreparationID}/":                                      {ResponseTypeName: getTypeName(&types.ValidPreparation{}), InputTypeName: ""},
	"GET /api/v1/valid_preparations/{validPreparationID}/":                                         {ResponseTypeName: getTypeName(&types.ValidPreparation{}), InputTypeName: ""},
	"POST /api/v1/valid_vessels/":                                                                  {ResponseTypeName: getTypeName(&types.ValidVessel{}), InputTypeName: ""},
	"GET /api/v1/valid_vessels/":                                                                   {ResponseTypeName: getTypeName(&types.ValidVessel{}), InputTypeName: ""},
	"GET /api/v1/valid_vessels/random":                                                             {ResponseTypeName: getTypeName(&types.ValidVessel{}), InputTypeName: ""},
	"GET /api/v1/valid_vessels/search":                                                             {ResponseTypeName: getTypeName(&types.ValidVessel{}), InputTypeName: ""},
	"GET /api/v1/valid_vessels/{validVesselID}/":                                                   {ResponseTypeName: getTypeName(&types.ValidVessel{}), InputTypeName: ""},
	"PUT /api/v1/valid_vessels/{validVesselID}/":                                                   {ResponseTypeName: getTypeName(&types.ValidVessel{}), InputTypeName: ""},
	"DELETE /api/v1/valid_vessels/{validVesselID}/":                                                {ResponseTypeName: getTypeName(&types.ValidVessel{}), InputTypeName: ""},
	"GET /api/v1/webhooks/":                                                                        {ResponseTypeName: getTypeName(&types.Webhook{}), InputTypeName: ""},
	"POST /api/v1/webhooks/":                                                                       {ResponseTypeName: getTypeName(&types.Webhook{}), InputTypeName: ""},
	"GET /api/v1/webhooks/{webhookID}/":                                                            {ResponseTypeName: getTypeName(&types.Webhook{}), InputTypeName: ""},
	"DELETE /api/v1/webhooks/{webhookID}/":                                                         {ResponseTypeName: getTypeName(&types.Webhook{}), InputTypeName: ""},
	"POST /api/v1/webhooks/{webhookID}/trigger_events":                                             {ResponseTypeName: getTypeName(&types.WebhookTriggerEvent{}), InputTypeName: ""},
	"DELETE /api/v1/webhooks/{webhookID}/trigger_events/{webhookTriggerEventID}/":                  {ResponseTypeName: getTypeName(&types.WebhookTriggerEvent{}), InputTypeName: ""},
	"POST /api/v1/workers/finalize_meal_plans":                                                     {ResponseTypeName: "", InputTypeName: ""},
	"POST /api/v1/workers/meal_plan_grocery_list_init":                                             {ResponseTypeName: "", InputTypeName: ""},
	"POST /api/v1/workers/meal_plan_tasks":                                                         {ResponseTypeName: "", InputTypeName: ""},
	"GET /auth/status":                                                                             {ResponseTypeName: getTypeName(&types.UserStatusResponse{}), InputTypeName: ""},
	"GET /auth/{auth_provider}":                                                                    {ResponseTypeName: "", InputTypeName: ""},
	"GET /auth/{auth_provider}/callback":                                                           {ResponseTypeName: "", InputTypeName: ""},
	"GET /oauth2/authorize":                                                                        {ResponseTypeName: "", InputTypeName: ""},
	"POST /oauth2/token":                                                                           {ResponseTypeName: "", InputTypeName: ""},
	"POST /users/":                                                                                 {ResponseTypeName: getTypeName(&types.User{}), InputTypeName: ""},
	"POST /users/email_address/verify":                                                             {ResponseTypeName: "", InputTypeName: ""},
	"POST /users/login":                                                                            {ResponseTypeName: "", InputTypeName: ""},
	"POST /users/login/admin":                                                                      {ResponseTypeName: "", InputTypeName: ""},
	"POST /users/logout":                                                                           {ResponseTypeName: "", InputTypeName: ""},
	"POST /users/password/reset":                                                                   {ResponseTypeName: getTypeName(&types.PasswordResetToken{}), InputTypeName: ""},
	"POST /users/password/reset/redeem":                                                            {ResponseTypeName: "", InputTypeName: ""},
	"POST /users/totp_secret/verify":                                                               {ResponseTypeName: "", InputTypeName: ""},
	"POST /users/username/reminder":                                                                {ResponseTypeName: "", InputTypeName: ""},
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
			ResponseType:  routeInfoMap[fmt.Sprintf("%s %s", route.Method, route.Path)].ResponseTypeName,
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
