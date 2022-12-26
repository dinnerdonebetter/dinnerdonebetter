package main

import (
	"regexp"

	"github.com/prixfixeco/backend/pkg/types"
)

var (
	filesToGenerate = map[string][]any{
		"admin.ts": {
			types.ModifyUserPermissionsInput{},
		},
		"apiClients.ts": {
			types.APIClient{},
			types.APIClientCreationRequestInput{},
			types.APIClientCreationResponse{},
		},
		"auth.ts": {
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
		"errors.ts": {
			types.APIError{},
		},
		"householdInvitations.ts": {
			types.HouseholdInvitation{},
			types.HouseholdInvitationUpdateRequestInput{},
			types.HouseholdInvitationCreationRequestInput{},
		},
		"households.ts": {
			types.Household{},
			types.HouseholdCreationRequestInput{},
			types.HouseholdUpdateRequestInput{},
			types.HouseholdOwnershipTransferInput{},
		},
		"householdUserMemberships.ts": {
			types.HouseholdUserMembership{},
			types.HouseholdUserMembershipWithUser{},
			types.HouseholdUserMembershipCreationRequestInput{},
		},
		"mealPlanEvents.ts": {
			types.MealPlanEvent{},
			types.MealPlanEventCreationRequestInput{},
			types.MealPlanEventUpdateRequestInput{},
		},
		"mealPlanGroceryListItems.ts": {
			types.MealPlanGroceryListItem{},
			types.MealPlanGroceryListItemCreationRequestInput{},
			types.MealPlanGroceryListItemUpdateRequestInput{},
		},
		"mealPlanOptions.ts": {
			types.MealPlanOption{},
			types.MealPlanOptionCreationRequestInput{},
			types.MealPlanOptionUpdateRequestInput{},
		},
		"mealPlanOptionVotes.ts": {
			types.MealPlanOptionVote{},
			types.MealPlanOptionVoteCreationInput{},
			types.MealPlanOptionVoteCreationRequestInput{},
			types.MealPlanOptionVoteUpdateRequestInput{},
		},
		"mealPlans.ts": {
			types.MealPlan{},
			types.MealPlanCreationRequestInput{},
			types.MealPlanUpdateRequestInput{},
		},
		"mealPlanTasks.ts": {
			types.MealPlanTask{},
			types.MealPlanTaskCreationRequestInput{},
			types.MealPlanTaskStatusChangeRequestInput{},
		},
		"meals.ts": {
			types.Meal{},
			types.MealCreationRequestInput{},
			types.MealUpdateRequestInput{},
		},
		"mealComponents.ts": {
			types.MealComponent{},
			types.MealComponentCreationRequestInput{},
			types.MealComponentUpdateRequestInput{},
		},
		"permissions.ts": {
			types.UserPermissionsRequestInput{},
			types.UserPermissionsResponse{},
		},
		"recipeMedia.ts": {
			types.RecipeMedia{},
			types.RecipeMediaCreationRequestInput{},
			types.RecipeMediaUpdateRequestInput{},
		},
		"recipePrepTasks.ts": {
			types.RecipePrepTask{},
			types.RecipePrepTaskCreationRequestInput{},
			types.RecipePrepTaskWithinRecipeCreationRequestInput{},
			types.RecipePrepTaskUpdateRequestInput{},
		},
		"recipePrepTaskSteps.ts": {
			types.RecipePrepTaskStep{},
			types.RecipePrepTaskStepWithinRecipeCreationRequestInput{},
			types.RecipePrepTaskStepCreationRequestInput{},
			types.RecipePrepTaskStepUpdateRequestInput{},
		},
		"recipeStepCompletionConditions.ts": {
			types.RecipeStepCompletionCondition{},
			types.RecipeStepCompletionConditionIngredient{},
			types.RecipeStepCompletionConditionCreationRequestInput{},
			types.RecipeStepCompletionConditionIngredientCreationRequestInput{},
			types.RecipeStepCompletionConditionUpdateRequestInput{},
		},
		"recipeStepIngredients.ts": {
			types.RecipeStepIngredient{},
			types.RecipeStepIngredientCreationRequestInput{},
			types.RecipeStepIngredientUpdateRequestInput{},
		},
		"recipeStepInstruments.ts": {
			types.RecipeStepInstrument{},
			types.RecipeStepInstrumentCreationRequestInput{},
			types.RecipeStepInstrumentUpdateRequestInput{},
		},
		"recipeStepProducts.ts": {
			types.RecipeStepProduct{},
			types.RecipeStepProductCreationRequestInput{},
			types.RecipeStepProductUpdateRequestInput{},
		},
		"recipeSteps.ts": {
			types.RecipeStep{},
			types.RecipeStepCreationRequestInput{},
			types.RecipeStepUpdateRequestInput{},
		},
		"recipes.ts": {
			types.Recipe{},
			types.RecipeCreationRequestInput{},
			types.RecipeUpdateRequestInput{},
		},
		"users.ts": {
			types.UserStatusResponse{},
			types.User{},
			types.UserRegistrationInput{},
			types.UserCreationResponse{},
			types.UserLoginInput{},
			types.UsernameReminderRequestInput{},
			types.UserAccountStatusUpdateInput{},
		},
		"validIngredientMeasurementUnits.ts": {
			types.ValidIngredientMeasurementUnit{},
			types.ValidIngredientMeasurementUnitCreationRequestInput{},
			types.ValidIngredientMeasurementUnitUpdateRequestInput{},
		},
		"validIngredientPreparations.ts": {
			types.ValidIngredientPreparation{},
			types.ValidIngredientPreparationCreationRequestInput{},
			types.ValidIngredientPreparationUpdateRequestInput{},
		},
		"validIngredientStates.ts": {
			types.ValidIngredientState{},
			types.ValidIngredientStateCreationRequestInput{},
			types.ValidIngredientStateUpdateRequestInput{},
		},
		"validIngredientStateIngredients.ts": {
			types.ValidIngredientStateIngredient{},
			types.ValidIngredientStateIngredientCreationRequestInput{},
			types.ValidIngredientStateIngredientUpdateRequestInput{},
		},
		"validIngredients.ts": {
			types.ValidIngredient{},
			types.ValidIngredientCreationRequestInput{},
			types.ValidIngredientUpdateRequestInput{},
		},
		"validInstruments.ts": {
			types.ValidInstrument{},
			types.ValidInstrumentCreationRequestInput{},
			types.ValidInstrumentUpdateRequestInput{},
		},
		"validMeasurementUnitConversions.ts": {
			types.ValidMeasurementUnitConversion{},
			types.ValidMeasurementUnitConversionCreationRequestInput{},
			types.ValidMeasurementUnitConversionUpdateRequestInput{},
		},
		"validMeasurementUnits.ts": {
			types.ValidMeasurementUnit{},
			types.ValidMeasurementUnitCreationRequestInput{},
			types.ValidMeasurementUnitUpdateRequestInput{},
		},
		"validPreparationInstruments.ts": {
			types.ValidPreparationInstrument{},
			types.ValidPreparationInstrumentCreationRequestInput{},
			types.ValidPreparationInstrumentUpdateRequestInput{},
		},
		"validPreparations.ts": {
			types.ValidPreparation{},
			types.ValidPreparationCreationRequestInput{},
			types.ValidPreparationUpdateRequestInput{},
		},
		"webhooks.ts": {
			types.Webhook{},
			types.WebhookTriggerEvent{},
			types.WebhookCreationRequestInput{},
		},
	}

	generatedDisclaimer = "// Code generated by gen_typescript. DO NOT EDIT.\n\n"
	customTypeMap       = map[string]string{
		"MealPlan.status":                                        "ValidMealPlanStatus",
		"MealPlan.electionMethod":                                "ValidMealPlanElectionMethod",
		"ValidIngredientState.attributeType":                     "ValidIngredientStateAttributeType",
		"RecipeStepProduct.type":                                 "ValidRecipeStepProductType",
		"MealComponent.componentType":                            "MealComponentType",
		"MealPlanCreationRequestInput.status":                    "ValidMealPlanStatus",
		"MealPlanCreationRequestInput.electionMethod":            "ValidMealPlanElectionMethod",
		"ValidIngredientStateCreationRequestInput.attributeType": "ValidIngredientStateAttributeType",
		"RecipeStepProductCreationRequestInput.type":             "ValidRecipeStepProductType",
		"MealComponentCreationRequestInput.componentType":        "MealComponentType",
		"MealPlanUpdateRequestInput.status":                      "ValidMealPlanStatus",
		"MealPlanUpdateRequestInput.electionMethod":              "ValidMealPlanElectionMethod",
		"ValidIngredientStateUpdateRequestInput.attributeType":   "ValidIngredientStateAttributeType",
		"RecipeStepProductUpdateRequestInput.type":               "ValidRecipeStepProductType",
		"MealComponentUpdateRequestInput.componentType":          "MealComponentType",
	}

	defaultValues = map[string]string{
		"ValidMealPlanStatus":               "'awaiting_votes'",
		"ValidMealPlanElectionMethod":       "'schulze'",
		"ValidIngredientStateAttributeType": "'other'",
		"ValidRecipeStepProductType":        "'ingredient'",
		"MealComponentType":                 "'unspecified'",
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
