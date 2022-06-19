package main

import (
	"io/ioutil"
	"os"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"

	"github.com/prixfixeco/api_server/pkg/types"
)

func writeFile(filename, content string) error {
	return ioutil.WriteFile(filename, []byte(content), os.ModePerm)
}

func buildConverterMap() map[string]*typescriptify.TypeScriptify {

	adminConverter := typescriptify.New().
		Add(types.UserAccountStatusUpdateInput{})

	apiClientConverter := typescriptify.New().
		Add(types.APIClient{}).
		Add(types.APIClientList{}).
		Add(types.APIClientCreationRequestInput{}).
		Add(types.APIClientCreationResponse{})

	// NOTE: Do not include SessionContext here.
	authConverter := typescriptify.New().
		Add(types.UserHouseholdMembershipInfo{}).
		Add(types.RequesterInfo{}).
		Add(types.UserStatusResponse{}).
		Add(types.ChangeActiveHouseholdInput{}).
		Add(types.PASETOCreationInput{}).
		Add(types.PASETOResponse{})

	householdConverter := typescriptify.New().
		Add(types.Household{}).
		Add(types.HouseholdList{}).
		Add(types.HouseholdCreationRequestInput{}).
		Add(types.HouseholdUpdateRequestInput{})

	householdInvitationConverter := typescriptify.New().
		Add(types.HouseholdInvitation{}).
		Add(types.HouseholdInvitationList{}).
		Add(types.HouseholdInvitationCreationRequestInput{}).
		Add(types.HouseholdInvitationUpdateRequestInput{})

	householdUserMembershipConverter := typescriptify.New().
		Add(types.HouseholdUserMembership{}).
		Add(types.HouseholdUserMembershipList{}).
		Add(types.HouseholdUserMembershipUpdateRequestInput{}).
		Add(types.HouseholdUserMembershipCreationRequestInput{}).
		Add(types.HouseholdOwnershipTransferInput{}).
		Add(types.ModifyUserPermissionsInput{})

	mealConverter := typescriptify.New().
		Add(types.Meal{}).
		Add(types.MealList{}).
		Add(types.MealCreationRequestInput{}).
		Add(types.MealUpdateRequestInput{})

	mealPlanConverter := typescriptify.New().
		Add(types.MealPlan{}).
		Add(types.MealPlanList{}).
		Add(types.MealPlanCreationRequestInput{}).
		Add(types.MealPlanUpdateRequestInput{})

	mealPlanOptionConverter := typescriptify.New().
		Add(types.MealPlanOption{}).
		Add(types.MealPlanOptionList{}).
		Add(types.MealPlanOptionCreationRequestInput{}).
		Add(types.MealPlanOptionUpdateRequestInput{})

	mealPlanOptionVoteConverter := typescriptify.New().
		Add(types.MealPlanOptionVote{}).
		Add(types.MealPlanOptionVoteList{}).
		Add(types.MealPlanOptionVoteCreationRequestInput{}).
		Add(types.MealPlanOptionVoteUpdateRequestInput{})

	queryFilterConverter := typescriptify.New().
		Add(types.QueryFilter{})

	recipeConverter := typescriptify.New().
		Add(types.Recipe{}).
		Add(types.RecipeList{}).
		Add(types.RecipeCreationRequestInput{}).
		Add(types.RecipeUpdateRequestInput{})

	recipeStepConverter := typescriptify.New().
		Add(types.RecipeStep{}).
		Add(types.RecipeStepList{}).
		Add(types.RecipeStepCreationRequestInput{}).
		Add(types.RecipeStepUpdateRequestInput{})

	recipeStepIngredientConverter := typescriptify.New().
		Add(types.RecipeStepIngredient{}).
		Add(types.RecipeStepIngredientList{}).
		Add(types.RecipeStepIngredientCreationRequestInput{}).
		Add(types.RecipeStepIngredientUpdateRequestInput{})

	recipeStepInstrumentConverter := typescriptify.New().
		Add(types.RecipeStepInstrument{}).
		Add(types.RecipeStepInstrumentList{}).
		Add(types.RecipeStepInstrumentCreationRequestInput{}).
		Add(types.RecipeStepInstrumentUpdateRequestInput{})

	recipeStepProductConverter := typescriptify.New().
		Add(types.RecipeStepProduct{}).
		Add(types.RecipeStepProductList{}).
		Add(types.RecipeStepProductCreationRequestInput{}).
		Add(types.RecipeStepProductUpdateRequestInput{})

	userConverter := typescriptify.New().
		Add(types.UserList{}).
		Add(types.UserRegistrationInput{}).
		Add(types.UserDatabaseCreationInput{}).
		Add(types.UserCreationResponse{}).
		Add(types.UserLoginInput{}).
		Add(types.PasswordUpdateInput{}).
		Add(types.TOTPSecretRefreshInput{}).
		Add(types.TOTPSecretVerificationInput{}).
		Add(types.TOTPSecretRefreshResponse{})

	validIngredientConverter := typescriptify.New().
		Add(types.ValidIngredient{}).
		Add(types.ValidIngredientList{}).
		Add(types.ValidIngredientCreationRequestInput{}).
		Add(types.ValidIngredientUpdateRequestInput{})

	validIngredientPreparationConverter := typescriptify.New().
		Add(types.ValidIngredientPreparation{}).
		Add(types.ValidIngredientPreparationList{}).
		Add(types.ValidIngredientPreparationCreationRequestInput{}).
		Add(types.ValidIngredientPreparationUpdateRequestInput{})

	validInstrumentConverter := typescriptify.New().
		Add(types.ValidInstrument{}).
		Add(types.ValidInstrumentList{}).
		Add(types.ValidInstrumentCreationRequestInput{}).
		Add(types.ValidInstrumentUpdateRequestInput{})

	validPreparationConverter := typescriptify.New().
		Add(types.ValidPreparation{}).
		Add(types.ValidPreparationList{}).
		Add(types.ValidPreparationCreationRequestInput{}).
		Add(types.ValidPreparationUpdateRequestInput{})

	webhookConverter := typescriptify.New().
		Add(types.Webhook{}).
		Add(types.WebhookCreationRequestInput{}).
		Add(types.WebhookList{})

	return map[string]*typescriptify.TypeScriptify{
		"artifacts/typescript/admin.ts":                       adminConverter,
		"artifacts/typescript/apiClients.ts":                  apiClientConverter,
		"artifacts/typescript/auth.ts":                        authConverter,
		"artifacts/typescript/households.ts":                  householdConverter,
		"artifacts/typescript/householdInvitations.ts":        householdInvitationConverter,
		"artifacts/typescript/householdUserMemberships.ts":    householdUserMembershipConverter,
		"artifacts/typescript/meals.ts":                       mealConverter,
		"artifacts/typescript/mealPlans.ts":                   mealPlanConverter,
		"artifacts/typescript/mealPlanOptions.ts":             mealPlanOptionConverter,
		"artifacts/typescript/mealPlanOptionVotes.ts":         mealPlanOptionVoteConverter,
		"artifacts/typescript/pagination.ts":                  queryFilterConverter,
		"artifacts/typescript/recipes.ts":                     recipeConverter,
		"artifacts/typescript/recipeSteps.ts":                 recipeStepConverter,
		"artifacts/typescript/recipeStepIngredients.ts":       recipeStepIngredientConverter,
		"artifacts/typescript/recipeStepInstruments.ts":       recipeStepInstrumentConverter,
		"artifacts/typescript/recipeStepProducts.ts":          recipeStepProductConverter,
		"artifacts/typescript/users.ts":                       userConverter,
		"artifacts/typescript/validIngredients.ts":            validIngredientConverter,
		"artifacts/typescript/validIngredientPreparations.ts": validIngredientPreparationConverter,
		"artifacts/typescript/validInstruments.ts":            validInstrumentConverter,
		"artifacts/typescript/validPreparations.ts":           validPreparationConverter,
		"artifacts/typescript/webhooks.ts":                    webhookConverter,
	}
}

func main() {
	converters := buildConverterMap()

	if err := os.MkdirAll("artifacts/typescript", os.ModeDir); err != nil {
		panic(err)
	}

	for filename, converter := range converters {
		output, convertErr := converter.Convert(nil)
		if convertErr != nil {
			panic(convertErr.Error())
		}

		if writeErr := writeFile(filename, output); writeErr != nil {
			panic(writeErr.Error())
		}
	}
}
