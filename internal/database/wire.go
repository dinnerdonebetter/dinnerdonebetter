package database

import (
	"github.com/google/wire"

	"github.com/prixfixeco/api_server/pkg/types"
)

var (
	// Providers represents what we provide to dependency injectors.
	Providers = wire.NewSet(
		ProvideHouseholdDataManager,
		ProvideHouseholdInvitationDataManager,
		ProvideHouseholdUserMembershipDataManager,
		ProvideValidInstrumentDataManager,
		ProvideValidIngredientDataManager,
		ProvideValidPreparationDataManager,
		ProvideValidIngredientPreparationDataManager,
		ProvideRecipeDataManager,
		ProvideRecipeStepDataManager,
		ProvideRecipeStepInstrumentDataManager,
		ProvideRecipeStepIngredientDataManager,
		ProvideMealDataManager,
		ProvideMealPlanDataManager,
		ProvideMealPlanOptionDataManager,
		ProvideMealPlanOptionVoteDataManager,
		ProvideUserDataManager,
		ProvideAdminUserDataManager,
		ProvideAPIClientDataManager,
		ProvideWebhookDataManager,
	)
)

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

// ProvideValidIngredientPreparationDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidIngredientPreparationDataManager(db DataManager) types.ValidIngredientPreparationDataManager {
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

// ProvideWebhookDataManager is an arbitrary function for dependency injection's sake.
func ProvideWebhookDataManager(db DataManager) types.WebhookDataManager {
	return db
}
