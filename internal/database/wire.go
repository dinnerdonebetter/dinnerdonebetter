package database

import (
	"github.com/google/wire"

	"github.com/prixfixeco/backend/pkg/types"
)

var (
	// Providers represents what we provide to dependency injectors.
	Providers = wire.NewSet(
		ProvideMealPlanTaskDataManager,
		ProvideHouseholdDataManager,
		ProvideHouseholdInvitationDataManager,
		ProvideHouseholdUserMembershipDataManager,
		ProvideValidInstrumentDataManager,
		ProvideValidIngredientDataManager,
		ProvideValidPreparationDataManager,
		ProvideValidMeasurementUnitDataManager,
		ProvideValidIngredientPreparationDataManager,
		ProvideValidPreparationInstrumentDataManager,
		ProvideValidIngredientMeasurementUnitDataManager,
		ProvideRecipeDataManager,
		ProvideRecipeStepDataManager,
		ProvideRecipeStepInstrumentDataManager,
		ProvideRecipeStepProductDataManager,
		ProvideRecipeStepIngredientDataManager,
		ProvideMealDataManager,
		ProvideMealPlanDataManager,
		ProvideMealPlanEventDataManager,
		ProvideMealPlanOptionDataManager,
		ProvideMealPlanOptionVoteDataManager,
		ProvideUserDataManager,
		ProvideAdminUserDataManager,
		ProvideAPIClientDataManager,
		ProvidePasswordResetTokenDataManager,
		ProvideWebhookDataManager,
		ProvideRecipePrepTaskDataManager,
		ProvideMealPlanGroceryListItemDataManager,
		ProvideValidMeasurementConversionDataManager,
		ProvideRecipeMediaDataManager,
		ProvideValidIngredientStateDataManager,
		ProvideRecipeStepCompletionConditionDataManager,
	)
)

// ProvideMealPlanTaskDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanTaskDataManager(db DataManager) types.MealPlanTaskDataManager {
	return db
}

// ProvideHouseholdDataManager is an arbitrary function for dependency injection's sake.
func ProvideHouseholdDataManager(db DataManager) types.HouseholdDataManager {
	return db
}

// ProvideHouseholdInvitationDataManager is an arbitrary function for dependency injection's sake.
func ProvideHouseholdInvitationDataManager(db DataManager) types.HouseholdInvitationDataManager {
	return db
}

// ProvideHouseholdUserMembershipDataManager is an arbitrary function for dependency injection's sake.
func ProvideHouseholdUserMembershipDataManager(db DataManager) types.HouseholdUserMembershipDataManager {
	return db
}

// ProvideValidInstrumentDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidInstrumentDataManager(db DataManager) types.ValidInstrumentDataManager {
	return db
}

// ProvideValidIngredientDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidIngredientDataManager(db DataManager) types.ValidIngredientDataManager {
	return db
}

// ProvideValidPreparationDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidPreparationDataManager(db DataManager) types.ValidPreparationDataManager {
	return db
}

// ProvideValidMeasurementUnitDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidMeasurementUnitDataManager(db DataManager) types.ValidMeasurementUnitDataManager {
	return db
}

// ProvideValidIngredientPreparationDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidIngredientPreparationDataManager(db DataManager) types.ValidIngredientPreparationDataManager {
	return db
}

// ProvideValidPreparationInstrumentDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidPreparationInstrumentDataManager(db DataManager) types.ValidPreparationInstrumentDataManager {
	return db
}

// ProvideValidIngredientMeasurementUnitDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidIngredientMeasurementUnitDataManager(db DataManager) types.ValidIngredientMeasurementUnitDataManager {
	return db
}

// ProvideRecipeDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeDataManager(db DataManager) types.RecipeDataManager {
	return db
}

// ProvideRecipeStepDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepDataManager(db DataManager) types.RecipeStepDataManager {
	return db
}

// ProvideRecipeStepInstrumentDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepInstrumentDataManager(db DataManager) types.RecipeStepInstrumentDataManager {
	return db
}

// ProvideRecipeStepProductDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepProductDataManager(db DataManager) types.RecipeStepProductDataManager {
	return db
}

// ProvideRecipeStepIngredientDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepIngredientDataManager(db DataManager) types.RecipeStepIngredientDataManager {
	return db
}

// ProvideMealDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealDataManager(db DataManager) types.MealDataManager {
	return db
}

// ProvideMealPlanDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanDataManager(db DataManager) types.MealPlanDataManager {
	return db
}

// ProvideMealPlanEventDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanEventDataManager(db DataManager) types.MealPlanEventDataManager {
	return db
}

// ProvideMealPlanOptionDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanOptionDataManager(db DataManager) types.MealPlanOptionDataManager {
	return db
}

// ProvideMealPlanOptionVoteDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanOptionVoteDataManager(db DataManager) types.MealPlanOptionVoteDataManager {
	return db
}

// ProvideUserDataManager is an arbitrary function for dependency injection's sake.
func ProvideUserDataManager(db DataManager) types.UserDataManager {
	return db
}

// ProvideAdminUserDataManager is an arbitrary function for dependency injection's sake.
func ProvideAdminUserDataManager(db DataManager) types.AdminUserDataManager {
	return db
}

// ProvideAPIClientDataManager is an arbitrary function for dependency injection's sake.
func ProvideAPIClientDataManager(db DataManager) types.APIClientDataManager {
	return db
}

// ProvidePasswordResetTokenDataManager is an arbitrary function for dependency injection's sake.
func ProvidePasswordResetTokenDataManager(db DataManager) types.PasswordResetTokenDataManager {
	return db
}

// ProvideWebhookDataManager is an arbitrary function for dependency injection's sake.
func ProvideWebhookDataManager(db DataManager) types.WebhookDataManager {
	return db
}

// ProvideRecipePrepTaskDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipePrepTaskDataManager(db DataManager) types.RecipePrepTaskDataManager {
	return db
}

// ProvideMealPlanGroceryListItemDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanGroceryListItemDataManager(db DataManager) types.MealPlanGroceryListItemDataManager {
	return db
}

// ProvideValidMeasurementConversionDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidMeasurementConversionDataManager(db DataManager) types.ValidMeasurementConversionDataManager {
	return db
}

// ProvideRecipeMediaDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeMediaDataManager(db DataManager) types.RecipeMediaDataManager {
	return db
}

// ProvideValidIngredientStateDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidIngredientStateDataManager(db DataManager) types.ValidIngredientStateDataManager {
	return db
}

// ProvideRecipeStepCompletionConditionDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeStepCompletionConditionDataManager(db DataManager) types.RecipeStepCompletionConditionDataManager {
	return db
}
