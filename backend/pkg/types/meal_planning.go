package types

type (
	MealPlanningDataManager interface {
		MealDataManager
		MealPlanDataManager
		MealPlanEventDataManager
		MealPlanOptionDataManager
		MealPlanOptionVoteDataManager
		MealPlanTaskDataManager
		MealPlanGroceryListItemDataManager
		UserIngredientPreferenceDataManager
		AccountInstrumentOwnershipDataManager
	}

	MealPlanningDataService interface {
		MealDataService
		MealPlanDataService
		MealPlanEventDataService
		MealPlanOptionDataService
		MealPlanOptionVoteDataService
		MealPlanTaskDataService
		MealPlanGroceryListItemDataService
		UserIngredientPreferenceDataService
		AccountInstrumentOwnershipDataService
	}
)
