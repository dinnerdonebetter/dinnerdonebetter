package codegen

import (
	"reflect"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

type TypeDefinition struct {
	Type        any
	Description string
}

func (d *TypeDefinition) Name() string {
	n := reflect.Indirect(reflect.ValueOf(d.Type)).Type().Name()
	return n
}

var (
	CustomTypeMap = map[string]string{
		"ValidIngredientState.attributeType":                     "ValidIngredientStateAttributeType",
		"ValidIngredientStateCreationRequestInput.attributeType": "ValidIngredientStateAttributeType",
		"ValidIngredientStateUpdateRequestInput.attributeType":   "ValidIngredientStateAttributeType",
		"RecipeStepProduct.type":                                 "ValidRecipeStepProductType",
		"RecipeStepProductCreationRequestInput.type":             "ValidRecipeStepProductType",
		"RecipeStepProductUpdateRequestInput.type":               "ValidRecipeStepProductType",
		"MealPlanGroceryListItem.status":                         "ValidMealPlanGroceryListItemStatus",
		"MealPlanGroceryListItemCreationRequestInput.status":     "ValidMealPlanGroceryListItemStatus",
		"MealPlanGroceryListItemUpdateRequestInput.status":       "ValidMealPlanGroceryListItemStatus",
		"MealComponent.componentType":                            "MealComponentType",
		"MealComponentCreationRequestInput.componentType":        "MealComponentType",
		"MealComponentUpdateRequestInput.componentType":          "MealComponentType",
		"MealPlanTaskStatusChangeRequestInput.status":            "MealPlanTaskStatus",
		"MealPlanTask.status":                                    "MealPlanTaskStatus",
		"MealPlanTaskCreationRequestInput.status":                "MealPlanTaskStatus",
		"MealPlanTaskUpdateRequestInput.status":                  "MealPlanTaskStatus",
		"MealPlan.status":                                        "ValidMealPlanStatus",
		"MealPlanCreationRequestInput.status":                    "ValidMealPlanStatus",
		"MealPlanUpdateRequestInput.status":                      "ValidMealPlanStatus",
		"MealPlan.electionMethod":                                "ValidMealPlanElectionMethod",
		"MealPlanUpdateRequestInput.electionMethod":              "ValidMealPlanElectionMethod",
		"MealPlanCreationRequestInput.electionMethod":            "ValidMealPlanElectionMethod",
		"ValidVessel.shape":                                      "ValidVesselShapeType",
		"ValidVesselUpdateRequestInput.shape":                    "ValidVesselShapeType",
		"ValidVesselCreationRequestInput.shape":                  "ValidVesselShapeType",
	}

	DefaultEnumValues = map[string]string{
		"ValidMealPlanStatus":                "'awaiting_votes'",
		"ValidMealPlanGroceryListItemStatus": "'unknown'",
		"ValidMealPlanElectionMethod":        "'schulze'",
		"ValidIngredientStateAttributeType":  "'other'",
		"ValidRecipeStepProductType":         "'ingredient'",
		"MealComponentType":                  "'unspecified'",
		"MealPlanTaskStatus":                 "'unfinished'",
		"ValidVesselShapeType":               "'other'",
	}

	TypeDefinitionFilesToGenerate = map[string][]any{
		"admin": {
			types.ModifyUserPermissionsInput{},
		},
		"oauth2Clients": {
			types.OAuth2Client{},
			types.OAuth2ClientCreationRequestInput{},
			types.OAuth2ClientCreationResponse{},
		},
		"auditLogEntries": {
			types.ChangeLog{},
			types.AuditLogEntry{},
		},
		"auth": {
			types.ChangeActiveHouseholdInput{},
			types.PasswordResetToken{},
			types.PasswordResetTokenCreationRequestInput{},
			types.PasswordResetTokenRedemptionRequestInput{},
			types.TOTPSecretRefreshInput{},
			types.TOTPSecretVerificationInput{},
			types.TOTPSecretRefreshResponse{},
			types.PasswordUpdateInput{},
			types.JWTResponse{},
		},
		"errors": {
			types.APIError{},
		},
		"householdInvitations": {
			types.HouseholdInvitation{},
			types.HouseholdInvitationUpdateRequestInput{},
			types.HouseholdInvitationCreationRequestInput{},
		},
		"households": {
			types.Household{},
			types.HouseholdCreationRequestInput{},
			types.HouseholdUpdateRequestInput{},
			types.HouseholdOwnershipTransferInput{},
		},
		"householdUserMemberships": {
			types.HouseholdUserMembership{},
			types.HouseholdUserMembershipWithUser{},
		},
		"mealPlanEvents": {
			types.MealPlanEvent{},
			types.MealPlanEventCreationRequestInput{},
			types.MealPlanEventUpdateRequestInput{},
		},
		"mealPlanGroceryListItems": {
			types.MealPlanGroceryListItem{},
			types.MealPlanGroceryListItemCreationRequestInput{},
			types.MealPlanGroceryListItemUpdateRequestInput{},
		},
		"mealPlanOptions": {
			types.MealPlanOption{},
			types.MealPlanOptionCreationRequestInput{},
			types.MealPlanOptionUpdateRequestInput{},
		},
		"mealPlanOptionVotes": {
			types.MealPlanOptionVote{},
			types.MealPlanOptionVoteCreationInput{},
			types.MealPlanOptionVoteCreationRequestInput{},
			types.MealPlanOptionVoteUpdateRequestInput{},
		},
		"mealPlans": {
			types.MealPlan{},
			types.MealPlanCreationRequestInput{},
			types.MealPlanUpdateRequestInput{},
		},
		"mealPlanTasks": {
			types.MealPlanTask{},
			types.MealPlanTaskCreationRequestInput{},
			types.MealPlanTaskStatusChangeRequestInput{},
		},
		"meals": {
			types.Meal{},
			types.MealCreationRequestInput{},
			types.MealUpdateRequestInput{},
		},
		"recipeRatings": {
			types.RecipeRating{},
			types.RecipeRatingCreationRequestInput{},
			types.RecipeRatingUpdateRequestInput{},
		},
		"householdInstrumentOwnerships": {
			types.HouseholdInstrumentOwnership{},
			types.HouseholdInstrumentOwnershipCreationRequestInput{},
			types.HouseholdInstrumentOwnershipUpdateRequestInput{},
		},
		"mealComponents": {
			types.MealComponent{},
			types.MealComponentCreationRequestInput{},
			types.MealComponentUpdateRequestInput{},
		},
		"permissions": {
			types.UserPermissionsRequestInput{},
			types.UserPermissionsResponse{},
		},
		"recipeMedia": {
			types.RecipeMedia{},
			types.RecipeMediaCreationRequestInput{},
			types.RecipeMediaUpdateRequestInput{},
		},
		"recipePrepTasks": {
			types.RecipePrepTask{},
			types.RecipePrepTaskCreationRequestInput{},
			types.RecipePrepTaskWithinRecipeCreationRequestInput{},
			types.RecipePrepTaskUpdateRequestInput{},
		},
		"recipePrepTaskSteps": {
			types.RecipePrepTaskStep{},
			types.RecipePrepTaskStepWithinRecipeCreationRequestInput{},
			types.RecipePrepTaskStepCreationRequestInput{},
			types.RecipePrepTaskStepUpdateRequestInput{},
		},
		"recipeStepCompletionConditions": {
			types.RecipeStepCompletionCondition{},
			types.RecipeStepCompletionConditionIngredient{},
			types.RecipeStepCompletionConditionUpdateRequestInput{},
			types.RecipeStepCompletionConditionCreationRequestInput{},
			types.RecipeStepCompletionConditionForExistingRecipeCreationRequestInput{},
			types.RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput{},
		},
		"recipeStepIngredients": {
			types.RecipeStepIngredient{},
			types.RecipeStepIngredientCreationRequestInput{},
			types.RecipeStepIngredientUpdateRequestInput{},
		},
		"recipeStepInstruments": {
			types.RecipeStepInstrument{},
			types.RecipeStepInstrumentCreationRequestInput{},
			types.RecipeStepInstrumentUpdateRequestInput{},
		},
		"recipeStepVessels": {
			types.RecipeStepVessel{},
			types.RecipeStepVesselCreationRequestInput{},
			types.RecipeStepVesselUpdateRequestInput{},
		},
		"recipeStepProducts": {
			types.RecipeStepProduct{},
			types.RecipeStepProductCreationRequestInput{},
			types.RecipeStepProductUpdateRequestInput{},
		},
		"recipeSteps": {
			types.RecipeStep{},
			types.RecipeStepCreationRequestInput{},
			types.RecipeStepUpdateRequestInput{},
		},
		"recipes": {
			types.Recipe{},
			types.RecipeCreationRequestInput{},
			types.RecipeUpdateRequestInput{},
		},
		"users": {
			types.UserStatusResponse{},
			types.User{},
			types.UserRegistrationInput{},
			types.UserCreationResponse{},
			types.UserLoginInput{},
			types.UsernameReminderRequestInput{},
			types.UserAccountStatusUpdateInput{},
			types.EmailAddressVerificationRequestInput{},
			types.AvatarUpdateInput{},
			types.UsernameUpdateInput{},
			types.UserDetailsUpdateRequestInput{},
			types.UserEmailAddressUpdateInput{},
			types.PasswordResetResponse{},
		},
		"userIngredientPreferences": {
			types.UserIngredientPreference{},
			types.UserIngredientPreferenceCreationRequestInput{},
			types.UserIngredientPreferenceUpdateRequestInput{},
		},
		"userNotifications": {
			types.UserNotification{},
			types.UserNotificationCreationRequestInput{},
			types.UserNotificationUpdateRequestInput{},
		},
		"validIngredientMeasurementUnits": {
			types.ValidIngredientMeasurementUnit{},
			types.ValidIngredientMeasurementUnitCreationRequestInput{},
			types.ValidIngredientMeasurementUnitUpdateRequestInput{},
		},
		"validIngredientPreparations": {
			types.ValidIngredientPreparation{},
			types.ValidIngredientPreparationCreationRequestInput{},
			types.ValidIngredientPreparationUpdateRequestInput{},
		},
		"validIngredientStates": {
			types.ValidIngredientState{},
			types.ValidIngredientStateCreationRequestInput{},
			types.ValidIngredientStateUpdateRequestInput{},
		},
		"validIngredientStateIngredients": {
			types.ValidIngredientStateIngredient{},
			types.ValidIngredientStateIngredientCreationRequestInput{},
			types.ValidIngredientStateIngredientUpdateRequestInput{},
		},
		"validIngredients": {
			types.ValidIngredient{},
			types.ValidIngredientCreationRequestInput{},
			types.ValidIngredientUpdateRequestInput{},
		},
		"validIngredientGroups": {
			types.ValidIngredientGroup{},
			types.ValidIngredientGroupCreationRequestInput{},
			types.ValidIngredientGroupUpdateRequestInput{},
			types.ValidIngredientGroupMember{},
			types.ValidIngredientGroupMemberCreationRequestInput{},
		},
		"validInstruments": {
			types.ValidInstrument{},
			types.ValidInstrumentCreationRequestInput{},
			types.ValidInstrumentUpdateRequestInput{},
		},
		"validMeasurementUnitConversions": {
			types.ValidMeasurementUnitConversion{},
			types.ValidMeasurementUnitConversionCreationRequestInput{},
			types.ValidMeasurementUnitConversionUpdateRequestInput{},
		},
		"validMeasurementUnits": {
			types.ValidMeasurementUnit{},
			types.ValidMeasurementUnitCreationRequestInput{},
			types.ValidMeasurementUnitUpdateRequestInput{},
		},
		"validPreparationInstruments": {
			types.ValidPreparationInstrument{},
			types.ValidPreparationInstrumentCreationRequestInput{},
			types.ValidPreparationInstrumentUpdateRequestInput{},
		},
		"validPreparationVessels": {
			types.ValidPreparationVessel{},
			types.ValidPreparationVesselCreationRequestInput{},
			types.ValidPreparationVesselUpdateRequestInput{},
		},
		"validPreparations": {
			types.ValidPreparation{},
			types.ValidPreparationCreationRequestInput{},
			types.ValidPreparationUpdateRequestInput{},
		},
		"validVessels": {
			types.ValidVessel{},
			types.ValidVesselCreationRequestInput{},
			types.ValidVesselUpdateRequestInput{},
		},
		"serviceSetting": {
			types.ServiceSetting{},
			types.ServiceSettingCreationRequestInput{},
		},
		"serviceSettingConfiguration": {
			types.ServiceSettingConfiguration{},
			types.ServiceSettingConfigurationCreationRequestInput{},
			types.ServiceSettingConfigurationUpdateRequestInput{},
		},
		"webhooks": {
			types.Webhook{},
			types.WebhookTriggerEvent{},
			types.WebhookCreationRequestInput{},
			types.WebhookTriggerEventCreationRequestInput{},
		},
		"workers": {
			types.FinalizeMealPlansRequest{},
			types.FinalizeMealPlansResponse{},
			types.InitializeMealPlanGroceryListRequest{},
			types.InitializeMealPlanGroceryListResponse{},
			types.CreateMealPlanTasksRequest{},
			types.CreateMealPlanTasksResponse{},
		},
	}
)
