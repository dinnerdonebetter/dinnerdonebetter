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
		IngredientPreferenceDataManager
		InstrumentOwnershipDataManager
	}

	MealPlanningDataService interface {
		MealDataService
		MealPlanDataService
		MealPlanEventDataService
		MealPlanOptionDataService
		MealPlanOptionVoteDataService
		MealPlanTaskDataService
		MealPlanGroceryListItemDataService
		IngredientPreferenceDataService
		InstrumentOwnershipDataService
	}
)
