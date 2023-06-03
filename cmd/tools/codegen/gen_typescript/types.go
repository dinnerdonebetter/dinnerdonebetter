package main

import (
	"regexp"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	filesToGenerate = map[string][]any{
		"admin": {
			types.ModifyUserPermissionsInput{},
		},
		"apiClients": {
			types.APIClient{},
			types.APIClientCreationRequestInput{},
			types.APIClientCreationResponse{},
		},
		"auth": {
			types.ChangeActiveHouseholdInput{},
			types.PASETOCreationInput{},
			types.PASETOResponse{},
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
		"mealRatings": {
			types.MealRating{},
			types.MealRatingCreationRequestInput{},
			types.MealRatingUpdateRequestInput{},
		},
		"householdInstrumentOwnerships": {
			types.HouseholdInstrumentOwnership{},
			types.HouseholdInstrumentOwnershipCreationRequestInput{},
			types.HouseholdInstrumentOwnershipUpdateRequestInput{},
		},
		"userFeedback": {
			types.UserFeedback{},
			types.UserFeedbackCreationRequestInput{},
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
		"validPreparations": {
			types.ValidPreparation{},
			types.ValidPreparationCreationRequestInput{},
			types.ValidPreparationUpdateRequestInput{},
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

	generatedDisclaimer = "// Code generated by gen_typescript. DO NOT EDIT.\n\n"
	customTypeMap       = map[string]string{
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
	}

	defaultValues = map[string]string{
		"ValidMealPlanStatus":                "'awaiting_votes'",
		"ValidMealPlanGroceryListItemStatus": "'unknown'",
		"ValidMealPlanElectionMethod":        "'schulze'",
		"ValidIngredientStateAttributeType":  "'other'",
		"ValidRecipeStepProductType":         "'ingredient'",
		"MealComponentType":                  "'unspecified'",
		"MealPlanTaskStatus":                 "'unfinished'",
	}

	// Times I've tried to optimize this regex before realizing it already accounts for
	// every edge case and there is no value (in either performance or readability terms)
	// in making it smaller: 1.
	numberMatcherRegex = regexp.MustCompile(`((u)?int(8|16|32|64)?|float(32|64))`)
)

func isCustomType(x string) bool {
	switch x {
	case "int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
		"float32",
		"float64",
		boolType,
		mapStringToBoolType,
		timeType,
		stringType:
		return false
	default:
		return true
	}
}
