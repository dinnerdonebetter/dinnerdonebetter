package mealplanning

type (
	MealPlanningDataManager interface {
		MealDataManager
		MealPlanDataManager
		MealPlanEventDataManager
		MealPlanOptionDataManager
		MealPlanOptionVoteDataManager
		MealPlanTaskDataManager
		MealPlanGroceryListItemDataManager
		MealPlanRecipeOptionSelectionDataManager
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
		MealPlanRecipeOptionSelectionDataService
		UserIngredientPreferenceDataService
		AccountInstrumentOwnershipDataService
	}
)
