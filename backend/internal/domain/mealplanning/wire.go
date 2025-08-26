package mealplanning

import (
	"github.com/google/wire"
)

var (
	Providers = wire.NewSet(
		ProvideMealDataManagerFromRepository,
		ProvideMealPlanDataManagerFromRepository,
		ProvideMealPlanEventDataManagerFromRepository,
		ProvideMealPlanGroceryListItemDataManagerFromRepository,
		ProvideMealPlanOptionDataManagerFromRepository,
		ProvideMealPlanOptionVoteDataManagerFromRepository,
		ProvideMealPlanTaskDataManagerFromRepository,
		ProvideMealPlanningDataManagerFromRepository,
		ProvideRecipeDataManagerFromRepository,
		ProvideRecipeManagementDataManagerFromRepository,
		ProvideRecipeMediaDataManagerFromRepository,
		ProvideRecipePrepTaskDataManagerFromRepository,
		ProvideRecipeRatingDataManagerFromRepository,
		ProvideRecipeStepDataManagerFromRepository,
		ProvideRecipeStepCompletionConditionDataManagerFromRepository,
		ProvideValidIngredientPreparationDataManagerFromRepository,
		ProvideRecipeStepIngredientDataManagerFromRepository,
		ProvideRecipeStepVesselDataManagerFromRepository,
		ProvideValidIngredientDataManagerFromRepository,
		ProvideRecipeStepProductDataManagerFromRepository,
		ProvideRecipeStepInstrumentDataManagerFromRepository,
		ProvideValidEnumerationDataManagerFromRepository,
		ProvideUserIngredientPreferenceDataManagerFromRepository,
		ProvideValidIngredientMeasurementUnitDataManagerFromRepository,
		ProvideValidIngredientGroupDataManagerFromRepository,
		ProvideValidIngredientStateDataManagerFromRepository,
		ProvideValidIngredientStateIngredientDataManagerFromRepository,
		ProvideValidInstrumentDataManagerFromRepository,
		ProvideValidMeasurementUnitDataManagerFromRepository,
		ProvideValidMeasurementUnitConversionDataManagerFromRepository,
		ProvideValidPreparationDataManagerFromRepository,
		ProvideValidPreparationInstrumentDataManagerFromRepository,
		ProvideValidPreparationVesselDataManagerFromRepository,
		ProvideValidVesselDataManagerFromRepository,
	)
)

func ProvideMealDataManagerFromRepository(r Repository) MealDataManager {
	return r
}

func ProvideMealPlanDataManagerFromRepository(r Repository) MealPlanDataManager {
	return r
}

func ProvideMealPlanEventDataManagerFromRepository(r Repository) MealPlanEventDataManager {
	return r
}

func ProvideMealPlanGroceryListItemDataManagerFromRepository(r Repository) MealPlanGroceryListItemDataManager {
	return r
}

func ProvideMealPlanOptionDataManagerFromRepository(r Repository) MealPlanOptionDataManager {
	return r
}

func ProvideMealPlanOptionVoteDataManagerFromRepository(r Repository) MealPlanOptionVoteDataManager {
	return r
}

func ProvideMealPlanTaskDataManagerFromRepository(r Repository) MealPlanTaskDataManager {
	return r
}

func ProvideMealPlanningDataManagerFromRepository(r Repository) MealPlanningDataManager {
	return r
}

func ProvideRecipeDataManagerFromRepository(r Repository) RecipeDataManager {
	return r
}

func ProvideRecipeManagementDataManagerFromRepository(r Repository) RecipeManagementDataManager {
	return r
}

func ProvideRecipeMediaDataManagerFromRepository(r Repository) RecipeMediaDataManager {
	return r
}

func ProvideRecipePrepTaskDataManagerFromRepository(r Repository) RecipePrepTaskDataManager {
	return r
}

func ProvideRecipeRatingDataManagerFromRepository(r Repository) RecipeRatingDataManager {
	return r
}

func ProvideRecipeStepDataManagerFromRepository(r Repository) RecipeStepDataManager {
	return r
}

func ProvideRecipeStepCompletionConditionDataManagerFromRepository(r Repository) RecipeStepCompletionConditionDataManager {
	return r
}

func ProvideValidIngredientPreparationDataManagerFromRepository(r Repository) ValidIngredientPreparationDataManager {
	return r
}

func ProvideRecipeStepIngredientDataManagerFromRepository(r Repository) RecipeStepIngredientDataManager {
	return r
}

func ProvideRecipeStepVesselDataManagerFromRepository(r Repository) RecipeStepVesselDataManager {
	return r
}

func ProvideValidIngredientDataManagerFromRepository(r Repository) ValidIngredientDataManager {
	return r
}

func ProvideRecipeStepProductDataManagerFromRepository(r Repository) RecipeStepProductDataManager {
	return r
}

func ProvideRecipeStepInstrumentDataManagerFromRepository(r Repository) RecipeStepInstrumentDataManager {
	return r
}

func ProvideValidEnumerationDataManagerFromRepository(r Repository) ValidEnumerationDataManager {
	return r
}

func ProvideUserIngredientPreferenceDataManagerFromRepository(r Repository) UserIngredientPreferenceDataManager {
	return r
}

func ProvideValidIngredientMeasurementUnitDataManagerFromRepository(r Repository) ValidIngredientMeasurementUnitDataManager {
	return r
}

func ProvideValidIngredientGroupDataManagerFromRepository(r Repository) ValidIngredientGroupDataManager {
	return r
}

func ProvideValidIngredientStateDataManagerFromRepository(r Repository) ValidIngredientStateDataManager {
	return r
}

func ProvideValidIngredientStateIngredientDataManagerFromRepository(r Repository) ValidIngredientStateIngredientDataManager {
	return r
}

func ProvideValidInstrumentDataManagerFromRepository(r Repository) ValidInstrumentDataManager {
	return r
}

func ProvideValidMeasurementUnitDataManagerFromRepository(r Repository) ValidMeasurementUnitDataManager {
	return r
}

func ProvideValidMeasurementUnitConversionDataManagerFromRepository(r Repository) ValidMeasurementUnitConversionDataManager {
	return r
}

func ProvideValidPreparationDataManagerFromRepository(r Repository) ValidPreparationDataManager {
	return r
}

func ProvideValidPreparationInstrumentDataManagerFromRepository(r Repository) ValidPreparationInstrumentDataManager {
	return r
}

func ProvideValidPreparationVesselDataManagerFromRepository(r Repository) ValidPreparationVesselDataManager {
	return r
}

func ProvideValidVesselDataManagerFromRepository(r Repository) ValidVesselDataManager {
	return r
}
