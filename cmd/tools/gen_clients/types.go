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
		"auth": {
			types.ChangeActiveHouseholdInput{},
			types.PasswordResetToken{},
			types.PasswordResetTokenCreationRequestInput{},
			types.PasswordResetTokenRedemptionRequestInput{},
			types.TOTPSecretRefreshInput{},
			types.TOTPSecretVerificationInput{},
			types.TOTPSecretRefreshResponse{},
			types.PasswordUpdateInput{},
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
			types.HouseholdUserMembershipCreationRequestInput{},
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
			types.RecipeStepCompletionConditionCreationRequestInput{},
			types.RecipeStepCompletionConditionIngredientCreationRequestInput{},
			types.RecipeStepCompletionConditionUpdateRequestInput{},
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
		},
		"userIngredientPreferences": {
			types.UserIngredientPreference{},
			types.UserIngredientPreferenceCreationRequestInput{},
			types.UserIngredientPreferenceUpdateRequestInput{},
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
			types.ServiceSettingUpdateRequestInput{},
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
		},
	}

	TypesWeCareAbout = []TypeDefinition{
		{
			Type: &types.ModifyUserPermissionsInput{},
		},
		{
			Type: &types.OAuth2Client{},
		},
		{
			Type: &types.OAuth2ClientCreationRequestInput{},
		},
		{
			Type: &types.OAuth2ClientCreationResponse{},
		},
		{
			Type: &types.ChangeActiveHouseholdInput{},
		},
		{
			Type: &types.PasswordResetToken{},
		},
		{
			Type: &types.PasswordResetTokenCreationRequestInput{},
		},
		{
			Type: &types.PasswordResetTokenRedemptionRequestInput{},
		},
		{
			Type: &types.TOTPSecretRefreshInput{},
		},
		{
			Type: &types.TOTPSecretVerificationInput{},
		},
		{
			Type: &types.TOTPSecretRefreshResponse{},
		},
		{
			Type: &types.PasswordUpdateInput{},
		},
		{
			Type: &types.APIError{},
		},
		{
			Type: &types.HouseholdInvitation{},
		},
		{
			Type: &types.HouseholdInvitationUpdateRequestInput{},
		},
		{
			Type: &types.HouseholdInvitationCreationRequestInput{},
		},
		{
			Type: &types.Household{},
		},
		{
			Type: &types.HouseholdCreationRequestInput{},
		},
		{
			Type: &types.HouseholdUpdateRequestInput{},
		},
		{
			Type: &types.HouseholdOwnershipTransferInput{},
		},
		{
			Type: &types.HouseholdUserMembership{},
		},
		{
			Type: &types.HouseholdUserMembershipWithUser{},
		},
		{
			Type: &types.HouseholdUserMembershipCreationRequestInput{},
		},
		{
			Type: &types.MealPlanEvent{},
		},
		{
			Type: &types.MealPlanEventCreationRequestInput{},
		},
		{
			Type: &types.MealPlanEventUpdateRequestInput{},
		},
		{
			Type: &types.MealPlanGroceryListItem{},
		},
		{
			Type: &types.MealPlanGroceryListItemCreationRequestInput{},
		},
		{
			Type: &types.MealPlanGroceryListItemUpdateRequestInput{},
		},
		{
			Type: &types.MealPlanOption{},
		},
		{
			Type: &types.MealPlanOptionCreationRequestInput{},
		},
		{
			Type: &types.MealPlanOptionUpdateRequestInput{},
		},
		{
			Type: &types.MealPlanOptionVote{},
		},
		{
			Type: &types.MealPlanOptionVoteCreationInput{},
		},
		{
			Type: &types.MealPlanOptionVoteCreationRequestInput{},
		},
		{
			Type: &types.MealPlanOptionVoteUpdateRequestInput{},
		},
		{
			Type: &types.MealPlan{},
		},
		{
			Type: &types.MealPlanCreationRequestInput{},
		},
		{
			Type: &types.MealPlanUpdateRequestInput{},
		},
		{
			Type: &types.MealPlanTask{},
		},
		{
			Type: &types.MealPlanTaskCreationRequestInput{},
		},
		{
			Type: &types.MealPlanTaskStatusChangeRequestInput{},
		},
		{
			Type: &types.Meal{},
		},
		{
			Type: &types.MealCreationRequestInput{},
		},
		{
			Type: &types.MealUpdateRequestInput{},
		},
		{
			Type: &types.RecipeRating{},
		},
		{
			Type: &types.RecipeRatingCreationRequestInput{},
		},
		{
			Type: &types.RecipeRatingUpdateRequestInput{},
		},
		{
			Type: &types.HouseholdInstrumentOwnership{},
		},
		{
			Type: &types.HouseholdInstrumentOwnershipCreationRequestInput{},
		},
		{
			Type: &types.HouseholdInstrumentOwnershipUpdateRequestInput{},
		},
		{
			Type: &types.MealComponent{},
		},
		{
			Type: &types.MealComponentCreationRequestInput{},
		},
		{
			Type: &types.MealComponentUpdateRequestInput{},
		},
		{
			Type: &types.UserPermissionsRequestInput{},
		},
		{
			Type: &types.UserPermissionsResponse{},
		},
		{
			Type: &types.RecipeMedia{},
		},
		{
			Type: &types.RecipeMediaCreationRequestInput{},
		},
		{
			Type: &types.RecipeMediaUpdateRequestInput{},
		},
		{
			Type: &types.RecipePrepTask{},
		},
		{
			Type: &types.RecipePrepTaskCreationRequestInput{},
		},
		{
			Type: &types.RecipePrepTaskWithinRecipeCreationRequestInput{},
		},
		{
			Type: &types.RecipePrepTaskUpdateRequestInput{},
		},
		{
			Type: &types.RecipePrepTaskStep{},
		},
		{
			Type: &types.RecipePrepTaskStepWithinRecipeCreationRequestInput{},
		},
		{
			Type: &types.RecipePrepTaskStepCreationRequestInput{},
		},
		{
			Type: &types.RecipePrepTaskStepUpdateRequestInput{},
		},
		{
			Type: &types.RecipeStepCompletionCondition{},
		},
		{
			Type: &types.RecipeStepCompletionConditionIngredient{},
		},
		{
			Type: &types.RecipeStepCompletionConditionCreationRequestInput{},
		},
		{
			Type: &types.RecipeStepCompletionConditionIngredientCreationRequestInput{},
		},
		{
			Type: &types.RecipeStepCompletionConditionUpdateRequestInput{},
		},
		{
			Type: &types.RecipeStepIngredient{},
		},
		{
			Type: &types.RecipeStepIngredientCreationRequestInput{},
		},
		{
			Type: &types.RecipeStepIngredientUpdateRequestInput{},
		},
		{
			Type: &types.RecipeStepInstrument{},
		},
		{
			Type: &types.RecipeStepInstrumentCreationRequestInput{},
		},
		{
			Type: &types.RecipeStepInstrumentUpdateRequestInput{},
		},
		{
			Type: &types.RecipeStepVessel{},
		},
		{
			Type: &types.RecipeStepVesselCreationRequestInput{},
		},
		{
			Type: &types.RecipeStepVesselUpdateRequestInput{},
		},
		{
			Type: &types.RecipeStepProduct{},
		},
		{
			Type: &types.RecipeStepProductCreationRequestInput{},
		},
		{
			Type: &types.RecipeStepProductUpdateRequestInput{},
		},
		{
			Type: &types.RecipeStep{},
		},
		{
			Type: &types.RecipeStepCreationRequestInput{},
		},
		{
			Type: &types.RecipeStepUpdateRequestInput{},
		},
		{
			Type: &types.Recipe{},
		},
		{
			Type: &types.RecipeCreationRequestInput{},
		},
		{
			Type: &types.RecipeUpdateRequestInput{},
		},
		{
			Type: &types.UserStatusResponse{},
		},
		{
			Type: &types.User{},
		},
		{
			Type: &types.UserRegistrationInput{},
		},
		{
			Type: &types.UserCreationResponse{},
		},
		{
			Type: &types.UserLoginInput{},
		},
		{
			Type: &types.UsernameReminderRequestInput{},
		},
		{
			Type: &types.UserAccountStatusUpdateInput{},
		},
		{
			Type: &types.EmailAddressVerificationRequestInput{},
		},
		{
			Type: &types.AvatarUpdateInput{},
		},
		{
			Type: &types.UserIngredientPreference{},
		},
		{
			Type: &types.UserIngredientPreferenceCreationRequestInput{},
		},
		{
			Type: &types.UserIngredientPreferenceUpdateRequestInput{},
		},
		{
			Type: &types.ValidIngredientMeasurementUnit{},
		},
		{
			Type: &types.ValidIngredientMeasurementUnitCreationRequestInput{},
		},
		{
			Type: &types.ValidIngredientMeasurementUnitUpdateRequestInput{},
		},
		{
			Type: &types.ValidIngredientPreparation{},
		},
		{
			Type: &types.ValidIngredientPreparationCreationRequestInput{},
		},
		{
			Type: &types.ValidIngredientPreparationUpdateRequestInput{},
		},
		{
			Type: &types.ValidIngredientState{},
		},
		{
			Type: &types.ValidIngredientStateCreationRequestInput{},
		},
		{
			Type: &types.ValidIngredientStateUpdateRequestInput{},
		},
		{
			Type: &types.ValidIngredientStateIngredient{},
		},
		{
			Type: &types.ValidIngredientStateIngredientCreationRequestInput{},
		},
		{
			Type: &types.ValidIngredientStateIngredientUpdateRequestInput{},
		},
		{
			Type: &types.ValidIngredient{},
		},
		{
			Type: &types.ValidIngredientCreationRequestInput{},
		},
		{
			Type: &types.ValidIngredientUpdateRequestInput{},
		},
		{
			Type: &types.ValidIngredientGroup{},
		},
		{
			Type: &types.ValidIngredientGroupCreationRequestInput{},
		},
		{
			Type: &types.ValidIngredientGroupUpdateRequestInput{},
		},
		{
			Type: &types.ValidIngredientGroupMember{},
		},
		{
			Type: &types.ValidIngredientGroupMemberCreationRequestInput{},
		},
		{
			Type: &types.ValidInstrument{},
		},
		{
			Type: &types.ValidInstrumentCreationRequestInput{},
		},
		{
			Type: &types.ValidInstrumentUpdateRequestInput{},
		},
		{
			Type: &types.ValidMeasurementUnitConversion{},
		},
		{
			Type: &types.ValidMeasurementUnitConversionCreationRequestInput{},
		},
		{
			Type: &types.ValidMeasurementUnitConversionUpdateRequestInput{},
		},
		{
			Type: &types.ValidMeasurementUnit{},
		},
		{
			Type: &types.ValidMeasurementUnitCreationRequestInput{},
		},
		{
			Type: &types.ValidMeasurementUnitUpdateRequestInput{},
		},
		{
			Type: &types.ValidPreparationInstrument{},
		},
		{
			Type: &types.ValidPreparationInstrumentCreationRequestInput{},
		},
		{
			Type: &types.ValidPreparationInstrumentUpdateRequestInput{},
		},
		{
			Type: &types.ValidPreparationVessel{},
		},
		{
			Type: &types.ValidPreparationVesselCreationRequestInput{},
		},
		{
			Type: &types.ValidPreparationVesselUpdateRequestInput{},
		},
		{
			Type: &types.ValidPreparation{},
		},
		{
			Type: &types.ValidPreparationCreationRequestInput{},
		},
		{
			Type: &types.ValidPreparationUpdateRequestInput{},
		},
		{
			Type: &types.ValidVessel{},
		},
		{
			Type: &types.ValidVesselCreationRequestInput{},
		},
		{
			Type: &types.ValidVesselUpdateRequestInput{},
		},
		{
			Type: &types.ServiceSetting{},
		},
		{
			Type: &types.ServiceSettingCreationRequestInput{},
		},
		{
			Type: &types.ServiceSettingUpdateRequestInput{},
		},
		{
			Type: &types.ServiceSettingConfiguration{},
		},
		{
			Type: &types.ServiceSettingConfigurationCreationRequestInput{},
		},
		{
			Type: &types.ServiceSettingConfigurationUpdateRequestInput{},
		},
		{
			Type: &types.Webhook{},
		},
		{
			Type: &types.WebhookTriggerEvent{},
		},
		{
			Type: &types.WebhookCreationRequestInput{},
		},
	}
)
