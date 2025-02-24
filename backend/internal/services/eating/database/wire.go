package database

import (
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/google/wire"
)

var (
	// DBProviders represents what we provide to dependency injectors.
	DBProviders = wire.NewSet(
		ProvideUserIngredientPreferenceDataManager,
		ProvideHouseholdInstrumentOwnershipDataManager,
		ProvideValidEnumerationDataManager,
		ProvideRecipeManagementDataManager,
		ProvideMealPlanningDataManager,
	)
)

// ProvideUserIngredientPreferenceDataManager is an arbitrary function for dependency injection's sake.
func ProvideUserIngredientPreferenceDataManager(db DataManager) types.IngredientPreferenceDataManager {
	return db
}

// ProvideHouseholdInstrumentOwnershipDataManager is an arbitrary function for dependency injection's sake.
func ProvideHouseholdInstrumentOwnershipDataManager(db DataManager) types.InstrumentOwnershipDataManager {
	return db
}

// ProvideValidEnumerationDataManager is an arbitrary function for dependency injection's sake.
func ProvideValidEnumerationDataManager(db DataManager) types.ValidEnumerationDataManager {
	return db
}

// ProvideRecipeManagementDataManager is an arbitrary function for dependency injection's sake.
func ProvideRecipeManagementDataManager(db DataManager) types.RecipeManagementDataManager {
	return db
}

// ProvideMealPlanningDataManager is an arbitrary function for dependency injection's sake.
func ProvideMealPlanningDataManager(db DataManager) types.MealPlanningDataManager {
	return db
}
